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

Core fields (required at top level, not nested under `meta`):
- `id`: Unique identifier, maps to file path
- `kind`: Node type (system, mechanic, feature, item, etc.)
- `version`: Auto-incremented on updates
- `status`: draft, review, approved, published, deprecated
- `title`: Human-readable name

References (`refs`):
- `uses`: Hard dependencies with context
- `related`: Informational links
- `emits_events`: Events this node produces
- `vocabulary`: Shared term definitions

Content (`content.sections`):
- Blocks of type: table, rule, param, mechanic, list, etc.

Issues (`issues`):
- Tracked TBDs with id, description, severity (low/medium/high/critical), location, resolved

Optional/extensible:
- `tags`, `summary`, `glossary`, `contracts`, `llm_context`, `constraints`, `reviewers`, `custom`

## CLI Commands

```
# Project setup
deco init [dir]              # Initialize project
deco create <id>             # Create new node with scaffolding

# Reading
deco list                    # List all nodes (--kind, --status, --tag)
deco show <id>               # Show node details + reverse refs
deco query [term]            # Search/filter nodes
deco validate                # Check schema + refs + constraints
deco stats                   # Project health overview
deco issues                  # List all open TBDs
deco graph                   # Output dependency graph (DOT/Mermaid)

# Modifying
deco set <id> <path> <value> # Set a field value
deco append <id> <path> <val># Append to array field
deco unset <id> <path>       # Remove a field
deco rm <id>                 # Delete a node
deco mv <old-id> <new-id>    # Rename with automatic ref updates
deco apply <id> <patch-file> # Apply structured patch (for AI)

# Review workflow
deco review submit <id>      # Submit for review
deco review approve <id>     # Approve node
deco review reject <id>      # Reject back to draft

# History
deco history [--node <id>]   # Show audit log
deco diff <id>               # Show before/after changes
deco sync                    # Detect manual edits, fix metadata
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
