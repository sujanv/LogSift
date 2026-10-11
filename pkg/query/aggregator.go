package query

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/sujanv/LogSift/pkg/model"
)

type TopKAggregator struct {
	mu     sync.Mutex
	field  string
	limit  int
	counts map[string]int64
}

func NewTopKAggregator(field string, limit int) *TopKAggregator {
	if limit <= 0 {
		limit = 10
	}
	return &TopKAggregator{
		field:  field,
		limit:  limit,
		counts: make(map[string]int64),
	}
}

func (t *TopKAggregator) Ingest(r *model.Record) {
	val := extractRecordValue(r, t.field)
	if val == "" {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts[val]++
}

type TopKItem struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

func (t *TopKAggregator) Result() []TopKItem {
	t.mu.Lock()
	defer t.mu.Unlock()

	items := make([]TopKItem, 0, len(t.counts))
	for k, c := range t.counts {
		items = append(items, TopKItem{Key: k, Count: c})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	if len(items) > t.limit {
		return items[:t.limit]
	}
	return items
}

func extractRecordValue(r *model.Record, field string) string {
	switch strings.ToLower(field) {
	case "path":
		return r.Path
	case "method":
		return r.Method
	case "ip":
		return r.IP
	case "level":
		return r.Level.String()
	default:
		if r.Fields != nil {
			if v, ok := r.Fields[field]; ok {
				return fmt.Sprintf("%v", v)
			}
		}
		return ""
	}
}
