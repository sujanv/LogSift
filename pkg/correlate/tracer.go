package correlate

import (
	"sort"
	"sync"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

// TraceEntry represents a single log record associated with a trace.
type TraceEntry struct {
	Timestamp  time.Time     `json:"timestamp"`
	Level      string        `json:"level"`
	Message    string        `json:"message"`
	Path       string        `json:"path"`
	StatusCode int           `json:"status_code"`
	Latency    time.Duration `json:"latency"`
}

// TraceCorrelator aggregates distributed logs matching trace IDs.
type TraceCorrelator struct {
	mu     sync.Mutex
	traces map[string][]*TraceEntry
}

// NewTraceCorrelator creates a TraceCorrelator.
func NewTraceCorrelator() *TraceCorrelator {
	return &TraceCorrelator{
		traces: make(map[string][]*TraceEntry),
	}
}

// Ingest adds a record into its trace bucket if TraceID is present.
func (tc *TraceCorrelator) Ingest(r *model.Record) {
	if r.TraceID == "" {
		return
	}

	tc.mu.Lock()
	defer tc.mu.Unlock()

	entry := &TraceEntry{
		Timestamp:  r.Timestamp,
		Level:      r.Level.String(),
		Message:    r.Message,
		Path:       r.Path,
		StatusCode: r.StatusCode,
		Latency:    r.Latency,
	}

	tc.traces[r.TraceID] = append(tc.traces[r.TraceID], entry)
}

// GetTrace returns chronological entries for a trace ID.
func (tc *TraceCorrelator) GetTrace(traceID string) []*TraceEntry {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	entries, exists := tc.traces[traceID]
	if !exists {
		return nil
	}

	res := make([]*TraceEntry, len(entries))
	copy(res, entries)

	sort.Slice(res, func(i, j int) bool {
		return res[i].Timestamp.Before(res[j].Timestamp)
	})

	return res
}
