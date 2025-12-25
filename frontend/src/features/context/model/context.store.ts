import type { domain } from '#wailsjs/go/models'
import { EventsOff, EventsOn } from '#wailsjs/runtime/runtime'
import { useLogger } from '@/composables/useLogger'
import { useProjectStore } from '@/stores/project.store'
import { useSettingsStore } from '@/stores/settings.store'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { contextApi } from '../api/context.api'
import { performAutoCleanup } from '../composables/useContextCleanup'
import {
    MAX_CONTEXT_SIZE,
    useContextContent,
    type ContextChunk
} from '../composables/useContextContent'
import { loadContextMetadata, saveContextMetadata } from '../composables/useContextMetadata'
import {
    createContextSummary,
    duplicateContextContent,
    mergeContextsContent,
    reorderContextList
} from '../composables/useContextOperations'
import { clearHighlightCache } from '../composables/useSyntaxHighlight'
import { generateSmartName } from '../lib/context-naming'

const logger = useLogger('ContextStore')

export interface ContextSummary {
    id: string
    name: string
    fileCount: number
    totalSize: number
    lineCount: number
    tokenCount?: number
    createdAt?: string
    projectPath?: string
    isFavorite?: boolean
    files?: string[]
    /** Original user-selected files (before expansion) */
    selectedFiles?: string[]
    metadata?: {
        warnings?: string[]
        skippedFiles?: string[]
        skippedReasons?: Record<string, string>
        projectPath?: string
    }
}

export { MAX_CONTEXT_SIZE, type ContextChunk }

