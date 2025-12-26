<template>
    <div class="tool-call-container" :class="{ 'tool-call-collapsed': !expanded }">
        <div class="tool-call-header" @click="expanded = !expanded">
            <!-- Category Icon -->
            <div class="tool-call-category" :title="categoryLabel">
                <span class="category-icon">{{ categoryIcon }}</span>
            </div>

            <!-- Status Icon -->
            <div class="tool-call-icon" :class="statusClass">
                <svg v-if="status === 'executing'" class="w-3.5 h-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                    </path>
                </svg>
                <svg v-else-if="status === 'completed'" class="w-3.5 h-3.5" fill="none" stroke="currentColor"
                    viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
            </div>

            <div class="tool-call-info">
                <span class="tool-call-name">{{ toolName }}</span>
                <span class="tool-call-brief" v-if="briefArgs">{{ briefArgs }}</span>
            </div>

            <div class="tool-call-meta">
                <span class="tool-call-duration" v-if="duration">{{ formatDuration(duration) }}</span>
                <svg class="w-3.5 h-3.5 text-gray-400 transition-transform" :class="{ 'rotate-180': expanded }" fill="none"
                    stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
            </div>
        </div>

        <div v-if="expanded" class="tool-call-details">
            <div v-if="Object.keys(args).length > 0" class="tool-call-section">
                <div class="tool-call-section-title">{{ t('toolCalls.arguments') }}</div>
                <pre class="tool-call-code">{{ formatArgs(args) }}</pre>
            </div>

            <div v-if="result" class="tool-call-section">
                <div class="tool-call-section-title">{{ t('toolCalls.result') }}</div>
                <pre class="tool-call-code tool-call-result" :class="{ 'tool-call-code-expanded': resultExpanded }">{{ displayResult }}</pre>
                <button v-if="isResultLong" @click.stop="resultExpanded = !resultExpanded" class="tool-call-expand-btn">
                    {{ resultExpanded ? t('toolCalls.showLess') : t('toolCalls.showMore') }}
                </button>
            </div>

            <div v-if="error" class="tool-call-section">
                <div class="tool-call-section-title text-red-400">{{ t('error.generic') }}</div>
                <pre class="tool-call-code tool-call-error">{{ error }}</pre>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';
import { computed, ref } from 'vue';

const { t } = useI18n()

// Tool categories for grouping and icons
const TOOL_CATEGORIES: Record<string, { icon: string; label: string; tools: string[] }> = {
    file: {
        icon: '📁',
        label: 'File',
        tools: ['search_files', 'search_content', 'read_file', 'list_directory', 'get_file_info', 'list_functions']
    },
    git: {
        icon: '🔀',
        label: 'Git',
        tools: ['git_status', 'git_diff', 'git_log', 'git_changed_files', 'git_co_changed', 'git_suggest_context']
    },
    symbol: {
        icon: '🔣',
        label: 'Symbol',
        tools: ['list_symbols', 'search_symbols', 'find_definition', 'find_references', 'get_symbol_info', 'get_class_hierarchy', 'get_imports', 'get_widget_tree']
    },
    analysis: {
        icon: '📊',
        label: 'Analysis',
        tools: ['get_callers', 'get_callees', 'get_impact', 'get_call_chain', 'get_dependent_files', 'get_change_risk']
    },
    project: {
        icon: '🏗️',
        label: 'Project',
        tools: ['detect_architecture', 'detect_frameworks', 'detect_conventions', 'get_project_structure', 'get_related_layers', 'suggest_related_files']
    },
    memory: {
        icon: '💾',
        label: 'Memory',
        tools: ['save_context', 'find_context', 'get_recent_contexts', 'set_preference', 'get_preferences']
    },
    semantic: {
        icon: '🔍',
        label: 'Semantic',
        tools: ['semantic_search']
    }
}

interface Props {
    toolName: string
    args: Record<string, unknown>
    status: 'executing' | 'completed' | 'failed'
    result?: string
    error?: string
    duration?: number
}

