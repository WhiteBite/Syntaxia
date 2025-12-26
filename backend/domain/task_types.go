// Package domain contains business entities and interfaces.
package domain

// TaskType represents the type of a task in the system.
// Used for guardrails, ephemeral mode, and task classification.
type TaskType string

// Task type constants for task classification and guardrails
const (
	// Core task types for ephemeral mode
	TaskTypeScaffold TaskType = "scaffold"
	TaskTypeDepsFix  TaskType = "deps_fix"

	// Task types for context analysis
	TaskTypeFeature       TaskType = "feature"
	TaskTypeBugFix        TaskType = "bug_fix"
	TaskTypeTest          TaskType = "test"
	TaskTypeRefactor      TaskType = "refactor"
	TaskTypeDocumentation TaskType = "documentation"

	// Default task type
	TaskTypeRegular TaskType = "regular"
)

// String returns the string representation of TaskType
func (t TaskType) String() string {
	return string(t)
}

// IsEphemeralAllowed returns true if this task type allows ephemeral mode
func (t TaskType) IsEphemeralAllowed() bool {
	return t == TaskTypeScaffold || t == TaskTypeDepsFix
}

// ParseTaskType parses a string into TaskType, returns TaskTypeRegular if unknown
func ParseTaskType(s string) TaskType {
	switch TaskType(s) {
	case TaskTypeScaffold, TaskTypeDepsFix, TaskTypeFeature,
		TaskTypeBugFix, TaskTypeTest, TaskTypeRefactor, TaskTypeDocumentation:
		return TaskType(s)
	default:
		return TaskTypeRegular
	}
}

// TaskTypeFromID determines task type from task ID by checking for known patterns
func TaskTypeFromID(taskID string) TaskType {
	patterns := map[string]TaskType{
		"scaffold": TaskTypeScaffold,
		"deps_fix": TaskTypeDepsFix,
		"feature":  TaskTypeFeature,
		"bug_fix":  TaskTypeBugFix,
		"test":     TaskTypeTest,
		"refactor": TaskTypeRefactor,
		"doc":      TaskTypeDocumentation,
	}

	for pattern, taskType := range patterns {
		if containsPattern(taskID, pattern) {
			return taskType
		}
	}
	return TaskTypeRegular
}

// containsPattern checks if s contains pattern (case-insensitive)
func containsPattern(s, pattern string) bool {
	// Simple contains check - avoid importing strings in domain
	for i := 0; i <= len(s)-len(pattern); i++ {
		match := true
		for j := 0; j < len(pattern); j++ {
			sc := s[i+j]
			pc := pattern[j]
			// Simple lowercase comparison for ASCII
			if sc >= 'A' && sc <= 'Z' {
				sc += 32
			}
			if pc >= 'A' && pc <= 'Z' {
				pc += 32
			}
			if sc != pc {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
