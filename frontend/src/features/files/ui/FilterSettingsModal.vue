<template>
  <BaseModal 
    :model-value="isOpen"
    @update:model-value="(value) => !value && $emit('close')"
    :title="t('quickFilters.settingsTitle')"
    size="md"
    :close-on-backdrop="true"
    @close="$emit('close')"
  >
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

    <template #footer>
      <button @click="$emit('reset')" class="btn btn-secondary btn-sm">
        {{ t('quickFilters.reset') }}
      </button>
      <button @click="$emit('close')" class="btn btn-primary btn-sm">
        {{ t('quickFilters.done') }}
      </button>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';
import type { QuickFilterConfig } from '@/stores/settings.store';
import BaseModal from '@/components/ui/BaseModal.vue';

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
.settings-section h4 { 
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-3);
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.settings-filter-card {
  padding: var(--space-3);
  border-radius: var(--radius-xl);
  margin-bottom: var(--space-2);
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  transition: all var(--transition-fast);
}

.settings-filter-card:hover {
  background: var(--bg-3);
}

.settings-filter-header { 
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-2);
}

.settings-toggle { 
  position: relative;
  display: inline-flex;
  cursor: pointer;
}

.settings-toggle input { 
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}

.settings-toggle-track {
  width: 2rem;
  height: 1rem;
  border-radius: var(--radius-full);
  background: var(--bg-3);
  transition: all var(--transition-normal);
  position: relative;
}

.settings-toggle input:checked + .settings-toggle-track { 
  background: var(--accent-indigo);
}

.settings-toggle-track::after {
  content: '';
  position: absolute;
  left: 2px;
  top: 2px;
  width: 0.75rem;
  height: 0.75rem;
  border-radius: var(--radius-full);
  background: white;
  transition: transform var(--transition-normal);
}

.settings-toggle input:checked + .settings-toggle-track::after { 
  transform: translateX(16px);
}

.settings-filter-name { 
  flex: 1;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.settings-filter-count { 
  font-size: var(--font-size-xs);
  color: var(--text-muted);
}

.settings-filter-inputs { 
  margin-top: var(--space-2);
}
</style>
