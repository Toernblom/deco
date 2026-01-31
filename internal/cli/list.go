package cli

import (
	"fmt"
	"strings"

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
  --status: Filter by status (draft, published, etc.)
  --tag:    Filter by tag (must have this tag)

Examples:
  deco list
  deco list --kind item
  deco list --status draft
  deco list --kind item --status published
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

	return cmd
}

func runList(flags *listFlags) error {
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
	if len(filtered) == 0 {
		fmt.Println("No nodes found")
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

	// Print header
	header := fmt.Sprintf("%-*s  %-*s  %-*s  %-*s",
		maxIDLen, "ID",
		maxKindLen, "KIND",
		maxStatusLen, "STATUS",
		maxTitleLen, "TITLE")
	fmt.Println(header)

	// Print separator
	separator := strings.Repeat("-", len(header))
	fmt.Println(separator)

	// Print rows
	for _, node := range nodes {
		title := node.Title
		if len(title) > maxTitleLen {
			title = title[:maxTitleLen-3] + "..."
		}

		fmt.Printf("%-*s  %-*s  %-*s  %-*s\n",
			maxIDLen, node.ID,
			maxKindLen, node.Kind,
			maxStatusLen, node.Status,
			maxTitleLen, title)
	}

	// Print summary
	fmt.Printf("\nTotal: %d node(s)\n", len(nodes))
}
