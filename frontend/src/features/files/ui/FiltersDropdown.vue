<template>
  <div class="relative">
    <BaseDropdown v-model="isOpen" placement="bottom-end">
      <template #trigger>
        <slot name="trigger">
          <BaseButton 
            variant="ghost" 
            size="xs" 
            icon-only 
            :class="{ 'text-indigo-400': hasActiveFilters, 'text-gray-500 hover:text-white': !hasActiveFilters }"
            title="Filter Files"
          >
            <FilterIcon class="w-4 h-4" />
          </BaseButton>
        </slot>
      </template>

      <template #content>
        <div class="w-56 py-1">
          <!-- View Modes -->
          <div class="px-3 py-2 text-xs font-semibold text-gray-400 uppercase tracking-wider border-b border-gray-700/50 mt-1 mb-1">
            View Modes
          </div>

          <button 
            @click="fileStore.toggleZenMode()"
            class="w-full text-left px-3 py-2 text-sm hover:bg-gray-700/30 transition-colors flex items-center justify-between group"
            :class="fileStore.isZenMode ? 'text-indigo-400' : 'text-gray-300'"
          >
            <div class="flex items-center gap-2">
              <Eye class="w-3.5 h-3.5" />
              <span>Zen Mode</span>
            </div>
            <span v-if="fileStore.isZenMode" class="w-1.5 h-1.5 rounded-full bg-indigo-500"></span>
          </button>

          <button 
            @click="fileStore.toggleSoloExpansionMode()"
            class="w-full text-left px-3 py-2 text-sm hover:bg-gray-700/30 transition-colors flex items-center justify-between group"
            :class="fileStore.isSoloExpansionMode ? 'text-indigo-400' : 'text-gray-300'"
          >
            <div class="flex items-center gap-2">
              <List class="w-3.5 h-3.5" />
              <span>Solo Expansion</span>
            </div>
            <span v-if="fileStore.isSoloExpansionMode" class="w-1.5 h-1.5 rounded-full bg-indigo-500"></span>
          </button>

          <button 
            @click="fileStore.toggleSelectedOnlyMode()"
            class="w-full text-left px-3 py-2 text-sm hover:bg-gray-700/30 transition-colors flex items-center justify-between group"
            :class="fileStore.isSelectedOnlyMode ? 'text-indigo-400' : 'text-gray-300'"
            :disabled="fileStore.selectedCount === 0"
          >
            <div class="flex items-center gap-2">
              <CheckSquare class="w-3.5 h-3.5" />
              <span>Selected Only</span>
            </div>
            <span v-if="fileStore.isSelectedOnlyMode" class="w-1.5 h-1.5 rounded-full bg-indigo-500"></span>
          </button>

          <div class="border-t border-gray-700/50 mt-1 mb-1"></div>

          <div class="px-3 py-2 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-1">
            Filters
          </div>
          
          <label class="flex items-center px-3 py-2 hover:bg-gray-700/30 cursor-pointer transition-colors group">
            <input 
              type="checkbox" 
              :model-value="(settingsStore.settings.fileExplorer as any).hideNodeModules"
              @update:model-value="(settingsStore.settings.fileExplorer as any).hideNodeModules = $event"
              class="form-checkbox w-3.5 h-3.5 text-indigo-500 rounded border-gray-600 bg-gray-800 focus:ring-offset-0 focus:ring-1 focus:ring-indigo-500 mr-2"
            />
            <span class="text-sm text-gray-300 group-hover:text-white transition-colors">Hide node_modules</span>
          </label>

          <label class="flex items-center px-3 py-2 hover:bg-gray-700/30 cursor-pointer transition-colors group">
            <input 
              type="checkbox" 
              :model-value="(settingsStore.settings.fileExplorer as any).hideHiddenFiles"
              @update:model-value="(settingsStore.settings.fileExplorer as any).hideHiddenFiles = $event"
              class="form-checkbox w-3.5 h-3.5 text-indigo-500 rounded border-gray-600 bg-gray-800 focus:ring-offset-0 focus:ring-1 focus:ring-indigo-500 mr-2"
            />
            <span class="text-sm text-gray-300 group-hover:text-white transition-colors">Hide dotfiles (.*)</span>
          </label>

          <label class="flex items-center px-3 py-2 hover:bg-gray-700/30 cursor-pointer transition-colors group">
            <input 
              type="checkbox" 
              :model-value="(settingsStore.settings.fileExplorer as any).hideTestFiles"
              @update:model-value="(settingsStore.settings.fileExplorer as any).hideTestFiles = $event"
              class="form-checkbox w-3.5 h-3.5 text-indigo-500 rounded border-gray-600 bg-gray-800 focus:ring-offset-0 focus:ring-1 focus:ring-indigo-500 mr-2"
            />
            <span class="text-sm text-gray-300 group-hover:text-white transition-colors">Hide tests</span>
          </label>
          
          <div class="border-t border-gray-700/50 mt-1 pt-1">
             <label class="flex items-center px-3 py-2 hover:bg-gray-700/30 cursor-pointer transition-colors group">
              <input 
                type="checkbox" 
                v-model="settingsStore.settings.fileExplorer.compactNestedFolders"
                class="form-checkbox w-3.5 h-3.5 text-indigo-500 rounded border-gray-600 bg-gray-800 focus:ring-offset-0 focus:ring-1 focus:ring-indigo-500 mr-2"
              />
              <span class="text-sm text-gray-300 group-hover:text-white transition-colors">Compact Folders</span>
            </label>
          </div>
        </div>
      </template>
    </BaseDropdown>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useSettingsStore } from '@/stores/settings.store'
import { useFileStore } from '../model/file.store'
import { BaseDropdown, BaseButton } from '@/components/ui'
import { Eye, List, CheckSquare, Filter as FilterIcon } from 'lucide-vue-next'

const settingsStore = useSettingsStore()
const fileStore = useFileStore()
const isOpen = ref(false)

const hasActiveFilters = computed(() => {
  const s = settingsStore.settings.fileExplorer
  return (s as any).hideNodeModules || (s as any).hideHiddenFiles || (s as any).hideTestFiles
})
</script>
