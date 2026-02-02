package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"os/user"

	"github.com/Toernblom/deco/internal/domain"
	"gopkg.in/yaml.v3"
)

// contentFields holds all fields that affect content hash.
// Excluded: ID (structural), Version (auto-incremented), Status (workflow),
// Reviewers (workflow), SourceFile (internal).
// Uses SortedStringMap and SortedInterfaceMap for deterministic hashing.
type contentFields struct {
	Kind        string                   `yaml:"kind"`
	Title       string                   `yaml:"title"`
	Summary     string                   `yaml:"summary"`
	Tags        []string                 `yaml:"tags,omitempty"`
	Refs        domain.Ref               `yaml:"refs,omitempty"`
	Issues      []domain.Issue           `yaml:"issues,omitempty"`
	Content     *domain.Content          `yaml:"content,omitempty"`
	Glossary    domain.SortedStringMap   `yaml:"glossary,omitempty"`
	Contracts   []domain.Contract        `yaml:"contracts,omitempty"`
	LLMContext  string                   `yaml:"llm_context,omitempty"`
	Constraints []domain.Constraint      `yaml:"constraints,omitempty"`
	Custom      domain.SortedInterfaceMap `yaml:"custom,omitempty"`
}

// ComputeContentHash computes a SHA-256 hash of the content fields.
// Returns 16 hex characters (first 64 bits of the hash).
// Used by all mutation commands to record content state in history.
//
// Included in hash: Kind, Title, Summary, Tags, Refs, Issues, Content,
// Glossary, Contracts, LLMContext, Constraints, Custom.
//
// Excluded from hash: ID, Version, Status, Reviewers (workflow metadata).
func ComputeContentHash(n domain.Node) string {
	fields := contentFields{
		Kind:        n.Kind,
		Title:       n.Title,
		Summary:     n.Summary,
		Tags:        n.Tags,
		Refs:        n.Refs,
		Issues:      n.Issues,
		Content:     n.Content,
		Glossary:    domain.SortedStringMap(n.Glossary),
		Contracts:   n.Contracts,
		LLMContext:  n.LLMContext,
		Constraints: n.Constraints,
		Custom:      domain.SortedInterfaceMap(n.Custom),
	}

	data, err := yaml.Marshal(fields)
	if err != nil {
		return ""
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:8])
}

// GetCurrentUser returns the current system username, or "unknown" if unavailable.
func GetCurrentUser() string {
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	return "unknown"
}

// CheckContentHash verifies that the current content hash matches the expected hash.
// Returns nil if expectHash is empty (no check), or if hashes match.
// Returns an ExitError with code 3 (conflict) if hashes don't match.
func CheckContentHash(n domain.Node, expectHash string) error {
	if expectHash == "" {
		return nil
	}

	currentHash := ComputeContentHash(n)
	if currentHash != expectHash {
		return NewExitErrorf(ExitCodeConflict,
			`Conflict detected on %s

  Expected hash:  %s
  Current hash:   %s

The node was modified since you last read it.

Options:
  1. Reload the node and reapply your changes
  2. Use --force to overwrite (loses concurrent changes)
  3. Use 'deco show %s' to see current state`,
			n.ID, expectHash, currentHash, n.ID)
	}
	return nil
}
