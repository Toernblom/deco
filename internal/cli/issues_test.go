package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func TestIssuesCommand_Structure(t *testing.T) {
	t.Run("creates command", func(t *testing.T) {
		cmd := NewIssuesCommand()
		if cmd == nil {
			t.Fatal("Expected issues command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "issues") {
			t.Errorf("Expected Use to start with 'issues', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewIssuesCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})
}

func TestIssuesCommand_ListsIssues(t *testing.T) {
	t.Run("lists all open issues", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "node1", []domain.Issue{
			{ID: "issue1", Description: "Fix this", Severity: "high", Location: "content"},
			{ID: "issue2", Description: "Check that", Severity: "low", Location: "refs"},
		})
		createNodeWithIssues(t, tmpDir, "node2", []domain.Issue{
			{ID: "issue3", Description: "Review this", Severity: "medium", Location: "meta"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"-d", tmpDir})

		// Capture output
		output := captureOutput(t, cmd)

		// All issues should appear
		if !strings.Contains(output, "issue1") {
			t.Error("Expected output to contain issue1")
		}
		if !strings.Contains(output, "issue2") {
			t.Error("Expected output to contain issue2")
		}
		if !strings.Contains(output, "issue3") {
			t.Error("Expected output to contain issue3")
		}
	})

	t.Run("excludes resolved issues", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "node1", []domain.Issue{
			{ID: "open-issue", Description: "Still open", Severity: "high", Location: "content"},
			{ID: "resolved-issue", Description: "Already done", Severity: "high", Location: "content", Resolved: true},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "open-issue") {
			t.Error("Expected output to contain open-issue")
		}
		if strings.Contains(output, "resolved-issue") {
			t.Error("Expected output to NOT contain resolved-issue")
		}
	})

	t.Run("shows no issues message when empty", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createTestNode(t, tmpDir, "clean-node")

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(strings.ToLower(output), "no") || !strings.Contains(strings.ToLower(output), "issue") {
			t.Errorf("Expected message about no issues, got: %s", output)
		}
	})
}

func TestIssuesCommand_Filters(t *testing.T) {
	t.Run("filters by severity", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "node1", []domain.Issue{
			{ID: "high-issue", Description: "Urgent", Severity: "high", Location: "content"},
			{ID: "low-issue", Description: "Minor", Severity: "low", Location: "refs"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"--severity", "high", "-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "high-issue") {
			t.Error("Expected output to contain high-issue")
		}
		if strings.Contains(output, "low-issue") {
			t.Error("Expected output to NOT contain low-issue")
		}
	})

	t.Run("filters by node", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "target-node", []domain.Issue{
			{ID: "target-issue", Description: "In target", Severity: "medium", Location: "content"},
		})
		createNodeWithIssues(t, tmpDir, "other-node", []domain.Issue{
			{ID: "other-issue", Description: "In other", Severity: "medium", Location: "content"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"--node", "target-node", "-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "target-issue") {
			t.Error("Expected output to contain target-issue")
		}
		if strings.Contains(output, "other-issue") {
			t.Error("Expected output to NOT contain other-issue")
		}
	})

	t.Run("combines severity and node filters", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "target", []domain.Issue{
			{ID: "target-high", Description: "High in target", Severity: "high", Location: "content"},
			{ID: "target-low", Description: "Low in target", Severity: "low", Location: "refs"},
		})
		createNodeWithIssues(t, tmpDir, "other", []domain.Issue{
			{ID: "other-high", Description: "High in other", Severity: "high", Location: "content"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"--severity", "high", "--node", "target", "-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "target-high") {
			t.Error("Expected output to contain target-high")
		}
		if strings.Contains(output, "target-low") {
			t.Error("Expected output to NOT contain target-low")
		}
		if strings.Contains(output, "other-high") {
			t.Error("Expected output to NOT contain other-high")
		}
	})
}

func TestIssuesCommand_ShowsContext(t *testing.T) {
	t.Run("shows node ID for each issue", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "my-node", []domain.Issue{
			{ID: "issue1", Description: "Test issue", Severity: "medium", Location: "content"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "my-node") {
			t.Errorf("Expected output to contain node ID 'my-node', got: %s", output)
		}
	})

	t.Run("shows location for each issue", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "node1", []domain.Issue{
			{ID: "issue1", Description: "Test issue", Severity: "high", Location: "content.sections[0]"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "content.sections[0]") {
			t.Errorf("Expected output to contain location, got: %s", output)
		}
	})

	t.Run("shows severity for each issue", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "node1", []domain.Issue{
			{ID: "issue1", Description: "Critical problem", Severity: "critical", Location: "content"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "critical") {
			t.Errorf("Expected output to contain severity 'critical', got: %s", output)
		}
	})
}

func TestIssuesCommand_Flags(t *testing.T) {
	t.Run("has severity flag", func(t *testing.T) {
		cmd := NewIssuesCommand()
		flag := cmd.Flags().Lookup("severity")
		if flag == nil {
			t.Fatal("Expected --severity flag")
		}
		if flag.Shorthand != "s" {
			t.Errorf("Expected shorthand 's', got %q", flag.Shorthand)
		}
	})

	t.Run("has node flag", func(t *testing.T) {
		cmd := NewIssuesCommand()
		flag := cmd.Flags().Lookup("node")
		if flag == nil {
			t.Fatal("Expected --node flag")
		}
		if flag.Shorthand != "n" {
			t.Errorf("Expected shorthand 'n', got %q", flag.Shorthand)
		}
	})
}

func TestIssuesCommand_RequiresProject(t *testing.T) {
	t.Run("errors without deco project", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"-d", tmpDir})
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

func createNodeWithIssues(t *testing.T, dir, id string, issues []domain.Issue) {
	t.Helper()

	n := domain.Node{
		ID:      id,
		Kind:    "test",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Issues:  issues,
	}

	nodePath := filepath.Join(dir, ".deco", "nodes", id+".yaml")
	data, err := yaml.Marshal(n)
	if err != nil {
		t.Fatalf("Failed to marshal node: %v", err)
	}
	if err := os.WriteFile(nodePath, data, 0644); err != nil {
		t.Fatalf("Failed to write node file: %v", err)
	}
}

func captureOutput(t *testing.T, cmd *cobra.Command) string {
	t.Helper()

	// Create a buffer to capture output
	var buf strings.Builder
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		// Include error in output for debugging
		buf.WriteString(err.Error())
	}

	return buf.String()
}
