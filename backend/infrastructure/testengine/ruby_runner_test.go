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

// mockRubyLogger implements domain.Logger for testing
type mockRubyLogger struct{}

func (m *mockRubyLogger) Debug(msg string)                          {}
func (m *mockRubyLogger) Info(msg string)                           {}
func (m *mockRubyLogger) Warning(msg string)                        {}
func (m *mockRubyLogger) Error(msg string)                          {}
func (m *mockRubyLogger) Fatal(msg string)                          {}
func (m *mockRubyLogger) Trace(msg string)                          {}
func (m *mockRubyLogger) Print(msg string)                          {}
func (m *mockRubyLogger) SetLevel(level string)                     {}
func (m *mockRubyLogger) WithField(key string, value interface{}) domain.Logger { return m }

func TestRubyTestRunner_GetLanguage(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	if got := runner.GetLanguage(); got != "ruby" {
		t.Errorf("GetLanguage() = %v, want %v", got, "ruby")
	}
}

func TestRubyTestRunner_DetectTestFramework(t *testing.T) {
	tests := []struct {
		name        string
		testPath    string
		projectPath string
		want        string
	}{
		{
			name:        "RSpec file by suffix",
			testPath:    "spec/models/user_spec.rb",
			projectPath: "/tmp/project",
			want:        "rspec",
		},
		{
			name:        "Minitest file by suffix",
			testPath:    "test/models/user_test.rb",
			projectPath: "/tmp/project",
			want:        "minitest",
		},
		{
			name:        "RSpec by directory",
			testPath:    "spec/something.rb",
			projectPath: "/tmp/project",
			want:        "rspec",
		},
	}

	runner := NewRubyTestRunner(&mockRubyLogger{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.detectTestFramework(tt.projectPath, tt.testPath)
			if got != tt.want {
				t.Errorf("detectTestFramework() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRubyTestRunner_IsTestFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"RSpec file", "user_spec.rb", true},
		{"Minitest file", "user_test.rb", true},
		{"Regular Ruby file", "user.rb", false},
		{"Non-Ruby file", "user.py", false},
		{"Spec helper", "spec_helper.rb", false},
	}

	runner := NewRubyTestRunner(&mockRubyLogger{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.isTestFile(tt.filename)
			if got != tt.want {
				t.Errorf("isTestFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestRubyTestRunner_ShouldSkipDirectory(t *testing.T) {
	tests := []struct {
		name    string
		dirName string
		want    bool
	}{
		{"vendor directory", "vendor", true},
		{".bundle directory", ".bundle", true},
		{".git directory", ".git", true},
		{"tmp directory", "tmp", true},
		{"log directory", "log", true},
		{"coverage directory", "coverage", true},
		{"spec directory", "spec", false},
		{"test directory", "test", false},
		{"app directory", "app", false},
	}

	runner := NewRubyTestRunner(&mockRubyLogger{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.shouldSkipDirectory(tt.dirName)
			if got != tt.want {
				t.Errorf("shouldSkipDirectory(%q) = %v, want %v", tt.dirName, got, tt.want)
			}
		})
	}
}

func TestRubyTestRunner_AnalyzeTestFile(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "ruby_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		filename string
		content  string
		want     string
	}{
		{
			name:     "Unit test",
			filename: "user_spec.rb",
			content:  "describe User do\n  it 'validates name' do\n  end\nend",
			want:     "unit",
		},
		{
			name:     "Integration test",
			filename: "integration_spec.rb",
			content:  "describe 'Integration' do\n  it 'works' do\n  end\nend",
			want:     "integration",
		},
		{
			name:     "Feature test with Capybara",
			filename: "login_spec.rb",
			content:  "require 'capybara'\nfeature 'Login' do\n  scenario 'user logs in' do\n  end\nend",
			want:     "feature",
		},
		{
			name:     "Request test",
			filename: "api_request_spec.rb",
			content:  "describe 'API' do\n  it 'returns 200' do\n  end\nend",
			want:     "request",
		},
	}

	runner := NewRubyTestRunner(&mockRubyLogger{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			filePath := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(filePath, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			got := runner.analyzeTestFile(filePath)
			if got != tt.want {
				t.Errorf("analyzeTestFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRubyTestRunner_HasBundler(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "ruby_bundler_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	runner := NewRubyTestRunner(&mockRubyLogger{})

	// Test without Gemfile
	if runner.hasBundler(tmpDir) {
		t.Error("hasBundler() should return false without Gemfile")
	}

	// Create Gemfile
	gemfilePath := filepath.Join(tmpDir, "Gemfile")
	if err := os.WriteFile(gemfilePath, []byte("source 'https://rubygems.org'\n"), 0o644); err != nil {
		t.Fatalf("Failed to write Gemfile: %v", err)
	}

	// Test with Gemfile
	if !runner.hasBundler(tmpDir) {
		t.Error("hasBundler() should return true with Gemfile")
	}
}

func TestRubyTestRunner_DetectProjectFramework(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "ruby_framework_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	runner := NewRubyTestRunner(&mockRubyLogger{})

	tests := []struct {
		name         string
		gemfileContent string
		createSpecDir bool
		want         string
	}{
		{
			name:           "RSpec in Gemfile",
			gemfileContent: "gem 'rspec'\n",
			createSpecDir:  false,
			want:           "rspec",
		},
		{
			name:           "Minitest in Gemfile",
			gemfileContent: "gem 'minitest'\n",
			createSpecDir:  false,
			want:           "minitest",
		},
		{
			name:           "No framework but spec dir exists",
			gemfileContent: "gem 'rails'\n",
			createSpecDir:  true,
			want:           "rspec",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			os.RemoveAll(filepath.Join(tmpDir, "Gemfile"))
			os.RemoveAll(filepath.Join(tmpDir, "spec"))

			// Create Gemfile
			gemfilePath := filepath.Join(tmpDir, "Gemfile")
			if err := os.WriteFile(gemfilePath, []byte(tt.gemfileContent), 0o644); err != nil {
				t.Fatalf("Failed to write Gemfile: %v", err)
			}

			// Create spec directory if needed
			if tt.createSpecDir {
				if err := os.MkdirAll(filepath.Join(tmpDir, "spec"), 0o755); err != nil {
					t.Fatalf("Failed to create spec dir: %v", err)
				}
			}

			got := runner.detectProjectFramework(tmpDir)
			if got != tt.want {
				t.Errorf("detectProjectFramework() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRubyTestRunner_ParseRSpecJSON(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	tests := []struct {
		name       string
		output     string
		wantCount  int
		wantPassed bool
	}{
		{
			name: "Successful RSpec JSON",
			output: `{"examples":[{"id":"./spec/user_spec.rb[1:1]","description":"validates name","full_description":"User validates name","status":"passed","file_path":"./spec/user_spec.rb","line_number":5,"run_time":0.001}]}`,
			wantCount:  1,
			wantPassed: true,
		},
		{
			name: "Failed RSpec JSON",
			output: `{"examples":[{"id":"./spec/user_spec.rb[1:1]","description":"validates name","full_description":"User validates name","status":"failed","file_path":"./spec/user_spec.rb","line_number":5,"run_time":0.001,"exception":{"class":"RSpec::Expectations::ExpectationNotMetError","message":"expected true, got false"}}]}`,
			wantCount:  1,
			wantPassed: false,
		},
		{
			name:       "Invalid JSON",
			output:     "not json",
			wantCount:  0,
			wantPassed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := runner.parseRSpecJSON(tt.output, 1.0)

			if len(results) != tt.wantCount {
				t.Errorf("parseRSpecJSON() returned %d results, want %d", len(results), tt.wantCount)
				return
			}

			if tt.wantCount > 0 && results[0].Success != tt.wantPassed {
				t.Errorf("parseRSpecJSON() result.Success = %v, want %v", results[0].Success, tt.wantPassed)
			}
		})
	}
}

func TestRubyTestRunner_ParseMinitestOutput(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	tests := []struct {
		name       string
		output     string
		testCount  int
		wantPassed bool
	}{
		{
			name:       "All tests passed",
			output:     "10 runs, 20 assertions, 0 failures, 0 errors",
			testCount:  2,
			wantPassed: true,
		},
		{
			name:       "Some tests failed",
			output:     "10 runs, 20 assertions, 2 failures, 0 errors",
			testCount:  2,
			wantPassed: false,
		},
		{
			name:       "Some tests errored",
			output:     "10 runs, 20 assertions, 0 failures, 1 errors",
			testCount:  2,
			wantPassed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testInfos := make([]*domain.TestInfo, tt.testCount)
			for i := 0; i < tt.testCount; i++ {
				testInfos[i] = &domain.TestInfo{
					Path: "test/test_file.rb",
					Name: "test_file.rb",
				}
			}

			results := runner.parseMinitestOutput(tt.output, testInfos, 1.0)

			if len(results) != tt.testCount {
				t.Errorf("parseMinitestOutput() returned %d results, want %d", len(results), tt.testCount)
				return
			}

			for _, result := range results {
				if result.Success != tt.wantPassed {
					t.Errorf("parseMinitestOutput() result.Success = %v, want %v", result.Success, tt.wantPassed)
				}
			}
		})
	}
}

func TestRubyTestRunner_DiscoverTests(t *testing.T) {
	// Create temp directory structure
	tmpDir, err := os.MkdirTemp("", "ruby_discover_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create spec directory
	specDir := filepath.Join(tmpDir, "spec")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("Failed to create spec dir: %v", err)
	}

	// Create test files
	testFiles := []string{
		"user_spec.rb",
		"order_spec.rb",
	}

	for _, f := range testFiles {
		filePath := filepath.Join(specDir, f)
		if err := os.WriteFile(filePath, []byte("describe 'Test' do\nend"), 0o644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
	}

	// Create non-test file
	nonTestFile := filepath.Join(specDir, "spec_helper.rb")
	if err := os.WriteFile(nonTestFile, []byte("# helper"), 0o644); err != nil {
		t.Fatalf("Failed to write non-test file: %v", err)
	}

	runner := NewRubyTestRunner(&mockRubyLogger{})
	tests, err := runner.DiscoverTests(context.Background(), tmpDir)

	if err != nil {
		t.Fatalf("DiscoverTests() error = %v", err)
	}

	if len(tests) != len(testFiles) {
		t.Errorf("DiscoverTests() found %d tests, want %d", len(tests), len(testFiles))
	}

	// Verify all discovered tests are spec files
	for _, test := range tests {
		if test.Metadata["framework"] != "rspec" {
			t.Errorf("Expected framework 'rspec', got %q", test.Metadata["framework"])
		}
	}
}

func TestRubyTestRunner_BuildRSpecCommand(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	tests := []struct {
		name       string
		testPath   string
		useBundler bool
		wantCmd    string
		wantArgs   int
	}{
		{
			name:       "Without bundler",
			testPath:   "spec/user_spec.rb",
			useBundler: false,
			wantCmd:    "rspec",
			wantArgs:   5, // --format json --format progress spec/user_spec.rb
		},
		{
			name:       "With bundler",
			testPath:   "spec/user_spec.rb",
			useBundler: true,
			wantCmd:    "bundle",
			wantArgs:   7, // exec rspec --format json --format progress spec/user_spec.rb
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &domain.TestConfig{
				ProjectPath: "/tmp/project",
			}

			cmdName, args := runner.buildRSpecCommand(tt.testPath, config, tt.useBundler)

			if cmdName != tt.wantCmd {
				t.Errorf("buildRSpecCommand() cmdName = %v, want %v", cmdName, tt.wantCmd)
			}

			if len(args) != tt.wantArgs {
				t.Errorf("buildRSpecCommand() args count = %v, want %v", len(args), tt.wantArgs)
			}
		})
	}
}

// TestRubyTestRunner_BuildMinitestCommand tests Minitest command building
func TestRubyTestRunner_BuildMinitestCommand(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	tests := []struct {
		name       string
		testPath   string
		useBundler bool
		wantCmd    string
	}{
		{
			name:       "Without bundler",
			testPath:   "test/user_test.rb",
			useBundler: false,
			wantCmd:    "ruby",
		},
		{
			name:       "With bundler",
			testPath:   "test/user_test.rb",
			useBundler: true,
			wantCmd:    "bundle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &domain.TestConfig{
				ProjectPath: "/tmp/project",
			}

			cmdName, args := runner.buildMinitestCommand(tt.testPath, config, tt.useBundler)

			if cmdName != tt.wantCmd {
				t.Errorf("buildMinitestCommand() cmdName = %v, want %v", cmdName, tt.wantCmd)
			}

			// Verify test path is in args
			found := false
			for _, arg := range args {
				if arg == tt.testPath {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("buildMinitestCommand() args should contain test path %v", tt.testPath)
			}
		})
	}
}

// TestRubyTestRunner_RunTest tests single test execution
func TestRubyTestRunner_RunTest(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	tmpDir, err := os.MkdirTemp("", "ruby_run_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	config := &domain.TestConfig{
		ProjectPath: tmpDir,
		Verbose:     false,
		Timeout:     60,
		EnvVars: map[string]string{
			"RAILS_ENV": "test",
		},
	}

	result, err := runner.RunTest(context.Background(), "spec/user_spec.rb", config)

	if err != nil {
		t.Fatalf("RunTest() error = %v", err)
	}

	if result == nil {
		t.Fatal("RunTest() returned nil result")
	}

	if result.TestPath != "spec/user_spec.rb" {
		t.Errorf("RunTest() TestPath = %v, want %v", result.TestPath, "spec/user_spec.rb")
	}

	if result.Language != "ruby" {
		t.Errorf("RunTest() Language = %v, want %v", result.Language, "ruby")
	}
}

// TestRubyTestRunner_RunTestSuite tests suite execution
func TestRubyTestRunner_RunTestSuite(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	tmpDir, err := os.MkdirTemp("", "ruby_suite_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	suite := &domain.TestSuite{
		Name:        "test_suite",
		Language:    "ruby",
		ProjectPath: tmpDir,
		Tests: []*domain.TestInfo{
			{Path: "spec/user_spec.rb", Name: "user_spec.rb", Type: "unit"},
			{Path: "spec/order_spec.rb", Name: "order_spec.rb", Type: "unit"},
		},
		Config: &domain.TestConfig{
			ProjectPath: tmpDir,
			Verbose:     false,
		},
	}

	results, err := runner.RunTestSuite(context.Background(), suite)

	if err != nil {
		t.Fatalf("RunTestSuite() error = %v", err)
	}

	if results == nil {
		t.Fatal("RunTestSuite() returned nil results")
	}

	if len(results) != 2 {
		t.Errorf("RunTestSuite() returned %d results, want 2", len(results))
	}
}

// TestRubyTestRunner_ParseRSpecJSON_WithBacktrace tests RSpec JSON parsing with backtrace
func TestRubyTestRunner_ParseRSpecJSON_WithBacktrace(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	output := `{"examples":[{"id":"./spec/user_spec.rb[1:1]","description":"validates name","full_description":"User validates name","status":"failed","file_path":"./spec/user_spec.rb","line_number":5,"run_time":0.001,"exception":{"class":"RSpec::Expectations::ExpectationNotMetError","message":"expected true, got false","backtrace":["./spec/user_spec.rb:10","./spec/spec_helper.rb:5"]}}]}`

	results := runner.parseRSpecJSON(output, 1.0)

	if len(results) != 1 {
		t.Fatalf("parseRSpecJSON() returned %d results, want 1", len(results))
	}

	if results[0].Success {
		t.Error("parseRSpecJSON() result.Success should be false")
	}

	if results[0].Error == "" {
		t.Error("parseRSpecJSON() result.Error should not be empty")
	}

	if results[0].Output == "" {
		t.Error("parseRSpecJSON() result.Output should contain backtrace")
	}
}

// TestRubyTestRunner_DiscoverTests_EmptyProject tests discovery in empty project
func TestRubyTestRunner_DiscoverTests_EmptyProject(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ruby_empty_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	runner := NewRubyTestRunner(&mockRubyLogger{})
	tests, err := runner.DiscoverTests(context.Background(), tmpDir)

	if err != nil {
		t.Fatalf("DiscoverTests() error = %v", err)
	}

	if len(tests) != 0 {
		t.Errorf("DiscoverTests() found %d tests, want 0", len(tests))
	}
}

// TestRubyTestRunner_DiscoverTests_TestDirectory tests discovery in test directory
func TestRubyTestRunner_DiscoverTests_TestDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ruby_test_dir_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test directory
	testDir := filepath.Join(tmpDir, "test")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}

	// Create minitest files
	testFiles := []string{
		"user_test.rb",
		"order_test.rb",
	}

	for _, f := range testFiles {
		filePath := filepath.Join(testDir, f)
		if err := os.WriteFile(filePath, []byte("class UserTest < Minitest::Test\nend"), 0o644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
	}

	runner := NewRubyTestRunner(&mockRubyLogger{})
	tests, err := runner.DiscoverTests(context.Background(), tmpDir)

	if err != nil {
		t.Fatalf("DiscoverTests() error = %v", err)
	}

	if len(tests) != len(testFiles) {
		t.Errorf("DiscoverTests() found %d tests, want %d", len(tests), len(testFiles))
	}

	// Verify framework is minitest
	for _, test := range tests {
		if test.Metadata["framework"] != "minitest" {
			t.Errorf("Expected framework 'minitest', got %q", test.Metadata["framework"])
		}
	}
}

// TestRubyTestRunner_BuildTestCommand tests test command building
func TestRubyTestRunner_BuildTestCommand(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	tests := []struct {
		name      string
		testPath  string
		framework string
		wantCmd   string
	}{
		{
			name:      "RSpec test",
			testPath:  "spec/user_spec.rb",
			framework: "rspec",
			wantCmd:   "rspec",
		},
		{
			name:      "Minitest test",
			testPath:  "test/user_test.rb",
			framework: "minitest",
			wantCmd:   "ruby",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &domain.TestConfig{
				ProjectPath: "/tmp/project",
			}

			cmdName, _ := runner.buildTestCommand(tt.testPath, config, tt.framework)

			if cmdName != tt.wantCmd {
				t.Errorf("buildTestCommand() cmdName = %v, want %v", cmdName, tt.wantCmd)
			}
		})
	}
}


// TestRubyTestAnalyzer_AnalyzeOutput tests output analysis
func TestRubyTestAnalyzer_AnalyzeOutput(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	tests := []struct {
		name          string
		output        string
		expectSuccess bool
		expectCount   int
	}{
		{
			name: "rspec_json_passed",
			output: `{"examples":[{"id":"./spec/user_spec.rb[1:1]","description":"test","full_description":"User test","status":"passed","file_path":"./spec/user_spec.rb","line_number":5,"run_time":0.001}]}`,
			expectSuccess: true,
			expectCount:   1,
		},
		{
			name:          "minitest_passed",
			output:        "10 runs, 20 assertions, 0 failures, 0 errors",
			expectSuccess: true,
			expectCount:   1,
		},
		{
			name:          "minitest_failed",
			output:        "10 runs, 20 assertions, 2 failures, 1 errors",
			expectSuccess: false,
			expectCount:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := analyzer.AnalyzeOutput(tt.output)
			if err != nil {
				t.Fatalf("AnalyzeOutput() error = %v", err)
			}

			if len(results) < tt.expectCount {
				t.Errorf("AnalyzeOutput() returned %d results, want at least %d", len(results), tt.expectCount)
			}

			if len(results) > 0 && results[0].Success != tt.expectSuccess {
				t.Errorf("AnalyzeOutput() success = %v, want %v", results[0].Success, tt.expectSuccess)
			}
		})
	}
}

// TestRubyTestAnalyzer_ExtractFailures tests failure extraction
func TestRubyTestAnalyzer_ExtractFailures(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	output := `Failure: UserTest#test_create [test/user_test.rb:10]:
Expected true, got false

Error: OrderTest#test_process [test/order_test.rb:20]:
NoMethodError: undefined method`

	failures := analyzer.ExtractFailures(output)

	if len(failures) < 2 {
		t.Errorf("ExtractFailures() returned %d failures, want at least 2", len(failures))
	}
}

// TestRubyTestAnalyzer_GetTestCoverage tests coverage extraction
func TestRubyTestAnalyzer_GetTestCoverage(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	tests := []struct {
		name     string
		output   string
		expected float64
	}{
		{
			name:     "simplecov_format",
			output:   "Coverage report generated... 85.5% covered",
			expected: 85.5,
		},
		{
			name:     "line_coverage_format",
			output:   "Line Coverage: 92.3%",
			expected: 92.3,
		},
		{
			name:     "no_coverage",
			output:   "Tests completed",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coverage := analyzer.GetTestCoverage(tt.output)
			if coverage.Percentage != tt.expected {
				t.Errorf("GetTestCoverage() = %v, want %v", coverage.Percentage, tt.expected)
			}
		})
	}
}

// TestRubyTestAnalyzer_AnalyzeTestDependencies tests dependency analysis
func TestRubyTestAnalyzer_AnalyzeTestDependencies(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	tmpDir, err := os.MkdirTemp("", "ruby_deps_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	content := `require 'rails_helper'
require_relative '../support/helpers'
require 'factory_bot'

describe User do
end`

	testFile := filepath.Join(tmpDir, "user_spec.rb")
	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	deps, err := analyzer.AnalyzeTestDependencies(context.Background(), testFile)
	if err != nil {
		t.Fatalf("AnalyzeTestDependencies() error = %v", err)
	}

	if len(deps) != 3 {
		t.Errorf("AnalyzeTestDependencies() returned %d deps, want 3", len(deps))
	}
}

// TestRubyTestAnalyzer_IsSmokeTest tests smoke test detection
func TestRubyTestAnalyzer_IsSmokeTest(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	tmpDir, err := os.MkdirTemp("", "ruby_smoke_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected bool
	}{
		{
			name:     "smoke_in_filename",
			fileName: "smoke_spec.rb",
			content:  "describe 'Smoke' do end",
			expected: true,
		},
		{
			name:     "smoke_tag",
			fileName: "api_spec.rb",
			content:  "describe 'API', :smoke do end",
			expected: true,
		},
		{
			name:     "regular_test",
			fileName: "user_spec.rb",
			content:  "describe User do end",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, tt.fileName)
			if err := os.WriteFile(testFile, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			result, err := analyzer.IsSmokeTest(context.Background(), testFile)
			if err != nil {
				t.Fatalf("IsSmokeTest() error = %v", err)
			}

			if result != tt.expected {
				t.Errorf("IsSmokeTest() = %v, want %v", result, tt.expected)
			}
		})
	}
}


// TestRubyTestAnalyzer_FindTestsForFile tests finding tests for a source file
func TestRubyTestAnalyzer_FindTestsForFile(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	tmpDir, err := os.MkdirTemp("", "ruby_find_tests_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file
	libDir := filepath.Join(tmpDir, "lib")
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		t.Fatalf("Failed to create lib dir: %v", err)
	}

	srcFile := filepath.Join(libDir, "user.rb")
	if err := os.WriteFile(srcFile, []byte("class User; end"), 0o644); err != nil {
		t.Fatalf("Failed to write source file: %v", err)
	}

	// Create spec directory and test file
	specDir := filepath.Join(tmpDir, "spec")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("Failed to create spec dir: %v", err)
	}

	testFile := filepath.Join(specDir, "user_spec.rb")
	if err := os.WriteFile(testFile, []byte("describe User do; end"), 0o644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	tests, err := analyzer.FindTestsForFile(context.Background(), "lib/user.rb", tmpDir)
	if err != nil {
		t.Fatalf("FindTestsForFile() error = %v", err)
	}

	if len(tests) < 1 {
		t.Errorf("FindTestsForFile() found %d tests, want at least 1", len(tests))
	}
}

// TestRubyTestAnalyzer_ExtractRequires tests require extraction
func TestRubyTestAnalyzer_ExtractRequires(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	content := `require 'rails_helper'
require_relative '../support/helpers'
require 'factory_bot'
require "json"

describe User do
end`

	requires := analyzer.extractRequires(content)

	if len(requires) != 4 {
		t.Errorf("extractRequires() returned %d requires, want 4", len(requires))
	}
}


// TestRubyTestAnalyzer_AnalyzeTestDependencies_NonExistent tests error handling
func TestRubyTestAnalyzer_AnalyzeTestDependencies_NonExistent(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	_, err := analyzer.AnalyzeTestDependencies(context.Background(), "/nonexistent/file.rb")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// TestRubyTestAnalyzer_IsSmokeTest_NonExistent tests error handling
func TestRubyTestAnalyzer_IsSmokeTest_NonExistent(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	_, err := analyzer.IsSmokeTest(context.Background(), "/nonexistent/file.rb")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// TestRubyTestAnalyzer_GetLanguage tests language identifier
func TestRubyTestAnalyzer_GetLanguage(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	if got := analyzer.GetLanguage(); got != "ruby" {
		t.Errorf("GetLanguage() = %v, want %v", got, "ruby")
	}
}

// TestRubyTestAnalyzer_ParseRSpecText tests RSpec text output parsing
func TestRubyTestAnalyzer_ParseRSpecText(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	tests := []struct {
		name          string
		output        string
		expectResults bool
	}{
		{
			name:          "all_passed",
			output:        "10 examples, 0 failures",
			expectResults: true,
		},
		{
			name:          "with_failures",
			output:        "rspec ./spec/user_spec.rb:10 # User validates name",
			expectResults: true,
		},
		{
			name:          "no_match",
			output:        "Running tests...",
			expectResults: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := analyzer.parseRSpecText(tt.output)
			if tt.expectResults && len(results) == 0 {
				t.Error("Expected results but got none")
			}
		})
	}
}

// TestRubyTestAnalyzer_CreateFallbackResult tests fallback result creation
func TestRubyTestAnalyzer_CreateFallbackResult(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	tests := []struct {
		name          string
		output        string
		expectSuccess bool
	}{
		{
			name:          "success_output",
			output:        "0 failures, 0 errors",
			expectSuccess: true,
		},
		{
			name:          "failure_output",
			output:        "1 failure occurred",
			expectSuccess: false,
		},
		{
			name:          "error_output",
			output:        "error: something went wrong",
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := analyzer.createFallbackResult(tt.output)
			if len(results) == 0 {
				t.Fatal("Expected at least one result")
			}
			if results[0].Success != tt.expectSuccess {
				t.Errorf("Success = %v, want %v", results[0].Success, tt.expectSuccess)
			}
		})
	}
}


// TestRubyTestRunner_parseRSpecJSON_EmptyOutput tests parsing empty RSpec JSON output
func TestRubyTestRunner_parseRSpecJSON_EmptyOutput(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	results := runner.parseRSpecJSON("", 1.0)
	assert.Empty(t, results)
}

// TestRubyTestRunner_parseRSpecJSON_MalformedJSON tests parsing malformed JSON
func TestRubyTestRunner_parseRSpecJSON_MalformedJSON(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	results := runner.parseRSpecJSON("{invalid json", 1.0)
	assert.Empty(t, results)
}

// TestRubyTestRunner_parseRSpecJSON_NoExamples tests parsing JSON with no examples
func TestRubyTestRunner_parseRSpecJSON_NoExamples(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	output := `{"version":"3.12.0","examples":[]}`
	results := runner.parseRSpecJSON(output, 1.0)
	assert.Empty(t, results)
}

// TestRubyTestRunner_parseMinitestOutput_EmptyOutput tests parsing empty Minitest output
func TestRubyTestRunner_parseMinitestOutput_EmptyOutput(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	testInfos := []*domain.TestInfo{
		{Path: "test/user_test.rb", Name: "user_test.rb"},
	}

	results := runner.parseMinitestOutput("", testInfos, 1.0)
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
}

// TestRubyTestRunner_parseMinitestOutput_NoSummary tests parsing output without summary
func TestRubyTestRunner_parseMinitestOutput_NoSummary(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	testInfos := []*domain.TestInfo{
		{Path: "test/user_test.rb", Name: "user_test.rb"},
	}

	results := runner.parseMinitestOutput("Running tests...", testInfos, 1.0)
	require.Len(t, results, 1)
	assert.False(t, results[0].Success)
}

// TestRubyTestRunner_analyzeTestFile_NonExistent tests analyzing non-existent file
func TestRubyTestRunner_analyzeTestFile_NonExistent(t *testing.T) {
	runner := NewRubyTestRunner(&mockRubyLogger{})

	result := runner.analyzeTestFile("/nonexistent/file.rb")
	assert.Equal(t, "unit", result)
}

// TestRubyTestRunner_analyzeTestFile_ControllerTest tests controller test detection
func TestRubyTestRunner_analyzeTestFile_ControllerTest(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ruby_controller_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	runner := NewRubyTestRunner(&mockRubyLogger{})

	content := `describe UsersController do
  it 'returns users' do
  end
end`
	testFile := filepath.Join(tmpDir, "users_controller_spec.rb")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	result := runner.analyzeTestFile(testFile)
	assert.Equal(t, "request", result)
}

// TestRubyTestRunner_detectProjectFramework_NoGemfile tests framework detection without Gemfile
func TestRubyTestRunner_detectProjectFramework_NoGemfile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ruby_no_gemfile_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	runner := NewRubyTestRunner(&mockRubyLogger{})

	// Without Gemfile and without spec directory
	result := runner.detectProjectFramework(tmpDir)
	assert.Equal(t, "minitest", result)
}

// TestRubyTestRunner_detectProjectFramework_WithSpecDir tests framework detection with spec directory
func TestRubyTestRunner_detectProjectFramework_WithSpecDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ruby_spec_dir_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	runner := NewRubyTestRunner(&mockRubyLogger{})

	// Create spec directory
	specDir := filepath.Join(tmpDir, "spec")
	err = os.MkdirAll(specDir, 0o755)
	require.NoError(t, err)

	result := runner.detectProjectFramework(tmpDir)
	assert.Equal(t, "rspec", result)
}

// TestRubyTestRunner_DiscoverTests_SkipCoverage tests that coverage directory is skipped
func TestRubyTestRunner_DiscoverTests_SkipCoverage(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ruby_skip_coverage_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	runner := NewRubyTestRunner(&mockRubyLogger{})

	// Create spec directory
	specDir := filepath.Join(tmpDir, "spec")
	err = os.MkdirAll(specDir, 0o755)
	require.NoError(t, err)

	// Create coverage directory inside spec (should be skipped)
	coverageDir := filepath.Join(specDir, "coverage")
	err = os.MkdirAll(coverageDir, 0o755)
	require.NoError(t, err)

	// Create test file in coverage (should not be found)
	coverageTest := filepath.Join(coverageDir, "coverage_spec.rb")
	err = os.WriteFile(coverageTest, []byte("describe 'Coverage' do end"), 0o644)
	require.NoError(t, err)

	// Create valid test file
	validTest := filepath.Join(specDir, "valid_spec.rb")
	err = os.WriteFile(validTest, []byte("describe 'Valid' do end"), 0o644)
	require.NoError(t, err)

	tests, err := runner.DiscoverTests(context.Background(), tmpDir)
	require.NoError(t, err)
	assert.Len(t, tests, 1)
	assert.Equal(t, "valid_spec.rb", tests[0].Name)
}

// TestRubyTestAnalyzer_AnalyzeOutput_EmptyOutput tests analyzing empty output
func TestRubyTestAnalyzer_AnalyzeOutput_EmptyOutput(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	results, err := analyzer.AnalyzeOutput("")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
}

// TestRubyTestAnalyzer_ExtractFailures_EmptyOutput tests extracting failures from empty output
func TestRubyTestAnalyzer_ExtractFailures_EmptyOutput(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	failures := analyzer.ExtractFailures("")
	assert.Empty(t, failures)
}

// TestRubyTestAnalyzer_GetTestCoverage_NoCoverage tests coverage extraction with no coverage info
func TestRubyTestAnalyzer_GetTestCoverage_NoCoverage(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	coverage := analyzer.GetTestCoverage("Tests completed successfully")
	assert.Equal(t, float64(0), coverage.Percentage)
}

// TestRubyTestAnalyzer_ExtractRequires_EmptyContent tests extracting requires from empty content
func TestRubyTestAnalyzer_ExtractRequires_EmptyContent(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	requires := analyzer.extractRequires("")
	assert.Empty(t, requires)
}

// TestRubyTestAnalyzer_ExtractRequires_DuplicateRequires tests duplicate require handling
func TestRubyTestAnalyzer_ExtractRequires_DuplicateRequires(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	content := `require 'rails_helper'
require 'rails_helper'
require 'factory_bot'`

	requires := analyzer.extractRequires(content)
	assert.Len(t, requires, 2)
}

// TestRubyTestAnalyzer_FindTestsForFile_NoTestDirs tests when no test directories exist
func TestRubyTestAnalyzer_FindTestsForFile_NoTestDirs(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	tmpDir, err := os.MkdirTemp("", "ruby_no_test_dirs_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create source file without test directories
	libDir := filepath.Join(tmpDir, "lib")
	err = os.MkdirAll(libDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(libDir, "user.rb")
	err = os.WriteFile(srcFile, []byte("class User; end"), 0o644)
	require.NoError(t, err)

	tests, err := analyzer.FindTestsForFile(context.Background(), "lib/user.rb", tmpDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestRubyTestAnalyzer_IsSmokeTest_EmptyFile tests smoke test detection with empty file
func TestRubyTestAnalyzer_IsSmokeTest_EmptyFile(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	tmpDir, err := os.MkdirTemp("", "ruby_empty_smoke_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "empty_spec.rb")
	err = os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	result, err := analyzer.IsSmokeTest(context.Background(), testFile)
	require.NoError(t, err)
	assert.False(t, result)
}

// TestRubyTestAnalyzer_ParseRSpecText_NoFailures tests parsing RSpec text with no failures
func TestRubyTestAnalyzer_ParseRSpecText_NoFailures(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	output := "10 examples, 0 failures\nFinished in 1.5 seconds"
	results := analyzer.parseRSpecText(output)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
}

// TestRubyTestAnalyzer_CreateFallbackResult_NoIndicators tests fallback with no indicators
func TestRubyTestAnalyzer_CreateFallbackResult_NoIndicators(t *testing.T) {
	analyzer := NewRubyTestAnalyzer(&mockRubyLogger{})

	results := analyzer.createFallbackResult("Running tests...")
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
}
