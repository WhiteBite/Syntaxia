package verification

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"syntaxia/domain"
)

// Default timeout and configuration constants
const (
	DefaultBuildTimeout    = 5 * time.Minute
	DefaultTestTimeout     = 10 * time.Minute
	DefaultAnalysisTimeout = 3 * time.Minute
	DefaultFastModeTimeout = 1 * time.Minute
)

// VerificationMode defines the verification mode
type VerificationMode string

const (
	// ModeFull runs full verification pipeline
	ModeFull VerificationMode = "full"
	// ModeFast runs only on changed files
	ModeFast VerificationMode = "fast"
	// ModeIncremental runs incremental verification (only new errors)
	ModeIncremental VerificationMode = "incremental"
)

// ChatIntegrationConfig configures verification for AI chat integration
type ChatIntegrationConfig struct {
	ProjectPath   string           `json:"projectPath"`
	Languages     []string         `json:"languages"`
	ChangedFiles  []string         `json:"changedFiles"`
	Mode          VerificationMode `json:"mode"`
	Timeout       time.Duration    `json:"timeout"`
	SkipTests     bool             `json:"skipTests"`
	SkipLinting   bool             `json:"skipLinting"`
	BaselineState *BaselineState   `json:"baselineState,omitempty"`
}

// BaselineState stores the state before changes for incremental verification
type BaselineState struct {
	BuildErrors    []IssueFingerprint `json:"buildErrors"`
	TypeErrors     []IssueFingerprint `json:"typeErrors"`
	LintIssues     []IssueFingerprint `json:"lintIssues"`
	TestFailures   []string           `json:"testFailures"`
	CapturedAt     time.Time          `json:"capturedAt"`
}

// IssueFingerprint uniquely identifies an issue for comparison
type IssueFingerprint struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// ChatVerificationResult contains verification results for AI chat
type ChatVerificationResult struct {
	Success       bool                `json:"success"`
	Mode          VerificationMode    `json:"mode"`
	Duration      time.Duration       `json:"duration"`
	NewErrors     []VerificationIssue `json:"newErrors"`
	FixedErrors   []VerificationIssue `json:"fixedErrors"`
	BuildResult   *StepResult         `json:"buildResult,omitempty"`
	TestResult    *StepResult         `json:"testResult,omitempty"`
	LintResult    *StepResult         `json:"lintResult,omitempty"`
	Summary       string              `json:"summary"`
	Timeout       bool                `json:"timeout"`
	ChangedFiles  []string            `json:"changedFiles"`
}

// VerificationIssue represents a single verification issue
type VerificationIssue struct {
	Type     string `json:"type"` // "build", "type", "lint", "test"
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Code     string `json:"code,omitempty"`
}

// StepResult contains result of a single verification step
type StepResult struct {
	Success  bool              `json:"success"`
	Duration time.Duration     `json:"duration"`
	Issues   []VerificationIssue `json:"issues"`
	Output   string            `json:"output,omitempty"`
	Skipped  bool              `json:"skipped"`
	Timeout  bool              `json:"timeout"`
}

// ChatIntegrationService provides verification for AI chat integration
type ChatIntegrationService struct {
	log            domain.Logger
	buildService   domain.IBuildService
	testService    domain.ITestService
	staticAnalyzer domain.IStaticAnalyzerService

	baselineMu sync.RWMutex
	baselines  map[string]*BaselineState // projectPath -> baseline
}

// NewChatIntegrationService creates a new chat integration service
func NewChatIntegrationService(
	log domain.Logger,
	buildService domain.IBuildService,
	testService domain.ITestService,
	staticAnalyzer domain.IStaticAnalyzerService,
) *ChatIntegrationService {
	return &ChatIntegrationService{
		log:            log,
		buildService:   buildService,
		testService:    testService,
		staticAnalyzer: staticAnalyzer,
		baselines:      make(map[string]*BaselineState),
	}
}

