# Wren Extensions Design

**Date:** 2026-04-16
**Deliverables:** `skills/wren`, `mcp-servers/wren`

---

## Overview

Two independent extensions for the `wren` todo CLI — one for agents with direct
shell access, one for agents in locked-down environments accessing wren via MCP.
Users are expected to install one or the other, not both.

Both extensions target vanilla `wren` (no wrapper scripts). The `ww` wrapper and
its extended functionality are explicitly out of scope.

---

## Repository Structure

```
docs/
  wren.md                        ← developer reference: wren data model,
                                    config, filename conventions, quirks
                                    (not deployed; used to author both skills)

skills/wren/
  SKILL.md                       ← fully self-contained skill for CLI agents
  README.md
  AGENTS.md

mcp-servers/wren/
  SKILL.md                       ← fully self-contained skill for MCP agents
  README.md
  AGENTS.md
  main.go                        ← MCP server implementation
  go.mod
  go.sum
```

Each second-level directory is fully self-contained. `mcp-servers/wren` contains
no references to files outside its own directory. `docs/wren.md` is a developer
artifact only — it is the upstream source used to derive both `SKILL.md` files
and keep them consistent, but it is never referenced by deployed agents.

---

## Developer Reference: `docs/wren.md`

Documents the following for use during authoring:

- Executable location and config path (`~/.config/wren/wren.json`)
- `notes_dir` and `done_dir` config keys
- Data model: tasks are plain files; filename = title; file content = body
- Filename conventions: cron prefix (5-field), date prefix (`YYYY-MM-DD`)
- Done directory: relative to `notes_dir`, default `done/`
- Behavior of each supported operation and any quirks
- Operations excluded and why

---

## `skills/wren` — CLI Skill

**Audience:** Agents with direct shell access to the `wren` binary and notes
directory.

**Content:**

- Executable: `wren` (PATH), fallback `/home/vagrant/.local/bin/wren`
- Config location and keys
- Task data model (filename = title, file content = body)
- Filename conventions: cron and date prefixes
- All supported operations with exact command syntax:
  - List active tasks: `wren -l [filter]`
  - List done tasks: `wren -d`
  - Create task: `wren "title"`
  - Complete task: `wren -d "title"` (exact name preferred to avoid
    interactive disambiguation)
  - Read task: `wren -r "keyword"`
  - Prepend to task: `wren --prepend "prefix" "task"`
  - Random task: `wren -o`
- Excluded operations: `-e` (interactive editor), `--summary` (requires
  `litellm`), `--telegram`/`--matrix`/`--http` (require external service setup)

The skill is authored from `docs/wren.md` but contains all necessary content
inline — no external file references.

---

## `mcp-servers/wren` — MCP Server + Skill

### Go Server

**Language:** Go  
**Third-party dependency:** `mark3labs/mcp-go`

Justified: the MCP protocol requires initialization handshake, capability
negotiation, and JSON-RPC dispatch. Implementing this from stdlib alone is
several hundred lines of protocol plumbing with no domain value. All wren
interaction and I/O use stdlib (`os/exec`, `encoding/json`, `bufio`).

**Transport:** stdio (standard for locally-installed MCP servers)

**Wren binary:** resolved via `WREN_BIN` environment variable; defaults to
`wren` on PATH. Allows server admins to point at a specific binary without code
changes.

**Tools:**

| Tool | wren command | Parameters |
|---|---|---|
| `list_tasks` | `wren -l [filter]` | `filter` (optional string) |
| `list_done_tasks` | `wren -d` | none |
| `create_task` | `wren "title"` | `title` (required string) |
| `complete_task` | `wren -d "title"` | `title` (required string, exact) |
| `read_task` | `wren -r "keyword"` | `keyword` (required string) |
| `prepend_task` | `wren --prepend "prefix" "task"` | `prefix`, `task` (both required strings) |
| `get_random_task` | `wren -o` | none |

Each tool invokes `wren` as a subprocess, captures stdout/stderr, and returns
the output as the tool result. Non-zero exit codes are surfaced as MCP tool
errors.

### Skill (`mcp-servers/wren/SKILL.md`)

**Audience:** Agents connecting to this MCP server from a locked-down
environment.

**Content:** mirrors the CLI skill's wren conceptual content (data model,
filename conventions, config) but describes MCP tool calls instead of CLI
commands. Fully self-contained — no references outside `mcp-servers/wren/`.

---

## Shared Development Approach

Both extensions are authored in the same development cycle using `docs/wren.md`
as a common reference. Conceptual content (data model, conventions, quirks) is
written once in `docs/wren.md` and then incorporated into each `SKILL.md`
independently. Duplication between the two skills is intentional: it preserves
the self-contained constraint for each deliverable.

---

## Out of Scope

- `ww` wrapper functionality (date-filtered listing, exact-match `-x` flag,
  interactive cron/future task creation)
- `wren --summary` (requires `litellm`)
- `wren --telegram`, `--matrix`, `--http` (require external service setup)
- `wren -e` (interactive editor, not suitable for agent use)
- Task body creation via MCP (wren has no non-interactive body flag; can be
  added later if wren gains one)
