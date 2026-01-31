package yaml

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"gopkg.in/yaml.v3"
)

// LocationTracker tracks the source locations (line/column) of YAML fields.
// This is essential for providing accurate error messages with file locations.
type LocationTracker struct {
	root     *yaml.Node
	filePath string
}

// NewLocationTracker creates a LocationTracker from YAML content.
func NewLocationTracker(content []byte) (*LocationTracker, error) {
	return NewLocationTrackerWithFile(content, "")
}

// NewLocationTrackerWithFile creates a LocationTracker with an associated file path.
func NewLocationTrackerWithFile(content []byte, filePath string) (*LocationTracker, error) {
	var root yaml.Node
	err := yaml.Unmarshal(content, &root)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &LocationTracker{
		root:     &root,
		filePath: filePath,
	}, nil
}

// GetLocation returns the location of a field specified by path.
// Path uses dot notation for nested fields (e.g., "metadata.author")
// and bracket notation for array indices (e.g., "tags[0]" or "[0].id").
// Returns a zero location if the path is not found.
func (t *LocationTracker) GetLocation(path string) domain.Location {
	if t.root == nil {
		return domain.Location{File: t.filePath}
	}

	node := t.findNode(path)
	if node == nil {
		return domain.Location{File: t.filePath}
	}

	return domain.Location{
		File:   t.filePath,
		Line:   node.Line,
		Column: node.Column,
	}
}

// GetValueLocation returns the location of a field's value.
// This finds the value node, not the key node.
func (t *LocationTracker) GetValueLocation(path string) domain.Location {
	if t.root == nil {
		return domain.Location{File: t.filePath}
	}

	if path == "" {
		return domain.Location{File: t.filePath}
	}

	// Parse the path into segments
	segments := parsePath(path)

	// Start from the root document node
	current := t.root
	if current.Kind == yaml.DocumentNode {
		if len(current.Content) == 0 {
			return domain.Location{File: t.filePath}
		}
		current = current.Content[0]
	}

	// Traverse all segments except the last one
	for i := 0; i < len(segments)-1; i++ {
		current = t.findChild(current, segments[i], false)
		if current == nil {
			return domain.Location{File: t.filePath}
		}
	}

	// For the last segment, find the value node
	lastSegment := segments[len(segments)-1]

	if lastSegment.isIndex {
		// For array index, the element itself is the value
		current = t.findChild(current, lastSegment, true)
		if current == nil {
			return domain.Location{File: t.filePath}
		}
		return domain.Location{
			File:   t.filePath,
			Line:   current.Line,
			Column: current.Column,
		}
	}

	// For object key, find the value node (not the key node)
	if current.Kind == yaml.MappingNode {
		for i := 0; i < len(current.Content); i += 2 {
			keyNode := current.Content[i]
			if keyNode.Value == lastSegment.key {
				// Return the value node location
				if i+1 < len(current.Content) {
					valueNode := current.Content[i+1]
					return domain.Location{
						File:   t.filePath,
						Line:   valueNode.Line,
						Column: valueNode.Column,
					}
				}
			}
		}
	}

	return domain.Location{File: t.filePath}
}

// findNode finds a node by path in the YAML tree.
func (t *LocationTracker) findNode(path string) *yaml.Node {
	if path == "" {
		return t.root
	}

	// Parse the path into segments
	segments := parsePath(path)

	// Start from the root document node
	current := t.root

	// The root is usually a DocumentNode, get its content
	if current.Kind == yaml.DocumentNode {
		if len(current.Content) == 0 {
			return nil
		}
		current = current.Content[0]
	}

	// Traverse the path
	for i, segment := range segments {
		isLastSegment := i == len(segments)-1
		current = t.findChild(current, segment, isLastSegment)
		if current == nil {
			return nil
		}
	}

	return current
}

// pathSegment represents a segment of a path (either a key or an array index)
type pathSegment struct {
	key   string
	index int
	isIndex bool
}

// parsePath parses a dot-notation path into segments.
// Examples:
//   "id" -> [{key: "id"}]
//   "metadata.author" -> [{key: "metadata"}, {key: "author"}]
//   "tags[0]" -> [{key: "tags"}, {index: 0, isIndex: true}]
//   "[0].id" -> [{index: 0, isIndex: true}, {key: "id"}]
func parsePath(path string) []pathSegment {
	var segments []pathSegment
	parts := strings.Split(path, ".")

	for _, part := range parts {
		// Check for array index notation
		if strings.Contains(part, "[") {
			// Split on '[' to separate key from index
			bracketIdx := strings.Index(part, "[")

			// Add the key part if it exists
			if bracketIdx > 0 {
				key := part[:bracketIdx]
				segments = append(segments, pathSegment{key: key})
			}

			// Extract and parse the index
			indexPart := part[bracketIdx+1:]
			if strings.HasSuffix(indexPart, "]") {
				indexStr := indexPart[:len(indexPart)-1]
				if idx, err := strconv.Atoi(indexStr); err == nil {
					segments = append(segments, pathSegment{index: idx, isIndex: true})
				}
			}
		} else {
			// Regular key
			segments = append(segments, pathSegment{key: part})
		}
	}

	return segments
}

// findChild finds a child node based on a path segment.
// If isLastSegment is true, returns the key node for location.
// Otherwise, returns the value node for further navigation.
func (t *LocationTracker) findChild(node *yaml.Node, segment pathSegment, isLastSegment bool) *yaml.Node {
	if node == nil {
		return nil
	}

	if segment.isIndex {
		// Handle array index
		if node.Kind == yaml.SequenceNode {
			if segment.index >= 0 && segment.index < len(node.Content) {
				return node.Content[segment.index]
			}
		}
		return nil
	}

	// Handle object key
	if node.Kind == yaml.MappingNode {
		// MappingNode content is [key1, value1, key2, value2, ...]
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			if keyNode.Value == segment.key {
				if isLastSegment {
					// Return the key node for location tracking
					return keyNode
				}
				// Return the value node for further navigation
				if i+1 < len(node.Content) {
					return node.Content[i+1]
				}
				return nil
			}
		}
	}

	return nil
}
