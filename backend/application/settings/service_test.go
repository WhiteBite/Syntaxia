package settings

import (
	"fmt"
	"syntaxia/domain"
	"strings"
	"sync"
	"testing"
)

// Mock logger - implements domain.Logger
type mockLogger struct{}

func (m *mockLogger) Debug(message string)   {}
func (m *mockLogger) Info(message string)    {}
func (m *mockLogger) Warning(message string) {}
func (m *mockLogger) Error(message string)   {}
func (m *mockLogger) Fatal(message string)   {}

// Mock event bus - implements domain.EventBus
type mockEventBus struct{}

func (m *mockEventBus) Emit(eventName string, data ...interface{}) {}

// Mock settings repository - implements domain.SettingsRepository
type mockSettingsRepo struct {
	mu                sync.RWMutex
	customIgnoreRules string
	customPromptRules string
	useGitignore      bool
	useCustomIgnore   bool
	selectedProvider  string
	openAIKey         string
	geminiKey         string
	openRouterKey     string
	localAIKey        string
	localAIHost       string
	localAIModelName  string
	qwenAPIKey        string
	qwenHost          string
	selectedModels    map[string]string
	availableModels   map[string][]string
	recentProjects    []domain.RecentProjectInfo
	saveError         error
}

func newMockSettingsRepo() *mockSettingsRepo {
	return &mockSettingsRepo{
		useGitignore:    true,
		useCustomIgnore: false,
		selectedModels:  make(map[string]string),
		availableModels: make(map[string][]string),
		recentProjects:  []domain.RecentProjectInfo{},
	}
}

func (m *mockSettingsRepo) GetSettingsDTO() (domain.SettingsDTO, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return domain.SettingsDTO{
		CustomIgnoreRules: m.customIgnoreRules,
		CustomPromptRules: m.customPromptRules,
		UseGitignore:      m.useGitignore,
		UseCustomIgnore:   m.useCustomIgnore,
		SelectedProvider:  m.selectedProvider,
		OpenAIAPIKey:      m.openAIKey,
		GeminiAPIKey:      m.geminiKey,
		OpenRouterAPIKey:  m.openRouterKey,
		LocalAIAPIKey:     m.localAIKey,
		LocalAIHost:       m.localAIHost,
		LocalAIModelName:  m.localAIModelName,
		QwenAPIKey:        m.qwenAPIKey,
		SelectedModels:    m.selectedModels,
		AvailableModels:   m.availableModels,
	}, nil
}

func (m *mockSettingsRepo) GetCustomIgnoreRules() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.customIgnoreRules
}

func (m *mockSettingsRepo) SetCustomIgnoreRules(rules string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.customIgnoreRules = rules
}

func (m *mockSettingsRepo) GetCustomPromptRules() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.customPromptRules
}

func (m *mockSettingsRepo) SetCustomPromptRules(rules string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.customPromptRules = rules
}

func (m *mockSettingsRepo) GetUseGitignore() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.useGitignore
}

func (m *mockSettingsRepo) SetUseGitignore(use bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.useGitignore = use
}

func (m *mockSettingsRepo) GetUseCustomIgnore() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.useCustomIgnore
}

func (m *mockSettingsRepo) SetUseCustomIgnore(use bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.useCustomIgnore = use
}

func (m *mockSettingsRepo) GetSelectedAIProvider() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.selectedProvider
}

func (m *mockSettingsRepo) SetSelectedAIProvider(provider string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.selectedProvider = provider
}

func (m *mockSettingsRepo) GetOpenAIKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.openAIKey
}

func (m *mockSettingsRepo) SetOpenAIKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.openAIKey = key
}

func (m *mockSettingsRepo) GetGeminiKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.geminiKey
}

func (m *mockSettingsRepo) SetGeminiKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.geminiKey = key
}

func (m *mockSettingsRepo) GetOpenRouterKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.openRouterKey
}

func (m *mockSettingsRepo) SetOpenRouterKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.openRouterKey = key
}

