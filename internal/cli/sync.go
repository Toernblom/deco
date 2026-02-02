package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type syncFlags struct {
	dryRun    bool
	quiet     bool
	targetDir string
}

// Exit codes for sync command
const (
	syncExitClean    = 0 // No changes needed
	syncExitModified = 1 // Files were modified, re-commit needed
	syncExitError    = 2 // Error occurred
)

// NewSyncCommand creates the sync subcommand
func NewSyncCommand() *cobra.Command {
	flags := &syncFlags{}

	cmd := &cobra.Command{
		Use:   "sync [directory]",
		Short: "Detect and fix unversioned node changes",
		Long: `Detect manually-edited nodes and fix their metadata.

When nodes are edited directly (bypassing the CLI), their version and
review status may become stale. This command detects such changes by
comparing content hashes and:

1. Bumps the version number
2. Resets status to "draft" (if was approved/review)
3. Clears reviewers
4. Logs the sync operation to history

Nodes without history are automatically baselined (their current state
is recorded without modification).

Exit codes:
  0 - No changes needed (or baseline only)
  1 - Files modified, re-commit needed
  2 - Error (invalid project)

Examples:
  deco sync              # Sync nodes in current directory
  deco sync --dry-run    # Show what would change
  deco sync /path/to/project`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				flags.targetDir = args[0]
			} else {
				flags.targetDir = "."
			}
			exitCode, err := runSync(flags)
			if err != nil {
				os.Exit(syncExitError)
			}
			os.Exit(exitCode)
			return nil
		},
	}

	cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Show what would change without applying")
	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")

	return cmd
}

// syncResult holds info about a synced node
type syncResult struct {
	nodeID     string
	oldVersion int
	newVersion int
	oldStatus  string
}

func runSync(flags *syncFlags) (int, error) {
	// Verify we're in a deco project
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return syncExitError, fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Discover all nodes
	nodesPath := config.ResolveNodesPath(cfg, flags.targetDir)
	discovery := node.NewFileDiscovery(nodesPath)
	nodePaths, err := discovery.DiscoverAll()
	if err != nil {
		return syncExitError, fmt.Errorf("failed to discover nodes: %w", err)
	}

	if len(nodePaths) == 0 {
		return syncExitClean, nil
	}

	// Load all latest content hashes in a single pass (O(history) instead of O(nodes × history))
	historyPath := config.ResolveHistoryPath(cfg, flags.targetDir)
	historyRepo := history.NewYAMLRepository(historyPath)
	latestHashes, err := historyRepo.QueryLatestHashes()
	if err != nil {
		return syncExitError, fmt.Errorf("failed to load history: %w", err)
	}

	nodeRepo := node.NewYAMLRepository(nodesPath)
	var syncResults []syncResult
	var baselinedNodes []string
	var errors []string

	for _, nodePath := range nodePaths {
		nodeID := discovery.PathToID(nodePath)
		if nodeID == "" {
			continue
		}

		currentNode, err := nodeRepo.Load(nodeID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to load %s: %v", nodeID, err))
			if !flags.quiet {
				fmt.Fprintf(os.Stderr, "Error: failed to load %s: %v\n", nodeID, err)
			}
			continue
		}

		currentHash := ComputeContentHash(currentNode)
		lastHash := latestHashes[nodeID]

		if lastHash == "" {
			// No history - baseline this node
			if !flags.dryRun {
				if err := logBaselineOperation(historyPath, nodeID, currentHash); err != nil {
					errors = append(errors, fmt.Sprintf("failed to baseline %s: %v", nodeID, err))
					if !flags.quiet {
						fmt.Fprintf(os.Stderr, "Error: failed to baseline %s: %v\n", nodeID, err)
					}
					continue
				}
			}
			baselinedNodes = append(baselinedNodes, nodeID)
			continue
		}

		if currentHash == lastHash {
			// No change
			continue
		}

		// Content changed - need to sync
		result := syncResult{
			nodeID:     currentNode.ID,
			oldVersion: currentNode.Version,
			newVersion: currentNode.Version + 1,
			oldStatus:  currentNode.Status,
		}

		if !flags.dryRun {
			if err := applySyncWithHash(historyPath, &currentNode, nodeRepo, currentHash); err != nil {
				errors = append(errors, fmt.Sprintf("failed to sync %s: %v", nodeID, err))
				if !flags.quiet {
					fmt.Fprintf(os.Stderr, "Error: failed to sync %s: %v\n", nodeID, err)
				}
				continue
			}
		}

		syncResults = append(syncResults, result)
	}

	// Output results
	if !flags.quiet {
		if len(baselinedNodes) > 0 {
			if flags.dryRun {
				fmt.Printf("Would baseline: %s (%d nodes)\n", strings.Join(baselinedNodes, ", "), len(baselinedNodes))
			} else {
				fmt.Printf("Baselined: %s (%d nodes)\n", strings.Join(baselinedNodes, ", "), len(baselinedNodes))
			}
		}

		if len(syncResults) > 0 {
			if flags.dryRun {
				fmt.Print("Would sync: ")
			} else {
				fmt.Print("Synced: ")
			}
			parts := make([]string, len(syncResults))
			for i, r := range syncResults {
				parts[i] = fmt.Sprintf("%s (v%d→v%d)", r.nodeID, r.oldVersion, r.newVersion)
			}
			fmt.Println(strings.Join(parts, ", "))
		}
	}

	// Report errors if any occurred
	if len(errors) > 0 {
		if !flags.quiet {
			fmt.Fprintf(os.Stderr, "\n%d error(s) occurred during sync\n", len(errors))
		}
		return syncExitError, fmt.Errorf("%d sync error(s)", len(errors))
	}

	// Return modified exit code if changes would be (or were) made
	if len(syncResults) > 0 {
		return syncExitModified, nil
	}

	return syncExitClean, nil
}

