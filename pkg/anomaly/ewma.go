package anomaly

import (
	"math"
	"sync"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type AnomalyScore struct {
	Timestamp   time.Time `json:"timestamp"`
	Metric      string    `json:"metric"`
	Current     float64   `json:"current"`
	Expected    float64   `json:"expected"`
	ZScore      float64   `json:"z_score"`
	IsAnomaly   bool      `json:"is_anomaly"`
	Description string    `json:"description"`
}

type EWMA struct {
	alpha    float64
	mean     float64
	variance float64
	count    int64
}

func NewEWMA(alpha float64) *EWMA {
	if alpha <= 0 || alpha >= 1 {
		alpha = 0.2
	}
	return &EWMA{alpha: alpha}
}

func (e *EWMA) Add(value float64) float64 {
	e.count++
	if e.count == 1 {
		e.mean = value
		e.variance = 0
		return 0
	}
	diff := value - e.mean
	incr := e.alpha * diff
	e.mean += incr
	e.variance = (1 - e.alpha) * (e.variance + diff*incr)
	return e.ZScore(value)
}

func (e *EWMA) Mean() float64 {
	return e.mean
}

func (e *EWMA) StdDev() float64 {
	if e.variance <= 0 {
		return 0
	}
	return math.Sqrt(e.variance)
}

func (e *EWMA) ZScore(value float64) float64 {
	sd := e.StdDev()
	if sd < 1e-6 {
		return 0
	}
	return (value - e.mean) / sd
}

type StreamingAnomalyDetector struct {
	mu           sync.Mutex
	errorRate    *EWMA
	threshold    float64
	windowLogs   int64
	windowErrors int64
	anomalies    []AnomalyScore
}

func NewStreamingAnomalyDetector(threshold float64) *StreamingAnomalyDetector {
	if threshold <= 0 {
		threshold = 3.0
	}
	return &StreamingAnomalyDetector{
		errorRate: NewEWMA(0.2),
		threshold: threshold,
		anomalies: make([]AnomalyScore, 0, 32),
	}
}

func (d *StreamingAnomalyDetector) Ingest(r *model.Record) *AnomalyScore {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.windowLogs++
	if r.Level >= model.LevelError {
		d.windowErrors++
	}

	if d.windowLogs >= 50 {
		rate := float64(d.windowErrors) / float64(d.windowLogs)
		expected := d.errorRate.Mean()
		z := d.errorRate.ZScore(rate)
		d.errorRate.Add(rate)
		isAnom := z >= d.threshold && d.errorRate.count > 3

		d.windowLogs = 0
		d.windowErrors = 0

		if isAnom {
			score := AnomalyScore{
				Timestamp:   r.Timestamp,
				Metric:      "error_rate",
				Current:     rate,
				Expected:    expected,
				ZScore:      z,
				IsAnomaly:   true,
				Description: "Sudden spike in error rate",
			}
			d.anomalies = append(d.anomalies, score)
			return &score
		}
	}
	return nil
}

func (d *StreamingAnomalyDetector) Anomalies() []AnomalyScore {
	d.mu.Lock()
	defer d.mu.Unlock()
	res := make([]AnomalyScore, len(d.anomalies))
	copy(res, d.anomalies)
	return res
}
