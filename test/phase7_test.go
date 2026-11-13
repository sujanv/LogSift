package test

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/cost"
	"github.com/sujanv/LogSift/pkg/diagnose"
	"github.com/sujanv/LogSift/pkg/dynlevel"
	"github.com/sujanv/LogSift/pkg/geoheat"
	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/otlprecv"
	"github.com/sujanv/LogSift/pkg/pseudonym"
)

func TestPhase7Integration(t *testing.T) {
	ce := cost.NewCostEstimator()
	ce.IngestService(&model.Record{Level: model.LevelInfo, Message: "test cost"}, 1024)
	rep := ce.Report(cost.ProviderAWSCloudWatch)
	if rep.TotalEvents != 1 {
		t.Errorf("cost event tracking mismatch")
	}

	fr := dynlevel.NewFlightRecorder(10)
	fr.Record(&model.Record{Level: model.LevelDebug, Message: "d1"})
	flushed := fr.FlushOnFault(&model.Record{Level: model.LevelError, Message: "err"})
	if len(flushed) != 2 {
		t.Errorf("flight recorder flush count mismatch")
	}

	ps := pseudonym.NewPseudonymizer("salt-secret")
	tok := ps.Tokenize("alice@domain.com")
	if tok == "" {
		t.Errorf("pseudonym token empty")
	}

	gh := geoheat.NewRegionalHeatmap(100 * time.Millisecond)
	gh.Record("DE", 50*time.Millisecond)
	if len(gh.Summary()) != 1 {
		t.Errorf("geoheat summary empty")
	}

	otlpLogs, err := otlprecv.ParseOTLPLogsJSON([]byte(`{"resourceLogs":[]}`))
	if err != nil || len(otlpLogs) != 0 {
		t.Errorf("otlp empty payload parse error: %v", err)
	}

	diag := diagnose.Diagnose(
		[]*model.Record{{Level: model.LevelInfo, Message: "ok"}},
		[]*model.Record{{Level: model.LevelError, Message: "crash"}},
	)
	if len(diag.TopCulprits) == 0 {
		t.Errorf("diagnose culprit identification failed")
	}
}
