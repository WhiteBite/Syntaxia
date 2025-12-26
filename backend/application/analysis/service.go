package analysis

import (
	"syntaxia/domain"
	"sort"
	"strings"
)

// Service handles task analysis and file suggestions.
type Service struct{}

// NewService creates a new Service instance.
func NewService() *Service {
	return &Service{}
}

// SuggestFiles analyzes task description and suggests relevant files.
func (s *Service) SuggestFiles(taskDescription string, files []*domain.FileNode) []domain.SuggestedFile {
	if taskDescription == "" || len(files) == 0 {
		return []domain.SuggestedFile{}
	}

	keywords := extractTaskKeywords(taskDescription)
	if len(keywords) == 0 {
		return []domain.SuggestedFile{}
	}

	flatFiles := flattenTree(files)

	suggestions := []domain.SuggestedFile{}
	for _, file := range flatFiles {
		if file.IsDir {
			continue
		}

		score, matchedKeywords := scoreFileByKeywords(file, keywords)
		if score > 0 {
			suggestions = append(suggestions, domain.SuggestedFile{
				Path:       file.Path,
				Reason:     buildSuggestionReason(matchedKeywords),
				Confidence: score,
			})
		}
	}

	sortByConfidence(suggestions)

	if len(suggestions) > 20 {
		suggestions = suggestions[:20]
	}

	return suggestions
}

func extractTaskKeywords(text string) []string {
	text = strings.ToLower(text)

	words := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == ',' || r == '.' || r == ';' || r == ':' || r == '\n' || r == '\t'
	})

	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "from": true, "as": true, "is": true, "was": true,
		"be": true, "have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "could": true, "should": true,
		"i": true, "you": true, "we": true, "they": true, "it": true, "this": true,
		"that": true, "these": true, "those": true, "my": true, "your": true,
	}

	keywords := []string{}
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 2 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

func flattenTree(nodes []*domain.FileNode) []domain.FileNode {
	result := []domain.FileNode{}

	var traverse func([]*domain.FileNode)
	traverse = func(nodeList []*domain.FileNode) {
		for _, node := range nodeList {
			result = append(result, *node)
			if node.Children != nil {
				traverse(node.Children)
			}
		}
	}

	traverse(nodes)
	return result
}

func scoreFileByKeywords(file domain.FileNode, keywords []string) (float64, []string) {
	fileLower := strings.ToLower(file.Path)
	fileName := strings.ToLower(file.Name)

	score := 0.0
	matched := []string{}

	for _, keyword := range keywords {
		if strings.Contains(fileLower, keyword) {
			score += 0.5
			matched = append(matched, keyword)
		}

		if strings.Contains(fileName, keyword) {
			score += 1.0
			if !containsStr(matched, keyword) {
				matched = append(matched, keyword)
			}
		}

		if strings.HasSuffix(fileName, "."+keyword) {
			score += 0.3
		}
	}

	if score > 0 {
		score /= float64(len(keywords))
		if score > 1.0 {
			score = 1.0
		}
	}

	return score, matched
}

func buildSuggestionReason(keywords []string) string {
	if len(keywords) == 0 {
		return "Matches search criteria"
	}
	return "Matches keywords: " + strings.Join(keywords, ", ")
}

func sortByConfidence(suggestions []domain.SuggestedFile) {
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Confidence > suggestions[j].Confidence
	})
}

func containsStr(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
