package tools

import (
	"fmt"
	"syntaxia/domain"
)

// HandlerRegistry manages tool handlers and routes tool calls
type HandlerRegistry struct {
	handlers []ToolHandler
	logger   domain.Logger
}

// NewHandlerRegistry creates a new handler registry
func NewHandlerRegistry(logger domain.Logger) *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make([]ToolHandler, 0),
		logger:   logger,
	}
}

// Register adds a handler to the registry
func (r *HandlerRegistry) Register(handler ToolHandler) {
	r.handlers = append(r.handlers, handler)
}

// GetAllTools returns all tools from all registered handlers
func (r *HandlerRegistry) GetAllTools() []domain.Tool {
	var tools []domain.Tool
	for _, h := range r.handlers {
		tools = append(tools, h.GetTools()...)
	}
	return tools
}

// Execute finds the appropriate handler and executes the tool
func (r *HandlerRegistry) Execute(toolName string, args map[string]any, projectRoot string) (string, error) {
	for _, h := range r.handlers {
		if h.CanHandle(toolName) {
			return h.Execute(toolName, args, projectRoot)
		}
	}
	return "", fmt.Errorf("no handler found for tool: %s", toolName)
}

// CanHandle returns true if any handler can handle the tool
func (r *HandlerRegistry) CanHandle(toolName string) bool {
	for _, h := range r.handlers {
		if h.CanHandle(toolName) {
			return true
		}
	}
	return false
}
