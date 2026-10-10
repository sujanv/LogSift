package reader

import (
	"context"
	"strings"
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/parser"
)

func TestPipeline(t *testing.T) {
	input := `{"msg":"line 1"}
{"msg":"line 2"}
{"msg":"line 3"}
{"msg":"line 4"}`

	sr := NewStreamReader(strings.NewReader(input))
	pl := NewPipeline(parser.NewJSONParser(), 2)

	recordChan := make(chan *model.Record, 10)
	ctx := context.Background()

	go func() {
		_ = pl.ProcessStream(ctx, sr, recordChan)
		close(recordChan)
	}()

	var count int
	for r := range recordChan {
		count++
		model.DefaultRecordPool.Put(r)
	}

	if count != 4 {
		t.Errorf("expected 4 records from pipeline, got %d", count)
	}
}
