package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"syntaxia/domain"
	"syntaxia/infrastructure/ai/common"
	"sort"
	"time"
)

// LocalAIProviderImpl реализует провайдер для LocalAI
type LocalAIProviderImpl struct {
	client *http.Client
	host   string
	apiKey string
	log    domain.Logger
}

// LocalAIRequest представляет запрос к LocalAI API
type LocalAIRequest struct {
	Model       string           `json:"model"`
	Messages    []LocalAIMessage `json:"messages"`
	Temperature float64          `json:"temperature,omitempty"`
	MaxTokens   int              `json:"max_tokens,omitempty"`
	TopP        float64          `json:"top_p,omitempty"`
	Grammar     *Grammar         `json:"grammar,omitempty"`
}

// LocalAIMessage представляет сообщение в LocalAI API
type LocalAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LocalAIResponse представляет ответ от LocalAI API
type LocalAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Grammar представляет грамматику для структурированного вывода
type Grammar struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// NewLocalAI создает новый провайдер LocalAI
func NewLocalAI(apiKey, host string, log domain.Logger) (domain.AIProvider, error) {
	if host == "" {
		host = domain.LocalAIDefaultHost
	}

	return &LocalAIProviderImpl{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		host:   host,
		apiKey: apiKey,
		log:    log,
	}, nil
}

// ListModels возвращает список доступных моделей
func (p *LocalAIProviderImpl) ListModels(ctx context.Context) ([]string, error) {
	p.log.Info("Requesting model list from LocalAI...")

	url := fmt.Sprintf("%s/models", p.host)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get models: %s, status: %d", string(body), resp.StatusCode)
	}

	var modelsResponse struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&modelsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode models response: %w", err)
	}

	models := make([]string, 0, len(modelsResponse.Data))
	for _, model := range modelsResponse.Data {
		models = append(models, model.ID)
	}

	sort.Strings(models)
	p.log.Info(fmt.Sprintf("Received %d models from LocalAI", len(models)))
	return models, nil
}

// Generate выполняет запрос к LocalAI
func (p *LocalAIProviderImpl) Generate(ctx context.Context, req domain.AIRequest) (domain.AIResponse, error) {
	startTime := time.Now()
	p.log.Info(fmt.Sprintf("Sending request to LocalAI with model: %s", req.Model))

	// Создаем сообщения для LocalAI
	messages := []LocalAIMessage{
		{
			Role:    "system",
			Content: req.SystemPrompt,
		},
		{
			Role:    "user",
			Content: req.UserPrompt,
		},
	}

	// Создаем запрос к LocalAI
	localAIReq := LocalAIRequest{
		Model:       req.Model,
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopP:        req.TopP,
	}

	// Добавляем грамматику если нужно
	if req.Grammar != "" {
		localAIReq.Grammar = &Grammar{
			Type:  "json",
			Value: req.Grammar,
		}
	}

	// Сериализуем запрос
	jsonData, err := json.Marshal(localAIReq)
	if err != nil {
		return domain.AIResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Создаем HTTP запрос
	url := fmt.Sprintf("%s/chat/completions", p.host)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return domain.AIResponse{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	// Выполняем запрос
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return domain.AIResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.AIResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return domain.AIResponse{}, fmt.Errorf("LocalAI request failed: %s, status: %d", string(body), resp.StatusCode)
	}

	// Парсим ответ
	var localAIResp LocalAIResponse
	if err := json.Unmarshal(body, &localAIResp); err != nil {
		return domain.AIResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Проверяем, что есть выбор
	if len(localAIResp.Choices) == 0 {
		return domain.AIResponse{}, errors.New("no choices in LocalAI response")
	}

	choice := localAIResp.Choices[0]
	duration := time.Since(startTime).Seconds()

	p.log.Info(fmt.Sprintf("LocalAI request completed in %.2fs", duration))

	return domain.AIResponse{
		Content:        choice.Message.Content,
		ModelUsed:      localAIResp.Model,
		TokensUsed:     localAIResp.Usage.TotalTokens,
		ProcessingTime: time.Duration(duration * float64(time.Second)),
		FinishReason:   choice.FinishReason,
	}, nil
}

// ValidateGrammar проверяет корректность грамматики
func (p *LocalAIProviderImpl) ValidateGrammar(grammar string) error {
	if grammar == "" {
		return nil
	}

	// Простая проверка - грамматика должна быть валидным JSON
	var test interface{}
	return json.Unmarshal([]byte(grammar), &test)
}

// GetSupportedGrammarTypes возвращает поддерживаемые типы грамматик
func (p *LocalAIProviderImpl) GetSupportedGrammarTypes() []string {
	return []string{"json", "gbnf"}
}

// GetProviderInfo возвращает информацию о провайдере
func (p *LocalAIProviderImpl) GetProviderInfo() domain.ProviderInfo {
	return domain.ProviderInfo{
		Name:            "LocalAI",
		Version:         "1.0",
		Capabilities:    []string{"local-inference", "grammar-guided", "json-structured"},
		Limitations:     []string{"requires-local-server", "limited-model-selection"},
		SupportedModels: []string{"local-model"},
	}
}

// ValidateRequest проверяет корректность запроса
func (p *LocalAIProviderImpl) ValidateRequest(req domain.AIRequest) error {
	if req.Model == "" {
		return fmt.Errorf("model is required")
	}
	if req.SystemPrompt == "" && req.UserPrompt == "" {
		return fmt.Errorf("at least one prompt is required")
	}
	if req.Grammar != "" {
		return p.ValidateGrammar(req.Grammar)
	}
	return nil
}

// EstimateTokens оценивает количество токенов в запросе
func (p *LocalAIProviderImpl) EstimateTokens(req domain.AIRequest) (int, error) {
	return common.EstimateTokens(req)
}

// GetPricing возвращает информацию о стоимости
func (p *LocalAIProviderImpl) GetPricing(model string) domain.PricingInfo {
	return domain.PricingInfo{
		InputTokensPer1K:  0.0, // LocalAI бесплатный
		OutputTokensPer1K: 0.0,
		Currency:          "USD",
		Model:             model,
	}
}

// GenerateStream implements streaming for LocalAI (fallback to non-streaming)
func (p *LocalAIProviderImpl) GenerateStream(ctx context.Context, req domain.AIRequest, onChunk func(chunk domain.StreamChunk)) error {
	// LocalAI may not support streaming, fallback to regular generation
	resp, err := p.Generate(ctx, req)
	if err != nil {
		onChunk(domain.StreamChunk{Done: true, Error: err.Error()})
		return err
	}

	// Send content in chunks to simulate streaming
	// Use runes to properly handle UTF-8 characters (like Cyrillic)
	runes := []rune(resp.Content)
	chunkSize := 50
	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		onChunk(domain.StreamChunk{
			Content: string(runes[i:end]),
			Done:    false,
		})
	}

	onChunk(domain.StreamChunk{
		Done:         true,
		TokensUsed:   resp.TokensUsed,
		FinishReason: "stop",
	})
	return nil
}
