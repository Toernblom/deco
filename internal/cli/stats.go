package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/validator"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type statsFlags struct {
	targetDir string
}

// NewStatsCommand creates the stats subcommand
func NewStatsCommand() *cobra.Command {
	flags := &statsFlags{}

	cmd := &cobra.Command{
		Use:   "stats [directory]",
		Short: "Show project overview and health statistics",
		Long: `Display a summary of project health including:
  - Node count by kind
  - Node count by status
  - Open issues by severity
  - Reference health (dangling refs)
  - Constraint violations

Examples:
  deco stats
  deco stats /path/to/project`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				flags.targetDir = args[0]
			} else {
				flags.targetDir = "."
			}
			return runStats(flags)
		},
	}

	return cmd
}

func runStats(flags *statsFlags) error {
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

	if len(nodes) == 0 {
		fmt.Println("No nodes found in project")
		return nil
	}

	// Gather statistics
	stats := gatherStats(nodes, flags.targetDir)

	// Print statistics
	printStats(stats)

	return nil
}

type projectStats struct {
	totalNodes          int
	nodesByKind         map[string]int
	nodesByStatus       map[string]int
	openIssuesBySev     map[string]int
	totalOpenIssues     int
	danglingRefs        int
	constraintViolations int
}

func gatherStats(nodes []domain.Node, targetDir string) projectStats {
	stats := projectStats{
		totalNodes:      len(nodes),
		nodesByKind:    make(map[string]int),
		nodesByStatus:  make(map[string]int),
		openIssuesBySev: make(map[string]int),
	}

	// Build set of existing node IDs for reference checking
	nodeIDs := make(map[string]bool)
	for _, n := range nodes {
		nodeIDs[n.ID] = true
	}

	// Gather node-level stats
	for _, n := range nodes {
		// Count by kind
		stats.nodesByKind[n.Kind]++

		// Count by status
		stats.nodesByStatus[n.Status]++

		// Count open issues by severity
		for _, issue := range n.Issues {
			if !issue.Resolved {
				stats.openIssuesBySev[issue.Severity]++
				stats.totalOpenIssues++
			}
		}

		// Count dangling references
		for _, ref := range n.Refs.Uses {
			if !nodeIDs[ref.Target] {
				stats.danglingRefs++
			}
		}
		for _, ref := range n.Refs.Related {
			if !nodeIDs[ref.Target] {
				stats.danglingRefs++
			}
		}
	}

	// Run validator to count constraint violations
	orchestrator := validator.NewOrchestrator()
	collector := orchestrator.ValidateAll(nodes)
	for _, err := range collector.Errors() {
		if err.Code == "E041" {
			stats.constraintViolations++
		}
	}

	return stats
}

func printStats(stats projectStats) {
	fmt.Println("PROJECT STATISTICS")
	fmt.Println(strings.Repeat("=", 50))

	// Total nodes
	fmt.Printf("\nTotal nodes: %d\n", stats.totalNodes)

	// Nodes by kind
	fmt.Println("\nNODES BY KIND")
	fmt.Println(strings.Repeat("-", 30))
	printSortedMap(stats.nodesByKind)

	// Nodes by status
	fmt.Println("\nNODES BY STATUS")
	fmt.Println(strings.Repeat("-", 30))
	printSortedMap(stats.nodesByStatus)

	// Open issues by severity
	fmt.Println("\nOPEN ISSUES BY SEVERITY")
	fmt.Println(strings.Repeat("-", 30))
	if stats.totalOpenIssues == 0 {
		fmt.Println("  No open issues")
	} else {
		// Print in severity order
		severityOrder := []string{"critical", "high", "medium", "low"}
		for _, sev := range severityOrder {
			if count, ok := stats.openIssuesBySev[sev]; ok && count > 0 {
				fmt.Printf("  %-12s %d\n", sev, count)
			}
		}
		fmt.Printf("  %-12s %d\n", "Total", stats.totalOpenIssues)
	}

	// Reference health
	fmt.Println("\nREFERENCE HEALTH")
	fmt.Println(strings.Repeat("-", 30))
	if stats.danglingRefs == 0 {
		fmt.Println("  All references valid")
	} else {
		fmt.Printf("  Dangling references: %d\n", stats.danglingRefs)
	}

	// Constraint violations
	fmt.Println("\nCONSTRAINT VIOLATIONS")
	fmt.Println(strings.Repeat("-", 30))
	if stats.constraintViolations == 0 {
		fmt.Println("  No violations")
	} else {
		fmt.Printf("  Violations: %d\n", stats.constraintViolations)
	}
}

func printSortedMap(m map[string]int) {
	if len(m) == 0 {
		fmt.Println("  (none)")
		return
	}

	// Sort keys alphabetically
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("  %-12s %d\n", k, m[k])
	}
}
