package cli

import (
	"fmt"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/refactor"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type mvFlags struct {
	quiet     bool
	targetDir string
	oldID     string
	newID     string
}

// NewMvCommand creates the mv subcommand
func NewMvCommand() *cobra.Command {
	flags := &mvFlags{}

	cmd := &cobra.Command{
		Use:   "mv <old-id> <new-id> [directory]",
		Short: "Rename a node and update all references",
		Long: `Rename a node ID and automatically update all references.

This command:
  - Changes the node's ID
  - Renames the node file
  - Updates all Uses and Related references pointing to it
  - Increments version on nodes whose references were updated
  - Records the operation in history

Examples:
  deco mv sword-001 blade-001
  deco mv old-character new-character ./my-project`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.oldID = args[0]
			flags.newID = args[1]
			if len(args) > 2 {
				flags.targetDir = args[2]
			} else {
				flags.targetDir = "."
			}
			return runMv(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")

	return cmd
}

func runMv(flags *mvFlags) error {
	// Validate inputs early
	if flags.oldID == "" {
		return fmt.Errorf("old ID cannot be empty")
	}
	if flags.newID == "" {
		return fmt.Errorf("new ID cannot be empty")
	}
	if flags.oldID == flags.newID {
		return fmt.Errorf("new ID must be different from old ID")
	}

	// Load config to verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load all nodes
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	nodes, err := nodeRepo.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}

	// Use the Renamer service
	renamer := refactor.NewRenamer()
	updatedNodes, err := renamer.Rename(nodes, flags.oldID, flags.newID)
	if err != nil {
		return fmt.Errorf("rename failed: %w", err)
	}

	// Build a map of original nodes for comparison
	originalByID := make(map[string]domain.Node)
	for _, n := range nodes {
		originalByID[n.ID] = n
	}

	// Find which nodes changed
	var modifiedCount int
	var renamedNode domain.Node
	for _, updated := range updatedNodes {
		// The renamed node: save with new ID, delete old file
		if updated.ID == flags.newID {
			renamedNode = updated
			// Save new file
			if err := nodeRepo.Save(updated); err != nil {
				return fmt.Errorf("failed to save renamed node: %w", err)
			}
			// Delete old file
			if err := nodeRepo.Delete(flags.oldID); err != nil {
				return fmt.Errorf("failed to delete old node file: %w", err)
			}
			modifiedCount++
			continue
		}

		// Check if this node's version changed (refs were updated)
		if orig, ok := originalByID[updated.ID]; ok {
			if updated.Version != orig.Version {
				if err := nodeRepo.Save(updated); err != nil {
					return fmt.Errorf("failed to save node %s: %w", updated.ID, err)
				}
				modifiedCount++
			}
		}
	}

	// Record history entry with content hash
	historyRepo := history.NewYAMLRepository(flags.targetDir)
	entry := domain.AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      flags.newID,
		Operation:   "move",
		User:        GetCurrentUser(),
		ContentHash: ComputeContentHash(renamedNode),
		Before:      map[string]interface{}{"id": flags.oldID},
		After:       map[string]interface{}{"id": flags.newID},
	}
	if err := historyRepo.Append(entry); err != nil {
		return fmt.Errorf("failed to record history: %w", err)
	}

	if !flags.quiet {
		fmt.Printf("Renamed %s -> %s (%d node(s) updated)\n", flags.oldID, flags.newID, modifiedCount)
	}

	return nil
}
