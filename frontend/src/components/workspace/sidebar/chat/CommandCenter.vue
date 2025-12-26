<template>
  <div class="command-footer">
    <!-- Context Chips -->
    <Transition name="chips">
      <div v-if="hasAnySelection" class="context-chips">
        <button class="chip chip--files" @click="toggleContextExpand">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          <span>{{ totalFileCount }} {{ t('context.filesShort') }}</span>
        </button>
        
        <div class="chip chip--tokens">
          <span>{{ formattedTokens }} {{ t('context.tokens') }}</span>
        </div>

        <button 
          v-if="hasContext" 
          class="chip-clear" 
          @click="clearContext"
        >
          {{ t('context.clear') }}
        </button>
      </div>
    </Transition>

    <!-- Expanded file list -->
    <Transition name="expand">
      <div v-if="isContextExpanded && hasAnySelection" class="context-list">
        <div v-for="file in displayedFiles" :key="file" class="context-item">
          <span class="context-item__name">{{ getFileName(file) }}</span>
          <span class="context-item__path">{{ getFilePath(file) }}</span>
          <button class="context-item__remove" @click="removeFile(file)">×</button>
        </div>
        <button v-if="hasMoreFiles" class="context-more" @click="showAllFiles = !showAllFiles">
          {{ showAllFiles ? t('common.showLess') : `+${hiddenCount}` }}
        </button>
      </div>
    </Transition>

    <!-- Input Capsule -->
    <div class="capsule-wrapper">
      <!-- Glow effect -->
      <div class="capsule-glow" :class="{ 'capsule-glow--active': isFocused }"></div>
      
      <!-- Main capsule -->
      <div class="capsule" :class="{ 'capsule--focused': isFocused }">
        <!-- Mention popup -->
        <Transition name="popup">
          <div v-if="showMentionPopup" class="mention-popup">
            <div class="mention-header">
              <span>{{ t('chat.mentions.title') }}</span>
              <kbd>↑↓</kbd>
            </div>
            <button 
              v-for="(item, idx) in mentionItems" 
              :key="item.id"
              class="mention-item"
              :class="{ 'mention-item--active': selectedMentionIndex === idx }"
              @click="selectMention(item)"
              @mouseenter="selectedMentionIndex = idx"
            >
              <span class="mention-icon" :class="`mention-icon--${item.id}`">
                <svg v-if="item.id === 'files'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
                </svg>
                <svg v-else-if="item.id === 'git'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
                <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                </svg>
              </span>
              <div class="mention-content">
                <span class="mention-label">{{ item.label }}</span>
                <span class="mention-hint">{{ item.hint }}</span>
              </div>
            </button>
          </div>
        </Transition>

        <!-- Textarea -->
        <textarea
          ref="inputRef"
          :value="modelValue"
          @input="handleInput"
          @keydown="handleKeydown"
          @focus="isFocused = true"
          @blur="handleBlur"
          class="capsule-input"
          :placeholder="t('chat.placeholder')"
          :disabled="isDisabled"
          rows="1"
        ></textarea>

        <!-- Toolbar -->
        <div class="capsule-toolbar">
          <div class="toolbar-left">
            <button class="toolbar-btn" @click="handleAttachClick" :title="t('chat.attachFiles')">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" 
                  d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
              </svg>
            </button>
            <button class="toolbar-btn toolbar-btn--mention" @click="handleMentionClick" :title="t('chat.hints.mention')">@</button>
          </div>

          <button
            class="send-btn"
            :class="{ 'send-btn--ready': hasText }"
            @click="handleSend"
            :disabled="!hasText || isDisabled"
          >
            <svg v-if="!isThinking" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18" />
            </svg>
            <svg v-else class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Hints -->
    <div class="hints">
      <span>↵ {{ t('chat.send') }}</span>
      <span>⇧↵ {{ t('chat.hints.newLine') }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useMentions } from '@/features/ai-chat/composables/useMentions'
import { useContextStore } from '@/features/context'
import { useFileStore } from '@/features/files'
import { computed, nextTick, ref, watch } from 'vue'

const props = defineProps<{
  modelValue: string
  isThinking: boolean
  isAnalyzing: boolean
  hasMessages: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  send: []
  clear: []
  attach: []
}>()

const { t } = useI18n()
const contextStore = useContextStore()
const fileStore = useFileStore()
const { processFilesMention, processGitMention, processProblemsMention } = useMentions()

const MAX_VISIBLE_FILES = 5

const inputRef = ref<HTMLTextAreaElement | null>(null)
const isFocused = ref(false)
const isContextExpanded = ref(false)
const showAllFiles = ref(false)
const showMentionPopup = ref(false)
const selectedMentionIndex = ref(0)

