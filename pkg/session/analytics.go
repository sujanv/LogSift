package session

import (
	"time"
)

type SessionAnalytics struct {
	TotalSessions int           `json:"total_sessions"`
	BounceRate    float64       `json:"bounce_rate"`
	AvgDuration   time.Duration `json:"avg_duration"`
	AvgRequests   float64       `json:"avg_requests"`
}

func AnalyzeSessions(sessions []*Session) SessionAnalytics {
	if len(sessions) == 0 {
		return SessionAnalytics{}
	}

	var bounces int
	var totalDur time.Duration
	var totalReqs int

	for _, s := range sessions {
		if s.Requests <= 1 {
			bounces++
		}
		totalDur += s.Duration
		totalReqs += s.Requests
	}

	return SessionAnalytics{
		TotalSessions: len(sessions),
		BounceRate:    (float64(bounces) / float64(len(sessions))) * 100.0,
		AvgDuration:   totalDur / time.Duration(len(sessions)),
		AvgRequests:   float64(totalReqs) / float64(len(sessions)),
	}
}
