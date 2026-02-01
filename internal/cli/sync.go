package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
review status may become stale. This command detects such changes and:

1. Bumps the version number
2. Resets status to "draft" (if was approved/review)
3. Clears reviewers
4. Logs the sync operation to history

Use as a pre-commit hook to ensure all changes are properly versioned.

Exit codes:
  0 - No changes needed
  1 - Files modified, re-commit needed
  2 - Error (not a git repo, invalid nodes)

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
	_, err := configRepo.Load()
	if err != nil {
		return syncExitError, fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Discover all nodes
	discovery := node.NewFileDiscovery(flags.targetDir)
	nodePaths, err := discovery.DiscoverAll()
	if err != nil {
		return syncExitError, fmt.Errorf("failed to discover nodes: %w", err)
	}

	if len(nodePaths) == 0 {
		return syncExitClean, nil
	}

	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	var syncResults []syncResult
	var baselinedNodes []string

	for _, nodePath := range nodePaths {
		nodeID := discovery.PathToID(nodePath)
		if nodeID == "" {
			continue
		}

		currentNode, err := nodeRepo.Load(nodeID)
		if err != nil {
			continue
		}

		currentHash := computeContentHash(currentNode)
		lastHash := getLastContentHash(flags.targetDir, nodeID)

		if lastHash == "" {
			// No history - baseline this node
			if !flags.dryRun {
				if err := logBaselineOperation(flags.targetDir, nodeID, currentHash); err != nil {
					if !flags.quiet {
						fmt.Fprintf(os.Stderr, "Warning: failed to baseline %s: %v\n", nodeID, err)
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
			if err := applySyncWithHash(flags.targetDir, &currentNode, nodeRepo, currentHash); err != nil {
				if !flags.quiet {
					fmt.Fprintf(os.Stderr, "Warning: failed to sync %s: %v\n", nodeID, err)
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

	if flags.dryRun {
		return syncExitClean, nil
	}

	if len(syncResults) > 0 {
		return syncExitModified, nil
	}

	return syncExitClean, nil
}

// logBaselineOperation records initial state for a node without modification
func logBaselineOperation(targetDir, nodeID, contentHash string) error {
	historyRepo := history.NewYAMLRepository(targetDir)

	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	entry := domain.AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      nodeID,
		Operation:   "baseline",
		User:        username,
		ContentHash: contentHash,
	}

	return historyRepo.Append(entry)
}

// applySyncWithHash applies sync changes and logs with content hash
func applySyncWithHash(targetDir string, n *domain.Node, nodeRepo *node.YAMLRepository, contentHash string) error {
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
	return logSyncOperationWithHash(targetDir, n.ID, oldVersion, n.Version, oldStatus, n.Status, contentHash)
}

// logSyncOperationWithHash adds a sync entry with content hash
func logSyncOperationWithHash(targetDir, nodeID string, oldVersion, newVersion int, oldStatus, newStatus, contentHash string) error {
	historyRepo := history.NewYAMLRepository(targetDir)

	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	entry := domain.AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      nodeID,
		Operation:   "sync",
		User:        username,
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

// contentFields holds only the fields that affect content hash
type contentFields struct {
	Title   string          `yaml:"title"`
	Summary string          `yaml:"summary"`
	Tags    []string        `yaml:"tags,omitempty"`
	Refs    domain.Ref      `yaml:"refs,omitempty"`
	Issues  []domain.Issue  `yaml:"issues,omitempty"`
	Content *domain.Content `yaml:"content,omitempty"`
}

// computeContentHash computes a SHA-256 hash of the content fields
// Returns 16 hex characters (first 64 bits of the hash)
func computeContentHash(n domain.Node) string {
	fields := contentFields{
		Title:   n.Title,
		Summary: n.Summary,
		Tags:    n.Tags,
		Refs:    n.Refs,
		Issues:  n.Issues,
		Content: n.Content,
	}

	data, err := yaml.Marshal(fields)
	if err != nil {
		return ""
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:8])
}

// getLastContentHash retrieves the most recent content hash for a node from history
// Returns empty string if no hash found
func getLastContentHash(targetDir, nodeID string) string {
	historyRepo := history.NewYAMLRepository(targetDir)

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
