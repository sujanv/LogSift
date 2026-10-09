package model

import (
	"time"
)

// Record represents a parsed, normalized structured log entry.
type Record struct {
	Timestamp  time.Time              `json:"timestamp"`
	Level      LogLevel               `json:"level"`
	Message    string                 `json:"message"`
	Latency    time.Duration          `json:"latency"`
	StatusCode int                    `json:"status_code,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Path       string                 `json:"path,omitempty"`
	IP         string                 `json:"ip,omitempty"`
	TraceID    string                 `json:"trace_id,omitempty"`
	SpanID     string                 `json:"span_id,omitempty"`
	UserID     string                 `json:"user_id,omitempty"`
	LineNum    int64                  `json:"line_num"`
	Offset     int64                  `json:"offset"`
	Raw        []byte                 `json:"-"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
}

// Reset clears all fields in the Record so it can be safely reused from a pool.
func (r *Record) Reset() {
	r.Timestamp = time.Time{}
	r.Level = LevelUnknown
	r.Message = ""
	r.Latency = 0
	r.StatusCode = 0
	r.Method = ""
	r.Path = ""
	r.IP = ""
	r.TraceID = ""
	r.SpanID = ""
	r.UserID = ""
	r.LineNum = 0
	r.Offset = 0
	r.Raw = nil
	if r.Fields != nil {
		for k := range r.Fields {
			delete(r.Fields, k)
		}
	}
}
