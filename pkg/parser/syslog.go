package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

// SyslogParser parses RFC 3164 and RFC 5424 syslog messages.
type SyslogParser struct{}

// NewSyslogParser creates a new SyslogParser.
func NewSyslogParser() *SyslogParser {
	return &SyslogParser{}
}

// Name returns format name.
func (p *SyslogParser) Name() string {
	return "syslog"
}

// Parse extracts fields from a syslog message line into record.
func (p *SyslogParser) Parse(line []byte, record *model.Record) error {
	record.Raw = line
	if record.Fields == nil {
		record.Fields = make(map[string]interface{}, 16)
	}

	s := strings.TrimSpace(string(line))
	if len(s) == 0 {
		return fmt.Errorf("empty line")
	}

	// Priority <PRI>
	if s[0] == '<' {
		priEnd := strings.IndexByte(s, '>')
		if priEnd > 1 {
			if pri, err := strconv.Atoi(s[1:priEnd]); err == nil {
				severity := pri & 7
				record.Fields["priority"] = pri
				record.Fields["severity_num"] = severity
				switch severity {
				case 0, 1, 2:
					record.Level = model.LevelFatal
				case 3:
					record.Level = model.LevelError
				case 4:
					record.Level = model.LevelWarn
				case 5, 6:
					record.Level = model.LevelInfo
				case 7:
					record.Level = model.LevelDebug
				}
			}
			s = s[priEnd+1:]
		}
	}

	// RFC 5424 check: <PRI>1 TIMESTAMP HOST APP-NAME PROCID MSGID STRUCTURED-DATA MSG
	if len(s) > 2 && s[0] == '1' && s[1] == ' ' {
		parts := strings.SplitN(s[2:], " ", 7)
		if len(parts) >= 1 {
			if t, err := time.Parse(time.RFC3339Nano, parts[0]); err == nil {
				record.Timestamp = t
			}
		}
		if len(parts) >= 7 {
			record.Message = parts[6]
		} else if len(parts) >= 6 {
			record.Message = parts[5]
		}
		return nil
	}

	// RFC 3164: Mmm dd hh:mm:ss hostname tag: msg
	if len(s) >= 15 {
		monthStr := s[:3]
		validMonths := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
		isRFC3164 := false
		for _, m := range validMonths {
			if monthStr == m {
				isRFC3164 = true
				break
			}
		}
		if isRFC3164 {
			timeStr := s[:15]
			now := time.Now()
			fullTimeStr := fmt.Sprintf("%d %s", now.Year(), timeStr)
			if t, err := time.Parse("2006 Jan _2 15:04:05", fullTimeStr); err == nil {
				record.Timestamp = t
			}
			msgStart := strings.Index(s[15:], ": ")
			if msgStart != -1 {
				record.Message = s[15+msgStart+2:]
			} else {
				record.Message = strings.TrimSpace(s[15:])
			}
			return nil
		}
	}

	record.Message = s
	return nil
}
