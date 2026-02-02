<template>
  <div class="qwen-cli-settings">
    <!-- Info Box -->
    <BaseAlert variant="info">
      <template #icon>
        <Info :size="16" />
      </template>
      <template #title>{{ t('settings.qwenCliTitle') }}</template>
      <p class="text-xs">{{ t('settings.qwenCliDescription') }}</p>
    </BaseAlert>

    <!-- Settings -->
    <div class="settings-group">
      <!-- YOLO Mode Toggle -->
      <label class="setting-row">
        <span class="setting-label">{{ t('settings.qwenCli.yoloMode') }}</span>
        <button
          type="button"
          :class="['toggle-btn', { active: settings.qwenCLISettings.yoloMode }]"
          @click="toggleSetting('yoloMode')"
        >
          <span class="toggle-slider" />
        </button>
      </label>
      <p class="setting-hint">{{ t('settings.qwenCli.yoloModeHint') }}</p>

      <!-- Approval Mode (shown when YOLO is off) -->
      <div v-if="!settings.qwenCLISettings.yoloMode" class="setting-row mt-3">
        <span class="setting-label">{{ t('settings.qwenCli.approvalMode') }}</span>
        <select
          :value="settings.qwenCLISettings.approvalMode"
          class="setting-select"
          @change="updateSelect('approvalMode', $event)"
        >
          <option value="default">{{ t('settings.qwenCli.approval.default') }}</option>
          <option value="auto-edit">{{ t('settings.qwenCli.approval.autoEdit') }}</option>
          <option value="plan">{{ t('settings.qwenCli.approval.plan') }}</option>
        </select>
      </div>

      <!-- Sandbox Mode Toggle -->
      <label class="setting-row mt-3">
        <span class="setting-label">{{ t('settings.qwenCli.sandboxMode') }}</span>
        <button
          type="button"
          :class="['toggle-btn', { active: settings.qwenCLISettings.sandboxMode }]"
          @click="toggleSetting('sandboxMode')"
        >
          <span class="toggle-slider" />
        </button>
      </label>
      <p class="setting-hint">{{ t('settings.qwenCli.sandboxModeHint') }}</p>

      <!-- Max Session Turns -->
      <div class="setting-row mt-3">
        <span class="setting-label">{{ t('settings.qwenCli.maxTurns') }}</span>
        <input
          type="number"
          :value="settings.qwenCLISettings.maxSessionTurns"
          min="0"
          max="100"
          class="setting-input"
          @change="updateNumber('maxSessionTurns', $event)"
        />
      </div>
      <p class="setting-hint">{{ t('settings.qwenCli.maxTurnsHint') }}</p>

      <!-- Debug Mode Toggle -->
      <label class="setting-row mt-3">
        <span class="setting-label">{{ t('settings.qwenCli.debugMode') }}</span>
        <button
          type="button"
          :class="['toggle-btn', { active: settings.qwenCLISettings.debugMode }]"
          @click="toggleSetting('debugMode')"
        >
          <span class="toggle-slider" />
        </button>
      </label>
      <p class="setting-hint">{{ t('settings.qwenCli.debugModeHint') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { Info } from 'lucide-vue-next'
import { BaseAlert } from '@/components/ui'

const props = defineProps<{
  settings: {
    qwenCLISettings: {
      approvalMode: string
      outputFormat: string
      yoloMode: boolean
      sandboxMode: boolean
      maxSessionTurns: number
      debugMode: boolean
    }
  }
}>()

const emit = defineEmits<{
  'update:settings': [key: string, value: string | boolean | number]
}>()

const { t } = useI18n()

function toggleSetting(key: string) {
  const current = props.settings.qwenCLISettings[key as keyof typeof props.settings.qwenCLISettings]
  emit('update:settings', key, !current)
}

function updateSelect(key: string, e: Event) {
  const value = (e.target as HTMLSelectElement).value
  emit('update:settings', key, value)
}

function updateNumber(key: string, e: Event) {
  const value = parseInt((e.target as HTMLInputElement).value, 10) || 0
  emit('update:settings', key, value)
}
</script>

<style scoped>
.qwen-cli-settings {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.settings-group {
  padding: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.setting-label {
  font-size: 12px;
  color: #d1d5db;
}

.setting-hint {
  font-size: 10px;
  color: #6b7280;
  margin-top: 4px;
}

.toggle-btn {
  position: relative;
  width: 36px;
  height: 20px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.2s;
}

.toggle-btn.active {
  background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
}

.toggle-slider {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 16px;
  height: 16px;
  background: white;
  border-radius: 50%;
  transition: transform 0.2s;
}

.toggle-btn.active .toggle-slider {
  transform: translateX(16px);
}

.setting-select {
  padding: 6px 10px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #e5e7eb;
  font-size: 12px;
  cursor: pointer;
}

.setting-select:focus {
  outline: none;
  border-color: rgba(139, 92, 246, 0.5);
}

.setting-select option {
  background: #1f2937;
  color: #e5e7eb;
}

.setting-input {
  width: 70px;
  padding: 6px 10px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #e5e7eb;
  font-size: 12px;
  text-align: center;
}

.setting-input:focus {
  outline: none;
  border-color: rgba(139, 92, 246, 0.5);
}

.setting-input::-webkit-inner-spin-button,
.setting-input::-webkit-outer-spin-button {
  opacity: 1;
}
</style>
