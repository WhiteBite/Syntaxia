/**
 * Auto-suggest composable for AI chat
 * Analyzes user input and suggests relevant files
 */

import { useI18n } from '@/composables/useI18n'
import { apiService, type SmartSuggestion } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { computed, ref, watch, type Ref } from 'vue'

const DEBOUNCE_MS = 500
const MIN_INPUT_LENGTH = 10
const MAX_SUGGESTIONS = 5
const FETCH_TIMEOUT_MS = 5000

export interface AutoSuggestOptions {
    inputText: Ref<string>
    currentFiles: Ref<string[]>
    onAddFiles: (files: string[]) => void
}

async function fetchWithTimeout<T>(promise: Promise<T>, timeoutMs: number): Promise<T | null> {
    const timeout = new Promise<null>((resolve) => setTimeout(() => resolve(null), timeoutMs))
    return Promise.race([promise, timeout])
}

export function useAutoSuggest(options: AutoSuggestOptions) {
    const { inputText, currentFiles, onAddFiles } = options
    const { t } = useI18n()
    const projectStore = useProjectStore()

    // State
    const suggestions = ref<SmartSuggestion[]>([])
    const selectedPaths = ref<Set<string>>(new Set())
    const isLoading = ref(false)
    const isVisible = ref(false)
    const isHidden = ref(false) // User manually hid the panel

    // Cache for suggestions
    const cache = ref<Map<string, SmartSuggestion[]>>(new Map())

    // Computed
    const projectPath = computed(() => projectStore.currentPath || '')
    const shouldShow = computed(() =>
        !isHidden.value &&
        inputText.value.length >= MIN_INPUT_LENGTH &&
        (isLoading.value || suggestions.value.length > 0)
    )
    const hasSelectedSuggestions = computed(() => selectedPaths.value.size > 0)

    // Debounce timer
    let debounceTimer: ReturnType<typeof setTimeout> | null = null

    // Generate cache key
    function getCacheKey(text: string): string {
        const normalized = text.toLowerCase().trim().slice(0, 100)
        return `${projectPath.value}:${normalized}`
    }

    // Fetch suggestions from API
    async function fetchSuggestions(text: string) {
        if (!projectPath.value || text.length < MIN_INPUT_LENGTH) {
            suggestions.value = []
            return
        }

        const cacheKey = getCacheKey(text)
        const cached = cache.value.get(cacheKey)
        if (cached) {
            suggestions.value = cached
            selectedPaths.value = new Set(cached.map(s => s.path))
            return
        }

        isLoading.value = true
        try {
            const result = await fetchWithTimeout(
                apiService.getSmartSuggestions(projectPath.value, currentFiles.value, text),
                FETCH_TIMEOUT_MS
            )
            if (result) {
                const filtered = result.suggestions
                    .filter(s => !currentFiles.value.includes(s.path))
                    .slice(0, MAX_SUGGESTIONS)

                suggestions.value = filtered
                selectedPaths.value = new Set(filtered.map(s => s.path))

                // Cache result
                cache.value.set(cacheKey, filtered)

                // Limit cache size
                if (cache.value.size > 50) {
                    const firstKey = cache.value.keys().next().value
                    if (firstKey) cache.value.delete(firstKey)
                }
            } else {
                suggestions.value = []
            }
        } catch {
            suggestions.value = []
        } finally {
            isLoading.value = false
        }
    }

    // Watch input text with debounce
    watch(inputText, (newText) => {
        if (debounceTimer) clearTimeout(debounceTimer)

        if (newText.length < MIN_INPUT_LENGTH) {
            suggestions.value = []
            isVisible.value = false
            return
        }

        isHidden.value = false
        isVisible.value = true

        debounceTimer = setTimeout(() => {
            fetchSuggestions(newText)
        }, DEBOUNCE_MS)
    })

    // Actions
    function toggleSelect(path: string) {
        if (selectedPaths.value.has(path)) {
            selectedPaths.value.delete(path)
        } else {
            selectedPaths.value.add(path)
        }
        selectedPaths.value = new Set(selectedPaths.value)
    }

    function addSelected() {
        if (selectedPaths.value.size > 0) {
            onAddFiles(Array.from(selectedPaths.value))
            suggestions.value = suggestions.value.filter(s => !selectedPaths.value.has(s.path))
            selectedPaths.value.clear()
        }
    }

    function hide() {
        isHidden.value = true
        isVisible.value = false
    }

    function clearCache() {
        cache.value.clear()
    }

    // Helpers
    function getFileName(path: string): string {
        return path.split('/').pop() || path
    }

    function getFilePath(path: string): string {
        const parts = path.split('/')
        if (parts.length <= 1) return ''
        return parts.slice(0, -1).join('/')
    }

    function getFileIconClass(path: string): string {
        const ext = path.split('.').pop()?.toLowerCase() || ''
        const classMap: Record<string, string> = {
            'vue': 'icon-vue',
            'ts': 'icon-ts',
            'tsx': 'icon-ts',
            'js': 'icon-js',
            'jsx': 'icon-js',
            'go': 'icon-go',
            'py': 'icon-py',
            'json': 'icon-json',
            'css': 'icon-css',
            'scss': 'icon-css',
        }
        return classMap[ext] || 'icon-default'
    }

    function getSourceBadgeClass(source: string): string {
        switch (source) {
            case 'git': return 'badge-git'
            case 'arch': return 'badge-arch'
            case 'semantic': return 'badge-semantic'
            default: return 'badge-default'
        }
    }

    function getSourceLabel(source: string): string {
        switch (source) {
            case 'git': return t('context.sourceGitShort')
            case 'arch': return t('context.sourceArchShort')
            case 'semantic': return t('context.sourceSemanticShort')
            default: return ''
        }
    }

    function getRelevancePercent(confidence: number): number {
        return Math.round(confidence * 100)
    }

    return {
        // State
        suggestions,
        selectedPaths,
        isLoading,
        isVisible,

        // Computed
        shouldShow,
        hasSelectedSuggestions,

        // Actions
        toggleSelect,
        addSelected,
        hide,
        clearCache,

        // Helpers
        getFileName,
        getFilePath,
        getFileIconClass,
        getSourceBadgeClass,
        getSourceLabel,
        getRelevancePercent,
    }
}
