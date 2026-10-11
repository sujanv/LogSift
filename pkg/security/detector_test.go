package security

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestThreatDetector_SQLi(t *testing.T) {
	td := NewThreatDetector()
	r := &model.Record{
		Timestamp: time.Now(),
		Path:      "/api/users?id=1'%20OR%201=1--",
		IP:        "198.51.100.1",
	}

	inc := td.InspectRecord(r)
	if inc == nil || inc.ThreatType != ThreatSQLInjection {
		t.Errorf("expected SQLi threat detection")
	}
}

func TestThreatDetector_PathTraversal(t *testing.T) {
	td := NewThreatDetector()
	r := &model.Record{
		Timestamp: time.Now(),
		Path:      "/static/../../etc/passwd",
		IP:        "198.51.100.2",
	}

	inc := td.InspectRecord(r)
	if inc == nil || inc.ThreatType != ThreatPathTraversal {
		t.Errorf("expected Path Traversal threat detection")
	}
}

func TestThreatDetector_ScannerBot(t *testing.T) {
	td := NewThreatDetector()
	r := &model.Record{IP: "198.51.100.3"}

	inc := td.InspectUserAgent(r, "Mozilla/5.0 (compatible; sqlmap/1.5.2#stable)")
	if inc == nil || inc.ThreatType != ThreatScannerBot {
		t.Errorf("expected Scanner Bot threat detection")
	}
}
