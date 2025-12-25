// Performance Optimization Utilities

import { useLogger } from '@/composables/useLogger'
import { computed, ref } from 'vue'

const logger = useLogger('Performance')

/**
 * Virtual Scrolling Utilities
 */

export interface VirtualScrollOptions {
  itemHeight: number
  containerHeight: number
  buffer: number
  threshold: number
}

export function useVirtualScroll<T>(
  items: T[],
  options: VirtualScrollOptions
) {
  const scrollTop = ref(0)
  const containerHeight = options.containerHeight
  const itemHeight = options.itemHeight
  const buffer = options.buffer || 5

  const visibleStart = computed(() => {
    return Math.max(0, Math.floor(scrollTop.value / itemHeight) - buffer)
  })

  const visibleEnd = computed(() => {
    const visibleCount = Math.ceil(containerHeight / itemHeight)
    return Math.min(items.length, visibleStart.value + visibleCount + buffer)
  })

  const visibleItems = computed(() => {
    return items.slice(visibleStart.value, visibleEnd.value)
  })

  const totalHeight = computed(() => {
    return items.length * itemHeight
  })

  const offsetY = computed(() => {
    return visibleStart.value * itemHeight
  })

  function handleScroll(event: Event) {
    const target = event.target as HTMLElement
    scrollTop.value = target.scrollTop
  }

  return {
    visibleItems,
    visibleStart,
    visibleEnd,
    totalHeight,
    offsetY,
    handleScroll
  }
}

/**
 * Debounce and Throttle Utilities
 */

export function debounce<T extends (...args: unknown[]) => unknown>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout | null = null

  return (...args: Parameters<T>) => {
    if (timeout) {
      clearTimeout(timeout)
    }
    timeout = setTimeout(() => func(...args), wait)
  }
}

export function throttle<T extends (...args: unknown[]) => unknown>(
  func: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle = false

  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args)
      inThrottle = true
      setTimeout(() => (inThrottle = false), limit)
    }
  }
}

/**
 * Lazy Loading Utilities
 */

export function useLazyLoading(threshold: number = 100) {
  const observer = ref<IntersectionObserver | null>(null)
  const loadedItems = ref(new Set<string>())

  function observe(element: HTMLElement, id: string, callback: () => void) {
    if (!observer.value) {
      observer.value = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              const itemId = entry.target.getAttribute('data-id')
              if (itemId && !loadedItems.value.has(itemId)) {
                loadedItems.value.add(itemId)
                callback()
                observer.value?.unobserve(entry.target)
              }
            }
          })
        },
        {
          rootMargin: `${threshold}px`
        }
      )
    }

    element.setAttribute('data-id', id)
    observer.value.observe(element)
  }

  function disconnect() {
    observer.value?.disconnect()
    observer.value = null
    loadedItems.value.clear()
  }

  return {
    observe,
    disconnect,
    loadedItems: computed(() => Array.from(loadedItems.value))
  }
}

/**
 * Memory Management Utilities
 */

interface CacheItem<T> {
  value: T
  timestamp: number
}

export class MemoryManager {
  private cache = new Map<string, CacheItem<unknown>>()
  private maxSize: number
  private ttl: number

  constructor(maxSize: number = 100, ttl: number = 5 * 60 * 1000) {
    this.maxSize = maxSize
    this.ttl = ttl
  }

  set<T>(key: string, value: T): void {
    // Remove oldest entries if cache is full
    if (this.cache.size >= this.maxSize) {
      const firstKey = this.cache.keys().next().value
      if (firstKey !== undefined) {
        this.cache.delete(firstKey)
      }
    }

    this.cache.set(key, {
      value,
      timestamp: Date.now()
    })
  }

  get<T>(key: string): T | null {
    const item = this.cache.get(key)
    if (!item) return null

    // Check if item has expired
    if (Date.now() - item.timestamp > this.ttl) {
      this.cache.delete(key)
      return null
    }

    return item.value as T
  }

  delete(key: string): boolean {
    return this.cache.delete(key)
  }

  clear(): void {
    this.cache.clear()
  }

  size(): number {
    return this.cache.size
  }

  cleanup(): void {
    const now = Date.now()
    for (const [key, item] of this.cache.entries()) {
      if (now - item.timestamp > this.ttl) {
        this.cache.delete(key)
      }
    }
  }
}

