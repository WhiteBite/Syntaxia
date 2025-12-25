package context

import (
	"context"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const benchContextProjectPath = "/test/project"

// Mock implementations for benchmarking
type mockFileReader struct {
	delayMs int
}

func (m *mockFileReader) ReadContents(ctx context.Context, filePaths []string, rootDir string, progress func(current, total int64)) (map[string]string, error) {
	// Simulate file reading delay
	if m.delayMs > 0 {
		time.Sleep(time.Duration(m.delayMs) * time.Millisecond)
	}

	contents := make(map[string]string)
	for _, path := range filePaths {
		contents[path] = "package main\n\nfunc main() {\n\tprintln(\"Hello World\")\n}"
	}
	return contents, nil
}

type mockTokenCounter struct {
	delayMs int
}

func (m *mockTokenCounter) CountTokens(text string) int {
	// Simulate token counting delay
	if m.delayMs > 0 {
		time.Sleep(time.Duration(m.delayMs) * time.Millisecond)
	}

	// Simple token estimation (4 characters per token)
	return len(text) / 4
}

type mockLogger struct{}

func (m *mockLogger) Debug(msg string)   {}
func (m *mockLogger) Info(msg string)    {}
func (m *mockLogger) Warning(msg string) {}
func (m *mockLogger) Error(msg string)   {}
func (m *mockLogger) Fatal(msg string)   {}

// Mock EventBus for benchmarking
type mockEventBus struct{}

func (m *mockEventBus) Emit(eventName string, data ...interface{}) {}

func BenchmarkService_BuildContext_Small(b *testing.B) {
	// Setup
	tempDir, err := os.MkdirTemp("", "context_benchmark")
	assert.NoError(b, err)
	defer os.RemoveAll(tempDir)

	mockFileReader := &mockFileReader{delayMs: 1}
	mockTokenCounter := &mockTokenCounter{delayMs: 1}
	mockLogger := &mockLogger{}
	mockBus := &mockEventBus{}

	service, err := NewService(mockFileReader, mockTokenCounter, mockBus, mockLogger)
	require.NoError(b, err)
	service.contextDir = tempDir

	b.ResetTimer()
	b.ReportAllocs()

	projectPath := benchContextProjectPath
	includedPaths := []string{"src/main.go"}
	options := &BuildOptions{
		MaxTokens:   10000,
		MaxMemoryMB: 100,
	}

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := service.BuildContext(ctx, projectPath, includedPaths, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_BuildContext_Medium(b *testing.B) {
	// Setup
	tempDir, err := os.MkdirTemp("", "context_benchmark")
	assert.NoError(b, err)
	defer os.RemoveAll(tempDir)

	mockFileReader := &mockFileReader{delayMs: 5}
	mockTokenCounter := &mockTokenCounter{delayMs: 2}
	mockLogger := &mockLogger{}
	mockBus := &mockEventBus{}

	service, err := NewService(mockFileReader, mockTokenCounter, mockBus, mockLogger)
	require.NoError(b, err)
	service.contextDir = tempDir

	b.ResetTimer()
	b.ReportAllocs()

	projectPath := benchContextProjectPath
	includedPaths := []string{
		"src/main.go",
		"src/util.go",
		"src/handler.go",
		"src/model.go",
		"src/service.go",
	}
	options := &BuildOptions{
		MaxTokens:   10000,
		MaxMemoryMB: 100,
	}

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := service.BuildContext(ctx, projectPath, includedPaths, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_BuildContext_Large(b *testing.B) {
	// Setup
	tempDir, err := os.MkdirTemp("", "context_benchmark")
	assert.NoError(b, err)
	defer os.RemoveAll(tempDir)

	mockFileReader := &mockFileReader{delayMs: 10}
	mockTokenCounter := &mockTokenCounter{delayMs: 5}
	mockLogger := &mockLogger{}
	mockBus := &mockEventBus{}

	service, err := NewService(mockFileReader, mockTokenCounter, mockBus, mockLogger)
	require.NoError(b, err)
	service.contextDir = tempDir

	b.ResetTimer()
	b.ReportAllocs()

	// Create 50 file paths to simulate a large project
	includedPaths := make([]string, 50)
	for i := 0; i < 50; i++ {
		includedPaths[i] = "src/file" + string(rune(i+'0')) + ".go"
	}

	projectPath := benchContextProjectPath
	options := &BuildOptions{
		MaxTokens:   10000,
		MaxMemoryMB: 100,
	}

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := service.BuildContext(ctx, projectPath, includedPaths, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_CreateStream(b *testing.B) {
	// Setup
	tempDir, err := os.MkdirTemp("", "context_benchmark")
	assert.NoError(b, err)
	defer os.RemoveAll(tempDir)

	mockFileReader := &mockFileReader{delayMs: 5}
	mockTokenCounter := &mockTokenCounter{delayMs: 2}
	mockLogger := &mockLogger{}
	mockBus := &mockEventBus{}

	service, err := NewService(mockFileReader, mockTokenCounter, mockBus, mockLogger)
	require.NoError(b, err)
	service.contextDir = tempDir

	b.ResetTimer()
	b.ReportAllocs()

	projectPath := benchContextProjectPath
	includedPaths := []string{
		"src/main.go",
		"src/util.go",
		"src/handler.go",
	}
	options := &BuildOptions{
		MaxTokens:   10000,
		MaxMemoryMB: 100,
	}

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := service.CreateStream(ctx, projectPath, includedPaths, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_GetContextLines(b *testing.B) {
	// Setup
	tempDir, err := os.MkdirTemp("", "context_benchmark")
	assert.NoError(b, err)
	defer os.RemoveAll(tempDir)

	mockFileReader := &mockFileReader{delayMs: 5}
	mockTokenCounter := &mockTokenCounter{delayMs: 2}
	mockLogger := &mockLogger{}
	mockBus := &mockEventBus{}

	service, err := NewService(mockFileReader, mockTokenCounter, mockBus, mockLogger)
	require.NoError(b, err)
	service.contextDir = tempDir

	// Create a stream first
	projectPath := benchContextProjectPath
	includedPaths := []string{"src/test.go"}

	// Create service with our test content
	testContent := strings.Repeat("Line content for testing.\n", 1000) // 1000 lines
	service.fileReader = &mockFileReaderWithContent{
		delayMs:       5,
		returnContent: map[string]string{"src/test.go": testContent},
	}

	options := &BuildOptions{
		MaxTokens:   10000,
		MaxMemoryMB: 100,
	}

	ctx := context.Background()
	stream, err := service.CreateStream(ctx, projectPath, includedPaths, options)
	assert.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := service.GetContextLines(ctx, stream.ID, 100, 200)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Extend mockFileReader to support returning specific content
type mockFileReaderWithContent struct {
	delayMs       int
	returnContent map[string]string
}

func (m *mockFileReaderWithContent) ReadContents(ctx context.Context, filePaths []string, rootDir string, progress func(current, total int64)) (map[string]string, error) {
	// Simulate file reading delay
	if m.delayMs > 0 {
		time.Sleep(time.Duration(m.delayMs) * time.Millisecond)
	}

	return m.returnContent, nil
}

// Add memory usage benchmark
func BenchmarkMemoryUsage_ContextBuilding(b *testing.B) {
	// Setup
	tempDir, err := os.MkdirTemp("", "context_benchmark")
	assert.NoError(b, err)
	defer os.RemoveAll(tempDir)

	mockFileReader := &mockFileReader{delayMs: 1}
	mockTokenCounter := &mockTokenCounter{delayMs: 1}
	mockLogger := &mockLogger{}
	mockBus := &mockEventBus{}

	service, err := NewService(mockFileReader, mockTokenCounter, mockBus, mockLogger)
	require.NoError(b, err)
	service.contextDir = tempDir

	b.ResetTimer()
	b.ReportAllocs()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startAllocs := m.TotalAlloc

	projectPath := benchContextProjectPath
	includedPaths := []string{"src/main.go"}
	options := &BuildOptions{
		MaxTokens:   10000,
		MaxMemoryMB: 100,
	}

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := service.BuildContext(ctx, projectPath, includedPaths, options)
		if err != nil {
			b.Fatal(err)
		}
	}

	runtime.ReadMemStats(&m)
	endAllocs := m.TotalAlloc
	b.ReportMetric(float64(endAllocs-startAllocs)/float64(b.N), "allocs/op")
}
