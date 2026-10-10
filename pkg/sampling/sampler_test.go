package sampling

import (
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestTailSampler_Guarantees(t *testing.T) {
	// Sample rate 0% but always errors and slow > 200ms
	sampler := NewTailSampler(0.0, 200*time.Millisecond, true)

	errRec := &model.Record{Level: model.LevelError}
	if !sampler.ShouldSample(errRec) {
		t.Errorf("error record must always be retained")
	}

	slowRec := &model.Record{Level: model.LevelInfo, Latency: 350 * time.Millisecond}
	if !sampler.ShouldSample(slowRec) {
		t.Errorf("slow record must always be retained")
	}

	normalRec := &model.Record{Level: model.LevelInfo, Latency: 50 * time.Millisecond}
	if sampler.ShouldSample(normalRec) {
		t.Errorf("normal record should be dropped at 0%% rate")
	}
}

func TestTokenBucket(t *testing.T) {
	tb := NewTokenBucket(5, 10)
	for i := 0; i < 5; i++ {
		if !tb.Allow() {
			t.Errorf("expected bucket token %d to be allowed", i+1)
		}
	}
	if tb.Allow() {
		t.Errorf("expected 6th token to exceed capacity")
	}
}
