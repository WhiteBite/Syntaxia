import { useSettingsStore } from '@/stores/settings.store'
import { computed, ref, shallowRef } from 'vue'

export interface ChunkInfo {
    index: number
    startLine: number
    endLine: number
    tokenCount: number
}

/**
 * Composable for managing context chunking logic.
 * Handles chunk calculation, navigation, and copy state.
 */
export function useChunking() {
    const settingsStore = useSettingsStore()

    // State
    const currentChunkIndex = ref(0)
    const chunks = shallowRef<ChunkInfo[]>([])
    const copiedChunks = ref<Set<number>>(new Set())

    // Settings
    const isEnabled = computed(() => settingsStore.settings.context.enableAutoSplit)
    const maxTokensPerChunk = computed(() => settingsStore.settings.context.maxTokensPerChunk)
    const strategy = computed(() => settingsStore.settings.context.splitStrategy)

    // Derived state
    const totalChunks = computed(() => chunks.value.length)
    const currentChunk = computed(() => currentChunkIndex.value + 1)
    const hasMultipleChunks = computed(() => totalChunks.value > 1)
    const allChunksCopied = computed(() => copiedChunks.value.size >= totalChunks.value)

    const currentChunkInfo = computed(() =>
        chunks.value[currentChunkIndex.value] ?? null
    )

    /**
     * Calculate chunk boundaries based on content and settings.
     * Uses smart strategy (file boundaries) or hard strategy (token limit).
     */
    function calculateChunks(
        lines: string[],
        fileMarkers: number[], // Line numbers where files start
        tokensPerLine: number[] // Estimated tokens per line
    ): ChunkInfo[] {
        if (!isEnabled.value || lines.length === 0) {
            return [{
                index: 0,
                startLine: 0,
                endLine: lines.length - 1,
                tokenCount: tokensPerLine.reduce((a, b) => a + b, 0)
            }]
        }

        const limit = maxTokensPerChunk.value
        const result: ChunkInfo[] = []

        if (strategy.value === 'smart' || strategy.value === 'file') {
            // Smart strategy: don't break files
            result.push(...calculateSmartChunks(lines, fileMarkers, tokensPerLine, limit))
        } else {
            // Hard/token strategy: strict token limit
            result.push(...calculateHardChunks(tokensPerLine, limit))
        }

        return result
    }

    function calculateSmartChunks(
        lines: string[],
        fileMarkers: number[],
        tokensPerLine: number[],
        limit: number
    ): ChunkInfo[] {
        const result: ChunkInfo[] = []
        let chunkStart = 0
        let chunkTokens = 0
        let chunkIndex = 0

        // Add end marker for easier iteration
        const markers = [...fileMarkers, lines.length]

        for (let i = 0; i < markers.length - 1; i++) {
            const fileStart = markers[i]
            const fileEnd = markers[i + 1] - 1

            // Calculate tokens for this file
            let fileTokens = 0
            for (let j = fileStart; j <= fileEnd; j++) {
                fileTokens += tokensPerLine[j] || 0
            }

            // If adding this file exceeds limit, start new chunk
            if (chunkTokens + fileTokens > limit && chunkTokens > 0) {
                result.push({
                    index: chunkIndex++,
                    startLine: chunkStart,
                    endLine: fileStart - 1,
                    tokenCount: chunkTokens
                })
                chunkStart = fileStart
                chunkTokens = 0
            }

            chunkTokens += fileTokens
        }

        // Add final chunk
        if (chunkTokens > 0 || result.length === 0) {
            result.push({
                index: chunkIndex,
                startLine: chunkStart,
                endLine: lines.length - 1,
                tokenCount: chunkTokens
            })
        }

        return result
    }

    function calculateHardChunks(
        tokensPerLine: number[],
        limit: number
    ): ChunkInfo[] {
        const result: ChunkInfo[] = []
        let chunkStart = 0
        let chunkTokens = 0
        let chunkIndex = 0

        for (let i = 0; i < tokensPerLine.length; i++) {
            const lineTokens = tokensPerLine[i] || 0

            if (chunkTokens + lineTokens > limit && chunkTokens > 0) {
                result.push({
                    index: chunkIndex++,
                    startLine: chunkStart,
                    endLine: i - 1,
                    tokenCount: chunkTokens
                })
                chunkStart = i
                chunkTokens = 0
            }

            chunkTokens += lineTokens
        }

        // Add final chunk
        if (chunkTokens > 0 || result.length === 0) {
            result.push({
                index: chunkIndex,
                startLine: chunkStart,
                endLine: tokensPerLine.length - 1,
                tokenCount: chunkTokens
            })
        }

        return result
    }

    /**
     * Set chunks from external calculation (e.g., from backend)
     */
    function setChunks(newChunks: ChunkInfo[]) {
        chunks.value = newChunks
        currentChunkIndex.value = 0
        copiedChunks.value = new Set()
    }

    /**
     * Navigate to specific chunk
     */
    function goToChunk(index: number) {
        if (index >= 0 && index < totalChunks.value) {
            currentChunkIndex.value = index
        }
    }

    function nextChunk() {
        if (currentChunkIndex.value < totalChunks.value - 1) {
            currentChunkIndex.value++
        }
    }

    function prevChunk() {
        if (currentChunkIndex.value > 0) {
            currentChunkIndex.value--
        }
    }

    /**
     * Mark chunk as copied and auto-advance
     */
    function markCopied(index: number, autoAdvance = true) {
        copiedChunks.value.add(index)

        if (autoAdvance && index < totalChunks.value - 1) {
            // Small delay for visual feedback
            setTimeout(() => {
                currentChunkIndex.value = index + 1
            }, 500)
        }
    }

    /**
     * Reset copy state
     */
    function resetCopyState() {
        copiedChunks.value = new Set()
        currentChunkIndex.value = 0
    }

    /**
     * Check if chunk is copied
     */
    function isChunkCopied(index: number): boolean {
        return copiedChunks.value.has(index)
    }

    /**
     * Get chunk boundaries for rendering cut lines
     */
    function getChunkBoundaries(): number[] {
        if (chunks.value.length <= 1) return []

        // Return end lines of all chunks except the last
        return chunks.value
            .slice(0, -1)
            .map(chunk => chunk.endLine)
    }

    return {
        // State
        currentChunkIndex,
        chunks,
        copiedChunks,

        // Computed
        isEnabled,
        totalChunks,
        currentChunk,
        hasMultipleChunks,
        allChunksCopied,
        currentChunkInfo,

        // Methods
        calculateChunks,
        setChunks,
        goToChunk,
        nextChunk,
        prevChunk,
        markCopied,
        resetCopyState,
        isChunkCopied,
        getChunkBoundaries,
    }
}
