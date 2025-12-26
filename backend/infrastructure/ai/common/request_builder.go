package common

import (
	"syntaxia/domain"

	"github.com/sashabaranov/go-openai"
)

// BuildChatMessages creates OpenAI chat messages from domain request
func BuildChatMessages(req domain.AIRequest) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: req.UserPrompt,
		},
	}
}

// BuildCompletionRequest creates OpenAI completion request from domain request
func BuildCompletionRequest(req domain.AIRequest, stream bool) openai.ChatCompletionRequest {
	completionReq := openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    BuildChatMessages(req),
		Temperature: float32(req.Temperature),
		MaxTokens:   req.MaxTokens,
		TopP:        float32(req.TopP),
		Stream:      stream,
	}

	if req.FrequencyPenalty != 0 {
		completionReq.FrequencyPenalty = float32(req.FrequencyPenalty)
	}
	if req.PresencePenalty != 0 {
		completionReq.PresencePenalty = float32(req.PresencePenalty)
	}

	return completionReq
}
