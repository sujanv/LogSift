package stats

import (
	"sort"
	"sync"
)

type Centroid struct {
	Mean   float64
	Weight float64
}

type TDigest struct {
	mu          sync.Mutex
	centroids   []Centroid
	maxClusters int
	totalWeight float64
}

func NewTDigest(maxClusters int) *TDigest {
	if maxClusters <= 0 {
		maxClusters = 100
	}
	return &TDigest{
		centroids:   make([]Centroid, 0, maxClusters),
		maxClusters: maxClusters,
	}
}

func (t *TDigest) Add(value float64) {
	t.AddWeighted(value, 1.0)
}

func (t *TDigest) AddWeighted(value, weight float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.totalWeight += weight
	t.centroids = append(t.centroids, Centroid{Mean: value, Weight: weight})
	if len(t.centroids) > t.maxClusters*2 {
		t.compress()
	}
}

func (t *TDigest) compress() {
	sort.Slice(t.centroids, func(i, j int) bool {
		return t.centroids[i].Mean < t.centroids[i].Mean
	})

	merged := make([]Centroid, 0, t.maxClusters)
	cur := t.centroids[0]

	for i := 1; i < len(t.centroids); i++ {
		c := t.centroids[i]
		if len(merged) < t.maxClusters && cur.Weight+c.Weight <= (t.totalWeight/float64(t.maxClusters))*2 {
			cur.Mean = (cur.Mean*cur.Weight + c.Mean*c.Weight) / (cur.Weight + c.Weight)
			cur.Weight += c.Weight
		} else {
			merged = append(merged, cur)
			cur = c
		}
	}
	merged = append(merged, cur)
	t.centroids = merged
}

func (t *TDigest) Quantile(q float64) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.centroids) == 0 {
		return 0
	}
	if q <= 0 {
		return t.centroids[0].Mean
	}
	if q >= 1.0 {
		return t.centroids[len(t.centroids)-1].Mean
	}

	targetWeight := q * t.totalWeight
	var curWeight float64

	for i, c := range t.centroids {
		curWeight += c.Weight
		if curWeight >= targetWeight {
			return c.Mean
		}
		if i == len(t.centroids)-1 {
			return c.Mean
		}
	}
	return t.centroids[len(t.centroids)-1].Mean
}
