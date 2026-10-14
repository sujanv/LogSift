package cluster

import (
	"regexp"
	"strings"
	"sync"
)

var (
	digitRegex = regexp.MustCompile(`\b\d+\b`)
	hexRegex   = regexp.MustCompile(`\b0x[0-9a-fA-F]+\b`)
	uuidRegex  = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
)

type LogCluster struct {
	ID        int      `json:"id"`
	Template  string   `json:"template"`
	Tokens    []string `json:"tokens"`
	Size      int64    `json:"size"`
	Variables []string `json:"variables,omitempty"`
}

type DrainClusterer struct {
	mu           sync.Mutex
	simThreshold float64
	clusters     []*LogCluster
}

func NewDrainClusterer(simThreshold float64) *DrainClusterer {
	if simThreshold <= 0 || simThreshold > 1.0 {
		simThreshold = 0.5
	}
	return &DrainClusterer{
		simThreshold: simThreshold,
		clusters:     make([]*LogCluster, 0, 64),
	}
}

func (dc *DrainClusterer) Add(msg string) *LogCluster {
	tokens := tokenizeMessage(msg)
	if len(tokens) == 0 {
		return nil
	}

	dc.mu.Lock()
	defer dc.mu.Unlock()

	var bestMatch *LogCluster
	var maxSim float64

	for _, c := range dc.clusters {
		if len(c.Tokens) != len(tokens) {
			continue
		}
		sim := computeSimilarity(c.Tokens, tokens)
		if sim > maxSim && sim >= dc.simThreshold {
			maxSim = sim
			bestMatch = c
		}
	}

	if bestMatch != nil {
		bestMatch.Size++
		for i := 0; i < len(tokens); i++ {
			if bestMatch.Tokens[i] != tokens[i] {
				bestMatch.Tokens[i] = "<*>"
			}
		}
		bestMatch.Template = strings.Join(bestMatch.Tokens, " ")
		return bestMatch
	}

	newID := len(dc.clusters) + 1
	clusterTokens := make([]string, len(tokens))
	copy(clusterTokens, tokens)
	cluster := &LogCluster{
		ID:       newID,
		Template: strings.Join(clusterTokens, " "),
		Tokens:   clusterTokens,
		Size:     1,
	}
	dc.clusters = append(dc.clusters, cluster)
	return cluster
}

func (dc *DrainClusterer) Clusters() []*LogCluster {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	res := make([]*LogCluster, len(dc.clusters))
	copy(res, dc.clusters)
	return res
}

func tokenizeMessage(msg string) []string {
	msg = uuidRegex.ReplaceAllString(msg, "<*UUID*>")
	msg = hexRegex.ReplaceAllString(msg, "<*HEX*>")
	msg = digitRegex.ReplaceAllString(msg, "<*NUM*>")
	return strings.Fields(msg)
}

func computeSimilarity(t1, t2 []string) float64 {
	if len(t1) == 0 || len(t2) == 0 || len(t1) != len(t2) {
		return 0
	}
	matches := 0
	for i := 0; i < len(t1); i++ {
		if t1[i] == t2[i] {
			matches++
		}
	}
	return float64(matches) / float64(len(t1))
}
