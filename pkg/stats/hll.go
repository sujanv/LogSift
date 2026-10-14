package stats

import (
	"hash/fnv"
	"math"
	"math/bits"
	"sync"
)

const (
	hllPrecision = 10
	hllRegisters = 1 << hllPrecision
	hllAlphaM    = 0.7213 / (1.0 + 1.079/float64(hllRegisters))
)

type HyperLogLog struct {
	mu        sync.Mutex
	registers [hllRegisters]uint8
}

func NewHyperLogLog() *HyperLogLog {
	return &HyperLogLog{}
}

func (h *HyperLogLog) AddString(s string) {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte(s))
	hash := hasher.Sum64()

	idx := hash >> (64 - hllPrecision)
	w := hash << hllPrecision
	leadingZeros := uint8(bits.LeadingZeros64(w)) + 1

	h.mu.Lock()
	defer h.mu.Unlock()
	if leadingZeros > h.registers[idx] {
		h.registers[idx] = leadingZeros
	}
}

func (h *HyperLogLog) Count() uint64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	var sum float64
	var zeros int
	for _, reg := range h.registers {
		sum += 1.0 / math.Pow(2.0, float64(reg))
		if reg == 0 {
			zeros++
		}
	}

	estimate := (hllAlphaM * float64(hllRegisters*hllRegisters)) / sum

	if estimate <= 2.5*float64(hllRegisters) && zeros > 0 {
		estimate = float64(hllRegisters) * math.Log(float64(hllRegisters)/float64(zeros))
	}
	return uint64(math.Round(estimate))
}
