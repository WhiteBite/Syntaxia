package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"syntaxia/domain"
	"strings"
)

// AgenticStreamEvent represents an event during agentic chat
type AgenticStreamEvent struct {
	Type      string `json:"type"` // "thinking", "tool_call", "tool_result", "content", "done", "error"
	Content   string `json:"content,omitempty"`
	ToolName  string `json:"toolName,omitempty"`
	ToolArgs  string `json:"toolArgs,omitempty"`
	Iteration int    `json:"iteration,omitempty"`
}

// AgenticStreamCallback is called for each event
type AgenticStreamCallback func(event AgenticStreamEvent)

// ChatStream performs agentic chat with streaming events
func (s *AgenticChatService) ChatStream(ctx context.Context, req AgenticChatRequest, callback AgenticStreamCallback) error {
	s.logger.Info(fmt.Sprintf("Starting streaming agentic chat: %s", req.Task))

	tools := s.toolExecutor.GetAvailableTools()
	toolsJSON := s.formatToolsForPrompt(tools)

	systemPrompt := fmt.Sprintf(`You are an expert code assistant with access to tools.

AVAILABLE TOOLS:
%s

INSTRUCTIONS:
1. Use tools by responding with JSON: {"tool_calls": [{"name": "tool_name", "arguments": {...}}]}
2. After receiving results, either call more tools or provide your final answer
3. When done, provide answer WITHOUT tool_calls
4. Respond in user's language

Be thorough but efficient.`, toolsJSON)

	messages := []domain.ChatMessage{
		{Role: domain.RoleSystem, Content: systemPrompt},
		{Role: domain.RoleUser, Content: req.Task},
	}

	for iterations := 0; iterations < s.maxIterations; iterations++ {
		callback(AgenticStreamEvent{Type: "thinking", Content: fmt.Sprintf("Iteration %d...", iterations+1), Iteration: iterations + 1})

		response, err := s.callAI(ctx, messages)
		if err != nil {
			callback(AgenticStreamEvent{Type: "error", Content: err.Error()})
			return err
		}

		toolCalls := s.parseToolCalls(response)
		if len(toolCalls) == 0 {
			callback(AgenticStreamEvent{Type: "content", Content: response, Iteration: iterations + 1})
			callback(AgenticStreamEvent{Type: "done", Iteration: iterations + 1})
			return nil
		}

		messages = append(messages, domain.ChatMessage{Role: domain.RoleAssistant, Content: response})

		var toolResults []string
		for _, call := range toolCalls {
			argsJSON, _ := json.Marshal(call.Arguments)
			callback(AgenticStreamEvent{Type: "tool_call", ToolName: call.Name, ToolArgs: string(argsJSON), Iteration: iterations + 1})

			result := s.toolExecutor.ExecuteTool(call, req.ProjectRoot)
			callback(AgenticStreamEvent{Type: "tool_result", ToolName: call.Name, Content: domain.TruncateString(result.Content, 200), Iteration: iterations + 1})

			toolResults = append(toolResults, fmt.Sprintf("Tool: %s\nResult:\n%s", call.Name, result.Content))
		}

		messages = append(messages, domain.ChatMessage{Role: domain.RoleTool, Content: strings.Join(toolResults, "\n\n---\n\n")})
	}

	callback(AgenticStreamEvent{Type: "error", Content: fmt.Sprintf("Max iterations (%d) reached", s.maxIterations)})
	return fmt.Errorf("max iterations reached")
}
