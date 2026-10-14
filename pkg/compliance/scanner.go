package compliance

import (
	"regexp"
	"strings"

	"github.com/sujanv/LogSift/pkg/model"
)

type ViolationType string

const (
	ViolationAWSAccessKey  ViolationType = "AWS_ACCESS_KEY"
	ViolationPrivateSSHKey ViolationType = "PRIVATE_KEY"
	ViolationGitHubToken   ViolationType = "GITHUB_TOKEN"
	ViolationPasswordLeak  ViolationType = "CLEARTEXT_PASSWORD"
	ViolationHIPAAData     ViolationType = "HEALTH_DATA_UNMASKED"
	ViolationGDPRData      ViolationType = "PERSONAL_DATA_UNMASKED"
)

type ComplianceFinding struct {
	LineNumber    int64         `json:"line_number"`
	ViolationType ViolationType `json:"violation_type"`
	Severity      string        `json:"severity"`
	Snippet       string        `json:"snippet"`
	Description   string        `json:"description"`
}

var (
	awsKeyRegex   = regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`)
	privKeyRegex  = regexp.MustCompile(`-----BEGIN (RSA|OPENSSH|DSA|EC)? ?PRIVATE KEY-----`)
	ghTokenRegex  = regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9_]{36,}\b`)
	pwdRegex      = regexp.MustCompile(`(?i)(password|secret|passwd|token)\s*[:=]\s*['"]?[^\s,;'"}\]]{6,}['"]?`)
	hipaaRegex    = regexp.MustCompile(`(?i)\b(diagnosis|patient_id|prescription|mrn|medical_record)\b`)
	gdprNameRegex = regexp.MustCompile(`(?i)\b(ssn|passport_number|driver_license|date_of_birth)\b`)
)

func ScanRecord(lineNum int64, r *model.Record) []ComplianceFinding {
	text := string(r.Raw) + " " + r.Message
	var findings []ComplianceFinding

	if awsKeyRegex.MatchString(text) {
		findings = append(findings, ComplianceFinding{
			LineNumber:    lineNum,
			ViolationType: ViolationAWSAccessKey,
			Severity:      "CRITICAL",
			Description:   "Hardcoded AWS Access Key ID exposed in log stream",
		})
	}

	if privKeyRegex.MatchString(text) {
		findings = append(findings, ComplianceFinding{
			LineNumber:    lineNum,
			ViolationType: ViolationPrivateSSHKey,
			Severity:      "CRITICAL",
			Description:   "Private cryptographic key material found in log payload",
		})
	}

	if ghTokenRegex.MatchString(text) {
		findings = append(findings, ComplianceFinding{
			LineNumber:    lineNum,
			ViolationType: ViolationGitHubToken,
			Severity:      "CRITICAL",
			Description:   "GitHub Personal Access Token leaked in log output",
		})
	}

	if pwdRegex.MatchString(text) && !strings.Contains(text, "password=***") {
		findings = append(findings, ComplianceFinding{
			LineNumber:    lineNum,
			ViolationType: ViolationPasswordLeak,
			Severity:      "HIGH",
			Description:   "Potential unmasked plaintext password or auth token parameter",
		})
	}

	if hipaaRegex.MatchString(text) {
		findings = append(findings, ComplianceFinding{
			LineNumber:    lineNum,
			ViolationType: ViolationHIPAAData,
			Severity:      "HIGH",
			Description:   "Protected Health Information (PHI/HIPAA) keyword detected without masking",
		})
	}

	if gdprNameRegex.MatchString(text) {
		findings = append(findings, ComplianceFinding{
			LineNumber:    lineNum,
			ViolationType: ViolationGDPRData,
			Severity:      "HIGH",
			Description:   "Sensitive personal identifier (GDPR) found in log",
		})
	}

	return findings
}
