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
import { forceMemoryCleanup } from './memory-cleanup';
import {
  type MemoryStats,
  type StoreMetrics,
  calculateMemoryTrend,
  dumpHeapSnapshot,
  getMemoryStats
} from './memory-stats';

const logger = useLogger('MemoryMonitor');

// Re-export types for backward compatibility
export type { MemoryStats, StoreMetrics };

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
    return getMemoryStats(this.isDevMode, this.storeImportsCache);
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
    forceMemoryCleanup();

    // Add a small delay to allow garbage collection to work
    setTimeout(() => {
      void this.getMemoryStats().then(stats => {
        if (stats) {
          logger.debug(`Memory after aggressive cleanup: ${stats.used}MB / ${stats.total}MB (${stats.percentage}%)`);
        }
      });
    }, 500);
  }

  /**
   * Dump heap snapshot for analysis in Chrome DevTools
   */
  public async dumpHeapSnapshot(reason: string): Promise<void> {
    await dumpHeapSnapshot(reason);
    const stats = await this.getMemoryStats();
    if (stats) {
      this.uiStore.addToast(`Heap snapshot saved: ${reason}`, 'info');
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
      logger.warn('Memory stats unavailable - performance.memory API not supported in this browser');
      return;
    }

    const { used, total, percentage } = stats;

    // Validate percentage is a valid number
    // Note: percentage can be 0 at startup or when memory API returns 0, this is normal
    if (!isFinite(percentage) || percentage < 0) {
      // Only warn for truly invalid values (NaN, Infinity, negative)
      if (import.meta.env.DEV && (isNaN(percentage) || percentage < 0)) {
        logger.warn('Invalid memory percentage calculated:', percentage);
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
        logger.error(`CRITICAL: Memory usage at ${percentage}% (${used}MB / ${total}MB)`);

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
            logger.error('Critical callback error:', e);
          }
        }

        this.criticalIssued = true;
      }
    }
    // Check if we've reached warning threshold or rapid growth
    else if (percentage >= warningPercentage || (memoryTrend > 10 && this.memoryHistory.length >= 5)) {
      if (!this.warningIssued) {
        logger.warn(`Memory usage at ${percentage}% (${used}MB / ${total}MB)`);

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
    return calculateMemoryTrend(
      this.memoryHistory,
      currentUsed,
      this.options.pollingInterval || 5000
    );
  }
}

// Types are declared in types/performance.d.ts

// Convenience function to get the monitor instance
export function useMemoryMonitor(options?: MemoryWarningOptions): MemoryMonitor {
  return MemoryMonitor.getInstance(options);
}

export default MemoryMonitor;