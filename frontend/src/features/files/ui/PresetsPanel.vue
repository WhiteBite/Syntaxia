<template>
  <div v-if="presets.length > 0 || !collapsed" class="presets-panel">
    <!-- Header -->
    <button 
      class="presets-header"
      @click="collapsed = !collapsed"
      :aria-expanded="!collapsed"
    >
      <BookmarkIcon class="w-4 h-4 text-indigo-400" />
      <span class="presets-title">{{ t('presets.title') }}</span>
      <span class="presets-count">{{ presets.length }}</span>
      <ChevronDownIcon 
        class="w-4 h-4 text-gray-400 transition-transform" 
        :class="{ 'rotate-180': collapsed }"
      />
    </button>

    <!-- Content -->
    <Transition name="collapse">
      <div v-if="!collapsed" class="presets-content">
        <!-- Save Current Selection Button -->
        <BaseButton
          v-if="fileStore.selectedCount > 0"
          variant="ghost"
          size="sm"
          class="w-full mb-2"
          @click="showSaveModal = true"
        >
          <template #icon>
            <PlusIcon class="w-4 h-4" />
          </template>
          {{ t('presets.saveCurrentShort') }}
        </BaseButton>

        <!-- Empty State -->
        <div v-if="presets.length === 0" class="presets-empty">
          <p class="text-xs text-gray-400">{{ t('presets.noPresets') }}</p>
          <p class="text-xs text-gray-500 mt-1">{{ t('presets.emptyHint') }}</p>
        </div>
        
        <!-- Presets List -->
        <div v-else class="presets-list">
          <div
            v-for="preset in presets"
            :key="preset.name"
            class="preset-item group"
          >
            <div class="preset-info" @click="handleLoad(preset)">
              <div class="preset-header-row">
                <span class="preset-name">{{ preset.name }}</span>
                <span class="preset-files-count">{{ preset.paths.length }} {{ t('presets.files') }}</span>
              </div>
              <p v-if="preset.description" class="preset-description">{{ preset.description }}</p>
              <span class="preset-date">{{ formatDate(preset.createdAt) }}</span>
            </div>
            <div class="preset-actions">
              <BaseButton
                variant="ghost"
                size="sm"
                icon-only
                @click="handleLoad(preset)"
                :title="t('presets.load')"
              >
                <template #icon>
                  <ArrowDownTrayIcon class="w-3.5 h-3.5" />
                </template>
              </BaseButton>
              <BaseButton
                variant="danger"
                size="sm"
                icon-only
                @click="handleDelete(preset)"
                :title="t('presets.delete')"
              >
                <template #icon>
                  <TrashIcon class="w-3.5 h-3.5" />
                </template>
              </BaseButton>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Save Preset Modal -->
    <BaseModal
      v-model="showSaveModal"
      :title="t('presets.savePresetTitle')"
      size="sm"
      :close-on-backdrop="true"
    >
      <div class="form-group">
        <label class="form-label">{{ t('presets.nameLabel') }}</label>
        <BaseInput
          v-model="newPresetName"
          :placeholder="t('presets.namePlaceholder')"
          @keydown.enter="handleSave"
        />
      </div>
      
      <div class="form-group">
        <label class="form-label">{{ t('presets.descriptionLabel') }}</label>
        <BaseInput
          v-model="newPresetDescription"
          :placeholder="t('presets.descriptionPlaceholder')"
        />
      </div>

      <div class="preset-info-box">
        <span class="text-xs text-gray-400">
          {{ t('presets.saveWithCount', { count: fileStore.selectedCount }) }}
        </span>
      </div>

      <template #footer>
        <BaseButton variant="ghost" @click="showSaveModal = false">
          {{ t('presets.cancel') }}
        </BaseButton>
        <BaseButton 
          variant="primary" 
          @click="handleSave"
          :disabled="!newPresetName.trim()"
        >
          {{ t('presets.save', { count: fileStore.selectedCount }) }}
        </BaseButton>
      </template>
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-model="showDeleteModal"
      :title="t('presets.confirmDelete')"
      size="sm"
      :close-on-backdrop="true"
    >
      <p class="text-sm text-gray-300">
        {{ t('presets.confirmDeleteMessage', { name: presetToDelete?.name || '' }) }}
      </p>

      <template #footer>
        <BaseButton variant="ghost" @click="showDeleteModal = false">
          {{ t('presets.cancel') }}
        </BaseButton>
        <BaseButton 
          variant="primary" 
          class="bg-red-600 hover:bg-red-700"
          @click="confirmDelete"
        >
          {{ t('presets.delete') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import BaseModal from '@/components/ui/BaseModal.vue'
import { useI18n } from '@/composables/useI18n'
import { useUIStore } from '@/stores/ui.store'
import { BaseButton, BaseInput } from '@/components/ui'
import { BookmarkIcon, ChevronDownIcon, PlusIcon, ArrowDownTrayIcon, TrashIcon } from '@heroicons/vue/24/solid'
import { computed, ref } from 'vue'
import { useFileStore } from '../model/file.store'
import type { SelectionPreset } from '@/composables/useFilePersistence'

const { t } = useI18n()
const fileStore = useFileStore()
const uiStore = useUIStore()

const collapsed = ref(false)
const showSaveModal = ref(false)
const showDeleteModal = ref(false)
const newPresetName = ref('')
const newPresetDescription = ref('')
const presetToDelete = ref<SelectionPreset | null>(null)

const presets = computed(() => fileStore.getPresets())

function formatDate(timestamp: number): string {
  const date = new Date(timestamp)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays === 0) {
    return 'Today'
  } else if (diffDays === 1) {
    return 'Yesterday'
  } else if (diffDays < 7) {
    return `${diffDays} days ago`
  } else {
    return date.toLocaleDateString()
  }
}

function handleSave() {
  const name = newPresetName.value.trim()
  if (!name) return

  try {
    fileStore.savePreset(name, newPresetDescription.value.trim())
    uiStore.addToast(t('presets.saveSuccess', { name }), 'success')
    
    // Reset form
    newPresetName.value = ''
    newPresetDescription.value = ''
    showSaveModal.value = false
  } catch (error) {
    uiStore.addToast('Failed to save preset', 'error')
  }
}

function handleLoad(preset: SelectionPreset) {
  try {
    const loadedCount = fileStore.loadPreset(preset.name)
    if (loadedCount > 0) {
      uiStore.addToast(
        t('presets.loaded') + ` "${preset.name}" (${loadedCount} ${t('presets.files')})`,
        'success'
      )
    } else {
      uiStore.addToast('No valid files found in preset', 'warning')
    }
  } catch (error) {
    uiStore.addToast('Failed to load preset', 'error')
  }
}

function handleDelete(preset: SelectionPreset) {
  presetToDelete.value = preset
  showDeleteModal.value = true
}

function confirmDelete() {
  if (!presetToDelete.value) return

  try {
    fileStore.deletePreset(presetToDelete.value.name)
    uiStore.addToast(t('presets.deleted'), 'success')
    showDeleteModal.value = false
    presetToDelete.value = null
  } catch (error) {
    uiStore.addToast('Failed to delete preset', 'error')
  }
}
</script>

<style scoped>
.presets-panel {
  border-bottom: 1px solid var(--border-default);
  background: rgba(255, 255, 255, 0.02);
}

.presets-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.5rem 0.75rem;
  text-align: left;
  background: transparent;
  border: none;
  cursor: pointer;
  transition: background 0.15s;
}