func (m *mockSettingsRepo) GetLocalAIKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.localAIKey
}

func (m *mockSettingsRepo) SetLocalAIKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.localAIKey = key
}

func (m *mockSettingsRepo) GetLocalAIHost() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.localAIHost
}

func (m *mockSettingsRepo) SetLocalAIHost(host string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.localAIHost = host
}

func (m *mockSettingsRepo) GetLocalAIModelName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.localAIModelName
}

func (m *mockSettingsRepo) SetLocalAIModelName(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.localAIModelName = name
}

func (m *mockSettingsRepo) GetQwenKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.qwenAPIKey
}

func (m *mockSettingsRepo) SetQwenKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.qwenAPIKey = key
}

func (m *mockSettingsRepo) GetQwenHost() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.qwenHost
}

func (m *mockSettingsRepo) SetQwenHost(host string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.qwenHost = host
}

func (m *mockSettingsRepo) GetSelectedModel(provider string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.selectedModels[provider]
}

func (m *mockSettingsRepo) SetSelectedModel(provider, model string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.selectedModels[provider] = model
}

func (m *mockSettingsRepo) GetModels(provider string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.availableModels[provider]
}

func (m *mockSettingsRepo) SetModels(provider string, models []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.availableModels[provider] = models
}

func (m *mockSettingsRepo) GetRecentProjects() []domain.RecentProjectInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.recentProjects
}

func (m *mockSettingsRepo) AddRecentProject(path, name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recentProjects = append([]domain.RecentProjectInfo{{Path: path, Name: name}}, m.recentProjects...)
}

func (m *mockSettingsRepo) RemoveRecentProject(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	filtered := make([]domain.RecentProjectInfo, 0)
	for _, p := range m.recentProjects {
		if p.Path != path {
			filtered = append(filtered, p)
		}
	}
	m.recentProjects = filtered
}

func (m *mockSettingsRepo) Save() error {
	return m.saveError
}

func (m *mockSettingsRepo) Load() error {
	return nil
}

// Mock AI cache invalidator
type mockAICacheInvalidator struct {
	invalidateCalled bool
}

func (m *mockAICacheInvalidator) InvalidateProviderCache() {
	m.invalidateCalled = true
}

func TestNewService(t *testing.T) {
	repo := newMockSettingsRepo()
	svc, err := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}
	if svc == nil {
		t.Fatal("NewService returned nil service")
	}
}

func TestGetSettingsDTO(t *testing.T) {
	repo := newMockSettingsRepo()
	repo.customIgnoreRules = "*.log"
	repo.useGitignore = true
	repo.selectedProvider = "openai"

	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	dto, err := svc.GetSettingsDTO()
	if err != nil {
		t.Fatalf("GetSettingsDTO returned error: %v", err)
	}

	if dto.CustomIgnoreRules != "*.log" {
		t.Errorf("Expected CustomIgnoreRules '*.log', got '%s'", dto.CustomIgnoreRules)
	}
	if dto.UseGitignore != true {
		t.Error("Expected UseGitignore true")
	}
	if dto.SelectedProvider != "openai" {
		t.Errorf("Expected SelectedProvider 'openai', got '%s'", dto.SelectedProvider)
	}
}

func TestSaveSettingsDTO(t *testing.T) {
	repo := newMockSettingsRepo()
	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	dto := domain.SettingsDTO{
		CustomIgnoreRules: "*.tmp",
		UseGitignore:      false,
		UseCustomIgnore:   true,
		SelectedProvider:  "gemini",
		GeminiAPIKey:      "test-key",
		SelectedModels:    map[string]string{"gemini": "gemini-pro"},
		AvailableModels:   map[string][]string{"gemini": {"gemini-pro", "gemini-flash"}},
	}

	err := svc.SaveSettingsDTO(dto)
	if err != nil {
		t.Fatalf("SaveSettingsDTO returned error: %v", err)
	}

	if repo.customIgnoreRules != "*.tmp" {
		t.Errorf("Expected customIgnoreRules '*.tmp', got '%s'", repo.customIgnoreRules)
	}
	if repo.useGitignore != false {
		t.Error("Expected useGitignore false")
	}
	if repo.selectedProvider != "gemini" {
		t.Errorf("Expected selectedProvider 'gemini', got '%s'", repo.selectedProvider)
	}
}

