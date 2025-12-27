/**
 * Types for @ mentions system in chat
 * Supports file, folder, symbol, git, and docs mentions
 */

export type MentionType = 'file' | 'folder' | 'symbol' | 'git' | 'docs'

export interface Mention {
    type: MentionType
    value: string       // e.g. "auth/AuthService.ts" or "diff"
    display: string     // e.g. "AuthService.ts" or "Git Diff"
    startIndex: number  // position in text
    endIndex: number
}

export interface MentionSuggestion {
    type: MentionType
    value: string
    display: string
    icon: string
    description?: string
    category?: string
}

export interface ParsedInput {
    text: string            // text without mentions
    mentions: Mention[]     // found mentions
    resolvedFiles: string[] // files for context
}

export interface MentionPattern {
    type: MentionType
    regex: RegExp
    prefix: string
}

/** Autocomplete state for mention dropdown */
export interface AutocompleteState {
    isActive: boolean
    query: string
    position: { top: number; left: number }
    triggerIndex: number
}

/** Category for grouping suggestions in dropdown */
export interface MentionCategory {
    id: MentionType | 'all'
    label: string
    icon: string
}

/** Constants for mention parsing */
export const MENTION_PATTERNS: MentionPattern[] = [
    { type: 'folder', regex: /@folder:([^\s]+)/g, prefix: '@folder:' },
    { type: 'symbol', regex: /@symbol:([^\s]+)/g, prefix: '@symbol:' },
    { type: 'git', regex: /@git:([^\s]+)/g, prefix: '@git:' },
    { type: 'docs', regex: /@docs\b/g, prefix: '@docs' },
    { type: 'file', regex: /@([^\s@:]+\.[a-zA-Z0-9]+)/g, prefix: '@' },
]

/** Git mention subtypes */
export type GitMentionSubtype = 'diff' | 'history' | 'staged' | 'branch'

export const GIT_MENTION_OPTIONS: Array<{ value: GitMentionSubtype; label: string; icon: string }> = [
    { value: 'diff', label: 'Current diff', icon: 'i-lucide-git-compare' },
    { value: 'history', label: 'Recent history', icon: 'i-lucide-history' },
    { value: 'staged', label: 'Staged changes', icon: 'i-lucide-git-commit' },
    { value: 'branch', label: 'Branch info', icon: 'i-lucide-git-branch' },
]

/** Docs mention options */
export const DOCS_MENTION_OPTIONS = [
    { value: 'readme', label: 'README', icon: 'i-lucide-book-open' },
    { value: 'docs', label: 'Documentation', icon: 'i-lucide-file-text' },
]
