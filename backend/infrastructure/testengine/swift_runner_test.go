package testengine

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"syntaxia/domain"
)

// TestSwiftTestRunner_GetLanguage tests language identifier
func TestSwiftTestRunner_GetLanguage(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	assert.Equal(t, "swift", runner.GetLanguage())
}

// TestSwiftTestRunner_extractTestName tests test name extraction from path
func TestSwiftTestRunner_extractTestName(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests := []struct {
		name     string
		testPath string
		expected string
	}{
		{
			name:     "simple_file",
			testPath: "MyTests.swift",
			expected: "MyTests",
		},
		{
			name:     "nested_file",
			testPath: "Tests/MyPackageTests/MyTests.swift",
			expected: "MyTests",
		},
		{
			name:     "xctest_file",
			testPath: "UnitTests/ViewModelTests.swift",
			expected: "ViewModelTests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.extractTestName(tt.testPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSwiftTestRunner_detectBuildSystem tests build system detection
func TestSwiftTestRunner_detectBuildSystem(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected string
	}{
		{
			name: "spm_project",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "Package.swift"), []byte("// swift-tools-version:5.5"), 0o644)
			},
			expected: "spm",
		},
		{
			name: "xcode_project",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, "MyApp.xcodeproj"), 0o755)
			},
			expected: "xcode",
		},
		{
			name: "xcode_workspace",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, "MyApp.xcworkspace"), 0o755)
			},
			expected: "xcode",
		},
		{
			name:     "no_build_system",
			setup:    func(dir string) error { return nil },
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "swift-build-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			err = tt.setup(tempDir)
			require.NoError(t, err)

			result := runner.detectBuildSystem(tempDir)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSwiftTestRunner_fileHasTests tests detection of XCTest in files
func TestSwiftTestRunner_fileHasTests(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "swift-runner-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "has_xctest_import",
			content: `import XCTest

class MyTests: XCTestCase {
    func testExample() {
        XCTAssertTrue(true)
    }
}
`,
			expected: true,
		},
		{
			name: "has_test_function",
			content: `import XCTest

class MyTests: XCTestCase {
    func testSomething() {
        XCTAssertEqual(1, 1)
    }
}
`,
			expected: true,
		},
		{
			name: "inherits_xctestcase",
			content: `import XCTest

final class ViewModelTests: XCTestCase {
    override func setUp() {}
}
`,
			expected: true,
		},
		{
			name: "no_tests",
			content: `import Foundation

struct MyModel {
    let name: String
}
`,
			expected: false,
		},
		{
			name: "test_in_comment",
			content: `// This is not a real test file
// It just has some comments
func notATest() {}
`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.name+".swift")
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result, err := runner.fileHasTests(testFile)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSwiftTestRunner_isSmokeTest tests smoke test detection
func TestSwiftTestRunner_isSmokeTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{
			name:     "smoke_test_file",
			filePath: "Tests/SmokeTests.swift",
			expected: true,
		},
		{
			name:     "smoke_in_name",
			filePath: "Tests/AppSmokeTest.swift",
			expected: true,
		},
		{
			name:     "regular_test",
			filePath: "Tests/ViewModelTests.swift",
			expected: false,
		},
		{
			name:     "unit_test",
			filePath: "Tests/UnitTests.swift",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.isSmokeTest(tt.filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSwiftTestRunner_parseTestOutput tests parsing of swift test output
func TestSwiftTestRunner_parseTestOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests := []struct {
		name           string
		output         string
		expectedPassed int
		expectedFailed int
	}{
		{
			name:           "all_passed",
			output:         "Test Suite 'All tests' passed at 2024-01-15.\nExecuted 5 tests, with 0 failures",
			expectedPassed: 5,
			expectedFailed: 0,
		},
		{
			name:           "some_failed",
			output:         "Test Suite 'All tests' failed at 2024-01-15.\nExecuted 10 tests, with 3 failures",
			expectedPassed: 7,
			expectedFailed: 3,
		},
		{
			name:           "single_test",
			output:         "Executed 1 test, with 0 failures",
			expectedPassed: 1,
			expectedFailed: 0,
		},
		{
			name:           "no_match",
			output:         "Running tests...",
			expectedPassed: 0,
			expectedFailed: 0,
		},
		{
			name:           "empty_output",
			output:         "",
			expectedPassed: 0,
			expectedFailed: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passed, failed, _ := runner.parseTestOutput(tt.output)
			assert.Equal(t, tt.expectedPassed, passed, "passed count mismatch")
			assert.Equal(t, tt.expectedFailed, failed, "failed count mismatch")
		})
	}
}

// TestSwiftTestRunner_DiscoverTests tests test discovery in Swift projects
func TestSwiftTestRunner_DiscoverTests(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "swift-discover-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create Package.swift
	packageSwift := `// swift-tools-version:5.5
import PackageDescription

let package = Package(
    name: "TestProject",
    targets: [
        .target(name: "TestProject"),
        .testTarget(name: "TestProjectTests", dependencies: ["TestProject"]),
    ]
)
`
	err = os.WriteFile(filepath.Join(tempDir, "Package.swift"), []byte(packageSwift), 0o644)
	require.NoError(t, err)

	// Create Tests directory with test files
	testsDir := filepath.Join(tempDir, "Tests", "TestProjectTests")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	// Create test file with XCTest
	testFile := `import XCTest
@testable import TestProject

final class TestProjectTests: XCTestCase {
    func testExample() {
        XCTAssertTrue(true)
    }
}
`
	err = os.WriteFile(filepath.Join(testsDir, "TestProjectTests.swift"), []byte(testFile), 0o644)
	require.NoError(t, err)

	// Create smoke test
	smokeTest := `import XCTest

final class SmokeTests: XCTestCase {
    func testSmoke() {
        XCTAssertTrue(true)
    }
}
`
	err = os.WriteFile(filepath.Join(testsDir, "SmokeTests.swift"), []byte(smokeTest), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 2)

	// Verify test types
	foundTypes := make(map[string]bool)
	for _, test := range tests {
		foundTypes[test.Type] = true
		assert.Equal(t, "swift", test.Metadata["language"])
	}

	assert.True(t, foundTypes["unit"] || foundTypes["smoke"], "Should find unit or smoke tests")
}

// TestSwiftTestRunner_DiscoverTests_NoPackage tests discovery without Package.swift
func TestSwiftTestRunner_DiscoverTests_NoPackage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "swift-no-package-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	_, err = runner.DiscoverTests(context.Background(), tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no Swift build system detected")
}

