/**
 * Automatic Memory Diagnostics
 * 
 * Collects detailed memory diagnostics and saves to files for AI analysis
 */

import { useLogger } from '@/composables/useLogger'
import type { MemoryStats } from './memory-monitor'

const logger = useLogger('MemoryDiagnostics')

export interface DiagnosticSnapshot {
    timestamp: string
    memory: MemoryStats | null
    stores: {
        fileStore: {
            nodesCount: number
            selectedCount: number
            memoryUsage: number
            searchQueryLength: number
            filterExtensionsCount: number
        }
        contextStore: {
            hasContext: boolean
            fileCount: number
            totalSize: number
            lineCount: number
            cacheSize: number
        }
        apiCache: {
            entries: number
            size: number
            sizeMB: string
            hitRate: string
            hits: number
            misses: number
        }
    }
    performance: {
        navigation: PerformanceNavigationTiming | null
        resources: number
        measures: number
    }
    browser: {
        userAgent: string
        language: string
        platform: string
        hardwareConcurrency: number
        deviceMemory?: number
    }
}

export interface DiagnosticReport {
    startTime: string
    endTime: string
    duration: number
    snapshots: DiagnosticSnapshot[]
    issues: DiagnosticIssue[]
    recommendations: string[]
}

export interface DiagnosticDetails {
    previous?: number
    current?: number
    growth?: number
    used?: number
    total?: number
    percentage?: number
    nodesCount?: number
    cacheSize?: number
    entries?: number
    size?: number
    sizeMB?: string
    hitRate?: string
    hits?: number
    misses?: number
}

export interface DiagnosticIssue {
    severity: 'critical' | 'warning' | 'info'
    category: 'memory' | 'performance' | 'store' | 'cache'
    message: string
    details: DiagnosticDetails | MemoryStats | DiagnosticSnapshot['stores']['fileStore'] | DiagnosticSnapshot['stores']['contextStore'] | DiagnosticSnapshot['stores']['apiCache']
    timestamp: string
}

interface DiagnosticsStoreImportsCache {
    useFileStore?: () => {
        nodes?: unknown[]
        selectedCount?: number
        getMemoryUsage?: () => number
        searchQuery?: string
        filterExtensions?: string[]
    }
    useContextStore?: () => {
        hasContext?: boolean
        fileCount?: number
        totalSize?: number
        lineCount?: number
        getMemoryUsage?: () => number
    }
    getCacheStats?: () => {
        entries?: number
        size?: number
        sizeMB?: string
        hitRate?: string
        hits?: number
        misses?: number
    }
}

class MemoryDiagnostics {
    private static instance: MemoryDiagnostics
    private snapshots: DiagnosticSnapshot[] = []
    private issues: DiagnosticIssue[] = []
    private startTime: Date = new Date()
    private isCollecting = false
    private collectionInterval: number | null = null
    private isDevMode = import.meta.env.DEV

    private constructor() { }

    public static getInstance(): MemoryDiagnostics {
        if (!MemoryDiagnostics.instance) {
            MemoryDiagnostics.instance = new MemoryDiagnostics()
        }
        return MemoryDiagnostics.instance
    }

    /**
     * Start automatic diagnostics collection
     * In dev mode, use longer intervals to reduce memory overhead from dynamic imports
     */
    public startCollection(intervalMs: number = 10000): void {
        if (this.isCollecting) return

        this.isCollecting = true
        this.startTime = new Date()
        this.snapshots = []
        this.issues = []

        // In dev mode, use much longer interval to reduce overhead
        const actualInterval = this.isDevMode ? Math.max(intervalMs, 60000) : intervalMs

        logger.debug(`Started collection (interval: ${actualInterval}ms, dev: ${this.isDevMode})`)

        // Collect initial snapshot
        void this.collectSnapshot()

        // Collect periodically
        this.collectionInterval = window.setInterval(() => {
            void this.collectSnapshot().then(() => {
                this.analyzeSnapshots()
            })
        }, actualInterval)
    }

    /**
     * Stop collection
     */
    public stopCollection(): void {
        if (this.collectionInterval) {
            clearInterval(this.collectionInterval)
            this.collectionInterval = null
        }
        this.isCollecting = false
        logger.debug('Stopped collection')
    }

