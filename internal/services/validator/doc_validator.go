package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
)

// DocValidator validates external doc references in nodes and doc blocks.
type DocValidator struct{}

// NewDocValidator creates a new doc validator.
func NewDocValidator() *DocValidator {
	return &DocValidator{}
}

// ValidateNodeDocs validates all node-level doc references.
func (dv *DocValidator) ValidateNodeDocs(node *domain.Node, projectRoot string, collector *errors.Collector) {
	if node == nil || len(node.Docs) == 0 {
		return
	}

	var location *domain.Location
	if node.SourceFile != "" {
		location = &domain.Location{File: node.SourceFile}
	}

	for _, doc := range node.Docs {
		dv.validateDocRef(doc.Path, doc.Keywords, node.ID, "", "", -1, projectRoot, location, collector)
	}
}

// ValidateDocBlock validates a single doc block's file reference and keywords.
func (dv *DocValidator) ValidateDocBlock(block *domain.Block, nodeID, sectionName string, blockIdx int, projectRoot string, collector *errors.Collector) {
	if block == nil {
		return
	}

	pathVal, ok := block.Data["path"]
	if !ok {
		return // Missing path is handled by block validator's requireField
	}

	path, ok := pathVal.(string)
	if !ok {
		return
	}

	var keywords []string
	if kw, ok := block.Data["keywords"]; ok {
		if kwList, ok := kw.([]interface{}); ok {
			for _, k := range kwList {
				if s, ok := k.(string); ok {
					keywords = append(keywords, s)
				}
			}
		}
	}

	dv.validateDocRef(path, keywords, nodeID, sectionName, "doc", blockIdx, projectRoot, nil, collector)
}

// validateDocRef checks that a doc file exists and contains required keywords.
func (dv *DocValidator) validateDocRef(path string, keywords []string, nodeID, sectionName, blockType string, blockIdx int, projectRoot string, location *domain.Location, collector *errors.Collector) {
	fullPath := filepath.Join(projectRoot, path)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		detail := fmt.Sprintf("node %q references doc file %q which does not exist", nodeID, path)
		if sectionName != "" {
			detail = fmt.Sprintf("in node %q, section %q, block %d: doc file %q not found", nodeID, sectionName, blockIdx, path)
		}
		collector.Add(domain.DecoError{
			Code:     "E055",
			Summary:  fmt.Sprintf("Doc file not found: %s", path),
			Detail:   detail,
			Location: location,
		})
		return // Skip keyword check if file doesn't exist
	}

	// Check keywords (case-insensitive substring matching)
	if len(keywords) == 0 {
		return
	}

	lowerContent := strings.ToLower(string(content))
	for _, keyword := range keywords {
		if !strings.Contains(lowerContent, strings.ToLower(keyword)) {
			detail := fmt.Sprintf("node %q: keyword %q not found in %s", nodeID, keyword, path)
			if sectionName != "" {
				detail = fmt.Sprintf("in node %q, section %q, block %d: keyword %q not found in %s", nodeID, sectionName, blockIdx, keyword, path)
			}
			collector.Add(domain.DecoError{
				Code:     "E056",
				Summary:  fmt.Sprintf("Missing keyword %q in %s", keyword, path),
				Detail:   detail,
				Location: location,
			})
		}
	}
}
