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

	"github.com/Toernblom/deco/internal/cli/style"
	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/query"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type listFlags struct {
	kind      string
	status    string
	tag       string
	quiet     bool
	targetDir string
}

// NewListCommand creates the list subcommand
func NewListCommand() *cobra.Command {
	flags := &listFlags{}

	cmd := &cobra.Command{
		Use:   "list [directory]",
		Short: "List all nodes in the project",
		Long: `List all nodes in the project with optional filtering.

Filters can be combined to narrow down results:
  --kind:   Filter by node type (item, character, quest, etc.)
  --status: Filter by status (draft, review, approved, etc.)
  --tag:    Filter by tag (must have this tag)

Examples:
  deco list
  deco list --kind item
  deco list --status draft
  deco list --kind item --status approved
  deco list --tag combat`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				flags.targetDir = args[0]
			} else {
				flags.targetDir = "."
			}
			return runList(flags)
		},
	}

	cmd.Flags().StringVarP(&flags.kind, "kind", "k", "", "Filter by node kind")
	cmd.Flags().StringVarP(&flags.status, "status", "s", "", "Filter by status")
	cmd.Flags().StringVarP(&flags.tag, "tag", "t", "", "Filter by tag")
	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Output node IDs only, one per line")

	return cmd
}

func runList(flags *listFlags) error {
	// Load config to verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load all nodes
	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))
	nodes, err := nodeRepo.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}

	// Validate filter values
	if err := validateStatus(flags.status); err != nil {
		return err
	}
	if err := validateKind(flags.kind, nodes); err != nil {
		return err
	}

	// Build filter criteria
	criteria := query.FilterCriteria{}
	if flags.kind != "" {
		criteria.Kind = &flags.kind
	}
	if flags.status != "" {
		criteria.Status = &flags.status
	}
	if flags.tag != "" {
		criteria.Tags = []string{flags.tag}
	}

	// Apply filters
	qe := query.New()
	filtered := qe.Filter(nodes, criteria)

	// Display results
	quiet := flags.quiet || globalConfig.Quiet
	if len(filtered) == 0 {
		if !quiet {
			fmt.Println("No nodes found")
		}
		return nil
	}

	if quiet {
		for _, n := range filtered {
			fmt.Println(n.ID)
		}
		return nil
	}

	printNodesTable(filtered)
	return nil
}

func printNodesTable(nodes []domain.Node) {
	// Calculate column widths
	maxIDLen := 2     // "ID"
	maxKindLen := 4   // "KIND"
	maxStatusLen := 6 // "STATUS"
	maxTitleLen := 5  // "TITLE"

	for _, node := range nodes {
		if len(node.ID) > maxIDLen {
			maxIDLen = len(node.ID)
		}
		if len(node.Kind) > maxKindLen {
			maxKindLen = len(node.Kind)
		}
		if len(node.Status) > maxStatusLen {
			maxStatusLen = len(node.Status)
		}
		if len(node.Title) > maxTitleLen {
			maxTitleLen = len(node.Title)
		}
	}

	// Limit maximum title width to keep table readable
	if maxTitleLen > 50 {
		maxTitleLen = 50
	}

	// Print header with styling
	header := fmt.Sprintf("%-*s  %-*s  %-*s  %-*s",
		maxIDLen, "ID",
		maxKindLen, "KIND",
		maxStatusLen, "STATUS",
		maxTitleLen, "TITLE")
	fmt.Println(style.Header.Sprint(header))

	// Print separator
	separator := strings.Repeat("─", len(header))
	fmt.Println(style.Muted.Sprint(separator))

	// Print rows
	for _, node := range nodes {
		title := node.Title
		if len(title) > maxTitleLen {
			title = title[:maxTitleLen-3] + "..."
		}

		// Color the status based on its value
		statusStr := fmt.Sprintf("%-*s", maxStatusLen, node.Status)
		if c := style.StatusColor(node.Status); c != nil {
			statusStr = c.Sprint(statusStr)
		}

		fmt.Printf("%-*s  %-*s  %s  %-*s\n",
			maxIDLen, node.ID,
			maxKindLen, style.Muted.Sprint(node.Kind),
			statusStr,
			maxTitleLen, title)
	}

	// Print summary
	fmt.Printf("\n%s %d node(s)\n", style.Muted.Sprint("Total:"), len(nodes))
}
