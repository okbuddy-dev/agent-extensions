# agent-extensions

`agent-extensions` is a monorepo for small agent-extension deliverables that do not justify separate repos.

## Organization

- `skills/` for skills
- `mcp-servers/` for MCP servers
- `plugins/` for plugins
- `hooks/` for hooks

Each deliverable lives in a second-level folder named for the deliverable, for example:

- `skills/wren`
- `mcp-servers/wren`
- `skills/showerthoughts`

Same-named folders across extension types are not automatically related.

## Baseline per Deliverable

Every deliverable starts with:

- `README.md`
- `AGENTS.md`

Additional files are added only when that extension type clearly needs them.
