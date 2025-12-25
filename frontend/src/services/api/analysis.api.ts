/**
 * Project analysis API
 * Handles static analysis, language detection
 */

import * as wails from '#wailsjs/go/main/App'
import type { domain } from '#wailsjs/go/models'
import { apiCall } from './base'

export const analysisApi = {
    analyzeProject: (path: string, analyzers: string[]): Promise<domain.StaticAnalysisReport> =>
        apiCall(() => wails.AnalyzeProject(path, analyzers), 'Failed to analyze project.', { logContext: 'analysis' }),

    analyzeFile: (projectPath: string, filePath: string): Promise<domain.StaticAnalysisResult> =>
        apiCall(
            () => wails.AnalyzeFile(projectPath, filePath),
            'Failed to analyze file.',
            { logContext: 'analysis' }
        ),

    detectLanguages: (projectPath: string): Promise<string[]> =>
        apiCall(
            () => wails.DetectLanguages(projectPath),
            'Failed to detect project languages.',
            { logContext: 'analysis' }
        ),

    getSupportedAnalyzers: (): Promise<string[]> =>
        apiCall(
            () => wails.GetSupportedAnalyzers(),
            'Failed to get supported analyzers.',
            { logContext: 'analysis' }
        ),
}
