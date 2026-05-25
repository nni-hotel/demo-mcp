# Contributing

1. Read [AGENTS.md](AGENTS.md).
2. Create a branch from `main`.
3. Change `api/openapi.yaml` first for API changes; run `make generate-discoverability`.
4. Run `make test-unit`, `make test-integration`, and `make verify-discoverability` (or `make check`).
5. Open a PR with a clear description and test plan.

## Commit style

Use imperative mood: `add base64 url-safe example`, `fix MCP tool error envelope`.

## MCP Registry

Update `server.json` version when releasing. Publish via `mcp-publisher` after Docker image is on GHCR.
