package cli

import (
	"fmt"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/spf13/cobra"
)

type historyFlags struct {
	nodeID    string
	limit     int
	targetDir string
}

// NewHistoryCommand creates the history subcommand
func NewHistoryCommand() *cobra.Command {
	flags := &historyFlags{}

	cmd := &cobra.Command{
		Use:   "history [directory]",
		Short: "Show audit log history",
		Long: `Show the audit log history for the project.

The audit log tracks all changes to nodes including creates, updates, and deletes.
Use filters to narrow down the results.

Examples:
  deco history                       # Show all history
  deco history --node sword-001      # Show history for specific node
  deco history --limit 10            # Show last 10 entries
  deco history -n hero-001 -l 5      # Combined filters`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				flags.targetDir = args[0]
			} else {
				flags.targetDir = "."
			}
			return runHistory(flags)
		},
	}

	cmd.Flags().StringVarP(&flags.nodeID, "node", "n", "", "Filter by node ID")
	cmd.Flags().IntVarP(&flags.limit, "limit", "l", 0, "Limit number of entries (0 = no limit)")

	return cmd
}

func runHistory(flags *historyFlags) error {
	// Load config to verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Query history
	historyRepo := history.NewYAMLRepository(flags.targetDir)
	filter := history.Filter{
		NodeID: flags.nodeID,
		Limit:  flags.limit,
	}
	entries, err := historyRepo.Query(filter)
	if err != nil {
		return fmt.Errorf("failed to query history: %w", err)
	}

	// Display results
	if len(entries) == 0 {
		fmt.Println("No history entries found")
		return nil
	}

	printHistoryTable(entries)
	return nil
}

type historyRow struct {
	time      string
	nodeID    string
	operation string
	user      string
}

func printHistoryTable(entries []domain.AuditEntry) {
	// Calculate column widths
	maxTimeLen := 4 // "TIME"
	maxNodeLen := 4 // "NODE"
	maxOpLen := 9   // "OPERATION"
	maxUserLen := 4 // "USER"

	var rows []historyRow
	for _, entry := range entries {
		row := historyRow{
			time:      entry.Timestamp.Format("2006-01-02 15:04"),
			nodeID:    entry.NodeID,
			operation: entry.Operation,
			user:      entry.User,
		}
		rows = append(rows, row)

		if len(row.time) > maxTimeLen {
			maxTimeLen = len(row.time)
		}
		if len(row.nodeID) > maxNodeLen {
			maxNodeLen = len(row.nodeID)
		}
		if len(row.operation) > maxOpLen {
			maxOpLen = len(row.operation)
		}
		if len(row.user) > maxUserLen {
			maxUserLen = len(row.user)
		}
	}

	// Print header
	header := fmt.Sprintf("%-*s  %-*s  %-*s  %-*s",
		maxTimeLen, "TIME",
		maxNodeLen, "NODE",
		maxOpLen, "OPERATION",
		maxUserLen, "USER")
	fmt.Println(header)

	// Print separator
	separator := strings.Repeat("-", len(header))
	fmt.Println(separator)

	// Print rows
	for _, row := range rows {
		fmt.Printf("%-*s  %-*s  %-*s  %-*s\n",
			maxTimeLen, row.time,
			maxNodeLen, row.nodeID,
			maxOpLen, row.operation,
			maxUserLen, row.user)
	}

	// Print summary
	fmt.Printf("\nTotal: %d entry/entries\n", len(entries))
}
