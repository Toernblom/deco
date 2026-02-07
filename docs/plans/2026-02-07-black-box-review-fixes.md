# Black-Box Review Fixes Design

Date: 2026-02-07

## Context

A black-box-only CLI review (no repo inspection) found 6 issues ranked by severity.
All 6 confirmed via codebase analysis. This plan addresses all of them.

## Canonical Status Lifecycle

```
draft → review → approved → deprecated → archived
```

"published" is removed from the codebase entirely.

---

## Fix 1: Sync Silent Output (High)

**Problem**: `deco sync` exits silently (exit 0) when no changes exist.

**Change**: In `runSync()`, when `!flags.quiet` and no changes were made, print:
```
All 12 nodes clean, no changes.
```

**Files**: `internal/cli/sync.go`

---

## Fix 2: Stats Misleading Health (High)

**Problem**: `deco stats` only counts E041 constraint violations, ignoring all other
validation errors. Reports "No violations" while `deco validate` finds errors.

**Change**:
- Replace `constraintViolations int` with `totalValidationErrors int` and
  `validationByCategory map[string]int` in `projectStats`
- Use `NewOrchestratorWithFullConfig()` to match what `deco validate` checks
- Count ALL collector errors grouped by category (schema, refs, validation)
- Rename section from "CONSTRAINT VIOLATIONS" to "VALIDATION HEALTH"
- Display total + per-category breakdown

**Files**: `internal/cli/stats.go`, `internal/cli/stats_test.go`

---

## Fix 3: Filter Validation (Medium)

**Problem**: Invalid `--status`, `--severity`, `--kind` values silently return empty
results. `--field` without `=` is silently ignored.

**Change**:
- New `internal/cli/filter_validation.go` with:
  - `validateStatus(status string) error` — hardcoded valid set
  - `validateSeverity(sev string) error` — low/medium/high/critical
  - `validateKind(kind string, validKinds []string) error` — schema-derived
  - `validateFieldFilter(field string) error` — requires key=value format
- Each returns user-friendly error with "did you mean?" via string distance
- Validation called in `RunE` of `list.go`, `query.go`, `issues.go` after flag
  parsing, before service calls
- For kind: extract valid kinds from loaded project schema

**Files**: new `internal/cli/filter_validation.go`, `list.go`, `query.go`, `issues.go`,
plus test files

---

## Fix 4: Remove "published" Status (Medium)

**Problem**: "published" used in code/tests/help but not in `validStatuses` map.

**Change**:
- `validator.go`: ContentValidator checks only `status != "approved"`
- `error_codes.go`: E046 message → "Content required for approved status"
- `style.go`: Remove `StatusPublished`, add colors for `review` (cyan) and
  `approved` (green)
- `help.go`: Example shows `draft, review, approved, deprecated, archived`
- `list.go`, `query.go`: Help text updated
- Tests (~15 files): Replace `"published"` with `"approved"` in test data

**Files**: `validator.go`, `validator_test.go`, `error_codes.go`, `style.go`,
`help.go`, `list.go`, `query.go`, `list_test.go`, `query_test.go`, `show_test.go`,
`validate_test.go`

---

## Fix 5: Diff Exit Code for Non-Existent Nodes (Medium)

**Problem**: `deco diff nope` returns exit 0 with "No changes found" — same as a real
node with no history.

**Change**: Before querying history, load the node from the repository. If not found,
return error (exit 1): `Error: node 'nope' not found`. If node exists but no history,
keep current message.

**Files**: `internal/cli/diff.go`, `internal/cli/diff_test.go`

---

## Fix 6: Quiet Flag Consistency (Low)

**Problem**: `--quiet` behavior varies. Global flag exists but isn't wired. `list` and
`query` lack it entirely.

**Change**:
- Wire global `--quiet` from `root.go` so commands respect it
- Add `--quiet` to `list` and `query`: output node IDs only, one per line
- Leave `stats` and `diff` without `--quiet` (add later if needed)

Standard quiet behavior:
- `validate`: suppress output, exit code only (existing)
- `sync`: suppress output, exit code only (existing)
- `issues`: counts only (existing)
- `list`: node IDs only (new)
- `query`: node IDs only (new)

**Files**: `internal/cli/root.go`, `list.go`, `query.go`, plus tests
