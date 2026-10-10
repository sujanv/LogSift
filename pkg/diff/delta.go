package diff

import (
	"fmt"
	"strings"
	"time"
)

// FormatDelta produces a human-readable comparison report.
func (d *DiffReport) FormatDelta() string {
	var sb strings.Builder
	sb.WriteString("=== LogSift Drift & Regression Analysis ===\n\n")

	// Error Rate
	rateSign := "+"
	if d.ErrorRateDiff < 0 {
		rateSign = ""
	}
	sb.WriteString(fmt.Sprintf("Error Rate Delta:  %s%.2f%%\n", rateSign, d.ErrorRateDiff))

	// Latency Deltas
	sb.WriteString(fmt.Sprintf("Latency p50 Delta: %s\n", formatDurationDelta(d.P50Delta)))
	sb.WriteString(fmt.Sprintf("Latency p95 Delta: %s\n", formatDurationDelta(d.P95Delta)))
	sb.WriteString(fmt.Sprintf("Latency p99 Delta: %s\n\n", formatDurationDelta(d.P99Delta)))

	if len(d.NewErrors) > 0 {
		sb.WriteString("⚠️  New Error Fingerprints Introduced:\n")
		for _, e := range d.NewErrors {
			sb.WriteString(fmt.Sprintf("  + %s\n", e))
		}
		sb.WriteString("\n")
	}

	if len(d.ResolvedErrors) > 0 {
		sb.WriteString("✅ Resolved Error Fingerprints:\n")
		for _, e := range d.ResolvedErrors {
			sb.WriteString(fmt.Sprintf("  - %s\n", e))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func formatDurationDelta(d time.Duration) string {
	if d > 0 {
		return fmt.Sprintf("+%v (Regression)", d)
	} else if d < 0 {
		return fmt.Sprintf("%v (Improvement)", d)
	}
	return "0s (Unchanged)"
}
