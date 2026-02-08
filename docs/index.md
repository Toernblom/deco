# Deco — Codebase & Project Index

This document is the single-source reference for understanding Deco. It describes what Deco is, how it's structured, its data model, CLI, validation system, and how to extend it. Read this before exploring the codebase.

## What Is Deco?

Deco is a Go CLI for managing complex documentation as structured, validated YAML. It treats documentation like source code — typed nodes, explicit references, schema validation, and a full audit trail.

**Use cases:** game design documents, API specifications, technical architecture, product requirements, knowledge bases.

**Key concepts:**

- **Nodes** — YAML files representing documentation units (a game mechanic, an API endpoint, a system spec)
- **References** — Typed links between nodes: `uses` (hard dependency), `related` (info), `emits_events`, `vocabulary`
- **Content blocks** — Structured data within nodes: `table`, `rule`, `param`, `mechanic`, `list`, `doc`, plus custom types
- **Custom block types** — Project-defined validation rules for blocks (required fields, types, enums, cross-refs)
- **Issues** — Tracked TBDs/questions with severity, embedded in nodes
- **Audit trail** — Append-only JSONL log tracking every change (who/what/when)
- **Status workflow** — `draft → review → approved → deprecated → archived`

**Quick facts:**

| | |
|-|-|
| Module | `github.com/Toernblom/deco` |
| Go version | 1.25.6 |
| Version | 0.9.0 |
| License | AGPL-3.0-or-later |
| Build | `go build -o deco ./cmd/deco` |
| Test | `go test -v -race ./...` |
| Install | `go install github.com/Toernblom/deco/cmd/deco@latest` |

---

## Architecture

```
CLI (Cobra commands)
  ↓
Services (validation, query, graph, refactor)
  ↓
Domain (core types, error definitions)
  ↓
Storage (YAML repos, append-only JSONL history)
```

| Layer | Package | Purpose |
|-------|---------|---------|
| Entry point | `cmd/deco` | Registers all commands, parses root flags |
| CLI | `internal/cli` | Command implementations, output formatting, style |
| Services | `internal/services/*` | Business logic (validation, query, graph, refactor) |
| Domain | `internal/domain` | Core types, error codes, error formatting |
| Errors | `internal/errors` | Error collector, YAML location tracking, suggestions |
| Storage | `internal/storage/*` | File I/O for config, nodes, history |
| Migrations | `internal/migrations` | Schema migration registry, executor, backup |

---

## Directory Tree

