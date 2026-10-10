package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config represents runtime configuration options for LogSift.
type Config struct {
	DefaultFormat string        `json:"default_format"`
	DefaultLevel  string        `json:"default_level"`
	BufferSize    int           `json:"buffer_size"`
	Masking       MaskConfig    `json:"masking"`
	Alerting      AlertConfig   `json:"alerting"`
	Sampling      SampleConfig  `json:"sampling"`
}

type MaskConfig struct {
	Enabled   bool     `json:"enabled"`
	Strategy  string   `json:"strategy"` // redact, mask, hash, truncate
	Patterns  []string `json:"patterns"`
}

type AlertConfig struct {
	Enabled            bool          `json:"enabled"`
	ErrorRateThreshold float64       `json:"error_rate_threshold"`
	P99Threshold       time.Duration `json:"p99_threshold"`
	Window             time.Duration `json:"window"`
	WebhookURL         string        `json:"webhook_url"`
}

type SampleConfig struct {
	Rate          float64       `json:"rate"` // 0.0 to 1.0
	AlwaysErrors  bool          `json:"always_errors"`
	SlowThreshold time.Duration `json:"slow_threshold"`
}

// DefaultConfig provides sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		DefaultFormat: "auto",
		DefaultLevel:  "INFO",
		BufferSize:    64 * 1024,
		Masking: MaskConfig{
			Enabled:  false,
			Strategy: "redact",
		},
		Alerting: AlertConfig{
			Enabled:            false,
			ErrorRateThreshold: 0.05,
			P99Threshold:       500 * time.Millisecond,
			Window:             time.Minute,
		},
		Sampling: SampleConfig{
			Rate:          1.0,
			AlwaysErrors:  true,
			SlowThreshold: 200 * time.Millisecond,
		},
	}
}

// LoadFromFile reads and parses a JSON config file.
func LoadFromFile(path string) (*Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
