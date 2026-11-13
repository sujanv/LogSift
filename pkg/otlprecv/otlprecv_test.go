package otlprecv

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestParseOTLPLogsJSON(t *testing.T) {
	payload := []byte(`{
		"resourceLogs": [{
			"scopeLogs": [{
				"logRecords": [{
					"timeUnixNano": "1731530000000000000",
					"severityNumber": 17,
					"severityText": "ERROR",
					"body": {"stringValue": "database pool exhausted"},
					"traceId": "4bf92f3577b34da6a3ce929d0e0e4736",
					"spanId": "00f067aa0ba902b7",
					"attributes": [
						{"key": "service.name", "value": {"stringValue": "auth-api"}}
					]
				}]
			}]
		}]
	}`)

	records, err := ParseOTLPLogsJSON(payload)
	if err != nil {
		t.Fatalf("unexpected error parsing OTLP: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	r := records[0]
	if r.Level != model.LevelError {
		t.Errorf("expected ERROR level, got %v", r.Level)
	}
	if r.Message != "database pool exhausted" {
		t.Errorf("message mismatch: %s", r.Message)
	}
	if r.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Errorf("traceID mismatch: %s", r.TraceID)
	}
	if r.Fields["service.name"] != "auth-api" {
		t.Errorf("service.name attribute missing")
	}
}
