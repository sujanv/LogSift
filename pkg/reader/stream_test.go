package reader

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"testing"
)

func TestStreamReader_Basic(t *testing.T) {
	input := "line 1\nline 2\nline 3"
	sr := NewStreamReader(strings.NewReader(input))

	expected := []string{"line 1", "line 2", "line 3"}
	for i, exp := range expected {
		line, lineNum, _, err := sr.NextLine()
		if err != nil {
			t.Fatalf("unexpected error at line %d: %v", i+1, err)
		}
		if int64(i+1) != lineNum {
			t.Errorf("expected lineNum %d, got %d", i+1, lineNum)
		}
		if string(line) != exp {
			t.Errorf("expected %q, got %q", exp, string(line))
		}
	}

	_, _, _, err := sr.NextLine()
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestStreamReader_GzipCompressed(t *testing.T) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write([]byte("compressed log 1\ncompressed log 2\n"))
	if err != nil {
		t.Fatalf("gzip write failed: %v", err)
	}
	gw.Close()

	sr := NewStreamReader(&buf)
	defer sr.Close()

	l1, num1, _, err := sr.NextLine()
	if err != nil || string(l1) != "compressed log 1" || num1 != 1 {
		t.Errorf("expected 'compressed log 1', got %q (%v)", l1, err)
	}

	l2, num2, _, err := sr.NextLine()
	if err != nil || string(l2) != "compressed log 2" || num2 != 2 {
		t.Errorf("expected 'compressed log 2', got %q (%v)", l2, err)
	}
}

func TestStreamReader_VeryLongLine(t *testing.T) {
	longStr := strings.Repeat("A", 128*1024)
	input := longStr + "\nsecond line"
	sr := NewStreamReader(strings.NewReader(input))

	line1, lineNum1, _, err := sr.NextLine()
	if err != nil {
		t.Fatalf("failed reading long line: %v", err)
	}
	if lineNum1 != 1 {
		t.Errorf("expected lineNum 1, got %d", lineNum1)
	}
	if len(line1) != len(longStr) {
		t.Errorf("expected len %d, got %d", len(longStr), len(line1))
	}
}
