package diagnose

import (
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
