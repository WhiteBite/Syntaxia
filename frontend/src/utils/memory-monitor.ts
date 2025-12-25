/**
 * Memory Monitor Utility
 *
 * This utility helps prevent Out Of Memory (OOM) errors by:
 * 1. Monitoring memory usage
 * 2. Providing early warnings
 * 3. Cleaning up large objects when memory usage is high
 * 4. Forcing garbage collection when possible
 */

import { useLogger } from '@/composables/useLogger';
import { useUIStore } from '@/stores/ui.store';

const logger = useLogger('MemoryMonitor');

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

export interface MemoryWarningOptions {
  warningThreshold?: number;   // Absolute MB threshold for warning (default: 15MB)
  criticalThreshold?: number;  // Absolute MB threshold for critical warning (default: 25MB)
  pollingInterval?: number;    // How often to check memory in ms (default: 1000ms)
  showToasts?: boolean;        // Whether to show toast notifications (default: true)
  autoCleanup?: boolean;       // Whether to auto cleanup when critical (default: true)
}

const DEFAULT_OPTIONS: MemoryWarningOptions = {
  warningThreshold: 150,  // FIXED: 150MB (was incorrectly 15MB)
  criticalThreshold: 250, // FIXED: 250MB (was incorrectly 25MB)
  pollingInterval: 10000,  // Check every 10 seconds to reduce overhead
  showToasts: true,
  autoCleanup: true
};

export class MemoryMonitor {
  private static instance: MemoryMonitor;
  private options: MemoryWarningOptions;
  private monitorInterval: number | null = null;
  private warningIssued = false;
  private criticalIssued = false;
  private uiStore = useUIStore();

  // Cache store imports to avoid repeated dynamic imports
  private storeImportsCache: StoreImportsCache = {};

  // Track if we're in dev mode to reduce monitoring overhead
  private isDevMode = import.meta.env.DEV;

  private constructor(options: MemoryWarningOptions = {}) {
    // In dev mode, use much longer polling interval to reduce overhead
    const devOptions = this.isDevMode ? {
      pollingInterval: 30000, // 30 seconds in dev (was 10s)
    } : {};
    this.options = { ...DEFAULT_OPTIONS, ...devOptions, ...options };
  }

  /**
   * Get the singleton instance of MemoryMonitor
   */
  public static getInstance(options?: MemoryWarningOptions): MemoryMonitor {
    if (!MemoryMonitor.instance) {
      MemoryMonitor.instance = new MemoryMonitor(options);
    } else if (options) {
      MemoryMonitor.instance.updateOptions(options);
    }
    return MemoryMonitor.instance;
  }

  /**
   * Update monitoring options
   */
  public updateOptions(options: Partial<MemoryWarningOptions>): void {
    this.options = { ...this.options, ...options };

    // If monitoring is active, restart it with new options
    if (this.monitorInterval) {
      this.stopMonitoring();
      this.startMonitoring();
    }
  }

  /**
   * Get current memory stats with store metrics
   */
  public async getMemoryStats(): Promise<MemoryStats | null> {
    if (!('performance' in window) || !performance.memory) {
      return null;
    }

    const memory = performance.memory;
    const used = Math.round(memory.usedJSHeapSize / (1024 * 1024));
    const total = Math.round(memory.jsHeapSizeLimit / (1024 * 1024));
    const percentage = Math.round((memory.usedJSHeapSize / memory.jsHeapSizeLimit) * 100);

    // Collect store metrics
    const storeMetrics = await this.collectStoreMetrics();

    return { used, total, percentage, storeMetrics };
  }

  /**
   * Collect memory metrics from all stores
   * In dev mode, skip detailed metrics to reduce memory overhead
   */
  private async collectStoreMetrics(): Promise<StoreMetrics> {
    const metrics: StoreMetrics = {
      fileStore: { nodesCount: 0, memoryEstimate: 0 },
      contextStore: { cacheSize: 0, chunkSize: 0 },
      apiCache: { entries: 0, totalSize: 0 }
    };

    // Skip detailed store metrics in dev mode to reduce dynamic import overhead
    if (this.isDevMode) {
      return metrics;
    }

    try {
      // Import stores once and cache them to avoid repeated dynamic imports
      if (!this.storeImportsCache.useFileStore) {
        const fileStoreModule = await import('@/features/files/model/file.store');
        this.storeImportsCache.useFileStore = fileStoreModule.useFileStore as StoreImportsCache['useFileStore'];
      }

      if (!this.storeImportsCache.useContextStore) {
        const contextStoreModule = await import('@/features/context/model/context.store');
        this.storeImportsCache.useContextStore = contextStoreModule.useContextStore as StoreImportsCache['useContextStore'];
      }

      if (!this.storeImportsCache.getCacheStats) {
        const apiCacheModule = await import('@/composables/useApiCache');
        this.storeImportsCache.getCacheStats = apiCacheModule.getCacheStats as StoreImportsCache['getCacheStats'];
      }

      const fileStore = this.storeImportsCache.useFileStore?.();
      const contextStore = this.storeImportsCache.useContextStore?.();
      const cacheStats = this.storeImportsCache.getCacheStats?.();

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
      console.warn('[MemoryMonitor] Could not collect store metrics:', e);
    }

    return metrics;
  }

