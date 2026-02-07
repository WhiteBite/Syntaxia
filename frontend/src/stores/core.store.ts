/**
 * Core Store - Consolidated UI, Settings, and Project stores
 * Uses composition pattern with separate composables for each domain
 */

import { useLogger } from '@/composables/useLogger'
import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'

// ============================================================================
// Types
// ============================================================================

// UI Types
export interface ToastAction {
    label: string
    icon?: string
    onClick: () => void
}

export interface Toast {
    id: string
    message: string
    type: 'success' | 'error' | 'info' | 'warning'
    duration?: number
    action?: ToastAction
}

// Settings Types
export type OutputFormat = 'markdown' | 'xml' | 'plain'

export interface ContextSettings {
    maxTokens: number
    enforceTokenLimit: boolean
    stripComments: boolean
    includeTests: boolean
    splitStrategy: 'smart' | 'file' | 'token'
    outputFormat: OutputFormat
    excludeTests: boolean
    collapseEmptyLines: boolean
    stripLicense: boolean
    compactDataFiles: boolean
    trimWhitespace: boolean
    includeFileTree: boolean
    includeLineNumbers: boolean
    enableAutoSplit: boolean
    maxTokensPerChunk: number
    applyTemplateOnCopy: boolean
    enableNoiseReduction: boolean
    noiseReductionCollapseImports: boolean
    noiseReductionRemoveComments: boolean
    noiseReductionRemoveTypeDefinitions: boolean
}

export interface QuickFilterConfig {
    id: string
    label: string
    extensions: string[]
    patterns: string[]
    enabled: boolean
    isDynamic?: boolean
    language?: string
}

export interface FileExplorerSettings {
    useGitignore: boolean
    useCustomIgnore: boolean
    autoSaveSelection: boolean
    compactNestedFolders: boolean
    showIgnoredFiles: boolean
    foldersFirst: boolean
    allowSelectBinary: boolean
    customIgnoreRules: string
    quickFilters: QuickFilterConfig[]
    warnHeavyFiles: boolean
    heavyFileThreshold: number
    showDependencyIndicators: boolean
    autoHighlightDependencies: boolean
    hideNodeModules: boolean
    hideHiddenFiles: boolean
    hideTestFiles: boolean
    showTokenBars: boolean
}

export interface ContextStorageSettings {
    maxContexts: number
    maxStorageMB: number
    autoCleanupDays: number
    autoCleanupOnLimit: boolean
}

export interface AppSettings {
    context: ContextSettings
    contextStorage: ContextStorageSettings
    fileExplorer: FileExplorerSettings
    aiModel: string
    theme: 'dark' | 'light' | 'auto'
    uiScale: number
}

// Project Types
export interface RecentProject {
    path: string
    name: string
    lastOpened: number
}

// ============================================================================
// Constants
// ============================================================================

const DEFAULT_SETTINGS: AppSettings = {
    context: {
        maxTokens: 100000,
        enforceTokenLimit: true,
        stripComments: false,
        includeTests: true,
        splitStrategy: 'smart',
        outputFormat: 'xml',
        excludeTests: false,
        collapseEmptyLines: false,
        stripLicense: false,
        compactDataFiles: false,
        trimWhitespace: false,
        includeFileTree: false,
        includeLineNumbers: false,
        enableAutoSplit: false,
        maxTokensPerChunk: 32000,
        applyTemplateOnCopy: true,
        enableNoiseReduction: false,
        noiseReductionCollapseImports: true,
        noiseReductionRemoveComments: false,
        noiseReductionRemoveTypeDefinitions: false
    },
    contextStorage: {
        maxContexts: 20,
        maxStorageMB: 100,
        autoCleanupDays: 30,
        autoCleanupOnLimit: true
    },
    fileExplorer: {
        useGitignore: true,
        useCustomIgnore: false,
        autoSaveSelection: true,
        compactNestedFolders: true,
        showIgnoredFiles: true,
        foldersFirst: true,
        allowSelectBinary: false,
        customIgnoreRules: '',
        quickFilters: [
            { id: 'source', label: 'Исходники', extensions: ['.ts', '.js', '.tsx', '.jsx', '.vue', '.go', '.py', '.java', '.cpp', '.c', '.rs'], patterns: [], enabled: true },
            { id: 'tests', label: 'Тесты', extensions: [], patterns: ['**/*.test.*', '**/*.spec.*', '**/test/**', '**/tests/**', '**/Test/**'], enabled: true },
            { id: 'config', label: 'Конфигурация', extensions: ['.json', '.yaml', '.yml', '.toml', '.ini', '.env'], patterns: [], enabled: true },
            { id: 'docs', label: 'Документация', extensions: ['.md', '.txt', '.rst', '.adoc'], patterns: [], enabled: true },
            { id: 'styles', label: 'Стили', extensions: ['.css', '.scss', '.sass', '.less'], patterns: [], enabled: true }
        ],
        warnHeavyFiles: true,
        heavyFileThreshold: 100000,
        showDependencyIndicators: true,
        autoHighlightDependencies: true,
        hideNodeModules: true,
        hideHiddenFiles: true,
        hideTestFiles: false,
        showTokenBars: true
    },
    aiModel: 'gpt-4',
    theme: 'dark',
    uiScale: 1.0
}

