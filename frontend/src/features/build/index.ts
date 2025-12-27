/**
 * Build feature module - Public API
 */

// Store
export { useBuildStore } from './model/build.store'

// API
export { buildFeatureApi } from './api/build.api'

// Types
export type {
    BuildError,
    BuildHistoryEntry,
    BuildResult,
    BuildStatus,
    BuildStatusType,
    BuildWarning
} from './types'

// UI Components
export { default as BuildHistory } from './ui/BuildHistory.vue'
export { default as BuildHistoryItem } from './ui/BuildHistoryItem.vue'
export { default as BuildOutput } from './ui/BuildOutput.vue'
export { default as BuildPanel } from './ui/BuildPanel.vue'

