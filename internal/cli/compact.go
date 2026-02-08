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

	"github.com/Toernblom/deco/internal/domain"
)

// renderNodeCompact renders a full node in compact, LLM-optimized format.
func renderNodeCompact(n domain.Node) string {
	var sb strings.Builder

	// Header: # {id} (v{version}, {status}) [{tags}]
	sb.WriteString(fmt.Sprintf("# %s (v%d, %s)", n.ID, n.Version, n.Status))
	if len(n.Tags) > 0 {
		sb.WriteString(fmt.Sprintf(" [%s]", strings.Join(n.Tags, ", ")))
	}
	sb.WriteString("\n")

	// Title line: {title} — {summary}
	sb.WriteString(n.Title)
	if n.Summary != "" {
		sb.WriteString(" — " + n.Summary)
	}
	sb.WriteString("\n")

	// Refs line: uses: a, b | related: c | emits: x, y | vocabulary: z
	var refParts []string
	if len(n.Refs.Uses) > 0 {
		targets := make([]string, len(n.Refs.Uses))
		for i, r := range n.Refs.Uses {
			targets[i] = r.Target
		}
		refParts = append(refParts, "uses: "+strings.Join(targets, ", "))
	}
	if len(n.Refs.Related) > 0 {
		targets := make([]string, len(n.Refs.Related))
		for i, r := range n.Refs.Related {
			targets[i] = r.Target
		}
		refParts = append(refParts, "related: "+strings.Join(targets, ", "))
	}
	if len(n.Refs.EmitsEvents) > 0 {
		refParts = append(refParts, "emits: "+strings.Join(n.Refs.EmitsEvents, ", "))
	}
	if len(n.Refs.Vocabulary) > 0 {
		refParts = append(refParts, "vocabulary: "+strings.Join(n.Refs.Vocabulary, ", "))
	}
	if len(refParts) > 0 {
		sb.WriteString(strings.Join(refParts, " | ") + "\n")
	}

	// LLM context
	if n.LLMContext != "" {
		sb.WriteString(n.LLMContext + "\n")
	}

	// Glossary
	if len(n.Glossary) > 0 {
		keys := make([]string, 0, len(n.Glossary))
		for k := range n.Glossary {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		pairs := make([]string, len(keys))
		for i, k := range keys {
			pairs[i] = k + "=" + n.Glossary[k]
		}
		sb.WriteString("glossary: " + strings.Join(pairs, ", ") + "\n")
	}

	// Content sections
	if n.Content != nil && len(n.Content.Sections) > 0 {
		for _, section := range n.Content.Sections {
			sb.WriteString("\n## " + section.Name + "\n")
			for _, block := range section.Blocks {
				sb.WriteString(renderBlockCompact(block))
			}
		}
	}

	// Issues
	if len(n.Issues) > 0 {
		sb.WriteString("\n## issues\n")
		for _, issue := range n.Issues {
			severity := issue.Severity
			if severity == "" {
				severity = "medium"
			}
			sb.WriteString(fmt.Sprintf("- %s (%s): %s\n", issue.ID, severity, issue.Description))
		}
	}

	// Contracts
	if len(n.Contracts) > 0 {
		sb.WriteString("\n## contracts\n")
		for _, c := range n.Contracts {
			sb.WriteString(fmt.Sprintf("- %s:", c.Name))
			if c.Scenario != "" {
				sb.WriteString(" " + c.Scenario)
			}
			if len(c.Given) > 0 {
				sb.WriteString(" Given " + strings.Join(c.Given, " And "))
			}
			if len(c.When) > 0 {
				sb.WriteString(" When " + strings.Join(c.When, " And "))
			}
			if len(c.Then) > 0 {
				sb.WriteString(" Then " + strings.Join(c.Then, " And "))
			}
			sb.WriteString("\n")
		}
	}

	// Custom fields
	if len(n.Custom) > 0 {
		keys := make([]string, 0, len(n.Custom))
		for k := range n.Custom {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		pairs := make([]string, len(keys))
		for i, k := range keys {
			pairs[i] = fmt.Sprintf("%s=%v", k, n.Custom[k])
		}
		sb.WriteString("custom: " + strings.Join(pairs, ", ") + "\n")
	}

	sb.WriteString("---\n")
	return sb.String()
}

// renderBlockCompact renders a single block in compact, single-line format.
func renderBlockCompact(block domain.Block) string {
	switch block.Type {
	case "rule":
		return renderRuleCompact(block)
	case "table":
		return renderTableCompact(block)
	case "param":
		return renderParamCompact(block)
	case "mechanic":
		return renderMechanicCompact(block)
	case "list":
		return renderListCompact(block)
	case "text", "note", "description":
		return renderTextCompact(block)
	default:
		return renderDefaultCompact(block)
	}
}

func renderRuleCompact(block domain.Block) string {
	name := getBlockString(block, "name")
	if name == "" {
		name = getBlockString(block, "id")
	}
	text := getBlockString(block, "text")
	if text == "" {
		text = getBlockString(block, "formula")
	}
	if name != "" && text != "" {
		return fmt.Sprintf("- [rule] %s: %s\n", name, text)
	}
	if text != "" {
		return fmt.Sprintf("- [rule] %s\n", text)
	}
	if name != "" {
		return fmt.Sprintf("- [rule] %s\n", name)
	}
	return ""
}

func renderTableCompact(block domain.Block) string {
	id := getBlockString(block, "id")

	columns, ok := block.Data["columns"].([]interface{})
	if !ok || len(columns) == 0 {
		if id != "" {
			return fmt.Sprintf("- [table] %s\n", id)
		}
		return ""
	}

	// Extract column display names
	var colNames []string
	var colKeys []string
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
		colNames = append(colNames, display)
		colKeys = append(colKeys, key)
	}

	rows, ok := block.Data["rows"].([]interface{})
	if !ok {
		rows = nil
	}

	var rowStrs []string
	for _, row := range rows {
		rm, ok := row.(map[string]interface{})
		if !ok {
			continue
		}
		var vals []string
		for _, key := range colKeys {
			vals = append(vals, fmt.Sprintf("%v", rm[key]))
		}
		rowStrs = append(rowStrs, "("+strings.Join(vals, ", ")+")")
	}

	label := id
	if label == "" {
		label = strings.Join(colNames, "/")
	}

	result := fmt.Sprintf("- [table] %s: columns(%s)", label, strings.Join(colNames, ", "))
	if len(rowStrs) > 0 {
		result += " rows: " + strings.Join(rowStrs, ", ")
	}
	return result + "\n"
}

func renderParamCompact(block domain.Block) string {
	name := getBlockString(block, "name")
	if name == "" {
		name = getBlockString(block, "id")
	}

	value := getBlockString(block, "value")
	if value == "" {
		value = getBlockString(block, "default")
	}

	var sb strings.Builder
	if value != "" {
		sb.WriteString(fmt.Sprintf("- [param] %s = %s", name, value))
	} else {
		sb.WriteString(fmt.Sprintf("- [param] %s", name))
	}

	// Constraints inline
	var constraints []string
	if min, ok := block.Data["min"]; ok {
		constraints = append(constraints, fmt.Sprintf("min=%v", min))
	}
	if max, ok := block.Data["max"]; ok {
		constraints = append(constraints, fmt.Sprintf("max=%v", max))
	}
	if def, ok := block.Data["default"]; ok && value != getBlockString(block, "default") {
		constraints = append(constraints, fmt.Sprintf("default=%v", def))
	}
	if unit := getBlockString(block, "unit"); unit != "" {
		constraints = append(constraints, fmt.Sprintf("unit=%s", unit))
	}
	if len(constraints) > 0 {
		sb.WriteString(" [" + strings.Join(constraints, ", ") + "]")
	}

	sb.WriteString("\n")
	return sb.String()
}

func renderMechanicCompact(block domain.Block) string {
	name := getBlockString(block, "name")
	desc := getBlockString(block, "description")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("- [mechanic] %s", name))
	if desc != "" {
		sb.WriteString(": " + desc)
	}

	// Append conditions/inputs/outputs inline if present
	var extras []string
	if conditions, ok := block.Data["conditions"].([]interface{}); ok && len(conditions) > 0 {
		strs := make([]string, len(conditions))
		for i, c := range conditions {
			strs[i] = fmt.Sprintf("%v", c)
		}
		extras = append(extras, "conditions=["+strings.Join(strs, "; ")+"]")
	}
	if inputs, ok := block.Data["inputs"].([]interface{}); ok && len(inputs) > 0 {
		strs := make([]string, len(inputs))
		for i, inp := range inputs {
			strs[i] = fmt.Sprintf("%v", inp)
		}
		extras = append(extras, "inputs=["+strings.Join(strs, "; ")+"]")
	}
	if outputs, ok := block.Data["outputs"].([]interface{}); ok && len(outputs) > 0 {
		strs := make([]string, len(outputs))
		for i, o := range outputs {
			strs[i] = fmt.Sprintf("%v", o)
		}
		extras = append(extras, "outputs=["+strings.Join(strs, "; ")+"]")
	}

	// Also handle trigger/effect pattern from Data
	if trigger := getBlockString(block, "trigger"); trigger != "" {
		extras = append(extras, "trigger="+trigger)
	}
	if effect := getBlockString(block, "effect"); effect != "" {
		extras = append(extras, "effect="+effect)
	}

	if len(extras) > 0 {
		if desc != "" {
			sb.WriteString(", ")
		} else {
			sb.WriteString(": ")
		}
		sb.WriteString(strings.Join(extras, ", "))
	}

	sb.WriteString("\n")
	return sb.String()
}

