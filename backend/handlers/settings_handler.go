package handlers

import (
	"encoding/json"
	"fmt"
	"syntaxia/application/settings"
	"syntaxia/domain"
	"sync"
)

// SettingsHandler handles all settings-related operations
type SettingsHandler struct {
	log             domain.Logger
	settingsService *settings.Service

	// Mutex for atomic settings updates
	mu sync.Mutex
}

// NewSettingsHandler creates a new settings handler
func NewSettingsHandler(
	log domain.Logger,
	settingsService *settings.Service,
) *SettingsHandler {
	return &SettingsHandler{
		log:             log,
		settingsService: settingsService,
	}
}

// GetSettings returns current settings
func (h *SettingsHandler) GetSettings() (domain.SettingsDTO, error) {
	return h.settingsService.GetSettingsDTO()
}

// SaveSettings saves settings from JSON
func (h *SettingsHandler) SaveSettings(settingsJSON string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var dto domain.SettingsDTO
	if err := json.Unmarshal([]byte(settingsJSON), &dto); err != nil {
		return fmt.Errorf("failed to parse settings JSON: %w", err)
	}
	return h.settingsService.SaveSettingsDTO(dto)
}

// RefreshAIModels refreshes AI models for a provider
func (h *SettingsHandler) RefreshAIModels(provider, apiKey string) error {
	return h.settingsService.RefreshModels(provider, apiKey)
}

// GetRecentProjects returns recent projects
func (h *SettingsHandler) GetRecentProjects() (string, error) {
	projects := h.settingsService.GetRecentProjects()
	result, err := json.Marshal(projects)
	if err != nil {
		return "", fmt.Errorf("failed to marshal recent projects: %w", err)
	}
	return string(result), nil
}

// AddRecentProject adds a project to recent list
func (h *SettingsHandler) AddRecentProject(path, name string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.settingsService.AddRecentProject(path, name)
	return h.settingsService.Save()
}

// RemoveRecentProject removes a project from recent list
func (h *SettingsHandler) RemoveRecentProject(path string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.settingsService.RemoveRecentProject(path)
	return h.settingsService.Save()
}

// GetCustomIgnoreRules returns custom ignore rules
func (h *SettingsHandler) GetCustomIgnoreRules() (string, error) {
	dto, err := h.settingsService.GetSettingsDTO()
	if err != nil {
		return "", fmt.Errorf("failed to get settings: %w", err)
	}
	return dto.CustomIgnoreRules, nil
}

// UpdateCustomIgnoreRules updates custom ignore rules
func (h *SettingsHandler) UpdateCustomIgnoreRules(rules string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	dto, err := h.settingsService.GetSettingsDTO()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}
	dto.CustomIgnoreRules = rules
	return h.settingsService.SaveSettingsDTO(dto)
}

// GetSLAPolicy returns current SLA policy
func (h *SettingsHandler) GetSLAPolicy() (string, error) {
	defaultPolicy := domain.SLAPolicy{
		Name:        "standard",
		Description: "Standard SLA policy",
		MaxTokens:   10000,
		MaxFiles:    50,
		MaxTime:     300,
		MaxMemory:   1024 * 1024 * 100,
		MaxRetries:  3,
		Timeout:     30,
	}

	policyJSON, err := json.Marshal(defaultPolicy)
	if err != nil {
		return "", fmt.Errorf("failed to marshal SLA policy: %w", err)
	}

	return string(policyJSON), nil
}

// SetSLAPolicy sets SLA policy
func (h *SettingsHandler) SetSLAPolicy(policyJSON string) error {
	var policy domain.SLAPolicy
	if err := json.Unmarshal([]byte(policyJSON), &policy); err != nil {
		return fmt.Errorf("failed to parse SLA policy JSON: %w", err)
	}

	h.log.Info(fmt.Sprintf("SLA Policy set: %s", policy.Name))
	return nil
}
