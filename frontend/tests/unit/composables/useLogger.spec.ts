import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

// Mock import.meta.env
const mockEnv = { DEV: true }
vi.mock('import.meta', () => ({
    env: mockEnv
}))

describe('useLogger', () => {
    let consoleSpy: {
        log: ReturnType<typeof vi.spyOn>
        info: ReturnType<typeof vi.spyOn>
        warn: ReturnType<typeof vi.spyOn>
        error: ReturnType<typeof vi.spyOn>
    }

    beforeEach(() => {
        consoleSpy = {
            log: vi.spyOn(console, 'log').mockImplementation(() => { }),
            info: vi.spyOn(console, 'info').mockImplementation(() => { }),
            warn: vi.spyOn(console, 'warn').mockImplementation(() => { }),
            error: vi.spyOn(console, 'error').mockImplementation(() => { }),
        }
    })

    afterEach(() => {
        vi.restoreAllMocks()
    })

    it('should create logger with context prefix', async () => {
        const { useLogger } = await import('@/composables/useLogger')
        const logger = useLogger('TestContext')

        expect(logger).toBeDefined()
        expect(typeof logger.debug).toBe('function')
        expect(typeof logger.info).toBe('function')
        expect(typeof logger.warn).toBe('function')
        expect(typeof logger.error).toBe('function')
    })

    it('should always log warnings regardless of environment', async () => {
        const { useLogger } = await import('@/composables/useLogger')
        const logger = useLogger('TestContext')

        logger.warn('test warning', { data: 123 })

        expect(consoleSpy.warn).toHaveBeenCalledWith('[TestContext]', 'test warning', { data: 123 })
    })

    it('should always log errors regardless of environment', async () => {
        const { useLogger } = await import('@/composables/useLogger')
        const logger = useLogger('TestContext')

        logger.error('test error', new Error('test'))

        expect(consoleSpy.error).toHaveBeenCalled()
        expect(consoleSpy.error.mock.calls[0][0]).toBe('[TestContext]')
    })

    it('should include context prefix in all log messages', async () => {
        const { useLogger } = await import('@/composables/useLogger')
        const logger = useLogger('MyModule')

        logger.warn('message')

        expect(consoleSpy.warn.mock.calls[0][0]).toBe('[MyModule]')
    })
})
