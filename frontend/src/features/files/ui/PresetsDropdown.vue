<template>
  <div class="relative">
    <BaseDropdown v-model="isOpen" placement="bottom-end">
      <template #trigger>
        <slot name="trigger">
          <BaseButton 
            variant="ghost" 
            size="xs" 
            icon-only 
            title="Saved Presets"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </BaseButton>
        </slot>
      </template>

      <template #content>
        <div class="w-64 py-1">
          <div class="px-3 py-2 text-xs font-semibold text-gray-400 uppercase tracking-wider border-b border-gray-700/50 mb-1">
            Saved Presets
          </div>

          <!-- Presets List -->
          <div v-if="fileStore.getPresets().length === 0" class="px-4 py-3 text-sm text-gray-500 text-center italic">
            No presets saved yet
          </div>
          
          <div v-else class="max-h-60 overflow-y-auto custom-scrollbar">
            <button
              v-for="preset in fileStore.getPresets()"
              :key="preset.name"
              class="w-full text-left px-3 py-2 text-sm text-gray-300 hover:bg-indigo-500/10 hover:text-indigo-300 transition-colors flex items-center justify-between group"
              @click="loadPreset(preset)"
            >
              <div class="flex items-center gap-2 truncate">
                <span class="truncate">{{ preset.name }}</span>
                <span class="text-xs text-gray-500">({{ preset.paths ? preset.paths.length : 0 }})</span>
              </div>
              
              <div class="opacity-0 group-hover:opacity-100 flex items-center gap-1 transition-opacity">
                <!-- Delete Preset -->
                <button 
                  @click.stop="deletePreset(preset.name)"
                  class="p-1 rounded hover:bg-red-500/20 text-gray-500 hover:text-red-400 transition-colors"
                  title="Delete preset"
                >
                  <TrashIcon class="w-3 h-3" />
                </button>
              </div>
            </button>
          </div>

          <div class="border-t border-gray-700/50 mt-1 pt-1">
            <!-- Save New Preset -->
            <button
              v-if="fileStore.selectedCount > 0"
              class="w-full text-left px-3 py-2 text-sm text-emerald-400 hover:bg-emerald-500/10 transition-colors flex items-center gap-2"
              @click="showSaveModal = true"
            >
              <PlusIcon class="w-4 h-4" />
              <span>Save current selection...</span>
            </button>
            <div v-else class="px-3 py-2 text-xs text-gray-500 italic text-center">
              Select files to create preset
            </div>
          </div>
        </div>
      </template>
    </BaseDropdown>

    <!-- Simple Save Modal (Inline) -->
    <div v-if="showSaveModal" class="fixed inset-0 z-[600] flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div class="bg-gray-800 border border-gray-700 rounded-lg shadow-xl p-4 w-80 animate-in fade-in zoom-in duration-200">
        <h3 class="text-sm font-semibold text-white mb-3">Save Preset</h3>
        <input 
          v-model="newPresetName"
          ref="nameInputRef"
          type="text" 
          class="w-full bg-gray-900 border border-gray-700 rounded px-3 py-2 text-sm text-white focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 mb-4"
          placeholder="Preset name (e.g. 'Frontend Core')"
          @keydown.enter="savePreset"
          @keydown.escape="showSaveModal = false"
        />
        <div class="flex justify-end gap-2">
          <button 
            class="px-3 py-1.5 text-xs font-medium text-gray-400 hover:text-white transition-colors"
            @click="showSaveModal = false"
          >
            Cancel
          </button>
          <button 
            class="px-3 py-1.5 text-xs font-medium bg-indigo-600 hover:bg-indigo-500 text-white rounded transition-colors"
            @click="savePreset"
            :disabled="!newPresetName.trim()"
          >
            Save
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'
import { useFileStore } from '../model/file.store'
import { useUIStore } from '@/stores/ui.store'
import { BaseDropdown, BaseButton } from '@/components/ui'
import { TrashIcon, PlusIcon } from 'lucide-vue-next'
import type { SelectionPreset } from '@/composables/useFilePersistence'

const fileStore = useFileStore()
const uiStore = useUIStore()

const showSaveModal = ref(false)
const newPresetName = ref('')
const nameInputRef = ref<HTMLInputElement | null>(null)
const isOpen = ref(false)

watch(showSaveModal, (val) => {
  if (val) {
    newPresetName.value = ''
    nextTick(() => {
      nameInputRef.value?.focus()
    })
  }
})

function loadPreset(preset: SelectionPreset) {
  fileStore.loadPreset(preset.name)
  uiStore.addToast(`Loaded preset: ${preset.name}`, 'success')
}

function deletePreset(name: string) {
  fileStore.deletePreset(name)
  uiStore.addToast('Preset deleted', 'info')
}

function savePreset() {
  if (!newPresetName.value.trim()) return
  
  fileStore.savePreset(newPresetName.value.trim(), '')
  uiStore.addToast('Preset saved', 'success')
  showSaveModal.value = false
}
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
