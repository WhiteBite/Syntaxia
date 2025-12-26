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
	"syntaxia/internal/executil"
	"time"
)

// RuboCop analyzer constants
const (
	rubocopDefaultTimeout = 300 // 5 minutes
)

// RuboCopAnalyzer implements StaticAnalyzer for Ruby using RuboCop
type RuboCopAnalyzer struct {
	log domain.Logger
}

// NewRuboCopAnalyzer creates a new RuboCop analyzer
func NewRuboCopAnalyzer(log domain.Logger) *RuboCopAnalyzer {
	return &RuboCopAnalyzer{
		log: log,
	}
}

// Analyze performs static analysis on Ruby code
func (a *RuboCopAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running RuboCop analysis for project: %s", config.ProjectPath))

	startTime := time.Now()

	// Check if RuboCop is installed
	if err := a.checkRuboCopInstalled(); err != nil {
		return &domain.StaticAnalysisResult{
			Success:     false,
			Language:    config.Language,
			ProjectPath: config.ProjectPath,
			Analyzer:    config.Analyzer,
			Error:       err.Error(),
			Issues:      []*domain.StaticIssue{},
			Summary:     &domain.StaticAnalysisSummary{},
		}, nil
	}

	// Build command
	args := a.buildCommand(config)

	// Check if bundle exec should be used
	cmdName := "rubocop"
	if a.hasBundler(config.ProjectPath) {
		cmdName = "bundle"
		args = append([]string{"exec", "rubocop"}, args...)
	}

	// Create command
	cmd := exec.CommandContext(ctx, cmdName, args...)
	executil.HideWindow(cmd)
	cmd.Dir = config.ProjectPath

	// Set environment variables
	cmd.Env = os.Environ()
	if config.EnvVars != nil {
		for key, value := range config.EnvVars {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Run command
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
		// RuboCop returns non-zero exit code when issues are found
		a.log.Warning(fmt.Sprintf("RuboCop found issues: %v", err))
	}

	// Parse JSON output
	issues, parseErr := a.parseRuboCopOutput(output)
	if parseErr != nil {
		a.log.Warning(fmt.Sprintf("Failed to parse RuboCop output: %v", parseErr))
		result.Error = fmt.Sprintf("Failed to parse output: %v", parseErr)
	} else {
		result.Issues = issues
	}

	// Generate summary
	result.Summary = generateSummary(issues)

	// Determine success based on error count
	result.Success = result.Summary.ErrorCount == 0

	a.log.Info(fmt.Sprintf("RuboCop analysis completed in %.2fs, found %d issues", duration, len(issues)))
	return result, nil
}

// GetSupportedLanguages returns supported languages
func (a *RuboCopAnalyzer) GetSupportedLanguages() []string {
	return []string{"ruby", "rb"}
}

// GetAnalyzerType returns the analyzer type
func (a *RuboCopAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypeRuboCop
}

// ValidateConfig validates the configuration
func (a *RuboCopAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	if config.Language != "ruby" && config.Language != "rb" {
		return fmt.Errorf("RuboCop analyzer only supports Ruby language")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Check for Ruby files
	hasRubyFiles := false
	err := filepath.Walk(config.ProjectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(path, ".rb") {
			hasRubyFiles = true
			return filepath.SkipAll
		}
		return nil
	})

	if err == nil && !hasRubyFiles {
		return fmt.Errorf("no Ruby files found in project path")
	}

	return nil
}

// buildCommand builds the RuboCop command arguments
func (a *RuboCopAnalyzer) buildCommand(config *domain.StaticAnalyzerConfig) []string {
	args := []string{"--format", "json"}

	// Add specific cops if rules are specified
	if len(config.Rules) > 0 {
		args = append(args, "--only", strings.Join(config.Rules, ","))
	}

	// Add excluded cops
	if len(config.ExcludeRules) > 0 {
		args = append(args, "--except", strings.Join(config.ExcludeRules, ","))
	}

	// Add target path (default to current directory)
	args = append(args, ".")

	return args
}

// checkRuboCopInstalled checks if RuboCop is installed
func (a *RuboCopAnalyzer) checkRuboCopInstalled() error {
	cmd := exec.Command("rubocop", "--version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("RuboCop not found, install with: gem install rubocop")
	}
	return nil
}

// hasBundler checks if project uses Bundler
func (a *RuboCopAnalyzer) hasBundler(projectPath string) bool {
	gemfilePath := filepath.Join(projectPath, "Gemfile")
	if _, err := os.Stat(gemfilePath); err == nil {
		// Check if rubocop is in Gemfile
		content, err := os.ReadFile(gemfilePath)
		if err == nil && strings.Contains(string(content), "rubocop") {
			return true
		}
	}
	return false
}

// ruboCopOutput represents RuboCop JSON output structure
type ruboCopOutput struct {
	Metadata struct {
		RubocopVersion string `json:"rubocop_version"`
		RubyEngine     string `json:"ruby_engine"`
		RubyVersion    string `json:"ruby_version"`
	} `json:"metadata"`
	Files []struct {
		Path     string `json:"path"`
		Offenses []struct {
			Severity    string `json:"severity"`
			Message     string `json:"message"`
			CopName     string `json:"cop_name"`
			Corrected   bool   `json:"corrected"`
			Correctable bool   `json:"correctable"`
			Location    struct {
				StartLine   int `json:"start_line"`
				StartColumn int `json:"start_column"`
				LastLine    int `json:"last_line"`
				LastColumn  int `json:"last_column"`
				Length      int `json:"length"`
				Line        int `json:"line"`
				Column      int `json:"column"`
			} `json:"location"`
		} `json:"offenses"`
	} `json:"files"`
	Summary struct {
		OffenseCount        int `json:"offense_count"`
		TargetFileCount     int `json:"target_file_count"`
		InspectedFileCount  int `json:"inspected_file_count"`
		CorrectedFileCount  int `json:"corrected_file_count,omitempty"`
		CorrectableCount    int `json:"correctable_count,omitempty"`
		UncorrectableCount  int `json:"uncorrectable_count,omitempty"`
	} `json:"summary"`
}

// parseRuboCopOutput parses RuboCop JSON output
func (a *RuboCopAnalyzer) parseRuboCopOutput(output []byte) ([]*domain.StaticIssue, error) {
	var report ruboCopOutput
	if err := json.Unmarshal(output, &report); err != nil {
		// Try to find JSON in output (might have other text before/after)
		jsonStart := strings.Index(string(output), "{")
		jsonEnd := strings.LastIndex(string(output), "}")
		if jsonStart != -1 && jsonEnd != -1 && jsonEnd > jsonStart {
			jsonStr := string(output)[jsonStart : jsonEnd+1]
			if err := json.Unmarshal([]byte(jsonStr), &report); err != nil {
				return nil, fmt.Errorf("failed to parse RuboCop JSON output: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to parse RuboCop JSON output: %w", err)
		}
	}

	var issues []*domain.StaticIssue

	for _, file := range report.Files {
		for _, offense := range file.Offenses {
			issue := &domain.StaticIssue{
				File:        file.Path,
				Line:        offense.Location.Line,
				Column:      offense.Location.Column,
				Severity:    a.convertSeverity(offense.Severity),
				Message:     offense.Message,
				Code:        offense.CopName,
				Category:    a.getCategory(offense.CopName),
				Suggestions: []string{},
			}

			// Add auto-correct suggestion if available
			if offense.Correctable {
				issue.Suggestions = append(issue.Suggestions,
					fmt.Sprintf("Run 'rubocop -a %s' to auto-correct this issue", file.Path))
			}

			issues = append(issues, issue)
		}
	}

	return issues, nil
}

// convertSeverity converts RuboCop severity to standard format
func (a *RuboCopAnalyzer) convertSeverity(ruboCopSeverity string) string {
	switch strings.ToLower(ruboCopSeverity) {
	case "fatal", "error":
		return severityError
	case "warning":
		return severityWarning
	case "convention", "refactor":
		return severityInfo
	case "info":
		return severityHint
	default:
		return severityWarning
	}
}

// getCategory determines the category based on cop name
func (a *RuboCopAnalyzer) getCategory(copName string) string {
	if copName == "" {
		return severityOther
	}

	// RuboCop cop naming convention: Department/CopName
	parts := strings.Split(copName, "/")
	if len(parts) < 1 {
		return severityOther
	}

	department := strings.ToLower(parts[0])

	switch department {
	case "layout":
		return "formatting"
	case "lint":
		return "lint"
	case "metrics":
		return "complexity"
	case "naming":
		return "naming"
	case "security":
		return "security"
	case "style":
		return "style"
	case "performance":
		return "performance"
	case "bundler":
		return "bundler"
	case "gemspec":
		return "gemspec"
	case "rails":
		return "rails"
	case "rspec":
		return "rspec"
	default:
		return department
	}
}
