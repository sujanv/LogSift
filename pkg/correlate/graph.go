package correlate

import (
	"time"
)

type SpanNode struct {
	SpanID    string        `json:"span_id"`
	ParentID  string        `json:"parent_id"`
	Operation string        `json:"operation"`
	Duration  time.Duration `json:"duration"`
	StartTime time.Time     `json:"start_time"`
	Children  []*SpanNode   `json:"children,omitempty"`
}

type TraceGraph struct {
	TraceID      string        `json:"trace_id"`
	Root         *SpanNode     `json:"root"`
	CriticalPath []*SpanNode   `json:"critical_path"`
	TotalLatency time.Duration `json:"total_latency"`
}

func BuildTraceGraph(traceID string, nodes []*SpanNode) *TraceGraph {
	nodeMap := make(map[string]*SpanNode)
	for _, n := range nodes {
		nodeMap[n.SpanID] = n
	}

	var root *SpanNode
	for _, n := range nodes {
		if n.ParentID == "" || nodeMap[n.ParentID] == nil {
			root = n
		} else {
			p := nodeMap[n.ParentID]
			p.Children = append(p.Children, n)
		}
	}

	tg := &TraceGraph{
		TraceID: traceID,
		Root:    root,
	}

	if root != nil {
		tg.CriticalPath = findCriticalPath(root)
		tg.TotalLatency = root.Duration
	}

	return tg
}

func findCriticalPath(node *SpanNode) []*SpanNode {
	if len(node.Children) == 0 {
		return []*SpanNode{node}
	}

	var longestChild *SpanNode
	for _, c := range node.Children {
		if longestChild == nil || c.Duration > longestChild.Duration {
			longestChild = c
		}
	}

	return append([]*SpanNode{node}, findCriticalPath(longestChild)...)
}
