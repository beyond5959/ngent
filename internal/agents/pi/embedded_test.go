package pi

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultRuntimeConfigReadsPiEnv(t *testing.T) {
	t.Setenv("PI_BIN", "/tmp/custom-pi")
	t.Setenv("PI_PROVIDER", "openai-codex")
	t.Setenv("PI_MODEL", "gpt-5.4-mini")
	t.Setenv("PI_SESSION_DIR", "/tmp/pi-sessions")
	t.Setenv("PI_DISABLE_GATE", "true")

	cfg := DefaultRuntimeConfig()
	if cfg.PiBin != "/tmp/custom-pi" {
		t.Fatalf("PiBin = %q, want %q", cfg.PiBin, "/tmp/custom-pi")
	}
	if cfg.DefaultProvider != "openai-codex" {
		t.Fatalf("DefaultProvider = %q, want %q", cfg.DefaultProvider, "openai-codex")
	}
	if cfg.DefaultModel != "gpt-5.4-mini" {
		t.Fatalf("DefaultModel = %q, want %q", cfg.DefaultModel, "gpt-5.4-mini")
	}
	if cfg.SessionDir != "/tmp/pi-sessions" {
		t.Fatalf("SessionDir = %q, want %q", cfg.SessionDir, "/tmp/pi-sessions")
	}
	if cfg.EnableGate {
		t.Fatalf("EnableGate = %v, want false", cfg.EnableGate)
	}
}

func TestPreflightUsesConfiguredBinary(t *testing.T) {
	cfg := DefaultRuntimeConfig()
	cfg.PiBin = filepath.Join(t.TempDir(), "missing-pi")

	err := Preflight(cfg)
	if err == nil {
		t.Fatal("Preflight() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "missing-pi") {
		t.Fatalf("Preflight() error = %q, want missing binary path", err.Error())
	}
}
