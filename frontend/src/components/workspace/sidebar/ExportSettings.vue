<template>
  <div class="inspector">
    <div class="inspector-content">
      <!-- Task Input -->
      <div class="inspector-section">
        <div class="task-wrapper">
          <span class="task-icon">✨</span>
          <textarea 
            v-model="task" 
            :placeholder="t('templates.taskPlaceholder')" 
            class="task-input" 
            rows="2" 
          />
        </div>
        <button class="template-row" @click="templateStore.openModal()">
          <span class="template-icon">{{ activeTemplate?.icon || '📝' }}</span>
          <span class="template-name">{{ activeTemplate?.name || t('templates.select') }}</span>
          <svg class="template-arrow" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>

      <div class="inspector-divider" />

      <!-- Output Options -->
      <div class="inspector-section">
        <div class="section-header">OUTPUT</div>
        <ToggleItem v-model="settings.applyTemplateOnCopy" :label="t('export.applyTemplate')" @update:model-value="update('applyTemplateOnCopy', $event)" />
        <ToggleItem v-model="settings.includeManifest" :label="t('export.includeManifest')" @update:model-value="update('includeManifest', $event)" />
        <ToggleItem v-model="settings.includeLineNumbers" :label="t('export.includeLineNumbers')" @update:model-value="update('includeLineNumbers', $event)" />
        <ToggleItem v-model="settings.stripComments" :label="t('export.stripComments')" @update:model-value="update('stripComments', $event)" />
      </div>

      <div class="inspector-divider" />

      <!-- Optimization -->
      <div class="inspector-section">
        <div class="section-header">OPTIMIZATION</div>
        <ToggleItem v-model="settings.excludeTests" :label="t('export.excludeTests')" @update:model-value="update('excludeTests', $event)" />
        <ToggleItem v-model="settings.stripLicense" :label="t('export.stripLicense')" @update:model-value="update('stripLicense', $event)" />
        <ToggleItem v-model="settings.compactDataFiles" :label="t('export.compactDataFiles')" @update:model-value="update('compactDataFiles', $event)" />
        <ToggleItem v-model="settings.trimWhitespace" :label="t('export.trimWhitespace')" @update:model-value="update('trimWhitespace', $event)" />
        <ToggleItem v-model="settings.collapseEmptyLines" :label="t('export.collapseEmptyLines')" @update:model-value="update('collapseEmptyLines', $event)" />
      </div>

      <div class="inspector-divider" />

      <!-- Chunking Section -->
      <div class="inspector-section">
        <div class="chunking-header">
          <span class="section-header">{{ t('export.chunking') }}</span>
          <button 
            class="chunking-toggle"
            :class="{ active: settings.enableAutoSplit }"
            @click="update('enableAutoSplit', !settings.enableAutoSplit)"
          >
            <span class="chunking-toggle-thumb" />
          </button>
        </div>
        
        <Transition name="accordion">
          <div v-if="settings.enableAutoSplit" class="chunking-content">
            <!-- Size: Presets + Custom inline -->
            <div class="chunk-row">
              <div class="chunk-presets">
                <button 
                  v-for="preset in chunkPresets" 
                  :key="preset.value"
                  class="preset-btn"
                  :class="{ active: isChunkPresetActive(preset.value) }"
                  @click="selectPreset(preset.value)"
                >{{ preset.label }}</button>
                <input 
                  v-if="showCustomInput"
                  type="text" 
                  v-model="customValue"
                  @blur="applyCustomValue"
                  @keyup.enter="applyCustomValue"
                  class="preset-custom"
                  placeholder="Custom"
                  ref="customInputRef"
                />
                <button 
                  v-else
                  class="preset-btn preset-btn--custom"
                  :class="{ active: isCustomActive }"
                  @click="enableCustomInput"
                >{{ isCustomActive ? formatChunkTokens(settings.maxTokensPerChunk) : '...' }}</button>
              </div>
            </div>

            <!-- Strategy: Segmented Control -->
            <div class="chunk-row">
              <div class="strategy-segment">
                <button 
                  class="strategy-btn"
                  :class="{ active: settings.splitStrategy === 'smart' || settings.splitStrategy === 'file' }"
                  @click="update('splitStrategy', 'smart')"
                  :title="t('export.strategy.smartDesc')"
                >
                  <span>🧠</span>
                  <span>Smart</span>
                </button>
                <button 
                  class="strategy-btn"
                  :class="{ active: settings.splitStrategy === 'token' }"
                  @click="update('splitStrategy', 'token')"
                  :title="t('export.strategy.hardDesc')"
                >
                  <span>✂️</span>
                  <span>Hard</span>
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </div>

    <TemplateModal />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { TemplateModal, useTemplateStore } from '@/features/templates'
