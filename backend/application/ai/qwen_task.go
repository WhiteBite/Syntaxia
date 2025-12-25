package ai

import (
	"context"
	"fmt"
	"path/filepath"
	"syntaxia/domain"
	"strings"
	"time"
)

// QwenTaskService orchestrates task execution with Qwen
type QwenTaskService struct {
log                 domain.Logger
aiService           *Service
smartContextService SmartContextProvider
settingsService     SettingsProvider
}

// SmartContextProvider interface for smart context collection
type SmartContextProvider interface {
CollectContext(ctx context.Context, req SmartContextRequest) (*SmartContextResult, error)
}

// SmartContextRequest represents a request for smart context collection
type SmartContextRequest struct {
ProjectRoot, Task, SelectedCode, SourceFile, Language string
SelectedFiles                                         []string
MaxTokens, MaxDepth                                   int
}

// SmartContextResult represents the result of smart context collection
type SmartContextResult struct {
Context         string
Files           []ContextFile
Symbols         []*domain.SymbolNode
CallStack       *CallStackResult
TokenEstimate   int
TruncatedFiles  []string
ExcludedFiles   []string
RelevanceScores map[string]float64
}

// ContextFile represents a file in the context
type ContextFile struct {
Path, Content, Reason string
Tokens                int
Relevance             float64
}

// CallStackResult mirrors the infrastructure type
type CallStackResult struct {
RootSymbol   *domain.SymbolNode
Callers      []*domain.SymbolNode
Callees      []*domain.SymbolNode
Dependencies []*domain.SymbolNode
RelatedFiles []string
TotalSymbols int
}

// NewQwenTaskService creates a new Qwen task service
func NewQwenTaskService(log domain.Logger, aiService *Service, smartContextService SmartContextProvider, settingsService SettingsProvider) *QwenTaskService {
return &QwenTaskService{log: log, aiService: aiService, smartContextService: smartContextService, settingsService: settingsService}
}

// TaskRequest represents a task to be executed with Qwen
type TaskRequest struct {
Task, ProjectRoot, SelectedCode, SourceFile, Model string
SelectedFiles                                      []string
MaxTokens                                          int
Temperature                                        float64
}

// TaskResponse contains the result of task execution
type TaskResponse struct {
Content, Model, Error              string
TokensUsed                         int
ProcessingTime                     time.Duration
ContextSummary                     ContextSummaryDTO
Success                            bool
}

// ContextSummaryDTO provides info about the context sent to the model
type ContextSummaryDTO struct {
TotalFiles, TotalTokens                          int
IncludedFiles, TruncatedFiles, ExcludedFiles     []string
}

// ExecuteTask executes a task using Qwen with smart context collection
func (s *QwenTaskService) ExecuteTask(ctx context.Context, req TaskRequest) (*TaskResponse, error) {
	startTime := time.Now()
	s.log.Info(fmt.Sprintf("Executing task with Qwen: %s", domain.TruncateString(req.Task, 100)))

if req.Model == "" {
req.Model = "qwen-coder-plus-latest"
}
if req.MaxTokens == 0 {
req.MaxTokens = 900000
}
if req.Temperature == 0 {
req.Temperature = 0.3
}

contextReq := SmartContextRequest{ProjectRoot: req.ProjectRoot, Task: req.Task, SelectedFiles: req.SelectedFiles, SelectedCode: req.SelectedCode, SourceFile: req.SourceFile, MaxTokens: req.MaxTokens, MaxDepth: 3}
smartContext, err := s.smartContextService.CollectContext(ctx, contextReq)
if err != nil {
return &TaskResponse{Success: false, Error: fmt.Sprintf("failed to collect context: %v", err)}, err
}

systemPrompt := s.buildSystemPrompt()
userPrompt := s.buildUserPrompt(req, smartContext)

response, err := s.aiService.GenerateCodeWithOptions(ctx, systemPrompt, userPrompt, GenerationOptions{Model: req.Model, Temperature: req.Temperature, MaxTokens: 32000, Timeout: 5 * time.Minute})
summary := ContextSummaryDTO{TotalFiles: len(smartContext.Files), TotalTokens: smartContext.TokenEstimate, IncludedFiles: s.getFilePaths(smartContext.Files), TruncatedFiles: smartContext.TruncatedFiles, ExcludedFiles: smartContext.ExcludedFiles}
if err != nil {
return &TaskResponse{Success: false, Error: fmt.Sprintf("AI generation failed: %v", err), ContextSummary: summary}, err
}

return &TaskResponse{Content: response, Model: req.Model, TokensUsed: smartContext.TokenEstimate, ProcessingTime: time.Since(startTime), Success: true, ContextSummary: summary}, nil
}

