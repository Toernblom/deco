package cli

import (
	"fmt"
	"os/user"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type reviewFlags struct {
	targetDir string
	nodeID    string
	quiet     bool
}

// NewReviewCommand creates the review command with subcommands
func NewReviewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Manage review workflow for nodes",
		Long: `Manage review workflow for nodes.

Subcommands:
  submit   - Submit a draft node for review
  approve  - Approve a node under review
  reject   - Reject a node back to draft
  status   - Show review status of a node`,
	}

	cmd.AddCommand(newSubmitCommand())
	cmd.AddCommand(newApproveCommand())
	cmd.AddCommand(newRejectCommand())

	return cmd
}

func newSubmitCommand() *cobra.Command {
	flags := &reviewFlags{}

	cmd := &cobra.Command{
		Use:   "submit <node-id> [directory]",
		Short: "Submit a draft node for review",
		Long:  `Submit a draft node for review. Changes status from 'draft' to 'review'.`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.nodeID = args[0]
			if len(args) > 1 {
				flags.targetDir = args[1]
			} else {
				flags.targetDir = "."
			}
			return runSubmit(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")

	return cmd
}

func runSubmit(flags *reviewFlags) error {
	// Load config
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	n, err := nodeRepo.Load(flags.nodeID)
	if err != nil {
		return fmt.Errorf("node %q not found: %w", flags.nodeID, err)
	}

	// Validate current status
	if n.Status != "draft" {
		return fmt.Errorf("cannot submit node %q: status is %q, must be 'draft'", flags.nodeID, n.Status)
	}

	// Update status
	oldStatus := n.Status
	n.Status = "review"

	// Save the node
	if err := nodeRepo.Save(n); err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	// Log submit operation
	if err := logReviewOperation(flags.targetDir, n.ID, "submit", oldStatus, n.Status, ""); err != nil {
		fmt.Printf("Warning: failed to log submit operation: %v\n", err)
	}

	if !flags.quiet {
		fmt.Printf("Submitted %s for review (status: draft -> review)\n", flags.nodeID)
	}

	return nil
}

func logReviewOperation(targetDir, nodeID, operation, oldStatus, newStatus, note string) error {
	historyRepo := history.NewYAMLRepository(targetDir)

	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    nodeID,
		Operation: operation,
		User:      username,
		Before:    map[string]interface{}{"status": oldStatus},
		After:     map[string]interface{}{"status": newStatus},
	}

	if note != "" {
		entry.After["note"] = note
	}

	return historyRepo.Append(entry)
}

type approveFlags struct {
	targetDir string
	nodeID    string
	note      string
	quiet     bool
}

func newApproveCommand() *cobra.Command {
	flags := &approveFlags{}

	cmd := &cobra.Command{
		Use:   "approve <node-id> [directory]",
		Short: "Approve a node under review",
		Long:  `Approve a node under review. Adds your approval and transitions to 'approved' if requirements are met.`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.nodeID = args[0]
			if len(args) > 1 {
				flags.targetDir = args[1]
			} else {
				flags.targetDir = "."
			}
			return runApprove(flags)
		},
	}

	cmd.Flags().StringVar(&flags.note, "note", "", "Optional approval note")
	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")

	return cmd
}

func runApprove(flags *approveFlags) error {
	// Load config
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	n, err := nodeRepo.Load(flags.nodeID)
	if err != nil {
		return fmt.Errorf("node %q not found: %w", flags.nodeID, err)
	}

	// Validate current status
	if n.Status != "review" {
		return fmt.Errorf("cannot approve node %q: status is %q, must be 'review'", flags.nodeID, n.Status)
	}

	// Get current user
	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	// Add reviewer
	reviewer := domain.Reviewer{
		Name:      username,
		Timestamp: time.Now(),
		Version:   n.Version,
		Note:      flags.note,
	}
	n.Reviewers = append(n.Reviewers, reviewer)

	// Count approvals for current version
	validApprovals := 0
	for _, r := range n.Reviewers {
		if r.Version == n.Version {
			validApprovals++
		}
	}

	// Transition to approved if requirements met
	oldStatus := n.Status
	if validApprovals >= cfg.RequiredApprovals {
		n.Status = "approved"
	}

	// Save the node
	if err := nodeRepo.Save(n); err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	// Log approve operation
	if err := logReviewOperation(flags.targetDir, n.ID, "approve", oldStatus, n.Status, flags.note); err != nil {
		fmt.Printf("Warning: failed to log approve operation: %v\n", err)
	}

	if !flags.quiet {
		if n.Status == "approved" {
			fmt.Printf("Approved %s (status: review -> approved, %d/%d approvals)\n", flags.nodeID, validApprovals, cfg.RequiredApprovals)
		} else {
			fmt.Printf("Added approval to %s (%d/%d approvals needed)\n", flags.nodeID, validApprovals, cfg.RequiredApprovals)
		}
	}

	return nil
}

type rejectFlags struct {
	targetDir string
	nodeID    string
	note      string
	quiet     bool
}

func newRejectCommand() *cobra.Command {
	flags := &rejectFlags{}

	cmd := &cobra.Command{
		Use:   "reject <node-id> [directory]",
		Short: "Reject a node back to draft",
		Long:  `Reject a node back to draft. Requires a note explaining the rejection reason.`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.nodeID = args[0]
			if len(args) > 1 {
				flags.targetDir = args[1]
			} else {
				flags.targetDir = "."
			}
			return runReject(flags)
		},
	}

	cmd.Flags().StringVar(&flags.note, "note", "", "Rejection reason (required)")
	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")

	return cmd
}

func runReject(flags *rejectFlags) error {
	// Validate note is provided
	if flags.note == "" {
		return fmt.Errorf("rejection note is required (use --note)")
	}

	// Load config
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	n, err := nodeRepo.Load(flags.nodeID)
	if err != nil {
		return fmt.Errorf("node %q not found: %w", flags.nodeID, err)
	}

	// Validate current status
	if n.Status != "review" {
		return fmt.Errorf("cannot reject node %q: status is %q, must be 'review'", flags.nodeID, n.Status)
	}

	// Update status
	oldStatus := n.Status
	n.Status = "draft"

	// Save the node
	if err := nodeRepo.Save(n); err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	// Log reject operation
	if err := logReviewOperation(flags.targetDir, n.ID, "reject", oldStatus, n.Status, flags.note); err != nil {
		fmt.Printf("Warning: failed to log reject operation: %v\n", err)
	}

	if !flags.quiet {
		fmt.Printf("Rejected %s (status: review -> draft)\nReason: %s\n", flags.nodeID, flags.note)
	}

	return nil
}
