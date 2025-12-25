import type { domain } from '#wailsjs/go/models'
import { apiService } from '@/services/api.service'

export interface ProjectContext {
    id: string
    name: string
    files: string[]
    fileCount?: number
    totalSize?: number
    lineCount?: number
    tokenCount?: number
    createdAt: string
    updatedAt?: string
    metadata?: {
        selectedFiles?: string[]
        warnings?: string[]
        skippedFiles?: string[]
    }
}

export class ContextApi {
    async buildContext(
        projectPath: string,
        files: string[],
        options: domain.ContextBuildOptions
    ): Promise<domain.ContextSummary> {
        try {
            return await apiService.buildContextFromRequest(projectPath, files, options)
        } catch (error) {
            console.error('[ContextApi] Failed to build context:', error)

            // Re-throw token limit errors with info preserved
            if (error instanceof Error && error.message === 'TOKEN_LIMIT_EXCEEDED') {
                throw error
            }

            throw new Error('Failed to build context. Please check your file selection.')
        }
    }

    async getContextContent(contextId: string): Promise<string> {
        try {
            return await apiService.getFullContextContent(contextId)
        } catch (error) {
            console.error('[ContextApi] Failed to get context content:', error)
            throw new Error('Failed to load context content.')
        }
    }

    async deleteContext(contextId: string): Promise<void> {
        try {
            await apiService.deleteContext(contextId)
        } catch (error) {
            console.error('[ContextApi] Failed to delete context:', error)
            throw new Error('Failed to delete context.')
        }
    }

    async exportContext(exportSettings: { format: string; includeLineNumbers?: boolean; includeMetadata?: boolean; maxTokens?: number }): Promise<domain.ExportResult> {
        try {
            return await apiService.exportContext(exportSettings)
        } catch (error) {
            console.error('[ContextApi] Failed to export context:', error)
            throw new Error('Failed to export context.')
        }
    }

    async getProjectContexts(projectPath: string): Promise<ProjectContext[]> {
        try {
            const result = await apiService.getProjectContexts(projectPath)
            const parsed = JSON.parse(result) as ProjectContext[] | null
            if (!Array.isArray(parsed)) return []

            // Map backend structure to frontend: files come from metadata.selectedFiles
            return parsed.map(ctx => ({
                ...ctx,
                files: ctx.metadata?.selectedFiles || ctx.files || []
            }))
        } catch (error) {
            console.error('[ContextApi] Failed to get project contexts:', error)
            throw new Error('Failed to load project contexts.')
        }
    }
}

export const contextApi = new ContextApi()
