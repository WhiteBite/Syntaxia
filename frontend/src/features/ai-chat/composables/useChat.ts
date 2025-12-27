/**
 * useChat composable for AI Chat feature
 * Provides streaming responses, inline diff preview, real-time token counting, and tool call handling
 */

import { useI18n } from '@/composables/useI18n'
import { useContextStore } from '@/features/context'
import { useFileStore } from '@/features/files'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, onUnmounted, ref, shallowRef } from 'vue'
import {
    chatApi,
    type StreamChunk,
    type StreamSubscription,
    type StreamToolCall,
    type TokenCountInfo
} from '../api/chat.api'
import { useMentions } from './useMentions'

/** Chat message structure */
export interface ChatMessage {
    id: string
    role: 'user' | 'assistant'
    content: string
    timestamp: string
    isStreaming?: boolean
    toolCalls?: ChatToolCall[]
    tokenCount?: number
    error?: string
    diffPreview?: DiffPreview
}

/** Tool call in chat */
export interface ChatToolCall {
    id: string
    tool: string
    arguments: string
    status: 'pending' | 'executing' | 'completed' | 'failed'
    result?: string
    error?: string
    duration?: number
}

/** Inline diff preview */
export interface DiffPreview {
    filePath: string
    hunks: DiffHunk[]
    additions: number
    deletions: number
}

/** Diff hunk for inline preview */
export interface DiffHunk {
    oldStart: number
    oldLines: number
    newStart: number
    newLines: number
    lines: DiffLine[]
}

/** Single diff line */
export interface DiffLine {
    type: 'add' | 'remove' | 'context'
    content: string
    oldLineNumber?: number
    newLineNumber?: number
}

/** Token budget state */
export interface TokenBudget {
    used: number
    limit: number
    inputTokens: number
    outputTokens: number
}

/** Chat state */
interface ChatState {
    messages: ChatMessage[]
    isStreaming: boolean
    isAnalyzing: boolean
    tokenBudget: TokenBudget
    currentStreamingMessageId: string | null
}

const DEFAULT_TOKEN_LIMIT = 128000
const TOKEN_WARNING_THRESHOLD = 0.8

