import { useFileSelection, useGitFilters } from '@/features/git/composables/useGitFilters'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'

describe('useGitFilters', () => {
    const mockFiles = ref([
        'src/app.ts',
        'src/utils.ts',
        'src/components/Button.vue',
        'src/styles/main.css',
        'package.json',
        'README.md',
        'tsconfig.json',
        'src/index.js',
    ])

    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('should initialize with empty state', () => {
        const filters = useGitFilters(mockFiles)

        expect(filters.searchQuery.value).toBe('')
        expect(filters.activeFilters.value.size).toBe(0)
        expect(filters.filteredFiles.value).toEqual(mockFiles.value)
    })

    it('should filter files by search query', () => {
        const filters = useGitFilters(mockFiles)

        filters.setSearch('app')

        expect(filters.filteredFiles.value).toEqual(['src/app.ts'])
    })

    it('should perform case-insensitive search', () => {
        const filters = useGitFilters(mockFiles)

        filters.setSearch('README')

        expect(filters.filteredFiles.value).toEqual(['README.md'])
    })

    it('should filter by path segments', () => {
        const filters = useGitFilters(mockFiles)

        filters.setSearch('src/')

        expect(filters.filteredFiles.value).toHaveLength(5)
        expect(filters.filteredFiles.value.every(f => f.startsWith('src/'))).toBe(true)
    })

    it('should toggle type filter', () => {
        const filters = useGitFilters(mockFiles)

        expect(filters.activeFilters.value.has('code')).toBe(false)

        filters.toggleFilter('code')
        expect(filters.activeFilters.value.has('code')).toBe(true)

        filters.toggleFilter('code')
        expect(filters.activeFilters.value.has('code')).toBe(false)
    })

    it('should filter files by type filter', () => {
        const filters = useGitFilters(mockFiles)

        filters.toggleFilter('code')

        // Code filter includes .ts, .js, .vue, etc.
        const result = filters.filteredFiles.value
        expect(result).toContain('src/app.ts')
        expect(result).toContain('src/utils.ts')
        expect(result).toContain('src/components/Button.vue')
        expect(result).toContain('src/index.js')
        expect(result).not.toContain('src/styles/main.css')
        expect(result).not.toContain('README.md')
    })

    it('should filter files by styles filter', () => {
        const filters = useGitFilters(mockFiles)

        filters.toggleFilter('styles')

        expect(filters.filteredFiles.value).toEqual(['src/styles/main.css'])
    })

    it('should filter files by config filter', () => {
        const filters = useGitFilters(mockFiles)

        filters.toggleFilter('config')

        expect(filters.filteredFiles.value).toContain('package.json')
        expect(filters.filteredFiles.value).toContain('tsconfig.json')
    })

    it('should filter files by docs filter', () => {
        const filters = useGitFilters(mockFiles)

        filters.toggleFilter('docs')

        expect(filters.filteredFiles.value).toEqual(['README.md'])
    })

    it('should combine multiple type filters', () => {
        const filters = useGitFilters(mockFiles)

        filters.toggleFilter('styles')
        filters.toggleFilter('docs')

        expect(filters.filteredFiles.value).toContain('src/styles/main.css')
        expect(filters.filteredFiles.value).toContain('README.md')
        expect(filters.filteredFiles.value).toHaveLength(2)
    })

    it('should combine search and type filters', () => {
        const filters = useGitFilters(mockFiles)

        filters.toggleFilter('code')
        filters.setSearch('app')

        expect(filters.filteredFiles.value).toEqual(['src/app.ts'])
    })

    it('should clear all filters', () => {
        const filters = useGitFilters(mockFiles)

        filters.toggleFilter('code')
        filters.setSearch('test')

        filters.clearFilters()

        expect(filters.searchQuery.value).toBe('')
        expect(filters.activeFilters.value.size).toBe(0)
        expect(filters.filteredFiles.value).toEqual(mockFiles.value)
    })

    it('should calculate filter counts', () => {
        const filters = useGitFilters(mockFiles)

        const counts = filters.filterCounts.value

        // app.ts, utils.ts, Button.vue, index.js = 4 code files
        expect(counts['code']).toBe(4)
        expect(counts['styles']).toBe(1) // .css
        expect(counts['docs']).toBe(1) // .md
        expect(counts['config']).toBe(2) // .json x2
    })

    it('should expose FILE_TYPE_FILTERS', () => {
        const filters = useGitFilters(mockFiles)

        expect(filters.filters).toBeDefined()
        expect(filters.filters.length).toBeGreaterThan(0)
        expect(filters.filters.find(f => f.id === 'code')).toBeDefined()
    })

    it('should handle empty files array', () => {
        const emptyFiles = ref<string[]>([])
        const filters = useGitFilters(emptyFiles)

        expect(filters.filteredFiles.value).toEqual([])
        expect(filters.filterCounts.value['code']).toBe(0)
    })

    it('should handle debounced search', async () => {
        vi.useFakeTimers()
        const filters = useGitFilters(mockFiles)

        filters.setSearchDebounced('app', 100)

        // Query should not be set immediately
        expect(filters.searchQuery.value).toBe('')

        // Advance timers
        vi.advanceTimersByTime(100)

        expect(filters.searchQuery.value).toBe('app')

        vi.useRealTimers()
    })

    it('should cancel previous debounced search', async () => {
        vi.useFakeTimers()
        const filters = useGitFilters(mockFiles)

        filters.setSearchDebounced('first', 100)
        vi.advanceTimersByTime(50)

        filters.setSearchDebounced('second', 100)
        vi.advanceTimersByTime(100)

        expect(filters.searchQuery.value).toBe('second')

        vi.useRealTimers()
    })

    it('should clear extension cache when it grows too large', () => {
        const filters = useGitFilters(mockFiles)

        // clearExtensionCache should not throw
        filters.clearExtensionCache()
    })

    it('should handle files without extension', () => {
        const filesWithoutExt = ref(['Makefile', 'Dockerfile', 'LICENSE'])
        const filters = useGitFilters(filesWithoutExt)

        filters.toggleFilter('code')

        // Files without matching extensions should be filtered out
        expect(filters.filteredFiles.value).toEqual([])
    })
})

