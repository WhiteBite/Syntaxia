package common

import (
	"context"
	"errors"
	"fmt"
	"syntaxia/domain"

	"github.com/sashabaranov/go-openai"
)

// StreamProcessor handles common stream processing logic for OpenAI-compatible APIs
func StreamProcessor(stream *openai.ChatCompletionStream, onChunk func(domain.StreamChunk), log domain.Logger) error {
	totalTokens := 0
	for {
		response, err := stream.Recv()
		if errors.Is(err, context.Canceled) {
			onChunk(domain.StreamChunk{Done: true, Error: "Request cancelled"})
			return err
		}
		if err != nil {
			if err.Error() == "EOF" {
				onChunk(domain.StreamChunk{Done: true, TokensUsed: totalTokens, FinishReason: "stop"})
				return nil
			}
			log.Error(fmt.Sprintf("Stream error: %v", err))
			onChunk(domain.StreamChunk{Done: true, Error: err.Error()})
			return err
		}

		if len(response.Choices) > 0 {
			content := response.Choices[0].Delta.Content
			finishReason := string(response.Choices[0].FinishReason)

			if content != "" {
				totalTokens += len(content) / 4
				onChunk(domain.StreamChunk{
					Content: content,
					Done:    false,
				})
			}

			if finishReason == "stop" || finishReason == "length" {
				onChunk(domain.StreamChunk{
					Done:         true,
					TokensUsed:   totalTokens,
					FinishReason: finishReason,
				})
				return nil
			}
		}
	}
}
