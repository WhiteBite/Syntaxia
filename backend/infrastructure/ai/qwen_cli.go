package ai

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syntaxia/domain"
	"time"
)

// QwenCLIProviderImpl implements domain.AIProvider for Qwen Code CLI
type QwenCLIProviderImpl struct {
	log      domain.Logger
	settings domain.QwenCLISettings
}

// NewQwenCLI creates a new Qwen CLI provider instance
func NewQwenCLI(log domain.Logger) (domain.AIProvider, error) {
	return NewQwenCLIWithSettings(log, domain.DefaultQwenCLISettings())
}

// NewQwenCLIWithSettings creates a new Qwen CLI provider with custom settings
func NewQwenCLIWithSettings(log domain.Logger, settings domain.QwenCLISettings) (domain.AIProvider, error) {
	// Check if qwen-coder-cli is available
	_, err := exec.LookPath("qwen-coder")
	if err != nil {
		// Try alternative names
		_, err = exec.LookPath("qwen")
		if err != nil {
			return nil, fmt.Errorf("qwen-coder CLI not found in PATH. Please install it first")
		}
	}

	return &QwenCLIProviderImpl{
		log:      log,
		settings: settings,
	}, nil
}

// getQwenCommand returns the available qwen command
func (p *QwenCLIProviderImpl) getQwenCommand() string {
	if _, err := exec.LookPath("qwen-coder"); err == nil {
		return "qwen-coder"
	}
	return "qwen"
}

// ListModels returns available Qwen CLI models
func (p *QwenCLIProviderImpl) ListModels(ctx context.Context) ([]string, error) {
	// Return known models that work with qwen-coder CLI
	return []string{
		"qwen-coder-plus-latest",
		"qwen-coder-turbo-latest",
		"qwen-turbo-latest",
		"qwen-plus-latest",
	}, nil
}

// Generate sends a request to Qwen CLI and returns the response
func (p *QwenCLIProviderImpl) Generate(ctx context.Context, req domain.AIRequest) (domain.AIResponse, error) {
	startTime := time.Now()
	p.log.Info(fmt.Sprintf("Sending request to Qwen CLI with model: %s", req.Model))

	// Build the prompt
	prompt := req.UserPrompt
	if req.SystemPrompt != "" {
		prompt = fmt.Sprintf("System: %s\n\nUser: %s", req.SystemPrompt, req.UserPrompt)
	}

	// Build command arguments using settings
	args := p.buildArgs(req.Model)

	// Create command with context
	cmdName := p.getQwenCommand()
	cmd := exec.CommandContext(ctx, cmdName, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Always pass prompt via stdin to avoid Windows command line length limits
	cmd.Stdin = strings.NewReader(prompt)

	// Run the command
	err := cmd.Run()
	if err != nil {
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = err.Error()
		}
		p.log.Error(fmt.Sprintf("Qwen CLI request failed: %s", errMsg))
		return domain.AIResponse{}, fmt.Errorf("qwen CLI error: %s", errMsg)
	}

	content := strings.TrimSpace(stdout.String())
	processingTime := time.Since(startTime)

	// Estimate tokens (rough approximation)
	tokensUsed := (len(prompt) + len(content)) / 4

	return domain.AIResponse{
		Content:        content,
		TokensUsed:     tokensUsed,
		ModelUsed:      req.Model,
		ProcessingTime: processingTime,
		FinishReason:   "stop",
		Confidence:     0.9,
	}, nil
}

// buildArgs constructs CLI arguments from settings
func (p *QwenCLIProviderImpl) buildArgs(model string) []string {
	args := []string{
		"--model", model,
		"--input-format", "text",
	}

	// Output format from settings
	outputFormat := p.settings.OutputFormat
	if outputFormat == "" {
		outputFormat = "text"
	}
	args = append(args, "--output-format", outputFormat)

	// YOLO mode and approval mode are mutually exclusive
	// --yolo is shorthand for --approval-mode=yolo, cannot use both
	if p.settings.YoloMode {
		args = append(args, "--approval-mode", "yolo")
	} else if p.settings.ApprovalMode != "" {
		args = append(args, "--approval-mode", p.settings.ApprovalMode)
	}

	// Sandbox mode (safe execution)
	if p.settings.SandboxMode {
		args = append(args, "--sandbox")
	}

	// Max session turns limit
	if p.settings.MaxSessionTurns > 0 {
		args = append(args, "--max-session-turns", fmt.Sprintf("%d", p.settings.MaxSessionTurns))
	}

	// Debug mode
	if p.settings.DebugMode {
		args = append(args, "--debug")
	}

	return args
}

// GetProviderInfo returns information about the Qwen CLI provider
func (p *QwenCLIProviderImpl) GetProviderInfo() domain.ProviderInfo {
	return domain.ProviderInfo{
		Name:         "Qwen Code CLI",
		Version:      "1.0",
		Capabilities: []string{"chat", "completion", "code_generation"},
		Limitations:  []string{"local_only", "no_streaming"},
		SupportedModels: []string{
			"qwen-coder-plus-latest",
			"qwen-coder-turbo-latest",
			"qwen-turbo-latest",
		},
	}
}

// ValidateRequest validates the request parameters
func (p *QwenCLIProviderImpl) ValidateRequest(req domain.AIRequest) error {
	if req.Model == "" {
		return fmt.Errorf("model is required")
	}
	if req.UserPrompt == "" {
		return fmt.Errorf("user prompt is required")
	}
	return nil
}

// EstimateTokens estimates the number of tokens in the request
func (p *QwenCLIProviderImpl) EstimateTokens(req domain.AIRequest) (int, error) {
	totalChars := len(req.SystemPrompt) + len(req.UserPrompt)
	estimatedTokens := totalChars / 4
	return estimatedTokens + 100, nil
}

// GetPricing returns pricing information (free for CLI)
func (p *QwenCLIProviderImpl) GetPricing(model string) domain.PricingInfo {
	return domain.PricingInfo{
		Model:             model,
		Currency:          "USD",
		InputTokensPer1K:  0,
		OutputTokensPer1K: 0,
	}
}

// GetMaxContextTokens returns the maximum context size for a model
func (p *QwenCLIProviderImpl) GetMaxContextTokens(model string) int {
	switch model {
	case "qwen-coder-plus-latest":
		return 131072
	case "qwen-coder-turbo-latest":
		return 131072
	default:
		return 32768
	}
}

// GenerateStream implements streaming for CLI (fallback to non-streaming)
func (p *QwenCLIProviderImpl) GenerateStream(ctx context.Context, req domain.AIRequest, onChunk func(chunk domain.StreamChunk)) error {
	// CLI doesn't support true streaming, so we generate and send as one chunk
	resp, err := p.Generate(ctx, req)
	if err != nil {
		onChunk(domain.StreamChunk{Done: true, Error: err.Error()})
		return err
	}

	// Send content in smaller chunks to simulate streaming
	content := resp.Content
	chunkSize := 50 // characters per chunk
	for i := 0; i < len(content); i += chunkSize {
		end := i + chunkSize
		if end > len(content) {
			end = len(content)
		}
		onChunk(domain.StreamChunk{
			Content: content[i:end],
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
