package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
	_ "modernc.org/sqlite"
)

// SQLiteIndexer indexes logs into an embedded SQLite database using batching.
type SQLiteIndexer struct {
	db        *sql.DB
	tx        *sql.Tx
	stmt      *sql.Stmt
	batchSize int
	pending   int
}

// OpenSQLite opens or creates an indexed SQLite database.
func OpenSQLite(dbPath string, batchSize int) (*SQLiteIndexer, error) {
	if batchSize <= 0 {
		batchSize = 10000
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite failed: %w", err)
	}

	// Performance pragmas
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA cache_size = -64000;", // 64MB cache
		"PRAGMA temp_store = MEMORY;",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return nil, err
		}
	}

	schema := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		line_num INTEGER,
		timestamp TEXT,
		level TEXT,
		latency_ms REAL,
		status_code INTEGER,
		method TEXT,
		path TEXT,
		ip TEXT,
		trace_id TEXT,
		message TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp);
	CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);
	CREATE INDEX IF NOT EXISTS idx_logs_status ON logs(status_code);
	CREATE INDEX IF NOT EXISTS idx_logs_path ON logs(path);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("create schema failed: %w", err)
	}

	idx := &SQLiteIndexer{
		db:        db,
		batchSize: batchSize,
	}

	if err := idx.beginTx(); err != nil {
		return nil, err
	}

	return idx, nil
}

func (idx *SQLiteIndexer) beginTx() error {
	tx, err := idx.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
		INSERT INTO logs (line_num, timestamp, level, latency_ms, status_code, method, path, ip, trace_id, message)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	idx.tx = tx
	idx.stmt = stmt
	idx.pending = 0
	return nil
}

// Index adds a record into the index buffer, committing transactions at batch boundaries.
func (idx *SQLiteIndexer) Index(r *model.Record) error {
	tsStr := ""
	if !r.Timestamp.IsZero() {
		tsStr = r.Timestamp.Format(time.RFC3339Nano)
	}

	latMs := float64(r.Latency.Milliseconds())

	_, err := idx.stmt.Exec(
		r.LineNum,
		tsStr,
		r.Level.String(),
		latMs,
		r.StatusCode,
		r.Method,
		r.Path,
		r.IP,
		r.TraceID,
		r.Message,
	)
	if err != nil {
		return err
	}

	idx.pending++
	if idx.pending >= idx.batchSize {
		if err := idx.Flush(); err != nil {
			return err
		}
		return idx.beginTx()
	}

	return nil
}

// Flush commits current transaction batch.
func (idx *SQLiteIndexer) Flush() error {
	if idx.stmt != nil {
		_ = idx.stmt.Close()
		idx.stmt = nil
	}
	if idx.tx != nil {
		if err := idx.tx.Commit(); err != nil {
			return err
		}
		idx.tx = nil
	}
	idx.pending = 0
	return nil
}

// Close flushes remaining records and closes SQLite handle.
func (idx *SQLiteIndexer) Close() error {
	_ = idx.Flush()
	return idx.db.Close()
}
