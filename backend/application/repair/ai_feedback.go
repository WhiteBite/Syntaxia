package repair

import (
	"fmt"
	"strings"
	"syntaxia/domain"
)

// AIFeedbackGenerator generates prompts for AI to fix errors
type AIFeedbackGenerator struct {
	templates map[domain.ErrorType]string
}

// NewAIFeedbackGenerator creates a new AIFeedbackGenerator
func NewAIFeedbackGenerator() *AIFeedbackGenerator {
	g := &AIFeedbackGenerator{
		templates: make(map[domain.ErrorType]string),
	}
	g.registerDefaultTemplates()
	return g
}

// registerDefaultTemplates registers default prompt templates for each error type
func (g *AIFeedbackGenerator) registerDefaultTemplates() {
	g.templates[domain.ErrorTypeImport] = `The following import error occurred after your changes:
File: %s:%d
Error: %s

Please fix this import error. Common solutions:
- Add missing import statement
- Fix import path
- Remove unused import

The relevant code is:
%s`

	g.templates[domain.ErrorTypeSyntax] = `The following syntax error occurred after your changes:
File: %s:%d
Error: %s

Please fix this syntax error. Check for:
- Missing brackets, parentheses, or braces
- Missing semicolons or commas
- Incorrect keyword usage

The relevant code is:
%s`

	g.templates[domain.ErrorTypeTypeCheck] = `The following type error occurred after your changes:
File: %s:%d
Error: %s

Please fix this type error. Common issues:
- Property does not exist on type
- Type mismatch in assignment
- Missing type annotation

The relevant code is:
%s`

	g.templates[domain.ErrorTypeCompilation] = `The following compilation error occurred after your changes:
File: %s:%d
Error: %s

Please fix this compilation error.

The relevant code is:
%s`

	g.templates[domain.ErrorTypeLinting] = `The following linting error occurred after your changes:
File: %s:%d
Error: %s
Tool: %s

Please fix this linting issue.

The relevant code is:
%s`

	g.templates[domain.ErrorTypeTesting] = `The following test failure occurred after your changes:
File: %s:%d
Error: %s

Please fix this test failure. Consider:
- Updating test expectations
- Fixing the implementation
- Checking test setup

The relevant code is:
%s`

	g.templates[domain.ErrorTypeDependency] = `The following dependency error occurred:
File: %s:%d
Error: %s

Please fix this dependency issue. Check:
- Package installation
- Version compatibility
- Module configuration

The relevant code is:
%s`

	g.templates[domain.ErrorTypeLogic] = `The following logic error was detected:
File: %s:%d
Error: %s

Please review and fix this logic issue.

The relevant code is:
%s`
}

// GeneratePrompt generates a prompt for AI to fix an error
func (g *AIFeedbackGenerator) GeneratePrompt(err *domain.ErrorDetails, codeContext string) string {
	template, ok := g.templates[err.ErrorType]
	if !ok {
		template = g.templates[domain.ErrorTypeCompilation]
	}

	lineNum := err.LineNumber
	if lineNum <= 0 {
		lineNum = 1
	}

	// Format based on error type
	if err.ErrorType == domain.ErrorTypeLinting {
		return fmt.Sprintf(template, err.SourceFile, lineNum, err.Message, err.Tool, codeContext)
	}

	return fmt.Sprintf(template, err.SourceFile, lineNum, err.Message, codeContext)
}

