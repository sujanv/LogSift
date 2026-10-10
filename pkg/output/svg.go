package output

import (
	"fmt"
	"strings"
	"time"
)

// RenderLatencySVG generates an embedded SVG bar chart representing latency percentiles.
func RenderLatencySVG(p50, p75, p90, p95, p99 time.Duration) string {
	maxVal := float64(p99.Milliseconds())
	if maxVal <= 0 {
		maxVal = 1
	}

	bars := []struct {
		Label string
		Val   time.Duration
		Color string
	}{
		{"p50", p50, "#38bdf8"},
		{"p75", p75, "#60a5fa"},
		{"p90", p90, "#facc15"},
		{"p95", p95, "#fb923c"},
		{"p99", p99, "#f43f5e"},
	}

	var sb strings.Builder
	sb.WriteString(`<svg width="500" height="180" xmlns="http://www.w3.org/2000/svg" style="background:#1e293b;border-radius:8px;">`)

	y := 25
	for _, b := range bars {
		w := (float64(b.Val.Milliseconds()) / maxVal) * 320.0
		if w < 4 {
			w = 4
		}
		sb.WriteString(fmt.Sprintf(`<text x="20" y="%d" fill="#94a3b8" font-size="12" font-family="sans-serif">%s</text>`, y+12, b.Label))
		sb.WriteString(fmt.Sprintf(`<rect x="60" y="%d" width="%.1f" height="16" rx="4" fill="%s"/>`, y, w, b.Color))
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%d" fill="#f8fafc" font-size="11" font-family="sans-serif"> %v</text>`, 70+w, y+12, b.Val))
		y += 30
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}
