/**
 * Logger composable for consistent logging across the application.
 * Only logs in development mode to avoid console pollution in production.
 */

const isDev = import.meta.env.DEV

export type LogLevel = 'debug' | 'info' | 'warn' | 'error'

export interface Logger {
    debug: (...args: unknown[]) => void
    info: (...args: unknown[]) => void
    warn: (...args: unknown[]) => void
    error: (...args: unknown[]) => void
}

/**
 * Creates a logger instance with a specific context prefix.
 * Debug and info logs are only shown in development mode.
 * Warn and error logs are always shown.
 * 
 * @param context - The context/module name to prefix log messages with
 * @returns Logger instance with debug, info, warn, error methods
 * 
 * @example
 * const logger = useLogger('FileExplorer')
 * logger.debug('Loading files...', { path: '/src' })
 * logger.error('Failed to load files', error)
 */
export function useLogger(context: string): Logger {
    const prefix = `[${context}]`

    return {
        debug: (...args: unknown[]) => {
            if (isDev) {
                console.log(prefix, ...args)
            }
        },
        info: (...args: unknown[]) => {
            if (isDev) {
                console.info(prefix, ...args)
            }
        },
        warn: (...args: unknown[]) => {
            console.warn(prefix, ...args)
        },
        error: (...args: unknown[]) => {
            console.error(prefix, ...args)
        },
    }
}
