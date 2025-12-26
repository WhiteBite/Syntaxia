package analysis

import (
	"context"
	"fmt"
	"regexp"
	"syntaxia/domain"
	"strings"
)

// KeywordAnalyzer is a simple ContextAnalyzer implementation for backward compatibility.
type KeywordAnalyzer struct {
	log domain.Logger
}

// NewKeywordAnalyzer creates a new KeywordAnalyzer instance.
func NewKeywordAnalyzer(log domain.Logger) *KeywordAnalyzer {
	return &KeywordAnalyzer{log: log}
}

// SuggestFiles: O(N) by file count — single regex for keywords.
func (a *KeywordAnalyzer) SuggestFiles(ctx context.Context, task string, allFiles []*domain.FileNode) ([]string, error) {
	keywords := a.extractKeywordsFromTask(task)
	if len(keywords) == 0 {
		a.log.Warning("Не удалось извлечь ключевые слова из задачи для авто-выбора.")
		return []string{}, nil
	}

	a.log.Info("Извлеченные ключевые слова для поиска: " + strings.Join(keywords, ", "))

	type pair struct {
		orig  string
		lower string
	}
	var files []pair

	var traverse func([]*domain.FileNode)
	traverse = func(nodes []*domain.FileNode) {
		for _, n := range nodes {
			if !n.IsDir {
				orig := n.RelPath
				files = append(files, pair{
					orig:  orig,
					lower: strings.ToLower(strings.ReplaceAll(orig, "\\", "/")),
				})
			}
			if len(n.Children) > 0 {
				traverse(n.Children)
			}
		}
	}
	traverse(allFiles)

	parts := make([]string, 0, len(keywords))
	for _, k := range keywords {
		parts = append(parts, regexp.QuoteMeta(strings.ToLower(k)))
	}
	pattern := "(?i)(" + strings.Join(parts, "|") + ")"
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	unique := make(map[string]struct{})
	for _, f := range files {
		if re.MatchString(f.lower) {
			unique[f.orig] = struct{}{}
		}
	}

	res := make([]string, 0, len(unique))
	for k := range unique {
		res = append(res, k)
	}
	a.log.Info(fmt.Sprintf("Предложено %d релевантных файлов.", len(res)))
	return res, nil
}

func (a *KeywordAnalyzer) extractKeywordsFromTask(task string) []string {
	re := regexp.MustCompile(`[^\w\s]`)
	cleanTask := re.ReplaceAllString(task, " ")
	words := strings.Fields(strings.ToLower(cleanTask))
	var keywords []string
	for _, word := range words {
		if len(word) > 3 {
			keywords = append(keywords, word)
		}
	}
	return keywords
}