```
deco/
├── cmd/deco/
│   └── main.go                          # CLI entry point, registers all commands
│
├── internal/
│   ├── cli/                             # Command implementations
│   │   ├── root.go                      # Root command, global flags (verbose, quiet, color)
│   │   ├── init.go                      # deco init — initialize projects
│   │   ├── validate.go                  # deco validate — schema/refs/constraints
│   │   ├── list.go                      # deco list — list nodes with filtering
│   │   ├── show.go                      # deco show — node details + reverse refs
│   │   ├── query.go                     # deco query — advanced search/filtering
│   │   ├── sync.go                      # deco sync — detect changes, bump versions
│   │   ├── review.go                    # deco review — submit/approve/reject/status
│   │   ├── history.go                   # deco history — view audit log
│   │   ├── diff.go                      # deco diff — before/after changes
│   │   ├── graph.go                     # deco graph — dependency graph (DOT/Mermaid/ASCII)
│   │   ├── stats.go                     # deco stats — project health statistics
│   │   ├── issues.go                    # deco issues — list open TBDs
│   │   ├── migrate.go                   # deco migrate — schema migrations
│   │   ├── new.go                       # deco new — scaffold a new node
│   │   ├── export.go                    # deco export — export nodes to markdown
│   │   ├── compact.go                   # Compact (LLM-optimized) node renderer
│   │   ├── follow.go                    # Ref-following logic for compact export
│   │   ├── help.go                      # deco llm-help — LLM documentation
│   │   ├── style/style.go               # CLI color/styling system
│   │   ├── errors.go                    # ExitError, CLI error handling
│   │   ├── filter_validation.go         # Filter validation helpers
│   │   ├── audit.go                     # Audit entry formatting
│   │   ├── templates.go                 # Command templates
│   │   └── *_test.go                    # Tests for each command (~17 files)
│   │
│   ├── domain/                          # Core domain models
│   │   ├── node.go                      # Node, Block, DocRef, Contract, Reviewer
│   │   ├── ref.go                       # Ref, RefLink
│   │   ├── graph.go                     # Graph (map-based node collection)
│   │   ├── issue.go                     # Issue (TBD/question tracking)
│   │   ├── audit.go                     # AuditEntry (change log)
│   │   ├── constraint.go                # Constraint (CEL expressions)
│   │   ├── contract.go                  # Contract (Gherkin-style scenarios + parsing)
│   │   ├── error.go                     # DecoError with Location, context, suggestions
│   │   ├── error_codes.go              # Error registry (E001–E120+)
│   │   ├── error_docs.go               # Error documentation and solutions
│   │   ├── error_formatter.go          # Rust-style CLI error formatting
│   │   └── *_test.go                    # Domain model tests
│   │
│   ├── errors/                          # Error handling utilities
│   │   ├── collector.go                 # Error collection and aggregation
│   │   ├── suggestions.go              # "Did you mean?" via edit distance
│   │   └── yaml/
│   │       ├── context.go               # YAML parsing context
│   │       └── location.go              # Line/column tracking in YAML files
│   │
│   ├── services/
│   │   ├── graph/
│   │   │   └── builder.go              # Build graph, topo sort, cycle detection
│   │   ├── validator/
│   │   │   ├── validator.go            # Schema validation orchestrator
│   │   │   ├── block_validator.go      # Custom block type validation
│   │   │   ├── doc_validator.go        # External doc reference validation
│   │   │   ├── crossref_validator.go   # Cross-reference field validation
│   │   │   ├── contract.go             # Contract/Gherkin validation
│   │   │   └── *_test.go
│   │   ├── query/
│   │   │   └── query.go                # Node filtering, block search, field follow
│   │   └── refactor/
│   │       └── rename.go               # Reference update on node rename
│   │
│   ├── storage/
│   │   ├── config/
│   │   │   ├── repository.go           # Config interface + types
│   │   │   └── yaml_repository.go      # .deco/config.yaml read/write
│   │   ├── node/
│   │   │   ├── repository.go           # Node storage interface
│   │   │   ├── yaml_repository.go      # .deco/nodes/**/*.yaml CRUD
│   │   │   └── discovery.go            # Find node files by ID
│   │   └── history/
│   │       ├── repository.go           # Audit log interface + Filter type
│   │       └── jsonl_repository.go     # .deco/history.jsonl (append-only)
│   │
│   └── migrations/                      # Schema migration system
│       ├── registry.go                  # Migration registry
│       ├── executor.go                  # Migration execution
│       ├── backup.go                    # Safety backups
│       ├── schema_hash.go              # Schema versioning/hashing
│       └── *_test.go
│
├── examples/
│   ├── snake/                           # Classic Snake game design
│   ├── space-invaders/                  # Space Invaders game design
│   └── api-spec/                        # API specification example
│
├── docs/
│   ├── index.md                         # This file
│   ├── SPEC.md                          # Full specification
│   ├── llm-reference.md                # LLM-friendly YAML reference
│   ├── cli-reference.md                # CLI command reference
│   ├── plans/                           # Development planning docs
│   └── ISSUE_TEMPLATE.md               # GitHub issue template
│
├── .github/workflows/
│   ├── ci.yml                           # CI: test, lint, validate examples
│   └── release.yml                      # Release automation
│
├── CLAUDE.md                            # Claude Code instructions
├── README.md                            # User documentation
├── go.mod / go.sum                      # Go module files
└── LICENSE                              # AGPL-3.0-or-later
```

---

## Data Model

### Node (the core unit)

Every piece of documentation is a **Node**, stored as a YAML file in `.deco/nodes/`.

```yaml
id: systems/combat          # Unique ID (maps to file path)
kind: system                 # Categorization
version: 3                   # Auto-incremented on change
status: approved             # draft|review|approved|deprecated|archived
title: "Combat System"
tags: [combat, core]
summary: "Handles all combat interactions"

glossary:
  dps: "Damage per second"

refs:
  uses:
    - id: items/weapons
      context: "Weapon stats for damage calculation"
  related:
    - id: systems/inventory
  emits_events: [damage_dealt, enemy_killed]
  vocabulary: [items/weapons]

content:
  - section: "Damage Calculation"
    blocks:
      - type: rule
        name: "Base Damage"
        formula: "weapon_damage * multiplier"
      - type: table
        columns: [Level, Multiplier]
        rows:
          - [1, 1.0]
          - [10, 2.5]
      - type: param
        name: "max_damage"
        value: 9999
      - type: mechanic
        name: "Critical Hit"
        trigger: "Random chance on attack"
        effect: "2x damage"

issues:
  - id: balance-review
    description: "Need to review damage scaling at high levels"
    severity: medium
    location: "Damage Calculation"

docs:
  - path: "docs/combat-guide.md"
    keywords: [combat, damage]
    context: "Detailed combat guide for players"

custom:
  complexity: high            # Per-kind custom fields (validated via schema_rules)
```

### Type Reference

