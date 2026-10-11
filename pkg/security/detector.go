package security

import (
	"sync"
	"time"
)

type ThreatType string

const (
	ThreatSQLInjection  ThreatType = "SQL_INJECTION"
	ThreatXSS           ThreatType = "CROSS_SITE_SCRIPTING"
	ThreatPathTraversal ThreatType = "PATH_TRAVERSAL"
	ThreatScannerBot    ThreatType = "VULNERABILITY_SCANNER"
	ThreatBruteForce    ThreatType = "BRUTE_FORCE"
)

type SecurityIncident struct {
	Timestamp   time.Time  `json:"timestamp"`
	ThreatType  ThreatType `json:"threat_type"`
	Severity    string     `json:"severity"`
	ClientIP    string     `json:"client_ip"`
	RequestURI  string     `json:"request_uri"`
	Description string     `json:"description"`
}

type ThreatDetector struct {
	mu        sync.Mutex
	incidents []SecurityIncident
	ipCounts  map[string]int
}

func NewThreatDetector() *ThreatDetector {
	return &ThreatDetector{
		incidents: make([]SecurityIncident, 0, 100),
		ipCounts:  make(map[string]int),
	}
}

func (td *ThreatDetector) Incidents() []SecurityIncident {
	td.mu.Lock()
	defer td.mu.Unlock()
	res := make([]SecurityIncident, len(td.incidents))
	copy(res, td.incidents)
	return res
}
