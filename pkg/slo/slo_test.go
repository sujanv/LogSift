package slo

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestSLOCalculator(t *testing.T) {
	calc := NewSLOCalculator(SLIConfig{
		TargetAvailability: 0.99,
		LatencyThreshold:   100 * time.Millisecond,
		TargetLatencyRatio: 0.90,
	})

	records := make([]*model.Record, 100)
	for i := 0; i < 100; i++ {
		records[i] = &model.Record{
			StatusCode: 200,
			Latency:    50 * time.Millisecond,
		}
	}
	for i := 0; i < 5; i++ {
		records[i].StatusCode = 500
	}

	metrics := calc.Evaluate(records)
	if metrics.TotalRequests != 100 {
		t.Errorf("expected 100 requests")
	}
	if metrics.AvailabilityRatio != 0.95 {
		t.Errorf("expected 0.95 availability, got %f", metrics.AvailabilityRatio)
	}
	if !metrics.IsBreached {
		t.Errorf("SLO must be marked breached when availability 95%% < target 99%%")
	}
	if metrics.BurnRate <= 1.0 {
		t.Errorf("burn rate must exceed 1.0 during active error breach, got %f", metrics.BurnRate)
	}

	alert := EvaluateBurnRate(metrics.BurnRate)
	if alert == nil || alert.Severity != "HIGH" {
		t.Errorf("expected high burn rate alert for 5x over budget: %+v", alert)
	}

	critAlert := EvaluateBurnRate(15.0)
	if critAlert == nil || critAlert.Severity != "CRITICAL" {
		t.Errorf("expected critical burn rate alert for 15x over budget: %+v", critAlert)
	}
}
