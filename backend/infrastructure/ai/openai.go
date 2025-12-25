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

type OpenAIProviderImpl struct {
	client *openai.Client
	log    domain.Logger
}

func NewOpenAI(apiKey, host string, log domain.Logger) (domain.AIProvider, error) {
	config := openai.DefaultConfig(apiKey)
	if host != "" {
		config.BaseURL = host
	}
	client := openai.NewClientWithConfig(config)
	return &OpenAIProviderImpl{
		client: client,
		log:    log,
	}, nil
}

func (p *OpenAIProviderImpl) ListModels(ctx context.Context) ([]string, error) {
	p.log.Info("Requesting model list from OpenAI compatible API...")
	resp, err := p.client.ListModels(ctx)
	if err != nil {
		p.log.Error(fmt.Sprintf("Error getting model list: %v", err))
		if domainErr := common.HandleOpenAIError(err); !errors.Is(domainErr, err) {
			return nil, domainErr
		}
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	models := make([]string, 0, len(resp.Models))
	for _, model := range resp.Models {
		models = append(models, model.ID)
	}
	sort.Strings(models)
	p.log.Info(fmt.Sprintf("Received %d models.", len(models)))
	return models, nil
}

func (p *OpenAIProviderImpl) Generate(ctx context.Context, req domain.AIRequest) (domain.AIResponse, error) {
	startTime := time.Now()
	p.log.Info(fmt.Sprintf("Sending request to OpenAI compatible API with model: %s", req.Model))

	completionReq := common.BuildCompletionRequest(req, false)
	resp, err := p.client.CreateChatCompletion(ctx, completionReq)

	if err != nil {
		p.log.Error(fmt.Sprintf("OpenAI API request failed: %v", err))
		return domain.AIResponse{}, err
	}

	if len(resp.Choices) == 0 {
		return domain.AIResponse{}, fmt.Errorf("no choices returned from OpenAI API")
	}

	processingTime := time.Since(startTime)

	// Подсчитываем токены (примерная оценка)
	tokensUsed := resp.Usage.TotalTokens

	return domain.AIResponse{
		Content:        resp.Choices[0].Message.Content,
		TokensUsed:     tokensUsed,
		ModelUsed:      req.Model,
		ProcessingTime: processingTime,
		FinishReason:   string(resp.Choices[0].FinishReason),
		Confidence:     0.9, // Базовая оценка
	}, nil
}

func (p *OpenAIProviderImpl) GetProviderInfo() domain.ProviderInfo {
	return domain.ProviderInfo{
		Name:            "OpenAI",
		Version:         "1.0",
		Capabilities:    []string{"chat", "completion", "embeddings"},
		Limitations:     []string{"rate_limited", "token_limited"},
		SupportedModels: []string{"gpt-4", "gpt-3.5-turbo", "gpt-4-turbo"},
	}
}

func (p *OpenAIProviderImpl) ValidateRequest(req domain.AIRequest) error {
	return common.ValidateRequestStrict(req)
}

func (p *OpenAIProviderImpl) EstimateTokens(req domain.AIRequest) (int, error) {
	return common.EstimateTokens(req)
}

func (p *OpenAIProviderImpl) GetPricing(model string) domain.PricingInfo {
	return common.GetPricingFromTable(model, "USD", common.OpenAIPricingTable, common.OpenAIDefaultPricing)
}

// GenerateStream sends a streaming request to OpenAI API
func (p *OpenAIProviderImpl) GenerateStream(ctx context.Context, req domain.AIRequest, onChunk func(chunk domain.StreamChunk)) error {
	p.log.Info(fmt.Sprintf("Starting streaming request to OpenAI API with model: %s", req.Model))

	completionReq := common.BuildCompletionRequest(req, true)
	stream, err := p.client.CreateChatCompletionStream(ctx, completionReq)
	if err != nil {
		p.log.Error(fmt.Sprintf("OpenAI API stream request failed: %v", err))
		onChunk(domain.StreamChunk{Done: true, Error: err.Error()})
		return err
	}
	defer stream.Close()
	return common.StreamProcessor(stream, onChunk, p.log)
}
