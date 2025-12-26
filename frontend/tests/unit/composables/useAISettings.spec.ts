import { useAISettings, type AIProvider, type AISettingsState } from '@/composables/useAISettings'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock dependencies
vi.mock('@/services/api.service', () => ({
    apiService: {
        getSettings: vi.fn(),
        saveSettings: vi.fn()
    }
}))

const mockAddToast = vi.fn()
vi.mock('@/stores/ui.store', () => ({
    useUIStore: () => ({
        addToast: mockAddToast
    })
}))

const mockUpdateModel = vi.fn()
vi.mock('@/stores/ai.store', () => ({
    useAIStore: () => ({
        updateModel: mockUpdateModel
    })
}))

vi.mock('@/composables/useI18n', () => ({
    useI18n: () => ({
        t: (key: string) => key
    })
}))

vi.mock('@/composables/useLogger', () => ({
    useLogger: () => ({
        debug: vi.fn(),
        info: vi.fn(),
        warn: vi.fn(),
        error: vi.fn()
    })
}))

import { apiService } from '@/services/api.service'

describe('useAISettings', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    describe('initial state', () => {
        it('should have default settings state', () => {
            const { settings } = useAISettings()

            expect(settings.selectedProvider).toBe('')
            expect(settings.openAIAPIKey).toBe('')
            expect(settings.geminiAPIKey).toBe('')
            expect(settings.qwenAPIKey).toBe('')
            expect(settings.openRouterAPIKey).toBe('')
            expect(settings.localAIAPIKey).toBe('')
            expect(settings.localAIHost).toBe('http://localhost:8080')
            expect(settings.qwenHost).toBe('https://dashscope.aliyuncs.com/compatible-mode/v1')
        })

        it('should have showApiKey as false', () => {
            const { showApiKey } = useAISettings()
            expect(showApiKey.value).toBe(false)
        })

        it('should have isSaving as false', () => {
            const { isSaving } = useAISettings()
            expect(isSaving.value).toBe(false)
        })

        it('should have empty statusMessage', () => {
            const { statusMessage } = useAISettings()
            expect(statusMessage.value).toBe('')
        })
    })

    describe('providers computed', () => {
        it('should return list of AI providers', () => {
            const { providers } = useAISettings()

            expect(providers.value).toHaveLength(6)
            expect(providers.value.map(p => p.id)).toEqual([
                'openai', 'gemini', 'qwen', 'qwen-cli', 'openrouter', 'localai'
            ])
        })

        it('should have correct provider structure', () => {
            const { providers } = useAISettings()
            const openai = providers.value.find(p => p.id === 'openai')

            expect(openai).toBeDefined()
            expect(openai?.name).toBe('OpenAI')
            expect(openai?.icon).toBe('🤖')
            expect(openai?.description).toBe('settings.provider.openai')
        })
    })

    describe('currentProvider computed', () => {
        it('should return undefined when no provider selected', () => {
            const { currentProvider } = useAISettings()
            expect(currentProvider.value).toBeUndefined()
        })

        it('should return selected provider', () => {
            const { settings, currentProvider } = useAISettings()
            settings.selectedProvider = 'openai'

            expect(currentProvider.value?.id).toBe('openai')
            expect(currentProvider.value?.name).toBe('OpenAI')
        })
    })

    describe('currentApiKey computed', () => {
        it('should return empty string when no provider selected', () => {
            const { currentApiKey } = useAISettings()
            expect(currentApiKey.value).toBe('')
        })

        it('should return correct API key for openai', () => {
            const { settings, currentApiKey } = useAISettings()
            settings.selectedProvider = 'openai'
            settings.openAIAPIKey = 'sk-test-key'

            expect(currentApiKey.value).toBe('sk-test-key')
        })

        it('should return correct API key for gemini', () => {
            const { settings, currentApiKey } = useAISettings()
            settings.selectedProvider = 'gemini'
            settings.geminiAPIKey = 'gemini-key'

            expect(currentApiKey.value).toBe('gemini-key')
        })

        it('should return correct API key for qwen', () => {
            const { settings, currentApiKey } = useAISettings()
            settings.selectedProvider = 'qwen'
            settings.qwenAPIKey = 'qwen-key'

            expect(currentApiKey.value).toBe('qwen-key')
        })
    })

    describe('availableModels computed', () => {
        it('should return default models for openai', () => {
            const { settings, availableModels } = useAISettings()
            settings.selectedProvider = 'openai'

            expect(availableModels.value).toContain('gpt-4o')
            expect(availableModels.value).toContain('gpt-4o-mini')
        })

        it('should return default models for gemini', () => {
            const { settings, availableModels } = useAISettings()
            settings.selectedProvider = 'gemini'

            expect(availableModels.value).toContain('gemini-1.5-pro')
            expect(availableModels.value).toContain('gemini-1.5-flash')
        })

        it('should return custom models if available', () => {
            const { settings, availableModels } = useAISettings()
            settings.selectedProvider = 'openai'
            settings.availableModels = { openai: ['custom-model-1', 'custom-model-2'] }

            expect(availableModels.value).toEqual(['custom-model-1', 'custom-model-2'])
        })
    })

    describe('selectedModel computed', () => {
        it('should return first available model when none selected', () => {
            const { settings, selectedModel } = useAISettings()
            settings.selectedProvider = 'openai'

            expect(selectedModel.value).toBe('gpt-4o')
        })

        it('should return selected model', () => {
            const { settings, selectedModel } = useAISettings()
            settings.selectedProvider = 'openai'
            settings.selectedModels = { openai: 'gpt-4-turbo' }

            expect(selectedModel.value).toBe('gpt-4-turbo')
        })
    })

    describe('needsApiKey computed', () => {
        it('should return falsy when no provider selected', () => {
            const { needsApiKey } = useAISettings()
            expect(needsApiKey.value).toBeFalsy()
        })

        it('should return true for openai', () => {
            const { settings, needsApiKey } = useAISettings()
            settings.selectedProvider = 'openai'

            expect(needsApiKey.value).toBe(true)
        })

        it('should return false for qwen-cli', () => {
            const { settings, needsApiKey } = useAISettings()
            settings.selectedProvider = 'qwen-cli'

            expect(needsApiKey.value).toBe(false)
        })
    })

    describe('needsHostUrl computed', () => {
        it('should return false for openai', () => {
            const { settings, needsHostUrl } = useAISettings()
            settings.selectedProvider = 'openai'

            expect(needsHostUrl.value).toBe(false)
        })

        it('should return true for localai', () => {
            const { settings, needsHostUrl } = useAISettings()
            settings.selectedProvider = 'localai'

            expect(needsHostUrl.value).toBe(true)
        })

        it('should return true for qwen', () => {
            const { settings, needsHostUrl } = useAISettings()
            settings.selectedProvider = 'qwen'

            expect(needsHostUrl.value).toBe(true)
        })
    })

    describe('currentHostUrl computed', () => {
        it('should return localAIHost for localai provider', () => {
            const { settings, currentHostUrl } = useAISettings()
            settings.selectedProvider = 'localai'
            settings.localAIHost = 'http://custom:8080'

            expect(currentHostUrl.value).toBe('http://custom:8080')
        })

        it('should return qwenHost for qwen provider', () => {
            const { settings, currentHostUrl } = useAISettings()
            settings.selectedProvider = 'qwen'
            settings.qwenHost = 'https://custom-qwen.com'

            expect(currentHostUrl.value).toBe('https://custom-qwen.com')
        })
    })

    describe('selectProvider', () => {
        it('should update selectedProvider', () => {
            const { settings, selectProvider } = useAISettings()

            selectProvider('gemini')

            expect(settings.selectedProvider).toBe('gemini')
        })
    })

    describe('updateApiKey', () => {
        it('should update openai API key', () => {
            const { settings, updateApiKey } = useAISettings()
            settings.selectedProvider = 'openai'

            updateApiKey('new-openai-key')

            expect(settings.openAIAPIKey).toBe('new-openai-key')
        })

        it('should update gemini API key', () => {
            const { settings, updateApiKey } = useAISettings()
            settings.selectedProvider = 'gemini'

            updateApiKey('new-gemini-key')

            expect(settings.geminiAPIKey).toBe('new-gemini-key')
        })

        it('should update qwen API key', () => {
            const { settings, updateApiKey } = useAISettings()
            settings.selectedProvider = 'qwen'

            updateApiKey('new-qwen-key')

            expect(settings.qwenAPIKey).toBe('new-qwen-key')
        })

        it('should update openrouter API key', () => {
            const { settings, updateApiKey } = useAISettings()
            settings.selectedProvider = 'openrouter'

            updateApiKey('new-openrouter-key')

            expect(settings.openRouterAPIKey).toBe('new-openrouter-key')
        })

        it('should update localai API key', () => {
            const { settings, updateApiKey } = useAISettings()
            settings.selectedProvider = 'localai'

            updateApiKey('new-localai-key')

            expect(settings.localAIAPIKey).toBe('new-localai-key')
        })
    })

    describe('clearApiKey', () => {
        it('should clear current provider API key', () => {
            const { settings, clearApiKey } = useAISettings()
            settings.selectedProvider = 'openai'
            settings.openAIAPIKey = 'existing-key'

            clearApiKey()

            expect(settings.openAIAPIKey).toBe('')
        })
    })

    describe('updateModel', () => {
        it('should update selected model for current provider', () => {
            const { settings, updateModel } = useAISettings()
            settings.selectedProvider = 'openai'

            updateModel('gpt-4-turbo')

            expect(settings.selectedModels['openai']).toBe('gpt-4-turbo')
        })
    })

    describe('updateHost', () => {
        it('should update localAIHost for localai provider', () => {
            const { settings, updateHost } = useAISettings()
            settings.selectedProvider = 'localai'

            updateHost('http://new-host:9000')

            expect(settings.localAIHost).toBe('http://new-host:9000')
        })

        it('should update qwenHost for qwen provider', () => {
            const { settings, updateHost } = useAISettings()
            settings.selectedProvider = 'qwen'

            updateHost('https://new-qwen-host.com')

            expect(settings.qwenHost).toBe('https://new-qwen-host.com')
        })
    })

    describe('toggleShowApiKey', () => {
        it('should toggle showApiKey value', () => {
            const { showApiKey, toggleShowApiKey } = useAISettings()

            expect(showApiKey.value).toBe(false)

            toggleShowApiKey()
            expect(showApiKey.value).toBe(true)

            toggleShowApiKey()
            expect(showApiKey.value).toBe(false)
        })
    })

    describe('loadSettings', () => {
        it('should load settings from API', async () => {
            const mockSettings = {
                selectedProvider: 'gemini',
                openAIAPIKey: 'openai-key',
                geminiAPIKey: 'gemini-key',
                qwenAPIKey: 'qwen-key',
                openRouterAPIKey: 'openrouter-key',
                localAIAPIKey: 'localai-key',
                localAIHost: 'http://custom:8080',
                qwenHost: 'https://custom-qwen.com',
                selectedModels: { gemini: 'gemini-1.5-flash' },
                availableModels: { gemini: ['gemini-1.5-pro', 'gemini-1.5-flash'] }
            }
            vi.mocked(apiService.getSettings).mockResolvedValue(mockSettings)

            const { settings, loadSettings } = useAISettings()
            await loadSettings()

            expect(settings.selectedProvider).toBe('gemini')
            expect(settings.openAIAPIKey).toBe('openai-key')
            expect(settings.geminiAPIKey).toBe('gemini-key')
            expect(settings.localAIHost).toBe('http://custom:8080')
            expect(settings.selectedModels).toEqual({ gemini: 'gemini-1.5-flash' })
        })

        it('should use defaults for missing values', async () => {
            vi.mocked(apiService.getSettings).mockResolvedValue({})

            const { settings, loadSettings } = useAISettings()
            await loadSettings()

            expect(settings.selectedProvider).toBe('openai')
            expect(settings.localAIHost).toBe('http://localhost:8080')
        })

        it('should handle API error gracefully', async () => {
            vi.mocked(apiService.getSettings).mockRejectedValue(new Error('API Error'))

            const { loadSettings } = useAISettings()

            await expect(loadSettings()).resolves.not.toThrow()
        })
    })

    describe('saveSettings', () => {
        it('should save settings to API', async () => {
            const existingSettings = {
                selectedProvider: 'openai',
                openAIAPIKey: '',
                selectedModels: {}
            }
            vi.mocked(apiService.getSettings).mockResolvedValue(existingSettings)
            vi.mocked(apiService.saveSettings).mockResolvedValue(undefined)

            const { settings, saveSettings } = useAISettings()
            settings.selectedProvider = 'gemini'
            settings.geminiAPIKey = 'new-gemini-key'
            settings.selectedModels = { gemini: 'gemini-1.5-pro' }

            await saveSettings()

            expect(apiService.saveSettings).toHaveBeenCalled()
            expect(mockUpdateModel).toHaveBeenCalledWith('gemini', 'gemini-1.5-pro')
            expect(mockAddToast).toHaveBeenCalledWith('settings.saved', 'success')
        })

        it('should set isSaving during save operation', async () => {
            let resolvePromise: () => void
            const promise = new Promise<void>((resolve) => {
                resolvePromise = resolve
            })
            vi.mocked(apiService.getSettings).mockResolvedValue({})
            vi.mocked(apiService.saveSettings).mockReturnValue(promise)

            const { isSaving, saveSettings } = useAISettings()

            const savePromise = saveSettings()
            expect(isSaving.value).toBe(true)

            resolvePromise!()
            await savePromise

            expect(isSaving.value).toBe(false)
        })

        it('should show error toast on save failure', async () => {
            vi.mocked(apiService.getSettings).mockResolvedValue({})
            vi.mocked(apiService.saveSettings).mockRejectedValue(new Error('Save failed'))

            const { saveSettings, statusMessage, statusType } = useAISettings()
            await saveSettings()

            expect(statusMessage.value).toBe('Save failed')
            expect(statusType.value).toBe('error')
            expect(mockAddToast).toHaveBeenCalledWith('settings.saveFailed', 'error')
        })

        it('should set success status message on successful save', async () => {
            vi.mocked(apiService.getSettings).mockResolvedValue({})
            vi.mocked(apiService.saveSettings).mockResolvedValue(undefined)

            const { saveSettings, statusMessage, statusType } = useAISettings()
            await saveSettings()

            expect(statusMessage.value).toBe('settings.saved')
            expect(statusType.value).toBe('success')
        })
    })

    describe('hostPlaceholder computed', () => {
        it('should return localai placeholder for localai provider', () => {
            const { settings, hostPlaceholder } = useAISettings()
            settings.selectedProvider = 'localai'

            expect(hostPlaceholder.value).toBe('http://localhost:8080')
        })

        it('should return qwen placeholder for other providers', () => {
            const { settings, hostPlaceholder } = useAISettings()
            settings.selectedProvider = 'qwen'

            expect(hostPlaceholder.value).toBe('https://dashscope.aliyuncs.com/compatible-mode/v1')
        })
    })

    describe('currentProviderHint computed', () => {
        it('should return hint for openai', () => {
            const { settings, currentProviderHint } = useAISettings()
            settings.selectedProvider = 'openai'

            expect(currentProviderHint.value).toBe('settings.hint.openai')
        })

        it('should return hint for gemini', () => {
            const { settings, currentProviderHint } = useAISettings()
            settings.selectedProvider = 'gemini'

            expect(currentProviderHint.value).toBe('settings.hint.gemini')
        })

        it('should return empty string for unknown provider', () => {
            const { settings, currentProviderHint } = useAISettings()
            settings.selectedProvider = 'unknown'

            expect(currentProviderHint.value).toBe('')
        })
    })

    describe('AIProvider type', () => {
        it('should correctly type AIProvider interface', () => {
            const provider: AIProvider = {
                id: 'test',
                name: 'Test Provider',
                icon: '🧪',
                description: 'Test description'
            }

            expect(provider.id).toBe('test')
            expect(provider.name).toBe('Test Provider')
        })
    })

    describe('AISettingsState type', () => {
        it('should correctly type AISettingsState interface', () => {
            const state: AISettingsState = {
                selectedProvider: 'openai',
                openAIAPIKey: 'key',
                geminiAPIKey: '',
                qwenAPIKey: '',
                openRouterAPIKey: '',
                localAIAPIKey: '',
                localAIHost: 'http://localhost:8080',
                qwenHost: 'https://example.com',
                selectedModels: { openai: 'gpt-4o' },
                availableModels: { openai: ['gpt-4o'] }
            }

            expect(state.selectedProvider).toBe('openai')
            expect(state.selectedModels['openai']).toBe('gpt-4o')
        })
    })
})
