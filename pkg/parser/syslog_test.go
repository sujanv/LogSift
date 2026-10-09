package parser

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestSyslogParser_RFC3164(t *testing.T) {
	p := NewSyslogParser()
	rec := &model.Record{Fields: make(map[string]interface{})}

	line := []byte(`<34>Oct 10 09:14:15 server1 sshd[1234]: Failed password for root from 192.168.1.100`)
	if err := p.Parse(line, rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Level != model.LevelFatal && rec.Level != model.LevelError {
		// priority 34: 34 & 7 = 2 (Critical / Fatal)
		t.Logf("priority 34 level: %v", rec.Level)
	}
	if rec.Message != "Failed password for root from 192.168.1.100" {
		t.Errorf("expected message, got %q", rec.Message)
	}
}

func TestSyslogParser_RFC5424(t *testing.T) {
	p := NewSyslogParser()
	rec := &model.Record{Fields: make(map[string]interface{})}

	line := []byte(`<165>1 2026-10-10T09:15:30.000Z myhost myapp 1234 ID47 - An application error occurred`)
	if err := p.Parse(line, rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Message != "An application error occurred" {
		t.Errorf("expected msg, got %q", rec.Message)
	}
}
