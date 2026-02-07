package query_test

import (
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/query"
	"github.com/Toernblom/deco/internal/storage/config"
)

// ===== FILTER OPERATION TESTS =====

// Test filtering by kind
func TestQueryEngine_FilterByKind(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "System Node"},
		{ID: "node2", Kind: "feature", Version: 1, Status: "draft", Title: "Feature Node"},
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Another System"},
	}

	criteria := query.FilterCriteria{
		Kind: strPtr("system"),
	}

	results := qe.Filter(nodes, criteria)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	for _, node := range results {
		if node.Kind != "system" {
			t.Errorf("expected kind 'system', got %q", node.Kind)
		}
	}
}

// Test filtering by status
func TestQueryEngine_FilterByStatus(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Draft Node"},
		{ID: "node2", Kind: "system", Version: 1, Status: "published", Title: "Published Node"},
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Another Draft"},
	}

	criteria := query.FilterCriteria{
		Status: strPtr("published"),
	}

	results := qe.Filter(nodes, criteria)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Status != "published" {
		t.Errorf("expected status 'published', got %q", results[0].Status)
	}
}

// Test filtering by single tag
func TestQueryEngine_FilterBySingleTag(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1", Tags: []string{"gameplay", "combat"}},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Node 2", Tags: []string{"ui", "menu"}},
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Node 3", Tags: []string{"gameplay", "inventory"}},
	}

	criteria := query.FilterCriteria{
		Tags: []string{"gameplay"},
	}

	results := qe.Filter(nodes, criteria)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	for _, node := range results {
		found := false
		for _, tag := range node.Tags {
			if tag == "gameplay" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected node %s to have tag 'gameplay'", node.ID)
		}
	}
}

// Test filtering by multiple tags (should match nodes with ALL tags)
func TestQueryEngine_FilterByMultipleTags(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1", Tags: []string{"gameplay", "combat", "pvp"}},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Node 2", Tags: []string{"gameplay", "combat"}},
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Node 3", Tags: []string{"gameplay"}},
	}

	criteria := query.FilterCriteria{
		Tags: []string{"gameplay", "combat"},
	}

	results := qe.Filter(nodes, criteria)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	for _, node := range results {
		hasGameplay := false
		hasCombat := false
		for _, tag := range node.Tags {
			if tag == "gameplay" {
				hasGameplay = true
			}
			if tag == "combat" {
				hasCombat = true
			}
		}
		if !hasGameplay || !hasCombat {
			t.Errorf("expected node %s to have both 'gameplay' and 'combat' tags", node.ID)
		}
	}
}

// Test combined filters (kind + status)
func TestQueryEngine_FilterCombinedKindAndStatus(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "System Draft"},
		{ID: "node2", Kind: "feature", Version: 1, Status: "published", Title: "Feature Published"},
		{ID: "node3", Kind: "system", Version: 1, Status: "published", Title: "System Published"},
		{ID: "node4", Kind: "system", Version: 1, Status: "draft", Title: "System Draft 2"},
	}

	criteria := query.FilterCriteria{
		Kind:   strPtr("system"),
		Status: strPtr("published"),
	}

	results := qe.Filter(nodes, criteria)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].ID != "node3" {
		t.Errorf("expected node3, got %s", results[0].ID)
	}
}

// Test combined filters (all criteria)
func TestQueryEngine_FilterCombinedAll(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1", Tags: []string{"gameplay", "combat"}},
		{ID: "node2", Kind: "feature", Version: 1, Status: "published", Title: "Node 2", Tags: []string{"gameplay", "combat"}},
		{ID: "node3", Kind: "system", Version: 1, Status: "published", Title: "Node 3", Tags: []string{"gameplay", "combat"}},
		{ID: "node4", Kind: "system", Version: 1, Status: "published", Title: "Node 4", Tags: []string{"gameplay"}},
	}

	criteria := query.FilterCriteria{
		Kind:   strPtr("system"),
		Status: strPtr("published"),
		Tags:   []string{"gameplay", "combat"},
	}

	results := qe.Filter(nodes, criteria)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].ID != "node3" {
		t.Errorf("expected node3, got %s", results[0].ID)
	}
}