.presets-header:hover {
  background: rgba(255, 255, 255, 0.05);
}

.presets-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.presets-count {
  font-size: 0.65rem;
  padding: 0.125rem 0.375rem;
  background: rgba(99, 102, 241, 0.2);
  color: rgb(99, 102, 241);
  border-radius: 9999px;
  margin-left: auto;
}

.presets-content {
  padding: 0.25rem 0.5rem 0.5rem;
}

.presets-empty {
  padding: 0.75rem;
  text-align: center;
}

.presets-list {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.preset-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  border-radius: 0.375rem;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  transition: all 0.15s;
}

.preset-item:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(99, 102, 241, 0.3);
}

.preset-info {
  flex: 1;
  cursor: pointer;
  min-width: 0;
}

.preset-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  margin-bottom: 0.25rem;
}

.preset-name {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preset-files-count {
  font-size: 0.65rem;
  padding: 0.125rem 0.375rem;
  background: rgba(99, 102, 241, 0.15);
  color: rgb(129, 140, 248);
  border-radius: 9999px;
  flex-shrink: 0;
}

.preset-description {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-bottom: 0.25rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preset-date {
  font-size: 0.6875rem;
  color: var(--text-muted);
}

.preset-actions {
  display: flex;
  gap: 0.25rem;
  opacity: 0;
  transition: opacity 0.15s;
}

.preset-item:hover .preset-actions {
  opacity: 1;
}

/* Form styles for modal content */
.form-group {
  margin-bottom: 1rem;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 0.5rem;
}

.preset-info-box {
  padding: 0.75rem;
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 0.5rem;
  margin-top: 1rem;
}

/* Collapse animation */
.collapse-enter-active,
.collapse-leave-active {
  transition: all 0.2s ease;
  overflow: hidden;
}

.collapse-enter-from,
.collapse-leave-to {
  opacity: 0;
  max-height: 0;
}

.collapse-enter-to,
.collapse-leave-from {
  opacity: 1;
  max-height: 30rem;
}
</style>
