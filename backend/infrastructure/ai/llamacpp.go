package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"syntaxia/domain"
	"strings"
	"time"
)

// LlamaCppClient представляет клиент для llama.cpp server
type LlamaCppClient struct {
	baseURL    string
	httpClient *http.Client
	log        domain.Logger
}

// LlamaCppConfig конфигурация для llama.cpp клиента
type LlamaCppConfig struct {
	BaseURL       string        `json:"base_url"`
	Timeout       time.Duration `json:"timeout"`
	MaxTokens     int           `json:"max_tokens"`
	Temperature   float64       `json:"temperature"`
	TopP          float64       `json:"top_p"`
	TopK          int           `json:"top_k"`
	RepeatPenalty float64       `json:"repeat_penalty"`
}

// LlamaCppRequest запрос к llama.cpp server
type LlamaCppRequest struct {
	Prompt        string                 `json:"prompt"`
	Stream        bool                   `json:"stream"`
	NPredict      int                    `json:"n_predict"`
	Temperature   float64                `json:"temperature"`
	TopP          float64                `json:"top_p"`
	TopK          int                    `json:"top_k"`
	RepeatPenalty float64                `json:"repeat_penalty"`
	Stop          []string               `json:"stop"`
	Grammar       string                 `json:"grammar,omitempty"`
	GrammarRules  []string               `json:"grammar_rules,omitempty"`
	Options       map[string]interface{} `json:"options,omitempty"`
}

// LlamaCppResponse ответ от llama.cpp server
type LlamaCppResponse struct {
	Content    string `json:"content"`
	Stop       bool   `json:"stop"`
	StopReason string `json:"stop_reason"`
	Timings    struct {
		PredN    int     `json:"pred_n"`
		PredMS   float64 `json:"pred_ms"`
		PromptN  int     `json:"prompt_n"`
		PromptMS float64 `json:"prompt_ms"`
	} `json:"timings"`
	TokensEvaluated int  `json:"tokens_evaluated"`
	TokensPredicted int  `json:"tokens_predicted"`
	Truncated       bool `json:"truncated"`
}

// LlamaCppStreamResponse потоковый ответ от llama.cpp server
type LlamaCppStreamResponse struct {
	Content    string `json:"content"`
	Stop       bool   `json:"stop"`
	StopReason string `json:"stop_reason"`
	Timings    struct {
		PredN    int     `json:"pred_n"`
		PredMS   float64 `json:"pred_ms"`
		PromptN  int     `json:"prompt_n"`
		PromptMS float64 `json:"prompt_ms"`
	} `json:"timings"`
	TokensEvaluated int  `json:"tokens_evaluated"`
	TokensPredicted int  `json:"tokens_predicted"`
	Truncated       bool `json:"truncated"`
}

// NewLlamaCppClient создает новый клиент для llama.cpp server
func NewLlamaCppClient(config LlamaCppConfig, log domain.Logger) *LlamaCppClient {
	if config.BaseURL == "" {
		config.BaseURL = domain.LlamaCppDefaultHost
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 2048
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.TopP == 0 {
		config.TopP = 0.9
	}
	if config.TopK == 0 {
		config.TopK = 40
	}
	if config.RepeatPenalty == 0 {
		config.RepeatPenalty = 1.1
	}

	return &LlamaCppClient{
		baseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		log: log,
	}
}

// GenerateText генерирует текст с помощью llama.cpp
func (c *LlamaCppClient) GenerateText(ctx context.Context, prompt string, options map[string]interface{}) (string, error) {
	request := LlamaCppRequest{
		Prompt:        prompt,
		Stream:        false,
		NPredict:      2048,
		Temperature:   0.7,
		TopP:          0.9,
		TopK:          40,
		RepeatPenalty: 1.1,
		Stop:          []string{"</s>", "Human:", "Assistant:"},
		Options:       options,
	}

	// Применяем пользовательские опции
	if temp, ok := options["temperature"].(float64); ok {
		request.Temperature = temp
	}
	if topP, ok := options["top_p"].(float64); ok {
		request.TopP = topP
	}
	if topK, ok := options["top_k"].(int); ok {
		request.TopK = topK
	}
	if nPredict, ok := options["n_predict"].(int); ok {
		request.NPredict = nPredict
	}

	response, err := c.makeRequest(ctx, request)
	if err != nil {
		return "", fmt.Errorf("llama.cpp request failed: %w", err)
	}

	return response.Content, nil
}

// GenerateStructuredJSON генерирует структурированный JSON с помощью GBNF грамматики
func (c *LlamaCppClient) GenerateStructuredJSON(ctx context.Context, prompt string, grammar string, options map[string]interface{}) ([]byte, error) {
	request := LlamaCppRequest{
		Prompt:        prompt,
		Stream:        false,
		NPredict:      2048,
		Temperature:   0.1, // Низкая температура для более детерминированного вывода
		TopP:          0.9,
		TopK:          40,
		RepeatPenalty: 1.1,
		Stop:          []string{"</s>", "Human:", "Assistant:"},
		Grammar:       grammar,
		Options:       options,
	}

	// Применяем пользовательские опции
	if temp, ok := options["temperature"].(float64); ok {
		request.Temperature = temp
	}
	if topP, ok := options["top_p"].(float64); ok {
		request.TopP = topP
	}
	if topK, ok := options["top_k"].(int); ok {
		request.TopK = topK
	}
	if nPredict, ok := options["n_predict"].(int); ok {
		request.NPredict = nPredict
	}

	response, err := c.makeRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("llama.cpp structured request failed: %w", err)
	}

	// Очищаем ответ от лишних символов
	content := strings.TrimSpace(response.Content)

	// Удаляем возможные префиксы
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// Проверяем, что это валидный JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(content), &jsonData); err != nil {
		return nil, fmt.Errorf("invalid JSON response from llama.cpp: %w", err)
	}

	return []byte(content), nil
}

