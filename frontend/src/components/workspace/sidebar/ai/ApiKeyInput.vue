<template>
  <div class="settings-field">
    <label class="settings-label">
      {{ t('settings.apiKey') }} ({{ providerName }})
    </label>
    <div class="settings-input-wrap">
      <input
        :type="showKey ? 'text' : 'password'"
        :value="modelValue"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        class="settings-input settings-input--password"
        :placeholder="t('settings.apiKeyPlaceholder')"
      />
      <div class="settings-input-actions">
        <button
          @click="$emit('toggle-visibility')"
          class="settings-input-btn"
          :title="showKey ? t('settings.hideKey') : t('settings.showKey')"
          :aria-label="showKey ? t('settings.hideKey') : t('settings.showKey')"
        >
          <EyeOff v-if="showKey" class="w-4 h-4" aria-hidden="true" />
          <Eye v-else class="w-4 h-4" aria-hidden="true" />
        </button>
        <button
          v-if="modelValue"
          @click="$emit('clear')"
          class="settings-input-btn settings-input-btn--danger"
          :title="t('settings.clearKey')"
          :aria-label="t('settings.clearKey')"
        >
          <X class="w-4 h-4" aria-hidden="true" />
        </button>
      </div>
    </div>
    <p v-if="hint" class="settings-hint">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';
import { Eye, EyeOff, X } from 'lucide-vue-next';

const { t } = useI18n()

defineProps<{
  modelValue: string
  providerName: string
  showKey: boolean
  hint?: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
  'toggle-visibility': []
  'clear': []
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

.settings-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
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

.settings-input--password {
  padding-right: 76px;
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

.settings-input-actions {
  position: absolute;
  right: 8px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.settings-input-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: rgba(255, 255, 255, 0.05);
  border: none;
  border-radius: 6px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.settings-input-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #e5e7eb;
}

.settings-input-btn--danger:hover {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
}

.settings-hint {
  font-size: 10px;
  color: #4b5563;
  margin-top: 6px;
  line-height: 1.4;
}
</style>
