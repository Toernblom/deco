package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/user"
	"path/filepath"
	"sort"

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
	Docs        []domain.DocRef          `yaml:"docs,omitempty"`
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
	return ComputeContentHashWithDir(n, "")
}

// ComputeContentHashWithDir computes a content hash that includes referenced doc files.
// When projectRoot is non-empty, the contents of referenced .md files are included
// in the hash so that changes to doc files trigger version bumps during sync.
func ComputeContentHashWithDir(n domain.Node, projectRoot string) string {
	fields := contentFields{
		Kind:        n.Kind,
		Title:       n.Title,
		Summary:     n.Summary,
		Tags:        n.Tags,
		Refs:        n.Refs,
		Docs:        n.Docs,
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

	h := sha256.New()
	h.Write(data)

	// Include referenced doc file contents in hash
	if projectRoot != "" {
		docPaths := collectDocPaths(n)
		for _, p := range docPaths {
			fullPath := filepath.Join(projectRoot, p)
			if content, err := os.ReadFile(fullPath); err == nil {
				h.Write([]byte(p)) // include path as separator
				h.Write(content)
			}
			// Missing files are silently skipped (validation catches them)
		}
	}

	return hex.EncodeToString(h.Sum(nil)[:8])
}

// collectDocPaths gathers all doc file paths from node-level docs and doc blocks,
// returning them in sorted order for deterministic hashing.
func collectDocPaths(n domain.Node) []string {
	seen := make(map[string]bool)

	// Node-level docs
	for _, doc := range n.Docs {
		if doc.Path != "" {
			seen[doc.Path] = true
		}
	}

	// Doc blocks in content
	if n.Content != nil {
		for _, section := range n.Content.Sections {
			for _, block := range section.Blocks {
				if block.Type == "doc" {
					if p, ok := block.Data["path"].(string); ok && p != "" {
						seen[p] = true
					}
				}
			}
		}
	}

	paths := make([]string, 0, len(seen))
	for p := range seen {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

// GetCurrentUser returns the current system username, or "unknown" if unavailable.
func GetCurrentUser() string {
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	return "unknown"
}
