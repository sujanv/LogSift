package test

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/anomaly"
	"github.com/sujanv/LogSift/pkg/cluster"
	"github.com/sujanv/LogSift/pkg/compliance"
	"github.com/sujanv/LogSift/pkg/crypto"
	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/slowlog"
)

func TestPhase4Integration(t *testing.T) {
	iqr := anomaly.NewIQROutlierDetector()
	for i := 1; i <= 50; i++ {
		iqr.Add(time.Duration(i) * time.Millisecond)
	}
	if !iqr.IsOutlier(500 * time.Millisecond) {
		t.Errorf("500ms must be flagged as outlier")
	}

	dc := cluster.NewDrainClusterer(0.6)
	c1 := dc.Add("connection dropped after 5000 ms to host 10.0.0.1")
	c2 := dc.Add("connection dropped after 1200 ms to host 10.0.0.2")
	if c1.ID != c2.ID {
		t.Errorf("expected same cluster for identical structures")
	}

	rec := &model.Record{Message: "Using AWS Key AKIAIOSFODNN7EXAMPLE"}
	findings := compliance.ScanRecord(1, rec)
	if len(findings) == 0 || findings[0].ViolationType != compliance.ViolationAWSAccessKey {
		t.Errorf("expected compliance violation for AWS key")
	}

	chain := crypto.NewTamperChain("test-secret")
	b1 := chain.Append([]byte("log1"))
	b2 := chain.Append([]byte("log2"))
	valid, _, err := chain.Verify([][]byte{[]byte("log1"), []byte("log2")})
	if !valid || err != nil || b1.Index != 1 || b2.Index != 2 {
		t.Errorf("tamper chain verification failed: %v", err)
	}

	pg := slowlog.ParsePostgresLine("2026-10-14 10:00:00 UTC [123]: LOG:  duration: 800.0 ms  statement: SELECT * FROM items WHERE id = 42;")
	if pg == nil || pg.NormalizedSQL != "SELECT * FROM ITEMS WHERE ID = ?;" {
		t.Errorf("slow query parse failed: %+v", pg)
	}
}
