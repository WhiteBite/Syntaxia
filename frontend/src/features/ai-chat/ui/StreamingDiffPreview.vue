<template>
    <div class="streaming-diff-preview">
        <!-- Header -->
        <div class="diff-header" @click="toggleExpand">
            <div class="diff-file-info">
                <span class="diff-operation" :class="operationClass">
                    {{ operationIcon }}
                </span>
                <span class="diff-filename">{{ fileName }}</span>
            </div>
            <div class="diff-stats">
                <span v-if="stats.additions > 0" class="stat-add">+{{ stats.additions }}</span>
                <span v-if="stats.deletions > 0" class="stat-remove">-{{ stats.deletions }}</span>
                <svg 
                    class="expand-icon" 
                    :class="{ 'rotate-180': isExpanded }" 
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24"
                >
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
            </div>
        </div>

        <!-- Diff Content -->
        <div class="diff-content" :class="{ 'diff-expanded': isExpanded }">
            <div 
                v-for="(line, idx) in displayLines" 
                :key="idx"
                class="diff-line"
                :class="lineClass(line.type)"
            >
                <span class="line-number">{{ formatLineNumber(line) }}</span>
                <span class="line-prefix">{{ linePrefix(line.type) }}</span>
                <span class="line-content">{{ line.content || ' ' }}</span>
            </div>
            
            <!-- Show more indicator -->
            <div v-if="!isExpanded && hasMoreLines" class="diff-more">
                {{ t('streamingDiff.moreLines', { count: hiddenLinesCount }) }}
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed, toRef } from 'vue'
import { useStreamingDiff, type DiffLine } from '../composables/useStreamingDiff'

interface Props {
    filePath: string
    content: string
    originalContent?: string
}

const props = defineProps<Props>()
const { t } = useI18n()

const contentRef = toRef(props, 'content')
const filePathRef = toRef(props, 'filePath')
const originalContentRef = toRef(props, 'originalContent')

const { 
    parsedDiff, 
    previewLines, 
    stats, 
    operation, 
    isExpanded, 
    toggleExpand 
} = useStreamingDiff(contentRef, filePathRef, originalContentRef)

// File name from path
const fileName = computed(() => {
    const parts = props.filePath.split('/')
    return parts[parts.length - 1] || props.filePath
})

// Operation styling
const operationClass = computed(() => ({
    'op-create': operation.value === 'create',
    'op-modify': operation.value === 'modify',
    'op-delete': operation.value === 'delete'
}))

const operationIcon = computed(() => {
    switch (operation.value) {
        case 'create': return '+'
        case 'delete': return '-'
        default: return '~'
    }
})

// Lines to display
const displayLines = computed(() => {
    return isExpanded.value ? parsedDiff.value.lines : previewLines.value
})

const hasMoreLines = computed(() => {
    return parsedDiff.value.lines.length > previewLines.value.length
})

const hiddenLinesCount = computed(() => {
    return parsedDiff.value.lines.length - previewLines.value.length
})

// Line styling helpers
function lineClass(type: DiffLine['type']) {
    return {
        'line-add': type === 'add',
        'line-remove': type === 'remove',
        'line-context': type === 'context'
    }
}

function linePrefix(type: DiffLine['type']): string {
    switch (type) {
        case 'add': return '+'
        case 'remove': return '-'
        default: return ' '
    }
}

function formatLineNumber(line: DiffLine): string {
    if (line.type === 'add') {
        return line.lineNumber?.toString().padStart(4, ' ') || '    '
    }
    if (line.type === 'remove') {
        return line.oldLineNumber?.toString().padStart(4, ' ') || '    '
    }
    return line.lineNumber?.toString().padStart(4, ' ') || '    '
}
</script>

<style scoped>
.streaming-diff-preview {
    @apply bg-gray-900/60 border border-gray-700/50 rounded-lg overflow-hidden my-2;
}

.diff-header {
    @apply flex items-center justify-between px-3 py-2 cursor-pointer hover:bg-gray-800/50 transition-colors;
}

.diff-file-info {
    @apply flex items-center gap-2 min-w-0;
}

.diff-operation {
    @apply w-5 h-5 rounded flex items-center justify-center text-xs font-bold flex-shrink-0;
}

.op-create {
    @apply bg-green-500/20 text-green-400;
}

.op-modify {
    @apply bg-yellow-500/20 text-yellow-400;
}

.op-delete {
    @apply bg-red-500/20 text-red-400;
}

.diff-filename {
    @apply text-sm text-gray-300 truncate;
}

.diff-stats {
    @apply flex items-center gap-2 flex-shrink-0;
}

.stat-add {
    @apply text-xs font-mono text-green-400;
}

.stat-remove {
    @apply text-xs font-mono text-red-400;
}

.expand-icon {
    @apply w-4 h-4 text-gray-400 transition-transform;
}

.diff-content {
    @apply border-t border-gray-700/50 font-mono text-xs overflow-hidden;
    max-height: 12rem;
    overflow-y: auto;
}

.diff-expanded {
    max-height: 50vh;
}

.diff-line {
    @apply flex;
    line-height: 1.5;
}

.line-number {
    @apply text-gray-500 select-none px-2 text-right flex-shrink-0;
    min-width: 3rem;
}

.line-prefix {
    @apply select-none flex-shrink-0 w-4 text-center;
}

.line-content {
    @apply flex-1 pr-2 whitespace-pre;
    overflow-x: auto;
}

.line-add {
    @apply bg-green-500/10;
}

.line-add .line-prefix,
.line-add .line-content {
    @apply text-green-400;
}

.line-remove {
    @apply bg-red-500/10;
}

.line-remove .line-prefix,
.line-remove .line-content {
    @apply text-red-400;
}

.line-context {
    @apply text-gray-400;
}

.diff-more {
    @apply text-center py-2 text-xs text-gray-500 bg-gray-800/30 border-t border-gray-700/50;
}
</style>
