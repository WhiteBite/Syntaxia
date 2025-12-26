package testengine

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"syntaxia/domain"
)

// TestDartTestRunner_GetLanguage tests language identifier
func TestDartTestRunner_GetLanguage(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	assert.Equal(t, "dart", runner.GetLanguage())
}

// TestDartTestRunner_isTestFile tests Dart test file detection
func TestDartTestRunner_isTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tests := []struct {
		name     string
		fileName string
		expected bool
	}{
		{
			name:     "valid_test_file",
			fileName: "user_service_test.dart",
			expected: true,
		},
		{
			name:     "widget_test",
			fileName: "widget_test.dart",
			expected: true,
		},
		{
			name:     "integration_test",
			fileName: "app_integration_test.dart",
			expected: true,
		},
		{
			name:     "regular_dart_file",
			fileName: "user_service.dart",
			expected: false,
		},
		{
			name:     "test_in_middle",
			fileName: "test_utils.dart",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.isTestFile(tt.fileName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDartTestRunner_shouldSkipDirectory tests directory skip logic
func TestDartTestRunner_shouldSkipDirectory(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{
			name:     "build_directory",
			dirName:  "build",
			expected: true,
		},
		{
			name:     "dart_tool_directory",
			dirName:  ".dart_tool",
			expected: true,
		},
		{
			name:     "packages_directory",
			dirName:  ".packages",
			expected: true,
		},
		{
			name:     "git_directory",
			dirName:  ".git",
			expected: true,
		},
		{
			name:     "ios_directory",
			dirName:  "ios",
			expected: true,
		},
		{
			name:     "android_directory",
			dirName:  "android",
			expected: true,
		},
		{
			name:     "web_directory",
			dirName:  "web",
			expected: true,
		},
		{
			name:     "lib_directory",
			dirName:  "lib",
			expected: false,
		},
		{
			name:     "test_directory",
			dirName:  "test",
			expected: false,
		},
		{
			name:     "src_directory",
			dirName:  "src",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.shouldSkipDirectory(tt.dirName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDartTestRunner_isFlutterProject tests Flutter project detection
func TestDartTestRunner_isFlutterProject(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tests := []struct {
		name           string
		pubspecContent string
		expected       bool
	}{
		{
			name: "flutter_project",
			pubspecContent: `name: my_app
dependencies:
  flutter:
    sdk: flutter
`,
			expected: true,
		},
		{
			name: "pure_dart_project",
			pubspecContent: `name: my_app
dependencies:
  http: ^1.0.0
`,
			expected: false,
		},
		{
			name: "flutter_in_dependencies",
			pubspecContent: `name: my_app
dependencies:
  flutter:
    sdk: flutter
dev_dependencies:
  test: ^1.0.0
`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "dart-flutter-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			err = os.WriteFile(filepath.Join(tempDir, "pubspec.yaml"), []byte(tt.pubspecContent), 0o644)
			require.NoError(t, err)

			result := runner.isFlutterProject(tempDir)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDartTestRunner_isFlutterProject_NoPubspec tests when pubspec.yaml doesn't exist
func TestDartTestRunner_isFlutterProject_NoPubspec(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tempDir, err := os.MkdirTemp("", "dart-no-pubspec-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	result := runner.isFlutterProject(tempDir)
	assert.False(t, result)
}

// TestDartTestRunner_analyzeTestFile tests test file type analysis
func TestDartTestRunner_analyzeTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tempDir, err := os.MkdirTemp("", "dart-analyze-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected string
	}{
		{
			name:     "widget_test_tester",
			fileName: "widget_test.dart",
			content: `import 'package:flutter_test/flutter_test.dart';
void main() {
  testWidgets('MyWidget test', (WidgetTester tester) async {
    await tester.pumpWidget(MyWidget());
  });
}`,
			expected: "widget",
		},
		{
			name:     "widget_test_pumpwidget",
			fileName: "button_test.dart",
			content: `void main() {
  test('button test', () async {
    await pumpWidget(MyButton());
  });
}`,
			expected: "widget",
		},
		{
			name:     "integration_test_keyword",
			fileName: "app_test.dart",
			content: `// This is an integration test
void main() {
  test('integration test', () {});
}`,
			expected: "integration",
		},
		{
			name:     "smoke_test",
			fileName: "smoke_test.dart",
			content: `void main() {
  test('smoke test', () {});
}`,
			expected: "smoke",
		},
		{
			name:     "unit_test",
			fileName: "service_test.dart",
			content: `void main() {
  test('service test', () {
    expect(1 + 1, equals(2));
  });
}`,
			expected: "unit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := runner.analyzeTestFile(testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDartTestRunner_analyzeTestFile_NonExistent tests analyzing non-existent file
func TestDartTestRunner_analyzeTestFile_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	result := runner.analyzeTestFile("/nonexistent/test.dart")
	assert.Equal(t, "unit", result)
}

// TestDartTestRunner_parseTestOutput tests parsing dart test output
func TestDartTestRunner_parseTestOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tests := []struct {
		name           string
		output         string
		testCount      int
		expectAllPass  bool
		expectSomeFail bool
	}{
		{
			name:          "all_passed",
			output:        "All tests passed!",
			testCount:     2,
			expectAllPass: true,
		},
		{
			name:           "some_failed",
			output:         "Some tests failed. user_service FAILED",
			testCount:      2,
			expectSomeFail: true,
		},
		{
			name:          "no_match",
			output:        "Running tests...",
			testCount:     1,
			expectAllPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testInfos := make([]*domain.TestInfo, tt.testCount)
			for i := 0; i < tt.testCount; i++ {
				testInfos[i] = &domain.TestInfo{
					Path: "test_" + string(rune('a'+i)) + "_test.dart",
					Name: "test_" + string(rune('a'+i)) + "_test.dart",
				}
			}

			results := runner.parseTestOutput(tt.output, testInfos, 1.0)

			assert.Len(t, results, tt.testCount)

			if tt.expectAllPass {
				for _, r := range results {
					assert.True(t, r.Success)
				}
			}
		})
	}
}

// TestDartTestRunner_DiscoverTests tests test discovery
func TestDartTestRunner_DiscoverTests(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dart-discover-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create pubspec.yaml
	pubspec := `name: test_project
version: 1.0.0
`
	err = os.WriteFile(filepath.Join(tempDir, "pubspec.yaml"), []byte(pubspec), 0o644)
	require.NoError(t, err)

	// Create test directory
	testDir := filepath.Join(tempDir, "test")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	// Create test files
	testFiles := map[string]string{
		"user_service_test.dart": `void main() { test('test', () {}); }`,
		"widget_test.dart":       `void main() { testWidgets('test', (tester) {}); }`,
		"helper.dart":            `void helper() {}`,
	}

	for name, content := range testFiles {
		err = os.WriteFile(filepath.Join(testDir, name), []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 2)

	foundFiles := make(map[string]bool)
	for _, test := range tests {
		foundFiles[test.Name] = true
	}

	assert.True(t, foundFiles["user_service_test.dart"])
	assert.True(t, foundFiles["widget_test.dart"])
	assert.False(t, foundFiles["helper.dart"])
}

// TestDartTestRunner_DiscoverTests_NoTestDirectory tests discovery when test/ doesn't exist
func TestDartTestRunner_DiscoverTests_NoTestDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dart-no-test-dir-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestDartTestRunner_DiscoverTests_NestedTests tests discovery in nested directories
func TestDartTestRunner_DiscoverTests_NestedTests(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dart-nested-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create nested test directories
	nestedDir := filepath.Join(tempDir, "test", "unit", "services")
	err = os.MkdirAll(nestedDir, 0o755)
	require.NoError(t, err)

	// Create test file in nested directory
	testFile := filepath.Join(nestedDir, "user_service_test.dart")
	err = os.WriteFile(testFile, []byte(`void main() { test('test', () {}); }`), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	assert.Len(t, tests, 1)
	assert.Contains(t, tests[0].Path, "user_service_test.dart")
}

// TestDartTestRunner_RunTest tests single test execution
func TestDartTestRunner_RunTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tempDir, err := os.MkdirTemp("", "dart-run-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	pubspec := `name: test_project`
	err = os.WriteFile(filepath.Join(tempDir, "pubspec.yaml"), []byte(pubspec), 0o644)
	require.NoError(t, err)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		Verbose:     true,
		EnvVars: map[string]string{
			"DART_VM_OPTIONS": "--enable-asserts",
		},
	}

	result, err := runner.RunTest(context.Background(), "test/user_service_test.dart", config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test/user_service_test.dart", result.TestPath)
	assert.Equal(t, "dart", result.Language)
}

// TestDartTestRunner_RunTestSuite tests running a test suite
func TestDartTestRunner_RunTestSuite(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tempDir, err := os.MkdirTemp("", "dart-suite-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	pubspec := `name: test_project`
	err = os.WriteFile(filepath.Join(tempDir, "pubspec.yaml"), []byte(pubspec), 0o644)
	require.NoError(t, err)

	suite := &domain.TestSuite{
		Name:        "test_suite",
		Language:    "dart",
		ProjectPath: tempDir,
		Tests: []*domain.TestInfo{
			{Path: "test/test1_test.dart", Name: "test1_test.dart", Type: "unit"},
			{Path: "test/test2_test.dart", Name: "test2_test.dart", Type: "unit"},
		},
		Config: &domain.TestConfig{
			ProjectPath: tempDir,
			Verbose:     false,
		},
	}

	results, err := runner.RunTestSuite(context.Background(), suite)

	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results, 2)
}

// TestDartTestRunner_parseTestOutput_EdgeCases tests edge cases in output parsing
func TestDartTestRunner_parseTestOutput_EdgeCases(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tests := []struct {
		name          string
		output        string
		testCount     int
		expectSuccess bool
	}{
		{
			name:          "explicit_failed_test",
			output:        "user_service FAILED\nSome tests failed.",
			testCount:     1,
			expectSuccess: false,
		},
		{
			name:          "no_tests_message",
			output:        "No tests found.",
			testCount:     1,
			expectSuccess: true,
		},
		{
			name:          "verbose_output_passed",
			output:        "00:00 +1: All tests passed!\nFinished in 0.5s",
			testCount:     1,
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testInfos := make([]*domain.TestInfo, tt.testCount)
			for i := 0; i < tt.testCount; i++ {
				testInfos[i] = &domain.TestInfo{
					Path: "test_" + string(rune('a'+i)) + "_test.dart",
					Name: "test_" + string(rune('a'+i)) + "_test.dart",
				}
			}

			results := runner.parseTestOutput(tt.output, testInfos, 1.0)

			assert.Len(t, results, tt.testCount)
			if tt.testCount > 0 {
				assert.Equal(t, tt.expectSuccess, results[0].Success)
			}
		})
	}
}

// TestDartTestRunner_analyzeTestFile_IntegrationTestBinding tests integration test binding detection
func TestDartTestRunner_analyzeTestFile_IntegrationTestBinding(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tempDir, err := os.MkdirTemp("", "dart-integration-binding-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// The function converts to lowercase, so "integration" keyword is detected
	// Note: widget tests are checked first, so we avoid widget keywords
	content := `// This is an integration test file
void main() {
  test('app integration test', () async {
    // test code without widget keywords
  });
}`
	testFile := filepath.Join(tempDir, "app_test.dart")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	result := runner.analyzeTestFile(testFile)
	assert.Equal(t, "integration", result)
}

// TestDartTestRunner_isFlutterProject_FlutterTest tests Flutter test dependency detection
func TestDartTestRunner_isFlutterProject_FlutterTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tempDir, err := os.MkdirTemp("", "dart-flutter-test-dep-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// flutter_test still contains "flutter:" so it should be detected
	pubspec := `name: my_app
dependencies:
  flutter:
    sdk: flutter
dev_dependencies:
  flutter_test:
    sdk: flutter
`
	err = os.WriteFile(filepath.Join(tempDir, "pubspec.yaml"), []byte(pubspec), 0o644)
	require.NoError(t, err)

	result := runner.isFlutterProject(tempDir)
	assert.True(t, result)
}

// TestDartTestRunner_DiscoverTests_SkipPlatformDirectories tests that platform directories are skipped
func TestDartTestRunner_DiscoverTests_SkipPlatformDirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dart-skip-platform-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test directory
	testDir := filepath.Join(tempDir, "test")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	// Create valid test file
	validTest := filepath.Join(testDir, "valid_test.dart")
	err = os.WriteFile(validTest, []byte(`void main() { test('test', () {}); }`), 0o644)
	require.NoError(t, err)

	// Create platform directories that should be skipped
	platformDirs := []string{"ios", "android", "web", "linux", "macos", "windows"}
	for _, dir := range platformDirs {
		fullDir := filepath.Join(testDir, dir)
		err = os.MkdirAll(fullDir, 0o755)
		require.NoError(t, err)

		testFile := filepath.Join(fullDir, "should_skip_test.dart")
		err = os.WriteFile(testFile, []byte(`void main() {}`), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Should only find valid_test.dart
	assert.Len(t, tests, 1)
	assert.Equal(t, "valid_test.dart", tests[0].Name)
}

// TestDartTestRunner_RunTest_WithMetadata tests that metadata is set correctly
func TestDartTestRunner_RunTest_WithMetadata(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tempDir, err := os.MkdirTemp("", "dart-metadata-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create Flutter project
	pubspec := `name: my_app
dependencies:
  flutter:
    sdk: flutter
`
	err = os.WriteFile(filepath.Join(tempDir, "pubspec.yaml"), []byte(pubspec), 0o644)
	require.NoError(t, err)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		Verbose:     false,
	}

	result, err := runner.RunTest(context.Background(), "test/widget_test.dart", config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Metadata)
	assert.Equal(t, true, result.Metadata["isFlutter"])
}

// TestDartTestRunner_shouldSkipDirectory_AllCases tests all skip directory cases
func TestDartTestRunner_shouldSkipDirectory_AllCases(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	skipDirs := []string{
		"build", ".dart_tool", ".packages", ".git",
		"ios", "android", "web", "linux", "macos", "windows",
	}

	for _, dir := range skipDirs {
		t.Run(dir, func(t *testing.T) {
			result := runner.shouldSkipDirectory(dir)
			assert.True(t, result, "Directory %s should be skipped", dir)
		})
	}

	allowedDirs := []string{"lib", "test", "src", "bin", "example"}
	for _, dir := range allowedDirs {
		t.Run(dir+"_allowed", func(t *testing.T) {
			result := runner.shouldSkipDirectory(dir)
			assert.False(t, result, "Directory %s should not be skipped", dir)
		})
	}
}


// TestDartTestAnalyzer_AnalyzeTestDependencies tests dependency analysis
func TestDartTestAnalyzer_AnalyzeTestDependencies(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "dart-analyzer-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `import 'package:flutter_test/flutter_test.dart';
import 'package:my_app/services/user_service.dart';
import '../helpers/test_helper.dart';

void main() {
  test('test', () {});
}
`
	testFile := filepath.Join(tempDir, "user_service_test.dart")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	deps, err := analyzer.AnalyzeTestDependencies(context.Background(), testFile)
	require.NoError(t, err)

	assert.Len(t, deps, 3)
	assert.Contains(t, deps, "package:flutter_test/flutter_test.dart")
	assert.Contains(t, deps, "package:my_app/services/user_service.dart")
	assert.Contains(t, deps, "../helpers/test_helper.dart")
}

// TestDartTestAnalyzer_IsSmokeTest tests smoke test detection
func TestDartTestAnalyzer_IsSmokeTest(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "dart-smoke-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected bool
	}{
		{
			name:     "smoke_in_filename",
			fileName: "smoke_test.dart",
			content:  `void main() {}`,
			expected: true,
		},
		{
			name:     "smoke_in_content",
			fileName: "app_test.dart",
			content:  `// smoke test\nvoid main() {}`,
			expected: true,
		},
		{
			name:     "regular_test",
			fileName: "user_service_test.dart",
			content:  `void main() {}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result, err := analyzer.IsSmokeTest(context.Background(), testFile)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDartTestAnalyzer_ExtractTestFunctions tests extracting test functions
func TestDartTestAnalyzer_ExtractTestFunctions(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "dart-extract-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `void main() {
  test('should create user', () {});
  test('should update user', () {});
  testWidgets('should render widget', (tester) async {});
}
`
	testFile := filepath.Join(tempDir, "user_test.dart")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	functions, err := analyzer.ExtractTestFunctions(testFile)
	require.NoError(t, err)

	assert.Len(t, functions, 3)
	assert.Contains(t, functions, "should create user")
	assert.Contains(t, functions, "should update user")
	assert.Contains(t, functions, "should render widget")
}

// TestDartTestAnalyzer_ExtractTestGroups tests extracting test groups
func TestDartTestAnalyzer_ExtractTestGroups(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "dart-extract-groups-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `void main() {
  group('UserService', () {
    test('should create', () {});
  });
  
  group('OrderService', () {
    test('should process', () {});
  });
}
`
	testFile := filepath.Join(tempDir, "services_test.dart")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	groups, err := analyzer.ExtractTestGroups(testFile)
	require.NoError(t, err)

	assert.Len(t, groups, 2)
	assert.Contains(t, groups, "UserService")
	assert.Contains(t, groups, "OrderService")
}

// TestDartTestAnalyzer_FindTestsForFile tests finding tests for a source file
func TestDartTestAnalyzer_FindTestsForFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "dart-find-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source file
	libDir := filepath.Join(tempDir, "lib")
	err = os.MkdirAll(libDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(libDir, "user_service.dart")
	err = os.WriteFile(srcFile, []byte("class UserService {}"), 0o644)
	require.NoError(t, err)

	// Create test directory and test file
	testDir := filepath.Join(tempDir, "test")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	testFile := filepath.Join(testDir, "user_service_test.dart")
	testContent := `import 'package:my_app/user_service.dart';
void main() { test('test', () {}); }`
	err = os.WriteFile(testFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForFile(context.Background(), "lib/user_service.dart", tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 1)
}

// TestDartTestAnalyzer_ExtractPackageName tests package name extraction
func TestDartTestAnalyzer_ExtractPackageName(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "simple_pubspec",
			content: `name: my_app
version: 1.0.0`,
			expected: "my_app",
		},
		{
			name: "pubspec_with_spaces",
			content: `name:   my_flutter_app  
version: 1.0.0`,
			expected: "my_flutter_app",
		},
		{
			name:     "no_name",
			content:  `version: 1.0.0`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.extractPackageName(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDartTestAnalyzer_ShouldSkipDirectory tests directory skip logic
func TestDartTestAnalyzer_ShouldSkipDirectory(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	skipDirs := []string{
		"build", ".dart_tool", ".packages", ".git",
		"ios", "android", "web", "linux", "macos", "windows",
	}

	for _, dir := range skipDirs {
		t.Run(dir+"_should_skip", func(t *testing.T) {
			result := analyzer.shouldSkipDirectory(dir)
			assert.True(t, result, "Directory %s should be skipped", dir)
		})
	}

	allowedDirs := []string{"lib", "test", "src", "bin"}
	for _, dir := range allowedDirs {
		t.Run(dir+"_allowed", func(t *testing.T) {
			result := analyzer.shouldSkipDirectory(dir)
			assert.False(t, result, "Directory %s should not be skipped", dir)
		})
	}
}


// TestDartTestAnalyzer_RemoveDuplicates tests duplicate removal
func TestDartTestAnalyzer_RemoveDuplicates(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tests := []struct {
		name     string
		input    []string
		expected int
	}{
		{
			name:     "no_duplicates",
			input:    []string{"a", "b", "c"},
			expected: 3,
		},
		{
			name:     "with_duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: 3,
		},
		{
			name:     "empty",
			input:    []string{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.removeDuplicates(tt.input)
			assert.Len(t, result, tt.expected)
		})
	}
}

// TestDartTestAnalyzer_FileExists tests file existence check
func TestDartTestAnalyzer_FileExists(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "dart-exists-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a file
	testFile := filepath.Join(tempDir, "exists.dart")
	err = os.WriteFile(testFile, []byte("void main() {}"), 0o644)
	require.NoError(t, err)

	assert.True(t, analyzer.fileExists(testFile))
	assert.False(t, analyzer.fileExists(filepath.Join(tempDir, "notexists.dart")))
}

// TestDartTestAnalyzer_ImportsFile tests file import detection
func TestDartTestAnalyzer_ImportsFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tests := []struct {
		name       string
		content    string
		importPath string
		filePath   string
		expected   bool
	}{
		{
			name:       "package_import",
			content:    `import 'package:my_app/user_service.dart';`,
			importPath: "package:my_app/user_service.dart",
			filePath:   "lib/user_service.dart",
			expected:   true,
		},
		{
			name:       "relative_import",
			content:    `import '../user_service.dart';`,
			importPath: "",
			filePath:   "lib/user_service.dart",
			expected:   true,
		},
		{
			name:       "no_import",
			content:    `import 'package:my_app/order_service.dart';`,
			importPath: "package:my_app/user_service.dart",
			filePath:   "lib/user_service.dart",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.importsFile(tt.content, tt.importPath, tt.filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}


// TestDartTestAnalyzer_AnalyzeTestDependencies_NonExistent tests error handling
func TestDartTestAnalyzer_AnalyzeTestDependencies_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	_, err := analyzer.AnalyzeTestDependencies(context.Background(), "/nonexistent/file.dart")
	assert.Error(t, err)
}

// TestDartTestAnalyzer_IsSmokeTest_NonExistent tests error handling
func TestDartTestAnalyzer_IsSmokeTest_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	_, err := analyzer.IsSmokeTest(context.Background(), "/nonexistent/file.dart")
	assert.Error(t, err)
}

// TestDartTestAnalyzer_ExtractTestFunctions_NonExistent tests error handling
func TestDartTestAnalyzer_ExtractTestFunctions_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	_, err := analyzer.ExtractTestFunctions("/nonexistent/file.dart")
	assert.Error(t, err)
}

// TestDartTestAnalyzer_ExtractTestGroups_NonExistent tests error handling
func TestDartTestAnalyzer_ExtractTestGroups_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	_, err := analyzer.ExtractTestGroups("/nonexistent/file.dart")
	assert.Error(t, err)
}

// TestDartTestAnalyzer_ExtractImports tests import extraction
func TestDartTestAnalyzer_ExtractImports(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	content := `import 'package:flutter_test/flutter_test.dart';
import 'package:my_app/services/user_service.dart';
import '../helpers/test_helper.dart';
// import 'commented.dart';

void main() {}`

	imports := analyzer.extractImports(content)

	assert.Len(t, imports, 3)
	assert.Contains(t, imports, "package:flutter_test/flutter_test.dart")
	assert.Contains(t, imports, "package:my_app/services/user_service.dart")
	assert.Contains(t, imports, "../helpers/test_helper.dart")
}


// TestDartTestAnalyzer_ExtractImports_EmptyContent tests extracting imports from empty content
func TestDartTestAnalyzer_ExtractImports_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	imports := analyzer.extractImports("")
	assert.Empty(t, imports)
}

// TestDartTestAnalyzer_ExtractImports_DuplicateImports tests duplicate import handling
func TestDartTestAnalyzer_ExtractImports_DuplicateImports(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	content := `import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:my_app/user.dart';`

	imports := analyzer.extractImports(content)
	assert.Len(t, imports, 2)
}

// TestDartTestAnalyzer_ExtractImports_SkipComments tests that commented imports are skipped
func TestDartTestAnalyzer_ExtractImports_SkipComments(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	content := `import 'package:flutter_test/flutter_test.dart';
// import 'package:commented/commented.dart';
import 'package:my_app/user.dart';`

	imports := analyzer.extractImports(content)
	assert.Len(t, imports, 2)
	assert.NotContains(t, imports, "package:commented/commented.dart")
}

// TestDartTestAnalyzer_ExtractPackageName_EmptyContent tests package name extraction with empty content
func TestDartTestAnalyzer_ExtractPackageName_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	result := analyzer.extractPackageName("")
	assert.Empty(t, result)
}

// TestDartTestAnalyzer_ExtractPackageName_NoName tests package name extraction with no name field
func TestDartTestAnalyzer_ExtractPackageName_NoName(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	content := `version: 1.0.0
description: A test package`

	result := analyzer.extractPackageName(content)
	assert.Empty(t, result)
}

// TestDartTestAnalyzer_FindTestsForFile_NoTestDir tests when test directory doesn't exist
func TestDartTestAnalyzer_FindTestsForFile_NoTestDir(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "dart-no-test-dir-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create lib directory without test directory
	libDir := filepath.Join(tempDir, "lib")
	err = os.MkdirAll(libDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(libDir, "user_service.dart")
	err = os.WriteFile(srcFile, []byte("class UserService {}"), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForFile(context.Background(), "lib/user_service.dart", tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestDartTestAnalyzer_IsSmokeTest_EmptyFile tests smoke test detection with empty file
func TestDartTestAnalyzer_IsSmokeTest_EmptyFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "dart-empty-smoke-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "empty_test.dart")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	result, err := analyzer.IsSmokeTest(context.Background(), testFile)
	require.NoError(t, err)
	assert.False(t, result)
}

// TestDartTestAnalyzer_ExtractTestFunctions_EmptyFile tests extracting functions from empty file
func TestDartTestAnalyzer_ExtractTestFunctions_EmptyFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "dart-empty-funcs-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "empty_test.dart")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	functions, err := analyzer.ExtractTestFunctions(testFile)
	require.NoError(t, err)
	assert.Empty(t, functions)
}

// TestDartTestAnalyzer_ExtractTestGroups_EmptyFile tests extracting groups from empty file
func TestDartTestAnalyzer_ExtractTestGroups_EmptyFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewDartTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "dart-empty-groups-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "empty_test.dart")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	groups, err := analyzer.ExtractTestGroups(testFile)
	require.NoError(t, err)
	assert.Empty(t, groups)
}

// TestDartTestRunner_parseTestOutput_EmptyOutput tests parsing empty output
func TestDartTestRunner_parseTestOutput_EmptyOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	testInfos := []*domain.TestInfo{
		{Path: "test/user_test.dart", Name: "user_test.dart"},
	}

	results := runner.parseTestOutput("", testInfos, 1.0)
	require.Len(t, results, 1)
	// Empty output defaults to success
	assert.True(t, results[0].Success)
}

// TestDartTestRunner_parseTestOutput_FailedKeyword tests parsing output with FAILED keyword
func TestDartTestRunner_parseTestOutput_FailedKeyword(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	testInfos := []*domain.TestInfo{
		{Path: "test/user_test.dart", Name: "user_test.dart"},
	}

	results := runner.parseTestOutput("user_test FAILED", testInfos, 1.0)
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
}

// TestDartTestRunner_analyzeTestFile_EmptyFile tests analyzing empty file
func TestDartTestRunner_analyzeTestFile_EmptyFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tempDir, err := os.MkdirTemp("", "dart-empty-analyze-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "empty_test.dart")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	result := runner.analyzeTestFile(testFile)
	assert.Equal(t, "unit", result)
}

// TestDartTestRunner_isFlutterProject_EmptyPubspec tests Flutter detection with empty pubspec
func TestDartTestRunner_isFlutterProject_EmptyPubspec(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tempDir, err := os.MkdirTemp("", "dart-empty-pubspec-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	err = os.WriteFile(filepath.Join(tempDir, "pubspec.yaml"), []byte(""), 0o644)
	require.NoError(t, err)

	result := runner.isFlutterProject(tempDir)
	assert.False(t, result)
}

// TestDartTestRunner_DiscoverTests_SkipPlatformDirs tests that platform directories are skipped
func TestDartTestRunner_DiscoverTests_SkipPlatformDirs(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dart-skip-platform-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test directory
	testDir := filepath.Join(tempDir, "test")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	// Create valid test file
	validTest := filepath.Join(testDir, "valid_test.dart")
	err = os.WriteFile(validTest, []byte("void main() { test('test', () {}); }"), 0o644)
	require.NoError(t, err)

	// Create android directory inside test (should be skipped)
	androidDir := filepath.Join(testDir, "android")
	err = os.MkdirAll(androidDir, 0o755)
	require.NoError(t, err)

	androidTest := filepath.Join(androidDir, "android_test.dart")
	err = os.WriteFile(androidTest, []byte("void main() {}"), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewDartTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Verify android test is not found
	for _, test := range tests {
		assert.NotEqual(t, "android_test.dart", test.Name)
	}
}
