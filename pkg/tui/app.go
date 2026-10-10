package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.isSearching {
			switch msg.String() {
			case "enter":
				m.isSearching = false
				m.applyFilters()
				return m, nil
			case "esc":
				m.isSearching = false
				m.searchBar.SetValue("")
				m.applyFilters()
				return m, nil
			default:
				m.searchBar, cmd = m.searchBar.Update(msg)
				m.applyFilters()
				return m, cmd
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "/":
			m.isSearching = true
			m.searchBar.Focus()
			return m, textinput.Blink
		case "1":
			m.selectedLevel = model.LevelUnknown
			m.applyFilters()
		case "2":
			m.selectedLevel = model.LevelDebug
			m.applyFilters()
		case "3":
			m.selectedLevel = model.LevelInfo
			m.applyFilters()
		case "4":
			m.selectedLevel = model.LevelWarn
			m.applyFilters()
		case "5":
			m.selectedLevel = model.LevelError
			m.applyFilters()
		case "enter", " ":
			if len(m.filtered) > 0 {
				if m.activeView == ViewDetail {
					m.activeView = ViewLogs
				} else {
					m.activeView = ViewDetail
				}
			}
		case "tab", "s":
			if m.activeView == ViewStats {
				m.activeView = ViewLogs
			} else {
				m.activeView = ViewStats
			}
		case "esc":
			m.activeView = ViewLogs
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

	header := HeaderStyle.Render("LogSift Interactive Log Explorer")
	help := HelpStyle.Render(" [j/k] Navigate  [/] Search  [1-5] Level  [Enter] Details  [Tab] Analytics  [q] Quit")

	var body string
	switch m.activeView {
	case ViewDetail:
		body = m.renderDetailView()
	case ViewStats:
		body = m.renderStatsView()
	default:
		if m.isSearching {
			body = m.searchBar.View() + "\n\n" + m.renderLogTable()
		} else {
			body = m.renderLogTable()
		}
	}

	status := StatusStyle.Render(fmt.Sprintf("Showing %d/%d records | Level: %s", len(m.filtered), len(m.records), m.selectedLevel.String()))

	return fmt.Sprintf("%s%s\n\n%s\n%s", header, help, body, status)
}
