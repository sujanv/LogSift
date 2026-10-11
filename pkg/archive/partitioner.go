package archive

import (
	"path/filepath"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type PartitionMode string

const (
	PartitionByDate  PartitionMode = "date"
	PartitionByHour  PartitionMode = "hour"
	PartitionByLevel PartitionMode = "level"
)

func DeterminePartitionPath(mode PartitionMode, r *model.Record) string {
	ts := r.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}

	switch mode {
	case PartitionByHour:
		return filepath.Join(ts.Format("2006-01-02"), ts.Format("15"))
	case PartitionByLevel:
		return r.Level.String()
	default:
		return ts.Format("2006-01-02")
	}
}
