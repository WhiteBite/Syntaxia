<template>
  <div class="filters-bar">
    <div class="filter-groups">
      <!-- Types Dropdown -->
      <div class="filter-dropdown" ref="typesDropdownRef">
        <button
          @click="toggleDropdown('types')"
          :class="['filter-trigger', { 'filter-trigger-active': hasActiveTypeFilters }]"
          :aria-label="t('quickFilters.filterByType')"
          :aria-expanded="openDropdown === 'types'"
          aria-haspopup="listbox"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
          </svg>
          <span>{{ t('quickFilters.types') }}</span>
          <span v-if="activeTypeFilters.length > 0" class="filter-badge">{{ activeTypeFilters.length }}</span>
          <ChevronIcon :open="openDropdown === 'types'" />
        </button>
        <FilterDropdownMenu
          :is-open="openDropdown === 'types'"
          :title="t('quickFilters.fileTypes')"
          :style="dropdownStyle"
          :has-active-filters="activeTypeFilters.length > 0"
          :hint="true"
          @clear="clearTypeFilters"
        >
          <FilterDropdownItem
            v-for="filter in typeFilters"
            :key="filter.id"
            :label="filter.label"
            :icon-type="'category'"
            :category="filter.category"
            :count="getFilterCount(filter)"
            :percentage="getFilterPercentage(filter)"
            :active="isFilterActive(filter.id)"
            :excluded="isFilterExcluded(filter.id)"
            @toggle="toggleFilter(filter, $event)"
          >
            <template #icon>
              <component :is="getFilterIcon(filter.category)" class="w-3.5 h-3.5" />
            </template>
          </FilterDropdownItem>
        </FilterDropdownMenu>
      </div>

      <!-- Languages Dropdown -->
      <div v-if="languageFilters.length > 0" class="filter-dropdown" ref="langsDropdownRef">
        <button
          @click="toggleDropdown('langs')"
          :class="['filter-trigger', { 'filter-trigger-active': hasActiveLanguageFilters }]"
          :aria-label="t('quickFilters.filterByLanguage')"
          :aria-expanded="openDropdown === 'langs'"
          aria-haspopup="listbox"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
          </svg>
          <span>{{ t('quickFilters.languages') }}</span>
          <span class="filter-badge">{{ languageFilters.length }}</span>
          <ChevronIcon :open="openDropdown === 'langs'" />
        </button>
        <FilterDropdownMenu
          :is-open="openDropdown === 'langs'"
          :title="t('quickFilters.projectLanguages')"
          :style="dropdownStyle"
          :has-active-filters="activeLanguageFilters.length > 0"
          :footer="t('quickFilters.autoDetected') + ' ✨'"
          @clear="clearLanguageFilters"
        >
          <FilterDropdownItem
            v-for="filter in languageFilters"
            :key="filter.id"
            :label="filter.language"
            :icon="filter.icon"
            :count="filter.fileCount"
            :percentage="filter.percentage"
            :active="isFilterActive(filter.id)"
            :primary="filter.primary"
            bar-class="bar-lang"
            @toggle="toggleFilter(filter, $event)"
          />
        </FilterDropdownMenu>
      </div>

      <!-- Smart Filters Dropdown -->
      <div v-if="smartFilters.length > 0" class="filter-dropdown" ref="smartDropdownRef">
        <button
          @click="toggleDropdown('smart')"
          :class="['filter-trigger filter-trigger-smart', { 'filter-trigger-active': hasActiveSmartFilters }]"
          :aria-label="t('quickFilters.smartFilters')"
          :aria-expanded="openDropdown === 'smart'"
          aria-haspopup="listbox"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
          </svg>
          <span>{{ t('quickFilters.smart') }}</span>
          <span class="filter-badge filter-badge-smart">{{ smartFilters.length }}</span>
          <ChevronIcon :open="openDropdown === 'smart'" />
        </button>
        <FilterDropdownMenu
          :is-open="openDropdown === 'smart'"
          :title="t('quickFilters.smartFilters')"
          :style="dropdownStyle"
          :has-active-filters="activeSmartFilters.length > 0"
          :footer="t('quickFilters.basedOnFramework') + ' 🧠'"
          menu-class="filter-menu-smart"
          header-class="filter-header-smart"
          footer-class="filter-footer-smart"
          @clear="clearSmartFilters"
        >
          <FilterDropdownItem
            v-for="filter in smartFilters"
            :key="filter.id"
            :label="filter.label"
            :icon="filter.icon"
            :subtitle="filter.framework"
            :count="getFilterCount(filter)"
            :percentage="getFilterPercentage(filter)"
            :active="isFilterActive(filter.id)"
            bar-class="bar-smart"
            @toggle="toggleFilter(filter, $event)"
          />
        </FilterDropdownMenu>
      </div>
    </div>

    <!-- Active & Excluded Chips -->
    <div v-if="allActiveFilters.length > 0 || allExcludedFilters.length > 0" class="active-filters">
      <TransitionGroup name="chip">
        <FilterChip
          v-for="filter in allActiveFilters"
          :key="filter.id"
          :label="filter.shortLabel || filter.label"
          :icon="filter.icon"
          :category="filter.category"
          @remove="removeFilter(filter)"
        />
        <FilterChip
          v-for="filter in allExcludedFilters"
          :key="'ex-' + filter.id"
          :label="filter.shortLabel || filter.label"
          :excluded="true"
          @remove="removeFilter(filter)"
        />
      </TransitionGroup>
    </div>

    <div class="flex-1"></div>

    <!-- Actions -->
    <div class="filter-actions">
      <button
        v-if="allActiveFilters.length > 0"
        @click="clearAllFilters"
        class="filter-clear-btn"
        :title="t('quickFilters.clearAll')"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
        <span>{{ t('quickFilters.clear') }}</span>
      </button>
      <button @click="showSettingsModal = true" class="filter-settings-btn" :title="t('quickFilters.settings')">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
        </svg>
      </button>
    </div>
  </div>

  <FilterSettingsModal
    :is-open="showSettingsModal"
    :filters="editableFilters"
    :get-count="getFilterCount"
    @close="showSettingsModal = false"
    @reset="resetFilters"
    @update-extensions="updateFilterExtensions"
  />
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { h, onMounted, onUnmounted } from 'vue'
import { useQuickFilters } from '../composables/useQuickFilters'
import FilterChip from './FilterChip.vue'
import FilterDropdownItem from './FilterDropdownItem.vue'
import FilterDropdownMenu from './FilterDropdownMenu.vue'
import FilterSettingsModal from './FilterSettingsModal.vue'

