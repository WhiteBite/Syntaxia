package context

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syntaxia/domain"
)

const (
	defaultMaxTokens   = 50000
	structureMaxTokens = 1000
	tokensPerChar      = 4 // Approximate: 1 token ≈ 4 chars
)

// SmartContextCollector собирает оптимальный контекст для AI чата
type SmartContextCollector struct {
	log         domain.Logger
	fileReader  domain.FileContentReader
	treeBuilder domain.TreeBuilder
}

// NewSmartContextCollector создает новый SmartContextCollector
func NewSmartContextCollector(
	log domain.Logger,
	fileReader domain.FileContentReader,
	treeBuilder domain.TreeBuilder,
) *SmartContextCollector {
	return &SmartContextCollector{
		log:         log,
		fileReader:  fileReader,
		treeBuilder: treeBuilder,
	}
}

// CollectContext собирает контекст для AI задачи
func (c *SmartContextCollector) CollectContext(ctx context.Context, req domain.SmartContextRequest) (*domain.SmartContextResult, error) {
	if req.MaxTokens <= 0 {
		req.MaxTokens = defaultMaxTokens
	}

	result := &domain.SmartContextResult{
		RelevantFiles: make([]domain.ContextFile, 0),
	}

	// 1. Build compact project structure
	result.ProjectStructure = c.buildCompactStructure(req.ProjectRoot)
	structureTokens := len(result.ProjectStructure) / tokensPerChar
	remainingTokens := req.MaxTokens - structureTokens

	// 2. Collect relevant files
	var files []string
	if len(req.SelectedFiles) > 0 {
		files = c.expandWithImports(req.SelectedFiles, req.ProjectRoot)
		result.Strategy = "selected_with_imports"
	} else {
		files = c.findFilesByTask(req.Task, req.ProjectRoot)
		result.Strategy = "task_analysis"
	}

	// 3. Read files and truncate to limit
	result.RelevantFiles = c.readAndTruncate(files, req.ProjectRoot, remainingTokens)

	// 4. Calculate total tokens
	result.TotalTokens = structureTokens
	for _, f := range result.RelevantFiles {
		result.TotalTokens += f.Tokens
	}

	c.log.Info(fmt.Sprintf("Smart context collected: %d files, %d tokens, strategy: %s",
		len(result.RelevantFiles), result.TotalTokens, result.Strategy))

	return result, nil
}

func (c *SmartContextCollector) buildCompactStructure(projectRoot string) string {
	var sb strings.Builder
	sb.WriteString("PROJECT STRUCTURE:\n")

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		if relPath == "." {
			return nil
		}

		// Skip common ignored directories
		name := info.Name()
		if info.IsDir() {
			if name == "node_modules" || name == ".git" || name == "vendor" ||
				name == "dist" || name == "build" || name == "__pycache__" {
				return filepath.SkipDir
			}
		}

		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > 3 {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		indent := strings.Repeat("  ", depth)
		if info.IsDir() {
			sb.WriteString(fmt.Sprintf("%s📁 %s/\n", indent, name))
		} else {
			sb.WriteString(fmt.Sprintf("%s📄 %s\n", indent, name))
		}

		if sb.Len() > structureMaxTokens*tokensPerChar {
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil && err != filepath.SkipAll {
		c.log.Warning("Error building structure: " + err.Error())
	}

	return sb.String()
}

func (c *SmartContextCollector) expandWithImports(files []string, _ string) []string {
	result := make([]string, 0, len(files)*2)
	seen := make(map[string]bool)

	for _, f := range files {
		if seen[f] {
			continue
		}
		seen[f] = true
		result = append(result, f)
	}

	// Limit to prevent explosion
	if len(result) > 30 {
		result = result[:30]
	}

	return result
}

func (c *SmartContextCollector) findFilesByTask(task string, projectRoot string) []string {
	keywords := c.extractKeywords(task)
	var matches []string

	_ = filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		name := info.Name()
		if name == "node_modules" || strings.HasPrefix(name, ".") {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		nameLower := strings.ToLower(name)

		for _, kw := range keywords {
			if strings.Contains(nameLower, kw) {
				matches = append(matches, relPath)
				break
			}
		}

		if len(matches) >= 20 {
			return filepath.SkipAll
		}

		return nil
	})

	return matches
}

func (c *SmartContextCollector) extractKeywords(task string) []string {
	task = strings.ToLower(task)

	// Common programming keywords to extract
	patterns := []string{
		`\b(auth|login|logout|register|signup)\b`,
		`\b(user|account|profile)\b`,
		`\b(button|modal|dialog|form|input)\b`,
		`\b(api|service|handler|controller)\b`,
		`\b(store|state|reducer)\b`,
		`\b(test|spec)\b`,
		`\b(config|settings)\b`,
	}

	keywords := make([]string, 0)
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		matches := re.FindAllString(task, -1)
		keywords = append(keywords, matches...)
	}

	// Also extract words > 3 chars
	words := regexp.MustCompile(`\b[a-z]{4,}\b`).FindAllString(task, -1)
	keywords = append(keywords, words...)

	// Deduplicate
	seen := make(map[string]bool)
	unique := make([]string, 0)
	for _, k := range keywords {
		if !seen[k] {
			seen[k] = true
			unique = append(unique, k)
		}
	}

	return unique
}

func (c *SmartContextCollector) readAndTruncate(files []string, projectRoot string, maxTokens int) []domain.ContextFile {
	result := make([]domain.ContextFile, 0, len(files))
	usedTokens := 0

	for _, relPath := range files {
		if usedTokens >= maxTokens {
			break
		}

		fullPath := filepath.Join(projectRoot, relPath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		contentStr := string(content)
		tokens := len(contentStr) / tokensPerChar

		// Truncate if too large
		remainingTokens := maxTokens - usedTokens
		if tokens > remainingTokens {
			maxChars := remainingTokens * tokensPerChar
			if maxChars > 0 && maxChars < len(contentStr) {
				contentStr = contentStr[:maxChars] + "\n... [truncated]"
				tokens = remainingTokens
			}
		}

		result = append(result, domain.ContextFile{
			Path:    relPath,
			Content: contentStr,
			Tokens:  tokens,
			Reason:  "selected",
		})

		usedTokens += tokens
	}

	return result
}
