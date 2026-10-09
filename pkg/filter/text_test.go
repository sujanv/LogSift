package filter

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestTextMatcher(t *testing.T) {
	rec := &model.Record{
		Message: "Connection reset by peer during handshake",
		Raw:     []byte(`{"msg":"Connection reset by peer during handshake"}`),
	}

	tm1 := NewTextMatcher("connection reset", false)
	if !tm1.Match(rec) {
		t.Errorf("expected case-insensitive match")
	}

	tm2 := NewTextMatcher("HANDSHAKE", true)
	if tm2.Match(rec) {
		t.Errorf("expected case-sensitive mismatch")
	}

	tm3 := NewTextMatcher("handshake", true)
	if !tm3.Match(rec) {
		t.Errorf("expected case-sensitive match")
	}
}

func TestRegexMatcher(t *testing.T) {
	rec := &model.Record{
		Message: "User id 48210 failed payment",
	}

	rm, err := NewRegexMatcher(`User id \d+ failed`)
	if err != nil {
		t.Fatalf("regex compile error: %v", err)
	}

	if !rm.Match(rec) {
		t.Errorf("expected regex match")
	}

	rmNoMatch, _ := NewRegexMatcher(`^Payment succeeded`)
	if rmNoMatch.Match(rec) {
		t.Errorf("expected regex NOT to match")
	}
}
