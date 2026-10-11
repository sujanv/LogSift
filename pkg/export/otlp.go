package export

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type OTLPLogRecord struct {
	Timestamp          string                 `json:"time_unix_nano"`
	SeverityNumber     int                    `json:"severity_number"`
	SeverityText       string                 `json:"severity_text"`
	Body               map[string]interface{} `json:"body"`
	Attributes         map[string]interface{} `json:"attributes"`
	TraceID            string                 `json:"trace_id,omitempty"`
	SpanID             string                 `json:"span_id,omitempty"`
}

func ConvertToOTLP(r *model.Record) *OTLPLogRecord {
	nano := r.Timestamp.UnixNano()
	if r.Timestamp.IsZero() {
		nano = time.Now().UnixNano()
	}

	attrs := make(map[string]interface{})
	if r.Path != "" {
		attrs["http.target"] = r.Path
	}
	if r.Method != "" {
		attrs["http.method"] = r.Method
	}
	if r.StatusCode > 0 {
		attrs["http.status_code"] = r.StatusCode
	}
	if r.Latency > 0 {
		attrs["duration_ms"] = float64(r.Latency.Milliseconds())
	}
	if r.IP != "" {
		attrs["net.peer.ip"] = r.IP
	}

	return &OTLPLogRecord{
		Timestamp:      fmt.Sprintf("%d", nano),
		SeverityNumber: int(r.Level) * 4,
		SeverityText:   r.Level.String(),
		Body:           map[string]interface{}{"stringValue": r.Message},
		Attributes:     attrs,
		TraceID:        r.TraceID,
		SpanID:         r.SpanID,
	}
}

func FormatOTLPJSON(r *model.Record) ([]byte, error) {
	otlp := ConvertToOTLP(r)
	return json.Marshal(otlp)
}
