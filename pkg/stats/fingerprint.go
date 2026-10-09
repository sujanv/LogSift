package stats

import (
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	uuidRegex   = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	ipRegex     = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	hexRegex    = regexp.MustCompile(`\b0x[0-9a-fA-F]+\b`)
	numRegex    = regexp.MustCompile(`\b\d+\b`)
	quotedRegex = regexp.MustCompile(`"[^"]*"`)
)

// NormalizeMessage masks dynamic variables to generate a structural error template.
func NormalizeMessage(msg string) string {
	msg = uuidRegex.ReplaceAllString(msg, "<UUID>")
	msg = ipRegex.ReplaceAllString(msg, "<IP>")
	msg = hexRegex.ReplaceAllString(msg, "<HEX>")
	msg = quotedRegex.ReplaceAllString(msg, `"<STR>"`)
	msg = numRegex.ReplaceAllString(msg, "<NUM>")
	return strings.TrimSpace(msg)
}

// ErrorCluster aggregates identical error patterns.
type ErrorCluster struct {
	Fingerprint string    `json:"fingerprint"`
	Count       int64     `json:"count"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
	Sample      string    `json:"sample"`
}

// ErrorSummarizer clusters errors by fingerprint.
type ErrorSummarizer struct {
	mu       sync.Mutex
	clusters map[string]*ErrorCluster
}

// NewErrorSummarizer creates an ErrorSummarizer.
func NewErrorSummarizer() *ErrorSummarizer {
	return &ErrorSummarizer{
		clusters: make(map[string]*ErrorCluster),
	}
}

// Add logs an error message occurrence into its cluster.
func (es *ErrorSummarizer) Add(msg string, ts time.Time) {
	if len(strings.TrimSpace(msg)) == 0 {
		return
	}

	fp := NormalizeMessage(msg)

	es.mu.Lock()
	defer es.mu.Unlock()

	c, exists := es.clusters[fp]
	if !exists {
		c = &ErrorCluster{
			Fingerprint: fp,
			Count:       0,
			FirstSeen:   ts,
			Sample:      msg,
		}
		es.clusters[fp] = c
	}

	c.Count++
	c.LastSeen = ts
}

// TopClusters returns the top N error clusters ranked by frequency.
func (es *ErrorSummarizer) TopClusters(n int) []*ErrorCluster {
	es.mu.Lock()
	defer es.mu.Unlock()

	list := make([]*ErrorCluster, 0, len(es.clusters))
	for _, c := range es.clusters {
		list = append(list, c)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Count > list[j].Count
	})

	if n > 0 && len(list) > n {
		return list[:n]
	}
	return list
}
