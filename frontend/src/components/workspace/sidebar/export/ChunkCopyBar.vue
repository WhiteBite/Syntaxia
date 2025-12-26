<template>
  <div v-if="contextStore.hasContext && needsSplitting" class="chunk-copy-bar">
    <!-- Header with token info -->
    <div class="chunk-header">
      <div class="chunk-header-info">
        <AlertTriangle class="w-3.5 h-3.5 text-amber-400" />
        <span class="chunk-header-text">
          {{ formatTokens(totalTokens) }} → {{ chunkCount }} {{ t('chunks.parts') }}
        </span>
      </div>
      <button 
        v-if="progress > 0" 
        class="chunk-reset-btn"
        :title="t('chunks.reset')"
        @click="handleReset"
      >
        <RotateCcw class="w-3 h-3" />
      </button>
    </div>

    <!-- Main action: Copy next or All done -->
    <button 
      v-if="!allCopied" 
      class="chunk-main-btn"
      :disabled="isProcessing"
      @click="handleCopyNext"
    >
      <Copy class="w-4 h-4" />
      <span>{{ t('chunks.copyPart') }} {{ nextUncopiedIndex + 1 }}/{{ chunkCount }}</span>
    </button>
    <div v-else class="chunk-done">
      <CheckCircle class="w-4 h-4" />
      <span>{{ t('chunks.allCopied') }}</span>
    </div>

    <!-- Chunk grid -->
    <div class="chunk-grid">
      <button
        v-for="chunk in chunks"
        :key="chunk.index"
        class="chunk-item"
        :class="{
          'chunk-item--copied': chunk.copied,
          'chunk-item--next': !chunk.copied && chunk.index === nextUncopiedIndex,
          'chunk-item--pending': !chunk.copied && chunk.index !== nextUncopiedIndex
        }"
        :title="getChunkTooltip(chunk)"
        @click="handleCopyChunk(chunk.index)"
      >
        <span class="chunk-item-num">
          <Check v-if="chunk.copied" class="w-3 h-3" />
          <span v-else>{{ chunk.index + 1 }}</span>
        </span>
        <span class="chunk-item-tokens">{{ formatTokens(chunk.tokens) }}</span>
      </button>
    </div>

    <!-- Progress indicator -->
    <div v-if="progress > 0" class="chunk-progress">
      <div class="chunk-progress-bar" :style="{ width: `${progress}%` }" />
      <span class="chunk-progress-text">{{ copiedCount }}/{{ chunkCount }}</span>
    </div>
  </div>

  <!-- Generate button when no chunks yet -->
  <button
    v-else-if="contextStore.hasContext && enableAutoSplit"
    class="chunk-generate-btn"
    :disabled="isProcessing"
    @click="generateChunks"
  >
    <Layers class="w-4 h-4" />
    <span>{{ t('chunks.generate') }}</span>
  </button>
</template>

<script setup lang="ts">
import { useChunkedCopy, type Chunk } from '@/composables/useChunkedCopy'
import { useI18n } from '@/composables/useI18n'
import { useContextStore } from '@/features/context/model/context.store'
import { AlertTriangle, Check, CheckCircle, Copy, Layers, RotateCcw } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'

const { t } = useI18n()
const contextStore = useContextStore()
const showCopiedFeedback = ref(false)

const {
  chunks,
  isProcessing,
  totalTokens,
  needsSplitting,
  chunkCount,
  allCopied,
  progress,
  enableAutoSplit,
  generateChunks,
  copyChunk,
  copyNext,
  resetCopied
} = useChunkedCopy()

const nextUncopiedIndex = computed(() => {
  const idx = chunks.value.findIndex(c => !c.copied)
  return idx === -1 ? 0 : idx
})

const copiedCount = computed(() => chunks.value.filter(c => c.copied).length)

function formatTokens(tokens: number): string {
  if (tokens >= 1000) return `${(tokens / 1000).toFixed(1)}K`
  return tokens.toString()
}

