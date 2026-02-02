/**
 * Memory Cleanup Utilities
 * 
 * Functions for aggressive memory cleanup when memory usage is critical.
 */

import { useLogger } from '@/composables/useLogger';

const logger = useLogger('MemoryCleanup');

/**
 * Force cleanup of memory with AGGRESSIVE strategies
 */
export function forceMemoryCleanup(): void {
    logger.warn('EMERGENCY: Forcing aggressive memory cleanup...');

    try {
        // Clear large arrays and objects
        window.largeObjects = [];
        window.cachedData = null;

        // Clear API caches
        import('@/composables/useApiCache').then(({ clearAllCaches }) => {
            clearAllCaches();
        }).catch(e => {
            logger.warn('Could not clear API caches:', e);
        });

        // Clear stores
        Promise.all([
            import('@/features/files/model/file.store'),
            import('@/features/context/model/context.store')
        ]).then(([{ useFileStore }, { useContextStore }]) => {
            const fileStore = useFileStore();
            const contextStore = useContextStore();

            fileStore.resetStore();
            contextStore.clearContext();
        }).catch(e => {
            logger.warn('Could not cleanup stores:', e);
        });

        // Clear Vue reactive caches if possible
        try {
            // Force cleanup of any global stores
            const stores = ['contextBuilder', 'fileTree', 'treeState', 'ui'] as const;
            stores.forEach(storeName => {
                const storeData = (window as unknown as Record<string, { cleanup?: () => void }>)[storeName];
                if (storeData && typeof storeData.cleanup === 'function') {
                    storeData.cleanup();
                }
            });
        } catch (e) {
            logger.warn('Could not cleanup Vue stores:', e);
        }

        // Force garbage collection multiple times if available
        if (window.gc) {
            try {
                window.gc();
                setTimeout(() => window.gc?.(), 100);
                setTimeout(() => window.gc?.(), 300);
                logger.debug('Multiple garbage collection cycles triggered');
            } catch (e) {
                logger.warn('Failed to trigger garbage collection', e);
            }
        }
    } catch (e) {
        logger.error('Error during emergency memory cleanup:', e);
    }
}

/**
 * Trigger garbage collection if available
 */
export function triggerGarbageCollection(): void {
    if (window.gc) {
        try {
            window.gc();
            logger.debug('Garbage collection triggered');
        } catch (e) {
            logger.warn('Failed to trigger garbage collection', e);
        }
    }
}
