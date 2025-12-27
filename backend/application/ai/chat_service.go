// Package ai provides AI chat service with streaming tool calling support.
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"syntaxia/domain"
)

// ChatService handles AI chat with native OpenAI function calling format
type ChatService struct {
	logger       domain.Logger
	aiService    *Service
	toolExecutor ToolExecutor
	config       ChatServiceConfig
}

// ChatServiceConfig configuration for chat service
type ChatServiceConfig struct {
	MaxIterations     int           // Maximum tool calling iterations
	MaxHistoryLength  int           // Maximum messages in history
	ToolTimeout       time.Duration // Timeout for tool execution
	StreamingEnabled  bool          // Enable streaming responses
	SandboxEnabled    bool          // Enable sandbox for write operations
	MaxToolsPerCall   int           // Maximum tools per single AI call
	MaxToolResultSize int           // Maximum size of tool result in bytes
}

// DefaultChatServiceConfig returns default configuration
func DefaultChatServiceConfig() ChatServiceConfig {
	return ChatServiceConfig{
		MaxIterations:     15,
		MaxHistoryLength:  50,
		ToolTimeout:       30 * time.Second,
		StreamingEnabled:  true,
		SandboxEnabled:    true,
		MaxToolsPerCall:   10,
		MaxToolResultSize: 50000,
	}
}

// NewChatService creates a new chat service
func NewChatService(
	logger domain.Logger,
	aiService *Service,
	toolExecutor ToolExecutor,
) *ChatService {
	return &ChatService{
		logger:       logger,
		aiService:    aiService,
		toolExecutor: toolExecutor,
		config:       DefaultChatServiceConfig(),
	}
}

// NewChatServiceWithConfig creates a new chat service with custom config
func NewChatServiceWithConfig(
	logger domain.Logger,
	aiService *Service,
	toolExecutor ToolExecutor,
	config ChatServiceConfig,
) *ChatService {
	return &ChatService{
		logger:       logger,
		aiService:    aiService,
		toolExecutor: toolExecutor,
		config:       config,
	}
}

// Chat performs a chat with tool calling support
func (s *ChatService) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	s.logger.Info("Starting chat with tool calling")

	messages := s.prepareMessages(req)
	tools := s.toolExecutor.GetAvailableTools()

	var allToolCalls []ToolCallResult
	iterations := 0
	totalTokens := 0

	for iterations < s.config.MaxIterations {
		iterations++
		s.logger.Info(fmt.Sprintf("Chat iteration %d", iterations))

		response, toolCalls, tokens, err := s.callAIWithTools(ctx, messages, tools, req)
		if err != nil {
			return nil, fmt.Errorf("AI call failed at iteration %d: %w", iterations, err)
		}
		totalTokens += tokens

		if len(toolCalls) == 0 {
			s.logger.Info("Chat completed with final response")
			return &ChatResponse{
				Message:    domain.ChatMessage{Role: domain.RoleAssistant, Content: response},
				ToolCalls:  allToolCalls,
				Iterations: iterations,
				TokensUsed: totalTokens,
				Finished:   true,
			}, nil
		}

		messages = append(messages, domain.ChatMessage{
			Role:      domain.RoleAssistant,
			Content:   response,
			ToolCalls: toolCalls,
		})

		toolResults := s.executeTools(ctx, toolCalls, req.ProjectRoot)
		for _, result := range toolResults {
			allToolCalls = append(allToolCalls, result)
			messages = append(messages, domain.ChatMessage{
				Role:       domain.RoleTool,
				Content:    result.Result,
				ToolCallID: result.ID,
			})
		}
	}

	return nil, fmt.Errorf("max iterations (%d) reached without final answer", s.config.MaxIterations)
}

