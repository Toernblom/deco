# Deco Specification

## What It Is

A Go CLI that replaces traditional GDDs entirely. Your game design lives as structured, validated YAML - no separate docs, wikis, or Notion pages needed. Deco is the GDD.

## Core Decisions

| Aspect | Decision |
|--------|----------|
| Language | Go |
| License | Open source (eventually) |
| Schema | Extensible - core fields required, custom sections allowed |
| File location | Configurable, default `.deco/` |
| History | Audit log (append-only: who, what, when) |
| AI workflow | Both patch operations and full file rewrites |
| Development | TDD - tests first, then implementation |

## Node Structure

Core fields (required):
- `meta`: id, kind, version, status, title, tags
- `refs`: uses, related, emits_events, vocabulary
- `content.sections`: blocks of type table, rule, param, mechanic, list
- `issues`: tracked TBDs/questions with severity and location

Optional/extensible:
- `summary`, `glossary`, `contracts`, `llm_context`
- Projects can add custom sections

## CLI Commands (MVP)

```
deco init                    # Initialize project
deco validate                # Check schema + refs, report errors
deco list                    # List all nodes
deco show <id>               # Show node details + reverse refs
deco query <filter>          # Search/filter nodes
deco set <id> <path> <value> # Patch a field
deco append <id> <path> <value>
deco unset <id> <path>
deco mv <old-id> <new-id>    # Rename with automatic ref updates
deco history [<id>]          # Show audit log
deco apply <patch-file>      # Apply structured patch (for AI)
```

## Design Principles

**No ambiguity in a complete GDD.** A design is "complete" when:
- Zero open `issues` (no TBDs, no unanswered questions)
- All refs resolve to existing nodes
- All required fields populated

**No contradictions.** Deco validates internal consistency through:
- Enum values must match their definitions
- Explicit constraints you define are enforced across nodes

Example constraint:
```yaml
constraints:
  - expr: "self.needs.food.threshold == refs['systems/food'].starvation_time"
    message: "Food threshold must match starvation time in food system"
```

Deco doesn't magically detect semantic contradictions - you declare what must be consistent, and Deco enforces it.

## What Deco Does NOT Do

- Gherkin compilation or test execution
- UI or editor plugins (CLI only)
- Game engine integration (it's just a design tool)

## Project Layout

```
.deco/
  config.yaml          # Project configuration
  history.jsonl        # Audit log
  nodes/
    systems/
      settlement/
        colonists.yaml
        housing.yaml
```

## AI Integration

Two modes:
1. **Patch mode**: AI outputs JSON patch operations, Deco validates and applies
2. **Rewrite mode**: AI rewrites full YAML files, Deco validates after

The engine is the source of truth, not the AI.
