/**
 * Base utilities for API modules
 * Provides common error handling, logging, and type-safe wrappers
 */

import { useLogger } from '@/composables/useLogger'

const logger = useLogger('API')

/**
 * Custom API Error with context
 */
export class ApiError extends Error {
    constructor(
        public readonly context: string,
        public readonly originalError?: unknown
    ) {
        const message = originalError instanceof Error
            ? `${context}: ${originalError.message}`
            : context
        super(message)
        this.name = 'ApiError'
    }
}

/**
 * Wrapper for API calls with consistent error handling
 * @param fn - Async function to execute
 * @param context - Error context for logging
 * @param options - Additional options
 */
export async function apiCall<T>(
    fn: () => Promise<T>,
    context: string,
    options?: {
        logContext?: string
        rethrow?: boolean
    }
): Promise<T> {
    const logPrefix = options?.logContext || context
    try {
        return await fn()
    } catch (error) {
        logger.error(`${logPrefix}:`, error)

        if (options?.rethrow && error instanceof Error) {
            throw error
        }

        throw new ApiError(context, error)
    }
}

/**
 * Wrapper for API calls that return a default value on error
 * Use for non-critical operations where failure is acceptable
 */
export async function apiCallSafe<T>(
    fn: () => Promise<T>,
    defaultValue: T,
    context?: string
): Promise<T> {
    try {
        return await fn()
    } catch (error) {
        if (context) {
            logger.warn(`${context} (using default):`, error)
        }
        return defaultValue
    }
}

/**
 * Wrapper for API calls that may return null on error
 */
export async function apiCallNullable<T>(
    fn: () => Promise<T>,
    context?: string
): Promise<T | null> {
    try {
        return await fn()
    } catch (error) {
        if (context) {
            logger.warn(`${context}:`, error)
        }
        return null
    }
}

/**
 * Parse JSON response with error handling
 */
export function parseJsonResponse<T>(json: string, context: string): T {
    try {
        return JSON.parse(json) as T
    } catch (error) {
        logger.error(`Failed to parse JSON (${context}):`, error)
        throw new ApiError(`Invalid JSON response: ${context}`, error)
    }
}

/**
 * Parse JSON response with default value on error
 */
export function parseJsonSafe<T>(json: string, defaultValue: T): T {
    try {
        return JSON.parse(json) as T
    } catch {
        return defaultValue
    }
}

// Re-export for backward compatibility
export const apiCallWithDefault = apiCallSafe
