package node_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/Toernblom/deco/internal/storage/node"
)

func TestFileDiscovery_DiscoverAll(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")

	// Create test files
	createFile(t, nodesDir, "systems/food.yaml")
	createFile(t, nodesDir, "mechanics/combat.yaml")
	createFile(t, nodesDir, "features/quest.yaml")

	discovery := node.NewFileDiscovery(tmpDir)

	// Discover all files
	files, err := discovery.DiscoverAll()
	if err != nil {
		t.Fatalf("DiscoverAll failed: %v", err)
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}

	// Sort for consistent comparison
	sort.Strings(files)

	expectedFiles := []string{
		filepath.Join(nodesDir, "features", "quest.yaml"),
		filepath.Join(nodesDir, "mechanics", "combat.yaml"),
		filepath.Join(nodesDir, "systems", "food.yaml"),
	}
	sort.Strings(expectedFiles)

	for i, expected := range expectedFiles {
		if files[i] != expected {
			t.Errorf("Expected file %q, got %q", expected, files[i])
		}
	}
}

func TestFileDiscovery_DiscoverAll_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	err := os.MkdirAll(nodesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	discovery := node.NewFileDiscovery(tmpDir)

	files, err := discovery.DiscoverAll()
	if err != nil {
		t.Fatalf("DiscoverAll failed on empty directory: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files in empty directory, got %d", len(files))
	}
}

func TestFileDiscovery_DiscoverAll_NoNodesDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	discovery := node.NewFileDiscovery(tmpDir)

	files, err := discovery.DiscoverAll()
	if err != nil {
		t.Fatalf("DiscoverAll failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files when nodes directory doesn't exist, got %d", len(files))
	}
}

func TestFileDiscovery_IgnoreNonYAMLFiles(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")

	// Create mix of YAML and non-YAML files
	createFile(t, nodesDir, "systems/food.yaml")
	createFile(t, nodesDir, "systems/README.md")
	createFile(t, nodesDir, "mechanics/combat.yaml")
	createFile(t, nodesDir, "mechanics/notes.txt")
	createFile(t, nodesDir, ".gitignore")

	discovery := node.NewFileDiscovery(tmpDir)

	files, err := discovery.DiscoverAll()
	if err != nil {
		t.Fatalf("DiscoverAll failed: %v", err)
	}

	// Should only find .yaml files
	if len(files) != 2 {
		t.Errorf("Expected 2 YAML files, got %d", len(files))
	}

	// Verify only .yaml files were found
	for _, file := range files {
		if filepath.Ext(file) != ".yaml" {
			t.Errorf("Non-YAML file found: %s", file)
		}
	}
}

func TestFileDiscovery_NestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")

	// Create deeply nested files
	createFile(t, nodesDir, "systems/economy/market/trading.yaml")
	createFile(t, nodesDir, "mechanics/combat/weapons/sword.yaml")
	createFile(t, nodesDir, "a/b/c/d/e/deep.yaml")

	discovery := node.NewFileDiscovery(tmpDir)

	files, err := discovery.DiscoverAll()
	if err != nil {
		t.Fatalf("DiscoverAll failed: %v", err)
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}
}

func TestFileDiscovery_PathToID(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")

	discovery := node.NewFileDiscovery(tmpDir)

	tests := []struct {
		name       string
		filePath   string
		expectedID string
	}{
		{
			name:       "simple file",
			filePath:   filepath.Join(nodesDir, "food.yaml"),
			expectedID: "food",
		},
		{
			name:       "one level deep",
			filePath:   filepath.Join(nodesDir, "systems", "food.yaml"),
			expectedID: "systems/food",
		},
		{
			name:       "two levels deep",
			filePath:   filepath.Join(nodesDir, "systems", "economy", "market.yaml"),
			expectedID: "systems/economy/market",
		},
		{
			name:       "three levels deep",
			filePath:   filepath.Join(nodesDir, "a", "b", "c", "d.yaml"),
			expectedID: "a/b/c/d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := discovery.PathToID(tt.filePath)
			if id != tt.expectedID {
				t.Errorf("PathToID(%q) = %q, expected %q", tt.filePath, id, tt.expectedID)
			}
		})
	}
}

func TestFileDiscovery_IDToPath(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")

	discovery := node.NewFileDiscovery(tmpDir)

	tests := []struct {
		name         string
		id           string
		expectedPath string
	}{
		{
			name:         "simple ID",
			id:           "food",
			expectedPath: filepath.Join(nodesDir, "food.yaml"),
		},
		{
			name:         "one level deep",
			id:           "systems/food",
			expectedPath: filepath.Join(nodesDir, "systems", "food.yaml"),
		},
		{
			name:         "two levels deep",
			id:           "systems/economy/market",
			expectedPath: filepath.Join(nodesDir, "systems", "economy", "market.yaml"),
		},
		{
			name:         "three levels deep",
			id:           "a/b/c/d",
			expectedPath: filepath.Join(nodesDir, "a", "b", "c", "d.yaml"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := discovery.IDToPath(tt.id)
			if path != tt.expectedPath {
				t.Errorf("IDToPath(%q) = %q, expected %q", tt.id, path, tt.expectedPath)
			}
		})
	}
}

func TestFileDiscovery_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	discovery := node.NewFileDiscovery(tmpDir)

	// Test that PathToID and IDToPath are inverses
	ids := []string{
		"food",
		"systems/food",
		"mechanics/combat",
		"systems/economy/market/trading",
	}

	for _, id := range ids {
		path := discovery.IDToPath(id)
		roundTripID := discovery.PathToID(path)

		if roundTripID != id {
			t.Errorf("Round trip failed: %q -> %q -> %q", id, path, roundTripID)
		}
	}
}

func TestFileDiscovery_DiscoverByPattern(t *testing.T) {
	tmpDir := t.TempDir()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")

	// Create files
	createFile(t, nodesDir, "systems/food.yaml")
	createFile(t, nodesDir, "systems/water.yaml")
	createFile(t, nodesDir, "mechanics/combat.yaml")
	createFile(t, nodesDir, "features/quest.yaml")

	discovery := node.NewFileDiscovery(tmpDir)

	// Discover files matching pattern
	files, err := discovery.DiscoverByPattern("systems")
	if err != nil {
		t.Fatalf("DiscoverByPattern failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files in systems/, got %d", len(files))
	}
}

// Helper function to create a file
func createFile(t *testing.T, baseDir, relativePath string) {
	t.Helper()

	filePath := filepath.Join(baseDir, relativePath)
	dir := filepath.Dir(filePath)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	err = os.WriteFile(filePath, []byte("# test file"), 0644)
	if err != nil {
		t.Fatalf("Failed to write file %s: %v", filePath, err)
	}
}
