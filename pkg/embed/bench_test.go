package embed

import "testing"

func BenchmarkTransform(b *testing.B) {
	v := NewVectorizer()
	v.Train([]string{"error timeout socket read redis connection pool"})
	msg := "error timeout on redis host 10.1.2.3"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = v.Transform(msg)
	}
}