// ChatStream performs streaming chat with tool calling
func (s *ChatService) ChatStream(ctx context.Context, req ChatRequest, callback ChatStreamCallback) error {
	s.logger.Info("Starting streaming chat with tool calling")

	messages := s.prepareMessages(req)
	tools := s.toolExecutor.GetAvailableTools()

	iterations := 0
	totalTokens := 0

	for iterations < s.config.MaxIterations {
		iterations++
		callback(ChatStreamEvent{
			Type:      "thinking",
			Content:   fmt.Sprintf("Processing iteration %d...", iterations),
			Iteration: iterations,
		})

		response, toolCalls, tokens, err := s.callAIWithToolsStream(ctx, messages, tools, req, callback)
		if err != nil {
			callback(ChatStreamEvent{Type: "error", Content: err.Error()})
			return err
		}
		totalTokens += tokens

		if len(toolCalls) == 0 {
			callback(ChatStreamEvent{
				Type:       "done",
				Content:    response,
				Iteration:  iterations,
				TokensUsed: totalTokens,
			})
			return nil
		}

		messages = append(messages, domain.ChatMessage{
			Role:      domain.RoleAssistant,
			Content:   response,
			ToolCalls: toolCalls,
		})

		toolResults := s.executeToolsWithCallback(ctx, toolCalls, req.ProjectRoot, callback, iterations)
		for _, result := range toolResults {
			messages = append(messages, domain.ChatMessage{
				Role:       domain.RoleTool,
				Content:    result.Result,
				ToolCallID: result.ID,
			})
		}
	}

	callback(ChatStreamEvent{
		Type:    "error",
		Content: fmt.Sprintf("Max iterations (%d) reached", s.config.MaxIterations),
	})
	return fmt.Errorf("max iterations reached")
}

// prepareMessages prepares messages with system prompt and context
func (s *ChatService) prepareMessages(req ChatRequest) []domain.ChatMessage {
	messages := make([]domain.ChatMessage, 0, len(req.Messages)+2)

	systemPrompt := s.buildSystemPrompt(req)
	messages = append(messages, domain.ChatMessage{
		Role:    domain.RoleSystem,
		Content: systemPrompt,
	})

	if req.SmartContext != nil {
		contextContent := formatSmartContext(req.SmartContext)
		if contextContent != "" {
			messages = append(messages, domain.ChatMessage{
				Role:    domain.RoleUser,
				Content: fmt.Sprintf("[Context]\n%s", contextContent),
			})
		}
	}

	userMessages := req.Messages
	if len(userMessages) > s.config.MaxHistoryLength {
		userMessages = userMessages[len(userMessages)-s.config.MaxHistoryLength:]
	}
	messages = append(messages, userMessages...)

	return messages
}

// buildSystemPrompt builds the system prompt with tool instructions
func (s *ChatService) buildSystemPrompt(req ChatRequest) string {
	if req.SystemPrompt != "" {
		return req.SystemPrompt
	}

	return defaultSystemPrompt
}

// callAIWithTools calls AI with tools in OpenAI function calling format
func (s *ChatService) callAIWithTools(
	ctx context.Context,
	messages []domain.ChatMessage,
	tools []domain.Tool,
	_ ChatRequest,
) (string, []domain.ToolCall, int, error) {
	systemPrompt, userPrompt := messagesToPrompts(messages)

	toolsJSON := formatToolsForPrompt(tools)
	fullSystemPrompt := fmt.Sprintf("%s\n\nAVAILABLE TOOLS:\n%s\n\n%s", systemPrompt, toolsJSON, toolCallInstructions)

	response, err := s.aiService.GenerateCode(ctx, fullSystemPrompt, userPrompt)
	if err != nil {
		return "", nil, 0, err
	}

	toolCalls := parseToolCalls(response, s.config.MaxToolsPerCall)
	cleanResponse := response
	if len(toolCalls) > 0 {
		cleanResponse = cleanResponseFromToolCalls(response)
	}

	return cleanResponse, toolCalls, 0, nil
}

