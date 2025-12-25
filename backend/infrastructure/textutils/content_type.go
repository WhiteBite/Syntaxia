// Package textutils provides utilities for detecting file content types.
package textutils

import (
	"bytes"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// ContentType represents the type of file content.
type ContentType int

const (
	ContentTypeText    ContentType = iota // Readable text file
	ContentTypeBinary                     // Binary file (images, executables, etc.)
	ContentTypeUnknown                    // Unknown, needs content inspection
)

// String returns string representation of ContentType.
func (ct ContentType) String() string {
	switch ct {
	case ContentTypeText:
		return "text"
	case ContentTypeBinary:
		return "binary"
	default:
		return "unknown"
	}
}

// knownTextExtensions are extensions that are always considered text.
var knownTextExtensions = map[string]bool{
	// Programming languages
	".go": true, ".js": true, ".ts": true, ".tsx": true, ".jsx": true,
	".py": true, ".rb": true, ".php": true, ".java": true, ".kt": true,
	".c": true, ".cpp": true, ".h": true, ".hpp": true, ".cs": true,
	".rs": true, ".swift": true, ".scala": true, ".clj": true,
	".lua": true, ".pl": true, ".pm": true, ".r": true, ".m": true,
	".dart": true, ".elm": true, ".ex": true, ".exs": true,
	".hs": true, ".ml": true, ".fs": true, ".fsx": true,
	".groovy": true, ".gradle": true, ".v": true, ".zig": true,

	// Web
	".html": true, ".htm": true, ".css": true, ".scss": true, ".sass": true,
	".less": true, ".vue": true, ".svelte": true, ".astro": true,

	// Data/Config
	".json": true, ".yaml": true, ".yml": true, ".toml": true,
	".xml": true, ".csv": true, ".tsv": true, ".ini": true,
	".conf": true, ".cfg": true, ".properties": true,
	".env": true, ".envrc": true, ".editorconfig": true,

	// Documentation
	".md": true, ".markdown": true, ".rst": true, ".txt": true,
	".adoc": true, ".asciidoc": true, ".tex": true, ".org": true,

	// Shell/Scripts
	".sh": true, ".bash": true, ".zsh": true, ".fish": true,
	".ps1": true, ".psm1": true, ".bat": true, ".cmd": true,

	// Build/Package
	".makefile": true, ".cmake": true, ".dockerfile": true,
	".mod": true, ".sum": true, ".lock": true,

	// Logs and misc
	".log": true, ".diff": true, ".patch": true, ".sql": true,
	".graphql": true, ".gql": true, ".proto": true,
}

// knownBinaryExtensions are extensions that are always considered binary.
var knownBinaryExtensions = map[string]bool{
	// Images
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".bmp": true, ".ico": true, ".icns": true, ".webp": true,
	".tiff": true, ".tif": true, ".psd": true, ".ai": true,
	".svg": true, // SVG is XML but often large and not useful for context

	// Fonts
	".ttf": true, ".otf": true, ".woff": true, ".woff2": true, ".eot": true,

	// Audio/Video
	".mp3": true, ".wav": true, ".ogg": true, ".flac": true, ".aac": true,
	".mp4": true, ".avi": true, ".mov": true, ".mkv": true, ".webm": true,

	// Archives
	".zip": true, ".rar": true, ".7z": true, ".tar": true,
	".gz": true, ".bz2": true, ".xz": true, ".tgz": true,

	// Executables/Libraries
	".exe": true, ".dll": true, ".so": true, ".dylib": true,
	".bin": true, ".app": true, ".msi": true, ".deb": true, ".rpm": true,

	// Compiled
	".class": true, ".jar": true, ".war": true, ".ear": true,
	".pyc": true, ".pyo": true, ".o": true, ".a": true, ".obj": true,

	// Documents (binary)
	".pdf": true, ".doc": true, ".docx": true,
	".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true,

	// Databases
	".db": true, ".sqlite": true, ".sqlite3": true, ".mdb": true,

	// Other
	".dat": true, ".pak": true, ".cache": true,
}

// magicSignatures maps file signatures to content type.
var magicSignatures = []struct {
	signature []byte
	isBinary  bool
}{
	// Binary signatures
	{[]byte{0x89, 0x50, 0x4E, 0x47}, true},             // PNG
	{[]byte{0xFF, 0xD8, 0xFF}, true},                   // JPEG
	{[]byte{0x47, 0x49, 0x46, 0x38}, true},             // GIF
	{[]byte{0x50, 0x4B, 0x03, 0x04}, true},             // ZIP/DOCX/XLSX
	{[]byte{0x50, 0x4B, 0x05, 0x06}, true},             // ZIP empty
	{[]byte{0x52, 0x61, 0x72, 0x21}, true},             // RAR
	{[]byte{0x1F, 0x8B}, true},                         // GZIP
	{[]byte{0x42, 0x5A, 0x68}, true},                   // BZIP2
	{[]byte{0x4D, 0x5A}, true},                         // EXE/DLL
	{[]byte{0x7F, 0x45, 0x4C, 0x46}, true},             // ELF
	{[]byte{0xCA, 0xFE, 0xBA, 0xBE}, true},             // Java class / Mach-O
	{[]byte{0x25, 0x50, 0x44, 0x46}, true},             // PDF
	{[]byte{0x00, 0x00, 0x01, 0x00}, true},             // ICO
	{[]byte{0x00, 0x01, 0x00, 0x00}, true},             // TTF
	{[]byte{0x4F, 0x54, 0x54, 0x4F}, true},             // OTF
	{[]byte{0x77, 0x4F, 0x46, 0x46}, true},             // WOFF
	{[]byte{0x77, 0x4F, 0x46, 0x32}, true},             // WOFF2
	{[]byte{0x49, 0x44, 0x33}, true},                   // MP3 with ID3
	{[]byte{0xFF, 0xFB}, true},                         // MP3
	{[]byte{0x52, 0x49, 0x46, 0x46}, true},             // WAV/AVI
	{[]byte{0x00, 0x00, 0x00, 0x1C, 0x66, 0x74}, true}, // MP4
	{[]byte{0x00, 0x00, 0x00, 0x20, 0x66, 0x74}, true}, // MP4
	{[]byte{0x1A, 0x45, 0xDF, 0xA3}, true},             // MKV/WebM
	{[]byte{0x53, 0x51, 0x4C, 0x69, 0x74, 0x65}, true}, // SQLite
}

// DetectByExtension returns content type based on file extension.
// Returns ContentTypeUnknown if extension is not recognized.
func DetectByExtension(filename string) ContentType {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		// Check for known extensionless files
		base := strings.ToLower(filepath.Base(filename))
		switch base {
		case "makefile", "dockerfile", "jenkinsfile", "vagrantfile",
			"gemfile", "rakefile", "procfile", "brewfile",
			".gitignore", ".gitattributes", ".dockerignore",
			".editorconfig", ".prettierrc", ".eslintrc",
			"license", "readme", "changelog", "authors", "contributors":
			return ContentTypeText
		}
		return ContentTypeUnknown
	}

	if knownTextExtensions[ext] {
		return ContentTypeText
	}
	if knownBinaryExtensions[ext] {
		return ContentTypeBinary
	}
	return ContentTypeUnknown
}

