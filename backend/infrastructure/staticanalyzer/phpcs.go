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

// PHPCS analyzer constants
const (
	phpcsDefaultTimeout = 300 // 5 minutes
)

// PHPCSAnalyzer implements StaticAnalyzer for PHP using PHP_CodeSniffer
type PHPCSAnalyzer struct {
	log domain.Logger
}

// NewPHPCSAnalyzer creates a new analyzer for PHP
func NewPHPCSAnalyzer(log domain.Logger) *PHPCSAnalyzer {
	return &PHPCSAnalyzer{
		log: log,
	}
}

// Analyze performs static analysis on PHP code
func (a *PHPCSAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running PHPCS analysis for project: %s", config.ProjectPath))

	startTime := time.Now()

	// Check if PHPCS is installed
	if err := a.checkPHPCSInstalled(config.ProjectPath); err != nil {
		a.log.Warning(fmt.Sprintf("PHPCS not available: %v", err))
		return &domain.StaticAnalysisResult{
			Success:     true,
			Language:    config.Language,
			ProjectPath: config.ProjectPath,
			Analyzer:    config.Analyzer,
			Duration:    time.Since(startTime).Seconds(),
			Issues:      []*domain.StaticIssue{},
			Summary:     &domain.StaticAnalysisSummary{},
			Error:       fmt.Sprintf("PHPCS not available: %v", err),
		}, nil
	}

	// Build PHPCS command
	phpcsPath := a.findPHPCSPath(config.ProjectPath)
	args := []string{"--report=json"}

	// Add standard if specified
	if len(config.Rules) > 0 {
		args = append(args, "--standard="+strings.Join(config.Rules, ","))
	} else {
		// Default to PSR-12
		args = append(args, "--standard=PSR12")
	}

	// Add exclude rules if specified
	if len(config.ExcludeRules) > 0 {
		args = append(args, "--exclude="+strings.Join(config.ExcludeRules, ","))
	}

	// Add project path
	args = append(args, config.ProjectPath)

	// Create command
	cmd := exec.CommandContext(ctx, phpcsPath, args...)
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
		// PHPCS returns non-zero exit code when issues are found
		a.log.Warning(fmt.Sprintf("PHPCS found issues: %v", err))
	}

	// Parse JSON output
	issues, parseErr := a.parsePHPCSOutput(output)
	if parseErr != nil {
		a.log.Warning(fmt.Sprintf("Failed to parse PHPCS output: %v", parseErr))
		result.Error = fmt.Sprintf("Failed to parse output: %v", parseErr)
	} else {
		result.Issues = issues
	}

	// Generate summary
	result.Summary = generateSummary(issues)

	// Determine success based on error count
	result.Success = result.Summary.ErrorCount == 0

	a.log.Info(fmt.Sprintf("PHPCS analysis completed in %.2fs, found %d issues", duration, len(issues)))
	return result, nil
}

// GetSupportedLanguages returns supported languages
func (a *PHPCSAnalyzer) GetSupportedLanguages() []string {
	return []string{"php"}
}

// GetAnalyzerType returns the analyzer type
func (a *PHPCSAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypePHPCS
}

// ValidateConfig validates the configuration
func (a *PHPCSAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	if config.Language != "php" {
		return fmt.Errorf("PHPCS analyzer only supports PHP language")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Check if project contains PHP files
	if !a.hasPHPFiles(config.ProjectPath) {
		return fmt.Errorf("no PHP files found in project path")
	}

	return nil
}

// checkPHPCSInstalled checks if PHPCS is installed
func (a *PHPCSAnalyzer) checkPHPCSInstalled(projectPath string) error {
	phpcsPath := a.findPHPCSPath(projectPath)
	cmd := exec.Command(phpcsPath, "--version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("PHPCS not found: %w", err)
	}
	return nil
}

// findPHPCSPath finds the PHPCS executable
func (a *PHPCSAnalyzer) findPHPCSPath(projectPath string) string {
	// Check for vendor/bin/phpcs (Composer)
	vendorPath := filepath.Join(projectPath, "vendor", "bin", "phpcs")
	if _, err := os.Stat(vendorPath); err == nil {
		return vendorPath
	}

	// Check for Windows variant
	vendorPathWin := filepath.Join(projectPath, "vendor", "bin", "phpcs.bat")
	if _, err := os.Stat(vendorPathWin); err == nil {
		return vendorPathWin
	}

	// Fallback to global phpcs
	return "phpcs"
}

// hasPHPFiles checks if project contains PHP files
func (a *PHPCSAnalyzer) hasPHPFiles(projectPath string) bool {
	found := false
	_ = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if info.Name() == "vendor" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(info.Name(), ".php") {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

// phpcsReport represents PHPCS JSON output structure
type phpcsReport struct {
	Totals struct {
		Errors   int `json:"errors"`
		Warnings int `json:"warnings"`
		Fixable  int `json:"fixable"`
	} `json:"totals"`
	Files map[string]struct {
		Errors   int `json:"errors"`
		Warnings int `json:"warnings"`
		Messages []struct {
			Message  string `json:"message"`
			Source   string `json:"source"`
			Severity int    `json:"severity"`
			Fixable  bool   `json:"fixable"`
			Type     string `json:"type"`
			Line     int    `json:"line"`
			Column   int    `json:"column"`
		} `json:"messages"`
	} `json:"files"`
}

// parsePHPCSOutput parses PHPCS JSON output
func (a *PHPCSAnalyzer) parsePHPCSOutput(output []byte) ([]*domain.StaticIssue, error) {
	var report phpcsReport
	if err := json.Unmarshal(output, &report); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var issues []*domain.StaticIssue
	for filePath, fileReport := range report.Files {
		for _, msg := range fileReport.Messages {
			severity := severityWarning
			if msg.Type == "ERROR" {
				severity = severityError
			}

			issue := &domain.StaticIssue{
				File:     filePath,
				Line:     msg.Line,
				Column:   msg.Column,
				Severity: severity,
				Message:  msg.Message,
				Code:     msg.Source,
				Category: a.getCategory(msg.Source),
			}

			if msg.Fixable {
				issue.Suggestions = append(issue.Suggestions, "This issue can be auto-fixed with phpcbf")
			}

			issues = append(issues, issue)
		}
	}

	return issues, nil
}

// getCategory determines the category based on the rule source
func (a *PHPCSAnalyzer) getCategory(source string) string {
	if source == "" {
		return severityOther
	}

	// Parse source like "PSR12.Files.FileHeader.SpacingAfterBlock"
	parts := strings.Split(source, ".")
	if len(parts) == 0 {
		return severityOther
	}

	standard := parts[0]
	switch standard {
	case "PSR1", "PSR2", "PSR12":
		return "style"
	case "Generic":
		if len(parts) > 1 {
			switch parts[1] {
			case "CodeAnalysis":
				return "complexity"
			case "Commenting":
				return "documentation"
			case "Files":
				return "style"
			case "Formatting":
				return "style"
			case "Functions":
				return "style"
			case "Metrics":
				return "complexity"
			case "NamingConventions":
				return "naming"
			case "PHP":
				return "correctness"
			case "Strings":
				return "style"
			case "WhiteSpace":
				return "style"
			}
		}
		return "generic"
	case "Squiz":
		return "squiz"
	case "PEAR":
		return "pear"
	case "Zend":
		return "zend"
	default:
		return severityOther
	}
}
