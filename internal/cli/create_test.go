package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCreateCommand_Structure(t *testing.T) {
	t.Run("creates command", func(t *testing.T) {
		cmd := NewCreateCommand()
		if cmd == nil {
			t.Fatal("Expected create command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "create") {
			t.Errorf("Expected Use to start with 'create', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewCreateCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("requires id argument", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"-d", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when no id provided")
		}
	})
}

func TestCreateCommand_CreatesNode(t *testing.T) {
	t.Run("creates node file", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"systems/combat", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "systems", "combat.yaml")
		if _, err := os.Stat(nodePath); os.IsNotExist(err) {
			t.Errorf("Expected node file at %s", nodePath)
		}
	})

	t.Run("creates node with required fields", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"sword-001", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "sword-001.yaml")
		content, err := os.ReadFile(nodePath)
		if err != nil {
			t.Fatalf("Failed to read node file: %v", err)
		}

		var node map[string]interface{}
		if err := yaml.Unmarshal(content, &node); err != nil {
			t.Fatalf("Failed to parse node YAML: %v", err)
		}

		// Check required fields
		if node["id"] != "sword-001" {
			t.Errorf("Expected id 'sword-001', got %v", node["id"])
		}
		if node["version"] != 1 {
			t.Errorf("Expected version 1, got %v", node["version"])
		}
		if node["status"] != "draft" {
			t.Errorf("Expected status 'draft', got %v", node["status"])
		}
		if node["kind"] == nil || node["kind"] == "" {
			t.Error("Expected kind to be set")
		}
		if node["title"] == nil || node["title"] == "" {
			t.Error("Expected title to be set")
		}
	})
}

func TestCreateCommand_Flags(t *testing.T) {
	t.Run("has kind flag", func(t *testing.T) {
		cmd := NewCreateCommand()
		flag := cmd.Flags().Lookup("kind")
		if flag == nil {
			t.Fatal("Expected --kind flag")
		}
		if flag.Shorthand != "k" {
			t.Errorf("Expected shorthand 'k', got %q", flag.Shorthand)
		}
	})

	t.Run("has title flag", func(t *testing.T) {
		cmd := NewCreateCommand()
		flag := cmd.Flags().Lookup("title")
		if flag == nil {
			t.Fatal("Expected --title flag")
		}
		if flag.Shorthand != "t" {
			t.Errorf("Expected shorthand 't', got %q", flag.Shorthand)
		}
	})

	t.Run("has status flag", func(t *testing.T) {
		cmd := NewCreateCommand()
		flag := cmd.Flags().Lookup("status")
		if flag == nil {
			t.Fatal("Expected --status flag")
		}
		if flag.Shorthand != "s" {
			t.Errorf("Expected shorthand 's', got %q", flag.Shorthand)
		}
	})

	t.Run("kind flag sets node kind", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"my-feature", "--kind", "feature", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		node := readNodeYAML(t, tmpDir, "my-feature")
		if node["kind"] != "feature" {
			t.Errorf("Expected kind 'feature', got %v", node["kind"])
		}
	})

	t.Run("title flag sets node title", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"combat-001", "--title", "Combat System", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		node := readNodeYAML(t, tmpDir, "combat-001")
		if node["title"] != "Combat System" {
			t.Errorf("Expected title 'Combat System', got %v", node["title"])
		}
	})

	t.Run("status flag sets node status", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"approved-001", "--status", "approved", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		node := readNodeYAML(t, tmpDir, "approved-001")
		if node["status"] != "approved" {
			t.Errorf("Expected status 'approved', got %v", node["status"])
		}
	})
}

func TestCreateCommand_Defaults(t *testing.T) {
	t.Run("default kind is system", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"default-test", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		node := readNodeYAML(t, tmpDir, "default-test")
		if node["kind"] != "system" {
			t.Errorf("Expected default kind 'system', got %v", node["kind"])
		}
	})

	t.Run("default status is draft", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"draft-test", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		node := readNodeYAML(t, tmpDir, "draft-test")
		if node["status"] != "draft" {
			t.Errorf("Expected default status 'draft', got %v", node["status"])
		}
	})

	t.Run("default title is derived from id", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"my-cool-feature", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		node := readNodeYAML(t, tmpDir, "my-cool-feature")
		title, ok := node["title"].(string)
		if !ok || title == "" {
			t.Error("Expected non-empty title")
		}
	})
}

func TestCreateCommand_ExistingNode(t *testing.T) {
	t.Run("errors if node already exists", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		// Create node first
		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"existing-node", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("First create failed: %v", err)
		}

		// Try to create again
		cmd = NewCreateCommand()
		cmd.SetArgs([]string{"existing-node", "-d", tmpDir})
		err = cmd.Execute()
		if err == nil {
			t.Error("Expected error when node already exists")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("Expected 'already exists' error, got: %v", err)
		}
	})

	t.Run("force flag overwrites existing node", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		// Create node with title "Original"
		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"overwrite-test", "--title", "Original", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("First create failed: %v", err)
		}

		// Overwrite with --force
		cmd = NewCreateCommand()
		cmd.SetArgs([]string{"overwrite-test", "--title", "Replaced", "--force", "-d", tmpDir})
		err = cmd.Execute()
		if err != nil {
			t.Fatalf("Force create failed: %v", err)
		}

		node := readNodeYAML(t, tmpDir, "overwrite-test")
		if node["title"] != "Replaced" {
			t.Errorf("Expected title 'Replaced', got %v", node["title"])
		}
	})
}

func TestCreateCommand_NestedPaths(t *testing.T) {
	t.Run("creates nested directory structure", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"systems/combat/melee", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "systems", "combat", "melee.yaml")
		if _, err := os.Stat(nodePath); os.IsNotExist(err) {
			t.Errorf("Expected nested node file at %s", nodePath)
		}
	})

	t.Run("node id includes full path", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"items/weapons/sword", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "items", "weapons", "sword.yaml")
		content, err := os.ReadFile(nodePath)
		if err != nil {
			t.Fatalf("Failed to read node: %v", err)
		}

		var node map[string]interface{}
		yaml.Unmarshal(content, &node)
		if node["id"] != "items/weapons/sword" {
			t.Errorf("Expected id 'items/weapons/sword', got %v", node["id"])
		}
	})
}

func TestCreateCommand_RequiresProject(t *testing.T) {
	t.Run("errors without deco project", func(t *testing.T) {
		tmpDir := t.TempDir() // No .deco directory

		cmd := NewCreateCommand()
		cmd.SetArgs([]string{"test-node", "-d", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error without deco project")
		}
		if !strings.Contains(err.Error(), ".deco") {
			t.Errorf("Expected error about .deco directory, got: %v", err)
		}
	})
}

// Helper functions

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
