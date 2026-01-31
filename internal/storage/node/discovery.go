package node

import (
	"os"
	"path/filepath"
	"strings"
)

// FileDiscovery handles discovering node files in the filesystem
type FileDiscovery struct {
	rootDir string
}

// NewFileDiscovery creates a new file discovery instance
func NewFileDiscovery(rootDir string) *FileDiscovery {
	return &FileDiscovery{
		rootDir: rootDir,
	}
}

// nodesDir returns the path to the nodes directory
func (d *FileDiscovery) nodesDir() string {
	return filepath.Join(d.rootDir, ".deco", "nodes")
}

// DiscoverAll finds all .yaml files in the nodes directory
func (d *FileDiscovery) DiscoverAll() ([]string, error) {
	nodesDir := d.nodesDir()

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
	nodesDir := d.nodesDir()
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
	nodesDir := d.nodesDir()

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
	nodesDir := d.nodesDir()

	// Convert forward slashes to OS-specific separator
	parts := strings.Split(id, "/")
	pathParts := append([]string{nodesDir}, parts...)
	path := filepath.Join(pathParts...)

	return path + ".yaml"
}
