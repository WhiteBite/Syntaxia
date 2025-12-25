/**
 * Application Configuration
 *
 * This file serves as a centralized configuration for component selection
 * throughout the application. This helps eliminate duplicated components
 * and clarify which version of a component should be used.
 */

export const APP_CONFIG = {
  // View Components - Define which workspace view is the primary one
  PRIMARY_WORKSPACE_VIEW: 'EnhancedWorkspaceView', // 'WorkspaceView' or 'EnhancedWorkspaceView'

  // UI Components - Define which UI components are the primary ones
  UI_COMPONENTS: {
    HEADER_BAR: 'HeaderBar', // 'HeaderBar' or 'CommandBar'
    BOTTOM_CONSOLE: 'BottomConsole',
    CONTEXT_VIEWER: 'ContextViewer',
    FILE_TREE: 'FilePanelModern' // 'FileTree' or 'FilePanelModern'
  },

  // Feature Flags - Enable/disable features
  FEATURES: {
    ENABLE_CONSOLE: true,
    ENABLE_SPLIT: true,
    ENABLE_VIRTUAL_MODE: true,
    ENABLE_RESIZE_HANDLES: true,
    ENABLE_MEMORY_MONITORING: true
  },

  // Layout Settings
  LAYOUT: {
    DEFAULT_CONTEXT_PANEL_WIDTH: 400,
    DEFAULT_RESULTS_PANEL_WIDTH: 500,
    MIN_PANEL_WIDTH: 200,
    MAX_PANEL_WIDTH_PERCENTAGE: 0.4, // 40% of window width
    MIN_CENTER_AREA_WIDTH: 300
  }
};

// Helper functions

/**
 * Get primary workspace view component name
 * @returns The primary workspace view component name
 */
export function getPrimaryWorkspaceView(): string {
  return APP_CONFIG.PRIMARY_WORKSPACE_VIEW;
}

/**
 * Get primary UI component
 * @param component The component key to get
 * @returns The primary component name
 */
export function getPrimaryComponent(component: keyof typeof APP_CONFIG.UI_COMPONENTS): string {
  return APP_CONFIG.UI_COMPONENTS[component];
}

/**
 * Check if a feature is enabled
 * @param feature The feature key to check
 * @returns True if the feature is enabled
 */
export function isFeatureEnabled(feature: keyof typeof APP_CONFIG.FEATURES): boolean {
  return APP_CONFIG.FEATURES[feature];
}

/**
 * Get layout setting
 * @param setting The layout setting key
 * @returns The layout setting value
 */
export function getLayoutSetting(setting: keyof typeof APP_CONFIG.LAYOUT): number {
  return APP_CONFIG.LAYOUT[setting];
}