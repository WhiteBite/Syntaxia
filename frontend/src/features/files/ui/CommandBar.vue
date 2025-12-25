<template>
  <div class="magic-bar-wrapper">
    <!-- MAGIC CONTROL BAR -->
    <div class="magic-bar">
      <!-- LEFT: Token Limit Selector -->
      <div class="limit-section" @click="toggleDropdown" ref="limitRef">
        <div class="limit-main">
          <span class="limit-value">{{ formatTokens(settings.maxTokens) }}</span>
          <svg class="limit-chevron" :class="{ open: showDropdown }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M19 9l-7 7-7-7" />
          </svg>
        </div>
        
        <!-- Dropdown -->
        <Transition name="dropdown">
          <div v-if="showDropdown" class="limit-dropdown">
            <button 
              v-for="preset in tokenPresets" 
              :key="preset.value"
              @click.stop="selectPreset(preset.value)"
              class="limit-option"
              :class="{ active: isPresetActive(preset.value) }"
            >
              <span class="limit-option-value">{{ preset.label }}</span>
              <span class="limit-option-model">{{ preset.model }}</span>
            </button>
            
            <!-- Custom Input -->
            <div class="limit-custom">
              <input
                ref="customInputRef"
                v-model="customTokenValue"
                type="text"
                inputmode="numeric"
                class="limit-custom-input"
                :placeholder="t('commandBar.customLimit')"
                @click.stop
                @keydown.enter="applyCustomLimit"
                @focus="isCustomFocused = true"
                @blur="isCustomFocused = false"
              />
              <span class="limit-custom-suffix">K</span>
              <button 
                v-if="customTokenValue"
                class="limit-custom-apply"
                @click.stop="applyCustomLimit"
              >
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
              </button>
            </div>
          </div>
        </Transition>
      </div>

      <!-- DIVIDER -->
      <div class="bar-divider"></div>

      <!-- RIGHT: Magic Build Button -->
      <button 
        class="build-section"
        :class="{ disabled: isDisabled && !isBuilding, loading: isBuilding }"
        :disabled="isButtonDisabled"
        @click="handleBuild"
      >
        <!-- Gradient Background -->
        <div class="build-bg"></div>
        
        <!-- Inner Ring (creates depth) -->
        <div class="build-ring"></div>
        
        <!-- Bottom Glow -->
        <div class="build-glow"></div>
        
        <!-- Shimmer -->
        <div class="build-shimmer"></div>
        
        <!-- Content -->
        <div class="build-content">
          <svg v-if="!isBuilding" class="build-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
          <svg v-else class="build-spinner" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
          </svg>
          <span class="build-text">{{ isBuilding ? t('context.building') : t('commandBar.build') }}</span>
        </div>
      </button>
    </div>

    <!-- File Counter -->
    <div v-if="selectedCount > 0" class="file-counter">
      {{ t('commandBar.selected') }}: <span class="file-count">{{ selectedCount }}</span> 
      <span class="token-estimate">~{{ estimatedTokens }}k tokens</span>
      <button 
        class="clear-btn"
        @click="handleClear"
        :title="t('files.clearSelection')"
      >
        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
        {{ t('files.clear') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useFileStore } from '@/features/files/model/file.store'
import { useSettingsStore } from '@/stores/settings.store'
import { computed, onMounted, onUnmounted, ref } from 'vue'

const props = defineProps<{
  selectedCount: number
  isBuilding: boolean
}>()

const emit = defineEmits<{
  (e: 'build'): void
}>()

const { t } = useI18n()
const settingsStore = useSettingsStore()
const fileStore = useFileStore()
const settings = computed(() => settingsStore.settings.context)

const showDropdown = ref(false)
const limitRef = ref<HTMLElement | null>(null)
const customTokenValue = ref('')
const isCustomFocused = ref(false)

const isDisabled = computed(() => props.selectedCount === 0)
const isButtonDisabled = computed(() => props.selectedCount === 0 || props.isBuilding)
const estimatedTokens = computed(() => Math.round(fileStore.estimatedTokenCount / 1000))

const tokenPresets = [
  { value: 32000, label: '32K', model: 'GPT-4' },
  { value: 128000, label: '128K', model: 'GPT-4 Turbo' },
  { value: 200000, label: '200K', model: 'Claude' },
  { value: 1000000, label: '1M', model: 'Gemini' },
]

function formatTokens(n: number): string {
  if (n >= 1000000) return `${(n / 1000000).toFixed(n % 1000000 === 0 ? 0 : 1)}M`
  if (n >= 1000) return `${Math.round(n / 1000)}K`
  return n.toString()
}

function isPresetActive(value: number): boolean {
  return Math.abs(settings.value.maxTokens - value) <= value * 0.05
}

function toggleDropdown() {
  showDropdown.value = !showDropdown.value
}

function selectPreset(value: number) {
  settingsStore.updateContextSettings({ maxTokens: value })
  showDropdown.value = false
}

function applyCustomLimit() {
  const input = customTokenValue.value.replace(/[^\d.]/g, '')
  const value = parseFloat(input)
  if (!isNaN(value) && value > 0) {
    const tokens = Math.round(value * 1000)
    const clamped = Math.min(Math.max(tokens, 1000), 10000000)
    settingsStore.updateContextSettings({ maxTokens: clamped })
    customTokenValue.value = ''
    showDropdown.value = false
  }
}

function handleBuild() {
  if (!isDisabled.value) {
    emit('build')
  }
}

function handleClear() {
  fileStore.clearSelection()
}

function handleClickOutside(event: MouseEvent) {
  if (showDropdown.value && limitRef.value && !limitRef.value.contains(event.target as Node)) {
    showDropdown.value = false
  }
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))
</script>

