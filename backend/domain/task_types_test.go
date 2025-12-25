package domain

import "testing"

func TestTaskType_String(t *testing.T) {
	tests := []struct {
		taskType TaskType
		expected string
	}{
		{TaskTypeScaffold, "scaffold"},
		{TaskTypeDepsFix, "deps_fix"},
		{TaskTypeFeature, "feature"},
		{TaskTypeBugFix, "bug_fix"},
		{TaskTypeRegular, "regular"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.taskType.String(); got != tt.expected {
				t.Errorf("TaskType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTaskType_IsEphemeralAllowed(t *testing.T) {
	tests := []struct {
		taskType TaskType
		expected bool
	}{
		{TaskTypeScaffold, true},
		{TaskTypeDepsFix, true},
		{TaskTypeFeature, false},
		{TaskTypeBugFix, false},
		{TaskTypeRegular, false},
	}

	for _, tt := range tests {
		t.Run(tt.taskType.String(), func(t *testing.T) {
			if got := tt.taskType.IsEphemeralAllowed(); got != tt.expected {
				t.Errorf("TaskType.IsEphemeralAllowed() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseTaskType(t *testing.T) {
	tests := []struct {
		input    string
		expected TaskType
	}{
		{"scaffold", TaskTypeScaffold},
		{"deps_fix", TaskTypeDepsFix},
		{"feature", TaskTypeFeature},
		{"bug_fix", TaskTypeBugFix},
		{"unknown", TaskTypeRegular},
		{"", TaskTypeRegular},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParseTaskType(tt.input); got != tt.expected {
				t.Errorf("ParseTaskType(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTaskTypeFromID(t *testing.T) {
	tests := []struct {
		taskID   string
		expected TaskType
	}{
		{"task-scaffold-001", TaskTypeScaffold},
		{"deps_fix_module", TaskTypeDepsFix},
		{"feature-login", TaskTypeFeature},
		{"bug_fix_crash", TaskTypeBugFix},
		{"test-unit-001", TaskTypeTest},
		{"refactor-api", TaskTypeRefactor},
		{"doc-readme", TaskTypeDocumentation},
		{"random-task", TaskTypeRegular},
		{"SCAFFOLD-UPPER", TaskTypeScaffold}, // case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.taskID, func(t *testing.T) {
			if got := TaskTypeFromID(tt.taskID); got != tt.expected {
				t.Errorf("TaskTypeFromID(%q) = %v, want %v", tt.taskID, got, tt.expected)
			}
		})
	}
}
