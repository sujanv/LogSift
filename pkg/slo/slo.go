package slo

import (
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type SLIConfig struct {
	TargetAvailability float64
	LatencyThreshold   time.Duration
	TargetLatencyRatio float64
}

type SLOMetrics struct {
	TotalRequests       int64   `json:"total_requests"`
	SuccessfulRequests  int64   `json:"successful_requests"`
	FastRequests        int64   `json:"fast_requests"`
	AvailabilityRatio   float64 `json:"availability_ratio"`
	LatencyRatio        float64 `json:"latency_ratio"`
	AvailabilityBudget  float64 `json:"availability_budget"`
	BurnRate            float64 `json:"burn_rate"`
	IsBreached          bool    `json:"is_breached"`
}

type SLOCalculator struct {
	cfg SLIConfig
}

func NewSLOCalculator(cfg SLIConfig) *SLOCalculator {
	if cfg.TargetAvailability <= 0 {
		cfg.TargetAvailability = 0.999
	}
	if cfg.TargetLatencyRatio <= 0 {
		cfg.TargetLatencyRatio = 0.95
	}
	if cfg.LatencyThreshold <= 0 {
		cfg.LatencyThreshold = 250 * time.Millisecond
	}
	return &SLOCalculator{cfg: cfg}
}

func (c *SLOCalculator) Evaluate(records []*model.Record) SLOMetrics {
	if len(records) == 0 {
		return SLOMetrics{}
	}

	var total, success, fast int64
	for _, r := range records {
		total++
		if r.StatusCode < 500 && r.Level < model.LevelError {
			success++
		}
		if r.Latency > 0 && r.Latency <= c.cfg.LatencyThreshold {
			fast++
		}
	}

	avail := float64(success) / float64(total)
	latRatio := float64(fast) / float64(total)

	allowedErrorRatio := 1.0 - c.cfg.TargetAvailability
	actualErrorRatio := 1.0 - avail

	burnRate := 0.0
	if allowedErrorRatio > 0 {
		burnRate = actualErrorRatio / allowedErrorRatio
	}

	remainingBudget := 1.0 - (actualErrorRatio / allowedErrorRatio)
	if remainingBudget < 0 {
		remainingBudget = 0
	}

	return SLOMetrics{
		TotalRequests:      total,
		SuccessfulRequests: success,
		FastRequests:       fast,
		AvailabilityRatio:  avail,
		LatencyRatio:       latRatio,
		AvailabilityBudget: remainingBudget,
		BurnRate:           burnRate,
		IsBreached:         avail < c.cfg.TargetAvailability || latRatio < c.cfg.TargetLatencyRatio,
	}
}
