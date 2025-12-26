package repair

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Command execution constants
const (
	defaultCommandTimeout = 30 * time.Second
	goimportsTimeout      = 60 * time.Second
	linterTimeout         = 120 * time.Second
)

// CommandResult represents the result of a command execution
type CommandResult struct {
	Output   string
	ExitCode int
	Success  bool
	Duration time.Duration
}

// runCommand executes a command with timeout and returns the result
func runCommand(ctx context.Context, dir string, name string, args ...string) (*CommandResult, error) {
	return runCommandWithTimeout(ctx, dir, defaultCommandTimeout, name, args...)
}

// runCommandWithTimeout executes a command with a specific timeout
func runCommandWithTimeout(ctx context.Context, dir string, timeout time.Duration, name string, args ...string) (*CommandResult, error) {
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	if dir != "" {
		cmd.Dir = dir
	}

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	result := &CommandResult{
		Output:   string(output),
		Duration: duration,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else if ctx.Err() == context.DeadlineExceeded {
			return result, fmt.Errorf("command timed out after %v: %s %s", timeout, name, strings.Join(args, " "))
		} else {
			return result, fmt.Errorf("failed to execute command %s: %w", name, err)
		}
		result.Success = false
	} else {
		result.Success = true
		result.ExitCode = 0
	}

	return result, nil
}

// isToolAvailable checks if a tool is available in PATH
func isToolAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// getAvailableFormatter returns the best available formatter for a language
func getAvailableFormatter(language string) (string, []string) {
	switch language {
	case langGo:
		if isToolAvailable("goimports") {
			return "goimports", []string{"-w"}
		}
		if isToolAvailable("gofmt") {
			return "gofmt", []string{"-w"}
		}
	case langTypeScript, langJavaScript:
		if isToolAvailable("prettier") {
			return "prettier", []string{"--write"}
		}
		// Try npx prettier as fallback
		if isToolAvailable("npx") {
			return "npx", []string{"prettier", "--write"}
		}
	case langPython:
		if isToolAvailable("black") {
			return "black", nil
		}
		if isToolAvailable("autopep8") {
			return "autopep8", []string{"--in-place"}
		}
	}
	return "", nil
}

// getAvailableLinter returns the best available linter for a language
func getAvailableLinter(language string) (string, []string) {
	switch language {
	case langGo:
		if isToolAvailable("golangci-lint") {
			return "golangci-lint", []string{"run", "--fix"}
		}
	case langTypeScript, langJavaScript:
		if isToolAvailable("eslint") {
			return "eslint", []string{"--fix"}
		}
		if isToolAvailable("npx") {
			return "npx", []string{"eslint", "--fix"}
		}
	case langPython:
		if isToolAvailable("ruff") {
			return "ruff", []string{"check", "--fix"}
		}
		if isToolAvailable("pylint") {
			return "pylint", nil // pylint doesn't have auto-fix
		}
	}
	return "", nil
}

// detectLanguageFromExtension returns the language based on file extension
func detectLanguageFromExtension(ext string) string {
	switch ext {
	case extGo:
		return langGo
	case extTS, extTSX:
		return langTypeScript
	case extJS, extJSX:
		return langJavaScript
	case extVue:
		return langTypeScript // Vue files typically use TypeScript
	case extPy:
		return langPython
	default:
		return ""
	}
}

// ToolChecker provides methods to check tool availability
type ToolChecker struct {
	cache map[string]bool
}

// NewToolChecker creates a new ToolChecker instance
func NewToolChecker() *ToolChecker {
	return &ToolChecker{
		cache: make(map[string]bool),
	}
}

// IsAvailable checks if a tool is available (with caching)
func (tc *ToolChecker) IsAvailable(name string) bool {
	if available, ok := tc.cache[name]; ok {
		return available
	}
	available := isToolAvailable(name)
	tc.cache[name] = available
	return available
}

// GetAvailableTools returns a list of available tools for a language
func (tc *ToolChecker) GetAvailableTools(language string) []string {
	var tools []string

	switch language {
	case langGo:
		candidates := []string{"goimports", "gofmt", "golangci-lint", "go"}
		for _, tool := range candidates {
			if tc.IsAvailable(tool) {
				tools = append(tools, tool)
			}
		}
	case langTypeScript, langJavaScript:
		candidates := []string{"prettier", "eslint", "tsc", "npx"}
		for _, tool := range candidates {
			if tc.IsAvailable(tool) {
				tools = append(tools, tool)
			}
		}
	case langPython:
		candidates := []string{"black", "autopep8", "ruff", "pylint", "mypy"}
		for _, tool := range candidates {
			if tc.IsAvailable(tool) {
				tools = append(tools, tool)
			}
		}
	}

	return tools
}
