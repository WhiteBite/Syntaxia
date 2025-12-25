package rag

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"syntaxia/domain"
)

// CallStackAnalyzerInterface defines the interface for call stack analysis
type CallStackAnalyzerInterface interface {
	AnalyzeCallStack(ctx context.Context, projectRoot, filePath, symbolName string, maxDepth int) (*CallStackResult, error)
	GetTransitiveDependencies(ctx context.Context, projectRoot, filePath, symbolName string, maxDepth int) ([]*domain.SymbolNode, error)
}

// CallStackResult mirrors the infrastructure type for decoupling
type CallStackResult struct {
	RootSymbol   *domain.SymbolNode
	Callers      []*domain.SymbolNode
	Callees      []*domain.SymbolNode
	Dependencies []*domain.SymbolNode
	RelatedFiles []string
	TotalSymbols int
}

// SmartContextRequest represents a request for smart context collection
type SmartContextRequest struct {
	ProjectRoot   string   `json:"projectRoot"`
	Task          string   `json:"task"`          // User's task description
	SelectedFiles []string `json:"selectedFiles"` // Files explicitly selected by user
	SelectedCode  string   `json:"selectedCode"`  // Code snippet selected by user
	SourceFile    string   `json:"sourceFile"`    // File where selection was made
	MaxTokens     int      `json:"maxTokens"`     // Maximum tokens for context (default: 900000 for Qwen)
	MaxDepth      int      `json:"maxDepth"`      // Max depth for call stack traversal
	Language      string   `json:"language"`      // Programming language
}

// SmartContextResult contains the collected context
type SmartContextResult struct {
	Context         string               `json:"context"`         // Full context string
	Files           []ContextFile        `json:"files"`           // Files included in context
	Symbols         []*domain.SymbolNode `json:"symbols"`         // Symbols analyzed
	CallStack       *CallStackResult     `json:"callStack"`       // Call stack analysis
	TokenEstimate   int                  `json:"tokenEstimate"`   // Estimated token count
	TruncatedFiles  []string             `json:"truncatedFiles"`  // Files that were truncated
	ExcludedFiles   []string             `json:"excludedFiles"`   // Files excluded due to token limit
	RelevanceScores map[string]float64   `json:"relevanceScores"` // File relevance scores
}

// ContextFile represents a file in the context
type ContextFile struct {
	Path      string  `json:"path"`
	Content   string  `json:"content"`
	Tokens    int     `json:"tokens"`
	Relevance float64 `json:"relevance"`
	Reason    string  `json:"reason"` // Why this file was included
}

// SymbolGraphServiceInterface defines the interface for symbol graph operations
type SymbolGraphServiceInterface interface {
	// Add methods as needed
}

// SmartContextService collects relevant context for AI tasks based on code analysis
type SmartContextService struct {
	log               domain.Logger
	fileReader        domain.FileContentReader
	symbolGraphSvc    SymbolGraphServiceInterface
	callStackAnalyzer CallStackAnalyzerInterface
}

// NewSmartContextService creates a new smart context service
func NewSmartContextService(
	log domain.Logger,
	fileReader domain.FileContentReader,
	symbolGraphSvc SymbolGraphServiceInterface,
	callStackAnalyzer CallStackAnalyzerInterface,
) *SmartContextService {
	return &SmartContextService{
		log:               log,
		fileReader:        fileReader,
		symbolGraphSvc:    symbolGraphSvc,
		callStackAnalyzer: callStackAnalyzer,
	}
}

// setRequestDefaults sets default values for the request
func (s *SmartContextService) setRequestDefaults(req *SmartContextRequest) {
	if req.MaxTokens == 0 {
		req.MaxTokens = 900000
	}
	if req.MaxDepth == 0 {
		req.MaxDepth = 3
	}
	if req.Language == "" {
		req.Language = s.detectLanguage(req.ProjectRoot, req.SelectedFiles)
	}
}

// analyzeSelectedCode analyzes selected code and returns call stack result
func (s *SmartContextService) analyzeSelectedCode(ctx context.Context, req SmartContextRequest) *CallStackResult {
	if req.SelectedCode == "" || req.SourceFile == "" {
		return nil
	}

	var callStackResult *CallStackResult
	symbols := s.extractSymbolsFromCode(req.SelectedCode)

	for _, symbolName := range symbols {
		csr, err := s.callStackAnalyzer.AnalyzeCallStack(ctx, req.ProjectRoot, req.SourceFile, symbolName, req.MaxDepth)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Failed to analyze call stack for %s: %v", symbolName, err))
			continue
		}
		if callStackResult == nil {
			callStackResult = csr
		} else {
			callStackResult.Callers = append(callStackResult.Callers, csr.Callers...)
			callStackResult.Callees = append(callStackResult.Callees, csr.Callees...)
			callStackResult.Dependencies = append(callStackResult.Dependencies, csr.Dependencies...)
			callStackResult.RelatedFiles = append(callStackResult.RelatedFiles, csr.RelatedFiles...)
		}
	}
	return callStackResult
}

