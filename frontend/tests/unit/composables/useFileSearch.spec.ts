import { useFileSearch } from '@/features/files/composables/useFileSearch'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock file store
const mockSetSearchQuery = vi.fn()
vi.mock('@/features/files/model/file.store', () => ({
    useFileStore: () => ({
        setSearchQuery: mockSetSearchQuery,
    }),
}))

// Mock constants
vi.mock('@/config/constants', () => ({
    FILE_TREE: {
        DEBOUNCE_MS: 0, // No debounce in tests
    },
}))

// Mock debounce to execute immediately
vi.mock('@/utils/performance', () => ({
    debounce: (fn: () => void) => fn,
}))

describe('useFileSearch', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('should initialize with empty query', () => {
        const search = useFileSearch()
        expect(search.query.value).toBe('')
        expect(search.isSearching.value).toBe(false)
    })

    it('should update query and call store on handleSearch', () => {
        const search = useFileSearch()
        search.query.value = 'test'
        search.handleSearch()

        expect(mockSetSearchQuery).toHaveBeenCalledWith('test')
    })

    it('should clear query and reset store', () => {
        const search = useFileSearch()
        search.query.value = 'test'
        search.clear()

        expect(search.query.value).toBe('')
        expect(mockSetSearchQuery).toHaveBeenCalledWith('')
    })

    it('should set query programmatically', () => {
        const search = useFileSearch()
        search.setQuery('new query')

        expect(search.query.value).toBe('new query')
    })

    it('should report active state correctly', () => {
        const search = useFileSearch()
        expect(search.isActive()).toBe(false)

        search.query.value = 'test'
        expect(search.isActive()).toBe(true)

        search.clear()
        expect(search.isActive()).toBe(false)
    })
})
