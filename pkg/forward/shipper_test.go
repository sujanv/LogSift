package forward

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

type mockShipperTransport struct{}

func (m *mockShipperTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("ok")),
	}, nil
}

func TestLogShipper(t *testing.T) {
	shipper := NewLogShipper("https://logs.example.com/api/v1/push", DestLoki)
	shipper.SetClient(&http.Client{Transport: &mockShipperTransport{}})

	records := []*model.Record{
		{Message: "test ship 1"},
		{Message: "test ship 2"},
	}

	err := shipper.ShipBatch(context.Background(), records)
	if err != nil {
		t.Errorf("ShipBatch failed: %v", err)
	}
}
