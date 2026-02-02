package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func setupMigrateTestProject(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	decoDir := filepath.Join(tmpDir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")

	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}

	// Create config with schema rules (no schema_version, so migration needed)
	configContent := `project_name: test
version: 1
nodes_path: .deco/nodes
schema_rules:
  quest:
    required_fields:
      - objective
`
	if err := os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	return tmpDir
}

func TestMigrateCommand_NoMigrationNeeded(t *testing.T) {
	tmpDir := t.TempDir()
	decoDir := filepath.Join(tmpDir, ".deco")

	if err := os.MkdirAll(decoDir, 0755); err != nil {
		t.Fatalf("failed to create .deco: %v", err)
	}

	// Create config without schema rules (empty schema = no migration needed)
	configContent := `project_name: test
version: 1
nodes_path: .deco/nodes
`
	if err := os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cmd := NewMigrateCommand()
	cmd.SetArgs([]string{tmpDir})

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
}

func TestMigrateCommand_DryRun(t *testing.T) {
	tmpDir := setupMigrateTestProject(t)

	cmd := NewMigrateCommand()
	cmd.SetArgs([]string{"--dry-run", tmpDir})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Verify config was NOT modified
	configPath := filepath.Join(tmpDir, ".deco", "config.yaml")
	content, _ := os.ReadFile(configPath)
	if bytes.Contains(content, []byte("schema_version:")) {
		t.Error("config should not be modified in dry run")
	}
}

func TestMigrateCommand_AppliesMigration(t *testing.T) {
	tmpDir := setupMigrateTestProject(t)

	cmd := NewMigrateCommand()
	cmd.SetArgs([]string{tmpDir})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Verify config was modified
	configPath := filepath.Join(tmpDir, ".deco", "config.yaml")
	content, _ := os.ReadFile(configPath)
	if !bytes.Contains(content, []byte("schema_version:")) {
		t.Error("config should have schema_version after migration")
	}
}

func TestMigrateCommand_CreatesBackup(t *testing.T) {
	tmpDir := setupMigrateTestProject(t)

	cmd := NewMigrateCommand()
	cmd.SetArgs([]string{tmpDir})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Check for backup directory
	entries, _ := os.ReadDir(filepath.Join(tmpDir, ".deco"))
	hasBackup := false
	for _, e := range entries {
		if e.IsDir() && len(e.Name()) > 7 && e.Name()[:7] == "backup-" {
			hasBackup = true
			break
		}
	}
	if !hasBackup {
		t.Error("expected backup directory to be created")
	}
}

func TestMigrateCommand_NoBackupFlag(t *testing.T) {
	tmpDir := setupMigrateTestProject(t)

	cmd := NewMigrateCommand()
	cmd.SetArgs([]string{"--no-backup", tmpDir})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Check that no backup directory was created
	entries, _ := os.ReadDir(filepath.Join(tmpDir, ".deco"))
	for _, e := range entries {
		if e.IsDir() && len(e.Name()) > 7 && e.Name()[:7] == "backup-" {
			t.Error("backup directory should not be created with --no-backup")
		}
	}
}
