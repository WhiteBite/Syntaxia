import { ref } from 'vue'

interface CacheEntry<T> {
    data: T
    timestamp: number
    ttl: number
}

const cache = new Map<string, CacheEntry<unknown>>()

// Default TTL: 5 minutes
const DEFAULT_TTL = 5 * 60 * 1000

/**
 * Git data caching composable for performance optimization
 * Caches API responses to reduce network requests
 */
export function useGitCache() {
    const isLoading = ref(false)

    function getCacheKey(type: string, ...args: string[]): string {
        return `git:${type}:${args.join(':')}`
    }

    function get<T>(key: string): T | null {
        const entry = cache.get(key) as CacheEntry<T> | undefined
        if (!entry) return null

        // Check if expired
        if (Date.now() - entry.timestamp > entry.ttl) {
            cache.delete(key)
            return null
        }

        return entry.data
    }

    function set<T>(key: string, data: T, ttl: number = DEFAULT_TTL): void {
        cache.set(key, {
            data,
            timestamp: Date.now(),
            ttl
        })
    }

    function invalidate(pattern?: string): void {
        if (!pattern) {
            cache.clear()
            return
        }

        for (const key of cache.keys()) {
            if (key.includes(pattern)) {
                cache.delete(key)
            }
        }
    }

    /**
     * Cached fetch with automatic cache management
     */
    async function cachedFetch<T>(
        key: string,
        fetcher: () => Promise<T>,
        ttl: number = DEFAULT_TTL
    ): Promise<T> {
        // Check cache first
        const cached = get<T>(key)
        if (cached !== null) {
            return cached
        }

        // Fetch and cache
        isLoading.value = true
        try {
            const data = await fetcher()
            set(key, data, ttl)
            return data
        } finally {
            isLoading.value = false
        }
    }

    // Specific cache methods for Git data
    async function cachedBranches(
        projectPath: string,
        fetcher: () => Promise<string[]>
    ): Promise<string[]> {
        const key = getCacheKey('branches', projectPath)
        return cachedFetch(key, fetcher, 2 * 60 * 1000) // 2 min TTL
    }

    async function cachedCommits(
        projectPath: string,
        ref: string,
        fetcher: () => Promise<unknown[]>
    ): Promise<unknown[]> {
        const key = getCacheKey('commits', projectPath, ref)
        return cachedFetch(key, fetcher, 1 * 60 * 1000) // 1 min TTL
    }

    async function cachedFiles(
        source: string,
        ref: string,
        fetcher: () => Promise<string[]>
    ): Promise<string[]> {
        const key = getCacheKey('files', source, ref)
        return cachedFetch(key, fetcher, 5 * 60 * 1000) // 5 min TTL
    }

    async function cachedFileContent(
        source: string,
        filePath: string,
        ref: string,
        fetcher: () => Promise<string>
    ): Promise<string> {
        const key = getCacheKey('content', source, filePath, ref)
        return cachedFetch(key, fetcher, 10 * 60 * 1000) // 10 min TTL
    }

    function getCacheStats() {
        let totalSize = 0
        let entryCount = 0
        let expiredCount = 0
        const now = Date.now()

        for (const [, entry] of cache) {
            entryCount++
            if (now - entry.timestamp > entry.ttl) {
                expiredCount++
            }
            // Rough size estimate
            totalSize += JSON.stringify(entry.data).length
        }

        return {
            entryCount,
            expiredCount,
            totalSizeKB: Math.round(totalSize / 1024)
        }
    }

    return {
        isLoading,
        get,
        set,
        invalidate,
        cachedFetch,
        cachedBranches,
        cachedCommits,
        cachedFiles,
        cachedFileContent,
        getCacheStats
    }
}

// Singleton instance for shared cache
let sharedCache: ReturnType<typeof useGitCache> | null = null

export function getGitCache() {
    if (!sharedCache) {
        sharedCache = useGitCache()
    }
    return sharedCache
}
