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

// TestCSharpTestRunner_GetLanguage tests language identifier
func TestCSharpTestRunner_GetLanguage(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	assert.Equal(t, "csharp", runner.GetLanguage())
}

// TestCSharpTestRunner_isCSharpTestFile tests C# test file detection
func TestCSharpTestRunner_isCSharpTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tempDir, err := os.MkdirTemp("", "csharp-test-file-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected bool
	}{
		{
			name:     "test_suffix",
			fileName: "UserServiceTest.cs",
			content:  `public class UserServiceTest {}`,
			expected: true,
		},
		{
			name:     "tests_suffix",
			fileName: "UserServiceTests.cs",
			content:  `public class UserServiceTests {}`,
			expected: true,
		},
		{
			name:     "spec_suffix",
			fileName: "UserServiceSpec.cs",
			content:  `public class UserServiceSpec {}`,
			expected: true,
		},
		{
			name:     "xunit_fact",
			fileName: "UserService.cs",
			content:  `public class UserService { [Fact] public void Test() {} }`,
			expected: true,
		},
		{
			name:     "xunit_theory",
			fileName: "Calculator.cs",
			content:  `public class Calculator { [Theory] public void Test() {} }`,
			expected: true,
		},
		{
			name:     "nunit_test",
			fileName: "Service.cs",
			content:  `public class Service { [Test] public void Test() {} }`,
			expected: true,
		},
		{
			name:     "mstest_testmethod",
			fileName: "Handler.cs",
			content:  `public class Handler { [TestMethod] public void Test() {} }`,
			expected: true,
		},
		{
			name:     "regular_csharp_file",
			fileName: "UserService.cs",
			content:  `public class UserService { public void DoWork() {} }`,
			expected: false,
		},
		{
			name:     "non_csharp_file",
			fileName: "UserServiceTest.java",
			content:  `public class UserServiceTest {}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := runner.isCSharpTestFile(tt.fileName, testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCSharpTestRunner_extractTestFilter tests test filter extraction
func TestCSharpTestRunner_extractTestFilter(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tests := []struct {
		name     string
		testPath string
		expected string
	}{
		{
			name:     "simple_file",
			testPath: "UserServiceTest.cs",
			expected: "FullyQualifiedName~UserServiceTest",
		},
		{
			name:     "nested_path",
			testPath: "Tests/Unit/UserServiceTest.cs",
			expected: "FullyQualifiedName~UserServiceTest",
		},
		{
			name:     "windows_path",
			testPath: "Tests\\Unit\\UserServiceTest.cs",
			expected: "FullyQualifiedName~UserServiceTest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.extractTestFilter(tt.testPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCSharpTestRunner_detectTestFramework tests framework detection
func TestCSharpTestRunner_detectTestFramework(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tempDir, err := os.MkdirTemp("", "csharp-framework-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected string
	}{
		{
			name:     "xunit_fact",
			fileName: "XunitTest.cs",
			content:  `public class XunitTest { [Fact] public void Test() {} }`,
			expected: "xunit",
		},
		{
			name:     "xunit_theory",
			fileName: "XunitTheory.cs",
			content:  `public class XunitTheory { [Theory] public void Test() {} }`,
			expected: "xunit",
		},
		{
			name:     "nunit_test",
			fileName: "NunitTest.cs",
			content:  `public class NunitTest { [Test] public void Test() {} }`,
			expected: "nunit",
		},
		{
			name:     "nunit_testcase",
			fileName: "NunitTestCase.cs",
			content:  `public class NunitTestCase { [TestCase] public void Test() {} }`,
			expected: "nunit",
		},
		{
			name:     "mstest_testmethod",
			fileName: "MsTest.cs",
			content:  `public class MsTest { [TestMethod] public void Test() {} }`,
			expected: "mstest",
		},
		{
			name:     "mstest_testclass",
			fileName: "MsTestClass.cs",
			content:  `[TestClass] public class MsTestClass { public void Test() {} }`,
			expected: "mstest",
		},
		{
			name:     "unknown_framework",
			fileName: "Unknown.cs",
			content:  `public class Unknown { public void Test() {} }`,
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := runner.detectTestFramework(testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCSharpTestRunner_analyzeTestFile tests test file type analysis
func TestCSharpTestRunner_analyzeTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tempDir, err := os.MkdirTemp("", "csharp-analyze-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected string
	}{
		{
			name:     "smoke_in_content",
			fileName: "SmokeTest.cs",
			content:  `// Smoke test for API public class SmokeTest {}`,
			expected: "smoke",
		},
		{
			name:     "smoke_in_filename",
			fileName: "UserSmokeTests.cs",
			content:  `public class UserSmokeTests {}`,
			expected: "smoke",
		},
		{
			name:     "integration_in_content",
			fileName: "ApiTest.cs",
			content:  `// Integration test public class ApiTest {}`,
			expected: "integration",
		},
		{
			name:     "integration_in_filename",
			fileName: "UserIntegrationTests.cs",
			content:  `public class UserIntegrationTests {}`,
			expected: "integration",
		},
		{
			name:     "webapplicationfactory",
			fileName: "ApiTests.cs",
			content:  `public class ApiTests : IClassFixture<WebApplicationFactory<Program>> {}`,
			expected: "integration",
		},
		{
			name:     "testserver",
			fileName: "ServerTests.cs",
			content:  `public class ServerTests { TestServer server; }`,
			expected: "integration",
		},
		{
			name:     "httpclient",
			fileName: "ClientTests.cs",
			content:  `public class ClientTests { HttpClient client; }`,
			expected: "integration",
		},
		{
			name:     "unit_test",
			fileName: "UserServiceTests.cs",
			content:  `public class UserServiceTests { [Fact] public void Test() {} }`,
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

// TestCSharpTestRunner_analyzeTestFile_NonExistent tests analyzing non-existent file
func TestCSharpTestRunner_analyzeTestFile_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	result := runner.analyzeTestFile("/nonexistent/UserServiceTest.cs")
	assert.Equal(t, "unit", result)
}

// TestCSharpTestRunner_extractClassName tests class name extraction
func TestCSharpTestRunner_extractClassName(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tempDir, err := os.MkdirTemp("", "csharp-classname-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected string
	}{
		{
			name:     "simple_class",
			fileName: "UserService.cs",
			content:  `public class UserService { }`,
			expected: "UserService",
		},
		{
			name:     "class_with_inheritance",
			fileName: "UserServiceTests.cs",
			content:  `public class UserServiceTests : TestBase { }`,
			expected: "UserServiceTests",
		},
		{
			name:     "no_class_declaration",
			fileName: "NoClass.cs",
			content:  `namespace MyApp { }`,
			expected: "NoClass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := runner.extractClassName(testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCSharpTestRunner_ParseTestOutput tests parsing dotnet test output
func TestCSharpTestRunner_ParseTestOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tests := []struct {
		name            string
		output          string
		expectedPassed  int
		expectedFailed  int
		expectedSkipped int
		expectedTotal   int
		expectedSuccess bool
	}{
		{
			name: "all_passed",
			output: `Test run for MyProject.Tests.dll
Passed: 10
Failed: 0
Skipped: 0
Total: 10`,
			expectedPassed:  10,
			expectedFailed:  0,
			expectedSkipped: 0,
			expectedTotal:   10,
			expectedSuccess: true,
		},
		{
			name: "some_failed",
			output: `Test run for MyProject.Tests.dll
Passed: 8
Failed: 2
Skipped: 1
Total: 11`,
			expectedPassed:  8,
			expectedFailed:  2,
			expectedSkipped: 1,
			expectedTotal:   11,
			expectedSuccess: false,
		},
		{
			name:            "no_match",
			output:          "Running tests...",
			expectedPassed:  0,
			expectedFailed:  0,
			expectedSkipped: 0,
			expectedTotal:   0,
			expectedSuccess: true,
		},
		{
			name:            "empty_output",
			output:          "",
			expectedPassed:  0,
			expectedFailed:  0,
			expectedSkipped: 0,
			expectedTotal:   0,
			expectedSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.ParseTestOutput(tt.output)
			assert.Equal(t, tt.expectedPassed, result.PassedTests, "passed count mismatch")
			assert.Equal(t, tt.expectedFailed, result.FailedTests, "failed count mismatch")
			assert.Equal(t, tt.expectedSkipped, result.SkippedTests, "skipped count mismatch")
			assert.Equal(t, tt.expectedTotal, result.TotalTests, "total count mismatch")
			assert.Equal(t, tt.expectedSuccess, result.Success, "success mismatch")
		})
	}
}

// TestCSharpTestRunner_DiscoverTests tests test discovery
func TestCSharpTestRunner_DiscoverTests(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "csharp-discover-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test project structure
	testDir := filepath.Join(tempDir, "MyProject.Tests")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	// Create .csproj file
	csproj := `<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
  </PropertyGroup>
</Project>`
	err = os.WriteFile(filepath.Join(testDir, "MyProject.Tests.csproj"), []byte(csproj), 0o644)
	require.NoError(t, err)

	// Create test files
	testFiles := map[string]string{
		"UserServiceTests.cs": `public class UserServiceTests { [Fact] public void Test() {} }`,
		"UserServiceIT.cs":    `public class UserServiceIT { [Fact] public void Integration() {} }`,
		"UserService.cs":      `public class UserService { public void DoWork() {} }`,
	}

	for name, content := range testFiles {
		err = os.WriteFile(filepath.Join(testDir, name), []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 2)

	foundFiles := make(map[string]bool)
	for _, test := range tests {
		foundFiles[test.Name] = true
	}

	assert.True(t, foundFiles["UserServiceTests.cs"], "UserServiceTests.cs should be found")
	assert.True(t, foundFiles["UserServiceIT.cs"], "UserServiceIT.cs should be found")
	assert.False(t, foundFiles["UserService.cs"], "UserService.cs should not be found")
}

// TestCSharpTestRunner_DiscoverTests_SkipDirectories tests that bin/obj are skipped
func TestCSharpTestRunner_DiscoverTests_SkipDirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "csharp-skip-dirs-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create directories that should be skipped
	skipDirs := []string{"bin", "obj", ".git"}
	for _, dir := range skipDirs {
		fullDir := filepath.Join(tempDir, dir)
		err = os.MkdirAll(fullDir, 0o755)
		require.NoError(t, err)

		testFile := filepath.Join(fullDir, "ShouldNotFind.cs")
		err = os.WriteFile(testFile, []byte(`[Fact] public void Test() {}`), 0o644)
		require.NoError(t, err)
	}

	// Create a valid test file
	validTest := filepath.Join(tempDir, "ValidTest.cs")
	err = os.WriteFile(validTest, []byte(`public class ValidTest { [Fact] public void Test() {} }`), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Should only find ValidTest.cs
	for _, test := range tests {
		assert.NotContains(t, test.Path, "bin")
		assert.NotContains(t, test.Path, "obj")
		assert.NotContains(t, test.Path, ".git")
	}
}

// TestCSharpTestRunner_RunTest tests single test execution
func TestCSharpTestRunner_RunTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tempDir, err := os.MkdirTemp("", "csharp-run-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	csproj := `<Project Sdk="Microsoft.NET.Sdk"></Project>`
	err = os.WriteFile(filepath.Join(tempDir, "Test.csproj"), []byte(csproj), 0o644)
	require.NoError(t, err)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		Verbose:     true,
		EnvVars: map[string]string{
			"DOTNET_CLI_TELEMETRY_OPTOUT": "1",
		},
	}

	result, err := runner.RunTest(context.Background(), "Tests/UserServiceTests.cs", config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Tests/UserServiceTests.cs", result.TestPath)
	assert.Equal(t, "csharp", result.Language)
}

// TestCSharpTestRunner_RunTestSuite tests running a test suite
func TestCSharpTestRunner_RunTestSuite(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tempDir, err := os.MkdirTemp("", "csharp-suite-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	csproj := `<Project Sdk="Microsoft.NET.Sdk"></Project>`
	err = os.WriteFile(filepath.Join(tempDir, "Test.csproj"), []byte(csproj), 0o644)
	require.NoError(t, err)

	suite := &domain.TestSuite{
		Name:        "test_suite",
		Language:    "csharp",
		ProjectPath: tempDir,
		Tests: []*domain.TestInfo{
			{Path: "Tests/Test1.cs", Name: "Test1.cs", Type: "unit"},
			{Path: "Tests/Test2.cs", Name: "Test2.cs", Type: "unit"},
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

// TestCSharpTestRunner_ParseTestOutput_FailedTests tests parsing failed test names
func TestCSharpTestRunner_ParseTestOutput_FailedTests(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	output := `Test run for MyProject.Tests.dll
Failed MyProject.Tests.UserServiceTests.TestMethod1
Failed MyProject.Tests.UserServiceTests.TestMethod2
Passed: 8
Failed: 2
Skipped: 0
Total: 10`

	result := runner.ParseTestOutput(output)

	assert.False(t, result.Success)
	assert.Equal(t, 2, result.FailedTests)
	assert.GreaterOrEqual(t, len(result.FailedTestPaths), 1)
}

// TestCSharpTestRunner_ParseTestOutput_SuccessRate tests success rate calculation
func TestCSharpTestRunner_ParseTestOutput_SuccessRate(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tests := []struct {
		name                string
		output              string
		expectedSuccessRate float64
	}{
		{
			name: "all_passed",
			output: `Passed: 10
Failed: 0
Skipped: 0
Total: 10`,
			expectedSuccessRate: 100.0,
		},
		{
			name: "half_passed",
			output: `Passed: 5
Failed: 5
Skipped: 0
Total: 10`,
			expectedSuccessRate: 50.0,
		},
		{
			name: "quarter_passed",
			output: `Passed: 25
Failed: 75
Skipped: 0
Total: 100`,
			expectedSuccessRate: 25.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.ParseTestOutput(tt.output)
			assert.Equal(t, tt.expectedSuccessRate, result.SuccessRate)
		})
	}
}

// TestCSharpTestRunner_isCSharpTestFile_ContentBased tests content-based test file detection
func TestCSharpTestRunner_isCSharpTestFile_ContentBased(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tempDir, err := os.MkdirTemp("", "csharp-content-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected bool
	}{
		{
			name:     "nunit_test_attribute",
			fileName: "Calculator.cs",
			content:  `public class Calculator { [Test] public void Add() {} }`,
			expected: true,
		},
		{
			name:     "interface_file",
			fileName: "IUserService.cs",
			content:  `public interface IUserService { void DoWork(); }`,
			expected: false,
		},
		{
			name:     "model_file",
			fileName: "User.cs",
			content:  `public class User { public string Name { get; set; } }`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := runner.isCSharpTestFile(tt.fileName, testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCSharpTestRunner_extractClassName_NoClass tests class name extraction when no class found
func TestCSharpTestRunner_extractClassName_NoClass(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tempDir, err := os.MkdirTemp("", "csharp-noclass-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// File with no class declaration
	content := `namespace MyApp
{
    public interface IService { }
}`
	testFile := filepath.Join(tempDir, "IService.cs")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	result := runner.extractClassName(testFile)
	assert.Equal(t, "IService", result) // Falls back to filename
}

// TestCSharpTestRunner_extractClassName_NonExistent tests class name extraction for non-existent file
func TestCSharpTestRunner_extractClassName_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	result := runner.extractClassName("/nonexistent/UserService.cs")
	assert.Equal(t, "UserService", result) // Falls back to filename
}

// TestCSharpTestRunner_detectTestFramework_NonExistent tests framework detection for non-existent file
func TestCSharpTestRunner_detectTestFramework_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	result := runner.detectTestFramework("/nonexistent/UserServiceTest.cs")
	assert.Equal(t, "unknown", result)
}

// TestCSharpTestRunner_DiscoverTests_NestedDirectories tests discovery in nested directories
func TestCSharpTestRunner_DiscoverTests_NestedDirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "csharp-nested-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create nested test directory structure
	nestedDir := filepath.Join(tempDir, "Tests", "Unit", "Services")
	err = os.MkdirAll(nestedDir, 0o755)
	require.NoError(t, err)

	// Create test file in nested directory
	testFile := filepath.Join(nestedDir, "UserServiceTests.cs")
	content := `public class UserServiceTests { [Fact] public void Test() {} }`
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 1)
	assert.Contains(t, tests[0].Path, "UserServiceTests.cs")
}


// TestCSharpTestAnalyzer_AnalyzeTestDependencies tests dependency analysis
func TestCSharpTestAnalyzer_AnalyzeTestDependencies(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "csharp-analyzer-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `using System;
using System.Collections.Generic;
using Xunit;
using MyApp.Services;

namespace MyApp.Tests
{
    public class UserServiceTests
    {
        [Fact]
        public void Test() {}
    }
}
`
	testFile := filepath.Join(tempDir, "UserServiceTests.cs")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	deps, err := analyzer.AnalyzeTestDependencies(context.Background(), testFile)
	require.NoError(t, err)

	assert.Len(t, deps, 4)
	assert.Contains(t, deps, "System")
	assert.Contains(t, deps, "System.Collections.Generic")
	assert.Contains(t, deps, "Xunit")
	assert.Contains(t, deps, "MyApp.Services")
}

// TestCSharpTestAnalyzer_IsSmokeTest tests smoke test detection
func TestCSharpTestAnalyzer_IsSmokeTest(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "csharp-smoke-test-*")
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
			fileName: "SmokeTests.cs",
			content:  `public class SmokeTests {}`,
			expected: true,
		},
		{
			name:     "smoke_category",
			fileName: "ApiTests.cs",
			content:  `[Category("smoke")] public class ApiTests {}`,
			expected: true,
		},
		{
			name:     "regular_test",
			fileName: "UserServiceTests.cs",
			content:  `public class UserServiceTests {}`,
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

// TestCSharpTestAnalyzer_AnalyzeTestFile tests test file analysis
func TestCSharpTestAnalyzer_AnalyzeTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "csharp-analyze-file-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `using Xunit;

public class UserServiceTests
{
    [Fact]
    public void TestCreate() {}
    
    [Fact]
    public void TestUpdate() {}
    
    [Theory]
    [InlineData(1)]
    public void TestDelete(int id) {}
}
`
	testFile := filepath.Join(tempDir, "UserServiceTests.cs")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	info, err := analyzer.AnalyzeTestFile(context.Background(), testFile)
	require.NoError(t, err)

	assert.Equal(t, "UserServiceTests.cs", info.Name)
	assert.Equal(t, "unit", info.Type)
	assert.Equal(t, "xunit", info.Metadata["framework"])
	assert.Equal(t, "UserServiceTests", info.Metadata["className"])
	assert.Equal(t, "3", info.Metadata["testCount"])
}

// TestCSharpTestAnalyzer_FindTestsForFile tests finding tests for a source file
func TestCSharpTestAnalyzer_FindTestsForFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "csharp-find-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source file
	srcDir := filepath.Join(tempDir, "src")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(srcDir, "UserService.cs")
	err = os.WriteFile(srcFile, []byte("public class UserService {}"), 0o644)
	require.NoError(t, err)

	// Create test directory and test file
	testDir := filepath.Join(tempDir, "tests")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	testFile := filepath.Join(testDir, "UserServiceTests.cs")
	testContent := `using Xunit;
public class UserServiceTests {
    [Fact]
    public void Test() { var svc = new UserService(); }
}`
	err = os.WriteFile(testFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForFile(context.Background(), "src/UserService.cs", tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 1)
}

// TestCSharpTestAnalyzer_CountTestMethods tests counting test methods
func TestCSharpTestAnalyzer_CountTestMethods(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name: "xunit_tests",
			content: `[Fact] public void Test1() {}
[Fact] public void Test2() {}
[Theory] public void Test3() {}`,
			expected: 3,
		},
		{
			name: "nunit_tests",
			content: `[Test] public void Test1() {}
[TestCase(1)] public void Test2() {}`,
			expected: 2,
		},
		{
			name: "mstest_tests",
			content: `[TestMethod] public void Test1() {}
[TestMethod] public void Test2() {}`,
			expected: 2,
		},
		{
			name:     "no_tests",
			content:  `public void Helper() {}`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := analyzer.countTestMethods(tt.content)
			assert.Equal(t, tt.expected, count)
		})
	}
}


// TestCSharpTestAnalyzer_RemoveDuplicates tests duplicate removal
func TestCSharpTestAnalyzer_RemoveDuplicates(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

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

// TestCSharpTestAnalyzer_ReferencesClass tests class reference detection
func TestCSharpTestAnalyzer_ReferencesClass(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tests := []struct {
		name      string
		content   string
		className string
		expected  bool
	}{
		{
			name:      "new_instance",
			content:   `var user = new UserService();`,
			className: "UserService",
			expected:  true,
		},
		{
			name:      "static_call",
			content:   `UserService.Create();`,
			className: "UserService",
			expected:  true,
		},
		{
			name:      "generic_type",
			content:   `List<UserService> services;`,
			className: "UserService",
			expected:  true,
		},
		{
			name:      "no_reference",
			content:   `var order = new OrderService();`,
			className: "UserService",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.referencesClass(tt.content, tt.className)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCSharpTestAnalyzer_DetermineTestType tests test type determination
func TestCSharpTestAnalyzer_DetermineTestType(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tests := []struct {
		name     string
		content  string
		fileName string
		expected string
	}{
		{
			name:     "smoke_test",
			content:  `// Smoke test`,
			fileName: "SmokeTests.cs",
			expected: "smoke",
		},
		{
			name:     "integration_test",
			content:  `WebApplicationFactory<Program>`,
			fileName: "ApiTests.cs",
			expected: "integration",
		},
		{
			name:     "unit_test",
			content:  `[Fact] public void Test() {}`,
			fileName: "UserTests.cs",
			expected: "unit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.determineTestType(tt.content, tt.fileName)
			assert.Equal(t, tt.expected, result)
		})
	}
}


// TestCSharpTestAnalyzer_AnalyzeTestDependencies_NonExistent tests error handling
func TestCSharpTestAnalyzer_AnalyzeTestDependencies_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	_, err := analyzer.AnalyzeTestDependencies(context.Background(), "/nonexistent/file.cs")
	assert.Error(t, err)
}

// TestCSharpTestAnalyzer_IsSmokeTest_NonExistent tests error handling
func TestCSharpTestAnalyzer_IsSmokeTest_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	_, err := analyzer.IsSmokeTest(context.Background(), "/nonexistent/file.cs")
	assert.Error(t, err)
}

// TestCSharpTestAnalyzer_AnalyzeTestFile_NonExistent tests error handling
func TestCSharpTestAnalyzer_AnalyzeTestFile_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	_, err := analyzer.AnalyzeTestFile(context.Background(), "/nonexistent/file.cs")
	assert.Error(t, err)
}

// TestCSharpTestAnalyzer_ExtractUsings tests using statement extraction
func TestCSharpTestAnalyzer_ExtractUsings(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	content := `using System;
using System.Collections.Generic;
using Xunit;
using MyApp.Services;

namespace MyApp.Tests {}`

	usings := analyzer.extractUsings(content)

	assert.Len(t, usings, 4)
	assert.Contains(t, usings, "System")
	assert.Contains(t, usings, "System.Collections.Generic")
	assert.Contains(t, usings, "Xunit")
	assert.Contains(t, usings, "MyApp.Services")
}

// TestCSharpTestAnalyzer_ExtractClassName tests class name extraction
func TestCSharpTestAnalyzer_ExtractClassName(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple_class",
			content:  `public class UserServiceTests {}`,
			expected: "UserServiceTests",
		},
		{
			name:     "class_with_inheritance",
			content:  `public class UserServiceTests : TestBase {}`,
			expected: "UserServiceTests",
		},
		{
			name:     "no_class",
			content:  `namespace MyApp {}`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.extractClassName(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCSharpTestAnalyzer_DetectTestFramework tests framework detection
func TestCSharpTestAnalyzer_DetectTestFramework(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "xunit",
			content:  `[Fact] public void Test() {}`,
			expected: "xunit",
		},
		{
			name:     "nunit",
			content:  `[Test] public void Test() {}`,
			expected: "nunit",
		},
		{
			name:     "mstest",
			content:  `[TestMethod] public void Test() {}`,
			expected: "mstest",
		},
		{
			name:     "unknown",
			content:  `public void Test() {}`,
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.detectTestFramework(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}


// TestCSharpTestAnalyzer_ExtractUsings_EmptyContent tests extracting usings from empty content
func TestCSharpTestAnalyzer_ExtractUsings_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	usings := analyzer.extractUsings("")
	assert.Empty(t, usings)
}

// TestCSharpTestAnalyzer_ExtractUsings_DuplicateUsings tests duplicate using handling
func TestCSharpTestAnalyzer_ExtractUsings_DuplicateUsings(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	content := `using System;
using System;
using System.Collections.Generic;`

	usings := analyzer.extractUsings(content)
	assert.Len(t, usings, 2)
}

// TestCSharpTestAnalyzer_DetectTestFramework_EmptyContent tests framework detection with empty content
func TestCSharpTestAnalyzer_DetectTestFramework_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	result := analyzer.detectTestFramework("")
	assert.Equal(t, "unknown", result)
}

// TestCSharpTestAnalyzer_ExtractClassName_EmptyContent tests class name extraction with empty content
func TestCSharpTestAnalyzer_ExtractClassName_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	result := analyzer.extractClassName("")
	assert.Empty(t, result)
}

// TestCSharpTestAnalyzer_CountTestMethods_EmptyContent tests counting test methods with empty content
func TestCSharpTestAnalyzer_CountTestMethods_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	count := analyzer.countTestMethods("")
	assert.Equal(t, 0, count)
}

