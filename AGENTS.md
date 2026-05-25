# Agent instructions (read first)

Project: **demo-mcp** — `github.com/nni-hotel/demo-mcp`

## Rules

1. **API-first**: Change [api/openapi.yaml](api/openapi.yaml) first, then `make generate`, then handlers.
2. **Single logic source**: Business logic only in `internal/tools/*`. Do not duplicate encode/decode in REST or MCP.
3. **New tools**: Add OpenAPI path + MCP tool + `docs/tools/<name>.md` + `docs/semantic-tags.yaml` + catalog entry + `docs/mcp/tools.json`.
4. **GEO**: Keep OpenAPI `x-mcp-tool-name`, MCP tool names, and `docs/mcp/tools.json` in sync. Run `make verify-discoverability` before PRs.
5. **Scope**: Do not add Redis/Postgres/API keys/billing unless explicitly requested.

## Commands

```bash
make generate-discoverability
make generate   # oapi-codegen (if installed)
make check
make verify-discoverability
go test ./...
```

## Layout

- `cmd/toolinfra` — CLI (`serve`, `mcp stdio`)
- `internal/discoverability` — `/openapi.json`, `/v1/tools`, `llms.txt`
- `internal/mcp` — MCP tool registration
- `docs/` — human + machine docs for GEO

## MCP tools (v0.1)

- `base64_encode` / `base64_decode`
