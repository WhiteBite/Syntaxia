import { useFilePersistence } from '@/composables/useFilePersistence'
import type { FileNode } from '@/types/domain'
import { beforeEach, describe, expect, it } from 'vitest'
import { ref, shallowRef } from 'vue'

describe('useFilePersistence - Presets', () => {
    let nodes: ReturnType<typeof ref<FileNode[]>>
    let selectedPaths: ReturnType<typeof shallowRef<Set<string>>>
    let rootPath: ReturnType<typeof ref<string>>
    let findNode: (path: string) => FileNode | null

    beforeEach(() => {
        localStorage.clear()
        nodes = ref<FileNode[]>([
            {
                name: 'src',
                path: '/project/src',
                isDir: true,
                isExpanded: false,
                children: [
                    { name: 'index.ts', path: '/project/src/index.ts', isDir: false },
                    { name: 'utils.ts', path: '/project/src/utils.ts', isDir: false },
                ],
            },
        ])
        selectedPaths = shallowRef(new Set<string>())
        rootPath = ref('/project')

        findNode = (path: string): FileNode | null => {
            const findInNodes = (nodeList: FileNode[]): FileNode | null => {
                for (const node of nodeList) {
                    if (node.path === path) return node
                    if (node.children) {
                        const found = findInNodes(node.children)
                        if (found) return found
                    }
                }
                return null
            }
            return findInNodes(nodes.value)
        }
    })

    it('should save and load presets', () => {
        const persistence = useFilePersistence({
            nodes,
            selectedPaths,
            rootPath,
            findNode,
        })

        selectedPaths.value.add('/project/src/index.ts')
        selectedPaths.value.add('/project/src/utils.ts')

        persistence.savePreset('Test Preset', 'Test description')

        const presets = persistence.getPresets()
        expect(presets).toHaveLength(1)
        expect(presets[0].name).toBe('Test Preset')
        expect(presets[0].description).toBe('Test description')
        expect(presets[0].paths).toEqual(['/project/src/index.ts', '/project/src/utils.ts'])
    })

    it('should load preset and restore selection', () => {
        const persistence = useFilePersistence({
            nodes,
            selectedPaths,
            rootPath,
            findNode,
        })

        selectedPaths.value.add('/project/src/index.ts')
        selectedPaths.value.add('/project/src/utils.ts')
        persistence.savePreset('Test Preset')

        selectedPaths.value.clear()
        expect(selectedPaths.value.size).toBe(0)

        const count = persistence.loadPreset('Test Preset')
        expect(count).toBe(2)
        expect(selectedPaths.value.has('/project/src/index.ts')).toBe(true)
        expect(selectedPaths.value.has('/project/src/utils.ts')).toBe(true)
    })

    it('should delete preset', () => {
        const persistence = useFilePersistence({
            nodes,
            selectedPaths,
            rootPath,
            findNode,
        })

        selectedPaths.value.add('/project/src/index.ts')
        persistence.savePreset('Test Preset')

        expect(persistence.getPresets()).toHaveLength(1)

        persistence.deletePreset('Test Preset')
        expect(persistence.getPresets()).toHaveLength(0)
    })

    it('should enforce max presets limit', () => {
        const persistence = useFilePersistence({
            nodes,
            selectedPaths,
            rootPath,
            findNode,
        })

        selectedPaths.value.add('/project/src/index.ts')

        // Save 20 presets (max limit)
        for (let i = 0; i < 20; i++) {
            persistence.savePreset(`Preset ${i}`)
        }

        expect(persistence.getPresets()).toHaveLength(20)

        // Try to save 21st preset
        expect(() => persistence.savePreset('Preset 21')).toThrow('Maximum 20 presets allowed')
    })

    it('should update existing preset with same name', () => {
        const persistence = useFilePersistence({
            nodes,
            selectedPaths,
            rootPath,
            findNode,
        })

        selectedPaths.value.add('/project/src/index.ts')
        persistence.savePreset('Test Preset', 'Original description')

        selectedPaths.value.add('/project/src/utils.ts')
        persistence.savePreset('Test Preset', 'Updated description')

        const presets = persistence.getPresets()
        expect(presets).toHaveLength(1)
        expect(presets[0].description).toBe('Updated description')
        expect(presets[0].paths).toHaveLength(2)
    })

    it('should skip non-existent files when loading preset', () => {
        const persistence = useFilePersistence({
            nodes,
            selectedPaths,
            rootPath,
            findNode,
        })

        selectedPaths.value.add('/project/src/index.ts')
        selectedPaths.value.add('/project/src/deleted.ts') // This file doesn't exist
        persistence.savePreset('Test Preset')

        selectedPaths.value.clear()
        const count = persistence.loadPreset('Test Preset')

        // Only 1 file should be loaded (deleted.ts is skipped)
        expect(count).toBe(1)
        expect(selectedPaths.value.has('/project/src/index.ts')).toBe(true)
        expect(selectedPaths.value.has('/project/src/deleted.ts')).toBe(false)
    })

    it('should clear all presets', () => {
        const persistence = useFilePersistence({
            nodes,
            selectedPaths,
            rootPath,
            findNode,
        })

        selectedPaths.value.add('/project/src/index.ts')
        persistence.savePreset('Preset 1')
        persistence.savePreset('Preset 2')

        expect(persistence.getPresets()).toHaveLength(2)

        persistence.clearAllPresets()
        expect(persistence.getPresets()).toHaveLength(0)
    })
})
