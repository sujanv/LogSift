package stats

import (
	"math"
	"testing"
)

func TestTDigest(t *testing.T) {
	td := NewTDigest(100)
	for i := 1; i <= 1000; i++ {
		td.Add(float64(i))
	}

	p50 := td.Quantile(0.50)
	if math.Abs(p50-500.0) > 30.0 {
		t.Errorf("expected p50 around 500, got %f", p50)
	}

	p99 := td.Quantile(0.99)
	if math.Abs(p99-990.0) > 30.0 {
		t.Errorf("expected p99 around 990, got %f", p99)
	}
}
