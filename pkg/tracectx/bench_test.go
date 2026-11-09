package tracectx

import "testing"

func BenchmarkParseTraceParent(b *testing.B) {
	header := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseTraceParent(header)
	}
}
