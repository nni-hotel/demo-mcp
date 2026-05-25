# Agent tool selection evals

`prompts.jsonl` lists natural-language prompts and the expected MCP tool name.

## Run (keyword baseline)

```bash
go test ./evals/agent_tool_selection/...
```

This uses simple intent matching from `docs/semantic-tags.yaml` as a baseline. Replace with LLM-judged evals in CI when ready.

## KPI target

Tool selection@1 > 90% on this set.
