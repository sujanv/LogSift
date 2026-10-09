package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

// JSONParser parses arbitrary structured JSON log entries.
type JSONParser struct{}

// NewJSONParser creates a new JSON parser instance.
func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

// Name returns the format name.
func (p *JSONParser) Name() string {
	return "json"
}

// Parse extracts fields from a JSON log line into record.
func (p *JSONParser) Parse(line []byte, record *model.Record) error {
	record.Raw = line

	if record.Fields == nil {
		record.Fields = make(map[string]interface{}, 16)
	}

	if err := json.Unmarshal(line, &record.Fields); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	for k, v := range record.Fields {
		lowerKey := strings.ToLower(k)

		switch lowerKey {
		case "timestamp", "time", "@timestamp", "ts", "datetime", "date":
			record.Timestamp = parseTimestamp(v)
		case "level", "lvl", "severity", "log_level":
			record.Level = parseLevelValue(v)
		case "message", "msg", "log", "text", "content":
			if s, ok := v.(string); ok {
				record.Message = s
			}
		case "latency", "latency_ms", "duration", "duration_ms", "elapsed", "response_time":
			record.Latency = parseLatency(v)
		case "status", "status_code", "code", "http_status":
			record.StatusCode = parseInt(v)
		case "method", "http_method":
			if s, ok := v.(string); ok {
				record.Method = strings.ToUpper(s)
			}
		case "path", "uri", "url", "route", "endpoint":
			if s, ok := v.(string); ok {
				record.Path = s
			}
		case "ip", "client_ip", "remote_ip", "remote_addr":
			if s, ok := v.(string); ok {
				record.IP = s
			}
		case "trace_id", "traceid", "trace":
			if s, ok := v.(string); ok {
				record.TraceID = s
			}
		case "span_id", "spanid":
			if s, ok := v.(string); ok {
				record.SpanID = s
			}
		case "user_id", "userid", "user":
			if s, ok := v.(string); ok {
				record.UserID = s
			}
		}
	}

	return nil
}

func parseTimestamp(val interface{}) time.Time {
	switch v := val.(type) {
	case string:
		formats := []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02 15:04:05.999999999",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
			time.RubyDate,
			time.ANSIC,
		}
		for _, f := range formats {
			if t, err := time.Parse(f, v); err == nil {
				return t
			}
		}
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return parseUnixTimestamp(n)
		}
	case float64:
		return parseUnixTimestamp(int64(v))
	case int64:
		return parseUnixTimestamp(v)
	}
	return time.Time{}
}

func parseUnixTimestamp(n int64) time.Time {
	switch {
	case n > 1e17: // nanoseconds
		return time.Unix(0, n)
	case n > 1e14: // microseconds
		return time.Unix(0, n*1000)
	case n > 1e11: // milliseconds
		return time.Unix(n/1000, (n%1000)*1e6)
	default: // seconds
		return time.Unix(n, 0)
	}
}

func parseLevelValue(val interface{}) model.LogLevel {
	switch v := val.(type) {
	case string:
		return model.ParseLevel(v)
	case float64:
		// Numeric log level e.g. bunyan/pino (10=trace, 20=debug, 30=info, 40=warn, 50=error, 60=fatal)
		lvl := int(v)
		switch {
		case lvl <= 20:
			return model.LevelDebug
		case lvl <= 30:
			return model.LevelInfo
		case lvl <= 40:
			return model.LevelWarn
		case lvl <= 50:
			return model.LevelError
		default:
			return model.LevelFatal
		}
	}
	return model.LevelUnknown
}

func parseLatency(val interface{}) time.Duration {
	switch v := val.(type) {
	case string:
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return time.Duration(f * float64(time.Millisecond))
		}
	case float64:
		// assume milliseconds
		return time.Duration(v * float64(time.Millisecond))
	case int64:
		return time.Duration(v) * time.Millisecond
	}
	return 0
}

func parseInt(val interface{}) int {
	switch v := val.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case string:
		n, _ := strconv.Atoi(v)
		return n
	}
	return 0
}
