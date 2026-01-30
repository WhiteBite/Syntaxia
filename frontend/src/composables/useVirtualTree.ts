/**
 * useVirtualTree - Flatten tree for virtualization
 * Converts hierarchical tree to flat list of visible nodes
 */

import { useSettingsStore } from '@/stores/settings.store'
import type { FileNode } from '@/types/domain'
import { computed, type Ref } from 'vue'

export interface FlattenedNode {
    id: string // Flat key for RecycleScroller (same as node.path)
    node: FileNode
    depth: number
    isLast: boolean
    ancestorHasMoreSiblings: boolean[]
    displayName?: string // For compact mode (e.g. "src/components")
}

export interface UseVirtualTreeOptions {
    nodes: Ref<FileNode[]>
}

/**
 * Flatten tree to visible nodes only (expanded folders)
 */
export function useVirtualTree(options: UseVirtualTreeOptions) {
    const { nodes } = options
    const settingsStore = useSettingsStore()

    const flattenedVisibleNodes = computed<FlattenedNode[]>(() => {
        const result: FlattenedNode[] = []
        const isCompactEnabled = settingsStore.settings.fileExplorer.compactNestedFolders

        function flatten(
            nodeList: FileNode[],
            depth: number,
            ancestorHasMoreSiblings: boolean[]
        ) {
            nodeList.forEach((node, index) => {
                const isLast = index === nodeList.length - 1
                let displayName = node.name
                let currentNode = node
                
                // Compact mode: merge folders with single subfolder
                if (isCompactEnabled && node.isDir && node.children?.length === 1 && node.children[0].isDir) {
                    let nextNode = node.children[0]
                    // We also need to check if the node is expanded to continue compacting
                    // Usually compacting shows the chain even if some are not expanded, 
                    // but the LAST one in the chain is what we interact with.
                    while (nextNode.isDir && nextNode.children?.length === 1 && nextNode.children[0].isDir) {
                        displayName += '/' + nextNode.name
                        nextNode = nextNode.children[0]
                    }
                    displayName += '/' + nextNode.name
                    currentNode = nextNode
                }

                result.push({
                    id: node.path, // Use original node path as ID
                    node: currentNode,
                    depth,
                    isLast,
                    ancestorHasMoreSiblings: [...ancestorHasMoreSiblings],
                    displayName
                })

                // Only recurse into expanded directories
                // Note: in compact mode, we interact with the deepest folder in the chain
                if (currentNode.isDir && currentNode.isExpanded && currentNode.children?.length) {
                    flatten(
                        currentNode.children,
                        depth + 1,
                        [...ancestorHasMoreSiblings, !isLast]
                    )
                }
            })
        }

        flatten(nodes.value, 0, [])
        return result
    })


    const totalVisibleCount = computed(() => flattenedVisibleNodes.value.length)

    return {
        flattenedVisibleNodes,
        totalVisibleCount,
    }
}
