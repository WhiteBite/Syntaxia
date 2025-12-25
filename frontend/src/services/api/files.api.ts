/**
 * File operations API
 * Handles file tree, reading files, file stats
 */

import * as wails from '#wailsjs/go/main/App'
import type { domain } from '#wailsjs/go/models'
import { apiCall } from './base'

export const filesApi = {
    listFiles: (path: string, useGitignore = true, useCustomIgnore = true): Promise<domain.FileNode[]> =>
        apiCall(() => wails.ListFiles(path, useGitignore, useCustomIgnore), 'Failed to load file tree.', { logContext: 'files' }),

    clearFileTreeCache: async (): Promise<void> => {
        try {
            await wails.ClearFileTreeCache()
        } catch (error) {
            console.error('[API:files] Error clearing file tree cache:', error)
        }
    },

    readFileContent: (projectPath: string, filePath: string): Promise<string> =>
        apiCall(() => wails.ReadFileContent(projectPath, filePath), 'Failed to read file content.', { logContext: 'files' }),

    getFileStats: (path: string): Promise<string> =>
        apiCall(() => wails.GetFileStats(path), 'Failed to get file statistics.', { logContext: 'files' }),
}
