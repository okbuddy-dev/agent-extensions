# Releases

This repository contains two independent projects with separate release cycles.

## Tag naming

| Project | Tag pattern | Example |
|---|---|---|
| MCP server (`mcp-servers/wren`) | `mcp-wren/v<version>` | `mcp-wren/v0.3.0` |
| AI skill (`skills/wren`) | `skill-wren/v<version>` | `skill-wren/v1.2.0` |

## Versioning

- **mcp-wren**: Follows semver. Bump on Go code changes, dependency updates, bug fixes, new MCP tools.
- **skill-wren**: Follows semver. Bump on documentation changes, front-matter updates, behavior clarifications.

## Making a release

### MCP server

```bash
cd mcp-servers/wren
git tag mcp-wren/v0.3.0
git push origin mcp-wren/v0.3.0
```

### AI skill

```bash
cd skills/wren
git tag skill-wren/v1.2.0
git push origin skill-wren/v1.2.0
```

## Release artifacts

| Project | GitHub Release | Contents |
|---|---|---|
| `mcp-wren/v*` | wren-mcp v\<version\> | Multi-platform Go binaries + checksums |
| `skill-wren/v*` | wren skill v\<version\> | SKILL.md only |