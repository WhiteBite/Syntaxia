// Package common provides shared utilities for AI providers
package common

import (
	"context"
	"errors"
	"net/http"
	"syntaxia/domain"

	"github.com/sashabaranov/go-openai"
)

// HandleOpenAIError converts OpenAI API errors to domain errors
// This is used by providers that use the go-openai client (OpenAI, Qwen, OpenRouter)
func HandleOpenAIError(err error) error {
	if err == nil {
		return nil
	}

	var apiErr *openai.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.HTTPStatusCode {
		case http.StatusUnauthorized:
			return domain.ErrInvalidAPIKey
		case http.StatusTooManyRequests:
			return domain.ErrRateLimitExceeded
		}
	}
	return err
}

// IsContextCanceled checks if the error is due to context cancellation
func IsContextCanceled(err error) bool {
	return errors.Is(err, context.Canceled)
}

// IsEOF checks if the error indicates end of stream
func IsEOF(err error) bool {
	return err != nil && err.Error() == "EOF"
}
