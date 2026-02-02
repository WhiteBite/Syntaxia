<template>
  <!-- Bottom Panel for Remote Files -->
  <div v-if="sourceType === 'remote' && remoteSelectedCount > 0" class="border-t border-gray-700 p-4">
    <div class="flex items-center justify-between mb-3">
      <span class="text-sm text-gray-300">{{ remoteSelectedCount }} {{ t('git.filesSelected') }}</span>
    </div>
    <BaseButton @click="$emit('build-remote')" :disabled="isBuilding" variant="primary" :loading="isBuilding" class="w-full">
      {{ t('git.buildContext') }} {{ remoteBranch }}
    </BaseButton>
  </div>

  <!-- Bottom Panel - Build Context (Local) -->
  <div v-if="sourceType === 'local' && localSelectedCount > 0" class="border-t border-gray-700 p-4">
    <div class="flex items-center justify-between mb-3">
      <span class="text-sm text-gray-300">{{ localSelectedCount }} {{ t('git.filesSelected') }}</span>
    </div>
    <BaseButton @click="$emit('build-local')" :disabled="isBuilding" variant="primary" :loading="isBuilding" class="w-full">
      {{ t('git.buildContext') }} {{ localRef?.slice(0, 7) || 'ref' }}
    </BaseButton>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import type { SourceType } from './GitSourceTabs.vue'
import { BaseButton } from '@/components/ui'

defineProps<{
  sourceType: SourceType
  localSelectedCount: number
  remoteSelectedCount: number
  localRef: string | null
  remoteBranch: string
  isBuilding: boolean
}>()

defineEmits<{
  'build-local': []
  'build-remote': []
}>()

const { t } = useI18n()
</script>
