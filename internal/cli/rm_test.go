package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/node"
	"gopkg.in/yaml.v3"
)

func TestRmCommand_Structure(t *testing.T) {
	t.Run("creates command", func(t *testing.T) {
		cmd := NewRmCommand()
		if cmd == nil {
			t.Fatal("Expected rm command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "rm") {
			t.Errorf("Expected Use to start with 'rm', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewRmCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("requires id argument", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewRmCommand()
		cmd.SetArgs([]string{"-d", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when no id provided")
		}
	})
}

func TestRmCommand_DeletesNode(t *testing.T) {
	t.Run("deletes existing node", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createTestNode(t, tmpDir, "delete-me")

		cmd := NewRmCommand()
		cmd.SetArgs([]string{"delete-me", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify node is gone
		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "delete-me.yaml")
		if _, err := os.Stat(nodePath); !os.IsNotExist(err) {
			t.Error("Expected node file to be deleted")
		}
	})

	t.Run("errors if node does not exist", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		cmd := NewRmCommand()
		cmd.SetArgs([]string{"nonexistent", "-d", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for nonexistent node")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Expected 'not found' error, got: %v", err)
		}
	})
}

func TestRmCommand_ReverseRefs(t *testing.T) {
	t.Run("warns about reverse references", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		// Create target node
		createTestNode(t, tmpDir, "target")

		// Create referencing node
		createTestNodeWithRefs(t, tmpDir, "referencer", []string{"target"})

		cmd := NewRmCommand()
		cmd.SetArgs([]string{"target", "-d", tmpDir})
		err := cmd.Execute()

		if err == nil {
			t.Error("Expected error when deleting referenced node")
		}
		if !strings.Contains(err.Error(), "referencer") {
			t.Errorf("Expected error to mention referencing node, got: %v", err)
		}
	})

	t.Run("force flag deletes despite reverse refs", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		// Create target node
		createTestNode(t, tmpDir, "target")

		// Create referencing node
		createTestNodeWithRefs(t, tmpDir, "referencer", []string{"target"})

		cmd := NewRmCommand()
		cmd.SetArgs([]string{"target", "--force", "-d", tmpDir})
		err := cmd.Execute()

		if err != nil {
			t.Fatalf("Expected no error with --force, got %v", err)
		}

		// Verify node is gone
		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "target.yaml")
		if _, err := os.Stat(nodePath); !os.IsNotExist(err) {
			t.Error("Expected node file to be deleted with --force")
		}
	})

	t.Run("deletes node with no references", func(t *testing.T) {
		tmpDir := setupDecoProject(t)

		// Create orphan node (no references to it)
		createTestNode(t, tmpDir, "orphan")

		cmd := NewRmCommand()
		cmd.SetArgs([]string{"orphan", "-d", tmpDir})
		err := cmd.Execute()

		if err != nil {
			t.Fatalf("Expected no error for unreferenced node, got %v", err)
		}
	})
}

func TestRmCommand_LogsDeletion(t *testing.T) {
	t.Run("logs deletion in history", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createTestNode(t, tmpDir, "logged-delete")

		cmd := NewRmCommand()
		cmd.SetArgs([]string{"logged-delete", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Check history file exists and contains deletion
		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		content, err := os.ReadFile(historyPath)
		if err != nil {
			t.Fatalf("Failed to read history: %v", err)
		}

		if !strings.Contains(string(content), "logged-delete") {
			t.Error("Expected history to contain node ID")
		}
		if !strings.Contains(string(content), "delete") {
			t.Error("Expected history to contain 'delete' operation")
		}
	})
}

func TestRmCommand_NestedPaths(t *testing.T) {
	t.Run("deletes nested node", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createTestNode(t, tmpDir, "systems/combat")

		cmd := NewRmCommand()
		cmd.SetArgs([]string{"systems/combat", "-d", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "systems", "combat.yaml")
		if _, err := os.Stat(nodePath); !os.IsNotExist(err) {
			t.Error("Expected nested node file to be deleted")
		}
	})
}

func TestRmCommand_Flags(t *testing.T) {
	t.Run("has force flag", func(t *testing.T) {
		cmd := NewRmCommand()
		flag := cmd.Flags().Lookup("force")
		if flag == nil {
			t.Fatal("Expected --force flag")
		}
		if flag.Shorthand != "f" {
			t.Errorf("Expected shorthand 'f', got %q", flag.Shorthand)
		}
	})
}

func TestRmCommand_RequiresProject(t *testing.T) {
	t.Run("errors without deco project", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewRmCommand()
		cmd.SetArgs([]string{"some-node", "-d", tmpDir})
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

func createTestNode(t *testing.T, dir, id string) {
	t.Helper()
	nodeRepo := node.NewYAMLRepository(dir)
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
