/**
 * useChatHistory composable
 * Manages chat history dropdown state and operations
 */

import { useI18n } from '@/composables/useI18n'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useChatStore } from '../model/chat.store'

export function useChatHistory() {
    const chatStore = useChatStore()
    const { t } = useI18n()

    // State
    const isOpen = ref(false)

    // Computed
    const chatHistory = computed(() => chatStore.chatHistory)
    const currentChatId = computed(() => chatStore.currentChatId)
    const historyCount = computed(() => chatStore.chatHistory.length)
    const hasHistory = computed(() => chatStore.chatHistory.length > 0)

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
     * Handle keyboard events (Escape to close)
     */
    function handleKeydown(e: KeyboardEvent): void {
        if (e.key === 'Escape' && isOpen.value) {
            isOpen.value = false
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

        // Computed
        chatHistory,
        currentChatId,
        historyCount,
        hasHistory,

        // Actions
        formatDate,
        handleNewChat,
        handleLoadChat,
        handleDeleteChat,
        toggleDropdown,
        closeDropdown
    }
}
