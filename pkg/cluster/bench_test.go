package cluster

import "testing"

func BenchmarkDrainClusterer(b *testing.B) {
	c := NewDrainClusterer(0.5)
	msg := "Database connection timeout after 3000 ms to host 10.0.0.1"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = c.Add(msg)
	}
}
