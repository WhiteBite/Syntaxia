<template>
  <BaseTabs :model-value="modelValue" @update:model-value="handleChange" class="px-2 pb-2">
    <BaseTab name="local" :label="t('git.localGit')">
      <template #label>
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
        </svg>
        {{ t('git.localGit') }}
      </template>
    </BaseTab>
    <BaseTab name="remote" :label="t('git.remoteUrl')">
      <template #label>
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
        </svg>
        {{ t('git.remoteUrl') }}
      </template>
    </BaseTab>
  </BaseTabs>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { BaseTabs, BaseTab } from '@/components/ui'

export type SourceType = 'local' | 'remote'

defineProps<{
  modelValue: SourceType
}>()

const emit = defineEmits<{
  change: [value: SourceType]
}>()

const { t } = useI18n()

function handleChange(value: string) {
  emit('change', value as SourceType)
}
</script>
