package correlate

import (
	"testing"
	"time"
)

func TestBuildTraceGraph(t *testing.T) {
	n1 := &SpanNode{SpanID: "root", Operation: "GET /checkout", Duration: 300 * time.Millisecond}
	n2 := &SpanNode{SpanID: "s1", ParentID: "root", Operation: "sql:query", Duration: 180 * time.Millisecond}
	n3 := &SpanNode{SpanID: "s2", ParentID: "root", Operation: "cache:get", Duration: 20 * time.Millisecond}

	tg := BuildTraceGraph("trace-1", []*SpanNode{n1, n2, n3})
	if tg.Root == nil || tg.Root.SpanID != "root" {
		t.Fatalf("expected root node")
	}
	if len(tg.CriticalPath) != 2 || tg.CriticalPath[1].SpanID != "s1" {
		t.Errorf("critical path must include slowest child s1")
	}

	bottleneck := tg.AnalyzeBottleneck()
	if bottleneck.Bottleneck != "GET /checkout" {
		t.Errorf("expected GET /checkout as overall longest span")
	}
}
