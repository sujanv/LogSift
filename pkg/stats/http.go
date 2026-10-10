package stats

import (
	"sort"
	"sync"
	"time"
)

// EndpointStats aggregates metrics per API endpoint.
type EndpointStats struct {
	Path       string          `json:"path"`
	Count      int64           `json:"count"`
	Errors     int64           `json:"errors"`
	Latency    *LatencyTracker `json:"-"`
	AvgLatency time.Duration   `json:"avg_latency"`
	P95Latency time.Duration   `json:"p95_latency"`
}

// HTTPAnalytics tracks HTTP status codes and endpoint performance.
type HTTPAnalytics struct {
	mu          sync.Mutex
	statusCodes map[int]int64
	endpoints   map[string]*EndpointStats
	totalReqs   int64
}

// NewHTTPAnalytics creates an HTTPAnalytics aggregator.
func NewHTTPAnalytics() *HTTPAnalytics {
	return &HTTPAnalytics{
		statusCodes: make(map[int]int64),
		endpoints:   make(map[string]*EndpointStats),
	}
}

// Add records an HTTP event.
func (h *HTTPAnalytics) Add(statusCode int, path string, latency time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.totalReqs++

	if statusCode > 0 {
		h.statusCodes[statusCode]++
	}

	if len(path) > 0 {
		ep, exists := h.endpoints[path]
		if !exists {
			ep = &EndpointStats{
				Path:    path,
				Latency: NewLatencyTracker(),
			}
			h.endpoints[path] = ep
		}
		ep.Count++
		if statusCode >= 400 {
			ep.Errors++
		}
		if latency > 0 {
			ep.Latency.Add(latency)
		}
	}
}

// StatusBreakdown returns counts of 2xx, 3xx, 4xx, 5xx.
func (h *HTTPAnalytics) StatusBreakdown() map[string]int64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	res := map[string]int64{"2xx": 0, "3xx": 0, "4xx": 0, "5xx": 0, "other": 0}
	for code, count := range h.statusCodes {
		switch {
		case code >= 200 && code < 300:
			res["2xx"] += count
		case code >= 300 && code < 400:
			res["3xx"] += count
		case code >= 400 && code < 500:
			res["4xx"] += count
		case code >= 500 && code < 600:
			res["5xx"] += count
		default:
			res["other"] += count
		}
	}
	return res
}

// TopSlowestEndpoints returns top N endpoints by P95 latency.
func (h *HTTPAnalytics) TopSlowestEndpoints(n int) []*EndpointStats {
	h.mu.Lock()
	defer h.mu.Unlock()

	list := make([]*EndpointStats, 0, len(h.endpoints))
	for _, ep := range h.endpoints {
		snap := ep.Latency.Snapshot()
		ep.AvgLatency = snap.Mean
		ep.P95Latency = snap.P95
		list = append(list, ep)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].P95Latency > list[j].P95Latency
	})

	if n > 0 && len(list) > n {
		return list[:n]
	}
	return list
}
