package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookChannel string

const (
	ChannelSlack   WebhookChannel = "slack"
	ChannelDiscord WebhookChannel = "discord"
	ChannelGeneric WebhookChannel = "generic"
)

type WebhookDispatcher struct {
	webhookURL string
	channel    WebhookChannel
	client     *http.Client
}

func NewWebhookDispatcher(url string, channel WebhookChannel) *WebhookDispatcher {
	return &WebhookDispatcher{
		webhookURL: url,
		channel:    channel,
		client:     &http.Client{Timeout: 5 * time.Second},
	}
}

func (wd *WebhookDispatcher) SetClient(c *http.Client) {
	wd.client = c
}

func (wd *WebhookDispatcher) SendAlert(ctx context.Context, title, message string) error {
	var payload interface{}

	switch wd.channel {
	case ChannelSlack:
		payload = map[string]interface{}{
			"text": fmt.Sprintf("*LogSift Alert: %s*\n%s", title, message),
		}
	case ChannelDiscord:
		payload = map[string]interface{}{
			"content": fmt.Sprintf("**LogSift Alert: %s**\n%s", title, message),
		}
	default:
		payload = map[string]interface{}{
			"title":     title,
			"message":   message,
			"timestamp": time.Now().Format(time.RFC3339),
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, wd.webhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := wd.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook responded with HTTP %d", resp.StatusCode)
	}
	return nil
}
