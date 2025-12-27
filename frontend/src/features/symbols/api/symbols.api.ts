/**
 * Symbols API - Wails backend integration
 */

import * as wails from '#wailsjs/go/main/App'
import { useLogger } from '@/composables/useLogger'
import { apiCall, apiCallSafe } from '@/services/api/base'
import type { Symbol, SymbolKind } from '../types'

const logger = useLogger('SymbolsApi')

interface BackendSymbolInfo {
    name: string
    kind: string
    filePath: string
    startLine: number
    endLine: number
    signature?: string
    parent?: string
    modifiers?: string[]
}

/**
 * Map backend symbol kind to frontend SymbolKind
 */
function mapSymbolKind(kind: string): SymbolKind {
    const kindMap: Record<string, SymbolKind> = {
        'class': 'class',
        'function': 'function',
        'method': 'method',
        'interface': 'interface',
        'variable': 'variable',
        'constant': 'constant',
        'type': 'type',
        'enum': 'enum',
        'property': 'property',
        'struct': 'class',
        'trait': 'interface',
    }
    return kindMap[kind.toLowerCase()] || 'function'
}

/**
 * Convert backend symbol to frontend Symbol type
 */
function mapSymbol(info: BackendSymbolInfo): Symbol {
    return {
        name: info.name,
        kind: mapSymbolKind(info.kind),
        file: info.filePath,
        line: info.startLine,
        endLine: info.endLine,
        signature: info.signature,
        parent: info.parent,
        modifiers: info.modifiers,
    }
}

export class SymbolsApi {
    /**
     * Get all symbols for a project
     */
    async getProjectSymbols(projectPath: string): Promise<Symbol[]> {
        logger.info(`Loading symbols for project: ${projectPath}`)

        // Search with empty query to get all symbols
        const result = await apiCallSafe(
            () => wails.SearchSymbols(projectPath, '*', ''),
            [] as BackendSymbolInfo[],
            'Failed to load project symbols'
        )

        return result.map(mapSymbol)
    }

    /**
     * Search symbols by query
     */
    async searchSymbols(
        query: string,
        projectPath: string,
        kindFilter?: SymbolKind
    ): Promise<Symbol[]> {
        if (!query || query.length < 2) {
            return []
        }

        logger.info(`Searching symbols: "${query}" in ${projectPath}`)

        const result = await apiCall(
            () => wails.SearchSymbols(projectPath, query, kindFilter || ''),
            'Failed to search symbols',
            { logContext: 'symbols' }
        )

        return (result || []).map(mapSymbol)
    }

    /**
     * Get symbols for a specific file
     */
    async getFileSymbols(projectPath: string, filePath: string): Promise<Symbol[]> {
        logger.info(`Loading symbols for file: ${filePath}`)

        const result = await apiCallSafe(
            () => wails.ListSymbols(projectPath, filePath),
            [] as BackendSymbolInfo[],
            'Failed to load file symbols'
        )

        return result.map(mapSymbol)
    }
}

export const symbolsApi = new SymbolsApi()
