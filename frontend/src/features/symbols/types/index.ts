/**
 * Symbol Browser Types
 */

export type SymbolKind = 'class' | 'function' | 'method' | 'interface' | 'variable' | 'constant' | 'type' | 'enum' | 'property'

export interface Symbol {
    name: string
    kind: SymbolKind
    file: string
    line: number
    endLine?: number
    signature?: string
    parent?: string
    modifiers?: string[]
    children?: Symbol[]
}

export interface SymbolFilter {
    query: string
    kinds: SymbolKind[]
}

export interface SymbolGroup {
    file: string
    symbols: Symbol[]
    expanded: boolean
}

export interface SymbolTreeNode {
    id: string
    symbol: Symbol
    children: SymbolTreeNode[]
    expanded: boolean
    level: number
}
