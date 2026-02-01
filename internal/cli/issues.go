package cli

import (
	"encoding/json"
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
	severity   string
	nodeID     string
	kind       string
	tag        string
	showAll    bool
	jsonOutput bool
	quiet      bool
	summary    bool
	targetDir  string
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

Severity levels (low to critical):
  low      - Minor improvements, nice-to-haves
  medium   - Should be addressed before approval
  high     - Significant gaps affecting design quality
  critical - Blockers that must be resolved immediately

Examples:
  deco issues                          # List all open issues
  deco issues --severity high          # Filter by severity
  deco issues --node systems/combat    # Filter by node
  deco issues --kind mechanic          # Filter by node kind
  deco issues --tag combat             # Filter by node tag
  deco issues --all                    # Include resolved issues
  deco issues --summary                # Show per-node rollup
  deco issues --json                   # Output as JSON
  deco issues -q                       # Quiet mode (counts only)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIssues(cmd.OutOrStdout(), flags)
		},
	}

	cmd.Flags().StringVarP(&flags.severity, "severity", "s", "", "Filter by severity (low, medium, high, critical)")
	cmd.Flags().StringVarP(&flags.nodeID, "node", "n", "", "Filter by node ID")
	cmd.Flags().StringVarP(&flags.kind, "kind", "k", "", "Filter by node kind")
	cmd.Flags().StringVarP(&flags.tag, "tag", "t", "", "Filter by node tag")
	cmd.Flags().BoolVarP(&flags.showAll, "all", "a", false, "Show all issues including resolved")
	cmd.Flags().BoolVarP(&flags.jsonOutput, "json", "j", false, "Output as JSON")
	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Quiet mode (show counts only)")
	cmd.Flags().BoolVar(&flags.summary, "summary", false, "Show per-node summary rollup")
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

	// Collect issues with filtering
	var results []IssueResult
	for _, n := range nodes {
		// Filter by node ID if specified
		if flags.nodeID != "" && n.ID != flags.nodeID {
			continue
		}

		// Filter by node kind if specified
		if flags.kind != "" && n.Kind != flags.kind {
			continue
		}

		// Filter by node tag if specified
		if flags.tag != "" && !hasTag(n.Tags, flags.tag) {
			continue
		}

		for _, issue := range n.Issues {
			// Skip resolved issues unless --all is specified
			if issue.Resolved && !flags.showAll {
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

	// Handle different output formats
	if flags.jsonOutput {
		return outputIssuesJSON(w, results, flags.summary)
	}

	if flags.quiet {
		return outputIssuesQuiet(w, results)
	}

	if flags.summary {
		return outputIssuesSummary(w, results)
	}

	return outputIssuesHuman(w, results)
}

// hasTag checks if a tag exists in the list
func hasTag(tags []string, target string) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}
	return false
}

// outputIssuesHuman outputs issues in human-readable format
func outputIssuesHuman(w io.Writer, results []IssueResult) error {
	if len(results) == 0 {
		fmt.Fprintln(w, "No open issues found.")
		return nil
	}

	fmt.Fprintf(w, "Found %d issue(s):\n\n", len(results))
	for _, r := range results {
		status := ""
		if r.Issue.Resolved {
			status = " [RESOLVED]"
		}
		fmt.Fprintf(w, "[%s] %s%s\n", r.Issue.Severity, r.Issue.ID, status)
		fmt.Fprintf(w, "  Node: %s\n", r.NodeID)
		fmt.Fprintf(w, "  Location: %s\n", r.Issue.Location)
		fmt.Fprintf(w, "  %s\n\n", r.Issue.Description)
	}

	return nil
}

// outputIssuesQuiet outputs just the count of issues
func outputIssuesQuiet(w io.Writer, results []IssueResult) error {
	// Count by severity
	counts := map[string]int{
		"critical": 0,
		"high":     0,
		"medium":   0,
		"low":      0,
	}
	resolved := 0

	for _, r := range results {
		if r.Issue.Resolved {
			resolved++
		} else {
			counts[r.Issue.Severity]++
		}
	}

	total := counts["critical"] + counts["high"] + counts["medium"] + counts["low"]
	if total == 0 && resolved == 0 {
		fmt.Fprintln(w, "0")
		return nil
	}

	fmt.Fprintf(w, "%d", total)
	if resolved > 0 {
		fmt.Fprintf(w, " (+%d resolved)", resolved)
	}
	fmt.Fprintln(w)
	return nil
}

