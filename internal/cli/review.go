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
	cmd.AddCommand(newStatusCommand())

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
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))
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
	historyPath := config.ResolveHistoryPath(cfg, flags.targetDir)
	if err := logReviewOperation(historyPath, n.ID, "submit", oldStatus, n.Status, ""); err != nil {
		fmt.Printf("Warning: failed to log submit operation: %v\n", err)
	}

	if !flags.quiet {
		fmt.Printf("Submitted %s for review (status: draft -> review)\n", flags.nodeID)
	}

	return nil
}

func logReviewOperation(historyPath, nodeID, operation, oldStatus, newStatus, note string) error {
	historyRepo := history.NewYAMLRepository(historyPath)

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
	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))
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
	historyPath := config.ResolveHistoryPath(cfg, flags.targetDir)
	if err := logReviewOperation(historyPath, n.ID, "approve", oldStatus, n.Status, flags.note); err != nil {
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
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))
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
	historyPath := config.ResolveHistoryPath(cfg, flags.targetDir)
	if err := logReviewOperation(historyPath, n.ID, "reject", oldStatus, n.Status, flags.note); err != nil {
		fmt.Printf("Warning: failed to log reject operation: %v\n", err)
	}

	if !flags.quiet {
		fmt.Printf("Rejected %s (status: review -> draft)\nReason: %s\n", flags.nodeID, flags.note)
	}

	return nil
}

type statusFlags struct {
	targetDir string
	nodeID    string
}

func newStatusCommand() *cobra.Command {
	flags := &statusFlags{}

	cmd := &cobra.Command{
		Use:   "status [node-id] [directory]",
		Short: "Show review status of a node or list all nodes in review",
		Long: `Show review status of a node including current status, approvals, and requirements.

When called without a node ID, lists all nodes currently in review status.`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 1 {
				flags.nodeID = args[0]
			}
			if len(args) >= 2 {
				flags.targetDir = args[1]
			} else {
				flags.targetDir = "."
			}
			if flags.nodeID == "" {
				return runStatusAll(cmd, flags)
			}
			return runStatus(cmd, flags)
		},
	}

	return cmd
}

func runStatusAll(cmd *cobra.Command, flags *statusFlags) error {
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))
	nodes, err := nodeRepo.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}

	out := cmd.OutOrStdout()
	var reviewNodes []domain.Node
	for _, n := range nodes {
		if n.Status == "review" {
			reviewNodes = append(reviewNodes, n)
		}
	}

	if len(reviewNodes) == 0 {
		fmt.Fprintln(out, "No nodes in review.")
		return nil
	}

	fmt.Fprintln(out, "Nodes in review:")
	for _, n := range reviewNodes {
		// Count current-version approvals
		approvals := 0
		submitter := ""
		for _, r := range n.Reviewers {
			if r.Version == n.Version {
				approvals++
			}
		}
		// Check history for submitter (best effort: use last reviewer or "unknown")
		if submitter == "" {
			submitter = "unknown"
		}
		fmt.Fprintf(out, "  %-25s v%d  %d approval(s)\n", n.ID, n.Version, approvals)
	}

	return nil
}

func runStatus(cmd *cobra.Command, flags *statusFlags) error {
	// Load config
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))
	n, err := nodeRepo.Load(flags.nodeID)
	if err != nil {
		return fmt.Errorf("node %q not found: %w", flags.nodeID, err)
	}

	// Count approvals for current version
	validApprovals := 0
	for _, r := range n.Reviewers {
		if r.Version == n.Version {
			validApprovals++
		}
	}

	// Print status
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Node: %s\n", n.ID)
	fmt.Fprintf(out, "Version: %d\n", n.Version)
	fmt.Fprintf(out, "Status: %s\n", n.Status)
	fmt.Fprintf(out, "Approvals: %d/%d\n", validApprovals, cfg.RequiredApprovals)

	if len(n.Reviewers) > 0 {
		fmt.Fprintf(out, "\nReviewers:\n")
		for _, r := range n.Reviewers {
			versionNote := ""
			if r.Version != n.Version {
				versionNote = " (stale - v" + fmt.Sprint(r.Version) + ")"
			}
			fmt.Fprintf(out, "  - %s at %s%s\n", r.Name, r.Timestamp.Format("2006-01-02 15:04"), versionNote)
			if r.Note != "" {
				fmt.Fprintf(out, "    Note: %s\n", r.Note)
			}
		}
	}

	return nil
}