// CaptureBaseline captures the current state before changes
func (s *ChatIntegrationService) CaptureBaseline(
	ctx context.Context,
	projectPath string,
	languages []string,
) (*BaselineState, error) {
	s.log.Info(fmt.Sprintf("Capturing baseline state for project: %s", projectPath))

	baseline := &BaselineState{
		CapturedAt: time.Now().UTC(),
	}

	// Capture build/type errors
	for _, lang := range languages {
		typeResult, err := s.buildService.TypeCheck(ctx, projectPath, lang)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Failed to capture type check baseline for %s: %v", lang, err))
			continue
		}
		for _, issue := range typeResult.Issues {
			baseline.TypeErrors = append(baseline.TypeErrors, IssueFingerprint{
				File:     issue.File,
				Line:     issue.Line,
				Code:     issue.Code,
				Message:  issue.Message,
				Severity: issue.Severity,
			})
		}
	}

	// Capture lint issues
	report, err := s.staticAnalyzer.AnalyzeProject(ctx, projectPath, languages)
	if err != nil {
		s.log.Warning(fmt.Sprintf("Failed to capture lint baseline: %v", err))
	} else if report != nil && report.Results != nil {
		for _, result := range report.Results {
			for _, issue := range result.Issues {
				baseline.LintIssues = append(baseline.LintIssues, IssueFingerprint{
					File:     issue.File,
					Line:     issue.Line,
					Code:     issue.Code,
					Message:  issue.Message,
					Severity: issue.Severity,
				})
			}
		}
	}

	// Store baseline
	s.baselineMu.Lock()
	s.baselines[projectPath] = baseline
	s.baselineMu.Unlock()

	s.log.Info(fmt.Sprintf("Baseline captured: %d type errors, %d lint issues",
		len(baseline.TypeErrors), len(baseline.LintIssues)))

	return baseline, nil
}

// GetBaseline returns the stored baseline for a project
func (s *ChatIntegrationService) GetBaseline(projectPath string) *BaselineState {
	s.baselineMu.RLock()
	defer s.baselineMu.RUnlock()
	return s.baselines[projectPath]
}

// ClearBaseline removes the stored baseline for a project
func (s *ChatIntegrationService) ClearBaseline(projectPath string) {
	s.baselineMu.Lock()
	delete(s.baselines, projectPath)
	s.baselineMu.Unlock()
}

