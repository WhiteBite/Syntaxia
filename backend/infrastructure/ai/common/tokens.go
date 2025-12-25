// Package common provides shared utilities for AI providers
package common

import "syntaxia/domain"

const (
	// DefaultCharsPerToken is the approximate number of characters per token
	// This is a rough estimate that works reasonably well for most models
	DefaultCharsPerToken = 4
	// DefaultTokenBuffer is added to estimates for safety margin
	DefaultTokenBuffer = 100
)

// EstimateTokens provides a simple token estimation based on character count
// This is used as a fallback when providers don't have their own tokenizer
func EstimateTokens(req domain.AIRequest) (int, error) {
	return EstimateTokensWithConfig(req, DefaultCharsPerToken, DefaultTokenBuffer)
}

// EstimateTokensWithConfig allows customizing the estimation parameters
func EstimateTokensWithConfig(req domain.AIRequest, charsPerToken, buffer int) (int, error) {
	totalChars := len(req.SystemPrompt) + len(req.UserPrompt)
	estimatedTokens := totalChars / charsPerToken
	return estimatedTokens + buffer, nil
}

// EstimateTokensFromText estimates tokens for a plain text string
func EstimateTokensFromText(text string) int {
	return len(text)/DefaultCharsPerToken + DefaultTokenBuffer
}
