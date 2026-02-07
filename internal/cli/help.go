// Copyright (C) 2026 Anton Törnblom
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewLLMHelpCommand creates the llm-help subcommand
func NewLLMHelpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "llm-help",
		Short: "Output LLM-optimized reference guide",
		Long:  `Output a concise reference guide optimized for LLM consumption.`,
		Run: func(cmd *cobra.Command, args []string) {
			printLLMReference()
		},
	}

	return cmd
}

func printLLMReference() {
	fmt.Print(llmReference)
}

const llmReference = `# Deco LLM Quick Reference

Deco manages structured documentation as validated YAML nodes with typed blocks, cross-references, and history.

## Workflow

Edit YAML files directly in .deco/nodes/, then run 'deco sync' to detect changes and update history.
Run 'deco validate' after edits to catch errors. All validation errors include error codes (E0xx).

## Commands

Reading & Querying:
  deco list [--kind X] [--status X] [--tag X]   List nodes
  deco show <id> [--json]                        Show node + reverse refs
  deco query [term] [--kind X] [--tag X]         Search/filter nodes
  deco query --block-type X [--field key=val]    Query blocks within nodes
  deco validate [--quiet]                        Check all nodes
  deco issues [--severity X] [--node X]          List open TBDs
  deco stats                                     Project health overview
  deco graph [--format dot|mermaid|ascii]        Show dependency graph

History & Sync:
  deco sync [--dry-run]                          Detect edits, bump versions, track history
  deco history [--node <id>]                     Show audit log
  deco diff <id> [--since 2h]                    Show changes over time

Review:
  deco review submit <id>                        Submit for review
  deco review approve <id> [--note "msg"]        Approve node
  deco review reject <id> --note "reason"        Reject node
  deco review status [<id>]                      Check review status

Setup:
  deco init                                      Initialize new project
  deco migrate                                   Migrate from older format

## Querying

Node-level (default): filters and searches across nodes.
  deco query "sword"                     # Search title/summary
  deco query --kind item --tag combat    # Filter by kind + tag
  deco query "sword" --kind item         # Combine search + filter

Block-level: filters blocks within node content. Activated by --block-type.
  deco query --block-type building                         # All building blocks
  deco query --block-type building --field age=bronze      # Filter by field value
  deco query --block-type building --field materials=Planks # List membership (contains)
  deco query --block-type recipe --field output=Bronze     # Multiple --field uses AND logic
  deco query --kind system --block-type building           # Combine node + block filters

Follow mode: traverses ref constraints to find related blocks.
  deco query --block-type building --follow materials                # Auto from ref config
  deco query --block-type building --field age=bronze --follow materials  # Filter + follow
  deco query --block-type recipe --follow inputs                    # Reverse: what resources?
  deco query --block-type building --follow materials:recipe.output # Explicit target

Returns block data with context: [node_id > section_name] type + all fields.
Follow mode groups results by value with reference counts.

## Node Structure (YAML)

Required fields:
  id: systems/example      # Must match file path (systems/example.yaml)
  kind: system             # system, component, feature, item, etc.
  version: 1               # Auto-incremented by sync
  status: draft            # draft, review, approved, published, deprecated
  title: Example System

Optional fields:
  tags: [core, gameplay]
  summary: |
    Multiline description.
  refs:
    uses:
      - target: other/node
        context: Why this dependency
    related:
      - target: another/node
    emits_events:
      - events/some_event
    vocabulary:
      - glossaries/terms
  content:
    sections:
      - name: Section Name
        blocks:
          - type: rule|table|param|mechanic|list|doc|<custom>
            # block-specific fields below
  docs:
    - path: narratives/chapter-1.md       # Relative to project root
      keywords: [protagonist, betrayal]   # Must appear in file (case-insensitive)
      context: Main narrative             # Optional description
  issues:
    - id: tbd_1
      description: Unresolved question
      severity: medium        # low, medium, high, critical
      location: content.sections[0]
      resolved: false
  contracts:
    - name: Contract name
      scenario: Description
      given: [preconditions]
      when: [actions]
      then: [expected results]
  constraints:
    - expr: "version > 0"        # CEL expression
      message: "Version must be positive"
      scope: all                 # all or specific kind
  glossary:
    term: Definition
  llm_context: Extra context for AI
  custom:
    any_key: any_value

## Built-in Block Types

rule:
  - type: rule
    text: The rule text

table:
  - type: table
    columns:
      - key: name
        type: string
        display: Name
    rows:
      - name: Value

param:
  - type: param
    name: Tick Rate
    datatype: int           # int, float, range_int, duration, string
    default: 200
    min: 100
    max: 500
    unit: ms

mechanic:
  - type: mechanic
    name: Collection
    description: What happens
    conditions: [when this occurs]
    outputs: [what results]

list:
  - type: list
    items: [item1, item2]

doc:
  - type: doc
    path: narratives/scene.md           # Required: relative to project root
    keywords: [storm, arrival]          # Optional: validated against file content
    context: Opening scene              # Optional: description

## External Doc References

Reference .md files from nodes (node-level) or blocks (doc block type).
Deco validates that files exist and contain declared keywords.

Node-level docs:
  docs:
    - path: narratives/chapter-1.md
      keywords: [protagonist, betrayal]
      context: Full chapter narrative

Block-level doc (in content sections):
  - type: doc
    path: narratives/opening-scene.md
    keywords: [storm, mysterious stranger]

Validation:
  E055 - Doc file not found (path doesn't resolve)
  E056 - Missing keyword in doc (case-insensitive substring match)

Sync integration:
  Changes to referenced .md files trigger version bumps and review resets,
  just like editing the YAML itself. This creates a feedback loop:
  .md changes -> sync bumps version -> review needed -> check keywords still match.

## Custom Block Types

Define in .deco/config.yaml. Two syntaxes supported:

Simple syntax (field presence validation only):
  custom_block_types:
    endpoint:
      required_fields: [method, path, response]
      optional_fields: [auth, rate_limit]

Advanced syntax (typed fields with constraints and cross-references):
  custom_block_types:
    building:
      fields:
        name: {type: string, required: true}
        age: {type: string, required: true, enum: [stone, bronze, iron]}
        category: {type: string, enum: [production, military, residential]}
        materials:
          type: list
          ref:
            - {block_type: resource, field: name}
            - {block_type: recipe, field: output}
    resource:
      fields:
        name: {type: string, required: true}
        tier: {type: number, required: true}
    recipe:
      fields:
        output: {type: string, required: true, ref: {block_type: resource, field: name}}
        inputs: {type: list, ref: {block_type: resource, field: name}}

Field definition options:
  type:      string, number, list, bool (validated at E052)
  required:  true/false - field must be present (validated at E047)
  enum:      [val1, val2] - restrict to allowed values (validated at E053, suggests typo fixes)
  ref:       {block_type: X, field: Y} - single ref target (validated at E054)
  ref:       [{block_type: X, field: Y}, ...] - union refs, OR logic (validated at E054)

Usage in nodes:
  - type: building
    name: Smithy
    age: bronze
    category: production
    materials: [Iron, Wood]
  - type: resource
    name: Iron
    tier: 2

Cross-references are validated across ALL nodes. If materials references
resource.name, then every value in materials must match a name field in some
resource block somewhere in the project. Union refs (array form) validate
against ANY of the listed targets using OR logic.

## Schema Rules

Define per-kind field requirements in .deco/config.yaml:
  schema_rules:
    requirement:
      required_fields: [priority, acceptance_criteria]
    component:
      required_fields: [owner]

## Validation Error Codes

Key errors you'll encounter:
  E008  Missing required node field (id, kind, version, status, title)
  E010  Unknown field in node or nested structure (typo detection with suggestions)
  E020  Reference target not found
  E047  Missing required block field
  E048  Unknown block type (not built-in or custom)
  E049  Unknown field in block (strict validation, no extra fields allowed)
  E051  Missing required schema rule field
  E052  Field type mismatch (e.g., string where number expected)
  E053  Invalid enum value (with did-you-mean suggestions)
  E054  Cross-reference not found (value doesn't exist in referenced block type)
  E055  Doc file not found (referenced .md file missing)
  E056  Missing keyword in doc (keyword not in .md file content)

## File Layout

.deco/
  config.yaml              # Project config, custom block types, schema rules
  history.jsonl            # Append-only audit log
  nodes/
    systems/core.yaml      # id: systems/core
    items/food.yaml        # id: items/food

## Tips for LLMs

1. Edit YAML directly: Create/modify files in .deco/nodes/, then run 'deco sync'
2. Always read before edit: Use 'deco show <id>' or read the YAML file
3. Validate after changes: Run 'deco validate' to catch errors early
4. Use block-level queries: 'deco query --block-type X --field k=v' to find specific data
   - List fields support membership: '--field materials=Planks' matches if Planks is in the list
5. Define typed blocks: Use advanced custom_block_types with type/enum/ref for data integrity
6. Cross-reference blocks: Use ref constraints to link block types (e.g., recipe -> resource)
   - Union refs: ref as array allows OR validation across multiple block types
7. Follow refs across types: 'deco query --block-type building --follow materials' traces supply chains
   - Use explicit targets for ad-hoc joins: '--follow materials:recipe.output'
8. Use doc references: Put prose in .md files, reference with docs or doc blocks
9. Use issues for TBDs: Don't leave unresolved questions in content
10. Reference other nodes: Use refs.uses for dependencies between nodes
11. Keep nodes focused: One concept per node, link related concepts
12. Match id to path: systems/auth.yaml must have id: systems/auth
`
