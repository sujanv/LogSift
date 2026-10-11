package archive

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestChunkArchiver(t *testing.T) {
	tmpDir := t.TempDir()
	archiver := NewChunkArchiver(tmpDir, PartitionByDate)

	r := &model.Record{
		Timestamp: time.Date(2026, 10, 11, 10, 0, 0, 0, time.UTC),
		Raw:       []byte(`{"msg":"test archive entry"}`),
	}

	if err := archiver.Archive(r); err != nil {
		t.Fatalf("Archive error: %v", err)
	}
	_ = archiver.Close()

	expectedFile := filepath.Join(tmpDir, "2026-10-11", "logs.gz")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("expected archived file at %s", expectedFile)
	}
}
