package cli

import (
	"fmt"

	"github.com/Toernblom/deco/internal/services/query"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type queryFlags struct {
	kind       string
	status     string
	tag        string
	targetDir  string
	searchTerm string
}

// NewQueryCommand creates the query subcommand
func NewQueryCommand() *cobra.Command {
	flags := &queryFlags{}

	cmd := &cobra.Command{
		Use:   "query [search-term] [directory]",
		Short: "Search and filter nodes",
		Long: `Search and filter nodes with optional text search and filtering.

The search term is optional and searches title and summary (case-insensitive).
Filters can be combined with search to narrow down results:
  --kind:   Filter by node type (item, character, quest, etc.)
  --status: Filter by status (draft, published, etc.)
  --tag:    Filter by tag (must have this tag)

All filters and search are combined with AND logic.

Examples:
  deco query sword                    # Search for "sword" in title/summary
  deco query --kind item              # List all items
  deco query sword --kind item        # Search "sword" in items only
  deco query --status draft --tag combat
  deco query "health potion" --kind item --status published`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse arguments: [search-term] [directory]
			switch len(args) {
			case 0:
				flags.targetDir = "."
			case 1:
				// Could be search term or directory
				// If it looks like a path (contains / or \, or is "."), treat as directory
				if isDirectory(args[0]) {
					flags.targetDir = args[0]
				} else {
					flags.searchTerm = args[0]
					flags.targetDir = "."
				}
			case 2:
				flags.searchTerm = args[0]
				flags.targetDir = args[1]
			}
			return runQuery(flags)
		},
	}

	cmd.Flags().StringVarP(&flags.kind, "kind", "k", "", "Filter by node kind")
	cmd.Flags().StringVarP(&flags.status, "status", "s", "", "Filter by status")
	cmd.Flags().StringVarP(&flags.tag, "tag", "t", "", "Filter by tag")

	return cmd
}

// isDirectory checks if a string looks like a directory path
func isDirectory(s string) bool {
	// Check for common directory indicators
	if s == "." || s == ".." {
		return true
	}
	// Check for path separators
	for _, c := range s {
		if c == '/' || c == '\\' {
			return true
		}
	}
	// Check if it starts with a drive letter (Windows)
	if len(s) >= 2 && s[1] == ':' {
		return true
	}
	return false
}

func runQuery(flags *queryFlags) error {
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

	// Apply filters and search
	qe := query.New()
	results := qe.Filter(nodes, criteria)
	if flags.searchTerm != "" {
		results = qe.Search(results, flags.searchTerm)
	}

	// Display results
	if len(results) == 0 {
		fmt.Println("No nodes found")
		return nil
	}

	printNodesTable(results)
	return nil
}
