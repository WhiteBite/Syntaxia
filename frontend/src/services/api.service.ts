/**
 * API Service Facade
 * Re-exports all API modules for backward compatibility
 *
 * This file serves as a facade that delegates to modular API files.
 * For new code, prefer importing directly from '@/services/api':
 *   import { projectApi, filesApi } from '@/services/api'
 *
 * Legacy usage (still supported):
 *   import { apiService } from '@/services/api.service'
 */

import { aiApi } from './api/ai.api'
import { analysisApi } from './api/analysis.api'
import { buildApi } from './api/build.api'
import { contextApi } from './api/context.api'
import { filesApi } from './api/files.api'
import { gitApi } from './api/git.api'
import { githubApi } from './api/github.api'
import { gitlabApi } from './api/gitlab.api'
import { memoryApi } from './api/memory.api'
import { projectApi } from './api/project.api'
import { reportsApi } from './api/reports.api'
import { semanticApi } from './api/semantic.api'
import { settingsApi } from './api/settings.api'
import { taskflowApi } from './api/taskflow.api'

/**
 * Unified API service object
 * Provides backward compatibility with the old monolithic ApiService class
 */
export const apiService = {
  // ============================================
  // Project Management
  // ============================================
  getRecentProjects: projectApi.getRecentProjects,
  addRecentProject: projectApi.addRecentProject,
  removeRecentProject: projectApi.removeRecentProject,
  selectDirectory: projectApi.selectDirectory,
  getCurrentDirectory: projectApi.getCurrentDirectory,
  pathExists: projectApi.pathExists,

  // ============================================
  // File Operations
  // ============================================
  listFiles: filesApi.listFiles,
  clearFileTreeCache: filesApi.clearFileTreeCache,
  readFileContent: filesApi.readFileContent,
  getFileStats: filesApi.getFileStats,

  // ============================================
  // Context
  // ============================================
  buildContext: contextApi.buildContext,
  buildContextFromRequest: contextApi.buildContextFromRequest,
  getContext: contextApi.getContext,
  deleteContext: contextApi.deleteContext,
  getProjectContexts: contextApi.getProjectContexts,
  exportContext: contextApi.exportContext,
  getFullContextContent: contextApi.getFullContextContent,
  suggestContextFiles: contextApi.suggestContextFiles,
  getSmartSuggestions: contextApi.getSmartSuggestions,
  getFileQuickInfo: contextApi.getFileQuickInfo,
  getImpactPreview: contextApi.getImpactPreview,
  analyzeTaskAndCollectContext: contextApi.analyzeTaskAndCollectContext,
  agenticChat: contextApi.agenticChat,

  // ============================================
  // AI and Code Generation
  // ============================================
  generateCode: aiApi.generateCode,
  generateCodeStream: aiApi.generateCodeStream,
  generateIntelligentCode: aiApi.generateIntelligentCode,
  listAvailableModels: aiApi.listAvailableModels,
  getProviderInfo: aiApi.getProviderInfo,
  qwenExecuteTask: aiApi.qwenExecuteTask,
  qwenPreviewContext: aiApi.qwenPreviewContext,
  qwenGetAvailableModels: aiApi.qwenGetAvailableModels,

  // ============================================
  // Analysis
  // ============================================
  analyzeProject: analysisApi.analyzeProject,
  analyzeFile: analysisApi.analyzeFile,
  detectLanguages: analysisApi.detectLanguages,
  getSupportedAnalyzers: analysisApi.getSupportedAnalyzers,

  // ============================================
  // Git Operations
  // ============================================
  getUncommittedFiles: gitApi.getUncommittedFiles,
  getBranches: gitApi.getBranches,
  getCurrentBranch: gitApi.getCurrentBranch,
  getRichCommitHistory: gitApi.getRichCommitHistory,
  isGitAvailable: gitApi.isGitAvailable,
  isGitRepository: gitApi.isGitRepository,
  cloneRepository: gitApi.cloneRepository,
  checkoutBranch: gitApi.checkoutBranch,
  checkoutCommit: gitApi.checkoutCommit,
  getCommitHistory: gitApi.getCommitHistory,
  getRemoteBranches: gitApi.getRemoteBranches,
  cleanupTempRepository: gitApi.cleanupTempRepository,
  listFilesAtRef: gitApi.listFilesAtRef,
  getFileAtRef: gitApi.getFileAtRef,
  buildContextAtRef: gitApi.buildContextAtRef,

  // ============================================
  // GitHub API
  // ============================================
  isGitHubURL: githubApi.isGitHubURL,
  gitHubGetDefaultBranch: githubApi.getDefaultBranch,
  gitHubGetBranches: githubApi.getBranches,
  gitHubGetCommits: githubApi.getCommits,
  gitHubListFiles: githubApi.listFiles,
  gitHubGetFileContent: githubApi.getFileContent,
  gitHubBuildContext: githubApi.buildContext,

  // ============================================
  // GitLab API
  // ============================================
  isGitLabURL: gitlabApi.isGitLabURL,
  gitLabGetDefaultBranch: gitlabApi.getDefaultBranch,
  gitLabGetBranches: gitlabApi.getBranches,
  gitLabGetCommits: gitlabApi.getCommits,
  gitLabListFiles: gitlabApi.listFiles,
  gitLabGetFileContent: gitlabApi.getFileContent,
  gitLabBuildContext: gitlabApi.buildContext,

  // ============================================
  // Settings
  // ============================================
  getSettings: settingsApi.getSettings,
  saveSettings: settingsApi.saveSettings,
  getGitignoreContent: settingsApi.getGitignoreContent,
  getCustomIgnoreRules: settingsApi.getCustomIgnoreRules,
  updateCustomIgnoreRules: settingsApi.updateCustomIgnoreRules,
  testIgnoreRules: settingsApi.testIgnoreRules,
  testIgnoreRulesDetailed: settingsApi.testIgnoreRulesDetailed,
  addToGitignore: settingsApi.addToGitignore,

  // ============================================
  // Build and Test
  // ============================================
  runTests: buildApi.runTests,
  discoverTests: buildApi.discoverTests,
  build: buildApi.build,
  typeCheck: buildApi.typeCheck,
  generateDiff: buildApi.generateDiff,
  applyEdits: buildApi.applyEdits,
  applySingleEdit: buildApi.applySingleEdit,

  // ============================================
  // Reports
  // ============================================
  generateReport: reportsApi.generateReport,
  listReports: reportsApi.listReports,
  getReport: reportsApi.getReport,
  exportProject: reportsApi.exportProject,

  // ============================================
  // Task Protocol and Guardrails
  // ============================================
  executeTaskProtocol: taskflowApi.executeTaskProtocol,
  getTaskProtocolConfiguration: taskflowApi.getTaskProtocolConfiguration,
  validatePath: taskflowApi.validatePath,
  getGuardrailPolicies: taskflowApi.getGuardrailPolicies,
  getBudgetPolicies: taskflowApi.getBudgetPolicies,

  // ============================================
  // Semantic Search
  // ============================================
  isSemanticSearchAvailable: semanticApi.isAvailable,
  semanticSearch: semanticApi.search,
  semanticFindSimilar: semanticApi.findSimilar,
  semanticIndexProject: semanticApi.indexProject,
  semanticIndexFile: semanticApi.indexFile,
  semanticGetStats: semanticApi.getStats,
  semanticIsIndexed: semanticApi.isIndexed,
  semanticRetrieveContext: semanticApi.retrieveContext,
  semanticHybridSearch: semanticApi.hybridSearch,

  // ============================================
  // Context Memory
  // ============================================
  getRecentContexts: memoryApi.getRecentContexts,
  findContextByTopic: memoryApi.findContextByTopic,
  saveContextMemory: memoryApi.saveContextMemory,
}

// Re-export all types for backward compatibility
export * from './types'

// Export modular APIs for new code
export {
  aiApi,
  analysisApi,
  buildApi,
  contextApi,
  filesApi,
  gitApi,
  githubApi,
  gitlabApi,
  memoryApi,
  projectApi,
  reportsApi,
  semanticApi,
  settingsApi,
  taskflowApi
}

