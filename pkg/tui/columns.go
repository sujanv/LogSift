package tui

type ColumnPicker struct {
	ShowTimestamp bool
	ShowLevel     bool
	ShowLatency   bool
	ShowPath      bool
	ShowStatus    bool
	ShowIP        bool
}

func DefaultColumnPicker() ColumnPicker {
	return ColumnPicker{
		ShowTimestamp: true,
		ShowLevel:     true,
		ShowLatency:   true,
		ShowPath:      true,
		ShowStatus:    true,
		ShowIP:        false,
	}
}
