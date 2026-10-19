package spool

import (
	"os"
	"path/filepath"
	"sync"
)

type DiskSpooler struct {
	mu       sync.Mutex
	dir      string
	file     *os.File
	count    int64
	capacity int64
}

func NewDiskSpooler(dir string, capacity int64) (*DiskSpooler, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(dir, "spool.dat"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	if capacity <= 0 {
		capacity = 50000
	}
	return &DiskSpooler{
		dir:      dir,
		file:     f,
		capacity: capacity,
	}, nil
}

func (ds *DiskSpooler) Push(entry []byte) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	line := append(entry, byte(10))
	_, err := ds.file.Write(line)
	if err == nil {
		ds.count++
	}
	return err
}

func (ds *DiskSpooler) Count() int64 {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	return ds.count
}

func (ds *DiskSpooler) Close() error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if ds.file != nil {
		return ds.file.Close()
	}
	return nil
}
