<template>
  <div
    class="group"
    :class="[message.role === 'user' ? 'message-user' : 'message-assistant']"
  >
    <div class="flex items-start gap-2">
      <div v-if="message.role === 'assistant'" class="message-avatar">
        <svg class="w-3.5 h-3.5 text-purple-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
        </svg>
      </div>
      <div class="flex-1 min-w-0">
        <p class="text-xs text-gray-300 whitespace-pre-wrap break-words">{{ message.content }}</p>
        
        <!-- Context attached info -->
        <div v-if="message.contextAttached" class="mt-2 flex items-center gap-1 text-[10px] text-indigo-400">
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
          </svg>
          {{ t('chat.contextAttached') }} ({{ message.tokenCount || 0 }} tokens)
        </div>
        
        <!-- Tool Calls Info -->
        <div v-if="message.toolCalls && message.toolCalls.length > 0" class="mt-2">
          <button
            @click="$emit('toggle-tools', index)"
            class="flex items-center gap-1 text-[10px] text-emerald-400 hover:text-emerald-300"
          >
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            🔧 {{ message.toolCalls.length }} tools ({{ message.iterations }} iter)
            <svg :class="['w-3 h-3 transition-transform', expandedToolCalls.has(index) ? 'rotate-180' : '']" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>
          <div v-if="expandedToolCalls.has(index)" class="mt-2 space-y-1 text-[10px]">
            <div
              v-for="(tc, tcIdx) in message.toolCalls"
              :key="tcIdx"
              class="p-1.5 bg-gray-900/50 rounded border border-gray-700/30"
            >
              <div class="font-medium text-emerald-400">{{ tc.tool }}</div>
              <div class="text-gray-400 truncate" :title="tc.arguments">{{ tc.arguments }}</div>
            </div>
          </div>
        </div>
        
        <!-- Error -->
        <div v-if="message.error" class="mt-2 text-[10px] text-red-400">
          {{ message.error }}
        </div>
      </div>
      
      <!-- Copy button -->
      <button
        v-if="message.role === 'assistant'"
        @click="$emit('copy')"
        class="icon-btn-sm opacity-0 group-hover:opacity-100 transition-opacity"
        :title="t('chat.copy')"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';
import type { Message as ChatMessage } from '@/features/ai-chat/composables/useChatMessages';

defineProps<{
  message: ChatMessage
  index: number
  expandedToolCalls: Set<number>
}>()

defineEmits<{
  copy: []
  'toggle-tools': [index: number]
}>()

const { t } = useI18n()
</script>

<style scoped>
/* User messages - right aligned */
.message-user {
  display: flex;
  justify-content: flex-end;
  padding: 4px 0;
}

.message-user .flex {
  flex-direction: row-reverse;
  max-width: 85%;
}

.message-user .flex-1 {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.3) 0%, rgba(139, 92, 246, 0.2) 100%);
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: 16px 16px 4px 16px;
  padding: 10px 14px;
}

.message-user .text-xs {
  color: #e0e7ff;
}

/* Assistant messages - left aligned */
.message-assistant {
  display: flex;
  justify-content: flex-start;
  padding: 4px 0;
}

.message-assistant .flex {
  max-width: 90%;
}

.message-assistant .flex-1 {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 4px 16px 16px 16px;
  padding: 10px 14px;
}

/* Avatar styling */
.message-avatar {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.2) 0%, rgba(99, 102, 241, 0.1) 100%);
  border: 1px solid rgba(139, 92, 246, 0.3);
  border-radius: 50%;
}
</style>
