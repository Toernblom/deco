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
	"sort"
	"strings"

	"github.com/Toernblom/deco/internal/cli/style"
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
		if !globalConfig.Quiet {
			fmt.Println("No nodes found in project")
		}
		return nil
	}

	// Gather statistics
	stats := gatherStats(nodes, cfg)

	// Print statistics
	if globalConfig.Quiet {
		printStatsQuiet(stats)
	} else {
		printStats(stats)
	}

	return nil
}

type projectStats struct {
	totalNodes            int
	nodesByKind           map[string]int
	nodesByStatus         map[string]int
	openIssuesBySev       map[string]int
	totalOpenIssues       int
	danglingRefs          int
	totalValidationErrors int
	validationByCategory  map[string]int
}

func gatherStats(nodes []domain.Node, cfg config.Config) projectStats {
	stats := projectStats{
		totalNodes:           len(nodes),
		nodesByKind:          make(map[string]int),
		nodesByStatus:        make(map[string]int),
		openIssuesBySev:      make(map[string]int),
		validationByCategory: make(map[string]int),
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

	// Run full validation to count all errors
	orchestrator := validator.NewOrchestratorWithFullConfig(
		cfg.RequiredApprovals,
		cfg.CustomBlockTypes,
		cfg.SchemaRules,
	)
	collector := orchestrator.ValidateAll(nodes)
	registry := domain.NewErrorCodeRegistry()
	for _, err := range collector.Errors() {
		stats.totalValidationErrors++
		if ec, ok := registry.Lookup(err.Code); ok {
			stats.validationByCategory[ec.Category]++
		} else {
			stats.validationByCategory["other"]++
		}
	}

	return stats
}

func printStats(stats projectStats) {
	fmt.Println(style.Header.Sprint("PROJECT STATISTICS"))
	fmt.Println(style.Muted.Sprint(strings.Repeat("═", 50)))

	// Total nodes
	fmt.Printf("\n%s %d\n", style.Muted.Sprint("Total nodes:"), stats.totalNodes)

	// Nodes by kind
	fmt.Printf("\n%s\n", style.Header.Sprint("NODES BY KIND"))
	fmt.Println(style.Muted.Sprint(strings.Repeat("─", 30)))
	printSortedMap(stats.nodesByKind)

	// Nodes by status
	fmt.Printf("\n%s\n", style.Header.Sprint("NODES BY STATUS"))
	fmt.Println(style.Muted.Sprint(strings.Repeat("─", 30)))
	printSortedMapWithStatus(stats.nodesByStatus)

	// Open issues by severity
	fmt.Printf("\n%s\n", style.Header.Sprint("OPEN ISSUES BY SEVERITY"))
	fmt.Println(style.Muted.Sprint(strings.Repeat("─", 30)))
	if stats.totalOpenIssues == 0 {
		fmt.Printf("  %s\n", style.Success.Sprint("No open issues"))
	} else {
		// Print in severity order
		severityOrder := []string{"critical", "high", "medium", "low"}
		for _, sev := range severityOrder {
			if count, ok := stats.openIssuesBySev[sev]; ok && count > 0 {
				sevColor := style.SeverityColor(sev)
				fmt.Printf("  %s %d\n", sevColor.Sprintf("%-12s", sev), count)
			}
		}
		fmt.Printf("  %-12s %d\n", style.Muted.Sprint("Total"), stats.totalOpenIssues)
	}

	// Reference health
	fmt.Printf("\n%s\n", style.Header.Sprint("REFERENCE HEALTH"))
	fmt.Println(style.Muted.Sprint(strings.Repeat("─", 30)))
	if stats.danglingRefs == 0 {
		fmt.Printf("  %s\n", style.Success.Sprint("All references valid"))
	} else {
		fmt.Printf("  %s %d\n", style.Warning.Sprint("Dangling references:"), stats.danglingRefs)
	}

	// Validation health
	fmt.Printf("\n%s\n", style.Header.Sprint("VALIDATION HEALTH"))
	fmt.Println(style.Muted.Sprint(strings.Repeat("─", 30)))
	if stats.totalValidationErrors == 0 {
		fmt.Printf("  %s\n", style.Success.Sprint("No errors"))
	} else {
		fmt.Printf("  %s %d\n", style.Error.Sprint("Errors:"), stats.totalValidationErrors)
		categoryOrder := []string{"schema", "refs", "validation"}
		for _, cat := range categoryOrder {
			if count, ok := stats.validationByCategory[cat]; ok && count > 0 {
				fmt.Printf("    %-12s %d\n", cat+":", count)
			}
		}
		// Print any extra categories not in the standard order
		for cat, count := range stats.validationByCategory {
			if cat != "schema" && cat != "refs" && cat != "validation" && count > 0 {
				fmt.Printf("    %-12s %d\n", cat+":", count)
			}
		}
	}
}

func printStatsQuiet(stats projectStats) {
	fmt.Printf("nodes=%d issues=%d dangling_refs=%d errors=%d\n",
		stats.totalNodes, stats.totalOpenIssues, stats.danglingRefs, stats.totalValidationErrors)
}

func printSortedMap(m map[string]int) {
	if len(m) == 0 {
		fmt.Printf("  %s\n", style.Muted.Sprint("(none)"))
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

func printSortedMapWithStatus(m map[string]int) {
	if len(m) == 0 {
		fmt.Printf("  %s\n", style.Muted.Sprint("(none)"))
		return
	}

	// Sort keys alphabetically
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		label := fmt.Sprintf("%-12s", k)
		if c := style.StatusColor(k); c != nil {
			label = c.Sprint(label)
		}
		fmt.Printf("  %s %d\n", label, m[k])
	}
}
