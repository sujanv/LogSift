package tui

import (
	"fmt"
	"strings"
)

func (m AppModel) renderStatsView() string {
	var sb strings.Builder

	sb.WriteString(HeaderStyle.Render("Live Analytics Dashboard (Press Tab or Esc to return)") + "\n\n")

	snap := m.latTracker.Snapshot()

	sb.WriteString("Latency Breakdown:\n")
	sb.WriteString(fmt.Sprintf("  Min:   %8s   |   Mean: %8s\n", snap.Min, snap.Mean))
	sb.WriteString(fmt.Sprintf("  p50:   %8s   |   p75:  %8s\n", snap.P50, snap.P75))
	sb.WriteString(fmt.Sprintf("  p90:   %8s   |   p95:  %8s\n", snap.P90, snap.P95))
	sb.WriteString(fmt.Sprintf("  p99:   %8s   |   Max:  %8s\n\n", snap.P99, snap.Max))

	var errCount int
	for _, r := range m.records {
		if r.Level >= 4 { // Error or Fatal
			errCount++
		}
	}

	errRate := 0.0
	if len(m.records) > 0 {
		errRate = (float64(errCount) / float64(len(m.records))) * 100.0
	}

	sb.WriteString(fmt.Sprintf("Total Logs: %d | Total Errors: %d (%.2f%%)\n", len(m.records), errCount, errRate))

	return sb.String()
}
