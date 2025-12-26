// Package textutils provides text processing utilities.
package textutils

// TruncateString truncates a string to maxLen characters (runes, not bytes).
// If the string is longer than maxLen, it appends "..." to indicate truncation.
// This is the single source of truth for string truncation in the application.
func TruncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

// TruncateStringNoEllipsis truncates a string to maxLen without adding ellipsis.
func TruncateStringNoEllipsis(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// TruncateLines truncates text to maxLines lines.
// If the text has more lines, it appends a message indicating truncation.
func TruncateLines(s string, maxLines int) string {
	if maxLines <= 0 {
		return ""
	}

	lines := 0
	lastNewline := -1

	for i, c := range s {
		if c == '\n' {
			lines++
			if lines >= maxLines {
				return s[:i] + "\n... (truncated)"
			}
			lastNewline = i
		}
	}

	// No truncation needed
	_ = lastNewline // silence unused warning
	return s
}