  /**
   * Start memory monitoring
   */
  public startMonitoring(): void {
    if (this.monitorInterval) {
      return; // Already monitoring
    }

    void this.checkMemory();
    this.monitorInterval = window.setInterval(
      () => void this.checkMemory(),
      this.options.pollingInterval
    );

    logger.debug(`Memory monitoring started (interval: ${this.options.pollingInterval}ms)`);
  }

  /**
   * Stop memory monitoring
   */
  public stopMonitoring(): void {
    if (this.monitorInterval) {
      clearInterval(this.monitorInterval);
      this.monitorInterval = null;
      logger.debug('Memory monitoring stopped');
    }
  }

  /**
   * Force cleanup of memory with AGGRESSIVE strategies
   */
  public forceCleanup(): void {
    logger.warn('EMERGENCY: Forcing aggressive memory cleanup...');

    // Attempt to free up memory by clearing large objects in memory
    try {
      // Clear large arrays and objects
      window.largeObjects = [];
      window.cachedData = null;

      // Clear API caches
      import('@/composables/useApiCache').then(({ clearAllCaches }) => {
        clearAllCaches();
      }).catch(e => {
        console.warn('Could not clear API caches:', e);
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
        console.warn('Could not cleanup stores:', e);
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
        console.warn('Could not cleanup Vue stores:', e);
      }

      // Force garbage collection multiple times if available
      if (window.gc) {
        try {
          window.gc();
          setTimeout(() => window.gc?.(), 100);
          setTimeout(() => window.gc?.(), 300);
          logger.debug('Multiple garbage collection cycles triggered');
        } catch (e) {
          console.warn('Failed to trigger garbage collection', e);
        }
      }

      // Add a small delay to allow garbage collection to work
      setTimeout(() => {
        void this.getMemoryStats().then(stats => {
          if (stats) {
            logger.debug(`Memory after aggressive cleanup: ${stats.used}MB / ${stats.total}MB (${stats.percentage}%)`);
          }
        });
      }, 500);
    } catch (e) {
      console.error('Error during emergency memory cleanup:', e);
    }
  }

  /**
   * Dump heap snapshot for analysis in Chrome DevTools
   */
  public async dumpHeapSnapshot(reason: string): Promise<void> {
    if ('performance' in window && performance.writeHeapSnapshot) {
      try {
        const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
        const filename = `heap-snapshot-${timestamp}-${reason}.heapsnapshot`;
        await performance.writeHeapSnapshot(filename);
        logger.debug(`Heap snapshot saved: ${filename}`);
        this.uiStore.addToast(`Heap snapshot saved: ${filename}`, 'info');
      } catch (e) {
        console.error('Failed to dump heap snapshot:', e);
      }
    } else {
      console.warn('Heap snapshot API not available. Run Chrome with --enable-precise-memory-info flag.');
    }
  }

  /**
   * Register callback for critical memory events
   * Returns unsubscribe function to prevent memory leaks
   */
  private criticalCallbacks: Set<() => void> = new Set();

  public onCritical(callback: () => void): () => void {
    this.criticalCallbacks.add(callback);
    // Return unsubscribe function
    return () => {
      this.criticalCallbacks.delete(callback);
    };
  }

  /**
   * Clear all critical callbacks
   */
  public clearCriticalCallbacks(): void {
    this.criticalCallbacks.clear();
  }

  /**
   * Cleanup and reset the monitor (call when unmounting)
   */
  public cleanup(): void {
    this.stopMonitoring();
    this.criticalCallbacks.clear();
    this.memoryHistory = [];
    this.warningIssued = false;
    this.criticalIssued = false;
  }

  /**
   * Check current memory usage with PERCENTAGE-BASED thresholds (more reliable)
   */
  private memoryHistory: number[] = [];

  private async checkMemory(): Promise<void> {
    const stats = await this.getMemoryStats();

    // Check if memory stats are available
    if (!stats) {
      console.warn('[MemoryMonitor] Memory stats unavailable - performance.memory API not supported in this browser');
      return;
    }

    const { used, total, percentage } = stats;

    // Validate percentage is a valid number
    // Note: percentage can be 0 at startup or when memory API returns 0, this is normal
    if (!isFinite(percentage) || percentage < 0) {
      // Only warn for truly invalid values (NaN, Infinity, negative)
      if (import.meta.env.DEV && (isNaN(percentage) || percentage < 0)) {
        console.warn('[MemoryMonitor] Invalid memory percentage calculated:', percentage);
      }
      return;
    }

    // Log memory usage to console (only in dev mode)
    if (import.meta.env.DEV) {
      console.debug(`Memory usage: ${used}MB / ${total}MB (${percentage}%)`);
    }

    // Use percentage-based thresholds (more reliable than absolute MB)
    const warningPercentage = 60; // Warning at 60% of heap limit
    const criticalPercentage = 80; // Critical at 80% of heap limit

    // Don't warn if using less than 500MB (even if high percentage of small heap)
    if (used < 500) {
      return;
    }

    // Track memory growth trend
    const memoryTrend = this.getMemoryTrend(used);

    // Check if we've reached critical threshold
    if (percentage >= criticalPercentage) {
      if (!this.criticalIssued) {
        console.error(`CRITICAL: Memory usage at ${percentage}% (${used}MB / ${total}MB)`);

        if (this.options.showToasts) {
          // Get current locale for localized message
          const locale = localStorage.getItem('app-locale') || 'ru';
          const message = locale === 'ru'
            ? `КРИТИЧЕСКОЕ использование памяти: ${used}МБ / ${total}МБ (${percentage}%). Браузер может зависнуть! Очистка...`
            : `CRITICAL memory usage: ${used}MB / ${total}MB (${percentage}%). Browser may crash! Cleaning up...`;

          this.uiStore.addToast(message, 'error', 15000);
        }

        // Dump heap snapshot for analysis
        this.dumpHeapSnapshot('critical-memory');

        // Force cleanup immediately
        if (this.options.autoCleanup) {
          this.forceCleanup();
          // Emergency additional cleanup
          setTimeout(() => this.forceCleanup(), 500);
        }

        // Trigger critical callbacks
        for (const cb of this.criticalCallbacks) {
          try {
            cb();
          } catch (e) {
            console.error('Critical callback error:', e);
          }
        }

        this.criticalIssued = true;
      }
    }
    // Check if we've reached warning threshold or rapid growth
    else if (percentage >= warningPercentage || (memoryTrend > 10 && this.memoryHistory.length >= 5)) {
      if (!this.warningIssued) {
        console.warn(`WARNING: Memory usage at ${percentage}% (${used}MB / ${total}MB)`);

        if (this.options.showToasts) {
          // Get current locale for localized message
          const locale = localStorage.getItem('app-locale') || 'ru';
          const message = locale === 'ru'
            ? `Использование памяти: ${used}МБ / ${total}МБ (${percentage}% от heap limit). Рекомендуется уменьшить выбор файлов.`
            : `Memory usage: ${used}MB / ${total}MB (${percentage}% of heap limit). Consider reducing file selection.`;

          this.uiStore.addToast(message, 'warning', 8000);
        }

        this.warningIssued = true;
      }
    }
    // Reset flags if memory usage drops below thresholds
    else {
      // Only reset if we drop significantly below thresholds
      if (percentage < warningPercentage - 10) {
        this.warningIssued = false;
        this.criticalIssued = false;
      }
    }
  }

  /**
   * Track memory trend to detect rapid growth
   */
  private getMemoryTrend(currentUsed: number): number {
    this.memoryHistory.push(currentUsed);

    // Keep only last 10 measurements
    if (this.memoryHistory.length > 10) {
      this.memoryHistory.shift();
    }

    // Calculate growth rate (MB per minute)
    if (this.memoryHistory.length < 2) {
      return 0;
    }

    const oldest = this.memoryHistory[0];
    const newest = this.memoryHistory[this.memoryHistory.length - 1];
    const timeDiff = (this.memoryHistory.length - 1) * (this.options.pollingInterval || 5000) / 60000; // minutes

    return (newest - oldest) / timeDiff;
  }
}

// Types are declared in types/performance.d.ts

// Convenience function to get the monitor instance
export function useMemoryMonitor(options?: MemoryWarningOptions): MemoryMonitor {
  return MemoryMonitor.getInstance(options);
}

export default MemoryMonitor;