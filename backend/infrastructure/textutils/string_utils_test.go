package textutils

import "testing"

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"empty string", "", 10, ""},
		{"short string", "hello", 10, "hello"},
		{"exact length", "hello", 5, "hello"},
		{"truncate with ellipsis", "hello world", 8, "hello..."},
		{"very short max", "hello", 3, "hel"},
		{"zero max", "hello", 0, ""},
		{"negative max", "hello", -1, ""},
		{"unicode string", "привет мир", 8, "приве..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateString(%q, %d) = %q, want %q",
					tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestTruncateStringNoEllipsis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"empty string", "", 10, ""},
		{"short string", "hello", 10, "hello"},
		{"exact length", "hello", 5, "hello"},
		{"truncate no ellipsis", "hello world", 8, "hello wo"},
		{"zero max", "hello", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateStringNoEllipsis(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateStringNoEllipsis(%q, %d) = %q, want %q",
					tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestTruncateLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLines int
		contains string
	}{
		{"empty string", "", 5, ""},
		{"single line", "hello", 5, "hello"},
		{"two lines under limit", "line1\nline2", 5, "line1\nline2"},
		{"truncate at limit", "line1\nline2\nline3", 2, "line1\nline2"},
		{"zero max", "hello", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateLines(tt.input, tt.maxLines)
			if tt.maxLines > 0 && len(result) > 0 && result[:len(tt.contains)] != tt.contains[:min(len(tt.contains), len(result))] {
				t.Errorf("TruncateLines result doesn't start with expected content")
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