// Test filter with no matches
func TestQueryEngine_FilterNoMatches(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "node2", Kind: "feature", Version: 1, Status: "draft", Title: "Node 2"},
	}

	criteria := query.FilterCriteria{
		Status: strPtr("published"),
	}

	results := qe.Filter(nodes, criteria)

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

// Test filter with empty criteria (should return all nodes)
func TestQueryEngine_FilterEmptyCriteria(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "node2", Kind: "feature", Version: 1, Status: "published", Title: "Node 2"},
	}

	criteria := query.FilterCriteria{}

	results := qe.Filter(nodes, criteria)

	if len(results) != 2 {
		t.Errorf("expected all 2 nodes, got %d", len(results))
	}
}

// Test filter with empty node list
func TestQueryEngine_FilterEmptyNodes(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{}

	criteria := query.FilterCriteria{
		Kind: strPtr("system"),
	}

	results := qe.Filter(nodes, criteria)

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

// Test filter with nodes that have no tags
func TestQueryEngine_FilterNodesWithoutTags(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1", Tags: []string{"gameplay"}},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Node 2", Tags: nil},
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Node 3", Tags: []string{}},
	}

	criteria := query.FilterCriteria{
		Tags: []string{"gameplay"},
	}

	results := qe.Filter(nodes, criteria)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].ID != "node1" {
		t.Errorf("expected node1, got %s", results[0].ID)
	}
}

// ===== SEARCH OPERATION TESTS =====

// Test searching in titles
func TestQueryEngine_SearchInTitle(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Combat System"},
		{ID: "node2", Kind: "feature", Version: 1, Status: "draft", Title: "Inventory System"},
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Player Movement"},
	}

	results := qe.Search(nodes, "System")

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	for _, node := range results {
		if node.ID != "node1" && node.ID != "node2" {
			t.Errorf("unexpected node %s in results", node.ID)
		}
	}
}

// Test searching in summary/content
func TestQueryEngine_SearchInSummary(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1", Summary: "Handles player combat"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Node 2", Summary: "Manages inventory items"},
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Node 3", Summary: "Player combat mechanics"},
	}

	results := qe.Search(nodes, "combat")

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

// Test case-insensitive search
func TestQueryEngine_SearchCaseInsensitive(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "COMBAT System"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "combat mechanics"},
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Combat Features"},
	}

	// Search with different cases
	results1 := qe.Search(nodes, "COMBAT")
	results2 := qe.Search(nodes, "combat")
	results3 := qe.Search(nodes, "CoMbAt")

	if len(results1) != 3 || len(results2) != 3 || len(results3) != 3 {
		t.Errorf("case-insensitive search failed: got %d, %d, %d results", len(results1), len(results2), len(results3))
	}
}

// Test partial matches
func TestQueryEngine_SearchPartialMatch(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Combat System"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Combative Behavior"},
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Player Movement"},
	}

	results := qe.Search(nodes, "Comba")

	if len(results) != 2 {
		t.Fatalf("expected 2 results for partial match, got %d", len(results))
	}
}

// Test search with no matches
func TestQueryEngine_SearchNoMatches(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Combat System"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Inventory System"},
	}

	results := qe.Search(nodes, "NonexistentTerm")

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

// Test search with empty term
func TestQueryEngine_SearchEmptyTerm(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Node 2"},
	}

	results := qe.Search(nodes, "")

	// Empty search should return all nodes (everything matches empty string)
	if len(results) != 2 {
		t.Errorf("expected all nodes for empty search, got %d", len(results))
	}
}

// Test search in both title and summary
func TestQueryEngine_SearchInTitleAndSummary(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Combat System", Summary: "Handles combat"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Inventory", Summary: "Contains combat loot"},
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Player Movement", Summary: "Movement mechanics"},
	}

	results := qe.Search(nodes, "combat")

	// Should find node1 (title + summary) and node2 (summary only)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

