# Deco — Design Engine for Game Documents

Deco is a CLI for managing game design documents (GDDs) as **structured, validated YAML**. Instead of treating design as scattered prose, Deco treats it like source code: typed nodes, explicit references, schema validation, and an audit trail.

## Why Deco?

If you've maintained a big GDD, you know how it ends:
- "TBDs" everywhere, contradictory rules, outdated pages
- Updating one system means hunting references across 30 documents
- An LLM helps write content, then loses context, and now you have two truths

Deco fixes this by making your GDD:
- **Validated** — Schema, reference, and constraint checking
- **Refactorable** — Rename nodes and all references update automatically
- **Auditable** — Every change tracked with who/what/when
- **AI-friendly** — Patch operations or full rewrites, the engine enforces correctness

## Quick Start

```bash
# Install
go install github.com/anthropics/deco/cmd/deco@latest

# Initialize a project
deco init my-game
cd my-game

# Create your first design node
cat > .deco/nodes/systems/combat.yaml << 'EOF'
id: systems/combat
kind: system
version: 1
status: draft
title: "Combat System"
tags: [core, gameplay]

summary: "Turn-based tactical combat with action points."

refs:
  uses:
    - target: systems/characters
      context: "Characters have combat stats"

content:
  sections:
    - name: overview
      blocks:
        - type: rule
          text: "Each turn, characters spend action points to perform moves."
EOF

# Validate your design
deco validate

# List all nodes
deco list

# Show node details with reverse references
deco show systems/combat
```

## Project Structure

```
.deco/
├── config.yaml          # Project configuration
├── history.jsonl        # Append-only audit log
└── nodes/               # Design documents (nested directories supported)
    ├── systems/
    │   ├── combat.yaml
    │   └── settlement/
    │       ├── colonists.yaml
    │       └── housing.yaml
    ├── items/
    │   └── weapons.yaml
    └── characters/
        └── hero.yaml
```

Node IDs map to file paths: `systems/settlement/colonists` → `.deco/nodes/systems/settlement/colonists.yaml`

## Node Format

Every node has required fields and optional extensions:

```yaml
# Required fields
id: systems/settlement/colonists
kind: system                      # Any type: system, mechanic, item, character, quest...
version: 1                        # Auto-incremented on updates
status: draft                     # draft, review, approved, published, deprecated
title: "Settlement: Colonists"

# Optional metadata
tags: [settlement, population]
summary: "Named individuals who power the settlement economy."

# References to other nodes
refs:
  uses:                           # Hard dependencies
    - target: systems/settlement/housing
      context: "Colonists need shelter"
  related:                        # Informational links
    - target: systems/settlement/morale
  emits_events:                   # Events this system produces
    - events/colonist/born
    - events/colonist/died

# Structured content
content:
  sections:
    - name: needs
      blocks:
        - type: table
          columns: [need, effect_if_unmet]
          rows:
            - need: food
              effect_if_unmet: "Starvation → Death"
        - type: rule
          text: "Unmet needs reduce morale daily."
        - type: param
          key: starting_population
          value: { min: 4, max: 5 }

# Tracked questions/TBDs
issues:
  - id: tbd_child_duration
    description: "Define child stage duration"
    severity: medium
    resolved: false

# BDD-style acceptance criteria
contracts:
  - name: "Starvation kills colonists"
    given:
      - colonist: { id: c1, needs: { food: 0 } }
    when:
      - tick: { hours: 48 }
    then:
      - expect: { colonist_state: dead, cause: starvation }

# CEL-based validation rules
constraints:
  - expr: "version > 0"
    message: "Version must be positive"

# Extensible custom fields
custom:
  designer_notes: "Balance this after playtesting"
```

## CLI Commands

### Project Setup

```bash
deco init [directory]        # Initialize a new project
deco init . --force          # Reinitialize existing project
```

### Querying & Reading

```bash
deco list                    # List all nodes
deco list --kind system      # Filter by type
deco list --status draft     # Filter by status
deco list --tag combat       # Filter by tag

deco show <node-id>          # Show node details + reverse references
deco show systems/combat --json   # Output as JSON

deco query sword             # Search title/summary for "sword"
deco query --kind item --status published   # Combined filters

deco validate                # Validate all nodes (schema, refs, constraints)
deco validate --quiet        # Exit code only (for CI)
```

### Modifying Nodes

```bash
# Set a field value
deco set systems/combat title "Tactical Combat"
deco set systems/combat status approved
deco set systems/combat tags[0] core-gameplay

# Append to arrays
deco append systems/combat tags stealth

# Remove fields
deco unset systems/combat summary
deco unset systems/combat tags[2]

# Batch operations (transactional)
deco apply systems/combat patch.json
deco apply systems/combat patch.json --dry-run
```

Patch file format:
```json
[
  {"op": "set", "path": "title", "value": "New Title"},
  {"op": "append", "path": "tags", "value": "new-tag"},
  {"op": "unset", "path": "summary"}
]
```

### Refactoring

