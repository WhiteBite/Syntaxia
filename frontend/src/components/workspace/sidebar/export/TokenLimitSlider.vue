<template>
  <div class="token-limit-container">
    <!-- Model Presets -->
    <div class="token-presets">
      <button
        v-for="preset in presets"
        :key="preset.id"
        class="token-preset"
        :class="{ 'token-preset-active': isPresetActive(preset) }"
        :title="preset.tooltip"
        @click="applyPreset(preset)"
      >
        <span class="token-preset-name">{{ preset.name }}</span>
        <span class="token-preset-value">{{ formatTokens(preset.tokens) }}</span>
      </button>
    </div>

    <!-- Custom Input -->
    <div class="token-custom">
      <div class="token-custom-header">
        <label class="text-xs text-gray-400">{{ t('export.tokenLimit') }}</label>
        <div class="flex items-center gap-2">
          <input 
            type="number"
            :value="modelValue"
            @input="handleInput"
            min="1000"
            max="1000000"
            step="1000"
            class="input text-xs py-0.5 px-2 w-20 text-right"
          />
          <span class="text-xs text-gray-400">{{ t('limits.tokens') }}</span>
        </div>
      </div>
      <input 
        type="range" 
        :value="modelValue"
        @input="emit('update:modelValue', Number(($event.target as HTMLInputElement).value))" 
        min="1000" 
        max="500000"
        step="1000" 
        class="token-slider" 
      />
      <div class="token-slider-labels">
        <span>1K</span>
        <span class="text-gray-400">{{ t('tokenLimit.hint') }}</span>
        <span>500K</span>
      </div>
    </div>

    <!-- Usage indicator -->
    <div v-if="contextTokens > 0" class="token-usage">
      <div class="token-usage-bar">
        <div 
          class="token-usage-fill"
          :class="usageClass"
          :style="{ width: `${usagePercent}%` }"
        />
      </div>
      <div class="token-usage-text">
        <span>{{ formatTokens(contextTokens) }} / {{ formatTokens(modelValue) }}</span>
        <span :class="usageClass">{{ usagePercent }}%</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useContextStore } from '@/features/context/model/context.store'
import { useI18n } from '@/composables/useI18n'
import { computed } from 'vue'

const { t } = useI18n()
const contextStore = useContextStore()

interface ModelPreset {
  id: string
  name: string
  tokens: number
  tooltip: string
}

const presets: ModelPreset[] = [
  { id: 'gpt4', name: 'GPT-4', tokens: 128000, tooltip: 'GPT-4 Turbo / GPT-4o (128K context)' },
  { id: 'claude', name: 'Claude', tokens: 200000, tooltip: 'Claude 3.5 Sonnet (200K context)' },
  { id: 'gemini', name: 'Gemini', tokens: 1000000, tooltip: 'Gemini 1.5 Pro (1M context)' },
  { id: 'small', name: '32K', tokens: 32000, tooltip: 'GPT-4 32K / Smaller models' },
]

const props = defineProps<{
  modelValue: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const contextTokens = computed(() => contextStore.stats?.tokens || 0)

const usagePercent = computed(() => {
  if (!props.modelValue) return 0
  return Math.min(100, Math.round((contextTokens.value / props.modelValue) * 100))
})

const usageClass = computed(() => {
  if (usagePercent.value >= 90) return 'usage-danger'
  if (usagePercent.value >= 70) return 'usage-warning'
  return 'usage-ok'
})

function formatTokens(tokens: number): string {
  if (tokens >= 1000000) return `${(tokens / 1000000).toFixed(0)}M`
  if (tokens >= 1000) return `${(tokens / 1000).toFixed(0)}K`
  return String(tokens)
}

function isPresetActive(preset: ModelPreset): boolean {
  // Allow 5% tolerance for "active" state
  const tolerance = preset.tokens * 0.05
  return Math.abs(props.modelValue - preset.tokens) <= tolerance
}

function applyPreset(preset: ModelPreset) {
  emit('update:modelValue', preset.tokens)
}

function handleInput(event: Event) {
  const value = Number((event.target as HTMLInputElement).value)
  const clamped = Math.max(1000, Math.min(1000000, value))
  emit('update:modelValue', clamped)
}
</script>

<style scoped>
.token-limit-container {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.token-presets {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.375rem;
}

.token-preset {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0.5rem 0.25rem;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 150ms;
}

.token-preset:hover {
  background: var(--bg-3);
  border-color: var(--color-primary);
}

.token-preset-active {
  background: rgba(99, 102, 241, 0.15);
  border-color: var(--color-primary);
}

.token-preset-name {
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--text-primary);
}

.token-preset-value {
  font-size: 0.625rem;
  color: var(--text-muted);
  margin-top: 0.125rem;
}

.token-preset-active .token-preset-value {
  color: var(--color-primary);
}

.token-custom {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.token-custom-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.token-slider {
  width: 100%;
  height: 0.375rem;
  background: var(--bg-3);
  border-radius: 0.25rem;
  appearance: none;
  cursor: pointer;
}

.token-slider::-webkit-slider-thumb {
  appearance: none;
  width: 0.875rem;
  height: 0.875rem;
  background: var(--color-primary);
  border-radius: 50%;
  cursor: pointer;
  transition: transform 150ms;
}

.token-slider::-webkit-slider-thumb:hover {
  transform: scale(1.2);
}

.token-slider-labels {
  display: flex;
  justify-content: space-between;
  font-size: 0.625rem;
  color: var(--text-muted);
}

.token-usage {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
  padding: 0.5rem 0.75rem;
  background: var(--bg-2);
  border-radius: 0.375rem;
}

.token-usage-bar {
  height: 0.375rem;
  background: var(--bg-3);
  border-radius: 0.25rem;
  overflow: hidden;
}

.token-usage-fill {
  height: 100%;
  border-radius: 0.25rem;
  transition: width 300ms ease-out;
}

.token-usage-text {
  display: flex;
  justify-content: space-between;
  font-size: 0.6875rem;
  color: var(--text-muted);
}

.usage-ok {
  background: #22c55e;
  color: #22c55e;
}

.usage-warning {
  background: #f59e0b;
  color: #f59e0b;
}

.usage-danger {
  background: #ef4444;
  color: #ef4444;
}
</style>
