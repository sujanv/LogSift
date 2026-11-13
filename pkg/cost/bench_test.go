package cost

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func BenchmarkCostIngest(b *testing.B) {
	ce := NewCostEstimator()
	rec := &model.Record{
		Level:   model.LevelInfo,
		Message: "benchmark cost line",
		Fields:  map[string]interface{}{"service": "orders"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ce.Ingest(rec, 256)
	}
}
