/**
 * Extended Performance interface with Chrome-specific memory API
 * @see https://developer.chrome.com/docs/devtools/memory-problems/
 */

interface PerformanceMemory {
    /** Total size of the JS heap in bytes */
    jsHeapSizeLimit: number
    /** Total allocated heap size in bytes */
    totalJSHeapSize: number
    /** Currently used heap size in bytes */
    usedJSHeapSize: number
}

interface PerformanceWithMemory extends Performance {
    /** Chrome-specific memory info (only available with --enable-precise-memory-info flag) */
    memory?: PerformanceMemory
    /** Write heap snapshot for debugging (Chrome DevTools) */
    writeHeapSnapshot?: (filename: string) => Promise<void>
}

declare global {
    interface Window {
        /** Force garbage collection (only available with --expose-gc flag) */
        gc?: () => void
        /** Large objects cache for cleanup */
        largeObjects?: unknown[]
        /** Cached data for cleanup */
        cachedData?: unknown
    }

    interface Performance {
        /** Chrome-specific memory info */
        memory?: PerformanceMemory
        /** Write heap snapshot for debugging */
        writeHeapSnapshot?: (filename: string) => Promise<void>
    }

    interface Navigator {
        /** Device memory in GB (Chrome-specific) */
        deviceMemory?: number
    }

}

export type { PerformanceMemory, PerformanceWithMemory }