// contextBuildState holds state during context building
type contextBuildState struct {
	builder       *strings.Builder
	currentTokens int
	result        *SmartContextResult
}

// addFileToContext adds a file to the context, handling truncation if needed
func (s *SmartContextService) addFileToContext(ctx context.Context, file *ContextFile, req SmartContextRequest, state *contextBuildState) bool {
	contents, err := s.fileReader.ReadContents(ctx, []string{file.Path}, req.ProjectRoot, nil)
	if err != nil {
		s.log.Warning(fmt.Sprintf("Failed to read file %s: %v", file.Path, err))
		return false
	}

	content, ok := contents[file.Path]
	if !ok {
		return false
	}

	fileTokens := s.estimateTokens(content)
	fileHeader := fmt.Sprintf("## %s\n```%s\n", file.Path, s.getFileExtension(file.Path))
	fileFooter := "\n```\n\n"
	headerTokens := s.estimateTokens(fileHeader + fileFooter)

	if state.currentTokens+fileTokens+headerTokens > req.MaxTokens {
		availableTokens := req.MaxTokens - state.currentTokens - headerTokens - 100
		if availableTokens > 500 {
			truncatedContent := s.truncateToTokens(content, availableTokens)
			file.Content = truncatedContent
			file.Tokens = s.estimateTokens(truncatedContent)
			state.result.TruncatedFiles = append(state.result.TruncatedFiles, file.Path)
			state.builder.WriteString(fileHeader + truncatedContent + "\n// ... truncated ...\n" + fileFooter)
			state.currentTokens += file.Tokens + headerTokens
		} else {
			state.result.ExcludedFiles = append(state.result.ExcludedFiles, file.Path)
			return false
		}
	} else {
		file.Content = content
		file.Tokens = fileTokens
		state.builder.WriteString(fileHeader + content + fileFooter)
		state.currentTokens += fileTokens + headerTokens
	}

	state.result.Files = append(state.result.Files, *file)
	state.result.RelevanceScores[file.Path] = file.Relevance
	return true
}

// CollectContext collects smart context based on the request
func (s *SmartContextService) CollectContext(ctx context.Context, req SmartContextRequest) (*SmartContextResult, error) {
	s.log.Info(fmt.Sprintf("Collecting smart context for task: %s", domain.TruncateString(req.Task, 50)))
	s.setRequestDefaults(&req)

	result := &SmartContextResult{
		Files: make([]ContextFile, 0), Symbols: make([]*domain.SymbolNode, 0),
		TruncatedFiles: make([]string, 0), ExcludedFiles: make([]string, 0),
		RelevanceScores: make(map[string]float64),
	}

	callStackResult := s.analyzeSelectedCode(ctx, req)
	result.CallStack = callStackResult

	filesToInclude := s.collectRelevantFiles(ctx, req, callStackResult)
	sort.Slice(filesToInclude, func(i, j int) bool { return filesToInclude[i].Relevance > filesToInclude[j].Relevance })

	var contextBuilder strings.Builder
	taskSection := fmt.Sprintf("# Task\n%s\n\n", req.Task)
	contextBuilder.WriteString(taskSection)
	currentTokens := s.estimateTokens(taskSection)

	if req.SelectedCode != "" {
		selectedSection := fmt.Sprintf("# Selected Code (from %s)\n```\n%s\n```\n\n", req.SourceFile, req.SelectedCode)
		contextBuilder.WriteString(selectedSection)
		currentTokens += s.estimateTokens(selectedSection)
	}

	contextBuilder.WriteString("# Project Files\n\n")
	currentTokens += s.estimateTokens("# Project Files\n\n")

	state := &contextBuildState{builder: &contextBuilder, currentTokens: currentTokens, result: result}
	for i := range filesToInclude {
		s.addFileToContext(ctx, &filesToInclude[i], req, state)
	}

	result.Context = contextBuilder.String()
	result.TokenEstimate = state.currentTokens

	if callStackResult != nil {
		result.Symbols = append(result.Symbols, callStackResult.Callers...)
		result.Symbols = append(result.Symbols, callStackResult.Callees...)
		result.Symbols = append(result.Symbols, callStackResult.Dependencies...)
	}

	s.log.Info(fmt.Sprintf("Smart context collected: %d files, ~%d tokens", len(result.Files), result.TokenEstimate))
	return result, nil
}

