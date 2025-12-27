package ai

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"syntaxia/domain"
)

const defaultSystemPrompt = `You are an expert code assistant with access to tools for exploring and modifying codebases.

CAPABILITIES:
- Read and analyze source code files
- Search for files and content patterns
- Explore project structure and symbols
- Write and modify files (changes go to sandbox for review)
- Access git history and status
- Analyze code dependencies and call graphs

GUIDELINES:
1. Use tools to gather information before making changes
2. Be thorough but efficient - don't read files unnecessarily
3. When modifying code, explain your changes clearly
4. Write operations go to sandbox - user must approve before applying
5. Respond in the user's language

IMPORTANT: Always provide clear, actionable responses. If you need more information, use the appropriate tools.`

const toolCallInstructions = `To use tools, respond with JSON: {"tool_calls": [{"name": "tool_name", "arguments": {...}}]}
When done, provide your final answer without tool_calls.`

// formatSmartContext formats smart context for the prompt
func formatSmartContext(sc *domain.SmartContextResult) string {
	if sc == nil {
		return ""
	}

	var sb strings.Builder

	if sc.ProjectStructure != "" {
		sb.WriteString("PROJECT STRUCTURE:\n")
		sb.WriteString(sc.ProjectStructure)
		sb.WriteString("\n\n")
	}

	if len(sc.RelevantFiles) > 0 {
		sb.WriteString("RELEVANT FILES:\n")
		for _, f := range sc.RelevantFiles {
			sb.WriteString(fmt.Sprintf("\n=== %s ===\n%s\n", f.Path, f.Content))
		}
	}

	return sb.String()
}

// messagesToPrompts converts messages to system and user prompts
func messagesToPrompts(messages []domain.ChatMessage) (string, string) {
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
			userPrompt.WriteString("Tool Result (")
			userPrompt.WriteString(msg.ToolCallID)
			userPrompt.WriteString("):\n")
			userPrompt.WriteString(msg.Content)
			userPrompt.WriteString("\n\n")
		}
	}
	userPrompt.WriteString("Assistant: ")

	return systemPrompt, userPrompt.String()
}

// formatToolsForPrompt formats tools for the prompt
func formatToolsForPrompt(tools []domain.Tool) string {
	var sb strings.Builder

	for _, tool := range tools {
		sb.WriteString(fmt.Sprintf("\n%s: %s\n", tool.Name, tool.Description))
		sb.WriteString("Parameters:\n")

		for name, prop := range tool.Parameters.Properties {
			required := ""
			for _, r := range tool.Parameters.Required {
				if r == name {
					required = " (required)"
					break
				}
			}
			sb.WriteString(fmt.Sprintf("  - %s (%s): %s%s\n", name, prop.Type, prop.Description, required))
		}
	}

	return sb.String()
}

// parseToolCalls parses tool calls from AI response
func parseToolCalls(response string, maxCalls int) []domain.ToolCall {
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 || end <= start {
		return nil
	}

	jsonStr := response[start : end+1]

	var parsed struct {
		ToolCalls []struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		} `json:"tool_calls"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil
	}

	if len(parsed.ToolCalls) == 0 {
		return nil
	}

	callCount := len(parsed.ToolCalls)
	if callCount > maxCalls {
		callCount = maxCalls
	}

	calls := make([]domain.ToolCall, 0, callCount)
	for i := 0; i < callCount; i++ {
		tc := parsed.ToolCalls[i]
		calls = append(calls, domain.ToolCall{
			ID:        fmt.Sprintf("call_%d_%d", time.Now().UnixNano(), i),
			Name:      tc.Name,
			Arguments: tc.Arguments,
		})
	}

	return calls
}

// cleanResponseFromToolCalls removes tool calls JSON from response
func cleanResponseFromToolCalls(response string) string {
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 {
		return response
	}

	before := strings.TrimSpace(response[:start])
	after := strings.TrimSpace(response[end+1:])

	result := before
	if after != "" {
		if result != "" {
			result += "\n"
		}
		result += after
	}

	return result
}
