package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/nni-hotel/demo-mcp/internal/config"
)

type HTTPServer struct {
	cfg    config.Config
	log    *slog.Logger
	server *http.Server
}

func NewHTTPServer(cfg config.Config, log *slog.Logger) *HTTPServer {
	return &HTTPServer{cfg: cfg, log: log}
}

func (s *HTTPServer) Start() error {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Mount("/", NewRouter(s.cfg))

	s.server = &http.Server{
		Addr:         s.cfg.APIAddr,
		Handler:      r,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}
	s.log.Info("starting REST API", "addr", s.cfg.APIAddr)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("REST API failed", "error", err)
		}
	}()
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownGrace)
	defer cancel()
	s.log.Info("shutting down REST API")
	return s.server.Shutdown(ctx)
}
