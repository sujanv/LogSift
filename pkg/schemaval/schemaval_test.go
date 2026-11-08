package schemaval

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestSchemaValidator(t *testing.T) {
	val := NewSchemaValidator([]string{"message", "status"}, map[string]FieldType{
		"status": TypeInt,
	})

	r1 := &model.Record{
		Message: "all good",
		Fields: map[string]interface{}{
			"status": 200,
		},
	}
	errs := val.Validate(r1)
	if len(errs) != 0 {
		t.Fatalf("expected valid record, got errs: %v", errs)
	}

	r2 := &model.Record{
		Message: "missing status",
	}
	errs2 := val.Validate(r2)
	if len(errs2) == 0 {
		t.Errorf("expected validation failure for missing required status")
	}

	r3 := &model.Record{
		Fields: map[string]interface{}{
			"new_tenant_id": "tenant-99",
		},
	}
	drifts := val.DetectDrift(r3)
	if len(drifts) != 1 || !drifts[0].IsNew {
		t.Errorf("expected new field drift event: %+v", drifts)
	}
}
