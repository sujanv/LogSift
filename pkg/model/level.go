package model

import (
	"bytes"
	"strings"
)

// LogLevel represents the normalized severity level of a log entry.
type LogLevel int

const (
	LevelUnknown LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// String returns the uppercase string representation of the level.
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Severity returns a numeric severity index useful for comparison.
func (l LogLevel) Severity() int {
	return int(l)
}

// ParseLevel converts a string into a normalized LogLevel.
func ParseLevel(s string) LogLevel {
	s = strings.TrimSpace(s)
	switch strings.ToLower(s) {
	case "debug", "dbg", "trace", "trc":
		return LevelDebug
	case "info", "inf", "notice":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error", "err", "eror":
		return LevelError
	case "fatal", "crit", "critical", "panic", "emerg", "emergency", "alert":
		return LevelFatal
	default:
		return LevelUnknown
	}
}

// ParseLevelBytes converts a byte slice into a normalized LogLevel without string allocations.
func ParseLevelBytes(b []byte) LogLevel {
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return LevelUnknown
	}
	switch {
	case bytes.EqualFold(b, []byte("debug")) || bytes.EqualFold(b, []byte("dbg")) || bytes.EqualFold(b, []byte("trace")):
		return LevelDebug
	case bytes.EqualFold(b, []byte("info")) || bytes.EqualFold(b, []byte("inf")):
		return LevelInfo
	case bytes.EqualFold(b, []byte("warn")) || bytes.EqualFold(b, []byte("warning")):
		return LevelWarn
	case bytes.EqualFold(b, []byte("error")) || bytes.EqualFold(b, []byte("err")):
		return LevelError
	case bytes.EqualFold(b, []byte("fatal")) || bytes.EqualFold(b, []byte("crit")) || bytes.EqualFold(b, []byte("panic")):
		return LevelFatal
	default:
		return LevelUnknown
	}
}
