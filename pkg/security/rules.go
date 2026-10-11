package security

import (
	"net/url"
	"regexp"

	"github.com/sujanv/LogSift/pkg/model"
)

var (
	sqliRegex = regexp.MustCompile(`(?i)(union\s+select|select\s+.*\s+from|'\s+or\s+1=1|;\s*drop\s+table|sleep\(\d+\)|benchmark\()`)
	xssRegex  = regexp.MustCompile(`(?i)(<script.*?>|javascript:|onload=|onerror=|alert\(|document\.cookie)`)
	pathRegex = regexp.MustCompile(`(\.\./|\.\.\\|\.\.%2f|%2e%2e%2f|/etc/passwd|win\.ini)`)
)

func (td *ThreatDetector) InspectRecord(r *model.Record) *SecurityIncident {
	target := r.Path + " " + r.Message
	if unescaped, err := url.QueryUnescape(target); err == nil {
		target += " " + unescaped
	}

	td.mu.Lock()
	defer td.mu.Unlock()

	if sqliRegex.MatchString(target) {
		inc := SecurityIncident{
			Timestamp:   r.Timestamp,
			ThreatType:  ThreatSQLInjection,
			Severity:    "HIGH",
			ClientIP:    r.IP,
			RequestURI:  r.Path,
			Description: "SQL Injection attack pattern detected in query parameter or URI",
		}
		td.incidents = append(td.incidents, inc)
		return &inc
	}

	if xssRegex.MatchString(target) {
		inc := SecurityIncident{
			Timestamp:   r.Timestamp,
			ThreatType:  ThreatXSS,
			Severity:    "HIGH",
			ClientIP:    r.IP,
			RequestURI:  r.Path,
			Description: "Cross-Site Scripting (XSS) payload detected in request",
		}
		td.incidents = append(td.incidents, inc)
		return &inc
	}

	if pathRegex.MatchString(target) {
		inc := SecurityIncident{
			Timestamp:   r.Timestamp,
			ThreatType:  ThreatPathTraversal,
			Severity:    "CRITICAL",
			ClientIP:    r.IP,
			RequestURI:  r.Path,
			Description: "Directory path traversal attempt detected",
		}
		td.incidents = append(td.incidents, inc)
		return &inc
	}

	return nil
}
