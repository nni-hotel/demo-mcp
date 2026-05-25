package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nni-hotel/demo-mcp/internal/config"
	mcppkg "github.com/nni-hotel/demo-mcp/internal/mcp"
)

type MCPHTTPServer struct {
	cfg    config.Config
	log    *slog.Logger
	server *http.Server
	mcp    *mcp.Server
}

func NewMCPHTTPServer(cfg config.Config, log *slog.Logger) *MCPHTTPServer {
	return &MCPHTTPServer{cfg: cfg, log: log, mcp: mcppkg.NewServer(cfg)}
}

func (s *MCPHTTPServer) Start() error {
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return s.mcp
	}, &mcp.StreamableHTTPOptions{JSONResponse: true})

	mux := http.NewServeMux()
	path := s.cfg.MCPHTTPPath
	if path == "" {
		path = "/mcp"
	}
	mux.Handle(path, handler)

	s.server = &http.Server{
		Addr:    s.cfg.MCPHTTPAddr,
		Handler: mux,
	}
	s.log.Info("starting MCP Streamable HTTP", "addr", s.cfg.MCPHTTPAddr, "path", path)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("MCP HTTP failed", "error", err)
		}
	}()
	return nil
}

func (s *MCPHTTPServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

func RunMCPStdio(cfg config.Config, log *slog.Logger) error {
	srv := mcppkg.NewServer(cfg)
	log.Info("starting MCP stdio")
	return srv.Run(context.Background(), &mcp.StdioTransport{})
}
