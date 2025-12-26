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

// TestTypeScriptTestRunner_GetLanguage tests language identifier
func TestTypeScriptTestRunner_GetLanguage(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewTypeScriptTestRunner(log)

	assert.Equal(t, "typescript", runner.GetLanguage())
}

// TestTypeScriptTestRunner_shouldSkipDir tests directory skip logic
func TestTypeScriptTestRunner_shouldSkipDir(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewTypeScriptTestRunner(log)

	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{
			name:     "node_modules",
			dirName:  "node_modules",
			expected: true,
		},
		{
			name:     "dist_directory",
			dirName:  "dist",
			expected: true,
		},
		{
			name:     "build_directory",
			dirName:  "build",
			expected: true,
		},
		{
			name:     "coverage_directory",
			dirName:  "coverage",
			expected: true,
		},
		{
			name:     "git_directory",
			dirName:  ".git",
			expected: true,
		},
		{
			name:     "next_directory",
			dirName:  ".next",
			expected: true,
		},
		{
			name:     "nuxt_directory",
			dirName:  ".nuxt",
			expected: true,
		},
		{
			name:     "out_directory",
			dirName:  "out",
			expected: true,
		},
		{
			name:     "hidden_directory",
			dirName:  ".hidden",
			expected: true,
		},
		{
			name:     "src_directory",
			dirName:  "src",
			expected: false,
		},
		{
			name:     "tests_directory",
			dirName:  "tests",
			expected: false,
		},
		{
			name:     "__tests__directory",
			dirName:  "__tests__",
			expected: false,
		},
		{
			name:     "current_directory",
			dirName:  ".",
			expected: false,
		},
		{
			name:     "parent_directory",
			dirName:  "..",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.shouldSkipDir(tt.dirName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestTypeScriptTestRunner_detectTestFramework tests framework detection
func TestTypeScriptTestRunner_detectTestFramework(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewTypeScriptTestRunner(log)

	tests := []struct {
		name         string
		packageJSON  string
		configFiles  map[string]string
		expected     string
	}{
		{
			name: "vitest_in_devDependencies",
			packageJSON: `{
				"name": "test-project",
				"devDependencies": {
					"vitest": "^0.34.0"
				}
			}`,
			expected: frameworkVitest,
		},
		{
			name: "jest_in_devDependencies",
			packageJSON: `{
				"name": "test-project",
				"devDependencies": {
					"jest": "^29.0.0"
				}
			}`,
			expected: frameworkJest,
		},
		{
			name: "mocha_in_devDependencies",
			packageJSON: `{
				"name": "test-project",
				"devDependencies": {
					"mocha": "^10.0.0"
				}
			}`,
			expected: frameworkMocha,
		},
		{
			name: "vitest_in_dependencies",
			packageJSON: `{
				"name": "test-project",
				"dependencies": {
					"vitest": "^0.34.0"
				}
			}`,
			expected: frameworkVitest,
		},
		{
			name: "vitest_config_file",
			packageJSON: `{
				"name": "test-project"
			}`,
			configFiles: map[string]string{
				"vitest.config.ts": "export default {}",
			},
			expected: frameworkVitest,
		},
		{
			name: "jest_config_file",
			packageJSON: `{
				"name": "test-project"
			}`,
			configFiles: map[string]string{
				"jest.config.js": "module.exports = {}",
			},
			expected: frameworkJest,
		},
		{
			name: "mocha_config_file",
			packageJSON: `{
				"name": "test-project"
			}`,
			configFiles: map[string]string{
				".mocharc.json": "{}",
			},
			expected: frameworkMocha,
		},
		{
			name: "no_framework_default_jest",
			packageJSON: `{
				"name": "test-project"
			}`,
			expected: frameworkJest,
		},
		{
			name:        "no_package_json",
			packageJSON: "",
			expected:    frameworkJest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir, err := os.MkdirTemp("", "ts-framework-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Create package.json if provided
			if tt.packageJSON != "" {
				err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(tt.packageJSON), 0o644)
				require.NoError(t, err)
			}

			// Create config files if provided
			for name, content := range tt.configFiles {
				err = os.WriteFile(filepath.Join(tempDir, name), []byte(content), 0o644)
				require.NoError(t, err)
			}

			result := runner.detectTestFramework(tempDir)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestTypeScriptTestRunner_buildTestCommand tests command building
func TestTypeScriptTestRunner_buildTestCommand(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewTypeScriptTestRunner(log)

	tests := []struct {
		name           string
		framework      string
		testPath       string
		config         *domain.TestConfig
		expectedArgs   []string
		expectError    bool
	}{
		{
			name:      "vitest_basic",
			framework: frameworkVitest,
			testPath:  "src/test.spec.ts",
			config: &domain.TestConfig{
				Coverage: false,
			},
			expectedArgs: []string{"vitest", "run", "src/test.spec.ts", "--reporter=json"},
			expectError:  false,
		},
		{
			name:      "vitest_with_coverage",
			framework: frameworkVitest,
			testPath:  "src/test.spec.ts",
			config: &domain.TestConfig{
				Coverage: true,
			},
			expectedArgs: []string{"vitest", "run", "src/test.spec.ts", "--coverage", "--reporter=json"},
			expectError:  false,
		},
		{
			name:      "jest_basic",
			framework: frameworkJest,
			testPath:  "src/test.spec.ts",
			config: &domain.TestConfig{
				Coverage: false,
				Verbose:  false,
			},
			expectedArgs: []string{"jest", "src/test.spec.ts", "--json"},
			expectError:  false,
		},
		{
			name:      "jest_with_coverage_and_verbose",
			framework: frameworkJest,
			testPath:  "src/test.spec.ts",
			config: &domain.TestConfig{
				Coverage: true,
				Verbose:  true,
			},
			expectedArgs: []string{"jest", "src/test.spec.ts", "--json", "--coverage", "--verbose"},
			expectError:  false,
		},
		{
			name:      "mocha_basic",
			framework: frameworkMocha,
			testPath:  "src/test.spec.ts",
			config: &domain.TestConfig{
				Timeout: 0,
			},
			expectedArgs: []string{"mocha", "src/test.spec.ts", "--reporter", "json"},
			expectError:  false,
		},
		{
			name:      "mocha_with_timeout",
			framework: frameworkMocha,
			testPath:  "src/test.spec.ts",
			config: &domain.TestConfig{
				Timeout: 5,
			},
			expectedArgs: []string{"mocha", "src/test.spec.ts", "--reporter", "json", "--timeout", "5000"},
			expectError:  false,
		},
		{
			name:        "unsupported_framework",
			framework:   "unknown",
			testPath:    "src/test.spec.ts",
			config:      &domain.TestConfig{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := runner.buildTestCommand(tt.framework, tt.testPath, tt.config)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedArgs, args)
			}
		})
	}
}

// TestTypeScriptTestRunner_analyzeTestFile tests test file type analysis
func TestTypeScriptTestRunner_analyzeTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewTypeScriptTestRunner(log)

	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "ts-analyze-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "smoke_test",
			content: `describe('smoke tests', () => {
				it('should pass smoke test', () => {
					expect(true).toBe(true);
				});
			});`,
			expected: "smoke",
		},
		{
			name: "integration_test",
			content: `describe('integration tests', () => {
				it('should test integration', () => {
					expect(true).toBe(true);
				});
			});`,
			expected: "integration",
		},
		{
			name: "e2e_test",
			content: `describe('e2e tests', () => {
				it('should test end-to-end', () => {
					expect(true).toBe(true);
				});
			});`,
			expected: "integration",
		},
		{
			name: "api_test",
			content: `import supertest from 'supertest';
			describe('api test', () => {
				it('should test api', () => {});
			});`,
			expected: "integration",
		},
		{
			name: "database_test",
			content: `describe('database tests', () => {
				it('should connect to database', () => {});
			});`,
			expected: "integration",
		},
		{
			name: "http_test",
			content: `describe('http tests', () => {
				it('should make http request', () => {});
			});`,
			expected: "integration",
		},
		{
			name: "fetch_test",
			content: `describe('fetch tests', () => {
				it('should fetch data', async () => {
					await fetch('http://example.com');
				});
			});`,
			expected: "integration",
		},
		{
			name: "axios_test",
			content: `import axios from 'axios';
			describe('axios tests', () => {
				it('should use axios', () => {});
			});`,
			expected: "integration",
		},
		{
			name: "unit_test",
			content: `describe('unit tests', () => {
				it('should add numbers', () => {
					expect(1 + 1).toBe(2);
				});
			});`,
			expected: "unit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, tt.name+".spec.ts")
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := runner.analyzeTestFile(testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestTypeScriptTestRunner_DiscoverTests tests test discovery
func TestTypeScriptTestRunner_DiscoverTests(t *testing.T) {
	// Create temp directory with test structure
	tempDir, err := os.MkdirTemp("", "ts-discover-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create package.json
	packageJSON := `{"name": "test-project", "devDependencies": {"jest": "^29.0.0"}}`
	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0o644)
	require.NoError(t, err)

	// Create test directory structure
	srcDir := filepath.Join(tempDir, "src")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	testsDir := filepath.Join(tempDir, "__tests__")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	// Create test files
	testFiles := map[string]string{
		"src/utils.test.ts":       "describe('utils', () => {});",
		"src/helper.spec.ts":      "describe('helper', () => {});",
		"src/main.ts":             "export const main = () => {};",
		"__tests__/app.test.tsx":  "describe('app', () => {});",
		"src/component.test.jsx":  "describe('component', () => {});",
		"src/service.spec.js":     "describe('service', () => {});",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		err = os.MkdirAll(filepath.Dir(fullPath), 0o755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	// Create node_modules (should be skipped)
	nodeModulesDir := filepath.Join(tempDir, "node_modules", "some-package")
	err = os.MkdirAll(nodeModulesDir, 0o755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(nodeModulesDir, "index.test.js"), []byte("test"), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewTypeScriptTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Should find test files but not files in node_modules or non-test files
	assert.GreaterOrEqual(t, len(tests), 5)

	// Verify test files are found
	foundFiles := make(map[string]bool)
	for _, test := range tests {
		foundFiles[test.Name] = true
	}

	assert.True(t, foundFiles["utils.test.ts"], "utils.test.ts should be found")
	assert.True(t, foundFiles["helper.spec.ts"], "helper.spec.ts should be found")
	assert.True(t, foundFiles["app.test.tsx"], "app.test.tsx should be found")
	assert.True(t, foundFiles["component.test.jsx"], "component.test.jsx should be found")
	assert.True(t, foundFiles["service.spec.js"], "service.spec.js should be found")
	assert.False(t, foundFiles["main.ts"], "main.ts should not be found")
	assert.False(t, foundFiles["index.test.js"], "node_modules test should not be found")
}

// TestTypeScriptTestRunner_RunTest tests single test execution
func TestTypeScriptTestRunner_RunTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewTypeScriptTestRunner(log)

	// Create temp directory with package.json
	tempDir, err := os.MkdirTemp("", "ts-run-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	packageJSON := `{"name": "test-project", "devDependencies": {"jest": "^29.0.0"}}`
	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0o644)
	require.NoError(t, err)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		Verbose:     true,
		Coverage:    false,
		EnvVars: map[string]string{
			"NODE_ENV": "test",
		},
	}

	// This will fail because npx/jest is not available, but we test the structure
	result, err := runner.RunTest(context.Background(), "test.spec.ts", config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test.spec.ts", result.TestPath)
	assert.Equal(t, "typescript", result.Language)
}

// TestTypeScriptTestRunner_RunTestSuite tests running a test suite
func TestTypeScriptTestRunner_RunTestSuite(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewTypeScriptTestRunner(log)

	// Create temp directory with package.json
	tempDir, err := os.MkdirTemp("", "ts-suite-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	packageJSON := `{"name": "test-project", "devDependencies": {"jest": "^29.0.0"}}`
	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0o644)
	require.NoError(t, err)

	suite := &domain.TestSuite{
		Name:        "test_suite",
		Language:    "typescript",
		ProjectPath: tempDir,
		Tests: []*domain.TestInfo{
			{Path: "test1.spec.ts", Name: "test1.spec.ts", Type: "unit"},
			{Path: "test2.spec.ts", Name: "test2.spec.ts", Type: "unit"},
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

// TestTypeScriptTestRunner_DiscoverTests_InTestsDirectory tests discovery in __tests__ directory
func TestTypeScriptTestRunner_DiscoverTests_InTestsDirectory(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "ts-tests-dir-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create __tests__ directory with regular .ts files (not .test.ts)
	testsDir := filepath.Join(tempDir, "__tests__")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	// Create a .ts file in __tests__ (should be detected as test)
	err = os.WriteFile(filepath.Join(testsDir, "app.ts"), []byte("describe('app', () => {});"), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewTypeScriptTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Should find the file in __tests__ directory
	assert.Len(t, tests, 1)
	assert.Equal(t, "app.ts", tests[0].Name)
}

// TestTypeScriptTestRunner_detectTestFramework_Priority tests framework detection priority
func TestTypeScriptTestRunner_detectTestFramework_Priority(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewTypeScriptTestRunner(log)

	// Create temp directory with both vitest and jest
	tempDir, err := os.MkdirTemp("", "ts-priority-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Package.json with both frameworks - vitest should be detected first
	packageJSON := `{
		"name": "test-project",
		"devDependencies": {
			"vitest": "^0.34.0",
			"jest": "^29.0.0"
		}
	}`
	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0o644)
	require.NoError(t, err)

	result := runner.detectTestFramework(tempDir)
	assert.Equal(t, frameworkVitest, result)
}

// TestTypeScriptTestRunner_analyzeTestFile_NonExistent tests analyzing non-existent file
func TestTypeScriptTestRunner_analyzeTestFile_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewTypeScriptTestRunner(log)

	result := runner.analyzeTestFile("/nonexistent/file.spec.ts")
	assert.Equal(t, "unit", result) // Default to unit
}
