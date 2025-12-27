<template>
  <div 
    class="h-full flex flex-col bg-gray-900"
    @dragenter="handleDragEnter"
    @dragover="handleDragOver"
    @dragleave="handleDragLeave"
    @drop="handleDrop"
  >
    <!-- Drop Zone Overlay -->
    <Transition name="fade">
      <div 
        v-if="isDragOver" 
        class="chat-drop-overlay"
      >
        <div class="chat-drop-zone">
          <svg class="w-12 h-12 text-purple-400 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
          <p class="text-lg font-medium text-white">{{ t('chat.dragDrop.dropHere') }}</p>
          <p class="text-sm text-gray-400">{{ t('chat.dragDrop.addToContext') }}</p>
        </div>
      </div>
    </Transition>

    <!-- Header -->
    <div class="section-header">
      <div class="section-title">
        <div class="section-icon section-icon-purple">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
          </svg>
        </div>
        <span class="section-title-text">{{ t('chat.title') }}</span>
        <span v-if="sandboxStore.hasChanges" class="badge badge-warning">
          {{ sandboxStore.changeCount }} {{ t('chat.pendingChanges') }}
        </span>
      </div>
      
      <div class="flex items-center gap-2">
        <ChatHistoryPanel />
        <button
          v-if="sandboxStore.hasChanges"
          @click="showChangesPanel = true"
          class="btn btn-primary btn-xs"
        >
          {{ t('chat.reviewChanges') }}
        </button>
        <button
          @click="chatStore.clearChat"
          :disabled="!chatStore.hasMessages"
          class="btn btn-secondary btn-xs"
        >
          {{ t('chat.clear') }}
        </button>
      </div>
    </div>

    <!-- Mode Selector & Token Budget Bar -->
    <div class="chat-toolbar">
      <ChatModeSelector />
      <TokenBudgetBar />
    </div>

    <!-- Messages -->
    <div ref="messagesContainer" class="flex-1 overflow-auto p-4 space-y-4 min-h-0">
      <div v-if="!chatStore.hasMessages" class="empty-state h-full">
        <div class="empty-state-icon !w-20 !h-20 !rounded-2xl">
          <svg class="!w-10 !h-10 text-purple-500/50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
        </div>
        <p class="empty-state-title !text-lg">{{ t('chat.emptyTitle') }}</p>
        <p class="empty-state-text mb-4">{{ t('chat.emptyDesc') }}</p>
        <div class="info-box text-left text-xs space-y-2">
          <p class="font-semibold text-gray-300">{{ t('chat.plannedFeatures') }}</p>
          <ul class="list-disc list-inside text-gray-400 space-y-1">
            <li>{{ t('chat.feature.realtime') }}</li>
            <li>{{ t('chat.feature.contextAware') }}</li>
            <li>{{ t('chat.feature.codeGen') }}</li>
            <li>{{ t('chat.feature.streaming') }}</li>
            <li>{{ t('chat.feature.history') }}</li>
          </ul>
        </div>
      </div>

      <MessageItem
        v-for="message in chatStore.messages"
        :key="message.id"
        :message="message"
        @delete="chatStore.deleteMessage"
        @edit="chatStore.editMessage"
      />

      <!-- Typing indicator -->
      <div v-if="chatStore.isStreaming" class="flex items-center gap-2 text-gray-400">
        <div class="flex gap-1">
          <div class="w-2 h-2 bg-gray-500 rounded-full animate-bounce" style="animation-delay: 0ms"></div>
          <div class="w-2 h-2 bg-gray-500 rounded-full animate-bounce" style="animation-delay: 150ms"></div>
          <div class="w-2 h-2 bg-gray-500 rounded-full animate-bounce" style="animation-delay: 300ms"></div>
        </div>
        <span class="text-sm">{{ t('chat.typing') }}</span>
      </div>
    </div>

    <!-- Execute Mode Preview Modal -->
    <ExecutePreviewModal
      v-if="showExecutePreview"
      :file-count="previewFileCount"
      @confirm="confirmExecute"
      @cancel="cancelExecute"
    />

    <!-- Input -->
    <div class="border-t border-gray-700 p-4">
      <!-- Auto-suggest Panel -->
      <AutoSuggestPanel
        :visible="autoSuggest.shouldShow.value"
        :suggestions="autoSuggest.suggestions.value"
        :selected-paths="autoSuggest.selectedPaths.value"
        :get-file-name="autoSuggest.getFileName"
        :get-file-path="autoSuggest.getFilePath"
        :get-source-badge-class="autoSuggest.getSourceBadgeClass"
        :get-source-label="autoSuggest.getSourceLabel"
        :get-relevance-percent="autoSuggest.getRelevancePercent"
        @toggle="autoSuggest.toggleSelect"
        @add="handleAddSuggestedFiles"
        @hide="autoSuggest.hide"
      />

      <div class="flex gap-2">
        <textarea
          v-model="inputMessage"
          :placeholder="inputPlaceholder"
          class="flex-1 px-4 py-3 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-purple-500 resize-none"
          rows="3"
          @keydown.ctrl.enter="handleSend"
        ></textarea>
        <button
          @click="handleSend"
          :disabled="!canSend"
          class="action-btn action-btn-accent px-6 py-3"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
          </svg>
        </button>
      </div>
      <p class="text-xs text-gray-400 mt-2">Ctrl+Enter {{ t('chat.toSend') }}</p>
    </div>

    <!-- Change Preview Modal -->
    <ChangePreviewModal
      v-if="showChangesPanel"
      @close="showChangesPanel = false"
      @applied="onChangesApplied"
    />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useFileStore } from '@/features/files'
