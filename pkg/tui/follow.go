package tui

type FollowMode struct {
	Enabled bool
}

func (f *FollowMode) Toggle() bool {
	f.Enabled = !f.Enabled
	return f.Enabled
}
