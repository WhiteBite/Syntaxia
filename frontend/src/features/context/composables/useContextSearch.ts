import { nextTick, ref, watch } from 'vue'

export interface SearchState {
    showSearch: boolean
    searchQuery: string
    searchResults: number[]
    currentSearchIndex: number
}

export function useContextSearch(getLines: () => { lines: string[]; startLine: number } | null) {
    const showSearch = ref(false)
    const searchQuery = ref('')
    const searchResults = ref<number[]>([])
    const currentSearchIndex = ref(0)
    const searchInput = ref<HTMLInputElement | null>(null)
    const lineRefs = ref<Map<number, HTMLElement>>(new Map())

    // Search when query changes
    watch(searchQuery, (query) => {
        const chunk = getLines()
        if (!query || !chunk?.lines) {
            searchResults.value = []
            currentSearchIndex.value = 0
            return
        }

        const results: number[] = []
        const lowerQuery = query.toLowerCase()

        chunk.lines.forEach((line, index) => {
            if (line.toLowerCase().includes(lowerQuery)) {
                results.push(chunk.startLine + index)
            }
        })

        searchResults.value = results
        currentSearchIndex.value = 0

        if (results.length > 0) {
            scrollToLine(results[0])
        }
    })

    // Focus input when search opens
    watch(showSearch, (show) => {
        if (show) {
            nextTick(() => searchInput.value?.focus())
        } else {
            searchQuery.value = ''
        }
    })

    function setLineRef(el: Element | null, lineNum: number | null) {
        if (el && lineNum !== null && lineNum !== undefined) {
            lineRefs.value.set(lineNum, el as HTMLElement)
        }
    }

    function isLineHighlighted(lineNum: number): boolean {
        return searchResults.value.includes(lineNum) &&
            searchResults.value[currentSearchIndex.value] === lineNum
    }

    function scrollToLine(lineNum: number) {
        nextTick(() => {
            const el = lineRefs.value.get(lineNum)
            if (el) {
                el.scrollIntoView({ behavior: 'smooth', block: 'center' })
            }
        })
    }

    function searchNext() {
        if (searchResults.value.length === 0) return
        currentSearchIndex.value = (currentSearchIndex.value + 1) % searchResults.value.length
        scrollToLine(searchResults.value[currentSearchIndex.value])
    }

    function searchPrev() {
        if (searchResults.value.length === 0) return
        currentSearchIndex.value = (currentSearchIndex.value - 1 + searchResults.value.length) % searchResults.value.length
        scrollToLine(searchResults.value[currentSearchIndex.value])
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
        searchInput,
        setLineRef,
        isLineHighlighted,
        searchNext,
        searchPrev,
        toggleSearch,
        closeSearch,
    }
}
