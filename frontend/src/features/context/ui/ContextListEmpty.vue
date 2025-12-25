<template>
  <div class="flex items-center justify-center flex-1 p-4">
    <div class="text-center max-w-xs">
      <div class="w-16 h-16 mx-auto mb-4 bg-gray-800/50 rounded-2xl flex items-center justify-center border border-gray-700/30">
        <svg class="w-8 h-8 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
        </svg>
      </div>
      <p class="text-sm font-medium text-gray-400 mb-1">{{ t('context.noSaved') }}</p>
      <p class="text-xs text-gray-400">{{ t('context.buildToSee') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';

const { t } = useI18n()
</script>
