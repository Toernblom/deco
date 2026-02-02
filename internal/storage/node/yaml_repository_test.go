package node_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/node"
)

func TestYAMLRepository_LoadAll(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test node files
	createTestNode(t, nodesDir, "systems/food.yaml", domain.Node{
		ID:      "systems/food",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Food System",
	})

	createTestNode(t, nodesDir, "mechanics/combat.yaml", domain.Node{
		ID:      "mechanics/combat",
		Kind:    "mechanic",
		Version: 1,
		Status:  "approved",
		Title:   "Combat Mechanics",
	})

	// Create repository
	repo := node.NewYAMLRepository(nodesDir)

	// Load all nodes
	nodes, err := repo.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(nodes))
	}

	// Verify nodes were loaded correctly
	nodesByID := make(map[string]domain.Node)
	for _, n := range nodes {
		nodesByID[n.ID] = n
	}

	if food, exists := nodesByID["systems/food"]; exists {
		if food.Title != "Food System" {
			t.Errorf("Expected Title 'Food System', got %q", food.Title)
		}
	} else {
		t.Error("Expected systems/food node to be loaded")
	}

	if combat, exists := nodesByID["mechanics/combat"]; exists {
		if combat.Title != "Combat Mechanics" {
			t.Errorf("Expected Title 'Combat Mechanics', got %q", combat.Title)
		}
	} else {
		t.Error("Expected mechanics/combat node to be loaded")
	}
}

func TestYAMLRepository_LoadAll_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	repo := node.NewYAMLRepository(nodesDir)

	nodes, err := repo.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed on empty directory: %v", err)
	}

	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes in empty directory, got %d", len(nodes))
	}
}

func TestYAMLRepository_Load(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a test node
	createTestNode(t, nodesDir, "systems/food.yaml", domain.Node{
		ID:      "systems/food",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Food System",
		Tags:    []string{"survival", "resource"},
	})

	repo := node.NewYAMLRepository(nodesDir)

	// Load the node
	n, err := repo.Load("systems/food")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if n.ID != "systems/food" {
		t.Errorf("Expected ID 'systems/food', got %q", n.ID)
	}
	if n.Title != "Food System" {
		t.Errorf("Expected Title 'Food System', got %q", n.Title)
	}
	if len(n.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(n.Tags))
	}
}

func TestYAMLRepository_Load_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	repo := node.NewYAMLRepository(nodesDir)

	// Try to load non-existent node
	_, err = repo.Load("nonexistent/node")
	if err == nil {
		t.Error("Expected error when loading non-existent node")
	}
}

func TestYAMLRepository_Load_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create invalid YAML file
	filePath := filepath.Join(nodesDir, "invalid.yaml")
	err = os.WriteFile(filePath, []byte("invalid: yaml: content: ["), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid YAML file: %v", err)
	}

	repo := node.NewYAMLRepository(nodesDir)

	// Try to load invalid YAML
	_, err = repo.Load("invalid")
	if err == nil {
		t.Error("Expected error when loading invalid YAML")
	}
}

func TestYAMLRepository_Save(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	repo := node.NewYAMLRepository(nodesDir)

	// Create and save a node
	n := domain.Node{
		ID:      "systems/water",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Water System",
		Tags:    []string{"survival"},
	}

	err = repo.Save(n)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(nodesDir, "systems", "water.yaml")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Expected node file to be created")
	}

	// Load and verify content
	loaded, err := repo.Load("systems/water")
	if err != nil {
		t.Fatalf("Failed to load saved node: %v", err)
	}

	if loaded.ID != n.ID {
		t.Errorf("Expected ID %q, got %q", n.ID, loaded.ID)
	}
	if loaded.Title != n.Title {
		t.Errorf("Expected Title %q, got %q", n.Title, loaded.Title)
	}
}

