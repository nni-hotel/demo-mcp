# demo-mcp

> Deterministic Base64 tools for LLMs — REST + MCP. UTF-8 text only; same input always yields same output.

[![MCP Registry](https://img.shields.io/badge/MCP-io.github.nni--hotel%2Fdemo--mcp-blue)](https://registry.modelcontextprotocol.io/)

## Quick start

```bash
# Docker (REST :8080 + MCP HTTP :8081)
docker compose -f deployments/docker-compose.yml up --build

# Local
go run ./cmd/toolinfra serve --api --mcp-http

# Encode via REST
curl -s -X POST http://localhost:8080/v1/tools/base64/encode \
  -H "Content-Type: application/json" \
  -d '{"input":"hello"}'
```

## Tools

| REST | MCP tool | Deterministic |
|------|----------|---------------|
| `POST /v1/tools/base64/encode` | `base64_encode` | Yes |
| `POST /v1/tools/base64/decode` | `base64_decode` | Yes |

## When to use / not use

**Use** for encoding or decoding UTF-8 text with Base64 (standard or URL-safe).

**Do not use** for binary files, images, or encryption.

## Cursor MCP (stdio)

```json
{
  "mcpServers": {
    "demo-mcp": {
      "command": "go",
      "args": ["run", "./cmd/toolinfra", "mcp", "stdio"],
      "cwd": "C:\\path\\to\\demo-mcp"
    }
  }
}
```

## Discoverability (GEO)

| Resource | URL |
|----------|-----|
| OpenAPI | [/openapi.json](http://localhost:8080/openapi.json) |
| Tool catalog | [/v1/tools](http://localhost:8080/v1/tools) |
| API catalog | [/.well-known/api-catalog](http://localhost:8080/.well-known/api-catalog) |
| llms.txt | [/llms.txt](http://localhost:8080/llms.txt) |
| Spec (repo) | [api/openapi.yaml](api/openapi.yaml) |

## Development

```bash
make check                  # fmt, unit + integration tests, validate spec
make test-unit              # package-level unit tests
make test-integration       # HTTP API integration tests (httptest)
make verify-discoverability # OpenAPI/MCP/docs consistency
make generate-discoverability
```

### Tests

| Layer | Location | Command |
|-------|----------|---------|
| Unit | `internal/*`, `evals/*` | `make test-unit` |
| Integration | `tests/integration` | `make test-integration` |

See [AGENTS.md](AGENTS.md) for agent constraints, [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md), [docs/VISION.md](docs/VISION.md).

## License

MIT
