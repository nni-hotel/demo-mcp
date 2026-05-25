# Base64 tools

Encode or decode **UTF-8 text** with deterministic output. Not for binary files.

## When to use

- User asks to encode text to Base64 (standard or URL-safe)
- User asks to decode a Base64 string to text
- Agent needs a reliable, stateless transform with predictable output

## When NOT to use

- Binary files, images, or encrypted payloads
- Password hashing or encryption (use dedicated crypto tools)
- Large payloads over 1 MiB (split input or use streaming tools)

## Tools

| MCP tool | REST | Deterministic |
|----------|------|---------------|
| `base64_encode` | `POST /v1/tools/base64/encode` | Yes |
| `base64_decode` | `POST /v1/tools/base64/decode` | Yes |

## Request

```json
{
  "input": "hello",
  "alphabet": "standard",
  "padding": true
}
```

- `alphabet`: `standard` (RFC 4648) or `url` (URL-safe)
- `padding`: include `=` padding (default `true`)

## Response

```json
{
  "data": { "output": "aGVsbG8=" },
  "meta": {
    "tool": "base64.encode",
    "input_bytes": 5,
    "output_bytes": 8,
    "duration_ms": 0.05
  }
}
```

## Examples

| Action | Input | Output |
|--------|-------|--------|
| Encode | `hello` | `aGVsbG8=` |
| Decode | `aGVsbG8=` | `hello` |
| URL encode (no padding) | `hello`, alphabet `url`, padding `false` | varies |

## Semantic tags

See [semantic-tags.yaml](../semantic-tags.yaml): `encoding`, `base64`, `text`, `utf8`.
