import { useI18n } from '@/composables/useI18n'
import { apiService } from '@/services/api.service'
import { useAIStore } from '@/stores/ai.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, reactive, ref } from 'vue'

export interface AIProvider {
    id: string
    name: string
    icon: string
    description: string
}

export interface AISettingsState {
    selectedProvider: string
    openAIAPIKey: string
    geminiAPIKey: string
    qwenAPIKey: string
    openRouterAPIKey: string
    localAIAPIKey: string
    localAIHost: string
    qwenHost: string
    selectedModels: Record<string, string>
    availableModels: Record<string, string[]>
}

const DEFAULT_MODELS: Record<string, string[]> = {
    openai: ['gpt-4o', 'gpt-4o-mini', 'gpt-4-turbo', 'gpt-4', 'gpt-3.5-turbo'],
    gemini: ['gemini-1.5-pro', 'gemini-1.5-flash', 'gemini-pro'],
    qwen: ['qwen-max', 'qwen-plus', 'qwen-turbo', 'qwen-long'],
    'qwen-cli': ['qwen-coder-plus-latest', 'qwen-coder-turbo-latest', 'qwen-turbo-latest'],
    openrouter: ['anthropic/claude-3.5-sonnet', 'openai/gpt-4o', 'google/gemini-pro'],
    localai: ['llama3', 'mistral', 'codellama']
}

