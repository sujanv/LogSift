package filter

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func BenchmarkFieldMatcher(b *testing.B) {
	m := NewFieldMatcher("status", OpGte, "500")
	rec := &model.Record{StatusCode: 503}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = m.Match(rec)
	}
}

func BenchmarkTextMatcher(b *testing.B) {
	m := NewTextMatcher("connection refused", false)
	rec := &model.Record{Message: "critical error: connection refused by database"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = m.Match(rec)
	}
}
