package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type issuesFlags struct {
	severity  string
	nodeID    string
	targetDir string
}

// IssueResult holds an issue with its parent node context
type IssueResult struct {
	NodeID   string
	NodeKind string
	Issue    domain.Issue
}

// NewIssuesCommand creates the issues subcommand
func NewIssuesCommand() *cobra.Command {
	flags := &issuesFlags{}

	cmd := &cobra.Command{
		Use:   "issues",
		Short: "List all open issues/TBDs across the design",
		Long: `List all open issues and TBDs across the entire design graph.

Issues are tracked problems, questions, or TBDs that need resolution.
Use filters to narrow down the list.

Examples:
  deco issues
  deco issues --severity high
  deco issues --node systems/combat
  deco issues -s critical -n mechanics/health`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIssues(cmd.OutOrStdout(), flags)
		},
	}

	cmd.Flags().StringVarP(&flags.severity, "severity", "s", "", "Filter by severity (low, medium, high, critical)")
	cmd.Flags().StringVarP(&flags.nodeID, "node", "n", "", "Filter by node ID")
	cmd.Flags().StringVarP(&flags.targetDir, "dir", "d", ".", "Project directory")

	return cmd
}

func runIssues(w io.Writer, flags *issuesFlags) error {
	if w == nil {
		w = os.Stdout
	}

	// Verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	nodeRepo := node.NewYAMLRepository(flags.targetDir)

	// Load all nodes
	nodes, err := nodeRepo.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}

	// Collect all open issues
	var results []IssueResult
	for _, n := range nodes {
		// Filter by node if specified
		if flags.nodeID != "" && n.ID != flags.nodeID {
			continue
		}

		for _, issue := range n.Issues {
			// Skip resolved issues
			if issue.Resolved {
				continue
			}

			// Filter by severity if specified
			if flags.severity != "" && issue.Severity != flags.severity {
				continue
			}

			results = append(results, IssueResult{
				NodeID:   n.ID,
				NodeKind: n.Kind,
				Issue:    issue,
			})
		}
	}

	// Sort by severity (critical > high > medium > low), then by node ID
	severityOrder := map[string]int{
		"critical": 0,
		"high":     1,
		"medium":   2,
		"low":      3,
	}
	sort.Slice(results, func(i, j int) bool {
		si, sj := severityOrder[results[i].Issue.Severity], severityOrder[results[j].Issue.Severity]
		if si != sj {
			return si < sj
		}
		return results[i].NodeID < results[j].NodeID
	})

	// Output results
	if len(results) == 0 {
		fmt.Fprintln(w, "No open issues found.")
		return nil
	}

	fmt.Fprintf(w, "Found %d open issue(s):\n\n", len(results))
	for _, r := range results {
		fmt.Fprintf(w, "[%s] %s\n", r.Issue.Severity, r.Issue.ID)
		fmt.Fprintf(w, "  Node: %s\n", r.NodeID)
		fmt.Fprintf(w, "  Location: %s\n", r.Issue.Location)
		fmt.Fprintf(w, "  %s\n\n", r.Issue.Description)
	}

	return nil
}
