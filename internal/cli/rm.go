package cli

import (
	"fmt"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/graph"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type rmFlags struct {
	force     bool
	targetDir string
}

// NewRmCommand creates the rm subcommand
func NewRmCommand() *cobra.Command {
	flags := &rmFlags{}

	cmd := &cobra.Command{
		Use:   "rm <id>",
		Short: "Delete a node",
		Long: `Delete a node from the project.

If other nodes reference this node, the command will fail unless --force is used.
The deletion is logged in the project history.

Examples:
  deco rm sword-001
  deco rm systems/combat --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRm(args[0], flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.force, "force", "f", false, "Delete even if other nodes reference this node")
	cmd.Flags().StringVarP(&flags.targetDir, "dir", "d", ".", "Project directory")

	return cmd
}

func runRm(id string, flags *rmFlags) error {
	// Verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	nodeRepo := node.NewYAMLRepository(flags.targetDir)

	// Check if node exists
	exists, err := nodeRepo.Exists(id)
	if err != nil {
		return fmt.Errorf("failed to check node existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("node %q not found", id)
	}

	// Load the node (for history logging)
	targetNode, err := nodeRepo.Load(id)
	if err != nil {
		return fmt.Errorf("failed to load node: %w", err)
	}

	// Check for reverse references
	if !flags.force {
		reverseRefs, err := findReverseRefs(nodeRepo, id)
		if err != nil {
			return fmt.Errorf("failed to check references: %w", err)
		}
		if len(reverseRefs) > 0 {
			return fmt.Errorf("node %q is referenced by: %v (use --force to delete anyway)", id, reverseRefs)
		}
	}

	// Delete the node
	if err := nodeRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	// Log deletion in history
	if err := logDeletion(flags.targetDir, targetNode); err != nil {
		// Node is already deleted, just warn
		fmt.Printf("Warning: failed to log deletion: %v\n", err)
	}

	fmt.Printf("Deleted node: %s\n", id)
	return nil
}

// findReverseRefs finds all nodes that reference the given node ID
func findReverseRefs(nodeRepo *node.YAMLRepository, targetID string) ([]string, error) {
	// Load all nodes
	nodes, err := nodeRepo.LoadAll()
	if err != nil {
		return nil, err
	}

	// Build graph
	builder := graph.NewBuilder()
	g, err := builder.Build(nodes)
	if err != nil {
		return nil, err
	}

	// Get reverse index
	reverseIndex := builder.BuildReverseIndex(g)

	// Return nodes that reference the target
	return reverseIndex[targetID], nil
}

// logDeletion adds a deletion entry to the history log with content hash
func logDeletion(targetDir string, deletedNode domain.Node) error {
	historyRepo := history.NewYAMLRepository(targetDir)

	entry := domain.AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      deletedNode.ID,
		Operation:   "delete",
		User:        GetCurrentUser(),
		ContentHash: ComputeContentHash(deletedNode),
		Before: map[string]interface{}{
			"id":      deletedNode.ID,
			"kind":    deletedNode.Kind,
			"version": deletedNode.Version,
			"status":  deletedNode.Status,
			"title":   deletedNode.Title,
		},
	}

	return historyRepo.Append(entry)
}