import { useSettingsStore, type ContextSettings } from '@/stores/settings.store'
import { computed, nextTick, ref } from 'vue'
import ToggleItem from './export/ToggleItem.vue'

const { t } = useI18n()
const settingsStore = useSettingsStore()
const templateStore = useTemplateStore()
const settings = computed(() => settingsStore.settings.context)
const activeTemplate = computed(() => templateStore.activeTemplate)
const task = computed({ get: () => templateStore.currentTask, set: (v: string) => templateStore.setTask(v) })

const chunkPresets = [
  { value: 32000, label: '32K' },
  { value: 64000, label: '64K' },
  { value: 128000, label: '128K' },
]

const showCustomInput = ref(false)
const customValue = ref('')
const customInputRef = ref<HTMLInputElement | null>(null)

const isCustomActive = computed(() => {
  return !chunkPresets.some(p => isChunkPresetActive(p.value))
})

function isChunkPresetActive(value: number): boolean {
  return Math.abs(settings.value.maxTokensPerChunk - value) <= value * 0.05
}

function selectPreset(value: number) {
  showCustomInput.value = false
  update('maxTokensPerChunk', value)
}

function enableCustomInput() {
  showCustomInput.value = true
  customValue.value = formatChunkTokens(settings.value.maxTokensPerChunk)
  nextTick(() => customInputRef.value?.focus())
}

function applyCustomValue() {
  const input = customValue.value.trim().toUpperCase()
  let value = 0
  if (input.endsWith('K')) {
    value = parseFloat(input.slice(0, -1)) * 1000
  } else if (input.endsWith('M')) {
    value = parseFloat(input.slice(0, -1)) * 1000000
  } else {
    value = parseFloat(input)
  }
  if (!isNaN(value) && value > 0) {
    update('maxTokensPerChunk', Math.round(value))
  }
  showCustomInput.value = false
}

function update<K extends keyof ContextSettings>(key: K, value: ContextSettings[K]) { 
  settingsStore.updateContextSettings({ [key]: value }) 
}

function formatChunkTokens(n: number): string {
  if (n >= 1000000) return `${(n / 1000000).toFixed(1)}M`
  if (n >= 1000) return `${Math.round(n / 1000)}K`
  return n.toString()
}
</script>

<style scoped>
.inspector {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: linear-gradient(180deg, #131620 0%, #0f111a 100%);
  font-size: 12px;
}

.inspector-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.1) transparent;
}

.inspector-content::-webkit-scrollbar { width: 4px; }
.inspector-content::-webkit-scrollbar-track { background: transparent; }
.inspector-content::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.1); border-radius: 2px; }
.inspector-content::-webkit-scrollbar-thumb:hover { background: rgba(255, 255, 255, 0.2); }

.inspector-section { padding: 0.5rem 0.75rem; }
.inspector-divider { height: 1px; background: rgba(255, 255, 255, 0.04); margin: 0 0.75rem; }

/* Section Header - improved */
.section-header {
  font-size: 10px;
  font-weight: 700;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  margin-bottom: 0.375rem;
}

