package livetail

import (
	"bufio"
	"context"
	"io"
	"os"
	"sync"
	"time"
)

type StreamEvent struct {
	Source    string
	Line      string
	Timestamp time.Time
	Err       error
}

type Multiplexer struct {
	mu      sync.Mutex
	sources map[string]io.Reader
	outCh   chan StreamEvent
}

func NewMultiplexer(bufSize int) *Multiplexer {
	if bufSize <= 0 {
		bufSize = 1000
	}
	return &Multiplexer{
		sources: make(map[string]io.Reader),
		outCh:   make(chan StreamEvent, bufSize),
	}
}

func (m *Multiplexer) AddReader(name string, r io.Reader) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sources[name] = r
}

func (m *Multiplexer) AddFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	m.AddReader(path, f)
	return nil
}

func (m *Multiplexer) Start(ctx context.Context) <-chan StreamEvent {
	m.mu.Lock()
	defer m.mu.Unlock()

	var wg sync.WaitGroup
	for name, r := range m.sources {
		wg.Add(1)
		go func(srcName string, reader io.Reader) {
			defer wg.Done()
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				select {
				case <-ctx.Done():
					return
				case m.outCh <- StreamEvent{
					Source:    srcName,
					Line:      scanner.Text(),
					Timestamp: time.Now(),
				}:
				}
			}
			if err := scanner.Err(); err != nil {
				select {
				case <-ctx.Done():
				case m.outCh <- StreamEvent{
					Source: srcName,
					Err:    err,
				}:
				}
			}
		}(name, r)
	}

	go func() {
		wg.Wait()
		close(m.outCh)
	}()

	return m.outCh
}
