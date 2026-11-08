package throttle

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu           sync.Mutex
	capacity     float64
	tokens       float64
	refillPerSec float64
	lastRefill   time.Time
}

func NewTokenBucket(capacity float64, refillPerSec float64) *TokenBucket {
	if capacity <= 0 {
		capacity = 100
	}
	if refillPerSec <= 0 {
		refillPerSec = 50
	}
	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity,
		refillPerSec: refillPerSec,
		lastRefill:   time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1.0)
}

func (tb *TokenBucket) AllowN(n float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.lastRefill = now

	tb.tokens += elapsed * tb.refillPerSec
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}
	return false
}
