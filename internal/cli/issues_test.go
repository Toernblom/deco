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

	t.Run("has kind flag", func(t *testing.T) {
		cmd := NewIssuesCommand()
		flag := cmd.Flags().Lookup("kind")
		if flag == nil {
			t.Fatal("Expected --kind flag")
		}
		if flag.Shorthand != "k" {
			t.Errorf("Expected shorthand 'k', got %q", flag.Shorthand)
		}
	})

	t.Run("has tag flag", func(t *testing.T) {
		cmd := NewIssuesCommand()
		flag := cmd.Flags().Lookup("tag")
		if flag == nil {
			t.Fatal("Expected --tag flag")
		}
		if flag.Shorthand != "t" {
			t.Errorf("Expected shorthand 't', got %q", flag.Shorthand)
		}
	})

	t.Run("has all flag", func(t *testing.T) {
		cmd := NewIssuesCommand()
		flag := cmd.Flags().Lookup("all")
		if flag == nil {
			t.Fatal("Expected --all flag")
		}
		if flag.Shorthand != "a" {
			t.Errorf("Expected shorthand 'a', got %q", flag.Shorthand)
		}
	})

	t.Run("has json flag", func(t *testing.T) {
		cmd := NewIssuesCommand()
		flag := cmd.Flags().Lookup("json")
		if flag == nil {
			t.Fatal("Expected --json flag")
		}
		if flag.Shorthand != "j" {
			t.Errorf("Expected shorthand 'j', got %q", flag.Shorthand)
		}
	})

	t.Run("has quiet flag", func(t *testing.T) {
		cmd := NewIssuesCommand()
		flag := cmd.Flags().Lookup("quiet")
		if flag == nil {
			t.Fatal("Expected --quiet flag")
		}
		if flag.Shorthand != "q" {
			t.Errorf("Expected shorthand 'q', got %q", flag.Shorthand)
		}
	})

	t.Run("has summary flag", func(t *testing.T) {
		cmd := NewIssuesCommand()
		flag := cmd.Flags().Lookup("summary")
		if flag == nil {
			t.Fatal("Expected --summary flag")
		}
	})
}

func TestIssuesCommand_KindFilter(t *testing.T) {
	t.Run("filters by node kind", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithKindAndIssues(t, tmpDir, "mech1", "mechanic", []domain.Issue{
			{ID: "mech-issue", Description: "Mechanic issue", Severity: "high", Location: "content"},
		})
		createNodeWithKindAndIssues(t, tmpDir, "sys1", "system", []domain.Issue{
			{ID: "sys-issue", Description: "System issue", Severity: "high", Location: "content"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"--kind", "mechanic", "-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "mech-issue") {
			t.Error("Expected output to contain mech-issue")
		}
		if strings.Contains(output, "sys-issue") {
			t.Error("Expected output to NOT contain sys-issue")
		}
	})
}

func TestIssuesCommand_TagFilter(t *testing.T) {
	t.Run("filters by node tag", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithTagsAndIssues(t, tmpDir, "node1", []string{"combat", "core"}, []domain.Issue{
			{ID: "combat-issue", Description: "Combat issue", Severity: "high", Location: "content"},
		})
		createNodeWithTagsAndIssues(t, tmpDir, "node2", []string{"ui"}, []domain.Issue{
			{ID: "ui-issue", Description: "UI issue", Severity: "high", Location: "content"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"--tag", "combat", "-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "combat-issue") {
			t.Error("Expected output to contain combat-issue")
		}
		if strings.Contains(output, "ui-issue") {
			t.Error("Expected output to NOT contain ui-issue")
		}
	})
}

func TestIssuesCommand_AllFlag(t *testing.T) {
	t.Run("shows resolved issues with --all", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "node1", []domain.Issue{
			{ID: "open-issue", Description: "Open", Severity: "high", Location: "content"},
			{ID: "resolved-issue", Description: "Resolved", Severity: "high", Location: "content", Resolved: true},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"--all", "-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "open-issue") {
			t.Error("Expected output to contain open-issue")
		}
		if !strings.Contains(output, "resolved-issue") {
			t.Error("Expected output to contain resolved-issue with --all")
		}
		if !strings.Contains(output, "RESOLVED") {
			t.Error("Expected resolved issues to be marked")
		}
	})
}

func TestIssuesCommand_JSONOutput(t *testing.T) {
	t.Run("outputs JSON format", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "node1", []domain.Issue{
			{ID: "issue1", Description: "Test issue", Severity: "high", Location: "content"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"--json", "-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.HasPrefix(strings.TrimSpace(output), "{") {
			t.Errorf("Expected JSON output to start with '{', got: %s", output)
		}
		if !strings.Contains(output, `"total"`) {
			t.Error("Expected JSON to contain 'total' field")
		}
		if !strings.Contains(output, `"issues"`) {
			t.Error("Expected JSON to contain 'issues' field")
		}
	})

	t.Run("json with summary shows by_node", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "node1", []domain.Issue{
			{ID: "issue1", Description: "Test issue", Severity: "high", Location: "content"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"--json", "--summary", "-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, `"by_node"`) {
			t.Error("Expected JSON summary to contain 'by_node' field")
		}
	})
}

func TestIssuesCommand_QuietOutput(t *testing.T) {
	t.Run("shows count only", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "node1", []domain.Issue{
			{ID: "issue1", Description: "Test", Severity: "high", Location: "content"},
			{ID: "issue2", Description: "Test", Severity: "low", Location: "content"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"--quiet", "-d", tmpDir})
		output := captureOutput(t, cmd)

		// Should show just the count
		output = strings.TrimSpace(output)
		if !strings.HasPrefix(output, "2") {
			t.Errorf("Expected quiet output to show '2', got: %s", output)
		}
	})

	t.Run("shows resolved count with --all", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "node1", []domain.Issue{
			{ID: "open", Description: "Open", Severity: "high", Location: "content"},
			{ID: "resolved", Description: "Resolved", Severity: "high", Location: "content", Resolved: true},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"--quiet", "--all", "-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "resolved") {
			t.Errorf("Expected quiet --all output to mention resolved, got: %s", output)
		}
	})
}

func TestIssuesCommand_SummaryOutput(t *testing.T) {
	t.Run("shows per-node rollup", func(t *testing.T) {
		tmpDir := setupDecoProject(t)
		createNodeWithIssues(t, tmpDir, "node1", []domain.Issue{
			{ID: "issue1", Description: "Test", Severity: "high", Location: "content"},
			{ID: "issue2", Description: "Test", Severity: "low", Location: "content"},
		})
		createNodeWithIssues(t, tmpDir, "node2", []domain.Issue{
			{ID: "issue3", Description: "Test", Severity: "critical", Location: "content"},
		})

		cmd := NewIssuesCommand()
		cmd.SetArgs([]string{"--summary", "-d", tmpDir})
		output := captureOutput(t, cmd)

		if !strings.Contains(output, "node1") {
			t.Error("Expected summary to contain node1")
		}
		if !strings.Contains(output, "node2") {
			t.Error("Expected summary to contain node2")
		}
		if !strings.Contains(output, "3 total") {
			t.Error("Expected summary to show total count")
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

func createNodeWithKindAndIssues(t *testing.T, dir, id, kind string, issues []domain.Issue) {
	t.Helper()

	n := domain.Node{
		ID:      id,
		Kind:    kind,
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

func createNodeWithTagsAndIssues(t *testing.T, dir, id string, tags []string, issues []domain.Issue) {
	t.Helper()

	n := domain.Node{
		ID:      id,
		Kind:    "test",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Tags:    tags,
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
