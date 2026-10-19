package tui

import "time"

type TimelineZoom struct {
	Start  time.Time
	End    time.Time
	Zoomed bool
}

func (tz *TimelineZoom) SetRange(start, end time.Time) {
	tz.Start = start
	tz.End = end
	tz.Zoomed = true
}

func (tz *TimelineZoom) Reset() {
	tz.Zoomed = false
}
