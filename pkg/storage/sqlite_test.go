package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestSQLiteIndexer(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sqlite_test_*")
	if err != nil {
		t.Fatalf("mktemp error: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	indexer, err := OpenSQLite(dbPath, 100)
	if err != nil {
		t.Fatalf("OpenSQLite error: %v", err)
	}

	rec := &model.Record{
		LineNum:    1,
		Timestamp:  time.Now(),
		Level:      model.LevelError,
		Latency:    250 * time.Millisecond,
		StatusCode: 500,
		Path:       "/api/pay",
		Message:    "payment failure",
	}

	if err := indexer.Index(rec); err != nil {
		t.Fatalf("Index error: %v", err)
	}

	if err := indexer.Close(); err != nil {
		t.Fatalf("Close error: %v", err)
	}
}
