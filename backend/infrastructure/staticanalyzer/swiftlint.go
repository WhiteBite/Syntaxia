package staticanalyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syntaxia/domain"
)

// SwiftLintAnalyzer реализует StaticAnalyzer для Swift с использованием SwiftLint
type SwiftLintAnalyzer struct {
	log domain.Logger
}

// NewSwiftLintAnalyzer создает новый анализатор для Swift
func NewSwiftLintAnalyzer(log domain.Logger) *SwiftLintAnalyzer {
	return &SwiftLintAnalyzer{
		log: log,
	}
}

// swiftLintIssue represents a SwiftLint JSON output issue
type swiftLintIssue struct {
	RuleID   string `json:"rule_id"`
	Reason   string `json:"reason"`
	Line     int    `json:"line"`
	Column   int    `json:"column,omitempty"`
	File     string `json:"file"`
	Severity string `json:"severity"`
}

// Analyze выполняет статический анализ Swift кода
func (a *SwiftLintAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running SwiftLint analysis on: %s", config.ProjectPath))

	result := &domain.StaticAnalysisResult{
		Success:     false,
		Language:    config.Language,
		ProjectPath: config.ProjectPath,
		Analyzer:    domain.StaticAnalyzerTypeSwiftLint,
		Issues:      []*domain.StaticIssue{},
		Summary:     &domain.StaticAnalysisSummary{},
	}

	// Check if swiftlint is available
	if _, err := exec.LookPath("swiftlint"); err != nil {
		result.Error = "swiftlint not found in PATH. Install with: brew install swiftlint"
		result.Summary.WarningCount = 1
		return result, nil
	}

	// Check if project has Swift files
	hasSwiftFiles, err := a.hasSwiftFiles(config.ProjectPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to check for Swift files: %v", err)
		return result, nil
	}
	if !hasSwiftFiles {
		result.Success = true
		return result, nil
	}

	// Build swiftlint command
	args := []string{"lint", "--reporter", "json"}

	// Add path if specified
	if config.ProjectPath != "" {
		args = append(args, "--path", config.ProjectPath)
	}

	// Add config file if exists
	configFile := a.findConfigFile(config.ProjectPath)
	if configFile != "" {
		args = append(args, "--config", configFile)
	}

	cmd := exec.CommandContext(ctx, "swiftlint", args...)
	cmd.Dir = config.ProjectPath

	output, err := cmd.CombinedOutput()

	// SwiftLint returns non-zero exit code when issues are found
	// Parse the JSON output regardless of exit code
	issues, parseErr := a.parseJSONOutput(string(output), config.ProjectPath)
	if parseErr != nil {
		// Try to parse as plain text if JSON fails
		issues = a.parsePlainOutput(string(output), config.ProjectPath)
	}

	result.Issues = issues

	// Calculate summary
	for _, issue := range issues {
		switch issue.Severity {
		case "error":
			result.Summary.ErrorCount++
		case "warning":
			result.Summary.WarningCount++
		default:
			result.Summary.InfoCount++
		}
	}
	result.Summary.TotalIssues = len(issues)

	// Success if no errors (warnings are acceptable)
	result.Success = result.Summary.ErrorCount == 0

	if err != nil && result.Summary.ErrorCount > 0 {
		result.Error = fmt.Sprintf("SwiftLint found %d errors", result.Summary.ErrorCount)
	}

	a.log.Info(fmt.Sprintf("SwiftLint analysis complete: %d issues found", len(issues)))
	return result, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (a *SwiftLintAnalyzer) GetSupportedLanguages() []string {
	return []string{"swift"}
}

// GetAnalyzerType возвращает тип анализатора
func (a *SwiftLintAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypeSwiftLint
}

// ValidateConfig проверяет корректность конфигурации
func (a *SwiftLintAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	if _, err := os.Stat(config.ProjectPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", config.ProjectPath)
	}

	return nil
}

// parseJSONOutput parses SwiftLint JSON output
func (a *SwiftLintAnalyzer) parseJSONOutput(output, projectPath string) ([]*domain.StaticIssue, error) {
	var swiftIssues []swiftLintIssue

	// Find JSON array in output (may have other text before/after)
	startIdx := strings.Index(output, "[")
	endIdx := strings.LastIndex(output, "]")

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return nil, fmt.Errorf("no JSON array found in output")
	}

	jsonStr := output[startIdx : endIdx+1]

	if err := json.Unmarshal([]byte(jsonStr), &swiftIssues); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var issues []*domain.StaticIssue
	for _, si := range swiftIssues {
		// Make path relative if possible
		relPath := si.File
		if rel, err := filepath.Rel(projectPath, si.File); err == nil {
			relPath = rel
		}

		issue := &domain.StaticIssue{
			File:     relPath,
			Line:     si.Line,
			Column:   si.Column,
			Message:  si.Reason,
			Severity: a.mapSeverity(si.Severity),
			Code:     si.RuleID,
		}
		issues = append(issues, issue)
	}

	return issues, nil
}

