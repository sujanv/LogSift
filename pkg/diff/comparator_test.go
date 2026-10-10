package diff

import (
	"strings"
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/stats"
)

func TestCompare(t *testing.T) {
	base := &Snapshot{
		TotalLogs:  1000,
		ErrorCount: 10,
		Latency: stats.LatencySummary{
			P50: 20 * time.Millisecond,
			P95: 50 * time.Millisecond,
		},
		Errors: []*stats.ErrorCluster{
			{Fingerprint: "old database timeout"},
		},
	}

	target := &Snapshot{
		TotalLogs:  1000,
		ErrorCount: 30, // increased to 3%
		Latency: stats.LatencySummary{
			P50: 30 * time.Millisecond, // +10ms regression
			P95: 100 * time.Millisecond, // +50ms regression
		},
		Errors: []*stats.ErrorCluster{
			{Fingerprint: "new connection pool exhaustion"},
		},
	}

	diff := Compare(base, target)

	if diff.ErrorRateDiff <= 0 {
		t.Errorf("expected positive error rate regression delta, got %f", diff.ErrorRateDiff)
	}
	if diff.P50Delta != 10*time.Millisecond {
		t.Errorf("expected +10ms p50 delta, got %v", diff.P50Delta)
	}
	if len(diff.NewErrors) != 1 || diff.NewErrors[0] != "new connection pool exhaustion" {
		t.Errorf("expected new error detected")
	}
	if len(diff.ResolvedErrors) != 1 || diff.ResolvedErrors[0] != "old database timeout" {
		t.Errorf("expected old error resolved")
	}

	formatted := diff.FormatDelta()
	if !strings.Contains(formatted, "Regression") {
		t.Errorf("expected regression notes in formatted delta: %s", formatted)
	}
}
