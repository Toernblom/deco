package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
)

func TestComputeContentHash(t *testing.T) {
	node := domain.Node{
		ID:      "test-001",
		Kind:    "item",
		Version: 1,
		Status:  "draft",
		Title:   "Test Title",
		Summary: "Test summary",
		Tags:    []string{"tag1", "tag2"},
	}

	t.Run("returns consistent hash for same content", func(t *testing.T) {
		hash1 := computeContentHash(node)
		hash2 := computeContentHash(node)
		if hash1 != hash2 {
			t.Errorf("Expected consistent hash, got %s and %s", hash1, hash2)
		}
		if len(hash1) != 16 {
			t.Errorf("Expected 16 char hash, got %d chars", len(hash1))
		}
	})

	t.Run("returns different hash for different content", func(t *testing.T) {
		modified := node
		modified.Title = "Different Title"
		hash1 := computeContentHash(node)
		hash2 := computeContentHash(modified)
		if hash1 == hash2 {
			t.Error("Expected different hash for different content")
		}
	})

	t.Run("ignores metadata fields", func(t *testing.T) {
		modified := node
		modified.Version = 99
		modified.Status = "approved"
		modified.Reviewers = []domain.Reviewer{{Name: "alice", Version: 1}}
		hash1 := computeContentHash(node)
		hash2 := computeContentHash(modified)
		if hash1 != hash2 {
			t.Error("Expected same hash when only metadata differs")
		}
	})
}

func TestSyncCommand_Structure(t *testing.T) {
	t.Run("creates sync command", func(t *testing.T) {
		cmd := NewSyncCommand()
		if cmd == nil {
			t.Fatal("Expected sync command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "sync") {
			t.Errorf("Expected Use to start with 'sync', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewSyncCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("has dry-run flag", func(t *testing.T) {
		cmd := NewSyncCommand()
		flag := cmd.Flags().Lookup("dry-run")
		if flag == nil {
			t.Fatal("Expected --dry-run flag to be defined")
		}
	})

	t.Run("has quiet flag", func(t *testing.T) {
		cmd := NewSyncCommand()
		flag := cmd.Flags().Lookup("quiet")
		if flag == nil {
			t.Fatal("Expected --quiet flag to be defined")
		}
		if flag.Shorthand != "q" {
			t.Errorf("Expected shorthand 'q', got %q", flag.Shorthand)
		}
	})
}

func TestContentChanged(t *testing.T) {
	baseNode := domain.Node{
		ID:      "test-001",
		Kind:    "mechanic",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Summary: "A test node",
		Tags:    []string{"tag1", "tag2"},
	}

	t.Run("returns false for identical nodes", func(t *testing.T) {
		old := baseNode
		new := baseNode
		if contentChanged(old, new) {
			t.Error("Expected no content change for identical nodes")
		}
	})

	t.Run("returns true when title changes", func(t *testing.T) {
		old := baseNode
		new := baseNode
		new.Title = "Modified Title"
		if !contentChanged(old, new) {
			t.Error("Expected content change when title differs")
		}
	})

	t.Run("returns true when summary changes", func(t *testing.T) {
		old := baseNode
		new := baseNode
		new.Summary = "Modified summary"
		if !contentChanged(old, new) {
			t.Error("Expected content change when summary differs")
		}
	})

	t.Run("returns true when tags change", func(t *testing.T) {
		old := baseNode
		new := baseNode
		new.Tags = []string{"tag1", "tag3"}
		if !contentChanged(old, new) {
			t.Error("Expected content change when tags differ")
		}
	})

	t.Run("returns true when tag added", func(t *testing.T) {
		old := baseNode
		new := baseNode
		new.Tags = []string{"tag1", "tag2", "tag3"}
		if !contentChanged(old, new) {
			t.Error("Expected content change when tag added")
		}
	})

	t.Run("ignores version changes", func(t *testing.T) {
		old := baseNode
		new := baseNode
		new.Version = 5
		if contentChanged(old, new) {
			t.Error("Expected no content change when only version differs")
		}
	})

	t.Run("ignores status changes", func(t *testing.T) {
		old := baseNode
		new := baseNode
		new.Status = "approved"
		if contentChanged(old, new) {
			t.Error("Expected no content change when only status differs")
		}
	})

	t.Run("ignores reviewer changes", func(t *testing.T) {
		old := baseNode
		new := baseNode
		new.Reviewers = []domain.Reviewer{{Name: "alice", Version: 1}}
		if contentChanged(old, new) {
			t.Error("Expected no content change when only reviewers differ")
		}
	})
}

func TestTagsEqual(t *testing.T) {
	t.Run("equal empty slices", func(t *testing.T) {
		if !tagsEqual([]string{}, []string{}) {
			t.Error("Expected empty slices to be equal")
		}
	})

	t.Run("equal non-empty slices", func(t *testing.T) {
		a := []string{"a", "b", "c"}
		b := []string{"a", "b", "c"}
		if !tagsEqual(a, b) {
			t.Error("Expected identical slices to be equal")
		}
	})

	t.Run("different lengths", func(t *testing.T) {
		a := []string{"a", "b"}
		b := []string{"a", "b", "c"}
		if tagsEqual(a, b) {
			t.Error("Expected different length slices to be unequal")
		}
	})

	t.Run("different values", func(t *testing.T) {
		a := []string{"a", "b", "c"}
		b := []string{"a", "x", "c"}
		if tagsEqual(a, b) {
			t.Error("Expected different value slices to be unequal")
		}
	})

	t.Run("nil vs empty", func(t *testing.T) {
		if !tagsEqual(nil, nil) {
			t.Error("Expected nil slices to be equal")
		}
	})
}

func TestExtractNodeID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{".deco/nodes/sword-001.yaml", "sword-001"},
		{".deco/nodes/test.yaml", "test"},
		{".deco/nodes/complex-id-123.yaml", "complex-id-123"},
		{"some/path/node.yaml", "node"},
		{"not-yaml.txt", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractNodeID(tt.input)
			if result != tt.expected {
				t.Errorf("extractNodeID(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsGitRepo(t *testing.T) {
	t.Run("returns true for git repo", func(t *testing.T) {
		tmpDir := t.TempDir()
		initGitRepo(t, tmpDir)

		if !isGitRepo(tmpDir) {
			t.Error("Expected isGitRepo to return true for initialized repo")
		}
	})

	t.Run("returns false for non-git directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		if isGitRepo(tmpDir) {
			t.Error("Expected isGitRepo to return false for non-git directory")
		}
	})
}

func TestRunSync_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		initGitRepo(t, tmpDir)

		flags := &syncFlags{targetDir: tmpDir}
		exitCode, err := runSync(flags)

		if err == nil {
			t.Error("Expected error for missing .deco directory")
		}
		if exitCode != syncExitError {
			t.Errorf("Expected exit code %d, got %d", syncExitError, exitCode)
		}
	})
}

