# Relational Blocks: Union Refs, List Queries, and Cross-Type Joins

Date: 2026-02-07

## Problem

Deco's block data model is flat when actual game data is relational. Three specific issues:

1. **Rigid refs**: `ref` only supports a single `{block_type, field}` target. Real data has union relationships (building.materials references both resource.name AND recipe.output).
2. **No list membership in queries**: `--field materials=Planks` doesn't work when materials is a list — it stringifies the whole list and compares.
3. **No cross-type joins**: Can query one block type at a time, but can't ask "what recipes produce materials used by bronze-age buildings?"

## Feature 1: Union Refs

### Config

`ref:` accepts both single object (backward compat) and array (new):

```yaml
# Existing syntax still works
materials:
  type: list
  ref:
    block_type: resource
    field: name

# New: array of targets (OR logic)
materials:
  type: list
  ref:
    - block_type: resource
      field: name
    - block_type: recipe
      field: output
```

### Go Model

Replace `Ref *RefConstraint` with `Refs []RefConstraint` on `FieldDef`. Custom `UnmarshalYAML` detects whether `ref:` is an object or array and normalizes to the slice form. Old single-object form becomes a one-element slice internally.

### Validation

`CrossRefValidator` builds ref sets for all targets, then validates with OR: a value is valid if it exists in **any** of the target ref sets. Error message lists all checked targets. Suggestions drawn from the union of all valid values.

### Backward Compat

JSON/YAML output always writes the array form. Old single-object configs parse correctly via the custom unmarshaler. No migration needed.

## Feature 2: List Membership in Queries

### Current Behavior

`matchesBlockCriteria` does `fmt.Sprintf("%v", fieldVal) != value`, which stringifies the whole value. A list `[Planks, Stone]` becomes `"[Planks Stone]"` and never matches `"Planks"`.

### Fix

Add a `matchesFieldValue` helper with type switch:

- **string**: exact match (same as today)
- **number**: string comparison via `fmt.Sprintf` (same as today)
- **[]interface{} (list)**: returns true if any element matches the value (contains semantics)
- **fallback**: stringify comparison (same as today)

`--field materials=Planks` matches a block where materials is `[Stone, Planks, Bronze Ingots]`.

AND logic across multiple `--field` flags preserved. No new CLI flags needed.

## Feature 3: Cross-Type Join Queries (`--follow`)

### CLI Syntax

```bash
# Auto mode: uses ref config to find targets
deco query --block-type building --field age=bronze --follow materials

# Explicit override: user specifies target
deco query --block-type building --field age=bronze --follow materials:recipe.output
```

### Execution Flow

1. Run the normal block query (filter by block type + field filters) -> get matched blocks
2. Extract the followed field's values from all matched blocks (expanding lists element-by-element)
3. Resolve targets: from ref config (`Refs []RefConstraint`) or from explicit `field:blocktype.field` override
4. For each target, find all blocks of that block type where the target field matches any extracted values
5. Group results by followed value, deduplicate, count referencing source blocks

### Output Format

```
Planks (referenced by 8 building blocks)
  recipe in systems/settlement/recipes > Crafting Recipes > block 2
    name: Plank Making
    inputs: [Logs]
    output: Planks
    tier: T1

Bronze Ingots (referenced by 3 building blocks)
  recipe in systems/settlement/recipes > Crafting Recipes > block 5
    name: Bronze Smelting
    inputs: [Copper, Tin]
    output: Bronze Ingots
    tier: T1
```

When a followed value matches across multiple target block types, both show under the same value header:

```
Stone (referenced by 12 building blocks)
  resource in systems/settlement/resources > Raw Resources > block 0
    name: Stone
    tier: T0
  recipe in systems/economy/recipes > Processing > block 3
    name: Stone Cutting
    inputs: [Raw Stone]
    output: Stone
    tier: T1
```

### Symmetry

`--follow` works from any starting block type. Direction comes from which block type you start with:

- `--block-type building --follow materials` -> finds recipes/resources that provide those materials
- `--block-type recipe --follow inputs` -> finds resources that provide those inputs

### Edge Cases

- Followed field doesn't exist on matched blocks -> error: `field "foo" not found in matched blocks`
- No ref config and no explicit target -> error: `field "materials" has no ref constraint; use --follow materials:blocktype.field`
- Followed value not found in any target -> listed with `(no matches found)`
- `--follow` without `--block-type` -> error, block-type is required for follow
- No chaining in v1 (one hop only). Multi-hop can come later.
- No reverse follows ("what references me") -- handled by list membership queries (`--field materials=Planks`).

## Implementation Scope

### Files Changed

**Config layer** (`internal/storage/config/repository.go`):
- Add `Refs []RefConstraint` to `FieldDef`, remove `Ref *RefConstraint`
- Custom `UnmarshalYAML` / `MarshalYAML` on `FieldDef`

**Cross-ref validator** (`internal/services/validator/crossref_validator.go`):
- Loop over `fieldDef.Refs` instead of single `fieldDef.Ref`
- Build union of valid values, OR logic
- Error message lists all checked targets

**Query engine** (`internal/services/query/query.go`):
- `matchesBlockCriteria`: type switch for list membership
- New `FollowField` / `FollowTarget` on `FilterCriteria`
- New `FollowBlocks` method and `FollowResult` type

**Query CLI** (`internal/cli/query.go`):
- New `--follow` flag
- Output formatter for grouped follow results

**Tests**: Update existing cross-ref and query tests. New tests for union refs, list membership, and follow queries.

### Not Touched

Node-level refs, graph, sync, review -- none of these change.
