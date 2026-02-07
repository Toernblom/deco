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

package migrations

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateBackup(t *testing.T) {
	// Create temp project directory
	tmpDir := t.TempDir()
	decoDir := filepath.Join(tmpDir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")

	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}

	// Create test config
	configContent := "project_name: test\nversion: 1\n"
	if err := os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Create test node
	nodeContent := "id: test-node\nkind: item\nversion: 1\nstatus: draft\ntitle: Test\n"
	if err := os.WriteFile(filepath.Join(nodesDir, "test-node.yaml"), []byte(nodeContent), 0644); err != nil {
		t.Fatalf("failed to write node: %v", err)
	}

	// Create backup
	result, err := CreateBackup(tmpDir)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Verify backup directory exists
	if _, err := os.Stat(result.BackupDir); os.IsNotExist(err) {
		t.Errorf("backup directory not created: %s", result.BackupDir)
	}

	// Verify config was backed up
	backupConfig := filepath.Join(result.BackupDir, "config.yaml")
	if _, err := os.Stat(backupConfig); os.IsNotExist(err) {
		t.Error("config.yaml not backed up")
	}

	// Verify node was backed up
	backupNode := filepath.Join(result.BackupDir, "nodes", "test-node.yaml")
	if _, err := os.Stat(backupNode); os.IsNotExist(err) {
		t.Error("node not backed up")
	}

	// Verify node count
	if result.NodeCount != 1 {
		t.Errorf("expected NodeCount=1, got %d", result.NodeCount)
	}
}

func TestRestoreBackup(t *testing.T) {
	// Create temp project directory
	tmpDir := t.TempDir()
	decoDir := filepath.Join(tmpDir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")

	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}

	// Create original content
	originalConfig := "project_name: original\nversion: 1\n"
	if err := os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(originalConfig), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	originalNode := "id: node-1\nkind: item\nversion: 1\nstatus: draft\ntitle: Original\n"
	if err := os.WriteFile(filepath.Join(nodesDir, "node-1.yaml"), []byte(originalNode), 0644); err != nil {
		t.Fatalf("failed to write node: %v", err)
	}

	// Create backup
	backup, err := CreateBackup(tmpDir)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Modify current state
	modifiedConfig := "project_name: modified\nversion: 2\n"
	if err := os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(modifiedConfig), 0644); err != nil {
		t.Fatalf("failed to modify config: %v", err)
	}

	// Restore backup
	if err := RestoreBackup(tmpDir, backup.BackupDir); err != nil {
		t.Fatalf("RestoreBackup failed: %v", err)
	}

	// Verify config was restored
	restoredConfig, err := os.ReadFile(filepath.Join(decoDir, "config.yaml"))
	if err != nil {
		t.Fatalf("failed to read restored config: %v", err)
	}
	if string(restoredConfig) != originalConfig {
		t.Errorf("config not restored correctly: got %q, want %q", string(restoredConfig), originalConfig)
	}
}

func TestListBackups(t *testing.T) {
	tmpDir := t.TempDir()
	decoDir := filepath.Join(tmpDir, ".deco")

	if err := os.MkdirAll(decoDir, 0755); err != nil {
		t.Fatalf("failed to create .deco: %v", err)
	}

	// Create some backup directories
	backupDirs := []string{
		"backup-20240101-120000",
		"backup-20240102-120000",
		"not-a-backup",
	}
	for _, dir := range backupDirs {
		if err := os.MkdirAll(filepath.Join(decoDir, dir), 0755); err != nil {
			t.Fatalf("failed to create %s: %v", dir, err)
		}
	}

	// List backups
	backups, err := ListBackups(tmpDir)
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}

	if len(backups) != 2 {
		t.Errorf("expected 2 backups, got %d", len(backups))
	}
}

func TestCreateBackup_NoNodes(t *testing.T) {
	tmpDir := t.TempDir()
	decoDir := filepath.Join(tmpDir, ".deco")

	if err := os.MkdirAll(decoDir, 0755); err != nil {
		t.Fatalf("failed to create .deco: %v", err)
	}

	// Create config only (no nodes directory)
	configContent := "project_name: test\nversion: 1\n"
	if err := os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Should succeed even without nodes
	result, err := CreateBackup(tmpDir)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	if result.NodeCount != 0 {
		t.Errorf("expected NodeCount=0, got %d", result.NodeCount)
	}
}
