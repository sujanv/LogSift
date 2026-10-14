package schema

import (
	"strings"
	"testing"
)

func TestSchemaInferenceAndDrift(t *testing.T) {
	s := NewInferredSchema()

	l1 := []byte(`{"level":"info","status":200,"path":"/index"}`)
	l2 := []byte(`{"level":"warn","status":404,"path":"/missing"}`)

	_ = s.IngestJSON(l1)
	_ = s.IngestJSON(l2)

	summary := s.Summary()
	if summary.FieldCount != 3 {
		t.Errorf("expected 3 fields, got %d", summary.FieldCount)
	}

	l3 := []byte(`{"level":"error","status":500,"path":"/err","unknown_extra":true}`)
	drifts := s.DetectDrift(l3)
	if len(drifts) == 0 || !strings.Contains(drifts[0], "unknown_extra") {
		t.Errorf("expected drift on new field")
	}

	l4 := []byte(`{"level":"info","status":"OK","path":"/index"}`)
	drifts2 := s.DetectDrift(l4)
	if len(drifts2) == 0 || !strings.Contains(drifts2[0], "Type drift") {
		t.Errorf("expected type drift on status: %v", drifts2)
	}
}
