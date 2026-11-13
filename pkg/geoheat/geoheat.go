package geoheat

import (
	"sort"
	"sync"
	"time"
)

type RegionMetric struct {
	Country    string        `json:"country"`
	TotalHits  int64         `json:"total_hits"`
	AvgLatency time.Duration `json:"avg_latency"`
	P95Latency time.Duration `json:"p95_latency"`
	Degraded   bool          `json:"degraded"`
}

type RegionalHeatmap struct {
	mu           sync.Mutex
	degradeLimit time.Duration
	latencies    map[string][]time.Duration
}

func NewRegionalHeatmap(degradeThreshold time.Duration) *RegionalHeatmap {
	if degradeThreshold <= 0 {
		degradeThreshold = 350 * time.Millisecond
	}
	return &RegionalHeatmap{
		degradeLimit: degradeThreshold,
		latencies:    make(map[string][]time.Duration),
	}
}

func (rh *RegionalHeatmap) Record(country string, lat time.Duration) {
	if country == "" {
		country = "UNKNOWN"
	}
	rh.mu.Lock()
	defer rh.mu.Unlock()

	rh.latencies[country] = append(rh.latencies[country], lat)
}

func (rh *RegionalHeatmap) Summary() []RegionMetric {
	rh.mu.Lock()
	defer rh.mu.Unlock()

	var metrics []RegionMetric
	for country, lats := range rh.latencies {
		if len(lats) == 0 {
			continue
		}
		sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })

		var sum time.Duration
		for _, l := range lats {
			sum += l
		}
		avg := sum / time.Duration(len(lats))
		p95Idx := int(float64(len(lats)-1) * 0.95)
		p95 := lats[p95Idx]

		metrics = append(metrics, RegionMetric{
			Country:    country,
			TotalHits:  int64(len(lats)),
			AvgLatency: avg,
			P95Latency: p95,
			Degraded:   p95 >= rh.degradeLimit,
		})
	}

	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].TotalHits > metrics[j].TotalHits
	})

	return metrics
}
