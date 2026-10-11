package replay

import (
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type ReplayComparison struct {
	Path            string        `json:"path"`
	OriginalStatus  int           `json:"original_status"`
	ReplayedStatus  int           `json:"replayed_status"`
	OriginalLatency time.Duration `json:"original_latency"`
	ReplayedLatency time.Duration `json:"replayed_latency"`
	StatusMatch     bool          `json:"status_match"`
}

func CompareResponse(original *model.Record, replayedStatus int, replayedLatency time.Duration) ReplayComparison {
	return ReplayComparison{
		Path:            original.Path,
		OriginalStatus:  original.StatusCode,
		ReplayedStatus:  replayedStatus,
		OriginalLatency: original.Latency,
		ReplayedLatency: replayedLatency,
		StatusMatch:     original.StatusCode == replayedStatus,
	}
}