/* Task Input */
.task-wrapper {
  position: relative;
  display: flex;
  align-items: flex-start;
}

.task-icon {
  position: absolute;
  left: 10px;
  top: 10px;
  font-size: 12px;
  color: #6b7280;
  pointer-events: none;
  transition: color 0.15s ease-out;
}

.task-wrapper:focus-within .task-icon {
  color: #a78bfa;
}

.task-input {
  width: 100%;
  padding: 0.5rem 0.625rem 0.5rem 2rem;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  color: #e5e7eb;
  font-size: 12px;
  font-family: inherit;
  resize: none;
  line-height: 1.4;
  transition: all 0.15s ease-out;
}

.task-input:focus {
  outline: none;
  border-color: #8b5cf6;
  background: rgba(139, 92, 246, 0.05);
  box-shadow: 0 0 0 2px rgba(139, 92, 246, 0.15);
}

.task-input::placeholder { color: #4b5563; }

/* Template Row */
.template-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  width: 100%;
  margin-top: 0.375rem;
  padding: 0.375rem 0.5rem;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.template-row:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
}

.template-icon { font-size: 13px; line-height: 1; }
.template-name { flex: 1; text-align: left; font-size: 11px; color: #d1d5db; }
.template-arrow { width: 10px; height: 10px; color: #6b7280; transition: transform 0.15s ease-out; }
.template-row:hover .template-arrow { transform: translateX(2px); }

/* Chunking */
.chunking-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chunking-toggle {
  position: relative;
  width: 36px;
  height: 20px;
  background: rgba(255, 255, 255, 0.08);
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.chunking-toggle.active {
  background: #8b5cf6;
}

.chunking-toggle-thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 16px;
  height: 16px;
  background: white;
  border-radius: 50%;
  transition: transform 0.15s ease-out;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.chunking-toggle.active .chunking-toggle-thumb {
  transform: translateX(16px);
}

.chunking-content {
  margin-top: 0.5rem;
  padding-top: 0.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.04);
}

.chunk-row {
  margin-bottom: 0.5rem;
}

.chunk-row:last-child {
  margin-bottom: 0;
}

/* Presets inline */
.chunk-presets {
  display: flex;
  gap: 4px;
}

.preset-btn {
  flex: 1;
  padding: 6px 8px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px;
  color: #9ca3af;
  font-size: 11px;
  font-weight: 600;
  font-family: ui-monospace, monospace;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.preset-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.preset-btn.active {
  background: #8b5cf6;
  border-color: #8b5cf6;
  color: white;
  box-shadow: 0 2px 8px rgba(139, 92, 246, 0.35);
}

.preset-btn--custom {
  min-width: 48px;
}

.preset-custom {
  flex: 1;
  min-width: 48px;
  padding: 6px 8px;
  background: rgba(139, 92, 246, 0.1);
  border: 1px solid #8b5cf6;
  border-radius: 6px;
  color: #e5e7eb;
  font-size: 11px;
  font-family: ui-monospace, monospace;
  text-align: center;
  outline: none;
}

/* Strategy Segment */
.strategy-segment {
  display: flex;
  gap: 4px;
  padding: 3px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
}

.strategy-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 6px 10px;
  background: transparent;
  border: none;
  border-radius: 5px;
  color: #6b7280;
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.strategy-btn:hover {
  color: #9ca3af;
  background: rgba(255, 255, 255, 0.03);
}

.strategy-btn.active {
  background: #8b5cf6;
  color: white;
  box-shadow: 0 2px 6px rgba(139, 92, 246, 0.35);
}

/* Accordion Transition */
.accordion-enter-active,
.accordion-leave-active {
  transition: all 0.15s ease-out;
  overflow: hidden;
}

.accordion-enter-from,
.accordion-leave-to {
  opacity: 0;
  max-height: 0;
  margin-top: 0;
  padding-top: 0;
}

.accordion-enter-to,
.accordion-leave-from {
  max-height: 150px;
}
</style>