/**
 * Performance Monitoring
 */

export class PerformanceMonitor {
  private measurements = new Map<string, number>()
  private marks = new Map<string, number>()

  startMeasurement(name: string): void {
    this.marks.set(name, performance.now())
  }

  endMeasurement(name: string): number {
    const startTime = this.marks.get(name)
    if (!startTime) {
      console.warn(`No start mark found for measurement: ${name}`)
      return 0
    }

    const duration = performance.now() - startTime
    this.measurements.set(name, duration)
    this.marks.delete(name)

    return duration
  }

  getMeasurement(name: string): number | undefined {
    return this.measurements.get(name)
  }

  getAllMeasurements(): Record<string, number> {
    return Object.fromEntries(this.measurements)
  }

  logMeasurement(name: string): void {
    const duration = this.getMeasurement(name)
    if (duration !== undefined) {
      logger.debug(`Performance [${name}]: ${duration.toFixed(2)}ms`)
    }
  }

  clear(): void {
    this.measurements.clear()
    this.marks.clear()
  }
}

/**
 * Component Performance Wrapper
 */

export function withPerformanceMonitoring<T extends (...args: unknown[]) => unknown>(
  func: T,
  name: string,
  monitor: PerformanceMonitor
): T {
  return ((...args: unknown[]) => {
    monitor.startMeasurement(name)
    const result = func(...args)

    if (result && typeof (result as Promise<unknown>).then === 'function') {
      // Handle promises
      return (result as Promise<unknown>).finally(() => {
        monitor.endMeasurement(name)
        monitor.logMeasurement(name)
      })
    } else {
      monitor.endMeasurement(name)
      monitor.logMeasurement(name)
      return result
    }
  }) as T
}

/**
 * Bundle Size Optimization
 */

interface FallbackComponent {
  template: string
}

export function asyncComponent(loader: () => Promise<unknown>): () => Promise<unknown | FallbackComponent> {
  return () => {
    return loader().catch(error => {
      console.error('Failed to load component:', error)
      // Return fallback component
      return {
        template: '<div class="error-loading">Failed to load component</div>'
      }
    })
  }
}

/**
 * Request Batching for API Calls
 */

interface BatchRequest<T, D> {
  resolve: (value: T) => void
  reject: (error: unknown) => void
  data: D
}

interface BatchEntry<T, D> {
  requests: Array<BatchRequest<T, D>>
  timeout: NodeJS.Timeout
  handler: (requests: D[]) => Promise<T[]>
}

export class RequestBatcher {
  private batches = new Map<string, BatchEntry<unknown, unknown>>()
  private batchDelay: number

  constructor(batchDelay: number = 50) {
    this.batchDelay = batchDelay
  }

  batch<T, D = unknown>(
    key: string,
    data: D,
    batchHandler: (requests: D[]) => Promise<T[]>
  ): Promise<T> {
    return new Promise((resolve, reject) => {
      let batch = this.batches.get(key) as BatchEntry<T, D> | undefined

      if (!batch) {
        batch = {
          requests: [],
          timeout: setTimeout(() => {
            this.executeBatch<T, D>(key)
          }, this.batchDelay),
          handler: batchHandler
        }
        this.batches.set(key, batch as BatchEntry<unknown, unknown>)
      }

      batch.requests.push({ resolve, reject, data })
    })
  }

  private async executeBatch<T, D>(key: string): Promise<void> {
    const batch = this.batches.get(key) as BatchEntry<T, D> | undefined
    if (!batch) return

    this.batches.delete(key)
    clearTimeout(batch.timeout)

    try {
      const requestData = batch.requests.map(req => req.data)
      const results = await batch.handler(requestData)

      batch.requests.forEach((request, index) => {
        request.resolve(results[index])
      })
    } catch (error) {
      batch.requests.forEach(request => {
        request.reject(error)
      })
    }
  }
}

// Global instances
export const globalMemoryManager = new MemoryManager()
export const globalPerformanceMonitor = new PerformanceMonitor()
export const globalRequestBatcher = new RequestBatcher()

// Cleanup function for component unmounting
export function cleanupPerformanceUtils(): void {
  globalMemoryManager.cleanup()
  globalPerformanceMonitor.clear()
}