export function useAISettings() {
    const { t } = useI18n()
    const uiStore = useUIStore()
    const aiStore = useAIStore()

    const settings = reactive<AISettingsState>({
        selectedProvider: '',
        openAIAPIKey: '',
        geminiAPIKey: '',
        qwenAPIKey: '',
        openRouterAPIKey: '',
        localAIAPIKey: '',
        localAIHost: 'http://localhost:8080',
        qwenHost: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
        selectedModels: {},
        availableModels: {}
    })

    const showApiKey = ref(false)
    const isSaving = ref(false)
    const statusMessage = ref('')
    const statusType = ref<'success' | 'error'>('success')

    const providers = computed<AIProvider[]>(() => [
        { id: 'openai', name: 'OpenAI', icon: '🤖', description: t('settings.provider.openai') },
        { id: 'gemini', name: 'Google Gemini', icon: '✨', description: t('settings.provider.gemini') },
        { id: 'qwen', name: 'Qwen (Alibaba)', icon: '🌐', description: t('settings.provider.qwen') },
        { id: 'qwen-cli', name: 'Qwen Code CLI', icon: '💻', description: t('settings.provider.qwenCli') },
        { id: 'openrouter', name: 'OpenRouter', icon: '🔀', description: t('settings.provider.openrouter') },
        { id: 'localai', name: 'LocalAI', icon: '🏠', description: t('settings.provider.localai') },
    ])

    const currentProvider = computed(() =>
        providers.value.find(p => p.id === settings.selectedProvider)
    )

    const currentProviderHint = computed(() => {
        const hints: Record<string, string> = {
            openai: t('settings.hint.openai'),
            gemini: t('settings.hint.gemini'),
            qwen: t('settings.hint.qwen'),
            openrouter: t('settings.hint.openrouter'),
            localai: t('settings.hint.localai')
        }
        return hints[settings.selectedProvider] || ''
    })

    const currentApiKey = computed(() => {
        const keys: Record<string, string> = {
            openai: settings.openAIAPIKey,
            gemini: settings.geminiAPIKey,
            qwen: settings.qwenAPIKey,
            openrouter: settings.openRouterAPIKey,
            localai: settings.localAIAPIKey
        }
        return keys[settings.selectedProvider] || ''
    })

    const availableModels = computed(() =>
        settings.availableModels[settings.selectedProvider] ||
        DEFAULT_MODELS[settings.selectedProvider] || []
    )

    const selectedModel = computed(() =>
        settings.selectedModels[settings.selectedProvider] || availableModels.value[0] || ''
    )

    const needsApiKey = computed(() =>
        settings.selectedProvider && settings.selectedProvider !== 'qwen-cli'
    )

    const needsHostUrl = computed(() =>
        settings.selectedProvider === 'localai' || settings.selectedProvider === 'qwen'
    )

    const currentHostUrl = computed(() =>
        settings.selectedProvider === 'localai' ? settings.localAIHost : settings.qwenHost
    )

    const hostPlaceholder = computed(() =>
        settings.selectedProvider === 'localai'
            ? 'http://localhost:8080'
            : 'https://dashscope.aliyuncs.com/compatible-mode/v1'
    )

    function selectProvider(providerId: string) {
        settings.selectedProvider = providerId
    }

    function updateApiKey(value: string) {
        switch (settings.selectedProvider) {
            case 'openai': settings.openAIAPIKey = value; break
            case 'gemini': settings.geminiAPIKey = value; break
            case 'qwen': settings.qwenAPIKey = value; break
            case 'openrouter': settings.openRouterAPIKey = value; break
            case 'localai': settings.localAIAPIKey = value; break
        }
    }

    function clearApiKey() {
        updateApiKey('')
    }

    function updateModel(model: string) {
        settings.selectedModels[settings.selectedProvider] = model
    }

    function updateHost(value: string) {
        if (settings.selectedProvider === 'localai') {
            settings.localAIHost = value
        } else if (settings.selectedProvider === 'qwen') {
            settings.qwenHost = value
        }
    }

    function toggleShowApiKey() {
        showApiKey.value = !showApiKey.value
    }

    async function loadSettings() {
        try {
            const dto = await apiService.getSettings()
            settings.selectedProvider = dto.selectedProvider || 'openai'
            settings.openAIAPIKey = dto.openAIAPIKey || ''
            settings.geminiAPIKey = dto.geminiAPIKey || ''
            settings.qwenAPIKey = dto.qwenAPIKey || ''
            settings.openRouterAPIKey = dto.openRouterAPIKey || ''
            settings.localAIAPIKey = dto.localAIAPIKey || ''
            settings.localAIHost = dto.localAIHost || 'http://localhost:8080'
            settings.qwenHost = dto.qwenHost || 'https://dashscope.aliyuncs.com/compatible-mode/v1'
            settings.selectedModels = dto.selectedModels || {}
            settings.availableModels = dto.availableModels || {}
        } catch (e) {
            console.error('Failed to load settings:', e)
        }
    }

    async function saveSettings() {
        isSaving.value = true
        statusMessage.value = ''

        try {
            const dto = await apiService.getSettings()

            dto.selectedProvider = settings.selectedProvider
            dto.openAIAPIKey = settings.openAIAPIKey
            dto.geminiAPIKey = settings.geminiAPIKey
            dto.qwenAPIKey = settings.qwenAPIKey
            dto.openRouterAPIKey = settings.openRouterAPIKey
            dto.localAIAPIKey = settings.localAIAPIKey
            dto.localAIHost = settings.localAIHost
            dto.qwenHost = settings.qwenHost
            dto.selectedModels = settings.selectedModels

            await apiService.saveSettings(JSON.stringify(dto))

            const model = settings.selectedModels[settings.selectedProvider] || ''
            aiStore.updateModel(settings.selectedProvider, model)

            statusMessage.value = t('settings.saved')
            statusType.value = 'success'
            uiStore.addToast(t('settings.saved'), 'success')

            setTimeout(() => { statusMessage.value = '' }, 3000)
        } catch (e) {
            const errorMsg = e instanceof Error ? e.message : 'Unknown error'
            statusMessage.value = errorMsg
            statusType.value = 'error'
            uiStore.addToast(t('settings.saveFailed'), 'error')
        } finally {
            isSaving.value = false
        }
    }

    return {
        // State
        settings,
        showApiKey,
        isSaving,
        statusMessage,
        statusType,
        // Computed
        providers,
        currentProvider,
        currentProviderHint,
        currentApiKey,
        availableModels,
        selectedModel,
        needsApiKey,
        needsHostUrl,
        currentHostUrl,
        hostPlaceholder,
        // Actions
        selectProvider,
        updateApiKey,
        clearApiKey,
        updateModel,
        updateHost,
        toggleShowApiKey,
        loadSettings,
        saveSettings
    }
}
