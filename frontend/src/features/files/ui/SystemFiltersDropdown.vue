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
            :title="t('filters.title')"
          >
            <FilterIcon class="w-4 h-4" />
          </BaseButton>
        </slot>
      </template>

      <template #content>
        <div class="w-56 py-1">
          <div class="px-3 py-2 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-1">
            {{ t('filters.systemFilters') }}
          </div>
          
          <label class="flex items-center px-3 py-2 hover:bg-gray-700/30 cursor-pointer transition-colors group">
            <input 
              type="checkbox" 
              :model-value="(settingsStore.settings.fileExplorer as any).hideNodeModules"
              @update:model-value="(settingsStore.settings.fileExplorer as any).hideNodeModules = $event"
              class="form-checkbox w-3.5 h-3.5 text-indigo-500 rounded border-gray-600 bg-gray-800 focus:ring-offset-0 focus:ring-1 focus:ring-indigo-500 mr-2"
            />
            <span class="text-sm text-gray-300 group-hover:text-white transition-colors">{{ t('filters.hideNodeModules') }}</span>
          </label>

          <label class="flex items-center px-3 py-2 hover:bg-gray-700/30 cursor-pointer transition-colors group">
            <input 
              type="checkbox" 
              :model-value="(settingsStore.settings.fileExplorer as any).hideHiddenFiles"
              @update:model-value="(settingsStore.settings.fileExplorer as any).hideHiddenFiles = $event"
              class="form-checkbox w-3.5 h-3.5 text-indigo-500 rounded border-gray-600 bg-gray-800 focus:ring-offset-0 focus:ring-1 focus:ring-indigo-500 mr-2"
            />
            <span class="text-sm text-gray-300 group-hover:text-white transition-colors">{{ t('filters.hideDotfiles') }}</span>
          </label>

          <label class="flex items-center px-3 py-2 hover:bg-gray-700/30 cursor-pointer transition-colors group">
            <input 
              type="checkbox" 
              :model-value="(settingsStore.settings.fileExplorer as any).hideTestFiles"
              @update:model-value="(settingsStore.settings.fileExplorer as any).hideTestFiles = $event"
              class="form-checkbox w-3.5 h-3.5 text-indigo-500 rounded border-gray-600 bg-gray-800 focus:ring-offset-0 focus:ring-1 focus:ring-indigo-500 mr-2"
            />
            <span class="text-sm text-gray-300 group-hover:text-white transition-colors">{{ t('filters.hideTests') }}</span>
          </label>

          <div class="border-t border-gray-700/50 mt-1 pt-1">
            <button 
              @click="openAdvancedFilters"
              class="w-full text-left px-3 py-2 text-sm text-indigo-400 hover:bg-gray-700/30 transition-colors flex items-center gap-2 group"
            >
              <PlusCircle class="w-3.5 h-3.5" />
              <span>{{ t('filters.addTypeFilter') }}...</span>
            </button>
          </div>
        </div>
      </template>
    </BaseDropdown>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useSettingsStore } from '@/stores/settings.store'
import { BaseDropdown, BaseButton } from '@/components/ui'
import { Filter as FilterIcon, PlusCircle } from 'lucide-vue-next'
import { useI18n } from '@/composables/useI18n'

const emit = defineEmits<{
  (e: 'open-advanced'): void
}>()

const settingsStore = useSettingsStore()
const { t } = useI18n()
const isOpen = ref(false)

const hasActiveFilters = computed(() => {
  const s = settingsStore.settings.fileExplorer
  return (s as any).hideNodeModules || (s as any).hideHiddenFiles || (s as any).hideTestFiles
})

function openAdvancedFilters() {
  emit('open-advanced')
}
</script>