// TestSwiftTestRunner_discoverSPMTests tests SPM test discovery
func TestSwiftTestRunner_discoverSPMTests(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "swift-spm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testsDir := filepath.Join(tempDir, "Tests", "MyTests")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	// Create file with tests
	fileWithTests := `import XCTest

class MyTests: XCTestCase {
    func testSomething() {}
}
`
	err = os.WriteFile(filepath.Join(testsDir, "MyTests.swift"), []byte(fileWithTests), 0o644)
	require.NoError(t, err)

	// Create file without tests
	fileWithoutTests := `import Foundation

struct Helper {}
`
	err = os.WriteFile(filepath.Join(testsDir, "Helper.swift"), []byte(fileWithoutTests), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests, err := runner.discoverSPMTests(tempDir)
	require.NoError(t, err)

	// Should only find file with tests
	assert.Len(t, tests, 1)
	assert.Equal(t, "MyTests.swift", tests[0].Name)
	assert.Equal(t, "swift", tests[0].Metadata["language"])
}

// TestSwiftTestRunner_RunTestSuite tests running a test suite
func TestSwiftTestRunner_RunTestSuite(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	suite := &domain.TestSuite{
		Name:        "test_suite",
		Language:    "swift",
		ProjectPath: "/nonexistent/path",
		Tests: []*domain.TestInfo{
			{Path: "Tests/MyTests.swift", Name: "MyTests.swift", Type: "unit"},
		},
		Config: &domain.TestConfig{
			ProjectPath: "/nonexistent/path",
			Verbose:     false,
		},
	}

	results, err := runner.RunTestSuite(context.Background(), suite)

	assert.NoError(t, err)
	assert.NotNil(t, results)
}

// TestSwiftTestRunner_RunTest tests single test execution
func TestSwiftTestRunner_RunTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	config := &domain.TestConfig{
		ProjectPath: "/nonexistent/path",
		Verbose:     false,
		Timeout:     60,
	}

	result, err := runner.RunTest(context.Background(), "Tests/MyTests.swift", config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Tests/MyTests.swift", result.TestPath)
	assert.Equal(t, "swift", result.Language)
	assert.False(t, result.Success) // Should fail because path doesn't exist
}

// TestSwiftTestRunner_getTestFilter tests test filter generation
func TestSwiftTestRunner_getTestFilter(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests := []struct {
		name     string
		testPath string
		expected string
	}{
		{
			name:     "simple_test",
			testPath: "MyTests.swift",
			expected: "MyTests",
		},
		{
			name:     "nested_test",
			testPath: "Tests/PackageTests/ViewModelTests.swift",
			expected: "ViewModelTests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.getTestFilter(tt.testPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSwiftTestRunner_extractFailureMessage tests failure message extraction
func TestSwiftTestRunner_extractFailureMessage(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name:     "error_message",
			output:   "error: cannot find 'MyType' in scope\nBuild failed",
			expected: "error: cannot find 'MyType' in scope",
		},
		{
			name:     "test_failed",
			output:   "Test Suite 'MyTests' failed at 2024-01-15",
			expected: "Test Suite 'MyTests' failed at 2024-01-15",
		},
		{
			name:     "no_failure",
			output:   "Build succeeded\nAll tests passed",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.extractFailureMessage(tt.output)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSwiftTestRunner_discoverXcodeTests tests Xcode test discovery
func TestSwiftTestRunner_discoverXcodeTests(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "swift-xcode-discover-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create Xcode project structure
	xcodeProj := filepath.Join(tempDir, "MyApp.xcodeproj")
	err = os.MkdirAll(xcodeProj, 0o755)
	require.NoError(t, err)

	// Create test directories
	testDirs := []string{"Tests", "UnitTests", "UITests"}
	for _, dir := range testDirs {
		fullDir := filepath.Join(tempDir, dir)
		err = os.MkdirAll(fullDir, 0o755)
		require.NoError(t, err)

		// Create test file with XCTest
		testFile := filepath.Join(fullDir, dir+"File.swift")
		content := `import XCTest
class ` + dir + `File: XCTestCase {
    func testExample() {}
}`
		err = os.WriteFile(testFile, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests, err := runner.discoverXcodeTests(tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 3)
}

// TestSwiftTestRunner_detectScheme tests Xcode scheme detection
func TestSwiftTestRunner_detectScheme(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected string
	}{
		{
			name: "with_xcodeproj",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, "MyApp.xcodeproj"), 0o755)
			},
			expected: "MyApp",
		},
		{
			name: "no_project",
			setup: func(dir string) error {
				return nil
			},
			expected: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "swift-scheme-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			err = tt.setup(tempDir)
			require.NoError(t, err)

			result := runner.detectScheme(tempDir)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSwiftTestRunner_runAllTests tests running all tests
func TestSwiftTestRunner_runAllTests(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tempDir, err := os.MkdirTemp("", "swift-all-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		Verbose:     false,
	}

	result, err := runner.runAllTests(context.Background(), tempDir, config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tempDir, result.TestPath)
	assert.Equal(t, "all", result.TestName)
	assert.Equal(t, "swift", result.Language)
	assert.False(t, result.Success) // Should fail because no build system
}

// TestSwiftTestRunner_applyEnvVars tests environment variable application
func TestSwiftTestRunner_applyEnvVars(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests := []struct {
		name    string
		envVars map[string]string
	}{
		{
			name:    "with_env_vars",
			envVars: map[string]string{"TEST_VAR": "test_value", "ANOTHER_VAR": "another_value"},
		},
		{
			name:    "nil_env_vars",
			envVars: nil,
		},
		{
			name:    "empty_env_vars",
			envVars: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &domain.TestConfig{
				EnvVars: tt.envVars,
			}

			cmd := &exec.Cmd{}
			runner.applyEnvVars(cmd, config)

			if tt.envVars != nil && len(tt.envVars) > 0 {
				assert.NotNil(t, cmd.Env)
			}
		})
	}
}

// TestSwiftTestRunner_DiscoverTests_XcodeWorkspace tests discovery with workspace
func TestSwiftTestRunner_DiscoverTests_XcodeWorkspace(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "swift-workspace-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create workspace
	workspace := filepath.Join(tempDir, "MyApp.xcworkspace")
	err = os.MkdirAll(workspace, 0o755)
	require.NoError(t, err)

	// Create Tests directory
	testsDir := filepath.Join(tempDir, "Tests")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	// Create test file
	testFile := `import XCTest
class MyTests: XCTestCase {
    func testExample() {}
}`
	err = os.WriteFile(filepath.Join(testsDir, "MyTests.swift"), []byte(testFile), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 1)
}

// TestSwiftTestRunner_fileHasTests_ReadError tests file read error handling
func TestSwiftTestRunner_fileHasTests_ReadError(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	result, err := runner.fileHasTests("/nonexistent/path/test.swift")

	assert.Error(t, err)
	assert.False(t, result)
}


// TestSwiftTestAnalyzer_AnalyzeTestDependencies tests dependency analysis
func TestSwiftTestAnalyzer_AnalyzeTestDependencies(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-analyzer-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `import XCTest
import Foundation

final class MyTests: XCTestCase {
    func testExample() {}
}
`
	testFile := filepath.Join(tempDir, "MyTests.swift")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	deps, err := analyzer.AnalyzeTestDependencies(context.Background(), testFile)
	require.NoError(t, err)

	assert.Len(t, deps, 2)
	assert.Contains(t, deps, "XCTest")
	assert.Contains(t, deps, "Foundation")
}

// TestSwiftTestAnalyzer_IsSmokeTest tests smoke test detection
func TestSwiftTestAnalyzer_IsSmokeTest(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-smoke-test-*")
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
			fileName: "SmokeTests.swift",
			content:  `import XCTest`,
			expected: true,
		},
		{
			name:     "smoke_in_content",
			fileName: "AppTests.swift",
			content:  `// smoke test\nimport XCTest`,
			expected: true,
		},
		{
			name:     "regular_test",
			fileName: "UserTests.swift",
			content:  `import XCTest`,
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

// TestSwiftTestAnalyzer_ExtractTestFunctions tests extracting test functions
func TestSwiftTestAnalyzer_ExtractTestFunctions(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-extract-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `import XCTest

final class MyTests: XCTestCase {
    func testExample() {}
    func testAnotherExample() {}
    func helperMethod() {}
    func testWithParameters(param: String) {}
}
`
	testFile := filepath.Join(tempDir, "MyTests.swift")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	functions, err := analyzer.ExtractTestFunctions(testFile)
	require.NoError(t, err)

	assert.Len(t, functions, 3)
	assert.Contains(t, functions, "testExample")
	assert.Contains(t, functions, "testAnotherExample")
	assert.Contains(t, functions, "testWithParameters")
}

// TestSwiftTestAnalyzer_ExtractTestClasses tests extracting test classes
func TestSwiftTestAnalyzer_ExtractTestClasses(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-extract-class-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `import XCTest

final class UserTests: XCTestCase {
    func testUser() {}
}

class OrderTest: XCTestCase {
    func testOrder() {}
}

class Helper {
    func help() {}
}
`
	testFile := filepath.Join(tempDir, "Tests.swift")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	classes, err := analyzer.ExtractTestClasses(testFile)
	require.NoError(t, err)

	assert.Len(t, classes, 2)
	assert.Contains(t, classes, "UserTests")
	assert.Contains(t, classes, "OrderTest")
}

// TestSwiftTestAnalyzer_FindTestsForFile tests finding tests for a source file
func TestSwiftTestAnalyzer_FindTestsForFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-find-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source file
	srcDir := filepath.Join(tempDir, "Sources", "MyApp")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(srcDir, "User.swift")
	err = os.WriteFile(srcFile, []byte("struct User {}"), 0o644)
	require.NoError(t, err)

	// Create test directory and test file
	testDir := filepath.Join(tempDir, "Tests")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	testFile := filepath.Join(testDir, "UserTests.swift")
	testContent := `import XCTest
@testable import MyApp

final class UserTests: XCTestCase {}`
	err = os.WriteFile(testFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForFile(context.Background(), "Sources/MyApp/User.swift", tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 1)
}

// TestSwiftTestAnalyzer_ShouldSkipDirectory tests directory skip logic
func TestSwiftTestAnalyzer_ShouldSkipDirectory(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	skipDirs := []string{".build", ".git", "DerivedData", "Pods", "Carthage", ".swiftpm"}
	for _, dir := range skipDirs {
		t.Run(dir+"_should_skip", func(t *testing.T) {
			result := analyzer.shouldSkipDirectory(dir)
			assert.True(t, result, "Directory %s should be skipped", dir)
		})
	}

	allowedDirs := []string{"Sources", "Tests", "Package"}
	for _, dir := range allowedDirs {
		t.Run(dir+"_allowed", func(t *testing.T) {
			result := analyzer.shouldSkipDirectory(dir)
			assert.False(t, result, "Directory %s should not be skipped", dir)
		})
	}
}


// TestSwiftTestAnalyzer_RemoveDuplicates tests duplicate removal
func TestSwiftTestAnalyzer_RemoveDuplicates(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

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

// TestSwiftTestAnalyzer_FileExists tests file existence check
func TestSwiftTestAnalyzer_FileExists(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-exists-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a file
	testFile := filepath.Join(tempDir, "exists.swift")
	err = os.WriteFile(testFile, []byte("import Foundation"), 0o644)
	require.NoError(t, err)

	assert.True(t, analyzer.fileExists(testFile))
	assert.False(t, analyzer.fileExists(filepath.Join(tempDir, "notexists.swift")))
}

// TestSwiftTestAnalyzer_ImportsModule tests module import detection
func TestSwiftTestAnalyzer_ImportsModule(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tests := []struct {
		name       string
		content    string
		moduleName string
		expected   bool
	}{
		{
			name:       "imports_module",
			content:    `import MyApp`,
			moduleName: "MyApp",
			expected:   true,
		},
		{
			name:       "testable_import",
			content:    `@testable import MyApp`,
			moduleName: "MyApp",
			expected:   true,
		},
		{
			name:       "no_import",
			content:    `import Foundation`,
			moduleName: "MyApp",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.importsModule(tt.content, tt.moduleName)
			assert.Equal(t, tt.expected, result)
		})
	}
}


// TestSwiftTestAnalyzer_AnalyzeTestDependencies_NonExistent tests error handling
func TestSwiftTestAnalyzer_AnalyzeTestDependencies_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	_, err := analyzer.AnalyzeTestDependencies(context.Background(), "/nonexistent/file.swift")
	assert.Error(t, err)
}

// TestSwiftTestAnalyzer_IsSmokeTest_NonExistent tests error handling
func TestSwiftTestAnalyzer_IsSmokeTest_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	_, err := analyzer.IsSmokeTest(context.Background(), "/nonexistent/file.swift")
	assert.Error(t, err)
}

// TestSwiftTestAnalyzer_ExtractTestFunctions_NonExistent tests error handling
func TestSwiftTestAnalyzer_ExtractTestFunctions_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	_, err := analyzer.ExtractTestFunctions("/nonexistent/file.swift")
	assert.Error(t, err)
}

// TestSwiftTestAnalyzer_ExtractTestClasses_NonExistent tests error handling
func TestSwiftTestAnalyzer_ExtractTestClasses_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	_, err := analyzer.ExtractTestClasses("/nonexistent/file.swift")
	assert.Error(t, err)
}

// TestSwiftTestAnalyzer_ExtractImports tests import extraction
func TestSwiftTestAnalyzer_ExtractImports(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	content := `import XCTest
import Foundation
import UIKit
// import Commented

final class MyTests: XCTestCase {}`

	imports := analyzer.extractImports(content)

	assert.Len(t, imports, 3)
	assert.Contains(t, imports, "XCTest")
	assert.Contains(t, imports, "Foundation")
	assert.Contains(t, imports, "UIKit")
}


// TestSwiftTestAnalyzer_FindTestTargets tests finding test targets
func TestSwiftTestAnalyzer_FindTestTargets(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-targets-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create Tests directory with subdirectories
	testsDir := filepath.Join(tempDir, "Tests")
	err = os.MkdirAll(filepath.Join(testsDir, "MyAppTests"), 0o755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(testsDir, "MyAppUITests"), 0o755)
	require.NoError(t, err)

	targets, err := analyzer.findTestTargets(tempDir)
	require.NoError(t, err)

	assert.Len(t, targets, 2)
	assert.Contains(t, targets, "MyAppTests")
	assert.Contains(t, targets, "MyAppUITests")
}

// TestSwiftTestAnalyzer_FindTestTargets_NoTestsDir tests when Tests directory doesn't exist
func TestSwiftTestAnalyzer_FindTestTargets_NoTestsDir(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-no-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	_, err = analyzer.findTestTargets(tempDir)
	assert.Error(t, err)
}

// TestSwiftTestAnalyzer_DetectModuleName tests module name detection
func TestSwiftTestAnalyzer_DetectModuleName(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-module-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create Package.swift
	packageContent := `// swift-tools-version:5.5
import PackageDescription

let package = Package(
    name: "MyApp",
    targets: [
        .target(name: "MyApp"),
    ]
)
`
	err = os.WriteFile(filepath.Join(tempDir, "Package.swift"), []byte(packageContent), 0o644)
	require.NoError(t, err)

	// Create Sources directory
	srcDir := filepath.Join(tempDir, "Sources", "MyApp")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(srcDir, "User.swift")
	err = os.WriteFile(srcFile, []byte("struct User {}"), 0o644)
	require.NoError(t, err)

	moduleName := analyzer.detectModuleName(srcFile, tempDir)
	assert.Equal(t, "MyApp", moduleName)
}


// TestSwiftTestRunner_parseTestOutput_EmptyOutput tests parsing empty output
func TestSwiftTestRunner_parseTestOutput_EmptyOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	passed, failed, skipped := runner.parseTestOutput("")
	assert.Equal(t, 0, passed)
	assert.Equal(t, 0, failed)
	assert.Equal(t, 0, skipped)
}

// TestSwiftTestRunner_parseTestOutput_MalformedOutput tests parsing malformed output
func TestSwiftTestRunner_parseTestOutput_MalformedOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	passed, failed, skipped := runner.parseTestOutput("Some random output without test results")
	assert.Equal(t, 0, passed)
	assert.Equal(t, 0, failed)
	assert.Equal(t, 0, skipped)
}

