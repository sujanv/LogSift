package tui

import (
	"fmt"
	"strings"

	"github.com/sujanv/LogSift/pkg/model"
)

func (m *AppModel) applyFilters() {
	query := strings.TrimSpace(m.searchBar.Value())
	m.filtered = m.filtered[:0]

	for _, r := range m.records {
		if m.selectedLevel != model.LevelUnknown && r.Level < m.selectedLevel {
			continue
		}

		if query != "" {
			msgMatch := strings.Contains(strings.ToLower(r.Message), strings.ToLower(query))
			pathMatch := strings.Contains(strings.ToLower(r.Path), strings.ToLower(query))
			if !msgMatch && !pathMatch {
				continue
			}
		}

		m.filtered = append(m.filtered, r)
	}

	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

func (m AppModel) renderLogTable() string {
	var sb strings.Builder

	visibleRows := m.height - 6
	if visibleRows < 5 {
		visibleRows = 5
	}

	start := max(0, m.cursor-visibleRows/2)
	end := min(len(m.filtered), start+visibleRows)
	if end-start < visibleRows && start > 0 {
		start = max(0, end-visibleRows)
	}

	for i := start; i < end; i++ {
		r := m.filtered[i]
		tsStr := r.Timestamp.Format("15:04:05.000")
		if r.Timestamp.IsZero() {
			tsStr = "  --:--:-- "
		}

		lvlBadge := m.renderLevelBadge(r.Level)
		latStr := fmt.Sprintf("%5s", r.Latency.Round(1).String())
		if r.Latency <= 0 {
			latStr = "     "
		}

		msg := r.Message
		if len(msg) > 60 {
			msg = msg[:57] + "..."
		}

		line := fmt.Sprintf("%s %s %s %s", tsStr, lvlBadge, latStr, msg)
		if i == m.cursor {
			sb.WriteString(SelectedRowStyle.Render(line) + "\n")
		} else {
			sb.WriteString(line + "\n")
		}
	}

	return sb.String()
}

func (m AppModel) renderLevelBadge(lvl model.LogLevel) string {
	s := fmt.Sprintf("[%-5s]", lvl.String())
	switch lvl {
	case model.LevelFatal:
		return LevelFatalStyle.Render(s)
	case model.LevelError:
		return LevelErrorStyle.Render(s)
	case model.LevelWarn:
		return LevelWarnStyle.Render(s)
	case model.LevelInfo:
		return LevelInfoStyle.Render(s)
	case model.LevelDebug:
		return LevelDebugStyle.Render(s)
	default:
		return HelpStyle.Render(s)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
