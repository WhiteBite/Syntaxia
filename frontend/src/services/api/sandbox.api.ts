/**
 * Sandbox API for AI file changes
 * Uses Wails generated bindings
 */

import * as wails from '#wailsjs/go/main/App'
import type { SandboxChangesResponse } from '../types'

export const sandboxApi = {
    getChanges: async (): Promise<SandboxChangesResponse> => {
        const result = await wails.GetSandboxChanges()
        return JSON.parse(result) as SandboxChangesResponse
    },

    getDiff: async (path: string): Promise<string> => {
        return await wails.GetSandboxDiff(path)
    },

    getAllDiffs: async (): Promise<string> => {
        return await wails.GetSandboxAllDiffs()
    },

    applyChanges: async (): Promise<void> => {
        await wails.ApplySandboxChanges()
    },

    discardChanges: async (): Promise<void> => {
        await wails.DiscardSandboxChanges()
    },

    discardFile: async (path: string): Promise<void> => {
        await wails.DiscardSandboxFile(path)
    },

    hasChanges: async (): Promise<boolean> => {
        return await wails.HasSandboxChanges()
    },

    getChangeCount: async (): Promise<number> => {
        return await wails.GetSandboxChangeCount()
    }
}
