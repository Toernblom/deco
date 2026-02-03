# Deco YAML Reference for LLMs

This document provides the authoritative reference for writing valid Deco YAML. All fields, block types, and validation rules are documented here.

## Node Structure

Every node file is a YAML document with the following structure.

### Required Fields

These fields **must** be present in every node:

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier, maps to file path (e.g., `systems/auth`) |
| `kind` | string | Node type (e.g., `system`, `component`, `feature`, `item`, `requirement`) |
| `version` | integer | Version number, must be > 0, auto-incremented on updates |
| `status` | string | One of: `draft`, `review`, `approved`, `published`, `deprecated` |
| `title` | string | Human-readable name |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `tags` | string[] | List of tags for categorization |
| `summary` | string | Brief description (can use multiline `\|` syntax) |
| `refs` | object | References to other nodes (see References section) |
| `content` | object | Structured content sections (see Content section) |
| `issues` | object[] | Tracked TBDs and questions (see Issues section) |
| `glossary` | map[string]string | Term definitions local to this node |
| `contracts` | object[] | BDD-style acceptance criteria (see Contracts section) |
| `constraints` | object[] | CEL validation expressions (see Constraints section) |
| `llm_context` | string | Additional context for AI assistants |
| `reviewers` | object[] | Approval records (managed by review workflow) |
| `custom` | map[string]any | Extensible custom fields |

### Minimal Valid Node

```yaml
id: systems/example
kind: system
version: 1
status: draft
title: Example System
```

### Complete Node Example

```yaml
id: systems/auth
kind: system
version: 3
status: approved
title: Authentication System
tags:
  - core
  - security

summary: |
  JWT-based authentication with refresh tokens.
  Handles login, logout, and session management.

refs:
  uses:
    - target: systems/users
      context: Authenticates user identities
  related:
    - target: systems/api
  emits_events:
    - events/auth/login
  vocabulary:
    - glossaries/security

content:
  sections:
    - name: Token Management
      blocks:
        - type: param
          name: Access Token TTL
          datatype: duration
          default: 15m
          description: Time before access token expires

glossary:
  jwt: JSON Web Token used for stateless authentication
  refresh_token: Long-lived token used to obtain new access tokens

contracts:
  - name: Valid login returns tokens
    scenario: User logs in with correct credentials
    given:
      - user exists with email user@example.com
      - password is correct
    when:
      - POST /auth/login is called
    then:
      - response contains access_token
      - response contains refresh_token

llm_context: This is the core authentication system. Changes affect all API security.

custom:
  owner: platform-team
  sla: 99.99%
```

---

## Content Sections

Content is organized into named sections, each containing blocks.

```yaml
content:
  sections:
    - name: Section Name
      blocks:
        - type: block_type
          # block-specific fields
```

---

## Block Types

Blocks have **strict field validation**. Only the fields listed below are allowed for each type. Unknown fields cause validation errors.

### `rule` Block

A textual rule or requirement.

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `type` | yes | string | Must be `rule` |
| `text` | yes | string | The rule text |
| `id` | no | string | Optional identifier |

```yaml
- type: rule
  id: no_direct_reverse
  text: Player cannot reverse direction directly into self.
```

### `table` Block

Tabular data with typed columns.

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `type` | yes | string | Must be `table` |
| `columns` | yes | object[] | Column definitions |
| `rows` | yes | object[] | Row data (keys must match column keys) |
| `id` | no | string | Optional identifier |

**Column fields** (strict - only these are allowed):

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `key` | yes | string | Column identifier used in row data |
| `type` | no | string | Data type hint (e.g., `string`, `int`, `enum`) |
| `enum` | no | string[] | Valid values if type is `enum` |
| `display` | no | string | Human-readable column header |

```yaml
- type: table
  id: food_types
  columns:
    - key: type
      display: Type
      type: string
    - key: points
      display: Points
      type: int
    - key: rarity
      display: Rarity
      type: enum
      enum: [common, rare, legendary]
  rows:
    - type: Apple
      points: 10
      rarity: common
    - type: Golden Apple
      points: 50
      rarity: legendary
```

### `param` Block

A configurable parameter or variable.

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `type` | yes | string | Must be `param` |
| `name` | yes | string | Parameter name |
| `datatype` | yes | string | Data type (e.g., `int`, `float`, `range_int`, `duration`, `string`) |
| `id` | no | string | Optional identifier |
| `min` | no | number | Minimum value (for numeric types) |
| `max` | no | number | Maximum value (for numeric types) |
| `default` | no | any | Default value |
| `unit` | no | string | Unit of measurement (e.g., `ms`, `px`, `%`) |
| `description` | no | string | Parameter description |

```yaml
- type: param
  id: tick_rate
  name: Tick Rate
  datatype: range_int
  min: 100
  max: 500
  default: 200
  unit: ms
  description: Time between game updates
```

### `mechanic` Block

