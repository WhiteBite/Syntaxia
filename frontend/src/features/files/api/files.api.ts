import type { domain } from '#wailsjs/go/models';
import { apiService } from '@/services/api.service';

export interface FileStats {
    size: number
    lines: number
    tokens: number
    language?: string
    modTime?: string
}

class FilesApi {
    private fileTreeCache: Map<string, { data: domain.FileNode[]; timestamp: number }> = new Map()
    private readonly CACHE_TTL = 60000 // 1 minute
    private readonly MAX_CACHE_ENTRIES = 20 // Limit cache entries

    async listFiles(path: string, useGitignore: boolean = true, useCustomIgnore: boolean = true): Promise<domain.FileNode[]> {
        const cacheKey = `${path}-${useGitignore}-${useCustomIgnore}`
        const cached = this.fileTreeCache.get(cacheKey)

        if (cached && Date.now() - cached.timestamp < this.CACHE_TTL) {
            // Move to end (LRU)
            this.fileTreeCache.delete(cacheKey)
            this.fileTreeCache.set(cacheKey, cached)
            return cached.data
        }

        try {
            const files = await apiService.listFiles(path, useGitignore, useCustomIgnore)

            // Evict oldest entry if cache is full
            if (this.fileTreeCache.size >= this.MAX_CACHE_ENTRIES) {
                const firstKey = this.fileTreeCache.keys().next().value
                if (firstKey) this.fileTreeCache.delete(firstKey)
            }

            this.fileTreeCache.set(cacheKey, { data: files, timestamp: Date.now() })
            return files
        } catch (error) {
            console.error('[FilesApi] Failed to list files:', error)
            throw new Error('Failed to load file tree.')
        }
    }

    async readFileContent(projectPath: string, filePath: string): Promise<string> {
        try {
            return await apiService.readFileContent(projectPath, filePath)
        } catch (error) {
            console.error('[FilesApi] Failed to read file:', error)
            throw new Error('Failed to read file content.')
        }
    }

    async getFileStats(path: string): Promise<FileStats> {
        try {
            const stats = await apiService.getFileStats(path)
            return JSON.parse(stats) as FileStats
        } catch (error) {
            console.error('[FilesApi] Failed to get file stats:', error)
            throw new Error('Failed to get file statistics.')
        }
    }

    clearCache() {
        this.fileTreeCache.clear()
    }
}

export const filesApi = new FilesApi()
