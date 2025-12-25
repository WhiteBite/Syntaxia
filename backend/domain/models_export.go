package domain

import "time"

type ExportMode string

const (
	ExportModeClipboard ExportMode = "clipboard"
	ExportModeAI        ExportMode = "ai"
	ExportModeHuman     ExportMode = "human"
)

type ExportSettings struct {
	Mode    ExportMode `json:"mode"`
	Context string     `json:"context"`

	// Project export fields for app.go compatibility
	ProjectPath string                 `json:"projectPath"`
	Format      string                 `json:"format"`
	Options     map[string]interface{} `json:"options,omitempty"`

	// Clipboard
	StripComments   bool   `json:"stripComments"`
	IncludeManifest bool   `json:"includeManifest"`
	ExportFormat    string `json:"exportFormat"` // "plain" | "manifest" | "json"

	// AI
	AIProfile       string `json:"aiProfile"`
	TokenLimit      int    `json:"tokenLimit"`
	FileSizeLimitKB int    `json:"fileSizeLimitKB"` // Max size of an individual file read, not context chunk

	// New AI Splitting Options
	EnableAutoSplit   bool   `json:"enableAutoSplit"`
	MaxTokensPerChunk int    `json:"maxTokensPerChunk"` // Max tokens for each generated chunk
	OverlapTokens     int    `json:"overlapTokens"`
	SplitStrategy     string `json:"splitStrategy"` // "token" | "file" | "smart"

	// Human
	Theme              string `json:"theme"`
	IncludeLineNumbers bool   `json:"includeLineNumbers"`
	IncludePageNumbers bool   `json:"includePageNumbers"`
}

type ExportResult struct {
	Mode       ExportMode `json:"mode"`
	Text       string     `json:"text,omitempty"`
	FileName   string     `json:"fileName,omitempty"`
	DataBase64 string     `json:"dataBase64,omitempty"`
	FilePath   string     `json:"filePath,omitempty"`  // NEW: для больших файлов
	IsLarge    bool       `json:"isLarge,omitempty"`   // NEW: флаг больших файлов
	SizeBytes  int64      `json:"sizeBytes,omitempty"` // NEW: размер файла
}

// SplitSettings для ContextSplitter
type SplitSettings struct {
	MaxTokensPerChunk int
	OverlapTokens     int
	SplitStrategy     string
}

// ExportHistoryItem represents a single export operation in history
type ExportHistoryItem struct {
	ID          string     `json:"id"`
	ProjectPath string     `json:"projectPath"`
	Mode        ExportMode `json:"mode"`
	Format      string     `json:"format"`
	FileName    string     `json:"fileName"`
	SizeBytes   int64      `json:"sizeBytes"`
	CreatedAt   time.Time  `json:"createdAt"`
	FilePath    string     `json:"filePath,omitempty"`
	Status      string     `json:"status"` // "success", "failed", "in_progress"
	Error       string     `json:"error,omitempty"`
}
