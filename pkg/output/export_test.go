package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestJSONExporter(t *testing.T) {
	var buf bytes.Buffer
	exp := NewJSONExporter(&buf)

	rec := &model.Record{
		Message: "test message",
		Level:   model.LevelInfo,
	}
	if err := exp.Export(rec); err != nil {
		t.Fatalf("export error: %v", err)
	}

	if !strings.Contains(buf.String(), `"message":"test message"`) {
		t.Errorf("expected json output with message, got %s", buf.String())
	}
}

func TestCSVExporter(t *testing.T) {
	var buf bytes.Buffer
	exp := NewCSVExporter(&buf)

	rec := &model.Record{
		Timestamp:  time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC),
		Level:      model.LevelError,
		StatusCode: 500,
		Message:    "internal server error",
	}
	if err := exp.Export(rec); err != nil {
		t.Fatalf("csv export error: %v", err)
	}
	exp.Flush()

	out := buf.String()
	if !strings.Contains(out, "timestamp,level,latency_ms,status") {
		t.Errorf("expected csv header, got %s", out)
	}
	if !strings.Contains(out, "ERROR") || !strings.Contains(out, "500") {
		t.Errorf("expected csv data row, got %s", out)
	}
}
