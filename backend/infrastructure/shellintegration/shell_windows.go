//go:build windows

package shellintegration

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	registryKeyDirectory  = `HKEY_CURRENT_USER\Software\Classes\Directory\shell\Syntaxia`
	registryKeyBackground = `HKEY_CURRENT_USER\Software\Classes\Directory\Background\shell\Syntaxia`
)

func (s *Service) isRegisteredOS() (bool, error) {
	cmd := exec.Command("reg", "query", registryKeyDirectory)
	err := cmd.Run()
	return err == nil, nil
}

// getRegisteredPath returns the exe path currently in registry
func (s *Service) getRegisteredPath() (string, error) {
	cmd := exec.Command("reg", "query", registryKeyDirectory+`\command`, "/ve")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// Parse output to extract path
	// Format: "    (Default)    REG_SZ    "C:\path\to\app.exe" "%V""
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "REG_SZ") {
			parts := strings.SplitN(line, "REG_SZ", 2)
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				// Extract path from "\"path\" \"%V\""
				if strings.HasPrefix(value, `"`) {
					endQuote := strings.Index(value[1:], `"`)
					if endQuote > 0 {
						return value[1 : endQuote+1], nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("could not parse registry value")
}

// NeedsUpdate checks if registered path differs from current exe path
func (s *Service) NeedsUpdate(currentExePath string) bool {
	registered, err := s.getRegisteredPath()
	if err != nil {
		return false
	}
	return !strings.EqualFold(registered, currentExePath)
}

func (s *Service) registerOS(exePath string) error {
	menuText := "Open in Syntaxia"
	iconPath := exePath

	commands := [][]string{
		{"reg", "add", registryKeyDirectory, "/ve", "/d", menuText, "/f"},
		{"reg", "add", registryKeyDirectory, "/v", "Icon", "/d", iconPath, "/f"},
		{"reg", "add", registryKeyDirectory + `\command`, "/ve", "/d", fmt.Sprintf(`"%s" "%%V"`, exePath), "/f"},
		{"reg", "add", registryKeyBackground, "/ve", "/d", menuText, "/f"},
		{"reg", "add", registryKeyBackground, "/v", "Icon", "/d", iconPath, "/f"},
		{"reg", "add", registryKeyBackground + `\command`, "/ve", "/d", fmt.Sprintf(`"%s" "%%V"`, exePath), "/f"},
	}

	for _, args := range commands {
		// #nosec G204 -- args are constructed from trusted constants, not user input
		cmd := exec.Command(args[0], args[1:]...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to execute %s: %w, output: %s", strings.Join(args, " "), err, string(output))
		}
	}

	return nil
}

func (s *Service) unregisterOS() error {
	commands := [][]string{
		{"reg", "delete", registryKeyDirectory, "/f"},
		{"reg", "delete", registryKeyBackground, "/f"},
	}

	var lastErr error
	for _, args := range commands {
		// #nosec G204 -- args are constructed from trusted constants, not user input
		cmd := exec.Command(args[0], args[1:]...)
		if err := cmd.Run(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}
