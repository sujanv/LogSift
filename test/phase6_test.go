package test

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/breaker"
	"github.com/sujanv/LogSift/pkg/embed"
	"github.com/sujanv/LogSift/pkg/logql"
	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/schemaval"
	"github.com/sujanv/LogSift/pkg/throttle"
	"github.com/sujanv/LogSift/pkg/tracectx"
)

func TestPhase6Integration(t *testing.T) {
	tc, err := tracectx.ParseTraceParent("00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	if err != nil || tc.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("failed trace context parsing: %v", err)
	}

	cb := breaker.NewCircuitBreaker("test-brk", 2, 10*time.Millisecond)
	_ = cb.Execute(func() error { return nil })
	if cb.State() != breaker.StateClosed {
		t.Errorf("circuit breaker should remain closed on success")
	}

	s := throttle.NewErrorPreservingSampler(0.01)
	errRec := &model.Record{Level: model.LevelError, StatusCode: 500}
	if !s.ShouldSample(errRec) {
		t.Errorf("error record must be preserved")
	}

	q, err := logql.Parse(`{level="error"} |= "fail"`)
	if err != nil || !q.Matches(&model.Record{Level: model.LevelError, Message: "job failed"}) {
		t.Errorf("LogQL match failed: %v", err)
	}

	val := schemaval.NewSchemaValidator([]string{"message"}, nil)
	errs := val.Validate(&model.Record{Message: "hello"})
	if len(errs) != 0 {
		t.Errorf("unexpected schema errors: %v", errs)
	}

	v := embed.NewVectorizer()
	v.Train([]string{"disk full on dev1", "memory limit exceeded"})
	v1 := v.Transform("disk full on dev1")
	if len(v1) == 0 {
		t.Errorf("vector transform failed")
	}
}
