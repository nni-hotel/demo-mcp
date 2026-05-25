package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/nni-hotel/demo-mcp/internal/api"
	"github.com/nni-hotel/demo-mcp/internal/api/gen"
	"github.com/nni-hotel/demo-mcp/internal/config"
	"github.com/nni-hotel/demo-mcp/internal/discoverability"
	"github.com/nni-hotel/demo-mcp/internal/platform/middleware"
)

// NewRouter builds the full HTTP handler (discoverability + tool API routes).
func NewRouter(cfg config.Config) http.Handler {
	apiHandler := api.NewHandler(cfg.MaxBodyBytes)
	disc := discoverability.NewHandler(cfg.PublicBaseURL)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(middleware.MaxBytes(cfg.MaxBodyBytes))
	if cfg.ReadTimeout > 0 {
		r.Use(middleware.Timeout(cfg.ReadTimeout))
	}

	disc.RegisterRoutes(r)
	gen.HandlerFromMux(apiHandler, r)
	return r
}

// TestConfig returns config suitable for httptest servers.
func TestConfig() config.Config {
	return config.Config{
		APIAddr:       ":8080",
		PublicBaseURL: "http://test.example",
		MaxBodyBytes:  1 << 20,
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  5 * time.Second,
		ShutdownGrace: 5 * time.Second,
		LogLevel:      "error",
		Version:       "0.1.0",
		ServiceName:   "demo-mcp",
	}
}
