// Package common provides shared utilities for AI providers
package common

import (
	"fmt"
	"syntaxia/domain"
)

// ValidationConfig allows customizing validation rules per provider
type ValidationConfig struct {
	RequireModel        bool
	RequireUserPrompt   bool
	RequireSystemPrompt bool
	MinTemperature      float64
	MaxTemperature      float64
	MinMaxTokens        int
}

// DefaultValidationConfig returns standard validation config
func DefaultValidationConfig() ValidationConfig {
	return ValidationConfig{
		RequireModel:        true,
		RequireUserPrompt:   true,
		RequireSystemPrompt: false,
		MinTemperature:      0,
		MaxTemperature:      2,
		MinMaxTokens:        1,
	}
}

// ValidateRequest validates an AI request with the given config
func ValidateRequest(req domain.AIRequest, cfg ValidationConfig) error {
	if cfg.RequireModel && req.Model == "" {
		return fmt.Errorf("model is required")
	}
	if cfg.RequireUserPrompt && req.UserPrompt == "" {
		return fmt.Errorf("user prompt is required")
	}
	if cfg.RequireSystemPrompt && req.SystemPrompt == "" {
		return fmt.Errorf("system prompt is required")
	}
	if req.Temperature < cfg.MinTemperature || req.Temperature > cfg.MaxTemperature {
		return fmt.Errorf("temperature must be between %.1f and %.1f", cfg.MinTemperature, cfg.MaxTemperature)
	}
	if req.MaxTokens < cfg.MinMaxTokens {
		return fmt.Errorf("max tokens must be greater than %d", cfg.MinMaxTokens-1)
	}
	return nil
}

// ValidateRequestStrict validates with full parameter checking (OpenAI, Qwen style)
func ValidateRequestStrict(req domain.AIRequest) error {
	return ValidateRequest(req, DefaultValidationConfig())
}

// ValidateRequestBasic validates only model and prompt (Gemini, LocalAI style)
func ValidateRequestBasic(req domain.AIRequest) error {
	cfg := ValidationConfig{
		RequireModel:      true,
		RequireUserPrompt: true,
	}
	return ValidateRequest(req, cfg)
}
