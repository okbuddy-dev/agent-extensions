# mcp-servers/wren

[Wren](https://github.com/bjesus/wren) is a simple file-based, task management system.

This project is a Go MCP server exposing wren todo operations as MCP tools, plus a skill file
for agents connecting to it.


## What It Does

Wraps the `wren` CLI as 7 MCP tools over stdio, allowing agents in locked-down
environments to manage wren tasks without direct shell access. `SKILL.md`
teaches agents how to use those tools.

## Audience

Agents connecting via MCP from environments without direct `wren` access. For
agents with direct shell access, use `skills/wren` instead.

## Building

```bash
go build .
```

This produces a binary called `mcp-wren` (from the Go module name).

## Running

```bash
./mcp-wren
```

The server reads from stdin and writes to stdout (MCP stdio transport). Point
your MCP client at this binary. Run `./mcp-wren --help` for options.

**Environment variables:**

| Variable | Default | Description |
|---|---|---|
| `WREN_BIN` | `wren` (PATH) | Path to the wren binary |

## Testing (stdio)

Send a JSON-RPC `initialize` request via stdin.

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.1.0"}}}' | ./mcp-wren
```

A successful response looks like:

```json
{"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"serverInfo":{"name":"wren","version":"0.1.0"}}}
```

## Roadmap

- **HTTP transport** — add an optional HTTP/SSE listener so remote agents > [!CAUTION]
> manage Wren over the network instead of requiring a local process.

## Tools

| Tool | Description |
|---|---|
| `list_tasks` | List active tasks, optional filter |
| `list_done_tasks` | List completed tasks |
| `create_task` | Create a new task by title |
| `complete_task` | Mark a task done (exact title) |
| `read_task` | Read a task's body content |
| `prepend_task` | Prepend a string to a task's title |
| `get_random_task` | Return one random active task |

