# Deco Codebase Structure

This document provides a map for AI agents working on the codebase.

## Directory Tree

```
deco/
├── cmd/deco/
│   └── main.go                     # CLI entry point, registers commands
│
├── internal/
│   ├── cli/                        # Command implementations
│   │   ├── root.go                 # Root command, global flags (verbose, quiet)
│   │   ├── init.go                 # Initialize new deco projects
│   │   ├── validate.go             # Validate schema, refs, constraints
│   │   ├── list.go                 # List nodes with filtering
│   │   ├── show.go                 # Display node details + reverse refs
│   │   ├── query.go                # Search/filter nodes
│   │   ├── sync.go                 # Detect changes, renames, deletions
│   │   ├── review.go               # Review workflow (submit/approve/reject)
│   │   ├── history.go              # View audit log
│   │   ├── diff.go                 # Show node change history
│   │   ├── graph.go                # Output dependency graph
│   │   ├── stats.go                # Project statistics
│   │   ├── issues.go               # List TBDs/issues
│   │   ├── migrate.go              # Schema migrations
│   │   └── *_test.go               # Tests for each command
│   │
│   ├── domain/                     # Core domain models
│   │   ├── node.go                 # Node: id, kind, version, status, title, tags, content, refs
│   │   ├── ref.go                  # Reference types: uses, related, emits_events, vocabulary
│   │   ├── graph.go                # Graph: collection of nodes with lookup
│   │   ├── issue.go                # Issue: TBDs/questions with severity
│   │   ├── audit.go                # AuditEntry: who/what/when tracking
│   │   ├── constraint.go           # Constraint: CEL expressions
│   │   ├── error.go                # DecoError with code/summary/detail
│   │   ├── error_codes.go          # Error codes (E001-E020+)
│   │   ├── error_docs.go           # Error documentation
│   │   ├── error_formatter.go      # CLI error formatting
│   │   └── *_test.go               # Tests for each model
│   │
│   ├── errors/                     # Error handling utilities
│   │   ├── collector.go            # Collect and manage errors
│   │   ├── suggestions.go          # Generate helpful suggestions
│   │   └── yaml/
│   │       ├── context.go          # YAML parsing context
│   │       └── location.go         # Error locations in YAML
│   │
│   ├── services/                   # Business logic
│   │   ├── graph/
│   │   │   └── builder.go          # Build graph, detect cycles
│   │   ├── validator/
│   │   │   └── validator.go        # Schema validation, required fields
│   │   ├── query/
│   │   │   └── query.go            # Node filtering/search
│   │   └── refactor/
│   │       └── rename.go           # Reference updates (used by sync)
│   │
│   └── storage/                    # Persistence layer
│       ├── config/
│       │   ├── repository.go       # Config storage interface
│       │   └── yaml_repository.go  # .deco/config.yaml implementation
│       ├── node/
│       │   ├── repository.go       # Node storage interface
│       │   ├── yaml_repository.go  # .deco/nodes/*.yaml implementation
│       │   └── discovery.go        # Find node files in directory tree
│       └── history/
│           ├── repository.go       # Audit log interface
│           └── jsonl_repository.go # .deco/history.jsonl implementation
│
├── examples/                       # Example game designs
│   ├── snake/                      # Classic Snake game
│   │   └── .deco/nodes/            # food, core, scoring
│   └── space-invaders/             # Space Invaders
│       └── .deco/nodes/            # aliens, player, core
│
├── docs/                           # Documentation
│   ├── index.md                    # This file (codebase structure)
│   ├── llm-reference.md            # Complete YAML reference for LLMs
│   └── cli-reference.md            # CLI command reference
│
├── .github/workflows/              # CI/CD
│   ├── ci.yml                      # Continuous Integration
│   └── release.yml                 # Release pipeline
│
├── .beads/                         # Issue tracking (bd)
│   ├── issues.jsonl                # Issues as JSON lines
│   └── config.yaml                 # Beads config
│
├── CLAUDE.md                       # Claude Code instructions
├── AGENTS.md                       # Agent instructions for bd
├── SPEC.md                         # Full Deco specification
├── README.md                       # Project documentation
├── go.mod                          # Go module (github.com/Toernblom/deco)
└── go.sum                          # Dependency checksums
```

## Architecture

```
CLI (cmd/deco) → Commands (internal/cli) → Services (internal/services) → Domain (internal/domain) → Storage (internal/storage)
```

| Layer | Purpose |
|-------|---------|
| `cmd/deco` | Entry point, command registration |
| `internal/cli` | Parse flags, orchestrate services, format output |
| `internal/services` | Business logic (validation, patching, queries) |
| `internal/domain` | Core types (Node, Graph, Ref, Issue, Error) |
| `internal/storage` | File I/O (YAML nodes, JSONL history) |

## Key Files by Task

| Task | Files |
|------|-------|
| Add new command | `cmd/deco/main.go`, `internal/cli/<cmd>.go` |
| Modify node structure | `internal/domain/node.go`, `internal/storage/node/yaml_repository.go` |
| Change validation rules | `internal/services/validator/validator.go` |
| Update reference handling | `internal/domain/ref.go`, `internal/services/graph/builder.go` |
| Add new error code | `internal/domain/error_codes.go`, `internal/domain/error_docs.go` |
| Change sync behavior | `internal/cli/sync.go`, `internal/services/refactor/rename.go` |