// TestCSharpTestAnalyzer_DetermineTestType_EmptyContent tests test type determination with empty content
func TestCSharpTestAnalyzer_DetermineTestType_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	result := analyzer.determineTestType("", "Test.cs")
	assert.Equal(t, "unit", result)
}

// TestCSharpTestAnalyzer_IsTestFile_NonExistent tests test file detection for non-existent file
func TestCSharpTestAnalyzer_IsTestFile_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	result := analyzer.isTestFile("/nonexistent/file.cs")
	assert.False(t, result)
}

// TestCSharpTestAnalyzer_ReferencesClass_EmptyContent tests class reference detection with empty content
func TestCSharpTestAnalyzer_ReferencesClass_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	result := analyzer.referencesClass("", "UserService")
	assert.False(t, result)
}

// TestCSharpTestAnalyzer_FindTestsForFile_NoTestDirs tests when no test directories exist
func TestCSharpTestAnalyzer_FindTestsForFile_NoTestDirs(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "csharp-no-test-dirs-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source file without test directories
	srcDir := filepath.Join(tempDir, "src")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(srcDir, "UserService.cs")
	err = os.WriteFile(srcFile, []byte("public class UserService {}"), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForFile(context.Background(), "src/UserService.cs", tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestCSharpTestAnalyzer_AnalyzeTestFile_FileNotFound tests analyzing non-existent file
func TestCSharpTestAnalyzer_AnalyzeTestFile_FileNotFound(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	_, err := analyzer.AnalyzeTestFile(context.Background(), "/nonexistent/file.cs")
	assert.Error(t, err)
}

// TestCSharpTestAnalyzer_FindTestsForSymbol tests finding tests for a symbol
func TestCSharpTestAnalyzer_FindTestsForSymbol(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "csharp-symbol-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test file that references the symbol
	testDir := filepath.Join(tempDir, "tests")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	testFile := filepath.Join(testDir, "UserServiceTests.cs")
	content := `using Xunit;
public class UserServiceTests {
    [Fact]
    public void Test() { var svc = new UserService(); }
}`
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForSymbol(context.Background(), "UserService", tempDir)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(tests), 1)
}

// TestCSharpTestAnalyzer_FindTestsForSymbol_NoMatches tests finding tests with no matches
func TestCSharpTestAnalyzer_FindTestsForSymbol_NoMatches(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewCSharpTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "csharp-no-symbol-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests, err := analyzer.FindTestsForSymbol(context.Background(), "NonExistentSymbol", tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestCSharpTestRunner_ParseTestOutput_WithFailedTestPaths tests parsing output with failed test paths
func TestCSharpTestRunner_ParseTestOutput_WithFailedTestPaths(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	output := `Test run for MyProject.Tests.dll
Failed MyProject.Tests.UserServiceTests.TestCreate
Failed MyProject.Tests.UserServiceTests.TestUpdate
Passed: 8
Failed: 2
Skipped: 0
Total: 10`

	result := runner.ParseTestOutput(output)
	assert.False(t, result.Success)
	assert.Equal(t, 2, result.FailedTests)
}

// TestCSharpTestRunner_isCSharpTestFile_EmptyFile tests test file detection with empty file
func TestCSharpTestRunner_isCSharpTestFile_EmptyFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewCSharpTestRunner(log)

	tempDir, err := os.MkdirTemp("", "csharp-empty-file-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "Empty.cs")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	result := runner.isCSharpTestFile("Empty.cs", testFile)
	assert.False(t, result)
}