// DetectByContent analyzes file content to determine if it's text or binary.
// Checks magic bytes first, then falls back to null byte and UTF-8 detection.
func DetectByContent(content []byte) ContentType {
	if len(content) == 0 {
		return ContentTypeText // Empty files are text
	}

	// Check magic signatures
	for _, sig := range magicSignatures {
		if len(content) >= len(sig.signature) && bytes.HasPrefix(content, sig.signature) {
			if sig.isBinary {
				return ContentTypeBinary
			}
			return ContentTypeText
		}
	}

	// Sample first 8KB for analysis
	sample := content
	if len(sample) > 8192 {
		sample = sample[:8192]
	}

	// Check for null bytes (strong indicator of binary)
	if bytes.Contains(sample, []byte{0}) {
		return ContentTypeBinary
	}

	// Check if valid UTF-8
	if !utf8.Valid(sample) {
		return ContentTypeBinary
	}

	// Check for high ratio of non-printable characters
	nonPrintable := 0
	for _, b := range sample {
		if b < 32 && b != '\t' && b != '\n' && b != '\r' {
			nonPrintable++
		}
	}
	if float64(nonPrintable)/float64(len(sample)) > 0.1 {
		return ContentTypeBinary
	}

	return ContentTypeText
}

// Detect combines extension and content detection.
// Uses extension first for speed, falls back to content analysis if unknown.
func Detect(filename string, content []byte) ContentType {
	extType := DetectByExtension(filename)
	if extType != ContentTypeUnknown {
		return extType
	}
	return DetectByContent(content)
}

// IsText returns true if the file is likely a text file.
func IsText(filename string, content []byte) bool {
	return Detect(filename, content) == ContentTypeText
}

// IsBinary returns true if the file is likely a binary file.
func IsBinary(filename string, content []byte) bool {
	return Detect(filename, content) == ContentTypeBinary
}
