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

Deco manages documentation as validated YAML nodes with references and history.

## Workflow

Edit YAML files directly, then run 'deco sync' to detect changes and update history.

## Commands

Reading:
  deco list [--kind X] [--status X]   List nodes
  deco show <id>                      Show node + reverse refs
  deco query [--kind X] [--tag X]     Search/filter nodes
  deco validate                       Check all nodes
  deco issues                         List open TBDs by severity
  deco stats                          Project health overview
  deco graph --ascii                  Show dependency graph

History:
  deco sync                           Detect edits, bump versions, track history
  deco history [<id>]                 Show audit log
  deco diff <id>                      Show changes over time

Review:
  deco review submit <id>             Submit for review
  deco review approve <id>            Approve node
  deco review reject <id>             Reject node
  deco review status [<id>]           Check review status

Setup:
  deco init                           Initialize new project
  deco migrate                        Migrate from older format

## Node Structure (YAML)

Required fields:
  id: systems/example      # Matches file path
  kind: system             # system, component, feature, item, etc.
  version: 1               # Auto-incremented
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
  content:
    sections:
      - name: Section Name
        blocks:
          - type: rule|table|param|mechanic|list
            # block-specific fields
  issues:
    - id: tbd_1
      description: Unresolved question
      severity: medium        # low, medium, high, critical
      location: content.sections[0]
      resolved: false
  glossary:
    term: Definition
  llm_context: Extra context for AI

## Block Types

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

## Custom Block Types

Define in .deco/config.yaml:

custom_block_types:
  endpoint:                    # Block type name
    required_fields:
      - method
      - path
      - response
    optional_fields:
      - auth
      - rate_limit

Usage in nodes:
  - type: endpoint
    method: POST
    path: /api/users
    response: User object
    auth: bearer

## File Layout

.deco/
  config.yaml
  history.jsonl
  nodes/
    systems/core.yaml      # id: systems/core
    items/food.yaml        # id: items/food

## Tips for LLMs

1. Edit YAML directly: Create/modify files in .deco/nodes/, then run 'deco sync'
2. Always read before edit: Use 'deco show <id>' or read the YAML file
3. Validate after changes: Run 'deco validate' to catch errors
4. Use issues for TBDs: Don't leave unresolved questions in content
5. Reference other nodes: Use refs.uses for dependencies
6. Keep nodes focused: One concept per node, link related concepts
7. Match id to path: systems/auth.yaml must have id: systems/auth
8. Sync detects changes: Version auto-increments, history is tracked
`
