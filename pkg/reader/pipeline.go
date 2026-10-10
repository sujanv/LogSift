package reader

import (
	"context"
	"io"
	"runtime"
	"sync"

	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/parser"
)

// Pipeline distributes parsing work over a concurrent worker pool.
type Pipeline struct {
	workers int
	p       parser.Parser
}

// NewPipeline creates a Pipeline with specified number of worker goroutines.
func NewPipeline(p parser.Parser, workers int) *Pipeline {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	return &Pipeline{
		workers: workers,
		p:       p,
	}
}

type job struct {
	line    []byte
	lineNum int64
	offset  int64
}

// ProcessStream reads lines from sr and sends parsed records to recordChan.
func (pl *Pipeline) ProcessStream(ctx context.Context, sr *StreamReader, recordChan chan<- *model.Record) error {
	jobs := make(chan job, pl.workers*64)
	var wg sync.WaitGroup

	// Spawn worker pool
	for i := 0; i < pl.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool := model.DefaultRecordPool

			for j := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}

				rec := pool.Get()
				rec.LineNum = j.lineNum
				rec.Offset = j.offset
				_ = pl.p.Parse(j.line, rec)

				recordChan <- rec
			}
		}()
	}

	// Line feeder goroutine
	for {
		select {
		case <-ctx.Done():
			break
		default:
		}

		line, lineNum, offset, err := sr.NextLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			close(jobs)
			wg.Wait()
			return err
		}

		// Copy line slice for worker safety
		lineCopy := make([]byte, len(line))
		copy(lineCopy, line)

		jobs <- job{line: lineCopy, lineNum: lineNum, offset: offset}
	}

	close(jobs)
	wg.Wait()
	return nil
}
