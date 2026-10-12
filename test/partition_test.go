package test

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/archive"
	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/session"
)

func TestPartitionAndSessionIntegration(t *testing.T) {
	st := session.NewSessionTracker(5 * time.Minute)
	now := time.Now()
	st.Ingest(&model.Record{IP: "10.0.0.1", Path: "/a", Timestamp: now})
	st.Ingest(&model.Record{IP: "10.0.0.1", Path: "/b", Timestamp: now.Add(time.Minute)})
	sessions := st.CompleteFinalizes()
	if len(sessions) != 1 || sessions[0].Requests != 2 {
		t.Errorf("session tracking mismatch")
	}

	p := archive.DeterminePartitionPath(archive.PartitionByLevel, &model.Record{Level: model.LevelError})
	if p != "ERROR" {
		t.Errorf("expected ERROR partition, got %s", p)
	}
}