// GenerateEditsJSON генерирует Edits JSON с помощью GBNF грамматики
func (c *LlamaCppClient) GenerateEditsJSON(ctx context.Context, prompt string, options map[string]interface{}) ([]byte, error) {
	// Загружаем GBNF грамматику для Edits JSON
	grammar, err := c.loadEditsGBNF()
	if err != nil {
		return nil, fmt.Errorf("failed to load Edits GBNF grammar: %w", err)
	}

	return c.GenerateStructuredJSON(ctx, prompt, grammar, options)
}

// StreamText генерирует текст в потоковом режиме
func (c *LlamaCppClient) StreamText(ctx context.Context, prompt string, options map[string]interface{}) (<-chan string, error) {
	request := LlamaCppRequest{
		Prompt:        prompt,
		Stream:        true,
		NPredict:      2048,
		Temperature:   0.7,
		TopP:          0.9,
		TopK:          40,
		RepeatPenalty: 1.1,
		Stop:          []string{"</s>", "Human:", "Assistant:"},
		Options:       options,
	}

	// Применяем пользовательские опции
	if temp, ok := options["temperature"].(float64); ok {
		request.Temperature = temp
	}
	if topP, ok := options["top_p"].(float64); ok {
		request.TopP = topP
	}
	if topK, ok := options["top_k"].(int); ok {
		request.TopK = topK
	}
	if nPredict, ok := options["n_predict"].(int); ok {
		request.NPredict = nPredict
	}

	response, err := c.makeStreamRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("llama.cpp stream request failed: %w", err)
	}

	// Создаем канал для потоковой передачи
	textChan := make(chan string, 100)

	go func() {
		defer close(textChan)
		for {
			select {
			case <-ctx.Done():
				return
			case chunk, ok := <-response:
				if !ok {
					return
				}
				textChan <- chunk
			}
		}
	}()

	return textChan, nil
}

// makeRequest выполняет HTTP запрос к llama.cpp server
func (c *LlamaCppClient) makeRequest(ctx context.Context, request LlamaCppRequest) (*LlamaCppResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/completion", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	c.log.Info(fmt.Sprintf("Making request to llama.cpp server: %s", c.baseURL))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("llama.cpp server returned status %d: %s", resp.StatusCode, string(body))
	}

	var response LlamaCppResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// createStreamHTTPRequest creates HTTP request for streaming
func (c *LlamaCppClient) createStreamHTTPRequest(ctx context.Context, request LlamaCppRequest) (*http.Response, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/completion", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	c.log.Info(fmt.Sprintf("Making stream request to llama.cpp server: %s", c.baseURL))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("llama.cpp server returned status %d: %s", resp.StatusCode, string(body))
	}
	return resp, nil
}

// streamResponseReader reads stream responses and sends to channel
func (c *LlamaCppClient) streamResponseReader(ctx context.Context, resp *http.Response, chunkChan chan<- string) {
	defer resp.Body.Close()
	defer close(chunkChan)

	decoder := json.NewDecoder(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var streamResp LlamaCppStreamResponse
			if err := decoder.Decode(&streamResp); err != nil {
				if err != io.EOF {
					c.log.Error(fmt.Sprintf("Failed to decode stream response: %v", err))
				}
				return
			}
			chunkChan <- streamResp.Content
			if streamResp.Stop {
				return
			}
		}
	}
}

// makeStreamRequest выполняет потоковый HTTP запрос к llama.cpp server
func (c *LlamaCppClient) makeStreamRequest(ctx context.Context, request LlamaCppRequest) (<-chan string, error) {
	resp, err := c.createStreamHTTPRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	chunkChan := make(chan string, 100)
	go c.streamResponseReader(ctx, resp, chunkChan)
	return chunkChan, nil
}

