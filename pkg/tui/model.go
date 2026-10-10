package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/stats"
)

type ViewMode int

const (
	ViewLogs ViewMode = iota
	ViewDetail
	ViewStats
)

// AppModel holds the state of the interactive LogSift TUI.
type AppModel struct {
	records       []*model.Record
	filtered      []*model.Record
	cursor        int
	viewport      viewport.Model
	searchBar     textinput.Model
	isSearching   bool
	activeView    ViewMode
	selectedLevel model.LogLevel
	errClusters   []*stats.ErrorCluster
	latTracker    *stats.LatencyTracker
	width         int
	height        int
	ready         bool
}

// NewAppModel initializes an AppModel.
func NewAppModel(records []*model.Record) AppModel {
	ti := textinput.New()
	ti.Placeholder = "Type filter query... (Enter to apply, Esc to cancel)"
	ti.CharLimit = 150
	ti.Width = 50

	m := AppModel{
		records:       records,
		filtered:      records,
		cursor:        0,
		searchBar:     ti,
		activeView:    ViewLogs,
		selectedLevel: model.LevelUnknown,
		latTracker:    stats.NewLatencyTracker(),
	}

	for _, r := range records {
		if r.Latency > 0 {
			m.latTracker.Add(r.Latency)
		}
	}

	return m
}

func (m AppModel) Init() tea.Cmd {
	return nil
}
