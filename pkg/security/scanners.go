package security

import (
	"strings"

	"github.com/sujanv/LogSift/pkg/model"
)

var knownScanners = []string{
	"nikto",
	"sqlmap",
	"gobuster",
	"dirbuster",
	"wpscan",
	"acunetix",
	"nessus",
	"nmap",
}

func (td *ThreatDetector) InspectUserAgent(r *model.Record, ua string) *SecurityIncident {
	lowerUA := strings.ToLower(ua)
	for _, scanner := range knownScanners {
		if strings.Contains(lowerUA, scanner) {
			td.mu.Lock()
			defer td.mu.Unlock()

			inc := SecurityIncident{
				Timestamp:   r.Timestamp,
				ThreatType:  ThreatScannerBot,
				Severity:    "MEDIUM",
				ClientIP:    r.IP,
				RequestURI:  r.Path,
				Description: "Vulnerability scanner detected: " + scanner,
			}
			td.incidents = append(td.incidents, inc)
			return &inc
		}
	}
	return nil
}
