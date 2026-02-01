package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
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

	// Verify we're in a git repository
	if !isGitRepo(flags.targetDir) {
		return syncExitError, fmt.Errorf("not a git repository")
	}

	// Get modified node files
	modifiedFiles, err := getModifiedNodeFiles(flags.targetDir)
	if err != nil {
		return syncExitError, fmt.Errorf("failed to get modified files: %w", err)
	}

	if len(modifiedFiles) == 0 {
		return syncExitClean, nil
	}

	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	var results []syncResult

	for _, filePath := range modifiedFiles {
		// Load current node from working tree
		nodeID := extractNodeID(filePath)
		if nodeID == "" {
			continue
		}

		currentNode, err := nodeRepo.Load(nodeID)
		if err != nil {
			continue // Skip nodes that can't be loaded
		}

		// Get HEAD version of the node
		headNode, err := getNodeFromHEAD(flags.targetDir, filePath)
		if err != nil {
			continue // New file or can't read HEAD version
		}

		// Compare content semantically
		if !contentChanged(headNode, currentNode) {
			continue // Only metadata changed, skip
		}

		// Content changed - need to sync
		result := syncResult{
			nodeID:     currentNode.ID,
			oldVersion: currentNode.Version,
			newVersion: currentNode.Version + 1,
			oldStatus:  currentNode.Status,
		}

		if !flags.dryRun {
			// Apply sync changes
			if err := applySync(flags.targetDir, &currentNode, nodeRepo); err != nil {
				if !flags.quiet {
					fmt.Fprintf(os.Stderr, "Warning: failed to sync %s: %v\n", nodeID, err)
				}
				continue
			}
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		return syncExitClean, nil
	}

	// Output results
	if !flags.quiet {
		if flags.dryRun {
			fmt.Print("Would sync: ")
		} else {
			fmt.Print("Synced: ")
		}

		parts := make([]string, len(results))
		for i, r := range results {
			parts[i] = fmt.Sprintf("%s (v%d→v%d)", r.nodeID, r.oldVersion, r.newVersion)
		}
		fmt.Println(strings.Join(parts, ", "))
	}

	if flags.dryRun {
		return syncExitClean, nil
	}

	return syncExitModified, nil
}

// isGitRepo checks if the directory is inside a git repository
func isGitRepo(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	return cmd.Run() == nil
}

// getModifiedNodeFiles returns list of node files modified since HEAD
func getModifiedNodeFiles(targetDir string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "HEAD")
	cmd.Dir = targetDir
	output, err := cmd.Output()
	if err != nil {
		// If HEAD doesn't exist (new repo), there are no changes to sync
		return nil, nil
	}

	var nodeFiles []string
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Only include .yaml files in .deco/nodes/
		if strings.HasPrefix(line, ".deco/nodes/") && strings.HasSuffix(line, ".yaml") {
			nodeFiles = append(nodeFiles, line)
		}
	}

	return nodeFiles, nil
}

// extractNodeID gets the node ID from a file path like .deco/nodes/sword-001.yaml
func extractNodeID(filePath string) string {
	base := filepath.Base(filePath)
	if !strings.HasSuffix(base, ".yaml") {
		return ""
	}
	return strings.TrimSuffix(base, ".yaml")
}

// getNodeFromHEAD loads a node from the HEAD commit
func getNodeFromHEAD(targetDir, filePath string) (domain.Node, error) {
	cmd := exec.Command("git", "show", "HEAD:"+filePath)
	cmd.Dir = targetDir
	output, err := cmd.Output()
	if err != nil {
		return domain.Node{}, err
	}

	var n domain.Node
	if err := yaml.Unmarshal(output, &n); err != nil {
		return domain.Node{}, err
	}

	return n, nil
}

// contentChanged compares content fields between two nodes
// Returns true if any content field differs (ignoring metadata)
func contentChanged(old, new domain.Node) bool {
	// Content fields that trigger sync:
	// title, summary, content, tags, refs, issues

	if old.Title != new.Title {
		return true
	}
	if old.Summary != new.Summary {
		return true
	}
	if !tagsEqual(old.Tags, new.Tags) {
		return true
	}
	if !refsEqual(old.Refs, new.Refs) {
		return true
	}
	if !issuesEqual(old.Issues, new.Issues) {
		return true
	}
	if !contentEqual(old.Content, new.Content) {
		return true
	}

	return false
}

// tagsEqual compares two tag slices
func tagsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// refsEqual compares two Ref structs
func refsEqual(a, b domain.Ref) bool {
	return reflect.DeepEqual(a, b)
}

// issuesEqual compares two Issue slices
func issuesEqual(a, b []domain.Issue) bool {
	return reflect.DeepEqual(a, b)
}

// contentEqual compares two Content structs
func contentEqual(a, b *domain.Content) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	// Use YAML serialization for deep comparison
	aYAML, err1 := yaml.Marshal(a)
	bYAML, err2 := yaml.Marshal(b)
	if err1 != nil || err2 != nil {
		return false
	}
	return bytes.Equal(aYAML, bYAML)
}

// applySync applies sync changes to a node
func applySync(targetDir string, n *domain.Node, nodeRepo *node.YAMLRepository) error {
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

	// Log to history
	return logSyncOperation(targetDir, n.ID, oldVersion, n.Version, oldStatus, n.Status)
}

// logSyncOperation adds a sync entry to the history log
func logSyncOperation(targetDir, nodeID string, oldVersion, newVersion int, oldStatus, newStatus string) error {
	historyRepo := history.NewYAMLRepository(targetDir)

	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    nodeID,
		Operation: "sync",
		User:      username,
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
