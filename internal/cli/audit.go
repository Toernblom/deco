package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"os/user"

	"github.com/Toernblom/deco/internal/domain"
	"gopkg.in/yaml.v3"
)

// contentFields holds only the fields that affect content hash
type contentFields struct {
	Title   string          `yaml:"title"`
	Summary string          `yaml:"summary"`
	Tags    []string        `yaml:"tags,omitempty"`
	Refs    domain.Ref      `yaml:"refs,omitempty"`
	Issues  []domain.Issue  `yaml:"issues,omitempty"`
	Content *domain.Content `yaml:"content,omitempty"`
}

// ComputeContentHash computes a SHA-256 hash of the content fields.
// Returns 16 hex characters (first 64 bits of the hash).
// Used by all mutation commands to record content state in history.
func ComputeContentHash(n domain.Node) string {
	fields := contentFields{
		Title:   n.Title,
		Summary: n.Summary,
		Tags:    n.Tags,
		Refs:    n.Refs,
		Issues:  n.Issues,
		Content: n.Content,
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