// RunAfterChanges runs verification after AI makes changes
func (s *ChatIntegrationService) RunAfterChanges(
	ctx context.Context,
	config *ChatIntegrationConfig,
) (*ChatVerificationResult, error) {
	s.log.Info(fmt.Sprintf("Running verification after changes: mode=%s, files=%d",
		config.Mode, len(config.ChangedFiles)))

	startTime := time.Now()
	result := &ChatVerificationResult{
		Mode:         config.Mode,
		ChangedFiles: config.ChangedFiles,
		NewErrors:    make([]VerificationIssue, 0),
		FixedErrors:  make([]VerificationIssue, 0),
	}

	// Set timeout
	timeout := config.Timeout
	if timeout == 0 {
		timeout = s.getDefaultTimeout(config.Mode)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Run verification steps based on mode
	var buildResult, testResult, lintResult *StepResult

	switch config.Mode {
	case ModeFast:
		buildResult, lintResult = s.runFastVerification(ctx, config)
	case ModeIncremental:
		buildResult, testResult, lintResult = s.runIncrementalVerification(ctx, config)
	default:
		buildResult, testResult, lintResult = s.runFullVerification(ctx, config)
	}

	result.BuildResult = buildResult
	result.TestResult = testResult
	result.LintResult = lintResult
	result.Duration = time.Since(startTime)

	// Check for timeout
	if ctx.Err() == context.DeadlineExceeded {
		result.Timeout = true
		result.Summary = "Verification timed out"
		return result, nil
	}

	// Collect all issues
	allIssues := s.collectAllIssues(buildResult, testResult, lintResult)

	// Filter new errors if baseline exists
	if config.BaselineState != nil || config.Mode == ModeIncremental {
		baseline := config.BaselineState
		if baseline == nil {
			baseline = s.GetBaseline(config.ProjectPath)
		}
		if baseline != nil {
			result.NewErrors, result.FixedErrors = s.filterNewErrors(allIssues, baseline)
		} else {
			result.NewErrors = allIssues
		}
	} else {
		result.NewErrors = allIssues
	}

	// Determine success
	result.Success = s.determineSuccess(result)
	result.Summary = s.generateSummary(result)

	s.log.Info(fmt.Sprintf("Verification completed: success=%t, new_errors=%d, fixed=%d",
		result.Success, len(result.NewErrors), len(result.FixedErrors)))

	return result, nil
}

// runFastVerification runs quick verification on changed files only
func (s *ChatIntegrationService) runFastVerification(
	ctx context.Context,
	config *ChatIntegrationConfig,
) (*StepResult, *StepResult) {
	s.log.Info("Running fast verification mode")

	buildResult := &StepResult{Skipped: true}
	lintResult := &StepResult{Skipped: true}

	if len(config.ChangedFiles) == 0 {
		return buildResult, lintResult
	}

	// Detect languages from changed files
	languages := s.detectLanguagesFromFiles(config.ChangedFiles, config.Languages)

	// Run type check (fast)
	if len(languages) > 0 {
		buildResult = s.runBuildStep(ctx, config.ProjectPath, languages)
	}

	// Run lint on changed files only
	if !config.SkipLinting && len(languages) > 0 {
		lintResult = s.runLintStepForFiles(ctx, config.ProjectPath, config.ChangedFiles, languages)
	}

	return buildResult, lintResult
}

// runIncrementalVerification runs verification and filters to new errors only
func (s *ChatIntegrationService) runIncrementalVerification(
	ctx context.Context,
	config *ChatIntegrationConfig,
) (*StepResult, *StepResult, *StepResult) {
	s.log.Info("Running incremental verification mode")

	buildResult := s.runBuildStep(ctx, config.ProjectPath, config.Languages)

	var testResult *StepResult
	if !config.SkipTests {
		testResult = s.runTestStepForChangedFiles(ctx, config)
	} else {
		testResult = &StepResult{Skipped: true}
	}

	var lintResult *StepResult
	if !config.SkipLinting {
		lintResult = s.runLintStep(ctx, config.ProjectPath, config.Languages)
	} else {
		lintResult = &StepResult{Skipped: true}
	}

	return buildResult, testResult, lintResult
}

// runFullVerification runs complete verification pipeline
func (s *ChatIntegrationService) runFullVerification(
	ctx context.Context,
	config *ChatIntegrationConfig,
) (*StepResult, *StepResult, *StepResult) {
	s.log.Info("Running full verification mode")

	buildResult := s.runBuildStep(ctx, config.ProjectPath, config.Languages)

	// Skip tests if build failed
	var testResult *StepResult
	if !config.SkipTests && buildResult.Success {
		testResult = s.runTestStep(ctx, config.ProjectPath, config.Languages)
	} else if config.SkipTests {
		testResult = &StepResult{Skipped: true}
	} else {
		testResult = &StepResult{Skipped: true, Output: "Skipped due to build failure"}
	}

	var lintResult *StepResult
	if !config.SkipLinting {
		lintResult = s.runLintStep(ctx, config.ProjectPath, config.Languages)
	} else {
		lintResult = &StepResult{Skipped: true}
	}

	return buildResult, testResult, lintResult
}

// runBuildStep executes build/type check step
func (s *ChatIntegrationService) runBuildStep(
	ctx context.Context,
	projectPath string,
	languages []string,
) *StepResult {
	startTime := time.Now()
	result := &StepResult{
		Issues: make([]VerificationIssue, 0),
	}

	for _, lang := range languages {
		select {
		case <-ctx.Done():
			result.Timeout = true
			result.Duration = time.Since(startTime)
			return result
		default:
		}

		typeResult, err := s.buildService.TypeCheck(ctx, projectPath, lang)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Type check failed for %s: %v", lang, err))
			continue
		}

		if typeResult != nil {
			result.Output += typeResult.Output + "\n"
			for _, issue := range typeResult.Issues {
				result.Issues = append(result.Issues, VerificationIssue{
					Type:     "type",
					File:     issue.File,
					Line:     issue.Line,
					Column:   issue.Column,
					Severity: issue.Severity,
					Message:  issue.Message,
					Code:     issue.Code,
				})
			}
		}
	}

	result.Success = len(result.Issues) == 0 || !s.hasErrors(result.Issues)
	result.Duration = time.Since(startTime)
	return result
}