func TestRunSync_NotGitRepo(t *testing.T) {
	t.Run("errors on non-git directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		flags := &syncFlags{targetDir: tmpDir}
		exitCode, err := runSync(flags)

		if err == nil {
			t.Error("Expected error for non-git directory")
		}
		if exitCode != syncExitError {
			t.Errorf("Expected exit code %d, got %d", syncExitError, exitCode)
		}
	})
}

func TestRunSync_NoChanges(t *testing.T) {
	t.Run("returns clean exit when no modified files", func(t *testing.T) {
		tmpDir := t.TempDir()
		initGitRepo(t, tmpDir)
		setupProjectForSync(t, tmpDir)
		commitAllChanges(t, tmpDir, "Initial commit")

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if exitCode != syncExitClean {
			t.Errorf("Expected exit code %d, got %d", syncExitClean, exitCode)
		}
	})
}

func TestRunSync_MetadataOnlyChange(t *testing.T) {
	t.Run("ignores metadata-only changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		initGitRepo(t, tmpDir)
		setupProjectForSync(t, tmpDir)
		commitAllChanges(t, tmpDir, "Initial commit")

		// Change only metadata (version, status) - should be ignored
		nodeYAML := `id: sword-001
kind: item
version: 5
status: approved
title: Iron Sword
summary: A basic iron sword
tags:
  - weapon
  - combat
`
		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "sword-001.yaml")
		if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
			t.Fatalf("Failed to modify node: %v", err)
		}

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if exitCode != syncExitClean {
			t.Errorf("Expected exit code %d (metadata-only change ignored), got %d", syncExitClean, exitCode)
		}
	})
}