// Test search with empty node list
func TestQueryEngine_SearchEmptyNodes(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{}

	results := qe.Search(nodes, "test")

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

// ===== BLOCK-LEVEL QUERY TESTS =====

func TestQueryEngine_FilterBlocks_ByBlockType(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{
			ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{"name": "Smithy", "age": "bronze"}},
							{Type: "building", Data: map[string]interface{}{"name": "Barracks", "age": "iron"}},
							{Type: "resource", Data: map[string]interface{}{"name": "Iron Ore", "tier": 2}},
						},
					},
				},
			},
		},
		{
			ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Combat",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Gear",
						Blocks: []domain.Block{
							{Type: "gear", Data: map[string]interface{}{"name": "Bronze Sword"}},
						},
					},
				},
			},
		},
	}

	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType: &blockType,
	}
	results := qe.FilterBlocks(nodes, criteria)

	if len(results) != 2 {
		t.Fatalf("expected 2 building blocks, got %d", len(results))
	}
	for _, r := range results {
		if r.Block.Type != "building" {
			t.Errorf("expected block type 'building', got %q", r.Block.Type)
		}
	}
}

func TestQueryEngine_FilterBlocks_ByBlockTypeWithFieldFilter(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{
			ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{"name": "Smithy", "age": "bronze"}},
							{Type: "building", Data: map[string]interface{}{"name": "Barracks", "age": "iron"}},
							{Type: "building", Data: map[string]interface{}{"name": "Farm", "age": "bronze"}},
						},
					},
				},
			},
		},
	}

	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType:    &blockType,
		FieldFilters: map[string]string{"age": "bronze"},
	}
	results := qe.FilterBlocks(nodes, criteria)

	if len(results) != 2 {
		t.Fatalf("expected 2 bronze buildings, got %d", len(results))
	}
	for _, r := range results {
		if r.Block.Data["age"] != "bronze" {
			t.Errorf("expected age 'bronze', got %v", r.Block.Data["age"])
		}
	}
}

func TestQueryEngine_FilterBlocks_MultipleFieldFilters(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{
			ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{"name": "Smithy", "age": "bronze", "category": "production"}},
							{Type: "building", Data: map[string]interface{}{"name": "Barracks", "age": "bronze", "category": "military"}},
							{Type: "building", Data: map[string]interface{}{"name": "Farm", "age": "iron", "category": "production"}},
						},
					},
				},
			},
		},
	}

	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType:    &blockType,
		FieldFilters: map[string]string{"age": "bronze", "category": "production"},
	}
	results := qe.FilterBlocks(nodes, criteria)

	if len(results) != 1 {
		t.Fatalf("expected 1 bronze production building, got %d", len(results))
	}
	if results[0].Block.Data["name"] != "Smithy" {
		t.Errorf("expected Smithy, got %v", results[0].Block.Data["name"])
	}
}

func TestQueryEngine_FilterBlocks_NoMatches(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{
			ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{"name": "Smithy", "age": "bronze"}},
						},
					},
				},
			},
		},
	}

	blockType := "vehicle"
	criteria := query.FilterCriteria{
		BlockType: &blockType,
	}
	results := qe.FilterBlocks(nodes, criteria)

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestQueryEngine_FilterBlocks_CombinedNodeAndBlockFilters(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{
			ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{"name": "Smithy", "age": "bronze"}},
						},
					},
				},
			},
		},
		{
			ID: "node2", Kind: "feature", Version: 1, Status: "draft", Title: "More Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{"name": "Castle", "age": "medieval"}},
						},
					},
				},
			},
		},
	}

	blockType := "building"
	kind := "system"
	criteria := query.FilterCriteria{
		Kind:      &kind,
		BlockType: &blockType,
	}
	results := qe.FilterBlocks(nodes, criteria)

	// Should only find buildings in "system" nodes
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].NodeID != "node1" {
		t.Errorf("expected node1, got %s", results[0].NodeID)
	}
}