import { computed, nextTick, ref, watch } from 'vue'
import { useChatStore } from '../model/chat.store'
import { useSandboxStore } from '@/stores/sandbox.store'
import { useChatDragDrop } from '../composables/useChatDragDrop'
import { useAutoSuggest } from '../composables/useAutoSuggest'
import MessageItem from './MessageItem.vue'
import ChangePreviewModal from './ChangePreviewModal.vue'
import ChatModeSelector from './ChatModeSelector.vue'
import ChatHistoryPanel from './ChatHistoryPanel.vue'
import TokenBudgetBar from './TokenBudgetBar.vue'
import ExecutePreviewModal from './ExecutePreviewModal.vue'
import AutoSuggestPanel from './AutoSuggestPanel.vue'

const { t } = useI18n()
const chatStore = useChatStore()
const sandboxStore = useSandboxStore()
const fileStore = useFileStore()
const { isDragOver, handleDragEnter, handleDragOver, handleDragLeave, handleDrop } = useChatDragDrop()
const inputMessage = ref('')
const messagesContainer = ref<HTMLElement>()
const showChangesPanel = ref(false)
const showExecutePreview = ref(false)
const previewFileCount = ref(0)
const pendingMessage = ref('')

// Auto-suggest composable
const selectedFilesList = computed(() => fileStore.selectedFilesList || [])
const autoSuggest = useAutoSuggest({
  inputText: inputMessage,
  currentFiles: selectedFilesList,
  onAddFiles: (files) => {
    files.forEach(path => fileStore.selectPath(path))
  }
})

function handleAddSuggestedFiles() {
  autoSuggest.addSelected()
}

const canSend = computed(() => {
  return inputMessage.value.trim().length > 0 && !chatStore.isStreaming
})

const inputPlaceholder = computed(() => {
  return chatStore.chatMode === 'explore' 
    ? t('chat.placeholder') 
    : t('chat.placeholderWithContext')
})

async function handleSend() {
  if (!canSend.value) return
  
  const message = inputMessage.value.trim()
  
  // In execute mode, show preview if >5 files selected
  if (chatStore.chatMode === 'execute') {
    const selectedCount = fileStore.selectedFilesList?.length || 0
    if (selectedCount > 5) {
      pendingMessage.value = message
      previewFileCount.value = selectedCount
      showExecutePreview.value = true
      return
    }
  }
  
  await sendMessageInternal(message)
}

async function sendMessageInternal(message: string) {
  inputMessage.value = ''
  
  // TODO: Pass exploreOnly flag to sendMessage when backend supports it
  // const exploreOnly = chatStore.chatMode === 'explore'
  
  await chatStore.sendMessage(message, undefined)
  scrollToBottom()
  
  // Refresh sandbox state after AI response
  await sandboxStore.refresh()
}

function confirmExecute() {
  showExecutePreview.value = false
  if (pendingMessage.value) {
    sendMessageInternal(pendingMessage.value)
    pendingMessage.value = ''
  }
}

function cancelExecute() {
  showExecutePreview.value = false
  pendingMessage.value = ''
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

function onChangesApplied() {
  showChangesPanel.value = false
}

watch(() => chatStore.messages.length, () => {
  scrollToBottom()
})
</script>

<style scoped>
.chat-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border-subtle);
  background: var(--bg-1);
  flex-wrap: wrap;
  gap: 0.5rem;
}

@media (max-width: 480px) {
  .chat-toolbar {
    flex-direction: column;
    align-items: stretch;
  }
}

.chat-drop-overlay {
  position: absolute;
  inset: 0;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
}

.chat-drop-zone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem 3rem;
  border: 2px dashed var(--purple-500);
  border-radius: 1rem;
  background: rgba(139, 92, 246, 0.1);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
