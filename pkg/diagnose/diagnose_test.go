package diagnose

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestDiagnose(t *testing.T) {
	baseline := []*model.Record{
		{Level: model.LevelInfo, Path: "/api/v1/users", Message: "user read ok"},
		{Level: model.LevelInfo, Path: "/api/v1/users", Message: "user read ok"},
		{Level: model.LevelInfo, Path: "/api/v1/orders", Message: "order created"},
	}

	incident := []*model.Record{
		{Level: model.LevelInfo, Path: "/api/v1/users", Message: "user read ok"},
		{Level: model.LevelError, Path: "/api/v1/pay", Message: "gateway 504 gateway timeout"},
		{Level: model.LevelError, Path: "/api/v1/pay", Message: "gateway 504 gateway timeout"},
	}

	report := Diagnose(baseline, incident)
	if len(report.TopCulprits) == 0 {
		t.Fatalf("expected anomalous culprits in diagnosis report")
	}

	top := report.TopCulprits[0]
	if !top.IsNewAnomaly {
		t.Errorf("504 timeout was not in baseline, should be marked IsNewAnomaly")
	}
	if top.RelativeRisk <= 1.0 {
		t.Errorf("relative risk must be high for incident culprit")
	}
}
