/**
 * Dependency Graph Backend API Client
 * 
 * Provides high-performance dependency queries using backend caching.
 * Replaces O(n) frontend traversals with O(1) cached lookups.
 */

import * as wails from '#wailsjs/go/main/App'
import type { domain } from '#wailsjs/go/models'

export type FileDependency = domain.FileDependency

export interface DependencyGraphStats {
    fileCount: number
    dependencyCount: number
    lastAnalyzed: string
    isCached: boolean
    cacheSize: number
}

/**
 * Dependency API - Backend-powered dependency graph queries
 */
export const dependencyApi = {
    /**
     * Build dependency graph on project open
     * This should be called once when opening a project
     */
    async buildGraph(projectRoot: string): Promise<void> {
        if (!projectRoot) {
            throw new Error('projectRoot is required')
        }
        await wails.GetProjectDependencyGraph(projectRoot, [])
    },

    /**
     * Get dependencies for multiple files (batch query)
     * This is the primary method for getting dependencies efficiently
     */
    async getDependenciesBatch(
        projectRoot: string,
        filePaths: string[]
    ): Promise<Record<string, domain.FileDependency[]>> {
        if (!projectRoot) {
            throw new Error('projectRoot is required')
        }
        if (!filePaths || filePaths.length === 0) {
            return {}
        }

        try {
            const result = await wails.GetFileDependenciesBatch(projectRoot, filePaths)
            return result || {}
        } catch (error) {
            console.error('Failed to get dependencies batch:', error)
            throw error
        }
    },

    /**
     * Get files that import this file (incoming dependencies)
     */
    async getIncoming(projectRoot: string, filePath: string): Promise<string[]> {
        if (!projectRoot || !filePath) {
            throw new Error('projectRoot and filePath are required')
        }

        try {
            const result = await wails.GetIncomingDependencies(projectRoot, filePath)
            return result || []
        } catch (error) {
            console.error('Failed to get incoming dependencies:', error)
            throw error
        }
    },

    /**
     * Get files that this file imports (outgoing dependencies)
     */
    async getOutgoing(projectRoot: string, filePath: string): Promise<string[]> {
        if (!projectRoot || !filePath) {
            throw new Error('projectRoot and filePath are required')
        }

        try {
            const result = await wails.GetOutgoingDependencies(projectRoot, filePath)
            return result || []
        } catch (error) {
            console.error('Failed to get outgoing dependencies:', error)
            throw error
        }
    },

    /**
     * Check if dependency graph is cached
     */
    async isCached(projectRoot: string): Promise<boolean> {
        if (!projectRoot) {
            return false
        }

        try {
            return await wails.IsDependencyGraphCached(projectRoot)
        } catch (error) {
            console.error('Failed to check cache status:', error)
            return false
        }
    },

    /**
     * Get dependency graph statistics
     */
    async getStats(projectRoot: string): Promise<DependencyGraphStats | null> {
        if (!projectRoot) {
            return null
        }

        try {
            const stats = await wails.GetDependencyGraphStats(projectRoot)
            return stats || null
        } catch (error) {
            console.error('Failed to get dependency stats:', error)
            return null
        }
    },

    /**
     * Update single file in cached graph (incremental update)
     * Call this when a file is modified
     */
    async updateFile(projectRoot: string, filePath: string): Promise<void> {
        if (!projectRoot || !filePath) {
            throw new Error('projectRoot and filePath are required')
        }

        try {
            await wails.UpdateDependencyFile(projectRoot, filePath)
        } catch (error) {
            console.error('Failed to update file in dependency graph:', error)
            throw error
        }
    },

    /**
     * Remove file from cached graph
     * Call this when a file is deleted
     */
    async removeFile(projectRoot: string, filePath: string): Promise<void> {
        if (!projectRoot || !filePath) {
            throw new Error('projectRoot and filePath are required')
        }

        try {
            await wails.RemoveDependencyFile(projectRoot, filePath)
        } catch (error) {
            console.error('Failed to remove file from dependency graph:', error)
            throw error
        }
    },

    /**
     * Clear entire dependency cache
     * Call this when switching projects or on major file changes
     */
    async clearCache(): Promise<void> {
        try {
            await wails.ClearDependencyCache()
        } catch (error) {
            console.error('Failed to clear dependency cache:', error)
            throw error
        }
    }
}
