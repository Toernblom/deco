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
	"fmt"
	"os"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
)

// ExecutorOptions configures the migration executor.
type ExecutorOptions struct {
	// DryRun if true, shows what would change without making changes.
	DryRun bool
	// NoBackup if true, skips creating a backup before migration.
	NoBackup bool
	// Quiet if true, suppresses non-essential output.
	Quiet bool
	// TargetDir is the project root directory.
	TargetDir string
}

// ExecutorResult contains the result of a migration execution.
type ExecutorResult struct {
	// NodesProcessed is the number of nodes processed.
	NodesProcessed int
	// NodesModified is the number of nodes that were actually changed.
	NodesModified int
	// BackupDir is the path to the backup (if created).
	BackupDir string
	// SourceHash is the schema hash before migration.
	SourceHash string
	// TargetHash is the schema hash after migration.
	TargetHash string
	// DryRun indicates if this was a dry run.
	DryRun bool
}

// Executor performs schema migrations.
type Executor struct {
	opts     ExecutorOptions
	registry *Registry
}

// NewExecutor creates a new migration executor.
func NewExecutor(opts ExecutorOptions, registry *Registry) *Executor {
	if registry == nil {
		registry = DefaultRegistry
	}
	return &Executor{
		opts:     opts,
		registry: registry,
	}
}

// Execute runs the migration from current schema to target schema.
func (e *Executor) Execute() (*ExecutorResult, error) {
	// Load config
	configRepo := config.NewYAMLRepository(e.opts.TargetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Compute current and target schema hashes
	sourceHash := cfg.SchemaVersion
	targetHash := ComputeSchemaHash(cfg)

	// Check if migration is needed
	if sourceHash == targetHash {
		return &ExecutorResult{
			SourceHash: sourceHash,
			TargetHash: targetHash,
			DryRun:     e.opts.DryRun,
		}, nil
	}

	// Find migration path
	migrations := e.registry.FindPath(sourceHash, targetHash)

	// If no registered migrations, use identity migration
	// (just update schema version without transforming nodes)
	if migrations == nil {
		migrations = []Migration{IdentityMigration(
			"auto-update",
			"Automatic schema version update",
			sourceHash,
			targetHash,
		)}
	}

	// Create backup unless disabled or dry run
	var backupDir string
	if !e.opts.NoBackup && !e.opts.DryRun {
		backup, err := CreateBackup(e.opts.TargetDir)
		if err != nil {
			return nil, fmt.Errorf("failed to create backup: %w", err)
		}
		backupDir = backup.BackupDir
	}

	// Load all nodes
	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, e.opts.TargetDir))
	nodes, err := nodeRepo.LoadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load nodes: %w", err)
	}

	// Apply migrations
	result := &ExecutorResult{
		NodesProcessed: len(nodes),
		BackupDir:      backupDir,
		SourceHash:     sourceHash,
		TargetHash:     targetHash,
		DryRun:         e.opts.DryRun,
	}

	modifiedNodes := make([]domain.Node, 0)

	for _, n := range nodes {
		modified := false
		currentNode := n

		// Apply each migration in sequence
		for _, m := range migrations {
			if m.Transform != nil {
				transformed, err := m.Transform(currentNode)
				if err != nil {
					return nil, fmt.Errorf("migration %q failed for node %s: %w", m.Name, n.ID, err)
				}
				// Check if node was actually modified
				if !nodesEqual(currentNode, transformed) {
					modified = true
					currentNode = transformed
				}
			}
		}

		if modified {
			modifiedNodes = append(modifiedNodes, currentNode)
			result.NodesModified++
		}
	}

	// Skip writes for dry run
	if e.opts.DryRun {
		return result, nil
	}

	// Save modified nodes
	for _, n := range modifiedNodes {
		// Increment version for modified nodes
		n.Version++
		if err := nodeRepo.Save(n); err != nil {
			return nil, fmt.Errorf("failed to save node %s: %w", n.ID, err)
		}
	}

	// Log migration to audit history
	historyRepo := history.NewYAMLRepository(config.ResolveHistoryPath(cfg, e.opts.TargetDir))
	user := getUser()

	for _, n := range modifiedNodes {
		entry := domain.AuditEntry{
			Timestamp: time.Now(),
			NodeID:    n.ID,
			Operation: "migrate",
			User:      user,
		}
		if err := historyRepo.Append(entry); err != nil {
			// Log error but don't fail migration
			fmt.Fprintf(os.Stderr, "Warning: failed to log audit entry for %s: %v\n", n.ID, err)
		}
	}

	// Update config with new schema version
	cfg.SchemaVersion = targetHash
	if err := configRepo.Save(cfg); err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	return result, nil
}

// nodesEqual compares two nodes for equality (simplified).
// A full implementation would do deep comparison.
func nodesEqual(a, b domain.Node) bool {
	// For now, compare by ID and version - transform should modify something
	// that we can detect. In practice, transforms that don't change anything
	// return the same node.
	return a.ID == b.ID && a.Version == b.Version && a.Title == b.Title
}

// getUser returns the current user for audit logging.
func getUser() string {
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	if user == "" {
		user = "unknown"
	}
	return user
}

// NeedsMigration checks if the project needs schema migration.
func NeedsMigration(targetDir string) (bool, string, string, error) {
	configRepo := config.NewYAMLRepository(targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return false, "", "", err
	}

	currentHash := cfg.SchemaVersion
	expectedHash := ComputeSchemaHash(cfg)

	return currentHash != expectedHash, currentHash, expectedHash, nil
}
