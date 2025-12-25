//go:build !windows

package executil

import "os/exec"

// HideWindow is a no-op on non-Windows platforms
func HideWindow(cmd *exec.Cmd) {
	// No-op on Linux/macOS
}
