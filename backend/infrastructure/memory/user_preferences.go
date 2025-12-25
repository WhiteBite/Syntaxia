package memory

import (
	"encoding/json"
	"strings"
)

// PreferenceKey defines known preference keys
type PreferenceKey string

const (
	PrefExcludeTests      PreferenceKey = "exclude_tests"
	PrefExcludeVendor     PreferenceKey = "exclude_vendor"
	PrefExcludeGenerated  PreferenceKey = "exclude_generated"
	PrefPreferredLanguage PreferenceKey = "preferred_language"
	PrefMaxContextFiles   PreferenceKey = "max_context_files"
	PrefIncludeComments   PreferenceKey = "include_comments"
	PrefCodeStyle         PreferenceKey = "code_style"
	PrefTestFramework     PreferenceKey = "test_framework"
	PrefExcludePatterns   PreferenceKey = "exclude_patterns"
	PrefIncludePatterns   PreferenceKey = "include_patterns"
)

// UserPreferences manages user preferences with parsing
type UserPreferences struct {
	memory ContextMemoryInterface
	cache  map[string]string
}

// ContextMemoryInterface defines the interface for context memory
type ContextMemoryInterface interface {
	SetPreference(key, value string) error
	GetPreference(key string) (string, error)
	GetAllPreferences() (map[string]string, error)
}

// NewUserPreferences creates a new user preferences manager
func NewUserPreferences(memory ContextMemoryInterface) *UserPreferences {
	return &UserPreferences{
		memory: memory,
		cache:  make(map[string]string),
	}
}

// Set sets a preference
func (up *UserPreferences) Set(key PreferenceKey, value string) error {
	up.cache[string(key)] = value
	return up.memory.SetPreference(string(key), value)
}

// Get gets a preference
func (up *UserPreferences) Get(key PreferenceKey) string {
	if v, ok := up.cache[string(key)]; ok {
		return v
	}
	v, _ := up.memory.GetPreference(string(key))
	up.cache[string(key)] = v
	return v
}

const boolTrue = "true"

// GetBool gets a boolean preference
func (up *UserPreferences) GetBool(key PreferenceKey) bool {
	v := strings.ToLower(up.Get(key))
	return v == boolTrue || v == "yes" || v == "1"
}

// GetInt gets an integer preference
func (up *UserPreferences) GetInt(key PreferenceKey, defaultVal int) int {
	v := up.Get(key)
	if v == "" {
		return defaultVal
	}
	var i int
	for _, c := range v {
		if c >= '0' && c <= '9' {
			i = i*10 + int(c-'0')
		}
	}
	if i == 0 {
		return defaultVal
	}
	return i
}

// GetStringList gets a list preference
func (up *UserPreferences) GetStringList(key PreferenceKey) []string {
	v := up.Get(key)
	if v == "" {
		return nil
	}
	var list []string
	_ = json.Unmarshal([]byte(v), &list)
	return list
}

// SetStringList sets a list preference
func (up *UserPreferences) SetStringList(key PreferenceKey, list []string) error {
	b, _ := json.Marshal(list)
	return up.Set(key, string(b))
}

// LoadAll loads all preferences into cache
func (up *UserPreferences) LoadAll() error {
	prefs, err := up.memory.GetAllPreferences()
	if err != nil {
		return err
	}
	for k, v := range prefs {
		up.cache[k] = v
	}
	return nil
}

// isTestFile checks if path looks like a test file
func isTestFile(path string) bool {
	testPatterns := []string{"_test.", ".test.", ".spec.", "/test/", "/tests/", "/__tests__/"}
	testPrefixes := []string{"test/", "tests/", "__tests__/"}
	for _, p := range testPatterns {
		if strings.Contains(path, p) {
			return true
		}
	}
	for _, p := range testPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// isVendorFile checks if path is in vendor/node_modules
func isVendorFile(path string) bool {
	return strings.Contains(path, "/vendor/") || strings.Contains(path, "/node_modules/") ||
		strings.HasPrefix(path, "vendor/") || strings.HasPrefix(path, "node_modules/")
}

// isGeneratedFile checks if path looks like generated code
func isGeneratedFile(path string) bool {
	return strings.Contains(path, ".gen.") || strings.Contains(path, ".generated.") ||
		strings.Contains(path, "_gen.") || strings.HasSuffix(path, ".pb.go")
}

// matchesPattern checks if path matches a glob-like pattern
func matchesPattern(path, pattern string) bool {
	pattern = strings.ToLower(pattern)
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(path, strings.TrimPrefix(pattern, "*"))
	}
	return strings.Contains(path, pattern)
}

// ShouldExcludeFile checks if a file should be excluded based on preferences
func (up *UserPreferences) ShouldExcludeFile(filePath string) bool {
	filePath = strings.ToLower(filePath)

	if up.GetBool(PrefExcludeTests) && isTestFile(filePath) {
		return true
	}
	if up.GetBool(PrefExcludeVendor) && isVendorFile(filePath) {
		return true
	}
	if up.GetBool(PrefExcludeGenerated) && isGeneratedFile(filePath) {
		return true
	}

	for _, pattern := range up.GetStringList(PrefExcludePatterns) {
		if matchesPattern(filePath, pattern) {
			return true
		}
	}
	return false
}

// ShouldIncludeFile checks if a file matches include patterns
func (up *UserPreferences) ShouldIncludeFile(filePath string) bool {
	includePatterns := up.GetStringList(PrefIncludePatterns)
	if len(includePatterns) == 0 {
		return true // No filter, include all
	}

	filePath = strings.ToLower(filePath)
	for _, pattern := range includePatterns {
		if strings.Contains(filePath, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// ParsePreferenceFromMessage extracts preferences from user message
func ParsePreferenceFromMessage(message string) (key PreferenceKey, value string, found bool) {
	message = strings.ToLower(message)

	// Patterns for preference extraction
	patterns := map[string]PreferenceKey{
		"не включай тесты":     PrefExcludeTests,
		"exclude tests":        PrefExcludeTests,
		"без тестов":           PrefExcludeTests,
		"ignore tests":         PrefExcludeTests,
		"не включай vendor":    PrefExcludeVendor,
		"exclude vendor":       PrefExcludeVendor,
		"без сгенерированного": PrefExcludeGenerated,
		"exclude generated":    PrefExcludeGenerated,
		"максимум файлов":      PrefMaxContextFiles,
		"max files":            PrefMaxContextFiles,
		"включай комментарии":  PrefIncludeComments,
		"include comments":     PrefIncludeComments,
	}

	for pattern, prefKey := range patterns {
		if strings.Contains(message, pattern) {
			// Determine value
			if strings.Contains(message, "не ") || strings.Contains(message, "exclude") ||
				strings.Contains(message, "без") || strings.Contains(message, "ignore") {
				return prefKey, boolTrue, true
			}
			if strings.Contains(message, "включай") || strings.Contains(message, "include") {
				return prefKey, boolTrue, true
			}
		}
	}

	return "", "", false
}