// runTestStep executes test step
func (s *ChatIntegrationService) runTestStep(
	ctx context.Context,
	projectPath string,
	languages []string,
) *StepResult {
	startTime := time.Now()
	result := &StepResult{
		Issues: make([]VerificationIssue, 0),
	}

	for _, lang := range languages {
		select {
		case <-ctx.Done():
			result.Timeout = true
			result.Duration = time.Since(startTime)
			return result
		default:
		}

		testResults, err := s.testService.RunSmokeTests(ctx, projectPath, lang)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Smoke tests failed for %s: %v", lang, err))
			continue
		}

		for _, tr := range testResults {
			if !tr.Success {
				result.Issues = append(result.Issues, VerificationIssue{
					Type:     "test",
					File:     tr.TestPath,
					Severity: "error",
					Message:  fmt.Sprintf("Test failed: %s - %s", tr.TestName, tr.Error),
				})
			}
			result.Output += tr.Output + "\n"
		}
	}

	result.Success = len(result.Issues) == 0
	result.Duration = time.Since(startTime)
	return result
}

// runTestStepForChangedFiles runs tests only for affected files
func (s *ChatIntegrationService) runTestStepForChangedFiles(
	ctx context.Context,
	config *ChatIntegrationConfig,
) *StepResult {
	startTime := time.Now()
	result := &StepResult{
		Issues: make([]VerificationIssue, 0),
	}

	if len(config.ChangedFiles) == 0 {
		result.Skipped = true
		return result
	}

	for _, lang := range config.Languages {
		select {
		case <-ctx.Done():
			result.Timeout = true
			result.Duration = time.Since(startTime)
			return result
		default:
		}

		testConfig := &domain.TestConfig{
			Language:    lang,
			ProjectPath: config.ProjectPath,
			Scope:       domain.TestScopeAffectedSmoke,
			Timeout:     int(DefaultTestTimeout.Seconds()),
		}

		testResults, err := s.testService.RunTargetedTests(ctx, testConfig, config.ChangedFiles)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Targeted tests failed for %s: %v", lang, err))
			continue
		}

		for _, tr := range testResults {
			if !tr.Success {
				result.Issues = append(result.Issues, VerificationIssue{
					Type:     "test",
					File:     tr.TestPath,
					Severity: "error",
					Message:  fmt.Sprintf("Test failed: %s - %s", tr.TestName, tr.Error),
				})
			}
			result.Output += tr.Output + "\n"
		}
	}

	result.Success = len(result.Issues) == 0
	result.Duration = time.Since(startTime)
	return result
}

// runLintStep executes static analysis step
func (s *ChatIntegrationService) runLintStep(
	ctx context.Context,
	projectPath string,
	languages []string,
) *StepResult {
	startTime := time.Now()
	result := &StepResult{
		Issues: make([]VerificationIssue, 0),
	}

	select {
	case <-ctx.Done():
		result.Timeout = true
		result.Duration = time.Since(startTime)
		return result
	default:
	}

	report, err := s.staticAnalyzer.AnalyzeProject(ctx, projectPath, languages)
	if err != nil {
		s.log.Warning(fmt.Sprintf("Static analysis failed: %v", err))
		result.Duration = time.Since(startTime)
		return result
	}

	if report != nil && report.Results != nil {
		for _, analysisResult := range report.Results {
			for _, issue := range analysisResult.Issues {
				result.Issues = append(result.Issues, VerificationIssue{
					Type:     "lint",
					File:     issue.File,
					Line:     issue.Line,
					Column:   issue.Column,
					Severity: issue.Severity,
					Message:  issue.Message,
					Code:     issue.Code,
				})
			}
		}
	}

	result.Success = !s.hasErrors(result.Issues)
	result.Duration = time.Since(startTime)
	return result
}