    /**
     * Collect a single diagnostic snapshot
     */
    public async collectSnapshot(): Promise<DiagnosticSnapshot> {
        const snapshot: DiagnosticSnapshot = {
            timestamp: new Date().toISOString(),
            memory: this.getMemoryStats(),
            stores: await this.getStoreStats(),
            performance: this.getPerformanceStats(),
            browser: this.getBrowserInfo()
        }

        this.snapshots.push(snapshot)

        // Keep only last 20 snapshots to reduce memory usage
        if (this.snapshots.length > 20) {
            this.snapshots.shift()
        }

        return snapshot
    }

    /**
     * Get memory stats
     */
    private getMemoryStats(): MemoryStats | null {
        if (!('performance' in window) || !performance.memory) {
            return null
        }

        const memory = performance.memory
        const used = Math.round(memory.usedJSHeapSize / (1024 * 1024))
        const total = Math.round(memory.jsHeapSizeLimit / (1024 * 1024))
        const percentage = Math.round((memory.usedJSHeapSize / memory.jsHeapSizeLimit) * 100)

        return { used, total, percentage }
    }

    // Cache store imports to avoid repeated dynamic imports causing memory leaks
    private storeImportsCache: DiagnosticsStoreImportsCache = {};

    /**
     * Get store statistics
     * In dev mode, skip to reduce dynamic import overhead
     */
    private async getStoreStats(): Promise<DiagnosticSnapshot['stores']> {
        const stats: DiagnosticSnapshot['stores'] = {
            fileStore: {
                nodesCount: 0,
                selectedCount: 0,
                memoryUsage: 0,
                searchQueryLength: 0,
                filterExtensionsCount: 0
            },
            contextStore: {
                hasContext: false,
                fileCount: 0,
                totalSize: 0,
                lineCount: 0,
                cacheSize: 0
            },
            apiCache: {
                entries: 0,
                size: 0,
                sizeMB: '0',
                hitRate: '0%',
                hits: 0,
                misses: 0
            }
        }

        // Skip detailed store stats in dev mode to prevent memory leaks from dynamic imports
        if (this.isDevMode) {
            return stats
        }

        try {
            // Cache imports to avoid repeated dynamic imports
            if (!this.storeImportsCache.useFileStore) {
                const { useFileStore } = await import('@/features/files/model/file.store')
                this.storeImportsCache.useFileStore = useFileStore as DiagnosticsStoreImportsCache['useFileStore']
            }
            if (!this.storeImportsCache.useContextStore) {
                const { useContextStore } = await import('@/features/context/model/context.store')
                this.storeImportsCache.useContextStore = useContextStore as DiagnosticsStoreImportsCache['useContextStore']
            }
            if (!this.storeImportsCache.getCacheStats) {
                const { getCacheStats } = await import('@/composables/useApiCache')
                this.storeImportsCache.getCacheStats = getCacheStats as DiagnosticsStoreImportsCache['getCacheStats']
            }

            const fileStore = this.storeImportsCache.useFileStore?.()
            if (fileStore) {
                stats.fileStore = {
                    nodesCount: fileStore.nodes?.length || 0,
                    selectedCount: fileStore.selectedCount || 0,
                    memoryUsage: fileStore.getMemoryUsage ? fileStore.getMemoryUsage() : 0,
                    searchQueryLength: fileStore.searchQuery?.length || 0,
                    filterExtensionsCount: fileStore.filterExtensions?.length || 0
                }
            }

            const contextStore = this.storeImportsCache.useContextStore?.()
            if (contextStore) {
                stats.contextStore = {
                    hasContext: contextStore.hasContext || false,
                    fileCount: contextStore.fileCount || 0,
                    totalSize: contextStore.totalSize || 0,
                    lineCount: contextStore.lineCount || 0,
                    cacheSize: contextStore.getMemoryUsage ? contextStore.getMemoryUsage() : 0
                }
            }

            const cacheStats = this.storeImportsCache.getCacheStats?.()
            if (cacheStats) {
                stats.apiCache = {
                    entries: cacheStats.entries || 0,
                    size: cacheStats.size || 0,
                    sizeMB: cacheStats.sizeMB || '0',
                    hitRate: cacheStats.hitRate || '0%',
                    hits: cacheStats.hits || 0,
                    misses: cacheStats.misses || 0
                }
            }
        } catch (e) {
            console.warn('[MemoryDiagnostics] Error collecting store stats:', e)
        }

        return stats
    }

