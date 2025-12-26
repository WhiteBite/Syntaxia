/**
 * useFileFilter - File filtering by extensions
 * Manages include/exclude extension filters and sorting
 */

import { useSettingsStore } from '@/stores/settings.store'
import type { FileNode } from '@/types/domain'
import { filterTreeByExtensions } from '@/utils/fileTreeUtils'
import { computed, ref, type Ref } from 'vue'

export interface UseFileFilterOptions {
    nodes: Ref<FileNode[]>
}

/**
 * Sort nodes with folders first, then alphabetically
 * Optimized to avoid unnecessary object creation
 */
function sortFoldersFirst(nodes: FileNode[]): FileNode[] {
    // Check if already sorted (optimization for repeated calls)
    let needsSort = false
    for (let i = 1; i < nodes.length; i++) {
        const prev = nodes[i - 1]
        const curr = nodes[i]
        // Check folder order
        if (!prev.isDir && curr.isDir) {
            needsSort = true
            break
        }
        // Check alphabetical order within same type
        if (prev.isDir === curr.isDir &&
            prev.name.toLowerCase().localeCompare(curr.name.toLowerCase()) > 0) {
            needsSort = true
            break
        }
    }

    const sorted = needsSort
        ? [...nodes].sort((a, b) => {
            if (a.isDir && !b.isDir) return -1
            if (!a.isDir && b.isDir) return 1
            return a.name.toLowerCase().localeCompare(b.name.toLowerCase())
        })
        : nodes

    // Recursively sort children only if they have children
    return sorted.map(node => {
        if (node.isDir && node.children && node.children.length > 1) {
            const sortedChildren = sortFoldersFirst(node.children)
            // Only create new object if children actually changed
            if (sortedChildren !== node.children) {
                return { ...node, children: sortedChildren }
            }
        }
        return node
    })
}

export function useFileFilter(options: UseFileFilterOptions) {
    const { nodes } = options
    const settingsStore = useSettingsStore()

    // State
    const filterExtensions = ref<string[]>([])
    const excludeExtensions = ref<string[]>([])

    // Computed
    const filteredNodes = computed(() => {
        let result = nodes.value

        // Apply extension filters
        if (filterExtensions.value.length > 0 || excludeExtensions.value.length > 0) {
            result = filterTreeByExtensions(
                result,
                filterExtensions.value,
                excludeExtensions.value
            )
        }

        // Apply folders first sorting
        if (settingsStore.settings.fileExplorer.foldersFirst) {
            result = sortFoldersFirst(result)
        }

        return result
    })

    // Actions
    function setFilterExtensions(include: string[], exclude: string[] = []) {
        filterExtensions.value = include
        excludeExtensions.value = exclude
    }

    function clearFilters() {
        filterExtensions.value = []
        excludeExtensions.value = []
    }

    function addIncludeExtension(ext: string) {
        if (!filterExtensions.value.includes(ext)) {
            filterExtensions.value.push(ext)
        }
    }

    function removeIncludeExtension(ext: string) {
        const index = filterExtensions.value.indexOf(ext)
        if (index > -1) {
            filterExtensions.value.splice(index, 1)
        }
    }

    function addExcludeExtension(ext: string) {
        if (!excludeExtensions.value.includes(ext)) {
            excludeExtensions.value.push(ext)
        }
    }

    function removeExcludeExtension(ext: string) {
        const index = excludeExtensions.value.indexOf(ext)
        if (index > -1) {
            excludeExtensions.value.splice(index, 1)
        }
    }

    return {
        // State
        filterExtensions,
        excludeExtensions,
        // Computed
        filteredNodes,
        // Actions
        setFilterExtensions,
        clearFilters,
        addIncludeExtension,
        removeIncludeExtension,
        addExcludeExtension,
        removeExcludeExtension,
    }
}
