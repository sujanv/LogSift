package stats

import (
	"testing"
	"time"
)

func TestNormalizeMessage(t *testing.T) {
	raw := "connection refused to 192.168.1.100:5432 for user 49201 session 123e4567-e89b-12d3-a456-426614174000"
	norm := NormalizeMessage(raw)

	expected := "connection refused to <IP>:<NUM> for user <NUM> session <UUID>"
	if norm != expected {
		t.Errorf("NormalizeMessage failed:\ngot:  %q\nwant: %q", norm, expected)
	}
}

func TestErrorSummarizer(t *testing.T) {
	s := NewErrorSummarizer()
	now := time.Now()

	s.Add("database timeout on query SELECT * FROM orders WHERE id = 12", now)
	s.Add("database timeout on query SELECT * FROM orders WHERE id = 99", now.Add(time.Second))
	s.Add("disk full on /var/log", now)

	top := s.TopClusters(5)
	if len(top) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(top))
	}

	if top[0].Count != 2 {
		t.Errorf("expected top cluster count 2, got %d", top[0].Count)
	}
}
