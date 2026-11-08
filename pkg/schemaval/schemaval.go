package schemaval

import (
	"fmt"
	"sync"

	"github.com/sujanv/LogSift/pkg/model"
)

type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeFloat   FieldType = "float"
	TypeBool    FieldType = "bool"
	TypeObject  FieldType = "object"
	TypeArray   FieldType = "array"
	TypeUnknown FieldType = "unknown"
)

type SchemaValidator struct {
	mu             sync.Mutex
	requiredFields []string
	expectedTypes  map[string]FieldType
	observedTypes  map[string]FieldType
}

func NewSchemaValidator(required []string, expected map[string]FieldType) *SchemaValidator {
	if expected == nil {
		expected = make(map[string]FieldType)
	}
	observed := make(map[string]FieldType)
	for k, v := range expected {
		observed[k] = v
	}
	return &SchemaValidator{
		requiredFields: required,
		expectedTypes:  expected,
		observedTypes:  observed,
	}
}

func DetectType(val interface{}) FieldType {
	if val == nil {
		return TypeUnknown
	}
	switch val.(type) {
	case string:
		return TypeString
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return TypeInt
	case float32, float64:
		return TypeFloat
	case bool:
		return TypeBool
	case map[string]interface{}:
		return TypeObject
	case []interface{}:
		return TypeArray
	default:
		return TypeUnknown
	}
}

func (sv *SchemaValidator) Validate(rec *model.Record) []string {
	var errs []string
	if rec == nil {
		return []string{"record is nil"}
	}

	for _, req := range sv.requiredFields {
		found := false
		switch req {
		case "timestamp":
			found = !rec.Timestamp.IsZero()
		case "level":
			found = rec.Level != model.LevelUnknown
		case "message":
			found = rec.Message != ""
		default:
			if rec.Fields != nil {
				_, found = rec.Fields[req]
			}
		}
		if !found {
			errs = append(errs, fmt.Sprintf("missing required field: %s", req))
		}
	}

	if rec.Fields != nil {
		for k, expectedType := range sv.expectedTypes {
			if val, ok := rec.Fields[k]; ok {
				actualType := DetectType(val)
				if actualType != expectedType && actualType != TypeUnknown {
					errs = append(errs, fmt.Sprintf("field %s expected type %s, got %s", k, expectedType, actualType))
				}
			}
		}
	}

	return errs
}

type DriftEvent struct {
	Field    string    `json:"field"`
	OldType  FieldType `json:"old_type"`
	NewType  FieldType `json:"new_type"`
	IsNew    bool      `json:"is_new"`
	Message  string    `json:"message"`
}

func (sv *SchemaValidator) DetectDrift(rec *model.Record) []DriftEvent {
	if rec == nil || rec.Fields == nil {
		return nil
	}

	sv.mu.Lock()
	defer sv.mu.Unlock()

	var drifts []DriftEvent
	for k, v := range rec.Fields {
		act := DetectType(v)
		prev, exists := sv.observedTypes[k]
		if !exists {
			sv.observedTypes[k] = act
			drifts = append(drifts, DriftEvent{
				Field:   k,
				NewType: act,
				IsNew:   true,
				Message: fmt.Sprintf("new uncataloged field introduced: %s (%s)", k, act),
			})
		} else if prev != act && act != TypeUnknown {
			sv.observedTypes[k] = act
			drifts = append(drifts, DriftEvent{
				Field:   k,
				OldType: prev,
				NewType: act,
				IsNew:   false,
				Message: fmt.Sprintf("schema type mutation on field %s: %s -> %s", k, prev, act),
			})
		}
	}

	return drifts
}
