# Design: Content Hash-Based Sync Detection

**Date:** 2026-02-01
**Status:** Approved

## Problem

`deco sync` currently relies on `git diff HEAD` to detect unversioned changes. This only catches uncommitted changes. If someone commits without the pre-commit hook running (or the hook isn't attached), the change slips through undetected.

## Solution

Store content hashes in history.jsonl. Sync compares current content hash against the last recorded hash to detect drift, regardless of git state.

## Content Hashing

**Fields hashed:** `title`, `summary`, `tags`, `refs`, `issues`, `content`
**Algorithm:** SHA-256, truncated to 16 hex characters
**Serialization:** Canonical YAML of content fields only

**Storage in history.jsonl:**
```json
{"timestamp":"...","node_id":"items/food","operation":"set","user":"anton","content_hash":"a1b2c3d4e5f67890","before":{...},"after":{...}}
```

## Sync Detection Logic

1. Load all nodes in the project
2. For each node, compute current content hash
3. Query history for the most recent entry for that node
4. Compare:
   - **Hash matches** → Clean, skip
   - **Hash differs** → Content changed outside CLI, needs sync
   - **No history entry** → Auto-baseline: record current hash with `baseline` operation, skip

## Baseline Operation

New audit operation type that records state without implying a change:
```json
{"timestamp":"...","node_id":"items/food","operation":"baseline","user":"anton","content_hash":"a1b2c3d4e5f67890"}
```

Used when a node has no prior hash recorded (new project, pre-feature nodes, cleared history).

## Implementation Changes

### Files to modify

1. **`internal/domain/audit.go`**
   - Add `ContentHash string` field to `AuditEntry` struct
   - Add `baseline` to valid operations list

2. **`internal/cli/sync.go`**
   - Remove git dependency (no more `git diff`, `git show`, `isGitRepo`)
   - New logic: load all nodes → compute hashes → compare against history
   - Add `computeContentHash(node)` function
   - Handle baseline case

3. **CLI commands that modify nodes** (set, append, unset, apply, create, mv)
   - Compute and include content hash when logging to history

### Removed code

- `getModifiedNodeFiles()` - no longer needed
- `getNodeFromHEAD()` - no longer needed
- `isGitRepo()` check - sync works without git

### New helper

```go
func computeContentHash(n domain.Node) string {
    // Serialize content fields to canonical YAML
    // Return SHA-256 truncated to 16 hex chars
}
```

## Output Format

```
$ deco sync
Baselined: items/food, systems/core (2 nodes)
Synced: items/weapon (v1→v2)

$ deco sync --dry-run
Would baseline: items/food, systems/core (2 nodes)
Would sync: items/weapon (v1→v2)
```

## Exit Codes

- 0 = clean (nothing to do, or baseline-only)
- 1 = files modified (sync happened)
- 2 = error

## Edge Cases

| Case | Behavior |
|------|----------|
| New project, no history | All nodes baselined on first run |
| Node deleted | Skipped (can't load) |
| Node created outside CLI | No history, gets baselined |
| Node deleted and recreated | Hash mismatch, sync bumps version |
