package testengine

import (
"context"
"os"
"path/filepath"
"testing"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

type testLogger struct{}

func (l *testLogger) Debug(msg string)   {}
func (l *testLogger) Info(msg string)    {}
func (l *testLogger) Warning(msg string) {}
func (l *testLogger) Error(msg string)   {}
func (l *testLogger) Fatal(msg string)   {}

func TestCppTestRunner_GetLanguage(t *testing.T) {
logger := &testLogger{}
runner := NewCppTestRunner(logger)
assert.Equal(t, "cpp", runner.GetLanguage())
}

func TestCppTestRunner_DetectBuildSystem(t *testing.T) {
logger := &testLogger{}
runner := NewCppTestRunner(logger)

tests := []struct {
name     string
setup    func(dir string) error
expected string
}{
{
name: "cmake_project",
setup: func(dir string) error {
return os.WriteFile(filepath.Join(dir, "CMakeLists.txt"), []byte("cmake"), 0o644)
},
expected: "cmake",
},
{
name: "make_project",
setup: func(dir string) error {
return os.WriteFile(filepath.Join(dir, "Makefile"), []byte("all:"), 0o644)
},
expected: "make",
},
{
name:     "direct_project",
setup:    func(dir string) error { return nil },
expected: "direct",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
tmpDir := t.TempDir()
err := tt.setup(tmpDir)
require.NoError(t, err)
result := runner.detectBuildSystem(tmpDir)
assert.Equal(t, tt.expected, result)
})
}
}

func TestCppTestRunner_IsTestFile(t *testing.T) {
logger := &testLogger{}
runner := NewCppTestRunner(logger)

tests := []struct {
name     string
filename string
expected bool
}{
{"test_suffix_cpp", "math_test.cpp", true},
{"test_suffix_cc", "math_test.cc", true},
{"regular_cpp", "math.cpp", false},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
result := runner.isTestFile(tt.filename)
assert.Equal(t, tt.expected, result)
})
}
}

func TestCppTestRunner_DiscoverTests(t *testing.T) {
logger := &testLogger{}
runner := NewCppTestRunner(logger)

tmpDir := t.TempDir()
testDir := filepath.Join(tmpDir, "test")
err := os.MkdirAll(testDir, 0o755)
require.NoError(t, err)

err = os.WriteFile(filepath.Join(testDir, "math_test.cpp"), []byte("// test"), 0o644)
require.NoError(t, err)

tests, err := runner.DiscoverTests(context.Background(), tmpDir)
require.NoError(t, err)
assert.Len(t, tests, 1)
}

func TestCppTestAnalyzer_ExtractIncludes(t *testing.T) {
logger := &testLogger{}
analyzer := NewCppTestAnalyzer(logger)

content := "#include <iostream>\n#include \"myheader.h\"\n"
includes := analyzer.extractIncludes(content)

assert.Len(t, includes, 2)
assert.Contains(t, includes, "iostream")
assert.Contains(t, includes, "myheader.h")
}

func TestCppTestAnalyzer_IsSmokeTest(t *testing.T) {
logger := &testLogger{}
analyzer := NewCppTestAnalyzer(logger)

tmpDir := t.TempDir()
smokeFile := filepath.Join(tmpDir, "smoke_test.cpp")
err := os.WriteFile(smokeFile, []byte("TEST(Smoke, Basic) { }"), 0o644)
require.NoError(t, err)

result, err := analyzer.IsSmokeTest(context.Background(), smokeFile)
require.NoError(t, err)
assert.True(t, result)
}
