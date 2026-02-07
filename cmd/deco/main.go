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

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Toernblom/deco/internal/cli"
)

func main() {
	root := cli.NewRootCommand()

	// Register subcommands
	root.AddCommand(cli.NewInitCommand())
	root.AddCommand(cli.NewIssuesCommand())
	root.AddCommand(cli.NewValidateCommand())
	root.AddCommand(cli.NewListCommand())
	root.AddCommand(cli.NewShowCommand())
	root.AddCommand(cli.NewQueryCommand())
	root.AddCommand(cli.NewHistoryCommand())
	root.AddCommand(cli.NewGraphCommand())
	root.AddCommand(cli.NewDiffCommand())
	root.AddCommand(cli.NewStatsCommand())
	root.AddCommand(cli.NewReviewCommand())
	root.AddCommand(cli.NewSyncCommand())
	root.AddCommand(cli.NewMigrateCommand())
	root.AddCommand(cli.NewLLMHelpCommand())
	root.AddCommand(cli.NewNewCommand())

	if err := root.Execute(); err != nil {
		// Check for ExitError with custom exit code
		// ExitErrors have already displayed their output, so just exit with code
		var exitErr *cli.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.Code)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
