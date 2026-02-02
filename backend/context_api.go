package main

import (
	"encoding/json"
	"fmt"
	"syntaxia/domain"
)

// RequestsyntaxiaContextGeneration generates context for selected files
func (a *App) RequestsyntaxiaContextGeneration(rootDir string, includedPaths []string) {
	a.projectHandler.GenerateContext(a.ctx, rootDir, includedPaths)
}

// BuildContext builds context and returns ContextSummary (OOM-safe)
func (a *App) BuildContext(projectPath string, includedPaths []string, optionsJson string) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	var options domain.ContextBuildOptions
	if optionsJson != "" {
		if err := json.Unmarshal([]byte(optionsJson), &options); err != nil {
			validationErr := domain.NewValidationError("failed to parse options JSON", map[string]interface{}{
				"originalError": err.Error(),
				"optionsJson":   optionsJson,
			})
			return "", a.transformError(validationErr)
		}
	}

	if len(includedPaths) == 0 {
		a.log.Warning("BuildContext called with empty includedPaths - this may include all project files")
	}

	summary, err := a.contextService.BuildContextSummary(a.ctx, projectPath, includedPaths, &options)
	if err != nil {
		return "", a.transformError(err)
	}

	contextJson, err := json.Marshal(summary)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context summary", err)
		return "", a.transformError(marshalErr)
	}

	return string(contextJson), nil
}

// BuildContextFromRequest builds context using provided options and returns ContextSummary
func (a *App) BuildContextFromRequest(projectPath string, includedPaths []string, options *domain.ContextBuildOptions) (*domain.ContextSummary, error) {
	if a.contextService == nil {
		return nil, a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	if len(includedPaths) == 0 {
		a.log.Warning("BuildContextFromRequest called with empty includedPaths")
	}

	if options == nil {
		options = &domain.ContextBuildOptions{}
	}

	a.log.Info(fmt.Sprintf("[BuildContextFromRequest] OutputFormat received: '%s', StripComments: %v, MaxTokens: %d",
		options.OutputFormat, options.StripComments, options.MaxTokens))

	summary, err := a.contextService.BuildContextSummary(a.ctx, projectPath, includedPaths, options)
	if err != nil {
		return nil, a.transformError(err)
	}

	return summary, nil
}

// BuildContextLegacy is deprecated - use BuildContext instead
func (a *App) BuildContextLegacy() (string, error) {
	return "", a.transformError(domain.NewConfigurationError("legacy context building is no longer supported", nil))
}

// CreateStreamingContext delegates to BuildContext to create a disk-backed context summary
func (a *App) CreateStreamingContext(projectPath string, includedPaths []string, optionsJson string) (string, error) {
	return a.BuildContext(projectPath, includedPaths, optionsJson)
}

// CollectSmartContext collects relevant context for AI task
func (a *App) CollectSmartContext(requestJson string) (string, error) {
	var req domain.SmartContextRequest
	if err := json.Unmarshal([]byte(requestJson), &req); err != nil {
		return "", fmt.Errorf("invalid request: %w", err)
	}

	if a.container.SmartContextCollector == nil {
		return "", fmt.Errorf("smart context collector not available")
	}

	result, err := a.container.SmartContextCollector.CollectContext(a.ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to collect context: %w", err)
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(resultJson), nil
}
