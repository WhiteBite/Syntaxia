import { useI18n } from '@/composables/useI18n'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export interface Message {
    id: string
    role: 'user' | 'assistant'
    content: string
    timestamp: string
    contextId?: string
    toolCalls?: ToolCallLog[]
}

export interface ToolCallLog {
    tool: string
    args?: Record<string, unknown>
    arguments?: string
    result?: string
    error?: string
    duration?: number
}

export interface ChatHistory {
    id: string
    title: string
    messages: Message[]
    createdAt: string
    updatedAt: string
}

const CHAT_HISTORY_KEY = 'chat-history'
const MAX_SAVED_CHATS = 10 // Chat-specific limit, not from CACHE

export const useChatStore = defineStore('chat', () => {
    const uiStore = useUIStore()
    const projectStore = useProjectStore()
    const { t } = useI18n()

    // State
    const messages = ref<Message[]>([])
    const isStreaming = ref(false)
    const currentModel = ref<string>('gpt-4')
    const chatHistory = ref<ChatHistory[]>([])
    const currentChatId = ref<string | null>(null)
    const streamingContent = ref<string>('')

    // Computed
    const hasMessages = computed(() => messages.value.length > 0)
    const lastMessage = computed(() => messages.value[messages.value.length - 1])

    // Get storage key based on project path
    function getStorageKey(): string {
        const projectPath = projectStore.currentPath || 'default'
        const safePath = projectPath.replace(/[^a-zA-Z0-9]/g, '_')
        return `${CHAT_HISTORY_KEY}-${safePath}`
    }

    // Actions
    async function sendMessage(content: string, contextId?: string): Promise<void> {
        const userMessage: Message = {
            id: `msg-${Date.now()}`,
            role: 'user',
            content,
            timestamp: new Date().toISOString(),
            contextId
        }

        messages.value.push(userMessage)

        try {
            const projectRoot = projectStore.currentPath || ''
            const response = await apiService.agenticChat(content, projectRoot)

            const assistantMessage: Message = {
                id: `msg-${Date.now()}-ai`,
                role: 'assistant',
                content: response.response,
                timestamp: new Date().toISOString(),
                toolCalls: response.toolCalls
            }

            messages.value.push(assistantMessage)
            await saveChat()
        } catch (error) {
            console.error('[ChatStore] sendMessage error:', error)
            uiStore.addToast(t('chat.error'), 'error')

            // Add error message
            messages.value.push({
                id: `msg-${Date.now()}-error`,
                role: 'assistant',
                content: t('error.generic'),
                timestamp: new Date().toISOString()
            })
        }
    }

    async function streamMessage(content: string, contextId?: string): Promise<void> {
        const userMessage: Message = {
            id: `msg-${Date.now()}`,
            role: 'user',
            content,
            timestamp: new Date().toISOString(),
            contextId
        }

        messages.value.push(userMessage)
        isStreaming.value = true
        streamingContent.value = ''

        // Add placeholder for streaming response
        const assistantMessage: Message = {
            id: `msg-${Date.now()}-ai`,
            role: 'assistant',
            content: '',
            timestamp: new Date().toISOString()
        }
        messages.value.push(assistantMessage)
        const streamingIndex = messages.value.length - 1

        // Subscribe to stream events
        const handleChunk = (event: CustomEvent) => {
            const chunk = event.detail
            if (chunk.content) {
                streamingContent.value += chunk.content
                messages.value[streamingIndex].content = streamingContent.value
            }
            if (chunk.done) {
                isStreaming.value = false
                window.removeEventListener('ai:stream:chunk', handleChunk as EventListener)
                saveChat()
            }
            if (chunk.error) {
                isStreaming.value = false
                uiStore.addToast(chunk.error, 'error')
                window.removeEventListener('ai:stream:chunk', handleChunk as EventListener)
            }
        }

        window.addEventListener('ai:stream:chunk', handleChunk as EventListener)

        try {
            const systemPrompt = 'You are a helpful coding assistant.'
            apiService.generateCodeStream(systemPrompt, content)
        } catch (error) {
            console.error('[ChatStore] streamMessage error:', error)
            isStreaming.value = false
            uiStore.addToast(t('chat.error'), 'error')
        }
    }

    function clearChat(): void {
        messages.value = []
        currentChatId.value = null
        streamingContent.value = ''
    }

    async function loadHistory(): Promise<void> {
        try {
            const key = getStorageKey()
            const saved = localStorage.getItem(key)
            if (saved) {
                const parsed = JSON.parse(saved) as ChatHistory[]
                chatHistory.value = parsed.slice(0, MAX_SAVED_CHATS)
            }
        } catch (error) {
            console.warn('[ChatStore] Failed to load history:', error)
        }
    }

    async function saveChat(): Promise<void> {
        if (messages.value.length === 0) return

        try {
            const key = getStorageKey()
            const chatId = currentChatId.value || `chat-${Date.now()}`
            const title = messages.value[0]?.content.slice(0, 50) || 'New Chat'

            const chat: ChatHistory = {
                id: chatId,
                title,
                messages: [...messages.value],
                createdAt: currentChatId.value
                    ? chatHistory.value.find(c => c.id === chatId)?.createdAt || new Date().toISOString()
                    : new Date().toISOString(),
                updatedAt: new Date().toISOString()
            }

            // Update or add chat
            const existingIndex = chatHistory.value.findIndex(c => c.id === chatId)
            if (existingIndex >= 0) {
                chatHistory.value[existingIndex] = chat
            } else {
                chatHistory.value.unshift(chat)
            }

            // Limit saved chats
            chatHistory.value = chatHistory.value.slice(0, MAX_SAVED_CHATS)
            currentChatId.value = chatId

            localStorage.setItem(key, JSON.stringify(chatHistory.value))
        } catch (error) {
            console.warn('[ChatStore] Failed to save chat:', error)
        }
    }

    async function loadChat(chatId: string): Promise<void> {
        const chat = chatHistory.value.find(c => c.id === chatId)
        if (chat) {
            messages.value = [...chat.messages]
            currentChatId.value = chatId
        }
    }

    function deleteMessage(messageId: string): void {
        messages.value = messages.value.filter(m => m.id !== messageId)
    }

    function editMessage(messageId: string, newContent: string): void {
        const message = messages.value.find(m => m.id === messageId)
        if (message) {
            message.content = newContent
        }
    }

    function stopStreaming(): void {
        isStreaming.value = false
        streamingContent.value = ''
    }

    // Delete chat from history
    function deleteChat(chatId: string): void {
        chatHistory.value = chatHistory.value.filter(c => c.id !== chatId)
        const key = getStorageKey()
        localStorage.setItem(key, JSON.stringify(chatHistory.value))

        if (currentChatId.value === chatId) {
            clearChat()
        }
    }

    return {
        // State
        messages,
        isStreaming,
        currentModel,
        chatHistory,
        currentChatId,
        streamingContent,
        // Computed
        hasMessages,
        lastMessage,
        // Actions
        sendMessage,
        streamMessage,
        clearChat,
        loadHistory,
        saveChat,
        loadChat,
        deleteMessage,
        editMessage,
        stopStreaming,
        deleteChat
    }
})
