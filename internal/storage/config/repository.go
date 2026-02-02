package config

// BlockTypeConfig defines validation rules for a custom block type.
type BlockTypeConfig struct {
	// RequiredFields lists field names that must be present in block.Data.
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
