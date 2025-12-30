package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"syntaxia/domain"
	"strings"
)

// AgenticChatService handles AI chat with tool use capabilities
type AgenticChatService struct {
	logger        domain.Logger
	aiService     *Service
	toolExecutor  ToolExecutor
	maxIterations int
}

// ToolExecutor interface for tool execution
type ToolExecutor interface {
	GetAvailableTools() []domain.Tool
	ExecuteTool(call domain.ToolCall, projectRoot string) domain.ToolResult
}

// NewAgenticChatService creates a new agentic chat service
func NewAgenticChatService(logger domain.Logger, aiService *Service, toolExecutor ToolExecutor) *AgenticChatService {
	return &AgenticChatService{logger: logger, aiService: aiService, toolExecutor: toolExecutor, maxIterations: 10}
}

// AgenticChatRequest represents a request for agentic chat
type AgenticChatRequest struct {
	Task         string                    `json:"task"`
	ProjectRoot  string                    `json:"projectRoot"`
	Context      []string                  `json:"context,omitempty"`      // Selected file paths to include as context
	SmartContext *domain.SmartContextResult `json:"smartContext,omitempty"` // Pre-collected smart context
	MaxTokens    int                       `json:"maxTokens,omitempty"`
}

// AgenticChatResponse represents the response from agentic chat
type AgenticChatResponse struct {
	Response   string        `json:"response"`
	ToolCalls  []ToolCallLog `json:"toolCalls"`
	Iterations int           `json:"iterations"`
	Context    []string      `json:"context"`
}

// ToolCallLog logs a tool call for transparency
type ToolCallLog struct {
	Tool      string `json:"tool"`
	Arguments string `json:"arguments"`
	Result    string `json:"result"`
}

// Chat performs an agentic chat with tool use
func (s *AgenticChatService) Chat(ctx context.Context, req AgenticChatRequest) (*AgenticChatResponse, error) {
	s.logger.Info(fmt.Sprintf("Starting agentic chat: %s", req.Task))

	tools := s.toolExecutor.GetAvailableTools()
	toolsJSON := s.formatToolsForPrompt(tools)

	// Build context section from SmartContext or legacy Context
	contextSection := ""
	if req.SmartContext != nil {
		contextSection = s.formatSmartContext(req.SmartContext)
		s.logger.Info(fmt.Sprintf("Using smart context: %d files, %d tokens",
			len(req.SmartContext.RelevantFiles), req.SmartContext.TotalTokens))
	} else if len(req.Context) > 0 {
		contextSection = s.readContextFiles(req.Context, req.ProjectRoot)
	}

	systemPrompt := fmt.Sprintf(`You are an expert code assistant helping with a software project. Your primary goal is to ANSWER USER QUESTIONS directly and helpfully.

AVAILABLE TOOLS:
%s

HOW TO RESPOND:
1. ALWAYS answer the user's question directly. Don't just say you're ready to help - actually help!

2. If the user asks about files, code, or the project:
   - If context is provided below, use it to answer
   - If you need more information, use tools to get it

3. To use tools, respond with JSON:
   {"tool_calls": [{"name": "tool_name", "arguments": {"arg1": "value1"}}]}

4. After getting tool results, provide a clear answer based on what you learned.

5. Use write_file tool to create or modify files. Changes go to sandbox for user review.

CRITICAL RULES:
- NEVER respond with generic phrases like "I'm ready to help" or "What would you like me to do?"
- ALWAYS provide substantive answers based on the context or tool results
- If you don't have enough information, use tools to get it, then answer
- Respond in the user's language (Russian if they write in Russian)
%s`, toolsJSON, contextSection)

	messages := []domain.ChatMessage{
		{Role: domain.RoleSystem, Content: systemPrompt},
		{Role: domain.RoleUser, Content: req.Task},
	}

	var toolCallLogs []ToolCallLog
	var readFiles []string
	iterations := 0

	for iterations < s.maxIterations {
		iterations++
		s.logger.Info(fmt.Sprintf("Agentic chat iteration %d", iterations))

		response, err := s.callAI(ctx, messages)
		if err != nil {
			return nil, fmt.Errorf("AI call failed: %w", err)
		}

		toolCalls := s.parseToolCalls(response)
		if len(toolCalls) == 0 {
			s.logger.Info("Agentic chat completed with final answer")
			return &AgenticChatResponse{Response: response, ToolCalls: toolCallLogs, Iterations: iterations, Context: readFiles}, nil
		}

		messages = append(messages, domain.ChatMessage{Role: domain.RoleAssistant, Content: response})

		var toolResults []string
		for _, call := range toolCalls {
			result := s.toolExecutor.ExecuteTool(call, req.ProjectRoot)
			argsJSON, _ := json.Marshal(call.Arguments)
			toolCallLogs = append(toolCallLogs, ToolCallLog{Tool: call.Name, Arguments: string(argsJSON), Result: domain.TruncateString(result.Content, 500)})

			if call.Name == "read_file" {
				if path, ok := call.Arguments["path"].(string); ok {
					readFiles = append(readFiles, path)
				}
			}
			toolResults = append(toolResults, fmt.Sprintf("Tool: %s\nResult:\n%s", call.Name, result.Content))
		}

		messages = append(messages, domain.ChatMessage{Role: domain.RoleTool, Content: strings.Join(toolResults, "\n\n---\n\n")})
	}

	return nil, fmt.Errorf("max iterations (%d) reached without final answer", s.maxIterations)
}