// collectRelevantFiles collects files relevant to the task
func (s *SmartContextService) collectRelevantFiles(
	_ context.Context,
	req SmartContextRequest,
	callStack *CallStackResult,
) []ContextFile {
	collector := &fileCollector{files: make([]ContextFile, 0), seen: make(map[string]bool)}

	collector.addFiles(req.SelectedFiles, 1.0, "explicitly selected")
	collector.addFile(req.SourceFile, 0.95, "source of selected code")
	s.addCallStackFiles(collector, callStack)

	return collector.files
}

// fileCollector helps collect unique files with relevance
type fileCollector struct {
	files []ContextFile
	seen  map[string]bool
}

func (c *fileCollector) addFile(path string, relevance float64, reason string) {
	if path == "" || c.seen[path] {
		return
	}
	c.seen[path] = true
	c.files = append(c.files, ContextFile{Path: path, Relevance: relevance, Reason: reason})
}

func (c *fileCollector) addFiles(paths []string, relevance float64, reason string) {
	for _, path := range paths {
		c.addFile(path, relevance, reason)
	}
}

func (c *fileCollector) addSymbolFiles(symbols []*domain.SymbolNode, relevance float64, reasonFmt string) {
	for _, sym := range symbols {
		c.addFile(sym.Path, relevance, fmt.Sprintf(reasonFmt, sym.Name))
	}
}

// addCallStackFiles adds files from call stack analysis
func (s *SmartContextService) addCallStackFiles(collector *fileCollector, callStack *CallStackResult) {
	if callStack == nil {
		return
	}
	collector.addFiles(callStack.RelatedFiles, 0.8, "call stack related")
	collector.addSymbolFiles(callStack.Callers, 0.7, "contains caller: %s")
	collector.addSymbolFiles(callStack.Callees, 0.75, "contains callee: %s")
	collector.addSymbolFiles(callStack.Dependencies, 0.6, "contains dependency: %s")
}

// extractSymbolsFromCode extracts symbol names from code snippet
func (s *SmartContextService) extractSymbolsFromCode(code string) []string {
	symbols := make([]string, 0)
	seen := make(map[string]bool)

	// Simple extraction: look for function/method names
	// This is a basic implementation - could be enhanced with proper parsing
	words := strings.FieldsFunc(code, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_')
	})

	for _, word := range words {
		// Skip common keywords and short words
		if len(word) < 3 || isCommonKeyword(word) {
			continue
		}
		// Look for PascalCase or camelCase (likely function/type names)
		if (word[0] >= 'A' && word[0] <= 'Z') || strings.Contains(word, "_") {
			if !seen[word] {
				seen[word] = true
				symbols = append(symbols, word)
			}
		}
	}

	return symbols
}

// detectLanguage detects the primary language of the project
func (s *SmartContextService) detectLanguage(_ string, files []string) string {
	extCount := make(map[string]int)

	for _, f := range files {
		ext := filepath.Ext(f)
		extCount[ext]++
	}

	// Map extensions to languages
	langMap := map[string]string{
		".go":   "go",
		".ts":   "typescript",
		".tsx":  "typescript",
		".js":   "javascript",
		".jsx":  "javascript",
		".py":   "python",
		".java": "java",
		".rs":   "rust",
		".cpp":  "cpp",
		".c":    "c",
	}

	maxCount := 0
	detectedLang := "go" // default

	for ext, count := range extCount {
		if count > maxCount {
			if lang, ok := langMap[ext]; ok {
				maxCount = count
				detectedLang = lang
			}
		}
	}

	return detectedLang
}

// estimateTokens estimates token count for text
func (s *SmartContextService) estimateTokens(text string) int {
	// Approximate: 1 token ≈ 4 characters
	return len(text) / 4
}

// truncateToTokens truncates text to fit within token limit
func (s *SmartContextService) truncateToTokens(text string, maxTokens int) string {
	maxChars := maxTokens * 4
	if len(text) <= maxChars {
		return text
	}

	// Try to truncate at a line boundary
	truncated := text[:maxChars]
	lastNewline := strings.LastIndex(truncated, "\n")
	if lastNewline > maxChars/2 {
		return truncated[:lastNewline]
	}

	return truncated
}

// getFileExtension returns the file extension for syntax highlighting
func (s *SmartContextService) getFileExtension(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return ""
	}
	return ext[1:] // Remove the dot
}

// Helper functions

func isCommonKeyword(word string) bool {
	keywords := map[string]bool{
		"func": true, "return": true, "if": true, "else": true,
		"for": true, "range": true, "var": true, "const": true,
		"type": true, "struct": true, "interface": true, "package": true,
		"import": true, "nil": true, "true": true, "false": true,
		"string": true, "int": true, "bool": true, "error": true,
		"make": true, "new": true, "len": true, "append": true,
		"function": true, "class": true, "this": true, "self": true,
		"def": true, "async": true, "await": true, "export": true,
	}
	return keywords[strings.ToLower(word)]
}
