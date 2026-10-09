package filter

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestParseExpression(t *testing.T) {
	expr, err := ParseExpression("level=ERROR AND (status>=500 OR latency>200ms)")
	if err != nil {
		t.Fatalf("ParseExpression failed: %v", err)
	}

	r1 := &model.Record{
		Level:      model.LevelError,
		StatusCode: 503,
		Latency:    50 * time.Millisecond,
	}
	if !expr.Match(r1) {
		t.Errorf("r1 (status 503) should match")
	}

	r2 := &model.Record{
		Level:      model.LevelError,
		StatusCode: 200,
		Latency:    300 * time.Millisecond,
	}
	if !expr.Match(r2) {
		t.Errorf("r2 (latency 300ms) should match")
	}

	r3 := &model.Record{
		Level:      model.LevelInfo,
		StatusCode: 500,
	}
	if expr.Match(r3) {
		t.Errorf("r3 (level INFO) should NOT match")
	}

	exprNot, err := ParseExpression("NOT level=DEBUG")
	if err != nil {
		t.Fatalf("ParseExpression NOT failed: %v", err)
	}
	if !exprNot.Match(r1) {
		t.Errorf("r1 (ERROR) should match NOT level=DEBUG")
	}
}
