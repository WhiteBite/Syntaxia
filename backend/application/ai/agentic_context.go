// Package ai provides AI service functionality including agentic context gathering.
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"syntaxia/domain"
)

const (
	// DefaultMaxIterations is the default maximum number of context gathering iterations
	DefaultMaxIterations = 5
	// DefaultIterationTimeout is the timeout for each iteration
	DefaultIterationTimeout = 30 * time.Second
	// DefaultMaxContextTokens is the default token budget for context
	DefaultMaxContextTokens = 100000
)

// AgenticContextService performs multi-step context gathering using LLM with tools.
// Similar to Sourcegraph Cody's approach - LLM analyzes the task and iteratively
// gathers relevant files using search and symbol tools.
type AgenticContextService struct {
	aiProvider  domain.AIProvider
	fileTools   FileToolsProvider
	symbolTools SymbolToolsProvider
	log         domain.Logger
}

// FileToolsProvider interface for file operations
type FileToolsProvider interface {
	SearchFiles(pattern, directory, projectRoot string) ([]string, error)
	SearchContent(pattern, filePattern, projectRoot string, maxResults int) ([]ContentMatch, error)
	ReadFile(path, projectRoot string) (string, error)
}

// SymbolToolsProvider interface for symbol operations
type SymbolToolsProvider interface {
	ListSymbols(path, projectRoot string) ([]SymbolInfo, error)
	SearchSymbols(query, kind, projectRoot string) ([]SymbolInfo, error)
}

// ContentMatch represents a content search match
type ContentMatch struct {
	FilePath string `json:"filePath"`
	Line     int    `json:"line"`
	Content  string `json:"content"`
}

// SymbolInfo represents symbol information
type SymbolInfo struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	FilePath  string `json:"filePath"`
	Line      int    `json:"line"`
	Signature string `json:"signature,omitempty"`
}

// AgenticContextRequest represents a request for agentic context gathering
type AgenticContextRequest struct {
	Task          string `json:"task"`
	ProjectRoot   string `json:"projectRoot"`
	MaxIterations int    `json:"maxIterations"` // default 5
	MaxTokens     int    `json:"maxTokens"`     // token budget
}

// AgenticContextResult represents the result of context gathering
type AgenticContextResult struct {
	Files       []ContextFileResult `json:"files"`
	Iterations  int                 `json:"iterations"`
	ToolCalls   []ToolCallLog       `json:"toolCalls"`
	TotalTokens int                 `json:"totalTokens"`
	Reasoning   string              `json:"reasoning,omitempty"`
}

// ContextFileResult represents a file selected for context
type ContextFileResult struct {
	Path      string  `json:"path"`
	Reason    string  `json:"reason"`
	Relevance float64 `json:"relevance"`
	Tokens    int     `json:"tokens"`
}

// NewAgenticContextService creates a new AgenticContextService
func NewAgenticContextService(
	aiProvider domain.AIProvider,
	fileTools FileToolsProvider,
	symbolTools SymbolToolsProvider,
	log domain.Logger,
) *AgenticContextService {
	return &AgenticContextService{
		aiProvider:  aiProvider,
		fileTools:   fileTools,
		symbolTools: symbolTools,
		log:         log,
	}
}

