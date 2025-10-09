package tree

import (
	"sort"
	"strings"
)

// SubtreeSelector selects optimal content subtrees based on content scores
type SubtreeSelector struct {
	minContentScore float64 // Minimum score for a subtree to be considered
	maxSubtrees     int     // Maximum number of subtrees to select
	avoidOverlap    bool    // Avoid selecting overlapping subtrees
	prioritizeDepth bool    // Prefer deeper (more specific) subtrees
}

// NewSubtreeSelector creates a new subtree selector with default settings
func NewSubtreeSelector() *SubtreeSelector {
	return &SubtreeSelector{
		minContentScore: 0.1, // 10% of max score
		maxSubtrees:     5,   // Select up to 5 top subtrees
		avoidOverlap:    true,
		prioritizeDepth: true,
	}
}

// SubtreeCandidate represents a potential subtree for selection
type SubtreeCandidate struct {
	Node             *TextNode
	ContentScore     float64
	WordCount        int
	Depth            int
	ContainsHeadings bool
	HeadingCount     int
	LinkDensity      float64
	SelectionScore   float64 // Combined score for selection ranking
}

// SelectSubtrees selects the best content subtrees from the tree
func (ss *SubtreeSelector) SelectSubtrees(root *TextNode) []*TextNode {
	if root == nil {
		return nil
	}

	// Find all potential subtree candidates
	candidates := ss.findSubtreeCandidates(root)

	// Calculate selection scores for each candidate
	ss.calculateSelectionScores(candidates)

	// Sort candidates by selection score
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].SelectionScore > candidates[j].SelectionScore
	})

	// Select the best non-overlapping subtrees
	return ss.selectBestSubtrees(candidates)
}

// findSubtreeCandidates identifies potential content subtrees
func (ss *SubtreeSelector) findSubtreeCandidates(node *TextNode) []*SubtreeCandidate {
	var candidates []*SubtreeCandidate

	if node == nil {
		return candidates
	}

	// Consider this node as a candidate if it meets minimum criteria
	if ss.isViableCandidate(node) {
		candidate := &SubtreeCandidate{
			Node:         node,
			ContentScore: node.ContentScore,
			WordCount:    ss.calculateSubtreeWordCount(node),
			Depth:        node.Depth,
		}

		// Calculate additional metrics
		candidate.ContainsHeadings, candidate.HeadingCount = ss.analyzeHeadings(node)
		candidate.LinkDensity = ss.calculateSubtreeLinkDensity(node)

		candidates = append(candidates, candidate)
	}

	// Recursively find candidates in children
	for _, child := range node.Children {
		childCandidates := ss.findSubtreeCandidates(child)
		candidates = append(candidates, childCandidates...)
	}

	return candidates
}

// isViableCandidate checks if a node is a viable subtree candidate
func (ss *SubtreeSelector) isViableCandidate(node *TextNode) bool {
	// Must meet minimum content score
	if node.ContentScore < ss.minContentScore {
		return false
	}

	// Must have meaningful word count
	wordCount := ss.calculateSubtreeWordCount(node)
	if wordCount < 10 {
		return false
	}

	// Avoid pure navigation or structural elements
	tag := node.Tag
	structuralTags := map[string]bool{
		"nav": true, "header": true, "footer": true,
		"aside": true, "menu": true,
	}
	if structuralTags[tag] {
		return false
	}

	// Prefer content containers
	contentTags := map[string]bool{
		"article": true, "main": true, "section": true,
		"div": true, "p": true, "blockquote": true,
	}
	if contentTags[tag] || tag == "#text" {
		return true
	}

	// Allow other elements if they have high content scores
	return node.ContentScore > 0.3
}

// calculateSubtreeWordCount calculates total word count for a subtree
func (ss *SubtreeSelector) calculateSubtreeWordCount(node *TextNode) int {
	if node == nil {
		return 0
	}

	wordCount := node.WordCount
	for _, child := range node.Children {
		wordCount += ss.calculateSubtreeWordCount(child)
	}

	return wordCount
}

// analyzeHeadings analyzes headings within a subtree
func (ss *SubtreeSelector) analyzeHeadings(node *TextNode) (bool, int) {
	headingCount := 0
	ss.countHeadings(node, &headingCount)
	return headingCount > 0, headingCount
}

// countHeadings recursively counts headings in a subtree
func (ss *SubtreeSelector) countHeadings(node *TextNode, count *int) {
	if node == nil {
		return
	}

	headings := map[string]bool{
		"h1": true, "h2": true, "h3": true,
		"h4": true, "h5": true, "h6": true,
	}
	if headings[node.Tag] {
		*count++
	}

	for _, child := range node.Children {
		ss.countHeadings(child, count)
	}
}

