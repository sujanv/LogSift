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

func (ce *CostEstimator) IngestService(rec *model.Record, sizeBytes int64) {
	ce.Ingest(rec, sizeBytes)
	if rec == nil {
		return
	}
	ce.mu.Lock()
	defer ce.mu.Unlock()

	svc := "default"
	if rec.Fields != nil {
		if s, ok := rec.Fields["service"].(string); ok && s != "" {
			svc = s
		} else if s, ok := rec.Fields["app"].(string); ok && s != "" {
			svc = s
		}
	}
	ce.bytesBySvc[svc] += sizeBytes
}

func (ce *CostEstimator) Report(provider Provider) CostBreakdown {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	gb := float64(ce.totalBytes) / (1024 * 1024 * 1024)
	var ratePerGB float64

	switch provider {
	case ProviderDatadog:
		ratePerGB = 0.65
	case ProviderSplunk:
		ratePerGB = 1.50
	case ProviderGCPLogging:
		ratePerGB = 0.50
	case ProviderAWSCloudWatch:
		fallthrough
	default:
		provider = ProviderAWSCloudWatch
		ratePerGB = 0.50
	}

	costUSD := gb * ratePerGB

	byLvl := make(map[string]float64)
	for lvl, b := range ce.bytesByLvl {
		byLvl[lvl.String()] = (float64(b) / (1024 * 1024 * 1024)) * ratePerGB
	}

	bySvc := make(map[string]float64)
	for s, b := range ce.bytesBySvc {
		bySvc[s] = (float64(b) / (1024 * 1024 * 1024)) * ratePerGB
	}

	var tips []string
	debugCost := byLvl[model.LevelDebug.String()]
	if debugCost > 0.05*costUSD && debugCost > 1.0 {
		tips = append(tips, "Dropping DEBUG logs in production could save ~30-50% on ingestion billing")
	}

	return CostBreakdown{
		Provider:         provider,
		TotalBytes:       ce.totalBytes,
		TotalGigabytes:   gb,
		TotalEvents:      ce.totalEvents,
		EstimatedCostUSD: costUSD,
		ByLevel:          byLvl,
		ByService:        bySvc,
		OptimizationTips: tips,
	}
}