| Type | File | Key Fields |
|------|------|------------|
| `Node` | domain/node.go | id, kind, version, status, title, tags, summary, glossary, refs, content, issues, docs, contracts, reviewers, custom |
| `Ref` | domain/ref.go | uses, related, emits_events, vocabulary |
| `RefLink` | domain/ref.go | id, context |
| `Block` | domain/node.go | type, Data (map[string]interface{}) |
| `Section` | domain/node.go | section (name), blocks |
| `Content` | domain/node.go | []Section |
| `DocRef` | domain/node.go | path, keywords, context |
| `Contract` | domain/contract.go | name, scenario (Gherkin given/when/then with @node refs) |
| `Reviewer` | domain/node.go | name, timestamp, version, note |
| `Issue` | domain/issue.go | id, description, severity, location, resolved |
| `Graph` | domain/graph.go | map[string]Node — Add/Get/Remove/Update/All/Count |
| `AuditEntry` | domain/audit.go | timestamp, node_id, operation, user, content_hash, before, after |
| `Constraint` | domain/constraint.go | expr (CEL), message, scope |
| `DecoError` | domain/error.go | code, summary, detail, location, suggestion, context |
| `Location` | domain/error.go | file, line, column |
| `Config` | storage/config/repository.go | project_name, nodes_path, history_path, version, required_approvals, custom_block_types, schema_rules, schema_version, custom |
| `BlockTypeConfig` | storage/config/repository.go | required_fields, optional_fields, fields (typed FieldDef) |
| `FieldDef` | storage/config/repository.go | type (string/number/list/bool), required, enum, refs |

### Built-in Block Types

| Type | Purpose | Common Fields |
|------|---------|---------------|
| `table` | Tabular data | columns, rows |
| `rule` | Game/business rules | name, formula/condition |
| `param` | Parameters/constants | name, value |
| `mechanic` | Gameplay mechanics | name, trigger, effect |
| `list` | Simple lists | items |
| `doc` | Prose documentation | text |

Custom block types are defined in `.deco/config.yaml` with required/optional fields, type constraints, enum values, and cross-reference validation.

### Relationships Diagram

```
Node ──refs.uses──→ Node          (hard dependency)
Node ──refs.related──→ Node       (informational)
Node ──refs.emits_events──→ [string events]
Node ──refs.vocabulary──→ Node    (shared terms)
Node ──docs──→ external .md files
Node ──contracts──→ Gherkin scenarios (reference @node.id)
Node ──issues──→ embedded TBDs
Node ──content──→ sections → blocks (typed structured data)
```

---

## CLI Commands

### Project Setup
```bash
deco init [dir]                         # Initialize new project
deco init --template game-design        # With template (game-design, api-spec)
deco new <id> --kind <k> --title <t>    # Scaffold a new node
deco migrate [dir]                      # Run schema migrations
```

### Reading & Querying
```bash
deco list                               # All nodes
deco list --kind system --status draft --tag core
deco show <id> [dir]                    # Node details + reverse refs
deco show <id> --json --full            # JSON output, all fields
deco query [term] [dir]                 # Text search + filters
deco query --block-type building --field age=bronze
deco query --block-type building --follow materials
deco stats [dir]                        # Project health overview
deco stats --quiet                      # Machine-readable
deco issues [dir]                       # List open TBDs
deco graph [dir]                        # Dependency graph
deco graph --format mermaid|dot|ascii
```

### Modification
```bash
deco sync [dir]                         # Detect changes, bump versions, update history
deco sync --dry-run                     # Preview only
deco sync --no-refactor                 # Skip auto-rename reference updates
```

### Review Workflow
```bash
deco review submit <id> [dir]           # draft → review
deco review approve <id> [dir]          # review → approved (tracks reviewer)
deco review reject <id> [dir]           # review → draft
deco review status [id] [dir]           # Show review status
```

### History & Diffing
```bash
deco history [dir]                      # Full audit log
deco history --node <id>                # Filter by node
deco diff <id> [dir]                    # Before/after changes
deco diff <id> --since 2h              # Changes within timeframe
```

### Export
```bash
deco export systems/combat              # Single node to markdown (stdout)
deco export --output docs/              # All nodes → .md files
deco export --compact --kind system              # All systems, LLM-compact
deco export --compact systems/combat --follow    # Node + its dependencies
deco export --compact --kind system --follow uses --depth 2
deco export --compact --output context.md --kind system --follow all
```

### Validation & Help
```bash
deco validate [dir]                     # Validate everything
deco validate --quiet                   # Exit code only
deco llm-help                           # LLM-friendly documentation
```

### Global Flags
```
--config, -c PATH    Config directory (.deco)
--verbose            Verbose output
--quiet, -q          Suppress output
--color              auto|always|never
```

---

## Validation System

Validation runs in layers, orchestrated by `internal/services/validator/validator.go`:

| Layer | Validator | What It Checks |
|-------|-----------|----------------|
| 1. Schema | `SchemaValidator` | Required fields: id, kind, version (>0), status (valid enum), title |
| 2. Schema rules | `SchemaRulesValidator` | Per-kind custom required fields (from config `schema_rules`) |
| 3. Blocks | `BlockValidator` | Block type exists (built-in or custom), required/optional fields, type checking, enums |
| 4. Cross-refs | `CrossRefValidator` | Block field values exist in target block type (via `refs` in FieldDef) |
| 5. References | `RefValidator` | All uses/related/vocabulary targets exist; typo suggestions via edit distance |
| 6. Constraints | `ConstraintValidator` | CEL expressions evaluate to true |
| 7. Docs | `DocValidator` | External doc file paths exist, keyword matching |
| 8. Contracts | `ContractValidator` | Gherkin structure valid, @node.id references resolve |

**Error codes** range from E001–E120+ across categories: schema, references, validation, I/O, graph, contract. Each code has documentation and suggested fixes in `error_codes.go` and `error_docs.go`.

**Exit codes:** 0 = success, 1 = validation failed, 2 = schema version mismatch (needs migration).

---

## Service Layer Detail

### graph/builder.go
- `Build(nodes)` — Construct graph with duplicate detection
- `BuildDependencyMap(g)` — node → []dependencies
- `DetectCycle(g)` — (bool, []cycle path)
- `TopologicalSort(g)` — []Node in dependency order

### query/query.go
- `Filter(nodes, criteria)` — Match by kind/status/tags/text
- `FindBlocksByType(nodes, blockType)` — All blocks of a type across nodes
- `FindBlocksByField(nodes, blockType, field, value)` — Blocks with specific field value
- `FollowRefs(nodes, field, target)` — Group blocks by reference target

### refactor/rename.go
- `UpdateReferences(nodes, oldID, newID)` — Batch rename all references when a node ID changes

---

## Persistence

| What | Format | Location | Access Pattern |
|------|--------|----------|----------------|
| Config | YAML | `.deco/config.yaml` | Read on startup, write on init/migrate |
| Nodes | YAML (one per node) | `.deco/nodes/**/*.yaml` | CRUD via `node.Repository` |
| History | JSONL (append-only) | `.deco/history.jsonl` | Append via `history.Repository`, query with filters |

**History operations:** create, update, delete, set, append, unset, move, submit, approve, reject, sync, baseline, migrate, rewrite.

---

## Migration System

Located in `internal/migrations/`:

- **Schema version** — SHA hash computed from `custom_block_types` + `schema_rules` in config
- Stored in `config.yaml` as `schema_version`
- `deco validate` checks version before validation (exit code 2 on mismatch)
- `deco migrate` runs registered migrations with automatic backup

---

## Testing & CI

**Running tests:**
```bash
go test -v -race ./...
```

**CI** (`.github/workflows/ci.yml`):
- Platforms: ubuntu-latest, macos-latest, windows-latest
- Go version: 1.24
- Steps: `gofmt` check → `go vet` → `go test -race ./...` → build binary → validate all examples

**Example validation:** CI builds the binary and runs `deco validate` on `examples/snake`, `examples/space-invaders`, and `examples/api-spec`.

---

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/spf13/cobra` | v1.10.2 | CLI framework |
| `github.com/google/cel-go` | v0.27.0 | CEL expression evaluation (constraints) |
| `gopkg.in/yaml.v3` | v3.0.1 | YAML parsing/marshaling |
| `github.com/fatih/color` | v1.18.0 | Terminal color output (indirect) |

---

## Task Mapping (What to Edit)

| Task | Files to Touch |
|------|---------------|
| New CLI command | `internal/cli/<cmd>.go` + register in `cmd/deco/main.go` |
| New error code | `internal/domain/error_codes.go` + `error_docs.go` |
| Node schema change | `internal/domain/node.go` + `internal/storage/node/yaml_repository.go` |
| New validation rule | `internal/services/validator/` (add or extend a validator) |
| Custom block types | `internal/storage/config/repository.go` + `internal/services/validator/block_validator.go` |
| Reference handling | `internal/domain/ref.go` + `internal/services/graph/builder.go` |
| Sync behavior | `internal/cli/sync.go` + `internal/services/refactor/rename.go` |
| Audit operations | `internal/storage/history/jsonl_repository.go` + `internal/domain/audit.go` |
| CLI output/styling | `internal/cli/style/style.go` |
| Schema migration | `internal/migrations/` |

---

## Example Projects

| Example | Path | Demonstrates |
|---------|------|-------------|
| Snake | `examples/snake/` | Custom block types (`powerup`), schema rules, game mechanics |
| Space Invaders | `examples/space-invaders/` | Multi-system game design |
| API Spec | `examples/api-spec/` | API documentation use case, schemas/endpoints |

All examples validate successfully in CI and serve as integration tests.
