import { useLogger } from '@/composables/useLogger'
import { apiService } from '@/services/api.service'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

const logger = useLogger('AIStore')

// Model context limits mapping
export const MODEL_LIMITS: Record<string, number> = {
    // OpenAI
    'gpt-4o': 128000,
    'gpt-4o-mini': 128000,
    'gpt-4-turbo': 128000,
    'gpt-4': 8192,
    'gpt-3.5-turbo': 16385,
    // Gemini
    'gemini-1.5-pro': 1000000,
    'gemini-1.5-flash': 1000000,
    'gemini-pro': 32000,
    // Qwen
    'qwen-max': 32000,
    'qwen-plus': 131072,
    'qwen-turbo': 131072,
    'qwen-long': 1000000,
    'qwen-coder-plus-latest': 131072,
    'qwen-coder-turbo-latest': 131072,
    'qwen-turbo-latest': 131072,
    // Claude (OpenRouter)
    'anthropic/claude-3.5-sonnet': 200000,
    'anthropic/claude-3-opus': 200000,
    'openai/gpt-4o': 128000,
    'google/gemini-pro': 32000,
    // LocalAI defaults
    'llama3': 8192,
    'mistral': 32000,
    'codellama': 16384,
}

// Pricing per 1K tokens (input)
export const MODEL_PRICING: Record<string, number> = {
    'gpt-4o': 0.005,
    'gpt-4o-mini': 0.00015,
    'gpt-4-turbo': 0.01,
    'gpt-4': 0.03,
    'gpt-3.5-turbo': 0.0005,
    'gemini-1.5-pro': 0.00125,
    'gemini-1.5-flash': 0.000075,
    'gemini-pro': 0.00025,
    'qwen-max': 0.004,
    'qwen-plus': 0.0004,
    'qwen-turbo': 0.0002,
    'qwen-coder-plus-latest': 0.0004,
    'anthropic/claude-3.5-sonnet': 0.003,
    'anthropic/claude-3-opus': 0.015,
}

export interface AIProviderInfo {
    name: string
    connected: boolean
    model: string
    provider: string
}

export const useAIStore = defineStore('ai', () => {
    // State
    const providerInfo = ref<AIProviderInfo>({
        name: 'Not configured',
        connected: false,
        model: 'gpt-4o',
        provider: 'openai'
    })
    const isLoading = ref(false)
    const lastError = ref<string | null>(null)

    // Computed
    const currentModel = computed(() => providerInfo.value.model)
    const currentProvider = computed(() => providerInfo.value.provider)
    const isConnected = computed(() => providerInfo.value.connected)

    const contextLimit = computed(() => {
        return MODEL_LIMITS[currentModel.value] || 32000
    })

    const pricePerKTokens = computed(() => {
        return MODEL_PRICING[currentModel.value] || 0.001
    })

    // Actions
    async function loadProviderInfo() {
        isLoading.value = true
        lastError.value = null

        try {
            // Get provider info from backend
            const infoJson = await apiService.getProviderInfo()
            const info = JSON.parse(infoJson)

            // Get settings to find selected model
            const settings = await apiService.getSettings()
            const selectedProvider = settings.selectedProvider || 'openai'
            const selectedModel = settings.selectedModels?.[selectedProvider] ||
                info.SupportedModels?.[0] || 'gpt-4o'

            providerInfo.value = {
                name: info.Name || 'AI Provider',
                connected: true,
                model: selectedModel,
                provider: selectedProvider
            }

            logger.debug('Provider info loaded:', providerInfo.value)
        } catch (e) {
            const errorMsg = e instanceof Error ? e.message : 'Unknown error'
            console.error('[AIStore] Failed to load provider info:', errorMsg)
            lastError.value = errorMsg

            providerInfo.value = {
                name: 'Not configured',
                connected: false,
                model: 'gpt-4o',
                provider: 'openai'
            }
        } finally {
            isLoading.value = false
        }
    }

    // Update model (called when settings change)
    function updateModel(provider: string, model: string) {
        providerInfo.value.provider = provider
        providerInfo.value.model = model
        logger.debug('Model updated:', provider, model)
    }

    // Calculate cost for given tokens
    function calculateCost(tokens: number): number {
        return (tokens / 1000) * pricePerKTokens.value
    }

    // Get usage percentage
    function getUsagePercent(usedTokens: number): number {
        const percent = Math.round((usedTokens / contextLimit.value) * 100)
        return Math.min(percent, 100)
    }

    // Format tokens for display
    function formatTokens(tokens: number): string {
        if (tokens >= 1000000) return (tokens / 1000000).toFixed(1) + 'M'
        if (tokens >= 1000) return (tokens / 1000).toFixed(1) + 'K'
        return tokens.toString()
    }

    return {
        // State
        providerInfo,
        isLoading,
        lastError,
        // Computed
        currentModel,
        currentProvider,
        isConnected,
        contextLimit,
        pricePerKTokens,
        // Actions
        loadProviderInfo,
        updateModel,
        calculateCost,
        getUsagePercent,
        formatTokens
    }
})
