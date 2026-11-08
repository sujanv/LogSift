package throttle

import (
	"math/rand"
	"sync"

	"github.com/sujanv/LogSift/pkg/model"
)

type ErrorPreservingSampler struct {
	mu         sync.Mutex
	sampleRate float64
	rng        *rand.Rand
}

func NewErrorPreservingSampler(sampleRate float64) *ErrorPreservingSampler {
	if sampleRate <= 0 || sampleRate > 1.0 {
		sampleRate = 0.10
	}
	return &ErrorPreservingSampler{
		sampleRate: sampleRate,
		rng:        rand.New(rand.NewSource(42)),
	}
}

func (s *ErrorPreservingSampler) SetSampleRate(rate float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if rate > 0 && rate <= 1.0 {
		s.sampleRate = rate
	}
}

func (s *ErrorPreservingSampler) ShouldSample(rec *model.Record) bool {
	if rec == nil {
		return false
	}
	if rec.Level >= model.LevelError || rec.StatusCode >= 500 {
		return true
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.rng.Float64() < s.sampleRate
}
