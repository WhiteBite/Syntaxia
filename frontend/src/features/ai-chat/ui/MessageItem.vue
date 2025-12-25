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

      <!-- Actions -->
      <div class="flex gap-2 mt-2">
        <button @click="copyToClipboard" class="text-xs text-gray-400 hover:text-gray-300 transition-colors"
          :title="t('chat.copy')">
          <svg v-if="!copied" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <svg v-else class="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </button>
        <button v-if="message.role === 'user'" @click="$emit('edit', message.id, message.content)"
          class="text-xs text-gray-400 hover:text-gray-300 transition-colors" title="Edit">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </button>
        <button @click="$emit('delete', message.id)" class="text-xs text-gray-400 hover:text-red-400 transition-colors"
          title="Delete">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useUIStore } from '@/stores/ui.store'
import { onUnmounted, ref } from 'vue'
import type { Message } from '../model/chat.store'

interface Props {
  message: Message
}

const props = defineProps<Props>()

defineEmits<{
  (e: 'delete', messageId: string): void
  (e: 'edit', messageId: string, content: string): void
}>()

const uiStore = useUIStore()
const { t } = useI18n()
const copied = ref(false)
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
