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
	nodesDir string
}

// NewYAMLRepository creates a new YAML-based node repository.
// nodesDir is the directory where node YAML files are stored.
// Use config.ResolveNodesPath() to get this from the project config.
func NewYAMLRepository(nodesDir string) *YAMLRepository {
	return &YAMLRepository{
		nodesDir: nodesDir,
	}
}

// nodesPath returns the path to the nodes directory
func (r *YAMLRepository) nodesPath() string {
	return r.nodesDir
}

// pathForNode returns the file path for a given node ID
func (r *YAMLRepository) pathForNode(id string) string {
	return filepath.Join(r.nodesPath(), id+".yaml")
}

// LoadAll loads all nodes from storage
func (r *YAMLRepository) LoadAll() ([]domain.Node, error) {
	nodesDir := r.nodesPath()

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

	// Store the source file path and raw content for error reporting
	node.SourceFile = path
	node.RawContent = data

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
