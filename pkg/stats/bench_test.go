package stats

import (
	"testing"
	"time"
)

func BenchmarkLatencyTracker(b *testing.B) {
	tracker := NewLatencyTracker()
	dur := 150 * time.Millisecond

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tracker.Add(dur)
	}
}

func BenchmarkNormalizeMessage(b *testing.B) {
	msg := "failed connection to 10.0.0.5:8080 for user 12948 uuid 123e4567-e89b-12d3-a456-426614174000"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = NormalizeMessage(msg)
	}
}
