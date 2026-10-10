package diff

import (
	"time"

	"github.com/sujanv/LogSift/pkg/stats"
)

// Snapshot represents a summary metrics capture of a log session.
type Snapshot struct {
	TotalLogs  int64                `json:"total_logs"`
	ErrorCount int64                `json:"error_count"`
	ErrorRate  float64              `json:"error_rate"`
	Latency    stats.LatencySummary `json:"latency"`
	Errors     []*stats.ErrorCluster `json:"errors"`
}

// DiffReport describes statistical differences and regressions between two runs.
type DiffReport struct {
	TotalDelta    int64         `json:"total_delta"`
	ErrorRateDiff float64       `json:"error_rate_diff"` // percentage points
	P50Delta      time.Duration `json:"p50_delta"`
	P95Delta      time.Duration `json:"p95_delta"`
	P99Delta      time.Duration `json:"p99_delta"`
	NewErrors     []string      `json:"new_errors"`
	ResolvedErrors []string     `json:"resolved_errors"`
}

// Compare computes regression metrics between base and target snapshots.
func Compare(base, target *Snapshot) *DiffReport {
	baseRate := 0.0
	if base.TotalLogs > 0 {
		baseRate = float64(base.ErrorCount) / float64(base.TotalLogs)
	}
	targetRate := 0.0
	if target.TotalLogs > 0 {
		targetRate = float64(target.ErrorCount) / float64(target.TotalLogs)
	}

	rep := &DiffReport{
		TotalDelta:    target.TotalLogs - base.TotalLogs,
		ErrorRateDiff: (targetRate - baseRate) * 100.0,
		P50Delta:      target.Latency.P50 - base.Latency.P50,
		P95Delta:      target.Latency.P95 - base.Latency.P95,
		P99Delta:      target.Latency.P99 - base.Latency.P99,
	}

	baseErrMap := make(map[string]bool)
	for _, e := range base.Errors {
		baseErrMap[e.Fingerprint] = true
	}

	targetErrMap := make(map[string]bool)
	for _, e := range target.Errors {
		targetErrMap[e.Fingerprint] = true
		if !baseErrMap[e.Fingerprint] {
			rep.NewErrors = append(rep.NewErrors, e.Fingerprint)
		}
	}

	for _, e := range base.Errors {
		if !targetErrMap[e.Fingerprint] {
			rep.ResolvedErrors = append(rep.ResolvedErrors, e.Fingerprint)
		}
	}

	return rep
}
