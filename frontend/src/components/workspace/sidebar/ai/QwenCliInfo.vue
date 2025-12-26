<template>
  <div class="info-box-purple">
    <div class="flex items-start gap-2">
      <Info class="w-4 h-4 text-purple-400 mt-0.5 flex-shrink-0" />
      <div>
        <p class="text-xs text-purple-300 font-medium mb-1">{{ t('settings.qwenCliTitle') }}</p>
        <p class="text-2xs text-gray-400">{{ t('settings.qwenCliDescription') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';
import { Info } from 'lucide-vue-next';

const { t } = useI18n()
</script>
