package compliance

type AuditReport struct {
	TotalScanned    int64               `json:"total_scanned"`
	Violations      []ComplianceFinding `json:"violations"`
	Passed          bool                `json:"passed"`
	ComplianceScore float64             `json:"compliance_score"`
}

func GenerateReport(totalScanned int64, findings []ComplianceFinding) AuditReport {
	score := 100.0
	if totalScanned > 0 && len(findings) > 0 {
		deduction := (float64(len(findings)) / float64(totalScanned)) * 100.0 * 10.0
		if deduction > 100.0 {
			deduction = 100.0
		}
		score = 100.0 - deduction
	}
	return AuditReport{
		TotalScanned:    totalScanned,
		Violations:      findings,
		Passed:          len(findings) == 0,
		ComplianceScore: score,
	}
}
