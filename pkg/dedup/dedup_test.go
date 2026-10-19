package dedup

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestWindowedDeduplicator(t *testing.T) {
	d := NewWindowedDeduplicator(500 * time.Millisecond)
	r1 := &model.Record{Level: model.LevelError, Message: "database timeout"}
	r2 := &model.Record{Level: model.LevelError, Message: "database timeout"}
	r3 := &model.Record{Level: model.LevelInfo, Message: "user login"}

	isNew1, ev1 := d.Ingest(r1)
	if !isNew1 || ev1.Count != 1 {
		t.Errorf("r1 must be marked as new")
	}

	isNew2, ev2 := d.Ingest(r2)
	if isNew2 || ev2.Count != 2 {
		t.Errorf("r2 must be deduplicated with count 2")
	}

	isNew3, _ := d.Ingest(r3)
	if !isNew3 {
		t.Errorf("r3 must be marked as new")
	}

	flushed := d.Flush()
	if len(flushed) != 2 {
		t.Errorf("expected 2 unique events flushed, got %d", len(flushed))
	}
}
