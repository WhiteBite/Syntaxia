/**
 * Memory Statistics Collection
 * 
 * Interfaces and functions for collecting memory usage statistics
 * from browser performance API and application stores.
 */

import { useLogger } from '@/composables/useLogger';

const logger = useLogger('MemoryStats');

export interface MemoryStats {
    used: number;        // Used memory in MB
    total: number;       // Total available memory in MB
    percentage: number;  // Usage percentage (0-100)
    storeMetrics?: StoreMetrics;
}

export interface StoreMetrics {
    fileStore: { nodesCount: number; memoryEstimate: number }
    contextStore: { cacheSize: number; chunkSize: number }
    apiCache: { entries: number; totalSize: number }
}

interface StoreImportsCache {
    useFileStore?: () => { nodes?: unknown[]; getMemoryUsage?: () => number }
    useContextStore?: () => { getMemoryUsage?: () => number; currentChunk?: { lines?: unknown[] } | null }
    getCacheStats?: () => { entries?: number; size?: number }
}

/**
 * Get current memory stats from browser performance API
 */
export async function getMemoryStats(
    isDevMode: boolean,
    storeImportsCache: StoreImportsCache
): Promise<MemoryStats | null> {
    if (!('performance' in window) || !performance.memory) {
        return null;
    }

    const memory = performance.memory;
    const used = Math.round(memory.usedJSHeapSize / (1024 * 1024));
    const total = Math.round(memory.jsHeapSizeLimit / (1024 * 1024));
    const percentage = Math.round((memory.usedJSHeapSize / memory.jsHeapSizeLimit) * 100);

    // Collect store metrics
    const storeMetrics = await collectStoreMetrics(isDevMode, storeImportsCache);

    return { used, total, percentage, storeMetrics };
}

/**
 * Collect memory metrics from all stores
 * In dev mode, skip detailed metrics to reduce memory overhead
 */
export async function collectStoreMetrics(
    isDevMode: boolean,
    storeImportsCache: StoreImportsCache
): Promise<StoreMetrics> {
    const metrics: StoreMetrics = {
        fileStore: { nodesCount: 0, memoryEstimate: 0 },
        contextStore: { cacheSize: 0, chunkSize: 0 },
        apiCache: { entries: 0, totalSize: 0 }
    };

    // Skip detailed store metrics in dev mode to reduce dynamic import overhead
    if (isDevMode) {
        return metrics;
    }

    try {
        // Import stores once and cache them to avoid repeated dynamic imports
        if (!storeImportsCache.useFileStore) {
            const fileStoreModule = await import('@/features/files/model/file.store');
            storeImportsCache.useFileStore = fileStoreModule.useFileStore as StoreImportsCache['useFileStore'];
        }

        if (!storeImportsCache.useContextStore) {
            const contextStoreModule = await import('@/features/context/model/context.store');
            storeImportsCache.useContextStore = contextStoreModule.useContextStore as StoreImportsCache['useContextStore'];
        }

        if (!storeImportsCache.getCacheStats) {
            const apiCacheModule = await import('@/composables/useApiCache');
            storeImportsCache.getCacheStats = apiCacheModule.getCacheStats as StoreImportsCache['getCacheStats'];
        }

        const fileStore = storeImportsCache.useFileStore?.();
        const contextStore = storeImportsCache.useContextStore?.();
        const cacheStats = storeImportsCache.getCacheStats?.();

        if (fileStore) {
            metrics.fileStore = {
                nodesCount: fileStore.nodes?.length || 0,
                memoryEstimate: fileStore.getMemoryUsage ? fileStore.getMemoryUsage() : 0
            };
        }

        if (contextStore) {
            metrics.contextStore = {
                cacheSize: contextStore.getMemoryUsage ? contextStore.getMemoryUsage() : 0,
                chunkSize: contextStore.currentChunk?.lines?.length || 0
            };
        }

        if (cacheStats) {
            metrics.apiCache = {
                entries: cacheStats.entries || 0,
                totalSize: cacheStats.size || 0
            };
        }
    } catch (e) {
        logger.warn('Could not collect store metrics:', e);
    }

    return metrics;
}

/**
 * Track memory trend to detect rapid growth
 * Returns growth rate in MB per minute
 */
export function calculateMemoryTrend(
    memoryHistory: number[],
    currentUsed: number,
    pollingInterval: number
): number {
    memoryHistory.push(currentUsed);

    // Keep only last 10 measurements
    if (memoryHistory.length > 10) {
        memoryHistory.shift();
    }

    // Calculate growth rate (MB per minute)
    if (memoryHistory.length < 2) {
        return 0;
    }

    const oldest = memoryHistory[0];
    const newest = memoryHistory[memoryHistory.length - 1];
    const timeDiff = (memoryHistory.length - 1) * pollingInterval / 60000; // minutes

    return (newest - oldest) / timeDiff;
}

/**
 * Dump heap snapshot for analysis in Chrome DevTools
 */
export async function dumpHeapSnapshot(reason: string): Promise<void> {
    if ('performance' in window && performance.writeHeapSnapshot) {
        try {
            const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
            const filename = `heap-snapshot-${timestamp}-${reason}.heapsnapshot`;
            await performance.writeHeapSnapshot(filename);
            logger.debug(`Heap snapshot saved: ${filename}`);
            return;
        } catch (e) {
            logger.error('Failed to dump heap snapshot:', e);
        }
    } else {
        logger.warn('Heap snapshot API not available. Run Chrome with --enable-precise-memory-info flag.');
    }
}
