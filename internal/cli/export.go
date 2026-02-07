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
	"os"
	"path/filepath"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type exportFlags struct {
	format    string
	output    string
	targetDir string
}

// NewExportCommand creates the export subcommand
func NewExportCommand() *cobra.Command {
	flags := &exportFlags{}

	cmd := &cobra.Command{
		Use:   "export [node-id] [directory]",
		Short: "Export nodes as markdown",
		Long: `Export nodes as markdown documents.

Without arguments, exports all nodes to stdout.
With a node ID, exports a single node.
With --output, writes one .md file per node to the specified directory.

Examples:
  deco export systems/combat              # Single node to stdout
  deco export                             # All nodes to stdout
  deco export --output docs/              # Write one .md per node to directory
  deco export --format markdown           # Explicit (markdown is default)`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var nodeID string
			if len(args) >= 1 {
				// Check if first arg looks like a directory (has .deco in it or is last arg)
				if len(args) == 2 {
					nodeID = args[0]
					flags.targetDir = args[1]
				} else if len(args) == 1 {
					// Could be a node ID or a directory
					// Try as node ID first (default target dir is ".")
					nodeID = args[0]
					flags.targetDir = "."
				}
			} else {
				flags.targetDir = "."
			}
			return runExport(nodeID, flags)
		},
	}

	cmd.Flags().StringVar(&flags.format, "format", "markdown", "Export format (markdown)")
	cmd.Flags().StringVar(&flags.output, "output", "", "Output directory (writes one .md per node)")

	return cmd
}

func runExport(nodeID string, flags *exportFlags) error {
	// Load config
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))

	var nodes []domain.Node
	if nodeID != "" {
		n, err := nodeRepo.Load(nodeID)
		if err != nil {
			return fmt.Errorf("node %q not found: %w", nodeID, err)
		}
		nodes = []domain.Node{n}
	} else {
		nodes, err = nodeRepo.LoadAll()
		if err != nil {
			return fmt.Errorf("failed to load nodes: %w", err)
		}
	}

	if len(nodes) == 0 {
		fmt.Println("No nodes found.")
		return nil
	}

	if flags.output != "" {
		return exportToDirectory(nodes, flags.output)
	}

	// Export to stdout
	for i, n := range nodes {
		if i > 0 {
			fmt.Println("---")
			fmt.Println()
		}
		fmt.Print(renderNodeMarkdown(n))
	}

	return nil
}

func exportToDirectory(nodes []domain.Node, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, n := range nodes {
		md := renderNodeMarkdown(n)
		outPath := filepath.Join(outputDir, n.ID+".md")

		// Create parent directories for nested IDs
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", n.ID, err)
		}

		if err := os.WriteFile(outPath, []byte(md), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", outPath, err)
		}
	}

	fmt.Printf("Exported %d node(s) to %s\n", len(nodes), outputDir)
	return nil
}

func renderNodeMarkdown(n domain.Node) string {
	var sb strings.Builder

	// H1 title with metadata
	sb.WriteString("# " + n.Title + "\n\n")
	sb.WriteString(fmt.Sprintf("**%s** | v%d | %s\n\n", n.Kind, n.Version, n.Status))

	// Summary
	if n.Summary != "" {
		sb.WriteString("> " + strings.TrimSpace(n.Summary) + "\n\n")
	}

	// Tags
	if len(n.Tags) > 0 {
		sb.WriteString("**Tags:** " + strings.Join(n.Tags, ", ") + "\n\n")
	}

	// References
	if len(n.Refs.Uses) > 0 || len(n.Refs.Related) > 0 {
		sb.WriteString("## References\n\n")
		if len(n.Refs.Uses) > 0 {
			sb.WriteString("**Uses:**\n")
			for _, ref := range n.Refs.Uses {
				if ref.Context != "" {
					sb.WriteString(fmt.Sprintf("- %s (%s)\n", ref.Target, ref.Context))
				} else {
					sb.WriteString(fmt.Sprintf("- %s\n", ref.Target))
				}
			}
			sb.WriteString("\n")
		}
		if len(n.Refs.Related) > 0 {
			sb.WriteString("**Related:**\n")
			for _, ref := range n.Refs.Related {
				if ref.Context != "" {
					sb.WriteString(fmt.Sprintf("- %s (%s)\n", ref.Target, ref.Context))
				} else {
					sb.WriteString(fmt.Sprintf("- %s\n", ref.Target))
				}
			}
			sb.WriteString("\n")
		}
	}

	// Content sections
	if n.Content != nil && len(n.Content.Sections) > 0 {
		for _, section := range n.Content.Sections {
			sb.WriteString("## " + section.Name + "\n\n")
			for _, block := range section.Blocks {
				sb.WriteString(RenderBlockMarkdown(block))
				sb.WriteString("\n")
			}
		}
	}

	// Issues
	if len(n.Issues) > 0 {
		sb.WriteString("## Issues\n\n")
		for _, issue := range n.Issues {
			severity := issue.Severity
			if severity == "" {
				severity = "medium"
			}
			icon := severityIcon(severity)
			sb.WriteString(fmt.Sprintf("- %s **[%s]** %s\n", icon, severity, issue.Description))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func severityIcon(severity string) string {
	switch severity {
	case "critical":
		return "!!!"
	case "high":
		return "!!"
	case "low":
		return "."
	default:
		return "!"
	}
}
