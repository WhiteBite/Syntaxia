/**
 * Mentions parser composable for chat input
 * Parses @file, @folder:, @symbol:, @git:, @docs mentions
 */

import { useI18n } from '@/composables/useI18n'
import { useContextStore } from '@/features/context'
import { useFileStore } from '@/features/files'
import { filesApi, gitApi, semanticApi } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, ref, watch, type Ref } from 'vue'
import type {
    Mention,
    MentionSuggestion,
    MentionType,
    ParsedInput,
} from '../types/mentions'
import { MENTION_PATTERNS } from '../types/mentions'

const DEBOUNCE_MS = 150

export interface UseMentionsOptions {
    inputText?: Ref<string>
}

export interface MentionResult {
    type: 'files' | 'git' | 'problems'
    success: boolean
    message: string
    filesAdded?: number
}

type TranslateFunc = (key: string, params?: Record<string, string | number>) => string

export function useMentions(options?: UseMentionsOptions) {
    const inputText = options?.inputText ?? ref('')
    const { t } = useI18n()
    const projectStore = useProjectStore()
    const fileStore = useFileStore()
    const contextStore = useContextStore()
    const uiStore = useUIStore()

    const isProcessing = ref(false)
    const suggestions = ref<MentionSuggestion[]>([])
    const parsedInput = ref<ParsedInput>({
        text: '',
        mentions: [],
        resolvedFiles: [],
    })

    /** Parse mentions from text */
    function parseMentions(text: string): Mention[] {
        const mentions: Mention[] = []

        for (const pattern of MENTION_PATTERNS) {
            const regex = new RegExp(pattern.regex.source, 'g')
            let match: RegExpExecArray | null

            while ((match = regex.exec(text)) !== null) {
                const value = match[1] || match[0].slice(1) // Remove @ prefix
                const display = getDisplayName(pattern.type, value)

                mentions.push({
                    type: pattern.type,
                    value,
                    display,
                    startIndex: match.index,
                    endIndex: match.index + match[0].length,
                })
            }
        }

        // Sort by position
        return mentions.sort((a, b) => a.startIndex - b.startIndex)
    }

    /** Get display name for mention */
    function getDisplayName(type: MentionType, value: string): string {
        switch (type) {
            case 'file':
                return value.split('/').pop() || value
            case 'folder':
                return value.endsWith('/') ? value : `${value}/`
            case 'symbol':
                return value
            case 'git':
                return `git:${value}`
            case 'docs':
                return t('mentions.docs')
            default:
                return value
        }
    }

    /** Remove mentions from text */
    function stripMentions(text: string, mentions: Mention[]): string {
        let result = text
        // Process in reverse order to preserve indices
        for (let i = mentions.length - 1; i >= 0; i--) {
            const mention = mentions[i]
            result = result.slice(0, mention.startIndex) + result.slice(mention.endIndex)
        }
        return result.replace(/\s+/g, ' ').trim()
    }

    /** Resolve mentions to file paths */
    async function resolveMentions(mentions: Mention[]): Promise<string[]> {
        const projectPath = projectStore.currentPath
        if (!projectPath) return []

        const resolvedFiles: string[] = []

        for (const mention of mentions) {
            const files = await resolveSingleMention(mention, projectPath)
            resolvedFiles.push(...files)
        }

        // Deduplicate
        return [...new Set(resolvedFiles)]
    }

    /** Resolve a single mention to files */
    async function resolveSingleMention(mention: Mention, projectPath: string): Promise<string[]> {
        switch (mention.type) {
            case 'file':
                return resolveFileMention(mention.value, projectPath)
            case 'folder':
                return resolveFolderMention(mention.value, projectPath)
            case 'symbol':
                return resolveSymbolMention(mention.value, projectPath)
            case 'git':
                return resolveGitMention(mention.value, projectPath)
            case 'docs':
                return resolveDocsMention(projectPath)
            default:
                return []
        }
    }

    /** Find files matching filename pattern */
    async function resolveFileMention(filename: string, projectPath: string): Promise<string[]> {
        try {
            const allFiles = await filesApi.listFiles(projectPath)
            const flatFiles = flattenFileTree(allFiles)
            const lowerFilename = filename.toLowerCase()

            return flatFiles
                .filter(f => f.toLowerCase().includes(lowerFilename))
                .slice(0, 10)
        } catch {
            return []
        }
    }

    /** Get all files in a folder */
    async function resolveFolderMention(folderPath: string, projectPath: string): Promise<string[]> {
        try {
            const allFiles = await filesApi.listFiles(projectPath)
            const flatFiles = flattenFileTree(allFiles)
            const normalizedPath = folderPath.replace(/\/$/, '')

            return flatFiles.filter(f => f.startsWith(normalizedPath + '/'))
        } catch {
            return []
        }
    }

    /** Find files containing symbol */
    async function resolveSymbolMention(symbolName: string, projectPath: string): Promise<string[]> {
        try {
            const result = await semanticApi.search({
                query: symbolName,
                projectRoot: projectPath,
                topK: 10,
            })
            return result.results?.map(r => r.chunk.filePath) || []
        } catch {
            return []
        }
    }

    /** Get git-related files */
    async function resolveGitMention(subtype: string, projectPath: string): Promise<string[]> {
        try {
            switch (subtype) {
                case 'diff':
                case 'staged': {
                    const uncommitted = await gitApi.getUncommittedFiles(projectPath)
                    return uncommitted?.map(f => f.path) || []
                }
                case 'history': {
                    const history = await gitApi.getRichCommitHistory(projectPath, '', 5)
                    const files: string[] = []
                    history?.forEach(commit => {
                        commit.files?.forEach(f => files.push(f))
                    })
                    return [...new Set(files)]
                }
                default:
                    return []
            }
        } catch {
            return []
        }
    }

    /** Get documentation files */
    async function resolveDocsMention(projectPath: string): Promise<string[]> {
        try {
            const allFiles = await filesApi.listFiles(projectPath)
            const flatFiles = flattenFileTree(allFiles)
            const docPatterns = [
                /readme\.md$/i,
                /docs?\//i,
                /\.md$/i,
                /changelog/i,
                /contributing/i,
            ]

            return flatFiles.filter(f =>
                docPatterns.some(pattern => pattern.test(f))
            ).slice(0, 20)
        } catch {
            return []
        }
    }

    /** Flatten file tree to array of paths */
    function flattenFileTree(nodes: { path: string; children?: unknown[] }[]): string[] {
        const result: string[] = []

        function traverse(items: { path: string; children?: unknown[] }[]) {
            for (const item of items) {
                if (!item.children || item.children.length === 0) {
                    result.push(item.path)
                } else {
                    traverse(item.children as { path: string; children?: unknown[] }[])
                }
            }
        }

        traverse(nodes)
        return result
    }

    /** Insert mention at cursor position */
    function insertMention(suggestion: MentionSuggestion, cursorPosition: number): string {
        const text = inputText.value
        const beforeCursor = text.slice(0, cursorPosition)
        const afterCursor = text.slice(cursorPosition)

        // Find the @ trigger position
        const atIndex = beforeCursor.lastIndexOf('@')
        if (atIndex === -1) return text

        const prefix = text.slice(0, atIndex)
        const mentionText = formatMentionText(suggestion)

        return `${prefix}${mentionText} ${afterCursor}`.trim()
    }

    /** Format mention for insertion */
    function formatMentionText(suggestion: MentionSuggestion): string {
        switch (suggestion.type) {
            case 'file':
                return `@${suggestion.value}`
            case 'folder':
                return `@folder:${suggestion.value}`
            case 'symbol':
                return `@symbol:${suggestion.value}`
            case 'git':
                return `@git:${suggestion.value}`
            case 'docs':
                return '@docs'
            default:
                return `@${suggestion.value}`
        }
    }

    /** Highlight mentions in text for display */
    function highlightMentions(text: string): string {
        const mentions = parseMentions(text)
        if (mentions.length === 0) return text

        let result = ''
        let lastIndex = 0

        for (const mention of mentions) {
            result += text.slice(lastIndex, mention.startIndex)
            result += `<span class="mention mention-${mention.type}">${text.slice(mention.startIndex, mention.endIndex)}</span>`
            lastIndex = mention.endIndex
        }

        result += text.slice(lastIndex)
        return result
    }

    /** Process input and resolve all mentions */
    async function processInput(): Promise<ParsedInput> {
        const text = inputText.value
        const mentions = parseMentions(text)
        const cleanText = stripMentions(text, mentions)

        isProcessing.value = true
        try {
            const resolvedFiles = await resolveMentions(mentions)

            parsedInput.value = {
                text: cleanText,
                mentions,
                resolvedFiles,
            }

            return parsedInput.value
        } finally {
            isProcessing.value = false
        }
    }

    /** Build context from resolved files */
    async function buildContextFromMentions(): Promise<boolean> {
        if (parsedInput.value.resolvedFiles.length === 0) {
            uiStore.addToast(t('mentions.noFilesResolved'), 'warning')
            return false
        }

        try {
            await contextStore.buildContext(parsedInput.value.resolvedFiles)
            uiStore.addToast(t('chat.contextAttached'), 'success')
            return true
        } catch {
            uiStore.addToast(t('chat.contextBuildFailed'), 'error')
            return false
        }
    }

    /** Check if text contains any mentions */
    const hasMentions = computed(() => {
        return parseMentions(inputText.value).length > 0
    })

    /** Get mention at cursor position */
    function getMentionAtCursor(cursorPosition: number): Mention | null {
        const mentions = parseMentions(inputText.value)
        return mentions.find(m =>
            cursorPosition >= m.startIndex && cursorPosition <= m.endIndex
        ) || null
    }

    // Watch input and parse mentions
    let debounceTimer: ReturnType<typeof setTimeout> | null = null
    watch(inputText, () => {
        if (debounceTimer) clearTimeout(debounceTimer)
        debounceTimer = setTimeout(() => {
            const mentions = parseMentions(inputText.value)
            parsedInput.value = {
                ...parsedInput.value,
                mentions,
                text: stripMentions(inputText.value, mentions),
            }
        }, DEBOUNCE_MS)
    })

    // ============================================
    // Legacy API for backward compatibility
    // ============================================

    /**
     * Process @files mention - build context from selected files
     */
    async function processFilesMention(translateFn: TranslateFunc): Promise<MentionResult> {
        const selectedFiles = fileStore.selectedFilesList

        if (selectedFiles.length === 0) {
            uiStore.addToast(translateFn('chat.selectFilesHint'), 'warning')
            return {
                type: 'files',
                success: false,
                message: translateFn('chat.selectFilesHint')
            }
        }

        try {
            isProcessing.value = true
            await contextStore.buildContext(selectedFiles)

            const message = translateFn('chat.contextAttached')
            uiStore.addToast(message, 'success')

            return {
                type: 'files',
                success: true,
                message,
                filesAdded: selectedFiles.length
            }
        } catch {
            const message = translateFn('chat.contextBuildFailed')
            uiStore.addToast(message, 'error')
            return {
                type: 'files',
                success: false,
                message
            }
        } finally {
            isProcessing.value = false
        }
    }

    /**
     * Process @git mention - get uncommitted files and build context
     */
    async function processGitMention(translateFn: TranslateFunc): Promise<MentionResult> {
        const projectPath = projectStore.currentPath

        if (!projectPath) {
            uiStore.addToast(translateFn('chat.noApiKey'), 'warning')
            return {
                type: 'git',
                success: false,
                message: 'No project selected'
            }
        }

        try {
            isProcessing.value = true

            const uncommittedFiles = await gitApi.getUncommittedFiles(projectPath)

            if (!uncommittedFiles || uncommittedFiles.length === 0) {
                const message = t('gitContext.noChanges')
                uiStore.addToast(message, 'info')
                return {
                    type: 'git',
                    success: true,
                    message,
                    filesAdded: 0
                }
            }

            const filePaths = uncommittedFiles.map(f => f.path)
            await contextStore.buildContext(filePaths)

            const message = `${translateFn('chat.contextAttached')} (${filePaths.length} git files)`
            uiStore.addToast(message, 'success')

            return {
                type: 'git',
                success: true,
                message,
                filesAdded: filePaths.length
            }
        } catch {
            const message = translateFn('chat.contextBuildFailed')
            uiStore.addToast(message, 'error')
            return {
                type: 'git',
                success: false,
                message
            }
        } finally {
            isProcessing.value = false
        }
    }

    /**
     * Process @problems mention - placeholder for future implementation
     */
    async function processProblemsMention(translateFn: TranslateFunc): Promise<MentionResult> {
        uiStore.addToast(translateFn('chat.comingSoon'), 'info')
        return {
            type: 'problems',
            success: false,
            message: translateFn('chat.comingSoon')
        }
    }

    /**
     * Parse message for mentions and process them (legacy API)
     */
    async function processMessageMentions(
        message: string,
        translateFn: TranslateFunc
    ): Promise<{ cleanedMessage: string; results: MentionResult[] }> {
        const results: MentionResult[] = []
        let cleanedMessage = message

        if (message.includes('@files')) {
            const result = await processFilesMention(translateFn)
            results.push(result)
            cleanedMessage = cleanedMessage.replace(/@files\s*/g, '')
        }

        if (message.includes('@git')) {
            const result = await processGitMention(translateFn)
            results.push(result)
            cleanedMessage = cleanedMessage.replace(/@git\s*/g, '')
        }

        if (message.includes('@problems')) {
            const result = await processProblemsMention(translateFn)
            results.push(result)
            cleanedMessage = cleanedMessage.replace(/@problems\s*/g, '')
        }

        return {
            cleanedMessage: cleanedMessage.trim(),
            results
        }
    }

    return {
        // State
        isProcessing,
        suggestions,
        parsedInput,
        hasMentions,
        // New API methods
        parseMentions,
        processInput,
        insertMention,
        highlightMentions,
        buildContextFromMentions,
        getMentionAtCursor,
        resolveMentions,
        // Legacy API methods
        processFilesMention,
        processGitMention,
        processProblemsMention,
        processMessageMentions,
    }
}
