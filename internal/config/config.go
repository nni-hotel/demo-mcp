package config

import (
	"os"
	"strconv"
	"time"
)

const (
	DefaultAPIAddr       = ":8080"
	DefaultMCPHTTPAddr   = ":8081"
	DefaultMCPHTTPPath   = "/mcp"
	DefaultMaxBodyBytes  = 1 << 20 // 1 MiB
	DefaultReadTimeout   = 10 * time.Second
	DefaultWriteTimeout  = 10 * time.Second
	DefaultShutdownGrace = 15 * time.Second
	DefaultLogLevel      = "info"
	DefaultVersion       = "0.1.0"
)

type Config struct {
	APIAddr        string
	PublicBaseURL  string
	MCPHTTPAddr    string
	MCPHTTPPath    string
	MaxBodyBytes   int64
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ShutdownGrace  time.Duration
	LogLevel       string
	Version        string
	ServiceName    string
}

func Load() Config {
	return Config{
		APIAddr:       envOr("TOOLINFRA_API_ADDR", DefaultAPIAddr),
		PublicBaseURL: envOr("TOOLINFRA_PUBLIC_BASE_URL", ""),
		MCPHTTPAddr:   envOr("TOOLINFRA_MCP_HTTP_ADDR", DefaultMCPHTTPAddr),
		MCPHTTPPath:   envOr("TOOLINFRA_MCP_HTTP_PATH", DefaultMCPHTTPPath),
		MaxBodyBytes:  envInt64Or("TOOLINFRA_MAX_BODY_BYTES", DefaultMaxBodyBytes),
		ReadTimeout:   envDurationOr("TOOLINFRA_READ_TIMEOUT", DefaultReadTimeout),
		WriteTimeout:  envDurationOr("TOOLINFRA_WRITE_TIMEOUT", DefaultWriteTimeout),
		ShutdownGrace: envDurationOr("TOOLINFRA_SHUTDOWN_GRACE", DefaultShutdownGrace),
		LogLevel:      envOr("TOOLINFRA_LOG_LEVEL", DefaultLogLevel),
		Version:       envOr("TOOLINFRA_VERSION", DefaultVersion),
		ServiceName:   envOr("TOOLINFRA_SERVICE_NAME", "demo-mcp"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt64Or(key string, fallback int64) int64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return fallback
}

func envDurationOr(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
