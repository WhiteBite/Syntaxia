/**
 * API types for frontend services
 * Centralized type definitions for all API modules
 */

// ============================================
// Git types
// ============================================

export interface CommitInfo {
    hash: string
    subject: string
    author: string
    date: string
}

export interface GitHubBranch {
    name: string
    commit: { sha: string }
}

export interface GitHubCommit {
    sha: string
    commit: {
        message: string
        author: { name: string; date: string }
    }
}

export interface GitLabBranch {
    name: string
    commit: { id: string }
    default: boolean
}

export interface GitLabCommit {
    id: string
    short_id: string
    title: string
    message: string
    author_name: string
    committed_date: string
}

// ============================================
// Agentic Chat types
// ============================================

export interface AgenticChatRequest {
    task: string
    projectRoot: string
    maxTokens?: number
}

export interface AgenticChatResponse {
    response: string
    toolCalls: ToolCallLog[]
    iterations: number
    context: string[]
}

export interface ToolCallLog {
    tool: string
    arguments: string
    result: string
}

// ============================================
// Smart Suggestions types
// ============================================

export interface SmartSuggestion {
    path: string
    source: 'git' | 'arch' | 'semantic'
    reason: string
    confidence: number
}

export interface SmartSuggestionsResult {
    suggestions: SmartSuggestion[]
    total: number
}

// ============================================
// File Quick Info types
// ============================================

export interface FileQuickInfo {
    symbolCount: number
    importCount: number
    dependentCount: number
    changeRisk: number
    riskLevel: 'low' | 'medium' | 'high'
}

// ============================================
// Impact Preview types
// ============================================

export interface ImpactPreviewResult {
    totalDependents: number
    aggregateRisk: number
    riskLevel: 'low' | 'medium' | 'high'
    affectedFiles: AffectedFile[]
    relatedTests: string[]
}

export interface AffectedFile {
    path: string
    type: 'direct' | 'transitive'
    dependents: number
}

// ============================================
// Context Memory types
// ============================================

export interface ContextMemoryEntry {
    id: string
    topic: string
    summary: string
    files: string[]
    createdAt: string
}

// ============================================
// Semantic Search types
// ============================================

export interface SemanticSearchRequest {
    query: string
    projectRoot: string
    topK?: number
    minScore?: number
    searchType?: 'semantic' | 'keyword' | 'hybrid'
    languages?: string[]
    chunkTypes?: string[]
}

export interface SemanticSearchResponse {
    results: SemanticSearchResult[]
    totalResults: number
    queryTime: number
    searchType: string
}

export interface SemanticSearchResult {
    chunk: CodeChunk
    score: number
    highlights?: string[]
    reason?: string
}

export interface CodeChunk {
    id: string
    filePath: string
    content: string
    startLine: number
    endLine: number
    chunkType: 'file' | 'function' | 'class' | 'method' | 'block'
    symbolName?: string
    symbolKind?: string
    language: string
    tokenCount: number
    hash: string
}

export interface FindSimilarRequest {
    filePath: string
    startLine: number
    endLine: number
    topK?: number
    minScore?: number
    excludeSelf?: boolean
}

export interface RetrieveContextRequest {
    query: string
    projectRoot: string
    maxTokens?: number
}

export interface SemanticIndexStats {
    totalChunks: number
    totalFiles: number
    totalTokens: number
    lastUpdated: string
    indexSize: number
    dimensions: number
}

// ============================================
// Qwen types
// ============================================

export interface QwenTaskRequest {
    task: string
    projectRoot: string
    selectedFiles?: string[]
    selectedCode?: string
    sourceFile?: string
    model?: string
    maxTokens?: number
    temperature?: number
}

export interface QwenTaskResponse {
    content: string
    model: string
    tokensUsed: number
    processingTime: string
    contextSummary: QwenContextSummary
    success: boolean
    error?: string
}

export interface QwenContextSummary {
    totalFiles: number
    totalTokens: number
    includedFiles: string[]
    truncatedFiles: string[]
    excludedFiles: string[]
}

export interface QwenContextPreview {
    totalFiles: number
    totalTokens: number
    files: QwenFilePreview[]
    truncatedFiles: string[]
    excludedFiles: string[]
    callStackInfo?: QwenCallStackInfo
    relevanceScores: Record<string, number>
}

export interface QwenFilePreview {
    path: string
    tokens: number
    relevance: number
    reason: string
}

export interface QwenCallStackInfo {
    rootSymbol: string
    callers: string[]
    callees: string[]
    dependencies: string[]
}

export interface QwenModelInfo {
    id: string
    name: string
    description: string
    maxContext: number
    recommended: boolean
}
