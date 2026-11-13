package dynlevel

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func BenchmarkFlightRecorderAppend(b *testing.B) {
	fr := NewFlightRecorder(128)
	rec := &model.Record{Level: model.LevelDebug, Message: "trace benchmark debug"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		fr.Record(rec)
	}
}
