/**
 * Composable for file quick info with caching
 * Phase 4: Quick File Info
 */
import * as wails from '#wailsjs/go/main/App'
import { ref } from 'vue'
import { useLogger } from './useLogger'

export interface FileQuickInfo {
    symbolCount: number
    importCount: number
    dependentCount: number
    changeRisk: number
    riskLevel: 'low' | 'medium' | 'high'
}

interface CacheEntry {
    info: FileQuickInfo
    timestamp: number
}

const CACHE_TTL = 5 * 60 * 1000 // 5 minutes

// Global cache shared across components
const cache = new Map<string, CacheEntry>()
const loading = ref<Set<string>>(new Set())

export function useFileQuickInfo() {
    const logger = useLogger('FileQuickInfo')

    /**
     * Get quick info for a file (with caching)
     */
    async function getInfo(projectPath: string, filePath: string): Promise<FileQuickInfo | null> {
        const cacheKey = `${projectPath}:${filePath}`

        // Check cache
        const cached = cache.get(cacheKey)
        if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
            return cached.info
        }

        // Prevent duplicate requests
        if (loading.value.has(cacheKey)) {
            return null
        }

        loading.value.add(cacheKey)
        try {
            // @ts-ignore - method may not exist in wails bindings yet
            const result = await wails.GetFileQuickInfo(projectPath, filePath)
            const info: FileQuickInfo = {
                symbolCount: result.symbolCount || 0,
                importCount: result.importCount || 0,
                dependentCount: result.dependentCount || 0,
                changeRisk: result.changeRisk || 0,
                riskLevel: (result.riskLevel as 'low' | 'medium' | 'high') || 'low'
            }

            cache.set(cacheKey, { info, timestamp: Date.now() })
            return info
        } catch (error) {
            logger.warn('Error fetching info:', error)
            return null
        } finally {
            loading.value.delete(cacheKey)
        }
    }

    /**
     * Check if info is being loaded
     */
    function isLoading(projectPath: string, filePath: string): boolean {
        return loading.value.has(`${projectPath}:${filePath}`)
    }

    /**
     * Clear cache for a specific file
     */
    function invalidate(projectPath: string, filePath: string): void {
        cache.delete(`${projectPath}:${filePath}`)
    }

    /**
     * Clear entire cache
     */
    function clearCache(): void {
        cache.clear()
    }

    /**
     * Get risk color based on level
     */
    function getRiskColor(level: string): string {
        switch (level) {
            case 'low': return 'text-green-400'
            case 'medium': return 'text-amber-400'
            case 'high': return 'text-red-400'
            default: return 'text-gray-400'
        }
    }

    /**
     * Get risk bar width percentage
     */
    function getRiskWidth(risk: number): string {
        return `${Math.round(risk * 100)}%`
    }

    return {
        getInfo,
        isLoading,
        invalidate,
        clearCache,
        getRiskColor,
        getRiskWidth
    }
}
