# Error codes (machine-readable)

All errors use this envelope:

```json
{
  "error": { "code": "INVALID_BASE64", "message": "...", "details": {} },
  "meta": { "request_id": "..." }
}
```

| Code | HTTP | Agent retry | When |
|------|------|-------------|------|
| `INVALID_REQUEST` | 400 | No | Malformed JSON or invalid parameters (e.g. unknown alphabet) |
| `INVALID_BASE64` | 400 | No | Decode input is not valid Base64; ask user to fix input |
| `PAYLOAD_TOO_LARGE` | 413 | No | Body exceeds `TOOLINFRA_MAX_BODY_BYTES` (default 1 MiB) |
| `INTERNAL_ERROR` | 500 | Once | Unexpected server failure; backoff and retry at most once |

OpenAPI enum: `components.schemas.ErrorCode` in [api/openapi.yaml](../api/openapi.yaml).