// PreviewContext returns a preview of what context would be collected
func (s *QwenTaskService) PreviewContext(ctx context.Context, req TaskRequest) (*SmartContextResult, error) {
if req.MaxTokens == 0 {
req.MaxTokens = 900000
}
return s.smartContextService.CollectContext(ctx, SmartContextRequest{ProjectRoot: req.ProjectRoot, Task: req.Task, SelectedFiles: req.SelectedFiles, SelectedCode: req.SelectedCode, SourceFile: req.SourceFile, MaxTokens: req.MaxTokens, MaxDepth: 3})
}

func (s *QwenTaskService) buildSystemPrompt() string {
return `You are an expert software developer assistant. Your task is to help implement features and fix issues in code.

IMPORTANT RULES:
1. Analyze the provided codebase context carefully before making changes
2. Follow the existing code style and patterns in the project
3. Provide complete, working code - not snippets or pseudocode
4. When modifying existing files, provide the changes in unified diff format
5. When creating new files, provide the complete file content
6. Explain your changes briefly but focus on the code
7. Consider edge cases and error handling
8. Maintain backward compatibility unless explicitly asked to break it

OUTPUT FORMAT:
- For file modifications, use this format:
  --- a/path/to/file.go
  +++ b/path/to/file.go
  @@ -line,count +line,count @@
  -old line
  +new line

- For new files, use this format:
  === NEW FILE: path/to/new_file.go ===
  <complete file content>
  === END FILE ===

- After code changes, provide a brief summary of what was changed and why.`
}

func (s *QwenTaskService) buildUserPrompt(req TaskRequest, smartContext *SmartContextResult) string {
var builder strings.Builder
builder.WriteString("# Task\n")
builder.WriteString(req.Task)
builder.WriteString("\n\n")

if req.SelectedCode != "" {
builder.WriteString("# Selected Code (focus area)\nFile: ")
builder.WriteString(req.SourceFile)
builder.WriteString("\n```\n")
builder.WriteString(req.SelectedCode)
builder.WriteString("\n```\n\n")
}

builder.WriteString(fmt.Sprintf("# Project Context\nTotal files in context: %d\nEstimated tokens: %d\n\n", len(smartContext.Files), smartContext.TokenEstimate))

if smartContext.CallStack != nil && smartContext.CallStack.RootSymbol != nil {
builder.WriteString(fmt.Sprintf("## Call Stack Analysis\nRoot symbol: %s (%s)\n", smartContext.CallStack.RootSymbol.Name, smartContext.CallStack.RootSymbol.Type))
if len(smartContext.CallStack.Callers) > 0 {
names := make([]string, 0, len(smartContext.CallStack.Callers))
for _, c := range smartContext.CallStack.Callers {
names = append(names, c.Name)
}
builder.WriteString("Called by: " + strings.Join(names, ", ") + "\n")
}
if len(smartContext.CallStack.Callees) > 0 {
names := make([]string, 0, len(smartContext.CallStack.Callees))
for _, c := range smartContext.CallStack.Callees {
names = append(names, c.Name)
}
builder.WriteString("Calls: " + strings.Join(names, ", ") + "\n")
}
builder.WriteString("\n")
}

builder.WriteString("## Files\n\n")
for _, file := range smartContext.Files {
builder.WriteString(fmt.Sprintf("### %s\nRelevance: %.0f%% - %s\n```%s\n%s\n```\n\n", file.Path, file.Relevance*100, file.Reason, getLanguageFromPath(file.Path), file.Content))
}

builder.WriteString("# Instructions\nBased on the task and the provided context, implement the required changes. Make sure to follow the existing code patterns and style.\n")
return builder.String()
}

func (s *QwenTaskService) getFilePaths(files []ContextFile) []string {
paths := make([]string, len(files))
for i, f := range files {
paths[i] = f.Path
}
return paths
}

var extToLang = map[string]string{".go": "go", ".ts": "typescript", ".tsx": "typescript", ".js": "javascript", ".jsx": "javascript", ".py": "python", ".java": "java", ".rs": "rust", ".vue": "vue", ".css": "css", ".scss": "css", ".html": "html", ".json": "json", ".yaml": "yaml", ".yml": "yaml"}

func getLanguageFromPath(path string) string {
if lang, ok := extToLang[filepath.Ext(path)]; ok {
return lang
}
return ""
}
