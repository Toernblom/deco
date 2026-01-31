package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
)

func TestIssue_Creation(t *testing.T) {
	issue := domain.Issue{
		ID:          "issue-001",
		Description: "Clarify damage calculation formula",
		Severity:    "high",
		Location:    "content.sections[0]",
		Resolved:    false,
	}

	if issue.ID != "issue-001" {
		t.Errorf("expected ID 'issue-001', got %q", issue.ID)
	}
	if issue.Description != "Clarify damage calculation formula" {
		t.Errorf("expected Description 'Clarify damage calculation formula', got %q", issue.Description)
	}
	if issue.Severity != "high" {
		t.Errorf("expected Severity 'high', got %q", issue.Severity)
	}
	if issue.Location != "content.sections[0]" {
		t.Errorf("expected Location 'content.sections[0]', got %q", issue.Location)
	}
	if issue.Resolved != false {
		t.Errorf("expected Resolved false, got %v", issue.Resolved)
	}
}

func TestIssue_SeverityLevels(t *testing.T) {
	severityLevels := []string{"low", "medium", "high", "critical"}

	for _, severity := range severityLevels {
		issue := domain.Issue{
			ID:          "test-issue",
			Description: "Test issue",
			Severity:    severity,
			Location:    "test.location",
			Resolved:    false,
		}

		if issue.Severity != severity {
			t.Errorf("expected Severity %q, got %q", severity, issue.Severity)
		}
	}
}

func TestIssue_LocationTracking(t *testing.T) {
	tests := []struct {
		name     string
		location string
	}{
		{
			name:     "section location",
			location: "content.sections[0]",
		},
		{
			name:     "nested field",
			location: "content.sections[2].blocks[1]",
		},
		{
			name:     "meta field",
			location: "meta.title",
		},
		{
			name:     "refs field",
			location: "refs.uses[3]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue := domain.Issue{
				ID:          "test-issue",
				Description: "Test description",
				Severity:    "medium",
				Location:    tt.location,
				Resolved:    false,
			}

			if issue.Location != tt.location {
				t.Errorf("expected Location %q, got %q", tt.location, issue.Location)
			}
		})
	}
}

func TestIssue_ResolvedState(t *testing.T) {
	// Create unresolved issue
	issue := domain.Issue{
		ID:          "test-issue",
		Description: "Test description",
		Severity:    "medium",
		Location:    "test.location",
		Resolved:    false,
	}

	if issue.Resolved {
		t.Errorf("expected new issue to be unresolved")
	}

	// Mark as resolved
	issue.Resolved = true

	if !issue.Resolved {
		t.Errorf("expected issue to be resolved after setting Resolved=true")
	}
}

func TestIssue_Serialization(t *testing.T) {
	original := domain.Issue{
		ID:          "issue-001",
		Description: "Define starvation threshold",
		Severity:    "high",
		Location:    "content.sections[1].blocks[0]",
		Resolved:    false,
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal Issue: %v", err)
	}

	// Unmarshal back
	var restored domain.Issue
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal Issue: %v", err)
	}

	// Compare
	if restored.ID != original.ID {
		t.Errorf("ID mismatch: got %q, want %q", restored.ID, original.ID)
	}
	if restored.Description != original.Description {
		t.Errorf("Description mismatch: got %q, want %q", restored.Description, original.Description)
	}
	if restored.Severity != original.Severity {
		t.Errorf("Severity mismatch: got %q, want %q", restored.Severity, original.Severity)
	}
	if restored.Location != original.Location {
		t.Errorf("Location mismatch: got %q, want %q", restored.Location, original.Location)
	}
	if restored.Resolved != original.Resolved {
		t.Errorf("Resolved mismatch: got %v, want %v", restored.Resolved, original.Resolved)
	}
}

func TestIssue_Validation(t *testing.T) {
	tests := []struct {
		name    string
		issue   domain.Issue
		wantErr bool
	}{
		{
			name: "valid issue",
			issue: domain.Issue{
				ID:          "issue-001",
				Description: "Test description",
				Severity:    "medium",
				Location:    "test.location",
				Resolved:    false,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			issue: domain.Issue{
				Description: "Test description",
				Severity:    "medium",
				Location:    "test.location",
				Resolved:    false,
			},
			wantErr: true,
		},
		{
			name: "missing description",
			issue: domain.Issue{
				ID:       "issue-001",
				Severity: "medium",
				Location: "test.location",
				Resolved: false,
			},
			wantErr: true,
		},
		{
			name: "invalid severity",
			issue: domain.Issue{
				ID:          "issue-001",
				Description: "Test description",
				Severity:    "invalid",
				Location:    "test.location",
				Resolved:    false,
			},
			wantErr: true,
		},
		{
			name: "missing location",
			issue: domain.Issue{
				ID:          "issue-001",
				Description: "Test description",
				Severity:    "medium",
				Resolved:    false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.issue.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Issue.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
