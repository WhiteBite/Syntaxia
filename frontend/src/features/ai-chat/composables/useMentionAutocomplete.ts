/**
 * Autocomplete composable for @ mentions
 * Handles cursor detection, suggestion filtering, and keyboard navigation
 */

import { useI18n } from '@/composables/useI18n'
import { filesApi, semanticApi } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { computed, ref, watch, type Ref } from 'vue'
import type {
    AutocompleteState,
    MentionCategory,
    MentionSuggestion,
    MentionType,
} from '../types/mentions'
import { GIT_MENTION_OPTIONS } from '../types/mentions'

const DEBOUNCE_MS = 200
const MAX_SUGGESTIONS = 10

export interface UseMentionAutocompleteOptions {
    inputText: Ref<string>
    cursorPosition: Ref<number>
    inputElement: Ref<HTMLTextAreaElement | HTMLInputElement | null>
}

export function useMentionAutocomplete(options: UseMentionAutocompleteOptions) {
    const { inputText, cursorPosition, inputElement } = options
    const { t } = useI18n()
    const projectStore = useProjectStore()

    // State
    const isOpen = ref(false)
    const suggestions = ref<MentionSuggestion[]>([])
    const selectedIndex = ref(0)
    const isLoading = ref(false)
    const activeCategory = ref<MentionType | 'all'>('all')

    const autocompleteState = ref<AutocompleteState>({
        isActive: false,
        query: '',
        position: { top: 0, left: 0 },
        triggerIndex: -1,
    })

    /** Categories for dropdown */
    const categories: MentionCategory[] = [
        { id: 'all', label: t('mentions.categories.all'), icon: 'i-lucide-layers' },
        { id: 'file', label: t('mentions.categories.files'), icon: 'i-lucide-file' },
        { id: 'folder', label: t('mentions.categories.folders'), icon: 'i-lucide-folder' },
        { id: 'symbol', label: t('mentions.categories.symbols'), icon: 'i-lucide-code' },
        { id: 'git', label: t('mentions.categories.git'), icon: 'i-lucide-git-branch' },
        { id: 'docs', label: t('mentions.categories.docs'), icon: 'i-lucide-book-open' },
    ]

    /** Detect @ trigger and extract query */
    function detectTrigger(): { triggered: boolean; query: string; triggerIndex: number } {
        const text = inputText.value
        const pos = cursorPosition.value

        // Find last @ before cursor
        const beforeCursor = text.slice(0, pos)
        const atIndex = beforeCursor.lastIndexOf('@')

        if (atIndex === -1) {
            return { triggered: false, query: '', triggerIndex: -1 }
        }

        // Check if @ is at start or after whitespace
        if (atIndex > 0 && !/\s/.test(text[atIndex - 1])) {
            return { triggered: false, query: '', triggerIndex: -1 }
        }

        // Extract query after @
        const query = beforeCursor.slice(atIndex + 1)

        // Don't trigger if there's a space in query (mention completed)
        if (query.includes(' ')) {
            return { triggered: false, query: '', triggerIndex: -1 }
        }

        return { triggered: true, query, triggerIndex: atIndex }
    }

    /** Calculate dropdown position */
    function calculatePosition(): { top: number; left: number } {
        const el = inputElement.value
        if (!el) return { top: 0, left: 0 }

        const rect = el.getBoundingClientRect()

        // Simple positioning below input
        return {
            top: rect.bottom + 4,
            left: rect.left,
        }
    }

    /** Fetch suggestions based on query */
    async function fetchSuggestions(query: string): Promise<MentionSuggestion[]> {
        const projectPath = projectStore.currentPath
        if (!projectPath) return getDefaultSuggestions()

        const results: MentionSuggestion[] = []

        // Determine what to search based on query prefix
        const lowerQuery = query.toLowerCase()

        // Check for typed prefixes
        if (lowerQuery.startsWith('folder:')) {
            const folderQuery = query.slice(7)
            results.push(...await searchFolders(folderQuery, projectPath))
        } else if (lowerQuery.startsWith('symbol:')) {
            const symbolQuery = query.slice(7)
            results.push(...await searchSymbols(symbolQuery, projectPath))
        } else if (lowerQuery.startsWith('git:')) {
            results.push(...getGitSuggestions(query.slice(4)))
        } else if (lowerQuery === 'docs' || lowerQuery.startsWith('docs')) {
            results.push(...getDocsSuggestions())
        } else {
            // General search - files first, then others
            const [files, folders, symbols] = await Promise.all([
                searchFiles(query, projectPath),
                searchFolders(query, projectPath),
                query.length >= 2 ? searchSymbols(query, projectPath) : Promise.resolve([]),
            ])

            results.push(...files, ...folders, ...symbols)

            // Add git and docs if query matches
            if ('git'.startsWith(lowerQuery)) {
                results.push(...getGitSuggestions(''))
            }
            if ('docs'.startsWith(lowerQuery)) {
                results.push(...getDocsSuggestions())
            }
        }

        return results.slice(0, MAX_SUGGESTIONS)
    }

    /** Search files by name */
    async function searchFiles(query: string, projectPath: string): Promise<MentionSuggestion[]> {
        try {
            const allFiles = await filesApi.listFiles(projectPath)
            const flatFiles = flattenFileTree(allFiles)
            const lowerQuery = query.toLowerCase()

            return flatFiles
                .filter(f => {
                    const filename = f.split('/').pop() || f
                    return filename.toLowerCase().includes(lowerQuery)
                })
                .slice(0, 5)
                .map(f => ({
                    type: 'file' as MentionType,
                    value: f,
                    display: f.split('/').pop() || f,
                    icon: getFileIcon(f),
                    description: f,
                    category: 'file',
                }))
        } catch {
            return []
        }
    }

    /** Search folders */
    async function searchFolders(query: string, projectPath: string): Promise<MentionSuggestion[]> {
        try {
            const allFiles = await filesApi.listFiles(projectPath)
            const folders = extractFolders(allFiles)
            const lowerQuery = query.toLowerCase()

            return folders
                .filter(f => f.toLowerCase().includes(lowerQuery))
                .slice(0, 5)
                .map(f => ({
                    type: 'folder' as MentionType,
                    value: f,
                    display: f.split('/').pop() || f,
                    icon: 'i-lucide-folder',
                    description: f,
                    category: 'folder',
                }))
        } catch {
            return []
        }
    }

    /** Search symbols using semantic search */
    async function searchSymbols(query: string, projectPath: string): Promise<MentionSuggestion[]> {
        try {
            const result = await semanticApi.search({
                query,
                projectRoot: projectPath,
                topK: 5,
            })

            return (result.results || []).map(r => ({
                type: 'symbol' as MentionType,
                value: r.chunk.symbolName || query,
                display: r.chunk.symbolName || query,
                icon: getSymbolIcon(r.chunk.symbolKind),
                description: r.chunk.filePath,
                category: 'symbol',
            }))
        } catch {
            return []
        }
    }

    /** Get git-related suggestions */
    function getGitSuggestions(query: string): MentionSuggestion[] {
        const lowerQuery = query.toLowerCase()
        return GIT_MENTION_OPTIONS
            .filter(opt => opt.value.includes(lowerQuery) || opt.label.toLowerCase().includes(lowerQuery))
            .map(opt => ({
                type: 'git' as MentionType,
                value: opt.value,
                display: opt.label,
                icon: opt.icon,
                description: t(`mentions.git.${opt.value}`),
                category: 'git',
            }))
    }

    /** Get docs suggestions */
    function getDocsSuggestions(): MentionSuggestion[] {
        return [{
            type: 'docs' as MentionType,
            value: 'docs',
            display: t('mentions.docs'),
            icon: 'i-lucide-book-open',
            description: t('mentions.docsDescription'),
            category: 'docs',
        }]
    }

    /** Get default suggestions when no query */
    function getDefaultSuggestions(): MentionSuggestion[] {
        return [
            {
                type: 'file',
                value: '',
                display: t('mentions.searchFiles'),
                icon: 'i-lucide-file-search',
                description: t('mentions.searchFilesHint'),
            },
            ...getGitSuggestions(''),
            ...getDocsSuggestions(),
        ]
    }

    /** Flatten file tree to paths */
    function flattenFileTree(nodes: { path: string; children?: unknown[] }[]): string[] {
        const result: string[] = []

        function traverse(items: { path: string; children?: unknown[] }[]) {
            for (const item of items) {
                if (!item.children || item.children.length === 0) {
                    result.push(item.path)
                } else {
                    traverse(item.children as { path: string; children?: unknown[] }[])
                }
            }
        }

        traverse(nodes)
        return result
    }

    /** Extract folder paths from file tree */
    function extractFolders(nodes: { path: string; children?: unknown[] }[]): string[] {
        const folders: string[] = []

        function traverse(items: { path: string; children?: unknown[] }[]) {
            for (const item of items) {
                if (item.children && item.children.length > 0) {
                    folders.push(item.path)
                    traverse(item.children as { path: string; children?: unknown[] }[])
                }
            }
        }

        traverse(nodes)
        return folders
    }

    /** Get icon for file type */
    function getFileIcon(filename: string): string {
        const ext = filename.split('.').pop()?.toLowerCase()
        const iconMap: Record<string, string> = {
            ts: 'i-lucide-file-code',
            tsx: 'i-lucide-file-code',
            js: 'i-lucide-file-code',
            jsx: 'i-lucide-file-code',
            vue: 'i-lucide-file-code',
            go: 'i-lucide-file-code',
            py: 'i-lucide-file-code',
            json: 'i-lucide-file-json',
            md: 'i-lucide-file-text',
            css: 'i-lucide-palette',
            scss: 'i-lucide-palette',
        }
        return iconMap[ext || ''] || 'i-lucide-file'
    }

    /** Get icon for symbol kind */
    function getSymbolIcon(kind?: string): string {
        const iconMap: Record<string, string> = {
            function: 'i-lucide-function-square',
            class: 'i-lucide-box',
            method: 'i-lucide-braces',
            interface: 'i-lucide-shapes',
            variable: 'i-lucide-variable',
            constant: 'i-lucide-hash',
        }
        return iconMap[kind?.toLowerCase() || ''] || 'i-lucide-code'
    }

    /** Open autocomplete dropdown */
    function open() {
        isOpen.value = true
        selectedIndex.value = 0
    }

    /** Close autocomplete dropdown */
    function close() {
        isOpen.value = false
        suggestions.value = []
        selectedIndex.value = 0
        autocompleteState.value.isActive = false
    }

    /** Select next suggestion */
    function selectNext() {
        if (suggestions.value.length === 0) return
        selectedIndex.value = (selectedIndex.value + 1) % suggestions.value.length
    }

    /** Select previous suggestion */
    function selectPrevious() {
        if (suggestions.value.length === 0) return
        selectedIndex.value = selectedIndex.value === 0
            ? suggestions.value.length - 1
            : selectedIndex.value - 1
    }

    /** Get currently selected suggestion */
    const selectedSuggestion = computed(() => {
        return suggestions.value[selectedIndex.value] || null
    })

    /** Filter suggestions by category */
    const filteredSuggestions = computed(() => {
        if (activeCategory.value === 'all') {
            return suggestions.value
        }
        return suggestions.value.filter(s => s.type === activeCategory.value)
    })

    /** Set active category */
    function setCategory(category: MentionType | 'all') {
        activeCategory.value = category
        selectedIndex.value = 0
    }

    // Watch for trigger detection
    let debounceTimer: ReturnType<typeof setTimeout> | null = null

    watch([inputText, cursorPosition], async () => {
        const { triggered, query, triggerIndex } = detectTrigger()

        if (!triggered) {
            close()
            return
        }

        autocompleteState.value = {
            isActive: true,
            query,
            position: calculatePosition(),
            triggerIndex,
        }

        // Debounce suggestion fetching
        if (debounceTimer) clearTimeout(debounceTimer)
        debounceTimer = setTimeout(async () => {
            isLoading.value = true
            try {
                suggestions.value = await fetchSuggestions(query)
                if (suggestions.value.length > 0 && !isOpen.value) {
                    open()
                }
            } finally {
                isLoading.value = false
            }
        }, DEBOUNCE_MS)
    })

    return {
        // State
        isOpen,
        isLoading,
        suggestions,
        selectedIndex,
        selectedSuggestion,
        filteredSuggestions,
        autocompleteState,
        activeCategory,
        categories,
        // Methods
        open,
        close,
        selectNext,
        selectPrevious,
        setCategory,
        fetchSuggestions,
    }
}
