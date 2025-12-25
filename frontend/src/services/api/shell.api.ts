/**
 * Shell Integration API
 * Handles OS context menu integration
 */

import * as wails from '#wailsjs/go/main/App'
import { apiCall } from './base'

export interface ShellIntegrationStatus {
    isRegistered: boolean
    currentOS: string
}

export const shellApi = {
    getStatus: (): Promise<ShellIntegrationStatus> =>
        apiCall(
            () => wails.GetShellIntegrationStatus(),
            'Failed to get shell integration status.',
            { logContext: 'shell' }
        ),

    register: (): Promise<void> =>
        apiCall(
            () => wails.RegisterShellIntegration(),
            'Failed to register shell integration.',
            { logContext: 'shell' }
        ),

    unregister: (): Promise<void> =>
        apiCall(
            () => wails.UnregisterShellIntegration(),
            'Failed to unregister shell integration.',
            { logContext: 'shell' }
        ),

    getStartupPath: (): Promise<string> =>
        apiCall(
            () => wails.GetStartupPath(),
            'Failed to get startup path.',
            { logContext: 'shell' }
        ),

    clearStartupPath: (): Promise<void> =>
        apiCall(
            () => wails.ClearStartupPath(),
            'Failed to clear startup path.',
            { logContext: 'shell' }
        ),
}
