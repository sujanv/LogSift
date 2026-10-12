package geodb

import "testing"

func BenchmarkResolverLookup(b *testing.B) {
	res := NewResolver()
	ip := "8.8.8.8"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = res.Lookup(ip)
	}
}
