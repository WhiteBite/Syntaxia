/**
 * Build feature API
 * Wraps existing build.api.ts from services
 */

import { useLogger } from '@/composables/useLogger'
import { buildApi as baseBuildApi } from '@/services/api/build.api'
import type { BuildHistoryEntry, BuildResult } from '../types'

const logger = useLogger('BuildFeatureApi')

// Storage key for build history
const BUILD_HISTORY_KEY = 'Syntaxia_build_history'
const MAX_HISTORY_ENTRIES = 10

/**
 * Generate unique ID for build
 */
function generateBuildId(): string {
    return `build-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`
}

/**
 * Build feature API
 */
export const buildFeatureApi = {
    /**
     * Run project build
     */
    async build(projectPath: string): Promise<BuildResult> {
        const startTime = new Date().toISOString()
        const buildId = generateBuildId()
        let output = ''

        try {
            logger.info(`Starting build for ${projectPath}`)

            const result = await baseBuildApi.build(projectPath, 'auto')

            const endTime = new Date().toISOString()
            const duration = new Date(endTime).getTime() - new Date(startTime).getTime()

            // Collect output
            output = result.output || ''
            if (result.error) {
                output += `\n\nError: ${result.error}`
            }
            if (result.warnings && result.warnings.length > 0) {
                output += `\n\nWarnings:\n${result.warnings.join('\n')}`
            }

            const buildResult: BuildResult = {
                id: buildId,
                status: result.success ? 'success' : 'failed',
                startTime,
                endTime,
                duration,
                output,
                errors: result.error ? [{
                    file: projectPath,
                    line: 0,
                    message: result.error
                }] : [],
                warnings: result.warnings?.map(w => ({
                    file: projectPath,
                    line: 0,
                    message: w
                })) || []
            }

            // Save to history
            saveBuildToHistory(projectPath, buildResult)

            return buildResult
        } catch (error) {
            logger.error('Build failed:', error)

            const endTime = new Date().toISOString()
            const duration = new Date(endTime).getTime() - new Date(startTime).getTime()

            const buildResult: BuildResult = {
                id: buildId,
                status: 'failed',
                startTime,
                endTime,
                duration,
                output: error instanceof Error ? error.message : 'Build failed',
                errors: [{
                    file: projectPath,
                    line: 0,
                    message: error instanceof Error ? error.message : 'Build failed'
                }],
                warnings: []
            }

            saveBuildToHistory(projectPath, buildResult)

            return buildResult
        }
    },

    /**
     * Clean build artifacts (simulated - actual implementation depends on backend)
     */
    async clean(projectPath: string): Promise<BuildResult> {
        const startTime = new Date().toISOString()
        const buildId = generateBuildId()

        try {
            logger.info(`Cleaning build for ${projectPath}`)

            // Simulate clean operation - in real implementation this would call backend
            await new Promise(resolve => setTimeout(resolve, 500))

            const endTime = new Date().toISOString()
            const duration = new Date(endTime).getTime() - new Date(startTime).getTime()

            return {
                id: buildId,
                status: 'success',
                startTime,
                endTime,
                duration,
                output: 'Build artifacts cleaned successfully.',
                errors: [],
                warnings: []
            }
        } catch (error) {
            logger.error('Clean failed:', error)

            const endTime = new Date().toISOString()
            const duration = new Date(endTime).getTime() - new Date(startTime).getTime()

            return {
                id: buildId,
                status: 'failed',
                startTime,
                endTime,
                duration,
                output: error instanceof Error ? error.message : 'Clean failed',
                errors: [{
                    file: projectPath,
                    line: 0,
                    message: error instanceof Error ? error.message : 'Clean failed'
                }],
                warnings: []
            }
        }
    },

    /**
     * Cancel current build (simulated)
     */
    async cancelBuild(): Promise<void> {
        logger.info('Build cancelled')
        // In real implementation, this would signal the backend to cancel
    },

    /**
     * Get build history for project
     */
    async getBuildHistory(projectPath: string): Promise<BuildHistoryEntry[]> {
        try {
            const historyKey = `${BUILD_HISTORY_KEY}_${hashPath(projectPath)}`
            const stored = localStorage.getItem(historyKey)

            if (stored) {
                return JSON.parse(stored)
            }

            return []
        } catch (error) {
            logger.warn('Failed to load build history:', error)
            return []
        }
    }
}

/**
 * Save build result to history
 */
function saveBuildToHistory(projectPath: string, result: BuildResult): void {
    try {
        const historyKey = `${BUILD_HISTORY_KEY}_${hashPath(projectPath)}`
        const stored = localStorage.getItem(historyKey)
        let history: BuildHistoryEntry[] = stored ? JSON.parse(stored) : []

        const entry: BuildHistoryEntry = {
            id: result.id,
            status: result.status,
            startTime: result.startTime,
            endTime: result.endTime,
            duration: result.duration,
            errorsCount: result.errors.length,
            warningsCount: result.warnings.length,
            output: result.output
        }

        // Add to beginning
        history.unshift(entry)

        // Keep only last N entries
        if (history.length > MAX_HISTORY_ENTRIES) {
            history = history.slice(0, MAX_HISTORY_ENTRIES)
        }

        localStorage.setItem(historyKey, JSON.stringify(history))
    } catch (error) {
        logger.warn('Failed to save build history:', error)
    }
}

/**
 * Simple hash for path to use as storage key
 */
function hashPath(path: string): string {
    let hash = 0
    for (let i = 0; i < path.length; i++) {
        const char = path.charCodeAt(i)
        hash = ((hash << 5) - hash) + char
        hash = hash & hash
    }
    return Math.abs(hash).toString(36)
}

export default buildFeatureApi
