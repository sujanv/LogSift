package dedup

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type DedupEvent struct {
	Original  *model.Record `json:"original"`
	Count     int64         `json:"count"`
	FirstSeen time.Time     `json:"first_seen"`
	LastSeen  time.Time     `json:"last_seen"`
	Hash      string        `json:"hash"`
}

type WindowedDeduplicator struct {
	mu           sync.Mutex
	window       time.Duration
	fingerprints map[string]*DedupEvent
	expired      []*DedupEvent
}

func NewWindowedDeduplicator(window time.Duration) *WindowedDeduplicator {
	if window <= 0 {
		window = 10 * time.Second
	}
	return &WindowedDeduplicator{
		window:       window,
		fingerprints: make(map[string]*DedupEvent),
		expired:      make([]*DedupEvent, 0, 64),
	}
}

func (d *WindowedDeduplicator) Ingest(r *model.Record) (isNew bool, event *DedupEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	h := sha256.New()
	h.Write([]byte(r.Level.String()))
	h.Write([]byte(":"))
	h.Write([]byte(r.Message))
	hash := hex.EncodeToString(h.Sum(nil))

	now := r.Timestamp
	if now.IsZero() {
		now = time.Now()
	}

	ev, exists := d.fingerprints[hash]
	if !exists || now.Sub(ev.FirstSeen) > d.window {
		if exists {
			d.expired = append(d.expired, ev)
		}
		newEvent := &DedupEvent{
			Original:  r,
			Count:     1,
			FirstSeen: now,
			LastSeen:  now,
			Hash:      hash,
		}
		d.fingerprints[hash] = newEvent
		return true, newEvent
	}

	ev.Count++
	ev.LastSeen = now
	return false, ev
}

func (d *WindowedDeduplicator) Flush() []*DedupEvent {
	d.mu.Lock()
	defer d.mu.Unlock()

	var result []*DedupEvent
	result = append(result, d.expired...)
	for _, ev := range d.fingerprints {
		result = append(result, ev)
	}
	d.fingerprints = make(map[string]*DedupEvent)
	d.expired = make([]*DedupEvent, 0, 64)
	return result
}
