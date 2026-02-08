// Copyright (C) 2026 Anton Törnblom
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package cli

import (
	"fmt"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
)

// FollowConfig controls which refs to follow and how deep.
type FollowConfig struct {
	RefTypes []string // which ref types to follow: "uses", "related", "vocabulary", "all"
	Depth    int      // max depth: 0 = unlimited, 1 = direct only (default), N = N levels
}

// FollowedNode is a node plus metadata about why it was included.
type FollowedNode struct {
	Node        domain.Node
	FollowedVia string // empty for root nodes, e.g. "uses from systems/combat" for followed
	Depth       int    // 0 for root, 1+ for followed
}

// bfsEntry is an internal type for the BFS queue.
type bfsEntry struct {
	node  domain.Node
	depth int
}

// collectFollowedNodes expands a set of root nodes by following their refs.
func collectFollowedNodes(rootNodes []domain.Node, allNodes []domain.Node, config FollowConfig) []FollowedNode {
	// Build lookup map from allNodes
	nodeByID := make(map[string]domain.Node, len(allNodes))
	for _, n := range allNodes {
		nodeByID[n.ID] = n
	}

	// Determine which ref types to follow
	followAll := false
	followUses := false
	followRelated := false
	followVocabulary := false
	for _, rt := range config.RefTypes {
		switch strings.ToLower(rt) {
		case "all":
			followAll = true
		case "uses":
			followUses = true
		case "related":
			followRelated = true
		case "vocabulary":
			followVocabulary = true
		}
	}
	if followAll {
		followUses = true
		followRelated = true
		followVocabulary = true
	}

	// Initialize result with root nodes
	seen := make(map[string]bool, len(rootNodes))
	result := make([]FollowedNode, 0, len(rootNodes))
	queue := make([]bfsEntry, 0, len(rootNodes))

	for _, n := range rootNodes {
		seen[n.ID] = true
		result = append(result, FollowedNode{
			Node:        n,
			FollowedVia: "",
			Depth:       0,
		})
		queue = append(queue, bfsEntry{node: n, depth: 0})
	}

	// BFS expansion
	for len(queue) > 0 {
		entry := queue[0]
		queue = queue[1:]

		// Check depth limit before collecting refs
		nextDepth := entry.depth + 1
		if config.Depth > 0 && nextDepth > config.Depth {
			continue
		}

		// Collect ref target IDs with their ref type
		type refTarget struct {
			id      string
			refType string
		}
		var targets []refTarget

		if followUses {
			for _, rl := range entry.node.Refs.Uses {
				targets = append(targets, refTarget{id: rl.Target, refType: "uses"})
			}
		}
		if followRelated {
			for _, rl := range entry.node.Refs.Related {
				targets = append(targets, refTarget{id: rl.Target, refType: "related"})
			}
		}
		if followVocabulary {
			for _, v := range entry.node.Refs.Vocabulary {
				targets = append(targets, refTarget{id: v, refType: "vocabulary"})
			}
		}

		for _, t := range targets {
			if seen[t.id] {
				continue
			}
			targetNode, exists := nodeByID[t.id]
			if !exists {
				continue
			}

			seen[t.id] = true
			result = append(result, FollowedNode{
				Node:        targetNode,
				FollowedVia: fmt.Sprintf("%s from %s", t.refType, entry.node.ID),
				Depth:       nextDepth,
			})

			// Enqueue for further expansion if depth allows
			if config.Depth == 0 || nextDepth < config.Depth {
				queue = append(queue, bfsEntry{node: targetNode, depth: nextDepth})
			}
		}
	}

	return result
}

// parseCompactFollowFlag parses the --follow flag for compact export.
// "" → nil (no following)
// "uses" → FollowConfig{RefTypes: ["uses"], Depth: 1}
// "related" → FollowConfig{RefTypes: ["related"], Depth: 1}
// "all" → FollowConfig{RefTypes: ["all"], Depth: 1}
// bare --follow (no value) defaults to "uses" — this is handled at the flag level
func parseCompactFollowFlag(follow string) *FollowConfig {
	if follow == "" {
		return nil
	}

	return &FollowConfig{
		RefTypes: []string{follow},
		Depth:    1,
	}
}
