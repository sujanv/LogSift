package anomaly

import (
	"sort"
	"time"
)

type IQROutlierDetector struct {
	samples []time.Duration
}

func NewIQROutlierDetector() *IQROutlierDetector {
	return &IQROutlierDetector{samples: make([]time.Duration, 0, 1000)}
}

func (iod *IQROutlierDetector) Add(lat time.Duration) {
	if lat > 0 {
		iod.samples = append(iod.samples, lat)
	}
}

func (iod *IQROutlierDetector) ComputeThresholds() (q1, q3, upperFence time.Duration) {
	if len(iod.samples) < 4 {
		return 0, 0, 0
	}
	sorted := make([]time.Duration, len(iod.samples))
	copy(sorted, iod.samples)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	n := len(sorted)
	q1 = sorted[n/4]
	q3 = sorted[(3*n)/4]
	iqr := q3 - q1
	upperFence = q3 + time.Duration(float64(iqr)*1.5)
	return q1, q3, upperFence
}

func (iod *IQROutlierDetector) IsOutlier(lat time.Duration) bool {
	_, _, upperFence := iod.ComputeThresholds()
	if upperFence == 0 {
		return false
	}
	return lat > upperFence
}
