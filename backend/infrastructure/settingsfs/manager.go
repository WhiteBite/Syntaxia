package settingsfs

import (
	"fmt"
	"syntaxia/domain"
	"sync"
	"time"
)

// appSettings stores settings that are safe to write to a JSON file.
// API keys are handled separately via the system's keyring.
type appSettings struct {
	CustomIgnoreRules string                     `json:"customIgnoreRules"`
	CustomPromptRules string                     `json:"customPromptRules"`
	UseGitignore      bool                       `json:"useGitignore"`
	UseCustomIgnore   bool                       `json:"useCustomIgnore"`
	LocalAIHost       string                     `json:"localAIHost,omitempty"`
	LocalAIModelName  string                     `json:"localAIModelName,omitempty"`
	QwenHost          string                     `json:"qwenHost,omitempty"`
	SelectedProvider  string                     `json:"selectedProvider"`
	SelectedModels    map[string]string          `json:"selectedModels"`
	AvailableModels   map[string][]string        `json:"availableModels"`
	RecentProjects    []domain.RecentProjectInfo `json:"recentProjects,omitempty"`
}

// secureSettings holds secrets that are stored in the system's keyring.
type secureSettings struct {
	openAIAPIKey     string
	geminiAPIKey     string
	openRouterAPIKey string
	localAIAPIKey    string
	qwenAPIKey       string
}

// Manager orchestrates settings persistence, separating file and keyring storage.
type Manager struct {
	log                domain.Logger
	mu                 sync.RWMutex
	storage            *storage
	settings           appSettings
	secure             secureSettings
	defaultIgnoreRules string
	defaultPromptRules string
}

// New creates a new Manager instance and loads settings.
func New(logger domain.Logger, defaultIgnore, defaultPrompt string) (domain.SettingsRepository, error) {
	s, err := newStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to create settings storage: %w", err)
	}
	m := &Manager{
		log:                logger,
		storage:            s,
		defaultIgnoreRules: defaultIgnore,
		defaultPromptRules: defaultPrompt,
	}
	if err := m.load(); err != nil {
		return nil, fmt.Errorf("failed to load settings: %w", err)
	}
	return m, nil
}

func (m *Manager) load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.loadDefaults()

	if err := m.storage.loadFromFile(&m.settings); err != nil {
		m.log.Error(fmt.Sprintf("Error reading settings file, using defaults: %v", err))
	} else {
		m.mergeWithDefaults()
	}

	if err := m.storage.loadKeysFromKeyring(&m.secure); err != nil {
		// Log the error but don't fail, as keys might not be critical on startup
		m.log.Warning(fmt.Sprintf("Could not load API keys from keyring: %v", err))
	}
	return nil
}

func (m *Manager) loadDefaults() {
	m.settings = appSettings{
		UseGitignore:     true,
		UseCustomIgnore:  true,
		SelectedProvider: "openai",
		LocalAIHost:      "http://localhost:1234/v1",
		QwenHost:         "https://dashscope.aliyuncs.com/compatible-mode/v1",
		SelectedModels: map[string]string{
			"openai":     "gpt-4o",
			"gemini":     "gemini-1.5-pro-latest",
			"openrouter": "google/gemini-flash-1.5",
			"localai":    "local-model",
			"qwen":       "qwen-coder-plus-latest",
		},
		AvailableModels: map[string][]string{
			"openai":     {"gpt-4o", "gpt-4-turbo", "gpt-3.5-turbo"},
			"gemini":     {"gemini-1.5-pro-latest", "gemini-1.5-flash-latest"},
			"openrouter": {"google/gemini-flash-1.5", "openai/gpt-4o", "meta-llama/llama-3-70b-instruct"},
			"localai":    {"local-model"},
			"qwen":       {"qwen-coder-plus-latest", "qwen-coder-plus", "qwen-plus-latest", "qwen-turbo-latest", "qwen-max"},
		},
	}
}

func (m *Manager) mergeWithDefaults() {
	if m.settings.CustomIgnoreRules == "" {
		m.settings.CustomIgnoreRules = m.defaultIgnoreRules
	}
	if m.settings.CustomPromptRules == "" {
		m.settings.CustomPromptRules = m.defaultPromptRules
	}
}

// Save persists all settings to their respective storage locations.
func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.storage.saveKeysToKeyring(&m.secure); err != nil {
		// Log the error but don't fail, as keyring is not always essential
		m.log.Warning(fmt.Sprintf("Could not save API keys to keyring: %v", err))
	}
	return m.storage.saveToFile(&m.settings)
}

// Getters and Setters for secure settings (API Keys)
func (m *Manager) GetOpenAIKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.secure.openAIAPIKey
}
func (m *Manager) GetGeminiKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.secure.geminiAPIKey
}
func (m *Manager) GetOpenRouterKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.secure.openRouterAPIKey
}
func (m *Manager) GetLocalAIKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.secure.localAIAPIKey
}
func (m *Manager) SetOpenAIKey(k string) { m.mu.Lock(); m.secure.openAIAPIKey = k; m.mu.Unlock() }
func (m *Manager) SetGeminiKey(k string) { m.mu.Lock(); m.secure.geminiAPIKey = k; m.mu.Unlock() }
func (m *Manager) SetOpenRouterKey(k string) {
	m.mu.Lock()
	m.secure.openRouterAPIKey = k
	m.mu.Unlock()
}
func (m *Manager) SetLocalAIKey(k string) { m.mu.Lock(); m.secure.localAIAPIKey = k; m.mu.Unlock() }
func (m *Manager) GetQwenKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.secure.qwenAPIKey
}
func (m *Manager) SetQwenKey(k string) { m.mu.Lock(); m.secure.qwenAPIKey = k; m.mu.Unlock() }