// TestSwiftTestRunner_extractFailureMessage_NoFailure tests extracting failure message when none exists
func TestSwiftTestRunner_extractFailureMessage_NoFailure(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	result := runner.extractFailureMessage("Build succeeded\nAll tests passed")
	assert.Empty(t, result)
}

// TestSwiftTestRunner_extractFailureMessage_MultipleErrors tests extracting first error
func TestSwiftTestRunner_extractFailureMessage_MultipleErrors(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	output := `error: first error
error: second error
Build failed`
	result := runner.extractFailureMessage(output)
	assert.Contains(t, result, "first error")
}

// TestSwiftTestRunner_isSmokeTest_CaseInsensitive tests case insensitive smoke test detection
func TestSwiftTestRunner_isSmokeTest_CaseInsensitive(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{"lowercase_smoke", "Tests/smoketests.swift", true},
		{"uppercase_smoke", "Tests/SMOKETESTS.swift", true},
		{"mixed_case_smoke", "Tests/SmokeTests.swift", true},
		{"no_smoke", "Tests/UnitTests.swift", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.isSmokeTest(tt.filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSwiftTestRunner_fileHasTests_EmptyFile tests checking empty file for tests
func TestSwiftTestRunner_fileHasTests_EmptyFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tempDir, err := os.MkdirTemp("", "swift-empty-file-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "Empty.swift")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	result, err := runner.fileHasTests(testFile)
	require.NoError(t, err)
	assert.False(t, result)
}

// TestSwiftTestRunner_discoverSPMTests_NoTestsDir tests SPM discovery without Tests directory
func TestSwiftTestRunner_discoverSPMTests_NoTestsDir(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tempDir, err := os.MkdirTemp("", "swift-no-tests-dir-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests, err := runner.discoverSPMTests(tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestSwiftTestRunner_discoverXcodeTests_NoTestDirs tests Xcode discovery without test directories
func TestSwiftTestRunner_discoverXcodeTests_NoTestDirs(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tempDir, err := os.MkdirTemp("", "swift-no-xcode-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests, err := runner.discoverXcodeTests(tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestSwiftTestRunner_detectScheme_NoProject tests scheme detection without project
func TestSwiftTestRunner_detectScheme_NoProject(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tempDir, err := os.MkdirTemp("", "swift-no-project-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	result := runner.detectScheme(tempDir)
	assert.Equal(t, "default", result)
}

// TestSwiftTestRunner_getTestFilter_NestedPath tests test filter with nested path
func TestSwiftTestRunner_getTestFilter_NestedPath(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tests := []struct {
		name     string
		testPath string
		expected string
	}{
		{"deeply_nested", "Tests/MyApp/Unit/Services/UserServiceTests.swift", "UserServiceTests"},
		{"simple", "UserTests.swift", "UserTests"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.getTestFilter(tt.testPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSwiftTestRunner_RunTestSuite_EmptySuite tests running empty test suite
func TestSwiftTestRunner_RunTestSuite_EmptySuite(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewSwiftTestRunner(log)

	tempDir, err := os.MkdirTemp("", "swift-empty-suite-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	suite := &domain.TestSuite{
		Name:        "empty_suite",
		Language:    "swift",
		ProjectPath: tempDir,
		Tests:       []*domain.TestInfo{},
		Config: &domain.TestConfig{
			ProjectPath: tempDir,
			Verbose:     false,
		},
	}

	results, err := runner.RunTestSuite(context.Background(), suite)
	require.NoError(t, err)
	// Empty suite runs all tests
	assert.Len(t, results, 1)
}

// TestSwiftTestAnalyzer_ExtractImports_EmptyContent tests extracting imports from empty content
func TestSwiftTestAnalyzer_ExtractImports_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	imports := analyzer.extractImports("")
	assert.Empty(t, imports)
}

// TestSwiftTestAnalyzer_ExtractImports_DuplicateImports tests duplicate import handling
func TestSwiftTestAnalyzer_ExtractImports_DuplicateImports(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	content := `import XCTest
import Foundation
import XCTest
import Foundation`

	imports := analyzer.extractImports(content)
	assert.Len(t, imports, 2)
}

// TestSwiftTestAnalyzer_ExtractImports_SkipComments tests that commented imports are skipped
func TestSwiftTestAnalyzer_ExtractImports_SkipComments(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	content := `import XCTest
// import Commented
import Foundation`

	imports := analyzer.extractImports(content)
	assert.Len(t, imports, 2)
	assert.NotContains(t, imports, "Commented")
}

// TestSwiftTestAnalyzer_DetectModuleName_NoPackageSwift tests module detection without Package.swift
func TestSwiftTestAnalyzer_DetectModuleName_NoPackageSwift(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-no-package-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	srcFile := filepath.Join(tempDir, "Sources", "MyApp", "User.swift")
	err = os.MkdirAll(filepath.Dir(srcFile), 0o755)
	require.NoError(t, err)
	err = os.WriteFile(srcFile, []byte("struct User {}"), 0o644)
	require.NoError(t, err)

	moduleName := analyzer.detectModuleName(srcFile, tempDir)
	assert.Empty(t, moduleName)
}

// TestSwiftTestAnalyzer_FindTestsForFile_NoTestsDir tests finding tests without Tests directory
func TestSwiftTestAnalyzer_FindTestsForFile_NoTestsDir(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-no-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source file without Tests directory
	srcDir := filepath.Join(tempDir, "Sources", "MyApp")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(srcDir, "User.swift")
	err = os.WriteFile(srcFile, []byte("struct User {}"), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForFile(context.Background(), "Sources/MyApp/User.swift", tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestSwiftTestAnalyzer_IsSmokeTest_EmptyFile tests smoke test detection with empty file
func TestSwiftTestAnalyzer_IsSmokeTest_EmptyFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-empty-smoke-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "EmptyTests.swift")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	result, err := analyzer.IsSmokeTest(context.Background(), testFile)
	require.NoError(t, err)
	assert.False(t, result)
}

// TestSwiftTestAnalyzer_ExtractTestFunctions_EmptyFile tests extracting functions from empty file
func TestSwiftTestAnalyzer_ExtractTestFunctions_EmptyFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-empty-funcs-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "EmptyTests.swift")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	functions, err := analyzer.ExtractTestFunctions(testFile)
	require.NoError(t, err)
	assert.Empty(t, functions)
}

// TestSwiftTestAnalyzer_ExtractTestClasses_EmptyFile tests extracting classes from empty file
func TestSwiftTestAnalyzer_ExtractTestClasses_EmptyFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "swift-empty-classes-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "EmptyTests.swift")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	classes, err := analyzer.ExtractTestClasses(testFile)
	require.NoError(t, err)
	assert.Empty(t, classes)
}

// TestSwiftTestAnalyzer_ImportsModule_TestableImport tests @testable import detection
func TestSwiftTestAnalyzer_ImportsModule_TestableImport(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	content := `@testable import MyApp
import XCTest`

	result := analyzer.importsModule(content, "MyApp")
	assert.True(t, result)
}

// TestSwiftTestAnalyzer_ImportsModule_NoMatch tests import detection with no match
func TestSwiftTestAnalyzer_ImportsModule_NoMatch(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewSwiftTestAnalyzer(log)

	content := `import Foundation
import XCTest`

	result := analyzer.importsModule(content, "MyApp")
	assert.False(t, result)
}