<style scoped>
.magic-bar-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

/* MAGIC BAR */
.magic-bar {
  display: flex;
  height: 56px;
  background: #1c1f2e;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 4px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3);
}

/* LIMIT SECTION */
.limit-section {
  position: relative;
  width: 35%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s;
  padding: 0 12px;
}

.limit-main {
  display: flex;
  align-items: center;
  gap: 6px;
}

.limit-section:hover .limit-value {
  color: white;
}

.limit-value {
  font-size: 20px;
  font-family: ui-monospace, monospace;
  font-weight: 700;
  color: #e5e7eb;
  transition: color 0.15s;
  line-height: 1;
}

.limit-chevron {
  width: 14px;
  height: 14px;
  color: #6b7280;
  transition: all 0.2s;
  flex-shrink: 0;
}

.limit-section:hover .limit-chevron {
  color: #9ca3af;
}

.limit-chevron.open {
  transform: rotate(180deg);
}

/* LIMIT DROPDOWN */
.limit-dropdown {
  position: absolute;
  bottom: calc(100% + 8px);
  left: -4px;
  width: calc(100% + 8px);
  background: #1c1f2e;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 12px;
  padding: 6px;
  z-index: 50;
  box-shadow: 0 -8px 32px rgba(0, 0, 0, 0.6), 0 0 0 1px rgba(0, 0, 0, 0.3);
}

.limit-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 8px 12px;
  background: transparent;
  border: none;
  border-radius: 8px;
  color: #9ca3af;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.1s;
}

.limit-option:hover {
  background: rgba(192, 132, 252, 0.1);
  color: white;
}

.limit-option.active {
  background: rgba(192, 132, 252, 0.2);
  color: #e9d5ff;
}

.limit-option-value {
  font-family: ui-monospace, monospace;
  font-weight: 600;
}

.limit-option-model {
  font-size: 11px;
  color: #6b7280;
}

