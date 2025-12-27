package ai

import "syntaxia/domain"

// ChatRequest represents a chat request with tool calling
type ChatRequest struct {
	Messages     []domain.ChatMessage       `json:"messages"`
	ProjectRoot  string                     `json:"projectRoot"`
	SmartContext *domain.SmartContextResult `json:"smartContext,omitempty"`
	SystemPrompt string                     `json:"systemPrompt,omitempty"`
	MaxTokens    int                        `json:"maxTokens,omitempty"`
	Temperature  float64                    `json:"temperature,omitempty"`
}

// ChatResponse represents the response from chat
type ChatResponse struct {
	Message    domain.ChatMessage `json:"message"`
	ToolCalls  []ToolCallResult   `json:"toolCalls,omitempty"`
	Iterations int                `json:"iterations"`
	TokensUsed int                `json:"tokensUsed"`
	Finished   bool               `json:"finished"`
}

// ToolCallResult represents a tool call with its result
type ToolCallResult struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
	Result    string `json:"result"`
	Error     string `json:"error,omitempty"`
	Duration  int64  `json:"duration"` // milliseconds
}

// ChatStreamEvent represents an event during streaming chat
type ChatStreamEvent struct {
	Type       string           `json:"type"` // "content", "tool_call", "tool_result", "thinking", "done", "error"
	Content    string           `json:"content,omitempty"`
	ToolCall   *ToolCallEvent   `json:"toolCall,omitempty"`
	ToolResult *ToolResultEvent `json:"toolResult,omitempty"`
	Iteration  int              `json:"iteration,omitempty"`
	TokensUsed int              `json:"tokensUsed,omitempty"`
}

// ToolCallEvent represents a tool call event
type ToolCallEvent struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// ToolResultEvent represents a tool result event
type ToolResultEvent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Result   string `json:"result"`
	Error    string `json:"error,omitempty"`
	Duration int64  `json:"duration"`
}

// ChatStreamCallback is called for each streaming event
type ChatStreamCallback func(event ChatStreamEvent)
