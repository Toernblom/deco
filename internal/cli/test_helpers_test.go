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

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/node"
	"gopkg.in/yaml.v3"
)

// setupDecoProject creates a temporary directory with an initialized deco project.
func setupDecoProject(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	// Initialize a deco project
	cmd := NewInitCommand()
	cmd.SetArgs([]string{tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Failed to initialize deco project: %v", err)
	}

	return tmpDir
}

// createTestNode creates a simple test node in the given directory.
func createTestNode(t *testing.T, dir, id string) {
	t.Helper()
	nodesDir := filepath.Join(dir, ".deco", "nodes")
	nodeRepo := node.NewYAMLRepository(nodesDir)
	n := domain.Node{
		ID:      id,
		Kind:    "test",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
	}
	if err := nodeRepo.Save(n); err != nil {
		t.Fatalf("Failed to create test node: %v", err)
	}
}

// createTestNodeWithRefs creates a test node with references to other nodes.
func createTestNodeWithRefs(t *testing.T, dir, id string, refs []string) {
	t.Helper()

	// Build refs
	uses := make([]domain.RefLink, len(refs))
	for i, ref := range refs {
		uses[i] = domain.RefLink{Target: ref}
	}

	n := domain.Node{
		ID:      id,
		Kind:    "test",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Refs: domain.Ref{
			Uses: uses,
		},
	}

	// Save using YAML directly to ensure refs are saved correctly
	nodePath := filepath.Join(dir, ".deco", "nodes", id+".yaml")
	data, err := yaml.Marshal(n)
	if err != nil {
		t.Fatalf("Failed to marshal node: %v", err)
	}
	if err := os.WriteFile(nodePath, data, 0644); err != nil {
		t.Fatalf("Failed to write node file: %v", err)
	}
}

// readNodeYAML reads and parses a node file as a raw map.
func readNodeYAML(t *testing.T, dir, id string) map[string]interface{} {
	t.Helper()
	nodePath := filepath.Join(dir, ".deco", "nodes", id+".yaml")
	content, err := os.ReadFile(nodePath)
	if err != nil {
		t.Fatalf("Failed to read node %s: %v", id, err)
	}

	var node map[string]interface{}
	if err := yaml.Unmarshal(content, &node); err != nil {
		t.Fatalf("Failed to parse node YAML: %v", err)
	}

	return node
}
