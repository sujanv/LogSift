package cluster

import (
	"strings"
	"testing"
)

func TestDrainClusterer(t *testing.T) {
	clusterer := NewDrainClusterer(0.6)

	m1 := "Failed password for root from 192.168.1.1 port 22 ssh2"
	m2 := "Failed password for admin from 192.168.1.2 port 22 ssh2"
	m3 := "Disk space warning on mount /data 95% full"

	c1 := clusterer.Add(m1)
	c2 := clusterer.Add(m2)
	c3 := clusterer.Add(m3)

	if c1.ID != c2.ID {
		t.Errorf("m1 and m2 should belong to the same cluster")
	}
	if c1.ID == c3.ID {
		t.Errorf("m3 should belong to a different cluster")
	}

	clusters := clusterer.Clusters()
	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(clusters))
	}

	if !strings.Contains(c1.Template, "<*>") {
		t.Errorf("expected wildcard in template: %s", c1.Template)
	}
}