    /**
     * Get performance statistics
     */
    private getPerformanceStats(): DiagnosticSnapshot['performance'] {
        const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming | null
        const resources = performance.getEntriesByType('resource').length
        const measures = performance.getEntriesByType('measure').length

        return { navigation, resources, measures }
    }

    /**
     * Get browser information
     */
    private getBrowserInfo(): DiagnosticSnapshot['browser'] {
        return {
            userAgent: navigator.userAgent,
            language: navigator.language,
            platform: navigator.platform,
            hardwareConcurrency: navigator.hardwareConcurrency,
            deviceMemory: navigator.deviceMemory
        }
    }

    /**
     * Analyze snapshots and detect issues
     */
    private analyzeSnapshots(): void {
        if (this.snapshots.length < 2) return

        const latest = this.snapshots[this.snapshots.length - 1]
        const previous = this.snapshots[this.snapshots.length - 2]

        // Check memory growth
        if (latest.memory && previous.memory) {
            const growth = latest.memory.used - previous.memory.used
            const growthPercent = (growth / previous.memory.used) * 100

            if (growth > 50) { // >50MB growth
                this.addIssue({
                    severity: 'warning',
                    category: 'memory',
                    message: `Memory grew by ${growth}MB (${growthPercent.toFixed(1)}%)`,
                    details: {
                        previous: previous.memory.used,
                        current: latest.memory.used,
                        growth
                    },
                    timestamp: latest.timestamp
                })
            }

            if (latest.memory.percentage >= 90) {
                this.addIssue({
                    severity: 'critical',
                    category: 'memory',
                    message: `Critical memory usage: ${latest.memory.percentage}%`,
                    details: latest.memory,
                    timestamp: latest.timestamp
                })
            }
        }

        // Check file store growth
        const nodeGrowth = latest.stores.fileStore.nodesCount - previous.stores.fileStore.nodesCount
        if (nodeGrowth > 1000) {
            this.addIssue({
                severity: 'warning',
                category: 'store',
                message: `File store grew by ${nodeGrowth} nodes`,
                details: {
                    previous: previous.stores.fileStore.nodesCount,
                    current: latest.stores.fileStore.nodesCount
                },
                timestamp: latest.timestamp
            })
        }

        // Check cache efficiency
        const hitRate = parseFloat(latest.stores.apiCache.hitRate)
        if (hitRate < 50 && latest.stores.apiCache.entries > 10) {
            this.addIssue({
                severity: 'info',
                category: 'cache',
                message: `Low cache hit rate: ${latest.stores.apiCache.hitRate}`,
                details: latest.stores.apiCache,
                timestamp: latest.timestamp
            })
        }

        // Check context cache size
        if (latest.stores.contextStore.cacheSize > 10 * 1024 * 1024) { // >10MB
            this.addIssue({
                severity: 'warning',
                category: 'store',
                message: `Large context cache: ${(latest.stores.contextStore.cacheSize / (1024 * 1024)).toFixed(1)}MB`,
                details: latest.stores.contextStore,
                timestamp: latest.timestamp
            })
        }
    }

    /**
     * Add diagnostic issue
     */
    private addIssue(issue: DiagnosticIssue): void {
        this.issues.push(issue)

        // Keep only last 50 issues
        if (this.issues.length > 50) {
            this.issues.shift()
        }

        logger.debug(`${issue.severity.toUpperCase()}: ${issue.message}`, issue.details)
    }

    /**
     * Generate diagnostic report
     */
    public generateReport(): DiagnosticReport {
        const endTime = new Date()
        const duration = endTime.getTime() - this.startTime.getTime()

        const recommendations = this.generateRecommendations()

        return {
            startTime: this.startTime.toISOString(),
            endTime: endTime.toISOString(),
            duration,
            snapshots: this.snapshots,
            issues: this.issues,
            recommendations
        }
    }

