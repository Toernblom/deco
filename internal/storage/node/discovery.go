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
	"os"
	"path/filepath"
	"strings"
)

// FileDiscovery handles discovering node files in the filesystem
type FileDiscovery struct {
	nodesDir string
}

// NewFileDiscovery creates a new file discovery instance.
// nodesDir is the directory where node YAML files are stored.
// Use config.ResolveNodesPath() to get this from the project config.
func NewFileDiscovery(nodesDir string) *FileDiscovery {
	return &FileDiscovery{
		nodesDir: nodesDir,
	}
}

// nodesPath returns the path to the nodes directory
func (d *FileDiscovery) nodesPath() string {
	return d.nodesDir
}

// DiscoverAll finds all .yaml files in the nodes directory
func (d *FileDiscovery) DiscoverAll() ([]string, error) {
	nodesDir := d.nodesPath()

	// Check if nodes directory exists
	if _, err := os.Stat(nodesDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	var files []string

	// Walk the directory tree
	err := filepath.Walk(nodesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only include .yaml files
		if strings.HasSuffix(path, ".yaml") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// DiscoverByPattern finds .yaml files matching a path pattern
func (d *FileDiscovery) DiscoverByPattern(pattern string) ([]string, error) {
	nodesDir := d.nodesPath()
	searchDir := filepath.Join(nodesDir, pattern)

	// Check if search directory exists
	if _, err := os.Stat(searchDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	var files []string

	// Walk the pattern directory
	err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only include .yaml files
		if strings.HasSuffix(path, ".yaml") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// PathToID converts a file path to a node ID
// Example: .deco/nodes/systems/food.yaml -> systems/food
func (d *FileDiscovery) PathToID(path string) string {
	nodesDir := d.nodesPath()

	// Get relative path from nodes directory
	relPath, err := filepath.Rel(nodesDir, path)
	if err != nil {
		return ""
	}

	// Remove .yaml extension
	id := strings.TrimSuffix(relPath, ".yaml")

	// Convert Windows backslashes to forward slashes for consistent IDs
	id = filepath.ToSlash(id)

	return id
}

// IDToPath converts a node ID to a file path
// Example: systems/food -> .deco/nodes/systems/food.yaml
func (d *FileDiscovery) IDToPath(id string) string {
	nodesDir := d.nodesPath()

	// Convert forward slashes to OS-specific separator
	parts := strings.Split(id, "/")
	pathParts := append([]string{nodesDir}, parts...)
	path := filepath.Join(pathParts...)

	return path + ".yaml"
}
