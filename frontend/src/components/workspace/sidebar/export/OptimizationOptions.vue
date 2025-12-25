<template>
  <div class="context-stats space-y-3">
    <SettingsToggle
      :model-value="excludeTests"
      :label="t('context.excludeTests')"
      :hint="t('export.hint.excludeTests')"
      :savings="testSavings"
      @update:model-value="$emit('update:excludeTests', $event)"
    />
    <SettingsToggle
      :model-value="stripLicense"
      :label="t('context.stripLicense')"
      :hint="t('export.hint.stripLicense')"
      :savings="licenseSavings"
      @update:model-value="$emit('update:stripLicense', $event)"
    />
    <SettingsToggle
      :model-value="compactDataFiles"
      :label="t('context.compactDataFiles')"
      :hint="t('export.hint.compactDataFiles')"
      :savings="dataSavings"
      @update:model-value="$emit('update:compactDataFiles', $event)"
    />
    <SettingsToggle
      :model-value="trimWhitespace"
      :label="t('context.trimWhitespace')"
      :hint="t('export.hint.trimWhitespace')"
      @update:model-value="$emit('update:trimWhitespace', $event)"
    />
    <SettingsToggle
      :model-value="collapseEmptyLines"
      :label="t('context.collapseEmptyLines')"
      :hint="t('export.hint.collapseEmptyLines')"
      @update:model-value="$emit('update:collapseEmptyLines', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { useContextStore } from '@/features/context/model/context.store'
import { useI18n } from '@/composables/useI18n'
import { computed } from 'vue'
import SettingsToggle from './SettingsToggle.vue'

const { t } = useI18n()
const contextStore = useContextStore()

const props = defineProps<{
  excludeTests: boolean
  stripLicense: boolean
  compactDataFiles: boolean
  trimWhitespace: boolean
  collapseEmptyLines: boolean
}>()

defineEmits<{
  'update:excludeTests': [value: boolean]
  'update:stripLicense': [value: boolean]
  'update:compactDataFiles': [value: boolean]
  'update:trimWhitespace': [value: boolean]
  'update:collapseEmptyLines': [value: boolean]
}>()

// Calculate estimated savings based on context stats
const totalTokens = computed(() => contextStore.tokenCount || 0)
const fileCount = computed(() => contextStore.fileCount || 0)

const testSavings = computed(() => {
  if (!fileCount.value) return undefined
  const estimatedTestFiles = Math.round(fileCount.value * 0.15)
  const tokens = estimatedTestFiles * 500
  return tokens > 0 ? formatSavings(tokens) : undefined
})

const licenseSavings = computed(() => {
  if (!fileCount.value) return undefined
  const tokens = Math.round(fileCount.value * 0.1 * 200)
  return tokens > 100 ? formatSavings(tokens) : undefined
})

const dataSavings = computed(() => {
  if (!totalTokens.value) return undefined
  const tokens = Math.round(totalTokens.value * 0.05 * 0.3)
  return tokens > 100 ? formatSavings(tokens) : undefined
})

function formatSavings(tokens: number): string {
  if (tokens >= 1000) {
    return `-${(tokens / 1000).toFixed(1)}K`
  }
  return `-${tokens}`
}
</script>
