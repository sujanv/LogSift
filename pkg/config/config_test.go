package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.DefaultFormat != "auto" {
		t.Errorf("expected default format 'auto', got %s", cfg.DefaultFormat)
	}
	if !cfg.Sampling.AlwaysErrors {
		t.Errorf("expected AlwaysErrors to be true by default")
	}
}

func TestLoadFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "logsift.json")

	data := []byte(`{"default_format":"json","buffer_size":131072,"masking":{"enabled":true,"strategy":"hash"}}`)
	if err := os.WriteFile(cfgFile, data, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	loaded, err := LoadFromFile(cfgFile)
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}
	if loaded.DefaultFormat != "json" {
		t.Errorf("expected format 'json', got %s", loaded.DefaultFormat)
	}
	if loaded.BufferSize != 131072 {
		t.Errorf("expected buffer size 131072, got %d", loaded.BufferSize)
	}
	if !loaded.Masking.Enabled || loaded.Masking.Strategy != "hash" {
		t.Errorf("masking config mismatch")
	}
}
