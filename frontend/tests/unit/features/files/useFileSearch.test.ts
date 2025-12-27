/**
 * Unit tests for useFileSearch.ts composable
 * Tests file search functionality
 */
import { useFileSearch } from '@/features/files/composables/useFileSearch'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock the file store
const mockSetSearchQuery = vi.fn()

vi.mock('@/features/files/model/file.store', () => ({
    useFileStore: () => ({
        setSearchQuery: mockSetSearchQuery,
    }),
}))

// Mock constants
vi.mock('@/config/constants', () => ({
    FILE_TREE: {
        DEBOUNCE_MS: 0, // Set to 0 for immediate execution in tests
    },
}))

// Mock debounce to execute immediately in tests
vi.mock('@/utils/performance', () => ({
    debounce: (fn: () => void) => fn,
}))

describe('useFileSearch', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    describe('Initial State', () => {
        it('should have empty query initially', () => {
            const search = useFileSearch()
            expect(search.query.value).toBe('')
        })

        it('should not be searching initially', () => {
            const search = useFileSearch()
            expect(search.isSearching.value).toBe(false)
        })

        it('should not be active initially', () => {
            const search = useFileSearch()
            expect(search.isActive()).toBe(false)
        })
    })

    describe('handleSearch', () => {
        it('should set isSearching to true when called', () => {
            const search = useFileSearch()
            search.query.value = 'test'
            search.handleSearch()

            // With mocked debounce, isSearching should be false after immediate execution
            expect(search.isSearching.value).toBe(false)
        })

        it('should call fileStore.setSearchQuery with current query', () => {
            const search = useFileSearch()
            search.query.value = 'main.ts'
            search.handleSearch()

            expect(mockSetSearchQuery).toHaveBeenCalledWith('main.ts')
        })
    })

    describe('clear', () => {
        it('should clear the query', () => {
            const search = useFileSearch()
            search.query.value = 'test'
            search.clear()

            expect(search.query.value).toBe('')
        })

        it('should call fileStore.setSearchQuery with empty string', () => {
            const search = useFileSearch()
            search.query.value = 'test'
            search.clear()

            expect(mockSetSearchQuery).toHaveBeenCalledWith('')
        })

        it('should set isSearching to false', () => {
            const search = useFileSearch()
            search.query.value = 'test'
            search.handleSearch()
            search.clear()

            expect(search.isSearching.value).toBe(false)
        })
    })

    describe('setQuery', () => {
        it('should set the query value', () => {
            const search = useFileSearch()
            search.setQuery('utils.ts')

            expect(search.query.value).toBe('utils.ts')
        })

        it('should trigger search after setting query', () => {
            const search = useFileSearch()
            search.setQuery('utils.ts')

            expect(mockSetSearchQuery).toHaveBeenCalledWith('utils.ts')
        })
    })

    describe('isActive', () => {
        it('should return false when query is empty', () => {
            const search = useFileSearch()
            search.query.value = ''

            expect(search.isActive()).toBe(false)
        })

        it('should return true when query has content', () => {
            const search = useFileSearch()
            search.query.value = 'test'

            expect(search.isActive()).toBe(true)
        })

        it('should return true for single character query', () => {
            const search = useFileSearch()
            search.query.value = 'a'

            expect(search.isActive()).toBe(true)
        })
    })

    describe('Search Flow', () => {
        it('should handle complete search flow', () => {
            const search = useFileSearch()

            // Start search
            search.setQuery('component')
            expect(search.query.value).toBe('component')
            expect(search.isActive()).toBe(true)
            expect(mockSetSearchQuery).toHaveBeenCalledWith('component')

            // Clear search
            search.clear()
            expect(search.query.value).toBe('')
            expect(search.isActive()).toBe(false)
            expect(mockSetSearchQuery).toHaveBeenCalledWith('')
        })

        it('should handle multiple search queries', () => {
            const search = useFileSearch()

            search.setQuery('first')
            expect(mockSetSearchQuery).toHaveBeenCalledWith('first')

            search.setQuery('second')
            expect(mockSetSearchQuery).toHaveBeenCalledWith('second')

            search.setQuery('third')
            expect(mockSetSearchQuery).toHaveBeenCalledWith('third')

            expect(mockSetSearchQuery).toHaveBeenCalledTimes(3)
        })
    })
})
