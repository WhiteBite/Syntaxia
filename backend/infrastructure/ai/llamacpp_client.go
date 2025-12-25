package ai

import (
	"context"
	"syntaxia/domain"
)

// LlamaCppClientAdapter адаптирует ai.LlamaCppClient к domain.LLMClient
type LlamaCppClientAdapter struct {
	client *LlamaCppClient
}

// NewLlamaCppClientAdapter создает новый адаптер для LlamaCppClient
func NewLlamaCppClientAdapter(client *LlamaCppClient) *LlamaCppClientAdapter {
	return &LlamaCppClientAdapter{
		client: client,
	}
}

// HealthCheck implements the LLMClient interface
func (a *LlamaCppClientAdapter) HealthCheck(ctx context.Context) error {
	return a.client.HealthCheck(ctx)
}

// GenerateWithGBNF implements the LLMClient interface
func (a *LlamaCppClientAdapter) GenerateWithGBNF(ctx context.Context, prompt string, grammar string, options map[string]interface{}) (*domain.LlamaCppResponse, error) {
	response, err := a.client.GenerateWithGBNF(ctx, prompt, grammar, options)
	if err != nil {
		return nil, err
	}

	// Convert ai.LlamaCppResponse to domain.LlamaCppResponse
	domainResponse := &domain.LlamaCppResponse{
		Content:    response.Content,
		Stop:       response.Stop,
		StopReason: response.StopReason,
		Timings: struct {
			PredN    int     `json:"pred_n"`
			PredMS   float64 `json:"pred_ms"`
			PromptN  int     `json:"prompt_n"`
			PromptMS float64 `json:"prompt_ms"`
		}{
			PredN:    response.Timings.PredN,
			PredMS:   response.Timings.PredMS,
			PromptN:  response.Timings.PromptN,
			PromptMS: response.Timings.PromptMS,
		},
		TokensEvaluated: response.TokensEvaluated,
		TokensPredicted: response.TokensPredicted,
		Truncated:       response.Truncated,
	}

	return domainResponse, nil
}

// GenerateEditsJSON implements the LLMClient interface
func (a *LlamaCppClientAdapter) GenerateEditsJSON(ctx context.Context, prompt string, options map[string]interface{}) ([]byte, error) {
	return a.client.GenerateEditsJSON(ctx, prompt, options)
}

// GetModelInfo implements the LLMClient interface
func (a *LlamaCppClientAdapter) GetModelInfo(ctx context.Context) (map[string]interface{}, error) {
	return a.client.GetModelInfo(ctx)
}
