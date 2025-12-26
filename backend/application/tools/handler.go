package tools

import (
	"syntaxia/domain"
)

// ToolHandler defines the interface for tool handlers
type ToolHandler interface {
	// GetTools returns the list of tools provided by this handler
	GetTools() []domain.Tool
	// Execute executes a tool and returns the result
	Execute(toolName string, args map[string]any, projectRoot string) (string, error)
	// CanHandle returns true if this handler can handle the given tool
	CanHandle(toolName string) bool
}

// BaseHandler provides common functionality for tool handlers
type BaseHandler struct {
	Logger domain.Logger
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(logger domain.Logger) BaseHandler {
	return BaseHandler{Logger: logger}
}
