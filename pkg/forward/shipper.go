package forward

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

type ForwardDestination string

const (
	DestLoki          ForwardDestination = "loki"
	DestElasticsearch ForwardDestination = "elasticsearch"
	DestHTTPGeneric   ForwardDestination = "generic"
)

type LogShipper struct {
	endpoint string
	dest     ForwardDestination
	client   *http.Client
}

func NewLogShipper(endpoint string, dest ForwardDestination) *LogShipper {
	return &LogShipper{
		endpoint: strings.TrimRight(endpoint, "/"),
		dest:     dest,
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (ls *LogShipper) SetClient(c *http.Client) {
	ls.client = c
}

func (ls *LogShipper) ShipBatch(ctx context.Context, records []*model.Record) error {
	if len(records) == 0 {
		return nil
	}

	payload, err := json.Marshal(records)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ls.endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ls.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("shipper received HTTP %d", resp.StatusCode)
	}
	return nil
}