/* Custom Input */
.limit-custom {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 6px;
  padding: 6px 8px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.limit-custom-input {
  flex: 1;
  width: 100%;
  background: transparent;
  border: none;
  color: white;
  font-family: ui-monospace, monospace;
  font-size: 13px;
  font-weight: 600;
  outline: none;
}

.limit-custom-input::placeholder {
  color: #6b7280;
  font-weight: 400;
}

.limit-custom-suffix {
  color: #6b7280;
  font-size: 12px;
  font-weight: 500;
}

.limit-custom-apply {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: rgba(34, 197, 94, 0.2);
  border: none;
  border-radius: 6px;
  color: #22c55e;
  cursor: pointer;
  transition: all 0.15s;
}

.limit-custom-apply:hover {
  background: rgba(34, 197, 94, 0.3);
}

/* DIVIDER */
.bar-divider {
  width: 1px;
  height: 60%;
  margin: auto 0;
  background: rgba(255, 255, 255, 0.1);
}

/* BUILD SECTION */
.build-section {
  position: relative;
  flex: 1;
  margin-left: 4px;
  border-radius: 12px;
  border: none;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.15s ease-out;
  box-shadow: 0 4px 20px rgba(147, 51, 234, 0.35);
}

.build-section:hover:not(:disabled) {
  box-shadow: 0 6px 28px rgba(147, 51, 234, 0.5);
}

.build-section:active:not(:disabled) {
  transform: scale(0.98);
}

.build-section:focus {
  outline: none;
}

.build-section.disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.build-section.disabled .build-bg {
  filter: grayscale(0.6);
}

.build-section.disabled .build-glow,
.build-section.disabled .build-shimmer {
  display: none;
}

/* Loading state - keep gradient with pulse animation, text stays white */
.build-section.loading {
  cursor: wait;
  opacity: 1 !important;
  animation: btn-pulse 2s ease-in-out infinite;
}

.build-section.loading .build-bg {
  opacity: 1;
  filter: none;
}

.build-section.loading .build-content {
  color: white;
  opacity: 1;
}

.build-section.loading .build-text {
  color: white;
}

.build-section.loading .build-shimmer {
  animation: shimmer 1s infinite;
}

.build-section.loading .build-glow {
  opacity: 0.6;
  animation: pulse-glow 1.5s ease-in-out infinite;
}

@keyframes btn-pulse {
  0%, 100% { 
    box-shadow: 0 4px 20px rgba(147, 51, 234, 0.35);
  }
  50% { 
    box-shadow: 0 4px 32px rgba(147, 51, 234, 0.6), 0 0 20px rgba(219, 39, 119, 0.3);
  }
}

@keyframes pulse-glow {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 0.9; }
}

/* Gradient Background - Pink to Purple */
.build-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #db2777 0%, #9333ea 50%, #7c3aed 100%);
  transition: transform 0.15s ease-out;
}

.build-section:hover:not(:disabled) .build-bg {
  transform: scale(1.05);
}

/* Inner Ring (depth effect) */
.build-ring {
  position: absolute;
  inset: 0;
  border-radius: 12px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

/* Bottom Glow */
.build-glow {
  position: absolute;
  bottom: -10px;
  left: 50%;
  transform: translateX(-50%);
  width: 70%;
  height: 24px;
  background: linear-gradient(90deg, #db2777, #9333ea);
  filter: blur(20px);
  opacity: 0.6;
  transition: opacity 0.15s ease-out;
}

.build-section:hover:not(:disabled) .build-glow {
  opacity: 0.9;
}

/* Shimmer */
.build-shimmer {
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.25) 50%, transparent 100%);
  transform: translateX(-100%) skewX(-12deg);
  animation: shimmer 3s infinite;
}

.build-section:hover:not(:disabled) .build-shimmer {
  animation-duration: 1.5s;
}

/* Content */
.build-content {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 100%;
  color: white;
  z-index: 1;
}

.build-icon {
  width: 20px;
  height: 20px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
}

.build-spinner {
  width: 20px;
  height: 20px;
  animation: spin 1s linear infinite;
}

.build-text {
  font-size: 14px;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

/* FILE COUNTER */
.file-counter {
  text-align: center;
  font-size: 10px;
  color: #6b7280;
}

.file-count {
  color: #9ca3af;
  font-weight: 600;
}

.token-estimate {
  color: #6b7280;
  margin-left: 4px;
}

.clear-btn {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  margin-left: 8px;
  padding: 2px 6px;
  background: transparent;
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 4px;
  color: #9ca3af;
  font-size: 10px;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.clear-btn:hover {
  background: rgba(239, 68, 68, 0.15);
  border-color: rgba(239, 68, 68, 0.5);
  color: #f87171;
}

/* DROPDOWN TRANSITION */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

/* ANIMATIONS */
@keyframes shimmer {
  0% { transform: translateX(-100%) skewX(-12deg); }
  100% { transform: translateX(200%) skewX(-12deg); }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
