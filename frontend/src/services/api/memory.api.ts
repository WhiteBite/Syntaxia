/**
 * Context Memory API
 * Handles context persistence and retrieval
 */

import * as wails from '#wailsjs/go/main/App'
import type { ContextMemoryEntry } from '../types'
import { apiCall, apiCallWithDefault } from './base'

export const memoryApi = {
    getRecentContexts: (projectPath: string, limit = 10): Promise<ContextMemoryEntry[]> =>
        apiCallWithDefault(
            // @ts-ignore
            () => wails.GetRecentContexts(projectPath, limit),
            [],
            'memory.getRecentContexts'
        ),

    findContextByTopic: (projectPath: string, topic: string): Promise<ContextMemoryEntry[]> =>
        apiCallWithDefault(
            // @ts-ignore
            () => wails.FindContextByTopic(projectPath, topic),
            [],
            'memory.findContextByTopic'
        ),

    saveContextMemory: (
        projectPath: string,
        topic: string,
        summary: string,
        files: string[]
    ): Promise<void> =>
        apiCall(
            // @ts-ignore
            () => wails.SaveContextMemory(projectPath, topic, summary, files),
            'Failed to save context.',
            { logContext: 'memory' }
        ),
}
