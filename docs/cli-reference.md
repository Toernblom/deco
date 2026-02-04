# Deco CLI Reference

Complete command reference for the Deco CLI.

## Workflow Overview

Deco is designed for LLM-assisted workflows where YAML files are edited directly:

```
Edit YAML → deco sync → deco validate → deco review approve
```

## Global Flags

| Flag | Description |
|------|-------------|
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

## Sync & Change Detection

### `deco sync`

Detect and process changes made directly to YAML files.

```bash
deco sync
deco sync --dry-run            # Show what would change
deco sync --no-refactor        # Skip automatic reference updates
```

| Flag | Description |
|------|-------------|
| `--dry-run` | Show what would change without applying |
| `--no-refactor` | Skip automatic reference updates for renames |
| `-q, --quiet` | Suppress output |

**What sync detects:**

1. **Modified nodes**: Content changed since last sync
   - Bumps version number
   - Resets status to `draft` if was `approved` or `review`
   - Clears reviewers
   - Logs to history

2. **New nodes**: Files with no history
   - Baselines current state (records in history)

3. **Renamed nodes**: File moved/renamed manually
   - Detects via content hash matching
   - Automatically updates references in other nodes
   - Logs move operation to history

4. **Deleted nodes**: Files removed
   - Logs deletion to history

**Exit codes:**
- `0` - No changes needed
- `1` - Files modified (re-commit needed)
- `2` - Error occurred

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

# Create node by writing YAML directly
cat > .deco/nodes/systems/auth.yaml << 'EOF'
id: systems/auth
kind: system
version: 1
status: draft
title: Authentication System
summary: Handles user login and session management
EOF

# Sync to register the new node
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

### Renaming Nodes

```bash
# Rename by moving the file
mv .deco/nodes/components/db.yaml .deco/nodes/components/database.yaml

# Edit the file to update the id field
# Then sync - references are updated automatically
deco sync

# Check what changed
deco history --limit 5
```

### Deleting Nodes

```bash
# Delete by removing the file
rm .deco/nodes/old-spec.yaml

# Sync to record the deletion
deco sync
```

### CI Integration

```bash
# Validate all nodes (quiet mode for CI)
deco validate --quiet && echo "Documentation valid"
```