const RECENT_PROJECTS_KEY = 'Syntaxia_recent_projects'
const MAX_RECENT = 10

// ============================================================================
// Composables
// ============================================================================

const uiLogger = useLogger('CoreStore:UI')
const settingsLogger = useLogger('CoreStore:Settings')
const projectLogger = useLogger('CoreStore:Project')

/**
 * UI State Composable
 */
function useUIState() {
    const toasts = ref<Toast[]>([])
    const nextToastId = ref(0)
    const showSettingsModal = ref(false)
    const showKeyboardShortcutsModal = ref(false)
    const isFileSearchModalOpen = ref(false)

    function addToast(message: string, type: Toast['type'] = 'info', duration = 3000, action?: ToastAction) {
        const toast: Toast = {
            id: `toast-${nextToastId.value++}`,
            message,
            type,
            duration,
            action
        }

        const logMessage = `[Toast ${type.toUpperCase()}] ${message}`
        switch (type) {
            case 'error':
                uiLogger.error(logMessage)
                break
            case 'warning':
                uiLogger.warn(logMessage)
                break
            case 'info':
                uiLogger.info(logMessage)
                break
            case 'success':
                uiLogger.debug(logMessage)
                break
        }

        toasts.value.push(toast)

        if (duration > 0) {
            setTimeout(() => {
                removeToast(toast.id)
            }, duration)
        }

        return toast.id
    }

    function removeToast(id: string) {
        const index = toasts.value.findIndex(t => t.id === id)
        if (index !== -1) {
            toasts.value.splice(index, 1)
        }
    }

    function clearToasts() {
        toasts.value = []
    }

    function openSettingsModal() {
        showSettingsModal.value = true
    }

    function openKeyboardShortcutsModal() {
        showKeyboardShortcutsModal.value = true
    }

    function openFileSearchModal() {
        isFileSearchModalOpen.value = true
    }

    function closeFileSearchModal() {
        isFileSearchModalOpen.value = false
    }

    return {
        toasts,
        showSettingsModal,
        showKeyboardShortcutsModal,
        isFileSearchModalOpen,
        addToast,
        removeToast,
        clearToasts,
        openSettingsModal,
        openKeyboardShortcutsModal,
        openFileSearchModal,
        closeFileSearchModal
    }
}

/**
 * Settings State Composable
 */
