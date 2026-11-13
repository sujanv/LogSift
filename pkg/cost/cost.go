package cost

import (
	"sync"

	"github.com/sujanv/LogSift/pkg/model"
)

type Provider string

const (
	ProviderAWSCloudWatch Provider = "aws"
	ProviderDatadog       Provider = "datadog"
	ProviderSplunk        Provider = "splunk"
	ProviderGCPLogging    Provider = "gcp"
)

type CostBreakdown struct {
	Provider          Provider           `json:"provider"`
	TotalBytes        int64              `json:"total_bytes"`
	TotalGigabytes    float64            `json:"total_gigabytes"`
	TotalEvents       int64              `json:"total_events"`
	EstimatedCostUSD  float64            `json:"estimated_cost_usd"`
	ByLevel           map[string]float64 `json:"by_level"`
	ByService         map[string]float64 `json:"by_service"`
	OptimizationTips  []string           `json:"optimization_tips"`
}

type CostEstimator struct {
	mu          sync.Mutex
	totalBytes  int64
	totalEvents int64
	bytesByLvl  map[model.LogLevel]int64
	bytesBySvc  map[string]int64
}

func NewCostEstimator() *CostEstimator {
	return &CostEstimator{
		bytesByLvl: make(map[model.LogLevel]int64),
		bytesBySvc: make(map[string]int64),
	}
}

func (ce *CostEstimator) Ingest(rec *model.Record, sizeBytes int64) {
	if rec == nil {
		return
	}
	if sizeBytes <= 0 {
		sizeBytes = int64(len(rec.Raw))
		if sizeBytes == 0 {
			sizeBytes = int64(len(rec.Message)) + 64
		}
	}

	ce.mu.Lock()
	defer ce.mu.Unlock()

	ce.totalBytes += sizeBytes
	ce.totalEvents++
	ce.bytesByLvl[rec.Level] += sizeBytes
}
