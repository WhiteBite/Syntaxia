/**
 * Unified API module
 * Re-exports all API modules for convenient imports
 *
 * Usage:
 *   import { projectApi, filesApi, contextApi } from '@/services/api'
 *   // or
 *   import { api } from '@/services/api'
 *   api.project.getRecentProjects()
 */

export { aiApi } from './ai.api'
export { analysisApi } from './analysis.api'
export { buildApi } from './build.api'
export { contextApi } from './context.api'
export { filesApi } from './files.api'
export { gitApi } from './git.api'
export { githubApi } from './github.api'
export { gitlabApi } from './gitlab.api'
export { memoryApi } from './memory.api'
export { projectApi } from './project.api'
export { reportsApi } from './reports.api'
export { semanticApi } from './semantic.api'
export { settingsApi } from './settings.api'
export { taskflowApi } from './taskflow.api'

// Unified API object for convenience
import { aiApi } from './ai.api'
import { analysisApi } from './analysis.api'
import { buildApi } from './build.api'
import { contextApi } from './context.api'
import { filesApi } from './files.api'
import { gitApi } from './git.api'
import { githubApi } from './github.api'
import { gitlabApi } from './gitlab.api'
import { memoryApi } from './memory.api'
import { projectApi } from './project.api'
import { reportsApi } from './reports.api'
import { semanticApi } from './semantic.api'
import { settingsApi } from './settings.api'
import { taskflowApi } from './taskflow.api'

export const api = {
    project: projectApi,
    files: filesApi,
    context: contextApi,
    ai: aiApi,
    git: gitApi,
    github: githubApi,
    gitlab: gitlabApi,
    analysis: analysisApi,
    settings: settingsApi,
    build: buildApi,
    reports: reportsApi,
    taskflow: taskflowApi,
    semantic: semanticApi,
    memory: memoryApi,
}

// Re-export types
export * from '../types'
