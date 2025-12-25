/**
 * Reports and export API
 * Handles report generation and project export
 */

import * as wails from '#wailsjs/go/main/App'
import { apiCall } from './base'

export const reportsApi = {
    generateReport: (contextId: string, format: string): Promise<string> =>
        apiCall(
            () => wails.GenerateReport(contextId, format),
            'Failed to generate report.',
            { logContext: 'reports' }
        ),

    listReports: (projectPath: string): Promise<string> =>
        apiCall(() => wails.ListReports(projectPath), 'Failed to list reports.', { logContext: 'reports' }),

    getReport: (reportId: string): Promise<string> =>
        apiCall(() => wails.GetReport(reportId), 'Failed to get report.', { logContext: 'reports' }),

    exportProject: (projectPath: string, format: string, outputPath: string): Promise<string> =>
        apiCall(
            () => wails.ExportProject(projectPath, format, outputPath),
            'Failed to export project.',
            { logContext: 'reports' }
        ),
}
