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
	"sort"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/query"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"gopkg.in/yaml.v3"
)

// runObsidianExport exports all nodes as an Obsidian vault to .deco/vault/.
func runObsidianExport(flags *exportFlags) error {
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))
	allNodes, err := nodeRepo.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}

	// Apply filters
	nodes := allNodes
	if flags.kind != "" || flags.status != "" || flags.tag != "" {
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
		qe := query.New()
		nodes = qe.Filter(allNodes, criteria)
	}

	if len(nodes) == 0 {
		fmt.Println("No nodes found.")
		return nil
	}

	// Build title lookup from ALL nodes (for wikilink resolution)
	titleMap := buildTitleMap(allNodes)

	// Resolve vault path
	vaultDir := filepath.Join(flags.targetDir, ".deco", "vault")

	// Wipe and recreate vault directory
	if err := os.RemoveAll(vaultDir); err != nil {
		return fmt.Errorf("failed to clean vault directory: %w", err)
	}

	for _, n := range nodes {
		md := renderNodeObsidian(n, titleMap)
		outPath := filepath.Join(vaultDir, n.ID+".md")

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", n.ID, err)
		}

		if err := os.WriteFile(outPath, []byte(md), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", outPath, err)
		}
	}

	fmt.Printf("Exported %d node(s) to %s\n", len(nodes), vaultDir)
	return nil
}

// buildTitleMap creates a map of node ID to title for wikilink resolution.
func buildTitleMap(nodes []domain.Node) map[string]string {
	m := make(map[string]string, len(nodes))
	for _, n := range nodes {
		m[n.ID] = n.Title
	}
	return m
}

// wikilink renders a node reference as an Obsidian wikilink.
// If the target exists in titleMap, uses [[id|Title]], otherwise [[id]].
func wikilink(target string, titleMap map[string]string) string {
	if title, ok := titleMap[target]; ok {
		return fmt.Sprintf("[[%s|%s]]", target, title)
	}
	return fmt.Sprintf("[[%s]]", target)
}

