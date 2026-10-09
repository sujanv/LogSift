package parser

import (
	"strings"

	"github.com/sujanv/LogSift/pkg/model"
)

// LogfmtParser parses logs formatted as key=value pairs.
type LogfmtParser struct{}

// NewLogfmtParser creates a new LogfmtParser.
func NewLogfmtParser() *LogfmtParser {
	return &LogfmtParser{}
}

// Name returns the format name.
func (p *LogfmtParser) Name() string {
	return "logfmt"
}

// Parse extracts key=value pairs from line into record.
func (p *LogfmtParser) Parse(line []byte, record *model.Record) error {
	record.Raw = line
	if record.Fields == nil {
		record.Fields = make(map[string]interface{}, 16)
	}

	s := string(line)
	n := len(s)
	i := 0

	for i < n {
		// Skip whitespace
		for i < n && (s[i] == ' ' || s[i] == '\t' || s[i] == '\r' || s[i] == '\n') {
			i++
		}
		if i >= n {
			break
		}

		// Read key
		keyStart := i
		for i < n && s[i] != '=' && s[i] != ' ' && s[i] != '\t' {
			i++
		}
		if i >= n || s[i] != '=' {
			// Bare key without value or end of line
			if i > keyStart {
				record.Fields[s[keyStart:i]] = true
			}
			continue
		}

		key := s[keyStart:i]
		i++ // skip '='

		// Read value
		var val string
		if i < n && s[i] == '"' {
			i++ // skip opening quote
			valStart := i
			var sb strings.Builder
			for i < n && s[i] != '"' {
				if s[i] == '\\' && i+1 < n {
					i++
					sb.WriteByte(s[i])
				} else {
					sb.WriteByte(s[i])
				}
				i++
			}
			if i < n && s[i] == '"' {
				i++ // skip closing quote
			}
			val = sb.String()
			_ = valStart
		} else {
			valStart := i
			for i < n && s[i] != ' ' && s[i] != '\t' && s[i] != '\r' && s[i] != '\n' {
				i++
			}
			val = s[valStart:i]
		}

		record.Fields[key] = val
		lowerKey := strings.ToLower(key)

		switch lowerKey {
		case "time", "ts", "timestamp", "date":
			record.Timestamp = parseTimestamp(val)
		case "level", "lvl", "severity":
			record.Level = model.ParseLevel(val)
		case "msg", "message", "log":
			record.Message = val
		case "latency", "duration", "elapsed", "latency_ms":
			record.Latency = parseLatency(val)
		case "status", "status_code", "code":
			record.StatusCode = parseInt(val)
		case "method":
			record.Method = strings.ToUpper(val)
		case "path", "uri", "url":
			record.Path = val
		case "ip", "client_ip":
			record.IP = val
		case "trace_id", "traceid":
			record.TraceID = val
		}
	}

	return nil
}