func TestQueryEngine_FilterBlocks_ResultContainsContext(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{
			ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Bronze Age",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{"name": "Smithy"}},
						},
					},
				},
			},
		},
	}

	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType: &blockType,
	}
	results := qe.FilterBlocks(nodes, criteria)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.NodeID != "node1" {
		t.Errorf("expected NodeID 'node1', got %q", r.NodeID)
	}
	if r.NodeTitle != "Buildings" {
		t.Errorf("expected NodeTitle 'Buildings', got %q", r.NodeTitle)
	}
	if r.SectionName != "Bronze Age" {
		t.Errorf("expected SectionName 'Bronze Age', got %q", r.SectionName)
	}
	if r.BlockIndex != 0 {
		t.Errorf("expected BlockIndex 0, got %d", r.BlockIndex)
	}
}

func TestQueryEngine_FilterBlocks_FieldFilterStringCoercion(t *testing.T) {
	// Field filters should match numeric values via string comparison
	qe := query.New()

	nodes := []domain.Node{
		{
			ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Resources",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Materials",
						Blocks: []domain.Block{
							{Type: "resource", Data: map[string]interface{}{"name": "Iron", "tier": 3}},
							{Type: "resource", Data: map[string]interface{}{"name": "Stone", "tier": 1}},
						},
					},
				},
			},
		},
	}

	blockType := "resource"
	criteria := query.FilterCriteria{
		BlockType:    &blockType,
		FieldFilters: map[string]string{"tier": "3"},
	}
	results := qe.FilterBlocks(nodes, criteria)

	if len(results) != 1 {
		t.Fatalf("expected 1 result for tier=3, got %d", len(results))
	}
	if results[0].Block.Data["name"] != "Iron" {
		t.Errorf("expected Iron, got %v", results[0].Block.Data["name"])
	}
}

func TestQueryEngine_FilterBlocks_NodesWithoutContent(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Empty Node", Content: nil},
		{
			ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Has Blocks",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Stuff",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{"name": "Smithy"}},
						},
					},
				},
			},
		},
	}

	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType: &blockType,
	}
	results := qe.FilterBlocks(nodes, criteria)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

// ===== LIST MEMBERSHIP QUERY TESTS =====

func TestQueryEngine_FilterBlocks_ListFieldMembership(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{
			ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{
								"name":      "Smithy",
								"materials": []interface{}{"Stone", "Planks", "Bronze Ingots"},
							}},
							{Type: "building", Data: map[string]interface{}{
								"name":      "Farm",
								"materials": []interface{}{"Wood", "Stone"},
							}},
							{Type: "building", Data: map[string]interface{}{
								"name":      "Barracks",
								"materials": []interface{}{"Iron", "Stone"},
							}},
						},
					},
				},
			},
		},
	}

	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType:    &blockType,
		FieldFilters: map[string]string{"materials": "Planks"},
	}
	results := qe.FilterBlocks(nodes, criteria)

	if len(results) != 1 {
		t.Fatalf("expected 1 building with Planks, got %d", len(results))
	}
	if results[0].Block.Data["name"] != "Smithy" {
		t.Errorf("expected Smithy, got %v", results[0].Block.Data["name"])
	}
}

func TestQueryEngine_FilterBlocks_ListFieldMembershipMultipleMatches(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{
			ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{
								"name":      "Smithy",
								"materials": []interface{}{"Stone", "Planks"},
							}},
							{Type: "building", Data: map[string]interface{}{
								"name":      "Farm",
								"materials": []interface{}{"Wood", "Stone"},
							}},
							{Type: "building", Data: map[string]interface{}{
								"name":      "Quarry",
								"materials": []interface{}{"Stone"},
							}},
						},
					},
				},
			},
		},
	}

	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType:    &blockType,
		FieldFilters: map[string]string{"materials": "Stone"},
	}
	results := qe.FilterBlocks(nodes, criteria)

	if len(results) != 3 {
		t.Fatalf("expected 3 buildings with Stone, got %d", len(results))
	}
}

func TestQueryEngine_FilterBlocks_ListFieldNoMatch(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{
			ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{
								"name":      "Smithy",
								"materials": []interface{}{"Stone", "Planks"},
							}},
						},
					},
				},
			},
		},
	}

	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType:    &blockType,
		FieldFilters: map[string]string{"materials": "Gold"},
	}
	results := qe.FilterBlocks(nodes, criteria)

	if len(results) != 0 {
		t.Errorf("expected 0 results for Gold, got %d", len(results))
	}
}

