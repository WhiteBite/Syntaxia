/**
 * useFilePersistence - Selection and expanded state persistence
 * Saves/loads state to localStorage
 */

import { useLogger } from '@/composables/useLogger'
import type { FileNode } from '@/types/domain'
import { walkTree } from '@/utils/fileTreeUtils'
import { type Ref, type ShallowRef } from 'vue'

const logger = useLogger('FilePersistence')

export interface UseFilePersistenceOptions {
    nodes: Ref<FileNode[]>
    selectedPaths: ShallowRef<Set<string>>
    rootPath: Ref<string>
    findNode: (path: string) => FileNode | null
}

const SELECTION_PREFIX = 'file-selection-'
const EXPANDED_PREFIX = 'file-expanded-'
const PRESETS_PREFIX = 'file-presets-'
const MAX_SAVED_SELECTIONS = 100
const MAX_PRESETS = 20

export interface SelectionPreset {
    name: string
    description?: string
    paths: string[]
    createdAt: number
}

export function useFilePersistence(options: UseFilePersistenceOptions) {
    const { nodes, selectedPaths, rootPath, findNode } = options

    // Debounce timers
    let saveExpandedStateTimer: ReturnType<typeof setTimeout> | null = null
    let saveSelectionTimer: ReturnType<typeof setTimeout> | null = null

    // Selection persistence
    function saveSelectionToStorage() {
        if (!rootPath.value) return

        try {
            const key = `${SELECTION_PREFIX}${rootPath.value}`
            const selection = Array.from(selectedPaths.value).slice(0, MAX_SAVED_SELECTIONS)
            localStorage.setItem(key, JSON.stringify(selection))
            logger.debug(`Saved selection: ${selection.length} files`)
        } catch (err) {
            logger.warn('Failed to save selection:', err)
        }
    }

    function debouncedSaveSelection(delay: number = 300) {
        if (saveSelectionTimer) {
            clearTimeout(saveSelectionTimer)
        }
        saveSelectionTimer = setTimeout(() => {
            saveSelectionToStorage()
            saveSelectionTimer = null
        }, delay)
    }

    function loadSelectionFromStorage(projectPath: string): string[] {
        try {
            const key = `${SELECTION_PREFIX}${projectPath}`
            const saved = localStorage.getItem(key)
            if (saved) {
                const selection = JSON.parse(saved) as string[]
                logger.debug(`Loaded selection: ${selection.length} files`)
                return selection
            }
        } catch (err) {
            logger.warn('Failed to load selection:', err)
        }
        return []
    }

    function clearSelectionHistory(projectPath?: string) {
        try {
            if (projectPath) {
                const key = `${SELECTION_PREFIX}${projectPath}`
                localStorage.removeItem(key)
            } else {
                const keys = Object.keys(localStorage).filter((k) =>
                    k.startsWith(SELECTION_PREFIX)
                )
                keys.forEach((k) => localStorage.removeItem(k))
            }
        } catch (err) {
            logger.warn('Failed to clear selection history:', err)
        }
    }

    function getSelectionStats(): Record<string, number> {
        const stats: Record<string, number> = {}
        try {
            const keys = Object.keys(localStorage).filter((k) =>
                k.startsWith(SELECTION_PREFIX)
            )
            keys.forEach((key) => {
                const projectPath = key.replace(SELECTION_PREFIX, '')
                const saved = localStorage.getItem(key)
                if (saved) {
                    const selection = JSON.parse(saved) as string[]
                    stats[projectPath] = selection.length
                }
            })
        } catch (err) {
            logger.warn('Failed to get selection stats:', err)
        }
        return stats
    }

    // Expanded state persistence
    function saveExpandedState() {
        if (!rootPath.value) return

        try {
            const expandedPaths: string[] = []
            walkTree(nodes.value, (node) => {
                if (node.isDir && node.isExpanded) {
                    expandedPaths.push(node.path)
                }
            })

            const key = `${EXPANDED_PREFIX}${rootPath.value}`
            localStorage.setItem(key, JSON.stringify(expandedPaths))
        } catch (err) {
            logger.warn('Failed to save expanded state:', err)
        }
    }

    function debouncedSaveExpandedState(delay: number = 500) {
        if (saveExpandedStateTimer) {
            clearTimeout(saveExpandedStateTimer)
        }
        saveExpandedStateTimer = setTimeout(() => {
            saveExpandedState()
            saveExpandedStateTimer = null
        }, delay)
    }

    function loadExpandedState(): string[] {
        if (!rootPath.value) return []

        try {
            const key = `${EXPANDED_PREFIX}${rootPath.value}`
            const saved = localStorage.getItem(key)
            if (saved) {
                const expandedPaths = JSON.parse(saved) as string[]
                expandedPaths.forEach((path) => {
                    const node = findNode(path)
                    if (node && node.isDir) {
                        node.isExpanded = true
                    }
                })
                return expandedPaths
            }
        } catch (err) {
            logger.warn('Failed to load expanded state:', err)
        }
        return []
    }

    function clearExpandedHistory(projectPath?: string) {
        try {
            if (projectPath) {
                const key = `${EXPANDED_PREFIX}${projectPath}`
                localStorage.removeItem(key)
            } else {
                const keys = Object.keys(localStorage).filter((k) =>
                    k.startsWith(EXPANDED_PREFIX)
                )
                keys.forEach((k) => localStorage.removeItem(k))
            }
        } catch (err) {
            logger.warn('Failed to clear expanded history:', err)
        }
    }

    // Preset management
    function getPresetsKey(): string {
        return `${PRESETS_PREFIX}${rootPath.value}`
    }

    function getPresets(): SelectionPreset[] {
        if (!rootPath.value) return []

        try {
            const key = getPresetsKey()
            const saved = localStorage.getItem(key)
            if (saved) {
                const presets = JSON.parse(saved) as SelectionPreset[]
                return presets.sort((a, b) => b.createdAt - a.createdAt)
            }
        } catch (err) {
            logger.warn('Failed to load presets:', err)
        }
        return []
    }

    function savePreset(name: string, description?: string): void {
        if (!rootPath.value) {
            throw new Error('No project loaded')
        }

        const trimmedName = name.trim()
        if (!trimmedName) {
            throw new Error('Preset name is required')
        }

        const presets = getPresets()

        // Check max limit
        if (presets.length >= MAX_PRESETS) {
            throw new Error(`Maximum ${MAX_PRESETS} presets allowed`)
        }

        // Check for duplicate name
        const existingIndex = presets.findIndex((p) => p.name === trimmedName)

        const newPreset: SelectionPreset = {
            name: trimmedName,
            description: description?.trim(),
            paths: Array.from(selectedPaths.value),
            createdAt: Date.now(),
        }

        if (existingIndex >= 0) {
            // Update existing preset
            presets[existingIndex] = newPreset
        } else {
            // Add new preset
            presets.push(newPreset)
        }

        try {
            const key = getPresetsKey()
            localStorage.setItem(key, JSON.stringify(presets))
            logger.debug(`Saved preset: ${trimmedName} (${newPreset.paths.length} files)`)
        } catch (err) {
            logger.error('Failed to save preset:', err)
            throw new Error('Failed to save preset')
        }
    }

    function loadPreset(name: string): number {
        const presets = getPresets()
        const preset = presets.find((p) => p.name === name)

        if (!preset) {
            throw new Error('Preset not found')
        }

        // Validate that files still exist
        const validPaths: string[] = []
        for (const path of preset.paths) {
            const node = findNode(path)
            if (node && !node.isDir) {
                validPaths.push(path)
            }
        }

        if (validPaths.length === 0) {
            logger.warn('No valid files found in preset')
            return 0
        }

        // Clear current selection and load preset
        selectedPaths.value.clear()
        validPaths.forEach((path) => selectedPaths.value.add(path))

        logger.debug(`Loaded preset: ${name} (${validPaths.length}/${preset.paths.length} files)`)
        return validPaths.length
    }

    function deletePreset(name: string): void {
        const presets = getPresets()
        const filteredPresets = presets.filter((p) => p.name !== name)

        if (filteredPresets.length === presets.length) {
            throw new Error('Preset not found')
        }

        try {
            const key = getPresetsKey()
            localStorage.setItem(key, JSON.stringify(filteredPresets))
            logger.debug(`Deleted preset: ${name}`)
        } catch (err) {
            logger.error('Failed to delete preset:', err)
            throw new Error('Failed to delete preset')
        }
    }

    function clearAllPresets(): void {
        if (!rootPath.value) return

        try {
            const key = getPresetsKey()
            localStorage.removeItem(key)
            logger.debug('Cleared all presets')
        } catch (err) {
            logger.warn('Failed to clear presets:', err)
        }
    }

    // Cleanup
    function dispose() {
        if (saveExpandedStateTimer) {
            clearTimeout(saveExpandedStateTimer)
            saveExpandedStateTimer = null
        }
        if (saveSelectionTimer) {
            clearTimeout(saveSelectionTimer)
            saveSelectionTimer = null
        }
    }

    return {
        // Selection persistence
        saveSelectionToStorage,
        debouncedSaveSelection,
        loadSelectionFromStorage,
        clearSelectionHistory,
        getSelectionStats,
        // Expanded state persistence
        saveExpandedState,
        debouncedSaveExpandedState,
        loadExpandedState,
        clearExpandedHistory,
        // Preset management
        getPresets,
        savePreset,
        loadPreset,
        deletePreset,
        clearAllPresets,
        // Cleanup
        dispose,
    }
}
