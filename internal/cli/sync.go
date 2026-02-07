package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/refactor"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type syncFlags struct {
	dryRun    bool
	quiet     bool
	noRefactor bool
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

Rename Detection:
When files are manually renamed (bypassing 'deco mv'), sync detects this
by matching content hashes of orphan nodes (no history) to missing nodes
(history but no file). Detected renames automatically update references
in other nodes. Use --no-refactor to skip automatic reference updates.

Nodes without history are automatically baselined (their current state
is recorded without modification).

Exit codes:
  0 - No changes needed (or baseline only)
  1 - Files modified, re-commit needed
  2 - Error (invalid project)

Examples:
  deco sync              # Sync nodes in current directory
  deco sync --dry-run    # Show what would change
  deco sync --no-refactor # Skip automatic reference updates for renames
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
	cmd.Flags().BoolVar(&flags.noRefactor, "no-refactor", false, "Skip automatic reference updates for detected renames")

	return cmd
}

// renameDetection holds info about a detected manual rename
type renameDetection struct {
	oldID       string
	newID       string
	contentHash string
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

	// Phase 1: Load all nodes and detect renames
	var allNodes []domain.Node
	loadedNodeIDs := make(map[string]bool)
	orphanNodes := make(map[string]domain.Node) // nodes with no history (keyed by content hash)
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

		allNodes = append(allNodes, currentNode)
		loadedNodeIDs[currentNode.ID] = true

		// Check if this node has history
		if latestHashes[currentNode.ID] == "" {
			// No history - this is either a new node or a renamed node
			contentHash := ComputeContentHashWithDir(currentNode, flags.targetDir)
			orphanNodes[contentHash] = currentNode
		}
	}

	// Find missing node IDs (history entries with no matching file)
	missingNodeHashes := make(map[string]string) // oldID -> last content hash
	for nodeID, hash := range latestHashes {
		if !loadedNodeIDs[nodeID] {
			missingNodeHashes[nodeID] = hash
		}
	}

	// Detect renames by matching orphan content hashes to missing node hashes
	var detectedRenames []renameDetection
	for oldID, oldHash := range missingNodeHashes {
		if orphanNode, found := orphanNodes[oldHash]; found {
			// Content hash matches - this is a rename!
			detectedRenames = append(detectedRenames, renameDetection{
				oldID:       oldID,
				newID:       orphanNode.ID,
				contentHash: oldHash,
			})
			// Remove from orphan map so it won't be baselined
			delete(orphanNodes, oldHash)
			// Remove from missing so it won't be treated as deletion
			delete(missingNodeHashes, oldID)
		}
	}

	// Detect deletions: remaining missing nodes (history but no file, not renamed)
	var deletedNodes []string
	for deletedID := range missingNodeHashes {
		deletedNodes = append(deletedNodes, deletedID)
		if !flags.dryRun {
			if err := logDeleteOperation(historyPath, deletedID); err != nil {
				errors = append(errors, fmt.Sprintf("failed to log deletion of %s: %v", deletedID, err))
				if !flags.quiet {
					fmt.Fprintf(os.Stderr, "Error: failed to log deletion of %s: %v\n", deletedID, err)
				}
			}
		}
	}

	// Phase 2: Apply rename refactoring
	var renameResults []string
	var refUpdatedNodes map[string]bool // track which nodes had refs updated

	if len(detectedRenames) > 0 && !flags.noRefactor {
		refUpdatedNodes = make(map[string]bool)
		renamer := refactor.NewRenamer()

		for _, rename := range detectedRenames {
			// Update references across all nodes
			updatedNodes, err := renamer.UpdateReferences(allNodes, rename.oldID, rename.newID)
			if err != nil {
				errors = append(errors, fmt.Sprintf("failed to update references for rename %s→%s: %v", rename.oldID, rename.newID, err))
				if !flags.quiet {
					fmt.Fprintf(os.Stderr, "Error: failed to update references for rename %s→%s: %v\n", rename.oldID, rename.newID, err)
				}
				continue
			}

			// Find which nodes were updated (version changed)
			var updatedIDs []string
			for i, updated := range updatedNodes {
				if updated.Version != allNodes[i].Version {
					updatedIDs = append(updatedIDs, updated.ID)
					refUpdatedNodes[updated.ID] = true
				}
			}

			// Apply updates to allNodes for subsequent renames
			allNodes = updatedNodes

			if !flags.dryRun {
				// Save nodes whose references were updated
				for _, updated := range updatedNodes {
					if refUpdatedNodes[updated.ID] {
						if err := nodeRepo.Save(updated); err != nil {
							errors = append(errors, fmt.Sprintf("failed to save %s: %v", updated.ID, err))
							if !flags.quiet {
								fmt.Fprintf(os.Stderr, "Error: failed to save %s: %v\n", updated.ID, err)
							}
						}
					}
				}

				// Log move operation (like deco mv does)
				if err := logMoveOperation(historyPath, rename.oldID, rename.newID, rename.contentHash); err != nil {
					errors = append(errors, fmt.Sprintf("failed to log rename %s→%s: %v", rename.oldID, rename.newID, err))
					if !flags.quiet {
						fmt.Fprintf(os.Stderr, "Error: failed to log rename %s→%s: %v\n", rename.oldID, rename.newID, err)
					}
				}
			}

			if len(updatedIDs) > 0 {
				renameResults = append(renameResults, fmt.Sprintf("%s→%s (updated refs in: %s)", rename.oldID, rename.newID, strings.Join(updatedIDs, ", ")))
			} else {
				renameResults = append(renameResults, fmt.Sprintf("%s→%s (no refs to update)", rename.oldID, rename.newID))
			}
		}
	} else if len(detectedRenames) > 0 && flags.noRefactor {
		// Just report detected renames without applying refactor
		for _, rename := range detectedRenames {
			renameResults = append(renameResults, fmt.Sprintf("%s→%s (--no-refactor: refs not updated)", rename.oldID, rename.newID))
		}
	}

	// Phase 3: Normal sync/baseline for remaining nodes
	var syncResults []syncResult
	var baselinedNodes []string

	for _, currentNode := range allNodes {
		currentHash := ComputeContentHashWithDir(currentNode, flags.targetDir)
		lastHash := latestHashes[currentNode.ID]

		if lastHash == "" {
			// No history - check if this was a detected rename (already handled)
			wasRenamed := false
			for _, rename := range detectedRenames {
				if rename.newID == currentNode.ID {
					wasRenamed = true
					break
				}
			}
			if wasRenamed {
				// Already logged as move operation, skip baseline
				continue
			}

			// Genuine new node - baseline it
			if !flags.dryRun {
				if err := logBaselineOperation(historyPath, currentNode.ID, currentHash); err != nil {
					errors = append(errors, fmt.Sprintf("failed to baseline %s: %v", currentNode.ID, err))
					if !flags.quiet {
						fmt.Fprintf(os.Stderr, "Error: failed to baseline %s: %v\n", currentNode.ID, err)
					}
					continue
				}
			}
			baselinedNodes = append(baselinedNodes, currentNode.ID)
			continue
		}

		// Skip nodes that had their refs updated by rename refactoring
		// (they've already been saved with version bump)
		if refUpdatedNodes != nil && refUpdatedNodes[currentNode.ID] {
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
			nodeCopy := currentNode // copy for modification
			if err := applySyncWithHash(historyPath, &nodeCopy, nodeRepo, currentHash); err != nil {
				errors = append(errors, fmt.Sprintf("failed to sync %s: %v", currentNode.ID, err))
				if !flags.quiet {
					fmt.Fprintf(os.Stderr, "Error: failed to sync %s: %v\n", currentNode.ID, err)
				}
				continue
			}
		}

		syncResults = append(syncResults, result)
	}

	// Output results
	if !flags.quiet {
		if len(deletedNodes) > 0 {
			if flags.dryRun {
				fmt.Printf("Would mark deleted: %s (%d nodes)\n", strings.Join(deletedNodes, ", "), len(deletedNodes))
			} else {
				fmt.Printf("Deleted: %s (%d nodes)\n", strings.Join(deletedNodes, ", "), len(deletedNodes))
			}
		}

		if len(renameResults) > 0 {
			if flags.dryRun {
				fmt.Println("Would apply renames:")
			} else {
				fmt.Println("Detected renames:")
			}
			for _, r := range renameResults {
				fmt.Printf("  %s\n", r)
			}
		}

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
	if len(syncResults) > 0 || len(renameResults) > 0 || len(deletedNodes) > 0 {
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

// logMoveOperation records a rename detected during sync (manual rename)
func logMoveOperation(historyPath, oldID, newID, contentHash string) error {
	historyRepo := history.NewYAMLRepository(historyPath)

	entry := domain.AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      newID,
		Operation:   "move",
		User:        GetCurrentUser(),
		ContentHash: contentHash,
		Before: map[string]interface{}{
			"id": oldID,
		},
		After: map[string]interface{}{
			"id": newID,
		},
	}

	return historyRepo.Append(entry)
}

// logDeleteOperation records a node deletion detected during sync
func logDeleteOperation(historyPath, nodeID string) error {
	historyRepo := history.NewYAMLRepository(historyPath)

	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    nodeID,
		Operation: "delete",
		User:      GetCurrentUser(),
	}

	return historyRepo.Append(entry)
}
