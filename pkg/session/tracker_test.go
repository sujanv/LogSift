package session

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestSessionTracker(t *testing.T) {
	st := NewSessionTracker(10 * time.Minute)
	t0 := time.Now()

	r1 := &model.Record{IP: "1.2.3.4", Path: "/home", Timestamp: t0}
	r2 := &model.Record{IP: "1.2.3.4", Path: "/pricing", Timestamp: t0.Add(2 * time.Minute)}
	r3 := &model.Record{IP: "5.6.7.8", Path: "/blog", Timestamp: t0}

	st.Ingest(r1)
	st.Ingest(r2)
	st.Ingest(r3)

	sessions := st.CompleteFinalizes()
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}

	analytics := AnalyzeSessions(sessions)
	if analytics.TotalSessions != 2 {
		t.Errorf("expected 2 sessions analyzed")
	}
	if analytics.BounceRate != 50.0 {
		t.Errorf("expected 50%% bounce rate, got %f", analytics.BounceRate)
	}
}
