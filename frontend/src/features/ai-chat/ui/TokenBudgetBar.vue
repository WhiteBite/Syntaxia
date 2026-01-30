<template>
  <div 
    class="token-budget-bar"
    :title="tooltipText"
  >
    <svg class="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
    </svg>
    
    <div class="token-progress-container">
      <div 
        class="token-progress-bar"
        :class="progressColorClass"
        :style="{ width: progressPercent + '%' }"
      ></div>
    </div>
    
    <span class="token-text">
      {{ formattedUsed }} / {{ formattedLimit }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { useChatStore } from '../model/chat.store'

const { t } = useI18n()
const chatStore = useChatStore()

const progressPercent = computed(() => {
  const { used, limit } = chatStore.tokenBudget
  if (limit === 0) return 0
  return Math.min((used / limit) * 100, 100)
})

const progressColorClass = computed(() => {
  const percent = progressPercent.value
  if (percent < 50) return 'progress-green'
  if (percent < 80) return 'progress-yellow'
  return 'progress-red'
})

function formatTokens(value: number): string {
  if (value >= 1000) {
    return Math.round(value / 1000) + 'k'
  }
  return value.toString()
}

const formattedUsed = computed(() => formatTokens(chatStore.tokenBudget.used))
const formattedLimit = computed(() => formatTokens(chatStore.tokenBudget.limit))

const tooltipText = computed(() => {
  const { used, limit } = chatStore.tokenBudget
  return `${t('chat.tokens.used')}: ${used.toLocaleString()} / ${t('chat.tokens.limit')}: ${limit.toLocaleString()}`
})
</script>

<style scoped>
.token-budget-bar {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-1) var(--space-3);
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-full);
  color: var(--text-muted);
  font-size: var(--font-size-xs);
  cursor: default;
}

.token-progress-container {
  width: calc(4rem * var(--ui-scale));
  height: calc(0.375rem * var(--ui-scale));
  background: var(--bg-3);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.token-progress-bar {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width 300ms ease-out, background-color 300ms ease-out;
}


.progress-green {
  background: var(--color-success);
  box-shadow: 0 0 8px rgba(74, 222, 128, 0.4);
}

.progress-yellow {
  background: var(--color-warning);
  box-shadow: 0 0 8px rgba(251, 191, 36, 0.4);
}

.progress-red {
  background: var(--color-danger);
  box-shadow: 0 0 8px rgba(248, 113, 113, 0.4);
}

.token-text {
  font-weight: 500;
  white-space: nowrap;
}
</style>