const hasText = computed(() => props.modelValue.trim().length > 0)
const isDisabled = computed(() => props.isThinking || props.isAnalyzing)
const hasContext = computed(() => contextStore.hasContext)
const fileCount = computed(() => contextStore.fileCount)
const files = computed(() => contextStore.summary?.files || [])
const selectedFilesCount = computed(() => fileStore.selectedCount || 0)
const totalFileCount = computed(() => hasContext.value ? fileCount.value : selectedFilesCount.value)
const hasAnySelection = computed(() => hasContext.value || selectedFilesCount.value > 0)

const formattedTokens = computed(() => {
  const tokens = contextStore.totalTokens
  if (tokens >= 1000) return `${(tokens / 1000).toFixed(1)}k`
  return tokens.toString()
})

const displayedFiles = computed(() => 
  showAllFiles.value ? files.value : files.value.slice(0, MAX_VISIBLE_FILES)
)
const hasMoreFiles = computed(() => files.value.length > MAX_VISIBLE_FILES)
const hiddenCount = computed(() => files.value.length - MAX_VISIBLE_FILES)

const mentionItems = [
  { id: 'files', label: '@files', hint: t('chat.mentions.files') },
  { id: 'git', label: '@git', hint: t('chat.mentions.git') },
  { id: 'problems', label: '@problems', hint: t('chat.mentions.problems') },
]

function handleInput(event: Event) {
  const textarea = event.target as HTMLTextAreaElement
  emit('update:modelValue', textarea.value)
  autoResize(textarea)
  
  const lastChar = textarea.value.slice(-1)
  if (lastChar === '@') {
    showMentionPopup.value = true
    selectedMentionIndex.value = 0
  } else if (showMentionPopup.value && !textarea.value.includes('@')) {
    showMentionPopup.value = false
  }
}

function handleKeydown(event: KeyboardEvent) {
  if (showMentionPopup.value) {
    if (event.key === 'ArrowDown') {
      event.preventDefault()
      selectedMentionIndex.value = (selectedMentionIndex.value + 1) % mentionItems.length
    } else if (event.key === 'ArrowUp') {
      event.preventDefault()
      selectedMentionIndex.value = (selectedMentionIndex.value - 1 + mentionItems.length) % mentionItems.length
    } else if (event.key === 'Enter' || event.key === 'Tab') {
      event.preventDefault()
      selectMention(mentionItems[selectedMentionIndex.value])
    } else if (event.key === 'Escape') {
      showMentionPopup.value = false
    }
    return
  }

  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    handleSend()
  }
}

function handleBlur() {
  setTimeout(() => {
    isFocused.value = false
    showMentionPopup.value = false
  }, 150)
}

function handleSend() {
  if (hasText.value && !isDisabled.value) {
    emit('send')
  }
}

async function selectMention(item: { id: string; label: string; hint: string }) {
  showMentionPopup.value = false
  
  // Execute the mention action immediately
  switch (item.id) {
    case 'files':
      await processFilesMention(t)
      break
    case 'git':
      await processGitMention(t)
      break
    case 'problems':
      await processProblemsMention(t)
      break
  }
  
  nextTick(() => inputRef.value?.focus())
}

function autoResize(textarea: HTMLTextAreaElement) {
  textarea.style.height = '24px'
  textarea.style.height = Math.min(textarea.scrollHeight, 160) + 'px'
}

function toggleContextExpand() {
  isContextExpanded.value = !isContextExpanded.value
}

function clearContext() {
  contextStore.clearContext()
}

function removeFile(file: string) {
  contextStore.removeFileFromContext(file)
}

function getFileName(path: string): string {
  return path.split('/').pop() || path
}

function getFilePath(path: string): string {
  const parts = path.split('/')
  return parts.length <= 1 ? '' : parts.slice(0, -1).join('/')
}

function handleMentionClick() {
  showMentionPopup.value = true
  selectedMentionIndex.value = 0
  nextTick(() => inputRef.value?.focus())
}

async function handleAttachClick() {
  await processFilesMention(t)
  nextTick(() => inputRef.value?.focus())
}

watch(() => props.isThinking, (thinking) => {
  if (!thinking) nextTick(() => inputRef.value?.focus())
})
</script>

<style scoped>
/* Footer container - fixed at bottom, no absolute */
.command-footer {
  flex: none;
  padding: 8px 12px 12px;
  background: linear-gradient(to top, rgba(15, 17, 26, 0.95) 0%, transparent 100%);
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* Context Chips */
.context-chips {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 2px;
}

.chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  font-size: 11px;
  font-weight: 500;
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.chip--files {
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid rgba(99, 102, 241, 0.2);
  color: #a5b4fc;
}

.chip--files:hover {
  background: rgba(99, 102, 241, 0.2);
  border-color: rgba(99, 102, 241, 0.4);
}

.chip--tokens {
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.2);
  color: #6ee7b7;
}

.chip-clear {
  margin-left: auto;
  padding: 4px 8px;
  font-size: 11px;
  color: #6b7280;
  background: transparent;
  border: none;
  cursor: pointer;
  transition: color 0.15s ease;
}

