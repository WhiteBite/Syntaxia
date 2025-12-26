/**
 * Project management API
 * Handles projects, directories, paths
 */

import * as wails from '#wailsjs/go/main/App'
import { apiCall, apiCallWithDefault } from './base'

export const projectApi = {
    getRecentProjects: () =>
        apiCall(() => wails.GetRecentProjects(), 'Failed to load recent projects.', { logContext: 'project' }),

    addRecentProject: (path: string, name: string) =>
        apiCall(() => wails.AddRecentProject(path, name), 'Failed to add recent project.', { logContext: 'project' }),

    removeRecentProject: (path: string) =>
        apiCall(() => wails.RemoveRecentProject(path), 'Failed to remove recent project.', { logContext: 'project' }),

    selectDirectory: () =>
        apiCall(() => wails.SelectDirectory(), 'Failed to select directory.', { logContext: 'project' }),

    getCurrentDirectory: () =>
        apiCall(() => wails.GetCurrentDirectory(), 'Failed to get current directory.', { logContext: 'project' }),

    pathExists: (path: string) =>
        apiCallWithDefault(() => wails.PathExists(path), false, 'project'),
}
