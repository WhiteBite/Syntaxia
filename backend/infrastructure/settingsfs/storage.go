package settingsfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const keyringService = "syntaxia"

// storage handles the low-level reading/writing of settings to files and the system keyring.
type storage struct {
	settingsFilePath string
}

// newStorage initializes a new storage helper.
func newStorage() (*storage, error) {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %w", err)
	}
	appSpecificDir := filepath.Join(appDataDir, "syntaxia")
	if err := os.MkdirAll(appSpecificDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create application data directory: %w", err)
	}
	return &storage{
		settingsFilePath: filepath.Join(appSpecificDir, "settings.json"),
	}, nil
}

// loadFromFile reads and unmarshals the settings from the JSON file.
func (s *storage) loadFromFile(settings *appSettings) error {
	data, err := os.ReadFile(s.settingsFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File not existing is not an error on first run
		}
		return fmt.Errorf("failed to read settings file: %w", err)
	}
	return json.Unmarshal(data, settings)
}

// saveToFile marshals and writes the settings to the JSON file.
func (s *storage) saveToFile(settings *appSettings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}
	return os.WriteFile(s.settingsFilePath, data, 0o600)
}

// loadKeysFromKeyring populates the API key fields from the system's keyring.
func (s *storage) loadKeysFromKeyring(settings *secureSettings) error {
	var err error
	settings.openAIAPIKey, err = keyring.Get(keyringService, "openai")
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return fmt.Errorf("failed to get openai key: %w", err)
	}
	settings.geminiAPIKey, err = keyring.Get(keyringService, "gemini")
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return fmt.Errorf("failed to get gemini key: %w", err)
	}
	settings.openRouterAPIKey, err = keyring.Get(keyringService, "openrouter")
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return fmt.Errorf("failed to get openrouter key: %w", err)
	}
	settings.localAIAPIKey, err = keyring.Get(keyringService, "localai")
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return fmt.Errorf("failed to get localai key: %w", err)
	}
	settings.qwenAPIKey, err = keyring.Get(keyringService, "qwen")
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return fmt.Errorf("failed to get qwen key: %w", err)
	}
	return nil
}

// saveKeysToKeyring persists the API keys to the system's keyring.
func (s *storage) saveKeysToKeyring(settings *secureSettings) error {
	if err := keyring.Set(keyringService, "openai", settings.openAIAPIKey); err != nil {
		return fmt.Errorf("failed to set openai key: %w", err)
	}
	if err := keyring.Set(keyringService, "gemini", settings.geminiAPIKey); err != nil {
		return fmt.Errorf("failed to set gemini key: %w", err)
	}
	if err := keyring.Set(keyringService, "openrouter", settings.openRouterAPIKey); err != nil {
		return fmt.Errorf("failed to set openrouter key: %w", err)
	}
	if err := keyring.Set(keyringService, "localai", settings.localAIAPIKey); err != nil {
		return fmt.Errorf("failed to set localai key: %w", err)
	}
	if err := keyring.Set(keyringService, "qwen", settings.qwenAPIKey); err != nil {
		return fmt.Errorf("failed to set qwen key: %w", err)
	}
	return nil
}
