# Product vision (architecture planning)

> Historical planning document. For running the project see [README.md](../README.md).

## Background

AI Tool Infrastructure provides **deterministic utility tools** for LLM, Agent, and MCP clients — not a traditional online tool website.

First tools: Base64 encode/decode. Future: timestamp, JSON format, regex, hash, UUID, URL encode, unit/timezone convert, etc.

## Goals

- API-first, MCP-ready
- Low latency, memory, CPU; high concurrency
- AI-friendly schema, deterministic output
- Observability, rate limiting (later), marketplace (later)

## Roadmap

| Version | Scope |
|---------|--------|
| v0.1 | Base64 REST + MCP + OpenAPI + Docker |
| v0.2 | Discoverability (Registry, llms.txt, catalog, evals) |
| v1.0 | API Key, rate limit, metrics, multi-tool SDK |
| v2.0 | Agent gateway, marketplace |

See [ARCHITECTURE.md](ARCHITECTURE.md) and the GEO discoverability plan for v0.2+ details.
