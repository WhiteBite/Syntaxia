import { useI18n } from '@/composables/useI18n'
import { apiService, type SmartSuggestion } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { computed, ref, watch, type Ref } from 'vue'

export interface AffectedFile {
    path: string
    type: 'direct' | 'transitive'
    dependents: number
}

export interface ImpactResult {
    totalDependents: number
    aggregateRisk: number
    riskLevel: 'low' | 'medium' | 'high'
    affectedFiles: AffectedFile[]
    relatedTests: string[]
}

interface UseAnalysisStatusOptions {
    selectedFiles: Ref<string[]>
    onAddFiles: (files: string[]) => void
}

const FETCH_TIMEOUT_MS = 5000

async function fetchWithTimeout<T>(promise: Promise<T>, timeoutMs: number): Promise<T | null> {
    const timeout = new Promise<null>((resolve) => setTimeout(() => resolve(null), timeoutMs))
    return Promise.race([promise, timeout])
}

export function useAnalysisStatus(options: UseAnalysisStatusOptions) {
    const { selectedFiles, onAddFiles } = options
    const { t } = useI18n()
    const projectStore = useProjectStore()

    // State
    const suggestions = ref<SmartSuggestion[]>([])
    const impactResult = ref<ImpactResult | null>(null)
    const selectedRelated = ref<Set<string>>(new Set())
    const isLoadingRelated = ref(false)
    const showRelatedPopup = ref(false)
    const showImpactPopup = ref(false)

    // Computed
    const projectPath = computed(() => projectStore.currentPath || '')
    const hasSelectedFiles = computed(() => selectedFiles.value.length > 0)
    const relatedCount = computed(() => suggestions.value.length)
    const dependentCount = computed(() => impactResult.value?.totalDependents || 0)

    const shouldShowBar = computed(() =>
        hasSelectedFiles.value && (isLoadingRelated.value || relatedCount.value > 0 || dependentCount.value > 0)
    )

    // Debounce timer
    let fetchTimer: ReturnType<typeof setTimeout> | null = null

    // Fetch related files
    async function fetchRelated() {
        if (!projectPath.value || selectedFiles.value.length === 0) {
            suggestions.value = []
            return
        }

        isLoadingRelated.value = true
        try {
            const result = await fetchWithTimeout(
                apiService.getSmartSuggestions(projectPath.value, selectedFiles.value),
                FETCH_TIMEOUT_MS
            )
            if (result) {
                suggestions.value = result.suggestions
                selectedRelated.value = new Set(result.suggestions.map(s => s.path))
            } else {
                suggestions.value = []
            }
        } catch {
            suggestions.value = []
        } finally {
            isLoadingRelated.value = false
        }
    }

    // Fetch impact analysis
    async function fetchImpact() {
        if (!projectPath.value || selectedFiles.value.length === 0) {
            impactResult.value = null
            return
        }

        try {
            const result = await fetchWithTimeout(
                apiService.getImpactPreview(projectPath.value, selectedFiles.value),
                FETCH_TIMEOUT_MS
            )
            impactResult.value = result
        } catch {
            impactResult.value = null
        }
    }

    // Watch selected files
    watch(selectedFiles, () => {
        if (fetchTimer) clearTimeout(fetchTimer)

        if (selectedFiles.value.length === 0) {
            suggestions.value = []
            impactResult.value = null
            return
        }

        fetchTimer = setTimeout(() => {
            fetchRelated()
            fetchImpact()
        }, 500)
    }, { immediate: true, deep: true })

    // Actions
    function toggleRelated(path: string) {
        if (selectedRelated.value.has(path)) {
            selectedRelated.value.delete(path)
        } else {
            selectedRelated.value.add(path)
        }
        selectedRelated.value = new Set(selectedRelated.value)
    }

    function addSelectedRelated() {
        if (selectedRelated.value.size > 0) {
            onAddFiles(Array.from(selectedRelated.value))
            suggestions.value = suggestions.value.filter(s => !selectedRelated.value.has(s.path))
            selectedRelated.value.clear()
            showRelatedPopup.value = false
        }
    }

    // Helpers
    function getSourceLabel(source: string): string {
        switch (source) {
            case 'git': return t('context.sourceGitShort')
            case 'arch': return t('context.sourceArchShort')
            case 'semantic': return t('context.sourceSemanticShort')
            default: return ''
        }
    }

    function getSourceBadgeClass(source: string): string {
        switch (source) {
            case 'git': return 'badge-git'
            case 'arch': return 'badge-arch'
            case 'semantic': return 'badge-semantic'
            default: return 'badge-default'
        }
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
            'rs': 'icon-rs',
            'json': 'icon-json',
            'yaml': 'icon-json',
            'yml': 'icon-json',
            'css': 'icon-css',
            'scss': 'icon-css',
            'html': 'icon-html',
        }
        return classMap[ext] || 'icon-default'
    }

    function getFileName(path: string): string {
        return path.split('/').pop() || path
    }

    function getFilePath(path: string): string {
        const parts = path.split('/')
        if (parts.length <= 1) return ''
        return parts.slice(0, -1).join('/')
    }

    function getRiskClass(level: string): string {
        switch (level) {
            case 'low': return 'risk-low'
            case 'medium': return 'risk-medium'
            case 'high': return 'risk-high'
            default: return ''
        }
    }

    function getRiskBarClass(level: string): string {
        switch (level) {
            case 'low': return 'bg-green-500'
            case 'medium': return 'bg-amber-500'
            case 'high': return 'bg-red-500'
            default: return 'bg-gray-500'
        }
    }

    function getRiskLabel(level: string): string {
        switch (level) {
            case 'low': return t('context.riskLow')
            case 'medium': return t('context.riskMedium')
            case 'high': return t('context.riskHigh')
            default: return ''
        }
    }

    return {
        // State
        suggestions,
        impactResult,
        selectedRelated,
        isLoadingRelated,
        showRelatedPopup,
        showImpactPopup,

        // Computed
        shouldShowBar,
        relatedCount,
        dependentCount,

        // Actions
        toggleRelated,
        addSelectedRelated,

        // Helpers
        getSourceLabel,
        getSourceBadgeClass,
        getFileIconClass,
        getFileName,
        getFilePath,
        getRiskClass,
        getRiskBarClass,
        getRiskLabel,
    }
}
