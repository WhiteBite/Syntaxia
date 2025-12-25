/**
 * Application constants - centralized magic numbers and configuration values
 */

// File Tree Configuration
export const FILE_TREE = {
    MAX_SEARCH_RESULTS: 100,
    MAX_FLATTEN_DEPTH: 5,
    DEBOUNCE_MS: 300,
    MAX_SAVED_SELECTIONS: 100,
    DEFAULT_EXPANDED_DEPTH: 2,
    VIRTUAL_ITEM_HEIGHT: 28,
    OVERSCAN_COUNT: 5,
} as const

// Cache Configuration
export const CACHE = {
    TTL_MS: 5 * 60 * 1000, // 5 minutes
    MAX_ENTRIES: 1000,
    STALE_TIME_MS: 30 * 1000, // 30 seconds
} as const

// UI Configuration
export const UI = {
    TOAST_DURATION_MS: 3000,
    TOAST_ERROR_DURATION_MS: 5000,
    ANIMATION_DURATION_MS: 200,
    COPY_FEEDBACK_DURATION_MS: 2000,
    MODAL_BACKDROP_OPACITY: 0.5,
} as const

// API Configuration
export const API = {
    RETRY_COUNT: 3,
    RETRY_DELAY_MS: 1000,
    TIMEOUT_MS: 30000,
    STREAM_CHUNK_SIZE: 1024,
} as const

// Context Configuration
export const CONTEXT = {
    MAX_FILES: 500,
    MAX_TOKENS: 100000,
    MAX_FILE_SIZE_BYTES: 1024 * 1024, // 1MB
    CHUNK_SIZE: 50000,
} as const

// Token Weight Thresholds (for file tree visualization)
export const TOKEN_THRESHOLDS = {
    MEDIUM: 5000,      // 5k tokens - yellow indicator
    HEAVY: 20000,      // 20k tokens - orange indicator  
    CRITICAL: 50000,   // 50k tokens - red indicator
    BYTES_PER_TOKEN: 4, // approximate bytes per token
} as const

// Git Configuration
export const GIT = {
    MAX_COMMITS: 100,
    MAX_DIFF_LINES: 1000,
    MAX_CHANGED_FILES: 50,
} as const

// Undo/Redo Configuration
export const UNDO_REDO = {
    MAX_HISTORY_SIZE: 50,
} as const

// Keyboard Shortcuts
export const SHORTCUTS = {
    UNDO: 'Ctrl+Z',
    REDO: 'Ctrl+Shift+Z',
    SEARCH: 'Ctrl+P',
    COMMAND_PALETTE: 'Ctrl+Shift+P',
    COPY: 'Ctrl+C',
    SELECT_ALL: 'Ctrl+A',
} as const

// LocalStorage Keys - centralized to avoid magic strings
export const STORAGE_KEYS = {
    // Project
    RECENT_PROJECTS: 'Syntaxia_recent_projects',
    AUTO_OPEN_LAST: 'Syntaxia_auto_open_last',
    // Settings
    APP_SETTINGS: 'app-settings',
    // Context
    CONTEXT_METADATA: 'context-metadata',
    // Chat
    CHAT_HISTORY_PREFIX: 'chat-history',
    // UI State
    LEFT_SIDEBAR_TAB: 'left-sidebar-tab',
    RIGHT_SIDEBAR_VISIBLE: 'right-sidebar-visible',
    RIGHT_SIDEBAR_TAB: 'right-sidebar-tab',
    TASK_PANEL_VISIBLE: 'task-panel-visible',
    // File Explorer
    EXPANDED_NODES_PREFIX: 'expanded-nodes',
    SELECTED_FILES_PREFIX: 'selected-files',
    FILTER_STATE_PREFIX: 'filter-state',
    // Workspace Layout
    WORKSPACE_LEFT_WIDTH: 'workspace-left-width',
    WORKSPACE_RIGHT_WIDTH: 'workspace-right-width',
    // Onboarding
    ONBOARDING_COMPLETED: 'Syntaxia-onboarding-completed',
} as const
