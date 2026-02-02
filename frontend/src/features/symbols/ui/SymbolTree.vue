<template>
  <div class="symbol-tree">
    <!-- Empty State -->
    <BaseEmptyState
      v-if="groups.length === 0"
      :title="t('symbols.noSymbols')"
      size="md"
    >
      <template #icon>
        <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" 
            d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
      </template>
    </BaseEmptyState>

    <!-- File Groups -->
    <div v-else class="tree-content">
      <div
        v-for="group in groups"
        :key="group.file"
        class="file-group"
      >
        <!-- File Header -->
        <div
          class="file-header"
          @click="toggleExpand(group.file)"
        >
          <svg
            class="expand-icon"
            :class="{ 'expand-icon--expanded': group.expanded }"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
          
          <svg class="file-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          
          <span class="file-name" :title="group.file">{{ getFileName(group.file) }}</span>
          <span class="symbol-count">{{ group.symbols.length }}</span>
        </div>

        <!-- Symbols List -->
        <div v-if="group.expanded" class="symbols-list">
          <SymbolItem
            v-for="symbol in group.symbols"
            :key="`${symbol.file}:${symbol.line}:${symbol.name}`"
            :symbol="symbol"
            :is-nested="!!symbol.parent"
            @click="handleSymbolClick"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { BaseEmptyState } from '@/components/ui'
import type { Symbol, SymbolGroup } from '../types'
import SymbolItem from './SymbolItem.vue'

defineProps<{
  groups: SymbolGroup[]
}>()

const emit = defineEmits<{
  (e: 'toggle-expand', file: string): void
  (e: 'symbol-click', symbol: Symbol): void
}>()

const { t } = useI18n()

function getFileName(filePath: string): string {
  const parts = filePath.replace(/\\/g, '/').split('/')
  return parts[parts.length - 1] || filePath
}

function toggleExpand(file: string) {
  emit('toggle-expand', file)
}

function handleSymbolClick(symbol: Symbol) {
  emit('symbol-click', symbol)
}
</script>

<style scoped>
.symbol-tree {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.tree-content {
  padding: 8px;
}

.file-group {
  margin-bottom: 4px;
}

.file-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.15s ease-out;
}

.file-header:hover {
  background: rgba(255, 255, 255, 0.06);
}

.expand-icon {
  flex-shrink: 0;
  width: 14px;
  height: 14px;
  color: #6b7280;
  transition: transform 0.15s ease-out;
}

.expand-icon--expanded {
  transform: rotate(90deg);
}

.file-icon {
  flex-shrink: 0;
  width: 16px;
  height: 16px;
  color: #6b7280;
}

.file-name {
  flex: 1;
  min-width: 0;
  font-size: 13px;
  font-weight: 500;
  color: #e5e7eb;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.symbol-count {
  flex-shrink: 0;
  font-size: 11px;
  padding: 2px 6px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  color: #9ca3af;
}

.symbols-list {
  padding-left: 20px;
  border-left: 1px solid rgba(255, 255, 255, 0.06);
  margin-left: 14px;
  margin-top: 4px;
}
</style>
