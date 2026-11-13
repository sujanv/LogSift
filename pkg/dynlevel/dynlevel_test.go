package dynlevel

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestFlightRecorder(t *testing.T) {
	fr := NewFlightRecorder(3)

	r1 := &model.Record{Level: model.LevelDebug, Message: "step 1"}
	r2 := &model.Record{Level: model.LevelDebug, Message: "step 2"}
	r3 := &model.Record{Level: model.LevelDebug, Message: "step 3"}
	r4 := &model.Record{Level: model.LevelDebug, Message: "step 4"}

	fr.Record(r1)
	fr.Record(r2)
	fr.Record(r3)
	fr.Record(r4)

	if fr.Size() != 3 {
		t.Fatalf("expected ring size 3, got %d", fr.Size())
	}

	fault := &model.Record{Level: model.LevelError, Message: "fatal nil pointer"}
	flushed := fr.FlushOnFault(fault)

	if len(flushed) != 4 {
		t.Fatalf("expected 4 flushed records (3 debug + 1 fault), got %d", len(flushed))
	}
	if flushed[0].Message != "step 2" {
		t.Errorf("expected step 2 as oldest preserved item, got %s", flushed[0].Message)
	}
	if flushed[3].Message != "fatal nil pointer" {
		t.Errorf("expected fault record at end")
	}
}
