/**
 * Error types for API and application errors
 */

export interface TokenInfo {
    used: number
    limit: number
    remaining: number
    actual?: number  // For token limit exceeded errors
}

/**
 * Domain error from backend (Go)
 */
export interface DomainError {
    cause?: string
    message?: string
}

/**
 * Token limit exceeded error
 */
export interface TokenLimitError extends Error {
    tokenInfo: { actual: number; limit: number }
}

export interface ApiErrorDetails {
    message: string
    code?: string
    tokenInfo?: TokenInfo
    cause?: Error
}

/**
 * Type guard to check if error has tokenInfo with actual and limit
 */
export function hasTokenInfo(error: unknown): error is { tokenInfo: { actual: number; limit: number } } {
    if (
        typeof error !== 'object' ||
        error === null ||
        !('tokenInfo' in error)
    ) {
        return false
    }
    const tokenInfo = (error as { tokenInfo: unknown }).tokenInfo
    return (
        typeof tokenInfo === 'object' &&
        tokenInfo !== null &&
        'actual' in tokenInfo &&
        'limit' in tokenInfo
    )
}

/**
 * Type guard to check if error has cause
 */
export function hasCause(error: unknown): error is { cause: Error } {
    return (
        typeof error === 'object' &&
        error !== null &&
        'cause' in error &&
        (error as { cause: unknown }).cause instanceof Error
    )
}

/**
 * Extract error message from unknown error
 */
export function getErrorMessage(error: unknown): string {
    if (error instanceof Error) {
        return error.message
    }
    if (typeof error === 'string') {
        return error
    }
    if (typeof error === 'object' && error !== null && 'message' in error) {
        return String((error as { message: unknown }).message)
    }
    return 'Unknown error'
}

/**
 * Type guard for DomainError (backend Go errors)
 */
export function isDomainError(error: unknown): error is DomainError {
    return (
        typeof error === 'object' &&
        error !== null &&
        ('cause' in error || 'message' in error)
    )
}

/**
 * Extract message from domain error or regular error
 */
export function getDomainErrorMessage(error: unknown): string {
    if (isDomainError(error)) {
        return error.cause || error.message || String(error)
    }
    return getErrorMessage(error)
}

/**
 * Create a TokenLimitError
 */
export function createTokenLimitError(actual: number, limit: number): TokenLimitError {
    const error = new Error('TOKEN_LIMIT_EXCEEDED') as TokenLimitError
    error.tokenInfo = { actual, limit }
    return error
}
