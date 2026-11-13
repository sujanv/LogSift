package test

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/cost"
	"github.com/sujanv/LogSift/pkg/dynlevel"
	"github.com/sujanv/LogSift/pkg/pseudonym"
)

func TestPhase7Regressions(t *testing.T) {
	ce := cost.NewCostEstimator()
	ce.Ingest(nil, 0)

	fr := dynlevel.NewFlightRecorder(0)
	fr.Record(nil)
	if len(fr.FlushOnFault(nil)) != 0 {
		t.Errorf("nil flight recorder flush should be empty")
	}

	ps := pseudonym.NewPseudonymizer("")
	ps.TokenizeRecord(nil)
	if ps.Tokenize("") != "" {
		t.Errorf("empty pseudonym string must return empty")
	}
}