// FormatErrorsForAI formats multiple errors into a single prompt for AI
func (g *AIFeedbackGenerator) FormatErrorsForAI(errors []*domain.ErrorDetails) string {
	if len(errors) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("The following errors occurred after your changes. Please fix all of them:\n\n")

	for i, err := range errors {
		sb.WriteString(fmt.Sprintf("## Error %d\n", i+1))
		sb.WriteString(fmt.Sprintf("- **Type**: %s\n", err.ErrorType))
		sb.WriteString(fmt.Sprintf("- **File**: %s", err.SourceFile))
		if err.LineNumber > 0 {
			sb.WriteString(fmt.Sprintf(":%d", err.LineNumber))
		}
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("- **Message**: %s\n", err.Message))
		if err.Tool != "" {
			sb.WriteString(fmt.Sprintf("- **Tool**: %s\n", err.Tool))
		}
		if len(err.Suggestions) > 0 {
			sb.WriteString("- **Suggestions**:\n")
			for _, suggestion := range err.Suggestions {
				sb.WriteString(fmt.Sprintf("  - %s\n", suggestion))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("Please provide fixes for all errors above. ")
	sb.WriteString("If errors are related, fix them together. ")
	sb.WriteString("Provide the corrected code for each affected file.\n")

	return sb.String()
}

// FormatErrorWithContext formats a single error with code context for AI
func (g *AIFeedbackGenerator) FormatErrorWithContext(err *domain.ErrorDetails, codeSnippet string) string {
	var sb strings.Builder

	sb.WriteString("The following error occurred after your changes:\n\n")
	sb.WriteString(fmt.Sprintf("**File**: %s", err.SourceFile))
	if err.LineNumber > 0 {
		sb.WriteString(fmt.Sprintf(":%d", err.LineNumber))
		if err.Column > 0 {
			sb.WriteString(fmt.Sprintf(":%d", err.Column))
		}
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("**Error Type**: %s\n", err.ErrorType))
	sb.WriteString(fmt.Sprintf("**Message**: %s\n", err.Message))

	if err.Tool != "" {
		sb.WriteString(fmt.Sprintf("**Detected by**: %s\n", err.Tool))
	}

	if err.Severity != "" {
		sb.WriteString(fmt.Sprintf("**Severity**: %s\n", err.Severity))
	}

	if codeSnippet != "" {
		sb.WriteString("\n**Relevant code**:\n```\n")
		sb.WriteString(codeSnippet)
		sb.WriteString("\n```\n")
	}

	if len(err.Suggestions) > 0 {
		sb.WriteString("\n**Analyzer suggestions**:\n")
		for _, suggestion := range err.Suggestions {
			sb.WriteString(fmt.Sprintf("- %s\n", suggestion))
		}
	}

	sb.WriteString("\nPlease fix this error. Provide the corrected code.\n")

	return sb.String()
}

// GenerateBatchFixPrompt generates a prompt for fixing multiple errors in batch
func (g *AIFeedbackGenerator) GenerateBatchFixPrompt(errors []*domain.ErrorDetails, fileContents map[string]string) string {
	var sb strings.Builder

	sb.WriteString("# Error Correction Request\n\n")
	sb.WriteString("Multiple errors were detected after your changes. Please fix all of them.\n\n")

	// Group errors by file
	errorsByFile := make(map[string][]*domain.ErrorDetails)
	for _, err := range errors {
		errorsByFile[err.SourceFile] = append(errorsByFile[err.SourceFile], err)
	}

	for file, fileErrors := range errorsByFile {
		sb.WriteString(fmt.Sprintf("## File: %s\n\n", file))

		sb.WriteString("### Errors:\n")
		for i, err := range fileErrors {
			sb.WriteString(fmt.Sprintf("%d. Line %d: [%s] %s\n", i+1, err.LineNumber, err.ErrorType, err.Message))
		}
		sb.WriteString("\n")

		if content, ok := fileContents[file]; ok && content != "" {
			sb.WriteString("### Current content:\n```\n")
			sb.WriteString(content)
			sb.WriteString("\n```\n\n")
		}
	}

	sb.WriteString("## Instructions\n")
	sb.WriteString("1. Fix all errors listed above\n")
	sb.WriteString("2. Provide the complete corrected code for each file\n")
	sb.WriteString("3. Ensure the fixes don't introduce new errors\n")

	return sb.String()
}

// SetTemplate allows setting a custom template for an error type
func (g *AIFeedbackGenerator) SetTemplate(errorType domain.ErrorType, template string) {
	g.templates[errorType] = template
}

// GetTemplate returns the template for an error type
func (g *AIFeedbackGenerator) GetTemplate(errorType domain.ErrorType) string {
	return g.templates[errorType]
}
