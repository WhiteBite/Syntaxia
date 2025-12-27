package verification

import (
	"fmt"
	"sync"
)

// ErrorFilter filters verification errors to show only new ones
type ErrorFilter struct {
	mu             sync.RWMutex
	baselineErrors map[string][]VerificationIssue // file -> errors before changes
	baselineKeys   map[string]bool                // fingerprint keys for quick lookup
}

// NewErrorFilter creates a new error filter
func NewErrorFilter() *ErrorFilter {
	return &ErrorFilter{
		baselineErrors: make(map[string][]VerificationIssue),
		baselineKeys:   make(map[string]bool),
	}
}

// SetBaseline sets the baseline errors from before changes
func (f *ErrorFilter) SetBaseline(errors []VerificationIssue) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.baselineErrors = make(map[string][]VerificationIssue)
	f.baselineKeys = make(map[string]bool)

	for _, err := range errors {
		f.baselineErrors[err.File] = append(f.baselineErrors[err.File], err)
		f.baselineKeys[f.issueKey(err)] = true
	}
}

// SetBaselineFromState sets baseline from BaselineState
func (f *ErrorFilter) SetBaselineFromState(state *BaselineState) {
	if state == nil {
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.baselineErrors = make(map[string][]VerificationIssue)
	f.baselineKeys = make(map[string]bool)

	// Convert type errors
	for _, fp := range state.TypeErrors {
		issue := VerificationIssue{
			Type:     "type",
			File:     fp.File,
			Line:     fp.Line,
			Severity: fp.Severity,
			Message:  fp.Message,
			Code:     fp.Code,
		}
		f.baselineErrors[fp.File] = append(f.baselineErrors[fp.File], issue)
		f.baselineKeys[f.issueKey(issue)] = true
	}

	// Convert lint issues
	for _, fp := range state.LintIssues {
		issue := VerificationIssue{
			Type:     "lint",
			File:     fp.File,
			Line:     fp.Line,
			Severity: fp.Severity,
			Message:  fp.Message,
			Code:     fp.Code,
		}
		f.baselineErrors[fp.File] = append(f.baselineErrors[fp.File], issue)
		f.baselineKeys[f.issueKey(issue)] = true
	}
}

// FilterNewErrors returns only errors that are new (not in baseline)
func (f *ErrorFilter) FilterNewErrors(currentErrors []VerificationIssue) []VerificationIssue {
	f.mu.RLock()
	defer f.mu.RUnlock()

	newErrors := make([]VerificationIssue, 0)
	for _, err := range currentErrors {
		if !f.baselineKeys[f.issueKey(err)] {
			newErrors = append(newErrors, err)
		}
	}
	return newErrors
}

// IsNewError checks if a single error is new
func (f *ErrorFilter) IsNewError(err VerificationIssue) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return !f.baselineKeys[f.issueKey(err)]
}

// GetFixedErrors returns errors that were in baseline but not in current
func (f *ErrorFilter) GetFixedErrors(currentErrors []VerificationIssue) []VerificationIssue {
	f.mu.RLock()
	defer f.mu.RUnlock()

	currentKeys := make(map[string]bool)
	for _, err := range currentErrors {
		currentKeys[f.issueKey(err)] = true
	}

	fixedErrors := make([]VerificationIssue, 0)
	for _, errors := range f.baselineErrors {
		for _, err := range errors {
			if !currentKeys[f.issueKey(err)] {
				fixedErrors = append(fixedErrors, err)
			}
		}
	}
	return fixedErrors
}

// GetBaselineCount returns the number of baseline errors
func (f *ErrorFilter) GetBaselineCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.baselineKeys)
}

// GetBaselineForFile returns baseline errors for a specific file
func (f *ErrorFilter) GetBaselineForFile(file string) []VerificationIssue {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if errors, ok := f.baselineErrors[file]; ok {
		result := make([]VerificationIssue, len(errors))
		copy(result, errors)
		return result
	}
	return nil
}

// Clear clears the baseline
func (f *ErrorFilter) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.baselineErrors = make(map[string][]VerificationIssue)
	f.baselineKeys = make(map[string]bool)
}

// issueKey creates a unique key for an issue
func (f *ErrorFilter) issueKey(issue VerificationIssue) string {
	return fmt.Sprintf("%s:%d:%s:%s", issue.File, issue.Line, issue.Code, issue.Message)
}