export const useContextStore = defineStore('context', () => {
    // Content management composable
    const contentManager = useContextContent()

    // State
    const contextId = ref<string | null>(null)
    const summary = ref<ContextSummary | null>(null)
    const currentChunk = ref<ContextChunk | null>(null)
    const isLoading = ref(false)
    const isBuilding = ref(false)
    const buildProgress = ref(0)
    const error = ref<string | null>(null)
    const contextList = ref<ContextSummary[]>([])
    const warnings = ref<string[]>([])
    const skippedFiles = ref<string[]>([])
    // Selected item from context list (for right panel preview)
    const selectedListItem = ref<ContextSummary | null>(null)

    // Computed
    const hasContext = computed(() => contextId.value !== null)
    const totalSize = computed(() => summary.value?.totalSize || 0)
    const fileCount = computed(() => summary.value?.fileCount || 0)
    const lineCount = computed(() => summary.value?.lineCount || 0)
    const tokenCount = computed(() => summary.value?.tokenCount || Math.round(lineCount.value * 2.5))
    const totalTokens = computed(() => tokenCount.value) // Alias for ContextIndicator
    const estimatedCost = computed(() => (tokenCount.value / 1000) * 0.002)

    // Stats object for components
    const stats = computed(() => ({
        tokens: tokenCount.value,
        totalFiles: fileCount.value,
        lines: lineCount.value,
        size: totalSize.value
    }))

    // Actions
    async function buildContext(filePaths: string[], options?: Partial<domain.ContextBuildOptions>) {
        if (filePaths.length === 0) {
            error.value = 'No files selected'
            return
        }

        isBuilding.value = true
        isLoading.value = true
        buildProgress.value = 5
        error.value = null

        // Subscribe to real progress events from backend
        let unsubscribeProgress: (() => void) | null = null
        const totalFiles = filePaths.length

        const setupProgressListener = () => {
            unsubscribeProgress = EventsOn('SyntaxiaContextGenerationProgress', (data: { current: number; total: number }) => {
                // Calculate real progress (5% to 85% range for file reading phase)
                const fileProgress = data.total > 0 ? (data.current / data.total) : 0
                buildProgress.value = Math.round(5 + fileProgress * 80)
            })
        }

        const cleanupProgressListener = () => {
            if (unsubscribeProgress) {
                unsubscribeProgress()
                unsubscribeProgress = null
            }
            EventsOff('SyntaxiaContextGenerationProgress')
        }

        try {
            const projectStore = useProjectStore()
            if (!projectStore.currentPath) {
                error.value = 'No project selected'
                throw new Error('No project selected')
            }

            const buildOptions = createBuildOptions(options)
            logger.debug('Building context with options:', {
                outputFormat: buildOptions.outputFormat,
                stripComments: buildOptions.stripComments,
                maxTokens: buildOptions.maxTokens,
                fileCount: totalFiles
            })

            setupProgressListener()
            const result = await contextApi.buildContext(projectStore.currentPath, filePaths, buildOptions)

            // Phase 2: Processing result (85% -> 90%)
            buildProgress.value = 90
            contextId.value = result.id || `ctx-${Date.now()}`
            summary.value = {
                id: contextId.value,
                name: generateSmartName(filePaths),
                fileCount: result.fileCount || filePaths.length,
                totalSize: result.totalSize || 0,
                lineCount: result.lineCount || 0,
                tokenCount: result.tokenCount,
                createdAt: new Date().toISOString(),
                files: filePaths,
                isFavorite: false
            }

            validateBuildResult(summary.value, result)

            warnings.value = result.metadata?.warnings || []
            skippedFiles.value = result.metadata?.skippedFiles || []

            // Phase 3: Loading preview (90% -> 95%)
            buildProgress.value = 95
            contentManager.clearCache()
            if (contextId.value && summary.value.fileCount > 0) {
                await loadContextContent(contextId.value, 0, 0)
            }

            // Phase 4: Complete
            cleanupProgressListener()
            buildProgress.value = 100

            // Small delay to show 100% before hiding
            await new Promise(resolve => setTimeout(resolve, 200))
        } catch (err) {
            cleanupProgressListener()
            handleBuildError(err)
            throw err
        } finally {
            cleanupProgressListener()
            isBuilding.value = false
            isLoading.value = false
            buildProgress.value = 0
        }
    }

    function createBuildOptions(options?: Partial<domain.ContextBuildOptions>): domain.ContextBuildOptions {
        const settingsStore = useSettingsStore()
        const contextSettings = settingsStore.settings.context

        return {
            maxTokens: options?.maxTokens || contextSettings.maxTokens,
            maxMemoryMB: options?.maxMemoryMB || 50,
            stripComments: options?.stripComments ?? contextSettings.stripComments,
            includeManifest: options?.includeManifest ?? true,
            includeLineNumbers: options?.includeLineNumbers ?? false,
            includeTests: options?.includeTests ?? true,
            splitStrategy: options?.splitStrategy || 'smart',
            forceStream: true,
            enableProgressEvents: true,
            outputFormat: options?.outputFormat || contextSettings.outputFormat,
            excludeTests: options?.excludeTests ?? contextSettings.excludeTests,
            collapseEmptyLines: options?.collapseEmptyLines ?? contextSettings.collapseEmptyLines,
            stripLicense: options?.stripLicense ?? contextSettings.stripLicense,
            compactDataFiles: options?.compactDataFiles ?? contextSettings.compactDataFiles,
            trimWhitespace: options?.trimWhitespace ?? contextSettings.trimWhitespace
        } as domain.ContextBuildOptions
    }

    function validateBuildResult(
        sum: ContextSummary,
        result: domain.ContextSummary
    ) {
        if (sum.totalSize > MAX_CONTEXT_SIZE) {
            const sizeMB = Math.round(sum.totalSize / (1024 * 1024))
            const limitMB = Math.round(MAX_CONTEXT_SIZE / (1024 * 1024))
            error.value = `Context size (${sizeMB}MB) exceeds maximum (${limitMB}MB)`
            clearContext()
            throw new Error(error.value)
        }

        if (sum.fileCount === 0 || sum.totalSize === 0 || sum.lineCount === 0) {
            const projectStore = useProjectStore()
            const skipped = result.metadata?.skippedFiles || []
            const skippedInfo = skipped.length > 0
                ? `\n\nSkipped files (${skipped.length}):\n${skipped.slice(0, 5).join('\n')}${skipped.length > 5 ? '\n...' : ''}`
                : ''

            error.value = `Context built with empty content.\n\nCurrent project: ${projectStore.currentPath}\n\nFiles may be outside the project directory.${skippedInfo}`
            throw new Error(error.value)
        }
    }

    async function handleBuildError(err: unknown) {
        if (err instanceof Error && err.message === 'TOKEN_LIMIT_EXCEEDED') {
            const { hasTokenInfo } = await import('@/types/errors')
            const settingsStore = useSettingsStore()
            const currentLimit = settingsStore.settings.context.maxTokens

            if (hasTokenInfo(err)) {
                error.value = `TOKEN_LIMIT_EXCEEDED:${err.tokenInfo.actual}:${err.tokenInfo.limit}`
            } else {
                error.value = `TOKEN_LIMIT_EXCEEDED:0:${currentLimit}`
            }
            return
        }

        let errorMessage = err instanceof Error ? err.message : 'Failed to build context'
        if (errorMessage.includes('path traversal') || errorMessage.includes('invalid path')) {
            errorMessage += ' Make sure all selected files are within the project directory.'
        }
        error.value = errorMessage
    }

    async function loadContextContent(ctxId: string, startLine = 0, lineCount = 0) {
        if (!ctxId) {
            error.value = 'No context ID provided'
            return
        }

        isLoading.value = true
        error.value = null

        try {
            const isNewContext = contextId.value !== ctxId
            if (isNewContext) {
                contextId.value = ctxId
                contentManager.clearCache()
            }
            // Always update summary from contextList to ensure files array is present
            await loadSummaryForContext(ctxId)

            currentChunk.value = await contentManager.loadChunk(ctxId, startLine, lineCount)
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Failed to load context'
            throw err
        } finally {
            isLoading.value = false
        }
    }

    async function loadSummaryForContext(ctxId: string) {
        const existing = contextList.value.find(c => c.id === ctxId)
        if (existing) {
            // Copy all fields including files array
            summary.value = { ...existing }
            return
        }

        try {
            const projectStore = useProjectStore()
            const contexts = await contextApi.getProjectContexts(projectStore.currentPath || '')
            const fetched = contexts.find((c: { id: string }) => c.id === ctxId)
            if (fetched) {
                summary.value = {
                    id: fetched.id,
                    name: fetched.name || generateSmartName(fetched.files || []),
                    fileCount: fetched.fileCount || 0,
                    totalSize: fetched.totalSize || 0,
                    lineCount: fetched.lineCount || 0,
                    tokenCount: fetched.tokenCount,
                    createdAt: fetched.createdAt,
                    files: fetched.files || [],
                    isFavorite: false
                }
            }
        } catch (err) {
            logger.warn('Failed to fetch summary from backend:', err)
        }
    }

    async function getFullContextContent(): Promise<string> {
        if (!contextId.value) throw new Error('No context ID available')
        return contentManager.loadFullContent(contextId.value)
    }

    async function deleteContext(ctxId: string) {
        try {
            await contextApi.deleteContext(ctxId)
            if (contextId.value === ctxId) clearContext()
            contextList.value = contextList.value.filter(c => c.id !== ctxId)
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Failed to delete context'
            throw err
        }
    }

    async function exportContext(ctxId: string) {
        if (!ctxId) throw new Error('No context ID provided')
        return contextApi.exportContext({ format: 'markdown', includeMetadata: true })
    }

    async function listProjectContexts() {
        isLoading.value = true
        error.value = null

        try {
            const projectStore = useProjectStore()
            if (!projectStore.currentPath) throw new Error('No project selected')

            const contexts = await contextApi.getProjectContexts(projectStore.currentPath)
            const contextArray = Array.isArray(contexts) ? contexts : []

            contextList.value = contextArray.map((ctx: { id: string; name?: string; files?: string[]; totalSize?: number; fileCount?: number; lineCount?: number; tokenCount?: number; createdAt?: string }) => ({
                id: ctx.id,
                name: ctx.name || generateSmartName(ctx.files || []),
                fileCount: ctx.fileCount || (ctx.files?.length || 0),
                totalSize: ctx.totalSize || 0,
                lineCount: ctx.lineCount || 0,
                tokenCount: ctx.tokenCount,
                createdAt: ctx.createdAt,
                files: ctx.files || [],
                isFavorite: false
            }))

            loadContextMetadata(contextList.value)
            await autoCleanup()
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Failed to list contexts'
            throw err
        } finally {
            isLoading.value = false
        }
    }

    function clearContext() {
        contextId.value = null
        summary.value = null
        currentChunk.value = null
        error.value = null
        contentManager.clearCache()
        clearHighlightCache() // Очистка кэша подсветки синтаксиса

        if (typeof window !== 'undefined' && window.gc) {
            try { window.gc() } catch { /* ignore */ }
        }
    }

    function setRawContext(content: string, fileCount: number) {
        const id = `git-ref-${Date.now()}`
        const lines = content.split('\n')

        contextId.value = id
        summary.value = createContextSummary(id, `Git ref [${fileCount}]`, content, fileCount)
        currentChunk.value = {
            lines: lines.slice(0, 100),
            startLine: 0,
            endLine: Math.min(100, lines.length),
            hasMore: lines.length > 100
        }
    }

    function getMemoryUsage(): number {
        let size = summary.value ? 1024 : 0
        if (currentChunk.value?.lines) {
            size += currentChunk.value.lines.length * 80 * 2
        }
        size += contentManager.getMemoryUsage()
        return size
    }

    function renameContext(ctxId: string, newName: string) {
        const ctx = contextList.value.find(c => c.id === ctxId)
        if (ctx) ctx.name = newName
        if (summary.value?.id === ctxId) summary.value.name = newName
        saveMetadata()
    }

    function toggleFavorite(ctxId: string) {
        const ctx = contextList.value.find(c => c.id === ctxId)
        if (ctx) ctx.isFavorite = !ctx.isFavorite
        if (summary.value?.id === ctxId) summary.value.isFavorite = !summary.value.isFavorite
        saveMetadata()
    }

    function saveMetadata() {
        saveContextMetadata(contextList.value)
    }

    async function autoCleanup() {
        await performAutoCleanup(contextList.value, deleteContext)
    }

    async function duplicateContext(ctxId: string): Promise<string | null> {
        const result = await duplicateContextContent(ctxId, contextList.value)
        if (!result) return null

        setRawContext(result.content, result.fileCount)
        if (summary.value) {
            summary.value.id = result.id
            summary.value.name = result.name
        }

        contextList.value.push(createContextSummary(result.id, result.name, result.content, result.fileCount))
        saveMetadata()
        return result.id
    }

    async function mergeContexts(contextIds: string[]): Promise<string | null> {
        const result = await mergeContextsContent(
            contextIds,
            contextList.value,
            async (ctxId) => {
                if (contextId.value !== ctxId) await loadContextContent(ctxId, 0, 0)
                return getFullContextContent()
            }
        )
        if (!result) return null

        setRawContext(result.content, result.fileCount)
        if (summary.value) {
            summary.value.id = result.id
            summary.value.name = result.name
        }

        contextList.value.push(createContextSummary(result.id, result.name, result.content, result.fileCount))
        saveMetadata()
        return result.id
    }

    function reorderContexts(fromIndex: number, toIndex: number) {
        contextList.value = reorderContextList(contextList.value, fromIndex, toIndex)
        saveMetadata()
    }

    // Select item from list for preview in right panel
    function selectListItem(ctxId: string | null) {
        if (!ctxId) {
            selectedListItem.value = null
            return
        }
        const item = contextList.value.find(c => c.id === ctxId)
        selectedListItem.value = item ? { ...item } : null
    }

    function clearListSelection() {
        selectedListItem.value = null
    }

    /**
     * Remove a file from current context and rebuild
     */
    async function removeFileFromContext(filePath: string) {
        if (!summary.value?.files) return

        const newFiles = summary.value.files.filter(f => f !== filePath)
        if (newFiles.length === 0) {
            clearContext()
            return
        }

        // Rebuild context without the removed file
        await buildContext(newFiles)
    }

    async function rebuildContext() {
        if (!summary.value?.files || summary.value.files.length === 0) {
            logger.warn('No files to rebuild context')
            return
        }

        const settingsStore = useSettingsStore()
        const options = {
            outputFormat: settingsStore.settings.context.outputFormat,
            stripComments: settingsStore.settings.context.stripComments,
            excludeTests: settingsStore.settings.context.excludeTests,
            collapseEmptyLines: settingsStore.settings.context.collapseEmptyLines,
            stripLicense: settingsStore.settings.context.stripLicense,
            compactDataFiles: settingsStore.settings.context.compactDataFiles,
            trimWhitespace: settingsStore.settings.context.trimWhitespace,
            maxTokens: settingsStore.settings.context.maxTokens,
        }

        await buildContext(summary.value.files, options)
    }

    return {
        // State
        contextId, summary, currentChunk, isLoading, isBuilding,
        buildProgress, error, contextList, warnings, skippedFiles,
        selectedListItem,
        // Computed
        hasContext, totalSize, fileCount, lineCount, tokenCount, totalTokens, estimatedCost, stats,
        // Actions
        buildContext, rebuildContext, loadContextContent, deleteContext, exportContext,
        listProjectContexts, clearContext, setRawContext, getFullContextContent,
        getMemoryUsage, renameContext, toggleFavorite, duplicateContext,
        autoCleanup, loadContextMetadata, saveContextMetadata, generateSmartName,
        mergeContexts, reorderContexts, selectListItem, clearListSelection,
        removeFileFromContext
    }
})

export { generateSmartName }

