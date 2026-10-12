package test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/sujanv/LogSift/pkg/geodb"
	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/replay"
	"github.com/sujanv/LogSift/pkg/security"
)

type mockIntegrationTransport struct{}

func (m *mockIntegrationTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("OK")),
	}, nil
}

func TestSecurityAndReplayIntegration(t *testing.T) {
	td := security.NewThreatDetector()
	recAttack := &model.Record{Path: "/login?user=admin' OR 1=1--", Message: "login attempt"}
	inc := td.InspectRecord(recAttack)
	if inc == nil || inc.ThreatType != security.ThreatSQLInjection {
		t.Errorf("security threat must be detected")
	}

	geo := geodb.NewResolver()
	info := geo.Lookup("192.168.1.1")
	if !info.IsPrivate {
		t.Errorf("LAN IP must be private")
	}

	tr := replay.NewTrafficReplayer("http://localhost:8080", 1.0, 2)
	tr.SetClient(&http.Client{Transport: &mockIntegrationTransport{}})
	code, _, err := tr.ReplaySingle(context.Background(), &model.Record{Path: "/test"})
	if err != nil || code != 200 {
		t.Errorf("replay failed: code=%d, err=%v", code, err)
	}
}
