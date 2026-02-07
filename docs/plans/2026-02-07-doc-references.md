# Doc References: External Markdown Files in Deco

## Problem

YAML is great for structured data but terrible for writing prose. Narrative design, user stories, lore, and long-form content are awkward as YAML strings. Authors need to write in proper `.md` files while keeping deco's validation and reference tracking.

## Design

### Data Model

Two levels of doc references:

**Node-level** - new top-level `docs` field:
```yaml
id: stories/chapter-1
title: The Beginning
docs:
  - path: narratives/chapter-1.md
    keywords: [protagonist, ancient temple, betrayal]
    context: Main narrative for chapter 1
```

**Block-level** - new built-in block type `doc`:
```yaml
content:
  sections:
    - name: Narrative
      blocks:
        - type: doc
          path: narratives/opening-scene.md
          keywords: [storm, arrival, mysterious stranger]
          context: Opening scene description
```

Both share the same fields:
- `path` (required) - relative to project root, must be `.md` file
- `keywords` (optional) - list of terms that must appear in the file
- `context` (optional) - describes what this document is for

### Validation

Two new error codes:

- **E055: Doc file not found** - referenced `.md` file doesn't exist at the specified path
- **E056: Missing keyword in doc** - keyword not found in the `.md` file content

Behavior:
- If `keywords` is omitted or empty, only file existence is checked
- Keywords are plain strings, case-insensitive substring matching
- Each missing keyword produces a separate E056 error naming the specific keyword
- Both node-level `docs` and block-level `doc` blocks use the same validation

### Sync & Change Tracking

- `deco sync` includes referenced `.md` file contents in the node's content hash
- If an `.md` file changes (but YAML doesn't), the hash still changes → version bumps, status resets to `draft`, reviewers cleared
- Missing `.md` files at sync time produce a warning (validation catches the hard error)
- History entries record which `.md` files changed (by path)
- `deco diff` flags doc files as modified (does not show full markdown diff)

**Review feedback loop:**
1. Author edits `.md` file
2. `deco sync` detects hash change → bumps version, resets to `draft`
3. `deco validate` may flag missing keywords if prose drifted
4. Author updates `keywords` in YAML if needed
5. `deco review submit` → reviewer checks both YAML and prose

### Show Integration

`deco show <id>` displays doc references:
```
stories/chapter-1 (v3, draft)
  Title: The Beginning
  Docs:
    - narratives/chapter-1.md (keywords: protagonist, ancient temple, betrayal)
```

### Query

No special query support for doc content initially. Existing node-level queries (kind, tag, status, search) work as before.

## Implementation Plan

### 1. Domain model changes
- Add `DocRef` struct to `internal/domain/node.go` with `Path`, `Keywords`, `Context` fields
- Add `Docs []DocRef` field to `Node` struct
- Add `doc` to built-in block types with required field `path`, optional fields `keywords`, `context`
- Register error codes E055 and E056 in `error_codes.go`

### 2. YAML storage
- Update `internal/storage/node/repository.go` to serialize/deserialize the `docs` field
- Add `docs` to the known fields list (avoid E010 unknown field errors)

### 3. Doc validator
- New file `internal/services/validator/doc_validator.go`
- `ValidateDocs(node, projectRoot)` - validates node-level docs
- `ValidateDocBlock(block, projectRoot)` - validates doc blocks
- File existence check (E055)
- Keyword substring matching, case-insensitive (E056)
- Integrate into validator orchestrator

### 4. Block validator updates
- Register `doc` as built-in block type in `block_validator.go`
- Required fields: `path`
- Optional fields: `keywords`, `context`

### 5. Sync changes
- Update content hash computation to include `.md` file contents
- Hash each referenced file (from both `docs` and `doc` blocks)
- If file missing, skip it in hash (validation handles the error)
- Include changed doc paths in history entries

### 6. Show command
- Update `internal/cli/show.go` to display `docs` field
- Show path and keywords for each doc reference

### 7. LLM help
- Update `internal/cli/help.go` with docs field and doc block type documentation

### 8. SPEC.md
- Document the `docs` field and `doc` block type
- Document E055 and E056 error codes
