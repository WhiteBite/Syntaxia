import { apiService } from '@/services/api.service'
import { getGitCache } from '../composables/useGitCache'
import type { GitBranch, GitCommit, GitProvider } from '../model/types'

/**
 * Git API layer - abstracts backend calls with caching
 */
export const gitApi = {
    // Provider detection
    async detectProvider(url: string): Promise<GitProvider> {
        const [isGitHub, isGitLab] = await Promise.all([
            apiService.isGitHubURL(url),
            apiService.isGitLabURL(url)
        ])

        if (isGitHub) return 'github'
        if (isGitLab) return 'gitlab'
        return 'unknown'
    },

    // Local repository
    async isGitRepository(projectPath: string): Promise<boolean> {
        return apiService.isGitRepository(projectPath)
    },

    async getCurrentBranch(projectPath: string): Promise<string> {
        return apiService.getCurrentBranch(projectPath)
    },

    async getBranches(projectPath: string): Promise<string[]> {
        const cache = getGitCache()
        return cache.cachedBranches(projectPath, async () => {
            const result = await apiService.getBranches(projectPath)
            return JSON.parse(result)
        })
    },

    async getCommits(projectPath: string, limit: number = 50): Promise<GitCommit[]> {
        const cache = getGitCache()
        return cache.cachedCommits(projectPath, 'HEAD', async () => {
            return apiService.getCommitHistory(projectPath, limit)
        }) as Promise<GitCommit[]>
    },

    async listFilesAtRef(projectPath: string, ref: string): Promise<string[]> {
        const cache = getGitCache()
        return cache.cachedFiles(projectPath, ref, async () => {
            const result = await apiService.listFilesAtRef(projectPath, ref)
            return Array.isArray(result) ? result : []
        })
    },

    async getFileAtRef(projectPath: string, filePath: string, ref: string): Promise<string> {
        const cache = getGitCache()
        return cache.cachedFileContent(projectPath, filePath, ref, async () => {
            return apiService.getFileAtRef(projectPath, filePath, ref)
        })
    },

    async buildContextAtRef(projectPath: string, files: string[], ref: string): Promise<string> {
        return apiService.buildContextAtRef(projectPath, files, ref)
    },

    // GitHub API
    async gitHubGetBranches(repoUrl: string): Promise<GitBranch[]> {
        const branches = await apiService.gitHubGetBranches(repoUrl)
        return branches.map(b => ({
            name: b.name,
            commit: { sha: b.commit.sha }
        }))
    },

    async gitHubGetDefaultBranch(repoUrl: string): Promise<string> {
        return apiService.gitHubGetDefaultBranch(repoUrl)
    },

    async gitHubListFiles(repoUrl: string, ref: string): Promise<string[]> {
        const cache = getGitCache()
        return cache.cachedFiles(repoUrl, ref, async () => {
            return apiService.gitHubListFiles(repoUrl, ref)
        })
    },

    async gitHubGetFileContent(repoUrl: string, filePath: string, ref: string): Promise<string> {
        const cache = getGitCache()
        return cache.cachedFileContent(repoUrl, filePath, ref, async () => {
            return apiService.gitHubGetFileContent(repoUrl, filePath, ref)
        })
    },

    async gitHubBuildContext(repoUrl: string, files: string[], ref: string): Promise<string> {
        return apiService.gitHubBuildContext(repoUrl, files, ref)
    },

    // GitLab API
    async gitLabGetBranches(repoUrl: string): Promise<GitBranch[]> {
        const branches = await apiService.gitLabGetBranches(repoUrl)
        return branches.map(b => ({
            name: b.name,
            commit: { sha: b.commit.id },
            isDefault: b.default
        }))
    },

    async gitLabGetDefaultBranch(repoUrl: string): Promise<string> {
        return apiService.gitLabGetDefaultBranch(repoUrl)
    },

    async gitLabListFiles(repoUrl: string, ref: string): Promise<string[]> {
        const cache = getGitCache()
        return cache.cachedFiles(repoUrl, ref, async () => {
            return apiService.gitLabListFiles(repoUrl, ref)
        })
    },

    async gitLabGetFileContent(repoUrl: string, filePath: string, ref: string): Promise<string> {
        const cache = getGitCache()
        return cache.cachedFileContent(repoUrl, filePath, ref, async () => {
            return apiService.gitLabGetFileContent(repoUrl, filePath, ref)
        })
    },

    async gitLabBuildContext(repoUrl: string, files: string[], ref: string): Promise<string> {
        return apiService.gitLabBuildContext(repoUrl, files, ref)
    },

    // Universal methods (auto-detect provider)
    async getRemoteBranches(repoUrl: string, provider: GitProvider): Promise<GitBranch[]> {
        switch (provider) {
            case 'github':
                return this.gitHubGetBranches(repoUrl)
            case 'gitlab':
                return this.gitLabGetBranches(repoUrl)
            default:
                throw new Error(`Unsupported provider: ${provider}`)
        }
    },

    async getRemoteDefaultBranch(repoUrl: string, provider: GitProvider): Promise<string> {
        switch (provider) {
            case 'github':
                return this.gitHubGetDefaultBranch(repoUrl)
            case 'gitlab':
                return this.gitLabGetDefaultBranch(repoUrl)
            default:
                throw new Error(`Unsupported provider: ${provider}`)
        }
    },

    async listRemoteFiles(repoUrl: string, ref: string, provider: GitProvider): Promise<string[]> {
        switch (provider) {
            case 'github':
                return this.gitHubListFiles(repoUrl, ref)
            case 'gitlab':
                return this.gitLabListFiles(repoUrl, ref)
            default:
                throw new Error(`Unsupported provider: ${provider}`)
        }
    },

    async getRemoteFileContent(repoUrl: string, filePath: string, ref: string, provider: GitProvider): Promise<string> {
        switch (provider) {
            case 'github':
                return this.gitHubGetFileContent(repoUrl, filePath, ref)
            case 'gitlab':
                return this.gitLabGetFileContent(repoUrl, filePath, ref)
            default:
                throw new Error(`Unsupported provider: ${provider}`)
        }
    },

    async buildRemoteContext(repoUrl: string, files: string[], ref: string, provider: GitProvider): Promise<string> {
        switch (provider) {
            case 'github':
                return this.gitHubBuildContext(repoUrl, files, ref)
            case 'gitlab':
                return this.gitLabBuildContext(repoUrl, files, ref)
            default:
                throw new Error(`Unsupported provider: ${provider}`)
        }
    },

    // Cache management
    invalidateCache(pattern?: string) {
        getGitCache().invalidate(pattern)
    },

    getCacheStats() {
        return getGitCache().getCacheStats()
    }
}
