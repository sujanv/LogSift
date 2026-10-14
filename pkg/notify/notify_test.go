package notify

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockNotifyTransport struct{}

func (m *mockNotifyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("ok")),
	}, nil
}

func TestWebhookDispatcher(t *testing.T) {
	disp := NewWebhookDispatcher("https://hooks.slack.com/services/test", ChannelSlack)
	disp.SetClient(&http.Client{Transport: &mockNotifyTransport{}})

	err := disp.SendAlert(context.Background(), "High Error Rate", "Errors exceeded 15% threshold")
	if err != nil {
		t.Errorf("SendAlert failed: %v", err)
	}
}
