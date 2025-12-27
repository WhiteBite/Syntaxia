// Package rag provides Retrieval Augmented Generation services
package rag

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"syntaxia/domain"
)

// SourceType represents the type of context source
type SourceType string

const (
	// SourceCode represents code files
	SourceCode SourceType = "code"
	// SourceGitDiff represents current git changes
	SourceGitDiff SourceType = "git_diff"
	// SourceGitHistory represents git commit history
	SourceGitHistory SourceType = "git_history"
	// SourceTests represents test files
	SourceTests SourceType = "tests"
	// SourceDocs represents documentation files
	SourceDocs SourceType = "docs"
	// SourceTerminal represents terminal output (errors)
	SourceTerminal SourceType = "terminal"
)

// SourceRequest represents a request for context from multiple sources
type SourceRequest struct {
	ProjectRoot string       `json:"projectRoot"`
	Sources     []SourceType `json:"sources"`
	Query       string       `json:"query"`
	MaxTokens   int          `json:"maxTokens"`
	// Optional filters
	FilePaths      []string `json:"filePaths,omitempty"`      // Specific files to include
	TerminalOutput string   `json:"terminalOutput,omitempty"` // Terminal output to include
	GitHistoryDays int      `json:"gitHistoryDays,omitempty"` // Days of git history (default: 7)
}

// SourceResult contains collected context from all sources
type SourceResult struct {
	Items       []SourceItem       `json:"items"`
	TotalTokens int                `json:"totalTokens"`
	Sources     map[SourceType]int `json:"sources"` // Token count per source
}

