import { useChatStore } from '@/features/ai-chat/model/chat.store'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock dependencies
vi.mock('@/services/api.service', () => ({
    apiService: {
        agenticChat: vi.fn().mockResolvedValue({
            response: 'AI response',
            toolCalls: []
        }),
        generateCodeStream: vi.fn()
    }
}))

vi.mock('@/stores/ui.store', () => ({
    useUIStore: () => ({
        addToast: vi.fn()
    })
}))

vi.mock('@/stores/project.store', () => ({
    useProjectStore: () => ({
        currentPath: '/test/project'
    })
}))

vi.mock('@/composables/useI18n', () => ({
    useI18n: () => ({
        t: (key: string) => key
    })
}))

// Mock localStorage
const localStorageMock = {
    store: {} as Record<string, string>,
    getItem: vi.fn((key: string) => localStorageMock.store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
        localStorageMock.store[key] = value
    }),
    clear: vi.fn(() => { localStorageMock.store = {} }),
}
Object.defineProperty(global, 'localStorage', { value: localStorageMock })

describe('chat.store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorageMock.clear()
        vi.clearAllMocks()
    })

    it('sendMessage adds user message', async () => {
        const store = useChatStore()

        await store.sendMessage('Hello AI')

        expect(store.messages.length).toBeGreaterThanOrEqual(1)
        expect(store.messages[0].role).toBe('user')
        expect(store.messages[0].content).toBe('Hello AI')
    })

    it('sendMessage calls backend API', async () => {
        const { apiService } = await import('@/services/api.service')
        const store = useChatStore()

        await store.sendMessage('Test message')

        expect(apiService.agenticChat).toHaveBeenCalledWith('Test message', '/test/project')
    })

    it('clearMessages clears history', async () => {
        const store = useChatStore()

        await store.sendMessage('Message 1')
        expect(store.messages.length).toBeGreaterThan(0)

        store.clearChat()

        expect(store.messages.length).toBe(0)
        expect(store.currentChatId).toBeNull()
    })

    it('saveChat saves to localStorage', async () => {
        const store = useChatStore()

        await store.sendMessage('Test message')

        expect(localStorageMock.setItem).toHaveBeenCalled()
        const savedKey = localStorageMock.setItem.mock.calls[0][0]
        expect(savedKey).toContain('chat-history')
    })

    it('loadHistory loads from localStorage', async () => {
        const store = useChatStore()
        const mockHistory = [{
            id: 'chat-1',
            title: 'Test Chat',
            messages: [{ id: 'msg-1', role: 'user', content: 'Hello', timestamp: new Date().toISOString() }],
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
        }]

        localStorageMock.store['chat-history-_test_project'] = JSON.stringify(mockHistory)

        await store.loadHistory()

        expect(store.chatHistory.length).toBe(1)
        expect(store.chatHistory[0].title).toBe('Test Chat')
    })
})
