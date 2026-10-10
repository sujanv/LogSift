package alert

import (
	"sync"
	"time"
)

// AlertRule defines threshold criteria for triggering alerts.
type AlertRule struct {
	Name               string        `json:"name"`
	ErrorRateThreshold float64       `json:"error_rate_threshold"` // e.g. 0.05 for 5%
	P99Threshold       time.Duration `json:"p99_threshold"`
	Window             time.Duration `json:"window"`
}

// AlertEvent describes an active triggered alert.
type AlertEvent struct {
	RuleName  string    `json:"rule_name"`
	Reason    string    `json:"reason"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// AlertEvaluator tracks metrics across sliding windows to detect anomalies.
type AlertEvaluator struct {
	mu           sync.Mutex
	rules        []AlertRule
	totalEvents  int64
	errorEvents  int64
	windowStart  time.Time
	activeAlerts []AlertEvent
}

// NewAlertEvaluator creates an AlertEvaluator.
func NewAlertEvaluator(rules []AlertRule) *AlertEvaluator {
	return &AlertEvaluator{
		rules:       rules,
		windowStart: time.Now(),
	}
}

// Record observes an event and evaluates alert rules.
func (ae *AlertEvaluator) Record(isError bool, lat time.Duration, ts time.Time) []AlertEvent {
	ae.mu.Lock()
	defer ae.mu.Unlock()

	ae.totalEvents++
	if isError {
		ae.errorEvents++
	}

	if ts.Sub(ae.windowStart) >= time.Minute && ae.totalEvents > 20 {
		rate := float64(ae.errorEvents) / float64(ae.totalEvents)
		var triggered []AlertEvent

		for _, r := range ae.rules {
			if r.ErrorRateThreshold > 0 && rate > r.ErrorRateThreshold {
				ev := AlertEvent{
					RuleName:  r.Name,
					Reason:    "Error rate exceeded threshold",
					Value:     rate,
					Timestamp: ts,
				}
				triggered = append(triggered, ev)
				ae.activeAlerts = append(ae.activeAlerts, ev)
			}
		}

		// Reset window
		ae.totalEvents = 0
		ae.errorEvents = 0
		ae.windowStart = ts
		return triggered
	}

	return nil
}

// ActiveAlerts returns all alerts triggered.
func (ae *AlertEvaluator) ActiveAlerts() []AlertEvent {
	ae.mu.Lock()
	defer ae.mu.Unlock()
	res := make([]AlertEvent, len(ae.activeAlerts))
	copy(res, ae.activeAlerts)
	return res
}
