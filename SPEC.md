# Deco Specification

## What It Is

A Go CLI for managing complex documentation as structured, validated YAML. Your specifications, designs, and requirements live as interconnected nodes with schema validation, reference tracking, and audit trails. Works with git, works with AI.

## Use Cases

| Domain | Example Nodes |
|--------|---------------|
| Game Design | systems, mechanics, items, characters, quests |
| API Specs | endpoints, schemas, authentication, versioning |
| Technical Architecture | components, interfaces, services, data flows |
| Product Requirements | features, user stories, epics, acceptance criteria |
| Knowledge Bases | concepts, definitions, relationships |

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
- `kind`: Node type (system, component, feature, requirement, etc.)
- `version`: Auto-incremented on updates
- `status`: draft, review, approved, published, deprecated
- `title`: Human-readable name

References (`refs`):
- `uses`: Hard dependencies with context
- `related`: Informational links
- `emits_events`: Events this node produces
- `vocabulary`: Shared term definitions

Content (`content.sections`):
- Blocks of type: table, rule, param, mechanic, list
- Custom block types can be defined in project config

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

**No ambiguity in complete documentation.** A spec is "complete" when:
- Zero open `issues` (no TBDs, no unanswered questions)
- All refs resolve to existing nodes
- All required fields populated

**No contradictions.** Deco validates internal consistency through:
- Enum values must match their definitions
- Explicit constraints you define are enforced across nodes

Example constraint:
```yaml
constraints:
  - expr: "self.rate_limit == refs['systems/api'].default_rate_limit"
    message: "Rate limit must match API default"
```

Deco doesn't magically detect semantic contradictions - you declare what must be consistent, and Deco enforces it.

## What Deco Does NOT Do

- Gherkin compilation or test execution
- UI or editor plugins (CLI only)
- Runtime integration (it's a documentation tool)

## Project Layout

```
.deco/
  config.yaml          # Project configuration
  history.jsonl        # Audit log
  nodes/
    systems/
      auth/
        core.yaml
        tokens.yaml
    components/
      database.yaml
```

## Project Configuration

The `.deco/config.yaml` file supports:

```yaml
project_name: my-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
version: 1
required_approvals: 2  # For review workflow

# Define custom block types with validation
custom_block_types:
  endpoint:
    required_fields:
      - method
      - path
      - response

# Define per-kind schema rules for nodes
schema_rules:
  requirement:
    required_fields:
      - priority
      - acceptance_criteria
  component:
    required_fields:
      - owner
      - dependencies
```

Custom block types extend the built-in types (rule, table, param, mechanic, list). When a custom type shares a name with a built-in type, both validations apply.

Schema rules enforce required custom fields per node kind. The `required_fields` must be present in the node's `custom:` section. Nodes with kinds not listed in schema_rules are not constrained.

## AI Integration

Two modes:
1. **Patch mode**: AI outputs JSON patch operations, Deco validates and applies
2. **Rewrite mode**: AI rewrites full YAML files, Deco validates after

The engine is the source of truth, not the AI.
