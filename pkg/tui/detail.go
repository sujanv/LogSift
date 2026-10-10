package tui

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (m AppModel) renderDetailView() string {
	if len(m.filtered) == 0 || m.cursor >= len(m.filtered) {
		return "No record selected. Press Esc to return."
	}

	r := m.filtered[m.cursor]
	var sb strings.Builder

	sb.WriteString(HeaderStyle.Render("Record Inspector (Press Esc or Enter to return)") + "\n\n")

	ts := "N/A"
	if !r.Timestamp.IsZero() {
		ts = r.Timestamp.Format(time.RFC3339Nano)
	}

	sb.WriteString(fmt.Sprintf("  Line Number:  %d\n", r.LineNum))
	sb.WriteString(fmt.Sprintf("  Timestamp:    %s\n", ts))
	sb.WriteString(fmt.Sprintf("  Severity:     %s\n", m.renderLevelBadge(r.Level)))
	if r.Latency > 0 {
		sb.WriteString(fmt.Sprintf("  Latency:      %v\n", r.Latency))
	}
	if r.StatusCode > 0 {
		sb.WriteString(fmt.Sprintf("  HTTP Status:  %d\n", r.StatusCode))
	}
	if r.Method != "" || r.Path != "" {
		sb.WriteString(fmt.Sprintf("  Endpoint:     %s %s\n", r.Method, r.Path))
	}
	if r.IP != "" {
		sb.WriteString(fmt.Sprintf("  Client IP:    %s\n", r.IP))
	}
	if r.TraceID != "" {
		sb.WriteString(fmt.Sprintf("  Trace ID:     %s\n", r.TraceID))
	}

	sb.WriteString(fmt.Sprintf("\n  Message:\n  %s\n\n", r.Message))

	if len(r.Fields) > 0 {
		sb.WriteString("  Structured Fields:\n")
		formattedJSON, err := json.MarshalIndent(r.Fields, "  ", "  ")
		if err == nil {
			sb.WriteString(string(formattedJSON) + "\n")
		}
	}

	return sb.String()
}
