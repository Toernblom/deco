package refactor_test

import (
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/refactor"
)

// ===== BASIC RENAME TESTS =====

// Test renaming a single node with no references
func TestRenamer_BasicRename(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "old-id", Kind: "system", Version: 1, Status: "draft", Title: "Test Node"},
	}

	result, err := r.Rename(nodes, "old-id", "new-id")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check old ID no longer exists
	for _, node := range result {
		if node.ID == "old-id" {
			t.Error("expected old ID to be removed from result")
		}
	}

	// Check new ID exists
	var found bool
	for _, node := range result {
		if node.ID == "new-id" {
			found = true
			if node.Title != "Test Node" {
				t.Errorf("expected title preserved, got %q", node.Title)
			}
			if node.Kind != "system" {
				t.Errorf("expected kind preserved, got %q", node.Kind)
			}
			if node.Status != "draft" {
				t.Errorf("expected status preserved, got %q", node.Status)
			}
		}
	}

	if !found {
		t.Fatal("expected node with new ID to exist")
	}
}

// Test renaming preserves version (renaming the node itself doesn't change version)
func TestRenamer_PreservesVersion(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "old-id", Kind: "system", Version: 5, Status: "draft", Title: "Test"},
	}

	result, err := r.Rename(nodes, "old-id", "new-id")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "new-id" {
			if node.Version != 5 {
				t.Errorf("expected version 5 preserved, got %d", node.Version)
			}
		}
	}
}

// Test renaming preserves all node properties
func TestRenamer_PreservesAllProperties(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{
			ID:      "old-id",
			Kind:    "mechanic",
			Version: 3,
			Status:  "approved",
			Title:   "Combat System",
			Tags:    []string{"combat", "core"},
			Summary: "Main combat mechanics",
		},
	}

	result, err := r.Rename(nodes, "old-id", "new-id")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "new-id" {
			if node.Kind != "mechanic" {
				t.Errorf("expected kind 'mechanic', got %q", node.Kind)
			}
			if node.Version != 3 {
				t.Errorf("expected version 3, got %d", node.Version)
			}
			if node.Status != "approved" {
				t.Errorf("expected status 'approved', got %q", node.Status)
			}
			if node.Title != "Combat System" {
				t.Errorf("expected title 'Combat System', got %q", node.Title)
			}
			if len(node.Tags) != 2 || node.Tags[0] != "combat" || node.Tags[1] != "core" {
				t.Errorf("expected tags [combat, core], got %v", node.Tags)
			}
			if node.Summary != "Main combat mechanics" {
				t.Errorf("expected summary preserved, got %q", node.Summary)
			}
		}
	}
}

// ===== REFERENCE UPDATE TESTS =====

// Test renaming updates Uses references
func TestRenamer_UpdatesUsesReferences(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{ID: "referrer", Kind: "mechanic", Version: 1, Status: "draft", Title: "Referrer",
			Refs: domain.Ref{
				Uses: []domain.RefLink{{Target: "target"}},
			}},
	}

	result, err := r.Rename(nodes, "target", "renamed-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check referrer's Uses now points to renamed-target
	for _, node := range result {
		if node.ID == "referrer" {
			if len(node.Refs.Uses) != 1 {
				t.Fatalf("expected 1 Uses ref, got %d", len(node.Refs.Uses))
			}
			if node.Refs.Uses[0].Target != "renamed-target" {
				t.Errorf("expected Uses target 'renamed-target', got %q", node.Refs.Uses[0].Target)
			}
		}
	}
}

// Test renaming updates Related references
func TestRenamer_UpdatesRelatedReferences(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{ID: "related-node", Kind: "feature", Version: 1, Status: "draft", Title: "Related",
			Refs: domain.Ref{
				Related: []domain.RefLink{{Target: "target", Context: "see also"}},
			}},
	}

	result, err := r.Rename(nodes, "target", "renamed-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check related-node's Related now points to renamed-target
	for _, node := range result {
		if node.ID == "related-node" {
			if len(node.Refs.Related) != 1 {
				t.Fatalf("expected 1 Related ref, got %d", len(node.Refs.Related))
			}
			if node.Refs.Related[0].Target != "renamed-target" {
				t.Errorf("expected Related target 'renamed-target', got %q", node.Refs.Related[0].Target)
			}
			// Context should be preserved
			if node.Refs.Related[0].Context != "see also" {
				t.Errorf("expected context preserved, got %q", node.Refs.Related[0].Context)
			}
		}
	}
}

