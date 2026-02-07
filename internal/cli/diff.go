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
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/spf13/cobra"
)

type diffFlags struct {
	since     string
	last      int
	targetDir string
}

// NewDiffCommand creates the diff subcommand
func NewDiffCommand() *cobra.Command {
	flags := &diffFlags{}

	cmd := &cobra.Command{
		Use:   "diff <id> [directory]",
		Short: "Show changes to a node over time",
		Long: `Show the change history for a specific node.

Displays before/after values for each change to the node.
Use filters to limit the output.

Examples:
  deco diff player-001                    # Show all changes to player-001
  deco diff player-001 --last 5           # Show last 5 changes
  deco diff player-001 --since 2024-01-01 # Changes since date
  deco diff player-001 --since 2h         # Changes in last 2 hours`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeID := args[0]
			if len(args) > 1 {
				flags.targetDir = args[1]
			} else {
				flags.targetDir = "."
			}
			return runDiff(nodeID, flags)
		},
	}

	cmd.Flags().StringVar(&flags.since, "since", "", "Show changes since timestamp (RFC3339 or relative: 2h, 1d, 1w)")
	cmd.Flags().IntVar(&flags.last, "last", 0, "Show only the last N changes (0 = all)")

	return cmd
}

func runDiff(nodeID string, flags *diffFlags) error {
	// Load config to verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Build filter
	filter := history.Filter{
		NodeID: nodeID,
	}

	// Parse --since if provided
	if flags.since != "" {
		since, err := parseSince(flags.since)
		if err != nil {
			return fmt.Errorf("invalid --since value: %w", err)
		}
		filter.Since = since.Unix()
	}

	// Query history
	historyRepo := history.NewYAMLRepository(config.ResolveHistoryPath(cfg, flags.targetDir))
	entries, err := historyRepo.Query(filter)
	if err != nil {
		return fmt.Errorf("failed to query history: %w", err)
	}

	if len(entries) == 0 {
		fmt.Printf("No changes found for node '%s'\n", nodeID)
		return nil
	}

	// Sort by timestamp descending (most recent first) for --last
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})

	// Apply --last limit
	if flags.last > 0 && len(entries) > flags.last {
		entries = entries[:flags.last]
	}

	// Reverse back to chronological order for display
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	// Display changes
	printDiff(nodeID, entries)
	return nil
}

// parseSince parses a since value that can be:
// - RFC3339 timestamp (2024-01-15T10:30:00Z)
// - Date only (2024-01-15)
// - Relative duration (2h, 1d, 1w)
func parseSince(s string) (time.Time, error) {
	// Try RFC3339 first
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	// Try date only
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}

	// Try relative duration
	return parseRelativeDuration(s)
}

// parseRelativeDuration parses strings like "2h", "1d", "1w"
func parseRelativeDuration(s string) (time.Time, error) {
	if len(s) < 2 {
		return time.Time{}, fmt.Errorf("invalid duration format")
	}

	// Parse the numeric part
	numStr := s[:len(s)-1]
	var num int
	_, err := fmt.Sscanf(numStr, "%d", &num)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid duration number: %s", numStr)
	}

	// Parse the unit
	unit := s[len(s)-1]
	var duration time.Duration
	switch unit {
	case 'h':
		duration = time.Duration(num) * time.Hour
	case 'd':
		duration = time.Duration(num) * 24 * time.Hour
	case 'w':
		duration = time.Duration(num) * 7 * 24 * time.Hour
	case 'm':
		duration = time.Duration(num) * time.Minute
	default:
		return time.Time{}, fmt.Errorf("unknown duration unit: %c (use h, d, w, or m)", unit)
	}

	return time.Now().Add(-duration), nil
}

func printDiff(nodeID string, entries []domain.AuditEntry) {
	fmt.Printf("Changes to %s (%d entries)\n", nodeID, len(entries))
	fmt.Println(strings.Repeat("=", 60))

	for i, entry := range entries {
		fmt.Printf("\n[%d] %s - %s by %s\n",
			i+1,
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Operation,
			entry.User)
		fmt.Println(strings.Repeat("-", 40))

		// Show before/after based on operation
		switch entry.Operation {
		case "create":
			printAfter(entry.After)
		case "delete":
			printBefore(entry.Before)
		default:
			// For update, set, append, unset, move - show both
			if len(entry.Before) > 0 || len(entry.After) > 0 {
				printBeforeAfter(entry.Before, entry.After)
			} else {
				fmt.Println("  (no details recorded)")
			}
		}
	}

	fmt.Println()
}

func printAfter(after map[string]interface{}) {
	if len(after) == 0 {
		fmt.Println("  (no details recorded)")
		return
	}
	fmt.Println("  Created with:")
	printFields(after, "+")
}

func printBefore(before map[string]interface{}) {
	if len(before) == 0 {
		fmt.Println("  (no details recorded)")
		return
	}
	fmt.Println("  Deleted state:")
	printFields(before, "-")
}

func printBeforeAfter(before, after map[string]interface{}) {
	// Collect all keys that changed
	allKeys := make(map[string]bool)
	for k := range before {
		allKeys[k] = true
	}
	for k := range after {
		allKeys[k] = true
	}

	// Sort keys for consistent output
	var keys []string
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		beforeVal, hasBefore := before[key]
		afterVal, hasAfter := after[key]

		if hasBefore && hasAfter {
			// Value changed
			fmt.Printf("  %s:\n", key)
			fmt.Printf("    - %v\n", formatValue(beforeVal))
			fmt.Printf("    + %v\n", formatValue(afterVal))
		} else if hasBefore {
			// Value removed
			fmt.Printf("  %s:\n", key)
			fmt.Printf("    - %v\n", formatValue(beforeVal))
		} else {
			// Value added
			fmt.Printf("  %s:\n", key)
			fmt.Printf("    + %v\n", formatValue(afterVal))
		}
	}
}

func printFields(fields map[string]interface{}, prefix string) {
	// Sort keys for consistent output
	var keys []string
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("  %s %s: %v\n", prefix, key, formatValue(fields[key]))
	}
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case []interface{}:
		if len(val) == 0 {
			return "[]"
		}
		var items []string
		for _, item := range val {
			items = append(items, fmt.Sprintf("%v", item))
		}
		return "[" + strings.Join(items, ", ") + "]"
	case map[string]interface{}:
		if len(val) == 0 {
			return "{}"
		}
		var items []string
		for k, v := range val {
			items = append(items, fmt.Sprintf("%s: %v", k, v))
		}
		return "{" + strings.Join(items, ", ") + "}"
	default:
		return fmt.Sprintf("%v", v)
	}
}
