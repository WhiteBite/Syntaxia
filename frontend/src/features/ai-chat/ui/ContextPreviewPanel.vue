<template>
  <BaseModal
    :model-value="visible"
    size="lg"
    :close-on-backdrop="true"
    :close-on-esc="true"
    @close="handleClose"
  >
    <template #header>
      <div class="flex items-center gap-3">
        <div class="section-icon section-icon-purple">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
          </svg>
        </div>
        <div>
          <h2 class="text-lg font-semibold text-white">
            {{ t('context.preview.title') }}
          </h2>
          <p class="text-sm text-gray-400">
            {{ t('context.preview.found', { count: files.length, tokens: formattedTotalTokens }) }}
          </p>
        </div>
      </div>
    </template>

          <!-- System Prompt Section -->
          <div class="system-prompt-section">
            <div 
              class="system-prompt-header"
              @click="contextPreview.togglePromptExpanded()"
              role="button"
              tabindex="0"
              @keydown.enter="contextPreview.togglePromptExpanded()"
              @keydown.space.prevent="contextPreview.togglePromptExpanded()"
            >
              <div class="system-prompt-info">
                <span class="system-prompt-icon">{{ contextPreview.systemPrompt.value.icon }}</span>
                <span class="system-prompt-name">{{ contextPreview.systemPrompt.value.name }}</span>
                <span class="system-prompt-tokens">~{{ formattedPromptTokens }} {{ t('context.tokens') }}</span>
                <span 
                  v-if="contextPreview.isPromptOverBudget.value" 
                  class="prompt-budget-warning"
                  :title="t('context.preview.promptBudgetWarning')"
                >
                  ⚠️
                </span>
              </div>
              <div class="system-prompt-actions">
                <button 
                  @click.stop="handleEditPrompt"
                  class="prompt-edit-btn"
                  :title="t('context.preview.editPrompt')"
                >
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <svg 
                  class="chevron-icon" 
                  :class="{ 'chevron-expanded': contextPreview.showPrompt.value }"
                  fill="none" 
                  stroke="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>
            <Transition name="expand">
              <div v-if="contextPreview.showPrompt.value" class="system-prompt-content">
                <pre class="system-prompt-text">{{ contextPreview.systemPrompt.value.content || t('context.preview.noPromptContent') }}</pre>
              </div>
            </Transition>
          </div>

          <!-- Actions Bar -->
          <div class="context-preview-actions">
            <button @click="contextPreview.selectAll()" class="btn btn-ghost btn-sm">
              {{ t('context.preview.selectAll') }}
            </button>
            <button @click="contextPreview.deselectAll()" class="btn btn-ghost btn-sm">
              {{ t('context.preview.deselectAll') }}
            </button>
          </div>

          <!-- File List -->
          <div class="context-preview-list">
            <ContextFileItem
              v-for="file in contextPreview.files.value"
              :key="file.path"
              :file="file"
              @toggle="contextPreview.toggleFile(file.path)"
            />
            
            <!-- Empty State -->
            <div v-if="files.length === 0" class="empty-state py-12">
              <p class="text-gray-400">{{ t('context.noResults') }}</p>
            </div>
          </div>

          <!-- Token Budget Bar -->
          <div class="context-preview-budget">
            <div class="budget-info">
              <span class="budget-label">
                {{ t('context.preview.selected', { 
                  count: contextPreview.selectedCount.value, 
                  tokens: formattedSelectedTokens 
                }) }}
              </span>
            </div>
            
            <!-- Token Breakdown -->
            <div class="token-breakdown">
              <span class="token-breakdown-item">
                {{ t('context.preview.promptTokens') }}: {{ formattedPromptTokens }}
              </span>
              <span class="token-breakdown-separator">+</span>
              <span class="token-breakdown-item">
                {{ t('context.preview.fileTokens') }}: {{ formattedSelectedTokens }}
              </span>
              <span class="token-breakdown-separator">=</span>
              <span class="token-breakdown-total" :class="{ 'token-over-limit': contextPreview.isOverLimit.value }">
                {{ formattedTotalWithPrompt }}
              </span>
            </div>
            
            <div class="budget-bar-container">
              <div 
                class="budget-bar"
                :class="{ 'budget-bar--over': contextPreview.isOverLimit.value }"
                :style="{ width: `${contextPreview.tokenUsagePercent.value}%` }"
              ></div>
            </div>
            <div class="budget-limits">
              <span>{{ formattedTotalWithPrompt }}</span>
              <span>/</span>
              <span>{{ formattedTokenLimit }}</span>
            </div>
          </div>

    <template #footer>
      <button @click="handleClose" class="btn btn-secondary">
        {{ t('context.preview.cancel') }}
      </button>
      <button 
        @click="handleConfirm" 
        :disabled="!contextPreview.canSend.value"
        class="btn btn-primary"
      >
        {{ t('context.preview.sendToAI') }}
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
        </svg>
      </button>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed, onMounted, onUnmounted, watch } from 'vue'
import { useContextPreview, type ContextFilePreview } from '../composables/useContextPreview'
import ContextFileItem from './ContextFileItem.vue'

interface Props {
  visible: boolean
  files: ContextFilePreview[]
  tokenLimit: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  confirm: [selectedFiles: string[]]
}>()

const { t } = useI18n()
const contextPreview = useContextPreview()

// Sync props to composable
watch(() => props.files, (newFiles) => {
  contextPreview.setFiles(newFiles)
}, { immediate: true })

watch(() => props.tokenLimit, (limit) => {
  contextPreview.setTokenLimit(limit)
}, { immediate: true })

