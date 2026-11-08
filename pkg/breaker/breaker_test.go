package breaker

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreakerTrips(t *testing.T) {
	cb := NewCircuitBreaker("test-exporter", 3, 50*time.Millisecond)

	fakeErr := errors.New("upstream timeout")

	for i := 0; i < 3; i++ {
		_ = cb.Execute(func() error {
			return fakeErr
		})
	}

	if cb.State() != StateOpen {
		t.Fatalf("expected breaker to be OPEN after 3 failures, got %s", cb.State())
	}

	err := cb.Execute(func() error {
		return nil
	})
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}

	time.Sleep(60 * time.Millisecond)
	if cb.State() != StateHalfOpen {
		t.Fatalf("expected breaker to transition to HALF-OPEN after cooldown, got %s", cb.State())
	}

	_ = cb.Execute(func() error { return nil })
	_ = cb.Execute(func() error { return nil })

	if cb.State() != StateClosed {
		t.Fatalf("expected breaker to be CLOSED after successful half-open probing, got %s", cb.State())
	}
}
