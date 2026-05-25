package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/nni-hotel/demo-mcp/internal/config"
	"github.com/nni-hotel/demo-mcp/internal/platform/logging"
	"github.com/nni-hotel/demo-mcp/internal/server"
)

var version = "0.1.0"

func main() {
	cfg := config.Load()
	log := logging.New(cfg.LogLevel)

	root := &cobra.Command{
		Use:     "toolinfra",
		Short:   "AI Tool Infrastructure — deterministic utilities for agents",
		Version: version,
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP API and/or MCP Streamable HTTP servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, _ := cmd.Flags().GetBool("api")
			mcpHTTP, _ := cmd.Flags().GetBool("mcp-http")
			if !api && !mcpHTTP {
				return fmt.Errorf("enable at least one of --api or --mcp-http")
			}

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			var httpSrv *server.HTTPServer
			var mcpSrv *server.MCPHTTPServer

			if api {
				httpSrv = server.NewHTTPServer(cfg, log)
				if err := httpSrv.Start(); err != nil {
					return err
				}
			}
			if mcpHTTP {
				mcpSrv = server.NewMCPHTTPServer(cfg, log)
				if err := mcpSrv.Start(); err != nil {
					return err
				}
			}

			log.Info("toolinfra ready", "version", version)
			<-ctx.Done()
			log.Info("shutting down")

			shutdownCtx := context.Background()
			if httpSrv != nil {
				_ = httpSrv.Shutdown(shutdownCtx)
			}
			if mcpSrv != nil {
				_ = mcpSrv.Shutdown(shutdownCtx)
			}
			return nil
		},
	}
	serveCmd.Flags().Bool("api", false, "Start REST API server")
	serveCmd.Flags().Bool("mcp-http", false, "Start MCP Streamable HTTP server")

	mcpCmd := &cobra.Command{
		Use:   "mcp",
		Short: "MCP transport commands",
	}
	mcpStdio := &cobra.Command{
		Use:   "stdio",
		Short: "Run MCP server over stdin/stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.RunMCPStdio(cfg, log)
		},
	}
	mcpCmd.AddCommand(mcpStdio)

	root.AddCommand(serveCmd, mcpCmd)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
