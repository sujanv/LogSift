package diagnose

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func BenchmarkDiagnose(b *testing.B) {
	base := make([]*model.Record, 100)
	inc := make([]*model.Record, 100)
	for i := 0; i < 100; i++ {
		base[i] = &model.Record{Level: model.LevelInfo, Message: "healthy operation"}
		inc[i] = &model.Record{Level: model.LevelError, Message: "sporadic error"}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = Diagnose(base, inc)
	}
}
