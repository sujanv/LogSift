package sampling

import (
	"math/rand"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

// TailSampler samples log entries while guaranteeing 100% retention of errors and slow requests.
type TailSampler struct {
	rate          float64
	slowThreshold time.Duration
	alwaysErrors  bool
	rnd           *rand.Rand
}

// NewTailSampler creates a TailSampler.
func NewTailSampler(rate float64, slowThreshold time.Duration, alwaysErrors bool) *TailSampler {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	return &TailSampler{
		rate:          rate,
		slowThreshold: slowThreshold,
		alwaysErrors:  alwaysErrors,
		rnd:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ShouldSample returns true if the record should be retained.
func (s *TailSampler) ShouldSample(r *model.Record) bool {
	// Always retain errors
	if s.alwaysErrors && (r.Level >= model.LevelError || r.StatusCode >= 500) {
		return true
	}

	// Always retain slow requests
	if s.slowThreshold > 0 && r.Latency >= s.slowThreshold {
		return true
	}

	// Probabilistic sample
	if s.rate >= 1.0 {
		return true
	}
	if s.rate <= 0.0 {
		return false
	}

	return s.rnd.Float64() < s.rate
}
