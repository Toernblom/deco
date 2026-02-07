# V1 Publish Readiness — Design

## Overview

Six changes to make deco usable out of the box for a new user. Items 1, 4, 6 are blockers; 2, 3, 5 are quick wins included for polish.

## Item 6: `deco new` — Node scaffolding

```bash
deco new systems/combat --kind system --title "Combat System"
deco new mechanics/stealth --kind mechanic --title "Stealth" --tags core,pvp --summary "..."
```

- Creates `.deco/nodes/<id>.yaml` with required fields populated
- `--kind` and `--title` required
- `--tags` optional (comma-separated), `--summary` optional
- `id` derived from positional argument
- Creates parent directories automatically
- Errors if node already exists (unless `--force`)
- Validates the scaffolded node before writing
- Logs creation to history

Scaffolded file:

```yaml
id: systems/combat
kind: system
version: 1
status: draft
title: "Combat System"
tags: []
summary: ""
content:
  sections: []
```

## Item 1: `deco show --full` — Expand content inline

Default behavior unchanged (section names with block counts). `--full` renders all blocks inline.

Block renderers by type:
- **table** — ASCII table using lipgloss
- **rule** — Quoted text
- **param** — `Name: datatype [constraints]`
- **list** — Bullet list
- **text/note/description** — Plain text
- **Unknown types** — YAML dump fallback

No pagination or truncation. Users pipe to `less` if needed.

## Item 5: `deco export` — Markdown export

```bash
deco export systems/combat              # Single node to stdout (markdown default)
deco export                             # All nodes to stdout
deco export --output docs/              # Write one .md per node to directory
deco export --format markdown           # Explicit (markdown is default)
```

Markdown output per node:
- H1 title with metadata line (kind, version, status)
- Blockquote summary
- Tags
- References section
- Content sections with block rendering (tables, rules, params as markdown)
- Issues section with severity indicators

Reuses block renderer logic from `show --full`, outputting markdown syntax instead of terminal formatting.

## Item 2: `deco review status` without arguments

When called with no node ID, lists all nodes with `status == "review"`:

```
Nodes in review:
  systems/combat        v3  submitted by ai    1 approval(s)
  systems/buildings     v2  submitted by ai    0 approval(s)
```

Existing single-node behavior (`deco review status <id>`) unchanged.

## Item 3: `diff` skip no-op fields

In `printBeforeAfter()`, compare old and new values. Skip fields where `old == new`. Skip entire history entries if all fields are identical.

## Item 4: Getting started experience

### A) `deco init --template <name>`

Two built-in templates (Go `embed`):

**`game-design`:**
- Config with block types: rule, param, table, list, building, unit, mechanic
- Starter nodes: `systems/core`, `mechanics/combat` (2-3 nodes)

**`api-spec`:**
- Config with block types: endpoint, schema, rule, param
- Starter nodes: `services/api`, `schemas/auth`

`deco init` without `--template` stays as-is (blank project).

### B) README quickstart

Replace current quickstart with zero-YAML-editing flow:

```bash
deco init my-game --template game-design
cd my-game
deco list
deco show systems/core --full
deco new mechanics/stealth --kind mechanic --title "Stealth System"
deco validate
deco query --kind mechanic
```

### C) Example projects

Existing `examples/` directory (snake, space-invaders, api-spec) stays as reference material. Templates serve the zero-to-productive path.

## Implementation Order

1. **Item 6** — `deco new` (foundation, unblocks item 4)
2. **Item 1** — `deco show --full` (block renderers needed for item 5)
3. **Item 5** — `deco export` (reuses block renderers)
4. **Item 3** — `diff` no-op fix (5 lines)
5. **Item 2** — `review status` no-arg (small change)
6. **Item 4** — Templates + quickstart (depends on `new` and `show --full` existing)
