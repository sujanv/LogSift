package slowlog

import (
	"sort"
	"time"
)

type QueryStat struct {
	Fingerprint string        `json:"fingerprint"`
	Count       int64         `json:"count"`
	TotalTime   time.Duration `json:"total_time"`
	AvgTime     time.Duration `json:"avg_time"`
}

type QueryAggregator struct {
	stats map[string]*QueryStat
}

func NewQueryAggregator() *QueryAggregator {
	return &QueryAggregator{stats: make(map[string]*QueryStat)}
}

func (qa *QueryAggregator) Ingest(sq *SlowQuery) {
	if sq == nil || sq.NormalizedSQL == "" {
		return
	}
	stat, exists := qa.stats[sq.NormalizedSQL]
	if !exists {
		stat = &QueryStat{Fingerprint: sq.NormalizedSQL}
		qa.stats[sq.NormalizedSQL] = stat
	}
	stat.Count++
	stat.TotalTime += sq.Duration
	stat.AvgTime = stat.TotalTime / time.Duration(stat.Count)
}

func (qa *QueryAggregator) TopSlowest(limit int) []*QueryStat {
	list := make([]*QueryStat, 0, len(qa.stats))
	for _, v := range qa.stats {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].TotalTime > list[j].TotalTime
	})
	if limit > 0 && len(list) > limit {
		return list[:limit]
	}
	return list
}