// Test renaming updates multiple references from different nodes
func TestRenamer_UpdatesMultipleReferences(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{ID: "referrer1", Kind: "mechanic", Version: 1, Status: "draft", Title: "Referrer 1",
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "target"}}}},
		{ID: "referrer2", Kind: "mechanic", Version: 1, Status: "draft", Title: "Referrer 2",
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "target"}}}},
		{ID: "referrer3", Kind: "feature", Version: 1, Status: "draft", Title: "Referrer 3",
			Refs: domain.Ref{Related: []domain.RefLink{{Target: "target"}}}},
	}

	result, err := r.Rename(nodes, "target", "new-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check all referrers have updated references
	for _, node := range result {
		switch node.ID {
		case "referrer1", "referrer2":
			if len(node.Refs.Uses) != 1 || node.Refs.Uses[0].Target != "new-target" {
				t.Errorf("expected %s Uses to point to 'new-target', got %v", node.ID, node.Refs.Uses)
			}
		case "referrer3":
			if len(node.Refs.Related) != 1 || node.Refs.Related[0].Target != "new-target" {
				t.Errorf("expected %s Related to point to 'new-target', got %v", node.ID, node.Refs.Related)
			}
		}
	}
}

// Test renaming updates both Uses and Related in same node
func TestRenamer_UpdatesBothUsesAndRelatedInSameNode(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{ID: "referrer", Kind: "mechanic", Version: 1, Status: "draft", Title: "Referrer",
			Refs: domain.Ref{
				Uses:    []domain.RefLink{{Target: "target"}},
				Related: []domain.RefLink{{Target: "target"}},
			}},
	}

	result, err := r.Rename(nodes, "target", "new-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "referrer" {
			if len(node.Refs.Uses) != 1 || node.Refs.Uses[0].Target != "new-target" {
				t.Errorf("expected Uses target 'new-target', got %v", node.Refs.Uses)
			}
			if len(node.Refs.Related) != 1 || node.Refs.Related[0].Target != "new-target" {
				t.Errorf("expected Related target 'new-target', got %v", node.Refs.Related)
			}
		}
	}
}

// Test renaming does not affect non-matching references
func TestRenamer_DoesNotAffectOtherReferences(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{ID: "other", Kind: "system", Version: 1, Status: "draft", Title: "Other"},
		{ID: "referrer", Kind: "mechanic", Version: 1, Status: "draft", Title: "Referrer",
			Refs: domain.Ref{
				Uses: []domain.RefLink{
					{Target: "target"},
					{Target: "other"},
				},
			}},
	}

	result, err := r.Rename(nodes, "target", "new-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "referrer" {
			if len(node.Refs.Uses) != 2 {
				t.Fatalf("expected 2 Uses refs, got %d", len(node.Refs.Uses))
			}
			// One should be renamed, one should stay the same
			targets := make(map[string]bool)
			for _, ref := range node.Refs.Uses {
				targets[ref.Target] = true
			}
			if !targets["new-target"] {
				t.Error("expected new-target in Uses refs")
			}
			if !targets["other"] {
				t.Error("expected other to remain in Uses refs")
			}
			if targets["target"] {
				t.Error("old target should not be in Uses refs")
			}
		}
	}
}

// Test renaming preserves RefLink context and resolved fields
func TestRenamer_PreservesRefLinkProperties(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{ID: "referrer", Kind: "mechanic", Version: 1, Status: "draft", Title: "Referrer",
			Refs: domain.Ref{
				Uses: []domain.RefLink{{Target: "target", Context: "depends on", Resolved: true}},
			}},
	}

	result, err := r.Rename(nodes, "target", "new-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "referrer" {
			ref := node.Refs.Uses[0]
			if ref.Target != "new-target" {
				t.Errorf("expected target 'new-target', got %q", ref.Target)
			}
			if ref.Context != "depends on" {
				t.Errorf("expected context preserved, got %q", ref.Context)
			}
			if !ref.Resolved {
				t.Error("expected resolved flag preserved")
			}
		}
	}
}

// ===== VERSION INCREMENT TESTS =====

// Test that referencing nodes get version incremented
func TestRenamer_IncrementsVersionOnReferencingNodes(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{ID: "referrer", Kind: "mechanic", Version: 3, Status: "draft", Title: "Referrer",
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "target"}}}},
	}

	result, err := r.Rename(nodes, "target", "new-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "referrer" {
			if node.Version != 4 {
				t.Errorf("expected referrer version to increment from 3 to 4, got %d", node.Version)
			}
		}
	}
}

