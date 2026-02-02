package cli

import (
	"fmt"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type graphFlags struct {
	format    string
	targetDir string
}

// NewGraphCommand creates the graph subcommand
func NewGraphCommand() *cobra.Command {
	flags := &graphFlags{}

	cmd := &cobra.Command{
		Use:   "graph [directory]",
		Short: "Output dependency graph",
		Long: `Output the node dependency graph in DOT or Mermaid format.

Edges are created from refs.uses and refs.related fields.

Formats:
  dot      Graphviz DOT format (default)
  mermaid  Mermaid flowchart for Markdown embedding

Examples:
  deco graph
  deco graph --format mermaid
  deco graph | dot -Tpng -o graph.png`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				flags.targetDir = args[0]
			} else {
				flags.targetDir = "."
			}
			return runGraph(flags)
		},
	}

	cmd.Flags().StringVarP(&flags.format, "format", "f", "dot", "Output format (dot, mermaid)")

	return cmd
}

func runGraph(flags *graphFlags) error {
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
		fmt.Println("No nodes found")
		return nil
	}

	// Build edges
	edges := buildEdges(nodes)

	// Output in requested format
	switch flags.format {
	case "dot":
		outputDOT(nodes, edges)
	case "mermaid":
		outputMermaid(nodes, edges)
	default:
		return fmt.Errorf("unknown format: %s (supported: dot, mermaid)", flags.format)
	}

	return nil
}

type edge struct {
	from    string
	to      string
	refType string // "uses" or "related"
}

func buildEdges(nodes []domain.Node) []edge {
	var edges []edge

	for _, n := range nodes {
		// Add edges from uses refs
		for _, ref := range n.Refs.Uses {
			edges = append(edges, edge{
				from:    n.ID,
				to:      ref.Target,
				refType: "uses",
			})
		}

		// Add edges from related refs
		for _, ref := range n.Refs.Related {
			edges = append(edges, edge{
				from:    n.ID,
				to:      ref.Target,
				refType: "related",
			})
		}
	}

	return edges
}

func outputDOT(nodes []domain.Node, edges []edge) {
	fmt.Println("digraph deco {")
	fmt.Println("  rankdir=LR;")
	fmt.Println("  node [shape=box];")
	fmt.Println()

	// Declare nodes with labels
	for _, n := range nodes {
		label := strings.ReplaceAll(n.Title, "\"", "\\\"")
		nodeID := dotID(n.ID)
		fmt.Printf("  %s [label=\"%s\\n(%s)\"];\n", nodeID, label, n.ID)
	}

	fmt.Println()

	// Declare edges
	for _, e := range edges {
		fromID := dotID(e.from)
		toID := dotID(e.to)
		style := ""
		if e.refType == "related" {
			style = " [style=dashed]"
		}
		fmt.Printf("  %s -> %s%s;\n", fromID, toID, style)
	}

	fmt.Println("}")
}

func outputMermaid(nodes []domain.Node, edges []edge) {
	fmt.Println("```mermaid")
	fmt.Println("flowchart LR")

	// Declare nodes with labels
	for _, n := range nodes {
		label := strings.ReplaceAll(n.Title, "\"", "'")
		nodeID := mermaidID(n.ID)
		fmt.Printf("  %s[\"%s\"]\n", nodeID, label)
	}

	fmt.Println()

	// Declare edges
	for _, e := range edges {
		fromID := mermaidID(e.from)
		toID := mermaidID(e.to)
		if e.refType == "related" {
			fmt.Printf("  %s -.-> %s\n", fromID, toID)
		} else {
			fmt.Printf("  %s --> %s\n", fromID, toID)
		}
	}

	fmt.Println("```")
}

// dotID converts a node ID to a valid DOT identifier
func dotID(id string) string {
	// Replace slashes and other special chars with underscores
	return "\"" + id + "\""
}

// mermaidID converts a node ID to a valid Mermaid identifier
func mermaidID(id string) string {
	// Replace slashes with underscores for valid identifiers
	return strings.ReplaceAll(id, "/", "_")
}
