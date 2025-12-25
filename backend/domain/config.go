package domain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Memory limits
const (
	DefaultMaxMemoryMB = 30
	DefaultMaxTokens   = 5000
)

// Worker pool limits
const (
	MaxWorkerCount = 16
	MinWorkerCount = 2
)

// Timeouts
const (
	DefaultOperationTimeout = 30 * time.Second
	DefaultAITimeout        = 60 * time.Second
	DefaultBuildTimeout     = 5 * time.Minute
)

// File limits
const (
	MaxFileSizeBytes = 10 * 1024 * 1024 // 10MB
	MaxFilesPerBatch = 100
)

// ConfigLoader interface for loading configuration
type ConfigLoader interface {
	Load(configPath string) (*Config, error)
	Save(config *Config, configPath string) error
}

// Config holds application configuration
type Config struct {
	// Tool execution settings
	Tools ToolsConfig `json:"tools"`

	// Semantic search settings
	SemanticSearch SemanticSearchConfig `json:"semanticSearch"`

	// Symbol index settings
	SymbolIndex SymbolIndexConfig `json:"symbolIndex"`

	// Call graph settings
	CallGraph CallGraphConfig `json:"callGraph"`

	// OpenAI settings
	OpenAI OpenAIConfig `json:"openai"`

	// Agentic chat settings
	AgenticChat AgenticChatConfig `json:"agenticChat"`
}

// ToolsConfig holds tool execution settings
type ToolsConfig struct {
	// Maximum number of search results
	MaxSearchResults int `json:"maxSearchResults"`

	// Maximum file size to read (bytes)
	MaxFileSize int64 `json:"maxFileSize"`

	// Maximum diff output length
	MaxDiffLength int `json:"maxDiffLength"`

	// Directories to skip during search
	SkipDirectories []string `json:"skipDirectories"`
}

// SemanticSearchConfig holds semantic search settings
type SemanticSearchConfig struct {
	// Batch size for embedding generation
	BatchSize int `json:"batchSize"`

	// Maximum tokens per chunk
	MaxChunkTokens int `json:"maxChunkTokens"`

	// Minimum chunk tokens (skip smaller chunks)
	MinChunkTokens int `json:"minChunkTokens"`

	// Default top-k results
	DefaultTopK int `json:"defaultTopK"`

	// Default minimum score threshold
	DefaultMinScore float64 `json:"defaultMinScore"`

	// RRF constant for hybrid search
	RRFConstant int `json:"rrfConstant"`
}

// SymbolIndexConfig holds symbol index settings
type SymbolIndexConfig struct {
	// Enable SQLite caching
	EnableCache bool `json:"enableCache"`

	// Cache directory (relative to project root)
	CacheDir string `json:"cacheDir"`

	// Maximum symbols to return in search
	MaxSearchResults int `json:"maxSearchResults"`
}

// CallGraphConfig holds call graph settings
type CallGraphConfig struct {
	// Maximum depth for transitive dependencies
	MaxTransitiveDepth int `json:"maxTransitiveDepth"`

	// Maximum nodes in Mermaid export
	MaxMermaidNodes int `json:"maxMermaidNodes"`

	// Maximum commits to analyze for co-changed files
	MaxCommitsForCoChanged int `json:"maxCommitsForCoChanged"`
}

// OpenAIConfig holds OpenAI API settings
type OpenAIConfig struct {
	// Maximum retries for API calls
	MaxRetries int `json:"maxRetries"`

	// Base delay for exponential backoff (milliseconds)
	BaseDelayMs int `json:"baseDelayMs"`

	// Maximum delay for exponential backoff (milliseconds)
	MaxDelayMs int `json:"maxDelayMs"`

	// Request timeout (seconds)
	TimeoutSeconds int `json:"timeoutSeconds"`
}

// AgenticChatConfig holds agentic chat settings
type AgenticChatConfig struct {
	// Maximum tool call iterations
	MaxIterations int `json:"maxIterations"`

	// Maximum tool log length in response
	MaxToolLogLength int `json:"maxToolLogLength"`
}

// DefaultConfig returns configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Tools: ToolsConfig{
			MaxSearchResults: 50,
			MaxFileSize:      10 * 1024 * 1024, // 10MB
			MaxDiffLength:    5000,
			SkipDirectories:  []string{"node_modules", ".git", "vendor", "dist", "build", ".idea", ".vscode"},
		},
		SemanticSearch: SemanticSearchConfig{
			BatchSize:       10,
			MaxChunkTokens:  512,
			MinChunkTokens:  20,
			DefaultTopK:     10,
			DefaultMinScore: 0.5,
			RRFConstant:     60,
		},
		SymbolIndex: SymbolIndexConfig{
			EnableCache:      true,
			CacheDir:         ".syntaxia_cache",
			MaxSearchResults: 30,
		},
		CallGraph: CallGraphConfig{
			MaxTransitiveDepth:     10,
			MaxMermaidNodes:        50,
			MaxCommitsForCoChanged: 50,
		},
		OpenAI: OpenAIConfig{
			MaxRetries:     3,
			BaseDelayMs:    1000,
			MaxDelayMs:     30000,
			TimeoutSeconds: 60,
		},
		AgenticChat: AgenticChatConfig{
			MaxIterations:    10,
			MaxToolLogLength: 2000,
		},
	}
}

// LoadConfig loads configuration from file, falling back to defaults
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // Use defaults if file doesn't exist
		}
		return nil, err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(config *Config, configPath string) error {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o600)
}

// LoadProjectConfig loads config from project's .syntaxia directory
func LoadProjectConfig(projectRoot string) (*Config, error) {
	configPath := filepath.Join(projectRoot, ".syntaxia", "config.json")
	return LoadConfig(configPath)
}
