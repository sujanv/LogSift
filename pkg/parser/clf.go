package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

// CLFParser parses Apache and Nginx access logs in Common or Combined format.
type CLFParser struct{}

// NewCLFParser creates a new CLFParser.
func NewCLFParser() *CLFParser {
	return &CLFParser{}
}

// Name returns format name.
func (p *CLFParser) Name() string {
	return "clf"
}

// Parse extracts fields from Apache/Nginx access log lines into record.
func (p *CLFParser) Parse(line []byte, record *model.Record) error {
	record.Raw = line
	if record.Fields == nil {
		record.Fields = make(map[string]interface{}, 16)
	}

	s := string(line)
	// Example: 127.0.0.1 - frank [10/Oct/2026:13:55:36 -0700] "GET /api/v1/users HTTP/1.1" 200 2326 "https://ref.com" "Mozilla/5.0"
	ipEnd := strings.IndexByte(s, ' ')
	if ipEnd == -1 {
		return fmt.Errorf("invalid clf: missing ip")
	}
	record.IP = s[:ipEnd]
	record.Fields["ip"] = record.IP

	// Find timestamp [date]
	timeStart := strings.IndexByte(s, '[')
	timeEnd := strings.IndexByte(s, ']')
	if timeStart == -1 || timeEnd == -1 || timeEnd <= timeStart {
		return fmt.Errorf("invalid clf: missing timestamp brackets")
	}
	timeStr := s[timeStart+1 : timeEnd]
	if t, err := time.Parse("02/Jan/2006:15:04:05 -0700", timeStr); err == nil {
		record.Timestamp = t
	}

	// Find request "METHOD /path HTTP/x.x"
	reqStart := strings.Index(s[timeEnd:], "\"")
	if reqStart == -1 {
		return fmt.Errorf("invalid clf: missing request quote")
	}
	reqStart += timeEnd + 1
	reqEnd := strings.Index(s[reqStart:], "\"")
	if reqEnd == -1 {
		return fmt.Errorf("invalid clf: unclosed request quote")
	}
	reqEnd += reqStart
	reqLine := s[reqStart:reqEnd]
	record.Message = reqLine
	record.Fields["request"] = reqLine

	parts := strings.Fields(reqLine)
	if len(parts) >= 1 {
		record.Method = parts[0]
	}
	if len(parts) >= 2 {
		record.Path = parts[1]
	}

	// Find status code and response size after request
	rest := strings.TrimSpace(s[reqEnd+1:])
	restParts := strings.Fields(rest)
	if len(restParts) >= 1 {
		if status, err := strconv.Atoi(restParts[0]); err == nil {
			record.StatusCode = status
			record.Fields["status"] = status
			switch {
			case status >= 500:
				record.Level = model.LevelError
			case status >= 400:
				record.Level = model.LevelWarn
			default:
				record.Level = model.LevelInfo
			}
		}
	}
	if len(restParts) >= 2 {
		if bytesSize, err := strconv.Atoi(restParts[1]); err == nil {
			record.Fields["bytes_sent"] = bytesSize
		}
	}

	return nil
}
