package tracectx

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type TraceContext struct {
	Version     string            `json:"version"`
	TraceID     string            `json:"trace_id"`
	SpanID      string            `json:"span_id"`
	TraceFlags  string            `json:"trace_flags"`
	Sampled     bool              `json:"sampled"`
	TraceState  map[string]string `json:"trace_state,omitempty"`
}

func ParseTraceParent(header string) (*TraceContext, error) {
	trimmed := strings.TrimSpace(header)
	parts := strings.Split(trimmed, "-")
	if len(parts) != 4 {
		return nil, errors.New("invalid traceparent format: expected 4 hyphen-delimited fields")
	}

	version := parts[0]
	traceID := parts[1]
	spanID := parts[2]
	flags := parts[3]

	if len(version) != 2 || len(traceID) != 32 || len(spanID) != 16 || len(flags) != 2 {
		return nil, errors.New("invalid field lengths in traceparent")
	}

	if traceID == strings.Repeat("0", 32) {
		return nil, errors.New("trace_id cannot be all zeros")
	}
	if spanID == strings.Repeat("0", 16) {
		return nil, errors.New("span_id cannot be all zeros")
	}

	if _, err := hex.DecodeString(version); err != nil {
		return nil, fmt.Errorf("invalid hex in version: %w", err)
	}
	if _, err := hex.DecodeString(traceID); err != nil {
		return nil, fmt.Errorf("invalid hex in trace_id: %w", err)
	}
	if _, err := hex.DecodeString(spanID); err != nil {
		return nil, fmt.Errorf("invalid hex in span_id: %w", err)
	}
	flagBytes, err := hex.DecodeString(flags)
	if err != nil {
		return nil, fmt.Errorf("invalid hex in trace_flags: %w", err)
	}

	sampled := (flagBytes[0] & 0x01) == 0x01

	return &TraceContext{
		Version:    version,
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: flags,
		Sampled:    sampled,
		TraceState: make(map[string]string),
	}, nil
}

func (tc *TraceContext) String() string {
	return fmt.Sprintf("%s-%s-%s-%s", tc.Version, tc.TraceID, tc.SpanID, tc.TraceFlags)
}

func ParseTraceState(header string) map[string]string {
	state := make(map[string]string)
	if strings.TrimSpace(header) == "" {
		return state
	}

	pairs := strings.Split(header, ",")
	for _, p := range pairs {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) == 2 {
			state[kv[0]] = kv[1]
		}
	}
	return state
}
