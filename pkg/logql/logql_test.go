package logql

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestLogQLParseAndMatch(t *testing.T) {
	raw := `{level="error", path=~"^/api.*"} |= "timeout" != "ignored"`
	q, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(q.Matchers) != 2 {
		t.Fatalf("expected 2 label matchers, got %d", len(q.Matchers))
	}
	if len(q.Stages) != 2 {
		t.Fatalf("expected 2 pipeline stages, got %d", len(q.Stages))
	}

	rec1 := &model.Record{
		Level:   model.LevelError,
		Path:    "/api/v1/checkout",
		Message: "downstream service timeout occurred",
	}
	if !q.Matches(rec1) {
		t.Errorf("expected rec1 to match LogQL query")
	}

	rec2 := &model.Record{
		Level:   model.LevelInfo,
		Path:    "/api/v1/checkout",
		Message: "downstream service timeout occurred",
	}
	if q.Matches(rec2) {
		t.Errorf("rec2 level info should not match")
	}

	rec3 := &model.Record{
		Level:   model.LevelError,
		Path:    "/api/v1/checkout",
		Message: "downstream service timeout ignored",
	}
	if q.Matches(rec3) {
		t.Errorf("rec3 should be excluded by pipeline stage")
	}
}

func TestLogQLMetric(t *testing.T) {
	raw := `rate({level="error"}[5m]) |= "database"`
	q, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if q.MetricFunc != "rate" {
		t.Errorf("expected MetricFunc rate, got %s", q.MetricFunc)
	}
	if q.Window != 5*time.Minute {
		t.Errorf("expected 5m window, got %v", q.Window)
	}
}
