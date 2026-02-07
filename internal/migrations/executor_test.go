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

	"github.com/Toernblom/deco/internal/storage/config"
)

func setupTestProject(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	decoDir := filepath.Join(tmpDir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")

	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}

	// Create config with schema rules
	configContent := `project_name: test
version: 1
nodes_path: .deco/nodes
schema_rules:
  item:
    required_fields:
      - cost
`
	if err := os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Create a test node
	nodeContent := `id: item-1
kind: item
version: 1
status: draft
title: Test Item
custom:
  cost: 100
`
	if err := os.WriteFile(filepath.Join(nodesDir, "item-1.yaml"), []byte(nodeContent), 0644); err != nil {
		t.Fatalf("failed to write node: %v", err)
	}

	return tmpDir
}

func TestExecutor_NoMigrationNeeded(t *testing.T) {
	tmpDir := setupTestProject(t)

	// Set schema version to match
	configRepo := config.NewYAMLRepository(tmpDir)
	cfg, _ := configRepo.Load()
	cfg.SchemaVersion = ComputeSchemaHash(cfg)
	configRepo.Save(cfg)

	executor := NewExecutor(ExecutorOptions{
		TargetDir: tmpDir,
	}, NewRegistry())

	result, err := executor.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodesProcessed != 0 {
		t.Errorf("expected 0 nodes processed when no migration needed, got %d", result.NodesProcessed)
	}
}

func TestExecutor_MigrationNeeded(t *testing.T) {
	tmpDir := setupTestProject(t)

	// Config has schema rules but no schema_version set, so migration needed
	executor := NewExecutor(ExecutorOptions{
		TargetDir: tmpDir,
	}, NewRegistry())

	result, err := executor.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.SourceHash != "" {
		t.Errorf("expected empty source hash, got %q", result.SourceHash)
	}
	if result.TargetHash == "" {
		t.Error("expected non-empty target hash")
	}

	// Verify config was updated
	configRepo := config.NewYAMLRepository(tmpDir)
	cfg, _ := configRepo.Load()
	if cfg.SchemaVersion != result.TargetHash {
		t.Errorf("config schema_version not updated: got %q, want %q", cfg.SchemaVersion, result.TargetHash)
	}
}

func TestExecutor_DryRun(t *testing.T) {
	tmpDir := setupTestProject(t)

	executor := NewExecutor(ExecutorOptions{
		TargetDir: tmpDir,
		DryRun:    true,
	}, NewRegistry())

	result, err := executor.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !result.DryRun {
		t.Error("expected DryRun=true in result")
	}

	// Verify config was NOT updated
	configRepo := config.NewYAMLRepository(tmpDir)
	cfg, _ := configRepo.Load()
	if cfg.SchemaVersion != "" {
		t.Errorf("config should not be updated in dry run, got schema_version=%q", cfg.SchemaVersion)
	}
}

func TestExecutor_WithBackup(t *testing.T) {
	tmpDir := setupTestProject(t)

	executor := NewExecutor(ExecutorOptions{
		TargetDir: tmpDir,
	}, NewRegistry())

	result, err := executor.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.BackupDir == "" {
		t.Error("expected backup to be created")
	}

	// Verify backup exists
	if _, err := os.Stat(result.BackupDir); os.IsNotExist(err) {
		t.Errorf("backup directory not found: %s", result.BackupDir)
	}
}

func TestExecutor_NoBackup(t *testing.T) {
	tmpDir := setupTestProject(t)

	executor := NewExecutor(ExecutorOptions{
		TargetDir: tmpDir,
		NoBackup:  true,
	}, NewRegistry())

	result, err := executor.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.BackupDir != "" {
		t.Errorf("expected no backup, got %s", result.BackupDir)
	}
}

func TestNeedsMigration(t *testing.T) {
	tmpDir := setupTestProject(t)

	// Initially needs migration (no schema_version set)
	needs, current, expected, err := NeedsMigration(tmpDir)
	if err != nil {
		t.Fatalf("NeedsMigration failed: %v", err)
	}
	if !needs {
		t.Error("expected migration to be needed")
	}
	if current != "" {
		t.Errorf("expected empty current hash, got %q", current)
	}
	if expected == "" {
		t.Error("expected non-empty expected hash")
	}

	// Set schema version
	configRepo := config.NewYAMLRepository(tmpDir)
	cfg, _ := configRepo.Load()
	cfg.SchemaVersion = expected
	configRepo.Save(cfg)

	// Now should not need migration
	needs, _, _, err = NeedsMigration(tmpDir)
	if err != nil {
		t.Fatalf("NeedsMigration failed: %v", err)
	}
	if needs {
		t.Error("expected no migration needed after setting version")
	}
}
