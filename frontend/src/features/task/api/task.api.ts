/**
 * Task API - Wails bridge for taskflow operations
 */

import type { domain } from '#wailsjs/go/models'
import { useLogger } from '@/composables/useLogger'
import { apiService } from '@/services/api.service'

const logger = useLogger('TaskApi')

export interface TaskAnalysisRequest {
    description: string
    projectPath: string
    fileTree: domain.FileNode[]
}

export interface TaskAnalysisResponse {
    suggestedFiles: string[]
    complexity: 'low' | 'medium' | 'high'
    estimatedTime: number
    recommendations: string[]
}

export interface TaskExecutionRequest {
    description: string
    selectedFiles: string[]
    projectPath: string
}

export interface TaskExecutionStatus {
    id: string
    status: 'pending' | 'running' | 'completed' | 'failed'
    progress: number
    currentStep?: string
    error?: string
}

export const taskApi = {
    /**
     * Analyze task and suggest relevant files
     */
    async analyzeTask(request: TaskAnalysisRequest): Promise<TaskAnalysisResponse> {
        logger.info('Analyzing task', { description: request.description })

        try {
            const suggestedFiles = await apiService.suggestContextFiles(
                request.description,
                request.fileTree
            )

            // Calculate complexity based on suggested files
            let complexity: 'low' | 'medium' | 'high' = 'low'
            if (suggestedFiles.length > 10) {
                complexity = 'high'
            } else if (suggestedFiles.length > 5) {
                complexity = 'medium'
            }

            // Estimate time based on complexity and file count
            const estimatedTime = suggestedFiles.length * 5 + (
                complexity === 'high' ? 30 : complexity === 'medium' ? 15 : 5
            )

            const recommendations = generateRecommendations(complexity, suggestedFiles.length)

            return {
                suggestedFiles,
                complexity,
                estimatedTime,
                recommendations
            }
        } catch (error) {
            logger.error('Failed to analyze task:', error)
            throw error
        }
    },

    /**
     * Execute task with AI assistance
     */
    async executeTask(request: TaskExecutionRequest): Promise<TaskExecutionStatus> {
        logger.info('Executing task', { description: request.description })

        // This would integrate with AI chat to execute the task
        // For now, return a placeholder status
        return {
            id: `task-${Date.now()}`,
            status: 'pending',
            progress: 0
        }
    },

    /**
     * Get task execution status
     */
    async getTaskStatus(taskId: string): Promise<TaskExecutionStatus> {
        logger.info('Getting task status', { taskId })

        // Placeholder - would query backend for actual status
        return {
            id: taskId,
            status: 'pending',
            progress: 0
        }
    }
}

function generateRecommendations(complexity: string, fileCount: number): string[] {
    const recommendations: string[] = []

    if (complexity === 'high') {
        recommendations.push('Consider breaking down into smaller tasks')
        recommendations.push('This is a complex task that may require multiple iterations')
    }

    if (fileCount > 15) {
        recommendations.push('Large number of files involved - ensure proper testing')
    }

    if (fileCount === 0) {
        recommendations.push('No relevant files found - try refining your task description')
    } else {
        recommendations.push(`${fileCount} relevant files identified`)
    }

    recommendations.push('Review existing similar implementations')
    recommendations.push('Add tests for new functionality')

    return recommendations
}
