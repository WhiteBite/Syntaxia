import { computed, ref } from 'vue'

import { estimateTokens } from '@/features/context/lib/context-utils'
import { useTemplateStore } from '@/features/templates/model/template.store'

export interface ContextFilePreview {
    path: string
    tokens: number
    relevance: number  // 0-1
    reason: string
    selected: boolean
}

export interface SystemPromptPreview {
    name: string
    icon: string
    content: string
    tokens: number
}

// localStorage key for prompt expanded state
const PROMPT_EXPANDED_KEY = 'context-preview-prompt-expanded'

export function useContextPreview() {
    const templateStore = useTemplateStore()

    // State
    const files = ref<ContextFilePreview[]>([])
    const tokenLimit = ref(128000)
    const showPrompt = ref(localStorage.getItem(PROMPT_EXPANDED_KEY) === 'true')

    // System Prompt computed
    const systemPrompt = computed<SystemPromptPreview>(() => {
        const template = templateStore.activeTemplate
        if (!template) {
            return {
                name: 'Default',
                icon: '📝',
                content: '',
                tokens: 0
            }
        }

        // Build prompt content preview (role + rules)
        const parts: string[] = []
        if (template.roleContent) {
            parts.push(`## Role\n${template.roleContent}`)
        }
        if (template.rulesContent) {
            parts.push(`## Rules\n${template.rulesContent}`)
        }
        if (templateStore.userRules) {
            parts.push(`## User Rules\n${templateStore.userRules}`)
        }

        const content = parts.join('\n\n')

        return {
            name: template.name,
            icon: template.icon,
            content,
            tokens: estimateTokens(content)
        }
    })

    const promptTokens = computed(() => systemPrompt.value.tokens)

    // Computed
    const selectedFiles = computed(() =>
        files.value.filter(f => f.selected)
    )

    const selectedTokens = computed(() =>
        selectedFiles.value.reduce((sum, f) => sum + f.tokens, 0)
    )

    const totalTokens = computed(() =>
        files.value.reduce((sum, f) => sum + f.tokens, 0)
    )

    const selectedCount = computed(() =>
        selectedFiles.value.length
    )

    // Total tokens including prompt
    const totalWithPrompt = computed(() =>
        selectedTokens.value + promptTokens.value
    )

    const canSend = computed(() =>
        selectedCount.value > 0 && totalWithPrompt.value <= tokenLimit.value
    )

    const tokenUsagePercent = computed(() =>
        Math.min((totalWithPrompt.value / tokenLimit.value) * 100, 100)
    )

    const isOverLimit = computed(() =>
        totalWithPrompt.value > tokenLimit.value
    )

    // Actions
    function setFiles(newFiles: ContextFilePreview[]) {
        files.value = newFiles.map(f => ({ ...f }))
    }

    function setTokenLimit(limit: number) {
        tokenLimit.value = limit
    }

    function toggleFile(path: string) {
        const file = files.value.find(f => f.path === path)
        if (file) {
            file.selected = !file.selected
        }
    }

    function selectAll() {
        files.value.forEach(f => f.selected = true)
    }

    function deselectAll() {
        files.value.forEach(f => f.selected = false)
    }

    function getSelectedPaths(): string[] {
        return selectedFiles.value.map(f => f.path)
    }

    function togglePromptExpanded() {
        showPrompt.value = !showPrompt.value
        localStorage.setItem(PROMPT_EXPANDED_KEY, String(showPrompt.value))
    }

    // Prompt budget warning threshold (20%)
    const PROMPT_BUDGET_WARNING_THRESHOLD = 0.2

    const promptBudgetPercent = computed(() =>
        tokenLimit.value > 0 ? promptTokens.value / tokenLimit.value : 0
    )

    const isPromptOverBudget = computed(() =>
        promptBudgetPercent.value > PROMPT_BUDGET_WARNING_THRESHOLD
    )

    function openTemplateSettings() {
        templateStore.openModal()
    }

    function reset() {
        files.value = []
        // Don't reset showPrompt - it's persisted in localStorage
    }

    return {
        // State
        files,
        tokenLimit,
        showPrompt,
        // System Prompt
        systemPrompt,
        promptTokens,
        // Computed
        selectedFiles,
        selectedTokens,
        totalTokens,
        totalWithPrompt,
        selectedCount,
        canSend,
        tokenUsagePercent,
        isOverLimit,
        promptBudgetPercent,
        isPromptOverBudget,
        // Actions
        setFiles,
        setTokenLimit,
        toggleFile,
        selectAll,
        deselectAll,
        getSelectedPaths,
        togglePromptExpanded,
        openTemplateSettings,
        reset,
    }
}