// parsePlainOutput parses SwiftLint plain text output as fallback
func (a *SwiftLintAnalyzer) parsePlainOutput(output, projectPath string) []*domain.StaticIssue {
	var issues []*domain.StaticIssue

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Format: /path/to/file.swift:10:5: warning: Message (rule_id)
		parts := strings.SplitN(line, ":", 4)
		if len(parts) < 4 {
			continue
		}

		file := parts[0]
		var lineNum, col int
		var severityMsg string

		// Try to parse line and column
		fmt.Sscanf(parts[1], "%d", &lineNum)
		if _, err := fmt.Sscanf(parts[2], "%d", &col); err != nil {
			severityMsg = parts[2] + ":" + parts[3]
		} else {
			severityMsg = parts[3]
		}

		// Parse severity and message
		severityMsg = strings.TrimSpace(severityMsg)
		severity := "warning"
		message := severityMsg

		if strings.HasPrefix(severityMsg, "error:") {
			severity = "error"
			message = strings.TrimPrefix(severityMsg, "error:")
		} else if strings.HasPrefix(severityMsg, "warning:") {
			severity = "warning"
			message = strings.TrimPrefix(severityMsg, "warning:")
		}

		// Extract rule ID from message (usually in parentheses at end)
		ruleID := ""
		if idx := strings.LastIndex(message, "("); idx != -1 {
			if endIdx := strings.LastIndex(message, ")"); endIdx > idx {
				ruleID = message[idx+1 : endIdx]
				message = strings.TrimSpace(message[:idx])
			}
		}

		// Make path relative
		relPath := file
		if rel, err := filepath.Rel(projectPath, file); err == nil {
			relPath = rel
		}

		issues = append(issues, &domain.StaticIssue{
			File:     relPath,
			Line:     lineNum,
			Column:   col,
			Message:  strings.TrimSpace(message),
			Severity: severity,
			Code:     ruleID,
		})
	}

	return issues
}

// mapSeverity maps SwiftLint severity to standard severity
func (a *SwiftLintAnalyzer) mapSeverity(severity string) string {
	switch strings.ToLower(severity) {
	case "error":
		return "error"
	case "warning":
		return "warning"
	default:
		return "info"
	}
}

// findConfigFile finds SwiftLint configuration file
func (a *SwiftLintAnalyzer) findConfigFile(projectPath string) string {
	configFiles := []string{
		".swiftlint.yml",
		".swiftlint.yaml",
		"swiftlint.yml",
		"swiftlint.yaml",
	}

	for _, configFile := range configFiles {
		path := filepath.Join(projectPath, configFile)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// hasSwiftFiles checks if project contains Swift files
func (a *SwiftLintAnalyzer) hasSwiftFiles(projectPath string) (bool, error) {
	found := false

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if name == ".build" || name == ".git" || name == "Pods" ||
				name == "Carthage" || name == "DerivedData" {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(info.Name(), ".swift") {
			found = true
			return filepath.SkipAll
		}

		return nil
	})

	return found, err
}
