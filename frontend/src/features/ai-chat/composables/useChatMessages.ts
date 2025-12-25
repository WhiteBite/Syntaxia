import { useContextStore } from '@/features/context'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useSettingsStore } from '@/stores/settings.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, ref } from 'vue'
import { useMentions } from './useMentions'

export interface Message {
    id: string
    role: 'user' | 'assistant'
    content: string
    timestamp: string
    toolCalls?: ToolCallLog[]
    contextAttached?: boolean
    tokenCount?: number
    iterations?: number
    error?: string
    /** Thinking process log for AI messages */
    thinkingLog?: ThinkingStep[]
    /** Suggested files from AI analysis */
    suggestedFiles?: SuggestedFile[]
}

export interface ToolCallLog {
    tool: string
    args?: Record<string, unknown>
    arguments?: string
    result?: string
    error?: string
    duration?: number
}

export interface ThinkingStep {
    action: string
    detail?: string
    status: 'pending' | 'done' | 'error'
    timestamp: number
}

export interface SuggestedFile {
    path: string
    reason: string
    relevance: number
}

export interface SmartContextPreview {
    files: SuggestedFile[]
    totalTokens: number
    query: string
}

type TranslateFunc = (key: string, params?: Record<string, string | number>) => string

export function useChatMessages() {
    const uiStore = useUIStore()
    const projectStore = useProjectStore()
    const settingsStore = useSettingsStore()
    const contextStore = useContextStore()
    const mentions = useMentions()

    // State
    const messages = ref<Message[]>([])
    const inputMessage = ref('')
    const isThinking = ref(false)
    const isAnalyzing = ref(false)
    const smartContextPreview = ref<SmartContextPreview | null>(null)
    const expandedToolCalls = ref<Set<number>>(new Set())
    const abortController = ref<AbortController | null>(null)

    // Computed
    const currentModel = computed(() => settingsStore.settings.aiModel || 'gpt-4')
    const providerName = computed(() => 'OpenAI')
    const isConnected = computed(() => true)
    const totalUsedTokens = computed(() => {
        return messages.value.reduce((acc, msg) => acc + (msg.content?.length || 0) / 4, 0)
    })

    // Actions

    /**
     * Unified smart message flow:
     * 1. Process any @mentions first (build context from files/git)
     * 2. If context exists - use it directly
     * 3. If no context - analyze and suggest files first
     * 4. User confirms -> AI generates response
     */
    async function sendSmartMessage(scrollToBottom: () => void, t: TranslateFunc) {
        if (!inputMessage.value.trim() || isThinking.value) return

        let content = inputMessage.value.trim()

        // Process mentions first (@files, @git, @problems)
        if (content.includes('@')) {
            const { cleanedMessage } = await mentions.processMessageMentions(content, t)
            content = cleanedMessage

            // If message was only mentions, don't send empty message
            if (!content) {
                inputMessage.value = ''
                return
            }
        }

        const hasExistingContext = contextStore.hasContext

        // If no context, first analyze and suggest files
        if (!hasExistingContext) {
            inputMessage.value = content
            await analyzeAndSuggestFiles(content, t)
            return
        }

        // Has context - send message directly
        inputMessage.value = content
        await executeChat(content, scrollToBottom)
    }

    /**
     * Analyze query and suggest relevant files
     */
    async function analyzeAndSuggestFiles(query: string, t: TranslateFunc) {
        isAnalyzing.value = true
        try {
            const result = await apiService.semanticSearch({
                query,
                projectRoot: projectStore.currentPath || '',
                topK: 10
            })

            if (result.results && result.results.length > 0) {
                const files: SuggestedFile[] = result.results.map((r) => ({
                    path: r.chunk.filePath,
                    reason: r.chunk.content?.slice(0, 100) || '',
                    relevance: r.score || 0.5
                }))

                const estimatedTokens = files.length * 500

                smartContextPreview.value = {
                    files,
                    totalTokens: estimatedTokens,
                    query,
                }
            } else {
                uiStore.addToast(t('chat.noFilesFound'), 'warning')
                // Still allow sending without context
                inputMessage.value = query
            }
        } catch {
            uiStore.addToast(t('chat.analysisFailed'), 'error')
        } finally {
            isAnalyzing.value = false
        }
    }

    /**
     * Execute chat with current context
     */
    async function executeChat(content: string, scrollToBottom: () => void) {
        inputMessage.value = ''

        const userMessage: Message = {
            id: `msg-${Date.now()}`,
            role: 'user',
            content,
            timestamp: new Date().toISOString(),
            contextAttached: contextStore.hasContext,
        }
        messages.value.push(userMessage)
        scrollToBottom()

        isThinking.value = true
        abortController.value = new AbortController()

        try {
            const projectRoot = projectStore.currentPath || ''
            const response = await apiService.agenticChat(content, projectRoot)

            const assistantMessage: Message = {
                id: `msg-${Date.now()}-ai`,
                role: 'assistant',
                content: response.response,
                timestamp: new Date().toISOString(),
                toolCalls: response.toolCalls,
            }
            messages.value.push(assistantMessage)
            scrollToBottom()
        } catch (error) {
            if ((error as Error).name !== 'AbortError') {
                messages.value.push({
                    id: `msg-${Date.now()}-error`,
                    role: 'assistant',
                    content: 'Error: Failed to get response',
                    timestamp: new Date().toISOString(),
                    error: (error as Error).message,
                })
            }
        } finally {
            isThinking.value = false
            abortController.value = null
        }
    }

    async function confirmSmartContext(selectedFiles: string[], t: TranslateFunc, scrollToBottom: () => void) {
        if (!smartContextPreview.value) return

        const query = smartContextPreview.value.query
        smartContextPreview.value = null

        // Build context from user-selected files
        try {
            await contextStore.buildContext(selectedFiles)

            // Now send the message with context
            await executeChat(query, scrollToBottom)
        } catch {
            uiStore.addToast(t('chat.contextBuildFailed'), 'error')
        }
    }

    function cancelSmartContext() {
        smartContextPreview.value = null
    }

    function clearChat() {
        messages.value = []
        expandedToolCalls.value.clear()
    }

    function toggleToolCalls(index: number) {
        if (expandedToolCalls.value.has(index)) {
            expandedToolCalls.value.delete(index)
        } else {
            expandedToolCalls.value.add(index)
        }
        expandedToolCalls.value = new Set(expandedToolCalls.value)
    }

    function stopGeneration() {
        if (abortController.value) {
            abortController.value.abort()
            abortController.value = null
        }
        isThinking.value = false
    }

    function copyMessage(content: string, t: TranslateFunc) {
        navigator.clipboard.writeText(content)
        uiStore.addToast(t('common.copied'), 'success')
    }

    function initialize() {
        // No mode to restore - unified smart flow
    }

    function cleanup() {
        stopGeneration()
    }

    return {
        // State
        messages,
        inputMessage,
        isThinking,
        isAnalyzing,
        smartContextPreview,
        expandedToolCalls,
        // Computed
        currentModel,
        providerName,
        isConnected,
        totalUsedTokens,
        // Actions
        sendSmartMessage,
        confirmSmartContext,
        cancelSmartContext,
        clearChat,
        toggleToolCalls,
        stopGeneration,
        copyMessage,
        initialize,
        cleanup,
    }
}