// GatherContext performs multi-step context gathering for a task
func (s *AgenticContextService) GatherContext(ctx context.Context, req AgenticContextRequest) (*AgenticContextResult, error) {
	if req.MaxIterations <= 0 {
		req.MaxIterations = DefaultMaxIterations
	}
	if req.MaxTokens <= 0 {
		req.MaxTokens = DefaultMaxContextTokens
	}

	s.log.Info(fmt.Sprintf("Starting agentic context gathering for task: %s", domain.TruncateString(req.Task, 100)))

	result := &AgenticContextResult{
		Files:     make([]ContextFileResult, 0),
		ToolCalls: make([]ToolCallLog, 0),
	}

	// Track requested files to avoid duplicates
	requestedFiles := make(map[string]bool)

	// Build system prompt with available tools
	systemPrompt := s.buildSystemPrompt()

	// Initialize conversation
	messages := []conversationMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: fmt.Sprintf("Task: %s\n\nProject root: %s", req.Task, req.ProjectRoot)},
	}

	for iteration := 1; iteration <= req.MaxIterations; iteration++ {
		result.Iterations = iteration
		s.log.Info(fmt.Sprintf("Context gathering iteration %d/%d", iteration, req.MaxIterations))

		// Create iteration context with timeout
		iterCtx, cancel := context.WithTimeout(ctx, DefaultIterationTimeout)

		// Call AI to decide next action
		response, err := s.callAI(iterCtx, messages)
		cancel()

		if err != nil {
			return nil, fmt.Errorf("AI call failed at iteration %d: %w", iteration, err)
		}

		// Parse tool calls from response
		toolCalls := s.parseToolCalls(response)

		// Check if context gathering is complete
		if s.isContextComplete(toolCalls) {
			s.log.Info("Context gathering completed by LLM decision")
			result.Reasoning = s.extractReasoning(response)
			break
		}

		// No tool calls means final answer
		if len(toolCalls) == 0 {
			s.log.Info("Context gathering completed - no more tool calls")
			result.Reasoning = response
			break
		}

		// Add assistant response to conversation
		messages = append(messages, conversationMessage{Role: "assistant", Content: response})

		// Execute tool calls
		var toolResults []string
		for _, call := range toolCalls {
			toolResult, err := s.executeTool(call, req.ProjectRoot, requestedFiles, result)
			if err != nil {
				s.log.Warning(fmt.Sprintf("Tool %s failed: %v", call.Name, err))
				toolResult = fmt.Sprintf("Error: %v", err)
			}

			argsJSON, _ := json.Marshal(call.Arguments)
			result.ToolCalls = append(result.ToolCalls, ToolCallLog{
				Tool:      call.Name,
				Arguments: string(argsJSON),
				Result:    domain.TruncateString(toolResult, 500),
			})

			toolResults = append(toolResults, fmt.Sprintf("Tool: %s\nResult:\n%s", call.Name, toolResult))
		}

		// Add tool results to conversation
		messages = append(messages, conversationMessage{Role: "tool", Content: strings.Join(toolResults, "\n\n---\n\n")})

		// Check token budget
		if result.TotalTokens >= req.MaxTokens {
			s.log.Info(fmt.Sprintf("Token budget reached: %d/%d", result.TotalTokens, req.MaxTokens))
			break
		}
	}

	s.log.Info(fmt.Sprintf("Context gathering completed: %d files, %d iterations, %d tokens",
		len(result.Files), result.Iterations, result.TotalTokens))

	return result, nil
}

func (s *AgenticContextService) buildSystemPrompt() string {
	return `You are a context gathering assistant. Your task is to find relevant files for the user's coding task.

Available tools:
- search_files_for_context(pattern): Search files by name pattern (glob or partial match)
  Parameters: {"pattern": "string", "directory": "string (optional)"}

- search_content_for_context(query): Search file contents (grep-like)
  Parameters: {"query": "string", "file_pattern": "string (optional)", "max_results": "number (optional, default 20)"}

- get_file_symbols(path): Get symbols (functions, classes, types) in a file
  Parameters: {"path": "string"}

- request_file(path, reason): Add a file to the context with explanation
  Parameters: {"path": "string", "reason": "string", "relevance": "number 0-1 (optional)"}

- mark_context_complete(reasoning): Signal that context gathering is complete
  Parameters: {"reasoning": "string"}

Process:
1. Analyze the user's task to understand what code areas are relevant
2. Search for files by name patterns related to the task
3. Search file contents for relevant keywords, function names, etc.
4. Examine file symbols to understand structure
5. Request files that are relevant with clear reasons
6. Mark complete when you have sufficient context

Guidelines:
- Be thorough but efficient - don't request unnecessary files
- Prioritize files directly related to the task
- Include related files (imports, interfaces, tests) when relevant
- Stop when you have enough context to understand the task area
- Explain your reasoning for file selections

Respond with JSON containing tool_calls:
{"tool_calls": [{"name": "tool_name", "arguments": {"arg1": "value1"}}]}

When done, call mark_context_complete with your reasoning.`
}

type conversationMessage struct {
	Role    string
	Content string
}

func (s *AgenticContextService) callAI(ctx context.Context, messages []conversationMessage) (string, error) {
	var systemPrompt string
	var userPrompt strings.Builder

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			systemPrompt = msg.Content
		case "user":
			userPrompt.WriteString("User: ")
			userPrompt.WriteString(msg.Content)
			userPrompt.WriteString("\n\n")
		case "assistant":
			userPrompt.WriteString("Assistant: ")
			userPrompt.WriteString(msg.Content)
			userPrompt.WriteString("\n\n")
		case "tool":
			userPrompt.WriteString("Tool Results:\n")
			userPrompt.WriteString(msg.Content)
			userPrompt.WriteString("\n\n")
		}
	}
	userPrompt.WriteString("Assistant: ")

	aiReq := domain.AIRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt.String(),
		Temperature:  0.3, // Lower temperature for more focused responses
		MaxTokens:    2000,
	}

	resp, err := s.aiProvider.Generate(ctx, aiReq)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

type parsedToolCall struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

func (s *AgenticContextService) parseToolCalls(response string) []parsedToolCall {
	// Find JSON in response
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 || end <= start {
		return nil
	}

	var parsed struct {
		ToolCalls []parsedToolCall `json:"tool_calls"`
	}

	if err := json.Unmarshal([]byte(response[start:end+1]), &parsed); err != nil {
		return nil
	}

	return parsed.ToolCalls
}

