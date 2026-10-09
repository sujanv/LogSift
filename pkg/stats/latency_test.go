package stats

import (
	"testing"
	"time"
)

func TestLatencyTracker(t *testing.T) {
	tracker := NewLatencyTracker()

	for i := 1; i <= 100; i++ {
		tracker.Add(time.Duration(i) * time.Millisecond)
	}

	s := tracker.Snapshot()
	if s.Count != 100 {
		t.Errorf("expected count 100, got %d", s.Count)
	}
	if s.Min != 1*time.Millisecond {
		t.Errorf("expected min 1ms, got %v", s.Min)
	}
	if s.Max != 100*time.Millisecond {
		t.Errorf("expected max 100ms, got %v", s.Max)
	}
	if s.P50 != 50*time.Millisecond {
		t.Errorf("expected p50 50ms, got %v", s.P50)
	}
	if s.P95 != 95*time.Millisecond {
		t.Errorf("expected p95 95ms, got %v", s.P95)
	}
	if s.P99 != 99*time.Millisecond {
		t.Errorf("expected p99 99ms, got %v", s.P99)
	}
}
