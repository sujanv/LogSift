package filter

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestFieldMatcher(t *testing.T) {
	rec := &model.Record{
		Level:      model.LevelError,
		StatusCode: 502,
		Latency:    350 * time.Millisecond,
		Path:       "/api/v1/orders",
		Fields:     map[string]interface{}{"env": "production"},
	}

	m1 := NewFieldMatcher("level", OpEq, "ERROR")
	if !m1.Match(rec) {
		t.Errorf("expected level=ERROR to match")
	}

	m2 := NewFieldMatcher("status", OpGte, "500")
	if !m2.Match(rec) {
		t.Errorf("expected status>=500 to match")
	}

	m3 := NewFieldMatcher("latency", OpGt, "200ms")
	if !m3.Match(rec) {
		t.Errorf("expected latency>200ms to match")
	}

	m4 := NewFieldMatcher("path", OpContains, "orders")
	if !m4.Match(rec) {
		t.Errorf("expected path contains orders to match")
	}

	m5 := NewFieldMatcher("env", OpEq, "production")
	if !m5.Match(rec) {
		t.Errorf("expected custom field env=production to match")
	}

	m6 := NewFieldMatcher("status", OpLt, "400")
	if m6.Match(rec) {
		t.Errorf("expected status<400 NOT to match")
	}
}