const props = defineProps<Props>()
const expanded = ref(false)
const resultExpanded = ref(false)

const MAX_RESULT_LENGTH = 300

const statusClass = computed(() => ({
    'tool-call-icon-executing': props.status === 'executing',
    'tool-call-icon-completed': props.status === 'completed',
    'tool-call-icon-failed': props.status === 'failed',
}))

// Get category for tool
const category = computed(() => {
    for (const [key, cat] of Object.entries(TOOL_CATEGORIES)) {
        if (cat.tools.includes(props.toolName)) return key
    }
    return 'other'
})

const categoryIcon = computed(() => TOOL_CATEGORIES[category.value]?.icon || '⚙️')
const categoryLabel = computed(() => TOOL_CATEGORIES[category.value]?.label || 'Other')

// Brief args for collapsed view
const briefArgs = computed(() => {
    const keys = Object.keys(props.args)
    if (keys.length === 0) return ''
    
    // Show first meaningful arg value
    const firstKey = keys[0]
    const val = props.args[firstKey]
    if (typeof val === 'string' && val.length > 0) {
        return val.length > 30 ? val.slice(0, 30) + '...' : val
    }
    return ''
})

const isResultLong = computed(() => (props.result?.length || 0) > MAX_RESULT_LENGTH)

const displayResult = computed(() => {
    if (!props.result) return ''
    if (resultExpanded.value || !isResultLong.value) return props.result
    return props.result.slice(0, MAX_RESULT_LENGTH) + '\n...'
})

function formatArgs(args: Record<string, unknown>): string {
    return JSON.stringify(args, null, 2)
}

function formatDuration(ms: number): string {
    if (ms < 1000) return `${ms}ms`
    return `${(ms / 1000).toFixed(1)}s`
}
</script>

<style scoped>
.tool-call-container {
    @apply bg-gray-800/50 border border-gray-700/50 rounded-lg overflow-hidden my-1;
}

.tool-call-collapsed {
    @apply bg-gray-800/30;
}

.tool-call-header {
    @apply flex items-center gap-2 px-2.5 py-1.5 cursor-pointer hover:bg-gray-700/30 transition-colors;
}

.tool-call-category {
    @apply flex-shrink-0;
}

.category-icon {
    @apply text-sm;
}

.tool-call-icon {
    @apply w-5 h-5 rounded flex items-center justify-center flex-shrink-0;
}

.tool-call-icon-executing {
    @apply bg-blue-500/20 text-blue-400;
}

.tool-call-icon-completed {
    @apply bg-emerald-500/20 text-emerald-400;
}

.tool-call-icon-failed {
    @apply bg-red-500/20 text-red-400;
}

.tool-call-info {
    @apply flex-1 flex items-center gap-2 min-w-0;
}

.tool-call-name {
    @apply text-sm font-medium text-gray-200;
}

.tool-call-brief {
    @apply text-xs text-gray-400 truncate max-w-32;
}

.tool-call-meta {
    @apply flex items-center gap-2 flex-shrink-0;
}

.tool-call-duration {
    @apply text-xs text-gray-400 font-mono;
}

.tool-call-details {
    @apply border-t border-gray-700/50 p-2.5 space-y-2;
}

.tool-call-section {
    @apply space-y-1;
}

.tool-call-section-title {
    @apply text-xs font-medium text-gray-400 uppercase tracking-wide;
}

.tool-call-code {
    @apply text-xs font-mono bg-gray-900/50 rounded p-2 overflow-x-auto text-gray-300 max-h-40 overflow-y-auto;
    white-space: pre-wrap;
    word-break: break-word;
}

.tool-call-code-expanded {
    @apply max-h-96;
}

.tool-call-result {
    @apply text-emerald-300/80;
}

.tool-call-error {
    @apply text-red-300/80;
}

.tool-call-expand-btn {
    @apply text-xs text-indigo-400 hover:text-indigo-300 mt-1 cursor-pointer;
}
</style>
