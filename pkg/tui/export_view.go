package tui

import "strings"

func RenderSparkline(values []int) string {
	if len(values) == 0 {
		return ""
	}
	bars := []rune{' ', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	max := 1
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	var sb strings.Builder
	for _, v := range values {
		idx := (v * (len(bars) - 1)) / max
		if idx < 0 {
			idx = 0
		}
		if idx >= len(bars) {
			idx = len(bars) - 1
		}
		sb.WriteRune(bars[idx])
	}
	return sb.String()
}
