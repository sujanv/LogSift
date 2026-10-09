package stats

import (
	"math"
	"sort"
	"sync"
	"time"
)

const (
	// MaxReservoirSamples keeps fixed bounded memory for percentiles (O(1) space).
	MaxReservoirSamples = 10000
)

// LatencyTracker computes streaming latency metrics without unbounded memory.
type LatencyTracker struct {
	mu      sync.Mutex
	count   int64
	sum     time.Duration
	min     time.Duration
	max     time.Duration
	samples []time.Duration
}

// NewLatencyTracker creates a LatencyTracker.
func NewLatencyTracker() *LatencyTracker {
	return &LatencyTracker{
		samples: make([]time.Duration, 0, MaxReservoirSamples),
	}
}

// Add records a new latency sample.
func (t *LatencyTracker) Add(d time.Duration) {
	if d <= 0 {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.count++
	t.sum += d

	if t.min == 0 || d < t.min {
		t.min = d
	}
	if d > t.max {
		t.max = d
	}

	// Reservoir sampling to maintain constant memory
	if len(t.samples) < MaxReservoirSamples {
		t.samples = append(t.samples, d)
	} else {
		// Replace existing sample probabilistically
		idx := int(t.count % int64(MaxReservoirSamples))
		t.samples[idx] = d
	}
}

// Summary represents the snapshot of latency statistics.
type LatencySummary struct {
	Count  int64         `json:"count"`
	Min    time.Duration `json:"min"`
	Max    time.Duration `json:"max"`
	Mean   time.Duration `json:"mean"`
	P50    time.Duration `json:"p50"`
	P75    time.Duration `json:"p75"`
	P90    time.Duration `json:"p90"`
	P95    time.Duration `json:"p95"`
	P99    time.Duration `json:"p99"`
	P999   time.Duration `json:"p999"`
}

// Snapshot calculates percentiles and returns LatencySummary.
func (t *LatencyTracker) Snapshot() LatencySummary {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.count == 0 || len(t.samples) == 0 {
		return LatencySummary{}
	}

	mean := time.Duration(int64(t.sum) / t.count)

	// Sort copy of samples to compute percentiles
	sorted := make([]time.Duration, len(t.samples))
	copy(sorted, t.samples)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	return LatencySummary{
		Count: t.count,
		Min:   t.min,
		Max:   t.max,
		Mean:  mean,
		P50:   percentile(sorted, 50.0),
		P75:   percentile(sorted, 75.0),
		P90:   percentile(sorted, 90.0),
		P95:   percentile(sorted, 95.0),
		P99:   percentile(sorted, 99.0),
		P999:  percentile(sorted, 99.9),
	}
}

func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(math.Ceil((p / 100.0) * float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
