package ai

import (
	"context"
	"errors"
	"fmt"
	"syntaxia/domain"
	"syntaxia/infrastructure/ai/common"
	"sort"
	"time"

	"github.com/sashabaranov/go-openai"
)

const (
	// QwenMaxContextTokens is the maximum context window for qwen-coder-plus (1M tokens)
	QwenMaxContextTokens = 1000000

	// Qwen model names
	qwenCoderPlusLatest  = "qwen-coder-plus-latest"
	qwenCoderTurboLatest = "qwen-coder-turbo-latest"
)

// QwenProviderImpl implements domain.AIProvider for Alibaba Qwen models
type QwenProviderImpl struct {
	client *openai.Client
	log    domain.Logger
	host   string
}

// NewQwen creates a new Qwen provider instance
func NewQwen(apiKey, host string, log domain.Logger) (domain.AIProvider, error) {
	if host == "" {
		host = domain.QwenDefaultHost
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = host

	client := openai.NewClientWithConfig(config)
	return &QwenProviderImpl{
		client: client,
		log:    log,
		host:   host,
	}, nil
}

// ListModels returns available Qwen models
func (p *QwenProviderImpl) ListModels(ctx context.Context) ([]string, error) {
	p.log.Info("Requesting model list from Qwen API...")

	// Qwen API may not support model listing, return known models
	// These are the main Qwen Coder models optimized for code generation
	models := []string{
		"qwen-coder-plus-latest",  // 1M context, best for large codebases
		"qwen-coder-plus",         // 1M context
		"qwen-coder-turbo-latest", // Faster, smaller context
		"qwen-coder-turbo",
		"qwen-plus-latest", // General purpose with good coding
		"qwen-plus",
		"qwen-turbo-latest", // Fast general purpose
		"qwen-turbo",
		"qwen-max", // Most capable general model
		"qwen-max-latest",
	}

	// Try to fetch from API if supported
	resp, err := p.client.ListModels(ctx)
	if err == nil && len(resp.Models) > 0 {
		apiModels := make([]string, 0, len(resp.Models))
		for _, model := range resp.Models {
			apiModels = append(apiModels, model.ID)
		}
		sort.Strings(apiModels)
		p.log.Info(fmt.Sprintf("Received %d models from Qwen API", len(apiModels)))
		return apiModels, nil
	}

	// Fallback to known models
	p.log.Info(fmt.Sprintf("Using %d known Qwen models", len(models)))
	return models, nil
}

// Generate sends a request to Qwen API and returns the response
func (p *QwenProviderImpl) Generate(ctx context.Context, req domain.AIRequest) (domain.AIResponse, error) {
	startTime := time.Now()
	p.log.Info(fmt.Sprintf("Sending request to Qwen API with model: %s", req.Model))

	completionReq := common.BuildCompletionRequest(req, false)
	resp, err := p.client.CreateChatCompletion(ctx, completionReq)
	if err != nil {
		p.log.Error(fmt.Sprintf("Qwen API request failed: %v", err))
		if domainErr := common.HandleOpenAIError(err); !errors.Is(domainErr, err) {
			return domain.AIResponse{}, domainErr
		}
		return domain.AIResponse{}, err
	}

	if len(resp.Choices) == 0 {
		return domain.AIResponse{}, fmt.Errorf("no choices returned from Qwen API")
	}

	processingTime := time.Since(startTime)
	tokensUsed := resp.Usage.TotalTokens

	return domain.AIResponse{
		Content:        resp.Choices[0].Message.Content,
		TokensUsed:     tokensUsed,
		ModelUsed:      req.Model,
		ProcessingTime: processingTime,
		FinishReason:   string(resp.Choices[0].FinishReason),
		Confidence:     0.9,
	}, nil
}

// GenerateStream sends a streaming request to Qwen API
func (p *QwenProviderImpl) GenerateStream(ctx context.Context, req domain.AIRequest, onChunk func(chunk domain.StreamChunk)) error {
	p.log.Info(fmt.Sprintf("Starting streaming request to Qwen API with model: %s", req.Model))

	completionReq := common.BuildCompletionRequest(req, true)
	stream, err := p.client.CreateChatCompletionStream(ctx, completionReq)
	if err != nil {
		p.log.Error(fmt.Sprintf("Qwen API stream request failed: %v", err))
		domainErr := common.HandleOpenAIError(err)
		onChunk(domain.StreamChunk{Done: true, Error: domainErr.Error()})
		return domainErr
	}
	defer stream.Close()
	return common.StreamProcessor(stream, onChunk, p.log)
}

// GetProviderInfo returns information about the Qwen provider
func (p *QwenProviderImpl) GetProviderInfo() domain.ProviderInfo {
	return domain.ProviderInfo{
		Name:         "Qwen",
		Version:      "1.0",
		Capabilities: []string{"chat", "completion", "code_generation", "large_context"},
		Limitations:  []string{"rate_limited"},
		SupportedModels: []string{
			"qwen-coder-plus-latest",
			"qwen-coder-plus",
			"qwen-plus-latest",
			"qwen-turbo-latest",
			"qwen-max",
		},
	}
}

// ValidateRequest validates the request parameters
func (p *QwenProviderImpl) ValidateRequest(req domain.AIRequest) error {
	return common.ValidateRequestStrict(req)
}

// EstimateTokens estimates the number of tokens in the request
func (p *QwenProviderImpl) EstimateTokens(req domain.AIRequest) (int, error) {
	return common.EstimateTokens(req)
}

// GetPricing returns pricing information for Qwen models
func (p *QwenProviderImpl) GetPricing(model string) domain.PricingInfo {
	return common.GetPricingFromTable(model, "CNY", common.QwenPricingTable, common.QwenDefaultPricing)
}

// GetMaxContextTokens returns the maximum context size for a model
func (p *QwenProviderImpl) GetMaxContextTokens(model string) int {
	switch model {
	case "qwen-coder-plus-latest", "qwen-coder-plus":
		return 1000000 // 1M tokens
	case "qwen-coder-turbo-latest", "qwen-coder-turbo":
		return 131072 // 128K tokens
	case "qwen-plus-latest", "qwen-plus":
		return 131072
	case "qwen-turbo-latest", "qwen-turbo":
		return 131072
	case "qwen-max", "qwen-max-latest":
		return 32768
	default:
		return 32768
	}
}
