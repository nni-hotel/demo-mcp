# Architecture

## Process model

Single binary `toolinfra`:

- `serve --api` — REST on `:8080` (Chi)
- `serve --mcp-http` — MCP Streamable HTTP on `:8081/mcp`
- `serve --api --mcp-http` — both (default in Docker Compose)
- `mcp stdio` — MCP over stdin/stdout for Cursor

## Layers

| Layer | Packages |
|-------|----------|
| Discoverability | `internal/discoverability` — OpenAPI, catalog, llms.txt, well-known |
| API | `internal/api`, `internal/api/gen` — OpenAPI-driven handlers |
| MCP | `internal/mcp` — official go-sdk tools |
| Tools | `internal/tools/*` — pure deterministic logic |
| Platform | `internal/platform/*` — errors, logging, middleware |

## Adding a new tool

1. Implement logic in `internal/tools/<name>/`
2. Add paths to `api/openapi.yaml` + `make generate`
3. Register MCP tools in `internal/mcp/server.go`
4. Update `docs/semantic-tags.yaml`, `docs/tools/<name>.md`
5. Update `internal/discoverability/catalog.go`
6. Run `make verify-discoverability`

## Ports

| Service | Default | Env |
|---------|---------|-----|
| REST | 8080 | `TOOLINFRA_API_ADDR` |
| MCP HTTP | 8081 | `TOOLINFRA_MCP_HTTP_ADDR` |
| MCP path | `/mcp` | `TOOLINFRA_MCP_HTTP_PATH` |
