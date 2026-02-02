package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"gopkg.in/yaml.v3"
)

func TestSetCommand_ExpectHash(t *testing.T) {
	t.Run("succeeds with matching hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		// Get the current hash
		hash := getNodeHash(t, tmpDir, "sword-001")

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"--expect-hash", hash, "sword-001", "title", "New Title", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with matching hash, got %v", err)
		}

		// Verify change was applied
		content := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(content, "New Title") {
			t.Errorf("Expected title to be changed, got: %s", content)
		}
	})

	t.Run("fails with mismatched hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"--expect-hash", "wronghash1234567", "sword-001", "title", "New Title", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected conflict error with mismatched hash, got nil")
		}

		// Check it's the right error type
		exitErr, ok := err.(*ExitError)
		if !ok {
			t.Fatalf("Expected ExitError, got %T: %v", err, err)
		}
		if exitErr.Code != ExitCodeConflict {
			t.Errorf("Expected exit code %d, got %d", ExitCodeConflict, exitErr.Code)
		}

		// Verify error message
		if !strings.Contains(exitErr.Message, "Conflict detected") {
			t.Errorf("Expected conflict message, got: %s", exitErr.Message)
		}
		if !strings.Contains(exitErr.Message, "Expected hash") {
			t.Errorf("Expected hash info in message, got: %s", exitErr.Message)
		}

		// Verify change was NOT applied
		content := readNodeFile(t, tmpDir, "sword-001")
		if strings.Contains(content, "New Title") {
			t.Errorf("Expected title to NOT be changed, got: %s", content)
		}
	})

	t.Run("force flag bypasses hash check", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"--expect-hash", "wronghash1234567", "--force", "sword-001", "title", "Forced Title", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected --force to bypass hash check, got error: %v", err)
		}

		// Verify change was applied despite wrong hash
		content := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(content, "Forced Title") {
			t.Errorf("Expected title to be changed with --force, got: %s", content)
		}
	})

	t.Run("no hash check without expect-hash flag", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"sword-001", "title", "Normal Update", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error without --expect-hash, got %v", err)
		}
	})
}

func TestAppendCommand_ExpectHash(t *testing.T) {
	t.Run("fails with mismatched hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{"--expect-hash", "wronghash1234567", "sword-001", "tags", "legendary", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected conflict error with mismatched hash")
		}

		exitErr, ok := err.(*ExitError)
		if !ok {
			t.Fatalf("Expected ExitError, got %T", err)
		}
		if exitErr.Code != ExitCodeConflict {
			t.Errorf("Expected exit code %d, got %d", ExitCodeConflict, exitErr.Code)
		}
	})

	t.Run("force flag bypasses hash check", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{"--expect-hash", "wronghash1234567", "--force", "sword-001", "tags", "legendary", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected --force to bypass hash check, got error: %v", err)
		}
	})
}

func TestUnsetCommand_ExpectHash(t *testing.T) {
	t.Run("fails with mismatched hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"--expect-hash", "wronghash1234567", "sword-001", "summary", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected conflict error with mismatched hash")
		}

		exitErr, ok := err.(*ExitError)
		if !ok {
			t.Fatalf("Expected ExitError, got %T", err)
		}
		if exitErr.Code != ExitCodeConflict {
			t.Errorf("Expected exit code %d, got %d", ExitCodeConflict, exitErr.Code)
		}
	})

	t.Run("force flag bypasses hash check", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"--expect-hash", "wronghash1234567", "--force", "sword-001", "summary", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected --force to bypass hash check, got error: %v", err)
		}
	})
}

func TestApplyCommand_ExpectHash(t *testing.T) {
	t.Run("fails with mismatched hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		// Create a patch file
		patchPath := filepath.Join(tmpDir, "patch.json")
		patchContent := `[{"op": "set", "path": "title", "value": "Patched Title"}]`
		os.WriteFile(patchPath, []byte(patchContent), 0644)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"--expect-hash", "wronghash1234567", "sword-001", patchPath, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected conflict error with mismatched hash")
		}

		exitErr, ok := err.(*ExitError)
		if !ok {
			t.Fatalf("Expected ExitError, got %T", err)
		}
		if exitErr.Code != ExitCodeConflict {
			t.Errorf("Expected exit code %d, got %d", ExitCodeConflict, exitErr.Code)
		}
	})

	t.Run("force flag bypasses hash check", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		// Create a patch file
		patchPath := filepath.Join(tmpDir, "patch.json")
		patchContent := `[{"op": "set", "path": "title", "value": "Forced Patch"}]`
		os.WriteFile(patchPath, []byte(patchContent), 0644)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"--expect-hash", "wronghash1234567", "--force", "sword-001", patchPath, tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected --force to bypass hash check, got error: %v", err)
		}

		// Verify patch was applied
		content := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(content, "Forced Patch") {
			t.Errorf("Expected patch to be applied with --force, got: %s", content)
		}
	})
}

func TestShowCommand_DisplaysHash(t *testing.T) {
	t.Run("json output includes content_hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		cmd := NewShowCommand()
		cmd.SetArgs([]string{"--json", "sword-001", tmpDir})
		err := cmd.Execute()

		w.Close()
		os.Stdout = old

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		buf := make([]byte, 4096)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if !strings.Contains(output, `"content_hash"`) {
			t.Errorf("Expected 'content_hash' in JSON output, got: %s", output)
		}
	})
}

func TestCheckContentHash(t *testing.T) {
	t.Run("returns nil for empty expectHash", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		n := loadTestNode(t, tmpDir, "sword-001")
		err := CheckContentHash(n, "")
		if err != nil {
			t.Errorf("Expected nil for empty hash, got %v", err)
		}
	})

	t.Run("returns nil for matching hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		n := loadTestNode(t, tmpDir, "sword-001")
		hash := ComputeContentHash(n)
		err := CheckContentHash(n, hash)
		if err != nil {
			t.Errorf("Expected nil for matching hash, got %v", err)
		}
	})

	t.Run("returns error for mismatched hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForConflict(t, tmpDir)

		n := loadTestNode(t, tmpDir, "sword-001")
		err := CheckContentHash(n, "wronghash")
		if err == nil {
			t.Fatal("Expected error for mismatched hash")
		}

		exitErr, ok := err.(*ExitError)
		if !ok {
			t.Fatalf("Expected ExitError, got %T", err)
		}
		if exitErr.Code != ExitCodeConflict {
			t.Errorf("Expected exit code %d, got %d", ExitCodeConflict, exitErr.Code)
		}
	})
}

// Test helpers

func setupProjectForConflict(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: conflict-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create a node to test with
	nodeYAML := `id: sword-001
kind: item
version: 1
status: draft
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

func getNodeHash(t *testing.T, dir, nodeID string) string {
	t.Helper()
	n := loadTestNode(t, dir, nodeID)
	return ComputeContentHash(n)
}

func loadTestNode(t *testing.T, dir, nodeID string) domain.Node {
	t.Helper()
	nodePath := filepath.Join(dir, ".deco", "nodes", nodeID+".yaml")
	content, err := os.ReadFile(nodePath)
	if err != nil {
		t.Fatalf("Failed to read node file: %v", err)
	}

	var n domain.Node
	if err := yaml.Unmarshal(content, &n); err != nil {
		t.Fatalf("Failed to parse node: %v", err)
	}
	return n
}
