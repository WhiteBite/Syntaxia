<template>
  <div :class="[
    'flex gap-3 p-4 rounded-lg',
    message.role === 'user' ? 'bg-blue-900/20' : 'bg-gray-800'
  ]">
    <!-- Avatar -->
    <div class="flex-shrink-0">
      <div :class="[
        'w-8 h-8 rounded-full flex items-center justify-center',
        message.role === 'user' ? 'bg-blue-600' : 'bg-purple-600'
      ]">
        <svg v-if="message.role === 'user'" class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
        </svg>
        <svg v-else class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
          <path d="M2 5a2 2 0 012-2h7a2 2 0 012 2v4a2 2 0 01-2 2H9l-3 3v-3H4a2 2 0 01-2-2V5z" />
          <path
            d="M15 7v2a4 4 0 01-4 4H9.828l-1.766 1.767c.28.149.599.233.938.233h2l3 3v-3h2a2 2 0 002-2V9a2 2 0 00-2-2h-1z" />
        </svg>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 min-w-0">
      <div class="flex items-center justify-between mb-1">
        <span class="text-sm font-semibold text-white">
          {{ message.role === 'user' ? 'You' : 'AI Assistant' }}
        </span>
        <span class="text-xs text-gray-400">{{ formatTime(message.timestamp) }}</span>
      </div>

      <div class="text-gray-300 text-sm whitespace-pre-wrap break-words">
        {{ message.content }}
      </div>

      <!-- Context Used Badge (only for assistant messages with context) -->
      <div v-if="message.role === 'assistant' && message.contextUsed" class="mt-3">
        <ContextUsedBadge
          :context="message.contextUsed"
          @toggle="contextExpanded = $event"
        />
        <ContextUsedDetails
          v-if="contextExpanded"
          :context="message.contextUsed"
          @add-context="$emit('addContext')"
        />
      </div>

      <!-- Actions -->
      <div class="flex gap-2 mt-2">
        <BaseButton
          variant="ghost"
          size="xs"
          icon-only
          @click="copyToClipboard"
          :title="t('chat.copy')"
          class="h-7 w-7"
        >
          <Copy v-if="!copied" class="w-3.5 h-3.5" />
          <Check v-else class="w-3.5 h-3.5 text-green-500" />
        </BaseButton>
        
        <BaseButton
          v-if="message.role === 'user'"
          variant="ghost"
          size="xs"
          icon-only
          @click="$emit('edit', message.id, message.content)"
          title="Edit"
          class="h-7 w-7"
        >
          <Pencil class="w-3.5 h-3.5" />
        </BaseButton>

        <BaseButton
          variant="ghost"
          size="xs"
          icon-only
          @click="$emit('delete', message.id)"
          title="Delete"
          class="text-gray-400 hover:text-red-400 h-7 w-7"
        >
          <Trash2 class="w-3.5 h-3.5" />
        </BaseButton>
      </div>
    </div>
  </div>
</template>


<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useUIStore } from '@/stores/ui.store'
import { onUnmounted, ref } from 'vue'
import type { Message } from '../model/chat.store'
import ContextUsedBadge from './ContextUsedBadge.vue'
import ContextUsedDetails from './ContextUsedDetails.vue'
import { BaseButton } from '@/components/ui'
import { Copy, Check, Pencil, Trash2 } from 'lucide-vue-next'
import type { ContextUsed } from './types'


interface MessageWithContext extends Message {
  contextUsed?: ContextUsed
}

interface Props {
  message: MessageWithContext
}

const props = defineProps<Props>()

defineEmits<{
  (e: 'delete', messageId: string): void
  (e: 'edit', messageId: string, content: string): void
  (e: 'addContext'): void
}>()

const uiStore = useUIStore()
const { t } = useI18n()
const copied = ref(false)
const contextExpanded = ref(false)
let copyTimeoutId: ReturnType<typeof setTimeout> | null = null

// Очистка таймера при размонтировании компонента
onUnmounted(() => {
  if (copyTimeoutId) {
    clearTimeout(copyTimeoutId)
    copyTimeoutId = null
  }
})

function formatTime(timestamp: string): string {
  const date = new Date(timestamp)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

async function copyToClipboard() {
  // Очищаем предыдущий таймер если есть
  if (copyTimeoutId) {
    clearTimeout(copyTimeoutId)
    copyTimeoutId = null
  }
  
  try {
    await navigator.clipboard.writeText(props.message.content)
    uiStore.addToast(t('chat.copied'), 'success')
    copied.value = true
    copyTimeoutId = setTimeout(() => {
      copied.value = false
      copyTimeoutId = null
    }, 2000)
  } catch (error) {
    // Fallback for older browsers
    const textarea = document.createElement('textarea')
    textarea.value = props.message.content
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select()
    try {
      document.execCommand('copy')
      uiStore.addToast(t('chat.copied'), 'success')
      copied.value = true
      copyTimeoutId = setTimeout(() => {
        copied.value = false
        copyTimeoutId = null
      }, 2000)
    } catch {
      uiStore.addToast(t('chat.copyFailed'), 'error')
    }
    document.body.removeChild(textarea)
  }
}
</script>
