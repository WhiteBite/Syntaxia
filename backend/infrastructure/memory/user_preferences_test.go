package memory

import (
	"testing"
)

// mockContextMemory implements ContextMemoryInterface for testing
type mockContextMemory struct {
	prefs map[string]string
}

func newMockContextMemory() *mockContextMemory {
	return &mockContextMemory{prefs: make(map[string]string)}
}

func (m *mockContextMemory) SetPreference(key, value string) error {
	m.prefs[key] = value
	return nil
}

func (m *mockContextMemory) GetPreference(key string) (string, error) {
	return m.prefs[key], nil
}

func (m *mockContextMemory) GetAllPreferences() (map[string]string, error) {
	return m.prefs, nil
}

func TestNewUserPreferences(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)
	if up == nil {
		t.Fatal("NewUserPreferences returned nil")
	}
}

func TestUserPreferences_SetGet(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)

	err := up.Set(PrefExcludeTests, "true")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val := up.Get(PrefExcludeTests)
	if val != boolTrue {
		t.Errorf("expected 'true', got %q", val)
	}
}

func TestUserPreferences_GetBool(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)

	tests := []struct {
		value    string
		expected bool
	}{
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"yes", true},
		{"YES", true},
		{"1", true},
		{"false", false},
		{"no", false},
		{"0", false},
		{"", false},
	}

	for _, tt := range tests {
		up.Set(PrefExcludeTests, tt.value)
		got := up.GetBool(PrefExcludeTests)
		if got != tt.expected {
			t.Errorf("GetBool(%q) = %v, want %v", tt.value, got, tt.expected)
		}
	}
}

func TestUserPreferences_GetInt(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)

	tests := []struct {
		value      string
		defaultVal int
		expected   int
	}{
		{"10", 5, 10},
		{"100", 5, 100},
		{"", 5, 5},
		{"abc", 5, 5},
	}

	for _, tt := range tests {
		up.Set(PrefMaxContextFiles, tt.value)
		got := up.GetInt(PrefMaxContextFiles, tt.defaultVal)
		if got != tt.expected {
			t.Errorf("GetInt(%q, %d) = %d, want %d", tt.value, tt.defaultVal, got, tt.expected)
		}
	}
}

func TestUserPreferences_StringList(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)

	list := []string{"*.log", "*.tmp", "vendor/"}
	err := up.SetStringList(PrefExcludePatterns, list)
	if err != nil {
		t.Fatalf("SetStringList failed: %v", err)
	}

	got := up.GetStringList(PrefExcludePatterns)
	if len(got) != len(list) {
		t.Errorf("expected %d items, got %d", len(list), len(got))
	}

	for i, v := range list {
		if got[i] != v {
			t.Errorf("item %d: expected %q, got %q", i, v, got[i])
		}
	}
}

func TestUserPreferences_GetStringList_Empty(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)

	got := up.GetStringList(PrefExcludePatterns)
	if got != nil {
		t.Errorf("expected nil for empty list, got %v", got)
	}
}

func TestUserPreferences_LoadAll(t *testing.T) {
	mock := newMockContextMemory()
	mock.prefs["key1"] = "value1"
	mock.prefs["key2"] = "value2"

	up := NewUserPreferences(mock)
	err := up.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if up.cache["key1"] != "value1" {
		t.Error("key1 not loaded into cache")
	}
	if up.cache["key2"] != "value2" {
		t.Error("key2 not loaded into cache")
	}
}

func TestUserPreferences_ShouldExcludeFile_Tests(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)
	up.Set(PrefExcludeTests, "true")

	tests := []struct {
		path     string
		expected bool
	}{
		{"main_test.go", true},
		{"app.test.ts", true},
		{"component.spec.js", true},
		{"tests/unit.go", true},
		{"__tests__/app.js", true},
		{"main.go", false},
		{"service.ts", false},
	}

	for _, tt := range tests {
		got := up.ShouldExcludeFile(tt.path)
		if got != tt.expected {
			t.Errorf("ShouldExcludeFile(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestUserPreferences_ShouldExcludeFile_Vendor(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)
	up.Set(PrefExcludeVendor, "true")

	tests := []struct {
		path     string
		expected bool
	}{
		{"vendor/lib/file.go", true},
		{"node_modules/pkg/index.js", true},
		{"src/main.go", false},
	}

	for _, tt := range tests {
		got := up.ShouldExcludeFile(tt.path)
		if got != tt.expected {
			t.Errorf("ShouldExcludeFile(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestUserPreferences_ShouldExcludeFile_Generated(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)
	up.Set(PrefExcludeGenerated, "true")

	tests := []struct {
		path     string
		expected bool
	}{
		{"model.gen.go", true},
		{"types.generated.ts", true},
		{"proto.pb.go", true},
		{"main.go", false},
	}

	for _, tt := range tests {
		got := up.ShouldExcludeFile(tt.path)
		if got != tt.expected {
			t.Errorf("ShouldExcludeFile(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestUserPreferences_ShouldExcludeFile_CustomPatterns(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)
	up.SetStringList(PrefExcludePatterns, []string{"*.log", "temp/"})

	tests := []struct {
		path     string
		expected bool
	}{
		{"app.log", true},
		{"temp/cache.txt", true},
		{"main.go", false},
	}

	for _, tt := range tests {
		got := up.ShouldExcludeFile(tt.path)
		if got != tt.expected {
			t.Errorf("ShouldExcludeFile(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestUserPreferences_ShouldIncludeFile(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)

	// No patterns - include all
	if !up.ShouldIncludeFile("any/file.go") {
		t.Error("should include all when no patterns set")
	}

	// With patterns
	up.SetStringList(PrefIncludePatterns, []string{".go", ".ts"})

	tests := []struct {
		path     string
		expected bool
	}{
		{"main.go", true},
		{"app.ts", true},
		{"style.css", false},
	}

	for _, tt := range tests {
		got := up.ShouldIncludeFile(tt.path)
		if got != tt.expected {
			t.Errorf("ShouldIncludeFile(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestParsePreferenceFromMessage(t *testing.T) {
	tests := []struct {
		message   string
		expectKey PreferenceKey
		expectVal string
		found     bool
	}{
		{"не включай тесты в контекст", PrefExcludeTests, "true", true},
		{"exclude tests from context", PrefExcludeTests, "true", true},
		{"без тестов пожалуйста", PrefExcludeTests, "true", true},
		{"ignore tests", PrefExcludeTests, "true", true},
		{"exclude vendor files", PrefExcludeVendor, "true", true},
		{"random message", "", "", false},
	}

	for _, tt := range tests {
		key, val, found := ParsePreferenceFromMessage(tt.message)
		if found != tt.found {
			t.Errorf("ParsePreferenceFromMessage(%q): found = %v, want %v", tt.message, found, tt.found)
		}
		if found && key != tt.expectKey {
			t.Errorf("ParsePreferenceFromMessage(%q): key = %v, want %v", tt.message, key, tt.expectKey)
		}
		if found && val != tt.expectVal {
			t.Errorf("ParsePreferenceFromMessage(%q): val = %q, want %q", tt.message, val, tt.expectVal)
		}
	}
}

func TestUserPreferences_Cache(t *testing.T) {
	mock := newMockContextMemory()
	up := NewUserPreferences(mock)

	// Set value
	up.Set(PrefExcludeTests, "true")

	// Should be in cache
	if up.cache[string(PrefExcludeTests)] != "true" {
		t.Error("value not cached after Set")
	}

	// Get should use cache
	val := up.Get(PrefExcludeTests)
	if val != "true" {
		t.Errorf("Get returned %q, want 'true'", val)
	}
}
