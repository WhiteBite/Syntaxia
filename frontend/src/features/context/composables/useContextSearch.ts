import { computed, nextTick, ref, watch, type Ref } from 'vue'

interface UseContextSearchOptions {
    lines: Ref<string[] | undefined>
    startLine: Ref<number>
    scrollToLine: (lineNum: number) => void
}

/**
 * Composable for search functionality in context panel.
 */
export function useContextSearch(options: UseContextSearchOptions) {
    const { lines, startLine, scrollToLine } = options

    const showSearch = ref(false)
    const searchQuery = ref('')
    const searchResults = ref<number[]>([])
    const currentSearchIndex = ref(0)

    // Highlighted lines for VirtualCodeView
    const highlightedLinesSet = computed(() => {
        if (searchResults.value.length === 0) return new Set<number>()
        const currentLine = searchResults.value[currentSearchIndex.value]
        return currentLine !== undefined ? new Set([currentLine]) : new Set<number>()
    })

    // Search when query changes
    watch(searchQuery, (query) => {
        if (!query || !lines.value) {
            searchResults.value = []
            currentSearchIndex.value = 0
            return
        }

        const results: number[] = []
        const lowerQuery = query.toLowerCase()

        lines.value.forEach((line, index) => {
            if (line.toLowerCase().includes(lowerQuery)) {
                results.push(startLine.value + index)
            }
        })

        searchResults.value = results
        currentSearchIndex.value = 0

        if (results.length > 0) {
            doScrollToLine(results[0])
        }
    })

    // Clear search when hidden
    watch(showSearch, (show) => {
        if (!show) {
            searchQuery.value = ''
        }
    })

    function doScrollToLine(lineNum: number) {
        nextTick(() => {
            scrollToLine(lineNum)
        })
    }

    function searchNext() {
        if (searchResults.value.length === 0) return
        currentSearchIndex.value = (currentSearchIndex.value + 1) % searchResults.value.length
        doScrollToLine(searchResults.value[currentSearchIndex.value])
    }

    function searchPrev() {
        if (searchResults.value.length === 0) return
        currentSearchIndex.value = (currentSearchIndex.value - 1 + searchResults.value.length) % searchResults.value.length
        doScrollToLine(searchResults.value[currentSearchIndex.value])
    }

    function toggleSearch() {
        showSearch.value = !showSearch.value
    }

    function closeSearch() {
        showSearch.value = false
    }

    return {
        showSearch,
        searchQuery,
        searchResults,
        currentSearchIndex,
        highlightedLinesSet,
        searchNext,
        searchPrev,
        toggleSearch,
        closeSearch
    }
}
