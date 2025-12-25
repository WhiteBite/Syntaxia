package domain

// Tool represents a tool that AI can call
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  ToolParameters `json:"parameters"`
}

// ToolParameters defines the JSON schema for tool parameters
type ToolParameters struct {
	Type       string                  `json:"type"`
	Properties map[string]ToolProperty `json:"properties"`
	Required   []string                `json:"required,omitempty"`
}

// ToolProperty defines a single parameter property
type ToolProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
	Default     any      `json:"default,omitempty"`
}

// ToolCall represents a tool call from AI
type ToolCall struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
	Error      string `json:"error,omitempty"`
}

// ToolExecutor executes tools and returns results
type ToolExecutor interface {
	// GetAvailableTools returns all available tools
	GetAvailableTools() []Tool

	// ExecuteTool executes a tool and returns the result
	ExecuteTool(call ToolCall, projectRoot string) ToolResult
}

// ChatMessage for tool use conversations
type ChatMessageRole string

const (
	RoleUser      ChatMessageRole = "user"
	RoleAssistant ChatMessageRole = "assistant"
	RoleTool      ChatMessageRole = "tool"
	RoleSystem    ChatMessageRole = "system"
)

type ChatMessage struct {
	Role       ChatMessageRole `json:"role"`
	Content    string          `json:"content,omitempty"`
	ToolCalls  []ToolCall      `json:"tool_calls,omitempty"`
	ToolCallID string          `json:"tool_call_id,omitempty"`
}
