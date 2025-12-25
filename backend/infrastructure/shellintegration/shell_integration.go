package shellintegration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Service handles OS shell integration (context menu)
type Service struct {
	appName string
}

// NewService creates a new shell integration service
func NewService(appName string) *Service {
	return &Service{appName: appName}
}

// IsRegistered checks if context menu is registered
func (s *Service) IsRegistered() (bool, error) {
	return s.isRegisteredOS()
}

// Register adds "Open folder in Syntaxia" to context menu
func (s *Service) Register() error {
	exePath, err := s.getExePath()
	if err != nil {
		return err
	}
	return s.registerOS(exePath)
}

// Unregister removes context menu entry
func (s *Service) Unregister() error {
	return s.unregisterOS()
}

// GetCurrentOS returns current operating system name
func (s *Service) GetCurrentOS() string {
	return runtime.GOOS
}

// UpdateIfNeeded checks if exe path changed and updates registry
func (s *Service) UpdateIfNeeded() (bool, error) {
	registered, err := s.IsRegistered()
	if err != nil || !registered {
		return false, err
	}

	exePath, err := s.getExePath()
	if err != nil {
		return false, err
	}

	if s.NeedsUpdate(exePath) {
		if err := s.registerOS(exePath); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (s *Service) getExePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	return exePath, nil
}