// renderNodeObsidian renders a single node as Obsidian-compatible markdown.
func renderNodeObsidian(n domain.Node, titleMap map[string]string) string {
	var sb strings.Builder

	// YAML frontmatter
	sb.WriteString(renderObsidianFrontmatter(n))

	// H1 title
	sb.WriteString("# " + n.Title + "\n\n")

	// Summary
	if n.Summary != "" {
		sb.WriteString(strings.TrimSpace(n.Summary) + "\n\n")
	}

	// Inline tags (Obsidian-native #tag format)
	if len(n.Tags) > 0 {
		tags := make([]string, len(n.Tags))
		for i, t := range n.Tags {
			tags[i] = "#" + t
		}
		sb.WriteString(strings.Join(tags, " ") + "\n\n")
	}

	// Glossary
	if len(n.Glossary) > 0 {
		sb.WriteString("## Glossary\n\n")
		keys := make([]string, 0, len(n.Glossary))
		for k := range n.Glossary {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("**%s**\n: %s\n\n", k, n.Glossary[k]))
		}
	}

	// Content sections
	if n.Content != nil && len(n.Content.Sections) > 0 {
		for _, section := range n.Content.Sections {
			sb.WriteString("## " + section.Name + "\n\n")
			for _, block := range section.Blocks {
				sb.WriteString(renderBlockObsidian(block))
				sb.WriteString("\n")
			}
		}
	}

	// References
	hasRefs := len(n.Refs.Uses) > 0 || len(n.Refs.Related) > 0 ||
		len(n.Refs.EmitsEvents) > 0 || len(n.Refs.Vocabulary) > 0
	if hasRefs {
		sb.WriteString("## References\n\n")

		if len(n.Refs.Uses) > 0 {
			sb.WriteString("**Uses:**\n")
			for _, ref := range n.Refs.Uses {
				link := wikilink(ref.Target, titleMap)
				if ref.Context != "" {
					sb.WriteString(fmt.Sprintf("- %s — %s\n", link, ref.Context))
				} else {
					sb.WriteString(fmt.Sprintf("- %s\n", link))
				}
			}
			sb.WriteString("\n")
		}

		if len(n.Refs.Related) > 0 {
			sb.WriteString("**Related:**\n")
			for _, ref := range n.Refs.Related {
				link := wikilink(ref.Target, titleMap)
				if ref.Context != "" {
					sb.WriteString(fmt.Sprintf("- %s — %s\n", link, ref.Context))
				} else {
					sb.WriteString(fmt.Sprintf("- %s\n", link))
				}
			}
			sb.WriteString("\n")
		}

		if len(n.Refs.Vocabulary) > 0 {
			sb.WriteString("**Vocabulary:**\n")
			for _, v := range n.Refs.Vocabulary {
				sb.WriteString(fmt.Sprintf("- %s\n", wikilink(v, titleMap)))
			}
			sb.WriteString("\n")
		}

		if len(n.Refs.EmitsEvents) > 0 {
			sb.WriteString("**Emits Events:**\n")
			for _, e := range n.Refs.EmitsEvents {
				sb.WriteString(fmt.Sprintf("- `%s`\n", e))
			}
			sb.WriteString("\n")
		}
	}

	// Issues as callouts
	if len(n.Issues) > 0 {
		sb.WriteString("## Issues\n\n")
		for _, issue := range n.Issues {
			if issue.Resolved {
				continue
			}
			calloutType := severityCallout(issue.Severity)
			sb.WriteString(fmt.Sprintf("> [!%s] %s\n", calloutType, issue.Description))
			if issue.Location != "" {
				sb.WriteString(fmt.Sprintf("> Location: `%s`\n", issue.Location))
			}
			sb.WriteString("\n")
		}
	}

	// Custom fields
	if len(n.Custom) > 0 {
		sb.WriteString("## Custom Fields\n\n")
		keys := make([]string, 0, len(n.Custom))
		for k := range n.Custom {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("**%s:** %v\n", k, n.Custom[k]))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderObsidianFrontmatter renders YAML frontmatter for a node.
func renderObsidianFrontmatter(n domain.Node) string {
	fm := make(map[string]interface{})
	fm["id"] = n.ID
	fm["kind"] = n.Kind
	fm["version"] = n.Version
	fm["status"] = n.Status

	if len(n.Tags) > 0 {
		fm["tags"] = n.Tags
	}
	if n.Summary != "" {
		fm["summary"] = strings.TrimSpace(n.Summary)
	}

	// Refs as flat lists for Dataview queryability
	if len(n.Refs.Uses) > 0 {
		uses := make([]string, len(n.Refs.Uses))
		for i, ref := range n.Refs.Uses {
			uses[i] = ref.Target
		}
		fm["uses"] = uses
	}
	if len(n.Refs.Related) > 0 {
		related := make([]string, len(n.Refs.Related))
		for i, ref := range n.Refs.Related {
			related[i] = ref.Target
		}
		fm["related"] = related
	}
	if len(n.Refs.EmitsEvents) > 0 {
		fm["emits_events"] = n.Refs.EmitsEvents
	}
	if len(n.Refs.Vocabulary) > 0 {
		fm["vocabulary"] = n.Refs.Vocabulary
	}

	data, err := yaml.Marshal(fm)
	if err != nil {
		return "---\n---\n"
	}
	return "---\n" + string(data) + "---\n\n"
}

// renderBlockObsidian renders a content block as Obsidian-compatible markdown.
func renderBlockObsidian(block domain.Block) string {
	switch block.Type {
	case "table":
		return renderTableBlockMarkdown(block)
	case "rule":
		return renderRuleBlockObsidian(block)
	case "param":
		return renderParamBlockMarkdown(block)
	case "list":
		return renderListBlockMarkdown(block)
	case "text", "note", "description":
		return renderTextBlockMarkdown(block)
	case "mechanic":
		return renderMechanicBlockMarkdown(block)
	default:
		return renderCustomBlockObsidian(block)
	}
}

// renderRuleBlockObsidian renders a rule block as an Obsidian callout.
func renderRuleBlockObsidian(block domain.Block) string {
	name := getBlockString(block, "name")
	text := getBlockString(block, "text")
	if text == "" {
		return ""
	}
	if name != "" {
		return fmt.Sprintf("> [!note] %s\n> %s\n\n", name, text)
	}
	return "> " + text + "\n"
}

// renderCustomBlockObsidian renders unknown/custom block types as structured definition lists.
func renderCustomBlockObsidian(block domain.Block) string {
	if len(block.Data) == 0 {
		return ""
	}

	var sb strings.Builder

	// Try to find a "name" or "id" field to use as a heading
	name := getBlockString(block, "name")
	if name == "" {
		name = getBlockString(block, "id")
	}

	if name != "" {
		sb.WriteString(fmt.Sprintf("**%s** *(%s)*\n", name, block.Type))
	} else {
		sb.WriteString(fmt.Sprintf("*(%s)*\n", block.Type))
	}

	// Render remaining fields as a definition list
	keys := make([]string, 0, len(block.Data))
	for k := range block.Data {
		if k == "name" || k == "id" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := block.Data[k]
		sb.WriteString(fmt.Sprintf("- **%s:** %v\n", k, v))
	}

	return sb.String()
}

// severityCallout maps issue severity to Obsidian callout types.
func severityCallout(severity string) string {
	switch severity {
	case "critical":
		return "bug"
	case "high":
		return "warning"
	case "medium":
		return "caution"
	case "low":
		return "info"
	default:
		return "info"
	}
}