function useSettingsState() {
    function loadSettings(): AppSettings {
        try {
            const saved = localStorage.getItem('app-settings')
            if (saved) {
                if (saved === 'undefined' || saved === 'null' || saved.trim() === '') {
                    localStorage.removeItem('app-settings')
                    return DEFAULT_SETTINGS
                }

                const parsed = JSON.parse(saved)

                if (typeof parsed !== 'object' || parsed === null) {
                    localStorage.removeItem('app-settings')
                    return DEFAULT_SETTINGS
                }

                return {
                    ...DEFAULT_SETTINGS,
                    ...parsed,
                    context: {
                        ...DEFAULT_SETTINGS.context,
                        ...(parsed.context || {})
                    },
                    contextStorage: {
                        ...DEFAULT_SETTINGS.contextStorage,
                        ...(parsed.contextStorage || {})
                    },
                    fileExplorer: {
                        ...DEFAULT_SETTINGS.fileExplorer,
                        ...(parsed.fileExplorer || {})
                    }
                }
            }
        } catch (err) {
            try {
                localStorage.removeItem('app-settings')
            } catch {
                // Ignore localStorage errors
            }
        }
        return DEFAULT_SETTINGS
    }

    function saveSettings(settings: AppSettings) {
        try {
            localStorage.setItem('app-settings', JSON.stringify(settings))
        } catch (err) {
            settingsLogger.warn('Failed to save settings:', err)
        }
    }

    const settings = ref<AppSettings>(loadSettings())

    watch(
        settings,
        (newSettings) => {
            saveSettings(newSettings)
        },
        { deep: true }
    )

    function resetToDefaults() {
        settings.value = JSON.parse(JSON.stringify(DEFAULT_SETTINGS))
    }

    function updateContextSettings(updates: Partial<ContextSettings>) {
        settings.value.context = {
            ...settings.value.context,
            ...updates
        }
    }

    function updateAIModel(model: string) {
        settings.value.aiModel = model
    }

    function updateTheme(theme: 'dark' | 'light' | 'auto') {
        settings.value.theme = theme
    }

    function updateUIScale(scale: number) {
        settings.value.uiScale = scale
    }

    function updateFileExplorerSettings(updates: Partial<FileExplorerSettings>) {
        settings.value.fileExplorer = {
            ...settings.value.fileExplorer,
            ...updates
        }
    }

    function getCustomIgnoreRules(): string {
        return settings.value.fileExplorer.customIgnoreRules
    }

    function setCustomIgnoreRules(rules: string) {
        settings.value.fileExplorer.customIgnoreRules = rules
    }

    function updateContextStorageSettings(updates: Partial<ContextStorageSettings>) {
        settings.value.contextStorage = {
            ...settings.value.contextStorage,
            ...updates
        }
    }

    return {
        settings,
        resetToDefaults,
        updateContextSettings,
        updateContextStorageSettings,
        updateAIModel,
        updateTheme,
        updateUIScale,
        updateFileExplorerSettings,
        getCustomIgnoreRules,
        setCustomIgnoreRules
    }
}

/**
 * Project State Composable
 */
