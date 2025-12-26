import { useContextSearch } from '@/features/context/composables/useContextSearch'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick, ref } from 'vue'

describe('useContextSearch', () => {
    const mockLinesData = [
        'function hello() {',
        '  return "Hello World";',
        '}',
        '',
        'function goodbye() {',
        '  return "Goodbye";',
        '}',
        'export { hello, goodbye };'
    ]

    function createSearch(linesData: string[] | undefined = mockLinesData, startLineValue = 0) {
        const lines = ref(linesData)
        const startLine = ref(startLineValue)
        const scrollToLine = vi.fn()

        const search = useContextSearch({
            lines,
            startLine,
            scrollToLine
        })

        return { search, lines, startLine, scrollToLine }
    }

    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('should initialize with default state', () => {
        const { search } = createSearch()

        expect(search.showSearch.value).toBe(false)
        expect(search.searchQuery.value).toBe('')
        expect(search.searchResults.value).toEqual([])
        expect(search.currentSearchIndex.value).toBe(0)
    })

    it('should toggle search visibility', () => {
        const { search } = createSearch()

        expect(search.showSearch.value).toBe(false)

        search.toggleSearch()
        expect(search.showSearch.value).toBe(true)

        search.toggleSearch()
        expect(search.showSearch.value).toBe(false)
    })

    it('should close search', async () => {
        const { search } = createSearch()

        search.showSearch.value = true
        search.searchQuery.value = 'test'

        search.closeSearch()
        await nextTick()

        expect(search.showSearch.value).toBe(false)
        // Query is cleared by the watch on showSearch
    })

    it('should find matching lines when searching', async () => {
        const { search } = createSearch()

        search.searchQuery.value = 'function'
        await nextTick()

        expect(search.searchResults.value).toEqual([0, 4])
        expect(search.currentSearchIndex.value).toBe(0)
    })

    it('should perform case-insensitive search', async () => {
        const { search } = createSearch()

        search.searchQuery.value = 'HELLO'
        await nextTick()

        // 'hello' appears in lines 0, 1, and 7 (export { hello, goodbye })
        expect(search.searchResults.value).toEqual([0, 1, 7])
    })

    it('should clear results when query is empty', async () => {
        const { search } = createSearch()

        search.searchQuery.value = 'function'
        await nextTick()
        expect(search.searchResults.value.length).toBeGreaterThan(0)

        search.searchQuery.value = ''
        await nextTick()
        expect(search.searchResults.value).toEqual([])
        expect(search.currentSearchIndex.value).toBe(0)
    })

    it('should navigate to next search result', async () => {
        const { search } = createSearch()

        search.searchQuery.value = 'function'
        await nextTick()

        expect(search.currentSearchIndex.value).toBe(0)

        search.searchNext()
        expect(search.currentSearchIndex.value).toBe(1)

        // Should wrap around
        search.searchNext()
        expect(search.currentSearchIndex.value).toBe(0)
    })

    it('should navigate to previous search result', async () => {
        const { search } = createSearch()

        search.searchQuery.value = 'function'
        await nextTick()

        expect(search.currentSearchIndex.value).toBe(0)

        // Should wrap to last result
        search.searchPrev()
        expect(search.currentSearchIndex.value).toBe(1)

        search.searchPrev()
        expect(search.currentSearchIndex.value).toBe(0)
    })

    it('should not navigate when no results', async () => {
        const { search } = createSearch()

        search.searchQuery.value = 'nonexistent'
        await nextTick()

        expect(search.searchResults.value).toEqual([])

        search.searchNext()
        expect(search.currentSearchIndex.value).toBe(0)

        search.searchPrev()
        expect(search.currentSearchIndex.value).toBe(0)
    })

    it('should highlight current search result line', async () => {
        const { search } = createSearch()

        search.searchQuery.value = 'function'
        await nextTick()

        // First result (line 0) should be highlighted
        expect(search.highlightedLinesSet.value.has(0)).toBe(true)
        expect(search.highlightedLinesSet.value.has(4)).toBe(false)

        // Navigate to next
        search.searchNext()
        expect(search.highlightedLinesSet.value.has(0)).toBe(false)
        expect(search.highlightedLinesSet.value.has(4)).toBe(true)
    })

    it('should handle undefined lines', async () => {
        const { search } = createSearch(undefined)

        search.searchQuery.value = 'test'
        await nextTick()

        expect(search.searchResults.value).toEqual([])
    })

    it('should handle empty lines array', async () => {
        const { search } = createSearch([])

        search.searchQuery.value = 'test'
        await nextTick()

        expect(search.searchResults.value).toEqual([])
    })

    it('should calculate correct line numbers with startLine offset', async () => {
        const { search } = createSearch(['first line', 'second line', 'third line'], 100)

        search.searchQuery.value = 'line'
        await nextTick()

        expect(search.searchResults.value).toEqual([100, 101, 102])
    })

    it('should clear query when search is closed', async () => {
        const { search } = createSearch()

        search.showSearch.value = true
        search.searchQuery.value = 'function'
        await nextTick()

        search.showSearch.value = false
        await nextTick()

        expect(search.searchQuery.value).toBe('')
    })

    it('should call scrollToLine when navigating results', async () => {
        const { search, scrollToLine } = createSearch()

        search.searchQuery.value = 'function'
        await nextTick()
        await nextTick() // Wait for scroll

        expect(scrollToLine).toHaveBeenCalledWith(0)

        search.searchNext()
        await nextTick()

        expect(scrollToLine).toHaveBeenCalledWith(4)
    })

    it('should return empty highlightedLinesSet when no results', () => {
        const { search } = createSearch()

        expect(search.highlightedLinesSet.value.size).toBe(0)
    })
})