// outputIssuesSummary outputs a per-node rollup
func outputIssuesSummary(w io.Writer, results []IssueResult) error {
	if len(results) == 0 {
		fmt.Fprintln(w, "No issues found.")
		return nil
	}

	// Group by node
	byNode := make(map[string][]IssueResult)
	for _, r := range results {
		byNode[r.NodeID] = append(byNode[r.NodeID], r)
	}

	// Get sorted node IDs
	nodeIDs := make([]string, 0, len(byNode))
	for id := range byNode {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)

	fmt.Fprintf(w, "Issues by node (%d total across %d nodes):\n\n", len(results), len(nodeIDs))

	for _, nodeID := range nodeIDs {
		issues := byNode[nodeID]
		counts := countBySeverity(issues)
		kind := issues[0].NodeKind

		fmt.Fprintf(w, "%s (%s): ", nodeID, kind)
		parts := []string{}
		if counts["critical"] > 0 {
			parts = append(parts, fmt.Sprintf("%d critical", counts["critical"]))
		}
		if counts["high"] > 0 {
			parts = append(parts, fmt.Sprintf("%d high", counts["high"]))
		}
		if counts["medium"] > 0 {
			parts = append(parts, fmt.Sprintf("%d medium", counts["medium"]))
		}
		if counts["low"] > 0 {
			parts = append(parts, fmt.Sprintf("%d low", counts["low"]))
		}
		if counts["resolved"] > 0 {
			parts = append(parts, fmt.Sprintf("%d resolved", counts["resolved"]))
		}
		for i, p := range parts {
			if i > 0 {
				fmt.Fprint(w, ", ")
			}
			fmt.Fprint(w, p)
		}
		fmt.Fprintln(w)
	}

	return nil
}

// countBySeverity counts issues by severity level
func countBySeverity(issues []IssueResult) map[string]int {
	counts := map[string]int{
		"critical": 0,
		"high":     0,
		"medium":   0,
		"low":      0,
		"resolved": 0,
	}
	for _, r := range issues {
		if r.Issue.Resolved {
			counts["resolved"]++
		} else {
			counts[r.Issue.Severity]++
		}
	}
	return counts
}

// IssuesJSONOutput is the structure for JSON output
type IssuesJSONOutput struct {
	Total   int                    `json:"total"`
	Open    int                    `json:"open"`
	Counts  map[string]int         `json:"counts"`
	Issues  []IssueResult          `json:"issues,omitempty"`
	ByNode  map[string]NodeSummary `json:"by_node,omitempty"`
}

// NodeSummary is a per-node issue summary
type NodeSummary struct {
	Kind   string         `json:"kind"`
	Counts map[string]int `json:"counts"`
}

// outputIssuesJSON outputs issues as JSON
func outputIssuesJSON(w io.Writer, results []IssueResult, summary bool) error {
	counts := map[string]int{
		"critical": 0,
		"high":     0,
		"medium":   0,
		"low":      0,
		"resolved": 0,
	}

	for _, r := range results {
		if r.Issue.Resolved {
			counts["resolved"]++
		} else {
			counts[r.Issue.Severity]++
		}
	}

	open := counts["critical"] + counts["high"] + counts["medium"] + counts["low"]

	output := IssuesJSONOutput{
		Total:  len(results),
		Open:   open,
		Counts: counts,
	}

	if summary {
		// Build per-node summary
		byNode := make(map[string]NodeSummary)
		for _, r := range results {
			ns, exists := byNode[r.NodeID]
			if !exists {
				ns = NodeSummary{
					Kind:   r.NodeKind,
					Counts: map[string]int{},
				}
			}
			if r.Issue.Resolved {
				ns.Counts["resolved"]++
			} else {
				ns.Counts[r.Issue.Severity]++
			}
			byNode[r.NodeID] = ns
		}
		output.ByNode = byNode
	} else {
		output.Issues = results
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