// runLintStepForFiles runs lint only on specific files
func (s *ChatIntegrationService) runLintStepForFiles(
	ctx context.Context,
	projectPath string,
	files []string,
	languages []string,
) *StepResult {
	startTime := time.Now()
	result := &StepResult{
		Issues: make([]VerificationIssue, 0),
	}

	// Create a set of changed files for quick lookup
	changedSet := make(map[string]bool)
	for _, f := range files {
		changedSet[f] = true
	}

	// Run full analysis but filter to changed files
	report, err := s.staticAnalyzer.AnalyzeProject(ctx, projectPath, languages)
	if err != nil {
		s.log.Warning(fmt.Sprintf("Static analysis failed: %v", err))
		result.Duration = time.Since(startTime)
		return result
	}

	if report != nil && report.Results != nil {
		for _, analysisResult := range report.Results {
			for _, issue := range analysisResult.Issues {
				// Only include issues from changed files
				if changedSet[issue.File] {
					result.Issues = append(result.Issues, VerificationIssue{
						Type:     "lint",
						File:     issue.File,
						Line:     issue.Line,
						Column:   issue.Column,
						Severity: issue.Severity,
						Message:  issue.Message,
						Code:     issue.Code,
					})
				}
			}
		}
	}

	result.Success = !s.hasErrors(result.Issues)
	result.Duration = time.Since(startTime)
	return result
}

// filterNewErrors filters issues to show only new ones (not in baseline)
func (s *ChatIntegrationService) filterNewErrors(
	current []VerificationIssue,
	baseline *BaselineState,
) (newErrors, fixedErrors []VerificationIssue) {
	// Build baseline fingerprint set
	baselineSet := make(map[string]bool)
	for _, fp := range baseline.TypeErrors {
		baselineSet[s.fingerprintKey(fp)] = true
	}
	for _, fp := range baseline.LintIssues {
		baselineSet[s.fingerprintKey(fp)] = true
	}

	// Find new errors
	currentSet := make(map[string]bool)
	for _, issue := range current {
		fp := IssueFingerprint{
			File:     issue.File,
			Line:     issue.Line,
			Code:     issue.Code,
			Message:  issue.Message,
			Severity: issue.Severity,
		}
		key := s.fingerprintKey(fp)
		currentSet[key] = true

		if !baselineSet[key] {
			newErrors = append(newErrors, issue)
		}
	}

	// Find fixed errors (in baseline but not in current)
	for _, fp := range baseline.TypeErrors {
		if !currentSet[s.fingerprintKey(fp)] {
			fixedErrors = append(fixedErrors, VerificationIssue{
				Type:     "type",
				File:     fp.File,
				Line:     fp.Line,
				Severity: fp.Severity,
				Message:  fp.Message,
				Code:     fp.Code,
			})
		}
	}
	for _, fp := range baseline.LintIssues {
		if !currentSet[s.fingerprintKey(fp)] {
			fixedErrors = append(fixedErrors, VerificationIssue{
				Type:     "lint",
				File:     fp.File,
				Line:     fp.Line,
				Severity: fp.Severity,
				Message:  fp.Message,
				Code:     fp.Code,
			})
		}
	}

	return newErrors, fixedErrors
}

// fingerprintKey creates a unique key for an issue fingerprint
func (s *ChatIntegrationService) fingerprintKey(fp IssueFingerprint) string {
	return fmt.Sprintf("%s:%d:%s:%s", fp.File, fp.Line, fp.Code, fp.Message)
}

// collectAllIssues collects all issues from step results
func (s *ChatIntegrationService) collectAllIssues(
	build, test, lint *StepResult,
) []VerificationIssue {
	var all []VerificationIssue

	if build != nil && !build.Skipped {
		all = append(all, build.Issues...)
	}
	if test != nil && !test.Skipped {
		all = append(all, test.Issues...)
	}
	if lint != nil && !lint.Skipped {
		all = append(all, lint.Issues...)
	}

	return all
}