// loadEditsGBNF загружает GBNF грамматику для Edits JSON
func (c *LlamaCppClient) loadEditsGBNF() (string, error) {
	// В реальной реализации здесь нужно загрузить файл docs/schemas/edits.gbnf
	// Для демонстрации используем упрощенную грамматику
	grammar := `
root ::= edits-json

edits-json ::= "{" space schema-version space "," space edits space "}"

schema-version ::= "\"schemaVersion\"" space ":" space "\"1.0\""

edits ::= "\"edits\"" space ":" space "[" space edit-list space "]"

edit-list ::= edit (space "," space edit)*

edit ::= "{" space edit-fields space "}"

edit-fields ::= edit-kind space "," space path space "," space language space "," space operation space "," space post

edit-kind ::= "\"kind\"" space ":" space kind-value

kind-value ::= "\"recipeOp\"" | "\"workspaceEdit\"" | "\"astOp\"" | "\"anchorPatch\"" | "\"fullFile\""

path ::= "\"path\"" space ":" space string

language ::= "\"language\"" space ":" space string

operation ::= "\"operation\"" space ":" space operation-value

operation-value ::= "{" space operation-fields space "}"

operation-fields ::= engine space "," space action space "," space params

engine ::= "\"engine\"" space ":" space string

action ::= "\"action\"" space ":" space string

params ::= "\"params\"" space ":" space "{" space param-list space "}"

param-list ::= param (space "," space param)*

param ::= string space ":" space param-value

param-value ::= string | number | boolean | "{" space param-list space "}" | "[" space param-list space "]"

post ::= "\"post\"" space ":" space "{" space post-fields space "}"

post-fields ::= formatters

formatters ::= "\"formatters\"" space ":" space "[" space formatter-list space "]"

formatter-list ::= string (space "," space string)*

string ::= "\"" [^"]* "\""

number ::= [0-9]+ ("." [0-9]+)?

boolean ::= "true" | "false"

space ::= [ \t\n\r]*
`

	return grammar, nil
}

// HealthCheck проверяет доступность llama.cpp server
func (c *LlamaCppClient) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("llama.cpp server health check failed with status %d", resp.StatusCode)
	}

	c.log.Info("Llama.cpp server health check passed")
	return nil
}

// GetModelInfo получает информацию о модели
func (c *LlamaCppClient) GetModelInfo(ctx context.Context) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/model", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create model info request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("model info request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("model info request failed with status %d", resp.StatusCode)
	}

	var modelInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&modelInfo); err != nil {
		return nil, fmt.Errorf("failed to decode model info: %w", err)
	}

	return modelInfo, nil
}

// LoadGBNFGrammar загружает GBNF грамматику из файла
func (c *LlamaCppClient) LoadGBNFGrammar(grammarPath string) (string, error) {
	content, err := os.ReadFile(grammarPath)
	if err != nil {
		return "", fmt.Errorf("failed to read GBNF grammar file %s: %w", grammarPath, err)
	}
	return string(content), nil
}

// GenerateWithGBNF генерирует ответ с использованием GBNF грамматики
func (c *LlamaCppClient) GenerateWithGBNF(ctx context.Context, prompt string, grammar string, options map[string]interface{}) (*LlamaCppResponse, error) {
	request := LlamaCppRequest{
		Prompt:        prompt,
		Stream:        false,
		NPredict:      2048,
		Temperature:   0.7,
		TopP:          0.9,
		TopK:          40,
		RepeatPenalty: 1.1,
		Grammar:       grammar,
		Options:       options,
	}

	return c.makeRequest(ctx, request)
}

// ValidateGBNFResponse проверяет корректность ответа согласно GBNF грамматике
func (c *LlamaCppClient) ValidateGBNFResponse(response string, grammar string) error {
	// Базовая валидация JSON структуры
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(response), &jsonData); err != nil {
		return fmt.Errorf("invalid JSON response: %w", err)
	}

	// Проверяем наличие обязательных полей для Edits JSON
	if edits, exists := jsonData["edits"]; !exists {
		return fmt.Errorf("missing required field 'edits' in response")
	} else if editsArray, ok := edits.([]interface{}); !ok {
		return fmt.Errorf("field 'edits' must be an array")
	} else if len(editsArray) == 0 {
		return fmt.Errorf("edits array cannot be empty")
	} else {
		// Проверяем структуру каждого edit
		for i, edit := range editsArray {
			if editMap, ok := edit.(map[string]interface{}); !ok {
				return fmt.Errorf("edit at index %d must be an object", i)
			} else {
				requiredFields := []string{"kind", "path", "language", "operation", "post"}
				for _, field := range requiredFields {
					if _, exists := editMap[field]; !exists {
						return fmt.Errorf("edit at index %d missing required field '%s'", i, field)
					}
				}
			}
		}
	}

	c.log.Info("GBNF response validation passed")
	return nil
}
