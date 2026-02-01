package domain

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Node represents a design node in the game design document.
// Nodes are the core building blocks of a GDD, containing metadata,
// references to other nodes, structured content, and tracked issues.
type Node struct {
	// Core metadata
	ID      string   `json:"id" yaml:"id"`
	Kind    string   `json:"kind" yaml:"kind"`       // mechanic, system, feature, etc.
	Version int      `json:"version" yaml:"version"` // Incremented on updates
	Status  string   `json:"status" yaml:"status"`   // draft, approved, deprecated, etc.
	Title   string   `json:"title" yaml:"title"`
	Tags    []string `json:"tags,omitempty" yaml:"tags,omitempty"`

	// Source location (set during load, not serialized)
	SourceFile string `json:"-" yaml:"-"`

	// References to other nodes
	Refs Ref `json:"refs,omitempty" yaml:"refs,omitempty"`

	// Structured content sections
	Content *Content `json:"content,omitempty" yaml:"content,omitempty"`

	// Tracked TBDs and questions
	Issues []Issue `json:"issues,omitempty" yaml:"issues,omitempty"`

	// Optional fields
	Summary     string                 `json:"summary,omitempty" yaml:"summary,omitempty"`
	Glossary    map[string]string      `json:"glossary,omitempty" yaml:"glossary,omitempty"`
	Contracts   []Contract             `json:"contracts,omitempty" yaml:"contracts,omitempty"`
	LLMContext  string                 `json:"llm_context,omitempty" yaml:"llm_context,omitempty"`
	Constraints []Constraint           `json:"constraints,omitempty" yaml:"constraints,omitempty"`
	Reviewers   []Reviewer             `json:"reviewers,omitempty" yaml:"reviewers,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty" yaml:"custom,omitempty"`
}

// Content holds the structured content sections of a node.
type Content struct {
	Sections []Section `json:"sections" yaml:"sections"`
}

// Section represents a logical section in the node content.
type Section struct {
	Name   string  `json:"name" yaml:"name"`
	Blocks []Block `json:"blocks" yaml:"blocks"`
}

// Block represents a content block within a section.
// Block fields (except 'type') are stored in Data map and written as inline YAML fields.
type Block struct {
	Type string                 `json:"type"` // table, rule, param, mechanic, list, etc.
	Data map[string]interface{} `json:"data,omitempty"`
}

// UnmarshalYAML implements custom YAML unmarshaling for Block.
// It captures all fields except 'type' into the Data map, preserving inline block fields.
func (b *Block) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping node for block, got %v", node.Kind)
	}

	b.Data = make(map[string]interface{})

	// Iterate through key-value pairs
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		key := keyNode.Value

		if key == "type" {
			b.Type = valueNode.Value
		} else if key == "data" {
			// If there's an explicit 'data' field, merge its contents
			var dataMap map[string]interface{}
			if err := valueNode.Decode(&dataMap); err != nil {
				return fmt.Errorf("failed to decode data field: %w", err)
			}
			for k, v := range dataMap {
				b.Data[k] = v
			}
		} else {
			// All other fields go into Data
			var value interface{}
			if err := valueNode.Decode(&value); err != nil {
				return fmt.Errorf("failed to decode field %s: %w", key, err)
			}
			b.Data[key] = value
		}
	}

	return nil
}

// MarshalYAML implements custom YAML marshaling for Block.
// It writes Data fields as inline block fields (not nested under 'data:').
func (b Block) MarshalYAML() (interface{}, error) {
	// Build a map with 'type' first, then all Data fields
	result := make(map[string]interface{})
	result["type"] = b.Type

	for k, v := range b.Data {
		result[k] = v
	}

	return result, nil
}

// Contract represents a Gherkin-style scenario or acceptance criteria.
type Contract struct {
	Name     string   `json:"name" yaml:"name"`
	Scenario string   `json:"scenario" yaml:"scenario"`
	Given    []string `json:"given,omitempty" yaml:"given,omitempty"`
	When     []string `json:"when,omitempty" yaml:"when,omitempty"`
	Then     []string `json:"then,omitempty" yaml:"then,omitempty"`
}

// Reviewer represents an approval record for a node version.
type Reviewer struct {
	Name      string    `json:"name" yaml:"name"`                         // reviewer email/username
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`               // when approved
	Version   int       `json:"version" yaml:"version"`                   // version that was approved
	Note      string    `json:"note,omitempty" yaml:"note,omitempty"`     // optional comment
}

// Validate checks that all required fields are present and valid.
func (n *Node) Validate() error {
	if n.ID == "" {
		return fmt.Errorf("node ID is required")
	}
	if n.Kind == "" {
		return fmt.Errorf("node Kind is required")
	}
	if n.Version == 0 {
		return fmt.Errorf("node Version is required and must be > 0")
	}
	if n.Status == "" {
		return fmt.Errorf("node Status is required")
	}
	if n.Title == "" {
		return fmt.Errorf("node Title is required")
	}
	return nil
}
