package anomaly

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func BenchmarkStreamingAnomalyDetector(b *testing.B) {
	det := NewStreamingAnomalyDetector(3.0)
	rec := &model.Record{Level: model.LevelInfo, Timestamp: time.Now()}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = det.Ingest(rec)
	}
}
