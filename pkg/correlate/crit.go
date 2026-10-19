package correlate

import "time"

type CriticalPathReport struct {
	TraceID       string        `json:"trace_id"`
	SpanCount     int           `json:"span_count"`
	Bottleneck    string        `json:"bottleneck"`
	BottleneckDur time.Duration `json:"bottleneck_duration"`
}

func (tg *TraceGraph) AnalyzeBottleneck() CriticalPathReport {
	rep := CriticalPathReport{
		TraceID:   tg.TraceID,
		SpanCount: len(tg.CriticalPath),
	}
	var maxDur time.Duration
	for _, s := range tg.CriticalPath {
		if s.Duration > maxDur {
			maxDur = s.Duration
			rep.Bottleneck = s.Operation
			rep.BottleneckDur = s.Duration
		}
	}
	return rep
}
