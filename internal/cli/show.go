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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Toernblom/deco/internal/cli/style"
	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/graph"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type showFlags struct {
	jsonOutput bool
	full       bool
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
	cmd.Flags().BoolVar(&flags.full, "full", false, "Expand content blocks inline")

	return cmd
}

func runShow(nodeID string, flags *showFlags) error {
	// Load config to verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load all nodes (needed for reverse references)
	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))
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

	outputHuman(targetNode, reverseIndex[nodeID], flags.full)
	return nil
}

func outputHuman(node *domain.Node, referencedBy []string, full bool) {
	fmt.Printf("%s %s\n", style.Header.Sprint("Node:"), node.ID)
	fmt.Println(style.Muted.Sprint(strings.Repeat("═", len("Node: "+node.ID))))
	fmt.Println()

	// Basic fields
	fmt.Printf("%s    %s\n", style.Muted.Sprint("Kind:"), node.Kind)
	fmt.Printf("%s %d\n", style.Muted.Sprint("Version:"), node.Version)

	// Color the status
	statusStr := node.Status
	if c := style.StatusColor(node.Status); c != nil {
		statusStr = c.Sprint(node.Status)
	}
	fmt.Printf("%s  %s\n", style.Muted.Sprint("Status:"), statusStr)
	fmt.Printf("%s   %s\n", style.Muted.Sprint("Title:"), node.Title)

	if node.Summary != "" {
		fmt.Printf("%s %s\n", style.Muted.Sprint("Summary:"), node.Summary)
	}

	// Tags
	if len(node.Tags) > 0 {
		fmt.Printf("%s    %s\n", style.Muted.Sprint("Tags:"), style.Info.Sprint(strings.Join(node.Tags, ", ")))
	}

	// Docs
	if len(node.Docs) > 0 {
		fmt.Println()
		fmt.Println(style.Header.Sprint("Docs:"))
		for _, doc := range node.Docs {
			if len(doc.Keywords) > 0 {
				fmt.Printf("  %s %s %s\n", style.SymbolBullet, doc.Path, style.Muted.Sprintf("(keywords: %s)", strings.Join(doc.Keywords, ", ")))
			} else {
				fmt.Printf("  %s %s\n", style.SymbolBullet, doc.Path)
			}
			if doc.Context != "" {
				fmt.Printf("    %s\n", style.Muted.Sprint(doc.Context))
			}
		}
	}

	// Content
	if node.Content != nil && len(node.Content.Sections) > 0 {
		fmt.Println()
		fmt.Println(style.Header.Sprint("Content:"))
		for _, section := range node.Content.Sections {
			fmt.Printf("  %s\n", style.Info.Sprintf("[%s]", section.Name))
			if full {
				for _, block := range section.Blocks {
					renderBlock(block)
				}
			} else if len(section.Blocks) > 0 {
				fmt.Printf("    %s\n", style.Muted.Sprintf("(%d block(s))", len(section.Blocks)))
			}
		}
	}

	// References (what this node uses/relates to)
	if len(node.Refs.Uses) > 0 || len(node.Refs.Related) > 0 {
		fmt.Println()
		fmt.Println(style.Header.Sprint("References:"))

		if len(node.Refs.Uses) > 0 {
			fmt.Printf("  %s\n", style.Muted.Sprint("Uses:"))
			for _, ref := range node.Refs.Uses {
				if ref.Context != "" {
					fmt.Printf("    %s %s %s\n", style.SymbolBullet, ref.Target, style.Muted.Sprintf("(%s)", ref.Context))
				} else {
					fmt.Printf("    %s %s\n", style.SymbolBullet, ref.Target)
				}
			}
		}

		if len(node.Refs.Related) > 0 {
			fmt.Printf("  %s\n", style.Muted.Sprint("Related:"))
			for _, ref := range node.Refs.Related {
				if ref.Context != "" {
					fmt.Printf("    %s %s %s\n", style.SymbolBullet, ref.Target, style.Muted.Sprintf("(%s)", ref.Context))
				} else {
					fmt.Printf("    %s %s\n", style.SymbolBullet, ref.Target)
				}
			}
		}
	}

	// Reverse references (what references this node)
	if len(referencedBy) > 0 {
		fmt.Println()
		fmt.Println(style.Header.Sprint("Referenced By:"))
		for _, refID := range referencedBy {
			fmt.Printf("  %s %s\n", style.SymbolBullet, refID)
		}
	}

	// Constraints
	if len(node.Constraints) > 0 {
		fmt.Println()
		fmt.Println(style.Header.Sprint("Constraints:"))
		for _, constraint := range node.Constraints {
			if constraint.Message != "" {
				fmt.Printf("  %s %s %s\n", style.SymbolBullet, style.Code.Sprint(constraint.Expr), style.Muted.Sprintf("(%s)", constraint.Message))
			} else {
				fmt.Printf("  %s %s\n", style.SymbolBullet, style.Code.Sprint(constraint.Expr))
			}
		}
	}

	// Custom fields
	if len(node.Custom) > 0 {
		fmt.Println()
		fmt.Println(style.Header.Sprint("Custom Fields:"))
		for key, value := range node.Custom {
			fmt.Printf("  %s: %v\n", style.Muted.Sprint(key), value)
		}
	}
}

