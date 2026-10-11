package export

import (
	"strings"
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/stats"
)

func TestFormatPrometheusMetrics(t *testing.T) {
	lat := stats.LatencySummary{
		P50:   25 * time.Millisecond,
		P99:   120 * time.Millisecond,
		Count: 500,
	}
	statusMap := map[string]int64{"2xx": 450, "5xx": 50}

	prom := FormatPrometheusMetrics(500, 50, lat, statusMap)

	if !strings.Contains(prom, "logsift_records_total 500") {
		t.Errorf("missing total records counter")
	}
	if !strings.Contains(prom, `quantile="0.99"} 0.1200`) {
		t.Errorf("missing p99 quantile metric")
	}
	if !strings.Contains(prom, `code="5xx"} 50`) {
		t.Errorf("missing 5xx status metric")
	}
}

func TestConvertToOTLP(t *testing.T) {
	r := &model.Record{
		Timestamp:  time.Now(),
		Level:      model.LevelError,
		Message:    "database connection refused",
		StatusCode: 500,
		Path:       "/api/pay",
		TraceID:    "trace-abc-123",
	}

	data, err := FormatOTLPJSON(r)
	if err != nil {
		t.Fatalf("FormatOTLPJSON error: %v", err)
	}

	str := string(data)
	if !strings.Contains(str, `"severity_text":"ERROR"`) || !strings.Contains(str, `"trace_id":"trace-abc-123"`) {
		t.Errorf("OTLP mapping failed: %s", str)
	}
}
