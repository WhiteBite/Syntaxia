import { useContextStore } from '@/features/context/model/context.store'
import { useSettingsStore } from '@/stores/settings.store'
import { computed, ref } from 'vue'

export interface Chunk {
    index: number
    content: string
    tokens: number
    files: number
    preview: string
    copied: boolean
}

/**
 * Composable for chunked context copying
 * Splits context into parts based on token limit for easy copying to LLMs
 */
export function useChunkedCopy() {
    const contextStore = useContextStore()
    const settingsStore = useSettingsStore()

    const chunks = ref<Chunk[]>([])
    const isProcessing = ref(false)
    const currentChunkIndex = ref(0)
    const lastCopiedIndex = ref(-1)

    // Settings from store
    const maxTokensPerChunk = computed(() => settingsStore.settings.context.maxTokensPerChunk)
    const splitStrategy = computed(() => settingsStore.settings.context.splitStrategy)
    const enableAutoSplit = computed(() => settingsStore.settings.context.enableAutoSplit)

    // Computed
    const totalTokens = computed(() => contextStore.tokenCount)
    const needsSplitting = computed(() => totalTokens.value > maxTokensPerChunk.value)
    const chunkCount = computed(() => chunks.value.length)
    const allCopied = computed(() => chunks.value.length > 0 && chunks.value.every(c => c.copied))
    const progress = computed(() => {
        if (chunks.value.length === 0) return 0
        const copied = chunks.value.filter(c => c.copied).length
        return Math.round((copied / chunks.value.length) * 100)
    })

    /**
     * Estimate tokens (1 token ≈ 4 chars)
     */
    function estimateTokens(text: string): number {
        return Math.ceil(text.length / 4)
    }

    /**
     * Split content by file headers (--- File: ... ---)
     */
    function splitByFiles(content: string, tokenLimit: number): Chunk[] {
        const fileRegex = /^--- File: .+? ---$/gm
        const matches = [...content.matchAll(fileRegex)]

        if (matches.length === 0) {
            return splitByTokens(content, tokenLimit)
        }

        const result: Chunk[] = []
        let currentContent = ''
        let currentTokens = 0
        let currentFiles = 0
        let chunkIndex = 0

        // Handle content before first file header (manifest)
        const firstMatch = matches[0]
        if (firstMatch.index && firstMatch.index > 0) {
            currentContent = content.slice(0, firstMatch.index)
            currentTokens = estimateTokens(currentContent)
        }

        for (let i = 0; i < matches.length; i++) {
            const start = matches[i].index!
            const end = i + 1 < matches.length ? matches[i + 1].index! : content.length
            const fileContent = content.slice(start, end)
            const fileTokens = estimateTokens(fileContent)

            if (currentTokens + fileTokens > tokenLimit && currentContent) {
                // Save current chunk
                result.push(createChunk(chunkIndex++, currentContent, currentTokens, currentFiles))
                currentContent = ''
                currentTokens = 0
                currentFiles = 0
            }

            currentContent += fileContent
            currentTokens += fileTokens
            currentFiles++
        }

        // Save last chunk
        if (currentContent) {
            result.push(createChunk(chunkIndex, currentContent, currentTokens, currentFiles))
        }

        return result
    }

    /**
     * Split content by token count
     */
    function splitByTokens(content: string, tokenLimit: number): Chunk[] {
        const result: Chunk[] = []
        const charLimit = tokenLimit * 4
        let pos = 0
        let chunkIndex = 0

        while (pos < content.length) {
            let end = Math.min(pos + charLimit, content.length)

            // Try to break at newline
            if (end < content.length) {
                const lastNewline = content.lastIndexOf('\n', end)
                if (lastNewline > pos) {
                    end = lastNewline + 1
                }
            }

            const chunkContent = content.slice(pos, end)
            const tokens = estimateTokens(chunkContent)
            const files = (chunkContent.match(/^--- File: /gm) || []).length

            result.push(createChunk(chunkIndex++, chunkContent, tokens, files))
            pos = end
        }

        return result
    }

    /**
     * Smart split: try files first, fallback to tokens
     */
    function splitSmart(content: string, tokenLimit: number): Chunk[] {
        const fileChunks = splitByFiles(content, tokenLimit)

        // Check if any chunk is still too large
        const hasOversized = fileChunks.some(c => c.tokens > tokenLimit * 1.2)

        if (hasOversized) {
            return splitByTokens(content, tokenLimit)
        }

        return fileChunks
    }

    /**
     * Create chunk object
     */
    function createChunk(index: number, content: string, tokens: number, files: number): Chunk {
        const lines = content.split('\n').slice(0, 3)
        const preview = lines.join('\n').slice(0, 100) + (content.length > 100 ? '...' : '')

        return {
            index,
            content,
            tokens,
            files,
            preview,
            copied: false
        }
    }

    /**
     * Generate chunks from current context
     */
    async function generateChunks(): Promise<void> {
        if (!contextStore.hasContext) return

        isProcessing.value = true
        chunks.value = []
        currentChunkIndex.value = 0
        lastCopiedIndex.value = -1

        try {
            const content = await contextStore.getFullContextContent()
            const tokenLimit = maxTokensPerChunk.value
            const strategy = splitStrategy.value

            let result: Chunk[]
            switch (strategy) {
                case 'file':
                    result = splitByFiles(content, tokenLimit)
                    break
                case 'token':
                    result = splitByTokens(content, tokenLimit)
                    break
                case 'smart':
                default:
                    result = splitSmart(content, tokenLimit)
            }

            chunks.value = result
        } catch (err) {
            console.error('[useChunkedCopy] Failed to generate chunks:', err)
        } finally {
            isProcessing.value = false
        }
    }

    /**
     * Copy specific chunk to clipboard
     */
    async function copyChunk(index: number): Promise<boolean> {
        if (index < 0 || index >= chunks.value.length) return false

        try {
            await navigator.clipboard.writeText(chunks.value[index].content)
            chunks.value[index].copied = true
            lastCopiedIndex.value = index
            currentChunkIndex.value = Math.min(index + 1, chunks.value.length - 1)
            return true
        } catch (err) {
            console.error('[useChunkedCopy] Failed to copy chunk:', err)
            return false
        }
    }

    /**
     * Copy next uncoped chunk
     */
    async function copyNext(): Promise<boolean> {
        const nextIndex = chunks.value.findIndex(c => !c.copied)
        if (nextIndex === -1) return false
        return copyChunk(nextIndex)
    }

    /**
     * Reset all copied states
     */
    function resetCopied(): void {
        chunks.value.forEach(c => c.copied = false)
        lastCopiedIndex.value = -1
        currentChunkIndex.value = 0
    }

    /**
     * Clear chunks
     */
    function clear(): void {
        chunks.value = []
        currentChunkIndex.value = 0
        lastCopiedIndex.value = -1
    }

    return {
        // State
        chunks,
        isProcessing,
        currentChunkIndex,
        lastCopiedIndex,

        // Computed
        totalTokens,
        maxTokensPerChunk,
        needsSplitting,
        chunkCount,
        allCopied,
        progress,
        enableAutoSplit,

        // Actions
        generateChunks,
        copyChunk,
        copyNext,
        resetCopied,
        clear
    }
}
