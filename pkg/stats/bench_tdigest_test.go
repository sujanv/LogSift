package stats

import "testing"

func BenchmarkTDigestQuantile(b *testing.B) {
	td := NewTDigest(100)
	for i := 1; i <= 500; i++ {
		td.Add(float64(i))
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = td.Quantile(0.99)
	}
}
