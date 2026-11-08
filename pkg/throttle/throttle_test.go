package throttle

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestTokenBucket(t *testing.T) {
	tb := NewTokenBucket(2, 100)
	if !tb.Allow() {
		t.Errorf("token 1 must be allowed")
	}
	if !tb.Allow() {
		t.Errorf("token 2 must be allowed")
	}
	if tb.Allow() {
		t.Errorf("token 3 should be rejected when empty")
	}

	time.Sleep(25 * time.Millisecond)
	if !tb.Allow() {
		t.Errorf("token should be allowed after refill")
	}
}

func TestErrorPreservingSampler(t *testing.T) {
	sampler := NewErrorPreservingSampler(0.0)

	errRec := &model.Record{Level: model.LevelError, StatusCode: 500}
	if !sampler.ShouldSample(errRec) {
		t.Errorf("error record MUST always be preserved")
	}

	infoRec := &model.Record{Level: model.LevelInfo, StatusCode: 200}
	if sampler.ShouldSample(infoRec) {
		t.Errorf("normal record should not be sampled at 0%% rate")
	}
}
