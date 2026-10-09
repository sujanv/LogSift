package parser

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestLogfmtParser(t *testing.T) {
	p := NewLogfmtParser()
	rec := &model.Record{Fields: make(map[string]interface{})}

	line := []byte(`ts=2026-10-09T09:00:00Z level=error msg="failed to connect to database" duration=350ms status=503 path=/api/v2/items client_ip=10.0.0.1 trace_id=xyz-999`)
	if err := p.Parse(line, rec); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if rec.Level != model.LevelError {
		t.Errorf("expected LevelError, got %v", rec.Level)
	}
	if rec.Message != "failed to connect to database" {
		t.Errorf("expected 'failed to connect to database', got %q", rec.Message)
	}
	if rec.Latency != 350*time.Millisecond {
		t.Errorf("expected 350ms, got %v", rec.Latency)
	}
	if rec.StatusCode != 503 {
		t.Errorf("expected 503, got %d", rec.StatusCode)
	}
	if rec.Path != "/api/v2/items" {
		t.Errorf("expected /api/v2/items, got %q", rec.Path)
	}
	if rec.TraceID != "xyz-999" {
		t.Errorf("expected xyz-999, got %q", rec.TraceID)
	}
}
