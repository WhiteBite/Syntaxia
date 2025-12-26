import { a as useLogger } from "./index-CCibvNLt.js";
import { r as ref } from "./vue-vendor-DLLjc1XI.js";
import "./ui-vendor-2hF-XK7S.js";
import "./utils-DVYXdora.js";
import "./icons-BGgRd9bo.js";
import "./highlighter-D41BBaGQ.js";
const logger = useLogger("ApiCache");
const cache = /* @__PURE__ */ new Map();
const MAX_CACHE_SIZE = 20 * 1024 * 1024;
const MAX_CACHE_ENTRIES = 50;
let currentCacheSize = 0;
let cacheHits = 0;
let cacheMisses = 0;
const DEFAULT_TTL = 2 * 60 * 1e3;
function estimateSize(data) {
  const str = JSON.stringify(data);
  return str.length * 2;
}
function evictIfNeeded(newEntrySize) {
  const now = Date.now();
  const expiredKeys = [];
  cache.forEach((entry, key) => {
    if (now - entry.timestamp > DEFAULT_TTL || now - entry.timestamp > 6e4) {
      expiredKeys.push(key);
    }
  });
  for (const key of expiredKeys) {
    const entry = cache.get(key);
    if (entry) {
      currentCacheSize -= entry.size;
    }
    cache.delete(key);
  }
  const entriesToLog = [];
  while ((currentCacheSize + newEntrySize > MAX_CACHE_SIZE || cache.size >= MAX_CACHE_ENTRIES) && cache.size > 0) {
    const firstKey = cache.keys().next().value;
    if (!firstKey) break;
    const entry = cache.get(firstKey);
    if (entry) {
      currentCacheSize -= entry.size;
      entriesToLog.push({ key: firstKey, size: entry.size });
    }
    cache.delete(firstKey);
  }
  if (entriesToLog.length > 0) {
    entriesToLog.sort((a, b) => b.size - a.size).slice(0, 5).map((e) => `${e.key}: ${(e.size / 1024).toFixed(1)}KB`);
  }
  if (currentCacheSize > MAX_CACHE_SIZE * 0.8) {
    const entriesToRemove = Math.ceil(cache.size * 0.3);
    let removed = 0;
    const keysToRemove = [];
    cache.forEach((_entry, key) => {
      if (removed >= entriesToRemove) return;
      keysToRemove.push(key);
      removed++;
    });
    for (const key of keysToRemove) {
      const entry = cache.get(key);
      if (entry) {
        currentCacheSize -= entry.size;
      }
      cache.delete(key);
    }
  }
}
function getCacheStats() {
  const hitRate = cacheHits + cacheMisses > 0 ? (cacheHits / (cacheHits + cacheMisses) * 100).toFixed(1) : "0";
  return {
    size: currentCacheSize,
    sizeMB: (currentCacheSize / (1024 * 1024)).toFixed(2),
    entries: cache.size,
    maxSize: MAX_CACHE_SIZE,
    maxEntries: MAX_CACHE_ENTRIES,
    hitRate: `${hitRate}%`,
    hits: cacheHits,
    misses: cacheMisses
  };
}
function getMemoryUsage() {
  return currentCacheSize;
}
function clearAllCaches() {
  cache.clear();
  currentCacheSize = 0;
  cacheHits = 0;
  cacheMisses = 0;
}
function useApiCache(key, fetcher, ttl = DEFAULT_TTL) {
  const data = ref(null);
  const isLoading = ref(false);
  const error = ref(null);
  const load = async (force = false) => {
    if (!force) {
      const cached = cache.get(key);
      if (cached && Date.now() - cached.timestamp < ttl) {
        cacheHits++;
        cache.delete(key);
        cache.set(key, cached);
        data.value = cached.data;
        return cached.data;
      } else if (cached) {
        cacheMisses++;
      }
    }
    isLoading.value = true;
    error.value = null;
    try {
      const result = await fetcher();
      data.value = result;
      const size = estimateSize(result);
      evictIfNeeded(size);
      const entry = {
        data: result,
        timestamp: Date.now(),
        size
      };
      cache.set(key, entry);
      currentCacheSize += size;
      return result;
    } catch (e) {
      error.value = e;
      logger.error(`API cache error for key "${key}":`, e);
      throw e;
    } finally {
      isLoading.value = false;
    }
  };
  const invalidate = () => {
    const entry = cache.get(key);
    if (entry) {
      currentCacheSize -= entry.size;
    }
    cache.delete(key);
    data.value = null;
  };
  const invalidateAll = () => {
    cache.clear();
    currentCacheSize = 0;
  };
  return {
    data,
    isLoading,
    error,
    load,
    invalidate,
    invalidateAll
  };
}
export {
  clearAllCaches,
  getCacheStats,
  getMemoryUsage,
  useApiCache
};
