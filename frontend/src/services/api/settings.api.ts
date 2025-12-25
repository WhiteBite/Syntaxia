/**
 * Settings and ignore rules API
 * Handles application settings and gitignore management
 */

import * as wails from '#wailsjs/go/main/App'
import type { domain } from '#wailsjs/go/models'
import { apiCall } from './base'

export interface DirCount {
    dir: string
    count: number
}

export interface IgnorePreviewResult {
    totalFiles: number
    byDirectory: Record<string, number>
    byRule: Record<string, number>
    sampleFiles: string[]
    topDirs: DirCount[]
}

export const settingsApi = {
    getSettings: (): Promise<domain.SettingsDTO> =>
        apiCall(() => wails.GetSettings(), 'Failed to load settings.', { logContext: 'settings' }),

    saveSettings: (settings: string): Promise<void> =>
        apiCall(() => wails.SaveSettings(settings), 'Failed to save settings.', { logContext: 'settings' }),

    // Ignore Rules
    getGitignoreContent: (projectPath: string): Promise<string> =>
        apiCall(
            () => wails.GetGitignoreContentForProject(projectPath),
            'Failed to load .gitignore content.',
            { logContext: 'settings' }
        ),

    getCustomIgnoreRules: (): Promise<string> =>
        apiCall(() => wails.GetCustomIgnoreRules(), 'Failed to load custom ignore rules.', { logContext: 'settings' }),

    updateCustomIgnoreRules: (rules: string): Promise<void> =>
        apiCall(
            () => wails.UpdateCustomIgnoreRules(rules),
            'Failed to update custom ignore rules.',
            { logContext: 'settings' }
        ),

    testIgnoreRules: (projectPath: string, rules: string): Promise<string[]> =>
        apiCall(
            () => wails.TestIgnoreRules(projectPath, rules),
            'Failed to test ignore rules.',
            { logContext: 'settings' }
        ),

    testIgnoreRulesDetailed: (projectPath: string, rules: string): Promise<IgnorePreviewResult> =>
        apiCall(
            () => wails.TestIgnoreRulesDetailed(projectPath, rules),
            'Failed to test ignore rules.',
            { logContext: 'settings' }
        ),

    addToGitignore: (projectPath: string, pattern: string): Promise<void> =>
        apiCall(
            () => wails.AddToGitignore(projectPath, pattern),
            'Failed to add to .gitignore.',
            { logContext: 'settings' }
        ),
}
