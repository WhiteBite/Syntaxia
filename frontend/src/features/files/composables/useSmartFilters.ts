/**
 * Composable for smart filters based on project structure
 */
import { useI18n } from '@/composables/useI18n'
import { useProjectStore } from '@/stores/project.store'
import type { ProjectStructure } from '@/types/dto'
import { computed, ref, watch } from 'vue'
import { GetProjectStructure } from '../../../../wailsjs/go/main/App'
import {
    frameworkFiltersConfig,
    languageConfig,
    languageExtensions,
} from '../constants/filterConfig'
import type { LanguageFilter, SmartFilter } from '../model/types'

const MIN_LANGUAGE_PERCENTAGE = 3
const MIN_LANGUAGE_FILE_COUNT = 2

export function useSmartFilters() {
    const { t } = useI18n()
    const projectStore = useProjectStore()
    const projectStructure = ref<ProjectStructure | null>(null)
    const isLoading = ref(false)

    const languageFilters = computed<LanguageFilter[]>(() => {
        if (!projectStructure.value?.languages?.length) return []

        return projectStructure.value.languages
            .filter(
                (lang) =>
                    lang.percentage > MIN_LANGUAGE_PERCENTAGE ||
                    lang.fileCount > MIN_LANGUAGE_FILE_COUNT
            )
            .map((lang) => ({
                id: `lang-${lang.name.toLowerCase().replace(/[^a-z0-9]/g, '')}`,
                label: lang.name,
                language: lang.name,
                icon: languageConfig[lang.name]?.icon || '📄',
                extensions: languageExtensions[lang.name] || [],
                fileCount: lang.fileCount,
                percentage: lang.percentage,
                primary: lang.primary,
                category: 'lang' as const,
                shortLabel: lang.name,
            }))
            .filter((f) => f.extensions.length > 0)
    })

    const smartFilters = computed<SmartFilter[]>(() => {
        if (!projectStructure.value?.frameworks?.length) return []

        const filters: SmartFilter[] = []
        for (const fw of projectStructure.value.frameworks) {
            const fwFilters = frameworkFiltersConfig[fw.name]
            if (fwFilters) {
                filters.push(
                    ...fwFilters.map((f) => ({
                        ...f,
                        label: t(f.labelKey) || f.labelKey,
                        shortLabel: t(f.labelKey) || f.labelKey,
                    }))
                )
            }
        }
        return filters
    })

    const detectedFrameworks = computed(() => {
        return projectStructure.value?.frameworks?.map((f) => f.name) || []
    })

    async function loadProjectStructure(): Promise<void> {
        if (!projectStore.currentPath) return

        isLoading.value = true
        try {
            projectStructure.value = (await GetProjectStructure(
                projectStore.currentPath
            )) as unknown as ProjectStructure
        } catch {
            // Silently fail - project structure is optional
            projectStructure.value = null
        } finally {
            isLoading.value = false
        }
    }

    // Auto-load on project change
    watch(
        () => projectStore.currentPath,
        (path) => {
            if (path) {
                loadProjectStructure()
            } else {
                projectStructure.value = null
            }
        },
        { immediate: true }
    )

    return {
        projectStructure,
        languageFilters,
        smartFilters,
        detectedFrameworks,
        isLoading,
        loadProjectStructure,
    }
}