function useProjectState() {
    const currentPath = ref<string | null>(null)
    const currentName = ref<string | null>(null)
    const recentProjects = ref<RecentProject[]>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)
    const autoOpenLast = ref<boolean>(false)

    const hasProject = computed(() => currentPath.value !== null)
    const projectName = computed(() => currentName.value || '')
    const projectPath = computed(() => currentPath.value || '')
    const hasRecentProjects = computed(() => recentProjects.value.length > 0)
    const lastProjectPath = computed(() => recentProjects.value[0]?.path || null)

    async function openProjectByPath(path: string): Promise<boolean> {
        // Import dynamically to avoid circular dependencies
        const { apiService } = await import('@/services/api.service')
        const { useFileStore } = await import('@/features/files')
        const { useContextStore } = await import('@/features/context')

        isLoading.value = true
        error.value = null

        try {
            const exists = await apiService.pathExists(path)
            if (!exists) {
                error.value = `Path does not exist: ${path}`
                recentProjects.value = recentProjects.value.filter(p => p.path !== path)
                saveRecentProjects()
                return false
            }

            currentPath.value = path
            currentName.value = path.split(/[\\/]/).pop() || path

            const fileStore = useFileStore()
            const contextStore = useContextStore()
            fileStore.resetStore()
            contextStore.clearContext()

            addToRecent(path)

            try {
                const name = path.split(/[\\/]/).pop() || path
                await apiService.addRecentProject(path, name)
            } catch (backendError) {
                projectLogger.error('Failed to save project to backend:', backendError)
            }

            saveRecentProjects()

            return true
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Failed to load project'
            return false
        } finally {
            isLoading.value = false
        }
    }

    function addToRecent(path: string) {
        const name = path.split(/[\\/]/).pop() || path
        const existing = recentProjects.value.findIndex(p => p.path === path)

        if (existing !== -1) {
            recentProjects.value.splice(existing, 1)
        }

        recentProjects.value.unshift({
            path,
            name,
            lastOpened: Date.now()
        })

        if (recentProjects.value.length > MAX_RECENT) {
            recentProjects.value = recentProjects.value.slice(0, MAX_RECENT)
        }
    }

    function loadRecentProjects() {
        try {
            const stored = localStorage.getItem(RECENT_PROJECTS_KEY)
            if (stored) {
                recentProjects.value = JSON.parse(stored)
            }
            const auto = localStorage.getItem('Syntaxia_auto_open_last')
            if (auto !== null) {
                autoOpenLast.value = auto === 'true'
            }
        } catch (err) {
            projectLogger.warn('Failed to load recent projects:', err)
        }
    }

    async function fetchRecentProjects() {
        const { apiService } = await import('@/services/api.service')

        try {
            isLoading.value = true
            error.value = null

            const projectsJson = await apiService.getRecentProjects()

            if (projectsJson) {
                const projects = JSON.parse(projectsJson)
                if (Array.isArray(projects)) {
                    recentProjects.value = projects
                    saveRecentProjects()
                }
            }
        } catch (err) {
            projectLogger.error('Failed to fetch recent projects from backend:', err)
            error.value = err instanceof Error ? err.message : 'Failed to load recent projects'
        } finally {
            isLoading.value = false
        }
    }

    function saveRecentProjects() {
        try {
            localStorage.setItem(RECENT_PROJECTS_KEY, JSON.stringify(recentProjects.value))
        } catch (err) {
            projectLogger.warn('Failed to save recent projects:', err)
        }
    }

    function setAutoOpenLast(value: boolean) {
        autoOpenLast.value = value
        try {
            localStorage.setItem('Syntaxia_auto_open_last', String(value))
        } catch (err) {
            projectLogger.warn('Failed to save auto-open setting:', err)
        }
    }

    async function maybeAutoOpenLastProject(): Promise<boolean> {
        if (currentPath.value) return true
        if (!autoOpenLast.value) return false
        if (!lastProjectPath.value) return false
        return await openProjectByPath(lastProjectPath.value)
    }

    async function removeFromRecent(path: string) {
        const { apiService } = await import('@/services/api.service')

        recentProjects.value = recentProjects.value.filter(p => p.path !== path)

        try {
            await apiService.removeRecentProject(path)
        } catch (backendError) {
            projectLogger.error('Failed to remove project from backend:', backendError)
        }

        saveRecentProjects()
    }

    function clearRecent() {
        recentProjects.value = []
        saveRecentProjects()
    }

    function clearProject() {
        currentPath.value = null
        currentName.value = null
        error.value = null

        // Note: Store cleanup is handled by the stores themselves
        // Tests should mock the stores or handle cleanup separately
    }

    loadRecentProjects()

    return {
        currentPath,
        currentName,
        recentProjects,
        isLoading,
        error,
        autoOpenLast,
        hasProject,
        projectName,
        projectPath,
        hasRecentProjects,
        lastProjectPath,
        openProjectByPath,
        fetchRecentProjects,
        removeFromRecent,
        clearRecent,
        clearProject,
        setAutoOpenLast,
        maybeAutoOpenLastProject
    }
}

// ============================================================================
// Store Definition
// ============================================================================

export const useCoreStore = defineStore('core', () => {
    const ui = useUIState()
    const settings = useSettingsState()
    const project = useProjectState()

    return {
        ui,
        settings,
        project
    }
})

// ============================================================================
// Backward Compatibility Aliases
// ============================================================================

export const useUIStore = () => useCoreStore().ui
export const useSettingsStore = () => useCoreStore().settings
export const useProjectStore = () => useCoreStore().project
