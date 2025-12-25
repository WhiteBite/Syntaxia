// Package domain provides core business types and utilities.
package domain

// TruncateString truncates a string to maxLen characters (runes, not bytes).
// If the string is longer than maxLen, it appends "..." to indicate truncation.
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