func TestYAMLRepository_Save_UpdateExisting(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	repo := node.NewYAMLRepository(nodesDir)

	// Create initial node
	n := domain.Node{
		ID:      "systems/food",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Food System",
	}

	err = repo.Save(n)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Update the node
	n.Version = 2
	n.Status = "approved"
	n.Title = "Food System v2"

	err = repo.Save(n)
	if err != nil {
		t.Fatalf("Update save failed: %v", err)
	}

	// Load and verify update
	loaded, err := repo.Load("systems/food")
	if err != nil {
		t.Fatalf("Failed to load updated node: %v", err)
	}

	if loaded.Version != 2 {
		t.Errorf("Expected Version 2, got %d", loaded.Version)
	}
	if loaded.Status != "approved" {
		t.Errorf("Expected Status 'approved', got %q", loaded.Status)
	}
	if loaded.Title != "Food System v2" {
		t.Errorf("Expected updated title, got %q", loaded.Title)
	}
}

func TestYAMLRepository_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test node
	createTestNode(t, nodesDir, "systems/food.yaml", domain.Node{
		ID:      "systems/food",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Food System",
	})

	repo := node.NewYAMLRepository(nodesDir)

	// Verify node exists
	exists, err := repo.Exists("systems/food")
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if !exists {
		t.Error("Expected node to exist before deletion")
	}

	// Delete the node
	err = repo.Delete("systems/food")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify node is gone
	exists, err = repo.Exists("systems/food")
	if err != nil {
		t.Fatalf("Exists check failed after deletion: %v", err)
	}
	if exists {
		t.Error("Expected node to not exist after deletion")
	}
}

func TestYAMLRepository_Delete_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	repo := node.NewYAMLRepository(nodesDir)

	// Try to delete non-existent node
	err = repo.Delete("nonexistent/node")
	if err == nil {
		t.Error("Expected error when deleting non-existent node")
	}
}

func TestYAMLRepository_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test node
	createTestNode(t, nodesDir, "systems/food.yaml", domain.Node{
		ID:      "systems/food",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Food System",
	})

	repo := node.NewYAMLRepository(nodesDir)

	// Check existing node
	exists, err := repo.Exists("systems/food")
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if !exists {
		t.Error("Expected node to exist")
	}

	// Check non-existent node
	exists, err = repo.Exists("nonexistent/node")
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if exists {
		t.Error("Expected node to not exist")
	}
}

func TestYAMLRepository_NestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	repo := node.NewYAMLRepository(nodesDir)

	// Create node with deeply nested ID
	n := domain.Node{
		ID:      "systems/economy/market/trading",
		Kind:    "mechanic",
		Version: 1,
		Status:  "draft",
		Title:   "Trading Mechanics",
	}

	err = repo.Save(n)
	if err != nil {
		t.Fatalf("Save failed for nested node: %v", err)
	}

	// Verify directory structure was created
	filePath := filepath.Join(nodesDir, "systems", "economy", "market", "trading.yaml")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Expected nested directory structure to be created")
	}

	// Load and verify
	loaded, err := repo.Load("systems/economy/market/trading")
	if err != nil {
		t.Fatalf("Failed to load nested node: %v", err)
	}

	if loaded.ID != n.ID {
		t.Errorf("Expected ID %q, got %q", n.ID, loaded.ID)
	}
}

// Helper function to create a test node file
func createTestNode(t *testing.T, baseDir, relativePath string, n domain.Node) {
	t.Helper()

	filePath := filepath.Join(baseDir, relativePath)
	dir := filepath.Dir(filePath)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	// Manually create YAML content to avoid yaml import in test
	content := ""
	content += "id: " + n.ID + "\n"
	content += "kind: " + n.Kind + "\n"
	content += "version: " + string(rune(n.Version+'0')) + "\n"
	content += "status: " + n.Status + "\n"
	content += "title: " + n.Title + "\n"

	if len(n.Tags) > 0 {
		content += "tags:\n"
		for _, tag := range n.Tags {
			content += "  - " + tag + "\n"
		}
	}

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test node file %s: %v", filePath, err)
	}
}
