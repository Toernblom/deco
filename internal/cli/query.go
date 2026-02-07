package cli

import (
	"fmt"
	"strings"

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
	blockType  string
	fields     []string // key=value pairs
	follow     string   // field name, or field:blocktype.field
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
  --kind:       Filter by node type (item, character, quest, etc.)
  --status:     Filter by status (draft, published, etc.)
  --tag:        Filter by tag (must have this tag)
  --block-type: Filter by custom block type within content
  --field:      Filter by block field value (key=value, repeatable)
  --follow:     Follow a field's refs to find related blocks

All filters and search are combined with AND logic.

Examples:
  deco query sword                              # Search for "sword" in title/summary
  deco query --kind item                        # List all items
  deco query sword --kind item                  # Search "sword" in items only
  deco query --status draft --tag combat
  deco query --block-type building              # List all building blocks
  deco query --block-type building --field age=bronze  # Bronze age buildings
  deco query --block-type building --field age=bronze --follow materials  # Follow refs
  deco query --block-type building --follow materials:recipe.output      # Explicit target`,
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
	cmd.Flags().StringVarP(&flags.blockType, "block-type", "b", "", "Filter by block type within content")
	cmd.Flags().StringArrayVarP(&flags.fields, "field", "f", nil, "Filter by block field (key=value, repeatable)")
	cmd.Flags().StringVar(&flags.follow, "follow", "", "Follow field refs to related blocks (field or field:blocktype.field)")

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
	if flags.blockType != "" {
		criteria.BlockType = &flags.blockType
	}
	if len(flags.fields) > 0 {
		criteria.FieldFilters = parseFieldFilters(flags.fields)
	}

	qe := query.New()

	// Validate --follow requires --block-type
	if flags.follow != "" && criteria.BlockType == nil {
		return fmt.Errorf("--follow requires --block-type")
	}

	// Block-level query mode
	if criteria.BlockType != nil {
		blockResults := qe.FilterBlocks(nodes, criteria)

		// Follow mode
		if flags.follow != "" {
			if len(blockResults) == 0 {
				fmt.Println("No blocks found to follow")
				return nil
			}
			followField, targets, err := parseFollowFlag(flags.follow)
			if err != nil {
				return err
			}
			followResults, err := qe.FollowBlocks(blockResults, followField, targets, nodes, cfg.CustomBlockTypes)
			if err != nil {
				return err
			}
			if len(followResults) == 0 {
				fmt.Println("No followed values found")
				return nil
			}
			printFollowResults(followResults)
			return nil
		}

		if len(blockResults) == 0 {
			fmt.Println("No blocks found")
			return nil
		}
		printBlocksTable(blockResults)
		return nil
	}

	// Node-level query mode
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

// parseFieldFilters parses key=value pairs into a map.
func parseFieldFilters(fields []string) map[string]string {
	result := make(map[string]string)
	for _, f := range fields {
		parts := splitFirst(f, '=')
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

// splitFirst splits s on the first occurrence of sep.
func splitFirst(s string, sep byte) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

// parseFollowFlag parses the --follow flag value.
// Supports: "fieldname" (auto from ref config) or "fieldname:blocktype.field" (explicit).
func parseFollowFlag(follow string) (string, []query.FollowTarget, error) {
	// Check for explicit target: field:blocktype.field
	colonIdx := strings.IndexByte(follow, ':')
	if colonIdx == -1 {
		// Auto mode: just the field name
		return follow, nil, nil
	}

	fieldName := follow[:colonIdx]
	targetSpec := follow[colonIdx+1:]

	// Parse blocktype.field
	dotIdx := strings.IndexByte(targetSpec, '.')
	if dotIdx == -1 {
		return "", nil, fmt.Errorf("invalid --follow target %q: expected blocktype.field", targetSpec)
	}

	target := query.FollowTarget{
		BlockType: targetSpec[:dotIdx],
		Field:     targetSpec[dotIdx+1:],
	}

	if target.BlockType == "" || target.Field == "" {
		return "", nil, fmt.Errorf("invalid --follow target %q: block type and field required", targetSpec)
	}

	return fieldName, []query.FollowTarget{target}, nil
}

// printFollowResults displays follow query results grouped by value.
func printFollowResults(results []query.FollowResult) {
	for i, r := range results {
		// Header: value (referenced by N block(s))
		fmt.Printf("%s (referenced by %d block(s))\n", r.Value, r.RefCount)

		if len(r.Matches) == 0 {
			fmt.Println("  (no matches found)")
		} else {
			for _, m := range r.Matches {
				fmt.Printf("  %s in %s > %s > block %d\n", m.Block.Type, m.NodeID, m.SectionName, m.BlockIndex)
				for k, v := range m.Block.Data {
					fmt.Printf("    %s: %v\n", k, v)
				}
			}
		}

		if i < len(results)-1 {
			fmt.Println()
		}
	}
}

// printBlocksTable displays block query results.
func printBlocksTable(blocks []query.BlockMatch) {
	fmt.Printf("Found %d block(s):\n\n", len(blocks))
	for _, b := range blocks {
		fmt.Printf("  [%s > %s] type: %s\n", b.NodeID, b.SectionName, b.Block.Type)
		for k, v := range b.Block.Data {
			fmt.Printf("    %s: %v\n", k, v)
		}
		fmt.Println()
	}
}
