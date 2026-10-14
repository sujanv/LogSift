package anomaly

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestEWMAAnomaly(t *testing.T) {
	ewma := NewEWMA(0.2)
	for i := 0; i < 20; i++ {
		ewma.Add(float64(10 + (i % 3)))
	}
	z := ewma.ZScore(50.0)
	if z <= 2.0 {
		t.Errorf("expected high z-score for sudden spike, got %f", z)
	}
}

func TestStreamingDetector(t *testing.T) {
	detector := NewStreamingAnomalyDetector(2.5)

	baselineErrors := []int{2, 5, 1, 4, 3, 2, 4}
	for _, numErrs := range baselineErrors {
		for i := 0; i < 50; i++ {
			rec := &model.Record{Level: model.LevelInfo, Timestamp: time.Now()}
			if i < numErrs {
				rec.Level = model.LevelError
			}
			detector.Ingest(rec)
		}
	}

	var anom *AnomalyScore
	for i := 0; i < 50; i++ {
		rec := &model.Record{Level: model.LevelInfo, Timestamp: time.Now()}
		if i < 45 {
			rec.Level = model.LevelError
		}
		res := detector.Ingest(rec)
		if res != nil {
			anom = res
		}
	}

	if anom == nil {
		t.Error("expected anomaly detection on high error batch")
	}
}

func TestIQROutlier(t *testing.T) {
	detector := NewIQROutlierDetector()
	for i := 1; i <= 100; i++ {
		detector.Add(time.Duration(i) * time.Millisecond)
	}
	if !detector.IsOutlier(500 * time.Millisecond) {
		t.Errorf("500ms should be an outlier")
	}
	if detector.IsOutlier(50 * time.Millisecond) {
		t.Errorf("50ms should not be an outlier")
	}
}
