package domain

import "fmt"

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
type Block struct {
	Type string                 `json:"type" yaml:"type"` // table, rule, param, mechanic, list, etc.
	Data map[string]interface{} `json:"data" yaml:"data"`
}

// Contract represents a Gherkin-style scenario or acceptance criteria.
type Contract struct {
	Name     string   `json:"name" yaml:"name"`
	Scenario string   `json:"scenario" yaml:"scenario"`
	Given    []string `json:"given,omitempty" yaml:"given,omitempty"`
	When     []string `json:"when,omitempty" yaml:"when,omitempty"`
	Then     []string `json:"then,omitempty" yaml:"then,omitempty"`
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
