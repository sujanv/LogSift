package slowlog

import (
	"testing"
	"time"
)

func TestParsePostgresLine(t *testing.T) {
	line := "2026-10-14 12:00:00 UTC [12345]: LOG:  duration: 450.25 ms  statement: SELECT * FROM orders WHERE user_id = 99 AND status = 'pending';"
	sq := ParsePostgresLine(line)
	if sq == nil {
		t.Fatalf("failed to parse postgres slow log line")
	}
	if sq.Duration < 450*time.Millisecond {
		t.Errorf("duration mismatch: %v", sq.Duration)
	}
	if sq.NormalizedSQL != "SELECT * FROM ORDERS WHERE USER_ID = ? AND STATUS = ?;" {
		t.Errorf("normalized SQL mismatch: %s", sq.NormalizedSQL)
	}
}

func TestParseMySQLHeader(t *testing.T) {
	header := "# Query_time: 2.500000  Lock_time: 0.001000 Rows_sent: 10  Rows_examined: 100000"
	query := "SELECT * FROM users WHERE age > 21;"
	sq := ParseMySQLHeader(header, query)
	if sq == nil {
		t.Fatalf("failed to parse mysql slow log line")
	}
	if sq.Duration != 2500*time.Millisecond {
		t.Errorf("duration mismatch: %v", sq.Duration)
	}
	if sq.RowsExamined != 100000 {
		t.Errorf("rows examined mismatch: %d", sq.RowsExamined)
	}
}
