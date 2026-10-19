package archive

import (
	"testing"
)

func TestChunkBloomFilter(t *testing.T) {
	bf := NewChunkBloomFilter(1024)
	bf.Add("order-12345")
	bf.Add("trace-abc")

	if !bf.Test("order-12345") {
		t.Errorf("order-12345 must match bloom filter")
	}
	if !bf.Test("trace-abc") {
		t.Errorf("trace-abc must match bloom filter")
	}
	if bf.Test("non-existent-key-999") {
		t.Errorf("unexpected false positive on small filter")
	}
}
