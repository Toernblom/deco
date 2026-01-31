package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/graph"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type showFlags struct {
	jsonOutput bool
	targetDir  string
}

// NewShowCommand creates the show subcommand
func NewShowCommand() *cobra.Command {
	flags := &showFlags{}

	cmd := &cobra.Command{
		Use:   "show <node-id> [directory]",
		Short: "Show detailed information about a node",
		Long: `Show detailed information about a specific node including:
  - All node fields (ID, kind, version, status, title, etc.)
  - Summary and description
  - Tags
  - References (uses/related)
  - Reverse references (what nodes reference this one)

Output can be formatted as human-readable text (default) or JSON.

Examples:
  deco show sword-001
  deco show character-hero --json
  deco show quest-001 /path/to/project`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeID := args[0]
			if len(args) > 1 {
				flags.targetDir = args[1]
			} else {
				flags.targetDir = "."
			}
			return runShow(nodeID, flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.jsonOutput, "json", "j", false, "Output as JSON")

	return cmd
}

func runShow(nodeID string, flags *showFlags) error {
	// Load config to verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load all nodes (needed for reverse references)
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	nodes, err := nodeRepo.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}

	// Find the requested node
	var targetNode *domain.Node
	for i, n := range nodes {
		if n.ID == nodeID {
			targetNode = &nodes[i]
			break
		}
	}

	if targetNode == nil {
		return fmt.Errorf("node '%s' not found", nodeID)
	}

	// Build graph and reverse index
	builder := graph.NewBuilder()
	g, err := builder.Build(nodes)
	if err != nil {
		// Even if there are graph errors (cycles, broken refs),
		// we can still show the node
		// Just note that reverse refs might be incomplete
	}

	reverseIndex := builder.BuildReverseIndex(g)

	// Output
	if flags.jsonOutput {
		return outputJSON(targetNode, reverseIndex[nodeID])
	}

	outputHuman(targetNode, reverseIndex[nodeID])
	return nil
}

func outputHuman(node *domain.Node, referencedBy []string) {
	fmt.Printf("Node: %s\n", node.ID)
	fmt.Println(strings.Repeat("=", len("Node: "+node.ID)))
	fmt.Println()

	// Basic fields
	fmt.Printf("Kind:    %s\n", node.Kind)
	fmt.Printf("Version: %d\n", node.Version)
	fmt.Printf("Status:  %s\n", node.Status)
	fmt.Printf("Title:   %s\n", node.Title)

	if node.Summary != "" {
		fmt.Printf("Summary: %s\n", node.Summary)
	}

	// Tags
	if len(node.Tags) > 0 {
		fmt.Printf("Tags:    %s\n", strings.Join(node.Tags, ", "))
	}

	// Content
	if node.Content != nil && len(node.Content.Sections) > 0 {
		fmt.Println()
		fmt.Println("Content:")
		for _, section := range node.Content.Sections {
			fmt.Printf("  [%s]\n", section.Name)
			// For now, just note that it has blocks
			// Full rendering would require block type-specific formatting
			if len(section.Blocks) > 0 {
				fmt.Printf("    (%d block(s))\n", len(section.Blocks))
			}
		}
	}

	// References (what this node uses/relates to)
	if len(node.Refs.Uses) > 0 || len(node.Refs.Related) > 0 {
		fmt.Println()
		fmt.Println("References:")

		if len(node.Refs.Uses) > 0 {
			fmt.Println("  Uses:")
			for _, ref := range node.Refs.Uses {
				if ref.Context != "" {
					fmt.Printf("    - %s (%s)\n", ref.Target, ref.Context)
				} else {
					fmt.Printf("    - %s\n", ref.Target)
				}
			}
		}

		if len(node.Refs.Related) > 0 {
			fmt.Println("  Related:")
			for _, ref := range node.Refs.Related {
				if ref.Context != "" {
					fmt.Printf("    - %s (%s)\n", ref.Target, ref.Context)
				} else {
					fmt.Printf("    - %s\n", ref.Target)
				}
			}
		}
	}

	// Reverse references (what references this node)
	if len(referencedBy) > 0 {
		fmt.Println()
		fmt.Println("Referenced By:")
		for _, refID := range referencedBy {
			fmt.Printf("  - %s\n", refID)
		}
	}

	// Constraints
	if len(node.Constraints) > 0 {
		fmt.Println()
		fmt.Println("Constraints:")
		for _, constraint := range node.Constraints {
			if constraint.Message != "" {
				fmt.Printf("  - %s (%s)\n", constraint.Expr, constraint.Message)
			} else {
				fmt.Printf("  - %s\n", constraint.Expr)
			}
		}
	}

	// Custom fields
	if len(node.Custom) > 0 {
		fmt.Println()
		fmt.Println("Custom Fields:")
		for key, value := range node.Custom {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}

func outputJSON(node *domain.Node, referencedBy []string) error {
	// Create output structure
	output := struct {
		*domain.Node
		ReferencedBy []string `json:"referenced_by"`
	}{
		Node:         node,
		ReferencedBy: referencedBy,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
