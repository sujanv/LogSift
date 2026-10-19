package spool

import (
	"testing"
)

func TestDiskSpooler(t *testing.T) {
	tmpDir := t.TempDir()
	spooler, err := NewDiskSpooler(tmpDir, 1000)
	if err != nil {
		t.Fatalf("NewDiskSpooler error: %v", err)
	}

	_ = spooler.Push([]byte(`{"msg":"first spooled log"}`))
	_ = spooler.Push([]byte(`{"msg":"second spooled log"}`))

	if spooler.Count() != 2 {
		t.Errorf("expected 2 spooled logs, got %d", spooler.Count())
	}

	lines, err := spooler.DrainAll()
	if err != nil || len(lines) != 2 {
		t.Fatalf("failed to drain spool lines: %v, count %d", err, len(lines))
	}
	_ = spooler.Close()
}