// Test that non-referencing nodes keep their version
func TestRenamer_DoesNotIncrementVersionOnNonReferencingNodes(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{ID: "unrelated", Kind: "mechanic", Version: 5, Status: "draft", Title: "Unrelated"},
	}

	result, err := r.Rename(nodes, "target", "new-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "unrelated" {
			if node.Version != 5 {
				t.Errorf("expected unrelated version to stay 5, got %d", node.Version)
			}
		}
	}
}

// ===== ERROR HANDLING TESTS =====

// Test renaming non-existent node returns error
func TestRenamer_ErrorOnNonExistentNode(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "existing", Kind: "system", Version: 1, Status: "draft", Title: "Existing"},
	}

	_, err := r.Rename(nodes, "non-existent", "new-id")

	if err == nil {
		t.Fatal("expected error for non-existent node, got nil")
	}
}

// Test renaming to existing ID returns error
func TestRenamer_ErrorOnIDCollision(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "node-a", Kind: "system", Version: 1, Status: "draft", Title: "Node A"},
		{ID: "node-b", Kind: "system", Version: 1, Status: "draft", Title: "Node B"},
	}

	_, err := r.Rename(nodes, "node-a", "node-b")

	if err == nil {
		t.Fatal("expected error for ID collision, got nil")
	}
}

// Test renaming with empty old ID returns error
func TestRenamer_ErrorOnEmptyOldID(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "existing", Kind: "system", Version: 1, Status: "draft", Title: "Existing"},
	}

	_, err := r.Rename(nodes, "", "new-id")

	if err == nil {
		t.Fatal("expected error for empty old ID, got nil")
	}
}

// Test renaming with empty new ID returns error
func TestRenamer_ErrorOnEmptyNewID(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "existing", Kind: "system", Version: 1, Status: "draft", Title: "Existing"},
	}

	_, err := r.Rename(nodes, "existing", "")

	if err == nil {
		t.Fatal("expected error for empty new ID, got nil")
	}
}

// Test renaming with nil nodes slice returns error
func TestRenamer_ErrorOnNilNodes(t *testing.T) {
	r := refactor.NewRenamer()

	_, err := r.Rename(nil, "old", "new")

	if err == nil {
		t.Fatal("expected error for nil nodes, got nil")
	}
}

// Test renaming with empty nodes slice returns error
func TestRenamer_ErrorOnEmptyNodes(t *testing.T) {
	r := refactor.NewRenamer()

	_, err := r.Rename([]domain.Node{}, "old", "new")

	if err == nil {
		t.Fatal("expected error for empty nodes, got nil")
	}
}

// Test renaming to same ID returns error (no-op is an error)
func TestRenamer_ErrorOnSameID(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "same-id", Kind: "system", Version: 1, Status: "draft", Title: "Same"},
	}

	_, err := r.Rename(nodes, "same-id", "same-id")

	if err == nil {
		t.Fatal("expected error when renaming to same ID, got nil")
	}
}

// ===== EDGE CASES =====

// Test renaming with no references at all
func TestRenamer_NoReferences(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "a", Kind: "system", Version: 1, Status: "draft", Title: "A"},
		{ID: "b", Kind: "system", Version: 1, Status: "draft", Title: "B"},
		{ID: "c", Kind: "system", Version: 1, Status: "draft", Title: "C"},
	}

	result, err := r.Rename(nodes, "a", "renamed-a")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(result))
	}

	// Check renamed node exists
	var found bool
	for _, node := range result {
		if node.ID == "renamed-a" {
			found = true
		}
	}
	if !found {
		t.Error("expected renamed-a to exist")
	}
}

// Test renaming when target node has its own references (outgoing refs should be preserved)
func TestRenamer_PreservesOutgoingRefs(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target",
			Refs: domain.Ref{
				Uses: []domain.RefLink{{Target: "dependency"}},
			}},
		{ID: "dependency", Kind: "system", Version: 1, Status: "draft", Title: "Dependency"},
	}

	result, err := r.Rename(nodes, "target", "new-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that new-target still references dependency
	for _, node := range result {
		if node.ID == "new-target" {
			if len(node.Refs.Uses) != 1 {
				t.Fatalf("expected 1 Uses ref on renamed node, got %d", len(node.Refs.Uses))
			}
			if node.Refs.Uses[0].Target != "dependency" {
				t.Errorf("expected outgoing ref to dependency preserved, got %q", node.Refs.Uses[0].Target)
			}
		}
	}
}

