package filter

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestTimeRangeMatcher(t *testing.T) {
	t1 := time.Date(2026, 10, 9, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)

	matcher := NewTimeRangeMatcher(t1, t2)

	recIn := &model.Record{Timestamp: time.Date(2026, 10, 9, 11, 0, 0, 0, time.UTC)}
	if !matcher.Match(recIn) {
		t.Errorf("expected 11:00 to match between 10:00 and 12:00")
	}

	recBefore := &model.Record{Timestamp: time.Date(2026, 10, 9, 9, 0, 0, 0, time.UTC)}
	if matcher.Match(recBefore) {
		t.Errorf("expected 09:00 NOT to match")
	}

	recAfter := &model.Record{Timestamp: time.Date(2026, 10, 9, 13, 0, 0, 0, time.UTC)}
	if matcher.Match(recAfter) {
		t.Errorf("expected 13:00 NOT to match")
	}
}

func TestParseSince(t *testing.T) {
	ref := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)

	from2h, err := ParseSince(ref, "2h")
	if err != nil || !from2h.Equal(ref.Add(-2*time.Hour)) {
		t.Errorf("ParseSince('2h') failed: %v", from2h)
	}

	from3d, err := ParseSince(ref, "3d")
	if err != nil || !from3d.Equal(ref.AddDate(0, 0, -3)) {
		t.Errorf("ParseSince('3d') failed: %v", from3d)
	}
}
