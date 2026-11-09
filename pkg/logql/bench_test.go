package logql

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func BenchmarkLogQLEval(b *testing.B) {
	q, _ := Parse(`{level="error"} |= "timeout" != "probe"`)
	rec := &model.Record{
		Level:   model.LevelError,
		Message: "connection timeout during socket read from redis",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = q.Matches(rec)
	}
}
