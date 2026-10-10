package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Palette
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	StatusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#353535")).
			Padding(0, 1)

	LevelDebugStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7D7"))
	LevelInfoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#5AF78E"))
	LevelWarnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F3F99D"))
	LevelErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	LevelFatalStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF55FF")).Bold(true)

	SelectedRowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#44475A")).
				Foreground(lipgloss.Color("#F8F8F2")).
				Bold(true)

	HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4"))
)
