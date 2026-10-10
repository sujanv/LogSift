package stats

import (
	"testing"
	"time"
)

func TestHTTPAnalytics(t *testing.T) {
	h := NewHTTPAnalytics()

	h.Add(200, "/api/v1/users", 20*time.Millisecond)
	h.Add(200, "/api/v1/users", 25*time.Millisecond)
	h.Add(500, "/api/v1/checkout", 500*time.Millisecond)
	h.Add(404, "/unknown", 5*time.Millisecond)

	bd := h.StatusBreakdown()
	if bd["2xx"] != 2 {
		t.Errorf("expected 2xx count 2, got %d", bd["2xx"])
	}
	if bd["5xx"] != 1 {
		t.Errorf("expected 5xx count 1, got %d", bd["5xx"])
	}
	if bd["4xx"] != 1 {
		t.Errorf("expected 4xx count 1, got %d", bd["4xx"])
	}

	slowest := h.TopSlowestEndpoints(1)
	if len(slowest) != 1 || slowest[0].Path != "/api/v1/checkout" {
		t.Errorf("expected /api/v1/checkout to be slowest")
	}
}
