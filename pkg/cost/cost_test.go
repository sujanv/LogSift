package cost

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestCostEstimator(t *testing.T) {
	ce := NewCostEstimator()

	for i := 0; i < 1000; i++ {
		ce.IngestService(&model.Record{
			Level:   model.LevelInfo,
			Message: "user checkout success",
			Fields:  map[string]interface{}{"service": "api-gateway"},
		}, 1024)
	}

	rep := ce.Report(ProviderAWSCloudWatch)
	if rep.TotalEvents != 1000 {
		t.Fatalf("expected 1000 events, got %d", rep.TotalEvents)
	}
	if rep.EstimatedCostUSD <= 0 {
		t.Errorf("cost must be greater than 0")
	}
	if rep.ByService["api-gateway"] <= 0 {
		t.Errorf("expected cost attribution for api-gateway")
	}
}
