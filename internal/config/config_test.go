package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("TOOLINFRA_API_ADDR", "")
	t.Setenv("TOOLINFRA_MAX_BODY_BYTES", "")
	cfg := Load()

	if cfg.APIAddr != DefaultAPIAddr {
		t.Fatalf("APIAddr: got %q", cfg.APIAddr)
	}
	if cfg.MaxBodyBytes != DefaultMaxBodyBytes {
		t.Fatalf("MaxBodyBytes: got %d", cfg.MaxBodyBytes)
	}
	if cfg.MCPHTTPPath != DefaultMCPHTTPPath {
		t.Fatalf("MCPHTTPPath: got %q", cfg.MCPHTTPPath)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("TOOLINFRA_API_ADDR", ":9090")
	t.Setenv("TOOLINFRA_MAX_BODY_BYTES", "2048")
	t.Setenv("TOOLINFRA_PUBLIC_BASE_URL", "https://api.example.com")
	t.Setenv("TOOLINFRA_READ_TIMEOUT", "3s")

	cfg := Load()

	if cfg.APIAddr != ":9090" {
		t.Fatalf("APIAddr: got %q", cfg.APIAddr)
	}
	if cfg.MaxBodyBytes != 2048 {
		t.Fatalf("MaxBodyBytes: got %d", cfg.MaxBodyBytes)
	}
	if cfg.PublicBaseURL != "https://api.example.com" {
		t.Fatalf("PublicBaseURL: got %q", cfg.PublicBaseURL)
	}
	if cfg.ReadTimeout != 3*time.Second {
		t.Fatalf("ReadTimeout: got %v", cfg.ReadTimeout)
	}
}
