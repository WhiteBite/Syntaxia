package embeddings

import (
	"context"
	"fmt"
	"syntaxia/domain"

	"github.com/sashabaranov/go-openai"
)

// OpenAIEmbeddingProvider implements EmbeddingProvider using OpenAI API
type OpenAIEmbeddingProvider struct {
	client *openai.Client
	model  domain.EmbeddingModel
	log    domain.Logger
}

// NewOpenAIEmbeddingProvider creates a new OpenAI embedding provider
func NewOpenAIEmbeddingProvider(apiKey string, model domain.EmbeddingModel, log domain.Logger) (*OpenAIEmbeddingProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	client := openai.NewClient(apiKey)

	return &OpenAIEmbeddingProvider{
		client: client,
		model:  model,
		log:    log,
	}, nil
}

// GenerateEmbeddings generates embeddings for the given texts
func (p *OpenAIEmbeddingProvider) GenerateEmbeddings(ctx context.Context, req domain.EmbeddingRequest) (*domain.EmbeddingResponse, error) {
	if err := p.ValidateRequest(req); err != nil {
		return nil, err
	}

	model := req.Model
	if model == "" {
		model = p.model
	}

	// Map domain model to OpenAI model
	openaiModel := mapToOpenAIModel(model)

	p.log.Info(fmt.Sprintf("Generating embeddings for %d texts using model %s", len(req.Texts), openaiModel))

	// OpenAI supports batch embedding
	resp, err := p.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: req.Texts,
		Model: openaiModel,
	})
	if err != nil {
		p.log.Error(fmt.Sprintf("Failed to generate embeddings: %v", err))
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Convert response
	embeddings := make([]domain.EmbeddingVector, len(resp.Data))
	for i, data := range resp.Data {
		embeddings[i] = data.Embedding
	}

	p.log.Info(fmt.Sprintf("Generated %d embeddings, tokens used: %d", len(embeddings), resp.Usage.TotalTokens))

	return &domain.EmbeddingResponse{
		Embeddings: embeddings,
		Model:      model,
		TokensUsed: resp.Usage.TotalTokens,
	}, nil
}

// GetModelInfo returns information about the embedding model
func (p *OpenAIEmbeddingProvider) GetModelInfo() domain.EmbeddingModelInfo {
	return domain.EmbeddingModelInfo{
		Model:      p.model,
		Dimensions: p.model.Dimensions(),
		MaxTokens:  8191, // OpenAI embedding models support up to 8191 tokens
		Provider:   "openai",
	}
}

// ValidateRequest validates the embedding request
func (p *OpenAIEmbeddingProvider) ValidateRequest(req domain.EmbeddingRequest) error {
	if len(req.Texts) == 0 {
		return fmt.Errorf("at least one text is required")
	}

	// OpenAI has a limit on batch size
	if len(req.Texts) > 2048 {
		return fmt.Errorf("batch size exceeds maximum of 2048")
	}

	// Check for empty texts
	for i, text := range req.Texts {
		if text == "" {
			return fmt.Errorf("text at index %d is empty", i)
		}
	}

	return nil
}

func mapToOpenAIModel(model domain.EmbeddingModel) openai.EmbeddingModel {
	switch model {
	case domain.EmbeddingModelOpenAI:
		return openai.AdaEmbeddingV2
	case domain.EmbeddingModelOpenAI3S:
		return openai.SmallEmbedding3
	case domain.EmbeddingModelOpenAI3L:
		return openai.LargeEmbedding3
	default:
		return openai.AdaEmbeddingV2
	}
}