// renderBlock renders a single content block based on its type.
func renderBlock(block domain.Block) {
	switch block.Type {
	case "table":
		renderTableBlock(block)
	case "rule":
		renderRuleBlock(block)
	case "param":
		renderParamBlock(block)
	case "list":
		renderListBlock(block)
	case "text", "note", "description":
		renderTextBlock(block)
	case "mechanic":
		renderMechanicBlock(block)
	default:
		renderFallbackBlock(block)
	}
}

func renderTableBlock(block domain.Block) {
	id := getBlockString(block, "id")
	if id != "" {
		fmt.Printf("    %s %s\n", style.Muted.Sprint("table:"), id)
	}

	columns, ok := block.Data["columns"].([]interface{})
	if !ok || len(columns) == 0 {
		return
	}
	rows, ok := block.Data["rows"].([]interface{})
	if !ok {
		return
	}

	// Extract column keys and display names
	type colInfo struct {
		key     string
		display string
	}
	var cols []colInfo
	for _, c := range columns {
		cm, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := cm["key"].(string)
		display, _ := cm["display"].(string)
		if display == "" {
			display = key
		}
		cols = append(cols, colInfo{key: key, display: display})
	}

	if len(cols) == 0 {
		return
	}

	// Compute column widths
	widths := make([]int, len(cols))
	for i, col := range cols {
		widths[i] = len(col.display)
	}
	for _, row := range rows {
		rm, ok := row.(map[string]interface{})
		if !ok {
			continue
		}
		for i, col := range cols {
			val := fmt.Sprintf("%v", rm[col.key])
			if len(val) > widths[i] {
				widths[i] = len(val)
			}
		}
	}

	// Print header
	var header, separator []string
	for i, col := range cols {
		header = append(header, padRight(col.display, widths[i]))
		separator = append(separator, strings.Repeat("-", widths[i]))
	}
	fmt.Printf("    | %s |\n", strings.Join(header, " | "))
	fmt.Printf("    | %s |\n", strings.Join(separator, " | "))

	// Print rows
	for _, row := range rows {
		rm, ok := row.(map[string]interface{})
		if !ok {
			continue
		}
		var vals []string
		for i, col := range cols {
			val := fmt.Sprintf("%v", rm[col.key])
			vals = append(vals, padRight(val, widths[i]))
		}
		fmt.Printf("    | %s |\n", strings.Join(vals, " | "))
	}
	fmt.Println()
}

func renderRuleBlock(block domain.Block) {
	text := getBlockString(block, "text")
	if text != "" {
		fmt.Printf("    %s %s\n", style.Muted.Sprint(">"), text)
	}
}

