package filter

import (
	"strconv"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

// TimeRangeMatcher filters records within [From, To].
type TimeRangeMatcher struct {
	From time.Time
	To   time.Time
}

// NewTimeRangeMatcher creates a TimeRangeMatcher.
func NewTimeRangeMatcher(from, to time.Time) *TimeRangeMatcher {
	return &TimeRangeMatcher{From: from, To: to}
}

// ParseSince parses a relative duration string (e.g. "15m", "2h", "7d") and returns the From timestamp.
func ParseSince(now time.Time, s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "d") {
		daysStr := strings.TrimSuffix(s, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return time.Time{}, err
		}
		return now.AddDate(0, 0, -days), nil
	}

	dur, err := time.ParseDuration(s)
	if err != nil {
		return time.Time{}, err
	}
	return now.Add(-dur), nil
}

// Match returns true if record timestamp is within range.
func (trm *TimeRangeMatcher) Match(record *model.Record) bool {
	if record.Timestamp.IsZero() {
		// If log has no timestamp, don't discard unless strict
		return true
	}

	if !trm.From.IsZero() && record.Timestamp.Before(trm.From) {
		return false
	}
	if !trm.To.IsZero() && record.Timestamp.After(trm.To) {
		return false
	}

	return true
}
