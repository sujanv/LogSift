package schema

import (
	"encoding/json"
	"fmt"
	"sync"
)

type FieldType string

const (
	TypeString  FieldType = "string"
	TypeNumber  FieldType = "number"
	TypeBoolean FieldType = "boolean"
	TypeObject  FieldType = "object"
	TypeArray   FieldType = "array"
	TypeNull    FieldType = "null"
)

type FieldSchema struct {
	Name  string              `json:"name"`
	Types map[FieldType]int64 `json:"types"`
	Count int64               `json:"count"`
}

type InferredSchema struct {
	mu        sync.Mutex
	TotalLogs int64
	Fields    map[string]*FieldSchema
}

func NewInferredSchema() *InferredSchema {
	return &InferredSchema{Fields: make(map[string]*FieldSchema)}
}

func (s *InferredSchema) IngestJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalLogs++

	for k, v := range m {
		fs, exists := s.Fields[k]
		if !exists {
			fs = &FieldSchema{
				Name:  k,
				Types: make(map[FieldType]int64),
			}
			s.Fields[k] = fs
		}
		fs.Count++
		t := detectType(v)
		fs.Types[t]++
	}
	return nil
}

func detectType(v interface{}) FieldType {
	if v == nil {
		return TypeNull
	}
	switch v.(type) {
	case string:
		return TypeString
	case float64, int, int64:
		return TypeNumber
	case bool:
		return TypeBoolean
	case map[string]interface{}:
		return TypeObject
	case []interface{}:
		return TypeArray
	default:
		return TypeString
	}
}

func (s *InferredSchema) DetectDrift(data []byte) []string {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return []string{"Malformed JSON"}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var warnings []string
	for k, v := range m {
		fs, exists := s.Fields[k]
		if !exists {
			warnings = append(warnings, fmt.Sprintf("Unexpected new field '%s'", k))
			continue
		}
		t := detectType(v)
		if _, ok := fs.Types[t]; !ok && len(fs.Types) > 0 {
			warnings = append(warnings, fmt.Sprintf("Type drift on '%s': saw %s, expected one of %v", k, t, keys(fs.Types)))
		}
	}
	return warnings
}

func keys(m map[FieldType]int64) []FieldType {
	res := make([]FieldType, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