func renderParamBlock(block domain.Block) {
	name := getBlockString(block, "name")
	datatype := getBlockString(block, "datatype")
	if name == "" {
		name = getBlockString(block, "id")
	}

	parts := []string{name + ":"}
	if datatype != "" {
		parts = append(parts, datatype)
	}

	// Build constraints string
	var constraints []string
	if min, ok := block.Data["min"]; ok {
		constraints = append(constraints, fmt.Sprintf("min=%v", min))
	}
	if max, ok := block.Data["max"]; ok {
		constraints = append(constraints, fmt.Sprintf("max=%v", max))
	}
	if def, ok := block.Data["default"]; ok {
		constraints = append(constraints, fmt.Sprintf("default=%v", def))
	}
	if unit := getBlockString(block, "unit"); unit != "" {
		constraints = append(constraints, unit)
	}
	if len(constraints) > 0 {
		parts = append(parts, "["+strings.Join(constraints, ", ")+"]")
	}

	desc := getBlockString(block, "description")
	if desc != "" {
		parts = append(parts, style.Muted.Sprintf("- %s", desc))
	}

	fmt.Printf("    %s\n", strings.Join(parts, " "))
}

func renderListBlock(block domain.Block) {
	id := getBlockString(block, "id")
	if id != "" {
		fmt.Printf("    %s\n", style.Muted.Sprint(id+":"))
	}
	items, ok := block.Data["items"].([]interface{})
	if !ok {
		return
	}
	for _, item := range items {
		fmt.Printf("      %s %v\n", style.SymbolBullet, item)
	}
}

func renderTextBlock(block domain.Block) {
	text := getBlockString(block, "text")
	if text == "" {
		text = getBlockString(block, "content")
	}
	if text != "" {
		fmt.Printf("    %s\n", text)
	}
}

func renderMechanicBlock(block domain.Block) {
	name := getBlockString(block, "name")
	desc := getBlockString(block, "description")

	fmt.Printf("    %s", style.Header.Sprint(name))
	if desc != "" {
		fmt.Printf(" - %s", desc)
	}
	fmt.Println()

	if conditions, ok := block.Data["conditions"].([]interface{}); ok && len(conditions) > 0 {
		fmt.Printf("      %s\n", style.Muted.Sprint("Conditions:"))
		for _, c := range conditions {
			fmt.Printf("        %s %v\n", style.SymbolBullet, c)
		}
	}
	if inputs, ok := block.Data["inputs"].([]interface{}); ok && len(inputs) > 0 {
		fmt.Printf("      %s\n", style.Muted.Sprint("Inputs:"))
		for _, inp := range inputs {
			fmt.Printf("        %s %v\n", style.SymbolBullet, inp)
		}
	}
	if outputs, ok := block.Data["outputs"].([]interface{}); ok && len(outputs) > 0 {
		fmt.Printf("      %s\n", style.Muted.Sprint("Outputs:"))
		for _, o := range outputs {
			fmt.Printf("        %s %v\n", style.SymbolBullet, o)
		}
	}
}

func renderFallbackBlock(block domain.Block) {
	fmt.Printf("    %s %s\n", style.Muted.Sprint("type:"), block.Type)
	data, err := yaml.Marshal(block.Data)
	if err != nil {
		return
	}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		fmt.Printf("      %s\n", line)
	}
}

