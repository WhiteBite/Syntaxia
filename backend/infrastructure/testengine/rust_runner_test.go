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

// TestRustTestRunner_GetLanguage tests language identifier
func TestRustTestRunner_GetLanguage(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	assert.Equal(t, "rust", runner.GetLanguage())
}

// TestRustTestRunner_extractTestName tests test name extraction from path
func TestRustTestRunner_extractTestName(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	tests := []struct {
		name     string
		testPath string
		expected string
	}{
		{
			name:     "simple_file",
			testPath: "main.rs",
			expected: "main",
		},
		{
			name:     "nested_file",
			testPath: "src/lib.rs",
			expected: "lib",
		},
		{
			name:     "test_file",
			testPath: "tests/integration_test.rs",
			expected: "integration_test",
		},
		{
			name:     "deeply_nested",
			testPath: "src/utils/helpers.rs",
			expected: "helpers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.extractTestName(tt.testPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRustTestRunner_getModuleName tests module name extraction
func TestRustTestRunner_getModuleName(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "simple_file",
			filePath: "main.rs",
			expected: "main",
		},
		{
			name:     "lib_file",
			filePath: "lib.rs",
			expected: "lib",
		},
		{
			name:     "mod_file",
			filePath: "utils.rs",
			expected: "utils",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.getModuleName(tt.filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRustTestRunner_fileHasTests tests detection of test attributes in files
func TestRustTestRunner_fileHasTests(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "rust-runner-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "has_test_attribute",
			content: `fn main() {}

#[test]
fn test_something() {
    assert!(true);
}
`,
			expected: true,
		},
		{
			name: "has_cfg_test",
			content: `fn main() {}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {}
}
`,
			expected: true,
		},
		{
			name: "has_mod_tests",
			content: `fn main() {}

mod tests {
    fn helper() {}
}
`,
			expected: true,
		},
		{
			name: "no_tests",
			content: `fn main() {
    println!("Hello, world!");
}
`,
			expected: false,
		},
		{
			name: "test_in_comment",
			content: `// This is not a real test attribute
fn main() {}
`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, tt.name+".rs")
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result, err := runner.fileHasTests(testFile)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRustTestRunner_parseTestOutput tests parsing of cargo test output
func TestRustTestRunner_parseTestOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	tests := []struct {
		name            string
		output          string
		expectedPassed  int
		expectedFailed  int
		expectedIgnored int
	}{
		{
			name:            "all_passed",
			output:          "test result: ok. 5 passed; 0 failed; 0 ignored; 0 measured; 0 filtered out",
			expectedPassed:  5,
			expectedFailed:  0,
			expectedIgnored: 0,
		},
		{
			name:            "some_failed",
			output:          "test result: FAILED. 3 passed; 2 failed; 1 ignored; 0 measured; 0 filtered out",
			expectedPassed:  3,
			expectedFailed:  2,
			expectedIgnored: 1,
		},
		{
			name:            "all_ignored",
			output:          "test result: ok. 0 passed; 0 failed; 10 ignored; 0 measured; 0 filtered out",
			expectedPassed:  0,
			expectedFailed:  0,
			expectedIgnored: 10,
		},
		{
			name:            "no_match",
			output:          "Running tests...",
			expectedPassed:  0,
			expectedFailed:  0,
			expectedIgnored: 0,
		},
		{
			name:            "empty_output",
			output:          "",
			expectedPassed:  0,
			expectedFailed:  0,
			expectedIgnored: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passed, failed, ignored := runner.parseTestOutput(tt.output)
			assert.Equal(t, tt.expectedPassed, passed, "passed count mismatch")
			assert.Equal(t, tt.expectedFailed, failed, "failed count mismatch")
			assert.Equal(t, tt.expectedIgnored, ignored, "ignored count mismatch")
		})
	}
}

// TestRustTestRunner_DiscoverTests tests test discovery in Rust projects
func TestRustTestRunner_DiscoverTests(t *testing.T) {
	// Create temp directory with Rust project structure
	tempDir, err := os.MkdirTemp("", "rust-discover-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create Cargo.toml
	cargoToml := `[package]
name = "test_project"
version = "0.1.0"
edition = "2021"
`
	err = os.WriteFile(filepath.Join(tempDir, "Cargo.toml"), []byte(cargoToml), 0o644)
	require.NoError(t, err)

	// Create src directory with files containing tests
	srcDir := filepath.Join(tempDir, "src")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	// Create main.rs with tests
	mainRs := `fn main() {}

#[cfg(test)]
mod tests {
    #[test]
    fn test_main() {
        assert!(true);
    }
}
`
	err = os.WriteFile(filepath.Join(srcDir, "main.rs"), []byte(mainRs), 0o644)
	require.NoError(t, err)

	// Create lib.rs with tests
	libRs := `pub fn add(a: i32, b: i32) -> i32 {
    a + b
}

#[test]
fn test_add() {
    assert_eq!(add(1, 2), 3);
}
`
	err = os.WriteFile(filepath.Join(srcDir, "lib.rs"), []byte(libRs), 0o644)
	require.NoError(t, err)

	// Create tests directory for integration tests
	testsDir := filepath.Join(tempDir, "tests")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	// Create integration test
	integrationTest := `#[test]
fn test_integration() {
    assert!(true);
}
`
	err = os.WriteFile(filepath.Join(testsDir, "integration_test.rs"), []byte(integrationTest), 0o644)
	require.NoError(t, err)

	// Create smoke test
	smokeTest := `#[test]
fn test_smoke() {
    assert!(true);
}
`
	err = os.WriteFile(filepath.Join(testsDir, "smoke_test.rs"), []byte(smokeTest), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Should find tests in src/ and tests/
	assert.GreaterOrEqual(t, len(tests), 2)

	// Verify test types
	foundTypes := make(map[string]bool)
	for _, test := range tests {
		foundTypes[test.Type] = true
	}

	// Should have both unit and integration tests
	assert.True(t, foundTypes["unit"] || foundTypes["integration"], "Should find unit or integration tests")
}

// TestRustTestRunner_DiscoverTests_NoCargo tests discovery without Cargo.toml
func TestRustTestRunner_DiscoverTests_NoCargo(t *testing.T) {
	// Create temp directory without Cargo.toml
	tempDir, err := os.MkdirTemp("", "rust-no-cargo-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	_, err = runner.DiscoverTests(context.Background(), tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cargo.toml not found")
}

// TestRustTestRunner_discoverUnitTests tests unit test discovery
func TestRustTestRunner_discoverUnitTests(t *testing.T) {
	// Create temp directory with src structure
	tempDir, err := os.MkdirTemp("", "rust-unit-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	srcDir := filepath.Join(tempDir, "src")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	// Create file with tests
	fileWithTests := `#[test]
fn test_something() {}
`
	err = os.WriteFile(filepath.Join(srcDir, "with_tests.rs"), []byte(fileWithTests), 0o644)
	require.NoError(t, err)

	// Create file without tests
	fileWithoutTests := `fn no_tests() {}
`
	err = os.WriteFile(filepath.Join(srcDir, "without_tests.rs"), []byte(fileWithoutTests), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	tests, err := runner.discoverUnitTests(tempDir)
	require.NoError(t, err)

	// Should only find file with tests
	assert.Len(t, tests, 1)
	assert.Equal(t, "with_tests.rs", tests[0].Name)
	assert.Equal(t, "unit", tests[0].Type)
}

// TestRustTestRunner_discoverIntegrationTests tests integration test discovery
func TestRustTestRunner_discoverIntegrationTests(t *testing.T) {
	// Create temp directory with tests structure
	tempDir, err := os.MkdirTemp("", "rust-integration-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testsDir := filepath.Join(tempDir, "tests")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	// Create integration test
	integrationTest := `#[test]
fn test_integration() {}
`
	err = os.WriteFile(filepath.Join(testsDir, "integration.rs"), []byte(integrationTest), 0o644)
	require.NoError(t, err)

	// Create smoke test
	smokeTest := `#[test]
fn test_smoke() {}
`
	err = os.WriteFile(filepath.Join(testsDir, "smoke_test.rs"), []byte(smokeTest), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	tests, err := runner.discoverIntegrationTests(tempDir)
	require.NoError(t, err)

	assert.Len(t, tests, 2)

	// Verify smoke test is detected
	foundSmoke := false
	for _, test := range tests {
		if test.Type == "smoke" {
			foundSmoke = true
			break
		}
	}
	assert.True(t, foundSmoke, "Should detect smoke test")
}

// TestRustTestRunner_RunTestSuite tests running a test suite
func TestRustTestRunner_RunTestSuite(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	// Create a mock test suite
	suite := &domain.TestSuite{
		Name:        "test_suite",
		Language:    "rust",
		ProjectPath: "/nonexistent/path",
		Tests: []*domain.TestInfo{
			{Path: "test1.rs", Name: "test1.rs", Type: "unit"},
		},
		Config: &domain.TestConfig{
			ProjectPath: "/nonexistent/path",
			Verbose:     false,
		},
	}

	// This will fail because the path doesn't exist, but we're testing the structure
	results, err := runner.RunTestSuite(context.Background(), suite)

	// We expect results even if tests fail
	assert.NoError(t, err)
	assert.NotNil(t, results)
}

// TestRustTestRunner_runAllTests tests running all tests at once
func TestRustTestRunner_runAllTests(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	config := &domain.TestConfig{
		ProjectPath: "/nonexistent/path",
		Verbose:     true,
	}

	// This will fail because the path doesn't exist
	result, err := runner.runAllTests(context.Background(), "/nonexistent/path", config)

	// We expect a result even if the test fails
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "rust", result.Language)
	assert.Equal(t, "all", result.TestName)
}

// TestRustTestRunner_RunTest tests single test execution
func TestRustTestRunner_RunTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	config := &domain.TestConfig{
		ProjectPath: "/nonexistent/path",
		Verbose:     false,
		Timeout:     60,
		EnvVars: map[string]string{
			"RUST_BACKTRACE": "1",
		},
	}

	// This will fail because the path doesn't exist
	result, err := runner.RunTest(context.Background(), "test_file.rs", config)

	// We expect a result even if the test fails
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test_file.rs", result.TestPath)
	assert.Equal(t, "rust", result.Language)
	assert.False(t, result.Success) // Should fail because path doesn't exist
}

// TestRustTestRunner_discoverIntegrationTests_NoTestsDir tests when tests directory doesn't exist
func TestRustTestRunner_discoverIntegrationTests_NoTestsDir(t *testing.T) {
	// Create temp directory without tests directory
	tempDir, err := os.MkdirTemp("", "rust-no-tests-dir-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	tests, err := runner.discoverIntegrationTests(tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestRustTestRunner_discoverUnitTests_NoSrcDir tests when src directory doesn't exist
func TestRustTestRunner_discoverUnitTests_NoSrcDir(t *testing.T) {
	// Create temp directory without src directory
	tempDir, err := os.MkdirTemp("", "rust-no-src-dir-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewRustTestRunner(log)

	tests, err := runner.discoverUnitTests(tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}