// calculateSubtreeLinkDensity calculates link density for a subtree
func (ss *SubtreeSelector) calculateSubtreeLinkDensity(node *TextNode) float64 {
	totalWords := ss.calculateSubtreeWordCount(node)
	if totalWords == 0 {
		return 0
	}

	linkWords := ss.calculateLinkWords(node)
	return float64(linkWords) / float64(totalWords)
}

// calculateLinkWords calculates word count within links in a subtree
func (ss *SubtreeSelector) calculateLinkWords(node *TextNode) int {
	if node == nil {
		return 0
	}

	linkWords := 0

	// If this is a link, count its words
	if node.Tag == "a" {
		linkWords += ss.calculateSubtreeWordCount(node)
	} else {
		// Otherwise, count link words in children
		for _, child := range node.Children {
			linkWords += ss.calculateLinkWords(child)
		}
	}

	return linkWords
}

// calculateSelectionScores calculates selection scores for all candidates
func (ss *SubtreeSelector) calculateSelectionScores(candidates []*SubtreeCandidate) {
	for _, candidate := range candidates {
		score := candidate.ContentScore

		// Bonus for word count (logarithmic scaling)
		if candidate.WordCount > 0 {
			wordBonus := 0.1 * (1.0 + 0.1*float64(candidate.WordCount))
			if candidate.WordCount > 100 {
				wordBonus = 0.2 // Cap the word bonus
			}
			score += wordBonus
		}

		// Bonus for containing headings
		if candidate.ContainsHeadings {
			headingBonus := 0.15 + 0.05*float64(candidate.HeadingCount)
			if headingBonus > 0.3 {
				headingBonus = 0.3 // Cap heading bonus
			}
			score += headingBonus
		}

		// Penalty for high link density
		if candidate.LinkDensity > 0.3 {
			linkPenalty := 0.2 * candidate.LinkDensity
			score -= linkPenalty
		}

		// Depth bonus/penalty based on preference
		if ss.prioritizeDepth {
			// Prefer moderately deep content (not too shallow, not too deep)
			optimalDepth := 5
			depthDiff := abs(candidate.Depth - optimalDepth)
			depthScore := 0.1 * (1.0 - float64(depthDiff)/10.0)
			if depthScore > 0 {
				score += depthScore
			}
		}

		// Ensure score is non-negative
		if score < 0 {
			score = 0
		}

		candidate.SelectionScore = score
	}
}

// selectBestSubtrees selects the best non-overlapping subtrees
func (ss *SubtreeSelector) selectBestSubtrees(candidates []*SubtreeCandidate) []*TextNode {
	var selected []*TextNode
	var selectedNodes []*TextNode

	for _, candidate := range candidates {
		if len(selected) >= ss.maxSubtrees {
			break
		}

		// Check for overlap with already selected subtrees
		if ss.avoidOverlap && ss.hasOverlap(candidate.Node, selectedNodes) {
			continue
		}

		selected = append(selected, candidate.Node)
		selectedNodes = append(selectedNodes, candidate.Node)
	}

	return selected
}

// hasOverlap checks if a node overlaps with any of the selected nodes
func (ss *SubtreeSelector) hasOverlap(node *TextNode, selected []*TextNode) bool {
	for _, selectedNode := range selected {
		if ss.isAncestorOrDescendant(node, selectedNode) {
			return true
		}
	}
	return false
}

// isAncestorOrDescendant checks if two nodes have ancestor-descendant relationship
func (ss *SubtreeSelector) isAncestorOrDescendant(node1, node2 *TextNode) bool {
	return ss.isAncestor(node1, node2) || ss.isAncestor(node2, node1)
}

// isAncestor checks if node1 is an ancestor of node2
func (ss *SubtreeSelector) isAncestor(ancestor, descendant *TextNode) bool {
	current := descendant.Parent
	for current != nil {
		if current == ancestor {
			return true
		}
		current = current.Parent
	}
	return false
}

// GetSubtreeContent extracts text content from selected subtrees
func (ss *SubtreeSelector) GetSubtreeContent(subtrees []*TextNode) []string {
	var content []string
	for _, subtree := range subtrees {
		text := ss.extractSubtreeText(subtree)
		if text != "" {
			content = append(content, text)
		}
	}
	return content
}

// extractSubtreeText extracts readable text from a subtree
func (ss *SubtreeSelector) extractSubtreeText(node *TextNode) string {
	if node == nil {
		return ""
	}

	if node.Tag == "#text" {
		return strings.TrimSpace(node.Text)
	}

	var parts []string
	for _, child := range node.Children {
		text := ss.extractSubtreeText(child)
		if text != "" {
			parts = append(parts, text)
		}
	}

	return strings.Join(parts, " ")
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
