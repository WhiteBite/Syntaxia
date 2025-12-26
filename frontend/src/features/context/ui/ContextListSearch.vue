<template>
  <div class="context-toolbar">
    <!-- Search -->
    <div class="context-search">
      <svg class="context-search-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      <input 
        :value="searchQuery" 
        @input="$emit('update:searchQuery', ($event.target as HTMLInputElement).value)"
        type="text" 
        :placeholder="t('context.search')" 
        class="context-search-input"
      />
      <!-- Favorites toggle inside search -->
      <button 
        @click="$emit('update:showFavoritesOnly', !showFavoritesOnly)"
        class="context-filter-btn"
        :class="{ active: showFavoritesOnly }"
        :title="t('context.showFavorites')"
      >
        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
        </svg>
      </button>
      <!-- Sort toggle -->
      <button 
        @click="$emit('cycle-sort')"
        class="context-filter-btn"
        :title="sortLabel"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
        </svg>
      </button>
    </div>
    
    <!-- Bulk Actions (when selected) -->
    <div v-if="selectedCount > 0" class="context-bulk-actions">
      <span class="context-bulk-count">{{ selectedCount }}</span>
      <button @click="$emit('copy-selected')" class="context-bulk-btn" :title="t('context.copyToClipboard')">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
      </button>
      <button @click="$emit('delete-selected')" class="context-bulk-btn context-bulk-btn--danger" :title="t('context.deleteSelected')">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'

const { t } = useI18n()

defineProps<{
  searchQuery: string
  showFavoritesOnly: boolean
  sortLabel: string
  selectedCount: number
}>()

defineEmits<{
  (e: 'update:searchQuery', value: string): void
  (e: 'update:showFavoritesOnly', value: boolean): void
  (e: 'cycle-sort'): void
  (e: 'copy-selected'): void
  (e: 'delete-selected'): void
}>()
</script>

<style scoped>
.context-toolbar {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.context-search {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 0 10px;
  height: 32px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  transition: all 0.15s ease-out;
}

.context-search:focus-within {
  border-color: #8b5cf6;
  background: rgba(139, 92, 246, 0.05);
}

.context-search-icon {
  width: 14px;
  height: 14px;
  color: #6b7280;
  flex-shrink: 0;
}

.context-search-input {
  flex: 1;
  background: none;
  border: none;
  outline: none;
  font-size: 12px;
  color: #e5e7eb;
}

.context-search-input::placeholder {
  color: #4b5563;
}

.context-filter-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 6px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.context-filter-btn:hover {
  color: #e5e7eb;
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.12);
}

.context-filter-btn.active {
  color: #facc15;
  background: rgba(250, 204, 21, 0.1);
  border-color: rgba(250, 204, 21, 0.2);
}

.context-bulk-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.context-bulk-count {
  font-size: 11px;
  font-weight: 600;
  color: #8b5cf6;
  padding: 2px 8px;
  background: rgba(139, 92, 246, 0.15);
  border-radius: 4px;
}

.context-bulk-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  background: rgba(255, 255, 255, 0.05);
  border: none;
  border-radius: 6px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.context-bulk-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.context-bulk-btn--danger:hover {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
}
</style>
