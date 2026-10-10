package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sujanv/LogSift/pkg/model"
)

// Run launches the interactive Terminal User Interface.
func Run(records []*model.Record) error {
	m := NewAppModel(records)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if !m.isSearching {
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	}
	return m, nil
}

func (m AppModel) View() string {
	if !m.ready {
		return "Initializing LogSift TUI..."
	}
	header := HeaderStyle.Render("LogSift Interactive Explorer")
	status := StatusStyle.Render(fmt.Sprintf("Records: %d | View: %d", len(m.records), m.activeView))
	return fmt.Sprintf("%s\n\nUse 'q' to quit\n\n%s", header, status)
}
