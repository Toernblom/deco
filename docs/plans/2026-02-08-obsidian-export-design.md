# Obsidian Export Design

## Overview

Add `--obsidian` flag to `deco export` that generates an Obsidian-compatible vault at `.deco/vault/`, mirroring the `.deco/nodes/` folder structure.

## Behavior

- `deco export --obsidian` exports all nodes to `.deco/vault/`
- Wipes and regenerates vault on each run (clean export)
- Mirrors `nodes/` directory structure exactly
- Filters (`--kind`, `--status`, `--tag`) still work for partial vaults
- `--output`, `--compact`, `--follow`, `--depth` are ignored when `--obsidian` is set

## Per-Node Markdown File

### Frontmatter (rich)

```yaml
---
id: systems/core
kind: system
version: 3
status: approved
tags: [core, multiplayer]
summary: Core game loop and state management
uses:
  - systems/networking
  - mechanics/combat
related:
  - systems/auth
emits_events:
  - game_started
  - round_ended
---
```

### Body

- H1 title
- Summary paragraph
- Tags as `#tag` inline (Obsidian-native)
- Glossary as definition list
- Content sections (H2) with blocks rendered per type
- References section with `[[node/id|Title]]` wikilinks and context
- Issues as Obsidian callouts (`[!warning]`, `[!info]`, `[!bug]`, etc.)
- Custom fields section

### Wikilinks

- All refs render as `[[node/id|Node Title]]`
- Broken refs (target not found) fall back to `[[node/id]]`

### Block Rendering

- `table` → markdown table
- `rule` → blockquote with bold name
- `param` → bold name, inline details (type, range, default, unit)
- `mechanic` → bold name with condition/input/output sub-lists
- `list` → bullet list
- `text/note/description` → plain prose
- Custom/unknown → structured definition list (bold field names, values)

### Issues

Rendered as Obsidian callouts with severity mapping:
- `critical` → `[!bug]`
- `high` → `[!warning]`
- `medium` → `[!caution]`
- `low` → `[!info]`

## Implementation

- New file: `internal/cli/obsidian.go`
- Add `--obsidian` flag to existing export command in `export.go`
- Tests in `internal/cli/obsidian_test.go`