func TestRunSync_ContentChange(t *testing.T) {
	t.Run("syncs nodes with content changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		initGitRepo(t, tmpDir)
		setupProjectForSync(t, tmpDir)
		commitAllChanges(t, tmpDir, "Initial commit")

		// Change content (title) - should trigger sync
		nodeYAML := `id: sword-001
kind: item
version: 1
status: approved
title: Golden Sword
summary: A basic iron sword
tags:
  - weapon
  - combat
reviewers:
  - name: alice@example.com
    timestamp: 2026-01-01T00:00:00Z
    version: 1
`
		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "sword-001.yaml")
		if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
			t.Fatalf("Failed to modify node: %v", err)
		}

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if exitCode != syncExitModified {
			t.Errorf("Expected exit code %d, got %d", syncExitModified, exitCode)
		}

		// Verify version was bumped
		content, _ := os.ReadFile(nodePath)
		if !strings.Contains(string(content), "version: 2") {
			t.Errorf("Expected version to be bumped to 2, got: %s", string(content))
		}

		// Verify status was reset to draft
		if !strings.Contains(string(content), "status: draft") {
			t.Errorf("Expected status to be reset to draft, got: %s", string(content))
		}

		// Verify reviewers were cleared
		if strings.Contains(string(content), "reviewers:") {
			t.Errorf("Expected reviewers to be cleared, got: %s", string(content))
		}
	})
}

func TestRunSync_DryRun(t *testing.T) {
	t.Run("dry-run does not modify files", func(t *testing.T) {
		tmpDir := t.TempDir()
		initGitRepo(t, tmpDir)
		setupProjectForSync(t, tmpDir)
		commitAllChanges(t, tmpDir, "Initial commit")

		// Change content
		nodeYAML := `id: sword-001
kind: item
version: 1
status: approved
title: Golden Sword
summary: A basic iron sword
tags:
  - weapon
  - combat
`
		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "sword-001.yaml")
		if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
			t.Fatalf("Failed to modify node: %v", err)
		}

		flags := &syncFlags{targetDir: tmpDir, dryRun: true, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Dry-run returns clean (0) since no actual changes were made
		if exitCode != syncExitClean {
			t.Errorf("Expected exit code %d for dry-run, got %d", syncExitClean, exitCode)
		}

		// Verify file was NOT modified
		content, _ := os.ReadFile(nodePath)
		if strings.Contains(string(content), "version: 2") {
			t.Error("Dry-run should not modify files")
		}
	})
}

func TestRunSync_HistoryLogging(t *testing.T) {
	t.Run("logs sync operation to history", func(t *testing.T) {
		tmpDir := t.TempDir()
		initGitRepo(t, tmpDir)
		setupProjectForSync(t, tmpDir)
		commitAllChanges(t, tmpDir, "Initial commit")

		// Change content
		nodeYAML := `id: sword-001
kind: item
version: 1
status: approved
title: Golden Sword
summary: A basic iron sword
tags:
  - weapon
  - combat
`
		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "sword-001.yaml")
		if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
			t.Fatalf("Failed to modify node: %v", err)
		}

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		_, err := runSync(flags)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify history entry was created
		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		content, err := os.ReadFile(historyPath)
		if err != nil {
			t.Fatalf("Failed to read history: %v", err)
		}

		historyStr := string(content)
		if !strings.Contains(historyStr, `"operation":"sync"`) {
			t.Errorf("Expected sync operation in history, got: %s", historyStr)
		}
		if !strings.Contains(historyStr, `"node_id":"sword-001"`) {
			t.Errorf("Expected node_id in history, got: %s", historyStr)
		}
	})
}

// Test helpers

func initGitRepo(t *testing.T, dir string) {
	t.Helper()

	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git user for commits
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = dir
	cmd.Run()

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	cmd.Run()
}

func commitAllChanges(t *testing.T, dir, message string) {
	t.Helper()

	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage changes: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
}

func setupProjectForSync(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: sync-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create a node
	nodeYAML := `id: sword-001
kind: item
version: 1
status: approved
title: Iron Sword
summary: A basic iron sword
tags:
  - weapon
  - combat
`
	nodePath := filepath.Join(nodesDir, "sword-001.yaml")
	if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
}