func TestQueryEngine_FilterBlocks_ListFieldCombinedWithScalarFilter(t *testing.T) {
	qe := query.New()

	nodes := []domain.Node{
		{
			ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{
								"name":      "Smithy",
								"age":       "bronze",
								"materials": []interface{}{"Stone", "Planks"},
							}},
							{Type: "building", Data: map[string]interface{}{
								"name":      "Farm",
								"age":       "stone",
								"materials": []interface{}{"Stone", "Wood"},
							}},
						},
					},
				},
			},
		},
	}

	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType:    &blockType,
		FieldFilters: map[string]string{"materials": "Stone", "age": "bronze"},
	}
	results := qe.FilterBlocks(nodes, criteria)

	if len(results) != 1 {
		t.Fatalf("expected 1 bronze building with Stone, got %d", len(results))
	}
	if results[0].Block.Data["name"] != "Smithy" {
		t.Errorf("expected Smithy, got %v", results[0].Block.Data["name"])
	}
}

// ===== FOLLOW QUERY TESTS =====

func buildFollowTestData() ([]domain.Node, map[string]config.BlockTypeConfig) {
	nodes := []domain.Node{
		{
			ID: "buildings", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Bronze Age",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{
								"name":      "Smithy",
								"age":       "bronze",
								"materials": []interface{}{"Stone", "Planks", "Bronze Ingots"},
							}},
							{Type: "building", Data: map[string]interface{}{
								"name":      "Carpenter",
								"age":       "bronze",
								"materials": []interface{}{"Wood", "Planks"},
							}},
						},
					},
					{
						Name: "Stone Age",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{
								"name":      "Hut",
								"age":       "stone",
								"materials": []interface{}{"Stone", "Wood"},
							}},
						},
					},
				},
			},
		},
		{
			ID: "resources", Kind: "system", Version: 1, Status: "draft", Title: "Resources",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Raw",
						Blocks: []domain.Block{
							{Type: "resource", Data: map[string]interface{}{"name": "Stone", "tier": 0}},
							{Type: "resource", Data: map[string]interface{}{"name": "Wood", "tier": 0}},
						},
					},
				},
			},
		},
		{
			ID: "recipes", Kind: "system", Version: 1, Status: "draft", Title: "Recipes",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Crafting",
						Blocks: []domain.Block{
							{Type: "recipe", Data: map[string]interface{}{
								"name":   "Plank Making",
								"output": "Planks",
								"inputs": []interface{}{"Wood"},
							}},
							{Type: "recipe", Data: map[string]interface{}{
								"name":   "Bronze Smelting",
								"output": "Bronze Ingots",
								"inputs": []interface{}{"Copper", "Tin"},
							}},
						},
					},
				},
			},
		},
	}

	blockTypes := map[string]config.BlockTypeConfig{
		"building": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
				"age":  {Type: "string", Required: true},
				"materials": {Type: "list", Refs: []config.RefConstraint{
					{BlockType: "resource", Field: "name"},
					{BlockType: "recipe", Field: "output"},
				}},
			},
		},
		"resource": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
				"tier": {Type: "number"},
			},
		},
		"recipe": {
			Fields: map[string]config.FieldDef{
				"name":   {Type: "string", Required: true},
				"output": {Type: "string", Required: true},
				"inputs": {Type: "list", Refs: []config.RefConstraint{
					{BlockType: "resource", Field: "name"},
				}},
			},
		},
	}

	return nodes, blockTypes
}

