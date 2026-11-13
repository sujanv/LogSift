package diagnose

import (
	"sort"
	"strings"

	"github.com/sujanv/LogSift/pkg/model"
)

type AnomalyCulprit struct {
	Signature      string  `json:"signature"`
	BaselineCount  int64   `json:"baseline_count"`
	IncidentCount  int64   `json:"incident_count"`
	RelativeRisk   float64 `json:"relative_risk"`
	IsNewAnomaly   bool    `json:"is_new_anomaly"`
}

type DiagnosisReport struct {
	TotalBaseline int64            `json:"total_baseline"`
	TotalIncident int64            `json:"total_incident"`
	TopCulprits   []AnomalyCulprit `json:"top_culprits"`
}

func getSignature(r *model.Record) string {
	msg := r.Message
	if len(msg) > 60 {
		msg = msg[:60]
	}
	return strings.TrimSpace(r.Level.String() + " | " + r.Path + " | " + msg)
}

func Diagnose(baseline, incident []*model.Record) DiagnosisReport {
	baseCounts := make(map[string]int64)
	incCounts := make(map[string]int64)

	for _, r := range baseline {
		if r == nil {
			continue
		}
		sig := getSignature(r)
		baseCounts[sig]++
	}

	for _, r := range incident {
		if r == nil {
			continue
		}
		sig := getSignature(r)
		incCounts[sig]++
	}

	totBase := float64(len(baseline))
	totInc := float64(len(incident))
	if totBase == 0 {
		totBase = 1.0
	}
	if totInc == 0 {
		totInc = 1.0
	}

	var culprits []AnomalyCulprit
	for sig, incCnt := range incCounts {
		baseCnt := baseCounts[sig]
		incFreq := float64(incCnt) / totInc
		baseFreq := float64(baseCnt) / totBase

		rr := 1.0
		isNew := false
		if baseCnt == 0 {
			rr = incFreq * 1000.0
			isNew = true
		} else {
			rr = incFreq / baseFreq
		}

		if rr > 1.5 || isNew {
			culprits = append(culprits, AnomalyCulprit{
				Signature:     sig,
				BaselineCount: baseCnt,
				IncidentCount: incCnt,
				RelativeRisk:  rr,
				IsNewAnomaly:  isNew,
			})
		}
	}

	sort.Slice(culprits, func(i, j int) bool {
		return culprits[i].RelativeRisk > culprits[j].RelativeRisk
	})

	if len(culprits) > 10 {
		culprits = culprits[:10]
	}

	return DiagnosisReport{
		TotalBaseline: int64(len(baseline)),
		TotalIncident: int64(len(incident)),
		TopCulprits:   culprits,
	}
}
