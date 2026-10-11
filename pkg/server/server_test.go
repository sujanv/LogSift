package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/stats"
)

func TestServerEndpoints(t *testing.T) {
	lat := stats.NewLatencyTracker()
	lat.Add(25 * time.Millisecond)

	errs := stats.NewErrorSummarizer()
	errs.Add("test database error", time.Now())

	srv := NewServer(":8080", lat, errs)
	handler := srv.Routes()

	// Test /health
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	// Test /api/v1/metrics
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/metrics", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("expected 200 OK on metrics, got %d", w2.Code)
	}
}
