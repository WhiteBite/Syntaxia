/**
 * Task execution composable
 */

import { useLogger } from '@/composables/useLogger'
import { ref } from 'vue'
import type { TaskExecutionStatus } from '../api/task.api'
import { taskApi } from '../api/task.api'

const logger = useLogger('useTaskExecution')

export function useTaskExecution() {
    const isExecuting = ref(false)
    const currentTaskId = ref<string | null>(null)
    const executionStatus = ref<TaskExecutionStatus | null>(null)
    const error = ref<string | null>(null)

    /**
     * Execute a task with AI assistance
     */
    async function executeTask(
        description: string,
        selectedFiles: string[],
        projectPath: string
    ): Promise<boolean> {
        if (isExecuting.value) {
            logger.warn('Task execution already in progress')
            return false
        }

        isExecuting.value = true
        error.value = null

        try {
            logger.info('Starting task execution', { description, fileCount: selectedFiles.length })

            const status = await taskApi.executeTask({
                description,
                selectedFiles,
                projectPath
            })

            currentTaskId.value = status.id
            executionStatus.value = status

            // Poll for status updates
            await pollTaskStatus(status.id)

            return executionStatus.value?.status === 'completed'
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Task execution failed'
            logger.error('Task execution failed:', err)
            return false
        } finally {
            isExecuting.value = false
        }
    }

    /**
     * Poll task status until completion
     */
    async function pollTaskStatus(taskId: string): Promise<void> {
        const maxAttempts = 60 // 5 minutes with 5s intervals
        let attempts = 0

        while (attempts < maxAttempts) {
            const status = await taskApi.getTaskStatus(taskId)
            executionStatus.value = status

            if (status.status === 'completed' || status.status === 'failed') {
                if (status.status === 'failed') {
                    error.value = status.error || 'Task execution failed'
                }
                break
            }

            // Wait 5 seconds before next poll
            await new Promise(resolve => setTimeout(resolve, 5000))
            attempts++
        }

        if (attempts >= maxAttempts) {
            error.value = 'Task execution timeout'
        }
    }

    /**
     * Cancel current task execution
     */
    function cancelExecution(): void {
        if (currentTaskId.value) {
            logger.info('Cancelling task execution', { taskId: currentTaskId.value })
            // Backend cancellation would be implemented here
            isExecuting.value = false
            currentTaskId.value = null
            executionStatus.value = null
        }
    }

    /**
     * Reset execution state
     */
    function reset(): void {
        isExecuting.value = false
        currentTaskId.value = null
        executionStatus.value = null
        error.value = null
    }

    return {
        // State
        isExecuting,
        currentTaskId,
        executionStatus,
        error,
        // Actions
        executeTask,
        cancelExecution,
        reset
    }
}
