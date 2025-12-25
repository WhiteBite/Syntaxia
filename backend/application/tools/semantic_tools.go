package tools

import (
	"fmt"
	"strings"

	"syntaxia/domain"
)

// SemanticSearcher interface for semantic search service
type SemanticSearcher interface {
	Search(query string, limit int) ([]domain.SemanticSearchResult, error)
}

// SemanticToolsHandler handles semantic search tools
type SemanticToolsHandler struct {
	BaseHandler
	searcher SemanticSearcher
}

// NewSemanticToolsHandler creates a new semantic tools handler
func NewSemanticToolsHandler(logger domain.Logger, searcher SemanticSearcher) *SemanticToolsHandler {
	return &SemanticToolsHandler{
		BaseHandler: NewBaseHandler(logger),
		searcher:    searcher,
	}
}

var semanticToolNames = map[string]bool{
	"semantic_search": true,
}

// CanHandle returns true if this handler can handle the given tool
func (h *SemanticToolsHandler) CanHandle(toolName string) bool {
	return semanticToolNames[toolName]
}

// GetTools returns the list of semantic tools
func (h *SemanticToolsHandler) GetTools() []domain.Tool {
	return []domain.Tool{
		{
			Name:        "semantic_search",
			Description: "Search code by meaning/intent using AI embeddings. Use for conceptual queries like 'error handling code' or 'user authentication logic'.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"query":        {Type: "string", Description: "Natural language description of what you're looking for"},
					"limit":        {Type: "integer", Description: "Maximum results (default: 10)"},
					"file_pattern": {Type: "string", Description: "Glob pattern to filter files (e.g., '*.go')"},
				},
				Required: []string{"query"},
			},
		},
	}
}

// Execute executes a semantic tool
func (h *SemanticToolsHandler) Execute(toolName string, args map[string]any, projectRoot string) (string, error) {
	switch toolName {
	case "semantic_search":
		return h.semanticSearch(args)
	default:
		return "", fmt.Errorf("unknown semantic tool: %s", toolName)
	}
}

func (h *SemanticToolsHandler) semanticSearch(args map[string]any) (string, error) {
	query, _ := args["query"].(string)
	if query == "" {
		return "", fmt.Errorf("query is required")
	}

	limit := 10
	if l, ok := args["limit"].(float64); ok && l > 0 {
		limit = int(l)
	}

	if h.searcher == nil {
		return "Semantic search is not available. Please configure embeddings first.", nil
	}

	results, err := h.searcher.Search(query, limit)
	if err != nil {
		return "", fmt.Errorf("semantic search failed: %w", err)
	}

	if len(results) == 0 {
		return fmt.Sprintf("No results found for: %s", query), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Semantic search results for '%s':\n\n", query))

	for i, r := range results {
		result.WriteString(fmt.Sprintf("%d. %s", i+1, r.Chunk.FilePath))
		if r.Chunk.StartLine > 0 {
			result.WriteString(fmt.Sprintf(":%d", r.Chunk.StartLine))
		}
		result.WriteString(fmt.Sprintf(" (score: %.2f)\n", r.Score))
		if r.Chunk.Content != "" {
			// Показываем первые 100 символов контента как сниппет
			snippet := r.Chunk.Content
			if len(snippet) > 100 {
				snippet = snippet[:100] + "..."
			}
			result.WriteString(fmt.Sprintf("   %s\n", snippet))
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}
