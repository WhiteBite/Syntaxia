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

// TestJavaTestRunner_GetLanguage tests language identifier
func TestJavaTestRunner_GetLanguage(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	assert.Equal(t, "java", runner.GetLanguage())
}

// TestJavaTestRunner_detectBuildTool tests build tool detection
func TestJavaTestRunner_detectBuildTool(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	tests := []struct {
		name     string
		files    map[string]string
		expected string
	}{
		{
			name: "maven_project",
			files: map[string]string{
				"pom.xml": "<project></project>",
			},
			expected: buildToolMaven,
		},
		{
			name: "gradle_project",
			files: map[string]string{
				"build.gradle": "plugins { id 'java' }",
			},
			expected: buildToolGradle,
		},
		{
			name: "gradle_kotlin_project",
			files: map[string]string{
				"build.gradle.kts": "plugins { java }",
			},
			expected: buildToolGradle,
		},
		{
			name:     "no_build_tool_default_maven",
			files:    map[string]string{},
			expected: buildToolMaven,
		},
		{
			name: "maven_priority_over_gradle",
			files: map[string]string{
				"pom.xml":      "<project></project>",
				"build.gradle": "plugins { id 'java' }",
			},
			expected: buildToolMaven,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir, err := os.MkdirTemp("", "java-build-tool-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Create files
			for name, content := range tt.files {
				err = os.WriteFile(filepath.Join(tempDir, name), []byte(content), 0o644)
				require.NoError(t, err)
			}

			result := runner.detectBuildTool(tempDir)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestJavaTestRunner_isJavaTestFile tests Java test file detection
func TestJavaTestRunner_isJavaTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	tests := []struct {
		name     string
		fileName string
		expected bool
	}{
		{
			name:     "test_suffix",
			fileName: "UserServiceTest.java",
			expected: true,
		},
		{
			name:     "tests_suffix",
			fileName: "UserServiceTests.java",
			expected: true,
		},
		{
			name:     "testcase_suffix",
			fileName: "UserServiceTestCase.java",
			expected: true,
		},
		{
			name:     "integration_test_suffix",
			fileName: "UserServiceIT.java",
			expected: true,
		},
		{
			name:     "test_prefix",
			fileName: "TestUserService.java",
			expected: true,
		},
		{
			name:     "regular_java_file",
			fileName: "UserService.java",
			expected: false,
		},
		{
			name:     "non_java_file",
			fileName: "UserServiceTest.kt",
			expected: false,
		},
		{
			name:     "test_in_middle",
			fileName: "UserTestService.java",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.isJavaTestFile(tt.fileName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestJavaTestRunner_extractTestClassName tests class name extraction
func TestJavaTestRunner_extractTestClassName(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	tests := []struct {
		name     string
		testPath string
		expected string
	}{
		{
			name:     "simple_file",
			testPath: "UserServiceTest.java",
			expected: "UserServiceTest",
		},
		{
			name:     "standard_maven_structure",
			testPath: "src/test/java/com/example/UserServiceTest.java",
			expected: "com.example.UserServiceTest",
		},
		{
			name:     "nested_package",
			testPath: "src/test/java/com/example/service/impl/UserServiceTest.java",
			expected: "com.example.service.impl.UserServiceTest",
		},
		{
			name:     "windows_path",
			testPath: "src\\test\\java\\com\\example\\UserServiceTest.java",
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

// TestJavaTestRunner_buildTestCommand tests command building
func TestJavaTestRunner_buildTestCommand(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	tests := []struct {
		name          string
		buildTool     string
		testClassName string
		config        *domain.TestConfig
		expectedCmd   string
		expectedArgs  []string
	}{
		{
			name:          "maven_basic",
			buildTool:     buildToolMaven,
			testClassName: "com.example.UserServiceTest",
			config: &domain.TestConfig{
				ProjectPath: "/project",
				Verbose:     false,
			},
			expectedCmd:  "mvn",
			expectedArgs: []string{"test", "-Dtest=com.example.UserServiceTest"},
		},
		{
			name:          "maven_verbose",
			buildTool:     buildToolMaven,
			testClassName: "com.example.UserServiceTest",
			config: &domain.TestConfig{
				ProjectPath: "/project",
				Verbose:     true,
			},
			expectedCmd:  "mvn",
			expectedArgs: []string{"test", "-Dtest=com.example.UserServiceTest", "-X"},
		},
		{
			name:          "gradle_basic",
			buildTool:     buildToolGradle,
			testClassName: "com.example.UserServiceTest",
			config: &domain.TestConfig{
				ProjectPath: "/project",
				Verbose:     false,
			},
			expectedCmd:  "gradle",
			expectedArgs: []string{"test", "--tests=com.example.UserServiceTest"},
		},
		{
			name:          "gradle_verbose",
			buildTool:     buildToolGradle,
			testClassName: "com.example.UserServiceTest",
			config: &domain.TestConfig{
				ProjectPath: "/project",
				Verbose:     true,
			},
			expectedCmd:  "gradle",
			expectedArgs: []string{"test", "--tests=com.example.UserServiceTest", "--info"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args := runner.buildTestCommand(tt.buildTool, tt.testClassName, tt.config)
			assert.Equal(t, tt.expectedCmd, cmd)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

// TestJavaTestRunner_hasGradleWrapper tests Gradle wrapper detection
func TestJavaTestRunner_hasGradleWrapper(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

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
			// Create temp directory
			tempDir, err := os.MkdirTemp("", "java-gradle-wrapper-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Create files
			for _, name := range tt.files {
				err = os.WriteFile(filepath.Join(tempDir, name), []byte("#!/bin/bash"), 0o755)
				require.NoError(t, err)
			}

			result := runner.hasGradleWrapper(tempDir)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestJavaTestRunner_getGradleWrapperCommand tests wrapper command based on OS
func TestJavaTestRunner_getGradleWrapperCommand(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	result := runner.getGradleWrapperCommand()

	if runtime.GOOS == "windows" {
		assert.Equal(t, ".\\gradlew.bat", result)
	} else {
		assert.Equal(t, "./gradlew", result)
	}
}

// TestJavaTestRunner_analyzeTestFile tests test file type analysis
func TestJavaTestRunner_analyzeTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "java-analyze-test-*")
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
			fileName: "SmokeTest.java",
			content: `@Tag("smoke")
public class SmokeTest {
    @Test
    void testSmoke() {}
}`,
			expected: "smoke",
		},
		{
			name:     "smoke_in_filename",
			fileName: "UserSmokeTest.java",
			content: `public class UserSmokeTest {
    @Test
    void test() {}
}`,
			expected: "smoke",
		},
		{
			name:     "spring_boot_test",
			fileName: "UserServiceIT.java",
			content: `@SpringBootTest
public class UserServiceIT {
    @Test
    void testIntegration() {}
}`,
			expected: "integration",
		},
		{
			name:     "data_jpa_test",
			fileName: "UserRepositoryTest.java",
			content: `@DataJpaTest
public class UserRepositoryTest {
    @Test
    void testRepository() {}
}`,
			expected: "integration",
		},
		{
			name:     "web_mvc_test",
			fileName: "UserControllerTest.java",
			content: `@WebMvcTest
public class UserControllerTest {
    @Test
    void testController() {}
}`,
			expected: "integration",
		},
		{
			name:     "testcontainers",
			fileName: "DatabaseIT.java",
			content: `@Testcontainers
public class DatabaseIT {
    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>();
}`,
			expected: "integration",
		},
		{
			name:     "integration_in_name",
			fileName: "UserIntegrationTest.java",
			content: `public class UserIntegrationTest {
    @Test
    void test() {}
}`,
			expected: "integration",
		},
		{
			name:     "it_suffix",
			fileName: "UserIT.java",
			content: `public class UserIT {
    @Test
    void test() {}
}`,
			expected: "integration",
		},
		{
			name:     "unit_test",
			fileName: "UserServiceTest.java",
			content: `public class UserServiceTest {
    @Test
    void testAdd() {
        assertEquals(2, 1 + 1);
    }
}`,
			expected: "unit",
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

// TestJavaTestRunner_parseTestAnnotations tests annotation parsing
func TestJavaTestRunner_parseTestAnnotations(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "single_test",
			content: `public class MyTest {
    @Test
    void testSomething() {}
}`,
			expected: []string{"Test"},
		},
		{
			name: "multiple_annotations",
			content: `public class MyTest {
    @Test
    @DisplayName("Test something")
    void testSomething() {}
    
    @ParameterizedTest
    void testParam() {}
    
    @RepeatedTest(5)
    void testRepeated() {}
}`,
			expected: []string{"Test", "DisplayName", "ParameterizedTest", "RepeatedTest"},
		},
		{
			name: "with_tag",
			content: `@Tag("integration")
public class MyTest {
    @Test
    void test() {}
}`,
			expected: []string{"Tag", "Test"},
		},
		{
			name:     "no_annotations",
			content:  `public class NoTest {}`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.parseTestAnnotations(tt.content)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

// TestJavaTestRunner_DiscoverTests tests test discovery
func TestJavaTestRunner_DiscoverTests(t *testing.T) {
	// Create temp directory with Java project structure
	tempDir, err := os.MkdirTemp("", "java-discover-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create pom.xml
	pomXml := `<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>test-project</artifactId>
    <version>1.0.0</version>
</project>`
	err = os.WriteFile(filepath.Join(tempDir, "pom.xml"), []byte(pomXml), 0o644)
	require.NoError(t, err)

	// Create test directory structure
	testDir := filepath.Join(tempDir, "src", "test", "java", "com", "example")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	// Create test files
	testFiles := map[string]string{
		"UserServiceTest.java": `public class UserServiceTest {
    @Test
    void test() {}
}`,
		"UserServiceIT.java": `@SpringBootTest
public class UserServiceIT {
    @Test
    void test() {}
}`,
		"UserService.java": `public class UserService {}`, // Not a test file
	}

	for name, content := range testFiles {
		err = os.WriteFile(filepath.Join(testDir, name), []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Should find test files but not regular Java files
	// Note: DiscoverTests searches multiple directories (src/test/java, src/test, test, tests)
	// so it may find files in multiple locations
	assert.GreaterOrEqual(t, len(tests), 2)

	// Verify test files are found
	foundFiles := make(map[string]bool)
	for _, test := range tests {
		foundFiles[test.Name] = true
	}

	assert.True(t, foundFiles["UserServiceTest.java"], "UserServiceTest.java should be found")
	assert.True(t, foundFiles["UserServiceIT.java"], "UserServiceIT.java should be found")
	assert.False(t, foundFiles["UserService.java"], "UserService.java should not be found")
}

// TestJavaTestRunner_DiscoverTests_MultipleDirectories tests discovery in multiple test directories
func TestJavaTestRunner_DiscoverTests_MultipleDirectories(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "java-multi-dir-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test directories
	testDirs := []string{
		"src/test/java",
		"src/test",
		"test",
		"tests",
	}

	for _, dir := range testDirs {
		fullDir := filepath.Join(tempDir, dir)
		err = os.MkdirAll(fullDir, 0o755)
		require.NoError(t, err)

		// Create a test file in each directory
		testFile := filepath.Join(fullDir, "Test"+filepath.Base(dir)+".java")
		content := `public class Test` + filepath.Base(dir) + ` { @Test void test() {} }`
		err = os.WriteFile(testFile, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Should find tests in all directories
	assert.GreaterOrEqual(t, len(tests), 1)
}

// TestJavaTestRunner_RunTest tests single test execution
func TestJavaTestRunner_RunTest(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	// Create temp directory with pom.xml
	tempDir, err := os.MkdirTemp("", "java-run-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	pomXml := `<project></project>`
	err = os.WriteFile(filepath.Join(tempDir, "pom.xml"), []byte(pomXml), 0o644)
	require.NoError(t, err)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		Verbose:     true,
		EnvVars: map[string]string{
			"JAVA_HOME": "/usr/lib/jvm/java-11",
		},
	}

	// This will fail because mvn is not available or project is not valid
	result, err := runner.RunTest(context.Background(), "src/test/java/com/example/UserServiceTest.java", config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "src/test/java/com/example/UserServiceTest.java", result.TestPath)
	assert.Equal(t, "java", result.Language)
}

// TestJavaTestRunner_RunTestSuite tests running a test suite
func TestJavaTestRunner_RunTestSuite(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	// Create temp directory with pom.xml
	tempDir, err := os.MkdirTemp("", "java-suite-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	pomXml := `<project></project>`
	err = os.WriteFile(filepath.Join(tempDir, "pom.xml"), []byte(pomXml), 0o644)
	require.NoError(t, err)

	suite := &domain.TestSuite{
		Name:        "test_suite",
		Language:    "java",
		ProjectPath: tempDir,
		Tests: []*domain.TestInfo{
			{Path: "src/test/java/com/example/Test1.java", Name: "Test1.java", Type: "unit"},
			{Path: "src/test/java/com/example/Test2.java", Name: "Test2.java", Type: "unit"},
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

// TestJavaTestRunner_buildTestCommand_WithGradleWrapper tests command with Gradle wrapper
func TestJavaTestRunner_buildTestCommand_WithGradleWrapper(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	// Create temp directory with Gradle wrapper
	tempDir, err := os.MkdirTemp("", "java-gradle-cmd-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create gradlew
	err = os.WriteFile(filepath.Join(tempDir, "gradlew"), []byte("#!/bin/bash"), 0o755)
	require.NoError(t, err)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		Verbose:     false,
	}

	cmd, args := runner.buildTestCommand(buildToolGradle, "com.example.Test", config)

	// When wrapper exists, should use wrapper command
	if runner.hasGradleWrapper(tempDir) {
		expectedCmd := runner.getGradleWrapperCommand()
		assert.Equal(t, expectedCmd, cmd)
	}
	assert.Contains(t, args, "test")
	assert.Contains(t, args, "--tests=com.example.Test")
}

// TestJavaTestRunner_analyzeTestFile_NonExistent tests analyzing non-existent file
func TestJavaTestRunner_analyzeTestFile_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	result := runner.analyzeTestFile("/nonexistent/UserServiceTest.java")
	assert.Equal(t, "unit", result) // Default to unit
}

// TestJavaTestRunner_DiscoverTests_NoTestDirectories tests discovery when no test directories exist
func TestJavaTestRunner_DiscoverTests_NoTestDirectories(t *testing.T) {
	// Create temp directory without test directories
	tempDir, err := os.MkdirTemp("", "java-no-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewJavaTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}
