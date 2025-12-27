/**
 * Symbols Store - Pinia store for symbol browser
 */

import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { symbolsApi } from '../api/symbols.api'
import type { Symbol, SymbolFilter, SymbolGroup, SymbolKind } from '../types'

const logger = useLogger('SymbolsStore')

export const useSymbolsStore = defineStore('symbols', () => {
    const projectStore = useProjectStore()
    const uiStore = useUIStore()
    const { t } = useI18n()

    // State
    const symbols = ref<Symbol[]>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)
    const filter = ref<SymbolFilter>({
        query: '',
        kinds: []
    })
    const expandedFiles = ref<Set<string>>(new Set())
    const lastLoadedProject = ref<string | null>(null)

    // Getters
    const filteredSymbols = computed(() => {
        let result = symbols.value

        // Filter by query
        if (filter.value.query) {
            const query = filter.value.query.toLowerCase()
            result = result.filter(s =>
                s.name.toLowerCase().includes(query) ||
                s.file.toLowerCase().includes(query)
            )
        }

        // Filter by kinds
        if (filter.value.kinds.length > 0) {
            result = result.filter(s => filter.value.kinds.includes(s.kind))
        }

        return result
    })

    const symbolsByFile = computed((): SymbolGroup[] => {
        const grouped = new Map<string, Symbol[]>()

        for (const symbol of filteredSymbols.value) {
            const existing = grouped.get(symbol.file) || []
            existing.push(symbol)
            grouped.set(symbol.file, existing)
        }

        // Sort symbols within each file by line number
        const result: SymbolGroup[] = []
        for (const [file, fileSymbols] of grouped) {
            result.push({
                file,
                symbols: fileSymbols.sort((a, b) => a.line - b.line),
                expanded: expandedFiles.value.has(file)
            })
        }

        // Sort groups by file path
        return result.sort((a, b) => a.file.localeCompare(b.file))
    })

    const symbolCount = computed(() => symbols.value.length)
    const filteredCount = computed(() => filteredSymbols.value.length)
    const fileCount = computed(() => symbolsByFile.value.length)

    const hasSymbols = computed(() => symbols.value.length > 0)
    const hasFilter = computed(() =>
        filter.value.query.length > 0 || filter.value.kinds.length > 0
    )

    // Symbol kind statistics
    const kindStats = computed(() => {
        const stats: Record<SymbolKind, number> = {
            class: 0,
            function: 0,
            method: 0,
            interface: 0,
            variable: 0,
            constant: 0,
            type: 0,
            enum: 0,
            property: 0
        }

        for (const symbol of symbols.value) {
            if (stats[symbol.kind] !== undefined) {
                stats[symbol.kind]++
            }
        }

        return stats
    })

    // Actions
    async function loadSymbols(): Promise<void> {
        const projectPath = projectStore.currentPath
        if (!projectPath) {
            logger.warn('No project selected')
            return
        }

        // Skip if already loaded for this project
        if (lastLoadedProject.value === projectPath && symbols.value.length > 0) {
            logger.info('Symbols already loaded for this project')
            return
        }

        isLoading.value = true
        error.value = null

        try {
            logger.info(`Loading symbols for project: ${projectPath}`)
            const result = await symbolsApi.getProjectSymbols(projectPath)
            symbols.value = result
            lastLoadedProject.value = projectPath

            // Auto-expand first few files
            const files = [...new Set(result.map(s => s.file))].slice(0, 3)
            expandedFiles.value = new Set(files)

            logger.info(`Loaded ${result.length} symbols from ${files.length} files`)
        } catch (err) {
            logger.error('Failed to load symbols:', err)
            error.value = err instanceof Error ? err.message : 'Failed to load symbols'
            uiStore.addToast(t('symbols.loadError'), 'error')
        } finally {
            isLoading.value = false
        }
    }

    async function searchSymbols(query: string): Promise<void> {
        const projectPath = projectStore.currentPath
        if (!projectPath) return

        if (!query || query.length < 2) {
            // Reset to full list
            filter.value.query = ''
            return
        }

        isLoading.value = true
        error.value = null

        try {
            const kindFilter = filter.value.kinds.length === 1
                ? filter.value.kinds[0]
                : undefined

            const result = await symbolsApi.searchSymbols(query, projectPath, kindFilter)
            symbols.value = result
            filter.value.query = query

            // Expand all files in search results
            expandedFiles.value = new Set(result.map(s => s.file))

            logger.info(`Search found ${result.length} symbols`)
        } catch (err) {
            logger.error('Failed to search symbols:', err)
            error.value = err instanceof Error ? err.message : 'Search failed'
        } finally {
            isLoading.value = false
        }
    }

    function setFilter(newFilter: Partial<SymbolFilter>): void {
        filter.value = { ...filter.value, ...newFilter }
    }

    function toggleKindFilter(kind: SymbolKind): void {
        const kinds = [...filter.value.kinds]
        const index = kinds.indexOf(kind)

        if (index === -1) {
            kinds.push(kind)
        } else {
            kinds.splice(index, 1)
        }

        filter.value.kinds = kinds
    }

    function clearFilter(): void {
        filter.value = { query: '', kinds: [] }
    }

    function toggleFileExpanded(file: string): void {
        if (expandedFiles.value.has(file)) {
            expandedFiles.value.delete(file)
        } else {
            expandedFiles.value.add(file)
        }
        // Trigger reactivity
        expandedFiles.value = new Set(expandedFiles.value)
    }

    function expandAll(): void {
        expandedFiles.value = new Set(symbolsByFile.value.map(g => g.file))
    }

    function collapseAll(): void {
        expandedFiles.value = new Set()
    }

    function reset(): void {
        symbols.value = []
        filter.value = { query: '', kinds: [] }
        expandedFiles.value = new Set()
        error.value = null
        lastLoadedProject.value = null
    }

    return {
        // State
        symbols,
        isLoading,
        error,
        filter,
        expandedFiles,

        // Getters
        filteredSymbols,
        symbolsByFile,
        symbolCount,
        filteredCount,
        fileCount,
        hasSymbols,
        hasFilter,
        kindStats,

        // Actions
        loadSymbols,
        searchSymbols,
        setFilter,
        toggleKindFilter,
        clearFilter,
        toggleFileExpanded,
        expandAll,
        collapseAll,
        reset
    }
})