export function useChat() {
    const { t } = useI18n()
    const uiStore = useUIStore()
    const projectStore = useProjectStore()
    const contextStore = useContextStore()
    const fileStore = useFileStore()
    const mentions = useMentions()

    // State
    const state = ref<ChatState>({
        messages: [],
        isStreaming: false,
        isAnalyzing: false,
        tokenBudget: {
            used: 0,
            limit: DEFAULT_TOKEN_LIMIT,
            inputTokens: 0,
            outputTokens: 0
        },
        currentStreamingMessageId: null
    })

    const inputMessage = ref('')
    const streamSubscription = shallowRef<StreamSubscription | null>(null)

    // Computed
    const messages = computed(() => state.value.messages)
    const isStreaming = computed(() => state.value.isStreaming)
    const isAnalyzing = computed(() => state.value.isAnalyzing)
    const tokenBudget = computed(() => state.value.tokenBudget)
    const hasMessages = computed(() => state.value.messages.length > 0)

    const tokenUsagePercent = computed(() => {
        const { used, limit } = state.value.tokenBudget
        return limit > 0 ? (used / limit) * 100 : 0
    })

    const isNearTokenLimit = computed(() => {
        return tokenUsagePercent.value >= TOKEN_WARNING_THRESHOLD * 100
    })

    const canSend = computed(() => {
        return inputMessage.value.trim().length > 0 && !state.value.isStreaming
    })

    // Real-time token counting for input
    const inputTokenCount = computed(() => {
        return chatApi.estimateTokens(inputMessage.value)
    })

    const estimatedTotalTokens = computed(() => {
        const contextTokens = contextStore.tokenCount || 0
        return inputTokenCount.value + contextTokens + state.value.tokenBudget.used
    })

    // Actions

    /**
     * Generate unique message ID
     */
    function generateMessageId(): string {
        return `msg-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`
    }

    /**
     * Add a message to the chat
     */
    function addMessage(message: Omit<ChatMessage, 'id' | 'timestamp'>): ChatMessage {
        const newMessage: ChatMessage = {
            ...message,
            id: generateMessageId(),
            timestamp: new Date().toISOString()
        }
        state.value.messages.push(newMessage)
        return newMessage
    }

    /**
     * Update a message by ID
     */
    function updateMessage(id: string, updates: Partial<ChatMessage>): void {
        const index = state.value.messages.findIndex(m => m.id === id)
        if (index !== -1) {
            state.value.messages[index] = {
                ...state.value.messages[index],
                ...updates
            }
        }
    }

    /**
     * Update token budget
     */
    function updateTokenBudget(tokens: number, isOutput = false): void {
        if (isOutput) {
            state.value.tokenBudget.outputTokens += tokens
        } else {
            state.value.tokenBudget.inputTokens = tokens
        }
        state.value.tokenBudget.used =
            state.value.tokenBudget.inputTokens + state.value.tokenBudget.outputTokens
    }

    /**
     * Parse diff from AI response for inline preview
     */
    function parseDiffFromResponse(content: string): DiffPreview | undefined {
        const diffMatch = content.match(/```diff\n([\s\S]*?)```/)
        if (!diffMatch) return undefined

        const diffContent = diffMatch[1]
        const lines = diffContent.split('\n')

        let filePath = ''
        const hunks: DiffHunk[] = []
        let currentHunk: DiffHunk | null = null
        let additions = 0
        let deletions = 0

        for (const line of lines) {
            if (line.startsWith('--- ') || line.startsWith('+++ ')) {
                const pathMatch = line.match(/^[+-]{3}\s+(?:a\/|b\/)?(.+)/)
                if (pathMatch) filePath = pathMatch[1]
                continue
            }

            if (line.startsWith('@@')) {
                const hunkMatch = line.match(/@@ -(\d+),?(\d*) \+(\d+),?(\d*) @@/)
                if (hunkMatch) {
                    if (currentHunk) hunks.push(currentHunk)
                    currentHunk = {
                        oldStart: parseInt(hunkMatch[1], 10),
                        oldLines: parseInt(hunkMatch[2] || '1', 10),
                        newStart: parseInt(hunkMatch[3], 10),
                        newLines: parseInt(hunkMatch[4] || '1', 10),
                        lines: []
                    }
                }
                continue
            }

            if (currentHunk) {
                if (line.startsWith('+') && !line.startsWith('+++')) {
                    currentHunk.lines.push({ type: 'add', content: line.slice(1) })
                    additions++
                } else if (line.startsWith('-') && !line.startsWith('---')) {
                    currentHunk.lines.push({ type: 'remove', content: line.slice(1) })
                    deletions++
                } else if (line.startsWith(' ') || line === '') {
                    currentHunk.lines.push({ type: 'context', content: line.slice(1) || '' })
                }
            }
        }

        if (currentHunk) hunks.push(currentHunk)

        if (hunks.length === 0) return undefined

        return { filePath, hunks, additions, deletions }
    }

    /**
     * Handle streaming chunk
     */
    function handleStreamChunk(chunk: StreamChunk): void {
        const messageId = state.value.currentStreamingMessageId
        if (!messageId) return

        const message = state.value.messages.find(m => m.id === messageId)
        if (!message) return

        if (chunk.content) {
            const newContent = message.content + chunk.content
            const diffPreview = parseDiffFromResponse(newContent)

            updateMessage(messageId, {
                content: newContent,
                diffPreview
            })
        }

        if (chunk.tokensUsed) {
            updateTokenBudget(chunk.tokensUsed, true)
            updateMessage(messageId, { tokenCount: chunk.tokensUsed })
        }
    }

    /**
     * Handle tool call during streaming
     */
    function handleToolCall(toolCall: StreamToolCall): void {
        const messageId = state.value.currentStreamingMessageId
        if (!messageId) return

        const message = state.value.messages.find(m => m.id === messageId)
        if (!message) return

        const existingCalls = message.toolCalls || []
        const existingIndex = existingCalls.findIndex(tc => tc.id === toolCall.id)

        const chatToolCall: ChatToolCall = {
            id: toolCall.id,
            tool: toolCall.tool,
            arguments: toolCall.arguments,
            status: toolCall.status,
            result: toolCall.result,
            error: toolCall.error,
            duration: toolCall.duration
        }

        if (existingIndex !== -1) {
            existingCalls[existingIndex] = chatToolCall
        } else {
            existingCalls.push(chatToolCall)
        }

        updateMessage(messageId, { toolCalls: [...existingCalls] })
    }

    /**
     * Send message with streaming response
     */
    async function sendStreamingMessage(scrollToBottom?: () => void): Promise<void> {
        if (!canSend.value) return

        let content = inputMessage.value.trim()
        const projectRoot = projectStore.currentPath || ''

        // Process mentions
        if (content.includes('@')) {
            const { cleanedMessage } = await mentions.processMessageMentions(content, t)
            content = cleanedMessage
            if (!content) {
                inputMessage.value = ''
                return
            }
        }

        // Add user message
        addMessage({
            role: 'user',
            content,
            tokenCount: chatApi.estimateTokens(content)
        })

        inputMessage.value = ''
        scrollToBottom?.()

        // Update input token count
        const contextTokens = contextStore.tokenCount || 0
        updateTokenBudget(chatApi.estimateTokens(content) + contextTokens)

        // Add placeholder for assistant response
        const assistantMessage = addMessage({
            role: 'assistant',
            content: '',
            isStreaming: true
        })

        state.value.isStreaming = true
        state.value.currentStreamingMessageId = assistantMessage.id

        // Collect smart context if needed
        let smartContext = undefined
        if (contextStore.hasContext) {
            try {
                const selectedFiles = fileStore.selectedFilesList
                smartContext = await chatApi.collectSmartContext(
                    content,
                    projectRoot,
                    selectedFiles.length > 0 ? selectedFiles : undefined
                )
            } catch {
                // Continue without smart context
            }
        }

        // Start streaming
        streamSubscription.value = chatApi.startStreamingChat(
            { task: content, projectRoot, smartContext },
            {
                onChunk: handleStreamChunk,
                onToolCall: handleToolCall,
                onTokenUpdate: (tokens) => updateTokenBudget(tokens, true),
                onComplete: () => {
                    state.value.isStreaming = false
                    state.value.currentStreamingMessageId = null
                    updateMessage(assistantMessage.id, { isStreaming: false })
                    scrollToBottom?.()
                },
                onError: (error) => {
                    state.value.isStreaming = false
                    state.value.currentStreamingMessageId = null
                    updateMessage(assistantMessage.id, {
                        isStreaming: false,
                        error,
                        content: error || t('chat.error')
                    })
                    uiStore.addToast(error || t('chat.error'), 'error')
                }
            }
        )

        scrollToBottom?.()
    }

    /**
     * Send message without streaming (fallback)
     */
    async function sendMessage(scrollToBottom?: () => void): Promise<void> {
        if (!canSend.value) return

        let content = inputMessage.value.trim()
        const projectRoot = projectStore.currentPath || ''

        // Process mentions
        if (content.includes('@')) {
            const { cleanedMessage } = await mentions.processMessageMentions(content, t)
            content = cleanedMessage
            if (!content) {
                inputMessage.value = ''
                return
            }
        }

        // Add user message
        addMessage({
            role: 'user',
            content,
            tokenCount: chatApi.estimateTokens(content)
        })

        inputMessage.value = ''
        scrollToBottom?.()

        state.value.isStreaming = true

        try {
            const selectedFiles = fileStore.selectedFilesList
            const smartContext = contextStore.hasContext
                ? await chatApi.collectSmartContext(
                    content,
                    projectRoot,
                    selectedFiles.length > 0 ? selectedFiles : undefined
                )
                : undefined

            const response = await chatApi.sendAgenticMessage(content, projectRoot, smartContext)

            const diffPreview = parseDiffFromResponse(response.response)

            addMessage({
                role: 'assistant',
                content: response.response,
                toolCalls: response.toolCalls?.map(tc => ({
                    id: `tc-${Date.now()}-${Math.random().toString(36).slice(2, 5)}`,
                    tool: tc.tool,
                    arguments: tc.arguments,
                    status: 'completed' as const,
                    result: tc.result
                })),
                diffPreview
            })

            scrollToBottom?.()
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : String(error)
            addMessage({
                role: 'assistant',
                content: t('chat.error'),
                error: errorMessage
            })
            uiStore.addToast(t('chat.error'), 'error')
        } finally {
            state.value.isStreaming = false
        }
    }

    /**
     * Stop current streaming
     */
    function stopStreaming(): void {
        if (streamSubscription.value) {
            streamSubscription.value.unsubscribe()
            streamSubscription.value = null
        }

        if (state.value.currentStreamingMessageId) {
            updateMessage(state.value.currentStreamingMessageId, { isStreaming: false })
        }

        state.value.isStreaming = false
        state.value.currentStreamingMessageId = null
    }

    /**
     * Clear all messages
     */
    function clearChat(): void {
        stopStreaming()
        state.value.messages = []
        state.value.tokenBudget = {
            used: 0,
            limit: DEFAULT_TOKEN_LIMIT,
            inputTokens: 0,
            outputTokens: 0
        }
    }

    /**
     * Delete a specific message
     */
    function deleteMessage(messageId: string): void {
        state.value.messages = state.value.messages.filter(m => m.id !== messageId)
    }

    /**
     * Copy message content to clipboard
     */
    async function copyMessage(content: string): Promise<void> {
        try {
            await navigator.clipboard.writeText(content)
            uiStore.addToast(t('chat.copied'), 'success')
        } catch {
            uiStore.addToast(t('chat.copyFailed'), 'error')
        }
    }

    /**
     * Set token limit
     */
    function setTokenLimit(limit: number): void {
        state.value.tokenBudget.limit = limit
    }

    /**
     * Get token count info for display
     */
    function getTokenCountInfo(): TokenCountInfo {
        return chatApi.calculateTokenCount(
            inputMessage.value,
            contextStore.tokenCount || 0
        )
    }

    // Cleanup on unmount
    onUnmounted(() => {
        stopStreaming()
    })

    return {
        // State
        messages,
        inputMessage,
        isStreaming,
        isAnalyzing,
        tokenBudget,
        hasMessages,

        // Computed
        tokenUsagePercent,
        isNearTokenLimit,
        canSend,
        inputTokenCount,
        estimatedTotalTokens,

        // Actions
        sendStreamingMessage,
        sendMessage,
        stopStreaming,
        clearChat,
        deleteMessage,
        copyMessage,
        setTokenLimit,
        getTokenCountInfo,
        updateMessage
    }
}
