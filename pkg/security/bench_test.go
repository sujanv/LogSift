package security

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func BenchmarkThreatDetector(b *testing.B) {
	td := NewThreatDetector()
	rec := &model.Record{
		Timestamp: time.Now(),
		Path:      "/api/v1/search?q=normal+query+parameter",
		Message:   "GET /api/v1/search HTTP/1.1",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = td.InspectRecord(rec)
	}
}