func (s *AgenticContextService) isContextComplete(toolCalls []parsedToolCall) bool {
	for _, call := range toolCalls {
		if call.Name == "mark_context_complete" {
			return true
		}
	}
	return false
}

func (s *AgenticContextService) extractReasoning(response string) string {
	// Try to extract reasoning from mark_context_complete call
	toolCalls := s.parseToolCalls(response)
	for _, call := range toolCalls {
		if call.Name == "mark_context_complete" {
			if reasoning, ok := call.Arguments["reasoning"].(string); ok {
				return reasoning
			}
		}
	}
	return ""
}

func (s *AgenticContextService) executeTool(
	call parsedToolCall,
	projectRoot string,
	requestedFiles map[string]bool,
	result *AgenticContextResult,
) (string, error) {
	switch call.Name {
	case "search_files_for_context":
		return s.executeSearchFiles(call.Arguments, projectRoot)

	case "search_content_for_context":
		return s.executeSearchContent(call.Arguments, projectRoot)

	case "get_file_symbols":
		return s.executeGetFileSymbols(call.Arguments, projectRoot)

	case "request_file":
		return s.executeRequestFile(call.Arguments, projectRoot, requestedFiles, result)

	case "mark_context_complete":
		return "Context gathering marked as complete.", nil

	default:
		return "", fmt.Errorf("unknown tool: %s", call.Name)
	}
}

func (s *AgenticContextService) executeSearchFiles(args map[string]any, projectRoot string) (string, error) {
	pattern, _ := args["pattern"].(string)
	directory, _ := args["directory"].(string)

	if pattern == "" {
		return "", fmt.Errorf("pattern is required")
	}

	files, err := s.fileTools.SearchFiles(pattern, directory, projectRoot)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return fmt.Sprintf("No files found matching pattern: %s", pattern), nil
	}

	return fmt.Sprintf("Found %d files:\n%s", len(files), strings.Join(files, "\n")), nil
}

func (s *AgenticContextService) executeSearchContent(args map[string]any, projectRoot string) (string, error) {
	query, _ := args["query"].(string)
	filePattern, _ := args["file_pattern"].(string)
	maxResults := 20
	if mr, ok := args["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	if query == "" {
		return "", fmt.Errorf("query is required")
	}

	matches, err := s.fileTools.SearchContent(query, filePattern, projectRoot, maxResults)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return fmt.Sprintf("No matches found for: %s", query), nil
	}

	var lines []string
	for _, m := range matches {
		lines = append(lines, fmt.Sprintf("%s:%d: %s", m.FilePath, m.Line, m.Content))
	}

	return fmt.Sprintf("Found %d matches:\n%s", len(matches), strings.Join(lines, "\n")), nil
}

func (s *AgenticContextService) executeGetFileSymbols(args map[string]any, projectRoot string) (string, error) {
	path, _ := args["path"].(string)

	if path == "" {
		return "", fmt.Errorf("path is required")
	}

	symbols, err := s.symbolTools.ListSymbols(path, projectRoot)
	if err != nil {
		return "", err
	}

	if len(symbols) == 0 {
		return fmt.Sprintf("No symbols found in: %s", path), nil
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("Symbols in %s:", path))
	for _, sym := range symbols {
		line := fmt.Sprintf("  [%s] %s", sym.Kind, sym.Name)
		if sym.Line > 0 {
			line += fmt.Sprintf(" (line %d)", sym.Line)
		}
		if sym.Signature != "" {
			line += fmt.Sprintf(" - %s", sym.Signature)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}

func (s *AgenticContextService) executeRequestFile(
	args map[string]any,
	projectRoot string,
	requestedFiles map[string]bool,
	result *AgenticContextResult,
) (string, error) {
	path, _ := args["path"].(string)
	reason, _ := args["reason"].(string)
	relevance := 0.8 // default relevance
	if rel, ok := args["relevance"].(float64); ok {
		relevance = rel
	}

	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	if reason == "" {
		reason = "Requested by context gatherer"
	}

	// Check if already requested
	if requestedFiles[path] {
		return fmt.Sprintf("File already in context: %s", path), nil
	}

	// Read file to estimate tokens
	content, err := s.fileTools.ReadFile(path, projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// Estimate tokens (rough: ~4 chars per token)
	tokens := len(content) / 4

	// Add to result
	result.Files = append(result.Files, ContextFileResult{
		Path:      path,
		Reason:    reason,
		Relevance: relevance,
		Tokens:    tokens,
	})
	result.TotalTokens += tokens
	requestedFiles[path] = true

	return fmt.Sprintf("Added to context: %s (%d tokens)\nReason: %s", path, tokens, reason), nil
}
