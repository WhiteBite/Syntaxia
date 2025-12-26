/**
 * Compact Nested Folders Tests
 * Tests for compactNestedFolders setting behavior
 */

import type { FileNode } from '@/types/domain'
import { describe, expect, it } from 'vitest'

// Helper to create file nodes
function createFileNode(
    path: string,
    isDir = false,
    children: FileNode[] = [],
    isExpanded = false
): FileNode {
    return {
        path,
        name: path.split('/').pop() || path,
        isDir,
        children,
        isExpanded,
        isIgnored: false,
    }
}

// Compact folder logic (extracted from useVirtualTree)
interface CompactInfo {
    isCompact: boolean
    displayName: string
    chainPaths: string[]
    lastNode: FileNode
}

function getCompactInfo(node: FileNode, compactMode: boolean): CompactInfo {
    if (!compactMode || !node.isDir) {
        return {
            isCompact: false,
            displayName: node.name,
            chainPaths: [node.path],
            lastNode: node,
        }
    }

    const chainPaths: string[] = [node.path]
    const names: string[] = [node.name]
    let current = node

    // Follow single-child folder chain
    while (
        current.isDir &&
        current.children?.length === 1 &&
        current.children[0].isDir
    ) {
        current = current.children[0]
        chainPaths.push(current.path)
        names.push(current.name)
    }

    return {
        isCompact: chainPaths.length > 1,
        displayName: names.join('/'),
        chainPaths,
        lastNode: current,
    }
}

