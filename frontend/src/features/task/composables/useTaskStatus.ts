/**
 * Task status tracking composable
 */

import { computed, ref } from 'vue'
import type { TaskStatus } from '../types/task.types'

export function useTaskStatus() {
    const status = ref<TaskStatus>('draft')
    const progress = ref(0)
    const currentStep = ref<string>('')

    /**
     * Status display properties
     */
    const statusColor = computed(() => {
        switch (status.value) {
            case 'draft':
                return 'gray'
            case 'analyzing':
                return 'blue'
            case 'ready':
                return 'green'
            case 'executing':
                return 'yellow'
            case 'completed':
                return 'green'
            case 'failed':
                return 'red'
            default:
                return 'gray'
        }
    })

    const statusIcon = computed(() => {
        switch (status.value) {
            case 'draft':
                return 'edit'
            case 'analyzing':
                return 'search'
            case 'ready':
                return 'check-circle'
            case 'executing':
                return 'play'
            case 'completed':
                return 'check'
            case 'failed':
                return 'x-circle'
            default:
                return 'circle'
        }
    })

    const statusLabel = computed(() => {
        switch (status.value) {
            case 'draft':
                return 'Draft'
            case 'analyzing':
                return 'Analyzing...'
            case 'ready':
                return 'Ready'
            case 'executing':
                return 'Executing...'
            case 'completed':
                return 'Completed'
            case 'failed':
                return 'Failed'
            default:
                return 'Unknown'
        }
    })

    const isInProgress = computed(() => {
        return status.value === 'analyzing' || status.value === 'executing'
    })

    const canExecute = computed(() => {
        return status.value === 'ready'
    })

    const isComplete = computed(() => {
        return status.value === 'completed' || status.value === 'failed'
    })

    /**
     * Update task status
     */
    function updateStatus(newStatus: TaskStatus, step?: string): void {
        status.value = newStatus
        if (step) {
            currentStep.value = step
        }

        // Auto-update progress based on status
        switch (newStatus) {
            case 'draft':
                progress.value = 0
                break
            case 'analyzing':
                progress.value = 25
                break
            case 'ready':
                progress.value = 50
                break
            case 'executing':
                progress.value = 75
                break
            case 'completed':
                progress.value = 100
                break
            case 'failed':
                progress.value = 0
                break
        }
    }

    /**
     * Update progress manually
     */
    function updateProgress(value: number): void {
        progress.value = Math.max(0, Math.min(100, value))
    }

    /**
     * Reset status to draft
     */
    function reset(): void {
        status.value = 'draft'
        progress.value = 0
        currentStep.value = ''
    }

    return {
        // State
        status,
        progress,
        currentStep,
        // Computed
        statusColor,
        statusIcon,
        statusLabel,
        isInProgress,
        canExecute,
        isComplete,
        // Actions
        updateStatus,
        updateProgress,
        reset
    }
}
