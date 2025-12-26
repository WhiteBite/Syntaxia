package staticanalyzer

import (
	"fmt"
	"os"
)

// parseInt parses string to int
func parseInt(s string) int {
	var val int
	_, _ = fmt.Sscanf(s, "%d", &val)
	return val
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
