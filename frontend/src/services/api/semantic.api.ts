/**
 * Semantic Search API
 * Handles semantic search, indexing, and RAG operations
 */

import * as wails from '#wailsjs/go/main/App'
import type {
    CodeChunk,
    FindSimilarRequest,
    RetrieveContextRequest,
    SemanticIndexStats,
    SemanticSearchRequest,
    SemanticSearchResponse,
} from '../types'
import { apiCall, apiCallWithDefault, parseJsonResponse } from './base'

export const semanticApi = {
    isAvailable: (): Promise<boolean> =>
        apiCallWithDefault(
            // @ts-ignore - method may not exist in wails bindings yet
            () => wails.IsSemanticSearchAvailable(),
            false,
            'semantic'
        ),

    search: async (request: SemanticSearchRequest): Promise<SemanticSearchResponse> => {
        const result = await apiCall(
            // @ts-ignore
            () => wails.SemanticSearch(JSON.stringify(request)),
            'Failed to perform semantic search.',
            { logContext: 'semantic' }
        )
        return parseJsonResponse(result, 'Failed to parse semantic search response.')
    },

    findSimilar: async (request: FindSimilarRequest): Promise<SemanticSearchResponse> => {
        const result = await apiCall(
            // @ts-ignore
            () => wails.SemanticFindSimilar(JSON.stringify(request)),
            'Failed to find similar code.',
            { logContext: 'semantic' }
        )
        return parseJsonResponse(result, 'Failed to parse similar code response.')
    },

    indexProject: (projectRoot: string): Promise<void> =>
        apiCall(
            // @ts-ignore
            () => wails.SemanticIndexProject(projectRoot),
            'Failed to index project.',
            { logContext: 'semantic' }
        ),

    indexFile: (projectRoot: string, filePath: string): Promise<void> =>
        apiCall(
            // @ts-ignore
            () => wails.SemanticIndexFile(projectRoot, filePath),
            'Failed to index file.',
            { logContext: 'semantic' }
        ),

    getStats: async (projectRoot: string): Promise<SemanticIndexStats> => {
        const result = await apiCall(
            // @ts-ignore
            () => wails.SemanticGetStats(projectRoot),
            'Failed to get semantic search statistics.',
            { logContext: 'semantic' }
        )
        return parseJsonResponse(result, 'Failed to parse semantic stats.')
    },

    isIndexed: (projectRoot: string): Promise<boolean> =>
        apiCallWithDefault(
            // @ts-ignore
            () => wails.SemanticIsIndexed(projectRoot),
            false,
            'semantic'
        ),

    retrieveContext: async (request: RetrieveContextRequest): Promise<CodeChunk[]> => {
        const result = await apiCall(
            // @ts-ignore
            () => wails.SemanticRetrieveContext(JSON.stringify(request)),
            'Failed to retrieve context.',
            { logContext: 'semantic' }
        )
        return parseJsonResponse(result, 'Failed to parse retrieved context.')
    },

    hybridSearch: async (request: SemanticSearchRequest): Promise<SemanticSearchResponse> => {
        const result = await apiCall(
            // @ts-ignore
            () => wails.SemanticHybridSearch(JSON.stringify(request)),
            'Failed to perform hybrid search.',
            { logContext: 'semantic' }
        )
        return parseJsonResponse(result, 'Failed to parse hybrid search response.')
    },
}
