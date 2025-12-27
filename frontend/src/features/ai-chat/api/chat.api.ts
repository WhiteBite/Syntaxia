/**
 * Chat API module for AI Chat feature
 * Handles streaming responses, tool calls, and token counting
 */

import * as wails from '#wailsjs/go/main/App'
import { EventsOff, EventsOn } from '#wailsjs/runtime/runtime'
import { useLogger } from '@/composables/useLogger'
import { apiCall, parseJsonResponse } from '@/services/api/base'
import type { AgenticChatResponse, SmartContextResult, ToolCallLog } from '@/services/types'

const logger = useLogger('API:chat')

/** Stream chunk from backend */
export interface StreamChunk {
    content: string
    done: boolean
    error?: string
    tokensUsed?: number
    finishReason?: string
}

/** Tool call event during streaming */
export interface StreamToolCall {
    id: string
    tool: string
    arguments: string
    status: 'pending' | 'executing' | 'completed' | 'failed'
    result?: string
    error?: string
    duration?: number
}

/** Streaming chat request */
export interface StreamChatRequest {
    task: string
    projectRoot: string
    smartContext?: SmartContextResult
    maxTokens?: number
}

/** Streaming chat callbacks */
export interface StreamChatCallbacks {
    onChunk: (chunk: StreamChunk) => void
    onToolCall?: (toolCall: StreamToolCall) => void
    onTokenUpdate?: (tokens: number) => void
    onComplete?: (response: AgenticChatResponse) => void
    onError?: (error: string) => void
}

/** Stream subscription handle */
export interface StreamSubscription {
    unsubscribe: () => void
}

/** Token count info */
export interface TokenCountInfo {
    inputTokens: number
    outputTokens: number
    totalTokens: number
    estimatedCost?: number
}

/**
 * Subscribe to AI stream events
 */
function subscribeToStream(callbacks: StreamChatCallbacks): StreamSubscription {
    const eventName = 'ai:stream:chunk'
    let totalTokens = 0
    let accumulatedContent = ''
    const toolCalls: ToolCallLog[] = []

    const unsubscribe = EventsOn(eventName, (chunk: StreamChunk) => {
        if (chunk.error) {
            callbacks.onError?.(chunk.error)
            return
        }

        if (chunk.content) {
            accumulatedContent += chunk.content
            callbacks.onChunk(chunk)
        }

        if (chunk.tokensUsed) {
            totalTokens = chunk.tokensUsed
            callbacks.onTokenUpdate?.(totalTokens)
        }

        if (chunk.done) {
            callbacks.onComplete?.({
                response: accumulatedContent,
                toolCalls,
                iterations: 1,
                context: []
            })
        }
    })

    return {
        unsubscribe: () => {
            unsubscribe()
            EventsOff(eventName)
        }
    }
}

/**
 * Subscribe to tool call events during streaming
 */
function subscribeToToolCalls(callback: (toolCall: StreamToolCall) => void): StreamSubscription {
    const eventName = 'ai:tool:call'

    const unsubscribe = EventsOn(eventName, (toolCall: StreamToolCall) => {
        callback(toolCall)
    })

    return {
        unsubscribe: () => {
            unsubscribe()
            EventsOff(eventName)
        }
    }
}

export const chatApi = {
    /**
     * Start streaming chat with AI
     */
    startStreamingChat: (
        request: StreamChatRequest,
        callbacks: StreamChatCallbacks
    ): StreamSubscription => {
        const streamSub = subscribeToStream(callbacks)
        const toolSub = callbacks.onToolCall
            ? subscribeToToolCalls(callbacks.onToolCall)
            : null

        // Start the stream via Wails
        try {
            wails.GenerateCodeStream(
                JSON.stringify(request.smartContext || {}),
                request.task
            )
        } catch (error) {
            logger.error('Error starting stream:', error)
            callbacks.onError?.(error instanceof Error ? error.message : String(error))
        }

        return {
            unsubscribe: () => {
                streamSub.unsubscribe()
                toolSub?.unsubscribe()
            }
        }
    },

    /**
     * Send agentic chat message (non-streaming)
     */
    sendAgenticMessage: async (
        task: string,
        projectRoot: string,
        smartContext?: SmartContextResult
    ): Promise<AgenticChatResponse> => {
        const request = { task, projectRoot, smartContext }
        const result = await apiCall(
            () => wails.AgenticChat(JSON.stringify(request)),
            'Failed to execute agentic chat.',
            { logContext: 'chat' }
        )
        return parseJsonResponse(result, 'Failed to parse agentic chat response.')
    },

    /**
     * Estimate tokens for a message
     */
    estimateTokens: (text: string): number => {
        // Rough estimation: ~4 characters per token for English
        // More accurate for code: ~3.5 characters per token
        const AVG_CHARS_PER_TOKEN = 3.5
        return Math.ceil(text.length / AVG_CHARS_PER_TOKEN)
    },

    /**
     * Calculate token count for context + message
     */
    calculateTokenCount: (
        message: string,
        contextTokens: number
    ): TokenCountInfo => {
        const messageTokens = chatApi.estimateTokens(message)
        const totalTokens = messageTokens + contextTokens

        return {
            inputTokens: totalTokens,
            outputTokens: 0,
            totalTokens,
            estimatedCost: undefined
        }
    },

    /**
     * Collect smart context for chat
     */
    collectSmartContext: async (
        task: string,
        projectRoot: string,
        selectedFiles?: string[],
        maxTokens?: number
    ): Promise<SmartContextResult> => {
        const request = { task, projectRoot, selectedFiles, maxTokens }
        const result = await apiCall(
            () => wails.CollectSmartContext(JSON.stringify(request)),
            'Failed to collect smart context.',
            { logContext: 'chat' }
        )
        return parseJsonResponse(result, 'Failed to parse smart context.')
    }
}
