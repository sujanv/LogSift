package filter

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

// Op represents comparison operators.
type Op string

const (
	OpEq       Op = "="
	OpNeq      Op = "!="
	OpGt       Op = ">"
	OpGte      Op = ">="
	OpLt       Op = "<"
	OpLte      Op = "<="
	OpContains Op = "contains"
	OpIn       Op = "in"
)

// FieldMatcher evaluates field predicates on a log Record.
type FieldMatcher struct {
	Field    string
	Operator Op
	Value    string

	// Cached parsed values
	numVal      float64
	durationVal time.Duration
	isNum       bool
	isDuration  bool
}

// NewFieldMatcher creates a new FieldMatcher.
func NewFieldMatcher(field string, op Op, val string) *FieldMatcher {
	m := &FieldMatcher{
		Field:    strings.TrimSpace(field),
		Operator: op,
		Value:    strings.TrimSpace(val),
	}

	if f, err := strconv.ParseFloat(m.Value, 64); err == nil {
		m.numVal = f
		m.isNum = true
	}
	if d, err := time.ParseDuration(m.Value); err == nil {
		m.durationVal = d
		m.isDuration = true
	}

	return m
}

// Match returns true if record satisfies the field comparison.
func (m *FieldMatcher) Match(record *model.Record) bool {
	fieldVal := m.extractField(record)

	switch m.Operator {
	case OpEq:
		if m.isNum {
			if vNum, ok := toFloat(fieldVal); ok {
				return vNum == m.numVal
			}
		}
		return strings.EqualFold(fmt.Sprintf("%v", fieldVal), m.Value)

	case OpNeq:
		if m.isNum {
			if vNum, ok := toFloat(fieldVal); ok {
				return vNum != m.numVal
			}
		}
		return !strings.EqualFold(fmt.Sprintf("%v", fieldVal), m.Value)

	case OpGt:
		if m.isDuration && record.Latency > 0 {
			return record.Latency > m.durationVal
		}
		if vNum, ok := toFloat(fieldVal); ok && m.isNum {
			return vNum > m.numVal
		}
		return fmt.Sprintf("%v", fieldVal) > m.Value

	case OpGte:
		if m.isDuration && record.Latency > 0 {
			return record.Latency >= m.durationVal
		}
		if vNum, ok := toFloat(fieldVal); ok && m.isNum {
			return vNum >= m.numVal
		}
		return fmt.Sprintf("%v", fieldVal) >= m.Value

	case OpLt:
		if m.isDuration && record.Latency > 0 {
			return record.Latency < m.durationVal
		}
		if vNum, ok := toFloat(fieldVal); ok && m.isNum {
			return vNum < m.numVal
		}
		return fmt.Sprintf("%v", fieldVal) < m.Value

	case OpLte:
		if m.isDuration && record.Latency > 0 {
			return record.Latency <= m.durationVal
		}
		if vNum, ok := toFloat(fieldVal); ok && m.isNum {
			return vNum <= m.numVal
		}
		return fmt.Sprintf("%v", fieldVal) <= m.Value

	case OpContains:
		return strings.Contains(strings.ToLower(fmt.Sprintf("%v", fieldVal)), strings.ToLower(m.Value))

	case OpIn:
		items := strings.Split(m.Value, ",")
		target := fmt.Sprintf("%v", fieldVal)
		for _, item := range items {
			if strings.EqualFold(strings.TrimSpace(item), target) {
				return true
			}
		}
		return false
	}

	return false
}

func (m *FieldMatcher) extractField(r *model.Record) interface{} {
	switch strings.ToLower(m.Field) {
	case "level":
		return r.Level.String()
	case "message", "msg":
		return r.Message
	case "status", "status_code":
		return r.StatusCode
	case "latency", "latency_ms":
		return float64(r.Latency.Milliseconds())
	case "method":
		return r.Method
	case "path":
		return r.Path
	case "ip":
		return r.IP
	case "trace_id":
		return r.TraceID
	}

	if r.Fields != nil {
		if val, exists := r.Fields[m.Field]; exists {
			return val
		}
	}

	return ""
}

func toFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}
