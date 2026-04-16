# Wren Developer Reference

Developer-only reference used to author `skills/wren/SKILL.md` and
`mcp-servers/wren/SKILL.md`. Not deployed or referenced by agents.

---

## Binary and Config

- **Executable:** `wren` on PATH; installed to `~/.local/bin/wren`
- **Config file:** `~/.config/wren/wren.json`
- **Config keys:**
  - `notes_dir` — absolute path to notes directory (default: `~/Documents/notes`)
  - `done_dir` — path to done directory, **relative to `notes_dir`** (default: `done`)

---

## Data Model

Tasks are plain files in `notes_dir`:

- **Title** = filename (no extension)
- **Body** = file contents (may be empty)
- No subdirectories — `notes_dir` is flat

Completed tasks are moved to `done_dir` (a subdirectory of `notes_dir`).

---

## Filename Conventions

### Date prefix (one-off future tasks)

```
YYYY-MM-DD task title
```

The task is hidden from listings until that date. Example: `2026-05-01 renew subscription`.

### Cron prefix (recurring tasks)

```
M H dom mon dow task title
```

Five space-separated cron fields followed by the task title. When marked done,
wren copies (not moves) the file to `done_dir` — the task reappears on the next
matching schedule.

Example: `0 9 * * 1 weekly review` (every Monday at 9am).

**Caveat:** wren replaces `*` with fullwidth `＊` (U+FF0A) in filenames. The
stored filename contains `＊`, not `*`. List tasks to get the exact filename
before completing.

---

## Supported Operations

### List active tasks

```bash
wren -l          # all active tasks
wren -l keyword  # filter by substring
```

Output: one task per line, prefixed with `➜ `.

### List done tasks

```bash
wren -d
```

### Create a task

```bash
wren "task title"
```

Creates an empty file in `notes_dir`. To include a body, write the file directly:

```bash
echo "body content" > ~/Documents/notes/"task title"
```

### Complete a task

```bash
wren -d "title"
```

Moves the file to `done_dir`. Partial title match may prompt for disambiguation
interactively — use the exact filename to avoid this. For cron tasks, remember
`*` → `＊` in stored filenames.

### Read a task body

```bash
wren -r "keyword"
```

Prints the file contents. Returns nothing if the file is empty.

### Prepend to a task title

```bash
wren --prepend "prefix" "task title"
```

Renames the file, prepending `prefix ` to the filename.

### Random task

```bash
wren -o
```

Prints one randomly selected active task title.

---

## Excluded Operations

| Flag | Reason |
|---|---|
| `-e` / `--edit` | Opens `$EDITOR` interactively — not usable by agents |
| `-s` / `--summary` | Requires `litellm` external dependency |
| `--telegram` | Requires external Telegram bot setup |
| `--matrix` | Requires external Matrix homeserver setup |
| `--http` | Starts an HTTP server; not a task operation |
