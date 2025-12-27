<template>
  <div class="chat-mode-selector">
    <button
      :class="['mode-btn', { 'mode-btn-active': chatStore.chatMode === 'explore' }]"
      @click="chatStore.setChatMode('explore')"
      :title="t('chat.mode.exploreDesc')"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      <span>{{ t('chat.mode.explore') }}</span>
    </button>
    
    <button
      :class="['mode-btn', { 'mode-btn-active': chatStore.chatMode === 'execute' }]"
      @click="chatStore.setChatMode('execute')"
      :title="t('chat.mode.executeDesc')"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
      </svg>
      <span>{{ t('chat.mode.execute') }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useChatStore } from '../model/chat.store'

const { t } = useI18n()
const chatStore = useChatStore()
</script>

<style scoped>
.chat-mode-selector {
  display: flex;
  gap: 0.5rem;
}

.mode-btn {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.75rem;
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: var(--radius-xl);
  background: var(--bg-2);
  color: var(--text-muted);
  border: 1px solid var(--border-default);
  transition: all 150ms ease-out;
  cursor: pointer;
}

.mode-btn:hover:not(.mode-btn-active) {
  background: var(--bg-3);
  color: var(--text-secondary);
  border-color: var(--border-strong);
}

.mode-btn-active {
  background: var(--accent-purple-bg);
  color: white;
  border-color: var(--accent-purple-border);
  box-shadow: 0 0 12px rgba(168, 85, 247, 0.3);
}

.mode-btn-active:hover {
  background: rgba(168, 85, 247, 0.4);
}
</style>
