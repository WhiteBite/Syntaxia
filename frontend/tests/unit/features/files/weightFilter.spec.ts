import { useFileStore } from '@/features/files/model/file.store'
import type { FileNode } from '@/types/domain'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

describe('Token Weight Filtering', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    const createFileNode = (path: string, name: string, size: number): FileNode => ({
        path,
        name,
        isDir: false,
        size,
        depth: 0,
    })

    const createDirNode = (path: string, name: string, children: FileNode[]): FileNode => ({
        path,
        name,
        isDir: true,
        depth: 0,
        children,
    })

    describe('Weight Filter State Management', () => {
        it('should initialize with no weight filter', () => {
            const store = useFileStore()
            expect(store.weightFilter).toBe('none')
        })

        it('should set weight filter to medium', () => {
            const store = useFileStore()
            store.setWeightFilter('medium')
            expect(store.weightFilter).toBe('medium')
        })

        it('should set weight filter to heavy', () => {
            const store = useFileStore()
            store.setWeightFilter('heavy')
            expect(store.weightFilter).toBe('heavy')
        })

        it('should set weight filter to critical', () => {
            const store = useFileStore()
            store.setWeightFilter('critical')
            expect(store.weightFilter).toBe('critical')
        })

        it('should clear weight filter', () => {
            const store = useFileStore()
            store.setWeightFilter('heavy')
            store.clearWeightFilter()
            expect(store.weightFilter).toBe('none')
        })
    })

    describe('Weight Filter Behavior', () => {
        it('should filter files by medium weight (10K+ tokens)', () => {
            const store = useFileStore()

            // Create files with different sizes
            // 4 bytes per token, so 10K tokens = 40,000 bytes
            const smallFile = createFileNode('/small.txt', 'small.txt', 20000) // 5K tokens
            const mediumFile = createFileNode('/medium.txt', 'medium.txt', 50000) // 12.5K tokens
            const heavyFile = createFileNode('/heavy.txt', 'heavy.txt', 250000) // 62.5K tokens

            store.setFileTree([smallFile, mediumFile, heavyFile])
            store.setWeightFilter('medium')

            const filtered = store.filteredNodes
            expect(filtered).toHaveLength(2)
            expect(filtered.find(n => n.name === 'small.txt')).toBeUndefined()
            expect(filtered.find(n => n.name === 'medium.txt')).toBeDefined()
            expect(filtered.find(n => n.name === 'heavy.txt')).toBeDefined()
        })

        it('should filter files by heavy weight (50K+ tokens)', () => {
            const store = useFileStore()

            const smallFile = createFileNode('/small.txt', 'small.txt', 20000) // 5K tokens
            const mediumFile = createFileNode('/medium.txt', 'medium.txt', 50000) // 12.5K tokens
            const heavyFile = createFileNode('/heavy.txt', 'heavy.txt', 250000) // 62.5K tokens
            const criticalFile = createFileNode('/critical.txt', 'critical.txt', 500000) // 125K tokens

            store.setFileTree([smallFile, mediumFile, heavyFile, criticalFile])
            store.setWeightFilter('heavy')

            const filtered = store.filteredNodes
            expect(filtered).toHaveLength(2)
            expect(filtered.find(n => n.name === 'small.txt')).toBeUndefined()
            expect(filtered.find(n => n.name === 'medium.txt')).toBeUndefined()
            expect(filtered.find(n => n.name === 'heavy.txt')).toBeDefined()
            expect(filtered.find(n => n.name === 'critical.txt')).toBeDefined()
        })

        it('should filter files by critical weight (100K+ tokens)', () => {
            const store = useFileStore()

            const mediumFile = createFileNode('/medium.txt', 'medium.txt', 50000) // 12.5K tokens
            const heavyFile = createFileNode('/heavy.txt', 'heavy.txt', 250000) // 62.5K tokens
            const criticalFile = createFileNode('/critical.txt', 'critical.txt', 500000) // 125K tokens

            store.setFileTree([mediumFile, heavyFile, criticalFile])
            store.setWeightFilter('critical')

            const filtered = store.filteredNodes
            expect(filtered).toHaveLength(1)
            expect(filtered.find(n => n.name === 'medium.txt')).toBeUndefined()
            expect(filtered.find(n => n.name === 'heavy.txt')).toBeUndefined()
            expect(filtered.find(n => n.name === 'critical.txt')).toBeDefined()
        })

        it('should show all files when filter is none', () => {
            const store = useFileStore()

            const smallFile = createFileNode('/small.txt', 'small.txt', 20000)
            const mediumFile = createFileNode('/medium.txt', 'medium.txt', 50000)
            const heavyFile = createFileNode('/heavy.txt', 'heavy.txt', 250000)

            store.setFileTree([smallFile, mediumFile, heavyFile])
            store.setWeightFilter('none')

            const filtered = store.filteredNodes
            expect(filtered).toHaveLength(3)
        })
    })

    describe('Weight Filter with Directories', () => {
        it('should include directories that contain files above threshold', () => {
            const store = useFileStore()

            const smallFile = createFileNode('/dir/small.txt', 'small.txt', 20000) // 5K tokens
            const heavyFile = createFileNode('/dir/heavy.txt', 'heavy.txt', 250000) // 62.5K tokens
            const dir = createDirNode('/dir', 'dir', [smallFile, heavyFile])

            store.setFileTree([dir])
            store.setWeightFilter('heavy')

            const filtered = store.filteredNodes
            expect(filtered).toHaveLength(1)
            expect(filtered[0].isDir).toBe(true)
            expect(filtered[0].children).toHaveLength(1)
            expect(filtered[0].children![0].name).toBe('heavy.txt')
        })

        it('should exclude directories with no files above threshold', () => {
            const store = useFileStore()

            const smallFile1 = createFileNode('/dir/small1.txt', 'small1.txt', 20000)
            const smallFile2 = createFileNode('/dir/small2.txt', 'small2.txt', 30000)
            const dir = createDirNode('/dir', 'dir', [smallFile1, smallFile2])

            store.setFileTree([dir])
            store.setWeightFilter('heavy')

            const filtered = store.filteredNodes
            expect(filtered).toHaveLength(0)
        })

        it('should filter nested directories correctly', () => {
            const store = useFileStore()

            const smallFile = createFileNode('/dir/subdir/small.txt', 'small.txt', 20000)
            const criticalFile = createFileNode('/dir/subdir/critical.txt', 'critical.txt', 500000)
            const subdir = createDirNode('/dir/subdir', 'subdir', [smallFile, criticalFile])
            const dir = createDirNode('/dir', 'dir', [subdir])

            store.setFileTree([dir])
            store.setWeightFilter('critical')

            const filtered = store.filteredNodes
            expect(filtered).toHaveLength(1)
            expect(filtered[0].isDir).toBe(true)
            expect(filtered[0].children).toHaveLength(1)
            expect(filtered[0].children![0].isDir).toBe(true)
            expect(filtered[0].children![0].children).toHaveLength(1)
            expect(filtered[0].children![0].children![0].name).toBe('critical.txt')
        })
    })

    describe('Weight Filter Toggle Behavior', () => {
        it('should toggle filter on and off when clicking same level', () => {
            const store = useFileStore()

            const mediumFile = createFileNode('/medium.txt', 'medium.txt', 50000)
            store.setFileTree([mediumFile])

            // First click: activate filter
            store.setWeightFilter('medium')
            expect(store.weightFilter).toBe('medium')
            expect(store.filteredNodes).toHaveLength(1)

            // Second click: clear filter (toggle off)
            store.clearWeightFilter()
            expect(store.weightFilter).toBe('none')
            expect(store.filteredNodes).toHaveLength(1)
        })

        it('should switch between different weight levels', () => {
            const store = useFileStore()

            const mediumFile = createFileNode('/medium.txt', 'medium.txt', 50000)
            const heavyFile = createFileNode('/heavy.txt', 'heavy.txt', 250000)
            const criticalFile = createFileNode('/critical.txt', 'critical.txt', 500000)

            store.setFileTree([mediumFile, heavyFile, criticalFile])

            // Set to medium
            store.setWeightFilter('medium')
            expect(store.filteredNodes).toHaveLength(3)

            // Switch to heavy
            store.setWeightFilter('heavy')
            expect(store.filteredNodes).toHaveLength(2)

            // Switch to critical
            store.setWeightFilter('critical')
            expect(store.filteredNodes).toHaveLength(1)
        })
    })

    describe('Weight Filter with Other Filters', () => {
        it('should work together with extension filters', () => {
            const store = useFileStore()

            const mediumTxt = createFileNode('/medium.txt', 'medium.txt', 50000)
            const mediumJs = createFileNode('/medium.js', 'medium.js', 50000)
            const heavyTxt = createFileNode('/heavy.txt', 'heavy.txt', 250000)

            store.setFileTree([mediumTxt, mediumJs, heavyTxt])

            // Apply both weight and extension filters
            store.setWeightFilter('medium')
            store.setFilterExtensions(['.txt'], [])

            const filtered = store.filteredNodes
            expect(filtered).toHaveLength(2)
            expect(filtered.find(n => n.name === 'medium.txt')).toBeDefined()
            expect(filtered.find(n => n.name === 'medium.js')).toBeUndefined()
            expect(filtered.find(n => n.name === 'heavy.txt')).toBeDefined()
        })
    })

    describe('Edge Cases', () => {
        it('should handle files with no size', () => {
            const store = useFileStore()

            const fileNoSize = createFileNode('/nosize.txt', 'nosize.txt', 0)
            const heavyFile = createFileNode('/heavy.txt', 'heavy.txt', 250000)

            store.setFileTree([fileNoSize, heavyFile])
            store.setWeightFilter('heavy')

            const filtered = store.filteredNodes
            expect(filtered).toHaveLength(1)
            expect(filtered[0].name).toBe('heavy.txt')
        })

        it('should handle empty file tree', () => {
            const store = useFileStore()

            store.setFileTree([])
            store.setWeightFilter('medium')

            expect(store.filteredNodes).toHaveLength(0)
        })

        it('should handle files at exact threshold boundaries', () => {
            const store = useFileStore()

            // Exactly 10K tokens = 40,000 bytes
            const exactMedium = createFileNode('/exact-medium.txt', 'exact-medium.txt', 40000)
            // Exactly 50K tokens = 200,000 bytes
            const exactHeavy = createFileNode('/exact-heavy.txt', 'exact-heavy.txt', 200000)
            // Exactly 100K tokens = 400,000 bytes
            const exactCritical = createFileNode('/exact-critical.txt', 'exact-critical.txt', 400000)

            store.setFileTree([exactMedium, exactHeavy, exactCritical])

            // Test medium threshold (should include all)
            store.setWeightFilter('medium')
            expect(store.filteredNodes).toHaveLength(3)

            // Test heavy threshold (should include heavy and critical)
            store.setWeightFilter('heavy')
            expect(store.filteredNodes).toHaveLength(2)

            // Test critical threshold (should include only critical)
            store.setWeightFilter('critical')
            expect(store.filteredNodes).toHaveLength(1)
        })
    })
})
