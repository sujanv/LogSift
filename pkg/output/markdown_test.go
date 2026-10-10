package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/stats"
)

func TestMarkdownReportGenerator(t *testing.T) {
	var buf bytes.Buffer
	g := NewMarkdownReportGenerator(&buf)

	lat := stats.LatencySummary{
		Count: 500,
		Min:   10 * time.Millisecond,
		Mean:  55 * time.Millisecond,
		P50:   50 * time.Millisecond,
		P95:   120 * time.Millisecond,
		P99:   250 * time.Millisecond,
		Max:   500 * time.Millisecond,
	}

	errClusters := []*stats.ErrorCluster{
		{Fingerprint: "database error on <IP>", Count: 42, Sample: "database error on 10.0.0.1"},
	}

	statusMap := map[string]int64{"2xx": 450, "5xx": 42}
	slowEndpoints := []*stats.EndpointStats{
		{Path: "/api/checkout", Count: 100, P95Latency: 120 * time.Millisecond},
	}

	err := g.Generate(500, time.Second, lat, errClusters, statusMap, slowEndpoints)
	if err != nil {
		t.Fatalf("markdown generation error: %v", err)
	}

	res := buf.String()
	if !strings.Contains(res, "# LogSift Analysis Report") {
		t.Errorf("missing header in report")
	}
	if !strings.Contains(res, "p99") || !strings.Contains(res, "database error on <IP>") {
		t.Errorf("missing stats in report")
	}
}
