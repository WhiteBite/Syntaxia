package repair

import (
	"context"
	"fmt"
	"math"
	"strings"
	"syntaxia/domain"
	"time"
)

// SelfCorrectionConfig holds configuration for self-correction mechanism
type SelfCorrectionConfig struct {
	MaxRetries         int           `json:"maxRetries"`
	InitialBackoff     time.Duration `json:"initialBackoff"`
	MaxBackoff         time.Duration `json:"maxBackoff"`
	BackoffMultiplier  float64       `json:"backoffMultiplier"`
	EnableAIFallback   bool          `json:"enableAIFallback"`
	CollectFullContext bool          `json:"collectFullContext"`
}

// DefaultSelfCorrectionConfig returns default configuration
func DefaultSelfCorrectionConfig() *SelfCorrectionConfig {
	return &SelfCorrectionConfig{
		MaxRetries:         5,
		InitialBackoff:     100 * time.Millisecond,
		MaxBackoff:         30 * time.Second,
		BackoffMultiplier:  2.0,
		EnableAIFallback:   true,
		CollectFullContext: true,
	}
}

// ErrorContext contains collected error context for AI assistance
type ErrorContext struct {
	OriginalError    *domain.ErrorDetails   `json:"originalError"`
	AttemptedFixes   []AttemptedFix         `json:"attemptedFixes"`
	SourceCode       string                 `json:"sourceCode,omitempty"`
	SurroundingCode  string                 `json:"surroundingCode,omitempty"`
	RelatedFiles     []string               `json:"relatedFiles,omitempty"`
	ProjectLanguages []string               `json:"projectLanguages,omitempty"`
	BuildOutput      string                 `json:"buildOutput,omitempty"`
	LintOutput       string                 `json:"lintOutput,omitempty"`
	TotalAttempts    int                    `json:"totalAttempts"`
	ElapsedTime      time.Duration          `json:"elapsedTime"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// AttemptedFix represents a single correction attempt
type AttemptedFix struct {
	Attempt     int                    `json:"attempt"`
	Action      domain.CorrectionAction `json:"action"`
	Description string                 `json:"description"`
	Success     bool                   `json:"success"`
	Result      string                 `json:"result"`
	Duration    time.Duration          `json:"duration"`
	Error       string                 `json:"error,omitempty"`
}

// SelfCorrectionResult represents the result of self-correction process
type SelfCorrectionResult struct {
	Success           bool                     `json:"success"`
	FinalError        *domain.ErrorDetails     `json:"finalError,omitempty"`
	TotalAttempts     int                      `json:"totalAttempts"`
	AttemptedFixes    []AttemptedFix           `json:"attemptedFixes"`
	RequiresAI        bool                     `json:"requiresAI"`
	AIContext         *ErrorContext            `json:"aiContext,omitempty"`
	FilesChanged      []string                 `json:"filesChanged"`
	TotalDuration     time.Duration            `json:"totalDuration"`
	CorrectionApplied *domain.CorrectionResult `json:"correctionApplied,omitempty"`
}

// AIAssistanceRequest represents a request for AI assistance
type AIAssistanceRequest struct {
	ErrorContext    *ErrorContext `json:"errorContext"`
	ProjectPath     string        `json:"projectPath"`
	RequestedAction string        `json:"requestedAction"`
	Priority        string        `json:"priority"`
}

// SelfCorrectionEngine implements the self-correction fallback mechanism
type SelfCorrectionEngine struct {
	log              domain.Logger
	errorAnalyzer    domain.ErrorAnalyzer
	correctionEngine domain.CorrectionEngine
	fileSystem       domain.FileSystemProvider
	config           *SelfCorrectionConfig
}

// NewSelfCorrectionEngine creates a new SelfCorrectionEngine instance
func NewSelfCorrectionEngine(
	log domain.Logger,
	errorAnalyzer domain.ErrorAnalyzer,
	correctionEngine domain.CorrectionEngine,
	fileSystem domain.FileSystemProvider,
	config *SelfCorrectionConfig,
) *SelfCorrectionEngine {
	if config == nil {
		config = DefaultSelfCorrectionConfig()
	}
	return &SelfCorrectionEngine{
		log:              log,
		errorAnalyzer:    errorAnalyzer,
		correctionEngine: correctionEngine,
		fileSystem:       fileSystem,
		config:           config,
	}
}

// AttemptCorrection attempts to correct an error with retry logic and exponential backoff
func (s *SelfCorrectionEngine) AttemptCorrection(
	ctx context.Context,
	errorOutput string,
	stage domain.ProtocolStage,
	projectPath string,
) (*SelfCorrectionResult, error) {
	startTime := time.Now()

	result := &SelfCorrectionResult{
		AttemptedFixes: make([]AttemptedFix, 0),
		FilesChanged:   make([]string, 0),
	}

	// Analyze the initial error
	errorDetails, err := s.errorAnalyzer.AnalyzeError(errorOutput, stage)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze error: %w", err)
	}

	s.log.Info(fmt.Sprintf("Starting self-correction for error type: %s in file: %s",
		errorDetails.ErrorType, errorDetails.SourceFile))

	// Attempt corrections with exponential backoff
	for attempt := 1; attempt <= s.config.MaxRetries; attempt++ {
		result.TotalAttempts = attempt

		s.log.Debug(fmt.Sprintf("Self-correction attempt %d/%d", attempt, s.config.MaxRetries))

		attemptResult := s.executeAttempt(ctx, attempt, errorDetails, projectPath)
		result.AttemptedFixes = append(result.AttemptedFixes, attemptResult)

		if attemptResult.Success {
			result.Success = true
			result.FilesChanged = append(result.FilesChanged, errorDetails.SourceFile)
			result.TotalDuration = time.Since(startTime)
			s.log.Info(fmt.Sprintf("Self-correction succeeded on attempt %d", attempt))
			return result, nil
		}

		// Check if we should continue retrying
		if attempt < s.config.MaxRetries {
			backoff := s.calculateBackoff(attempt)
			s.log.Debug(fmt.Sprintf("Waiting %v before next attempt", backoff))

			select {
			case <-ctx.Done():
				result.TotalDuration = time.Since(startTime)
				return result, ctx.Err()
			case <-time.After(backoff):
				// Continue to next attempt
			}
		}
	}

	// All simple fixes failed - prepare for AI fallback
	result.Success = false
	result.FinalError = errorDetails
	result.TotalDuration = time.Since(startTime)

	if s.config.EnableAIFallback {
		result.RequiresAI = true
		result.AIContext = s.collectErrorContext(ctx, errorDetails, result.AttemptedFixes, projectPath, startTime)
		s.log.Info("Simple fixes exhausted, AI assistance required")
	}

	return result, nil
}

// executeAttempt executes a single correction attempt
func (s *SelfCorrectionEngine) executeAttempt(
	ctx context.Context,
	attempt int,
	errorDetails *domain.ErrorDetails,
	projectPath string,
) AttemptedFix {
	startTime := time.Now()

	fix := AttemptedFix{
		Attempt: attempt,
	}

	// Get correction suggestions
	corrections, err := s.errorAnalyzer.SuggestCorrections(errorDetails)
	if err != nil {
		fix.Success = false
		fix.Error = fmt.Sprintf("failed to get corrections: %v", err)
		fix.Duration = time.Since(startTime)
		return fix
	}

	if len(corrections) == 0 {
		fix.Success = false
		fix.Error = "no corrections suggested"
		fix.Duration = time.Since(startTime)
		return fix
	}

	// Try to apply the first applicable correction
	for _, correction := range corrections {
		fix.Action = correction.Action
		fix.Description = correction.Description

		// Check if correction engine can handle this error
		if !s.correctionEngine.CanHandle(errorDetails) {
			continue
		}

		result, applyErr := s.correctionEngine.ApplyCorrection(ctx, correction, projectPath)
		if applyErr != nil {
			fix.Error = fmt.Sprintf("failed to apply correction: %v", applyErr)
			continue
		}

		fix.Success = result.Success
		fix.Result = result.Message
		fix.Duration = time.Since(startTime)

		if result.Success {
			return fix
		}
	}

	fix.Duration = time.Since(startTime)
	return fix
}

// calculateBackoff calculates the backoff duration for a given attempt
func (s *SelfCorrectionEngine) calculateBackoff(attempt int) time.Duration {
	backoff := float64(s.config.InitialBackoff) * math.Pow(s.config.BackoffMultiplier, float64(attempt-1))
	if backoff > float64(s.config.MaxBackoff) {
		backoff = float64(s.config.MaxBackoff)
	}
	return time.Duration(backoff)
}

// collectErrorContext collects comprehensive error context for AI assistance
func (s *SelfCorrectionEngine) collectErrorContext(
	ctx context.Context,
	errorDetails *domain.ErrorDetails,
	attemptedFixes []AttemptedFix,
	projectPath string,
	startTime time.Time,
) *ErrorContext {
	errCtx := &ErrorContext{
		OriginalError:  errorDetails,
		AttemptedFixes: attemptedFixes,
		TotalAttempts:  len(attemptedFixes),
		ElapsedTime:    time.Since(startTime),
		Metadata:       make(map[string]interface{}),
	}

	if !s.config.CollectFullContext {
		return errCtx
	}

	// Collect source code if available
	if errorDetails.SourceFile != "" {
		errCtx.SourceCode = s.readSourceFile(ctx, projectPath, errorDetails.SourceFile)
		errCtx.SurroundingCode = s.extractSurroundingCode(
			errCtx.SourceCode,
			errorDetails.LineNumber,
			contextLinesBefore,
			contextLinesAfter,
		)
	}

	// Collect related files based on error type
	errCtx.RelatedFiles = s.findRelatedFiles(errorDetails, projectPath)

	// Add metadata
	errCtx.Metadata["projectPath"] = projectPath
	errCtx.Metadata["errorType"] = string(errorDetails.ErrorType)
	errCtx.Metadata["stage"] = string(errorDetails.Stage)
	errCtx.Metadata["tool"] = errorDetails.Tool

	return errCtx
}

// Context extraction constants
const (
	contextLinesBefore = 10
	contextLinesAfter  = 10
)

// readSourceFile reads the source file content
func (s *SelfCorrectionEngine) readSourceFile(ctx context.Context, projectPath, sourceFile string) string {
	if s.fileSystem == nil {
		return ""
	}

	filePath := sourceFile
	if projectPath != "" && !strings.HasPrefix(sourceFile, projectPath) {
		filePath = projectPath + "/" + sourceFile
	}

	content, err := s.fileSystem.ReadFile(filePath)
	if err != nil {
		s.log.Debug(fmt.Sprintf("Failed to read source file %s: %v", filePath, err))
		return ""
	}

	return string(content)
}

// extractSurroundingCode extracts code around the error line
func (s *SelfCorrectionEngine) extractSurroundingCode(source string, lineNumber, before, after int) string {
	if source == "" || lineNumber <= 0 {
		return ""
	}

	lines := strings.Split(source, "\n")
	if lineNumber > len(lines) {
		return ""
	}

	startLine := lineNumber - before - 1
	if startLine < 0 {
		startLine = 0
	}

	endLine := lineNumber + after
	if endLine > len(lines) {
		endLine = len(lines)
	}

	var result strings.Builder
	for i := startLine; i < endLine; i++ {
		lineNum := i + 1
		marker := "  "
		if lineNum == lineNumber {
			marker = "> "
		}
		result.WriteString(fmt.Sprintf("%s%4d | %s\n", marker, lineNum, lines[i]))
	}

	return result.String()
}

// findRelatedFiles finds files related to the error
func (s *SelfCorrectionEngine) findRelatedFiles(errorDetails *domain.ErrorDetails, projectPath string) []string {
	relatedFiles := make([]string, 0)

	// Add the source file itself
	if errorDetails.SourceFile != "" {
		relatedFiles = append(relatedFiles, errorDetails.SourceFile)
	}

	// Based on error type, suggest related files
	switch errorDetails.ErrorType {
	case domain.ErrorTypeImport:
		// For import errors, might need to check go.mod, package.json, etc.
		relatedFiles = append(relatedFiles, "go.mod", "go.sum", "package.json")
	case domain.ErrorTypeTypeCheck:
		// For type errors, might need to check type definition files
		relatedFiles = append(relatedFiles, "types.go", "types.ts", "index.d.ts")
	}

	return relatedFiles
}

// FormatAIPrompt formats the error context into a prompt for AI assistance
func (s *SelfCorrectionEngine) FormatAIPrompt(errCtx *ErrorContext) string {
	var prompt strings.Builder

	prompt.WriteString("## Error Correction Request\n\n")

	// Error details
	prompt.WriteString("### Error Details\n")
	prompt.WriteString(fmt.Sprintf("- **Type**: %s\n", errCtx.OriginalError.ErrorType))
	prompt.WriteString(fmt.Sprintf("- **Stage**: %s\n", errCtx.OriginalError.Stage))
	prompt.WriteString(fmt.Sprintf("- **File**: %s\n", errCtx.OriginalError.SourceFile))
	if errCtx.OriginalError.LineNumber > 0 {
		prompt.WriteString(fmt.Sprintf("- **Line**: %d\n", errCtx.OriginalError.LineNumber))
	}
	prompt.WriteString(fmt.Sprintf("- **Message**: %s\n\n", errCtx.OriginalError.Message))

	// Attempted fixes
	if len(errCtx.AttemptedFixes) > 0 {
		prompt.WriteString("### Previously Attempted Fixes\n")
		for _, fix := range errCtx.AttemptedFixes {
			status := "❌ Failed"
			if fix.Success {
				status = "✅ Success"
			}
			prompt.WriteString(fmt.Sprintf("- Attempt %d: %s - %s\n", fix.Attempt, fix.Description, status))
			if fix.Error != "" {
				prompt.WriteString(fmt.Sprintf("  Error: %s\n", fix.Error))
			}
		}
		prompt.WriteString("\n")
	}

	// Surrounding code
	if errCtx.SurroundingCode != "" {
		prompt.WriteString("### Code Context\n")
		prompt.WriteString("```\n")
		prompt.WriteString(errCtx.SurroundingCode)
		prompt.WriteString("```\n\n")
	}

	// Suggestions from error analysis
	if len(errCtx.OriginalError.Suggestions) > 0 {
		prompt.WriteString("### Analyzer Suggestions\n")
		for _, suggestion := range errCtx.OriginalError.Suggestions {
			prompt.WriteString(fmt.Sprintf("- %s\n", suggestion))
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("### Request\n")
	prompt.WriteString("Please analyze this error and provide a fix. ")
	prompt.WriteString("The automated correction attempts have failed. ")
	prompt.WriteString("Provide the corrected code or specific instructions to fix this issue.\n")

	return prompt.String()
}

// CreateAIAssistanceRequest creates a request for AI assistance
func (s *SelfCorrectionEngine) CreateAIAssistanceRequest(
	result *SelfCorrectionResult,
	projectPath string,
) *AIAssistanceRequest {
	if !result.RequiresAI || result.AIContext == nil {
		return nil
	}

	priority := "normal"
	if result.AIContext.OriginalError.Severity == "error" {
		priority = "high"
	}

	return &AIAssistanceRequest{
		ErrorContext:    result.AIContext,
		ProjectPath:     projectPath,
		RequestedAction: "fix_error",
		Priority:        priority,
	}
}

// ShouldRetry determines if another retry should be attempted
func (s *SelfCorrectionEngine) ShouldRetry(attempt int, err error) bool {
	if attempt >= s.config.MaxRetries {
		return false
	}

	// Don't retry on context cancellation
	if err == context.Canceled || err == context.DeadlineExceeded {
		return false
	}

	return true
}

// GetConfig returns the current configuration
func (s *SelfCorrectionEngine) GetConfig() *SelfCorrectionConfig {
	return s.config
}

// UpdateConfig updates the configuration
func (s *SelfCorrectionEngine) UpdateConfig(config *SelfCorrectionConfig) {
	if config != nil {
		s.config = config
	}
}
