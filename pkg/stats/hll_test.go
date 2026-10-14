package stats

import (
	"fmt"
	"math"
	"testing"
)

func TestHyperLogLogAccuracy(t *testing.T) {
	hll := NewHyperLogLog()
	n := 5000
	for i := 0; i < n; i++ {
		hll.AddString(fmt.Sprintf("user-%d-ip-192.168.1.%d", i, i%256))
	}

	est := hll.Count()
	errRate := math.Abs(float64(est)-float64(n)) / float64(n)
	if errRate > 0.05 {
		t.Errorf("HLL error rate too high: estimated %d, actual %d, error %f", est, n, errRate)
	}
}