func getBlockString(block domain.Block, key string) string {
	v, ok := block.Data[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return s
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// RenderBlockMarkdown renders a single content block as markdown.
// Used by both show --full (terminal) and export (markdown) commands.
func RenderBlockMarkdown(block domain.Block) string {
	switch block.Type {
	case "table":
		return renderTableBlockMarkdown(block)
	case "rule":
		return renderRuleBlockMarkdown(block)
	case "param":
		return renderParamBlockMarkdown(block)
	case "list":
		return renderListBlockMarkdown(block)
	case "text", "note", "description":
		return renderTextBlockMarkdown(block)
	case "mechanic":
		return renderMechanicBlockMarkdown(block)
	default:
		return renderFallbackBlockMarkdown(block)
	}
}

func renderTableBlockMarkdown(block domain.Block) string {
	columns, ok := block.Data["columns"].([]interface{})
	if !ok || len(columns) == 0 {
		return ""
	}
	rows, ok := block.Data["rows"].([]interface{})
	if !ok {
		return ""
	}

	type colInfo struct {
		key     string
		display string
	}
	var cols []colInfo
	for _, c := range columns {
		cm, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := cm["key"].(string)
		display, _ := cm["display"].(string)
		if display == "" {
			display = key
		}
		cols = append(cols, colInfo{key: key, display: display})
	}

	var sb strings.Builder
	// Header
	var headers, seps []string
	for _, col := range cols {
		headers = append(headers, col.display)
		seps = append(seps, "---")
	}
	sb.WriteString("| " + strings.Join(headers, " | ") + " |\n")
	sb.WriteString("| " + strings.Join(seps, " | ") + " |\n")

	for _, row := range rows {
		rm, ok := row.(map[string]interface{})
		if !ok {
			continue
		}
		var vals []string
		for _, col := range cols {
			vals = append(vals, fmt.Sprintf("%v", rm[col.key]))
		}
		sb.WriteString("| " + strings.Join(vals, " | ") + " |\n")
	}
	return sb.String()
}

func renderRuleBlockMarkdown(block domain.Block) string {
	text := getBlockString(block, "text")
	if text == "" {
		return ""
	}
	return "> " + text + "\n"
}

func renderParamBlockMarkdown(block domain.Block) string {
	name := getBlockString(block, "name")
	datatype := getBlockString(block, "datatype")
	if name == "" {
		name = getBlockString(block, "id")
	}

	parts := []string{"**" + name + "**:"}
	if datatype != "" {
		parts = append(parts, "`"+datatype+"`")
	}

	var constraints []string
	if min, ok := block.Data["min"]; ok {
		constraints = append(constraints, fmt.Sprintf("min=%v", min))
	}
	if max, ok := block.Data["max"]; ok {
		constraints = append(constraints, fmt.Sprintf("max=%v", max))
	}
	if def, ok := block.Data["default"]; ok {
		constraints = append(constraints, fmt.Sprintf("default=%v", def))
	}
	if unit := getBlockString(block, "unit"); unit != "" {
		constraints = append(constraints, unit)
	}
	if len(constraints) > 0 {
		parts = append(parts, "["+strings.Join(constraints, ", ")+"]")
	}

	desc := getBlockString(block, "description")
	if desc != "" {
		parts = append(parts, "- "+desc)
	}

	return "- " + strings.Join(parts, " ") + "\n"
}

func renderListBlockMarkdown(block domain.Block) string {
	items, ok := block.Data["items"].([]interface{})
	if !ok {
		return ""
	}
	var sb strings.Builder
	for _, item := range items {
		sb.WriteString(fmt.Sprintf("- %v\n", item))
	}
	return sb.String()
}

func renderTextBlockMarkdown(block domain.Block) string {
	text := getBlockString(block, "text")
	if text == "" {
		text = getBlockString(block, "content")
	}
	if text == "" {
		return ""
	}
	return text + "\n"
}

func renderMechanicBlockMarkdown(block domain.Block) string {
	name := getBlockString(block, "name")
	desc := getBlockString(block, "description")

	var sb strings.Builder
	sb.WriteString("**" + name + "**")
	if desc != "" {
		sb.WriteString(": " + desc)
	}
	sb.WriteString("\n\n")

	if conditions, ok := block.Data["conditions"].([]interface{}); ok && len(conditions) > 0 {
		sb.WriteString("*Conditions:*\n")
		for _, c := range conditions {
			sb.WriteString(fmt.Sprintf("- %v\n", c))
		}
		sb.WriteString("\n")
	}
	if inputs, ok := block.Data["inputs"].([]interface{}); ok && len(inputs) > 0 {
		sb.WriteString("*Inputs:*\n")
		for _, inp := range inputs {
			sb.WriteString(fmt.Sprintf("- %v\n", inp))
		}
		sb.WriteString("\n")
	}
	if outputs, ok := block.Data["outputs"].([]interface{}); ok && len(outputs) > 0 {
		sb.WriteString("*Outputs:*\n")
		for _, o := range outputs {
			sb.WriteString(fmt.Sprintf("- %v\n", o))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func renderFallbackBlockMarkdown(block domain.Block) string {
	data, err := yaml.Marshal(block.Data)
	if err != nil {
		return ""
	}
	return "```yaml\ntype: " + block.Type + "\n" + string(data) + "```\n"
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
