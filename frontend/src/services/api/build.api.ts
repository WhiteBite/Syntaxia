/**
 * Build, test, and diff API
 * Handles project building, testing, type checking, and diff operations
 */

import * as wails from '#wailsjs/go/main/App'
import type { domain } from '#wailsjs/go/models'
import { apiCall } from './base'

export const buildApi = {
    // Testing
    runTests: (config: domain.TestConfig): Promise<domain.TestResult[]> =>
        apiCall(() => wails.RunTests(config), 'Failed to run tests.', { logContext: 'build' }),

    discoverTests: (projectPath: string, language: string): Promise<domain.TestSuite> =>
        apiCall(
            () => wails.DiscoverTests(projectPath, language),
            'Failed to discover tests.',
            { logContext: 'build' }
        ),

    // Build
    build: (projectPath: string, language: string): Promise<domain.BuildResult> =>
        apiCall(() => wails.Build(projectPath, language), 'Failed to build project.', { logContext: 'build' }),

    typeCheck: (projectPath: string, language: string): Promise<domain.TypeCheckResult> =>
        apiCall(() => wails.TypeCheck(projectPath, language), 'Failed to type check.', { logContext: 'build' }),

    // Diff and Apply
    generateDiff: (original: string, modified: string, format: string): Promise<domain.DiffResult> =>
        apiCall(
            () => wails.GenerateDiff(original, modified, format),
            'Failed to generate diff.',
            { logContext: 'build' }
        ),

    applyEdits: (edits: domain.EditsJSON): Promise<domain.ApplyResult[]> =>
        apiCall(() => wails.ApplyEdits(edits), 'Failed to apply edits.', { logContext: 'build' }),

    applySingleEdit: (edit: domain.Edit): Promise<domain.ApplyResult> =>
        apiCall(() => wails.ApplySingleEdit(edit), 'Failed to apply edit.', { logContext: 'build' }),
}
