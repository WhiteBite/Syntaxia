<template>
  <div class="context-used-details" role="region" :aria-label="t('contextUsed.title')">
    <!-- File list -->
    <div class="context-used-details__list">
      <div
        v-for="file in displayedFiles"
        :key="`${file.path}-${file.operation}`"
        class="context-used-details__item"
      >
        <span class="context-used-details__icon" :title="getOperationLabel(file.operation)">
          {{ getOperationIcon(file.operation) }}
        </span>
        <span class="context-used-details__path" :title="file.path">
          {{ file.path }}
        </span>
        <span class="context-used-details__operation">
          ({{ getOperationLabel(file.operation) }})
        </span>
        <span v-if="file.tokens" class="context-used-details__tokens">
          {{ formatTokens(file.tokens) }}
        </span>
      </div>

      <!-- Show more/less -->
      <button
        v-if="context.files.length > maxVisibleFiles"
        class="context-used-details__toggle"
        @click="toggleShowAll"
      >
        {{ showAll ? t('contextUsed.showLess') : t('contextUsed.showMore', { count: hiddenCount }) }}
      </button>
    </div>

    <!-- Actions -->
    <div class="context-used-details__actions">
      <button class="context-used-details__action" @click="addMoreContext">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {{ t('contextUsed.addMore') }}
      </button>
      <button class="context-used-details__action" @click="copyFileList">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
        {{ copied ? t('contextUsed.copied') : t('contextUsed.copyList') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useUIStore } from '@/stores/ui.store'
import { computed, onUnmounted, ref } from 'vue'
import type { ContextOperation, ContextUsed } from './types'

interface Props {
  context: ContextUsed
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'addContext'): void
}>()

const { t } = useI18n()
const uiStore = useUIStore()

const maxVisibleFiles = 5
const showAll = ref(false)
const copied = ref(false)
let copyTimeoutId: ReturnType<typeof setTimeout> | null = null

onUnmounted(() => {
  if (copyTimeoutId) {
    clearTimeout(copyTimeoutId)
    copyTimeoutId = null
  }
})

const displayedFiles = computed(() => {
  if (showAll.value) {
    return props.context.files
  }
  return props.context.files.slice(0, maxVisibleFiles)
})

const hiddenCount = computed(() => {
  return Math.max(0, props.context.files.length - maxVisibleFiles)
})

function toggleShowAll() {
  showAll.value = !showAll.value
}

function getOperationIcon(operation: ContextOperation): string {
  const icons: Record<ContextOperation, string> = {
    read: '📖',
    write: '✏️',
    search: '🔍',
    created: '✨',
    deleted: '🗑️'
  }
  return icons[operation] || '📄'
}

function getOperationLabel(operation: ContextOperation): string {
  return t(`contextUsed.operation.${operation}`)
}

function formatTokens(tokens: number): string {
  if (tokens >= 1000) {
    return `${(tokens / 1000).toFixed(1)}k`
  }
  return String(tokens)
}

function addMoreContext() {
  emit('addContext')
}

async function copyFileList() {
  if (copyTimeoutId) {
    clearTimeout(copyTimeoutId)
    copyTimeoutId = null
  }

  const fileList = props.context.files
    .map(f => `${f.path} (${f.operation})`)
    .join('\n')

  try {
    await navigator.clipboard.writeText(fileList)
    copied.value = true
    uiStore.addToast(t('contextUsed.copied'), 'success')
    copyTimeoutId = setTimeout(() => {
      copied.value = false
      copyTimeoutId = null
    }, 2000)
  } catch {
    uiStore.addToast(t('contextUsed.copyFailed'), 'error')
  }
}
</script>

<style scoped>
.context-used-details {
  margin-top: 0.5rem;
  padding: 0.75rem;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: 0.5rem;
  animation: slideDown 200ms ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.context-used-details__list {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.context-used-details__item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.375rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.8125rem;
  transition: background 150ms ease-out;
}

.context-used-details__item:hover {
  background: var(--bg-2);
}

.context-used-details__icon {
  flex-shrink: 0;
  font-size: 0.875rem;
}

.context-used-details__path {
  flex: 1;
  min-width: 0;
  color: var(--text-primary);
  font-family: var(--font-mono, monospace);
  font-size: 0.75rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.context-used-details__operation {
  flex-shrink: 0;
  color: var(--text-muted);
  font-size: 0.6875rem;
}

.context-used-details__tokens {
  flex-shrink: 0;
  padding: 0.125rem 0.375rem;
  background: var(--bg-3);
  border-radius: 9999px;
  color: var(--text-secondary);
  font-size: 0.625rem;
  font-weight: 500;
}

.context-used-details__toggle {
  margin-top: 0.25rem;
  padding: 0.25rem 0.5rem;
  background: transparent;
  border: none;
  color: var(--accent-purple);
  font-size: 0.75rem;
  cursor: pointer;
  transition: color 150ms ease-out;
}

.context-used-details__toggle:hover {
  color: var(--accent-purple-hover);
  text-decoration: underline;
}

.context-used-details__actions {
  display: flex;
  gap: 0.75rem;
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--border-subtle);
}

.context-used-details__action {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.625rem;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  border-radius: 0.375rem;
  color: var(--text-secondary);
  font-size: 0.75rem;
  cursor: pointer;
  transition: all 150ms ease-out;
}

.context-used-details__action:hover {
  background: var(--bg-3);
  border-color: var(--border-strong);
  color: var(--text-primary);
}
</style>
