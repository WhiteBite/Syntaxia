<template>
  <button
    class="context-used-badge"
    :class="{ 'context-used-badge--expanded': isExpanded }"
    :aria-expanded="isExpanded"
    role="button"
    @click="toggleExpand"
  >
    <span class="context-used-badge__icon">📎</span>
    <span class="context-used-badge__text">
      {{ context.files.length }} {{ t('contextUsed.files') }}, {{ formatTokens(context.totalTokens) }}
    </span>
    <svg
      class="context-used-badge__chevron"
      :class="{ 'context-used-badge__chevron--rotated': isExpanded }"
      width="16"
      height="16"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
    >
      <polyline points="6 9 12 15 18 9" />
    </svg>
  </button>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { ref } from 'vue'
import type { ContextUsed } from './types'

interface Props {
  context: ContextUsed
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'toggle', expanded: boolean): void
}>()

const { t } = useI18n()
const isExpanded = ref(false)

function toggleExpand() {
  isExpanded.value = !isExpanded.value
  emit('toggle', isExpanded.value)
}

function formatTokens(tokens: number): string {
  if (tokens >= 1000) {
    return `${Math.round(tokens / 1000)}k tokens`
  }
  return `${tokens} tokens`
}
</script>

<style scoped>
.context-used-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.375rem 0.75rem;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  border-radius: 9999px;
  color: var(--text-secondary);
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 150ms ease-out;
}

.context-used-badge:hover {
  background: var(--bg-3);
  border-color: var(--border-strong);
  color: var(--text-primary);
}

.context-used-badge--expanded {
  background: var(--accent-purple-bg);
  border-color: var(--accent-purple-border);
  color: white;
}

.context-used-badge--expanded:hover {
  background: var(--accent-purple-bg);
  opacity: 0.9;
}

.context-used-badge__icon {
  font-size: 0.875rem;
}

.context-used-badge__text {
  white-space: nowrap;
}

.context-used-badge__chevron {
  width: 14px;
  height: 14px;
  transition: transform 200ms ease-out;
}

.context-used-badge__chevron--rotated {
  transform: rotate(180deg);
}
</style>
