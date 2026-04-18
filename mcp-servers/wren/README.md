# mcp-servers/wren

[Wren](https://github.com/bjesus/wren) is a simple file-based, task management system.

This project is a Go MCP server exposing wren todo operations as MCP tools, plus a skill file
for agents connecting to it.


## What It Does

Wraps the `wren` CLI as 7 MCP tools over stdio or HTTP, allowing agents in
locked-down environments to manage wren tasks without direct shell access.
`SKILL.md` teaches agents how to use those tools.

## Audience

Agents connecting via MCP from environments without direct `wren` access. For
agents with direct shell access, use `skills/wren` instead.

## Building

```bash
go build .
```

This produces a binary called `mcp-wren` (from the Go module name).

## Running

The server supports two transports: stdio (default) and HTTP.

### stdio (default)

```bash
./mcp-wren
```

The server reads from stdin and writes to stdout (MCP stdio transport). Point
your MCP client at this binary.

### HTTP

```bash
TRANSPORT=http ./mcp-wren
```

The server listens on `MCP_WREN_ADDR` (default `:8080`) at the `/mcp` endpoint.

**Flags:**

| Flag | Description |
|---|---|
| `--addr <address>` | Listen address (overrides `MCP_WREN_ADDR`, e.g., `127.0.0.1:8080`) |

**Environment variables:**

| Variable | Default | Description |
|---|---|---|
| `WREN_BIN` | `wren` (PATH) | Path to the wren binary |
| `TRANSPORT` | `stdio` | Transport mode: `stdio` or `http` |
| `MCP_WREN_ADDR` | `:8080` | HTTP listen address |

**Priority:** `--addr` flag > `MCP_WREN_ADDR` env var > `:8080` default

Run `./mcp-wren --help` for full options.

## Testing

### stdio

Send a JSON-RPC `initialize` request via stdin.

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.1.0"}}}' | ./mcp-wren
```

A successful response looks like:

```json
{"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"serverInfo":{"name":"wren","version":"0.2.0"}}}
```

### HTTP

Start the server with `TRANSPORT=http`, then send an initialize request:

```bash
TRANSPORT=http ./mcp-wren &
sleep 1
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.1.0"}}}'
```

A successful response looks like:

```json
{"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"serverInfo":{"name":"wren","version":"0.2.0"}}}
```

## Security

When deploying with HTTP transport, consider:

- **No authentication** — Anyone who can reach the endpoint can perform all wren
  operations. Use a reverse proxy with authentication for production deployments.
- **No TLS** — Plain HTTP. Terminate TLS at a reverse proxy for production.
- **Default binds all interfaces** — Use `--addr 127.0.0.1:8080` for localhost-only,
  or restrict access via firewall/reverse proxy.
- **No rate limiting** — Apply at the reverse proxy layer if needed.
- **Sensitive task data** — Network-level access controls are the operator's
  responsibility if tasks contain sensitive information.

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

