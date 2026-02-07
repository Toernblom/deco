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
	"strings"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type newFlags struct {
	kind      string
	title     string
	tags      string
	summary   string
	force     bool
	quiet     bool
	targetDir string
}

// NewNewCommand creates the new subcommand
func NewNewCommand() *cobra.Command {
	flags := &newFlags{}

	cmd := &cobra.Command{
		Use:   "new <node-id> [directory]",
		Short: "Scaffold a new node",
		Long: `Create a new node YAML file with required fields populated.

The node ID is derived from the positional argument and supports
nested paths (e.g., systems/combat, mechanics/stealth).

Examples:
  deco new systems/combat --kind system --title "Combat System"
  deco new mechanics/stealth --kind mechanic --title "Stealth" --tags core,pvp --summary "..."`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				flags.targetDir = args[1]
			} else {
				flags.targetDir = "."
			}
			return runNew(args[0], flags)
		},
	}

	cmd.Flags().StringVar(&flags.kind, "kind", "", "Node kind (required)")
	cmd.Flags().StringVar(&flags.title, "title", "", "Node title (required)")
	cmd.Flags().StringVar(&flags.tags, "tags", "", "Comma-separated tags")
	cmd.Flags().StringVar(&flags.summary, "summary", "", "Node summary")
	cmd.Flags().BoolVar(&flags.force, "force", false, "Overwrite existing node")
	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")
	cmd.MarkFlagRequired("kind")
	cmd.MarkFlagRequired("title")

	return cmd
}

func runNew(nodeID string, flags *newFlags) error {
	// Load config to verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Check if node already exists
	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))
	exists, err := nodeRepo.Exists(nodeID)
	if err != nil {
		return fmt.Errorf("failed to check node existence: %w", err)
	}
	if exists && !flags.force {
		return fmt.Errorf("node %q already exists (use --force to overwrite)", nodeID)
	}

	// Parse tags
	var tags []string
	if flags.tags != "" {
		for _, tag := range strings.Split(flags.tags, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}

	// Build the node
	n := domain.Node{
		ID:      nodeID,
		Kind:    flags.kind,
		Version: 1,
		Status:  "draft",
		Title:   flags.title,
		Tags:    tags,
		Summary: flags.summary,
		Content: &domain.Content{
			Sections: []domain.Section{},
		},
	}

	// Validate the node
	if err := n.Validate(); err != nil {
		return fmt.Errorf("invalid node: %w", err)
	}

	// Save the node
	if err := nodeRepo.Save(n); err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	// Log creation to history
	historyPath := config.ResolveHistoryPath(cfg, flags.targetDir)
	historyRepo := history.NewYAMLRepository(historyPath)
	entry := domain.AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      nodeID,
		Operation:   "create",
		User:        GetCurrentUser(),
		ContentHash: ComputeContentHash(n),
		After: map[string]interface{}{
			"kind":    n.Kind,
			"title":   n.Title,
			"version": n.Version,
			"status":  n.Status,
		},
	}
	if err := historyRepo.Append(entry); err != nil {
		if !flags.quiet {
			fmt.Printf("Warning: failed to log creation: %v\n", err)
		}
	}

	if !flags.quiet {
		fmt.Printf("Created node: %s (%s)\n", nodeID, flags.kind)
	}

	return nil
}
