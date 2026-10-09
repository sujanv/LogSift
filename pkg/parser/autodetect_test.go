package parser

import (
	"testing"
)

func TestAutoDetect(t *testing.T) {
	jsonSamples := [][]byte{
		[]byte(`{"time":"2026-10-10T10:00:00Z","level":"info","msg":"server started"}`),
		[]byte(`{"time":"2026-10-10T10:00:01Z","level":"error","msg":"connection refused"}`),
	}
	p := AutoDetect(jsonSamples)
	if p.Name() != "json" {
		t.Errorf("expected json parser, got %s", p.Name())
	}

	clfSamples := [][]byte{
		[]byte(`127.0.0.1 - - [10/Oct/2026:10:00:00 +0000] "GET /health HTTP/1.1" 200 45`),
		[]byte(`127.0.0.1 - - [10/Oct/2026:10:00:01 +0000] "POST /api/login HTTP/1.1" 401 128`),
	}
	p2 := AutoDetect(clfSamples)
	if p2.Name() != "clf" {
		t.Errorf("expected clf parser, got %s", p2.Name())
	}

	logfmtSamples := [][]byte{
		[]byte(`ts=2026-10-10T10:00:00Z level=warn msg="rate limit approaching" client=123`),
		[]byte(`ts=2026-10-10T10:00:01Z level=info msg="job finished" duration=45ms`),
	}
	p3 := AutoDetect(logfmtSamples)
	if p3.Name() != "logfmt" {
		t.Errorf("expected logfmt parser, got %s", p3.Name())
	}
}
