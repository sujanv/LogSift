package model

import (
	"testing"
	"time"
)

func TestLogLevelParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"DEBUG", LevelDebug},
		{"debug", LevelDebug},
		{"INFO", LevelInfo},
		{"info", LevelInfo},
		{"WARN", LevelWarn},
		{"warning", LevelWarn},
		{"ERROR", LevelError},
		{"err", LevelError},
		{"FATAL", LevelFatal},
		{"critical", LevelFatal},
		{"unknown", LevelUnknown},
		{"", LevelUnknown},
	}

	for _, tt := range tests {
		got := ParseLevel(tt.input)
		if got != tt.expected {
			t.Errorf("ParseLevel(%q) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

func TestLogLevelBytes(t *testing.T) {
	if ParseLevelBytes([]byte("info")) != LevelInfo {
		t.Errorf("ParseLevelBytes('info') failed")
	}
	if ParseLevelBytes([]byte("ERROR")) != LevelError {
		t.Errorf("ParseLevelBytes('ERROR') failed")
	}
}

func TestRecordResetAndPool(t *testing.T) {
	pool := NewRecordPool()
	r := pool.Get()
	r.Timestamp = time.Now()
	r.Level = LevelError
	r.Message = "database connection timed out"
	r.Latency = 250 * time.Millisecond
	r.StatusCode = 500
	r.Fields["query"] = "SELECT 1"

	pool.Put(r)

	r2 := pool.Get()
	if r2.Message != "" {
		t.Errorf("expected empty Message after pool Get, got %q", r2.Message)
	}
	if r2.Level != LevelUnknown {
		t.Errorf("expected LevelUnknown, got %v", r2.Level)
	}
	if len(r2.Fields) != 0 {
		t.Errorf("expected empty fields map, got len %d", len(r2.Fields))
	}
}