// Getters and Setters for file-based settings
func (m *Manager) GetCustomIgnoreRules() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings.CustomIgnoreRules
}
func (m *Manager) GetCustomPromptRules() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings.CustomPromptRules
}
func (m *Manager) GetUseGitignore() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings.UseGitignore
}
func (m *Manager) GetUseCustomIgnore() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings.UseCustomIgnore
}
func (m *Manager) GetLocalAIHost() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings.LocalAIHost
}
func (m *Manager) GetLocalAIModelName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings.LocalAIModelName
}
func (m *Manager) GetQwenHost() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.settings.QwenHost == "" {
		return "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}
	return m.settings.QwenHost
}
func (m *Manager) SetQwenHost(h string) { m.mu.Lock(); m.settings.QwenHost = h; m.mu.Unlock() }
func (m *Manager) GetSelectedAIProvider() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings.SelectedProvider
}
func (m *Manager) GetSelectedModel(p string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings.SelectedModels[p]
}
func (m *Manager) GetModels(p string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings.AvailableModels[p]
}
func (m *Manager) SetCustomIgnoreRules(r string) {
	m.mu.Lock()
	m.settings.CustomIgnoreRules = r
	m.mu.Unlock()
}
func (m *Manager) SetCustomPromptRules(r string) {
	m.mu.Lock()
	m.settings.CustomPromptRules = r
	m.mu.Unlock()
}
func (m *Manager) SetUseGitignore(e bool) { m.mu.Lock(); m.settings.UseGitignore = e; m.mu.Unlock() }
func (m *Manager) SetUseCustomIgnore(e bool) {
	m.mu.Lock()
	m.settings.UseCustomIgnore = e
	m.mu.Unlock()
}
func (m *Manager) SetLocalAIHost(h string) { m.mu.Lock(); m.settings.LocalAIHost = h; m.mu.Unlock() }
func (m *Manager) SetLocalAIModelName(n string) {
	m.mu.Lock()
	m.settings.LocalAIModelName = n
	m.mu.Unlock()
}
func (m *Manager) SetSelectedAIProvider(p string) {
	m.mu.Lock()
	m.settings.SelectedProvider = p
	m.mu.Unlock()
}

func (m *Manager) SetSelectedModel(p, mdl string) {
	m.mu.Lock()
	if m.settings.SelectedModels == nil {
		m.settings.SelectedModels = make(map[string]string)
	}
	m.settings.SelectedModels[p] = mdl
	m.mu.Unlock()
}

func (m *Manager) SetModels(p string, mdls []string) {
	m.mu.Lock()
	if m.settings.AvailableModels == nil {
		m.settings.AvailableModels = make(map[string][]string)
	}
	m.settings.AvailableModels[p] = mdls
	m.mu.Unlock()
}

func (m *Manager) GetSettingsDTO() (domain.SettingsDTO, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create copies of maps to avoid race conditions on the frontend if the DTO is held onto
	selectedModelsCopy := make(map[string]string)
	for k, v := range m.settings.SelectedModels {
		selectedModelsCopy[k] = v
	}
	availableModelsCopy := make(map[string][]string)
	for k, v := range m.settings.AvailableModels {
		availableModelsCopy[k] = append([]string(nil), v...)
	}

	qwenHost := m.settings.QwenHost
	if qwenHost == "" {
		qwenHost = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}

	return domain.SettingsDTO{
		CustomIgnoreRules: m.settings.CustomIgnoreRules,
		CustomPromptRules: m.settings.CustomPromptRules,
		OpenAIAPIKey:      m.secure.openAIAPIKey,
		GeminiAPIKey:      m.secure.geminiAPIKey,
		OpenRouterAPIKey:  m.secure.openRouterAPIKey,
		LocalAIAPIKey:     m.secure.localAIAPIKey,
		LocalAIHost:       m.settings.LocalAIHost,
		LocalAIModelName:  m.settings.LocalAIModelName,
		QwenAPIKey:        m.secure.qwenAPIKey,
		QwenHost:          qwenHost,
		SelectedProvider:  m.settings.SelectedProvider,
		SelectedModels:    selectedModelsCopy,
		AvailableModels:   availableModelsCopy,
		UseGitignore:      m.settings.UseGitignore,
		UseCustomIgnore:   m.settings.UseCustomIgnore,
		RecentProjects:    m.settings.RecentProjects,
	}, nil
}

// GetRecentProjects returns the list of recent projects
func (m *Manager) GetRecentProjects() []domain.RecentProjectInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]domain.RecentProjectInfo, len(m.settings.RecentProjects))
	copy(result, m.settings.RecentProjects)
	return result
}

// AddRecentProject adds or updates a project in the recent list
func (m *Manager) AddRecentProject(path, name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Current timestamp
	now := time.Now().Format(time.RFC3339)

	// Remove if exists
	filtered := []domain.RecentProjectInfo{}
	for _, p := range m.settings.RecentProjects {
		if p.Path != path {
			filtered = append(filtered, p)
		}
	}

	// Add to beginning
	newEntry := domain.RecentProjectInfo{
		Path:         path,
		Name:         name,
		LastOpenedAt: now,
	}
	m.settings.RecentProjects = append([]domain.RecentProjectInfo{newEntry}, filtered...)

	// Keep only last 10
	if len(m.settings.RecentProjects) > 10 {
		m.settings.RecentProjects = m.settings.RecentProjects[:10]
	}
}

// RemoveRecentProject removes a project from the recent list
func (m *Manager) RemoveRecentProject(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	filtered := []domain.RecentProjectInfo{}
	for _, p := range m.settings.RecentProjects {
		if p.Path != path {
			filtered = append(filtered, p)
		}
	}
	m.settings.RecentProjects = filtered
}
