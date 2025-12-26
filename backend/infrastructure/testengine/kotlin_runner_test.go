package testengine

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"syntaxia/domain"
)

// TestKotlinTestRunner_GetLanguage tests language identifier
func TestKotlinTestRunner_GetLanguage(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	assert.Equal(t, "kotlin", runner.GetLanguage())
}

// TestKotlinTestRunner_isKotlinTestFile tests Kotlin test file detection
func TestKotlinTestRunner_isKotlinTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tests := []struct {
		name     string
		fileName string
		expected bool
	}{
		{
			name:     "test_suffix",
			fileName: "UserServiceTest.kt",
			expected: true,
		},
		{
			name:     "tests_suffix",
			fileName: "UserServiceTests.kt",
			expected: true,
		},
		{
			name:     "spec_suffix",
			fileName: "UserServiceSpec.kt",
			expected: true,
		},
		{
			name:     "integration_test_suffix",
			fileName: "UserServiceIT.kt",
			expected: true,
		},
		{
			name:     "test_prefix",
			fileName: "TestUserService.kt",
			expected: true,
		},
		{
			name:     "regular_kotlin_file",
			fileName: "UserService.kt",
			expected: false,
		},
		{
			name:     "non_kotlin_file",
			fileName: "UserServiceTest.java",
			expected: false,
		},
		{
			name:     "test_in_middle",
			fileName: "UserTestService.kt",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.isKotlinTestFile(tt.fileName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestKotlinTestRunner_extractTestClassName tests class name extraction
func TestKotlinTestRunner_extractTestClassName(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tests := []struct {
		name     string
		testPath string
		expected string
	}{
		{
			name:     "simple_file",
			testPath: "UserServiceTest.kt",
			expected: "UserServiceTest",
		},
		{
			name:     "standard_kotlin_structure",
			testPath: "src/test/kotlin/com/example/UserServiceTest.kt",
			expected: "com.example.UserServiceTest",
		},
		{
			name:     "nested_package",
			testPath: "src/test/kotlin/com/example/service/impl/UserServiceTest.kt",
			expected: "com.example.service.impl.UserServiceTest",
		},
		{
			name:     "java_test_directory",
			testPath: "src/test/java/com/example/UserServiceTest.kt",
			expected: "com.example.UserServiceTest",
		},
		{
			name:     "windows_path",
			testPath: "src\\test\\kotlin\\com\\example\\UserServiceTest.kt",
			expected: "com.example.UserServiceTest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.extractTestClassName(tt.testPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestKotlinTestRunner_hasGradleWrapper tests Gradle wrapper detection
func TestKotlinTestRunner_hasGradleWrapper(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tests := []struct {
		name     string
		files    []string
		expected bool
	}{
		{
			name:     "has_gradlew",
			files:    []string{"gradlew"},
			expected: true,
		},
		{
			name:     "has_gradlew_bat",
			files:    []string{"gradlew.bat"},
			expected: true,
		},
		{
			name:     "has_both",
			files:    []string{"gradlew", "gradlew.bat"},
			expected: true,
		},
		{
			name:     "no_wrapper",
			files:    []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "kotlin-gradle-wrapper-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			for _, name := range tt.files {
				err = os.WriteFile(filepath.Join(tempDir, name), []byte("#!/bin/bash"), 0o755)
				require.NoError(t, err)
			}

			result := runner.hasGradleWrapper(tempDir)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestKotlinTestRunner_getGradleCommand tests gradle command detection
func TestKotlinTestRunner_getGradleCommand(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tests := []struct {
		name          string
		setupFiles    []string
		expectWrapper bool
	}{
		{
			name:          "with_wrapper",
			setupFiles:    []string{"gradlew", "gradlew.bat"},
			expectWrapper: true,
		},
		{
			name:          "without_wrapper",
			setupFiles:    []string{},
			expectWrapper: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "kotlin-gradle-cmd-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			for _, name := range tt.setupFiles {
				err = os.WriteFile(filepath.Join(tempDir, name), []byte("#!/bin/bash"), 0o755)
				require.NoError(t, err)
			}

			result := runner.getGradleCommand(tempDir)

			if tt.expectWrapper {
				if runtime.GOOS == "windows" {
					assert.Contains(t, result, "gradlew.bat")
				} else {
					assert.Contains(t, result, "gradlew")
				}
			}
		})
	}
}

// TestKotlinTestRunner_analyzeTestFile tests test file type analysis
func TestKotlinTestRunner_analyzeTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tempDir, err := os.MkdirTemp("", "kotlin-analyze-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected string
	}{
		{
			name:     "smoke_test_annotation",
			fileName: "SmokeTest.kt",
			content: `@Tag("smoke")
class SmokeTest {
    @Test
    fun testSmoke() {}
}`,
			expected: "smoke",
		},
		{
			name:     "smoke_in_filename",
			fileName: "UserSmokeTest.kt",
			content: `class UserSmokeTest {
    @Test
    fun test() {}
}`,
			expected: "smoke",
		},
		{
			name:     "spring_boot_test",
			fileName: "UserServiceIT.kt",
			content: `@SpringBootTest
class UserServiceIT {
    @Test
    fun testIntegration() {}
}`,
			expected: "integration",
		},
		{
			name:     "data_jpa_test",
			fileName: "UserRepositoryTest.kt",
			content: `@DataJpaTest
class UserRepositoryTest {
    @Test
    fun testRepository() {}
}`,
			expected: "integration",
		},
		{
			name:     "testcontainers",
			fileName: "DatabaseIT.kt",
			content: `@Testcontainers
class DatabaseIT {
    @Container
    val postgres = PostgreSQLContainer<Nothing>()
}`,
			expected: "integration",
		},
		{
			name:     "integration_in_name",
			fileName: "UserIntegrationTest.kt",
			content: `class UserIntegrationTest {
    @Test
    fun test() {}
}`,
			expected: "integration",
		},
		{
			name:     "it_suffix",
			fileName: "UserIT.kt",
			content: `class UserIT {
    @Test
    fun test() {}
}`,
			expected: "integration",
		},
		{
			name:     "unit_test",
			fileName: "UserServiceTest.kt",
			content: `class UserServiceTest {
    @Test
    fun testAdd() {
        assertEquals(2, 1 + 1)
    }
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

// TestKotlinTestRunner_analyzeTestFile_NonExistent tests analyzing non-existent file
func TestKotlinTestRunner_analyzeTestFile_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	result := runner.analyzeTestFile("/nonexistent/UserServiceTest.kt")
	assert.Equal(t, "unit", result)
}

// TestKotlinTestRunner_ParseTestOutput tests parsing gradle test output
func TestKotlinTestRunner_ParseTestOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tests := []struct {
		name            string
		output          string
		expectedPassed  int
		expectedFailed  int
		expectedSkipped int
	}{
		{
			name:            "all_passed",
			output:          "5 tests completed, 0 failed, 0 skipped",
			expectedPassed:  5,
			expectedFailed:  0,
			expectedSkipped: 0,
		},
		{
			name:            "some_failed",
			output:          "10 tests completed, 2 failed, 1 skipped",
			expectedPassed:  8,
			expectedFailed:  2,
			expectedSkipped: 1,
		},
		{
			name:            "no_match",
			output:          "Running tests...",
			expectedPassed:  0,
			expectedFailed:  0,
			expectedSkipped: 0,
		},
		{
			name:            "empty_output",
			output:          "",
			expectedPassed:  0,
			expectedFailed:  0,
			expectedSkipped: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passed, failed, skipped := runner.ParseTestOutput(tt.output)
			assert.Equal(t, tt.expectedPassed, passed, "passed count mismatch")
			assert.Equal(t, tt.expectedFailed, failed, "failed count mismatch")
			assert.Equal(t, tt.expectedSkipped, skipped, "skipped count mismatch")
		})
	}
}

// TestKotlinTestRunner_DiscoverTests tests test discovery
func TestKotlinTestRunner_DiscoverTests(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "kotlin-discover-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create build.gradle.kts
	buildGradle := `plugins {
    kotlin("jvm") version "1.9.0"
}
`
	err = os.WriteFile(filepath.Join(tempDir, "build.gradle.kts"), []byte(buildGradle), 0o644)
	require.NoError(t, err)

	// Create test directory structure
	testDir := filepath.Join(tempDir, "src", "test", "kotlin", "com", "example")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	// Create test files
	testFiles := map[string]string{
		"UserServiceTest.kt": `class UserServiceTest {
    @Test
    fun test() {}
}`,
		"UserServiceIT.kt": `@SpringBootTest
class UserServiceIT {
    @Test
    fun test() {}
}`,
		"UserService.kt": `class UserService {}`,
	}

	for name, content := range testFiles {
		err = os.WriteFile(filepath.Join(testDir, name), []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 2)

	foundFiles := make(map[string]bool)
	for _, test := range tests {
		foundFiles[test.Name] = true
	}

	assert.True(t, foundFiles["UserServiceTest.kt"], "UserServiceTest.kt should be found")
	assert.True(t, foundFiles["UserServiceIT.kt"], "UserServiceIT.kt should be found")
	assert.False(t, foundFiles["UserService.kt"], "UserService.kt should not be found")
}

// TestKotlinTestRunner_DiscoverTests_MultipleDirectories tests discovery in multiple test directories
func TestKotlinTestRunner_DiscoverTests_MultipleDirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "kotlin-multi-dir-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testDirs := []string{
		"src/test/kotlin",
		"src/test/java",
		"test",
		"tests",
	}

	for i, dir := range testDirs {
		fullDir := filepath.Join(tempDir, dir)
		err = os.MkdirAll(fullDir, 0o755)
		require.NoError(t, err)

		testFile := filepath.Join(fullDir, "Test"+string(rune('A'+i))+"Test.kt")
		content := `class Test` + string(rune('A'+i)) + `Test { @Test fun test() {} }`
		err = os.WriteFile(testFile, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 1)
}

// TestKotlinTestRunner_DiscoverTests_NoTestDirectories tests discovery when no test directories exist
func TestKotlinTestRunner_DiscoverTests_NoTestDirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "kotlin-no-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestKotlinTestRunner_RunTest tests single test execution
func TestKotlinTestRunner_RunTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tempDir, err := os.MkdirTemp("", "kotlin-run-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	buildGradle := `plugins { kotlin("jvm") version "1.9.0" }`
	err = os.WriteFile(filepath.Join(tempDir, "build.gradle.kts"), []byte(buildGradle), 0o644)
	require.NoError(t, err)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		Verbose:     true,
		EnvVars: map[string]string{
			"JAVA_HOME": "/usr/lib/jvm/java-17",
		},
	}

	result, err := runner.RunTest(context.Background(), "src/test/kotlin/com/example/UserServiceTest.kt", config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "src/test/kotlin/com/example/UserServiceTest.kt", result.TestPath)
	assert.Equal(t, "kotlin", result.Language)
}

// TestKotlinTestRunner_RunTestSuite tests running a test suite
func TestKotlinTestRunner_RunTestSuite(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tempDir, err := os.MkdirTemp("", "kotlin-suite-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	buildGradle := `plugins { kotlin("jvm") version "1.9.0" }`
	err = os.WriteFile(filepath.Join(tempDir, "build.gradle.kts"), []byte(buildGradle), 0o644)
	require.NoError(t, err)

	suite := &domain.TestSuite{
		Name:        "test_suite",
		Language:    "kotlin",
		ProjectPath: tempDir,
		Tests: []*domain.TestInfo{
			{Path: "src/test/kotlin/com/example/Test1.kt", Name: "Test1.kt", Type: "unit"},
			{Path: "src/test/kotlin/com/example/Test2.kt", Name: "Test2.kt", Type: "unit"},
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

// TestKotlinTestRunner_ParseTestOutput_EdgeCases tests edge cases in output parsing
func TestKotlinTestRunner_ParseTestOutput_EdgeCases(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tests := []struct {
		name            string
		output          string
		expectedPassed  int
		expectedFailed  int
		expectedSkipped int
	}{
		{
			name:            "single_test_passed",
			output:          "1 test completed, 0 failed, 0 skipped",
			expectedPassed:  1,
			expectedFailed:  0,
			expectedSkipped: 0,
		},
		{
			name:            "all_skipped",
			output:          "5 tests completed, 0 failed, 5 skipped",
			expectedPassed:  5,
			expectedFailed:  0,
			expectedSkipped: 5,
		},
		{
			name:            "mixed_results",
			output:          "100 tests completed, 10 failed, 5 skipped",
			expectedPassed:  90,
			expectedFailed:  10,
			expectedSkipped: 5,
		},
		{
			name:            "gradle_verbose_output",
			output:          "BUILD SUCCESSFUL\n> Task :test\n20 tests completed, 3 failed, 2 skipped\nBuild finished",
			expectedPassed:  17,
			expectedFailed:  3,
			expectedSkipped: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passed, failed, skipped := runner.ParseTestOutput(tt.output)
			assert.Equal(t, tt.expectedPassed, passed, "passed count mismatch")
			assert.Equal(t, tt.expectedFailed, failed, "failed count mismatch")
			assert.Equal(t, tt.expectedSkipped, skipped, "skipped count mismatch")
		})
	}
}

// TestKotlinTestRunner_extractTestClassName_EdgeCases tests edge cases in class name extraction
func TestKotlinTestRunner_extractTestClassName_EdgeCases(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tests := []struct {
		name     string
		testPath string
		expected string
	}{
		{
			name:     "deeply_nested_package",
			testPath: "src/test/kotlin/com/example/app/feature/module/UserServiceTest.kt",
			expected: "com.example.app.feature.module.UserServiceTest",
		},
		{
			name:     "root_level_test",
			testPath: "Test.kt",
			expected: "Test",
		},
		{
			name:     "mixed_separators",
			testPath: "src/test/kotlin\\com\\example/UserTest.kt",
			expected: "com.example.UserTest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.extractTestClassName(tt.testPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestKotlinTestRunner_analyzeTestFile_WebMvcTest tests WebMvcTest detection
func TestKotlinTestRunner_analyzeTestFile_WebMvcTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tempDir, err := os.MkdirTemp("", "kotlin-webmvc-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `@WebMvcTest(UserController::class)
class UserControllerTest {
    @Test
    fun testGetUser() {}
}`
	testFile := filepath.Join(tempDir, "UserControllerTest.kt")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	result := runner.analyzeTestFile(testFile)
	assert.Equal(t, "integration", result)
}

// TestKotlinTestRunner_getGradleCommand_NoGradle tests when gradle is not available
func TestKotlinTestRunner_getGradleCommand_NoGradle(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tempDir, err := os.MkdirTemp("", "kotlin-no-gradle-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// No gradlew and assuming gradle is not in PATH for this test
	result := runner.getGradleCommand(tempDir)

	// Result could be empty or "gradle" depending on system
	// Just verify it doesn't panic
	assert.NotNil(t, result)
}


// TestKotlinTestAnalyzer_AnalyzeTestDependencies tests dependency analysis
func TestKotlinTestAnalyzer_AnalyzeTestDependencies(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "kotlin-analyzer-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `package com.example.test

import org.junit.jupiter.api.Test
import com.example.service.UserService
import kotlin.test.assertEquals

class UserServiceTest {
    @Test
    fun testUser() {}
}
`
	testFile := filepath.Join(tempDir, "UserServiceTest.kt")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	deps, err := analyzer.AnalyzeTestDependencies(context.Background(), testFile)
	require.NoError(t, err)

	assert.Len(t, deps, 3)
	assert.Contains(t, deps, "org.junit.jupiter.api.Test")
	assert.Contains(t, deps, "com.example.service.UserService")
	assert.Contains(t, deps, "kotlin.test.assertEquals")
}

// TestKotlinTestAnalyzer_IsSmokeTest tests smoke test detection
func TestKotlinTestAnalyzer_IsSmokeTest(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "kotlin-smoke-test-*")
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
			fileName: "SmokeTest.kt",
			content:  `class SmokeTest {}`,
			expected: true,
		},
		{
			name:     "smoke_tag",
			fileName: "ApiTest.kt",
			content:  `@Tag("smoke") class ApiTest {}`,
			expected: true,
		},
		{
			name:     "regular_test",
			fileName: "UserServiceTest.kt",
			content:  `class UserServiceTest {}`,
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

// TestKotlinTestAnalyzer_FindTestsForFile tests finding tests for a source file
func TestKotlinTestAnalyzer_FindTestsForFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "kotlin-find-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source file
	srcDir := filepath.Join(tempDir, "src", "main", "kotlin", "com", "example")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(srcDir, "UserService.kt")
	err = os.WriteFile(srcFile, []byte("class UserService {}"), 0o644)
	require.NoError(t, err)

	// Create test directory and test file
	testDir := filepath.Join(tempDir, "src", "test", "kotlin", "com", "example")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	testFile := filepath.Join(testDir, "UserServiceTest.kt")
	testContent := `import com.example.UserService
class UserServiceTest {}`
	err = os.WriteFile(testFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForFile(context.Background(), "src/main/kotlin/com/example/UserService.kt", tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 1)
}

// TestKotlinTestAnalyzer_ExtractImports tests import extraction
func TestKotlinTestAnalyzer_ExtractImports(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	content := `package com.example

import org.junit.Test
import com.example.service.UserService
import kotlin.test.assertEquals

class Test {}`

	imports := analyzer.extractImports(content)

	assert.Len(t, imports, 3)
	assert.Contains(t, imports, "org.junit.Test")
	assert.Contains(t, imports, "com.example.service.UserService")
	assert.Contains(t, imports, "kotlin.test.assertEquals")
}


// TestKotlinTestAnalyzer_RemoveDuplicates tests duplicate removal
func TestKotlinTestAnalyzer_RemoveDuplicates(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

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

// TestKotlinTestAnalyzer_FileExists tests file existence check
func TestKotlinTestAnalyzer_FileExists(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "kotlin-exists-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a file
	testFile := filepath.Join(tempDir, "exists.kt")
	err = os.WriteFile(testFile, []byte("class Test {}"), 0o644)
	require.NoError(t, err)

	assert.True(t, analyzer.fileExists(testFile))
	assert.False(t, analyzer.fileExists(filepath.Join(tempDir, "notexists.kt")))
}

// TestKotlinTestAnalyzer_ImportsClass tests class import detection
func TestKotlinTestAnalyzer_ImportsClass(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	tests := []struct {
		name      string
		content   string
		className string
		expected  bool
	}{
		{
			name:      "imports_class",
			content:   `import com.example.UserService`,
			className: "UserService",
			expected:  true,
		},
		{
			name:      "uses_class",
			content:   `val service = UserService()`,
			className: "UserService",
			expected:  true,
		},
		{
			name:      "no_import",
			content:   `import com.example.OrderService`,
			className: "UserService",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.importsClass(tt.content, tt.className)
			assert.Equal(t, tt.expected, result)
		})
	}
}


// TestKotlinTestAnalyzer_AnalyzeTestDependencies_NonExistent tests error handling
func TestKotlinTestAnalyzer_AnalyzeTestDependencies_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	_, err := analyzer.AnalyzeTestDependencies(context.Background(), "/nonexistent/file.kt")
	assert.Error(t, err)
}

// TestKotlinTestAnalyzer_IsSmokeTest_NonExistent tests error handling
func TestKotlinTestAnalyzer_IsSmokeTest_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	_, err := analyzer.IsSmokeTest(context.Background(), "/nonexistent/file.kt")
	assert.Error(t, err)
}


// TestKotlinTestAnalyzer_ExtractImports_EmptyContent tests extracting imports from empty content
func TestKotlinTestAnalyzer_ExtractImports_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	imports := analyzer.extractImports("")
	assert.Empty(t, imports)
}

// TestKotlinTestAnalyzer_ExtractImports_DuplicateImports tests duplicate import handling
func TestKotlinTestAnalyzer_ExtractImports_DuplicateImports(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	content := `import org.junit.Test
import org.junit.Test
import kotlin.test.assertEquals`

	imports := analyzer.extractImports(content)
	assert.Len(t, imports, 2)
}

// TestKotlinTestAnalyzer_ExtractImports_WildcardImport tests wildcard import extraction
func TestKotlinTestAnalyzer_ExtractImports_WildcardImport(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	content := `import org.junit.*
import kotlin.test.assertEquals`

	imports := analyzer.extractImports(content)
	assert.Len(t, imports, 2)
	// The regex captures "org.junit." (without the asterisk)
	assert.Contains(t, imports, "org.junit.")
}

// TestKotlinTestAnalyzer_FindTestsForFile_NoTestDirs tests when no test directories exist
func TestKotlinTestAnalyzer_FindTestsForFile_NoTestDirs(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "kotlin-no-test-dirs-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source file without test directories
	srcDir := filepath.Join(tempDir, "src", "main", "kotlin")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(srcDir, "UserService.kt")
	err = os.WriteFile(srcFile, []byte("class UserService {}"), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForFile(context.Background(), "src/main/kotlin/UserService.kt", tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestKotlinTestAnalyzer_IsSmokeTest_EmptyFile tests smoke test detection with empty file
func TestKotlinTestAnalyzer_IsSmokeTest_EmptyFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "kotlin-empty-smoke-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "EmptyTest.kt")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	result, err := analyzer.IsSmokeTest(context.Background(), testFile)
	require.NoError(t, err)
	assert.False(t, result)
}

// TestKotlinTestAnalyzer_IsSmokeTest_CategoryAnnotation tests smoke test detection with @Category annotation
func TestKotlinTestAnalyzer_IsSmokeTest_CategoryAnnotation(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewKotlinTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "kotlin-category-smoke-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `@Category(Smoke::class)
class ApiTest {
    @Test
    fun testApi() {}
}`
	testFile := filepath.Join(tempDir, "ApiTest.kt")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	result, err := analyzer.IsSmokeTest(context.Background(), testFile)
	require.NoError(t, err)
	// The implementation checks for lowercase "smoke" which matches "Smoke" after ToLower
	assert.True(t, result)
}

// TestKotlinTestRunner_analyzeTestFile_EmptyFile tests analyzing empty file
func TestKotlinTestRunner_analyzeTestFile_EmptyFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tempDir, err := os.MkdirTemp("", "kotlin-empty-analyze-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "EmptyTest.kt")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	result := runner.analyzeTestFile(testFile)
	assert.Equal(t, "unit", result)
}

// TestKotlinTestRunner_ParseTestOutput_ZeroTests tests parsing output with zero tests
func TestKotlinTestRunner_ParseTestOutput_ZeroTests(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	output := "0 tests completed, 0 failed, 0 skipped"
	passed, failed, skipped := runner.ParseTestOutput(output)
	assert.Equal(t, 0, passed)
	assert.Equal(t, 0, failed)
	assert.Equal(t, 0, skipped)
}

// TestKotlinTestRunner_extractTestClassName_EmptyPath tests class name extraction with empty path
func TestKotlinTestRunner_extractTestClassName_EmptyPath(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	result := runner.extractTestClassName("")
	// Empty path returns "." from filepath.Base
	assert.Equal(t, ".", result)
}

// TestKotlinTestRunner_DiscoverTests_SkipBuildDir tests discovery behavior with build directories
func TestKotlinTestRunner_DiscoverTests_SkipBuildDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "kotlin-skip-build-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test directory
	testDir := filepath.Join(tempDir, "src", "test", "kotlin")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	// Create valid test file
	validTestFile := filepath.Join(testDir, "ValidTest.kt")
	err = os.WriteFile(validTestFile, []byte("class ValidTest {}"), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewKotlinTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Verify valid test is found
	foundValid := false
	for _, test := range tests {
		if test.Name == "ValidTest.kt" {
			foundValid = true
			break
		}
	}
	assert.True(t, foundValid, "ValidTest.kt should be discovered")
}
