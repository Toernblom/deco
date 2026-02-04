package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type graphFlags struct {
	format    string
	ascii     bool
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
  ascii    Simple ASCII layered view grouped by category

Examples:
  deco graph
  deco graph --format mermaid
  deco graph --ascii
  deco graph | dot -Tpng -o graph.png`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				flags.targetDir = args[0]
			} else {
				flags.targetDir = "."
			}
			if flags.ascii {
				flags.format = "ascii"
			}
			return runGraph(flags)
		},
	}

	cmd.Flags().StringVarP(&flags.format, "format", "f", "dot", "Output format (dot, mermaid, ascii)")
	cmd.Flags().BoolVar(&flags.ascii, "ascii", false, "Shorthand for --format ascii")

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
	case "ascii":
		outputASCII(nodes, edges)
	default:
		return fmt.Errorf("unknown format: %s (supported: dot, mermaid, ascii)", flags.format)
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

// outputASCII renders an ASCII graph grouped by category with branching connectors
func outputASCII(nodes []domain.Node, edges []edge) {
	if len(nodes) == 0 {
		return
	}

	// Group nodes by category
	categories := make(map[string][]string)
	for _, n := range nodes {
		cat := getCategory(n.ID)
		categories[cat] = append(categories[cat], n.ID)
	}
	for cat := range categories {
		sort.Strings(categories[cat])
	}

	// Build inter-category edges
	catEdges := make(map[string]map[string]bool)
	for _, e := range edges {
		fromCat := getCategory(e.from)
		toCat := getCategory(e.to)
		if fromCat != toCat {
			if catEdges[fromCat] == nil {
				catEdges[fromCat] = make(map[string]bool)
			}
			catEdges[fromCat][toCat] = true
		}
	}

	// Topological sort of categories
	layers := topoSortCategories(categories, catEdges)

	// Render
	fmt.Println()
	for i, layer := range layers {
		// Render category header
		renderCategoryHeader(layer)

		// Render nodes in each category
		renderCategoryNodes(layer, categories)

		// Render connectors
		if i < len(layers)-1 {
			renderCatConnectors(layer, layers[i+1], catEdges)
		}
	}
	fmt.Println()
}

// getCategory extracts the category (folder) from a node ID
func getCategory(id string) string {
	parts := strings.Split(id, "/")
	if len(parts) > 1 {
		return parts[0]
	}
	return id
}

// shortLabel extracts a short display name from node ID
func shortLabel(id string) string {
	parts := strings.Split(id, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return id
}

// topoSortCategories arranges categories into layers
func topoSortCategories(categories map[string][]string, catEdges map[string]map[string]bool) [][]string {
	remaining := make(map[string]bool)
	for cat := range categories {
		remaining[cat] = true
	}

	var layers [][]string
	for len(remaining) > 0 {
		var layer []string
		for cat := range remaining {
			hasIncoming := false
			for other := range remaining {
				if other != cat && catEdges[other] != nil && catEdges[other][cat] {
					hasIncoming = true
					break
				}
			}
			if !hasIncoming {
				layer = append(layer, cat)
			}
		}

		if len(layer) == 0 {
			for cat := range remaining {
				layer = append(layer, cat)
			}
		}

		sort.Strings(layer)
		for _, cat := range layer {
			delete(remaining, cat)
		}
		layers = append(layers, layer)
	}

	return layers
}

// renderCategoryHeader prints category names
func renderCategoryHeader(layer []string) {
	var names []string
	for _, cat := range layer {
		names = append(names, strings.Title(cat))
	}
	line := strings.Join(names, "     ")
	pad := (50 - len(line)) / 2
	if pad < 0 {
		pad = 0
	}
	fmt.Printf("%s%s\n", strings.Repeat(" ", pad), line)
}

// renderCategoryNodes prints nodes under each category
func renderCategoryNodes(layer []string, categories map[string][]string) {
	// Collect all node labels for this layer
	var allLabels []string
	for _, cat := range layer {
		for _, id := range categories[cat] {
			allLabels = append(allLabels, shortLabel(id))
		}
	}

	// If too many nodes, show count per category
	if len(allLabels) > 6 {
		var parts []string
		for _, cat := range layer {
			parts = append(parts, fmt.Sprintf("(%d)", len(categories[cat])))
		}
		line := strings.Join(parts, "     ")
		pad := (50 - len(line)) / 2
		if pad < 0 {
			pad = 0
		}
		fmt.Printf("%s%s\n", strings.Repeat(" ", pad), line)
	} else {
		// Show all node names
		line := strings.Join(allLabels, "   ")
		pad := (50 - len(line)) / 2
		if pad < 0 {
			pad = 0
		}
		fmt.Printf("%s%s\n", strings.Repeat(" ", pad), line)
	}
}

// renderCatConnectors draws connectors between category layers
func renderCatConnectors(fromLayer, toLayer []string, catEdges map[string]map[string]bool) {
	// Count connections
	hasConnection := false
	for _, from := range fromLayer {
		if catEdges[from] != nil {
			for _, to := range toLayer {
				if catEdges[from][to] {
					hasConnection = true
					break
				}
			}
		}
		if hasConnection {
			break
		}
	}

	if !hasConnection {
		fmt.Println()
		return
	}

	fromCount := len(fromLayer)
	toCount := len(toLayer)
	center := 25

	if fromCount == 1 && toCount == 1 {
		// Straight line
		fmt.Printf("%s|\n", strings.Repeat(" ", center))
		fmt.Printf("%sv\n", strings.Repeat(" ", center))
	} else if fromCount == 1 && toCount > 1 {
		// Fan out: one source to multiple targets
		fmt.Printf("%s|\n", strings.Repeat(" ", center))
		half := (toCount - 1) * 4
		fmt.Printf("%s/%s\\\n", strings.Repeat(" ", center-half), strings.Repeat("-", half*2-1))
		arrows := ""
		for i := 0; i < toCount; i++ {
			arrows += "v"
			if i < toCount-1 {
				arrows += strings.Repeat(" ", 7)
			}
		}
		pad := center - len(arrows)/2
		if pad < 0 {
			pad = 0
		}
		fmt.Printf("%s%s\n", strings.Repeat(" ", pad), arrows)
	} else if fromCount > 1 && toCount == 1 {
		// Fan in: multiple sources to one target
		bars := ""
		for i := 0; i < fromCount; i++ {
			bars += "|"
			if i < fromCount-1 {
				bars += strings.Repeat(" ", 7)
			}
		}
		pad := center - len(bars)/2
		if pad < 0 {
			pad = 0
		}
		fmt.Printf("%s%s\n", strings.Repeat(" ", pad), bars)
		half := (fromCount - 1) * 4
		fmt.Printf("%s\\%s/\n", strings.Repeat(" ", center-half), strings.Repeat("-", half*2-1))
		fmt.Printf("%sv\n", strings.Repeat(" ", center))
	} else {
		// Complex case: just show simple connector
		fmt.Printf("%s|\n", strings.Repeat(" ", center))
		fmt.Printf("%sv\n", strings.Repeat(" ", center))
	}
}
