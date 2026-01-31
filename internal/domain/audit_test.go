package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Toernblom/deco/internal/domain"
)

func TestAuditEntry_Creation(t *testing.T) {
	now := time.Now()
	entry := domain.AuditEntry{
		Timestamp: now,
		NodeID:    "systems/food",
		Operation: "update",
		User:      "alice",
		Before:    map[string]interface{}{"status": "draft"},
		After:     map[string]interface{}{"status": "approved"},
	}

	if !entry.Timestamp.Equal(now) {
		t.Errorf("expected Timestamp %v, got %v", now, entry.Timestamp)
	}
	if entry.NodeID != "systems/food" {
		t.Errorf("expected NodeID 'systems/food', got %q", entry.NodeID)
	}
	if entry.Operation != "update" {
		t.Errorf("expected Operation 'update', got %q", entry.Operation)
	}
	if entry.User != "alice" {
		t.Errorf("expected User 'alice', got %q", entry.User)
	}
}

func TestAuditEntry_TimestampHandling(t *testing.T) {
	tests := []struct {
		name      string
		timestamp time.Time
	}{
		{
			name:      "current time",
			timestamp: time.Now(),
		},
		{
			name:      "past time",
			timestamp: time.Now().Add(-24 * time.Hour),
		},
		{
			name:      "UTC time",
			timestamp: time.Now().UTC(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := domain.AuditEntry{
				Timestamp: tt.timestamp,
				NodeID:    "test/node",
				Operation: "create",
				User:      "test",
			}

			if !entry.Timestamp.Equal(tt.timestamp) {
				t.Errorf("Timestamp mismatch: got %v, want %v", entry.Timestamp, tt.timestamp)
			}
		})
	}
}

func TestAuditEntry_OperationTypes(t *testing.T) {
	operations := []string{"create", "update", "delete", "set", "append", "unset", "move"}

	for _, op := range operations {
		t.Run(op, func(t *testing.T) {
			entry := domain.AuditEntry{
				Timestamp: time.Now(),
				NodeID:    "test/node",
				Operation: op,
				User:      "test",
			}

			if entry.Operation != op {
				t.Errorf("expected Operation %q, got %q", op, entry.Operation)
			}
		})
	}
}

func TestAuditEntry_BeforeAfterCapture(t *testing.T) {
	before := map[string]interface{}{
		"status":  "draft",
		"title":   "Old Title",
		"version": 1,
	}

	after := map[string]interface{}{
		"status":  "approved",
		"title":   "New Title",
		"version": 2,
	}

	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    "test/node",
		Operation: "update",
		User:      "alice",
		Before:    before,
		After:     after,
	}

	// Verify Before state
	if entry.Before["status"] != "draft" {
		t.Errorf("expected Before status 'draft', got %v", entry.Before["status"])
	}
	if entry.Before["title"] != "Old Title" {
		t.Errorf("expected Before title 'Old Title', got %v", entry.Before["title"])
	}
	if entry.Before["version"] != 1 {
		t.Errorf("expected Before version 1, got %v", entry.Before["version"])
	}

	// Verify After state
	if entry.After["status"] != "approved" {
		t.Errorf("expected After status 'approved', got %v", entry.After["status"])
	}
	if entry.After["title"] != "New Title" {
		t.Errorf("expected After title 'New Title', got %v", entry.After["title"])
	}
	if entry.After["version"] != 2 {
		t.Errorf("expected After version 2, got %v", entry.After["version"])
	}
}

func TestAuditEntry_CreateOperation(t *testing.T) {
	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    "systems/new-system",
		Operation: "create",
		User:      "bob",
		Before:    nil,
		After: map[string]interface{}{
			"id":      "systems/new-system",
			"kind":    "system",
			"version": 1,
			"status":  "draft",
			"title":   "New System",
		},
	}

	if entry.Operation != "create" {
		t.Errorf("expected Operation 'create', got %q", entry.Operation)
	}
	if entry.Before != nil {
		t.Errorf("expected Before to be nil for create operation")
	}
	if entry.After == nil {
		t.Errorf("expected After to be set for create operation")
	}
}

func TestAuditEntry_DeleteOperation(t *testing.T) {
	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    "systems/old-system",
		Operation: "delete",
		User:      "alice",
		Before: map[string]interface{}{
			"id":      "systems/old-system",
			"status":  "deprecated",
			"version": 5,
		},
		After: nil,
	}

	if entry.Operation != "delete" {
		t.Errorf("expected Operation 'delete', got %q", entry.Operation)
	}
	if entry.Before == nil {
		t.Errorf("expected Before to be set for delete operation")
	}
	if entry.After != nil {
		t.Errorf("expected After to be nil for delete operation")
	}
}

func TestAuditEntry_Serialization(t *testing.T) {
	now := time.Now().UTC()
	original := domain.AuditEntry{
		Timestamp: now,
		NodeID:    "test/node",
		Operation: "update",
		User:      "alice",
		Before:    map[string]interface{}{"status": "draft"},
		After:     map[string]interface{}{"status": "approved"},
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal AuditEntry: %v", err)
	}

	// Unmarshal back
	var restored domain.AuditEntry
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal AuditEntry: %v", err)
	}

	// Compare
	if !restored.Timestamp.Equal(original.Timestamp) {
		t.Errorf("Timestamp mismatch: got %v, want %v", restored.Timestamp, original.Timestamp)
	}
	if restored.NodeID != original.NodeID {
		t.Errorf("NodeID mismatch: got %q, want %q", restored.NodeID, original.NodeID)
	}
	if restored.Operation != original.Operation {
		t.Errorf("Operation mismatch: got %q, want %q", restored.Operation, original.Operation)
	}
	if restored.User != original.User {
		t.Errorf("User mismatch: got %q, want %q", restored.User, original.User)
	}
}

func TestAuditEntry_Validation(t *testing.T) {
	tests := []struct {
		name    string
		entry   domain.AuditEntry
		wantErr bool
	}{
		{
			name: "valid entry",
			entry: domain.AuditEntry{
				Timestamp: time.Now(),
				NodeID:    "test/node",
				Operation: "update",
				User:      "alice",
			},
			wantErr: false,
		},
		{
			name: "missing NodeID",
			entry: domain.AuditEntry{
				Timestamp: time.Now(),
				Operation: "update",
				User:      "alice",
			},
			wantErr: true,
		},
		{
			name: "missing Operation",
			entry: domain.AuditEntry{
				Timestamp: time.Now(),
				NodeID:    "test/node",
				User:      "alice",
			},
			wantErr: true,
		},
		{
			name: "missing User",
			entry: domain.AuditEntry{
				Timestamp: time.Now(),
				NodeID:    "test/node",
				Operation: "update",
			},
			wantErr: true,
		},
		{
			name: "zero timestamp",
			entry: domain.AuditEntry{
				Timestamp: time.Time{},
				NodeID:    "test/node",
				Operation: "update",
				User:      "alice",
			},
			wantErr: true,
		},
		{
			name: "invalid operation",
			entry: domain.AuditEntry{
				Timestamp: time.Now(),
				NodeID:    "test/node",
				Operation: "invalid_op",
				User:      "alice",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AuditEntry.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
