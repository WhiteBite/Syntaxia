<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="modal-overlay" @click.self="$emit('close')">
        <div class="modal-content filter-settings-modal">
          <div class="modal-header">
            <h3>{{ t('quickFilters.settingsTitle') }}</h3>
            <button @click="$emit('close')" class="icon-btn" :aria-label="t('common.close')">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="settings-section">
              <h4>{{ t('quickFilters.customFilters') }}</h4>
              <div v-for="filter in filters" :key="filter.id" class="settings-filter-card">
                <div class="settings-filter-header">
                  <label class="settings-toggle">
                    <input type="checkbox" v-model="filter.enabled" />
                    <span class="settings-toggle-track"></span>
                  </label>
                  <span class="settings-filter-name">{{ filter.label }}</span>
                  <span class="settings-filter-count">{{ getCount(filter) }}</span>
                </div>
                <div class="settings-filter-inputs">
                  <input
                    type="text"
                    :value="filter.extensions?.join(', ')"
                    @change="$emit('updateExtensions', filter, ($event.target as HTMLInputElement).value)"
                    class="input input-sm"
                    :placeholder="t('quickFilters.extensionsPlaceholder')"
                  />
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button @click="$emit('reset')" class="btn btn-secondary btn-sm">
              {{ t('quickFilters.reset') }}
            </button>
            <button @click="$emit('close')" class="btn btn-primary btn-sm">
              {{ t('quickFilters.done') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';
import type { QuickFilterConfig } from '@/stores/settings.store';

const { t } = useI18n()

defineProps<{
  isOpen: boolean
  filters: QuickFilterConfig[]
  getCount: (filter: QuickFilterConfig) => number
}>()

defineEmits<{
  close: []
  reset: []
  updateExtensions: [filter: QuickFilterConfig, value: string]
}>()
</script>

<style scoped>
.filter-settings-modal { @apply max-w-md; }
.modal-header {
  @apply flex items-center justify-between px-4 py-3;
  border-bottom: 1px solid var(--border-default);
}
.modal-header h3 { @apply text-lg font-semibold; color: var(--text-primary); }
.modal-body { @apply p-4 max-h-[60vh] overflow-y-auto; }
.settings-section h4 { @apply text-xs font-semibold mb-3; color: var(--text-muted); }
.settings-filter-card {
  @apply p-3 rounded-xl mb-2;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
}
.settings-filter-header { @apply flex items-center gap-3 mb-2; }
.settings-toggle { @apply relative inline-flex cursor-pointer; }
.settings-toggle input { @apply sr-only; }
.settings-toggle-track {
  @apply w-8 h-4 rounded-full;
  background: var(--bg-3);
  transition: all 200ms ease-out;
}
.settings-toggle input:checked + .settings-toggle-track { background: var(--accent-indigo); }
.settings-toggle-track::after {
  content: '';
  @apply absolute left-0.5 top-0.5 w-3 h-3 rounded-full bg-white;
  transition: transform 200ms ease-out;
}
.settings-toggle input:checked + .settings-toggle-track::after { transform: translateX(16px); }
.settings-filter-name { @apply flex-1 text-sm font-medium; color: var(--text-primary); }
.settings-filter-count { @apply text-xs; color: var(--text-muted); }
.settings-filter-inputs { @apply mt-2; }
.modal-footer {
  @apply flex justify-end gap-2 px-4 py-3;
  border-top: 1px solid var(--border-default);
}
.modal-enter-active, .modal-leave-active { transition: all 200ms ease-out; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
.modal-enter-from .modal-content, .modal-leave-to .modal-content { transform: scale(0.95); }
</style>
