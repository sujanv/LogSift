package tracectx

import (
	"testing"
)

func TestParseTraceParent(t *testing.T) {
	raw := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	tc, err := ParseTraceParent(raw)
	if err != nil {
		t.Fatalf("unexpected error parsing valid traceparent: %v", err)
	}
	if tc.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Errorf("traceID mismatch: %s", tc.TraceID)
	}
	if tc.SpanID != "00f067aa0ba902b7" {
		t.Errorf("spanID mismatch: %s", tc.SpanID)
	}
	if !tc.Sampled {
		t.Errorf("expected sampled=true")
	}
	if tc.String() != raw {
		t.Errorf("String() mismatch: %s != %s", tc.String(), raw)
	}

	state := ParseTraceState("rojo=123,congo=456")
	if state["rojo"] != "123" || state["congo"] != "456" {
		t.Errorf("tracestate parse failure: %+v", state)
	}
}
