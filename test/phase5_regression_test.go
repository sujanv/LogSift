package test

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/archive"
	"github.com/sujanv/LogSift/pkg/slo"
)

func TestPhase5Regressions(t *testing.T) {
	alert := slo.EvaluateBurnRate(15.0)
	if alert == nil || alert.Severity != "CRITICAL" {
		t.Errorf("expected critical alert on 15x burn rate")
	}

	bf := archive.NewChunkBloomFilter(256)
	bf.Add("req-id-1")
	if !bf.Test("req-id-1") {
		t.Errorf("bloom filter regression")
	}
}
