package tui

type FlightRecorderView struct {
	Visible bool
	Buffer  []string
}

func (frv *FlightRecorderView) Toggle() bool {
	frv.Visible = !frv.Visible
	return frv.Visible
}