function getChunkTooltip(chunk: Chunk): string {
  const status = chunk.copied ? `✓ ${t('chunks.copied')}` : t('chunks.clickToCopy')
  return `${t('chunks.part')} ${chunk.index + 1}: ${formatTokens(chunk.tokens)} ${t('context.tokens')}\n${chunk.files} ${t('context.files')}\n${status}`
}

async function handleCopyChunk(index: number) {
  const success = await copyChunk(index)
  if (success) {
    showFeedback()
  }
}

async function handleCopyNext() {
  const success = await copyNext()
  if (success) {
    showFeedback()
  }
}

function handleReset() {
  resetCopied()
}

function showFeedback() {
  showCopiedFeedback.value = true
  setTimeout(() => {
    showCopiedFeedback.value = false
  }, 1500)
}

// Auto-generate chunks when context changes and auto-split is enabled
watch(
  () => [contextStore.contextId, enableAutoSplit.value],
  async ([ctxId, autoSplit]) => {
    if (ctxId && autoSplit && needsSplitting.value) {
      await generateChunks()
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.chunk-copy-bar {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 0.75rem;
  background: var(--bg-2);
  border-radius: 0.5rem;
  border: 1px solid var(--border-default);
}

.chunk-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
}

.chunk-header-info {
  display: flex;
  align-items: center;
  gap: 0.375rem;
}

.chunk-header-text {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.chunk-reset-btn {
  padding: 0.25rem;
  border-radius: 0.25rem;
  color: var(--text-muted);
  transition: all 150ms;
}

.chunk-reset-btn:hover {
  background: var(--bg-3);
  color: var(--text-primary);
}

.chunk-main-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: var(--color-primary);
  color: white;
  border-radius: 0.375rem;
  font-size: 0.8125rem;
  font-weight: 500;
  transition: all 150ms;
}

.chunk-main-btn:hover:not(:disabled) {
  filter: brightness(1.1);
}

.chunk-main-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.chunk-done {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
  border-radius: 0.375rem;
  font-size: 0.8125rem;
  font-weight: 500;
}

.chunk-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 0.375rem;
}

.chunk-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-width: 3rem;
  padding: 0.375rem 0.5rem;
  border-radius: 0.375rem;
  border: 1px solid var(--border-default);
  background: var(--bg-1);
  cursor: pointer;
  transition: all 150ms;
}

.chunk-item:hover {
  border-color: var(--color-primary);
  background: var(--bg-2);
}

.chunk-item--copied {
  background: rgba(34, 197, 94, 0.1);
  border-color: rgba(34, 197, 94, 0.3);
  color: #22c55e;
}

.chunk-item--next {
  background: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.4);
  box-shadow: 0 0 0 1px rgba(99, 102, 241, 0.2);
}

.chunk-item--pending {
  opacity: 0.6;
}

.chunk-item--pending:hover {
  opacity: 1;
}

.chunk-item-num {
  font-size: 0.75rem;
  font-weight: 600;
  line-height: 1;
}

.chunk-item-tokens {
  font-size: 0.625rem;
  color: var(--text-muted);
  margin-top: 0.125rem;
}

.chunk-progress {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  height: 1rem;
}

.chunk-progress-bar {
  flex: 1;
  height: 0.25rem;
  background: var(--color-primary);
  border-radius: 0.125rem;
  transition: width 300ms ease-out;
}

.chunk-progress-text {
  font-size: 0.625rem;
  color: var(--text-muted);
  min-width: 2rem;
  text-align: right;
}

.chunk-generate-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.5rem 0.75rem;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  border-radius: 0.375rem;
  color: var(--text-secondary);
  font-size: 0.8125rem;
  transition: all 150ms;
}

.chunk-generate-btn:hover:not(:disabled) {
  background: var(--bg-3);
  border-color: var(--color-primary);
  color: var(--text-primary);
}

.chunk-generate-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
