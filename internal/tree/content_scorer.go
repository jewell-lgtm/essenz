package tree

import (
	"container/list"
	"strings"
)

// ContentScorer implements BFS-based content scoring with parent propagation
type ContentScorer struct {
	// Scoring weights for different content types
	textContentWeight float64
	headingWeights    map[string]float64
	linkPenalty       float64
	parentPropagation float64

	// Node removal criteria
	removePresentational bool
	removeHidden         bool
}

// NewContentScorer creates a new content scorer with default settings
func NewContentScorer() *ContentScorer {
	return &ContentScorer{
		textContentWeight: 1.0,
		headingWeights: map[string]float64{
			"h1": 5.0,
			"h2": 4.0,
			"h3": 3.0,
			"h4": 2.5,
			"h5": 2.0,
			"h6": 1.5,
		},
		linkPenalty:          0.3, // Reduce score for link-heavy content
		parentPropagation:    0.7, // 70% of child score propagates to parent
		removePresentational: true,
		removeHidden:         true,
	}
}

// ScoreContent performs BFS-based content scoring on the tree
func (cs *ContentScorer) ScoreContent(root *TextNode) *TextNode {
	if root == nil {
		return nil
	}

	// Step 1: Clean the tree by removing presentational/hidden nodes
	cleaned := cs.cleanTree(root)
	if cleaned == nil {
		return nil
	}

	// Step 2: Calculate initial content scores for each node
	cs.calculateInitialScores(cleaned)

	// Step 3: Propagate scores from children to parents using BFS
	cs.propagateScoresBFS(cleaned)

	// Step 4: Normalize scores for better comparison
	cs.normalizeScores(cleaned)

	return cleaned
}

// cleanTree removes presentational and hidden nodes from the tree
func (cs *ContentScorer) cleanTree(root *TextNode) *TextNode {
	if root == nil {
		return nil
	}

	// Check if this node should be removed
	if cs.shouldRemoveNode(root) {
		return nil
	}

	// Clean children and keep only valid ones
	var cleanedChildren []*TextNode
	for _, child := range root.Children {
		if cleaned := cs.cleanTree(child); cleaned != nil {
			cleaned.Parent = root
			cleanedChildren = append(cleanedChildren, cleaned)
		}
	}

	// Update children list
	root.Children = cleanedChildren

	// Update indices
	for i, child := range root.Children {
		child.Index = i
	}

	return root
}

// shouldRemoveNode determines if a node should be removed from the tree
func (cs *ContentScorer) shouldRemoveNode(node *TextNode) bool {
	if !cs.removePresentational && !cs.removeHidden {
		return false
	}

	tag := strings.ToLower(node.Tag)

	// Always remove these tags
	alwaysRemove := map[string]bool{
		"script": true, "style": true, "noscript": true,
		"meta": true, "link": true, "title": true,
	}
	if alwaysRemove[tag] {
		return true
	}

	// Remove presentational elements
	if cs.removePresentational {
		presentational := map[string]bool{
			"hr": true, "br": true, "wbr": true,
		}
		if presentational[tag] {
			return true
		}
	}

	// Remove hidden elements
	if cs.removeHidden {
		if cs.isHiddenElement(node) {
			return true
		}
	}

	// Remove empty containers with no content
	if cs.isEmptyContainer(node) {
		return true
	}

	return false
}

// isHiddenElement checks if an element is visually hidden
func (cs *ContentScorer) isHiddenElement(node *TextNode) bool {
	// Check for hidden classes
	if class, exists := node.Attributes["class"]; exists {
		hiddenClasses := []string{"hidden", "invisible", "sr-only", "screen-reader", "visually-hidden"}
		classLower := strings.ToLower(class)
		for _, hiddenClass := range hiddenClasses {
			if strings.Contains(classLower, hiddenClass) {
				return true
			}
		}
	}

	// Check for display:none or visibility:hidden in style
	if style, exists := node.Attributes["style"]; exists {
		styleLower := strings.ToLower(style)
		if strings.Contains(styleLower, "display:none") ||
			strings.Contains(styleLower, "display: none") ||
			strings.Contains(styleLower, "visibility:hidden") ||
			strings.Contains(styleLower, "visibility: hidden") {
			return true
		}
	}

	return false
}

// isEmptyContainer checks if a node is an empty container with no meaningful content
func (cs *ContentScorer) isEmptyContainer(node *TextNode) bool {
	// Text nodes are never empty containers
	if node.Tag == "#text" {
		return strings.TrimSpace(node.Text) == ""
	}

	// If it has no children, it's empty
	if len(node.Children) == 0 {
		return true
	}

	// Check if all children would be removed or are empty
	hasContent := false
	for _, child := range node.Children {
		if !cs.shouldRemoveNode(child) && !cs.isEmptyContainer(child) {
			hasContent = true
			break
		}
	}

	return !hasContent
}

