package stats

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// TimeBucket holds record and error counts for a specific interval.
type TimeBucket struct {
	Timestamp time.Time `json:"timestamp"`
	Total     int64     `json:"total"`
	Errors    int64     `json:"errors"`
}

// TimeHistogram aggregates logs into discrete time intervals with ASCII rendering.
type TimeHistogram struct {
	mu       sync.Mutex
	interval time.Duration
	buckets  map[int64]*TimeBucket
}

// NewTimeHistogram creates a TimeHistogram with specified bucket interval.
func NewTimeHistogram(interval time.Duration) *TimeHistogram {
	if interval <= 0 {
		interval = time.Minute
	}
	return &TimeHistogram{
		interval: interval,
		buckets:  make(map[int64]*TimeBucket),
	}
}

// Add records an entry into the corresponding time bucket.
func (th *TimeHistogram) Add(ts time.Time, isError bool) {
	if ts.IsZero() {
		return
	}

	th.mu.Lock()
	defer th.mu.Unlock()

	bucketKey := ts.Truncate(th.interval).Unix()
	b, exists := th.buckets[bucketKey]
	if !exists {
		b = &TimeBucket{
			Timestamp: time.Unix(bucketKey, 0),
		}
		th.buckets[bucketKey] = b
	}

	b.Total++
	if isError {
		b.Errors++
	}
}

// Buckets returns sorted time buckets.
func (th *TimeHistogram) Buckets() []*TimeBucket {
	th.mu.Lock()
	defer th.mu.Unlock()

	list := make([]*TimeBucket, 0, len(th.buckets))
	for _, b := range th.buckets {
		list = append(list, b)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Timestamp.Before(list[j].Timestamp)
	})

	return list
}

// RenderASCII produces a Unicode sparkline/bar graph of the distribution.
func (th *TimeHistogram) RenderASCII(maxWidth int) string {
	buckets := th.Buckets()
	if len(buckets) == 0 {
		return "No time data available"
	}

	var maxVal int64 = 1
	for _, b := range buckets {
		if b.Total > maxVal {
			maxVal = b.Total
		}
	}

	bars := []rune{' ', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	var sb strings.Builder

	for _, b := range buckets {
		ratio := float64(b.Total) / float64(maxVal)
		barIdx := int(ratio * float64(len(bars)-1))
		if barIdx >= len(bars) {
			barIdx = len(bars) - 1
		}
		sb.WriteString(fmt.Sprintf("%s | %6d reqs %s\n", b.Timestamp.Format("15:04:05"), b.Total, string(bars[barIdx])))
	}

	return sb.String()
}
