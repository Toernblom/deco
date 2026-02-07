# Deco — Structured Documentation for Complex Systems

Deco is a CLI for managing complex documentation as **structured, validated YAML**. Instead of treating specs as scattered prose, Deco treats them like source code: typed nodes, explicit references, schema validation, and an audit trail.

## Why Deco?

If you've maintained complex documentation—system specs, API designs, game design docs, product requirements—you know how it ends:
- "TBDs" everywhere, contradictory rules, outdated pages
- Updating one component means hunting references across 30 documents
- An LLM helps write content, then loses context, and now you have two truths

Deco fixes this by making your documentation:
- **Validated** — Schema, reference, and constraint checking
- **Refactorable** — Rename nodes and all references update automatically
- **Auditable** — Every change tracked with who/what/when
- **AI-friendly** — Patch operations or full rewrites, the engine enforces correctness

## Use Cases

- **Game Design Documents** — Systems, mechanics, items, characters with interconnected rules
- **API Specifications** — Endpoints, schemas, authentication, versioning with dependency tracking
- **Technical Architecture** — Components, interfaces, data flows with validation
- **Product Requirements** — Features, user stories, acceptance criteria with status workflow
- **Knowledge Bases** — Interconnected concepts with reference integrity

## Quick Start

```bash
# Install
go install github.com/anthropics/deco/cmd/deco@latest

# Initialize a project
deco init my-project
cd my-project

# Create your first node
cat > .deco/nodes/systems/auth.yaml << 'EOF'
id: systems/auth
kind: system
version: 1
status: draft
title: "Authentication System"
tags: [core, security]

summary: "JWT-based authentication with refresh tokens."

refs:
  uses:
    - target: systems/users
      context: "Authenticates user identities"

content:
  sections:
    - name: overview
      blocks:
        - type: rule
          text: "All API endpoints require valid JWT except /auth/login and /auth/refresh."
EOF

# Validate your documentation
deco validate

# List all nodes
deco list

# Show node details with reverse references
deco show systems/auth
```

## Project Structure

```
.deco/
├── config.yaml          # Project configuration
├── history.jsonl        # Append-only audit log
└── nodes/               # Documentation nodes (nested directories supported)
    ├── systems/
    │   ├── auth.yaml
    │   └── api/
    │       ├── endpoints.yaml
    │       └── schemas.yaml
    ├── components/
    │   └── database.yaml
    └── requirements/
        └── user-stories.yaml
```

Node IDs map to file paths: `systems/api/endpoints` → `.deco/nodes/systems/api/endpoints.yaml`

## Node Format

Every node has required fields and optional extensions:

```yaml
# Required fields
id: systems/api/endpoints
kind: system                      # Any type: system, component, feature, requirement...
version: 1                        # Auto-incremented on updates
status: draft                     # draft, review, approved, published, deprecated
title: "API Endpoints"

# Optional metadata
tags: [api, rest]
summary: "REST API endpoint definitions and contracts."

# References to other nodes
refs:
  uses:                           # Hard dependencies
    - target: systems/auth
      context: "All endpoints require authentication"
  related:                        # Informational links
    - target: components/database
  emits_events:                   # Events this system produces
    - events/api/request-logged
  vocabulary:                     # Shared term definitions
    - glossaries/api-terms

# Structured content
content:
  sections:
    - name: endpoints
      blocks:
        - type: table
          columns: [method, path, description]
          rows:
            - method: GET
              path: /users/{id}
              description: "Retrieve user by ID"
        - type: rule
          text: "All endpoints return JSON with consistent error format."
        - type: param
          name: "Rate Limit"
          datatype: range_int
          min: 100
          max: 1000

# Tracked questions/TBDs
issues:
  - id: tbd_pagination
    description: "Define pagination strategy for list endpoints"
    severity: medium
    resolved: false

# BDD-style acceptance criteria
contracts:
  - name: "Unauthorized access returns 401"
    scenario: "Request without valid token"
    given:
      - "no Authorization header present"
    when:
      - "GET /users/123 is called"
    then:
      - "response status is 401"
      - "response body contains error message"

# CEL-based validation rules
constraints:
  - expr: "version > 0"
    message: "Version must be positive"
    scope: all                    # Applies to all nodes

# Extensible custom fields
custom:
  owner: "platform-team"
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
deco list --tag security     # Filter by tag

deco show <node-id>          # Show node details + reverse references
deco show systems/auth --json   # Output as JSON

deco query auth              # Search title/summary for "auth"
deco query --kind system --status published   # Combined filters

deco validate                # Validate all nodes (schema, refs, constraints)
deco validate --quiet        # Exit code only (for CI)

deco stats                   # Project overview: nodes by kind/status, open issues
deco issues                  # List all open TBDs/issues across nodes
deco issues --severity high  # Filter by severity

deco graph                   # Output dependency graph (DOT format)
deco graph --format mermaid  # Output as Mermaid for Markdown
```

### Modifying Nodes

Edit YAML files directly in `.deco/nodes/`, then run sync to detect changes:

