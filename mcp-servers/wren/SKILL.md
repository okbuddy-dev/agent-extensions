## Overview

`wren` is a plain-file todo manager. Tasks are files in a flat notes directory;
the filename is the task title and the file content (if any) is the body.

You interact with wren through MCP tools that wrap the `wren` CLI. All tools
invoke `wren` as a subprocess and return its output.

**Config** (on the server machine):
- `~/.config/wren/wren.json` — `notes_dir` and `done_dir` keys
- `notes_dir` default: `~/Documents/notes`
- `done_dir` default: `done` (relative to `notes_dir`)

The MCP server URL is provided by the agent runtime. This skill describes the tools available once connected.

---

## Tools

### `list_tasks`

List active tasks. Optionally filter by keyword.

```
list_tasks()
list_tasks(filter: "keyword")
```

Output lines are prefixed with `➜ `. Strip that prefix when presenting to the user.

### `list_done_tasks`

List completed tasks.

```
list_done_tasks()
```

### `create_task`

Create a new task. The title becomes the filename.

```
create_task(title: "buy milk")
```

### `complete_task`

Mark a task as done. Use the exact task title to avoid ambiguity.

```
complete_task(title: "buy milk")
```

**Cron task caveat:** wren stores `*` as fullwidth `＊` (U+FF0A) in filenames.
Use `list_tasks` to get the exact stored filename before completing a cron task.

### `read_task`

Read the body content of a task.

```
read_task(keyword: "milk")
```

Returns nothing if the task body is empty.

### `prepend_task`

Prepend a string to a task's title (renames the file).

```
prepend_task(prefix: "2026-05-01", task: "renew subscription")
```

### `get_random_task`

Return one randomly selected active task.

```
get_random_task()
```

---

## Scheduled Tasks

Task titles follow special filename conventions that wren interprets.

### Future (one-off)

Prefix with a date: `create_task(title: "2026-05-01 renew subscription")`.
The task is hidden until that date.

### Recurring (cron)

Prefix with five cron fields: `create_task(title: "0 9 * * 1 weekly review")`.
When marked done, wren copies (not moves) the task — it reappears on the next
matching schedule.

Cron field order: `minute hour day-of-month month day-of-week`

---

## What Is Not Supported

Task body creation via MCP is out of scope — wren has no non-interactive body
flag. Interactive (`-e`), AI-summary (`-s`), and external-service operations
(`--telegram`, `--matrix`, `--http`) are not exposed.

