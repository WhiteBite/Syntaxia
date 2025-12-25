<template>
  <div class="settings-field">
    <label class="settings-label">{{ t('settings.hostUrl') }}</label>
    <input
      type="text"
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      class="settings-input"
      :placeholder="placeholder"
    />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';

const { t } = useI18n()

defineProps<{
  modelValue: string
  placeholder?: string
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

.settings-input {
  width: 100%;
  padding: 10px 12px;
  background: #0a0c12;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  font-size: 13px;
  font-family: ui-monospace, monospace;
  color: #e5e7eb;
  outline: none;
  transition: all 0.15s ease-out;
}

.settings-input::placeholder {
  color: #4b5563;
  font-family: inherit;
}

.settings-input:hover {
  border-color: rgba(255, 255, 255, 0.15);
}

.settings-input:focus {
  border-color: #8b5cf6;
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.15);
}
</style>
