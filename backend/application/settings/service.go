package settings

import (
	"fmt"
	"syntaxia/domain"
	"sync"
)

// AIProviderCacheInvalidator interface for invalidating AI provider cache
type AIProviderCacheInvalidator interface {
	InvalidateProviderCache()
}

// Service отвечает за управление настройками приложения.
type Service struct {
	log                           domain.Logger
	bus                           domain.EventBus
	settingsRepo                  domain.SettingsRepository
	modelFetchers                 domain.ModelFetcherRegistry
	aiCacheInvalidator            AIProviderCacheInvalidator
	onIgnoreRulesChangedCallbacks []func() error
	muCallbacks                   sync.RWMutex
}

// NewService создает новый экземпляр Service.
func NewService(
	log domain.Logger,
	bus domain.EventBus,
	settingsRepo domain.SettingsRepository,
	modelFetchers domain.ModelFetcherRegistry,
) (*Service, error) {
	s := &Service{
		log:           log,
		bus:           bus,
		settingsRepo:  settingsRepo,
		modelFetchers: modelFetchers,
	}
	return s, nil
}

// GetSettingsDTO собирает все настройки в один DTO для передачи на фронтенд.
func (s *Service) GetSettingsDTO() (domain.SettingsDTO, error) {
	return s.settingsRepo.GetSettingsDTO()
}

// SaveSettingsDTO принимает DTO с фронтенда и обновляет настройки.
func (s *Service) SaveSettingsDTO(dto domain.SettingsDTO) error {
	// Track if AI-related settings changed
	oldDTO, _ := s.settingsRepo.GetSettingsDTO()
	aiSettingsChanged := oldDTO.SelectedProvider != dto.SelectedProvider ||
		oldDTO.OpenAIAPIKey != dto.OpenAIAPIKey ||
		oldDTO.GeminiAPIKey != dto.GeminiAPIKey ||
		oldDTO.OpenRouterAPIKey != dto.OpenRouterAPIKey ||
		oldDTO.LocalAIAPIKey != dto.LocalAIAPIKey ||
		oldDTO.LocalAIHost != dto.LocalAIHost ||
		oldDTO.QwenAPIKey != dto.QwenAPIKey

	s.settingsRepo.SetCustomIgnoreRules(dto.CustomIgnoreRules)
	s.settingsRepo.SetCustomPromptRules(dto.CustomPromptRules)
	s.settingsRepo.SetOpenAIKey(dto.OpenAIAPIKey)
	s.settingsRepo.SetGeminiKey(dto.GeminiAPIKey)
	s.settingsRepo.SetOpenRouterKey(dto.OpenRouterAPIKey)
	s.settingsRepo.SetLocalAIKey(dto.LocalAIAPIKey)
	s.settingsRepo.SetLocalAIHost(dto.LocalAIHost)
	s.settingsRepo.SetLocalAIModelName(dto.LocalAIModelName)
	s.settingsRepo.SetSelectedAIProvider(dto.SelectedProvider)
	s.settingsRepo.SetUseGitignore(dto.UseGitignore)
	s.settingsRepo.SetUseCustomIgnore(dto.UseCustomIgnore)

	for provider, model := range dto.SelectedModels {
		s.settingsRepo.SetSelectedModel(provider, model)
	}
	for provider, models := range dto.AvailableModels {
		s.settingsRepo.SetModels(provider, models)
	}

	if err := s.settingsRepo.Save(); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	// Invalidate AI provider cache if AI settings changed
	if aiSettingsChanged && s.aiCacheInvalidator != nil {
		s.aiCacheInvalidator.InvalidateProviderCache()
	}

	s.notifyIgnoreRulesChanged()
	return nil
}

// RefreshModels обновляет список моделей для указанного провайдера.
func (s *Service) RefreshModels(provider, apiKey string) error {
	fetcher, ok := s.modelFetchers[provider]
	if !ok {
		return fmt.Errorf("no model fetcher for provider: %s", provider)
	}

	s.log.Info("Refreshing model list for: " + provider)
	models, err := fetcher(apiKey)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to refresh models for %s: %v", provider, err))
		return err
	}

	s.settingsRepo.SetModels(provider, models)

	currentModel := s.settingsRepo.GetSelectedModel(provider)
	isCurrentModelValid := false
	for _, m := range models {
		if m == currentModel {
			isCurrentModelValid = true
			break
		}
	}
	if !isCurrentModelValid && len(models) > 0 {
		s.log.Info(fmt.Sprintf("Selected model '%s' not found in new list. Selecting '%s' by default.", currentModel, models[0]))
		s.settingsRepo.SetSelectedModel(provider, models[0])
	}

	if err := s.settingsRepo.Save(); err != nil {
		return fmt.Errorf("failed to save new model list: %w", err)
	}

	s.log.Info(fmt.Sprintf("Model list for %s refreshed successfully.", provider))
	return nil
}

// OnIgnoreRulesChanged регистрирует коллбэк.
func (s *Service) OnIgnoreRulesChanged(callback func() error) {
	s.muCallbacks.Lock()
	defer s.muCallbacks.Unlock()
	s.onIgnoreRulesChangedCallbacks = append(s.onIgnoreRulesChangedCallbacks, callback)
}

func (s *Service) notifyIgnoreRulesChanged() {
	s.muCallbacks.RLock()
	defer s.muCallbacks.RUnlock()
	for _, cb := range s.onIgnoreRulesChangedCallbacks {
		if err := cb(); err != nil {
			s.log.Error(fmt.Sprintf("Error executing OnIgnoreRulesChanged callback: %v", err))
		}
	}
}

// GetRecentProjects returns the list of recent projects
func (s *Service) GetRecentProjects() []domain.RecentProjectInfo {
	return s.settingsRepo.GetRecentProjects()
}

// AddRecentProject adds a project to the recent list
func (s *Service) AddRecentProject(path, name string) {
	s.settingsRepo.AddRecentProject(path, name)
}

// RemoveRecentProject removes a project from the recent list
func (s *Service) RemoveRecentProject(path string) {
	s.settingsRepo.RemoveRecentProject(path)
}

// Save persists settings to disk
func (s *Service) Save() error {
	return s.settingsRepo.Save()
}

// SetAICacheInvalidator sets the AI cache invalidator (called after AIService is created)
func (s *Service) SetAICacheInvalidator(invalidator AIProviderCacheInvalidator) {
	s.aiCacheInvalidator = invalidator
}

// GetCustomIgnoreRules returns custom ignore rules
func (s *Service) GetCustomIgnoreRules() string {
	return s.settingsRepo.GetCustomIgnoreRules()
}

// SetCustomIgnoreRules updates custom ignore rules
func (s *Service) SetCustomIgnoreRules(rules string) {
	s.settingsRepo.SetCustomIgnoreRules(rules)
	s.notifyIgnoreRulesChanged()
}

// GetUseGitignore returns whether to use .gitignore
func (s *Service) GetUseGitignore() bool {
	return s.settingsRepo.GetUseGitignore()
}

// SetUseGitignore updates whether to use .gitignore
func (s *Service) SetUseGitignore(use bool) {
	s.settingsRepo.SetUseGitignore(use)
	s.notifyIgnoreRulesChanged()
}

// GetUseCustomIgnore returns whether to use custom ignore rules
func (s *Service) GetUseCustomIgnore() bool {
	return s.settingsRepo.GetUseCustomIgnore()
}

// SetUseCustomIgnore updates whether to use custom ignore rules
func (s *Service) SetUseCustomIgnore(use bool) {
	s.settingsRepo.SetUseCustomIgnore(use)
	s.notifyIgnoreRulesChanged()
}