func TestSaveSettingsDTO_InvalidatesAICache(t *testing.T) {
	repo := newMockSettingsRepo()
	repo.selectedProvider = "openai"
	repo.openAIKey = "old-key"

	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	invalidator := &mockAICacheInvalidator{}
	svc.SetAICacheInvalidator(invalidator)

	dto := domain.SettingsDTO{
		SelectedProvider: "openai",
		OpenAIAPIKey:     "new-key", // Changed
	}

	err := svc.SaveSettingsDTO(dto)
	if err != nil {
		t.Fatalf("SaveSettingsDTO returned error: %v", err)
	}

	if !invalidator.invalidateCalled {
		t.Error("Expected AI cache to be invalidated when API key changes")
	}
}

func TestOnIgnoreRulesChanged(t *testing.T) {
	repo := newMockSettingsRepo()
	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	callbackCalled := false
	svc.OnIgnoreRulesChanged(func() error {
		callbackCalled = true
		return nil
	})

	svc.SetCustomIgnoreRules("*.bak")

	if !callbackCalled {
		t.Error("Expected OnIgnoreRulesChanged callback to be called")
	}
}

func TestRecentProjects(t *testing.T) {
	repo := newMockSettingsRepo()
	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	// Add projects
	svc.AddRecentProject("/path/to/project1", "project1")
	svc.AddRecentProject("/path/to/project2", "project2")

	projects := svc.GetRecentProjects()
	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}

	// Most recent should be first
	if projects[0].Name != "project2" {
		t.Errorf("Expected first project to be 'project2', got '%s'", projects[0].Name)
	}

	// Remove project
	svc.RemoveRecentProject("/path/to/project1")
	projects = svc.GetRecentProjects()
	if len(projects) != 1 {
		t.Errorf("Expected 1 project after removal, got %d", len(projects))
	}
}

func TestGetSetUseGitignore(t *testing.T) {
	repo := newMockSettingsRepo()
	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	// Default should be true
	if !svc.GetUseGitignore() {
		t.Error("Expected default UseGitignore to be true")
	}

	svc.SetUseGitignore(false)
	if svc.GetUseGitignore() {
		t.Error("Expected UseGitignore to be false after setting")
	}
}

func TestGetSetUseCustomIgnore(t *testing.T) {
	repo := newMockSettingsRepo()
	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	// Default should be false
	if svc.GetUseCustomIgnore() {
		t.Error("Expected default UseCustomIgnore to be false")
	}

	svc.SetUseCustomIgnore(true)
	if !svc.GetUseCustomIgnore() {
		t.Error("Expected UseCustomIgnore to be true after setting")
	}
}

func TestGetSetCustomIgnoreRules(t *testing.T) {
	repo := newMockSettingsRepo()
	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	rules := "*.log\n*.tmp\nnode_modules/"
	svc.SetCustomIgnoreRules(rules)

	if svc.GetCustomIgnoreRules() != rules {
		t.Errorf("Expected rules '%s', got '%s'", rules, svc.GetCustomIgnoreRules())
	}
}

func TestSaveSettingsDTO_SaveError(t *testing.T) {
	repo := newMockSettingsRepo()
	repo.saveError = fmt.Errorf("disk full")

	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	dto := domain.SettingsDTO{
		CustomIgnoreRules: "*.tmp",
	}

	err := svc.SaveSettingsDTO(dto)
	if err == nil {
		t.Error("Expected error when save fails")
	}
	if !strings.Contains(err.Error(), "failed to save settings") {
		t.Errorf("Expected 'failed to save settings' error, got: %v", err)
	}
}

