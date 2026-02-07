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

package config

import (
	"fmt"
	"path/filepath"
)

// RefConstraint declares that a field references values from another block type.
type RefConstraint struct {
	BlockType string `yaml:"block_type" json:"block_type"` // the referenced block type
	Field     string `yaml:"field" json:"field"`           // the field name to match against
}

// FieldDef defines a field with type, required, and constraint information.
type FieldDef struct {
	Type     string          `yaml:"type" json:"type"`                       // string, number, list, bool
	Required bool            `yaml:"required" json:"required"`               // whether the field must be present
	Enum     []string        `yaml:"enum,omitempty" json:"enum,omitempty"`   // allowed values (for string fields)
	Refs     []RefConstraint `yaml:"refs,omitempty" json:"refs,omitempty"`   // cross-reference constraints (OR logic)
}

// UnmarshalYAML implements custom unmarshaling for FieldDef.
// Supports both the old single-object "ref:" syntax and the new array "ref:" syntax,
// normalizing both into the Refs slice.
func (fd *FieldDef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Use a raw struct to avoid infinite recursion
	var raw struct {
		Type     string      `yaml:"type"`
		Required bool        `yaml:"required"`
		Enum     []string    `yaml:"enum,omitempty"`
		Ref      interface{} `yaml:"ref,omitempty"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	fd.Type = raw.Type
	fd.Required = raw.Required
	fd.Enum = raw.Enum

	if raw.Ref == nil {
		return nil
	}

	// Determine if ref is a single object or an array
	switch v := raw.Ref.(type) {
	case map[string]interface{}:
		// Single ref object: ref: {block_type: x, field: y}
		rc, err := mapToRefConstraint(v)
		if err != nil {
			return fmt.Errorf("invalid ref constraint: %w", err)
		}
		fd.Refs = []RefConstraint{rc}
	case []interface{}:
		// Array of ref objects: ref: [{block_type: x, field: y}, ...]
		for i, item := range v {
			m, ok := item.(map[string]interface{})
			if !ok {
				return fmt.Errorf("ref[%d]: expected object, got %T", i, item)
			}
			rc, err := mapToRefConstraint(m)
			if err != nil {
				return fmt.Errorf("ref[%d]: %w", i, err)
			}
			fd.Refs = append(fd.Refs, rc)
		}
	default:
		return fmt.Errorf("ref: expected object or array, got %T", raw.Ref)
	}

	return nil
}

// mapToRefConstraint converts a raw map to a RefConstraint.
func mapToRefConstraint(m map[string]interface{}) (RefConstraint, error) {
	rc := RefConstraint{}
	if bt, ok := m["block_type"]; ok {
		rc.BlockType = fmt.Sprintf("%v", bt)
	}
	if f, ok := m["field"]; ok {
		rc.Field = fmt.Sprintf("%v", f)
	}
	if rc.BlockType == "" || rc.Field == "" {
		return rc, fmt.Errorf("ref constraint requires both block_type and field")
	}
	return rc, nil
}

// BlockTypeConfig defines validation rules for a custom block type.
type BlockTypeConfig struct {
	// RequiredFields lists field names that must be present in block.Data.
	RequiredFields []string `yaml:"required_fields" json:"required_fields"`
	// OptionalFields lists additional allowed fields for the block type.
	OptionalFields []string `yaml:"optional_fields,omitempty" json:"optional_fields,omitempty"`
	// Fields defines typed field definitions with constraints.
	// When set, provides type checking, enum validation, and required enforcement.
	Fields map[string]FieldDef `yaml:"fields,omitempty" json:"fields,omitempty"`
}

// SchemaRuleConfig defines validation rules for nodes of a specific kind.
type SchemaRuleConfig struct {
	// RequiredFields lists field names that must be present in the node's custom data.
	RequiredFields []string `yaml:"required_fields" json:"required_fields"`
}

// Config represents the project configuration.
// It defines where nodes are stored, project metadata, and other settings.
type Config struct {
	// ProjectName is the name of the game design project.
	ProjectName string `yaml:"project_name" json:"project_name"`

	// NodesPath is the directory path where node YAML files are stored.
	// Defaults to ".deco/nodes".
	NodesPath string `yaml:"nodes_path" json:"nodes_path"`

	// HistoryPath is the file path for the append-only audit log.
	// Defaults to ".deco/history.jsonl".
	HistoryPath string `yaml:"history_path" json:"history_path"`

	// Version is the config file format version.
	Version int `yaml:"version" json:"version"`

	// RequiredApprovals is the number of approvals needed for a node to be approved.
	// Defaults to 1.
	RequiredApprovals int `yaml:"required_approvals" json:"required_approvals"`

	// CustomBlockTypes defines additional block types beyond the built-in ones.
	// Keys are type names, values define validation rules.
	CustomBlockTypes map[string]BlockTypeConfig `yaml:"custom_block_types,omitempty" json:"custom_block_types,omitempty"`

	// SchemaRules defines per-kind validation rules for nodes.
	// Keys are kind names (e.g., "character", "quest"), values define required fields.
	SchemaRules map[string]SchemaRuleConfig `yaml:"schema_rules,omitempty" json:"schema_rules,omitempty"`

	// SchemaVersion is a hash of the schema configuration (CustomBlockTypes + SchemaRules).
	// Used to detect when schema changes require migration.
	SchemaVersion string `yaml:"schema_version,omitempty" json:"schema_version,omitempty"`

	// Custom allows projects to add arbitrary configuration fields.
	Custom map[string]interface{} `yaml:"custom,omitempty" json:"custom,omitempty"`
}

// Repository defines the interface for configuration persistence.
type Repository interface {
	// Load reads the project configuration from storage.
	// Returns the config or an error if not found or invalid.
	Load() (Config, error)

	// Save writes the project configuration to storage.
	Save(config Config) error
}

// DefaultNodesPath is the default directory for node YAML files.
const DefaultNodesPath = ".deco/nodes"

// DefaultHistoryPath is the default file path for the audit log.
const DefaultHistoryPath = ".deco/history.jsonl"

// ResolveNodesPath returns the absolute nodes directory path.
// Uses config.NodesPath if set, otherwise DefaultNodesPath.
func ResolveNodesPath(cfg Config, rootDir string) string {
	path := cfg.NodesPath
	if path == "" {
		path = DefaultNodesPath
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(rootDir, path)
}

// ResolveHistoryPath returns the absolute history file path.
// Uses config.HistoryPath if set, otherwise DefaultHistoryPath.
func ResolveHistoryPath(cfg Config, rootDir string) string {
	path := cfg.HistoryPath
	if path == "" {
		path = DefaultHistoryPath
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(rootDir, path)
}
