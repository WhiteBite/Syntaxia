<template>
  <div class="flex gap-1 px-2 pb-2">
    <button 
      @click="$emit('change', 'local')" 
      :class="['tab-btn', modelValue === 'local' ? 'tab-btn-active tab-btn-active-indigo' : 'tab-btn-inactive']"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
      </svg>
      {{ t('git.localGit') }}
    </button>
    <button 
      @click="$emit('change', 'remote')" 
      :class="['tab-btn', modelValue === 'remote' ? 'tab-btn-active tab-btn-active-purple' : 'tab-btn-inactive']"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
      </svg>
      {{ t('git.remoteUrl') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'

export type SourceType = 'local' | 'remote'

defineProps<{
  modelValue: SourceType
}>()

defineEmits<{
  change: [value: SourceType]
}>()

const { t } = useI18n()
</script>