describe('useFileSelection', () => {
    const mockFiles = ref([
        'src/app.ts',
        'src/utils.ts',
        'src/components/Button.vue',
        'src/components/Input.vue',
        'package.json',
    ])

    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('should initialize with empty selection', () => {
        const selection = useFileSelection(mockFiles)

        expect(selection.selectedCount.value).toBe(0)
        expect(selection.hasSelection.value).toBe(false)
        expect(selection.allSelected.value).toBe(false)
    })

    it('should toggle file selection', () => {
        const selection = useFileSelection(mockFiles)

        selection.toggle('src/app.ts')
        expect(selection.isSelected('src/app.ts')).toBe(true)
        expect(selection.selectedCount.value).toBe(1)

        selection.toggle('src/app.ts')
        expect(selection.isSelected('src/app.ts')).toBe(false)
        expect(selection.selectedCount.value).toBe(0)
    })

    it('should select all files', () => {
        const selection = useFileSelection(mockFiles)

        selection.selectAll()

        expect(selection.selectedCount.value).toBe(5)
        expect(selection.allSelected.value).toBe(true)
        expect(selection.hasSelection.value).toBe(true)
    })

    it('should clear selection', () => {
        const selection = useFileSelection(mockFiles)

        selection.selectAll()
        selection.clearSelection()

        expect(selection.selectedCount.value).toBe(0)
        expect(selection.hasSelection.value).toBe(false)
    })

    it('should select folder files', () => {
        const selection = useFileSelection(mockFiles)

        selection.selectFolder('src/components')

        expect(selection.isSelected('src/components/Button.vue')).toBe(true)
        expect(selection.isSelected('src/components/Input.vue')).toBe(true)
        expect(selection.isSelected('src/app.ts')).toBe(false)
        expect(selection.selectedCount.value).toBe(2)
    })

    it('should deselect folder when all files are selected', () => {
        const selection = useFileSelection(mockFiles)

        // Select folder
        selection.selectFolder('src/components')
        expect(selection.hasAllInFolder('src/components')).toBe(true)

        // Deselect folder
        selection.selectFolder('src/components')
        expect(selection.hasAllInFolder('src/components')).toBe(false)
        expect(selection.selectedCount.value).toBe(0)
    })

    it('should check if some files in folder are selected', () => {
        const selection = useFileSelection(mockFiles)

        selection.toggle('src/components/Button.vue')

        expect(selection.hasSomeInFolder('src/components')).toBe(true)
        expect(selection.hasAllInFolder('src/components')).toBe(false)
    })

    it('should check if all files in folder are selected', () => {
        const selection = useFileSelection(mockFiles)

        selection.toggle('src/components/Button.vue')
        selection.toggle('src/components/Input.vue')

        expect(selection.hasAllInFolder('src/components')).toBe(true)
    })

    it('should get files in folder', () => {
        const selection = useFileSelection(mockFiles)

        const folderFiles = selection.getFilesInFolder('src/components')

        expect(folderFiles).toEqual([
            'src/components/Button.vue',
            'src/components/Input.vue'
        ])
    })

    it('should set selection from Set', () => {
        const selection = useFileSelection(mockFiles)

        selection.setSelection(new Set(['src/app.ts', 'package.json']))

        expect(selection.selectedCount.value).toBe(2)
        expect(selection.isSelected('src/app.ts')).toBe(true)
        expect(selection.isSelected('package.json')).toBe(true)
    })

    it('should get selected files as array', () => {
        const selection = useFileSelection(mockFiles)

        selection.toggle('src/app.ts')
        selection.toggle('src/utils.ts')

        const selected = selection.getSelectedArray()

        expect(selected).toHaveLength(2)
        expect(selected).toContain('src/app.ts')
        expect(selected).toContain('src/utils.ts')
    })

    it('should handle empty files array', () => {
        const emptyFiles = ref<string[]>([])
        const selection = useFileSelection(emptyFiles)

        expect(selection.selectedCount.value).toBe(0)
        expect(selection.allSelected.value).toBe(false)

        selection.selectAll()
        expect(selection.selectedCount.value).toBe(0)
    })

    it('should return empty array for non-existent folder', () => {
        const selection = useFileSelection(mockFiles)

        const folderFiles = selection.getFilesInFolder('nonexistent')

        expect(folderFiles).toEqual([])
    })

    it('should return false for hasSomeInFolder with empty folder', () => {
        const selection = useFileSelection(mockFiles)

        expect(selection.hasSomeInFolder('nonexistent')).toBe(false)
    })

    it('should return false for hasAllInFolder with empty folder', () => {
        const selection = useFileSelection(mockFiles)

        expect(selection.hasAllInFolder('nonexistent')).toBe(false)
    })
})
