import { getErrorMessage, handleError, withRetry } from '@/composables/useErrorHandler'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock dependencies
const mockAddToast = vi.fn()
vi.mock('@/stores/ui.store', () => ({
    useUIStore: () => ({
        addToast: mockAddToast
    })
}))

vi.mock('@/composables/useI18n', () => ({
    useI18n: () => ({
        t: (key: string) => key
    })
}))

describe('useErrorHandler', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    describe('handleError', () => {
        it('shows toast notification', () => {
            const error = new Error('Test error')

            handleError(error, 'TestContext')

            expect(mockAddToast).toHaveBeenCalledWith('Test error', 'error')
        })

        it('does not show toast when silent=true', () => {
            const error = new Error('Silent error')

            handleError(error, 'TestContext', true)

            expect(mockAddToast).not.toHaveBeenCalled()
        })
    })

    describe('getErrorMessage', () => {
        it('returns user-friendly message for network errors', () => {
            const error = new Error('Network request failed')

            const message = getErrorMessage(error)

            expect(message).toBe('error.networkError')
        })

        it('returns user-friendly message for not found errors', () => {
            const error = new Error('Resource not found')

            const message = getErrorMessage(error)

            expect(message).toBe('error.notFound')
        })

        it('returns generic message for unknown errors', () => {
            const message = getErrorMessage({})

            expect(message).toBe('error.generic')
        })
    })

    describe('withRetry', () => {
        it('retries on failure', async () => {
            let attempts = 0
            const fn = vi.fn().mockImplementation(async () => {
                attempts++
                if (attempts < 3) throw new Error('Fail')
                return 'success'
            })

            const result = await withRetry(fn, 3, 10)

            expect(result).toBe('success')
            expect(fn).toHaveBeenCalledTimes(3)
        })

        it('throws after max retries', async () => {
            const fn = vi.fn().mockRejectedValue(new Error('Always fails'))

            await expect(withRetry(fn, 3, 10)).rejects.toThrow('Always fails')
            expect(fn).toHaveBeenCalledTimes(3)
        })
    })
})
