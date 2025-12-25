package analyzers

// findBlockEndLine finds the end line of a block by matching braces
// Used by Dart, Java, Kotlin and other C-style languages
func findBlockEndLine(lines []string, startLineIdx int) int {
	if startLineIdx < 0 || startLineIdx >= len(lines) {
		return startLineIdx + 1
	}

	braceCount := 0
	started := false

	for i := startLineIdx; i < len(lines); i++ {
		line := lines[i]
		for _, ch := range line {
			if ch == '{' {
				braceCount++
				started = true
			} else if ch == '}' {
				braceCount--
			}
		}
		if started && braceCount == 0 {
			return i + 1
		}
	}

	endLine := startLineIdx + 20
	if endLine > len(lines) {
		endLine = len(lines)
	}
	return endLine
}
