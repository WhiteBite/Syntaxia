package handlers

import (
	"context"
	"syntaxia/application/ai"
	"syntaxia/domain"
	"time"
)

// QwenHandler handles Qwen-related API requests
type QwenHandler struct {
	log             domain.Logger
	qwenTaskService *ai.QwenTaskService
}

// NewQwenHandler creates a new Qwen handler
func NewQwenHandler(
	log domain.Logger,
	qwenTaskService *ai.QwenTaskService,
) *QwenHandler {
	return &QwenHandler{
		log:             log,
		qwenTaskService: qwenTaskService,
	}
}

// ExecuteTaskRequest is the request for executing a task
type ExecuteTaskRequest struct {
	Task          string   `json:"task"`
	ProjectRoot   string   `json:"projectRoot"`
	SelectedFiles []string `json:"selectedFiles"`
	SelectedCode  string   `json:"selectedCode"`
	SourceFile    string   `json:"sourceFile"`
	Model         string   `json:"model"`
	MaxTokens     int      `json:"maxTokens"`
	Temperature   float64  `json:"temperature"`
}

// ExecuteTaskResponse is the response from task execution
type ExecuteTaskResponse struct {
	Content        string               `json:"content"`
	Model          string               `json:"model"`
	TokensUsed     int                  `json:"tokensUsed"`
	ProcessingTime string               `json:"processingTime"`
	ContextSummary ai.ContextSummaryDTO `json:"contextSummary"`
	Success        bool                 `json:"success"`
	Error          string               `json:"error,omitempty"`
}

// PreviewContextResponse is the response for context preview
type PreviewContextResponse struct {
	TotalFiles      int                `json:"totalFiles"`
	TotalTokens     int                `json:"totalTokens"`
	Files           []FilePreview      `json:"files"`
	TruncatedFiles  []string           `json:"truncatedFiles"`
	ExcludedFiles   []string           `json:"excludedFiles"`
	CallStackInfo   *CallStackInfo     `json:"callStackInfo,omitempty"`
	RelevanceScores map[string]float64 `json:"relevanceScores"`
}

// FilePreview represents a file in the preview
type FilePreview struct {
	Path      string  `json:"path"`
	Tokens    int     `json:"tokens"`
	Relevance float64 `json:"relevance"`
	Reason    string  `json:"reason"`
}

// CallStackInfo contains call stack analysis info
type CallStackInfo struct {
	RootSymbol   string   `json:"rootSymbol"`
	Callers      []string `json:"callers"`
	Callees      []string `json:"callees"`
	Dependencies []string `json:"dependencies"`
}

// ExecuteTask executes a task with Qwen
func (h *QwenHandler) ExecuteTask(req ExecuteTaskRequest) ExecuteTaskResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	taskReq := ai.TaskRequest{
		Task:          req.Task,
		ProjectRoot:   req.ProjectRoot,
		SelectedFiles: req.SelectedFiles,
		SelectedCode:  req.SelectedCode,
		SourceFile:    req.SourceFile,
		Model:         req.Model,
		MaxTokens:     req.MaxTokens,
		Temperature:   req.Temperature,
	}

	result, err := h.qwenTaskService.ExecuteTask(ctx, taskReq)
	if err != nil {
		return ExecuteTaskResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ExecuteTaskResponse{
		Content:        result.Content,
		Model:          result.Model,
		TokensUsed:     result.TokensUsed,
		ProcessingTime: result.ProcessingTime.String(),
		ContextSummary: result.ContextSummary,
		Success:        result.Success,
		Error:          result.Error,
	}
}

// PreviewContext returns a preview of the context that would be collected
func (h *QwenHandler) PreviewContext(req ExecuteTaskRequest) PreviewContextResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	taskReq := ai.TaskRequest{
		Task:          req.Task,
		ProjectRoot:   req.ProjectRoot,
		SelectedFiles: req.SelectedFiles,
		SelectedCode:  req.SelectedCode,
		SourceFile:    req.SourceFile,
		MaxTokens:     req.MaxTokens,
	}

	result, err := h.qwenTaskService.PreviewContext(ctx, taskReq)
	if err != nil {
		h.log.Error("Failed to preview context: " + err.Error())
		return PreviewContextResponse{}
	}

	// Convert files to preview format
	files := make([]FilePreview, len(result.Files))
	for i, f := range result.Files {
		files[i] = FilePreview{
			Path:      f.Path,
			Tokens:    f.Tokens,
			Relevance: f.Relevance,
			Reason:    f.Reason,
		}
	}

	response := PreviewContextResponse{
		TotalFiles:      len(result.Files),
		TotalTokens:     result.TokenEstimate,
		Files:           files,
		TruncatedFiles:  result.TruncatedFiles,
		ExcludedFiles:   result.ExcludedFiles,
		RelevanceScores: result.RelevanceScores,
	}

	// Add call stack info if available
	if result.CallStack != nil && result.CallStack.RootSymbol != nil {
		callStackInfo := &CallStackInfo{
			RootSymbol:   result.CallStack.RootSymbol.Name,
			Callers:      make([]string, 0),
			Callees:      make([]string, 0),
			Dependencies: make([]string, 0),
		}

		for _, c := range result.CallStack.Callers {
			callStackInfo.Callers = append(callStackInfo.Callers, c.Name)
		}
		for _, c := range result.CallStack.Callees {
			callStackInfo.Callees = append(callStackInfo.Callees, c.Name)
		}
		for _, d := range result.CallStack.Dependencies {
			callStackInfo.Dependencies = append(callStackInfo.Dependencies, d.Name)
		}

		response.CallStackInfo = callStackInfo
	}

	return response
}

// GetAvailableModels returns available Qwen models
func (h *QwenHandler) GetAvailableModels() []ModelInfo {
	return []ModelInfo{
		{
			ID:          "qwen-coder-plus-latest",
			Name:        "Qwen Coder Plus (Latest)",
			Description: "Best for large codebases, 1M token context",
			MaxContext:  1000000,
			Recommended: true,
		},
		{
			ID:          "qwen-coder-plus",
			Name:        "Qwen Coder Plus",
			Description: "Stable version, 1M token context",
			MaxContext:  1000000,
			Recommended: false,
		},
		{
			ID:          "qwen-coder-turbo-latest",
			Name:        "Qwen Coder Turbo (Latest)",
			Description: "Faster, 128K token context",
			MaxContext:  131072,
			Recommended: false,
		},
		{
			ID:          "qwen-plus-latest",
			Name:        "Qwen Plus (Latest)",
			Description: "General purpose with good coding, 128K context",
			MaxContext:  131072,
			Recommended: false,
		},
		{
			ID:          "qwen-turbo-latest",
			Name:        "Qwen Turbo (Latest)",
			Description: "Fast general purpose, 128K context",
			MaxContext:  131072,
			Recommended: false,
		},
		{
			ID:          "qwen-max",
			Name:        "Qwen Max",
			Description: "Most capable, 32K context",
			MaxContext:  32768,
			Recommended: false,
		},
	}
}

// ModelInfo contains information about a model
type ModelInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxContext  int    `json:"maxContext"`
	Recommended bool   `json:"recommended"`
}
