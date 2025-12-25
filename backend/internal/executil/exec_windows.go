//go:build windows

package executil

import (
	"os/exec"
	"syscall"
)

// HideWindow sets SysProcAttr to hide console window on Windows
func HideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
}
