package replay

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type mockTransport struct{}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Path == "/api/test" {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("OK")),
		}, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("Not Found")),
	}, nil
}

func TestTrafficReplayer(t *testing.T) {
	replayer := NewTrafficReplayer("http://localhost:8080", 1.0, 2)
	replayer.SetClient(&http.Client{Transport: &mockTransport{}})
	rec := &model.Record{Method: "GET", Path: "/api/test", StatusCode: 200, Latency: 15 * time.Millisecond}

	statusCode, elapsed, err := replayer.ReplaySingle(context.Background(), rec)
	if err != nil {
		t.Fatalf("ReplaySingle error: %v", err)
	}
	if statusCode != 200 {
		t.Errorf("expected 200 OK, got %d", statusCode)
	}

	comp := CompareResponse(rec, statusCode, elapsed)
	if !comp.StatusMatch {
		t.Errorf("status must match")
	}
}
