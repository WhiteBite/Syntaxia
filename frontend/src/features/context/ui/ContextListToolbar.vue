<template>
  <div class="p-2 border-b border-gray-700/30 space-y-2">
    <!-- Search -->
    <div class="relative">
      <svg class="absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      <input
        :value="searchQuery"
        @input="$emit('update:searchQuery', ($event.target as HTMLInputElement).value)"
        type="text"
        :placeholder="t('context.search')"
        class="input pl-8 text-sm"
      />
    </div>
    
    <!-- Sort & Filter -->
    <div class="flex items-center gap-2 flex-wrap">
      <select 
        :value="sortBy" 
        @change="$emit('update:sortBy', ($event.target as HTMLSelectElement).value)"
        class="input text-xs py-1 px-2 w-auto"
      >
        <option value="date">{{ t('context.sortByDate') }}</option>
        <option value="name">{{ t('context.sortByName') }}</option>
        <option value="size">{{ t('context.sortBySize') }}</option>
      </select>
      
      <label class="flex items-center gap-1.5 text-xs text-gray-400 cursor-pointer">
        <input 
          type="checkbox" 
          :checked="showFavoritesOnly"
          @change="$emit('update:showFavoritesOnly', ($event.target as HTMLInputElement).checked)"
          class="rounded border-gray-600 bg-gray-700 text-purple-500 focus:ring-purple-500 w-3 h-3" 
        />
        <svg class="w-3 h-3 text-yellow-400" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
        </svg>
      </label>
      
      <div class="flex-1"></div>
      
      <!-- Multi-select actions -->
      <template v-if="selectedCount > 0">
        <span class="text-xs text-gray-400">{{ selectedCount }} {{ t('files.selected') }}</span>
        <button 
          @click="$emit('copySelected')" 
          class="text-xs text-indigo-400 hover:text-indigo-300"
          :title="t('context.shortcutCopy')"
        >
          {{ t('context.copyToClipboard') }}
        </button>
        <button 
          v-if="selectedCount >= 2"
          @click="$emit('mergeSelected')" 
          class="text-xs text-purple-400 hover:text-purple-300"
          :title="t('context.shortcutMerge')"
        >
          {{ t('context.mergeSelected') }}
        </button>
        <button 
          @click="$emit('deleteSelected')" 
          class="text-xs text-red-400 hover:text-red-300" 
          :title="t('context.deleteAction')"
        >
          {{ t('context.deleteSelected') }}
        </button>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';

const { t } = useI18n()

defineProps<{
  searchQuery: string
  sortBy: string
  showFavoritesOnly: boolean
  selectedCount: number
}>()

defineEmits<{
  (e: 'update:searchQuery', value: string): void
  (e: 'update:sortBy', value: string): void
  (e: 'update:showFavoritesOnly', value: boolean): void
  (e: 'copySelected'): void
  (e: 'mergeSelected'): void
  (e: 'deleteSelected'): void
}>()
</script>
