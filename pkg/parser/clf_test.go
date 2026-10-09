package parser

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestCLFParser_Combined(t *testing.T) {
	p := NewCLFParser()
	rec := &model.Record{Fields: make(map[string]interface{})}

	line := []byte(`192.168.1.50 - admin [10/Oct/2026:13:55:36 -0700] "POST /api/v1/auth/login HTTP/1.1" 500 1024 "https://example.com" "Mozilla/5.0"`)
	if err := p.Parse(line, rec); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if rec.IP != "192.168.1.50" {
		t.Errorf("expected IP 192.168.1.50, got %q", rec.IP)
	}
	if rec.Method != "POST" {
		t.Errorf("expected Method POST, got %q", rec.Method)
	}
	if rec.Path != "/api/v1/auth/login" {
		t.Errorf("expected Path /api/v1/auth/login, got %q", rec.Path)
	}
	if rec.StatusCode != 500 {
		t.Errorf("expected StatusCode 500, got %d", rec.StatusCode)
	}
	if rec.Level != model.LevelError {
		t.Errorf("expected LevelError for 500 status, got %v", rec.Level)
	}
}
