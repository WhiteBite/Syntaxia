package main

import (
	"encoding/json"
	"fmt"
	"syntaxia/domain"
	"syntaxia/handlers"
)

// GenerateCode generates code using AI
func (a *App) GenerateCode(systemPrompt, userPrompt string) (string, error) {
	return a.aiHandler.GenerateCode(a.ctx, systemPrompt, userPrompt)
}

// GenerateCodeStream generates code with streaming response via Wails events
func (a *App) GenerateCodeStream(systemPrompt, userPrompt string) {
	go func() {
		if err := a.aiHandler.GenerateCodeStream(a.ctx, systemPrompt, userPrompt, func(chunk domain.StreamChunk) {
			a.bridge.Emit("ai:stream:chunk", chunk)
		}); err != nil {
			a.bridge.Emit("ai:stream:chunk", domain.StreamChunk{Done: true, Error: err.Error()})
		}
	}()
}

// GenerateIntelligentCode performs intelligent code generation
func (a *App) GenerateIntelligentCode(task, context, optionsJson string) (string, error) {
	return a.aiHandler.GenerateIntelligentCode(a.ctx, task, context, optionsJson)
}

// GenerateCodeWithOptions generates code with additional options
func (a *App) GenerateCodeWithOptions(systemPrompt, userPrompt, optionsJson string) (string, error) {
	return a.aiHandler.GenerateCodeWithOptions(a.ctx, systemPrompt, userPrompt, optionsJson)
}

// GetProviderInfo returns information about the current AI provider
func (a *App) GetProviderInfo() (string, error) {
	return a.aiHandler.GetProviderInfo(a.ctx)
}

// ListAvailableModels returns list of available AI models
func (a *App) ListAvailableModels() ([]string, error) {
	return a.aiHandler.ListAvailableModels(a.ctx)
}

// SuggestContextFiles suggests relevant files for a task
func (a *App) SuggestContextFiles(task string, allFiles []*domain.FileNode) ([]string, error) {
	return a.aiHandler.SuggestContextFiles(a.ctx, task, allFiles)
}

// AnalyzeTaskAndCollectContext analyzes a task and collects relevant context
func (a *App) AnalyzeTaskAndCollectContext(task string, allFilesJson string, rootDir string) (string, error) {
	return a.aiHandler.AnalyzeTaskAndCollectContext(a.ctx, task, allFilesJson, rootDir)
}

// AgenticChat performs agentic chat with tool use
func (a *App) AgenticChat(requestJson string) (string, error) {
	return a.aiHandler.AgenticChat(a.ctx, requestJson)
}

// ==================== Qwen Task Methods ====================

// QwenExecuteTask executes a task using Qwen with smart context collection
func (a *App) QwenExecuteTask(requestJson string) (string, error) {
	var req handlers.ExecuteTaskRequest
	if err := json.Unmarshal([]byte(requestJson), &req); err != nil {
		return "", fmt.Errorf("failed to parse request: %w", err)
	}

	result := a.qwenHandler.ExecuteTask(req)

	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(resultJson), nil
}

// QwenPreviewContext returns a preview of the context that would be collected
func (a *App) QwenPreviewContext(requestJson string) (string, error) {
	var req handlers.ExecuteTaskRequest
	if err := json.Unmarshal([]byte(requestJson), &req); err != nil {
		return "", fmt.Errorf("failed to parse request: %w", err)
	}

	result := a.qwenHandler.PreviewContext(req)

	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(resultJson), nil
}

// QwenGetAvailableModels returns available Qwen models
func (a *App) QwenGetAvailableModels() (string, error) {
	models := a.qwenHandler.GetAvailableModels()

	resultJson, err := json.Marshal(models)
	if err != nil {
		return "", fmt.Errorf("failed to marshal models: %w", err)
	}

	return string(resultJson), nil
}

// ==================== Semantic Search Methods ====================

// SemanticSearch performs semantic search on the project
func (a *App) SemanticSearch(requestJson string) (string, error) {
	if a.container.SemanticHandler == nil {
		return "", fmt.Errorf("semantic search not available: embedding provider not configured")
	}
	return a.container.SemanticHandler.Search(a.ctx, requestJson)
}

// SemanticFindSimilar finds similar code
func (a *App) SemanticFindSimilar(requestJson string) (string, error) {
	if a.container.SemanticHandler == nil {
		return "", fmt.Errorf("semantic search not available: embedding provider not configured")
	}
	return a.container.SemanticHandler.FindSimilar(a.ctx, requestJson)
}

// SemanticIndexProject indexes a project for semantic search
func (a *App) SemanticIndexProject(projectRoot string) error {
	if a.container.SemanticHandler == nil {
		return fmt.Errorf("semantic search not available: embedding provider not configured")
	}
	return a.container.SemanticHandler.IndexProject(a.ctx, projectRoot)
}

// SemanticIndexFile indexes a single file
func (a *App) SemanticIndexFile(projectRoot, filePath string) error {
	if a.container.SemanticHandler == nil {
		return fmt.Errorf("semantic search not available: embedding provider not configured")
	}
	return a.container.SemanticHandler.IndexFile(a.ctx, projectRoot, filePath)
}

// SemanticGetStats returns semantic search index statistics
func (a *App) SemanticGetStats(projectRoot string) (string, error) {
	if a.container.SemanticHandler == nil {
		return "", fmt.Errorf("semantic search not available: embedding provider not configured")
	}
	return a.container.SemanticHandler.GetStats(a.ctx, projectRoot)
}

// SemanticIsIndexed checks if a project is indexed
func (a *App) SemanticIsIndexed(projectRoot string) bool {
	if a.container.SemanticHandler == nil {
		return false
	}
	return a.container.SemanticHandler.IsIndexed(a.ctx, projectRoot)
}

// SemanticRetrieveContext retrieves relevant context using RAG
func (a *App) SemanticRetrieveContext(requestJson string) (string, error) {
	if a.container.SemanticHandler == nil {
		return "", fmt.Errorf("semantic search not available: embedding provider not configured")
	}
	return a.container.SemanticHandler.RetrieveContext(a.ctx, requestJson)
}

// SemanticHybridSearch performs hybrid keyword + semantic search
func (a *App) SemanticHybridSearch(requestJson string) (string, error) {
	if a.container.SemanticHandler == nil {
		return "", fmt.Errorf("semantic search not available: embedding provider not configured")
	}
	return a.container.SemanticHandler.HybridSearch(a.ctx, requestJson)
}

// IsSemanticSearchAvailable checks if semantic search is configured
func (a *App) IsSemanticSearchAvailable() bool {
	return a.container.SemanticHandler != nil
}
