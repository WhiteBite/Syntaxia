import {
    MODEL_LIMITS,
    MODEL_PRICING,
    useAIStore,
    type AIProviderInfo
} from '@/stores/ai.store'
import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

// Mock apiService
vi.mock('@/services/api.service', () => ({
    apiService: {
        getProviderInfo: vi.fn(),
        getSettings: vi.fn()
    }
}))

// Mock logger
vi.mock('@/composables/useLogger', () => ({
    useLogger: () => ({
        debug: vi.fn(),
        info: vi.fn(),
        warn: vi.fn(),
        error: vi.fn()
    })
}))

import { apiService } from '@/services/api.service'

describe('AIStore', () => {
    let store: ReturnType<typeof useAIStore>
    let consoleSpy: {
        error: ReturnType<typeof vi.spyOn>
    }

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useAIStore()
        vi.clearAllMocks()

        consoleSpy = {
            error: vi.spyOn(console, 'error').mockImplementation(() => { })
        }
    })

    afterEach(() => {
        vi.restoreAllMocks()
    })

    describe('initial state', () => {
        it('should have default provider info', () => {
            expect(store.providerInfo).toEqual({
                name: 'Not configured',
                connected: false,
                model: 'gpt-4o',
                provider: 'openai'
            })
        })

        it('should have isLoading as false', () => {
            expect(store.isLoading).toBe(false)
        })

        it('should have lastError as null', () => {
            expect(store.lastError).toBeNull()
        })
    })

    describe('computed properties', () => {
        it('should return currentModel from providerInfo', () => {
            expect(store.currentModel).toBe('gpt-4o')
        })

        it('should return currentProvider from providerInfo', () => {
            expect(store.currentProvider).toBe('openai')
        })

        it('should return isConnected from providerInfo', () => {
            expect(store.isConnected).toBe(false)
        })

        it('should return contextLimit for known model', () => {
            store.updateModel('openai', 'gpt-4o')
            expect(store.contextLimit).toBe(128000)
        })

        it('should return default contextLimit for unknown model', () => {
            store.updateModel('custom', 'unknown-model')
            expect(store.contextLimit).toBe(32000)
        })

        it('should return pricePerKTokens for known model', () => {
            store.updateModel('openai', 'gpt-4o')
            expect(store.pricePerKTokens).toBe(0.005)
        })

        it('should return default pricePerKTokens for unknown model', () => {
            store.updateModel('custom', 'unknown-model')
            expect(store.pricePerKTokens).toBe(0.001)
        })
    })

    describe('MODEL_LIMITS constant', () => {
        it('should have correct limits for OpenAI models', () => {
            expect(MODEL_LIMITS['gpt-4o']).toBe(128000)
            expect(MODEL_LIMITS['gpt-4']).toBe(8192)
            expect(MODEL_LIMITS['gpt-3.5-turbo']).toBe(16385)
        })

        it('should have correct limits for Gemini models', () => {
            expect(MODEL_LIMITS['gemini-1.5-pro']).toBe(1000000)
            expect(MODEL_LIMITS['gemini-pro']).toBe(32000)
        })

        it('should have correct limits for Qwen models', () => {
            expect(MODEL_LIMITS['qwen-max']).toBe(32000)
            expect(MODEL_LIMITS['qwen-long']).toBe(1000000)
        })
    })

    describe('MODEL_PRICING constant', () => {
        it('should have correct pricing for OpenAI models', () => {
            expect(MODEL_PRICING['gpt-4o']).toBe(0.005)
            expect(MODEL_PRICING['gpt-4o-mini']).toBe(0.00015)
        })

        it('should have correct pricing for Gemini models', () => {
            expect(MODEL_PRICING['gemini-1.5-pro']).toBe(0.00125)
            expect(MODEL_PRICING['gemini-1.5-flash']).toBe(0.000075)
        })
    })

    describe('updateModel', () => {
        it('should update provider and model', () => {
            store.updateModel('gemini', 'gemini-1.5-pro')

            expect(store.currentProvider).toBe('gemini')
            expect(store.currentModel).toBe('gemini-1.5-pro')
        })

        it('should update contextLimit after model change', () => {
            store.updateModel('openai', 'gpt-4')
            expect(store.contextLimit).toBe(8192)

            store.updateModel('gemini', 'gemini-1.5-pro')
            expect(store.contextLimit).toBe(1000000)
        })
    })

    describe('calculateCost', () => {
        it('should calculate cost correctly', () => {
            store.updateModel('openai', 'gpt-4o') // price: 0.005 per 1K
            const cost = store.calculateCost(10000)
            expect(cost).toBe(0.05) // 10K tokens * 0.005
        })

        it('should calculate cost for zero tokens', () => {
            const cost = store.calculateCost(0)
            expect(cost).toBe(0)
        })

        it('should calculate cost for small token count', () => {
            store.updateModel('openai', 'gpt-4o') // price: 0.005 per 1K
            const cost = store.calculateCost(500)
            expect(cost).toBe(0.0025) // 0.5K tokens * 0.005
        })
    })

    describe('getUsagePercent', () => {
        it('should calculate usage percentage correctly', () => {
            store.updateModel('openai', 'gpt-4o') // limit: 128000
            const percent = store.getUsagePercent(64000)
            expect(percent).toBe(50)
        })

        it('should cap at 100%', () => {
            store.updateModel('openai', 'gpt-4o') // limit: 128000
            const percent = store.getUsagePercent(200000)
            expect(percent).toBe(100)
        })

        it('should return 0 for zero tokens', () => {
            const percent = store.getUsagePercent(0)
            expect(percent).toBe(0)
        })

        it('should round percentage', () => {
            store.updateModel('openai', 'gpt-4') // limit: 8192
            const percent = store.getUsagePercent(1000)
            expect(percent).toBe(12) // ~12.2% rounded
        })
    })

    describe('formatTokens', () => {
        it('should format millions', () => {
            expect(store.formatTokens(1500000)).toBe('1.5M')
            expect(store.formatTokens(1000000)).toBe('1.0M')
        })

        it('should format thousands', () => {
            expect(store.formatTokens(15000)).toBe('15.0K')
            expect(store.formatTokens(1500)).toBe('1.5K')
        })

        it('should return raw number for small values', () => {
            expect(store.formatTokens(500)).toBe('500')
            expect(store.formatTokens(0)).toBe('0')
        })
    })

    describe('loadProviderInfo', () => {
        it('should load provider info successfully', async () => {
            vi.mocked(apiService.getProviderInfo).mockResolvedValue(
                JSON.stringify({
                    Name: 'OpenAI',
                    SupportedModels: ['gpt-4o', 'gpt-4']
                })
            )
            vi.mocked(apiService.getSettings).mockResolvedValue({
                selectedProvider: 'openai',
                selectedModels: { openai: 'gpt-4o' }
            })

            await store.loadProviderInfo()

            expect(store.providerInfo.name).toBe('OpenAI')
            expect(store.providerInfo.connected).toBe(true)
            expect(store.providerInfo.model).toBe('gpt-4o')
            expect(store.providerInfo.provider).toBe('openai')
            expect(store.isLoading).toBe(false)
            expect(store.lastError).toBeNull()
        })

        it('should set isLoading during request', async () => {
            let resolvePromise: (value: string) => void
            const promise = new Promise<string>((resolve) => {
                resolvePromise = resolve
            })
            vi.mocked(apiService.getProviderInfo).mockReturnValue(promise)

            const loadPromise = store.loadProviderInfo()
            expect(store.isLoading).toBe(true)

            resolvePromise!(JSON.stringify({ Name: 'Test' }))
            vi.mocked(apiService.getSettings).mockResolvedValue({})
            await loadPromise

            expect(store.isLoading).toBe(false)
        })

        it('should handle error and reset to defaults', async () => {
            vi.mocked(apiService.getProviderInfo).mockRejectedValue(
                new Error('Connection failed')
            )

            await store.loadProviderInfo()

            expect(store.providerInfo.connected).toBe(false)
            expect(store.providerInfo.name).toBe('Not configured')
            expect(store.lastError).toBe('Connection failed')
            // Logger is mocked, so we just verify the error state is set correctly
        })

        it('should use default model when not in settings', async () => {
            vi.mocked(apiService.getProviderInfo).mockResolvedValue(
                JSON.stringify({
                    Name: 'Provider',
                    SupportedModels: ['model-1', 'model-2']
                })
            )
            vi.mocked(apiService.getSettings).mockResolvedValue({})

            await store.loadProviderInfo()

            expect(store.providerInfo.model).toBe('model-1')
        })
    })

    describe('AIProviderInfo type', () => {
        it('should correctly type AIProviderInfo interface', () => {
            const info: AIProviderInfo = {
                name: 'Test Provider',
                connected: true,
                model: 'test-model',
                provider: 'test'
            }

            expect(info.name).toBe('Test Provider')
            expect(info.connected).toBe(true)
        })
    })
})