func TestSaveSettingsDTO_NoAICacheInvalidation(t *testing.T) {
	repo := newMockSettingsRepo()
	repo.selectedProvider = "openai"
	repo.openAIKey = "same-key"

	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	invalidator := &mockAICacheInvalidator{}
	svc.SetAICacheInvalidator(invalidator)

	// Same settings - no change
	dto := domain.SettingsDTO{
		SelectedProvider: "openai",
		OpenAIAPIKey:     "same-key",
	}

	err := svc.SaveSettingsDTO(dto)
	if err != nil {
		t.Fatalf("SaveSettingsDTO returned error: %v", err)
	}

	if invalidator.invalidateCalled {
		t.Error("AI cache should NOT be invalidated when settings unchanged")
	}
}

func TestRefreshModels_NoFetcher(t *testing.T) {
	repo := newMockSettingsRepo()
	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	err := svc.RefreshModels("unknown-provider", "api-key")
	if err == nil {
		t.Error("Expected error for unknown provider")
	}
	if !strings.Contains(err.Error(), "no model fetcher") {
		t.Errorf("Expected 'no model fetcher' error, got: %v", err)
	}
}

func TestRefreshModels_Success(t *testing.T) {
	repo := newMockSettingsRepo()

	fetchers := domain.ModelFetcherRegistry{
		"openai": func(apiKey string) ([]string, error) {
			return []string{"gpt-4", "gpt-3.5-turbo"}, nil
		},
	}

	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, fetchers)

	err := svc.RefreshModels("openai", "test-key")
	if err != nil {
		t.Fatalf("RefreshModels returned error: %v", err)
	}

	models := repo.GetModels("openai")
	if len(models) != 2 {
		t.Errorf("Expected 2 models, got %d", len(models))
	}
}

func TestRefreshModels_InvalidCurrentModel(t *testing.T) {
	repo := newMockSettingsRepo()
	repo.SetSelectedModel("openai", "old-model")

	fetchers := domain.ModelFetcherRegistry{
		"openai": func(apiKey string) ([]string, error) {
			return []string{"gpt-4", "gpt-3.5-turbo"}, nil
		},
	}

	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, fetchers)

	err := svc.RefreshModels("openai", "test-key")
	if err != nil {
		t.Fatalf("RefreshModels returned error: %v", err)
	}

	// Should auto-select first model when current is invalid
	selected := repo.GetSelectedModel("openai")
	if selected != "gpt-4" {
		t.Errorf("Expected selected model 'gpt-4', got '%s'", selected)
	}
}

func TestRefreshModels_FetcherError(t *testing.T) {
	repo := newMockSettingsRepo()

	fetchers := domain.ModelFetcherRegistry{
		"openai": func(apiKey string) ([]string, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, fetchers)

	err := svc.RefreshModels("openai", "test-key")
	if err == nil {
		t.Error("Expected error when fetcher fails")
	}
}

func TestMultipleIgnoreRulesCallbacks(t *testing.T) {
	repo := newMockSettingsRepo()
	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	callback1Called := false
	callback2Called := false

	svc.OnIgnoreRulesChanged(func() error {
		callback1Called = true
		return nil
	})

	svc.OnIgnoreRulesChanged(func() error {
		callback2Called = true
		return nil
	})

	svc.SetUseGitignore(false)

	if !callback1Called {
		t.Error("Expected callback1 to be called")
	}
	if !callback2Called {
		t.Error("Expected callback2 to be called")
	}
}

func TestSave(t *testing.T) {
	repo := newMockSettingsRepo()
	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	err := svc.Save()
	if err != nil {
		t.Errorf("Save returned error: %v", err)
	}
}

func TestSave_Error(t *testing.T) {
	repo := newMockSettingsRepo()
	repo.saveError = fmt.Errorf("write error")

	svc, _ := NewService(&mockLogger{}, &mockEventBus{}, repo, nil)

	err := svc.Save()
	if err == nil {
		t.Error("Expected error when save fails")
	}
}
