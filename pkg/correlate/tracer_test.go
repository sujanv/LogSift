package correlate

import (
	"strings"
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestTraceCorrelator(t *testing.T) {
	tc := NewTraceCorrelator()
	t0 := time.Now()

	r1 := &model.Record{TraceID: "tx-42", Timestamp: t0, Message: "gateway request received"}
	r2 := &model.Record{TraceID: "tx-42", Timestamp: t0.Add(50 * time.Millisecond), Message: "database query executed"}
	r3 := &model.Record{TraceID: "other", Timestamp: t0, Message: "unrelated log"}

	tc.Ingest(r1)
	tc.Ingest(r2)
	tc.Ingest(r3)

	trace := tc.GetTrace("tx-42")
	if len(trace) != 2 {
		t.Fatalf("expected 2 trace entries, got %d", len(trace))
	}

	rendered := FormatTraceTimeline("tx-42", trace)
	if !strings.Contains(rendered, "gateway request received") || !strings.Contains(rendered, "database query executed") {
		t.Errorf("trace rendering failed:\n%s", rendered)
	}
}
