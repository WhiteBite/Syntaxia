package analysis

import (
	"syntaxia/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuggestFiles_EmptyTask(t *testing.T) {
	service := NewService()
	files := []*domain.FileNode{
		{Name: "main.go", Path: "/project/main.go", IsDir: false},
	}

	result := service.SuggestFiles("", files)

	assert.Empty(t, result, "Should return empty for empty task")
}

func TestSuggestFiles_NoFiles(t *testing.T) {
	service := NewService()

	result := service.SuggestFiles("fix auth bug", []*domain.FileNode{})

	assert.Empty(t, result, "Should return empty when no files provided")
}

func TestSuggestFiles_MatchingFiles(t *testing.T) {
	service := NewService()
	files := []*domain.FileNode{
		{Name: "auth.go", Path: "/project/auth.go", IsDir: false},
		{Name: "auth_test.go", Path: "/project/auth_test.go", IsDir: false},
		{Name: "user.go", Path: "/project/user.go", IsDir: false},
		{Name: "config.go", Path: "/project/config.go", IsDir: false},
	}

	result := service.SuggestFiles("fix authentication bug in auth module", files)

	assert.NotEmpty(t, result, "Should return suggestions")

	found := false
	for _, suggestion := range result {
		if suggestion.Path == "/project/auth.go" || suggestion.Path == "/project/auth_test.go" {
			found = true
			assert.Greater(t, suggestion.Confidence, 0.0, "Confidence should be greater than 0")
			assert.NotEmpty(t, suggestion.Reason, "Reason should not be empty")
		}
	}
	assert.True(t, found, "Should suggest auth-related files")
}

func TestSuggestFiles_IgnoresDirectories(t *testing.T) {
	service := NewService()
	files := []*domain.FileNode{
		{Name: "auth", Path: "/project/auth", IsDir: true},
		{Name: "auth.go", Path: "/project/auth.go", IsDir: false},
	}

	result := service.SuggestFiles("auth", files)

	assert.Len(t, result, 1, "Should only suggest files, not directories")
	assert.Equal(t, "/project/auth.go", result[0].Path)
}

func TestSuggestFiles_WithNestedFiles(t *testing.T) {
	service := NewService()
	files := []*domain.FileNode{
		{
			Name:  "src",
			Path:  "/project/src",
			IsDir: true,
			Children: []*domain.FileNode{
				{Name: "auth.go", Path: "/project/src/auth.go", IsDir: false},
				{Name: "user.go", Path: "/project/src/user.go", IsDir: false},
			},
		},
	}

	result := service.SuggestFiles("modify auth", files)

	assert.NotEmpty(t, result, "Should find files in nested structure")

	found := false
	for _, suggestion := range result {
		if suggestion.Path == "/project/src/auth.go" {
			found = true
		}
	}
	assert.True(t, found, "Should suggest nested auth file")
}

func TestExtractTaskKeywords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Simple text",
			input:    "fix authentication bug",
			expected: 3,
		},
		{
			name:     "With stop words",
			input:    "fix the authentication bug in the system",
			expected: 4,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "Only stop words",
			input:    "the and or but",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTaskKeywords(tt.input)
			assert.Len(t, result, tt.expected, "Keyword count should match")
		})
	}
}

func TestScoreFileByKeywords(t *testing.T) {
	tests := []struct {
		name     string
		file     domain.FileNode
		keywords []string
		minScore float64
	}{
		{
			name:     "Exact filename match",
			file:     domain.FileNode{Name: "auth.go", Path: "/project/auth.go"},
			keywords: []string{"auth"},
			minScore: 0.5,
		},
		{
			name:     "Path contains keyword",
			file:     domain.FileNode{Name: "service.go", Path: "/project/auth/service.go"},
			keywords: []string{"auth"},
			minScore: 0.1,
		},
		{
			name:     "No match",
			file:     domain.FileNode{Name: "config.go", Path: "/project/config.go"},
			keywords: []string{"auth"},
			minScore: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := scoreFileByKeywords(tt.file, tt.keywords)
			if tt.minScore > 0 {
				assert.GreaterOrEqual(t, score, tt.minScore, "Score should meet minimum")
			} else {
				assert.Equal(t, 0.0, score, "Score should be zero for no match")
			}
		})
	}
}

func TestSortByConfidence(t *testing.T) {
	suggestions := []domain.SuggestedFile{
		{Path: "a.go", Confidence: 0.3},
		{Path: "b.go", Confidence: 0.9},
		{Path: "c.go", Confidence: 0.5},
	}

	sortByConfidence(suggestions)

	assert.Equal(t, 0.9, suggestions[0].Confidence, "First should have highest confidence")
	assert.Equal(t, 0.5, suggestions[1].Confidence, "Second should have medium confidence")
	assert.Equal(t, 0.3, suggestions[2].Confidence, "Third should have lowest confidence")
}
