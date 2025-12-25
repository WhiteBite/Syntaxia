/**
 * Composable for caching API requests with LRU eviction
 * Reduces unnecessary network calls and improves performance
 */

import { useLogger } from '@/composables/useLogger'
import { ref } from 'vue'

const logger = useLogger('ApiCache')

interface CacheEntry<T> {
  data: T
  timestamp: number
  size: number // Estimated size in bytes
}

// Global cache with LRU eviction (reduced limits for memory safety)
const cache = new Map<string, CacheEntry<unknown>>()
const MAX_CACHE_SIZE = 20 * 1024 * 1024 // 20 MB max cache size (reduced from 100MB)
const MAX_CACHE_ENTRIES = 50 // Max number of entries (reduced from 100)
let currentCacheSize = 0

// Cache statistics (separate from Map to avoid any types)
let cacheHits = 0
let cacheMisses = 0

const DEFAULT_TTL = 2 * 60 * 1000 // 2 minutes (reduced from 5 minutes)

// Estimate size of data in bytes
function estimateSize(data: unknown): number {
  const str = JSON.stringify(data)
  return str.length * 2 // UTF-16 uses 2 bytes per character
}

// Evict oldest entries when cache is full (more aggressive)
function evictIfNeeded(newEntrySize: number) {
  // Remove expired entries first
  const now = Date.now()
  const expiredKeys: string[] = []

  cache.forEach((entry, key) => {
    // More aggressive: also remove entries not used in last 1 minute
    if (now - entry.timestamp > DEFAULT_TTL || now - entry.timestamp > 60000) {
      expiredKeys.push(key)
    }
  })

  for (const key of expiredKeys) {
    const entry = cache.get(key)
    if (entry) {
      currentCacheSize -= entry.size
    }
    cache.delete(key)
  }

  // If still over limit, remove oldest entries (LRU)
  const entriesToLog: Array<{ key: string; size: number }> = []

  while (
    (currentCacheSize + newEntrySize > MAX_CACHE_SIZE || cache.size >= MAX_CACHE_ENTRIES) &&
    cache.size > 0
  ) {
    const firstKey = cache.keys().next().value
    if (!firstKey) break

    const entry = cache.get(firstKey)
    if (entry) {
      currentCacheSize -= entry.size
      entriesToLog.push({ key: firstKey, size: entry.size })
    }
    cache.delete(firstKey)
  }

  // Log top 5 largest evicted entries
  if (entriesToLog.length > 0) {
    const topEntries = entriesToLog
      .sort((a, b) => b.size - a.size)
      .slice(0, 5)
      .map(e => `${e.key}: ${(e.size / 1024).toFixed(1)}KB`)
    logger.debug(`Evicted ${entriesToLog.length} entries. Top 5:`, topEntries)
  }

  // Auto-cleanup at 80% threshold
  if (currentCacheSize > MAX_CACHE_SIZE * 0.8) {
    const entriesToRemove = Math.ceil(cache.size * 0.3)
    let removed = 0
    const keysToRemove: string[] = []

    cache.forEach((_entry, key) => {
      if (removed >= entriesToRemove) return
      keysToRemove.push(key)
      removed++
    })

    for (const key of keysToRemove) {
      const entry = cache.get(key)
      if (entry) {
        currentCacheSize -= entry.size
      }
      cache.delete(key)
    }
    logger.debug(`Auto-cleanup: removed ${removed} entries`)
  }
}

// Get cache stats for monitoring
export function getCacheStats() {
  const hitRate = cacheHits + cacheMisses > 0
    ? (cacheHits / (cacheHits + cacheMisses) * 100).toFixed(1)
    : '0'

  return {
    size: currentCacheSize,
    sizeMB: (currentCacheSize / (1024 * 1024)).toFixed(2),
    entries: cache.size,
    maxSize: MAX_CACHE_SIZE,
    maxEntries: MAX_CACHE_ENTRIES,
    hitRate: `${hitRate}%`,
    hits: cacheHits,
    misses: cacheMisses
  }
}

// Get memory usage
export function getMemoryUsage(): number {
  return currentCacheSize
}

// Clear all caches (exported for emergency cleanup)
export function clearAllCaches() {
  cache.clear()
  currentCacheSize = 0
  cacheHits = 0
  cacheMisses = 0
  logger.debug('All caches cleared')
}

export function useApiCache<T>(
  key: string,
  fetcher: () => Promise<T>,
  ttl = DEFAULT_TTL
) {
  const data = ref<T | null>(null)
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  const load = async (force = false): Promise<T | null> => {
    // Check cache first (unless force refresh)
    if (!force) {
      const cached = cache.get(key)
      if (cached && Date.now() - cached.timestamp < ttl) {
        // Track hit
        cacheHits++
        // Move to end (LRU)
        cache.delete(key)
        cache.set(key, cached)
        data.value = cached.data
        return cached.data as T
      } else if (cached) {
        // Track miss (expired)
        cacheMisses++
      }
    }

    // Fetch fresh data
    isLoading.value = true
    error.value = null

    try {
      const result = await fetcher()
      data.value = result

      // Estimate size and evict if needed
      const size = estimateSize(result)
      evictIfNeeded(size)

      // Store in cache
      const entry: CacheEntry<T> = {
        data: result,
        timestamp: Date.now(),
        size
      }
      cache.set(key, entry)
      currentCacheSize += size

      return result
    } catch (e) {
      error.value = e as Error
      console.error(`API cache error for key "${key}":`, e)
      throw e
    } finally {
      isLoading.value = false
    }
  }

  const invalidate = () => {
    const entry = cache.get(key)
    if (entry) {
      currentCacheSize -= entry.size
    }
    cache.delete(key)
    data.value = null
  }

  const invalidateAll = () => {
    cache.clear()
    currentCacheSize = 0
  }

  return {
    data,
    isLoading,
    error,
    load,
    invalidate,
    invalidateAll
  }
}
