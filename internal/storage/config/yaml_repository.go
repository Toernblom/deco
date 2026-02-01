package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// YAMLRepository implements Repository using a YAML file
type YAMLRepository struct {
	rootDir string
}

// NewYAMLRepository creates a new YAML-based config repository
func NewYAMLRepository(rootDir string) *YAMLRepository {
	return &YAMLRepository{
		rootDir: rootDir,
	}
}

// configPath returns the path to the config file
func (r *YAMLRepository) configPath() string {
	return filepath.Join(r.rootDir, ".deco", "config.yaml")
}

// Load reads the project configuration from storage
func (r *YAMLRepository) Load() (Config, error) {
	path := r.configPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("config file not found: %s", path)
		}
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	// Default required_approvals to 1 if not set
	if cfg.RequiredApprovals == 0 {
		cfg.RequiredApprovals = 1
	}

	return cfg, nil
}

// Save writes the project configuration to storage
func (r *YAMLRepository) Save(cfg Config) error {
	path := r.configPath()

	// Create parent directory if it doesn't exist
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	// Write to file
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
