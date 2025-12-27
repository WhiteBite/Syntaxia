package repair

import (
	"context"
	"fmt"
	"strings"
	"syntaxia/domain"
)

// CorrectionAttempt represents a single correction attempt with its outcome
type CorrectionAttempt struct {
	Error           *domain.ErrorDetails `json:"error"`
	AutoFixApplied  bool                 `json:"autoFixApplied"`
	AutoFixSuccess  bool                 `json:"autoFixSuccess"`
	NeedsAIHelp     bool                 `json:"needsAIHelp"`
	SuggestedPrompt string               `json:"suggestedPrompt,omitempty"`
	AttemptNumber   int                  `json:"attemptNumber"`
}

// CorrectionResult represents the result of attempting to correct multiple errors
type CorrectionResult struct {
	Attempts        []CorrectionAttempt    `json:"attempts"`
	AllFixed        bool                   `json:"allFixed"`
	RemainingErrors []*domain.ErrorDetails `json:"remainingErrors"`
	AIPrompts       []string               `json:"aiPrompts"`
	TotalAttempts   int                    `json:"totalAttempts"`
}

// SelfCorrectionService provides high-level self-correction with AI feedback
type SelfCorrectionService struct {
	correctionEngine domain.CorrectionEngine
	errorAnalyzer    domain.ErrorAnalyzer
	feedbackGen      *AIFeedbackGenerator
	fileSystem       domain.FileSystemProvider
	maxRetries       int
	log              domain.Logger
}

// NewSelfCorrectionService creates a new SelfCorrectionService
func NewSelfCorrectionService(
	log domain.Logger,
	correctionEngine domain.CorrectionEngine,
	errorAnalyzer domain.ErrorAnalyzer,
	fileSystem domain.FileSystemProvider,
	maxRetries int,
) *SelfCorrectionService {
	if maxRetries <= 0 {
		maxRetries = 3
	}
	return &SelfCorrectionService{
		correctionEngine: correctionEngine,
		errorAnalyzer:    errorAnalyzer,
		feedbackGen:      NewAIFeedbackGenerator(),
		fileSystem:       fileSystem,
		maxRetries:       maxRetries,
		log:              log,
	}
}

// AttemptCorrection attempts to correct errors, returning results with AI prompts if needed
func (s *SelfCorrectionService) AttemptCorrection(ctx context.Context, errors []*domain.ErrorDetails, projectPath string) (*CorrectionResult, error) {
	result := &CorrectionResult{
		Attempts:        make([]CorrectionAttempt, 0, len(errors)),
		RemainingErrors: make([]*domain.ErrorDetails, 0),
		AIPrompts:       make([]string, 0),
	}

	for _, err := range errors {
		attempt := s.attemptSingleCorrection(ctx, err, projectPath)
		result.Attempts = append(result.Attempts, attempt)
		result.TotalAttempts++

		if !attempt.AutoFixSuccess {
			result.RemainingErrors = append(result.RemainingErrors, err)
			if attempt.SuggestedPrompt != "" {
				result.AIPrompts = append(result.AIPrompts, attempt.SuggestedPrompt)
			}
		}
	}

	result.AllFixed = len(result.RemainingErrors) == 0
	return result, nil
}

// attemptSingleCorrection attempts to correct a single error
func (s *SelfCorrectionService) attemptSingleCorrection(ctx context.Context, err *domain.ErrorDetails, projectPath string) CorrectionAttempt {
	attempt := CorrectionAttempt{
		Error:         err,
		AttemptNumber: 1,
	}

	// Check if correction engine can handle this error
	if s.correctionEngine == nil || !s.correctionEngine.CanHandle(err) {
		attempt.NeedsAIHelp = true
		attempt.SuggestedPrompt = s.GenerateAIPrompt(err)
		return attempt
	}

	// Get correction suggestions
	corrections, suggestErr := s.errorAnalyzer.SuggestCorrections(err)
	if suggestErr != nil || len(corrections) == 0 {
		attempt.NeedsAIHelp = true
		attempt.SuggestedPrompt = s.GenerateAIPrompt(err)
		return attempt
	}

	// Try to apply corrections
	for i := 0; i < s.maxRetries; i++ {
		attempt.AttemptNumber = i + 1
		for _, correction := range corrections {
			result, applyErr := s.correctionEngine.ApplyCorrection(ctx, correction, projectPath)
			if applyErr != nil {
				continue
			}

			attempt.AutoFixApplied = true
			if result.Success {
				attempt.AutoFixSuccess = true
				return attempt
			}
		}
	}

	// Auto-fix failed, need AI help
	attempt.NeedsAIHelp = true
	attempt.SuggestedPrompt = s.GenerateAIPrompt(err)
	return attempt
}

// GenerateAIPrompt generates a prompt for AI to fix an error
func (s *SelfCorrectionService) GenerateAIPrompt(err *domain.ErrorDetails) string {
	codeContext := ""
	if s.fileSystem != nil && err.SourceFile != "" {
		content, readErr := s.fileSystem.ReadFile(err.SourceFile)
		if readErr == nil {
			codeContext = extractCodeContext(string(content), err.LineNumber, 5, 5)
		}
	}
	return s.feedbackGen.GeneratePrompt(err, codeContext)
}

// ShouldRetryWithAI determines if the correction result warrants AI retry
func (s *SelfCorrectionService) ShouldRetryWithAI(result *CorrectionResult) bool {
	if result.AllFixed {
		return false
	}
	return len(result.AIPrompts) > 0
}

// GetCombinedAIPrompt returns a combined prompt for all remaining errors
func (s *SelfCorrectionService) GetCombinedAIPrompt(result *CorrectionResult) string {
	if len(result.RemainingErrors) == 0 {
		return ""
	}
	return s.feedbackGen.FormatErrorsForAI(result.RemainingErrors)
}

// extractCodeContext extracts code around a specific line
func extractCodeContext(source string, lineNumber, before, after int) string {
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
