package dynlevel

import (
	"sync"

	"github.com/sujanv/LogSift/pkg/model"
)

type FlightRecorder struct {
	mu       sync.Mutex
	capacity int
	ring     []*model.Record
	head     int
	size     int
}

func NewFlightRecorder(capacity int) *FlightRecorder {
	if capacity <= 0 {
		capacity = 64
	}
	return &FlightRecorder{
		capacity: capacity,
		ring:     make([]*model.Record, capacity),
	}
}

func (fr *FlightRecorder) Record(rec *model.Record) {
	if rec == nil {
		return
	}
	fr.mu.Lock()
	defer fr.mu.Unlock()

	fr.ring[fr.head] = rec
	fr.head = (fr.head + 1) % fr.capacity
	if fr.size < fr.capacity {
		fr.size++
	}
}

func (fr *FlightRecorder) Size() int {
	fr.mu.Lock()
	defer fr.mu.Unlock()
	return fr.size
}

func (fr *FlightRecorder) Clear() {
	fr.mu.Lock()
	defer fr.mu.Unlock()
	fr.size = 0
	fr.head = 0
}

func (fr *FlightRecorder) FlushOnFault(faultRec *model.Record) []*model.Record {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	out := make([]*model.Record, 0, fr.size+1)
	start := (fr.head - fr.size + fr.capacity) % fr.capacity
	for i := 0; i < fr.size; i++ {
		idx := (start + i) % fr.capacity
		if fr.ring[idx] != nil {
			out = append(out, fr.ring[idx])
		}
	}

	if faultRec != nil {
		out = append(out, faultRec)
	}

	fr.size = 0
	fr.head = 0
	return out
}
