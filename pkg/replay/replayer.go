package replay

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type ReplayStats struct {
	TotalReplayed int64
	Success       int64
	Failed        int64
	AvgLatency    time.Duration
}

type TrafficReplayer struct {
	targetBaseURL string
	client        *http.Client
	speedFactor   float64
	concurrency   int
}

func NewTrafficReplayer(targetBaseURL string, speedFactor float64, concurrency int) *TrafficReplayer {
	if speedFactor <= 0 {
		speedFactor = 1.0
	}
	if concurrency <= 0 {
		concurrency = 8
	}
	return &TrafficReplayer{
		targetBaseURL: strings.TrimRight(targetBaseURL, "/"),
		speedFactor:   speedFactor,
		concurrency:   concurrency,
		client:        &http.Client{Timeout: 10 * time.Second},
	}
}

func (tr *TrafficReplayer) SetClient(c *http.Client) {
	tr.client = c
}

func (tr *TrafficReplayer) ReplaySingle(ctx context.Context, r *model.Record) (int, time.Duration, error) {
	targetURL := tr.targetBaseURL + r.Path
	method := r.Method
	if method == "" {
		method = http.MethodGet
	}

	req, err := http.NewRequestWithContext(ctx, method, targetURL, nil)
	if err != nil {
		return 0, 0, err
	}

	t0 := time.Now()
	resp, err := tr.client.Do(req)
	elapsed := time.Since(t0)

	if err != nil {
		return 0, elapsed, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, elapsed, nil
}
