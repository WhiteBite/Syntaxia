<template>
  <div class="filter-popover p-2 w-[340px]">
    <div class="flex items-center justify-between mb-3 px-1">
      <h3 class="text-xs font-bold text-gray-400 uppercase tracking-wider">{{ t('quickFilters.title') }}</h3>
      <BaseButton
        v-if="allActiveFilters.length > 0"
        variant="ghost"
        size="xs"
        @click="clearAllFilters"
        class="text-xs h-6"
      >
        {{ t('quickFilters.clearAll') }}
      </BaseButton>
    </div>

    <!-- Active Filters Chips -->
    <div v-if="allActiveFilters.length > 0" class="flex flex-wrap gap-1.5 mb-3 px-1">
      <BaseChip
        v-for="filter in allActiveFilters"
        :key="filter.id"
        :variant="getCategoryVariant(filter.category)"
        removable
        size="sm"
        @remove="removeFilter(filter)"
      >
        <template v-if="filter.icon" #icon>
          <span>{{ filter.icon }}</span>
        </template>
        {{ filter.shortLabel || filter.label }}
      </BaseChip>
    </div>

    <!-- Filter Groups -->
    <div class="space-y-4">
      <!-- File Types -->
      <div>
        <div class="text-[10px] font-semibold text-gray-500 uppercase tracking-wider px-1 mb-1.5">{{ t('quickFilters.fileTypes') }}</div>
        <div class="grid grid-cols-2 gap-1">
          <button
            v-for="filter in typeFilters"
            :key="filter.id"
            @click="toggleFilter(filter, $event)"
            class="filter-btn"
            :class="{ active: isFilterActive(filter.id), excluded: isFilterExcluded(filter.id) }"
            :title="filter.label"
          >
            <div class="flex items-center gap-1.5 min-w-0">
              <component :is="getFilterIcon(filter.category)" class="w-3.5 h-3.5 opacity-70" />
              <span class="truncate">{{ filter.shortLabel || filter.label }}</span>
            </div>
            <span class="count-badge">{{ getFilterCount(filter) }}</span>
          </button>
        </div>
      </div>

      <!-- Languages -->
      <div v-if="languageFilters.length > 0">
        <div class="text-[10px] font-semibold text-gray-500 uppercase tracking-wider px-1 mb-1.5">{{ t('quickFilters.projectLanguages') }}</div>
        <div class="grid grid-cols-2 gap-1">
          <button
            v-for="filter in languageFilters"
            :key="filter.id"
            @click="toggleFilter(filter, $event)"
            class="filter-btn"
            :class="{ active: isFilterActive(filter.id), excluded: isFilterExcluded(filter.id) }"
          >
            <div class="flex items-center gap-1.5 min-w-0">
              <span class="text-xs">{{ filter.icon }}</span>
              <span class="truncate">{{ filter.language }}</span>
            </div>
            <span class="count-badge">{{ filter.fileCount }}</span>
          </button>
        </div>
      </div>

      <!-- Smart Filters -->
      <div v-if="smartFilters.length > 0">
        <div class="text-[10px] font-semibold text-gray-500 uppercase tracking-wider px-1 mb-1.5">{{ t('quickFilters.smartFilters') }}</div>
        <div class="flex flex-col gap-1">
          <button
            v-for="filter in smartFilters"
            :key="filter.id"
            @click="toggleFilter(filter, $event)"
            class="filter-btn w-full"
            :class="{ active: isFilterActive(filter.id), excluded: isFilterExcluded(filter.id) }"
          >
            <div class="flex items-center gap-1.5 min-w-0">
              <span class="text-xs">{{ filter.icon }}</span>
              <div class="flex flex-col items-start truncate">
                <span class="truncate">{{ filter.label }}</span>
                <span class="text-[9px] text-gray-500">{{ filter.framework }}</span>
              </div>
            </div>
            <span class="count-badge">{{ getFilterCount(filter) }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Hints -->
    <div class="mt-3 pt-2 border-t border-gray-700/50 text-[10px] text-gray-500 flex justify-between px-1">
      <span>Click to toggle</span>
      <span>Shift+Click to exclude</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { h } from 'vue'
import { useQuickFilters } from '../composables/useQuickFilters'
import { BaseChip, BaseButton } from '@/components/ui'
import { useI18n } from '@/composables/useI18n'

const { t } = useI18n()

const {
  typeFilters, languageFilters, smartFilters,
  allActiveFilters,
  toggleFilter, removeFilter,
  isFilterActive, isFilterExcluded,
  clearAllFilters,
  getFilterCount
} = useQuickFilters()

const iconPaths: Record<string, string> = {
  code: 'M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4',
  test: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
  config: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z',
  docs: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
  styles: 'M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01'
}

function getFilterIcon(category: string) {
  return () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
    h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: iconPaths[category] || iconPaths.code })
  ])
}

function getCategoryVariant(category?: string): 'default' | 'primary' | 'success' | 'warning' | 'danger' {
  switch (category) {
    case 'code': return 'primary'
    case 'test': return 'success'
    case 'config': return 'warning'
    case 'docs': return 'primary'
    case 'styles': return 'primary'
    default: return 'default'
  }
}
</script>

<style scoped>
.filter-btn {
  @apply flex items-center justify-between px-2 py-1.5 rounded text-xs text-gray-400 bg-gray-800/50 border border-transparent transition-all;
}

.filter-btn:hover {
  @apply bg-gray-700/50 text-gray-200 border-gray-600/50;
}

.filter-btn.active {
  @apply bg-indigo-500/10 text-indigo-300 border-indigo-500/30;
}

.filter-btn.excluded {
  @apply bg-red-500/10 text-red-300 border-red-500/30 line-through decoration-red-500/50;
}

.count-badge {
  @apply text-[10px] font-mono opacity-50 ml-2;
}
</style>
