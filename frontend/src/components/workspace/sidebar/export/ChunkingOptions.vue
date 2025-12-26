<template>
  <div class="space-y-3">
    <SettingsToggle
      :model-value="enableAutoSplit"
      :label="t('export.autoSplit')"
      :hint="t('export.hint.autoSplit')"
      @update:model-value="$emit('update:enableAutoSplit', $event)"
    />

    <div v-if="enableAutoSplit" class="space-y-3">
      <!-- Settings -->
      <div class="chunking-settings">
        <div>
          <label class="text-2xs text-gray-400 mb-1 block">{{ t('export.tokensPerChunk') }}</label>
          <input 
            type="number" 
            :value="maxTokensPerChunk"
            @input="$emit('update:maxTokensPerChunk', Number(($event.target as HTMLInputElement).value))" 
            min="500"
            max="32000" 
            class="input input-xs" 
          />
        </div>
        <div>
          <label class="text-2xs text-gray-400 mb-1 block">{{ t('export.strategy') }}</label>
          <select 
            :value="splitStrategy"
            @change="$emit('update:splitStrategy', ($event.target as HTMLSelectElement).value as SplitStrategy)"
            class="input input-xs"
          >
            <option value="smart">{{ t('export.strategySmart') }}</option>
            <option value="file">{{ t('export.strategyFile') }}</option>
            <option value="token">{{ t('export.strategyToken') }}</option>
          </select>
        </div>
        <!-- Compact chunk estimate -->
        <div v-if="hasContext && estimatedChunks > 0" class="chunk-estimate">
          → {{ estimatedChunks }} {{ t('export.chunksEstimated') }} (~{{ avgTokensPerChunk }} {{ t('export.each') }})
        </div>
      </div>

      <!-- Chunk Copy UI -->
      <ChunkCopyBar />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useContextStore } from '@/features/context/model/context.store'
import { useI18n } from '@/composables/useI18n'
import { computed } from 'vue'
import ChunkCopyBar from './ChunkCopyBar.vue'
import SettingsToggle from './SettingsToggle.vue'

const { t } = useI18n()
const contextStore = useContextStore()

export type SplitStrategy = 'smart' | 'file' | 'token'

const props = defineProps<{
  enableAutoSplit: boolean
  maxTokensPerChunk: number
  splitStrategy: SplitStrategy
}>()

defineEmits<{
  'update:enableAutoSplit': [value: boolean]
  'update:maxTokensPerChunk': [value: number]
  'update:splitStrategy': [value: SplitStrategy]
}>()

const hasContext = computed(() => contextStore.hasContext)
const totalTokens = computed(() => contextStore.stats?.tokens || 0)

const estimatedChunks = computed(() => {
  if (!totalTokens.value || !props.maxTokensPerChunk) return 0
  return Math.ceil(totalTokens.value / props.maxTokensPerChunk)
})

const avgTokensPerChunk = computed(() => {
  if (!estimatedChunks.value) return '0'
  const avg = Math.round(totalTokens.value / estimatedChunks.value)
  return avg >= 1000 ? `${(avg / 1000).toFixed(1)}K` : String(avg)
})
</script>

<style scoped>
.chunking-settings {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding-left: 0.5rem;
  border-left: 2px solid var(--accent-indigo-border);
}

.chunk-estimate {
  font-size: 0.6875rem;
  color: var(--text-muted);
  padding: 0.25rem 0;
}
</style>
