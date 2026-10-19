package slo

import "fmt"

type BurnRateAlert struct {
	Window   string
	BurnRate float64
	Severity string
	Message  string
}

func EvaluateBurnRate(burnRate float64) *BurnRateAlert {
	if burnRate >= 14.4 {
		return &BurnRateAlert{
			Window:   "1h",
			BurnRate: burnRate,
			Severity: "CRITICAL",
			Message:  fmt.Sprintf("Critical error budget burn rate (%.1fx): 2%% budget consumed in 1 hour", burnRate),
		}
	}
	if burnRate >= 3.0 {
		return &BurnRateAlert{
			Window:   "6h",
			BurnRate: burnRate,
			Severity: "HIGH",
			Message:  fmt.Sprintf("High error budget burn rate (%.1fx): elevated budget consumption", burnRate),
		}
	}
	return nil
}
