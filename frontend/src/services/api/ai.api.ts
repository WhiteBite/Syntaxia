/**
 * AI and code generation API
 * Handles AI models, code generation, Qwen tasks
 */

import * as wails from '#wailsjs/go/main/App'
import type {
    QwenContextPreview,
    QwenModelInfo,
    QwenTaskRequest,
    QwenTaskResponse,
} from '../types'
import { apiCall, parseJsonResponse } from './base'

export const aiApi = {
    generateCode: (context: string, task: string): Promise<string> =>
        apiCall(() => wails.GenerateCode(context, task), 'Failed to generate code.', { logContext: 'ai' }),

    generateCodeStream: (context: string, task: string): void => {
        try {
            wails.GenerateCodeStream(context, task)
        } catch (error) {
            console.error('[API:ai] Error starting code stream:', error)
        }
    },

    generateIntelligentCode: (context: string, task: string, options: string): Promise<string> =>
        apiCall(
            () => wails.GenerateIntelligentCode(context, task, options),
            'Failed to generate code.',
            { logContext: 'ai' }
        ),

    listAvailableModels: (): Promise<string[]> =>
        apiCall(() => wails.ListAvailableModels(), 'Failed to list AI models.', { logContext: 'ai' }),

    getProviderInfo: (): Promise<string> =>
        apiCall(() => wails.GetProviderInfo(), 'Failed to get provider information.', { logContext: 'ai' }),

    // Qwen Task Execution
    qwenExecuteTask: async (request: QwenTaskRequest): Promise<QwenTaskResponse> => {
        const result = await apiCall(
            () => wails.QwenExecuteTask(JSON.stringify(request)),
            'Failed to execute task with Qwen.',
            { logContext: 'ai' }
        )
        return parseJsonResponse(result, 'Failed to parse Qwen response.')
    },

    qwenPreviewContext: async (request: QwenTaskRequest): Promise<QwenContextPreview> => {
        const result = await apiCall(
            () => wails.QwenPreviewContext(JSON.stringify(request)),
            'Failed to preview context.',
            { logContext: 'ai' }
        )
        return parseJsonResponse(result, 'Failed to parse context preview.')
    },

    qwenGetAvailableModels: async (): Promise<QwenModelInfo[]> => {
        const result = await apiCall(
            () => wails.QwenGetAvailableModels(),
            'Failed to get Qwen models.',
            { logContext: 'ai' }
        )
        return parseJsonResponse(result, 'Failed to parse Qwen models.')
    },
}