func renderListCompact(block domain.Block) string {
	id := getBlockString(block, "id")
	items, ok := block.Data["items"].([]interface{})
	if !ok || len(items) == 0 {
		if id != "" {
			return fmt.Sprintf("- [list] %s\n", id)
		}
		return ""
	}

	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = fmt.Sprintf("%v", item)
	}

	if id != "" {
		return fmt.Sprintf("- [list] %s: %s\n", id, strings.Join(strs, ", "))
	}
	return fmt.Sprintf("- [list] %s\n", strings.Join(strs, ", "))
}

func renderTextCompact(block domain.Block) string {
	text := getBlockString(block, "text")
	if text == "" {
		text = getBlockString(block, "content")
	}
	if text == "" {
		return ""
	}
	return fmt.Sprintf("- [text] %s\n", text)
}

func renderDefaultCompact(block domain.Block) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("- [%s]", block.Type))

	if len(block.Data) > 0 {
		keys := make([]string, 0, len(block.Data))
		for k := range block.Data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		pairs := make([]string, len(keys))
		for i, k := range keys {
			pairs[i] = fmt.Sprintf("%s=%v", k, block.Data[k])
		}
		sb.WriteString(" " + strings.Join(pairs, ", "))
	}

	sb.WriteString("\n")
	return sb.String()
}
