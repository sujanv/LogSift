package query

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestParsePipeline(t *testing.T) {
	query := `filter(level = ERROR) | topk(path, 5)`
	stages, err := ParsePipeline(query)
	if err != nil {
		t.Fatalf("ParsePipeline error: %v", err)
	}
	if len(stages) != 2 {
		t.Fatalf("expected 2 stages, got %d", len(stages))
	}
	if stages[0].Command != "filter" || stages[0].Args[0] != "level = ERROR" {
		t.Errorf("filter stage mismatch: %+v", stages[0])
	}
	if stages[1].Command != "topk" || stages[1].Args[0] != "path" || stages[1].Args[1] != "5" {
		t.Errorf("topk stage mismatch: %+v", stages[1])
	}
}

func TestPipelineExecution(t *testing.T) {
	stages, _ := ParsePipeline("filter(status>=500) | topk(path, 2)")
	exec, err := CompilePipeline(stages)
	if err != nil {
		t.Fatalf("CompilePipeline error: %v", err)
	}

	r1 := &model.Record{StatusCode: 500, Path: "/api/checkout"}
	r2 := &model.Record{StatusCode: 502, Path: "/api/checkout"}
	r3 := &model.Record{StatusCode: 500, Path: "/api/login"}
	r4 := &model.Record{StatusCode: 200, Path: "/api/checkout"}

	exec.Process(r1)
	exec.Process(r2)
	exec.Process(r3)
	exec.Process(r4)

	res := exec.TopKResult()
	if len(res) != 2 {
		t.Fatalf("expected 2 topk items, got %d", len(res))
	}
	if res[0].Key != "/api/checkout" || res[0].Count != 2 {
		t.Errorf("top 1 should be /api/checkout with 2, got %+v", res[0])
	}
}
