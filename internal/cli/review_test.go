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
	"bytes"
	"fmt"
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

func TestReviewCommand_Approve(t *testing.T) {
	t.Run("approve adds reviewer to node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "review")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"approve", "test/node", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "reviewers:") {
			t.Error("Expected reviewers field to be added")
		}
	})

	t.Run("approve with note includes note in reviewer", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "review")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"approve", "test/node", "--note", "LGTM", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "LGTM") {
			t.Errorf("Expected note 'LGTM' in node, got: %s", nodeYAML)
		}
	})

	t.Run("approve transitions to approved when requirements met", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir) // requires 1 approval
		createNodeWithStatus(t, tmpDir, "test/node", "review")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"approve", "test/node", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "status: approved") {
			t.Errorf("Expected status 'approved', got: %s", nodeYAML)
		}
	})

	t.Run("approve fails if not in review status", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "draft")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"approve", "test/node", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when approving non-review node")
		}
	})
}

func TestReviewCommand_Reject(t *testing.T) {
	t.Run("reject changes status from review to draft", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "review")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"reject", "test/node", "--note", "Needs more detail", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "status: draft") {
			t.Errorf("Expected status 'draft', got: %s", nodeYAML)
		}
	})

	t.Run("reject requires note", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "review")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"reject", "test/node", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when rejecting without note")
		}
	})

	t.Run("reject fails if not in review status", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "draft")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"reject", "test/node", "--note", "reason", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when rejecting non-review node")
		}
	})
}

func TestReviewCommand_Status(t *testing.T) {
	t.Run("status shows review state", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithReviewers(t, tmpDir, "test/node", "review", 1)

		cmd := NewReviewCommand()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetArgs([]string{"status", "test/node", tmpDir})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "review") {
			t.Errorf("Expected output to contain status, got: %s", output)
		}
	})
}

func TestReviewWorkflow_Integration(t *testing.T) {
	t.Run("full workflow: draft -> review -> approved", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createDraftNode(t, tmpDir, "test/node")

		// 1. Submit for review
		submitCmd := NewReviewCommand()
		submitCmd.SetArgs([]string{"submit", "test/node", tmpDir})
		if err := submitCmd.Execute(); err != nil {
			t.Fatalf("Submit failed: %v", err)
		}

		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "status: review") {
			t.Fatalf("Expected status 'review' after submit")
		}

		// 2. Approve
		approveCmd := NewReviewCommand()
		approveCmd.SetArgs([]string{"approve", "test/node", "--note", "LGTM", tmpDir})
		if err := approveCmd.Execute(); err != nil {
			t.Fatalf("Approve failed: %v", err)
		}

		nodeYAML = readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "status: approved") {
			t.Fatalf("Expected status 'approved' after approve")
		}

		// Note: Edit->draft workflow is now tested via sync_test.go
		// (file edit + deco sync resets status to draft)
	})
}

func createNodeWithReviewers(t *testing.T, tmpDir, nodeID, status string, numReviewers int) {
	t.Helper()
	nodeContent := `id: ` + nodeID + `
kind: mechanic
version: 1
status: ` + status + `
title: Test Node
`
	if numReviewers > 0 {
		nodeContent += "reviewers:\n"
		for i := 0; i < numReviewers; i++ {
			nodeContent += fmt.Sprintf(`  - name: reviewer%d@example.com
    timestamp: 2026-01-01T00:00:00Z
    version: 1
`, i+1)
		}
	}
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	nodePath := filepath.Join(nodesDir, nodeID+".yaml")
	os.MkdirAll(filepath.Dir(nodePath), 0755)
	os.WriteFile(nodePath, []byte(nodeContent), 0644)
}
