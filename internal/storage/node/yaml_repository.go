package node

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"gopkg.in/yaml.v3"
)

// YAMLRepository implements Repository using YAML files on the filesystem
type YAMLRepository struct {
	rootDir string
}

// NewYAMLRepository creates a new YAML-based node repository
// rootDir should be the project root (containing .deco/)
func NewYAMLRepository(rootDir string) *YAMLRepository {
	return &YAMLRepository{
		rootDir: rootDir,
	}
}

// nodesDir returns the path to the nodes directory
func (r *YAMLRepository) nodesDir() string {
	return filepath.Join(r.rootDir, ".deco", "nodes")
}

// pathForNode returns the file path for a given node ID
func (r *YAMLRepository) pathForNode(id string) string {
	return filepath.Join(r.nodesDir(), id+".yaml")
}

// LoadAll loads all nodes from storage
func (r *YAMLRepository) LoadAll() ([]domain.Node, error) {
	nodesDir := r.nodesDir()

	// Check if nodes directory exists
	if _, err := os.Stat(nodesDir); os.IsNotExist(err) {
		return []domain.Node{}, nil
	}

	var nodes []domain.Node

	// Walk the nodes directory recursively
	err := filepath.Walk(nodesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .yaml files
		if !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		// Load the node
		node, err := r.loadFromFile(path)
		if err != nil {
			return fmt.Errorf("failed to load %s: %w", path, err)
		}

		nodes = append(nodes, node)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// Load retrieves a single node by its ID
func (r *YAMLRepository) Load(id string) (domain.Node, error) {
	path := r.pathForNode(id)
	return r.loadFromFile(path)
}

// loadFromFile loads a node from a specific file path
func (r *YAMLRepository) loadFromFile(path string) (domain.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.Node{}, fmt.Errorf("node not found: %s", path)
		}
		return domain.Node{}, fmt.Errorf("failed to read file: %w", err)
	}

	var node domain.Node
	err = yaml.Unmarshal(data, &node)
	if err != nil {
		return domain.Node{}, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return node, nil
}

// Save persists a node to storage
func (r *YAMLRepository) Save(node domain.Node) error {
	path := r.pathForNode(node.ID)

	// Create parent directories if they don't exist
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Marshal node to YAML
	data, err := yaml.Marshal(&node)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Write to file
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Delete removes a node from storage by ID
func (r *YAMLRepository) Delete(id string) error {
	path := r.pathForNode(id)

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("node not found: %s", id)
	}

	// Delete the file
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	return nil
}

// Exists checks if a node with the given ID exists in storage
func (r *YAMLRepository) Exists(id string) (bool, error) {
	path := r.pathForNode(id)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