A game mechanic or behavioral rule.

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `type` | yes | string | Must be `mechanic` |
| `name` | yes | string | Mechanic name |
| `description` | yes | string | What the mechanic does |
| `id` | no | string | Optional identifier |
| `conditions` | no | string[] | When this mechanic triggers |
| `inputs` | no | string[] | Required inputs |
| `outputs` | no | string[] | Results/effects |

```yaml
- type: mechanic
  id: food_collection
  name: Food Collection
  description: Snake collects food by moving head onto food position
  conditions:
    - snake head position equals food position
  inputs:
    - snake position
    - food position
  outputs:
    - snake length increases by 1
    - score increases by food point value
    - food respawns at new location
```

### `list` Block

A simple list of items.

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `type` | yes | string | Must be `list` |
| `items` | yes | string[] | List items |
| `id` | no | string | Optional identifier |

```yaml
- type: list
  id: controls
  items:
    - Arrow keys to change direction
    - Space to pause
    - Enter to restart
```

---

## References

The `refs` section links nodes together.

### Reference Types

| Type | Purpose | Format |
|------|---------|--------|
| `uses` | Hard dependencies | object[] with `target` and optional `context` |
| `related` | Informational links | object[] with `target` and optional `context` |
| `emits_events` | Events this node produces | string[] of node IDs |
| `vocabulary` | Shared term definitions | string[] of glossary node IDs |

### RefLink Fields

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `target` | yes | string | Target node ID |
| `context` | no | string | Why this reference exists |

```yaml
refs:
  uses:
    - target: items/food
      context: Food spawns on grid and snake collects it
    - target: systems/collision
      context: Handles wall and self-collision detection
  related:
    - target: systems/audio
  emits_events:
    - events/game/score_changed
    - events/game/game_over
  vocabulary:
    - glossaries/game-terms
```

---

## Issues (TBDs)

Track unresolved questions directly in nodes.

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `id` | yes | string | Unique identifier for this issue |
| `description` | yes | string | What needs to be resolved |
| `severity` | yes | string | One of: `low`, `medium`, `high`, `critical` |
| `location` | yes | string | Path to affected field (e.g., `content.sections[0]`) |
| `resolved` | yes | boolean | Whether the issue is resolved |

```yaml
issues:
  - id: tbd_difficulty_scaling
    description: How should difficulty increase as score grows?
    severity: medium
    location: content.sections.difficulty
    resolved: false
  - id: tbd_powerup_duration
    description: How long should speed boost last?
    severity: low
    location: content.sections.powerups
    resolved: false
```

---

## Contracts

BDD-style acceptance criteria using Given/When/Then.

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `name` | yes | string | Contract name |
| `scenario` | yes | string | Scenario description |
| `given` | no | string[] | Preconditions |
| `when` | no | string[] | Actions/triggers |
| `then` | no | string[] | Expected outcomes |

```yaml
contracts:
  - name: Eating food grows snake
    scenario: Snake head collides with food
    given:
      - snake has length 3
      - food exists at position (5, 5)
    when:
      - snake head moves to (5, 5)
    then:
      - snake length becomes 4
      - score increases by food point value
      - new food spawns at random empty cell
```

---

## Constraints

CEL expressions for custom validation rules.

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `expr` | yes | string | CEL expression that must evaluate to true |
| `message` | yes | string | Error message if constraint fails |
| `scope` | no | string | Which nodes this applies to (`all` or a specific kind) |

```yaml
constraints:
  - expr: "version > 0"
    message: Version must be positive
    scope: all
  - expr: "size(tags) > 0"
    message: At least one tag required
    scope: requirement
```

---

## Custom Block Types

Projects can define custom block types in `.deco/config.yaml`:

```yaml
custom_block_types:
  endpoint:
    required_fields:
      - method
      - path
      - response
    optional_fields:
      - auth
      - rate_limit
```

Custom types allow `required_fields` + `optional_fields` + `id`. If a custom type shares a name with a built-in type, both validations apply.

---

## Validation Errors

Common validation errors and how to fix them:

| Code | Error | Fix |
|------|-------|-----|
| E008 | Missing required field | Add the required field to the node |
| E047 | Block missing required field | Add the field (e.g., `text` for rule, `columns` for table) |
| E048 | Unknown block type | Use valid type: `rule`, `table`, `param`, `mechanic`, `list` |
| E049 | Unknown field in block | Remove the field or check spelling |
| E050 | Table column missing key | Add `key` field to column definition |

---

## File Organization

Nodes are stored in `.deco/nodes/` with paths matching IDs:

```
.deco/
├── config.yaml
├── history.jsonl
└── nodes/
    ├── systems/
    │   ├── core.yaml        # id: systems/core
    │   └── scoring.yaml     # id: systems/scoring
    ├── items/
    │   └── food.yaml        # id: items/food
    └── entities/
        └── player.yaml      # id: entities/player
```

The `id` field must match the file path relative to `nodes/` (without `.yaml`).
