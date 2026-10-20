package dedup

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func BenchmarkDeduplicator(b *testing.B) {
	d := NewWindowedDeduplicator(10 * time.Second)
	r := &model.Record{Level: model.LevelError, Message: "connection refused: 127.0.0.1:5432"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = d.Ingest(r)
	}
}
