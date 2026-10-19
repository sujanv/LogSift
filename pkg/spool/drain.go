package spool

import (
	"bufio"
	"os"
	"path/filepath"
)

func (ds *DiskSpooler) DrainAll() ([][]byte, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	path := filepath.Join(ds.dir, "spool.dat")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines [][]byte
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		b := make([]byte, len(scanner.Bytes()))
		copy(b, scanner.Bytes())
		lines = append(lines, b)
	}

	return lines, scanner.Err()
}
