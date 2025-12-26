/**
 * GitHub API (no clone required)
 * Handles GitHub repository operations via API
 */

import * as wails from '#wailsjs/go/main/App'
import type { GitHubBranch, GitHubCommit } from '../types'
import { apiCall, apiCallWithDefault, parseJsonResponse } from './base'

export const githubApi = {
    isGitHubURL: (url: string): Promise<boolean> =>
        apiCallWithDefault(() => wails.IsGitHubURL(url), false, 'github.isGitHubURL'),

    getDefaultBranch: (repoURL: string): Promise<string> =>
        apiCall(() => wails.GitHubGetDefaultBranch(repoURL), 'Failed to get default branch.', { logContext: 'github' }),

    getBranches: async (repoURL: string): Promise<GitHubBranch[]> => {
        const result = await apiCall(
            () => wails.GitHubGetBranches(repoURL),
            'Failed to get GitHub branches.',
            { logContext: 'github' }
        )
        return parseJsonResponse(result, 'Failed to parse GitHub branches.')
    },

    getCommits: async (repoURL: string, branch: string, limit = 50): Promise<GitHubCommit[]> => {
        const result = await apiCall(
            () => wails.GitHubGetCommits(repoURL, branch, limit),
            'Failed to get GitHub commits.',
            { logContext: 'github' }
        )
        return parseJsonResponse(result, 'Failed to parse GitHub commits.')
    },

    listFiles: async (repoURL: string, ref: string): Promise<string[]> => {
        const result = await apiCall(
            () => wails.GitHubListFiles(repoURL, ref),
            'Failed to list GitHub files.',
            { logContext: 'github' }
        )
        const parsed = parseJsonResponse<unknown>(result, 'Failed to parse GitHub files.')
        return Array.isArray(parsed) ? parsed : []
    },

    getFileContent: (repoURL: string, filePath: string, ref: string): Promise<string> =>
        apiCall(
            () => wails.GitHubGetFileContent(repoURL, filePath, ref),
            'Failed to get GitHub file content.',
            { logContext: 'github' }
        ),

    buildContext: (repoURL: string, files: string[], ref: string): Promise<string> =>
        apiCall(
            () => wails.GitHubBuildContext(repoURL, files, ref),
            'Failed to build GitHub context.',
            { logContext: 'github' }
        ),
}
