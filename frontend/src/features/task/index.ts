/**
 * Task feature module - Public API
 */

// Store
export { useTaskStore } from './model/task.store'

// API
export { taskApi } from './api/task.api'
export type {
    TaskAnalysisRequest,
    TaskAnalysisResponse,
    TaskExecutionRequest,
    TaskExecutionStatus
} from './api/task.api'

// Composables
export { useTaskExecution } from './composables/useTaskExecution'
export { useTaskStatus } from './composables/useTaskStatus'

// Types
export type {
    AnalysisResult,
    Task,
    TaskComplexity,
    TaskDraft,
    TaskExecutionResult,
    TaskStatus,
    TaskSuggestion,
    TaskType
} from './types/task.types'

// UI Components
export { default as TaskPanel } from './ui/TaskPanel.vue'