    /**
     * Generate recommendations based on collected data
     */
    private generateRecommendations(): string[] {
        const recommendations: string[] = []

        if (this.snapshots.length === 0) {
            return ['No data collected yet']
        }

        const latest = this.snapshots[this.snapshots.length - 1]

        // Memory recommendations
        if (latest.memory && latest.memory.percentage > 75) {
            recommendations.push('Memory usage is high. Consider reducing file selection or context size.')
        }

        // File store recommendations
        if (latest.stores.fileStore.nodesCount > 5000) {
            recommendations.push('Large file tree detected. Consider using file filters or working in subdirectories.')
        }

        if (latest.stores.fileStore.selectedCount > 100) {
            recommendations.push('Many files selected. This may impact performance.')
        }

        // Context recommendations
        if (latest.stores.contextStore.cacheSize > 15 * 1024 * 1024) {
            recommendations.push('Context cache is large. Clear context when not needed.')
        }

        // Cache recommendations
        const hitRate = parseFloat(latest.stores.apiCache.hitRate)
        if (hitRate < 50 && latest.stores.apiCache.entries > 10) {
            recommendations.push('Low cache hit rate. Cache may not be effective.')
        }

        // Growth trend
        if (this.snapshots.length >= 5) {
            const first = this.snapshots[0]
            const memoryGrowth = latest.memory && first.memory
                ? latest.memory.used - first.memory.used
                : 0

            if (memoryGrowth > 100) {
                recommendations.push(`Memory grew by ${memoryGrowth}MB since start. Possible memory leak.`)
            }
        }

        if (recommendations.length === 0) {
            recommendations.push('No issues detected. Memory usage is normal.')
        }

        return recommendations
    }

    /**
     * Save report to file (as JSON for AI to read)
     */
    public async saveReportToFile(): Promise<string> {
        const report = this.generateReport()
        const filename = `memory-diagnostics-${Date.now()}.json`

        try {
            // Create blob and download
            const blob = new Blob([JSON.stringify(report, null, 2)], { type: 'application/json' })
            const url = URL.createObjectURL(blob)
            const a = document.createElement('a')
            a.href = url
            a.download = filename
            a.click()
            URL.revokeObjectURL(url)

            logger.debug('Report saved:', filename)
            return filename
        } catch (e) {
            console.error('[MemoryDiagnostics] Failed to save report:', e)
            throw e
        }
    }

    /**
     * Get current status summary
     */
    public getStatusSummary(): string {
        if (this.snapshots.length === 0) {
            return 'No diagnostics collected yet'
        }

        const latest = this.snapshots[this.snapshots.length - 1]
        const criticalIssues = this.issues.filter(i => i.severity === 'critical').length
        const warnings = this.issues.filter(i => i.severity === 'warning').length

        let summary = `Diagnostics Summary:\n`
        summary += `- Snapshots collected: ${this.snapshots.length}\n`
        summary += `- Critical issues: ${criticalIssues}\n`
        summary += `- Warnings: ${warnings}\n`

        if (latest.memory) {
            summary += `- Current memory: ${latest.memory.used}MB / ${latest.memory.total}MB (${latest.memory.percentage}%)\n`
        }

        summary += `- File nodes: ${latest.stores.fileStore.nodesCount}\n`
        summary += `- Selected files: ${latest.stores.fileStore.selectedCount}\n`
        summary += `- Context cache: ${(latest.stores.contextStore.cacheSize / (1024 * 1024)).toFixed(1)}MB\n`
        summary += `- API cache: ${latest.stores.apiCache.entries} entries (${latest.stores.apiCache.sizeMB}MB)\n`

        return summary
    }

    /**
     * Get latest snapshot
     */
    public getLatestSnapshot(): DiagnosticSnapshot | null {
        return this.snapshots.length > 0 ? this.snapshots[this.snapshots.length - 1] : null
    }

    /**
     * Get all issues
     */
    public getIssues(): DiagnosticIssue[] {
        return this.issues
    }

    /**
     * Clear all data
     */
    public clear(): void {
        this.snapshots = []
        this.issues = []
        this.startTime = new Date()
    }
}

export function useMemoryDiagnostics(): MemoryDiagnostics {
    return MemoryDiagnostics.getInstance()
}

export default MemoryDiagnostics
