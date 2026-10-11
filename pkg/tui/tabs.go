package tui

import "github.com/charmbracelet/lipgloss"

var (
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	InactiveTabStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#282A36")).
			Foreground(lipgloss.Color("#6272A4")).
			Padding(0, 1)
)

func RenderTabs(active int, titles []string) string {
	var rendered []string
	for i, t := range titles {
		if i == active {
			rendered = append(rendered, ActiveTabStyle.Render(t))
		} else {
			rendered = append(rendered, InactiveTabStyle.Render(t))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}
