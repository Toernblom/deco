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
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewCommand_Structure(t *testing.T) {
	t.Run("creates new command", func(t *testing.T) {
		cmd := NewNewCommand()
		if cmd == nil {
			t.Fatal("Expected new command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "new") {
			t.Errorf("Expected Use to start with 'new', got %q", cmd.Use)
		}
	})

	t.Run("requires kind and title flags", func(t *testing.T) {
		cmd := NewNewCommand()
		kindFlag := cmd.Flags().Lookup("kind")
		if kindFlag == nil {
			t.Fatal("Expected --kind flag")
		}
		titleFlag := cmd.Flags().Lookup("title")
		if titleFlag == nil {
			t.Fatal("Expected --title flag")
		}
	})
}

func TestNewCommand_CreatesNode(t *testing.T) {
	t.Run("creates a basic node", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewNewCommand()
		cmd.SetArgs([]string{"systems/combat", "--kind", "system", "--title", "Combat System", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify the node file exists
		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "systems", "combat.yaml")
		if _, err := os.Stat(nodePath); os.IsNotExist(err) {
			t.Fatal("Expected node file to be created")
		}

		// Verify content
		data, err := os.ReadFile(nodePath)
		if err != nil {
			t.Fatalf("Failed to read node file: %v", err)
		}

		var node map[string]interface{}
		if err := yaml.Unmarshal(data, &node); err != nil {
			t.Fatalf("Failed to parse YAML: %v", err)
		}

		if node["id"] != "systems/combat" {
			t.Errorf("Expected id 'systems/combat', got %v", node["id"])
		}
		if node["kind"] != "system" {
			t.Errorf("Expected kind 'system', got %v", node["kind"])
		}
		if node["version"] != 1 {
			t.Errorf("Expected version 1, got %v", node["version"])
		}
		if node["status"] != "draft" {
			t.Errorf("Expected status 'draft', got %v", node["status"])
		}
		if node["title"] != "Combat System" {
			t.Errorf("Expected title 'Combat System', got %v", node["title"])
		}
	})

	t.Run("creates node with tags and summary", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewNewCommand()
		cmd.SetArgs([]string{"mechanics/stealth", "--kind", "mechanic", "--title", "Stealth",
			"--tags", "core,pvp", "--summary", "Stealth mechanics", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		node := readNodeYAML(t, tmpDir, "mechanics/stealth")
		if node["summary"] != "Stealth mechanics" {
			t.Errorf("Expected summary 'Stealth mechanics', got %v", node["summary"])
		}

		tags, ok := node["tags"].([]interface{})
		if !ok {
			t.Fatalf("Expected tags to be a list, got %T", node["tags"])
		}
		if len(tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(tags))
		}
		if tags[0] != "core" || tags[1] != "pvp" {
			t.Errorf("Expected tags [core, pvp], got %v", tags)
		}
	})

	t.Run("creates parent directories automatically", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewNewCommand()
		cmd.SetArgs([]string{"deep/nested/node", "--kind", "system", "--title", "Deep Node", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "deep", "nested", "node.yaml")
		if _, err := os.Stat(nodePath); os.IsNotExist(err) {
			t.Fatal("Expected node file to be created in nested directory")
		}
	})
}

func TestNewCommand_Errors(t *testing.T) {
	t.Run("errors if node already exists", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createTestNode(t, tmpDir, "existing-node")

		cmd := NewNewCommand()
		cmd.SetArgs([]string{"existing-node", "--kind", "system", "--title", "Duplicate", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected error for existing node")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("Expected 'already exists' error, got %q", err.Error())
		}
	})

	t.Run("force overwrites existing node", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createTestNode(t, tmpDir, "existing-node")

		cmd := NewNewCommand()
		cmd.SetArgs([]string{"existing-node", "--kind", "system", "--title", "Overwritten", "--force", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with --force, got %v", err)
		}

		node := readNodeYAML(t, tmpDir, "existing-node")
		if node["title"] != "Overwritten" {
			t.Errorf("Expected title 'Overwritten', got %v", node["title"])
		}
	})

	t.Run("errors on missing project", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewNewCommand()
		cmd.SetArgs([]string{"test-node", "--kind", "system", "--title", "Test", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected error for missing project")
		}
	})
}

func TestNewCommand_LogsHistory(t *testing.T) {
	t.Run("logs creation to history", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewNewCommand()
		cmd.SetArgs([]string{"test-node", "--kind", "system", "--title", "Test", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		data, err := os.ReadFile(historyPath)
		if err != nil {
			t.Fatalf("Failed to read history: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "test-node") {
			t.Error("Expected history to contain node ID")
		}
		if !strings.Contains(content, "create") {
			t.Error("Expected history to contain 'create' operation")
		}
	})
}