// SourceItem represents a single item from a context source
type SourceItem struct {
	Source    SourceType        `json:"source"`
	Path      string            `json:"path"`
	Content   string            `json:"content"`
	Tokens    int               `json:"tokens"`
	Relevance float64           `json:"relevance"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// ContextSourceManager collects context from multiple sources
type ContextSourceManager interface {
	// CollectFromSources collects context from specified sources
	CollectFromSources(ctx context.Context, req SourceRequest) (*SourceResult, error)
	// GetAvailableSources returns all available source types
	GetAvailableSources() []SourceType
}

// ContextSourceManagerImpl implements ContextSourceManager
type ContextSourceManagerImpl struct {
	log        domain.Logger
	fileReader domain.FileContentReader
	gitRepo    domain.GitRepository

	mu              sync.RWMutex
	terminalOutputs map[string]string // sessionID -> output
}

// NewContextSourceManager creates a new ContextSourceManager
func NewContextSourceManager(
	log domain.Logger,
	fileReader domain.FileContentReader,
	gitRepo domain.GitRepository,
) *ContextSourceManagerImpl {
	return &ContextSourceManagerImpl{
		log:             log,
		fileReader:      fileReader,
		gitRepo:         gitRepo,
		terminalOutputs: make(map[string]string),
	}
}

// GetAvailableSources returns all available source types
func (m *ContextSourceManagerImpl) GetAvailableSources() []SourceType {
	return []SourceType{
		SourceCode,
		SourceGitDiff,
		SourceGitHistory,
		SourceTests,
		SourceDocs,
		SourceTerminal,
	}
}

// CollectFromSources collects context from specified sources
func (m *ContextSourceManagerImpl) CollectFromSources(ctx context.Context, req SourceRequest) (*SourceResult, error) {
	if req.MaxTokens <= 0 {
		req.MaxTokens = 100000 // Default 100k tokens
	}

	m.log.Info(fmt.Sprintf("Collecting context from %d sources, maxTokens=%d", len(req.Sources), req.MaxTokens))

	result := &SourceResult{
		Items:   make([]SourceItem, 0),
		Sources: make(map[SourceType]int),
	}

	// Calculate token budget per source
	sourceBudget := m.calculateSourceBudgets(req.Sources, req.MaxTokens)

	// Collect from each source
	for _, source := range req.Sources {
		budget := sourceBudget[source]
		if budget <= 0 {
			continue
		}

		items, err := m.collectFromSource(ctx, source, req, budget)
		if err != nil {
			m.log.Warning(fmt.Sprintf("Failed to collect from source %s: %v", source, err))
			continue
		}

		sourceTokens := 0
		for _, item := range items {
			result.Items = append(result.Items, item)
			sourceTokens += item.Tokens
		}
		result.Sources[source] = sourceTokens
		result.TotalTokens += sourceTokens
	}

	// Sort by relevance
	sort.Slice(result.Items, func(i, j int) bool {
		return result.Items[i].Relevance > result.Items[j].Relevance
	})

	m.log.Info(fmt.Sprintf("Collected %d items, %d tokens from %d sources",
		len(result.Items), result.TotalTokens, len(result.Sources)))

	return result, nil
}

// calculateSourceBudgets distributes token budget across sources
func (m *ContextSourceManagerImpl) calculateSourceBudgets(sources []SourceType, maxTokens int) map[SourceType]int {
	budgets := make(map[SourceType]int)

	if len(sources) == 0 {
		return budgets
	}

	// Priority weights for different sources
	weights := map[SourceType]float64{
		SourceCode:       0.50, // Code gets most budget
		SourceTests:      0.15,
		SourceGitDiff:    0.15,
		SourceGitHistory: 0.10,
		SourceDocs:       0.05,
		SourceTerminal:   0.05,
	}

	totalWeight := 0.0
	for _, source := range sources {
		if w, ok := weights[source]; ok {
			totalWeight += w
		}
	}

	if totalWeight == 0 {
		// Equal distribution if no weights
		perSource := maxTokens / len(sources)
		for _, source := range sources {
			budgets[source] = perSource
		}
		return budgets
	}

	// Distribute based on weights
	for _, source := range sources {
		if w, ok := weights[source]; ok {
			budgets[source] = int(float64(maxTokens) * w / totalWeight)
		}
	}

	return budgets
}

// collectFromSource collects items from a specific source
func (m *ContextSourceManagerImpl) collectFromSource(
	ctx context.Context,
	source SourceType,
	req SourceRequest,
	budget int,
) ([]SourceItem, error) {
	switch source {
	case SourceCode:
		return m.collectCodeFiles(ctx, req, budget)
	case SourceGitDiff:
		return m.collectGitDiff(ctx, req, budget)
	case SourceGitHistory:
		return m.collectGitHistory(ctx, req, budget)
	case SourceTests:
		return m.collectTestFiles(ctx, req, budget)
	case SourceDocs:
		return m.collectDocs(ctx, req, budget)
	case SourceTerminal:
		return m.collectTerminalOutput(req, budget)
	default:
		return nil, fmt.Errorf("unknown source type: %s", source)
	}
}

// collectCodeFiles collects code files
func (m *ContextSourceManagerImpl) collectCodeFiles(ctx context.Context, req SourceRequest, budget int) ([]SourceItem, error) {
	items := make([]SourceItem, 0)
	currentTokens := 0

	// Use specified files or find relevant ones
	filePaths := req.FilePaths
	if len(filePaths) == 0 {
		// Get all tracked files
		if m.gitRepo != nil && m.gitRepo.IsGitRepository(req.ProjectRoot) {
			var err error
			filePaths, err = m.gitRepo.GetAllFiles(req.ProjectRoot)
			if err != nil {
				m.log.Warning(fmt.Sprintf("Failed to get git files: %v", err))
			}
		}
	}

	// Filter to code files only
	codeFiles := filterCodeFiles(filePaths)

	// Read and add files within budget
	for _, filePath := range codeFiles {
		if currentTokens >= budget {
			break
		}

		fullPath := filepath.Join(req.ProjectRoot, filePath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		tokens := estimateTokens(string(content))
		if currentTokens+tokens > budget {
			// Try to truncate
			availableTokens := budget - currentTokens
			if availableTokens < 100 {
				continue
			}
			content = []byte(truncateToTokens(string(content), availableTokens))
			tokens = availableTokens
		}

		items = append(items, SourceItem{
			Source:    SourceCode,
			Path:      filePath,
			Content:   string(content),
			Tokens:    tokens,
			Relevance: calculateRelevance(filePath, req.Query),
		})
		currentTokens += tokens
	}

	return items, nil
}

// collectGitDiff collects current git changes
func (m *ContextSourceManagerImpl) collectGitDiff(ctx context.Context, req SourceRequest, budget int) ([]SourceItem, error) {
	if m.gitRepo == nil || !m.gitRepo.IsGitRepository(req.ProjectRoot) {
		return nil, nil
	}

	// Get uncommitted changes
	uncommitted, err := m.gitRepo.GetUncommittedFiles(req.ProjectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to get uncommitted files: %w", err)
	}

	if len(uncommitted) == 0 {
		return nil, nil
	}

	// Get diff
	diff, err := m.gitRepo.GenerateDiff(req.ProjectRoot)
	if err != nil {
		// Try to build diff from uncommitted files
		diff = m.buildUncommittedSummary(uncommitted)
	}

	tokens := estimateTokens(diff)
	if tokens > budget {
		diff = truncateToTokens(diff, budget)
		tokens = budget
	}

	items := []SourceItem{
		{
			Source:    SourceGitDiff,
			Path:      "git_diff",
			Content:   diff,
			Tokens:    tokens,
			Relevance: 0.9, // High relevance for current changes
			Metadata: map[string]string{
				"changedFiles": fmt.Sprintf("%d", len(uncommitted)),
			},
		},
	}

	return items, nil
}

// buildUncommittedSummary builds a summary of uncommitted files
func (m *ContextSourceManagerImpl) buildUncommittedSummary(files []domain.FileStatus) string {
	var sb strings.Builder
	sb.WriteString("# Uncommitted Changes\n\n")

	for _, f := range files {
		status := "modified"
		switch f.Status {
		case "A":
			status = "added"
		case "D":
			status = "deleted"
		case "U":
			status = "untracked"
		}
		sb.WriteString(fmt.Sprintf("- %s: %s\n", f.Path, status))
	}

	return sb.String()
}

// collectGitHistory collects git commit history
func (m *ContextSourceManagerImpl) collectGitHistory(ctx context.Context, req SourceRequest, budget int) ([]SourceItem, error) {
	if m.gitRepo == nil || !m.gitRepo.IsGitRepository(req.ProjectRoot) {
		return nil, nil
	}

	limit := 20
	if req.GitHistoryDays > 0 {
		limit = req.GitHistoryDays * 3 // Approximate commits per day
	}

	commits, err := m.gitRepo.GetCommitHistory(req.ProjectRoot, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit history: %w", err)
	}

	if len(commits) == 0 {
		return nil, nil
	}

	// Format commit history
	var sb strings.Builder
	sb.WriteString("# Recent Git History\n\n")

	for _, commit := range commits {
		sb.WriteString(fmt.Sprintf("## %s\n", commit.Hash[:8]))
		sb.WriteString(fmt.Sprintf("- **Subject**: %s\n", commit.Subject))
		sb.WriteString(fmt.Sprintf("- **Author**: %s\n", commit.Author))
		sb.WriteString(fmt.Sprintf("- **Date**: %s\n\n", commit.Date))
	}

	content := sb.String()
	tokens := estimateTokens(content)
	if tokens > budget {
		content = truncateToTokens(content, budget)
		tokens = budget
	}

	items := []SourceItem{
		{
			Source:    SourceGitHistory,
			Path:      "git_history",
			Content:   content,
			Tokens:    tokens,
			Relevance: 0.6,
			Metadata: map[string]string{
				"commits": fmt.Sprintf("%d", len(commits)),
			},
		},
	}

	return items, nil
}

// collectTestFiles collects test files
func (m *ContextSourceManagerImpl) collectTestFiles(ctx context.Context, req SourceRequest, budget int) ([]SourceItem, error) {
	items := make([]SourceItem, 0)
	currentTokens := 0

	// Find test files
	testFiles := make([]string, 0)
	err := filepath.Walk(req.ProjectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if name == "node_modules" || name == ".git" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, _ := filepath.Rel(req.ProjectRoot, path)
		relPath = filepath.ToSlash(relPath)

		if isTestFile(relPath) {
			testFiles = append(testFiles, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	// Read test files within budget
	for _, filePath := range testFiles {
		if currentTokens >= budget {
			break
		}

		fullPath := filepath.Join(req.ProjectRoot, filePath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		tokens := estimateTokens(string(content))
		if currentTokens+tokens > budget {
			availableTokens := budget - currentTokens
			if availableTokens < 100 {
				continue
			}
			content = []byte(truncateToTokens(string(content), availableTokens))
			tokens = availableTokens
		}

		items = append(items, SourceItem{
			Source:    SourceTests,
			Path:      filePath,
			Content:   string(content),
			Tokens:    tokens,
			Relevance: calculateRelevance(filePath, req.Query),
		})
		currentTokens += tokens
	}

	return items, nil
}

// collectDocs collects documentation files
func (m *ContextSourceManagerImpl) collectDocs(ctx context.Context, req SourceRequest, budget int) ([]SourceItem, error) {
	items := make([]SourceItem, 0)
	currentTokens := 0

	// Find documentation files
	docFiles := make([]string, 0)
	err := filepath.Walk(req.ProjectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if name == "node_modules" || name == ".git" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, _ := filepath.Rel(req.ProjectRoot, path)
		relPath = filepath.ToSlash(relPath)

		if isDocFile(relPath) {
			docFiles = append(docFiles, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	// Prioritize README files
	sort.Slice(docFiles, func(i, j int) bool {
		iReadme := strings.Contains(strings.ToLower(docFiles[i]), "readme")
		jReadme := strings.Contains(strings.ToLower(docFiles[j]), "readme")
		if iReadme != jReadme {
			return iReadme
		}
		return docFiles[i] < docFiles[j]
	})

	// Read doc files within budget
	for _, filePath := range docFiles {
		if currentTokens >= budget {
			break
		}

		fullPath := filepath.Join(req.ProjectRoot, filePath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		tokens := estimateTokens(string(content))
		if currentTokens+tokens > budget {
			availableTokens := budget - currentTokens
			if availableTokens < 100 {
				continue
			}
			content = []byte(truncateToTokens(string(content), availableTokens))
			tokens = availableTokens
		}

		items = append(items, SourceItem{
			Source:    SourceDocs,
			Path:      filePath,
			Content:   string(content),
			Tokens:    tokens,
			Relevance: calculateDocRelevance(filePath),
		})
		currentTokens += tokens
	}

	return items, nil
}

// collectTerminalOutput collects terminal output
func (m *ContextSourceManagerImpl) collectTerminalOutput(req SourceRequest, budget int) ([]SourceItem, error) {
	if req.TerminalOutput == "" {
		return nil, nil
	}

	content := req.TerminalOutput
	tokens := estimateTokens(content)
	if tokens > budget {
		content = truncateToTokens(content, budget)
		tokens = budget
	}

	items := []SourceItem{
		{
			Source:    SourceTerminal,
			Path:      "terminal_output",
			Content:   content,
			Tokens:    tokens,
			Relevance: 0.95, // Very high relevance for error output
			Metadata: map[string]string{
				"type": "error_output",
			},
		},
	}

	return items, nil
}

// SetTerminalOutput stores terminal output for a session
func (m *ContextSourceManagerImpl) SetTerminalOutput(sessionID, output string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.terminalOutputs[sessionID] = output
}

// GetTerminalOutput retrieves terminal output for a session
func (m *ContextSourceManagerImpl) GetTerminalOutput(sessionID string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.terminalOutputs[sessionID]
}

// ClearTerminalOutput clears terminal output for a session
func (m *ContextSourceManagerImpl) ClearTerminalOutput(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.terminalOutputs, sessionID)
}


// Helper functions

// filterCodeFiles filters file paths to only include code files
func filterCodeFiles(paths []string) []string {
	codeExtensions := map[string]bool{
		".go": true, ".ts": true, ".tsx": true, ".js": true, ".jsx": true,
		".py": true, ".java": true, ".kt": true, ".rs": true, ".c": true,
		".cpp": true, ".h": true, ".hpp": true, ".cs": true, ".rb": true,
		".php": true, ".swift": true, ".vue": true, ".dart": true,
	}

	result := make([]string, 0, len(paths))
	for _, path := range paths {
		ext := strings.ToLower(filepath.Ext(path))
		if codeExtensions[ext] && !isTestFile(path) {
			result = append(result, path)
		}
	}
	return result
}

// isTestFile checks if a file is a test file
func isTestFile(path string) bool {
	lower := strings.ToLower(path)
	return strings.Contains(lower, "_test.") ||
		strings.Contains(lower, ".test.") ||
		strings.Contains(lower, ".spec.") ||
		strings.Contains(lower, "/test/") ||
		strings.Contains(lower, "/tests/") ||
		strings.Contains(lower, "/__tests__/")
}

// isDocFile checks if a file is a documentation file
func isDocFile(path string) bool {
	lower := strings.ToLower(path)
	ext := filepath.Ext(lower)

	if ext == ".md" || ext == ".txt" || ext == ".rst" {
		return true
	}

	baseName := strings.ToLower(filepath.Base(path))
	docNames := []string{"readme", "changelog", "contributing", "license", "authors", "history"}
	for _, name := range docNames {
		if strings.Contains(baseName, name) {
			return true
		}
	}

	return strings.Contains(lower, "/docs/") || strings.Contains(lower, "/doc/")
}

// estimateTokens estimates the number of tokens in text
func estimateTokens(text string) int {
	// Approximate: 1 token ≈ 4 characters
	return len(text) / 4
}

// truncateToTokens truncates text to fit within token limit
func truncateToTokens(text string, maxTokens int) string {
	maxChars := maxTokens * 4
	if len(text) <= maxChars {
		return text
	}

	truncated := text[:maxChars]
	lastNewline := strings.LastIndex(truncated, "\n")
	if lastNewline > maxChars/2 {
		return truncated[:lastNewline] + "\n// ... truncated ..."
	}

	return truncated + "\n// ... truncated ..."
}

// calculateRelevance calculates relevance score for a file
func calculateRelevance(filePath, query string) float64 {
	if query == "" {
		return 0.5
	}

	score := 0.5
	pathLower := strings.ToLower(filePath)
	queryLower := strings.ToLower(query)

	// Check for query terms in path
	queryTerms := strings.Fields(queryLower)
	for _, term := range queryTerms {
		if strings.Contains(pathLower, term) {
			score += 0.15
		}
	}

	// Boost for main/index files
	baseName := strings.ToLower(filepath.Base(filePath))
	if strings.Contains(baseName, "main") || strings.Contains(baseName, "index") {
		score += 0.1
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateDocRelevance calculates relevance for documentation files
func calculateDocRelevance(filePath string) float64 {
	lower := strings.ToLower(filePath)

	if strings.Contains(lower, "readme") {
		return 0.9
	}
	if strings.Contains(lower, "contributing") || strings.Contains(lower, "changelog") {
		return 0.7
	}
	return 0.5
}

// FormatSourceResult formats a SourceResult for display
func FormatSourceResult(result *SourceResult) string {
	var sb strings.Builder

	// Group items by source
	bySource := make(map[SourceType][]SourceItem)
	for _, item := range result.Items {
		bySource[item.Source] = append(bySource[item.Source], item)
	}

	// Format each source section
	sourceOrder := []SourceType{SourceTerminal, SourceGitDiff, SourceCode, SourceTests, SourceDocs, SourceGitHistory}
	sourceNames := map[SourceType]string{
		SourceCode:       "Code Files",
		SourceGitDiff:    "Git Changes",
		SourceGitHistory: "Git History",
		SourceTests:      "Test Files",
		SourceDocs:       "Documentation",
		SourceTerminal:   "Terminal Output",
	}

	for _, source := range sourceOrder {
		items, ok := bySource[source]
		if !ok || len(items) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("# %s\n\n", sourceNames[source]))

		for _, item := range items {
			if item.Path != "" && item.Path != "terminal_output" && item.Path != "git_diff" && item.Path != "git_history" {
				sb.WriteString(fmt.Sprintf("## %s\n", item.Path))
			}
			sb.WriteString(item.Content)
			sb.WriteString("\n\n")
		}
	}

	return sb.String()
}
