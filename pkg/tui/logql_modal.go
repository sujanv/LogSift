package tui

type LogQLModal struct {
	Visible bool
	Input   string
	History []string
}

func (m *LogQLModal) Toggle() bool {
	m.Visible = !m.Visible
	return m.Visible
}
