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

package domain

import (
	"fmt"
	"sort"
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
	// Raw YAML content (set during load, not serialized) - used for line number tracking
	RawContent []byte `json:"-" yaml:"-"`

	// References to other nodes
	Refs Ref `json:"refs,omitempty" yaml:"refs,omitempty"`

	// Structured content sections
	Content *Content `json:"content,omitempty" yaml:"content,omitempty"`

	// Tracked TBDs and questions
	Issues []Issue `json:"issues,omitempty" yaml:"issues,omitempty"`

	// External doc references
	Docs []DocRef `json:"docs,omitempty" yaml:"docs,omitempty"`

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
// Keys are sorted for deterministic output (important for content hashing).
func (b Block) MarshalYAML() (interface{}, error) {
	// Use yaml.Node to control field order: 'type' first, then Data keys sorted
	node := &yaml.Node{
		Kind: yaml.MappingNode,
	}

	// Add 'type' first
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "type"},
		&yaml.Node{Kind: yaml.ScalarNode, Value: b.Type},
	)

	// Sort Data keys for deterministic ordering
	keys := make([]string, 0, len(b.Data))
	for k := range b.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Add Data fields in sorted order
	for _, k := range keys {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: k}
		var valueNode yaml.Node
		if err := valueNode.Encode(b.Data[k]); err != nil {
			return nil, fmt.Errorf("failed to encode field %s: %w", k, err)
		}
		node.Content = append(node.Content, keyNode, &valueNode)
	}

	return node, nil
}

// DocRef represents a reference to an external markdown file.
type DocRef struct {
	Path     string   `json:"path" yaml:"path"`
	Keywords []string `json:"keywords,omitempty" yaml:"keywords,omitempty"`
	Context  string   `json:"context,omitempty" yaml:"context,omitempty"`
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
	Name      string    `json:"name" yaml:"name"`                     // reviewer email/username
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`           // when approved
	Version   int       `json:"version" yaml:"version"`               // version that was approved
	Note      string    `json:"note,omitempty" yaml:"note,omitempty"` // optional comment
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

// SortedStringMap wraps map[string]string for deterministic YAML output.
// Keys are sorted alphabetically when marshaling to YAML.
type SortedStringMap map[string]string

func (m SortedStringMap) MarshalYAML() (interface{}, error) {
	if len(m) == 0 {
		return nil, nil
	}
	node := &yaml.Node{Kind: yaml.MappingNode}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: k},
			&yaml.Node{Kind: yaml.ScalarNode, Value: m[k]},
		)
	}
	return node, nil
}

// SortedInterfaceMap wraps map[string]interface{} for deterministic YAML output.
// Keys are sorted alphabetically when marshaling to YAML.
type SortedInterfaceMap map[string]interface{}

func (m SortedInterfaceMap) MarshalYAML() (interface{}, error) {
	if len(m) == 0 {
		return nil, nil
	}
	node := &yaml.Node{Kind: yaml.MappingNode}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: k}
		var valueNode yaml.Node
		if err := valueNode.Encode(m[k]); err != nil {
			return nil, fmt.Errorf("failed to encode field %s: %w", k, err)
		}
		node.Content = append(node.Content, keyNode, &valueNode)
	}
	return node, nil
}

// nodeForMarshal is an internal type for marshaling Node with sorted map keys
type nodeForMarshal struct {
	ID          string             `yaml:"id"`
	Kind        string             `yaml:"kind"`
	Version     int                `yaml:"version"`
	Status      string             `yaml:"status"`
	Title       string             `yaml:"title"`
	Tags        []string           `yaml:"tags,omitempty"`
	Refs        Ref                `yaml:"refs,omitempty"`
	Docs        []DocRef           `yaml:"docs,omitempty"`
	Content     *Content           `yaml:"content,omitempty"`
	Issues      []Issue            `yaml:"issues,omitempty"`
	Summary     string             `yaml:"summary,omitempty"`
	Glossary    SortedStringMap    `yaml:"glossary,omitempty"`
	Contracts   []Contract         `yaml:"contracts,omitempty"`
	LLMContext  string             `yaml:"llm_context,omitempty"`
	Constraints []Constraint       `yaml:"constraints,omitempty"`
	Reviewers   []Reviewer         `yaml:"reviewers,omitempty"`
	Custom      SortedInterfaceMap `yaml:"custom,omitempty"`
}

// MarshalYAML implements custom YAML marshaling for Node.
// It ensures Glossary and Custom maps are serialized with sorted keys.
func (n Node) MarshalYAML() (interface{}, error) {
	return nodeForMarshal{
		ID:          n.ID,
		Kind:        n.Kind,
		Version:     n.Version,
		Status:      n.Status,
		Title:       n.Title,
		Tags:        n.Tags,
		Refs:        n.Refs,
		Docs:        n.Docs,
		Content:     n.Content,
		Issues:      n.Issues,
		Summary:     n.Summary,
		Glossary:    SortedStringMap(n.Glossary),
		Contracts:   n.Contracts,
		LLMContext:  n.LLMContext,
		Constraints: n.Constraints,
		Reviewers:   n.Reviewers,
		Custom:      SortedInterfaceMap(n.Custom),
	}, nil
}