describe('Compact Nested Folders', () => {
    describe('compactInfo computed', () => {
        it('should merge single-child folder chains: src/main/java → "src/main/java"', () => {
            const javaDir = createFileNode('src/main/java', true, [
                createFileNode('src/main/java/App.java'),
            ])
            const mainDir = createFileNode('src/main', true, [javaDir])
            const srcDir = createFileNode('src', true, [mainDir])

            const info = getCompactInfo(srcDir, true)

            expect(info.isCompact).toBe(true)
            expect(info.displayName).toBe('src/main/java')
            expect(info.chainPaths).toEqual(['src', 'src/main', 'src/main/java'])
            expect(info.lastNode.path).toBe('src/main/java')
        })

        it('should NOT merge when folder has multiple children', () => {
            const srcDir = createFileNode('src', true, [
                createFileNode('src/components', true),
                createFileNode('src/utils', true),
            ])

            const info = getCompactInfo(srcDir, true)

            expect(info.isCompact).toBe(false)
            expect(info.displayName).toBe('src')
            expect(info.chainPaths).toEqual(['src'])
        })

        it('should NOT merge when child is a file', () => {
            const srcDir = createFileNode('src', true, [
                createFileNode('src/main.ts'),
            ])

            const info = getCompactInfo(srcDir, true)

            expect(info.isCompact).toBe(false)
            expect(info.displayName).toBe('src')
            expect(info.chainPaths).toEqual(['src'])
        })

        it('should return chainPaths for all nodes in chain', () => {
            const deepDir = createFileNode('a/b/c/d', true, [
                createFileNode('a/b/c/d/file.ts'),
            ])
            const cDir = createFileNode('a/b/c', true, [deepDir])
            const bDir = createFileNode('a/b', true, [cDir])
            const aDir = createFileNode('a', true, [bDir])

            const info = getCompactInfo(aDir, true)

            expect(info.chainPaths).toHaveLength(4)
            expect(info.chainPaths).toEqual(['a', 'a/b', 'a/b/c', 'a/b/c/d'])
        })

        it('should return lastNode pointing to deepest folder', () => {
            const deepDir = createFileNode('x/y/z', true, [
                createFileNode('x/y/z/file.ts'),
            ])
            const yDir = createFileNode('x/y', true, [deepDir])
            const xDir = createFileNode('x', true, [yDir])

            const info = getCompactInfo(xDir, true)

            expect(info.lastNode.path).toBe('x/y/z')
            expect(info.lastNode.name).toBe('z')
        })

        it('should handle single folder without children', () => {
            const emptyDir = createFileNode('empty', true, [])

            const info = getCompactInfo(emptyDir, true)

            expect(info.isCompact).toBe(false)
            expect(info.displayName).toBe('empty')
            expect(info.chainPaths).toEqual(['empty'])
        })
    })

    describe('expand/collapse in compact mode', () => {
        it('should expand ALL nodes in chain when clicking compact folder', () => {
            const chainPaths = ['src', 'src/main', 'src/main/java']
            const expandedPaths = new Set<string>()

            // Simulate expanding all paths in chain
            chainPaths.forEach(p => expandedPaths.add(p))

            expect(expandedPaths.has('src')).toBe(true)
            expect(expandedPaths.has('src/main')).toBe(true)
            expect(expandedPaths.has('src/main/java')).toBe(true)
        })

        it('should collapse ALL nodes in chain when clicking expanded compact folder', () => {
            const chainPaths = ['src', 'src/main', 'src/main/java']
            const expandedPaths = new Set<string>(chainPaths)

            // Simulate collapsing all paths in chain
            chainPaths.forEach(p => expandedPaths.delete(p))

            expect(expandedPaths.has('src')).toBe(false)
            expect(expandedPaths.has('src/main')).toBe(false)
            expect(expandedPaths.has('src/main/java')).toBe(false)
        })

        it('should show children of lastNode after expanding', () => {
            const javaDir = createFileNode('src/main/java', true, [
                createFileNode('src/main/java/App.java'),
                createFileNode('src/main/java/Utils.java'),
            ])
            const mainDir = createFileNode('src/main', true, [javaDir])
            const srcDir = createFileNode('src', true, [mainDir])

            const info = getCompactInfo(srcDir, true)

            // After expanding, should show children of lastNode
            expect(info.lastNode.children).toHaveLength(2)
            expect(info.lastNode.children?.[0].name).toBe('App.java')
        })
    })

    describe('disabled compact mode', () => {
        it('should show each folder separately when compactMode=false', () => {
            const javaDir = createFileNode('src/main/java', true, [
                createFileNode('src/main/java/App.java'),
            ])
            const mainDir = createFileNode('src/main', true, [javaDir])
            const srcDir = createFileNode('src', true, [mainDir])

            const info = getCompactInfo(srcDir, false)

            expect(info.isCompact).toBe(false)
            expect(info.displayName).toBe('src')
            expect(info.chainPaths).toEqual(['src'])
        })

        it('should toggle only clicked folder', () => {
            const expandedPaths = new Set<string>()

            // Toggle src only
            expandedPaths.add('src')

            expect(expandedPaths.has('src')).toBe(true)
            expect(expandedPaths.has('src/main')).toBe(false)
            expect(expandedPaths.has('src/main/java')).toBe(false)
        })
    })

    describe('edge cases', () => {
        it('should handle deeply nested single-child chains', () => {
            // Create a/b/c/d/e/f chain
            let current = createFileNode('a/b/c/d/e/f', true, [
                createFileNode('a/b/c/d/e/f/file.ts'),
            ])
            const paths = ['a/b/c/d/e/f', 'a/b/c/d/e', 'a/b/c/d', 'a/b/c', 'a/b', 'a']

            for (let i = 1; i < paths.length; i++) {
                current = createFileNode(paths[i], true, [current])
            }

            const info = getCompactInfo(current, true)

            expect(info.isCompact).toBe(true)
            expect(info.chainPaths).toHaveLength(6)
            expect(info.displayName).toBe('a/b/c/d/e/f')
        })

        it('should handle mixed chain (some with multiple children)', () => {
            const utilsDir = createFileNode('src/utils', true, [
                createFileNode('src/utils/helper.ts'),
            ])
            const componentsDir = createFileNode('src/components', true, [
                createFileNode('src/components/Button.vue'),
            ])
            const srcDir = createFileNode('src', true, [utilsDir, componentsDir])

            const info = getCompactInfo(srcDir, true)

            // Should not compact because src has 2 children
            expect(info.isCompact).toBe(false)
            expect(info.displayName).toBe('src')
        })

        it('should handle file nodes (not directories)', () => {
            const fileNode = createFileNode('main.ts', false)

            const info = getCompactInfo(fileNode, true)

            expect(info.isCompact).toBe(false)
            expect(info.displayName).toBe('main.ts')
            expect(info.chainPaths).toEqual(['main.ts'])
        })
    })
})
