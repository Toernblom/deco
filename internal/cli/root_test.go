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
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	t.Run("returns root command", func(t *testing.T) {
		cmd := NewRootCommand()
		if cmd == nil {
			t.Fatal("Expected root command, got nil")
		}
		if cmd.Use != "deco" {
			t.Errorf("Expected Use 'deco', got %q", cmd.Use)
		}
	})

	t.Run("has short description", func(t *testing.T) {
		cmd := NewRootCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("has long description", func(t *testing.T) {
		cmd := NewRootCommand()
		if cmd.Long == "" {
			t.Error("Expected non-empty Long description")
		}
	})
}

func TestVersionFlag(t *testing.T) {
	t.Run("shows version with --version", func(t *testing.T) {
		cmd := NewRootCommand()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"--version"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "deco") {
			t.Errorf("Expected version output to contain 'deco', got %q", output)
		}
	})

	t.Run("shows version with -v", func(t *testing.T) {
		cmd := NewRootCommand()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"-v"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "deco") {
			t.Errorf("Expected version output to contain 'deco', got %q", output)
		}
	})
}

func TestHelpOutput(t *testing.T) {
	t.Run("shows help with --help", func(t *testing.T) {
		cmd := NewRootCommand()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"--help"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "Usage:") {
			t.Errorf("Expected help to contain 'Usage:', got %q", output)
		}
		if !strings.Contains(output, "deco") {
			t.Errorf("Expected help to contain 'deco', got %q", output)
		}
	})

	t.Run("shows help with -h", func(t *testing.T) {
		cmd := NewRootCommand()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"-h"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "Usage:") {
			t.Errorf("Expected help to contain 'Usage:', got %q", output)
		}
	})

	t.Run("shows help when no args", func(t *testing.T) {
		cmd := NewRootCommand()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "Usage:") {
			t.Errorf("Expected help to contain 'Usage:', got %q", output)
		}
	})
}

func TestGlobalFlags(t *testing.T) {
	t.Run("has config flag", func(t *testing.T) {
		cmd := NewRootCommand()
		flag := cmd.PersistentFlags().Lookup("config")
		if flag == nil {
			t.Fatal("Expected --config flag to be defined")
		}
		if flag.Shorthand != "c" {
			t.Errorf("Expected shorthand 'c', got %q", flag.Shorthand)
		}
		if flag.DefValue != ".deco" {
			t.Errorf("Expected default '.deco', got %q", flag.DefValue)
		}
	})

	t.Run("has verbose flag", func(t *testing.T) {
		cmd := NewRootCommand()
		flag := cmd.PersistentFlags().Lookup("verbose")
		if flag == nil {
			t.Fatal("Expected --verbose flag to be defined")
		}
		if flag.Shorthand != "" {
			t.Errorf("Expected no shorthand, got %q", flag.Shorthand)
		}
		if flag.DefValue != "false" {
			t.Errorf("Expected default 'false', got %q", flag.DefValue)
		}
	})

	t.Run("has quiet flag", func(t *testing.T) {
		cmd := NewRootCommand()
		flag := cmd.PersistentFlags().Lookup("quiet")
		if flag == nil {
			t.Fatal("Expected --quiet flag to be defined")
		}
		if flag.Shorthand != "q" {
			t.Errorf("Expected shorthand 'q', got %q", flag.Shorthand)
		}
		if flag.DefValue != "false" {
			t.Errorf("Expected default 'false', got %q", flag.DefValue)
		}
	})
}

func TestSubcommandRegistration(t *testing.T) {
	t.Run("accepts subcommands", func(t *testing.T) {
		root := NewRootCommand()

		// Create a mock subcommand
		mockCmd := &cobra.Command{
			Use:   "test",
			Short: "Test subcommand",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		root.AddCommand(mockCmd)

		found := false
		for _, cmd := range root.Commands() {
			if cmd.Use == "test" {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected subcommand 'test' to be registered")
		}
	})

	t.Run("has cobra.Command type", func(t *testing.T) {
		cmd := NewRootCommand()
		// This will fail to compile if not *cobra.Command
		_ = cmd.Commands()
	})
}
