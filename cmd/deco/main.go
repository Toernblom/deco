package main

import (
	"fmt"
	"os"

	"github.com/Toernblom/deco/internal/cli"
)

func main() {
	root := cli.NewRootCommand()

	// Register subcommands
	root.AddCommand(cli.NewInitCommand())
	root.AddCommand(cli.NewValidateCommand())

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
