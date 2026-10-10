package correlate

import (
	"fmt"
	"strings"
	"time"
)

// FormatTraceTimeline formats a correlated trace waterfall into an ASCII timeline.
func FormatTraceTimeline(traceID string, entries []*TraceEntry) string {
	if len(entries) == 0 {
		return fmt.Sprintf("No logs found for trace %s", traceID)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== Trace Waterfall: %s (%d events) ===\n\n", traceID, len(entries)))

	t0 := entries[0].Timestamp
	for i, e := range entries {
		offset := e.Timestamp.Sub(t0).Round(time.Millisecond)
		latStr := ""
		if e.Latency > 0 {
			latStr = fmt.Sprintf(" [%v]", e.Latency)
		}
		pathStr := ""
		if e.Path != "" {
			pathStr = fmt.Sprintf(" (%s %d)", e.Path, e.StatusCode)
		}

		sb.WriteString(fmt.Sprintf("  +%-6v | [%-5s] %s%s%s\n", offset, e.Level, e.Message, pathStr, latStr))
		_ = i
	}

	return sb.String()
}
