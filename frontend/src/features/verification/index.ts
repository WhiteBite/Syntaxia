/**
 * Verification feature module - Public API
 */

// Store
export { useVerificationStore } from './model/verification.store'
export type { VerificationStatus, VerificationStatusType } from './model/verification.store'

// API
export { verificationApi } from './api/verification.api'
export type {
    Severity,
    VerificationError,
    VerificationResult,
    VerificationType
} from './api/verification.api'

// Composables
export { useVerification } from './composables/useVerification'

// Types
export type {
    VerificationFilter,
    VerificationStats
} from './types'

// UI Components
export { default as ErrorItem } from './ui/ErrorItem.vue'
export { default as ErrorList } from './ui/ErrorList.vue'
export { default as VerificationPanel } from './ui/VerificationPanel.vue'
export { default as VerificationStatusBar } from './ui/VerificationStatusBar.vue'

