---
name: mcp-wren
description: >
  How to use the 'wren' CLI to interact with wren on another
  server. Use whenever creating, listing, deleting or managing
  todos, or things to do.
---

## Overview

`wren` is a plain-file todo manager. Tasks are files in a flat notes directory;
the filename is the task title and the file content (if any) is the body.

- **Executable:** `wren` (PATH); fallback `~/.local/bin/wren`
- **Config:** `~/.config/wren/wren.json`
  - `notes_dir` — notes directory (default: `~/Documents/notes`)
  - `done_dir` — done subdirectory, relative to `notes_dir` (default: `done`)

---

## Commands

### List tasks

```bash
wren -l              # all active tasks
wren -l "keyword"    # filter by substring
wren -d              # list done tasks
```

Output lines are prefixed with `➜ `. Strip that prefix when presenting to the user.

### Create a task

```bash
wren "task title"
```

Creates an empty file. The title becomes the filename.

To create a task with body content, write the file directly:

```bash
echo "body content" > ~/Documents/notes/"task title"
```

### Complete a task

```bash
wren -d "exact task title"
```

Moves the file to `done_dir`. Use the exact filename to avoid interactive
disambiguation prompts.

**Cron task caveat:** wren stores `*` as fullwidth `＊` (U+FF0A) in filenames.
Always use `wren -l` to get the exact stored filename before completing a cron
task.

### Read a task

```bash
wren -r "keyword"
```

Prints file contents. Returns nothing if the task body is empty.

### Prepend to a task title

```bash
wren --prepend "prefix" "task title"
```

Renames the file by prepending `prefix ` to the filename.

### Random task

```bash
wren -o
```

Prints one randomly selected active task.

---

## Scheduled Tasks

### Future (one-off)

Prefix the title with a date: `wren "2026-05-01 renew subscription"`.
The task is hidden until that date.

### Recurring (cron)

Prefix the title with five cron fields:

```bash
wren "0 9 * * 1 weekly review"
```

When marked done, wren copies (not moves) the file — the task reappears on
the next matching schedule.

Cron field order: `minute hour day-of-month month day-of-week`

---

## What Not to Use

| Flag | Why |
|---|---|
| `-e` / `--edit` | Opens `$EDITOR` interactively — unusable in agent context |
| `-s` / `--summary` | Requires `litellm` |
| `--telegram` / `--matrix` / `--http` | Require external service setup |