```bash
# Rename a node (automatically updates all references)
deco mv old-node-id new-node-id
deco mv systems/combat systems/tactical-combat
```

### Audit Trail

```bash
deco history                 # Show all changes
deco history --node systems/combat   # Filter by node
deco history --limit 10      # Limit entries
```

Output shows: TIME, NODE, OPERATION, USER

## Validation

Deco validates your design graph across three dimensions:

### Schema Validation
Every node must have: `id`, `kind`, `version`, `status`, `title`

```bash
$ deco validate
ERROR [E008] Missing required field: status
  → systems/combat.yaml:1
  Suggestion: Add 'status: draft' to the node
```

### Reference Validation
All references must resolve to existing nodes:

```bash
$ deco validate
ERROR [E020] Reference not found: systems/combaat
  → systems/weapons.yaml:12
  Did you mean: systems/combat?
```

### Constraint Validation
CEL expressions enforce custom rules:

```yaml
constraints:
  - expr: "version > 0"
    message: "Version must be positive"
  - expr: "status in ['draft', 'approved', 'published']"
    message: "Invalid status"
  - expr: "size(tags) > 0"
    message: "At least one tag required"
```

## Contracts

Define testable acceptance criteria using BDD-style scenarios:

```yaml
contracts:
  - name: "Hunger damages health"
    scenario: "A colonist with no food loses health over time"
    given:
      - colonist: { id: c1, health: 100, needs: { food: 0 } }
    when:
      - tick: { hours: 24 }
    then:
      - expect: { health: 80 }
      - expect_event: { id: events/colonist/starving }

  - name: "Death from starvation"
    given:
      - colonist: { id: c1, health: 10, needs: { food: 0 } }
    when:
      - tick: { hours: 24 }
    then:
      - expect: { colonist_state: dead, cause: starvation }
```

Contracts serve as executable specifications — your design documents what should happen, and tests verify it actually does.

## Issues: Tracking TBDs

Mark unresolved questions directly in nodes:

```yaml
issues:
  - id: tbd_damage_formula
    description: "Finalize damage calculation formula"
    severity: high          # low, medium, high, critical
    location: "content.sections.damage"
    resolved: false
```

A design is complete when all issues are resolved, all references exist, and all constraints pass.

## AI Integration

Deco is designed for AI-assisted design workflows. The engine validates all changes, so AI can propose updates without breaking consistency.

### Two Update Modes

**Patch Mode** — Surgical changes via operations:
```json
[
  {"op": "set", "path": "status", "value": "approved"},
  {"op": "append", "path": "tags", "value": "balanced"},
  {"op": "unset", "path": "issues[0]"}
]
```

**Rewrite Mode** — AI generates complete YAML, Deco validates before saving.

### LLM Context Control

Configure what context AI sees per node:

```yaml
llm_context:
  include:
    - summary
    - glossary
    - refs
    - content.sections
  exclude:
    - custom.internal_notes
```

### Workflow Example

1. Designer describes intent in natural language
2. AI generates structured YAML or patch operations
3. Deco validates schema, references, and constraints
4. Changes applied only if valid
5. Audit log records who/what/when

The engine — not the AI — is the source of truth.

## Typical Workflow

```bash
# Start a new game design
deco init my-rpg && cd my-rpg

# Create core systems
# (manually or via AI-generated YAML)

# Validate as you go
deco validate

# Explore your design
deco list --kind system
deco show systems/combat
deco query "damage"

# Make changes
deco set systems/combat status approved
deco append systems/combat tags balanced

# Refactor when needed
deco mv items/sword items/weapons/sword

# Review history
deco history --node systems/combat

# CI integration
deco validate --quiet && echo "Design valid"
```

## Architecture

```
cmd/deco/main.go          # CLI entry point

internal/
├── domain/               # Core types: Node, Ref, Graph, AuditEntry
├── storage/              # Persistence layer
│   ├── config/           # Project config (YAML)
│   ├── node/             # Node storage (YAML files)
│   └── history/          # Audit log (JSONL)
├── services/             # Business logic
│   ├── graph/            # Dependency graph, cycle detection, reverse indexing
│   ├── validator/        # Schema, reference, and constraint validation
│   ├── patcher/          # Set/append/unset operations
│   ├── query/            # Filtering and search
│   └── refactor/         # Rename with reference updates
└── cli/                  # Command implementations
```

### Key Design Decisions

- **YAML for nodes** — Human-readable, diffable, merge-friendly
- **JSONL for history** — Append-only, line-per-entry for streaming
- **CEL for constraints** — Google's Common Expression Language for validation rules
- **No database** — Everything is files, works with git

## Installation

### From Source

```bash
git clone https://github.com/anthropics/deco.git
cd deco
go build -o deco ./cmd/deco
./deco --help
```

### Go Install

```bash
go install github.com/anthropics/deco/cmd/deco@latest
```

## Requirements

- Go 1.21+

## Documentation

- [SPEC.md](SPEC.md) — Full specification and design decisions

## License

Open source (license TBD)
