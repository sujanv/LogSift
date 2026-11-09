package test

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/breaker"
	"github.com/sujanv/LogSift/pkg/logql"
	"github.com/sujanv/LogSift/pkg/tracectx"
)

func TestPhase6Regressions(t *testing.T) {
	_, err := tracectx.ParseTraceParent("malformed-header")
	if err == nil {
		t.Errorf("expected error on malformed traceparent")
	}

	_, err = logql.Parse("")
	if err == nil {
		t.Errorf("expected error on empty logql query")
	}

	cb := breaker.NewCircuitBreaker("brk", 1, 1)
	cb.Reset()
	if cb.State() != breaker.StateClosed {
		t.Errorf("expected breaker to be closed after reset")
	}
}
