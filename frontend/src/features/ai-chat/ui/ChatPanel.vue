<template>
  <div class="h-full flex flex-col bg-gray-900">
    <!-- Header -->
    <div class="section-header">
      <div class="section-title">
        <div class="section-icon section-icon-purple">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
          </svg>
        </div>
        <span class="section-title-text">{{ t('chat.title') }}</span>
        <span class="badge badge-primary">{{ t('chat.comingSoon') }}</span>
      </div>
      
      <button
        @click="chatStore.clearChat"
        :disabled="!chatStore.hasMessages"
        class="btn btn-secondary btn-xs"
      >
        {{ t('chat.clear') }}
      </button>
    </div>

    <!-- Messages -->
    <div ref="messagesContainer" class="flex-1 overflow-auto p-4 space-y-4">
      <div v-if="!chatStore.hasMessages" class="empty-state h-full">
        <div class="empty-state-icon !w-20 !h-20 !rounded-2xl">
          <svg class="!w-10 !h-10 text-purple-500/50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
        </div>
        <p class="empty-state-title !text-lg">{{ t('chat.comingSoonTitle') }}</p>
        <p class="empty-state-text mb-4">{{ t('chat.comingSoonDesc') }}</p>
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

    <!-- Input -->
    <div class="border-t border-gray-700 p-4">
      <div class="flex gap-2">
        <textarea
          v-model="inputMessage"
          :placeholder="`${t('chat.placeholder')} (${t('chat.comingSoon')})`"
          disabled
          class="flex-1 px-4 py-3 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-purple-500 resize-none disabled:opacity-50 disabled:cursor-not-allowed"
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
      <p class="text-xs text-gray-400 mt-2">Ctrl+Enter to send</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed, nextTick, ref, watch } from 'vue'
import { useChatStore } from '../model/chat.store'
import MessageItem from './MessageItem.vue'

const { t } = useI18n()
const chatStore = useChatStore()
const inputMessage = ref('')
const messagesContainer = ref<HTMLElement>()

const canSend = computed(() => {
  return inputMessage.value.trim().length > 0 && !chatStore.isStreaming
})

async function handleSend() {
  if (!canSend.value) return
  
  const message = inputMessage.value.trim()
  inputMessage.value = ''
  
  await chatStore.sendMessage(message)
  scrollToBottom()
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

watch(() => chatStore.messages.length, () => {
  scrollToBottom()
})
</script>
