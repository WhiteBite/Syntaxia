<template>
  <div class="panel-header-unified">
    <div class="panel-header-unified-title">
      <div class="section-icon section-icon-indigo">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
      </div>
      <span>{{ t('context.preview') }}</span>
    </div>

    <div v-if="hasContext" class="flex items-center gap-1.5 ml-2 overflow-x-auto">
      <!-- Search toggle -->
      <button 
        @click="$emit('toggle-search')" 
        class="icon-btn"
        :class="{ 'text-indigo-400 bg-indigo-500/10': showSearch }" 
        :title="t('context.search')"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </button>

      <!-- Stats (clickable badges - open popover) -->
      <button 
        @click="$emit('show-stats')"
        class="hidden xl:flex items-center gap-1.5 flex-shrink-0 stats-chips-btn"
        :title="t('stats.clickToExpand')"
      >
        <span class="chip-unified chip-unified-accent">{{ fileCount }} {{ t('context.files') }}</span>
        <span class="chip-unified chip-unified-accent">{{ lineCount }} {{ t('context.lines') }}</span>
        <span class="chip-unified chip-unified-accent">{{ tokenCount }} {{ t('context.tokens') }}</span>
      </button>

      <!-- Format Selector Dropdown -->
      <FormatDropdown
        :model-value="outputFormat"
        @update:model-value="$emit('format-change', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import FormatDropdown from '@/components/workspace/sidebar/export/FormatDropdown.vue'
import { useI18n } from '@/composables/useI18n'
import type { OutputFormat } from '@/stores/settings.store'

defineProps<{
  hasContext: boolean
  showSearch: boolean
  fileCount: number
  lineCount: number
  tokenCount: number
  outputFormat: OutputFormat
}>()

defineEmits<{
  'toggle-search': []
  'show-stats': []
  'format-change': [format: OutputFormat]
}>()

const { t } = useI18n()
</script>

<style scoped>
.stats-chips-btn {
  background: none;
  border: none;
  padding: 2px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.stats-chips-btn:hover {
  background: rgba(139, 92, 246, 0.1);
}

.stats-chips-btn:hover .chip-unified {
  border-color: rgba(139, 92, 246, 0.4);
}
</style>