.chip-clear:hover {
  color: #f87171;
}

/* Context List */
.context-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
  background: rgba(15, 17, 24, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  max-height: 120px;
  overflow-y: auto;
}

.context-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 8px;
  border-radius: 6px;
}

.context-item:hover {
  background: rgba(255, 255, 255, 0.03);
}

.context-item__name {
  font-size: 12px;
  font-weight: 500;
  color: #e5e7eb;
}

.context-item__path {
  flex: 1;
  font-size: 11px;
  color: #6b7280;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.context-item__remove {
  padding: 2px 6px;
  font-size: 14px;
  color: #6b7280;
  background: transparent;
  border: none;
  cursor: pointer;
  opacity: 0;
  transition: all 0.15s ease;
}

.context-item:hover .context-item__remove {
  opacity: 1;
}

.context-item__remove:hover {
  color: #f87171;
}

.context-more {
  padding: 4px 8px;
  font-size: 11px;
  color: #9ca3af;
  background: transparent;
  border: none;
  cursor: pointer;
}

/* Capsule Wrapper */
.capsule-wrapper {
  position: relative;
}

.capsule-glow {
  position: absolute;
  inset: -1px;
  background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
  border-radius: 12px;
  opacity: 0;
  filter: blur(6px);
  transition: opacity 0.3s ease;
  pointer-events: none;
}

.capsule-glow--active {
  opacity: 0.12;
}

/* Main Capsule */
.capsule {
  position: relative;
  display: flex;
  flex-direction: column;
  background: rgba(28, 31, 46, 0.8);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.25);
  overflow: visible;
  transition: all 0.2s ease;
}

.capsule--focused {
  border-color: rgba(139, 92, 246, 0.5);
  box-shadow: 
    0 4px 16px rgba(0, 0, 0, 0.3),
    0 0 0 1px rgba(139, 92, 246, 0.2);
}

/* Input */
.capsule-input {
  width: 100%;
  min-height: 40px;
  max-height: 120px;
  padding: 10px 12px;
  background: transparent;
  border: none;
  font-size: 13px;
  line-height: 1.4;
  color: #f3f4f6;
  resize: none;
  outline: none;
}

.capsule-input::placeholder {
  color: #6b7280;
}

.capsule-input:disabled {
  opacity: 0.5;
}

/* Toolbar */
.capsule-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 8px 6px;
}

.toolbar-left {
  display: flex;
  gap: 2px;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s ease;
}

.toolbar-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #e5e7eb;
}

.toolbar-btn--mention {
  font-size: 14px;
  font-weight: 600;
}

.toolbar-btn--mention:hover {
  color: #a78bfa;
  background: rgba(139, 92, 246, 0.1);
}

/* Send Button */
.send-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: rgba(139, 92, 246, 0.2);
  border: none;
  border-radius: 8px;
  color: #a78bfa;
  cursor: pointer;
  transition: all 0.15s ease;
}

.send-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.send-btn--ready {
  background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
}

.send-btn--ready:hover:not(:disabled) {
  transform: scale(1.05);
  box-shadow: 0 6px 20px rgba(139, 92, 246, 0.5);
}

/* Hints */
.hints {
  display: flex;
  justify-content: center;
  gap: 12px;
  font-size: 9px;
  color: #4b5563;
}

/* Mention Popup */
.mention-popup {
  position: absolute;
  bottom: 100%;
  left: 0;
  right: 0;
  margin-bottom: 8px;
  padding: 8px;
  background: rgba(22, 27, 34, 0.98);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 14px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5);
  z-index: 50;
}

.mention-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 10px 10px;
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: #6b7280;
}

.mention-header kbd {
  padding: 2px 6px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  font-size: 9px;
  color: #9ca3af;
}

.mention-item {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 10px 12px;
  background: transparent;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.15s ease;
  text-align: left;
}

.mention-item:hover,
.mention-item--active {
  background: rgba(139, 92, 246, 0.15);
}

.mention-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
}

.mention-icon--files {
  background: rgba(139, 92, 246, 0.15);
  color: #a78bfa;
}

.mention-icon--git {
  background: rgba(249, 115, 22, 0.15);
  color: #fb923c;
}

.mention-icon--problems {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
}

.mention-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.mention-label {
  font-size: 13px;
  font-weight: 600;
  color: #e5e7eb;
}

.mention-hint {
  font-size: 11px;
  color: #6b7280;
}

/* Transitions */
.chips-enter-active,
.chips-leave-active {
  transition: all 0.2s ease;
}

.chips-enter-from,
.chips-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

.expand-enter-active,
.expand-leave-active {
  transition: all 0.2s ease;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
}

.popup-enter-active,
.popup-leave-active {
  transition: all 0.15s ease;
}

.popup-enter-from,
.popup-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
