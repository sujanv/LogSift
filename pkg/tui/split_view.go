package tui

type SplitViewMode int

const (
	SplitSingle SplitViewMode = iota
	SplitSideBySide
	SplitStacked
)

type SplitViewController struct {
	Mode SplitViewMode
}

func (s *SplitViewController) Toggle() SplitViewMode {
	if s.Mode == SplitSingle {
		s.Mode = SplitSideBySide
	} else {
		s.Mode = SplitSingle
	}
	return s.Mode
}
