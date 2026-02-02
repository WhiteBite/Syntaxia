<template>
  <div class="symbol-browser">
    <!-- Header -->
    <div class="panel-header">
      <div class="panel-header-unified">
        <div class="panel-header-unified-title">
          <div class="panel-header-unified-icon panel-header-unified-icon-violet">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </div>
          <h2>{{ t('symbols.title') }}</h2>
        </div>

        <div class="flex items-center gap-1">
          <!-- Stats Badge -->
          <div class="stats-badge" :title="statsTooltip">
            <span class="stats-count">{{ store.filteredCount }}</span>
            <span class="stats-label">{{ t('symbols.symbols') }}</span>
          </div>

          <!-- Toolbar Buttons -->
          <BaseButton
            variant="ghost"
            size="sm"
            icon-only
            :title="t('symbols.expandAll')"
            @click="store.expandAll"
          >
            <template #icon>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </template>
          </BaseButton>
          <BaseButton
            variant="ghost"
            size="sm"
            icon-only
            :title="t('symbols.collapseAll')"
            @click="store.collapseAll"
          >
            <template #icon>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
              </svg>
            </template>
          </BaseButton>
          <BaseButton
            variant="ghost"
            size="sm"
            icon-only
            :title="t('symbols.refresh')"
            :disabled="store.isLoading"
            :loading="store.isLoading"
            @click="handleRefresh"
          >
            <template #icon>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </template>
          </BaseButton>
        </div>
      </div>
    </div>

    <!-- Search -->
    <SymbolSearch
      v-model="searchQuery"
      :selected-kinds="store.filter.kinds"
      :kind-stats="store.kindStats"
      @search="handleSearch"
      @toggle-kind="store.toggleKindFilter"
      @clear="handleClearSearch"
    />

    <!-- Content -->
    <div class="panel-content">
      <!-- Loading State -->
      <div v-if="store.isLoading && !store.hasSymbols" class="loading-state">
        <svg class="loading-spinner" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" opacity="0.25"/>
          <path d="M12 2a10 10 0 0110 10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        <span class="loading-text">{{ t('symbols.loading') }}</span>
      </div>

      <!-- Error State -->
      <div v-else-if="store.error" class="error-state">
        <svg class="error-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <p class="error-text">{{ store.error }}</p>
        <BaseButton variant="ghost" @click="handleRefresh">
          {{ t('symbols.retry') }}
        </BaseButton>
      </div>

      <!-- No Project State -->
      <div v-else-if="!projectStore.hasProject" class="empty-state">
        <svg class="empty-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" 
            d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
        </svg>
        <p class="empty-text">{{ t('symbols.selectProject') }}</p>
      </div>

      <!-- Symbol Tree -->
      <SymbolTree
        v-else
        :groups="store.symbolsByFile"
        @toggle-expand="store.toggleFileExpanded"
        @symbol-click="handleSymbolClick"
      />
    </div>

    <!-- Footer -->
    <div v-if="store.hasSymbols" class="panel-footer">
      <span class="footer-stats">
        {{ store.fileCount }} {{ t('symbols.files') }} · {{ store.filteredCount }} {{ t('symbols.symbols') }}
      </span>
      <BaseButton
        v-if="store.hasFilter"
        variant="ghost"
        size="sm"
        @click="handleClearFilter"
      >
        {{ t('symbols.clearFilter') }}
      </BaseButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '@/components/ui/BaseButton.vue'
import { useI18n } from '@/composables/useI18n'
import { useProjectStore } from '@/stores/project.store'
import { computed, onMounted, ref, watch } from 'vue'
import { useSymbolsStore } from '../model/symbols.store'
import type { Symbol } from '../types'
import SymbolSearch from './SymbolSearch.vue'
import SymbolTree from './SymbolTree.vue'

const { t } = useI18n()
const store = useSymbolsStore()
const projectStore = useProjectStore()

const searchQuery = ref('')

const statsTooltip = computed(() => {
  const stats = store.kindStats
  const parts: string[] = []
  
  if (stats.class > 0) parts.push(`${stats.class} classes`)
  if (stats.function > 0) parts.push(`${stats.function} functions`)
  if (stats.method > 0) parts.push(`${stats.method} methods`)
  if (stats.interface > 0) parts.push(`${stats.interface} interfaces`)
  
  return parts.join(', ') || 'No symbols'
})

// Watch for project changes
watch(() => projectStore.currentPath, (newPath, oldPath) => {
  if (newPath && newPath !== oldPath) {
    store.reset()
    store.loadSymbols()
  }
})

onMounted(() => {
  if (projectStore.hasProject) {
    store.loadSymbols()
  }
})

function handleSearch(query: string) {
  if (query.length >= 2) {
    store.searchSymbols(query)
  } else if (query.length === 0) {
    store.loadSymbols()
  }
}

function handleClearSearch() {
  searchQuery.value = ''
  store.clearFilter()
  store.loadSymbols()
}

function handleClearFilter() {
  searchQuery.value = ''
  store.clearFilter()
}

function handleRefresh() {
  store.reset()
  store.loadSymbols()
}

function handleSymbolClick(symbol: Symbol) {
  // Emit event for navigation
  window.dispatchEvent(new CustomEvent('symbols:go-to-file', {
    detail: { file: symbol.file, line: symbol.line }
  }))
}
</script>

<style scoped>
.symbol-browser {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-primary, #0f1117);
  border-radius: 12px;
  overflow: hidden;
}

.panel-header {
  flex-shrink: 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.stats-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: rgba(139, 92, 246, 0.1);
  border-radius: 6px;
  font-size: 11px;
  margin-right: 8px;
}

.stats-count {
  font-weight: 600;
  color: #a78bfa;
}

.stats-label {
  color: #9ca3af;
}

.panel-content {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.loading-state,
.error-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 32px;
  text-align: center;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  color: #a78bfa;
  animation: spin 1s linear infinite;
}

.loading-text {
  margin-top: 12px;
  font-size: 13px;
  color: #9ca3af;
}

.error-icon,
.empty-icon {
  width: 48px;
  height: 48px;
  color: #4b5563;
  margin-bottom: 12px;
}

.error-icon {
  color: #ef4444;
}

.error-text,
.empty-text {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 16px;
}

.panel-footer {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(0, 0, 0, 0.2);
}

.footer-stats {
  font-size: 11px;
  color: #6b7280;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
