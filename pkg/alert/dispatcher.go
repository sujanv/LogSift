package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Dispatcher posts alert events to remote webhooks (Slack, Discord, generic HTTP).
type Dispatcher struct {
	webhookURL string
	client     *http.Client
}

// NewDispatcher creates an alert dispatcher.
func NewDispatcher(webhookURL string) *Dispatcher {
	return &Dispatcher{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 5 * time.Second},
	}
}

// Dispatch sends alert event to the webhook.
func (d *Dispatcher) Dispatch(event AlertEvent) error {
	if d.webhookURL == "" {
		return nil
	}

	payload := map[string]interface{}{
		"text": fmt.Sprintf("🚨 *[LogSift Alert]* `%s`: %s (Value: %.2f%%)",
			event.RuleName, event.Reason, event.Value*100.0),
		"event": event,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, d.webhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook responded with HTTP %d", resp.StatusCode)
	}
	return nil
}
