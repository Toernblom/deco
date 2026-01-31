package query_test

import (
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/query"
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

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}
