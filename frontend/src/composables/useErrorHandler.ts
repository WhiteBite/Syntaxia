import { useI18n } from '@/composables/useI18n'
import { useUIStore } from '@/stores/ui.store'

export interface ErrorContext {
    context: string
    silent?: boolean
    retry?: () => Promise<void>
}

/**
 * Maps technical errors to user-friendly messages
 */
export function getErrorMessage(error: unknown): string {
    const { t } = useI18n()

    if (error instanceof Error) {
        const msg = error.message.toLowerCase()

        if (msg.includes('network') || msg.includes('fetch')) {
            return t('error.networkError')
        }
        if (msg.includes('not found') || msg.includes('404')) {
            return t('error.notFound')
        }
        if (msg.includes('timeout')) {
            return t('error.networkError')
        }
        if (msg.includes('invalid') || msg.includes('parse')) {
            return t('error.invalidData')
        }

        return error.message
    }

    if (typeof error === 'string') {
        return error
    }

    return t('error.generic')
}

/**
 * Handles errors with logging and toast notification
 */
export function handleError(error: unknown, context: string, silent = false): void {
    const uiStore = useUIStore()
    const message = getErrorMessage(error)

    console.error(`[${context}]`, error)

    if (!silent) {
        uiStore.addToast(message, 'error')
    }
}

/**
 * Retry wrapper with exponential backoff
 */
export async function withRetry<T>(
    fn: () => Promise<T>,
    retries = 3,
    delayMs = 1000
): Promise<T> {
    let lastError: unknown

    for (let i = 0; i < retries; i++) {
        try {
            return await fn()
        } catch (error) {
            lastError = error
            if (i < retries - 1) {
                await sleep(delayMs * (i + 1))
            }
        }
    }

    throw lastError
}

function sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms))
}

/**
 * Composable for error handling
 */
export function useErrorHandler() {
    return {
        handleError,
        getErrorMessage,
        withRetry
    }
}
