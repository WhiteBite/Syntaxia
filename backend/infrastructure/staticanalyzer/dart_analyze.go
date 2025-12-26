package staticanalyzer

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"time"
)

// DartAnalyzeAnalyzer implements StaticAnalyzer for Dart using dart analyze
type DartAnalyzeAnalyzer struct {
	log domain.Logger
}

// NewDartAnalyzeAnalyzer creates a new Dart static analyzer
func NewDartAnalyzeAnalyzer(log domain.Logger) *DartAnalyzeAnalyzer {
	return &DartAnalyzeAnalyzer{log: log}
}

// Analyze performs static analysis on Dart code
func (a *DartAnalyzeAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running dart analyze for project: %s", config.ProjectPath))

	startTime := time.Now()

	// Check if dart is installed
	if err := a.checkDartInstalled(); err != nil {
		a.log.Warning(fmt.Sprintf("Dart not available: %v", err))
		return &domain.StaticAnalysisResult{
			Success:     true,
			Language:    config.Language,
			ProjectPath: config.ProjectPath,
			Analyzer:    config.Analyzer,
			Duration:    time.Since(startTime).Seconds(),
			Issues:      []*domain.StaticIssue{},
			Summary:     &domain.StaticAnalysisSummary{},
			Error:       fmt.Sprintf("Dart not available: %v", err),
		}, nil
	}

	// Build command
	args := []string{"analyze", "--format=machine"}
	args = append(args, config.ProjectPath)

	cmd := exec.CommandContext(ctx, "dart", args...)
	executil.HideWindow(cmd)
	cmd.Dir = config.ProjectPath

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime).Seconds()

	result := &domain.StaticAnalysisResult{
		Success:     true,
		Language:    config.Language,
		ProjectPath: config.ProjectPath,
		Analyzer:    config.Analyzer,
		Duration:    duration,
		Issues:      []*domain.StaticIssue{},
		Summary:     &domain.StaticAnalysisSummary{},
	}

	if err != nil {
		// dart analyze returns non-zero exit code when issues are found
		a.log.Warning(fmt.Sprintf("dart analyze found issues: %v", err))
	}

	// Parse output
	issues, parseErr := a.parseDartAnalyzeOutput(output)
	if parseErr != nil {
		a.log.Warning(fmt.Sprintf("Failed to parse dart analyze output: %v", parseErr))
		result.Error = fmt.Sprintf("Failed to parse output: %v", parseErr)
	} else {
		result.Issues = issues
	}

	result.Summary = generateSummary(issues)

	a.log.Info(fmt.Sprintf("dart analyze completed in %.2fs, found %d issues", duration, len(issues)))
	return result, nil
}

// GetSupportedLanguages returns supported languages
func (a *DartAnalyzeAnalyzer) GetSupportedLanguages() []string {
	return []string{"dart"}
}

// GetAnalyzerType returns the analyzer type
func (a *DartAnalyzeAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypeDartAnalyze
}

// ValidateConfig validates the configuration
func (a *DartAnalyzeAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	if config.Language != "dart" {
		return fmt.Errorf("DartAnalyzeAnalyzer only supports Dart language")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Check for pubspec.yaml
	pubspecPath := filepath.Join(config.ProjectPath, "pubspec.yaml")
	if !fileExists(pubspecPath) {
		return fmt.Errorf("pubspec.yaml not found in project path")
	}

	return nil
}

// checkDartInstalled checks if dart is installed
func (a *DartAnalyzeAnalyzer) checkDartInstalled() error {
	cmd := exec.Command("dart", "--version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dart not found: %w", err)
	}
	return nil
}

// parseDartAnalyzeOutput parses dart analyze machine format output
func (a *DartAnalyzeAnalyzer) parseDartAnalyzeOutput(output []byte) ([]*domain.StaticIssue, error) {
	lines := strings.Split(string(output), "\n")
	issues := make([]*domain.StaticIssue, 0, len(lines))

	// Machine format: SEVERITY|TYPE|ERROR_CODE|FILE|LINE|COLUMN|LENGTH|MESSAGE
	machineRe := regexp.MustCompile(`^(ERROR|WARNING|INFO)\|([^|]+)\|([^|]+)\|([^|]+)\|(\d+)\|(\d+)\|(\d+)\|(.+)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if matches := machineRe.FindStringSubmatch(line); len(matches) >= 9 {
			issue := &domain.StaticIssue{
				Severity: strings.ToLower(matches[1]),
				Category: matches[2],
				Code:     matches[3],
				File:     matches[4],
				Line:     parseInt(matches[5]),
				Column:   parseInt(matches[6]),
				Message:  matches[8],
			}
			issues = append(issues, issue)
		}
	}

	return issues, nil
}
