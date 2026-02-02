/**
 * Dependency Graph Composable
 * Detects and manages file dependencies for visual indicators
 */

import type { FileNode } from '@/types/domain'

export interface FileDependency {
    path: string
    type: 'import' | 'style' | 'test' | 'type' | 'related'
    confidence: 'high' | 'medium' | 'low'
    direction: 'incoming' | 'outgoing' // incoming: files that import this file, outgoing: files this file imports
}

/**
 * Pattern-based dependency detection
 * Phase 1: Uses naming conventions to detect related files
 * Phase 2: Will integrate with ProjectStructure backend API
 */
export function useDependencyGraph() {
    /**
     * Find dependencies for a given file
     * Returns files that the given file likely depends on (outgoing dependencies)
     */
    function findDependencies(
        filePath: string,
        allNodes: FileNode[]
    ): FileDependency[] {
        const deps: FileDependency[] = []
        const dir = filePath.substring(0, filePath.lastIndexOf('/'))
        const basename = getBasename(filePath)

        // Find related files by naming convention
        for (const node of allNodes) {
            if (node.isDir || node.path === filePath) continue

            const nodeDir = node.path.substring(0, node.path.lastIndexOf('/'))
            const nodeBasename = getBasename(node.path)

            // Same directory, similar name
            if (nodeDir === dir && nodeBasename === basename) {
                if (isStyleFile(node.path)) {
                    deps.push({
                        path: node.path,
                        type: 'style',
                        confidence: 'high',
                        direction: 'outgoing'
                    })
                } else if (isTestFile(node.path)) {
                    deps.push({
                        path: node.path,
                        type: 'test',
                        confidence: 'high',
                        direction: 'outgoing'
                    })
                } else if (isTypeFile(node.path)) {
                    deps.push({
                        path: node.path,
                        type: 'type',
                        confidence: 'high',
                        direction: 'outgoing'
                    })
                } else {
                    deps.push({
                        path: node.path,
                        type: 'related',
                        confidence: 'medium',
                        direction: 'outgoing'
                    })
                }
            }
        }

        return deps
    }

    /**
     * Find incoming dependencies (files that depend on this file)
     * Returns files that import/use the given file
     */
    function findIncomingDependencies(
        filePath: string,
        allNodes: FileNode[]
    ): FileDependency[] {
        const deps: FileDependency[] = []
        const basename = getBasename(filePath)

        // Find files that might depend on this file
        for (const node of allNodes) {
            if (node.isDir || node.path === filePath) continue

            const nodeBasename = getBasename(node.path)

            // If this is a test/style/type file, find the main file
            if (isTestFile(filePath) || isStyleFile(filePath) || isTypeFile(filePath)) {
                const dir = filePath.substring(0, filePath.lastIndexOf('/'))
                const nodeDir = node.path.substring(0, node.path.lastIndexOf('/'))

                if (nodeDir === dir && nodeBasename === basename && !isTestFile(node.path) && !isStyleFile(node.path) && !isTypeFile(node.path)) {
                    deps.push({
                        path: node.path,
                        type: 'related',
                        confidence: 'high',
                        direction: 'incoming'
                    })
                }
            }
        }

        return deps
    }

    /**
     * Find all dependencies (both incoming and outgoing) for a given file
     */
    function findAllDependencies(
        filePath: string,
        allNodes: FileNode[]
    ): { incoming: FileDependency[], outgoing: FileDependency[] } {
        return {
            incoming: findIncomingDependencies(filePath, allNodes),
            outgoing: findDependencies(filePath, allNodes)
        }
    }

    /**
     * Get basename without extension
     */
    function getBasename(path: string): string {
        const name = path.substring(path.lastIndexOf('/') + 1)

        // Handle compound extensions
        const compounds = [
            '.test.ts', '.test.js', '.test.tsx', '.test.jsx',
            '.spec.ts', '.spec.js', '.spec.tsx', '.spec.jsx',
            '.types.ts', '.d.ts',
            '.module.css', '.module.scss',
            '.stories.ts', '.stories.js', '.stories.tsx', '.stories.jsx'
        ]

        for (const ext of compounds) {
            if (name.endsWith(ext)) {
                return name.substring(0, name.length - ext.length)
            }
        }

        // Simple extension
        const dotIndex = name.indexOf('.')
        return dotIndex > 0 ? name.substring(0, dotIndex) : name
    }

    function isStyleFile(path: string): boolean {
        return /\.(css|scss|sass|less)$/.test(path)
    }

    function isTestFile(path: string): boolean {
        return /\.(test|spec)\.(ts|js|tsx|jsx)$/.test(path)
    }

    function isTypeFile(path: string): boolean {
        return /\.(types\.ts|d\.ts)$/.test(path)
    }

    /**
     * Find related files for a given file (for smart selection)
     * Returns files that are related by naming patterns (styles, tests, types)
     */
    function findRelatedFiles(
        filePath: string,
        allNodes: FileNode[]
    ): string[] {
        const related: string[] = []
        const baseName = getBasename(filePath)
        const dir = filePath.substring(0, filePath.lastIndexOf('/'))

        // Find files with same base name
        for (const node of allNodes) {
            if (node.isDir || node.path === filePath) continue

            const nodeDir = node.path.substring(0, node.path.lastIndexOf('/'))
            const nodeBaseName = getBasename(node.path)

            // Same directory, similar base name
            if (nodeDir === dir) {
                // Exact match or one contains the other
                if (
                    nodeBaseName === baseName ||
                    nodeBaseName.toLowerCase().includes(baseName.toLowerCase()) ||
                    baseName.toLowerCase().includes(nodeBaseName.toLowerCase())
                ) {
                    related.push(node.path)
                    continue
                }
            }

            // Common patterns across directories
            if (filePath.endsWith('.vue')) {
                // Component → store, types, styles
                if (
                    nodeBaseName === baseName &&
                    (isStyleFile(node.path) ||
                        isTestFile(node.path) ||
                        isTypeFile(node.path) ||
                        node.path.includes('.store.'))
                ) {
                    related.push(node.path)
                }
            } else if (filePath.match(/\.(ts|tsx|js|jsx)$/)) {
                // TypeScript/JavaScript → test file
                if (nodeBaseName === baseName && isTestFile(node.path)) {
                    related.push(node.path)
                }
            }
        }

        return [...new Set(related)]
    }

    return {
        findDependencies,
        findIncomingDependencies,
        findAllDependencies,
        findRelatedFiles
    }
}