// logBaselineOperation records initial state for a node without modification
func logBaselineOperation(historyPath, nodeID, contentHash string) error {
	historyRepo := history.NewYAMLRepository(historyPath)

	entry := domain.AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      nodeID,
		Operation:   "baseline",
		User:        GetCurrentUser(),
		ContentHash: contentHash,
	}

	return historyRepo.Append(entry)
}

// applySyncWithHash applies sync changes and logs with content hash
func applySyncWithHash(historyPath string, n *domain.Node, nodeRepo *node.YAMLRepository, contentHash string) error {
	oldVersion := n.Version
	oldStatus := n.Status

	// Bump version
	n.Version++

	// Reset status if was approved or review
	if n.Status == "approved" || n.Status == "review" {
		n.Status = "draft"
	}

	// Clear reviewers
	n.Reviewers = nil

	// Save the node
	if err := nodeRepo.Save(*n); err != nil {
		return err
	}

	// Log to history with content hash
	return logSyncOperationWithHash(historyPath, n.ID, oldVersion, n.Version, oldStatus, n.Status, contentHash)
}

// logSyncOperationWithHash adds a sync entry with content hash
func logSyncOperationWithHash(historyPath, nodeID string, oldVersion, newVersion int, oldStatus, newStatus, contentHash string) error {
	historyRepo := history.NewYAMLRepository(historyPath)

	entry := domain.AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      nodeID,
		Operation:   "sync",
		User:        GetCurrentUser(),
		ContentHash: contentHash,
		Before: map[string]interface{}{
			"version": oldVersion,
			"status":  oldStatus,
		},
		After: map[string]interface{}{
			"version": newVersion,
			"status":  newStatus,
		},
	}

	return historyRepo.Append(entry)
}

// getLastContentHash retrieves the most recent content hash for a node from history
// Returns empty string if no hash found
// historyPath should be the full path to the history file (use config.ResolveHistoryPath)
func getLastContentHash(historyPath, nodeID string) string {
	historyRepo := history.NewYAMLRepository(historyPath)

	entries, err := historyRepo.Query(history.Filter{NodeID: nodeID})
	if err != nil || len(entries) == 0 {
		return ""
	}

	// Entries are sorted oldest first, get the last one
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].ContentHash != "" {
			return entries[i].ContentHash
		}
	}

	return ""
}
