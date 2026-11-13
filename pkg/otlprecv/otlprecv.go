package otlprecv

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type otlpKeyValue struct {
	Key   string          `json:"key"`
	Value otlpValueHolder `json:"value"`
}

type otlpValueHolder struct {
	StringValue *string  `json:"stringValue,omitempty"`
	IntValue    *string  `json:"intValue,omitempty"`
	DoubleValue *float64 `json:"doubleValue,omitempty"`
	BoolValue   *bool    `json:"boolValue,omitempty"`
}

type otlpLogRecord struct {
	TimeUnixNano   string          `json:"timeUnixNano"`
	SeverityNumber int             `json:"severityNumber"`
	SeverityText   string          `json:"severityText"`
	Body           otlpValueHolder `json:"body"`
	Attributes     []otlpKeyValue  `json:"attributes"`
	TraceID        string          `json:"traceId"`
	SpanID         string          `json:"spanId"`
}

type otlpScopeLogs struct {
	LogRecords []otlpLogRecord `json:"logRecords"`
}

type otlpResourceLogs struct {
	ScopeLogs []otlpScopeLogs `json:"scopeLogs"`
}

type otlpPayload struct {
	ResourceLogs []otlpResourceLogs `json:"resourceLogs"`
}

func ParseOTLPLogsJSON(body []byte) ([]*model.Record, error) {
	if len(body) == 0 {
		return nil, errors.New("empty OTLP JSON payload")
	}

	var payload otlpPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	var records []*model.Record
	for _, rl := range payload.ResourceLogs {
		for _, sl := range rl.ScopeLogs {
			for _, lr := range sl.LogRecords {
				rec := &model.Record{
					Fields: make(map[string]interface{}),
				}

				if lr.TimeUnixNano != "" {
					if nano, err := strconv.ParseInt(lr.TimeUnixNano, 10, 64); err == nil {
						rec.Timestamp = time.Unix(0, nano)
					}
				}

				switch {
				case lr.SeverityNumber >= 17:
					rec.Level = model.LevelError
				case lr.SeverityNumber >= 13:
					rec.Level = model.LevelWarn
				case lr.SeverityNumber >= 9:
					rec.Level = model.LevelInfo
				case lr.SeverityNumber >= 1:
					rec.Level = model.LevelDebug
				default:
					rec.Level = model.ParseLevel(lr.SeverityText)
				}

				if lr.Body.StringValue != nil {
					rec.Message = *lr.Body.StringValue
				}
				rec.TraceID = lr.TraceID
				rec.SpanID = lr.SpanID

				for _, attr := range lr.Attributes {
					if attr.Value.StringValue != nil {
						rec.Fields[attr.Key] = *attr.Value.StringValue
					} else if attr.Value.DoubleValue != nil {
						rec.Fields[attr.Key] = *attr.Value.DoubleValue
					} else if attr.Value.BoolValue != nil {
						rec.Fields[attr.Key] = *attr.Value.BoolValue
					}
				}

				records = append(records, rec)
			}
		}
	}

	return records, nil
}
