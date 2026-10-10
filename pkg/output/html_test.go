package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/stats"
)

func TestHTMLReportGenerator(t *testing.T) {
	var buf bytes.Buffer
	g := NewHTMLReportGenerator(&buf)

	lat := stats.LatencySummary{
		P50: 15 * time.Millisecond,
		P95: 80 * time.Millisecond,
		P99: 150 * time.Millisecond,
	}
	errs := []*stats.ErrorCluster{
		{Fingerprint: "db connection failed", Count: 12, Sample: "db connection failed on 10.0.0.1"},
	}

	err := g.Generate(1000, 2*time.Second, lat, errs, nil)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "<!DOCTYPE html>") || !strings.Contains(out, "LogSift Analytics Report") {
		t.Errorf("expected HTML structure")
	}
}

func TestRenderLatencySVG(t *testing.T) {
	svg := RenderLatencySVG(10*time.Millisecond, 20*time.Millisecond, 50*time.Millisecond, 90*time.Millisecond, 150*time.Millisecond)
	if !strings.HasPrefix(svg, "<svg") || !strings.HasSuffix(svg, "</svg>") {
		t.Errorf("expected valid SVG container")
	}
	if !strings.Contains(svg, "p99") {
		t.Errorf("expected p99 label in SVG")
	}
}
