/**
 * Types for Quick Filters feature
 */
import type { QuickFilterConfig } from '@/stores/settings.store'

export type FilterCategory = 'code' | 'test' | 'config' | 'docs' | 'styles'

export interface TypeFilter extends QuickFilterConfig {
    category: FilterCategory
    shortLabel: string
    icon?: string
}

export interface LanguageFilter {
    id: string
    label: string
    language: string
    icon: string
    extensions: string[]
    fileCount: number
    percentage: number
    primary: boolean
    category: 'lang'
    shortLabel: string
}

export interface SmartFilter {
    id: string
    label: string
    shortLabel: string
    icon: string
    extensions: string[]
    patterns: string[]
    framework: string
    category: 'smart'
}

export type AnyFilter = TypeFilter | LanguageFilter | SmartFilter

export type DropdownType = 'types' | 'langs' | 'smart'

export interface FilterState {
    active: Set<string>
    excluded: Set<string>
}