// Test renaming in a chain: A -> B -> C, rename B
func TestRenamer_ChainedReferences(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "c", Kind: "system", Version: 1, Status: "draft", Title: "C"},
		{ID: "b", Kind: "system", Version: 1, Status: "draft", Title: "B",
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "c"}}}},
		{ID: "a", Kind: "system", Version: 1, Status: "draft", Title: "A",
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "b"}}}},
	}

	result, err := r.Rename(nodes, "b", "b-renamed")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		switch node.ID {
		case "a":
			// A should now reference b-renamed
			if len(node.Refs.Uses) != 1 || node.Refs.Uses[0].Target != "b-renamed" {
				t.Errorf("expected A to reference b-renamed, got %v", node.Refs.Uses)
			}
		case "b-renamed":
			// b-renamed should still reference c
			if len(node.Refs.Uses) != 1 || node.Refs.Uses[0].Target != "c" {
				t.Errorf("expected b-renamed to still reference c, got %v", node.Refs.Uses)
			}
		}
	}
}

// Test renaming with self-reference (node references itself)
func TestRenamer_SelfReference(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "self", Kind: "system", Version: 1, Status: "draft", Title: "Self",
			Refs: domain.Ref{
				Related: []domain.RefLink{{Target: "self"}},
			}},
	}

	result, err := r.Rename(nodes, "self", "new-self")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Self reference should also be updated
	for _, node := range result {
		if node.ID == "new-self" {
			if len(node.Refs.Related) != 1 || node.Refs.Related[0].Target != "new-self" {
				t.Errorf("expected self-reference updated to new-self, got %v", node.Refs.Related)
			}
		}
	}
}

// Test renaming with multiple refs to same target (deduplication not needed, just update all)
func TestRenamer_MultipleRefsToSameTarget(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{ID: "referrer", Kind: "mechanic", Version: 1, Status: "draft", Title: "Referrer",
			Refs: domain.Ref{
				Uses: []domain.RefLink{
					{Target: "target", Context: "first"},
					{Target: "target", Context: "second"},
				},
			}},
	}

	result, err := r.Rename(nodes, "target", "new-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "referrer" {
			if len(node.Refs.Uses) != 2 {
				t.Fatalf("expected 2 Uses refs, got %d", len(node.Refs.Uses))
			}
			for i, ref := range node.Refs.Uses {
				if ref.Target != "new-target" {
					t.Errorf("expected Uses[%d] target to be new-target, got %q", i, ref.Target)
				}
			}
		}
	}
}

// Test that EmitsEvents and Vocabulary are preserved
func TestRenamer_PreservesEmitsEventsAndVocabulary(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target",
			Refs: domain.Ref{
				EmitsEvents: []string{"event1", "event2"},
				Vocabulary:  []string{"term1", "term2"},
			}},
	}

	result, err := r.Rename(nodes, "target", "new-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "new-target" {
			if len(node.Refs.EmitsEvents) != 2 {
				t.Errorf("expected 2 EmitsEvents, got %d", len(node.Refs.EmitsEvents))
			}
			if len(node.Refs.Vocabulary) != 2 {
				t.Errorf("expected 2 Vocabulary, got %d", len(node.Refs.Vocabulary))
			}
		}
	}
}

// Test that EmitsEvents references pointing to renamed node are updated
func TestRenamer_UpdatesEmitsEventsReferences(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "events/player-died", Kind: "event", Version: 1, Status: "draft", Title: "Player Died Event"},
		{ID: "systems/combat", Kind: "system", Version: 1, Status: "draft", Title: "Combat System",
			Refs: domain.Ref{
				EmitsEvents: []string{"events/player-died", "events/damage-dealt"},
			}},
	}

	result, err := r.Rename(nodes, "events/player-died", "events/player-death")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "systems/combat" {
			if len(node.Refs.EmitsEvents) != 2 {
				t.Fatalf("expected 2 EmitsEvents, got %d", len(node.Refs.EmitsEvents))
			}
			// First event should be updated, second should remain unchanged
			if node.Refs.EmitsEvents[0] != "events/player-death" {
				t.Errorf("expected EmitsEvents[0] to be 'events/player-death', got %q", node.Refs.EmitsEvents[0])
			}
			if node.Refs.EmitsEvents[1] != "events/damage-dealt" {
				t.Errorf("expected EmitsEvents[1] to remain 'events/damage-dealt', got %q", node.Refs.EmitsEvents[1])
			}
			// Version should be incremented
			if node.Version != 2 {
				t.Errorf("expected version to be incremented to 2, got %d", node.Version)
			}
		}
	}
}

