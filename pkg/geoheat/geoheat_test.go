package geoheat

import (
	"testing"
	"time"
)

func TestRegionalHeatmap(t *testing.T) {
	hm := NewRegionalHeatmap(200 * time.Millisecond)

	for i := 0; i < 100; i++ {
		hm.Record("US", 50*time.Millisecond)
	}
	for i := 0; i < 20; i++ {
		hm.Record("JP", 450*time.Millisecond)
	}

	summary := hm.Summary()
	if len(summary) != 2 {
		t.Fatalf("expected 2 countries in summary, got %d", len(summary))
	}

	var jpMetric *RegionMetric
	for _, m := range summary {
		if m.Country == "JP" {
			jpMetric = &m
			break
		}
	}

	if jpMetric == nil || !jpMetric.Degraded {
		t.Errorf("expected JP to be marked as degraded (>200ms latency)")
	}
}