// hasErrors checks if any issues are errors
func (s *ChatIntegrationService) hasErrors(issues []VerificationIssue) bool {
	for _, issue := range issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

// determineSuccess determines overall success based on results
func (s *ChatIntegrationService) determineSuccess(result *ChatVerificationResult) bool {
	// Check for timeout
	if result.Timeout {
		return false
	}

	// Check build result
	if result.BuildResult != nil && !result.BuildResult.Skipped && !result.BuildResult.Success {
		return false
	}

	// Check test result
	if result.TestResult != nil && !result.TestResult.Skipped && !result.TestResult.Success {
		return false
	}

	// Check for new errors (only errors, not warnings)
	for _, issue := range result.NewErrors {
		if issue.Severity == "error" {
			return false
		}
	}

	return true
}

// generateSummary generates a human-readable summary
func (s *ChatIntegrationService) generateSummary(result *ChatVerificationResult) string {
	var parts []string

	if result.Timeout {
		parts = append(parts, "⏱️ Verification timed out")
	}

	if result.BuildResult != nil && !result.BuildResult.Skipped {
		if result.BuildResult.Success {
			parts = append(parts, "✅ Build passed")
		} else {
			parts = append(parts, fmt.Sprintf("❌ Build failed (%d issues)", len(result.BuildResult.Issues)))
		}
	}

	if result.TestResult != nil && !result.TestResult.Skipped {
		if result.TestResult.Success {
			parts = append(parts, "✅ Tests passed")
		} else {
			parts = append(parts, fmt.Sprintf("❌ Tests failed (%d failures)", len(result.TestResult.Issues)))
		}
	}

	if result.LintResult != nil && !result.LintResult.Skipped {
		if result.LintResult.Success {
			parts = append(parts, "✅ Lint passed")
		} else {
			errorCount := 0
			for _, issue := range result.LintResult.Issues {
				if issue.Severity == "error" {
					errorCount++
				}
			}
			if errorCount > 0 {
				parts = append(parts, fmt.Sprintf("⚠️ Lint: %d errors, %d warnings",
					errorCount, len(result.LintResult.Issues)-errorCount))
			} else {
				parts = append(parts, fmt.Sprintf("⚠️ Lint: %d warnings", len(result.LintResult.Issues)))
			}
		}
	}

	if len(result.NewErrors) > 0 {
		parts = append(parts, fmt.Sprintf("🆕 %d new issues", len(result.NewErrors)))
	}

	if len(result.FixedErrors) > 0 {
		parts = append(parts, fmt.Sprintf("🔧 %d issues fixed", len(result.FixedErrors)))
	}

	if len(parts) == 0 {
		return "No verification steps executed"
	}

	return strings.Join(parts, " | ")
}

// getDefaultTimeout returns default timeout based on mode
func (s *ChatIntegrationService) getDefaultTimeout(mode VerificationMode) time.Duration {
	switch mode {
	case ModeFast:
		return DefaultFastModeTimeout
	case ModeIncremental:
		return DefaultBuildTimeout + DefaultTestTimeout
	default:
		return DefaultBuildTimeout + DefaultTestTimeout + DefaultAnalysisTimeout
	}
}

// detectLanguagesFromFiles detects languages from file extensions
func (s *ChatIntegrationService) detectLanguagesFromFiles(
	files []string,
	fallback []string,
) []string {
	langSet := make(map[string]bool)

	for _, file := range files {
		lang := s.detectLanguageFromFile(file)
		if lang != "" {
			langSet[lang] = true
		}
	}

	if len(langSet) == 0 {
		return fallback
	}

	languages := make([]string, 0, len(langSet))
	for lang := range langSet {
		languages = append(languages, lang)
	}
	return languages
}

// detectLanguageFromFile detects language from file extension
func (s *ChatIntegrationService) detectLanguageFromFile(file string) string {
	switch {
	case strings.HasSuffix(file, ".go"):
		return "go"
	case strings.HasSuffix(file, ".ts"), strings.HasSuffix(file, ".tsx"):
		return "typescript"
	case strings.HasSuffix(file, ".js"), strings.HasSuffix(file, ".jsx"):
		return "javascript"
	case strings.HasSuffix(file, ".py"):
		return "python"
	case strings.HasSuffix(file, ".java"):
		return "java"
	case strings.HasSuffix(file, ".rs"):
		return "rust"
	case strings.HasSuffix(file, ".kt"):
		return "kotlin"
	case strings.HasSuffix(file, ".cs"):
		return "csharp"
	case strings.HasSuffix(file, ".cpp"), strings.HasSuffix(file, ".cc"), strings.HasSuffix(file, ".c"):
		return "cpp"
	case strings.HasSuffix(file, ".rb"):
		return "ruby"
	case strings.HasSuffix(file, ".php"):
		return "php"
	case strings.HasSuffix(file, ".swift"):
		return "swift"
	case strings.HasSuffix(file, ".dart"):
		return "dart"
	default:
		return ""
	}
}