// callAIWithToolsStream calls AI with streaming and tool support
func (s *ChatService) callAIWithToolsStream(
	ctx context.Context,
	messages []domain.ChatMessage,
	tools []domain.Tool,
	_ ChatRequest,
	callback ChatStreamCallback,
) (string, []domain.ToolCall, int, error) {
	systemPrompt, userPrompt := messagesToPrompts(messages)

	toolsJSON := formatToolsForPrompt(tools)
	fullSystemPrompt := fmt.Sprintf("%s\n\nAVAILABLE TOOLS:\n%s\n\n%s", systemPrompt, toolsJSON, toolCallInstructions)

	var responseBuilder strings.Builder
	var totalTokens int

	err := s.aiService.GenerateCodeStream(ctx, fullSystemPrompt, userPrompt, func(chunk domain.StreamChunk) {
		if chunk.Error != "" {
			callback(ChatStreamEvent{Type: "error", Content: chunk.Error})
			return
		}

		if chunk.Content != "" {
			responseBuilder.WriteString(chunk.Content)
			if !strings.Contains(chunk.Content, `"tool_calls"`) {
				callback(ChatStreamEvent{Type: "content", Content: chunk.Content})
			}
		}

		if chunk.TokensUsed > 0 {
			totalTokens = chunk.TokensUsed
		}
	})

	if err != nil {
		return "", nil, 0, err
	}

	response := responseBuilder.String()
	toolCalls := parseToolCalls(response, s.config.MaxToolsPerCall)
	cleanResponse := response
	if len(toolCalls) > 0 {
		cleanResponse = cleanResponseFromToolCalls(response)
	}

	return cleanResponse, toolCalls, totalTokens, nil
}

// executeTools executes tool calls and returns results
func (s *ChatService) executeTools(ctx context.Context, toolCalls []domain.ToolCall, projectRoot string) []ToolCallResult {
	results := make([]ToolCallResult, 0, len(toolCalls))
	for _, call := range toolCalls {
		result := s.executeSingleTool(ctx, call, projectRoot)
		results = append(results, result)
	}
	return results
}

// executeToolsWithCallback executes tools with streaming callbacks
func (s *ChatService) executeToolsWithCallback(
	ctx context.Context,
	toolCalls []domain.ToolCall,
	projectRoot string,
	callback ChatStreamCallback,
	iteration int,
) []ToolCallResult {
	results := make([]ToolCallResult, 0, len(toolCalls))

	for _, call := range toolCalls {
		callback(ChatStreamEvent{
			Type:      "tool_call",
			Iteration: iteration,
			ToolCall: &ToolCallEvent{
				ID:        call.ID,
				Name:      call.Name,
				Arguments: call.Arguments,
			},
		})

		result := s.executeSingleTool(ctx, call, projectRoot)
		results = append(results, result)

		callback(ChatStreamEvent{
			Type:      "tool_result",
			Iteration: iteration,
			ToolResult: &ToolResultEvent{
				ID:       result.ID,
				Name:     result.Name,
				Result:   domain.TruncateString(result.Result, 500),
				Error:    result.Error,
				Duration: result.Duration,
			},
		})
	}

	return results
}

// executeSingleTool executes a single tool with timeout
func (s *ChatService) executeSingleTool(ctx context.Context, call domain.ToolCall, projectRoot string) ToolCallResult {
	argsJSON, _ := json.Marshal(call.Arguments)

	result := ToolCallResult{
		ID:        call.ID,
		Name:      call.Name,
		Arguments: string(argsJSON),
	}

	toolCtx, cancel := context.WithTimeout(ctx, s.config.ToolTimeout)
	defer cancel()

	start := time.Now()

	resultChan := make(chan domain.ToolResult, 1)
	go func() {
		toolResult := s.toolExecutor.ExecuteTool(call, projectRoot)
		resultChan <- toolResult
	}()

	select {
	case <-toolCtx.Done():
		result.Error = "tool execution timed out"
		result.Result = fmt.Sprintf("Error: tool '%s' execution timed out after %v", call.Name, s.config.ToolTimeout)
	case toolResult := <-resultChan:
		if toolResult.Error != "" {
			result.Error = toolResult.Error
		}
		if len(toolResult.Content) > s.config.MaxToolResultSize {
			result.Result = toolResult.Content[:s.config.MaxToolResultSize] + "\n... (truncated)"
		} else {
			result.Result = toolResult.Content
		}
	}

	result.Duration = time.Since(start).Milliseconds()
	s.logger.Info(fmt.Sprintf("Tool %s executed in %dms", call.Name, result.Duration))

	return result
}