```bash
deco sync                    # Detect edits, bump versions, track history
```

The sync command:
- Detects manual file edits
- Auto-increments version numbers
- Records changes in history
- Normalizes YAML formatting (sorted keys, expanded arrays)

### Review Workflow

```bash
deco review submit <id>      # Submit a draft node for review
deco review approve <id>     # Approve a node (sets status to approved)
deco review reject <id>      # Reject back to draft
deco review status <id>      # Show review status
```

### Audit Trail

```bash
deco history                 # Show all changes
deco history --node systems/auth   # Filter by node
deco history --limit 10      # Limit entries

deco diff <id>               # Show before/after for all changes
deco diff <id> --last 5      # Last 5 changes
deco diff <id> --since 2h    # Changes in last 2 hours
```

### Sync (for manual edits)

```bash
deco sync                    # Detect manual edits, bump version, reset status
deco sync --dry-run          # Show what would change
```

When nodes are edited directly (bypassing CLI), `sync` detects changes by content hash and:
- Bumps the version number
- Resets status to "draft" if it was approved/review
- Logs the sync operation to history

## Validation

Deco validates your documentation graph across three dimensions:

### Schema Validation
Every node must have: `id`, `kind`, `version`, `status`, `title`

```bash
$ deco validate
ERROR [E008] Missing required field: status
  → systems/auth.yaml:1
  Suggestion: Add 'status: draft' to the node
```

### Reference Validation
All references must resolve to existing nodes:

```bash
$ deco validate
ERROR [E020] Reference not found: systems/auht
  → systems/api.yaml:12
  Did you mean: systems/auth?
```

### Constraint Validation
CEL expressions enforce custom rules:

```yaml
constraints:
  - expr: "version > 0"
    message: "Version must be positive"
    scope: all
  - expr: "status in ['draft', 'approved', 'published']"
    message: "Invalid status"
    scope: all
  - expr: "size(tags) > 0"
    message: "At least one tag required"
    scope: requirement  # Only applies to requirement nodes
```

### Block Validation
Blocks use strict field allowlists. Unknown block fields (and table column keys outside `key`, `type`, `enum`, `display`) are validation errors with suggestions.

## Contracts

Define testable acceptance criteria using BDD-style scenarios:

```yaml
contracts:
  - name: "Rate limiting enforced"
    scenario: "Too many requests triggers rate limit"
    given:
      - "client has made 100 requests in 1 minute"
    when:
      - "client makes another request"
    then:
      - "response status is 429"
      - "Retry-After header is present"

  - name: "Token refresh extends session"
    scenario: "Valid refresh token issues new access token"
    given:
      - "user has valid refresh token"
    when:
      - "POST /auth/refresh is called"
    then:
      - "new access token is returned"
      - "refresh token is rotated"
```

Contracts serve as executable specifications — your documentation defines what should happen, and tests verify it actually does.

## Issues: Tracking TBDs

Mark unresolved questions directly in nodes:

```yaml
issues:
  - id: tbd_error_format
    description: "Finalize error response schema"
    severity: high          # low, medium, high, critical
    location: "content.sections.errors"
    resolved: false
```

Documentation is complete when all issues are resolved, all references exist, and all constraints pass.

## AI Integration

Deco is designed for AI-assisted documentation workflows. The engine validates all changes, so AI can propose updates without breaking consistency.

### Two Update Modes

**Patch Mode** — Surgical changes via operations:
```json
[
  {"op": "set", "path": "status", "value": "approved"},
  {"op": "append", "path": "tags", "value": "reviewed"},
  {"op": "unset", "path": "issues[0]"}
]
```

**Rewrite Mode** — AI generates complete YAML, Deco validates before saving.

### LLM Context

Add context for AI assistants directly in the node:

```yaml
llm_context: "This is the core authentication system. Changes affect all API security."
```

### Workflow Example

1. Author describes intent in natural language
2. AI generates structured YAML or patch operations
3. Deco validates schema, references, and constraints
4. Changes applied only if valid
5. Audit log records who/what/when

The engine — not the AI — is the source of truth.

## Typical Workflow

```bash
# Start a new documentation project
deco init my-api && cd my-api

# Create nodes by editing YAML files in .deco/nodes/
# (manually or via AI-generated YAML)

# Sync changes to track history
deco sync

# Validate as you go
deco validate

# Explore your documentation
deco list --kind system
deco show systems/auth
deco query "endpoint"

# Edit YAML files directly, then sync
deco sync

# Review history
deco history --node systems/auth

# CI integration
deco validate --quiet && echo "Documentation valid"
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
│   └── query/            # Filtering and search
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

Deco is licensed under the [GNU Affero General Public License v3.0](LICENSE) (AGPL-3.0-or-later).

You are free to use, modify, and distribute this software under the terms of the AGPL. If you modify Deco and offer it as a network service, you must make your source code available to users of that service.

A commercial license is available for organizations that cannot comply with AGPL requirements. See [NOTICE](NOTICE) for details.
