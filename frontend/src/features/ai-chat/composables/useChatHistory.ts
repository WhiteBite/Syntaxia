/**
 * useChatHistory composable
 * Manages chat history dropdown state and operations
 */

import { useI18n } from '@/composables/useI18n'
import { useDebounceFn } from '@vueuse/core'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useChatStore } from '../model/chat.store'

export function useChatHistory() {
    const chatStore = useChatStore()
    const { t } = useI18n()

    // State
    const isOpen = ref(false)
    const searchQuery = ref('')
    const debouncedQuery = ref('')

    // Debounce search input (200ms)
    const updateDebouncedQuery = useDebounceFn((value: string) => {
        debouncedQuery.value = value
    }, 200)

    watch(searchQuery, (value) => updateDebouncedQuery(value))

    // Computed
    const chatHistory = computed(() => chatStore.chatHistory)
    const currentChatId = computed(() => chatStore.currentChatId)
    const historyCount = computed(() => chatStore.chatHistory.length)
    const hasHistory = computed(() => chatStore.chatHistory.length > 0)

    // Filtered history based on search query
    const filteredHistory = computed(() => {
        const query = debouncedQuery.value.toLowerCase().trim()
        if (!query) return chatStore.chatHistory

        return chatStore.chatHistory.filter(chat => {
            // Search in title
            if (chat.title.toLowerCase().includes(query)) return true
            // Search in first user message
            const firstUserMsg = chat.messages.find(m => m.role === 'user')
            if (firstUserMsg?.content.toLowerCase().includes(query)) return true
            return false
        })
    })

    const hasFilteredResults = computed(() => filteredHistory.value.length > 0)

    /**
     * Format date for display with relative time
     */
    function formatDate(dateStr: string): string {
        const date = new Date(dateStr)
        const now = new Date()
        const diffMs = now.getTime() - date.getTime()
        const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

        if (diffDays === 0) {
            return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
        } else if (diffDays === 1) {
            return t('chat.history.yesterday')
        } else if (diffDays < 7) {
            return `${diffDays} ${t('chat.history.daysAgo')}`
        } else {
            return date.toLocaleDateString([], { month: 'short', day: 'numeric' })
        }
    }

    /**
     * Create new chat
     */
    function handleNewChat(): void {
        chatStore.clearChat()
        isOpen.value = false
    }

    /**
     * Load chat by ID
     */
    function handleLoadChat(chatId: string): void {
        chatStore.loadChat(chatId)
        isOpen.value = false
    }

    /**
     * Delete chat by ID
     */
    function handleDeleteChat(chatId: string): void {
        chatStore.deleteChat(chatId)
    }

    /**
     * Toggle dropdown visibility
     */
    function toggleDropdown(): void {
        isOpen.value = !isOpen.value
    }

    /**
     * Close dropdown
     */
    function closeDropdown(): void {
        isOpen.value = false
    }

    /**
     * Clear search query
     */
    function clearSearch(): void {
        searchQuery.value = ''
        debouncedQuery.value = ''
    }

    /**
     * Handle keyboard events (Escape to close/clear search)
     */
    function handleKeydown(e: KeyboardEvent): void {
        if (e.key === 'Escape' && isOpen.value) {
            if (searchQuery.value) {
                clearSearch()
            } else {
                isOpen.value = false
            }
        }
    }

    // Lifecycle
    onMounted(() => {
        chatStore.loadHistory()
        window.addEventListener('keydown', handleKeydown)
    })

    onUnmounted(() => {
        window.removeEventListener('keydown', handleKeydown)
    })

    return {
        // State
        isOpen,
        searchQuery,

        // Computed
        chatHistory,
        filteredHistory,
        currentChatId,
        historyCount,
        hasHistory,
        hasFilteredResults,

        // Actions
        formatDate,
        handleNewChat,
        handleLoadChat,
        handleDeleteChat,
        toggleDropdown,
        closeDropdown,
        clearSearch
    }
}
