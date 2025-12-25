/**
 * Context building API
 * Handles context creation, export, suggestions
 */

import * as wails from '#wailsjs/go/main/App'
import type { domain } from '#wailsjs/go/models'
import type {
    AgenticChatResponse,
    FileQuickInfo,
    ImpactPreviewResult,
    SmartSuggestionsResult,
} from '../types'
import { apiCall, parseJsonResponse } from './base'

export const contextApi = {
    buildContext: (projectPath: string, files: string[], task: string): Promise<string> =>
        apiCall(() => wails.BuildContext(projectPath, files, task), 'Failed to build context.', { logContext: 'context' }),

    buildContextFromRequest: async (
        projectPath: string,
        files: string[],
        options: domain.ContextBuildOptions
    ): Promise<domain.ContextSummary> => {
        try {
            return await wails.BuildContextFromRequest(projectPath, files, options)
        } catch (error) {
            console.error('[API:context] Error building context from request:', error)

            let errorMsg = ''
            if (error && typeof error === 'object') {
                const err = error as { cause?: string; message?: string }
                errorMsg = err.cause || err.message || String(error)
            } else if (error instanceof Error) {
                errorMsg = error.message
            } else {
                errorMsg = String(error)
            }

            if (errorMsg.includes('token limit')) {
                const match = errorMsg.match(/(\d+)\s*[>\u003e]\s*(\d+)/)
                if (match) {
                    const actual = Number(match[1])
                    const limit = Number(match[2])
                    const tokenError = new Error('TOKEN_LIMIT_EXCEEDED') as Error & {
                        tokenInfo: { actual: number; limit: number }
                    }
                    tokenError.tokenInfo = { actual, limit }
                    throw tokenError
                }
                throw new Error('TOKEN_LIMIT_EXCEEDED')
            }

            throw new Error('Failed to build context.')
        }
    },

    getContext: (contextId: string): Promise<string> =>
        apiCall(() => wails.GetContext(contextId), 'Failed to get context.', { logContext: 'context' }),

    deleteContext: (contextId: string): Promise<void> =>
        apiCall(() => wails.DeleteContext(contextId), 'Failed to delete context.', { logContext: 'context' }),

    getProjectContexts: (projectPath: string): Promise<string> =>
        apiCall(() => wails.GetProjectContexts(projectPath), 'Failed to get project contexts.', { logContext: 'context' }),

    exportContext: (exportSettings: unknown): Promise<domain.ExportResult> =>
        apiCall(
            () => wails.ExportContext(JSON.stringify(exportSettings)),
            'Failed to export context.',
            { logContext: 'context' }
        ),

    getFullContextContent: (contextId: string): Promise<string> =>
        apiCall(
            () => wails.GetFullContextContent(contextId),
            'Failed to get full context content.',
            { logContext: 'context' }
        ),

    suggestContextFiles: (taskDescription: string, files: domain.FileNode[]): Promise<string[]> =>
        apiCall(
            () => wails.SuggestContextFiles(taskDescription, files),
            'Failed to suggest context files.',
            { logContext: 'context' }
        ),

    getSmartSuggestions: async (
        projectPath: string,
        currentFiles: string[],
        task = ''
    ): Promise<SmartSuggestionsResult> => {
        try {
            const result = await wails.GetSmartSuggestions(projectPath, currentFiles, task)
            return {
                suggestions:
                    result.suggestions?.map((s) => ({
                        path: s.path,
                        source: s.source as 'git' | 'arch' | 'semantic',
                        reason: s.reason,
                        confidence: s.confidence,
                    })) || [],
                total: result.total || 0,
            }
        } catch (error) {
            console.error('[API:context] Error getting smart suggestions:', error)
            return { suggestions: [], total: 0 }
        }
    },

    getFileQuickInfo: async (projectPath: string, filePath: string): Promise<FileQuickInfo> => {
        try {
            // @ts-ignore - method may not exist in wails bindings yet
            const result = await wails.GetFileQuickInfo(projectPath, filePath)
            return {
                symbolCount: result.symbolCount,
                importCount: result.importCount,
                dependentCount: result.dependentCount,
                changeRisk: result.changeRisk,
                riskLevel: result.riskLevel as 'high' | 'medium' | 'low',
            }
        } catch {
            return { symbolCount: 0, importCount: 0, dependentCount: 0, changeRisk: 0, riskLevel: 'low' }
        }
    },

    getImpactPreview: async (projectPath: string, filePaths: string[]): Promise<ImpactPreviewResult> => {
        try {
            // @ts-ignore - method may not exist in wails bindings yet
            const result = await wails.GetImpactPreview(projectPath, filePaths)
            return {
                totalDependents: result.totalDependents,
                aggregateRisk: result.aggregateRisk,
                riskLevel: result.riskLevel as 'high' | 'medium' | 'low',
                affectedFiles: (result.affectedFiles || []).map((f: { path: string; type: string; dependents: number }) => ({
                    path: f.path,
                    type: f.type as 'direct' | 'transitive',
                    dependents: f.dependents,
                })),
                relatedTests: result.relatedTests || [],
            }
        } catch {
            return { totalDependents: 0, aggregateRisk: 0, riskLevel: 'low', affectedFiles: [], relatedTests: [] }
        }
    },

    analyzeTaskAndCollectContext: (task: string, allFilesJson: string, rootDir: string): Promise<string> =>
        apiCall(
            () => wails.AnalyzeTaskAndCollectContext(task, allFilesJson, rootDir),
            'Failed to analyze task and collect context.',
            { logContext: 'context' }
        ),

    agenticChat: async (task: string, projectRoot: string): Promise<AgenticChatResponse> => {
        const request = { task, projectRoot }
        const result = await apiCall(
            () => wails.AgenticChat(JSON.stringify(request)),
            'Failed to execute agentic chat.',
            { logContext: 'context' }
        )
        return parseJsonResponse(result, 'Failed to parse agentic chat response.')
    },
}
