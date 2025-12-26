import { useContextSearch } from '@/features/context/composables/useContextSearch'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick } from 'vue'

describe('useContextSearch', () => {
    const mockLines = {
        lines: [
            'function hello() {',
            '  return "Hello World";',
            '}',
            '',
            'function goodbye() {',
            '  return "Goodbye";',
            '}',
            'export { hello, goodbye };'
        ],
        startLine: 1
    }

    let getLinesMock: () => { lines: string[]; startLine: number } | null

    beforeEach(() => {
        vi.clearAllMocks()
        getLinesMock = () => mockLines
    })

    it('should initialize with default state', () => {
        const search = useContextSearch(getLinesMock)

        expect(search.showSearch.value).toBe(false)
        expect(search.searchQuery.value).toBe('')
        expect(search.searchResults.value).toEqual([])
        expect(search.currentSearchIndex.value).toBe(0)
    })

    it('should toggle search visibility', () => {
        const search = useContextSearch(getLinesMock)

        expect(search.showSearch.value).toBe(false)

        search.toggleSearch()
        expect(search.showSearch.value).toBe(true)

        search.toggleSearch()
        expect(search.showSearch.value).toBe(false)
    })

    it('should close search', async () => {
        const search = useContextSearch(getLinesMock)

        search.showSearch.value = true
        search.searchQuery.value = 'test'

        search.closeSearch()
        await nextTick()

        expect(search.showSearch.value).toBe(false)
        // Query is cleared by the watch on showSearch
    })

    it('should find matching lines when searching', async () => {
        const search = useContextSearch(getLinesMock)

        search.searchQuery.value = 'function'
        await nextTick()

        expect(search.searchResults.value).toEqual([1, 5])
        expect(search.currentSearchIndex.value).toBe(0)
    })

    it('should perform case-insensitive search', async () => {
        const search = useContextSearch(getLinesMock)

        search.searchQuery.value = 'HELLO'
        await nextTick()

        // 'hello' appears in lines 1, 2, and 8 (export { hello, goodbye })
        expect(search.searchResults.value).toEqual([1, 2, 8])
    })

    it('should clear results when query is empty', async () => {
        const search = useContextSearch(getLinesMock)

        search.searchQuery.value = 'function'
        await nextTick()
        expect(search.searchResults.value.length).toBeGreaterThan(0)

        search.searchQuery.value = ''
        await nextTick()
        expect(search.searchResults.value).toEqual([])
        expect(search.currentSearchIndex.value).toBe(0)
    })

    it('should navigate to next search result', async () => {
        const search = useContextSearch(getLinesMock)

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
        const search = useContextSearch(getLinesMock)

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
        const search = useContextSearch(getLinesMock)

        search.searchQuery.value = 'nonexistent'
        await nextTick()

        expect(search.searchResults.value).toEqual([])

        search.searchNext()
        expect(search.currentSearchIndex.value).toBe(0)

        search.searchPrev()
        expect(search.currentSearchIndex.value).toBe(0)
    })

    it('should highlight current search result line', async () => {
        const search = useContextSearch(getLinesMock)

        search.searchQuery.value = 'function'
        await nextTick()

        // First result (line 1) should be highlighted
        expect(search.isLineHighlighted(1)).toBe(true)
        expect(search.isLineHighlighted(5)).toBe(false)

        // Navigate to next
        search.searchNext()
        expect(search.isLineHighlighted(1)).toBe(false)
        expect(search.isLineHighlighted(5)).toBe(true)
    })

    it('should handle null getLines', async () => {
        const search = useContextSearch(() => null)

        search.searchQuery.value = 'test'
        await nextTick()

        expect(search.searchResults.value).toEqual([])
    })

    it('should handle empty lines array', async () => {
        const search = useContextSearch(() => ({ lines: [], startLine: 1 }))

        search.searchQuery.value = 'test'
        await nextTick()

        expect(search.searchResults.value).toEqual([])
    })

    it('should calculate correct line numbers with startLine offset', async () => {
        const offsetLines = {
            lines: ['first line', 'second line', 'third line'],
            startLine: 100
        }
        const search = useContextSearch(() => offsetLines)

        search.searchQuery.value = 'line'
        await nextTick()

        expect(search.searchResults.value).toEqual([100, 101, 102])
    })

    it('should set line refs correctly', () => {
        const search = useContextSearch(getLinesMock)
        const mockElement = document.createElement('div')

        search.setLineRef(mockElement, 5)

        // Line ref should be stored (internal state)
        expect(search.isLineHighlighted(5)).toBe(false) // Not highlighted until search
    })

    it('should not set line ref for null element', () => {
        const search = useContextSearch(getLinesMock)

        // Should not throw
        search.setLineRef(null, 5)
    })

    it('should not set line ref for null line number', () => {
        const search = useContextSearch(getLinesMock)
        const mockElement = document.createElement('div')

        // Should not throw
        search.setLineRef(mockElement, null)
    })

    it('should clear query when search is closed', async () => {
        const search = useContextSearch(getLinesMock)

        search.showSearch.value = true
        search.searchQuery.value = 'function'
        await nextTick()

        search.showSearch.value = false
        await nextTick()

        expect(search.searchQuery.value).toBe('')
    })
})
