package alert

import (
	"testing"
	"time"
)

func TestAlertEvaluator(t *testing.T) {
	rule := AlertRule{
		Name:               "HighErrorRate",
		ErrorRateThreshold: 0.10, // 10%
		Window:             time.Minute,
	}
	evaluator := NewAlertEvaluator([]AlertRule{rule})

	t0 := time.Now()
	// Record 25 events, 5 errors (20% error rate)
	for i := 0; i < 20; i++ {
		evaluator.Record(false, 10*time.Millisecond, t0)
	}
	for i := 0; i < 5; i++ {
		evaluator.Record(true, 50*time.Millisecond, t0)
	}

	// Trigger window boundary
	triggered := evaluator.Record(false, 10*time.Millisecond, t0.Add(65*time.Second))
	if len(triggered) != 1 {
		t.Fatalf("expected 1 triggered alert, got %d", len(triggered))
	}
	if triggered[0].RuleName != "HighErrorRate" {
		t.Errorf("expected HighErrorRate alert, got %s", triggered[0].RuleName)
	}
}
