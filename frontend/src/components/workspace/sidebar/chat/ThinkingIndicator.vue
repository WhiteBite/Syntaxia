<template>
  <div class="thinking-indicator">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <div class="flex gap-1">
          <span class="thinking-dot animate-bounce" style="animation-delay: 0ms"></span>
          <span class="thinking-dot animate-bounce" style="animation-delay: 150ms"></span>
          <span class="thinking-dot animate-bounce" style="animation-delay: 300ms"></span>
        </div>
        <span class="text-xs text-purple-300">{{ t('chat.thinking') }}</span>
      </div>
      <button
        @click="$emit('stop')"
        class="flex items-center gap-1 px-2 py-1 text-xs text-red-400 hover:text-red-300 hover:bg-red-500/10 rounded-lg transition-all"
        :title="t('chat.stop')"
      >
        <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 24 24">
          <rect x="6" y="6" width="12" height="12" rx="2" />
        </svg>
        {{ t('chat.stop') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';

defineEmits<{
  stop: []
}>()

const { t } = useI18n()
</script>
