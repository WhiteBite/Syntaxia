package ai

import (
	"context"
	"errors"
	"fmt"
	"syntaxia/domain"
	"syntaxia/infrastructure/ai/common"
	"sort"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GeminiProviderImpl struct {
	log    domain.Logger
	apiKey string
}

func NewGemini(apiKey, host string, log domain.Logger) (domain.AIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("gemini API key is required")
	}
	// The host parameter is ignored for the official Gemini client, but kept for signature compatibility.
	return &GeminiProviderImpl{
		log:    log,
		apiKey: apiKey,
	}, nil
}

func (p *GeminiProviderImpl) ListModels(ctx context.Context) ([]string, error) {
	p.log.Info("Requesting model list from Gemini API...")
	client, err := genai.NewClient(ctx, option.WithAPIKey(p.apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client for model listing: %w", err)
	}
	defer client.Close()

	var models []string
	iter := client.ListModels(ctx)
	for {
		m, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), "API_KEY_INVALID") {
				return nil, domain.ErrInvalidAPIKey
			}
			return nil, fmt.Errorf("failed to iterate models: %w", err)
		}

		isSupported := false
		for _, method := range m.SupportedGenerationMethods {
			if method == "generateContent" {
				isSupported = true
				break
			}
		}
		if isSupported {
			models = append(models, m.Name)
		}
	}

	sort.Strings(models)
	p.log.Info(fmt.Sprintf("Received %d supported models from Gemini.", len(models)))
	return models, nil
}

func (p *GeminiProviderImpl) Generate(ctx context.Context, req domain.AIRequest) (domain.AIResponse, error) {
	startTime := time.Now()
	p.log.Info(fmt.Sprintf("Sending request to Gemini API with model: %s", req.Model))

	client, err := genai.NewClient(ctx, option.WithAPIKey(p.apiKey))
	if err != nil {
		return domain.AIResponse{}, fmt.Errorf("failed to create gemini client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(req.Model)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(req.SystemPrompt)},
	}

	// Настраиваем параметры генерации
	if req.Temperature > 0 {
		temp := float32(req.Temperature)
		model.Temperature = &temp
	}
	if req.MaxTokens > 0 {
		maxTokens := int32(req.MaxTokens) //nolint:gosec // MaxTokens is always positive and within int32 range
		model.MaxOutputTokens = &maxTokens
	}
	if req.TopP > 0 {
		topP := float32(req.TopP)
		model.TopP = &topP
	}

	resp, err := model.GenerateContent(ctx, genai.Text(req.UserPrompt))
	if err != nil {
		p.log.Error(fmt.Sprintf("Gemini API request failed: %v", err))
		return domain.AIResponse{}, err
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return domain.AIResponse{}, fmt.Errorf("no content returned from Gemini API")
	}

	firstPart := resp.Candidates[0].Content.Parts[0]
	if text, ok := firstPart.(genai.Text); ok {
		processingTime := time.Since(startTime)

		// Подсчитываем токены (примерная оценка)
		tokensUsed := len(string(text)) / 4

		return domain.AIResponse{
			Content:        string(text),
			TokensUsed:     tokensUsed,
			ModelUsed:      req.Model,
			ProcessingTime: processingTime,
			FinishReason:   "stop",
			Confidence:     0.9,
		}, nil
	}

	return domain.AIResponse{}, fmt.Errorf("unsupported content type returned from Gemini: %T", firstPart)
}

func (p *GeminiProviderImpl) GetProviderInfo() domain.ProviderInfo {
	return domain.ProviderInfo{
		Name:            "Google Gemini",
		Version:         "1.0",
		Capabilities:    []string{"chat", "completion", "embeddings"},
		Limitations:     []string{"rate_limited", "token_limited"},
		SupportedModels: []string{"gemini-pro", "gemini-pro-vision", "gemini-1.5-pro"},
	}
}

func (p *GeminiProviderImpl) ValidateRequest(req domain.AIRequest) error {
	return common.ValidateRequestBasic(req)
}

func (p *GeminiProviderImpl) EstimateTokens(req domain.AIRequest) (int, error) {
	return common.EstimateTokens(req)
}

func (p *GeminiProviderImpl) GetPricing(model string) domain.PricingInfo {
	return common.GetPricingFromTable(model, "USD", common.GeminiPricingTable, common.GeminiDefaultPricing)
}

// GenerateStream implements streaming for Gemini
func (p *GeminiProviderImpl) GenerateStream(ctx context.Context, req domain.AIRequest, onChunk func(chunk domain.StreamChunk)) error {
	p.log.Info(fmt.Sprintf("Starting streaming request to Gemini API with model: %s", req.Model))

	client, err := genai.NewClient(ctx, option.WithAPIKey(p.apiKey))
	if err != nil {
		onChunk(domain.StreamChunk{Done: true, Error: err.Error()})
		return fmt.Errorf("failed to create gemini client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(req.Model)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(req.SystemPrompt)},
	}

	if req.Temperature > 0 {
		temp := float32(req.Temperature)
		model.Temperature = &temp
	}
	if req.MaxTokens > 0 {
		maxTokens := int32(req.MaxTokens) //nolint:gosec // MaxTokens is always positive and within int32 range
		model.MaxOutputTokens = &maxTokens
	}
	if req.TopP > 0 {
		topP := float32(req.TopP)
		model.TopP = &topP
	}

	iter := model.GenerateContentStream(ctx, genai.Text(req.UserPrompt))
	totalTokens := 0

	for {
		resp, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			onChunk(domain.StreamChunk{Done: true, TokensUsed: totalTokens, FinishReason: "stop"})
			return nil
		}
		if err != nil {
			p.log.Error(fmt.Sprintf("Gemini stream error: %v", err))
			onChunk(domain.StreamChunk{Done: true, Error: err.Error()})
			return err
		}

		if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
			for _, part := range resp.Candidates[0].Content.Parts {
				if text, ok := part.(genai.Text); ok {
					content := string(text)
					totalTokens += len(content) / 4
					onChunk(domain.StreamChunk{
						Content: content,
						Done:    false,
					})
				}
			}
		}
	}
}
