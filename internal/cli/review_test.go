package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReviewCommand_Structure(t *testing.T) {
	t.Run("creates review command with subcommands", func(t *testing.T) {
		cmd := NewReviewCommand()
		if cmd == nil {
			t.Fatal("Expected review command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "review") {
			t.Errorf("Expected Use to start with 'review', got %q", cmd.Use)
		}

		// Should have submit subcommand
		submitCmd, _, _ := cmd.Find([]string{"submit"})
		if submitCmd == nil || submitCmd.Use == "review" {
			t.Error("Expected submit subcommand")
		}
	})
}

func TestReviewCommand_Submit(t *testing.T) {
	t.Run("submit changes status from draft to review", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createDraftNode(t, tmpDir, "test/node")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"submit", "test/node", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify status changed
		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "status: review") {
			t.Errorf("Expected status to be 'review', got: %s", nodeYAML)
		}
	})

	t.Run("submit fails if not in draft status", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "approved")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"submit", "test/node", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when submitting non-draft node")
		}
	})
}

func setupReviewProject(t *testing.T, tmpDir string) {
	t.Helper()
	decoDir := filepath.Join(tmpDir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	os.MkdirAll(nodesDir, 0755)

	configContent := `project_name: TestProject
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
version: 1
required_approvals: 1
`
	os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644)
}

func createDraftNode(t *testing.T, tmpDir, nodeID string) {
	t.Helper()
	createNodeWithStatus(t, tmpDir, nodeID, "draft")
}

func createNodeWithStatus(t *testing.T, tmpDir, nodeID, status string) {
	t.Helper()
	nodeContent := `id: ` + nodeID + `
kind: mechanic
version: 1
status: ` + status + `
title: Test Node
`
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	nodePath := filepath.Join(nodesDir, nodeID+".yaml")
	os.MkdirAll(filepath.Dir(nodePath), 0755)
	os.WriteFile(nodePath, []byte(nodeContent), 0644)
}

func readNodeFileByID(t *testing.T, tmpDir, nodeID string) string {
	t.Helper()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	nodePath := filepath.Join(nodesDir, nodeID+".yaml")
	content, err := os.ReadFile(nodePath)
	if err != nil {
		t.Fatalf("Failed to read node file: %v", err)
	}
	return string(content)
}
