package pseudonym

import "testing"

func BenchmarkTokenize(b *testing.B) {
	p := NewPseudonymizer("benchmark-salt")
	input := "customer-sensitive-id-981247"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = p.Tokenize(input)
	}
}
