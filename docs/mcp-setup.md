# MCP setup

Registry name: `io.github.nni-hotel/demo-mcp`

## stdio (Cursor / Claude Desktop)

Build or install `toolinfra`, then add to Cursor MCP settings:

```json
{
  "mcpServers": {
    "demo-mcp": {
      "command": "toolinfra",
      "args": ["mcp", "stdio"]
    }
  }
}
```

Windows (development):

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

## Streamable HTTP

Run the server with MCP HTTP enabled:

```bash
toolinfra serve --api --mcp-http
```

Default endpoint: `http://localhost:8081/mcp` (set `TOOLINFRA_MCP_HTTP_ADDR`, `TOOLINFRA_MCP_HTTP_PATH`).

## Tools

- `base64_encode` — encode UTF-8 text to Base64
- `base64_decode` — decode Base64 to UTF-8 text
