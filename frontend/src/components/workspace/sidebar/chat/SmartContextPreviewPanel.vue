<template>
  <div class="smart-context-panel">
    <div class="smart-context-panel__header">
      <div class="smart-context-panel__title">
        <svg class="w-4 h-4 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
        </svg>
        <span>{{ t('chat.suggestedFiles') }}</span>
      </div>
      <span class="smart-context-panel__stats">
        ~{{ formattedTokens }} tokens
      </span>
    </div>
    
    <!-- File list with reasons -->
    <div class="smart-context-panel__files">
      <div 
        v-for="file in displayedFiles" 
        :key="file.path"
        class="suggested-file"
      >
        <div class="suggested-file__main">
          <input 
            type="checkbox" 
            :checked="selectedFiles.has(file.path)"
            @change="toggleFile(file.path)"
            class="suggested-file__checkbox"
          />
          <span class="suggested-file__name" :title="file.path">
            {{ getFileName(file.path) }}
          </span>
          <span class="suggested-file__relevance" :style="{ opacity: file.relevance }">
            {{ Math.round(file.relevance * 100) }}%
          </span>
        </div>
        <div v-if="file.reason" class="suggested-file__reason">
          {{ truncateReason(file.reason) }}
        </div>
      </div>
      
      <button 
        v-if="hasMoreFiles"
        class="smart-context-panel__more"
        @click="showAll = !showAll"
      >
        {{ showAll ? t('common.showLess') : `+${hiddenCount} ${t('common.more')}` }}
      </button>
    </div>
    
    <!-- Actions -->
    <div class="smart-context-panel__actions">
      <button 
        class="btn btn-sm btn-primary flex-1"
        @click="handleConfirm"
        :disabled="selectedFiles.size === 0"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        {{ t('chat.useSelectedFiles', { count: selectedFiles.size }) }}
      </button>
      <button 
        class="btn btn-sm btn-ghost"
        @click="$emit('cancel')"
      >
        {{ t('common.cancel') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed, ref, watch } from 'vue'

interface SuggestedFile {
  path: string
  reason: string
  relevance: number
}

interface SmartContextPreview {
  files: SuggestedFile[]
  totalTokens: number
  query: string
}

const props = defineProps<{
  preview: SmartContextPreview
}>()

const emit = defineEmits<{
  confirm: [files: string[]]
  cancel: []
}>()

const { t } = useI18n()

const MAX_VISIBLE = 5
const MAX_REASON_LENGTH = 80

const showAll = ref(false)
const selectedFiles = ref<Set<string>>(new Set())

// Initialize all files as selected
watch(() => props.preview.files, (files) => {
  selectedFiles.value = new Set(files.map(f => f.path))
}, { immediate: true })

const displayedFiles = computed(() => {
  if (showAll.value) return props.preview.files
  return props.preview.files.slice(0, MAX_VISIBLE)
})

const hasMoreFiles = computed(() => props.preview.files.length > MAX_VISIBLE)
const hiddenCount = computed(() => props.preview.files.length - MAX_VISIBLE)

const formattedTokens = computed(() => {
  const tokens = props.preview.totalTokens
  if (tokens >= 1000) return `${(tokens / 1000).toFixed(1)}k`
  return tokens.toString()
})

function getFileName(path: string): string {
  return path.split('/').pop() || path
}

function truncateReason(reason: string): string {
  if (reason.length <= MAX_REASON_LENGTH) return reason
  return reason.slice(0, MAX_REASON_LENGTH) + '...'
}

function toggleFile(path: string) {
  if (selectedFiles.value.has(path)) {
    selectedFiles.value.delete(path)
  } else {
    selectedFiles.value.add(path)
  }
  selectedFiles.value = new Set(selectedFiles.value)
}

function handleConfirm() {
  emit('confirm', Array.from(selectedFiles.value))
}
</script>

<style scoped>
.smart-context-panel {
  margin: 0 12px 8px;
  padding: 12px;
  background: rgba(139, 92, 246, 0.08);
  border: 1px solid rgba(139, 92, 246, 0.25);
  border-radius: 12px;
}

.smart-context-panel__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.smart-context-panel__title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 500;
  color: var(--color-purple-300, #c4b5fd);
}

.smart-context-panel__stats {
  font-size: 10px;
  color: var(--color-gray-400, #9ca3af);
}

.smart-context-panel__files {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 180px;
  overflow-y: auto;
  margin-bottom: 10px;
}

.suggested-file {
  padding: 6px 8px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 6px;
}

.suggested-file__main {
  display: flex;
  align-items: center;
  gap: 8px;
}

.suggested-file__checkbox {
  width: 14px;
  height: 14px;
  border-radius: 3px;
  border: 1px solid rgba(139, 92, 246, 0.4);
  background: transparent;
  cursor: pointer;
  accent-color: var(--color-purple-500, #8b5cf6);
}

.suggested-file__name {
  flex: 1;
  font-size: 11px;
  color: var(--color-gray-200, #e5e7eb);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.suggested-file__relevance {
  font-size: 10px;
  color: var(--color-purple-400, #a78bfa);
  font-weight: 500;
}

.suggested-file__reason {
  margin-top: 4px;
  padding-left: 22px;
  font-size: 10px;
  color: var(--color-gray-500, #6b7280);
  line-height: 1.3;
}

.smart-context-panel__more {
  padding: 4px 8px;
  font-size: 10px;
  color: var(--color-gray-400, #9ca3af);
  background: transparent;
  border: none;
  cursor: pointer;
  text-align: left;
}

.smart-context-panel__more:hover {
  color: var(--color-gray-300, #d1d5db);
}

.smart-context-panel__actions {
  display: flex;
  gap: 8px;
}
</style>
