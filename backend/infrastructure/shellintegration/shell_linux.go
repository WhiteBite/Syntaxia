//go:build linux

package shellintegration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const desktopFileName = "syntaxia-folder.desktop"

func (s *Service) getDesktopFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "nautilus", "scripts", "Open in Syntaxia")
}

func (s *Service) getDesktopEntryPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "applications", desktopFileName)
}

func (s *Service) isRegisteredOS() (bool, error) {
	_, err := os.Stat(s.getDesktopFilePath())
	if err == nil {
		return true, nil
	}
	_, err = os.Stat(s.getDesktopEntryPath())
	return err == nil, nil
}

func (s *Service) registerOS(exePath string) error {
	nautilusDir := filepath.Dir(s.getDesktopFilePath())
	if err := os.MkdirAll(nautilusDir, 0755); err != nil {
		return fmt.Errorf("failed to create nautilus scripts directory: %w", err)
	}

	nautilusScript := fmt.Sprintf(`#!/bin/bash
"%s" "$NAUTILUS_SCRIPT_CURRENT_URI"
`, exePath)

	if err := os.WriteFile(s.getDesktopFilePath(), []byte(nautilusScript), 0755); err != nil {
		return fmt.Errorf("failed to write nautilus script: %w", err)
	}

	appsDir := filepath.Dir(s.getDesktopEntryPath())
	if err := os.MkdirAll(appsDir, 0755); err != nil {
		return fmt.Errorf("failed to create applications directory: %w", err)
	}

	desktopEntry := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=Open in Syntaxia
Exec="%s" %%f
Icon=%s
Terminal=false
Categories=Development;
MimeType=inode/directory;
`, exePath, exePath)

	if err := os.WriteFile(s.getDesktopEntryPath(), []byte(desktopEntry), 0644); err != nil {
		return fmt.Errorf("failed to write desktop entry: %w", err)
	}

	return nil
}

func (s *Service) unregisterOS() error {
	var lastErr error

	if err := os.Remove(s.getDesktopFilePath()); err != nil && !os.IsNotExist(err) {
		lastErr = err
	}

	if err := os.Remove(s.getDesktopEntryPath()); err != nil && !os.IsNotExist(err) {
		lastErr = err
	}

	return lastErr
}

// NeedsUpdate checks if registered path differs from current exe path
func (s *Service) NeedsUpdate(currentExePath string) bool {
	content, err := os.ReadFile(s.getDesktopFilePath())
	if err != nil {
		return false
	}
	return !strings.Contains(string(content), currentExePath)
}