// Formatted values
const formattedTotalTokens = computed(() => formatTokens(contextPreview.totalTokens.value))
const formattedSelectedTokens = computed(() => formatTokens(contextPreview.selectedTokens.value))
const formattedTokenLimit = computed(() => formatTokens(props.tokenLimit))
const formattedPromptTokens = computed(() => formatTokens(contextPreview.promptTokens.value))
const formattedTotalWithPrompt = computed(() => formatTokens(contextPreview.totalWithPrompt.value))

function formatTokens(tokens: number): string {
  if (tokens >= 1000) {
    return `${(tokens / 1000).toFixed(1)}k`
  }
  return String(tokens)
}

// Handlers
function handleClose() {
  emit('close')
}

function handleConfirm() {
  if (contextPreview.canSend.value) {
    emit('confirm', contextPreview.getSelectedPaths())
  }
}

function handleEditPrompt() {
  contextPreview.openTemplateSettings()
}

// Keyboard handling
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.visible) {
    handleClose()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  contextPreview.reset()
})
</script>

<style scoped>
.context-preview-overlay {
  @apply fixed inset-0 flex items-center justify-center p-4;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
  z-index: var(--z-modal-backdrop, 50);
}

.context-preview-modal {
  @apply flex flex-col rounded-2xl overflow-hidden;
  width: min(700px, 90vw);
  max-height: 85vh;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  box-shadow: var(--shadow-xl);
}

.context-preview-header {
  @apply flex items-center justify-between p-4;
  background: var(--bg-2);
  border-bottom: 1px solid var(--border-default);
}

.context-preview-actions {
  @apply flex items-center gap-2 px-4 py-2;
  background: var(--bg-1);
  border-bottom: 1px solid var(--border-subtle);
}

/* System Prompt Section */
.system-prompt-section {
  background: var(--bg-2);
  border-bottom: 1px solid var(--border-default);
}

.system-prompt-header {
  @apply flex items-center justify-between px-4 py-3 cursor-pointer;
  transition: background 0.15s ease;
}

.system-prompt-header:hover {
  background: var(--bg-3);
}

.system-prompt-header:focus {
  outline: 2px solid var(--color-primary);
  outline-offset: -2px;
}

.system-prompt-info {
  @apply flex items-center gap-2;
}

.system-prompt-actions {
  @apply flex items-center gap-2;
}

.prompt-edit-btn {
  @apply p-1.5 rounded-md transition-all duration-150;
  color: var(--text-muted);
  background: transparent;
  border: none;
  cursor: pointer;
}

.prompt-edit-btn:hover {
  color: var(--text-primary);
  background: var(--bg-3);
}

.prompt-budget-warning {
  @apply text-xs ml-1;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.system-prompt-icon {
  @apply text-base;
}

.system-prompt-name {
  @apply font-medium text-sm;
  color: var(--text-primary);
}

.system-prompt-tokens {
  @apply text-xs px-2 py-0.5 rounded-full;
  background: var(--bg-3);
  color: var(--text-muted);
}

.chevron-icon {
  @apply w-4 h-4 transition-transform duration-200;
  color: var(--text-muted);
}

.chevron-expanded {
  transform: rotate(180deg);
}

.system-prompt-content {
  @apply px-4 pb-3;
  overflow: hidden;
}

.system-prompt-text {
  @apply text-xs p-3 rounded-lg overflow-auto whitespace-pre-wrap;
  background: var(--bg-1);
  color: var(--text-secondary);
  max-height: 25vh;
  font-family: var(--font-mono, monospace);
  border: 1px solid var(--border-subtle);
}

/* Expand transition */
.expand-enter-active,
.expand-leave-active {
  transition: all 0.2s ease;
  overflow: hidden;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
  padding-top: 0;
  padding-bottom: 0;
}

.context-preview-list {
  @apply flex-1 overflow-y-auto;
  min-height: 0;
  max-height: 50vh;
}

.context-preview-budget {
  @apply px-4 py-3;
  background: var(--bg-2);
  border-top: 1px solid var(--border-default);
}

.budget-info {
  @apply flex items-center justify-between mb-2;
}

.budget-label {
  @apply text-sm;
  color: var(--text-secondary);
}

/* Token Breakdown */
.token-breakdown {
  @apply flex items-center gap-2 text-xs mb-2 flex-wrap;
  color: var(--text-muted);
}

.token-breakdown-item {
  @apply px-2 py-0.5 rounded;
  background: var(--bg-3);
}

.token-breakdown-separator {
  color: var(--text-muted);
}

.token-breakdown-total {
  @apply px-2 py-0.5 rounded font-medium;
  background: var(--accent-indigo-bg);
  color: var(--accent-indigo);
}

.token-over-limit {
  background: rgba(239, 68, 68, 0.15);
  color: var(--color-danger);
}

.budget-bar-container {
  @apply w-full h-2 rounded-full overflow-hidden;
  background: var(--bg-3);
}

.budget-bar {
  @apply h-full rounded-full transition-all duration-300;
  background: var(--gradient-primary);
}

.budget-bar--over {
  background: var(--color-danger);
}

.budget-limits {
  @apply flex items-center justify-end gap-1 mt-1 text-xs font-mono;
  color: var(--text-muted);
}

.context-preview-footer {
  @apply flex items-center justify-end gap-3 p-4;
  background: var(--bg-2);
  border-top: 1px solid var(--border-default);
}

/* Modal Transition */
.modal-enter-active,
.modal-leave-active {
  transition: all 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .context-preview-modal,
.modal-leave-to .context-preview-modal {
  transform: scale(0.95) translateY(-10px);
}

.modal-enter-active .context-preview-modal,
.modal-leave-active .context-preview-modal {
  transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
</style>
