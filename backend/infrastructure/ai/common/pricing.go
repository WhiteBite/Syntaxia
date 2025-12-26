// Package common provides shared utilities for AI providers
package common

import "syntaxia/domain"

// ModelPricing contains pricing info for a specific model
type ModelPricing struct {
	InputPer1K  float64
	OutputPer1K float64
}

// PricingTable maps model names to their pricing
type PricingTable map[string]ModelPricing

// GetPricingFromTable returns pricing info for a model using a pricing table
func GetPricingFromTable(model, currency string, table PricingTable, defaultPricing ModelPricing) domain.PricingInfo {
	pricing := domain.PricingInfo{
		Model:    model,
		Currency: currency,
	}

	if p, ok := table[model]; ok {
		pricing.InputTokensPer1K = p.InputPer1K
		pricing.OutputTokensPer1K = p.OutputPer1K
	} else {
		pricing.InputTokensPer1K = defaultPricing.InputPer1K
		pricing.OutputTokensPer1K = defaultPricing.OutputPer1K
	}

	return pricing
}

// OpenAIPricingTable contains pricing for OpenAI models
var OpenAIPricingTable = PricingTable{
	"gpt-4":         {InputPer1K: 0.03, OutputPer1K: 0.06},
	"gpt-4-turbo":   {InputPer1K: 0.01, OutputPer1K: 0.03},
	"gpt-3.5-turbo": {InputPer1K: 0.0015, OutputPer1K: 0.002},
}

// OpenAIDefaultPricing is the fallback pricing for unknown OpenAI models
var OpenAIDefaultPricing = ModelPricing{InputPer1K: 0.01, OutputPer1K: 0.02}

// GeminiPricingTable contains pricing for Gemini models
var GeminiPricingTable = PricingTable{
	"gemini-pro":        {InputPer1K: 0.0005, OutputPer1K: 0.0015},
	"gemini-pro-vision": {InputPer1K: 0.0005, OutputPer1K: 0.0015},
	"gemini-1.5-pro":    {InputPer1K: 0.00375, OutputPer1K: 0.015},
}

// GeminiDefaultPricing is the fallback pricing for unknown Gemini models
var GeminiDefaultPricing = ModelPricing{InputPer1K: 0.001, OutputPer1K: 0.002}

// QwenPricingTable contains pricing for Qwen models (in CNY)
var QwenPricingTable = PricingTable{
	"qwen-coder-plus-latest":  {InputPer1K: 0.004, OutputPer1K: 0.012},
	"qwen-coder-plus":         {InputPer1K: 0.004, OutputPer1K: 0.012},
	"qwen-coder-turbo-latest": {InputPer1K: 0.002, OutputPer1K: 0.006},
	"qwen-coder-turbo":        {InputPer1K: 0.002, OutputPer1K: 0.006},
	"qwen-plus-latest":        {InputPer1K: 0.004, OutputPer1K: 0.012},
	"qwen-plus":               {InputPer1K: 0.004, OutputPer1K: 0.012},
	"qwen-turbo-latest":       {InputPer1K: 0.002, OutputPer1K: 0.006},
	"qwen-turbo":              {InputPer1K: 0.002, OutputPer1K: 0.006},
	"qwen-max":                {InputPer1K: 0.02, OutputPer1K: 0.06},
	"qwen-max-latest":         {InputPer1K: 0.02, OutputPer1K: 0.06},
}

// QwenDefaultPricing is the fallback pricing for unknown Qwen models
var QwenDefaultPricing = ModelPricing{InputPer1K: 0.004, OutputPer1K: 0.012}
