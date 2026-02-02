package config

import "path/filepath"

// BlockTypeConfig defines validation rules for a custom block type.
type BlockTypeConfig struct {
	// RequiredFields lists field names that must be present in block.Data.
	RequiredFields []string `yaml:"required_fields" json:"required_fields"`
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
