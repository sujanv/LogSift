package test

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/correlate"
	"github.com/sujanv/LogSift/pkg/diff"
	"github.com/sujanv/LogSift/pkg/mask"
	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/sampling"
)

func TestFeatureSuite(t *testing.T) {
	// 1. Masking
	m := mask.NewMasker(mask.RedactStrategy{})
	r := &model.Record{Message: "customer email is test@company.com"}
	m.MaskRecord(r)
	if r.Message == "customer email is test@company.com" {
		t.Errorf("masking failed")
	}

	// 2. Sampling
	sampler := sampling.NewTailSampler(0.0, 100*time.Millisecond, true)
	if !sampler.ShouldSample(&model.Record{Level: model.LevelError}) {
		t.Errorf("sampling must keep error")
	}

	// 3. Diff
	s1 := &diff.Snapshot{TotalLogs: 100, ErrorCount: 5}
	s2 := &diff.Snapshot{TotalLogs: 100, ErrorCount: 15}
	d := diff.Compare(s1, s2)
	if d.ErrorRateDiff <= 0 {
		t.Errorf("diff error rate calculation incorrect")
	}

	// 4. Correlate
	tc := correlate.NewTraceCorrelator()
	tc.Ingest(&model.Record{TraceID: "t-1", Timestamp: time.Now(), Message: "auth success"})
	if len(tc.GetTrace("t-1")) != 1 {
		t.Errorf("trace correlation failed")
	}
}