// calculateInitialScores assigns initial content scores to each node
func (cs *ContentScorer) calculateInitialScores(node *TextNode) {
	if node == nil {
		return
	}

	// Calculate base content score
	node.ContentScore = cs.calculateNodeContentScore(node)

	// Recursively calculate for children
	for _, child := range node.Children {
		cs.calculateInitialScores(child)
	}
}

// calculateNodeContentScore calculates the content score for a single node
func (cs *ContentScorer) calculateNodeContentScore(node *TextNode) float64 {
	score := 0.0

	if node.Tag == "#text" {
		// Text nodes get score based on word count
		wordCount := len(strings.Fields(strings.TrimSpace(node.Text)))
		score = float64(wordCount) * cs.textContentWeight
	} else {
		// Element nodes get score based on tag type
		tag := strings.ToLower(node.Tag)

		// Heading bonus
		if weight, isHeading := cs.headingWeights[tag]; isHeading {
			score = weight
		}

		// Content semantic elements
		contentTags := map[string]float64{
			"article": 3.0, "main": 3.0, "section": 2.0,
			"p": 1.5, "blockquote": 2.0, "pre": 1.5,
			"li": 1.0, "dd": 1.0, "dt": 1.0,
		}
		if weight, isContent := contentTags[tag]; isContent {
			score = weight
		}

		// Navigation/structural penalty
		structuralTags := map[string]float64{
			"nav": -1.0, "header": -0.5, "footer": -0.5,
			"aside": -0.5, "form": -0.3,
		}
		if penalty, isStructural := structuralTags[tag]; isStructural {
			score = penalty
		}
	}

	// Apply link density penalty
	if node.LinkDensity > 0.3 {
		score *= (1.0 - cs.linkPenalty)
	}

	// Ensure minimum score of 0
	if score < 0 {
		score = 0
	}

	return score
}

// propagateScoresBFS propagates content scores from children to parents using BFS
func (cs *ContentScorer) propagateScoresBFS(root *TextNode) {
	if root == nil {
		return
	}

	// Use BFS to process nodes level by level, from bottom to top
	visited := make(map[*TextNode]bool)
	queue := list.New()

	// First, find all leaf nodes and add them to queue
	cs.findLeafNodes(root, queue)

	for queue.Len() > 0 {
		element := queue.Front()
		queue.Remove(element)
		node := element.Value.(*TextNode)

		if visited[node] {
			continue
		}
		visited[node] = true

		// Calculate total child score
		childScore := 0.0
		for _, child := range node.Children {
			childScore += child.ContentScore
		}

		// Add propagated child score to this node's own score
		node.ContentScore += childScore * cs.parentPropagation

		// Add parent to queue if all its children have been processed
		if node.Parent != nil && cs.allChildrenProcessed(node.Parent, visited) {
			queue.PushBack(node.Parent)
		}
	}
}

// findLeafNodes finds all leaf nodes (nodes with no children) and adds them to the queue
func (cs *ContentScorer) findLeafNodes(node *TextNode, queue *list.List) {
	if node == nil {
		return
	}

	if len(node.Children) == 0 {
		queue.PushBack(node)
	} else {
		for _, child := range node.Children {
			cs.findLeafNodes(child, queue)
		}
	}
}

// allChildrenProcessed checks if all children of a node have been processed
func (cs *ContentScorer) allChildrenProcessed(node *TextNode, visited map[*TextNode]bool) bool {
	for _, child := range node.Children {
		if !visited[child] {
			return false
		}
	}
	return true
}

// normalizeScores normalizes all scores in the tree to a 0-1 range
func (cs *ContentScorer) normalizeScores(root *TextNode) {
	maxScore := cs.findMaxScore(root)
	if maxScore <= 0 {
		return
	}

	cs.normalizeNodeScores(root, maxScore)
}

// findMaxScore finds the maximum content score in the tree
func (cs *ContentScorer) findMaxScore(node *TextNode) float64 {
	if node == nil {
		return 0
	}

	maxScore := node.ContentScore
	for _, child := range node.Children {
		childMax := cs.findMaxScore(child)
		if childMax > maxScore {
			maxScore = childMax
		}
	}

	return maxScore
}

// normalizeNodeScores normalizes scores for a node and its children
func (cs *ContentScorer) normalizeNodeScores(node *TextNode, maxScore float64) {
	if node == nil {
		return
	}

	node.ContentScore = node.ContentScore / maxScore

	for _, child := range node.Children {
		cs.normalizeNodeScores(child, maxScore)
	}
}

// GetTopContentNodes returns nodes with content scores above the threshold
func (cs *ContentScorer) GetTopContentNodes(root *TextNode, threshold float64) []*TextNode {
	var topNodes []*TextNode
	cs.collectTopContentNodes(root, threshold, &topNodes)
	return topNodes
}

// collectTopContentNodes recursively collects nodes above the content threshold
func (cs *ContentScorer) collectTopContentNodes(node *TextNode, threshold float64, result *[]*TextNode) {
	if node == nil {
		return
	}

	if node.ContentScore >= threshold {
		*result = append(*result, node)
	}

	for _, child := range node.Children {
		cs.collectTopContentNodes(child, threshold, result)
	}
}
