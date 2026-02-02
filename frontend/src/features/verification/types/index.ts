/**
 * Verification feature types
 */

import type { Severity, VerificationType } from '../api/verification.api'

// Re-export types from API for convenience
export type {
    Severity,
    VerificationError,
    VerificationResult,
    VerificationType
} from '../api/verification.api'

// Re-export types from store
export type {
    VerificationStatus,
    VerificationStatusType
} from '../model/verification.store'

// Additional UI types
export interface VerificationFilter {
    severity?: Severity[]
    source?: VerificationType[]
    file?: string
}

export interface VerificationStats {
    total: number
    errors: number
    warnings: number
    info: number
    bySource: Record<VerificationType, number>
}
