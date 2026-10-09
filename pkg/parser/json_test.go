package parser

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestJSONParser_Standard(t *testing.T) {
	p := NewJSONParser()
	rec := &model.Record{Fields: make(map[string]interface{})}

	raw := []byte(`{"timestamp":"2026-10-09T08:15:30Z","level":"info","msg":"user checkout","latency_ms":145.5,"status":200,"path":"/api/checkout","method":"POST","trace_id":"trace-12345"}`)
	if err := p.Parse(raw, rec); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if rec.Level != model.LevelInfo {
		t.Errorf("expected LevelInfo, got %v", rec.Level)
	}
	if rec.Message != "user checkout" {
		t.Errorf("expected 'user checkout', got %q", rec.Message)
	}
	if rec.Latency != time.Duration(145.5*float64(time.Millisecond)) {
		t.Errorf("expected 145.5ms, got %v", rec.Latency)
	}
	if rec.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", rec.StatusCode)
	}
	if rec.Path != "/api/checkout" {
		t.Errorf("expected /api/checkout, got %q", rec.Path)
	}
	if rec.Method != "POST" {
		t.Errorf("expected POST, got %q", rec.Method)
	}
	if rec.TraceID != "trace-12345" {
		t.Errorf("expected trace-12345, got %q", rec.TraceID)
	}
}

func TestJSONParser_NumericLevelAndDurationString(t *testing.T) {
	p := NewJSONParser()
	rec := &model.Record{Fields: make(map[string]interface{})}

	raw := []byte(`{"time":1760000000000,"lvl":50,"message":"internal error","duration":"2.5s","status_code":"500"}`)
	if err := p.Parse(raw, rec); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if rec.Level != model.LevelError {
		t.Errorf("expected LevelError for lvl=50, got %v", rec.Level)
	}
	if rec.Latency != 2500*time.Millisecond {
		t.Errorf("expected 2.5s duration, got %v", rec.Latency)
	}
	if rec.StatusCode != 500 {
		t.Errorf("expected status 500, got %d", rec.StatusCode)
	}
}

func BenchmarkJSONParser(b *testing.B) {
	p := NewJSONParser()
	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	line := []byte(`{"ts":"2026-10-09T08:15:30Z","level":"warn","msg":"cache miss for key","latency_ms":12.3,"status":200,"path":"/v1/items","client_ip":"192.168.1.10"}`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec.Reset()
		_ = p.Parse(line, rec)
	}
}
