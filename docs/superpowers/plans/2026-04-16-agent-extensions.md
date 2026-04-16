# Agent Extensions Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create a lightweight monorepo for small agent-extension deliverables, organized by extension type, to reduce cognitive load without over-scoping the repository.

**Architecture:** The repository uses an extension-first layout with top-level folders for each extension type and second-level folders for deliverables. Deliverables stay independent by default, even when they share a name across extension types. The initial phase is intentionally minimal: define the repo contract in docs and create only the top-level folder skeleton after plan approval.

**Tech Stack:** Markdown docs, empty directory skeletons, and repo-level conventions enforced by written guidance rather than automation.

---

## Decisions Locked In

- Repository name: `agent-extensions`
- Primary organization axis: extension type first
- Top-level folders only for the first pass
- Second-level folder names are deliverable names, not implied relationships
- Initial deliverable baseline: `README.md` and `AGENTS.md`
- Initial phase scope: plan document plus top-level folder skeleton only
- No deliverable templates in the first pass

## Task 1: Write the Repo Contract

**Files:**
- Create: `README.md`
- Create: `AGENTS.md`

- [ ] **Step 1: Draft the repository-purpose README**

```markdown
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
```

- [ ] **Step 2: Draft the agent instructions file**

```markdown
# AGENTS.md

This repository contains multiple small agent-extension deliverables.

## Repo Rules

- Keep the structure extension-first.
- Use top-level folders only for extension types.
- Name second-level folders after the deliverable.
- Do not infer relationships between same-named deliverables in different extension-type folders.
- Keep deliverable layouts minimal unless a specific type needs more structure.

## Current Repo Scope

The first pass focuses on docs and a top-level folder skeleton only.
```

- [ ] **Step 3: Review the repo contract for consistency**

Run:

```bash
sed -n '1,220p' README.md
sed -n '1,220p' AGENTS.md
```

Expected: both files match the decisions locked in above and do not introduce templates or extra structure.

## Task 2: Create the Top-Level Skeleton

**Files:**
- Create: `skills/.gitkeep`
- Create: `mcp-servers/.gitkeep`
- Create: `plugins/.gitkeep`
- Create: `hooks/.gitkeep`

- [ ] **Step 1: Create the four top-level extension folders**

Run:

```bash
mkdir -p skills mcp-servers plugins hooks
touch skills/.gitkeep mcp-servers/.gitkeep plugins/.gitkeep hooks/.gitkeep
```

- [ ] **Step 2: Verify the skeleton**

Run:

```bash
find . -maxdepth 2 -type f | sort
find . -maxdepth 1 -type d | sort
```

Expected: only the four top-level extension folders exist at the repository root, and the repo contains no deliverable templates yet.

## Task 3: Record the Initial Deliverable Conventions

**Files:**
- Create: `docs/README.md`
- Create: `docs/superpowers/README.md`
- Create: `docs/superpowers/plans/README.md`

- [ ] **Step 1: Document where future plans live**

```markdown
# Documentation Index

This repository keeps planning and design artifacts under `docs/superpowers/`.

- `docs/superpowers/plans/` contains implementation plans.
- Future design notes can live alongside plans if the repository later needs them.
```

- [ ] **Step 2: Verify the documentation path**

Run:

```bash
find docs -maxdepth 3 -type f | sort
```

Expected: the docs path exists and clearly points future contributors to the plan artifacts.

## Task 4: Final Sanity Check

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/README.md`

- [ ] **Step 1: Scan for scope creep**

Run:

```bash
rg -n "template|scaffold|TODO|TBD|later" README.md AGENTS.md docs
```

Expected: no placeholder language and no deliverable templates.

- [ ] **Step 2: Confirm the plan matches the repo shape**

Run:

```bash
git status --short
```

Expected: only the intended docs and skeleton files are present, with no unrelated changes.
