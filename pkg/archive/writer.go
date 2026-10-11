package archive

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"sync"

	"github.com/sujanv/LogSift/pkg/model"
)

type ChunkArchiver struct {
	mu      sync.Mutex
	baseDir string
	mode    PartitionMode
	writers map[string]*gzip.Writer
	files   map[string]*os.File
}

func NewChunkArchiver(baseDir string, mode PartitionMode) *ChunkArchiver {
	return &ChunkArchiver{
		baseDir: baseDir,
		mode:    mode,
		writers: make(map[string]*gzip.Writer),
		files:   make(map[string]*os.File),
	}
}

func (ca *ChunkArchiver) Archive(r *model.Record) error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	part := DeterminePartitionPath(ca.mode, r)
	gzWriter, exists := ca.writers[part]

	if !exists {
		dir := filepath.Join(ca.baseDir, part)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		filePath := filepath.Join(dir, "logs.gz")
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		gzWriter = gzip.NewWriter(f)
		ca.writers[part] = gzWriter
		ca.files[part] = f
	}

	line := append(r.Raw, byte(10))
	_, err := gzWriter.Write(line)
	return err
}

func (ca *ChunkArchiver) Close() error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	for k, gw := range ca.writers {
		_ = gw.Close()
		_ = ca.files[k].Close()
	}
	return nil
}