func (s *AgenticChatService) formatToolsForPrompt(tools []domain.Tool) string {
	lines := make([]string, 0, len(tools))
	for _, tool := range tools {
		params := []string{}
		for name, prop := range tool.Parameters.Properties {
			required := ""
			for _, r := range tool.Parameters.Required {
				if r == name {
					required = " (required)"
					break
				}
			}
			params = append(params, fmt.Sprintf("  - %s: %s%s", name, prop.Description, required))
		}
		lines = append(lines, fmt.Sprintf("\n%s: %s\nParameters:\n%s", tool.Name, tool.Description, strings.Join(params, "\n")))
	}
	return strings.Join(lines, "\n")
}

func (s *AgenticChatService) callAI(ctx context.Context, messages []domain.ChatMessage) (string, error) {
	var systemPrompt string
	var userPrompt strings.Builder

	for _, msg := range messages {
		switch msg.Role {
		case domain.RoleSystem:
			systemPrompt = msg.Content
		case domain.RoleUser:
			userPrompt.WriteString("User: ")
			userPrompt.WriteString(msg.Content)
			userPrompt.WriteString("\n\n")
		case domain.RoleAssistant:
			userPrompt.WriteString("Assistant: ")
			userPrompt.WriteString(msg.Content)
			userPrompt.WriteString("\n\n")
		case domain.RoleTool:
			userPrompt.WriteString("Tool Results:\n")
			userPrompt.WriteString(msg.Content)
			userPrompt.WriteString("\n\n")
		}
	}
	userPrompt.WriteString("Assistant: ")

	return s.aiService.GenerateCode(ctx, systemPrompt, userPrompt.String())
}

func (s *AgenticChatService) parseToolCalls(response string) []domain.ToolCall {
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 || end <= start {
		return nil
	}

	var parsed struct {
		ToolCalls []struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		} `json:"tool_calls"`
	}

	if err := json.Unmarshal([]byte(response[start:end+1]), &parsed); err != nil {
		return nil
	}

	calls := make([]domain.ToolCall, 0, len(parsed.ToolCalls))
	for i, tc := range parsed.ToolCalls {
		calls = append(calls, domain.ToolCall{ID: fmt.Sprintf("call_%d", i), Name: tc.Name, Arguments: tc.Arguments})
	}
	return calls
}

func (s *AgenticChatService) formatSmartContext(sc *domain.SmartContextResult) string {
	var sb strings.Builder

	if sc.ProjectStructure != "" {
		sb.WriteString("\n\n")
		sb.WriteString(sc.ProjectStructure)
	}

	if len(sc.RelevantFiles) > 0 {
		sb.WriteString("\n\nRELEVANT FILES:\n")
		for _, f := range sc.RelevantFiles {
			sb.WriteString(fmt.Sprintf("\n=== %s ===\n%s\n", f.Path, f.Content))
		}
	}

	return sb.String()
}

func (s *AgenticChatService) readContextFiles(files []string, projectRoot string) string {
	var contextParts []string
	for _, filePath := range files {
		result := s.toolExecutor.ExecuteTool(domain.ToolCall{
			ID:        "context_read",
			Name:      "read_file",
			Arguments: map[string]any{"path": filePath},
		}, projectRoot)
		if result.Error == "" {
			contextParts = append(contextParts, fmt.Sprintf("=== %s ===\n%s", filePath, result.Content))
		}
	}
	if len(contextParts) > 0 {
		return fmt.Sprintf("\n\nSELECTED FILES CONTEXT:\n%s\n", strings.Join(contextParts, "\n\n"))
	}
	return ""
}