func TestQueryEngine_FollowBlocks_AutoFromRefConfig(t *testing.T) {
	qe := query.New()
	nodes, blockTypes := buildFollowTestData()

	// Query bronze-age buildings, follow materials
	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType:    &blockType,
		FieldFilters: map[string]string{"age": "bronze"},
	}
	sourceMatches := qe.FilterBlocks(nodes, criteria)

	results, err := qe.FollowBlocks(sourceMatches, "materials", nil, nodes, blockTypes)
	if err != nil {
		t.Fatal(err)
	}

	// Should find: Bronze Ingots, Planks, Stone, Wood
	if len(results) != 4 {
		t.Fatalf("expected 4 followed values, got %d", len(results))
	}

	// Results are sorted alphabetically
	expectedValues := []string{"Bronze Ingots", "Planks", "Stone", "Wood"}
	for i, expected := range expectedValues {
		if results[i].Value != expected {
			t.Errorf("result[%d]: expected value %q, got %q", i, expected, results[i].Value)
		}
	}

	// Bronze Ingots: referenced by 1 building (Smithy), matched by recipe
	if results[0].RefCount != 1 {
		t.Errorf("Bronze Ingots: expected refcount 1, got %d", results[0].RefCount)
	}
	if len(results[0].Matches) != 1 {
		t.Errorf("Bronze Ingots: expected 1 match (recipe), got %d", len(results[0].Matches))
	}

	// Planks: referenced by 2 buildings (Smithy, Carpenter), matched by recipe
	if results[1].RefCount != 2 {
		t.Errorf("Planks: expected refcount 2, got %d", results[1].RefCount)
	}
	if len(results[1].Matches) != 1 {
		t.Errorf("Planks: expected 1 match (recipe), got %d", len(results[1].Matches))
	}

	// Stone: referenced by 1 building (Smithy), matched by resource
	if results[2].RefCount != 1 {
		t.Errorf("Stone: expected refcount 1, got %d", results[2].RefCount)
	}
	if len(results[2].Matches) != 1 {
		t.Errorf("Stone: expected 1 match (resource), got %d", len(results[2].Matches))
	}
}

func TestQueryEngine_FollowBlocks_ExplicitTarget(t *testing.T) {
	qe := query.New()
	nodes, blockTypes := buildFollowTestData()

	// Query all buildings, follow materials with explicit target recipe.output only
	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType: &blockType,
	}
	sourceMatches := qe.FilterBlocks(nodes, criteria)

	targets := []query.FollowTarget{{BlockType: "recipe", Field: "output"}}
	results, err := qe.FollowBlocks(sourceMatches, "materials", targets, nodes, blockTypes)
	if err != nil {
		t.Fatal(err)
	}

	// Should find values but only match against recipes, not resources
	// Stone and Wood should have 0 matches (they're resources, not recipe outputs)
	for _, r := range results {
		if r.Value == "Stone" || r.Value == "Wood" {
			if len(r.Matches) != 0 {
				t.Errorf("%s: expected 0 matches when following only recipe.output, got %d", r.Value, len(r.Matches))
			}
		}
		if r.Value == "Planks" {
			if len(r.Matches) != 1 {
				t.Errorf("Planks: expected 1 match, got %d", len(r.Matches))
			}
		}
	}
}

func TestQueryEngine_FollowBlocks_ReverseDirection(t *testing.T) {
	qe := query.New()
	nodes, blockTypes := buildFollowTestData()

	// Query recipes where output=Planks, follow inputs
	blockType := "recipe"
	criteria := query.FilterCriteria{
		BlockType:    &blockType,
		FieldFilters: map[string]string{"output": "Planks"},
	}
	sourceMatches := qe.FilterBlocks(nodes, criteria)

	if len(sourceMatches) != 1 {
		t.Fatalf("expected 1 Planks recipe, got %d", len(sourceMatches))
	}

	results, err := qe.FollowBlocks(sourceMatches, "inputs", nil, nodes, blockTypes)
	if err != nil {
		t.Fatal(err)
	}

	// Planks recipe inputs: [Wood] -> should find Wood resource
	if len(results) != 1 {
		t.Fatalf("expected 1 followed value (Wood), got %d", len(results))
	}
	if results[0].Value != "Wood" {
		t.Errorf("expected value 'Wood', got %q", results[0].Value)
	}
	if len(results[0].Matches) != 1 {
		t.Errorf("expected 1 match for Wood, got %d", len(results[0].Matches))
	}
	if results[0].Matches[0].Block.Type != "resource" {
		t.Errorf("expected resource block, got %s", results[0].Matches[0].Block.Type)
	}
}

