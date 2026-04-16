# AGENTS.md — mcp-servers/wren

This deliverable contains a Go MCP server (`main.go`, `go.mod`, `go.sum`) and
a skill file (`SKILL.md`).

## Rules

- `SKILL.md` must be fully self-contained: no references outside `mcp-servers/wren/`.
- `main.go` invokes `wren` as a subprocess. All wren interaction uses stdlib
  (`os/exec`). Only `mark3labs/mcp-go` is a third-party dependency.
- Resolve the wren binary via `WREN_BIN` env var, falling back to `wren` on PATH.
- Non-zero wren exit codes must be surfaced as MCP tool errors, not panics.
- Do not add interactive or external-service wren operations.
- Task body creation is out of scope until wren gains a non-interactive body flag.
