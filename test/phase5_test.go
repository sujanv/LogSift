package test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/archive"
	"github.com/sujanv/LogSift/pkg/correlate"
	"github.com/sujanv/LogSift/pkg/dedup"
	"github.com/sujanv/LogSift/pkg/enrich"
	"github.com/sujanv/LogSift/pkg/forward"
	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/slo"
	"github.com/sujanv/LogSift/pkg/spool"
	"github.com/sujanv/LogSift/pkg/stats"
)

type mockPhase5Transport struct{}

func (m *mockPhase5Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("ok")),
	}, nil
}

func TestPhase5Integration(t *testing.T) {
	d := dedup.NewWindowedDeduplicator(1 * time.Second)
	r1 := &model.Record{Level: model.LevelError, Message: "timeout"}
	r2 := &model.Record{Level: model.LevelError, Message: "timeout"}
	d.Ingest(r1)
	_, ev := d.Ingest(r2)
	if ev.Count != 2 {
		t.Errorf("expected count 2, got %d", ev.Count)
	}

	ua := enrich.ParseUserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) Mobile Safari/604.1")
	if ua.Device != "Mobile" || ua.OS != "iOS" {
		t.Errorf("expected Mobile iOS, got %+v", ua)
	}

	calc := slo.NewSLOCalculator(slo.SLIConfig{TargetAvailability: 0.99})
	m := calc.Evaluate([]*model.Record{{StatusCode: 200}, {StatusCode: 500}})
	if m.AvailabilityRatio != 0.5 {
		t.Errorf("expected 0.5 availability")
	}

	root := &correlate.SpanNode{SpanID: "root", Duration: 100 * time.Millisecond}
	c1 := &correlate.SpanNode{SpanID: "c1", ParentID: "root", Duration: 80 * time.Millisecond}
	tg := correlate.BuildTraceGraph("tr-1", []*correlate.SpanNode{root, c1})
	if len(tg.CriticalPath) != 2 {
		t.Errorf("expected critical path length 2")
	}

	bf := archive.NewChunkBloomFilter(512)
	bf.Add("k1")
	if !bf.Test("k1") {
		t.Errorf("bloom filter test failed")
	}

	tmpDir := t.TempDir()
	sp, _ := spool.NewDiskSpooler(tmpDir, 100)
	_ = sp.Push([]byte("entry"))
	if sp.Count() != 1 {
		t.Errorf("spool count mismatch")
	}
	_ = sp.Close()

	shipper := forward.NewLogShipper("https://localhost/push", forward.DestHTTPGeneric)
	shipper.SetClient(&http.Client{Transport: &mockPhase5Transport{}})
	err := shipper.ShipBatch(context.Background(), []*model.Record{r1})
	if err != nil {
		t.Errorf("shipper batch failed: %v", err)
	}

	td := stats.NewTDigest(50)
	td.Add(10.0)
	td.Add(20.0)
	td.Add(30.0)
	q := td.Quantile(0.5)
	if q < 10.0 || q > 30.0 {
		t.Errorf("unexpected quantile %f", q)
	}
}
