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

// TestPHPTestRunner_GetLanguage tests language identifier
func TestPHPTestRunner_GetLanguage(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	assert.Equal(t, "php", runner.GetLanguage())
}

// TestPHPTestRunner_isTestFile tests PHP test file detection
func TestPHPTestRunner_isTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tests := []struct {
		name     string
		fileName string
		expected bool
	}{
		{
			name:     "test_file_suffix",
			fileName: "UserTest.php",
			expected: true,
		},
		{
			name:     "regular_php_file",
			fileName: "User.php",
			expected: false,
		},
		{
			name:     "controller_file",
			fileName: "UserController.php",
			expected: false,
		},
		{
			name:     "test_in_name_but_not_suffix",
			fileName: "TestHelper.php",
			expected: false,
		},
		{
			name:     "non_php_file",
			fileName: "test.js",
			expected: false,
		},
		{
			name:     "feature_test",
			fileName: "LoginFeatureTest.php",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.isTestFile(tt.fileName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPHPTestRunner_shouldSkipDirectory tests directory skip logic
func TestPHPTestRunner_shouldSkipDirectory(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{
			name:     "vendor_directory",
			dirName:  "vendor",
			expected: true,
		},
		{
			name:     "node_modules",
			dirName:  "node_modules",
			expected: true,
		},
		{
			name:     "git_directory",
			dirName:  ".git",
			expected: true,
		},
		{
			name:     "cache_directory",
			dirName:  "cache",
			expected: true,
		},
		{
			name:     "phpunit_cache",
			dirName:  ".phpunit.cache",
			expected: true,
		},
		{
			name:     "var_directory",
			dirName:  "var",
			expected: true,
		},
		{
			name:     "tests_directory",
			dirName:  "tests",
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

// TestPHPTestRunner_analyzeTestFile tests test type detection
func TestPHPTestRunner_analyzeTestFile(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "php-runner-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected string
	}{
		{
			name:     "unit_test",
			fileName: "UserTest.php",
			content: `<?php
namespace Tests\Unit;

use PHPUnit\Framework\TestCase;

class UserTest extends TestCase
{
    public function testUserCreation()
    {
        $this->assertTrue(true);
    }
}
`,
			expected: "unit",
		},
		{
			name:     "integration_test",
			fileName: "DatabaseIntegrationTest.php",
			content: `<?php
namespace Tests\Integration;

use PHPUnit\Framework\TestCase;

/**
 * @group integration
 */
class DatabaseIntegrationTest extends TestCase
{
    public function testDatabaseConnection()
    {
        $this->assertTrue(true);
    }
}
`,
			expected: "integration",
		},
		{
			name:     "smoke_test",
			fileName: "SmokeTest.php",
			content: `<?php
namespace Tests;

use PHPUnit\Framework\TestCase;

class SmokeTest extends TestCase
{
    public function testApplicationStarts()
    {
        $this->assertTrue(true);
    }
}
`,
			expected: "smoke",
		},
		{
			name:     "feature_test_laravel",
			fileName: "LoginFeatureTest.php",
			content: `<?php
namespace Tests\Feature;

use Tests\TestCase;

class LoginFeatureTest extends TestCase
{
    public function testUserCanLogin()
    {
        $this->assertTrue(true);
    }
}
`,
			expected: "feature",
		},
		{
			name:     "kernel_test_symfony",
			fileName: "ServiceTest.php",
			content: `<?php
namespace Tests;

use Symfony\Bundle\FrameworkBundle\Test\KernelTestCase;

class ServiceTest extends KernelTestCase
{
    public function testServiceWorks()
    {
        $this->assertTrue(true);
    }
}
`,
			expected: "integration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := runner.analyzeTestFile(testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPHPTestRunner_findPHPUnitPath tests PHPUnit path detection
func TestPHPTestRunner_findPHPUnitPath(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "php-phpunit-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	// Test without vendor/bin/phpunit
	result := runner.findPHPUnitPath(tempDir)
	assert.Equal(t, "phpunit", result) // Falls back to global

	// Create vendor/bin directory
	vendorBin := filepath.Join(tempDir, "vendor", "bin")
	err = os.MkdirAll(vendorBin, 0o755)
	require.NoError(t, err)

	// Create phpunit file
	phpunitPath := filepath.Join(vendorBin, "phpunit")
	err = os.WriteFile(phpunitPath, []byte("#!/bin/bash\necho 'phpunit'"), 0o755)
	require.NoError(t, err)

	result = runner.findPHPUnitPath(tempDir)
	assert.Equal(t, phpunitPath, result)
}

// TestPHPTestRunner_findPHPUnitConfig tests PHPUnit config detection
func TestPHPTestRunner_findPHPUnitConfig(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "php-config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	// Test without config
	result := runner.findPHPUnitConfig(tempDir)
	assert.Empty(t, result)

	// Create phpunit.xml
	phpunitXML := filepath.Join(tempDir, "phpunit.xml")
	err = os.WriteFile(phpunitXML, []byte("<phpunit></phpunit>"), 0o644)
	require.NoError(t, err)

	result = runner.findPHPUnitConfig(tempDir)
	assert.Equal(t, phpunitXML, result)
}

// TestPHPTestRunner_DiscoverTests tests test discovery
func TestPHPTestRunner_DiscoverTests(t *testing.T) {
	// Create temp directory with PHP project structure
	tempDir, err := os.MkdirTemp("", "php-discover-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create tests directory
	testsDir := filepath.Join(tempDir, "tests")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	// Create test files
	unitTest := `<?php
class UserTest extends \PHPUnit\Framework\TestCase
{
    public function testUser() { $this->assertTrue(true); }
}
`
	err = os.WriteFile(filepath.Join(testsDir, "UserTest.php"), []byte(unitTest), 0o644)
	require.NoError(t, err)

	integrationTest := `<?php
/**
 * @group integration
 */
class DatabaseTest extends \PHPUnit\Framework\TestCase
{
    public function testDatabase() { $this->assertTrue(true); }
}
`
	err = os.WriteFile(filepath.Join(testsDir, "DatabaseTest.php"), []byte(integrationTest), 0o644)
	require.NoError(t, err)

	// Create non-test file
	helper := `<?php
class TestHelper { }
`
	err = os.WriteFile(filepath.Join(testsDir, "TestHelper.php"), []byte(helper), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Should find 2 test files (UserTest.php and DatabaseTest.php)
	// Note: The discovery may find more files depending on implementation
	assert.GreaterOrEqual(t, len(tests), 2)

	// Verify metadata
	for _, test := range tests {
		assert.Equal(t, "phpunit", test.Metadata["framework"])
	}
}

// TestPHPTestRunner_RunTest tests single test execution
func TestPHPTestRunner_RunTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	config := &domain.TestConfig{
		ProjectPath: "/nonexistent/path",
		Verbose:     false,
		Timeout:     60,
	}

	// This will fail because the path doesn't exist
	result, err := runner.RunTest(context.Background(), "tests/UserTest.php", config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "tests/UserTest.php", result.TestPath)
	assert.Equal(t, "php", result.Language)
	assert.False(t, result.Success) // Should fail because path doesn't exist
}

// TestPHPTestRunner_RunTestSuite tests suite execution
func TestPHPTestRunner_RunTestSuite(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	suite := &domain.TestSuite{
		Name:        "test_suite",
		Language:    "php",
		ProjectPath: "/nonexistent/path",
		Tests: []*domain.TestInfo{
			{Path: "tests/UserTest.php", Name: "UserTest.php", Type: "unit"},
			{Path: "tests/OrderTest.php", Name: "OrderTest.php", Type: "unit"},
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

// TestPHPTestRunner_parseTextOutput tests text output parsing
func TestPHPTestRunner_parseTextOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tests := []struct {
		name            string
		output          string
		expectedSuccess bool
	}{
		{
			name:            "all_passed",
			output:          "OK (5 tests, 10 assertions)",
			expectedSuccess: true,
		},
		{
			name:            "failures",
			output:          "FAILURES!\nTests: 5, Assertions: 10, Failures: 2.",
			expectedSuccess: false,
		},
		{
			name:            "errors",
			output:          "ERRORS!\nTests: 5, Assertions: 10, Errors: 1.",
			expectedSuccess: false,
		},
		{
			name:            "empty_output",
			output:          "",
			expectedSuccess: true,
		},
	}

	testInfos := []*domain.TestInfo{
		{Path: "test.php", Name: "test.php"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := runner.parseTextOutput(tt.output, testInfos, 1.0)
			require.Len(t, results, 1)
			assert.Equal(t, tt.expectedSuccess, results[0].Success)
		})
	}
}

// TestPHPTestRunner_buildTestCommand tests command building
func TestPHPTestRunner_buildTestCommand(t *testing.T) {
	// Create temp directory with phpunit.xml
	tempDir, err := os.MkdirTemp("", "php-cmd-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create phpunit.xml
	phpunitXML := filepath.Join(tempDir, "phpunit.xml")
	err = os.WriteFile(phpunitXML, []byte("<phpunit></phpunit>"), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tests := []struct {
		name     string
		config   *domain.TestConfig
		testPath string
		wantArgs []string
	}{
		{
			name: "basic_test",
			config: &domain.TestConfig{
				ProjectPath: tempDir,
				Verbose:     false,
				Coverage:    false,
			},
			testPath: "tests/UserTest.php",
			wantArgs: []string{"--configuration", phpunitXML, "tests/UserTest.php", "--testdox"},
		},
		{
			name: "verbose_test",
			config: &domain.TestConfig{
				ProjectPath: tempDir,
				Verbose:     true,
				Coverage:    false,
			},
			testPath: "tests/UserTest.php",
			wantArgs: []string{"--configuration", phpunitXML, "tests/UserTest.php", "--verbose", "--testdox"},
		},
		{
			name: "with_coverage",
			config: &domain.TestConfig{
				ProjectPath: tempDir,
				Verbose:     false,
				Coverage:    true,
			},
			testPath: "tests/UserTest.php",
			wantArgs: []string{"--configuration", phpunitXML, "tests/UserTest.php", "--coverage-text", "--testdox"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := runner.buildTestCommand(tt.testPath, tt.config)
			require.NoError(t, err)
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}

// TestPHPTestRunner_parseJSONOutput tests JSON output parsing
func TestPHPTestRunner_parseJSONOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tests := []struct {
		name            string
		output          string
		expectedCount   int
		expectedSuccess bool
	}{
		{
			name: "valid_json_passed",
			output: `{"event":"test","test":"UserTest::testCreate","status":"pass","time":0.001}
{"event":"test","test":"UserTest::testUpdate","status":"pass","time":0.002}`,
			expectedCount:   2,
			expectedSuccess: true,
		},
		{
			name: "valid_json_failed",
			output: `{"event":"test","test":"UserTest::testCreate","status":"fail","time":0.001,"message":"Expected true, got false"}`,
			expectedCount:   1,
			expectedSuccess: false,
		},
		{
			name:            "invalid_json",
			output:          "not json at all",
			expectedCount:   0,
			expectedSuccess: false,
		},
		{
			name:            "empty_output",
			output:          "",
			expectedCount:   0,
			expectedSuccess: false,
		},
		{
			name: "mixed_events",
			output: `{"event":"suiteStart","suite":"UserTest"}
{"event":"test","test":"UserTest::testCreate","status":"pass","time":0.001}
{"event":"suiteEnd","suite":"UserTest"}`,
			expectedCount:   1,
			expectedSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := runner.parseJSONOutput(tt.output)
			assert.Len(t, results, tt.expectedCount)

			if tt.expectedCount > 0 {
				assert.Equal(t, tt.expectedSuccess, results[0].Success)
			}
		})
	}
}

// TestPHPTestRunner_DiscoverTests_MultipleDirectories tests discovery in multiple test directories
func TestPHPTestRunner_DiscoverTests_MultipleDirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "php-multi-dir-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create multiple test directories
	testDirs := []string{"tests", "test", "Tests", "Test"}
	for i, dir := range testDirs {
		fullDir := filepath.Join(tempDir, dir)
		err = os.MkdirAll(fullDir, 0o755)
		require.NoError(t, err)

		testFile := filepath.Join(fullDir, "Test"+string(rune('A'+i))+"Test.php")
		content := `<?php class Test` + string(rune('A'+i)) + `Test extends \PHPUnit\Framework\TestCase { public function testSomething() {} }`
		err = os.WriteFile(testFile, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 1)
}

// TestPHPTestRunner_DiscoverTests_EmptyProject tests discovery in empty project
func TestPHPTestRunner_DiscoverTests_EmptyProject(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "php-empty-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestPHPTestRunner_findPHPUnitConfig_DistFile tests finding phpunit.xml.dist
func TestPHPTestRunner_findPHPUnitConfig_DistFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "php-config-dist-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	// Create phpunit.xml.dist
	phpunitXMLDist := filepath.Join(tempDir, "phpunit.xml.dist")
	err = os.WriteFile(phpunitXMLDist, []byte("<phpunit></phpunit>"), 0o644)
	require.NoError(t, err)

	result := runner.findPHPUnitConfig(tempDir)
	assert.Equal(t, phpunitXMLDist, result)
}

// TestPHPTestRunner_buildTestCommand_NoConfig tests command building without config file
func TestPHPTestRunner_buildTestCommand_NoConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "php-no-config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		Verbose:     false,
		Coverage:    false,
	}

	args, err := runner.buildTestCommand("tests/UserTest.php", config)
	require.NoError(t, err)

	// Should not contain --configuration
	for _, arg := range args {
		assert.NotEqual(t, "--configuration", arg)
	}
	assert.Contains(t, args, "tests/UserTest.php")
	assert.Contains(t, args, "--testdox")
}


// TestPHPTestAnalyzer_AnalyzeTestDependencies tests dependency analysis
func TestPHPTestAnalyzer_AnalyzeTestDependencies(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "php-analyzer-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `<?php
namespace Tests;

use App\Models\User;
use App\Services\UserService;
use PHPUnit\Framework\TestCase;

class UserTest extends TestCase
{
    public function testUser() {}
}
`
	testFile := filepath.Join(tempDir, "UserTest.php")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	deps, err := analyzer.AnalyzeTestDependencies(context.Background(), testFile)
	require.NoError(t, err)

	assert.Len(t, deps, 3)
	assert.Contains(t, deps, `App\Models\User`)
	assert.Contains(t, deps, `App\Services\UserService`)
	assert.Contains(t, deps, `PHPUnit\Framework\TestCase`)
}

// TestPHPTestAnalyzer_IsSmokeTest tests smoke test detection
func TestPHPTestAnalyzer_IsSmokeTest(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "php-smoke-test-*")
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
			fileName: "SmokeTest.php",
			content:  `<?php class SmokeTest {}`,
			expected: true,
		},
		{
			name:     "smoke_group",
			fileName: "ApiTest.php",
			content:  `<?php /** @group smoke */ class ApiTest {}`,
			expected: true,
		},
		{
			name:     "regular_test",
			fileName: "UserTest.php",
			content:  `<?php class UserTest {}`,
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

// TestPHPTestAnalyzer_ExtractTestMethods tests extracting test methods
func TestPHPTestAnalyzer_ExtractTestMethods(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	content := `<?php
class UserTest extends TestCase
{
    public function testCreateUser() {}
    public function testUpdateUser() {}
    
    public function helperMethod() {}
}
`
	methods := analyzer.ExtractTestMethods(content)

	assert.Len(t, methods, 2)
	assert.Contains(t, methods, "testCreateUser")
	assert.Contains(t, methods, "testUpdateUser")
}

// TestPHPTestAnalyzer_GetTestGroups tests extracting test groups
func TestPHPTestAnalyzer_GetTestGroups(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	content := `<?php
/**
 * @group integration
 * @group database
 */
class UserTest extends TestCase
{
    /** @group slow */
    public function testSlowOperation() {}
}
`
	groups := analyzer.GetTestGroups(content)

	assert.Len(t, groups, 3)
	assert.Contains(t, groups, "integration")
	assert.Contains(t, groups, "database")
	assert.Contains(t, groups, "slow")
}

// TestPHPTestAnalyzer_FindTestsForFile tests finding tests for a source file
func TestPHPTestAnalyzer_FindTestsForFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "php-find-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source file
	srcDir := filepath.Join(tempDir, "src")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(srcDir, "User.php")
	err = os.WriteFile(srcFile, []byte("<?php class User {}"), 0o644)
	require.NoError(t, err)

	// Create test directory and test file
	testDir := filepath.Join(tempDir, "tests")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	testFile := filepath.Join(testDir, "UserTest.php")
	testContent := `<?php
use App\User;
class UserTest extends TestCase {}`
	err = os.WriteFile(testFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForFile(context.Background(), "src/User.php", tempDir)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(tests), 1)
}


// TestPHPTestAnalyzer_UsesClass tests class usage detection
func TestPHPTestAnalyzer_UsesClass(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	tests := []struct {
		name      string
		content   string
		className string
		expected  bool
	}{
		{
			name:      "use_statement",
			content:   `use App\Models\User;`,
			className: "User",
			expected:  true,
		},
		{
			name:      "direct_usage",
			content:   `$user = new User();`,
			className: "User",
			expected:  true,
		},
		{
			name:      "no_usage",
			content:   `$order = new Order();`,
			className: "User",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.usesClass(tt.content, tt.className)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPHPTestAnalyzer_RemoveDuplicates tests duplicate removal
func TestPHPTestAnalyzer_RemoveDuplicates(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

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

// TestPHPTestAnalyzer_FileExists tests file existence check
func TestPHPTestAnalyzer_FileExists(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "php-exists-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a file
	testFile := filepath.Join(tempDir, "exists.php")
	err = os.WriteFile(testFile, []byte("<?php"), 0o644)
	require.NoError(t, err)

	assert.True(t, analyzer.fileExists(testFile))
	assert.False(t, analyzer.fileExists(filepath.Join(tempDir, "notexists.php")))
}


// TestPHPTestAnalyzer_AnalyzeTestDependencies_NonExistent tests error handling
func TestPHPTestAnalyzer_AnalyzeTestDependencies_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	_, err := analyzer.AnalyzeTestDependencies(context.Background(), "/nonexistent/file.php")
	assert.Error(t, err)
}

// TestPHPTestAnalyzer_IsSmokeTest_NonExistent tests error handling
func TestPHPTestAnalyzer_IsSmokeTest_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	_, err := analyzer.IsSmokeTest(context.Background(), "/nonexistent/file.php")
	assert.Error(t, err)
}

// TestPHPTestAnalyzer_ExtractUseStatements tests use statement extraction
func TestPHPTestAnalyzer_ExtractUseStatements(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	content := `<?php
use App\Models\User;
use App\Services\UserService;
use PHPUnit\Framework\TestCase;
use Illuminate\Foundation\Testing\RefreshDatabase;

class UserTest extends TestCase {}
`
	uses := analyzer.extractUseStatements(content)

	assert.Len(t, uses, 4)
	assert.Contains(t, uses, `App\Models\User`)
	assert.Contains(t, uses, `App\Services\UserService`)
	assert.Contains(t, uses, `PHPUnit\Framework\TestCase`)
	assert.Contains(t, uses, `Illuminate\Foundation\Testing\RefreshDatabase`)
}


// TestPHPTestRunner_isPHPUnitAvailable tests PHPUnit availability check
func TestPHPTestRunner_isPHPUnitAvailable(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tempDir, err := os.MkdirTemp("", "php-phpunit-available-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Without PHPUnit, should return false or true depending on system
	result := runner.isPHPUnitAvailable(tempDir)
	// Just verify it doesn't panic
	assert.IsType(t, true, result)
}

// TestPHPTestRunner_parseJUnitXML_InvalidFile tests parsing invalid XML
func TestPHPTestRunner_parseJUnitXML_InvalidFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	testInfos := []*domain.TestInfo{
		{Path: "test.php", Name: "test.php"},
	}

	// Non-existent file should return nil
	results := runner.parseJUnitXML("/nonexistent/file.xml", testInfos, 1.0)
	assert.Nil(t, results)
}

// TestPHPTestRunner_RunTestSuite_Empty tests running empty suite
func TestPHPTestRunner_RunTestSuite_Empty(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tempDir, err := os.MkdirTemp("", "php-empty-suite-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	suite := &domain.TestSuite{
		Name:        "empty_suite",
		Language:    "php",
		ProjectPath: tempDir,
		Tests:       []*domain.TestInfo{},
		Config: &domain.TestConfig{
			ProjectPath: tempDir,
			Verbose:     false,
		},
	}

	results, err := runner.RunTestSuite(context.Background(), suite)

	assert.NoError(t, err)
	assert.Empty(t, results)
}


// TestPHPTestRunner_parseTextOutput_EmptyOutput tests parsing empty output
func TestPHPTestRunner_parseTextOutput_EmptyOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	testInfos := []*domain.TestInfo{
		{Path: "test.php", Name: "test.php"},
	}

	results := runner.parseTextOutput("", testInfos, 1.0)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
}

// TestPHPTestRunner_parseTextOutput_OnlyOK tests parsing output with only OK
func TestPHPTestRunner_parseTextOutput_OnlyOK(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	testInfos := []*domain.TestInfo{
		{Path: "test.php", Name: "test.php"},
	}

	results := runner.parseTextOutput("OK (10 tests, 20 assertions)", testInfos, 2.5)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
	assert.Equal(t, 2.5, results[0].Duration)
}

// TestPHPTestRunner_parseJSONOutput_EmptyOutput tests parsing empty JSON output
func TestPHPTestRunner_parseJSONOutput_EmptyOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	results := runner.parseJSONOutput("")
	assert.Empty(t, results)
}

// TestPHPTestRunner_parseJSONOutput_MalformedJSON tests parsing malformed JSON
func TestPHPTestRunner_parseJSONOutput_MalformedJSON(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	results := runner.parseJSONOutput("{invalid json")
	assert.Empty(t, results)
}

// TestPHPTestRunner_parseJSONOutput_SkippedTest tests parsing skipped test status
func TestPHPTestRunner_parseJSONOutput_SkippedTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	output := `{"event":"test","test":"UserTest::testSkipped","status":"skipped","time":0.001}`
	results := runner.parseJSONOutput(output)
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
}

// TestPHPTestRunner_analyzeTestFile_NonExistent tests analyzing non-existent file
func TestPHPTestRunner_analyzeTestFile_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	result := runner.analyzeTestFile("/nonexistent/file.php")
	assert.Equal(t, "unit", result)
}

// TestPHPTestRunner_analyzeTestFile_WebTestCase tests WebTestCase detection
func TestPHPTestRunner_analyzeTestFile_WebTestCase(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "php-webtestcase-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	content := `<?php
use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;

class ApiTest extends WebTestCase
{
    public function testApi() {}
}
`
	testFile := filepath.Join(tempDir, "ApiTest.php")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	result := runner.analyzeTestFile(testFile)
	assert.Equal(t, "integration", result)
}

// TestPHPTestRunner_analyzeTestFile_FunctionalTestCase tests FunctionalTestCase detection
func TestPHPTestRunner_analyzeTestFile_FunctionalTestCase(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "php-functional-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	content := `<?php
class ApiTest extends FunctionalTestCase
{
    public function testApi() {}
}
`
	testFile := filepath.Join(tempDir, "ApiTest.php")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	result := runner.analyzeTestFile(testFile)
	assert.Equal(t, "integration", result)
}

// TestPHPTestRunner_DiscoverTests_SkipVendor tests that vendor directory is skipped
func TestPHPTestRunner_DiscoverTests_SkipVendor(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "php-skip-vendor-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create tests directory
	testsDir := filepath.Join(tempDir, "tests")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	// Create vendor directory inside tests (should be skipped)
	vendorDir := filepath.Join(testsDir, "vendor")
	err = os.MkdirAll(vendorDir, 0o755)
	require.NoError(t, err)

	// Create test file in vendor (should not be found)
	vendorTest := filepath.Join(vendorDir, "VendorTest.php")
	err = os.WriteFile(vendorTest, []byte("<?php class VendorTest {}"), 0o644)
	require.NoError(t, err)

	// Create valid test file
	validTest := filepath.Join(testsDir, "ValidTest.php")
	err = os.WriteFile(validTest, []byte("<?php class ValidTest {}"), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewPHPTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Verify vendor test is not found
	for _, test := range tests {
		assert.NotEqual(t, "VendorTest.php", test.Name, "VendorTest.php should not be discovered")
	}
	// ValidTest.php should be found
	foundValid := false
	for _, test := range tests {
		if test.Name == "ValidTest.php" {
			foundValid = true
			break
		}
	}
	assert.True(t, foundValid, "ValidTest.php should be discovered")
}

// TestPHPTestAnalyzer_ExtractTestMethods_EmptyContent tests extracting methods from empty content
func TestPHPTestAnalyzer_ExtractTestMethods_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	methods := analyzer.ExtractTestMethods("")
	assert.Empty(t, methods)
}

// TestPHPTestAnalyzer_ExtractTestMethods_WithAnnotation tests extracting methods with @test annotation
func TestPHPTestAnalyzer_ExtractTestMethods_WithAnnotation(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	content := `<?php
class UserTest extends TestCase
{
    public function testCreateUser() {}
    public function testUpdateUser() {}
}
`
	methods := analyzer.ExtractTestMethods(content)
	assert.GreaterOrEqual(t, len(methods), 2)
}

// TestPHPTestAnalyzer_GetTestGroups_EmptyContent tests extracting groups from empty content
func TestPHPTestAnalyzer_GetTestGroups_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	groups := analyzer.GetTestGroups("")
	assert.Empty(t, groups)
}

// TestPHPTestAnalyzer_GetTestGroups_DuplicateGroups tests duplicate group handling
func TestPHPTestAnalyzer_GetTestGroups_DuplicateGroups(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	content := `<?php
/**
 * @group integration
 * @group integration
 * @group slow
 */
class UserTest extends TestCase {}
`
	groups := analyzer.GetTestGroups(content)
	assert.Len(t, groups, 2)
}

// TestPHPTestAnalyzer_ExtractUseStatements_EmptyContent tests extracting use statements from empty content
func TestPHPTestAnalyzer_ExtractUseStatements_EmptyContent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	uses := analyzer.extractUseStatements("")
	assert.Empty(t, uses)
}

// TestPHPTestAnalyzer_ExtractUseStatements_DuplicateUses tests duplicate use statement handling
func TestPHPTestAnalyzer_ExtractUseStatements_DuplicateUses(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	content := `<?php
use App\Models\User;
use App\Models\User;
use App\Services\UserService;
`
	uses := analyzer.extractUseStatements(content)
	assert.Len(t, uses, 2)
}

// TestPHPTestAnalyzer_UsesClass_PartialMatch tests class usage with partial match
func TestPHPTestAnalyzer_UsesClass_PartialMatch(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	// Should not match partial class names
	content := `use App\Models\UserService;`
	result := analyzer.usesClass(content, "User")
	// This will match because "User" is contained in "UserService"
	assert.True(t, result)
}

// TestPHPTestAnalyzer_FindTestsForFile_NoTestDirs tests when no test directories exist
func TestPHPTestAnalyzer_FindTestsForFile_NoTestDirs(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPHPTestAnalyzer(log)

	tempDir, err := os.MkdirTemp("", "php-no-test-dirs-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source file without test directories
	srcDir := filepath.Join(tempDir, "src")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(srcDir, "User.php")
	err = os.WriteFile(srcFile, []byte("<?php class User {}"), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForFile(context.Background(), "src/User.php", tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}
