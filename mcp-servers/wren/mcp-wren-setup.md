# mcp-wren Setup

`mcp-wren` exposes wren task management as MCP tools over HTTP. The server is
deployed separately from the agent; agents connect to it via a URL.

## URL Format

```
http://<host>:<port>/mcp
```

The `<host>:<port>` is determined by whoever deployed the server. Consult
their deployment documentation.

## Agent Configuration

Add the server to your agent's MCP configuration:

```json
{
  "mcpServers": {
    "wren": {
      "url": "http://<host>:<port>/mcp"
    }
  }
}
```

## Security

The server has no built-in authentication. If you are connecting over a
network you do not control, use a reverse proxy that enforces authentication
(mTLS, basic auth, OAuth proxy, etc.) before the server.