func TestQueryEngine_FollowBlocks_FieldNotFound(t *testing.T) {
	qe := query.New()
	nodes, blockTypes := buildFollowTestData()

	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType: &blockType,
	}
	sourceMatches := qe.FilterBlocks(nodes, criteria)

	_, err := qe.FollowBlocks(sourceMatches, "nonexistent", nil, nodes, blockTypes)
	if err == nil {
		t.Fatal("expected error for nonexistent field")
	}
}

func TestQueryEngine_FollowBlocks_NoRefConfig(t *testing.T) {
	qe := query.New()
	nodes, _ := buildFollowTestData()

	// Use block types without ref config
	noRefTypes := map[string]config.BlockTypeConfig{
		"building": {
			Fields: map[string]config.FieldDef{
				"name":      {Type: "string", Required: true},
				"materials": {Type: "list"}, // No refs
			},
		},
	}

	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType: &blockType,
	}
	sourceMatches := qe.FilterBlocks(nodes, criteria)

	_, err := qe.FollowBlocks(sourceMatches, "materials", nil, nodes, noRefTypes)
	if err == nil {
		t.Fatal("expected error when field has no ref config")
	}
	if !strings.Contains(err.Error(), "no ref constraint") {
		t.Errorf("expected 'no ref constraint' in error, got: %s", err.Error())
	}
}

func TestQueryEngine_FollowBlocks_ValueNotFound(t *testing.T) {
	qe := query.New()

	// Building uses "Diamond" which no recipe/resource produces
	nodes := []domain.Node{
		{
			ID: "buildings", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Special",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{
								"name":      "Palace",
								"materials": []interface{}{"Diamond"},
							}},
						},
					},
				},
			},
		},
	}

	blockTypes := map[string]config.BlockTypeConfig{
		"building": {
			Fields: map[string]config.FieldDef{
				"name":      {Type: "string", Required: true},
				"materials": {Type: "list", Refs: []config.RefConstraint{{BlockType: "resource", Field: "name"}}},
			},
		},
	}

	blockType := "building"
	criteria := query.FilterCriteria{BlockType: &blockType}
	sourceMatches := qe.FilterBlocks(nodes, criteria)

	results, err := qe.FollowBlocks(sourceMatches, "materials", nil, nodes, blockTypes)
	if err != nil {
		t.Fatal(err)
	}

	// Should still return the value, just with 0 matches
	if len(results) != 1 {
		t.Fatalf("expected 1 result (Diamond), got %d", len(results))
	}
	if results[0].Value != "Diamond" {
		t.Errorf("expected value 'Diamond', got %q", results[0].Value)
	}
	if len(results[0].Matches) != 0 {
		t.Errorf("expected 0 matches for Diamond, got %d", len(results[0].Matches))
	}
}

func TestQueryEngine_FollowBlocks_DeduplicatesGrouping(t *testing.T) {
	qe := query.New()
	nodes, blockTypes := buildFollowTestData()

	// Query all buildings (not just bronze) — Stone is used by Smithy, Hut
	blockType := "building"
	criteria := query.FilterCriteria{
		BlockType: &blockType,
	}
	sourceMatches := qe.FilterBlocks(nodes, criteria)

	results, err := qe.FollowBlocks(sourceMatches, "materials", nil, nodes, blockTypes)
	if err != nil {
		t.Fatal(err)
	}

	// Find Stone in results
	var stoneResult *query.FollowResult
	for i := range results {
		if results[i].Value == "Stone" {
			stoneResult = &results[i]
			break
		}
	}

	if stoneResult == nil {
		t.Fatal("expected Stone in results")
	}

	// Stone is used by Smithy and Hut = 2 source refs
	if stoneResult.RefCount != 2 {
		t.Errorf("Stone: expected refcount 2, got %d", stoneResult.RefCount)
	}

	// But only 1 resource block matches (deduplicated)
	if len(stoneResult.Matches) != 1 {
		t.Errorf("Stone: expected 1 deduplicated match, got %d", len(stoneResult.Matches))
	}
}

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}