// Test that Vocabulary references pointing to renamed node are updated
func TestRenamer_UpdatesVocabularyReferences(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "glossary/game-terms", Kind: "glossary", Version: 1, Status: "draft", Title: "Game Terms"},
		{ID: "systems/combat", Kind: "system", Version: 1, Status: "draft", Title: "Combat System",
			Refs: domain.Ref{
				Vocabulary: []string{"glossary/game-terms", "glossary/combat-terms"},
			}},
	}

	result, err := r.Rename(nodes, "glossary/game-terms", "glossary/common-terms")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "systems/combat" {
			if len(node.Refs.Vocabulary) != 2 {
				t.Fatalf("expected 2 Vocabulary refs, got %d", len(node.Refs.Vocabulary))
			}
			// First term should be updated, second should remain unchanged
			if node.Refs.Vocabulary[0] != "glossary/common-terms" {
				t.Errorf("expected Vocabulary[0] to be 'glossary/common-terms', got %q", node.Refs.Vocabulary[0])
			}
			if node.Refs.Vocabulary[1] != "glossary/combat-terms" {
				t.Errorf("expected Vocabulary[1] to remain 'glossary/combat-terms', got %q", node.Refs.Vocabulary[1])
			}
			// Version should be incremented
			if node.Version != 2 {
				t.Errorf("expected version to be incremented to 2, got %d", node.Version)
			}
		}
	}
}

// Test that all reference types are updated together
func TestRenamer_UpdatesAllReferenceTypes(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{ID: "referrer", Kind: "system", Version: 1, Status: "draft", Title: "Referrer",
			Refs: domain.Ref{
				Uses:        []domain.RefLink{{Target: "target"}},
				Related:     []domain.RefLink{{Target: "target"}},
				EmitsEvents: []string{"target"},
				Vocabulary:  []string{"target"},
			}},
	}

	result, err := r.Rename(nodes, "target", "new-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, node := range result {
		if node.ID == "referrer" {
			// All reference types should point to new-target
			if node.Refs.Uses[0].Target != "new-target" {
				t.Errorf("expected Uses[0] to be 'new-target', got %q", node.Refs.Uses[0].Target)
			}
			if node.Refs.Related[0].Target != "new-target" {
				t.Errorf("expected Related[0] to be 'new-target', got %q", node.Refs.Related[0].Target)
			}
			if node.Refs.EmitsEvents[0] != "new-target" {
				t.Errorf("expected EmitsEvents[0] to be 'new-target', got %q", node.Refs.EmitsEvents[0])
			}
			if node.Refs.Vocabulary[0] != "new-target" {
				t.Errorf("expected Vocabulary[0] to be 'new-target', got %q", node.Refs.Vocabulary[0])
			}
			// Version should be incremented only once
			if node.Version != 2 {
				t.Errorf("expected version to be 2, got %d", node.Version)
			}
		}
	}
}

// ===== RESULT COUNT TESTS =====

// Test that rename returns same number of nodes
func TestRenamer_ReturnsCorrectNodeCount(t *testing.T) {
	r := refactor.NewRenamer()

	nodes := []domain.Node{
		{ID: "a", Kind: "system", Version: 1, Status: "draft", Title: "A"},
		{ID: "b", Kind: "system", Version: 1, Status: "draft", Title: "B"},
		{ID: "c", Kind: "system", Version: 1, Status: "draft", Title: "C"},
	}

	result, err := r.Rename(nodes, "a", "renamed-a")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != len(nodes) {
		t.Errorf("expected %d nodes, got %d", len(nodes), len(result))
	}
}

// Test that rename doesn't modify original nodes slice
func TestRenamer_DoesNotModifyOriginal(t *testing.T) {
	r := refactor.NewRenamer()

	original := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{ID: "referrer", Kind: "mechanic", Version: 1, Status: "draft", Title: "Referrer",
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "target"}}}},
	}

	// Save original values
	origTargetID := original[0].ID
	origReferrerRef := original[1].Refs.Uses[0].Target

	_, err := r.Rename(original, "target", "new-target")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check originals are unchanged
	if original[0].ID != origTargetID {
		t.Errorf("original target ID modified from %q to %q", origTargetID, original[0].ID)
	}
	if original[1].Refs.Uses[0].Target != origReferrerRef {
		t.Errorf("original referrer ref modified from %q to %q", origReferrerRef, original[1].Refs.Uses[0].Target)
	}
}