const { t } = useI18n()
const {
  openDropdown, dropdownStyle, showSettingsModal,
  typeFilters, languageFilters, smartFilters, editableFilters,
  activeTypeFilters, activeLanguageFilters, activeSmartFilters,
  allActiveFilters, allExcludedFilters,
  hasActiveTypeFilters, hasActiveLanguageFilters, hasActiveSmartFilters,
  toggleDropdown, toggleFilter, removeFilter,
  isFilterActive, isFilterExcluded,
  clearTypeFilters, clearLanguageFilters, clearSmartFilters, clearAllFilters,
  getFilterCount, getFilterPercentage,
  updateFilterExtensions, resetFilters,
  setupEventListeners, cleanupEventListeners,
} = useQuickFilters()

// Icon components
const ChevronIcon = (props: { open: boolean }) => h('svg', {
  class: ['w-3 h-3 transition-transform', { 'rotate-180': props.open }],
  fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24'
}, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M19 9l-7 7-7-7' })])

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

onMounted(() => setupEventListeners())
onUnmounted(() => cleanupEventListeners())
</script>

<style scoped>
.filters-bar {
  @apply flex items-center gap-2 px-2 py-1.5 flex-wrap;
  border-bottom: 1px solid var(--border-default);
  background: var(--bg-1);
  min-height: 40px;
}
.filter-groups { @apply flex items-center gap-1.5 flex-wrap; }
.filter-actions { @apply flex items-center gap-1; }
.active-filters { @apply flex items-center gap-1 flex-wrap; }
.filter-dropdown { @apply relative; }

.filter-trigger {
  @apply flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-full;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  color: var(--text-secondary);
  transition: all 150ms ease-out;
}
.filter-trigger:hover {
  background: var(--bg-3);
  border-color: var(--border-strong);
  color: var(--text-primary);
}
.filter-trigger-active {
  background: var(--accent-indigo-bg);
  border-color: var(--accent-indigo-border);
  color: white;
}
.filter-trigger-smart {
  background: linear-gradient(135deg, var(--bg-2), rgba(16, 185, 129, 0.1));
}

.filter-badge {
  @apply px-1.5 py-0.5 text-[10px] font-bold rounded-full min-w-[18px] text-center;
  background: var(--accent-indigo);
  color: white;
}
.filter-badge-smart { background: var(--accent-emerald); }

.filter-clear-btn {
  @apply flex items-center gap-1 px-2 py-1 text-xs font-medium rounded-lg;
  color: var(--color-danger);
  background: var(--color-danger-soft);
  border: 1px solid var(--color-danger-border);
  transition: all 150ms ease-out;
}
.filter-clear-btn:hover { background: rgba(248, 113, 113, 0.25); }

.filter-settings-btn {
  @apply p-1.5 rounded-lg;
  color: var(--text-muted);
  transition: all 150ms ease-out;
}
.filter-settings-btn:hover {
  color: var(--text-primary);
  background: var(--bg-2);
}

/* Smart filter menu variants */
:deep(.filter-menu-smart) { border-color: var(--accent-emerald-border); }
:deep(.filter-header-smart) { background: linear-gradient(135deg, var(--bg-2), rgba(16, 185, 129, 0.1)); }
:deep(.filter-footer-smart) { background: linear-gradient(135deg, var(--bg-2), rgba(16, 185, 129, 0.1)); }

/* Chip animations */
.chip-enter-active, .chip-leave-active { transition: all 200ms ease-out; }
.chip-enter-from { opacity: 0; transform: scale(0.8); }
.chip-leave-to { opacity: 0; transform: scale(0.8) translateX(-10px); }
.chip-move { transition: transform 200ms ease-out; }
</style>
