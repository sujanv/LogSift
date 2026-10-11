package tui

import (
	"regexp"

	"github.com/charmbracelet/lipgloss"
)

var matchHighlightStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#FFB86C")).
	Foreground(lipgloss.Color("#282A36")).
	Bold(true)

// HighlightMatch applies visual background styling to regex matches within text.
func HighlightMatch(text string, query string) string {
	if query == "" {
		return text
	}

	re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(query))
	if err != nil {
		return text
	}

	return re.ReplaceAllStringFunc(text, func(m string) string {
		return matchHighlightStyle.Render(m)
	})
}
