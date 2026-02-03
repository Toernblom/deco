# Deco CLI Reference

Complete command reference for the Deco CLI.

## Global Flags

| Flag | Description |
|------|-------------|
| `-c, --config` | Path to deco project directory (default `.deco`) |
| `-q, --quiet` | Suppress non-error output |
| `--verbose` | Enable verbose output |
| `-v, --version` | Show version |

---

## Project Setup

### `deco init`

Initialize a new Deco project.

```bash
deco init [directory]
deco init .                    # Initialize in current directory
deco init my-project           # Create and initialize new directory
deco init . --force            # Reinitialize existing project
```

### `deco create`

Create a new node with required fields.

```bash
deco create <node-id>
deco create systems/auth
deco create systems/auth --kind system --title "Authentication"
```

| Flag | Description |
|------|-------------|
| `--kind` | Node kind (default: inferred from path) |
| `--title` | Node title |

---

## Reading & Querying

### `deco list`

List all nodes in the project.

```bash
deco list
deco list --kind system        # Filter by kind
deco list --status draft       # Filter by status
deco list --tag security       # Filter by tag
```

| Flag | Description |
|------|-------------|
| `--kind` | Filter by node kind |
| `--status` | Filter by status |
| `--tag` | Filter by tag |

### `deco show`

Show detailed information about a node.

```bash
deco show <node-id>
deco show systems/auth
deco show systems/auth --json  # Output as JSON
```

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON |

### `deco query`

Search and filter nodes by text.

```bash
deco query <term>
deco query auth                # Search title/summary for "auth"
deco query --kind system       # Filter by kind
deco query --status approved   # Filter by status
```

### `deco validate`

Validate all nodes (schema, references, constraints, blocks).

```bash
deco validate
deco validate [directory]
deco validate --quiet          # Exit code only (for CI)
```

Returns exit code 0 if valid, non-zero if errors found.

### `deco stats`

Show project overview and health statistics.

```bash
deco stats
deco stats [directory]
```

Shows: node counts by kind/status, open issues, reference statistics.

### `deco issues`

List all open issues/TBDs across the design.

```bash
deco issues
deco issues --severity high    # Filter by severity
```

### `deco graph`

Output dependency graph.

```bash
deco graph                     # DOT format (default)
deco graph --format mermaid    # Mermaid format for Markdown
deco graph --format dot        # Graphviz DOT format
```

---

## Modifying Nodes

### `deco set`

Set a field value on a node.

```bash
deco set <node-id> <path> <value>
deco set systems/auth title "New Title"
deco set systems/auth status approved
deco set systems/auth tags[0] security
```

Supports dot notation and array indexing.

### `deco append`

Append a value to an array field.

```bash
deco append <node-id> <path> <value>
deco append systems/auth tags oauth
deco append systems/auth refs.uses.target items/token
```

### `deco unset`

Remove a field or array element.

```bash
deco unset <node-id> <path>
deco unset systems/auth summary
deco unset systems/auth tags[2]
```

### `deco apply`

Apply a batch of patch operations from a JSON file.

```bash
deco apply <node-id> <patch-file>
deco apply systems/auth changes.json
deco apply systems/auth changes.json --dry-run
```

Patch file format:
```json
[
  {"op": "set", "path": "title", "value": "New Title"},
  {"op": "append", "path": "tags", "value": "new-tag"},
  {"op": "unset", "path": "summary"}
]
```

| Flag | Description |
|------|-------------|
| `--dry-run` | Show what would change without applying |

### `deco rewrite`

Replace entire node content from a YAML file.

```bash
deco rewrite <node-id> <yaml-file>
deco rewrite systems/auth new-content.yaml
```

Replaces the node content completely. Validates before saving.

### `deco rm`

Delete a node.

```bash
deco rm <node-id>
deco rm old-spec               # Fails if other nodes reference it
deco rm old-spec --force       # Delete even with references
```

| Flag | Description |
|------|-------------|
| `--force` | Delete even if other nodes reference this node |

### `deco mv`

Rename a node and update all references.

```bash
deco mv <old-id> <new-id>
deco mv systems/auth systems/authentication
```

Automatically updates all nodes that reference the renamed node.

---

## Review Workflow

Nodes progress through status: `draft` → `review` → `approved`

### `deco review submit`

Submit a draft node for review.

```bash
deco review submit <node-id>
deco review submit systems/auth
```

Changes status from `draft` to `review`. Node must be in `draft` status.

### `deco review approve`

Approve a node under review.

```bash
deco review approve <node-id>
deco review approve systems/auth
deco review approve systems/auth --note "LGTM"
```

| Flag | Description |
|------|-------------|
| `--note` | Optional approval note |

Adds your approval. Status changes to `approved` when `required_approvals` threshold is met (configured in `.deco/config.yaml`).

### `deco review reject`

Reject a node back to draft.

```bash
deco review reject <node-id> --note "Needs more detail"
```

| Flag | Description |
|------|-------------|
| `--note` | Rejection reason (required) |

Changes status from `review` back to `draft`.

### `deco review status`

Show review status of a node.

```bash
deco review status <node-id>
deco review status systems/auth
```

Shows: current status, version, approval count, reviewer list.

---

## History & Audit

### `deco history`

Show audit log history.

```bash
deco history
deco history --node systems/auth   # Filter by node
deco history --limit 10            # Limit entries
```

| Flag | Description |
|------|-------------|
| `--node` | Filter by node ID |
| `--limit` | Maximum entries to show |

### `deco diff`

Show changes to a node over time.

```bash
deco diff <node-id>
deco diff systems/auth
deco diff systems/auth --last 5    # Last 5 changes
deco diff systems/auth --since 2h  # Changes in last 2 hours
```

| Flag | Description |
|------|-------------|
| `--last` | Show last N changes |
| `--since` | Show changes since duration (e.g., `2h`, `1d`) |

### `deco sync`

Detect and fix unversioned node changes.

```bash
deco sync
deco sync --dry-run            # Show what would change
```

When nodes are edited directly (bypassing CLI), `sync` detects changes by content hash and:
- Bumps the version number
- Resets status to `draft` if it was `approved` or `review`
- Logs the sync operation to history

---

## Migration

### `deco migrate`

Migrate nodes to current schema version.

```bash
deco migrate
deco migrate [directory]
```

Updates nodes to match the current schema structure.

---

## Configuration

Project configuration in `.deco/config.yaml`:

```yaml
project_name: my-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
version: 1
required_approvals: 2          # Approvals needed for review → approved

# Custom block types
custom_block_types:
  endpoint:
    required_fields:
      - method
      - path
    optional_fields:
      - auth

# Per-kind schema rules
schema_rules:
  requirement:
    required_fields:
      - priority
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Validation errors or command failure |

---

## Examples

### Typical Workflow

```bash
# Start a new project
deco init my-docs && cd my-docs

# Create nodes
deco create systems/auth --kind system --title "Authentication"

# Edit the YAML file directly, then sync
deco sync

# Validate
deco validate

# Submit for review
deco review submit systems/auth

# Approve (repeat until required_approvals met)
deco review approve systems/auth --note "LGTM"

# Check status
deco review status systems/auth
```

### CI Integration

```bash
# Validate all nodes (quiet mode for CI)
deco validate --quiet && echo "Documentation valid"
```

### Refactoring

```bash
# Rename a node (updates all references)
deco mv components/db components/database

# Check what changed
deco history --limit 5
```
