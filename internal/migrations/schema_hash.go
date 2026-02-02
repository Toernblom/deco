package migrations

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/Toernblom/deco/internal/storage/config"
)

// ComputeSchemaHash computes a deterministic SHA-256 hash of the schema configuration.
// It includes CustomBlockTypes and SchemaRules, using sorted keys for determinism.
// Returns the first 16 hex characters (64 bits) for readability.
// Returns empty string if no schema constraints are defined.
func ComputeSchemaHash(cfg config.Config) string {
	// If no schema constraints, return empty
	if len(cfg.CustomBlockTypes) == 0 && len(cfg.SchemaRules) == 0 {
		return ""
	}

	// Build canonical representation with sorted keys
	canonical := buildCanonicalSchema(cfg)

	// Compute SHA-256 hash
	hash := sha256.Sum256(canonical)

	// Return first 16 hex chars (64 bits)
	return hex.EncodeToString(hash[:8])
}

// buildCanonicalSchema creates a deterministic JSON representation of the schema.
func buildCanonicalSchema(cfg config.Config) []byte {
	// Create a structure with sorted keys
	schema := map[string]interface{}{}

	// Add custom block types (sorted)
	if len(cfg.CustomBlockTypes) > 0 {
		blockTypes := make(map[string]interface{})
		for name, bt := range cfg.CustomBlockTypes {
			// Sort required fields for determinism
			fields := make([]string, len(bt.RequiredFields))
			copy(fields, bt.RequiredFields)
			sort.Strings(fields)
			blockTypes[name] = map[string]interface{}{
				"required_fields": fields,
			}
		}
		schema["custom_block_types"] = sortedMap(blockTypes)
	}

	// Add schema rules (sorted)
	if len(cfg.SchemaRules) > 0 {
		rules := make(map[string]interface{})
		for kind, rule := range cfg.SchemaRules {
			// Sort required fields for determinism
			fields := make([]string, len(rule.RequiredFields))
			copy(fields, rule.RequiredFields)
			sort.Strings(fields)
			rules[kind] = map[string]interface{}{
				"required_fields": fields,
			}
		}
		schema["schema_rules"] = sortedMap(rules)
	}

	// Marshal to JSON (deterministic with sorted keys)
	data, _ := json.Marshal(sortedMap(schema))
	return data
}

// sortedMap creates a sorted representation for JSON marshaling.
// This ensures deterministic output regardless of Go map iteration order.
func sortedMap(m map[string]interface{}) []sortedEntry {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]sortedEntry, len(keys))
	for i, k := range keys {
		entries[i] = sortedEntry{Key: k, Value: m[k]}
	}
	return entries
}

// sortedEntry represents a key-value pair for deterministic JSON output.
type sortedEntry struct {
	Key   string      `json:"k"`
	Value interface{} `json:"v"`
}

// SchemaVersionMatches checks if the config's SchemaVersion matches the computed hash.
// Returns true if they match or if there are no schema constraints.
func SchemaVersionMatches(cfg config.Config) bool {
	computed := ComputeSchemaHash(cfg)

	// Both empty means no schema constraints - always matches
	if computed == "" && cfg.SchemaVersion == "" {
		return true
	}

	return computed == cfg.SchemaVersion
}
