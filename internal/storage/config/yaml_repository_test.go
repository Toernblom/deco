package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Toernblom/deco/internal/storage/config"
)

func TestYAMLRepository_Load(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test config file
	configPath := filepath.Join(tmpDir, ".deco", "config.yaml")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	configContent := `project_name: "Test Game"
nodes_path: ".deco/nodes"
history_path: ".deco/history.jsonl"
version: 1
custom_block_types:
  quest:
    required_fields:
      - name
    optional_fields:
      - reward
custom:
  author: "Test Author"
  tags:
    - rpg
    - strategy
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	repo := config.NewYAMLRepository(tmpDir)

	// Load config
	cfg, err := repo.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify fields
	if cfg.ProjectName != "Test Game" {
		t.Errorf("Expected ProjectName 'Test Game', got %q", cfg.ProjectName)
	}
	if cfg.NodesPath != ".deco/nodes" {
		t.Errorf("Expected NodesPath '.deco/nodes', got %q", cfg.NodesPath)
	}
	if cfg.HistoryPath != ".deco/history.jsonl" {
		t.Errorf("Expected HistoryPath '.deco/history.jsonl', got %q", cfg.HistoryPath)
	}
	if cfg.Version != 1 {
		t.Errorf("Expected Version 1, got %d", cfg.Version)
	}
	if cfg.CustomBlockTypes == nil {
		t.Fatal("Expected custom block types to be loaded")
	}
	questCfg, ok := cfg.CustomBlockTypes["quest"]
	if !ok {
		t.Fatal("Expected quest block type to be present")
	}
	if len(questCfg.RequiredFields) != 1 || questCfg.RequiredFields[0] != "name" {
		t.Error("Expected custom block type required_fields to be loaded")
	}
	if len(questCfg.OptionalFields) != 1 || questCfg.OptionalFields[0] != "reward" {
		t.Error("Expected custom block type optional_fields to be loaded")
	}

	// Verify custom fields
	if cfg.Custom == nil {
		t.Error("Expected Custom map to be populated")
	} else {
		if author, ok := cfg.Custom["author"].(string); ok {
			if author != "Test Author" {
				t.Errorf("Expected author 'Test Author', got %q", author)
			}
		} else {
			t.Error("Expected author in custom fields")
		}
	}
}

func TestYAMLRepository_Load_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	repo := config.NewYAMLRepository(tmpDir)

	// Try to load non-existent config
	_, err := repo.Load()
	if err == nil {
		t.Error("Expected error when loading non-existent config")
	}
}

func TestYAMLRepository_Load_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	// Create invalid YAML config
	configPath := filepath.Join(tmpDir, ".deco", "config.yaml")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	err = os.WriteFile(configPath, []byte("invalid: yaml: [content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	repo := config.NewYAMLRepository(tmpDir)

	// Try to load invalid config
	_, err = repo.Load()
	if err == nil {
		t.Error("Expected error when loading invalid YAML config")
	}
}

func TestYAMLRepository_Save(t *testing.T) {
	tmpDir := t.TempDir()

	repo := config.NewYAMLRepository(tmpDir)

	// Create and save config
	cfg := config.Config{
		ProjectName: "New Game",
		NodesPath:   ".deco/nodes",
		HistoryPath: ".deco/history.jsonl",
		Version:     1,
		Custom: map[string]interface{}{
			"author": "Alice",
			"genre":  "RPG",
		},
	}

	err := repo.Save(cfg)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file was created
	configPath := filepath.Join(tmpDir, ".deco", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}

	// Load and verify
	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loaded.ProjectName != cfg.ProjectName {
		t.Errorf("Expected ProjectName %q, got %q", cfg.ProjectName, loaded.ProjectName)
	}
	if loaded.NodesPath != cfg.NodesPath {
		t.Errorf("Expected NodesPath %q, got %q", cfg.NodesPath, loaded.NodesPath)
	}
	if loaded.Version != cfg.Version {
		t.Errorf("Expected Version %d, got %d", cfg.Version, loaded.Version)
	}
}

func TestYAMLRepository_Save_UpdateExisting(t *testing.T) {
	tmpDir := t.TempDir()

	repo := config.NewYAMLRepository(tmpDir)

	// Create initial config
	cfg := config.Config{
		ProjectName: "Game v1",
		NodesPath:   ".deco/nodes",
		HistoryPath: ".deco/history.jsonl",
		Version:     1,
	}

	err := repo.Save(cfg)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Update config
	cfg.ProjectName = "Game v2"
	cfg.Version = 2

	err = repo.Save(cfg)
	if err != nil {
		t.Fatalf("Update save failed: %v", err)
	}

	// Load and verify update
	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	if loaded.ProjectName != "Game v2" {
		t.Errorf("Expected updated ProjectName, got %q", loaded.ProjectName)
	}
	if loaded.Version != 2 {
		t.Errorf("Expected Version 2, got %d", loaded.Version)
	}
}

func TestYAMLRepository_DefaultValues(t *testing.T) {
	tmpDir := t.TempDir()

	repo := config.NewYAMLRepository(tmpDir)

	// Save minimal config
	cfg := config.Config{
		ProjectName: "Minimal Game",
		Version:     1,
	}

	err := repo.Save(cfg)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load and verify defaults
	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Empty strings should be preserved
	if loaded.NodesPath != "" {
		t.Errorf("Expected empty NodesPath, got %q", loaded.NodesPath)
	}
	if loaded.HistoryPath != "" {
		t.Errorf("Expected empty HistoryPath, got %q", loaded.HistoryPath)
	}
}

func TestYAMLRepository_EmptyCustom(t *testing.T) {
	tmpDir := t.TempDir()

	repo := config.NewYAMLRepository(tmpDir)

	// Save config without custom fields
	cfg := config.Config{
		ProjectName: "Simple Game",
		NodesPath:   ".deco/nodes",
		HistoryPath: ".deco/history.jsonl",
		Version:     1,
	}

	err := repo.Save(cfg)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load and verify
	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Custom should be nil or empty
	if loaded.Custom != nil && len(loaded.Custom) > 0 {
		t.Error("Expected empty Custom map")
	}
}

func TestConfig_RequiredApprovals(t *testing.T) {
	t.Run("loads required_approvals from config", func(t *testing.T) {
		tmpDir := t.TempDir()
		decoDir := filepath.Join(tmpDir, ".deco")
		os.MkdirAll(decoDir, 0755)

		configContent := `project_name: TestProject
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
version: 1
required_approvals: 2
`
		os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644)

		repo := config.NewYAMLRepository(tmpDir)
		cfg, err := repo.Load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		if cfg.RequiredApprovals != 2 {
			t.Errorf("Expected RequiredApprovals=2, got %d", cfg.RequiredApprovals)
		}
	})

	t.Run("defaults to 1 if not specified", func(t *testing.T) {
		tmpDir := t.TempDir()
		decoDir := filepath.Join(tmpDir, ".deco")
		os.MkdirAll(decoDir, 0755)

		configContent := `project_name: TestProject
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
version: 1
`
		os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644)

		repo := config.NewYAMLRepository(tmpDir)
		cfg, err := repo.Load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		if cfg.RequiredApprovals != 1 {
			t.Errorf("Expected RequiredApprovals=1 (default), got %d", cfg.RequiredApprovals)
		}
	})
}
