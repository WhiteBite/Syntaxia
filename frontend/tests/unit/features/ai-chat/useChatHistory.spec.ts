import { useChatHistory } from '@/features/ai-chat/composables/useChatHistory'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock dependencies
vi.mock('@/composables/useI18n', () => ({
    useI18n: () => ({
        t: (key: string) => {
            const translations: Record<string, string> = {
                'chat.history.yesterday': 'Yesterday',
                'chat.history.daysAgo': 'days ago'
            }
            return translations[key] || key
        }
    })
}))

vi.mock('@/features/ai-chat/model/chat.store', () => ({
    useChatStore: () => ({
        chatHistory: [],
        currentChatId: null,
        clearChat: vi.fn(),
        loadChat: vi.fn(),
        deleteChat: vi.fn(),
        loadHistory: vi.fn()
    })
}))

// Mock localStorage
const localStorageMock = {
    store: {} as Record<string, string>,
    getItem: vi.fn((key: string) => localStorageMock.store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
        localStorageMock.store[key] = value
    }),
    clear: vi.fn(() => {
        localStorageMock.store = {}
    }),
}
Object.defineProperty(global, 'localStorage', { value: localStorageMock })

describe('useChatHistory', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorageMock.clear()
        vi.clearAllMocks()
    })

    it('should initialize with closed state', () => {
        const { isOpen } = useChatHistory()
        expect(isOpen.value).toBe(false)
    })

    it('should format today date as time', () => {
        const { formatDate } = useChatHistory()
        const now = new Date().toISOString()
        const result = formatDate(now)
        // Should match HH:MM format
        expect(result).toMatch(/\d{1,2}:\d{2}/)
    })

    it('should format yesterday date', () => {
        const { formatDate } = useChatHistory()
        const yesterday = new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString()
        const result = formatDate(yesterday)
        // Should contain "Yesterday" from mocked translation
        expect(result).toBe('Yesterday')
    })

    it('should format dates within a week as days ago', () => {
        const { formatDate } = useChatHistory()
        const threeDaysAgo = new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString()
        const result = formatDate(threeDaysAgo)
        expect(result).toBe('3 days ago')
    })

    it('should format older dates as month/day', () => {
        const { formatDate } = useChatHistory()
        const twoWeeksAgo = new Date(Date.now() - 14 * 24 * 60 * 60 * 1000).toISOString()
        const result = formatDate(twoWeeksAgo)
        // Should be formatted as short date (e.g., "Jan 15")
        expect(result).toBeTruthy()
        expect(result).not.toMatch(/\d{1,2}:\d{2}/) // Not a time
        expect(result).not.toBe('Yesterday')
    })

    it('should close dropdown on new chat', () => {
        const { isOpen, handleNewChat } = useChatHistory()
        isOpen.value = true
        handleNewChat()
        expect(isOpen.value).toBe(false)
    })

    it('should toggle dropdown visibility', () => {
        const { isOpen, toggleDropdown } = useChatHistory()
        expect(isOpen.value).toBe(false)

        toggleDropdown()
        expect(isOpen.value).toBe(true)

        toggleDropdown()
        expect(isOpen.value).toBe(false)
    })

    it('should close dropdown explicitly', () => {
        const { isOpen, closeDropdown } = useChatHistory()
        isOpen.value = true

        closeDropdown()
        expect(isOpen.value).toBe(false)
    })

    it('should close dropdown on load chat', () => {
        const { isOpen, handleLoadChat } = useChatHistory()
        isOpen.value = true

        handleLoadChat('chat-123')
        expect(isOpen.value).toBe(false)
    })

    it('should have computed properties from store', () => {
        const { chatHistory, currentChatId, historyCount, hasHistory } = useChatHistory()

        expect(chatHistory.value).toEqual([])
        expect(currentChatId.value).toBeNull()
        expect(historyCount.value).toBe(0)
        expect(hasHistory.value).toBe(false)
    })
})
