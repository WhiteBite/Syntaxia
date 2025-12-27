<template>
  <div class="chat-history-panel">
    <!-- Toggle Button -->
    <button
      @click="toggleDropdown"
      class="history-toggle-btn"
      :title="t('chat.history.title')"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span v-if="hasHistory" class="history-count">
        {{ historyCount }}
      </span>
    </button>

    <!-- Dropdown Panel -->
    <Transition name="dropdown">
      <div v-if="isOpen" class="history-dropdown">
        <!-- Header -->
        <div class="history-header">
          <span class="history-title">{{ t('chat.history.title') }}</span>
          <button
            @click="handleNewChat"
            class="btn btn-primary btn-xs"
          >
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            {{ t('chat.history.newChat') }}
          </button>
        </div>

        <!-- Chat List -->
        <div class="history-list">
          <div v-if="!hasHistory" class="history-empty">
            <svg class="w-8 h-8 text-gray-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
            <p class="text-gray-500 text-xs">{{ t('chat.history.empty') }}</p>
          </div>

          <div
            v-for="chat in chatHistory"
            :key="chat.id"
            :class="['history-item', { 'history-item-active': chat.id === currentChatId }]"
            @click="handleLoadChat(chat.id)"
          >
            <div class="history-item-content">
              <div class="history-item-title" :title="chat.title">
                {{ chat.title }}
              </div>
              <div class="history-item-meta">
                <span class="history-item-date">{{ formatDate(chat.updatedAt) }}</span>
                <span class="history-item-messages">
                  {{ chat.messages.length }} {{ t('chat.history.messages') }}
                </span>
              </div>
            </div>
            <button
              @click.stop="handleDeleteChat(chat.id)"
              class="history-item-delete"
              :title="t('chat.history.delete')"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Backdrop -->
    <div v-if="isOpen" class="history-backdrop" @click="closeDropdown"></div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useChatHistory } from '../composables/useChatHistory'

const { t } = useI18n()

const {
  isOpen,
  chatHistory,
  currentChatId,
  historyCount,
  hasHistory,
  formatDate,
  handleNewChat,
  handleLoadChat,
  handleDeleteChat,
  toggleDropdown,
  closeDropdown
} = useChatHistory()
</script>

<style scoped>
.chat-history-panel {
  position: relative;
}

.history-toggle-btn {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.375rem 0.5rem;
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: var(--radius-lg);
  background: var(--bg-2);
  color: var(--text-muted);
  border: 1px solid var(--border-default);
  transition: all 150ms ease-out;
  cursor: pointer;
}

.history-toggle-btn:hover {
  background: var(--bg-3);
  color: var(--text-secondary);
  border-color: var(--border-strong);
}

.history-count {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 1.25rem;
  height: 1.25rem;
  padding: 0 0.25rem;
  font-size: 0.625rem;
  font-weight: 600;
  background: var(--accent-purple-bg);
  color: white;
  border-radius: var(--radius-full);
}

.history-dropdown {
  position: absolute;
  top: calc(100% + 0.5rem);
  right: 0;
  z-index: 50;
  width: min(20rem, 90vw);
  max-height: min(24rem, 60vh);
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.history-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem;
  border-bottom: 1px solid var(--border-subtle);
  background: var(--bg-2);
}

.history-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-primary);
}

.history-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.history-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem 1rem;
  text-align: center;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem;
  border-bottom: 1px solid var(--border-subtle);
  cursor: pointer;
  transition: background 150ms ease-out;
}

.history-item:hover {
  background: var(--bg-2);
}

.history-item-active {
  background: var(--accent-purple-bg);
  border-left: 2px solid var(--accent-purple);
}

.history-item-content {
  flex: 1;
  min-width: 0;
}

.history-item-title {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.history-item-meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.25rem;
  font-size: 0.6875rem;
  color: var(--text-muted);
}

.history-item-date {
  color: var(--text-muted);
}

.history-item-messages {
  color: var(--text-muted);
}

.history-item-delete {
  flex-shrink: 0;
  padding: 0.25rem;
  color: var(--text-muted);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  opacity: 0;
  transition: all 150ms ease-out;
}

.history-item:hover .history-item-delete {
  opacity: 1;
}

.history-item-delete:hover {
  color: var(--accent-red);
  background: var(--bg-3);
}

.history-backdrop {
  position: fixed;
  inset: 0;
  z-index: 40;
}

/* Dropdown animation */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 150ms ease-out;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-0.5rem);
}
</style>
