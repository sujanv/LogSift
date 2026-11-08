package livetail

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestMultiplexer(t *testing.T) {
	mux := NewMultiplexer(10)
	r1 := strings.NewReader("line1 from r1\nline2 from r1\n")
	r2 := strings.NewReader("line1 from r2\n")

	mux.AddReader("r1", r1)
	mux.AddReader("r2", r2)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := mux.Start(ctx)
	var events []StreamEvent
	for ev := range ch {
		events = append(events, ev)
	}

	if len(events) != 3 {
		t.Fatalf("expected 3 multiplexed lines, got %d", len(events))
	}
}
