/**
 * useVirtualTree - Flatten tree for virtualization
 * Converts hierarchical tree to flat list of visible nodes
 */

import type { FileNode } from '@/types/domain'
import { computed, type Ref } from 'vue'

export interface FlattenedNode {
    id: string // Flat key for RecycleScroller (same as node.path)
    node: FileNode
    depth: number
    isLast: boolean
    ancestorHasMoreSiblings: boolean[]
}

export interface UseVirtualTreeOptions {
    nodes: Ref<FileNode[]>
}

/**
 * Flatten tree to visible nodes only (expanded folders)
 */
export function useVirtualTree(options: UseVirtualTreeOptions) {
    const { nodes } = options

    const flattenedVisibleNodes = computed<FlattenedNode[]>(() => {
        const result: FlattenedNode[] = []

        function flatten(
            nodeList: FileNode[],
            depth: number,
            ancestorHasMoreSiblings: boolean[]
        ) {
            nodeList.forEach((node, index) => {
                const isLast = index === nodeList.length - 1

                result.push({
                    id: node.path, // Flat key for RecycleScroller
                    node,
                    depth,
                    isLast,
                    ancestorHasMoreSiblings: [...ancestorHasMoreSiblings],
                })

                // Only recurse into expanded directories
                if (node.isDir && node.isExpanded && node.children?.length) {
                    flatten(
                        node.children,
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
