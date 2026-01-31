package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCommand_Structure(t *testing.T) {
	t.Run("creates init command", func(t *testing.T) {
		cmd := NewInitCommand()
		if cmd == nil {
			t.Fatal("Expected init command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "init") {
			t.Errorf("Expected Use to start with 'init', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewInitCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})
}

func TestInitCommand_CreateStructure(t *testing.T) {
	t.Run("creates .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewInitCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		decoDir := filepath.Join(tmpDir, ".deco")
		if _, err := os.Stat(decoDir); os.IsNotExist(err) {
			t.Errorf("Expected .deco directory to be created at %s", decoDir)
		}
	})

	t.Run("creates config.yaml", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewInitCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		configPath := filepath.Join(tmpDir, ".deco", "config.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Errorf("Expected config.yaml to be created at %s", configPath)
		}

		// Verify it's valid YAML
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("Failed to read config.yaml: %v", err)
		}
		if len(content) == 0 {
			t.Error("Expected config.yaml to have content")
		}
	})

	t.Run("creates nodes directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewInitCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		if _, err := os.Stat(nodesDir); os.IsNotExist(err) {
			t.Errorf("Expected nodes directory to be created at %s", nodesDir)
		}
	})

	t.Run("creates all structure in one command", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewInitCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify all parts exist
		decoDir := filepath.Join(tmpDir, ".deco")
		configPath := filepath.Join(decoDir, "config.yaml")
		nodesDir := filepath.Join(decoDir, "nodes")

		for _, path := range []string{decoDir, configPath, nodesDir} {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("Expected %s to exist", path)
			}
		}
	})
}

func TestInitCommand_ExistingProject(t *testing.T) {
	t.Run("detects existing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create existing .deco directory
		decoDir := filepath.Join(tmpDir, ".deco")
		if err := os.Mkdir(decoDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		cmd := NewInitCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()

		if err == nil {
			t.Error("Expected error when .deco already exists, got nil")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, "already initialized") &&
			!strings.Contains(errMsg, "already exists") &&
			!strings.Contains(errMsg, "existing project") {
			t.Errorf("Expected error about existing project, got %q", errMsg)
		}
	})

	t.Run("suggests --force flag on existing project", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create existing .deco directory
		decoDir := filepath.Join(tmpDir, ".deco")
		if err := os.Mkdir(decoDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		cmd := NewInitCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()

		if err == nil {
			t.Error("Expected error, got nil")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, "--force") {
			t.Errorf("Expected error to suggest --force flag, got %q", errMsg)
		}
	})
}

func TestInitCommand_ForceFlag(t *testing.T) {
	t.Run("has force flag", func(t *testing.T) {
		cmd := NewInitCommand()
		flag := cmd.Flags().Lookup("force")
		if flag == nil {
			t.Fatal("Expected --force flag to be defined")
		}
		if flag.Shorthand != "f" {
			t.Errorf("Expected shorthand 'f', got %q", flag.Shorthand)
		}
		if flag.DefValue != "false" {
			t.Errorf("Expected default 'false', got %q", flag.DefValue)
		}
	})

	t.Run("force flag reinitializes existing project", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create existing .deco with a test file
		decoDir := filepath.Join(tmpDir, ".deco")
		if err := os.Mkdir(decoDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
		testFile := filepath.Join(decoDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("existing"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		cmd := NewInitCommand()
		cmd.SetArgs([]string{"--force", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with --force, got %v", err)
		}

		// Verify structure was recreated
		configPath := filepath.Join(decoDir, "config.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Expected config.yaml to be created with --force")
		}

		nodesDir := filepath.Join(decoDir, "nodes")
		if _, err := os.Stat(nodesDir); os.IsNotExist(err) {
			t.Error("Expected nodes directory to be created with --force")
		}
	})

	t.Run("force flag short version works", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create existing .deco directory
		decoDir := filepath.Join(tmpDir, ".deco")
		if err := os.Mkdir(decoDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		cmd := NewInitCommand()
		cmd.SetArgs([]string{"-f", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with -f, got %v", err)
		}
	})
}

func TestInitCommand_ConfigContent(t *testing.T) {
	t.Run("config contains project metadata", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewInitCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		configPath := filepath.Join(tmpDir, ".deco", "config.yaml")
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("Failed to read config.yaml: %v", err)
		}

		contentStr := string(content)
		expectedFields := []string{"version", "project_name"}
		for _, field := range expectedFields {
			if !strings.Contains(contentStr, field) {
				t.Errorf("Expected config to contain %q field, got:\n%s", field, contentStr)
			}
		}
	})
}

func TestInitCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()

		root := NewRootCommand()
		init := NewInitCommand()
		root.AddCommand(init)

		root.SetArgs([]string{"init", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify initialization happened
		decoDir := filepath.Join(tmpDir, ".deco")
		if _, err := os.Stat(decoDir); os.IsNotExist(err) {
			t.Error("Expected .deco directory to be created via root command")
		}
	})
}
