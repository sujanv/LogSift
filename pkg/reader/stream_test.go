package reader

import (
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

func TestStreamReader_VeryLongLine(t *testing.T) {
	// Line longer than standard 64KB bufio buffer
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

	line2, lineNum2, _, err := sr.NextLine()
	if err != nil {
		t.Fatalf("failed reading line 2: %v", err)
	}
	if lineNum2 != 2 || string(line2) != "second line" {
		t.Errorf("expected 'second line', got %q", string(line2))
	}
}

func TestStreamReader_EmptyLines(t *testing.T) {
	input := "\n\nline\n"
	sr := NewStreamReader(strings.NewReader(input))

	l1, _, _, err := sr.NextLine()
	if err != nil || len(l1) != 0 {
		t.Errorf("expected empty line 1, got %q", l1)
	}
	l2, _, _, err := sr.NextLine()
	if err != nil || len(l2) != 0 {
		t.Errorf("expected empty line 2, got %q", l2)
	}
	l3, _, _, err := sr.NextLine()
	if err != nil || string(l3) != "line" {
		t.Errorf("expected 'line', got %q", l3)
	}
}
