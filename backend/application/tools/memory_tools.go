package tools

import (
	"fmt"
	"syntaxia/domain"
)

// MemoryToolsHandler handles memory-related tools
type MemoryToolsHandler struct {
	BaseHandler
	ContextMemory domain.ContextMemory
}

// NewMemoryToolsHandler creates a new memory tools handler
func NewMemoryToolsHandler(logger domain.Logger, contextMemory domain.ContextMemory) *MemoryToolsHandler {
	return &MemoryToolsHandler{
		BaseHandler:   NewBaseHandler(logger),
		ContextMemory: contextMemory,
	}
}

var memoryToolNames = map[string]bool{
	"save_context":        true,
	"find_context":        true,
	"get_recent_contexts": true,
}

// CanHandle returns true if this handler can handle the given tool
func (h *MemoryToolsHandler) CanHandle(toolName string) bool {
	return memoryToolNames[toolName]
}

// GetTools returns the list of memory tools
func (h *MemoryToolsHandler) GetTools() []domain.Tool {
	return []domain.Tool{
		{
			Name:        "save_context",
			Description: "Save current context for later retrieval",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"topic":   {Type: "string", Description: "Topic/name for the context"},
					"summary": {Type: "string", Description: "Brief summary"},
					"files":   {Type: "array", Description: "List of file paths"},
				},
				Required: []string{"topic"},
			},
		},
		{
			Name:        "find_context",
			Description: "Find saved contexts by topic",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"topic": {Type: "string", Description: "Topic to search for"},
				},
				Required: []string{"topic"},
			},
		},
		{
			Name:        "get_recent_contexts",
			Description: "Get recently accessed contexts",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"limit": {Type: "integer", Description: "Maximum results (default: 10)"},
				},
			},
		},
	}
}

// Execute executes a memory tool
func (h *MemoryToolsHandler) Execute(toolName string, args map[string]any, projectRoot string) (string, error) {
	switch toolName {
	case "save_context":
		return h.saveContext(args, projectRoot)
	case "find_context":
		return h.findContext(args, projectRoot)
	case "get_recent_contexts":
		return h.getRecentContexts(args, projectRoot)
	default:
		return "", fmt.Errorf("unknown memory tool: %s", toolName)
	}
}

func (h *MemoryToolsHandler) saveContext(args map[string]any, projectRoot string) (string, error) {
	if h.ContextMemory == nil {
		return "", fmt.Errorf("context memory not initialized")
	}

	topic, _ := args["topic"].(string)
	summary, _ := args["summary"].(string)
	var files []string
	if f, ok := args["files"].([]any); ok {
		for _, file := range f {
			if s, ok := file.(string); ok {
				files = append(files, s)
			}
		}
	}

	ctx := &domain.ConversationContext{
		ProjectRoot: projectRoot,
		Topic:       topic,
		Summary:     summary,
		Files:       files,
	}

	if err := h.ContextMemory.SaveContext(ctx); err != nil {
		return "", fmt.Errorf("failed to save context: %w", err)
	}

	return fmt.Sprintf("Context saved: %s", topic), nil
}

func (h *MemoryToolsHandler) findContext(args map[string]any, projectRoot string) (string, error) {
	if h.ContextMemory == nil {
		return "", fmt.Errorf("context memory not initialized")
	}

	topic, _ := args["topic"].(string)
	if topic == "" {
		return "", fmt.Errorf("topic is required")
	}

	contexts, err := h.ContextMemory.FindContextByTopic(projectRoot, topic)
	if err != nil {
		return "", fmt.Errorf("failed to find context: %w", err)
	}

	if len(contexts) == 0 {
		return fmt.Sprintf("No contexts found for topic: %s", topic), nil
	}

	var result string
	result = fmt.Sprintf("Contexts matching '%s':\n\n", topic)
	for _, c := range contexts {
		result += fmt.Sprintf("  - %s: %s (%d files)\n", c.Topic, c.Summary, len(c.Files))
	}
	return result, nil
}

func (h *MemoryToolsHandler) getRecentContexts(args map[string]any, projectRoot string) (string, error) {
	if h.ContextMemory == nil {
		return "", fmt.Errorf("context memory not initialized")
	}

	limit := 10
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	contexts, err := h.ContextMemory.GetRecentContexts(projectRoot, limit)
	if err != nil {
		return "", fmt.Errorf("failed to get recent contexts: %w", err)
	}

	if len(contexts) == 0 {
		return "No recent contexts found", nil
	}

	var result string
	result = "Recent contexts:\n\n"
	for _, c := range contexts {
		result += fmt.Sprintf("  - %s: %s (last accessed: %s)\n", c.Topic, c.Summary, c.LastAccessed.Format("2006-01-02"))
	}
	return result, nil
}
