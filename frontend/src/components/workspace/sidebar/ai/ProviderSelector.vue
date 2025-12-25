<template>
  <div class="settings-field">
    <label class="settings-label">{{ t('settings.selectProvider') }}</label>
    <div class="settings-select-wrap">
      <select
        :value="modelValue"
        @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
        class="settings-select"
      >
        <option v-for="provider in providers" :key="provider.id" :value="provider.id">
          {{ provider.icon }} {{ provider.name }}
        </option>
      </select>
      <div class="settings-select-icon">
        <ChevronDown class="w-4 h-4" />
      </div>
    </div>
    <p v-if="description" class="settings-hint">{{ description }}</p>
  </div>
</template>

<script setup lang="ts">
import type { AIProvider } from '@/composables/useAISettings';
import { useI18n } from '@/composables/useI18n';
import { ChevronDown } from 'lucide-vue-next';

const { t } = useI18n()

defineProps<{
  modelValue: string
  providers: AIProvider[]
  description?: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<style scoped>
.settings-field {
  display: flex;
  flex-direction: column;
}

.settings-label {
  font-size: 11px;
  font-weight: 500;
  color: #9ca3af;
  margin-bottom: 6px;
}

.settings-select-wrap {
  position: relative;
}

.settings-select {
  width: 100%;
  padding: 10px 36px 10px 12px;
  background: #0a0c12;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  font-size: 13px;
  color: #e5e7eb;
  cursor: pointer;
  appearance: none;
  outline: none;
  transition: all 0.15s ease-out;
}

.settings-select:hover {
  border-color: rgba(255, 255, 255, 0.15);
}

.settings-select:focus {
  border-color: #8b5cf6;
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.15);
}

.settings-select option {
  background: #1c1f2e;
  color: #e5e7eb;
  padding: 8px;
}

.settings-select-icon {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #6b7280;
  pointer-events: none;
}

.settings-hint {
  font-size: 10px;
  color: #4b5563;
  margin-top: 6px;
  line-height: 1.4;
}
</style>
