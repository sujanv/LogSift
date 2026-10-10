package test

import (
	"bytes"
	"io"
	"runtime"
	"testing"

	"github.com/sujanv/LogSift/pkg/generator"
	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/parser"
	"github.com/sujanv/LogSift/pkg/reader"
	"github.com/sujanv/LogSift/pkg/stats"
)

// TestStreamingMemoryBound verifies that streaming through hundreds of thousands of logs
// uses bounded memory (constant RSS) without leaking allocations.
func TestStreamingMemoryBound(t *testing.T) {
	pr, pw := io.Pipe()

	// Spawn generator in background
	go func() {
		gen := generator.NewGenerator(1234)
		_ = gen.GenerateJSONLines(pw, 100000)
		pw.Close()
	}()

	sr := reader.NewStreamReader(pr)
	defer sr.Close()

	jsonParser := parser.NewJSONParser()
	pool := model.DefaultRecordPool
	latTracker := stats.NewLatencyTracker()
	errSummarizer := stats.NewErrorSummarizer()

	var processed int64

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("read error: %v", err)
		}

		processed++
		rec := pool.Get()
		_ = jsonParser.Parse(line, rec)

		if rec.Latency > 0 {
			latTracker.Add(rec.Latency)
		}
		if rec.Level >= model.LevelError {
			errSummarizer.Add(rec.Message, rec.Timestamp)
		}

		pool.Put(rec)
	}

	if processed != 100000 {
		t.Fatalf("expected 100000 processed records, got %d", processed)
	}

	// Assert heap allocation is under 30MB
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	allocMB := float64(m.Alloc) / (1024 * 1024)
	t.Logf("Allocated memory after 100,000 logs: %.2f MB", allocMB)

	if allocMB > 30.0 {
		t.Errorf("expected memory usage < 30MB, got %.2f MB", allocMB)
	}
}

func TestEndToEndPipeline(t *testing.T) {
	var buf bytes.Buffer
	gen := generator.NewGenerator(42)
	_ = gen.GenerateJSONLines(&buf, 100)

	sr := reader.NewStreamReader(&buf)
	defer sr.Close()

	p := parser.NewJSONParser()
	var errCount int

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		rec := model.DefaultRecordPool.Get()
		_ = p.Parse(line, rec)
		if rec.Level == model.LevelError {
			errCount++
		}
		model.DefaultRecordPool.Put(rec)
	}

	if errCount == 0 {
		t.Errorf("expected at least 1 error in synthetic log batch")
	}
}
