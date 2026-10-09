package reader

import (
	"bufio"
	"compress/gzip"
	"io"
)

const (
	// DefaultBufferSize is 64KB for high throughput streaming.
	DefaultBufferSize = 64 * 1024
)

// StreamReader reads lines from an io.Reader efficiently with transparent gzip decompression.
type StreamReader struct {
	reader     *bufio.Reader
	closer     io.Closer
	lineNum    int64
	byteOffset int64
	bytesRead  int64
	buf        []byte
}

// Option configures StreamReader.
type Option func(*StreamReader)

// WithBufferSize sets the internal buffer size.
func WithBufferSize(size int) Option {
	return func(s *StreamReader) {
		if size > 0 {
			s.buf = make([]byte, 0, size)
		}
	}
}

// NewStreamReader initializes a streaming reader with buffered reading and auto gzip detection.
func NewStreamReader(r io.Reader, opts ...Option) *StreamReader {
	bufReader := bufio.NewReaderSize(r, DefaultBufferSize)
	var closer io.Closer
	if c, ok := r.(io.Closer); ok {
		closer = c
	}

	// Sniff for gzip magic bytes (0x1f, 0x8b)
	peek, err := bufReader.Peek(2)
	if err == nil && len(peek) == 2 && peek[0] == 0x1f && peek[1] == 0x8b {
		gzReader, gzErr := gzip.NewReader(bufReader)
		if gzErr == nil {
			bufReader = bufio.NewReaderSize(gzReader, DefaultBufferSize)
			closer = gzReader
		}
	}

	sr := &StreamReader{
		reader: bufReader,
		closer: closer,
		buf:    make([]byte, 0, 4096),
	}
	for _, opt := range opts {
		opt(sr)
	}
	return sr
}

// NextLine reads the next line from the stream.
// It returns the line bytes, line number, byte offset, and error (if any).
// Lines of arbitrary length are supported without buffer overflow.
func (s *StreamReader) NextLine() ([]byte, int64, int64, error) {
	currentOffset := s.byteOffset
	s.buf = s.buf[:0]

	for {
		chunk, isPrefix, err := s.reader.ReadLine()
		if err != nil {
			if len(s.buf) > 0 && err == io.EOF {
				s.lineNum++
				s.byteOffset += int64(len(s.buf))
				s.bytesRead += int64(len(s.buf))
				return s.buf, s.lineNum, currentOffset, nil
			}
			return nil, s.lineNum, currentOffset, err
		}

		s.buf = append(s.buf, chunk...)
		s.byteOffset += int64(len(chunk))
		s.bytesRead += int64(len(chunk))

		if !isPrefix {
			s.byteOffset++
			s.bytesRead++
			break
		}
	}

	s.lineNum++
	return s.buf, s.lineNum, currentOffset, nil
}

// LineNum returns the current line number.
func (s *StreamReader) LineNum() int64 {
	return s.lineNum
}

// ByteOffset returns the current byte offset.
func (s *StreamReader) ByteOffset() int64 {
	return s.byteOffset
}

// BytesRead returns total bytes processed.
func (s *StreamReader) BytesRead() int64 {
	return s.bytesRead
}

// Close closes the underlying reader if it implements io.Closer.
func (s *StreamReader) Close() error {
	if s.closer != nil {
		return s.closer.Close()
	}
	return nil
}
