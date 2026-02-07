# Deco

**Structured documentation for complex systems.**

Deco is a CLI that treats documentation like source code: typed YAML nodes, explicit references between them, schema validation, and a full audit trail. No more scattered TBDs, broken cross-references, or outdated specs.

## Why Deco?

If you've maintained complex documentation — system specs, API designs, game design docs, product requirements — you know how it ends:

- TBDs everywhere, contradictory rules, outdated pages
- Updating one component means hunting references across 30 documents
- An LLM helps write content, then loses context, and now you have two truths

Deco fixes this:

| Problem | Deco's answer |
|---|---|
| Broken cross-references | Reference validation with typo suggestions |
| Outdated specs | Sync detection, version bumping, status workflow |
| No accountability | Append-only audit log with who/what/when |
| AI-generated drift | Engine validates all changes — AI proposes, Deco enforces |

## Quick Start

```bash
# Install
go install github.com/Toernblom/deco/cmd/deco@latest

# Create a project
deco init my-project && cd my-project

# Scaffold a node
deco new systems/auth --kind system --title "Authentication System"

# Edit the generated YAML in .deco/nodes/systems/auth.yaml, then:
deco sync              # Detect changes, bump version, record history
deco validate          # Check schema + references + constraints
deco show systems/auth # View node with reverse references
```

## Installation

```bash
# Go install (requires Go 1.25+)
go install github.com/Toernblom/deco/cmd/deco@latest
```

**From source:**

```bash
git clone https://github.com/Toernblom/deco.git
cd deco
go build -o deco ./cmd/deco
```

## How It Works

Documentation lives as YAML files in `.deco/nodes/`. Each file is a **node** with required metadata, optional references to other nodes, structured content, and tracked issues.

```
.deco/
├── config.yaml        # Project configuration
├── history.jsonl      # Append-only audit log
└── nodes/
    ├── systems/
    │   ├── auth.yaml
    │   └── api.yaml
    ├── components/
    │   └── database.yaml
    └── requirements/
        └── user-stories.yaml
```

Node IDs map to file paths: `systems/auth` → `.deco/nodes/systems/auth.yaml`

### Node Format

```yaml
id: systems/auth
kind: system
version: 1
status: draft                        # draft → review → approved → deprecated → archived
title: "Authentication System"
tags: [core, security]
summary: "JWT-based authentication with refresh tokens."

refs:
  uses:
    - target: systems/users
      context: "Authenticates user identities"
  related:
    - target: components/database

content:
  sections:
    - name: overview
      blocks:
        - type: rule
          text: "All API endpoints require valid JWT."
        - type: param
          name: "Token TTL"
          datatype: duration
          value: "15m"

issues:
  - id: tbd_refresh
    description: "Define refresh token rotation policy"
    severity: medium
    resolved: false

contracts:
  - name: "Unauthorized access returns 401"
    given: ["no Authorization header present"]
    when: ["GET /users/123 is called"]
    then: ["response status is 401"]
```

See [docs/SPEC.md](docs/SPEC.md) for the full node format, including constraints, custom fields, and LLM context.

## Validation

Deco validates across three dimensions:

**Schema** — Required fields, valid block types, field allowlists:
```
ERROR [E008] Missing required field: status
  → systems/auth.yaml:1
  Suggestion: Add 'status: draft' to the node
```

**References** — All refs must resolve, with typo suggestions:
```
ERROR [E020] Reference not found: systems/auht
  → systems/api.yaml:12
  Did you mean: systems/auth?
```

**Constraints** — Custom CEL expressions enforce project rules:
```yaml
constraints:
  - expr: "size(tags) > 0"
    message: "At least one tag required"
    scope: requirement
```

## CLI Reference

### Setup

```bash
deco init [directory]                # Initialize a new project
deco new <id> --kind <k> --title <t> # Scaffold a node
deco migrate                         # Migrate nodes to current schema
```

### Query

```bash
deco list                            # List all nodes
deco list --kind system --status draft --tag security
deco show <id>                       # Node details + reverse references
deco show <id> --full                # Expand content blocks inline
deco query <text>                    # Search titles and summaries
deco stats                           # Project health overview
deco issues                          # Open TBDs across all nodes
deco graph                           # Dependency graph (DOT format)
deco graph --format mermaid          # Mermaid for Markdown embedding
```

### Modify

```bash
deco sync                            # Detect edits, bump versions, track history
deco sync --dry-run                  # Preview without changes
```

### Review

```bash
deco review submit <id>              # Submit for review
deco review approve <id>             # Approve (draft → approved)
deco review reject <id>              # Reject back to draft
deco review status                   # List nodes in review
```

### History

```bash
deco history                         # Full audit log
deco history --node <id>             # Filter by node
deco diff <id>                       # Before/after for all changes
deco diff <id> --since 2h            # Changes in the last 2 hours
```

### Export

```bash
deco export <id>                     # Single node to stdout (markdown)
deco export --output docs/           # All nodes → one .md per node
```

## Use Cases

- **Game Design Documents** — Systems, mechanics, items with interconnected rules
- **API Specifications** — Endpoints, schemas, versioning with dependency tracking
- **Technical Architecture** — Components, interfaces, data flows with validation
- **Product Requirements** — Features, user stories, acceptance criteria with status workflow
- **Knowledge Bases** — Interconnected concepts with reference integrity

## Design

```
cmd/deco/          CLI entry point
internal/
├── domain/        Core types: Node, Ref, Graph, AuditEntry
├── storage/       Persistence (YAML nodes, JSONL history, YAML config)
├── services/      Graph analysis, validation, query engine
└── cli/           Command implementations
```

Key decisions:
- **YAML nodes** — human-readable, diffable, merge-friendly
- **JSONL history** — append-only, streamable
- **CEL constraints** — Google's Common Expression Language for custom rules
- **No database** — everything is files, works natively with git

## Documentation

- [Full Specification](docs/SPEC.md) — node format, validation rules, CLI details, and design decisions
- [Examples](examples/) — sample projects (game design, API spec)

## License

[AGPL-3.0-or-later](LICENSE). A commercial license is available for organizations that cannot comply with AGPL requirements — see [NOTICE](NOTICE) for details.
