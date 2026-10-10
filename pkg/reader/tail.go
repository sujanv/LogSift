package reader

import (
	"context"
	"io"
	"os"
	"time"
)

// TailFollower provides live log tailing analogous to `tail -f`.
type TailFollower struct {
	filePath string
	pollInt  time.Duration
}

// NewTailFollower creates a TailFollower.
func NewTailFollower(filePath string, pollInterval time.Duration) *TailFollower {
	if pollInterval <= 0 {
		pollInterval = 250 * time.Millisecond
	}
	return &TailFollower{
		filePath: filePath,
		pollInt:  pollInterval,
	}
}

// Follow streams new lines as they are appended to the file.
func (tf *TailFollower) Follow(ctx context.Context, lineChan chan<- []byte) error {
	f, err := os.Open(tf.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end
	offset, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	sr := NewStreamReader(f)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		line, _, _, err := sr.NextLine()
		if err == nil {
			lineCopy := make([]byte, len(line))
			copy(lineCopy, line)
			lineChan <- lineCopy
			continue
		}

		if err == io.EOF {
			// Check if file was truncated/rotated
			stat, statErr := os.Stat(tf.filePath)
			if statErr == nil && stat.Size() < offset {
				// Reopen rotated file
				_ = f.Close()
				f, err = os.Open(tf.filePath)
				if err == nil {
					offset = 0
					sr = NewStreamReader(f)
					continue
				}
			}

			time.Sleep(tf.pollInt)
			continue
		}

		return err
	}
}
