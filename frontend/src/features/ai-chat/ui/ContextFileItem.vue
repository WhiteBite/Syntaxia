<template>
  <div
    :class="[
      'context-file-item',
      { 'context-file-item--selected': file.selected }
    ]"
    @click="$emit('toggle')"
    @keydown.enter="$emit('toggle')"
    @keydown.space.prevent="$emit('toggle')"
    tabindex="0"
    role="checkbox"
    :aria-checked="file.selected"
  >
    <!-- Checkbox -->
    <div class="context-file-checkbox">
      <div :class="['checkbox-box', { 'checkbox-box--checked': file.selected }]">
        <svg v-if="file.selected" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
        </svg>
      </div>
    </div>

    <!-- File Icon -->
    <div class="context-file-icon">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
          d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
    </div>

    <!-- File Info -->
    <div class="context-file-info">
      <div class="context-file-name">{{ fileName }}</div>
      <div class="context-file-path">{{ filePath }}</div>
    </div>

    <!-- Relevance Badge -->
    <div 
      :class="['context-file-relevance', relevanceClass]"
      :title="file.reason"
    >
      {{ relevancePercent }}%
    </div>

    <!-- Token Count -->
    <div class="context-file-tokens">
      {{ formattedTokens }}
    </div>

    <!-- Reason Tooltip Icon -->
    <div 
      class="context-file-reason"
      :title="file.reason"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ContextFilePreview } from '../composables/useContextPreview'

interface Props {
  file: ContextFilePreview
}

const props = defineProps<Props>()

defineEmits<{
  toggle: []
}>()

const fileName = computed(() => {
  const parts = props.file.path.split('/')
  return parts[parts.length - 1] || props.file.path
})

const filePath = computed(() => {
  const parts = props.file.path.split('/')
  if (parts.length > 1) {
    return parts.slice(0, -1).join('/')
  }
  return ''
})

const relevancePercent = computed(() => 
  Math.round(props.file.relevance * 100)
)

const relevanceClass = computed(() => {
  const percent = relevancePercent.value
  if (percent >= 90) return 'relevance-high'
  if (percent >= 70) return 'relevance-medium'
  return 'relevance-low'
})

const formattedTokens = computed(() => {
  const tokens = props.file.tokens
  if (tokens >= 1000) {
    return `${(tokens / 1000).toFixed(1)}k`
  }
  return String(tokens)
})
</script>

<style scoped>
.context-file-item {
  @apply flex items-center gap-3 px-4 py-3;
  @apply cursor-pointer transition-all duration-150;
  @apply border-b;
  border-color: var(--border-subtle);
  background: var(--bg-1);
}

.context-file-item:hover {
  background: var(--bg-2);
}

.context-file-item:focus {
  @apply outline-none;
  background: var(--bg-2);
  box-shadow: inset 0 0 0 2px var(--accent-primary);
}

.context-file-item--selected {
  background: var(--accent-purple-bg);
}

.context-file-item--selected:hover {
  background: rgba(139, 92, 246, 0.2);
}

.context-file-checkbox {
  @apply flex-shrink-0;
}

.checkbox-box {
  @apply w-5 h-5 rounded-md flex items-center justify-center;
  @apply transition-all duration-150;
  background: var(--bg-2);
  border: 2px solid var(--border-strong);
}

.checkbox-box--checked {
  background: var(--accent-primary);
  border-color: var(--accent-primary);
  color: white;
}

.context-file-icon {
  @apply flex-shrink-0;
  color: var(--text-muted);
}

.context-file-info {
  @apply flex-1 min-w-0;
}

.context-file-name {
  @apply text-sm font-medium truncate;
  color: var(--text-primary);
}

.context-file-path {
  @apply text-xs truncate;
  color: var(--text-muted);
}

.context-file-relevance {
  @apply flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-semibold;
  @apply cursor-help;
}

.relevance-high {
  background: var(--color-success-soft);
  color: var(--color-success);
}

.relevance-medium {
  background: var(--color-warning-soft);
  color: var(--color-warning);
}

.relevance-low {
  background: var(--bg-3);
  color: var(--text-secondary);
}

.context-file-tokens {
  @apply flex-shrink-0 text-xs font-mono px-2 py-0.5 rounded;
  background: var(--bg-2);
  color: var(--text-secondary);
}

.context-file-reason {
  @apply flex-shrink-0 cursor-help;
  color: var(--text-muted);
  @apply transition-colors duration-150;
}

.context-file-reason:hover {
  color: var(--text-secondary);
}
</style>
