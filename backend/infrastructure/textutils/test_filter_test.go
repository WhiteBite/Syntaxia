package textutils

import (
	"testing"
)

func TestTestFilter_IsTestFile(t *testing.T) {
	f := NewTestFilter()

	tests := []struct {
		path     string
		expected bool
		desc     string
	}{
		// Go tests
		{"main_test.go", true, "Go test file"},
		{"service_test.go", true, "Go test file with underscore"},
		{"main.go", false, "Go source file"},
		{"testutils.go", false, "Go file with test in name but not test file"},

		// JavaScript/TypeScript tests
		{"app.test.js", true, "JS test file"},
		{"app.spec.js", true, "JS spec file"},
		{"app.test.ts", true, "TS test file"},
		{"app.spec.ts", true, "TS spec file"},
		{"app.test.tsx", true, "TSX test file"},
		{"app.spec.tsx", true, "TSX spec file"},
		{"app.js", false, "JS source file"},
		{"app.ts", false, "TS source file"},

		// Python tests
		{"test_main.py", true, "Python test file with prefix"},
		{"main_test.py", true, "Python test file with suffix"},
		{"main_tests.py", true, "Python tests file"},
		{"main.py", false, "Python source file"},
		{"testing.py", false, "Python file with test in name"},

		// Java tests
		{"MainTest.java", true, "Java test file"},
		{"TestMain.java", true, "Java test file with prefix"},
		{"Main.java", false, "Java source file"},

		// C# tests
		{"MainTests.cs", true, "C# test file"},
		{"Main.cs", false, "C# source file"},

		// Ruby tests
		{"main_spec.rb", true, "Ruby spec file"},
		{"main.rb", false, "Ruby source file"},

		// Storybook
		{"Button.stories.tsx", true, "Storybook file"},
		{"Button.stories.js", true, "Storybook JS file"},

		// Directory-based detection
		{"tests/main.go", true, "File in tests directory"},
		{"test/main.py", true, "File in test directory"},
		{"__tests__/app.js", true, "File in __tests__ directory"},
		{"src/__tests__/app.js", true, "File in nested __tests__ directory"},
		{"spec/models/user_spec.rb", true, "File in spec directory"},
		{"e2e/login.spec.ts", true, "File in e2e directory"},
		{"integration/api_test.go", true, "File in integration directory"},
		{"__mocks__/api.js", true, "File in __mocks__ directory"},
		{"fixtures/data.json", true, "File in fixtures directory"},

		// Edge cases
		{"", false, "Empty path"},
		{"src/main.go", false, "Source file in src"},
		{"lib/utils.ts", false, "Lib file"},
		{"testdata/input.txt", true, "File in testdata directory"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := f.IsTestFile(tt.path)
			if result != tt.expected {
				t.Errorf("IsTestFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestTestFilter_FilterTestFiles(t *testing.T) {
	f := NewTestFilter()

	input := []string{
		"main.go",
		"main_test.go",
		"service.go",
		"service_test.go",
		"utils.go",
		"tests/helper.go",
	}

	result := f.FilterTestFiles(input)

	expected := []string{
		"main.go",
		"service.go",
		"utils.go",
	}

	if len(result) != len(expected) {
		t.Errorf("FilterTestFiles() returned %d files, want %d", len(result), len(expected))
	}

	for i, path := range result {
		if path != expected[i] {
			t.Errorf("FilterTestFiles()[%d] = %q, want %q", i, path, expected[i])
		}
	}
}

func TestTestFilter_CountTestFiles(t *testing.T) {
	f := NewTestFilter()

	input := []string{
		"main.go",
		"main_test.go",
		"service.go",
		"service_test.go",
		"utils.go",
	}

	count := f.CountTestFiles(input)
	if count != 2 {
		t.Errorf("CountTestFiles() = %d, want 2", count)
	}
}

func TestTestFilter_EmptyInput(t *testing.T) {
	f := NewTestFilter()

	result := f.FilterTestFiles([]string{})
	if len(result) != 0 {
		t.Error("FilterTestFiles([]) should return empty slice")
	}

	result = f.FilterTestFiles(nil)
	if result != nil {
		t.Error("FilterTestFiles(nil) should return nil")
	}
}

func TestGlobalFunctions(t *testing.T) {
	// Test global IsTestFile
	if !IsTestFile("main_test.go") {
		t.Error("Global IsTestFile should detect test file")
	}
	if IsTestFile("main.go") {
		t.Error("Global IsTestFile should not detect source file as test")
	}

	// Test global FilterTestFiles
	result := FilterTestFiles([]string{"main.go", "main_test.go"})
	if len(result) != 1 || result[0] != "main.go" {
		t.Error("Global FilterTestFiles should filter test files")
	}
}

func BenchmarkIsTestFile_Simple(b *testing.B) {
	f := NewTestFilter()
	path := "main_test.go"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.IsTestFile(path)
	}
}

func BenchmarkIsTestFile_DeepPath(b *testing.B) {
	f := NewTestFilter()
	path := "src/components/features/auth/__tests__/login.spec.tsx"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.IsTestFile(path)
	}
}

func BenchmarkFilterTestFiles_100(b *testing.B) {
	f := NewTestFilter()
	paths := make([]string, 100)
	for i := 0; i < 100; i++ {
		if i%3 == 0 {
			paths[i] = "file_test.go"
		} else {
			paths[i] = "file.go"
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.FilterTestFiles(paths)
	}
}
