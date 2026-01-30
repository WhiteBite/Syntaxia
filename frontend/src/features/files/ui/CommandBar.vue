<template>
  <div class="magic-bar-wrapper">
    <!-- MAGIC CONTROL BAR -->
    <div class="magic-bar">
      <!-- LEFT: Token Limit Selector -->
      <div class="limit-section" ref="limitRef">
        <BaseButton 
          variant="ghost" 
          size="sm" 
          class="limit-trigger"
          @click="toggleDropdown"
        >
          <span class="limit-value">{{ formatTokens(settings.maxTokens) }}</span>
          <ChevronDown class="limit-chevron" :class="{ open: showDropdown }" />
        </BaseButton>
        
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
              <BaseButton 
                v-if="customTokenValue"
                variant="success"
                size="xs"
                icon-only
                class="limit-custom-apply"
                @click.stop="applyCustomLimit"
              >
                <Check class="w-3.5 h-3.5" />
              </BaseButton>
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
          <Zap v-if="!isBuilding" class="build-icon" />
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
      {{ t('commandBar.selected') }}: 
      <BaseBadge variant="primary" size="xs" class="mx-1">
        {{ selectedCount }}
      </BaseBadge>
      <span class="token-estimate">~{{ estimatedTokens }}k tokens</span>
      <BaseButton 
        variant="ghost" 
        size="xs" 
        class="clear-btn ml-2"
        @click="handleClear"
        :title="t('files.clearSelection')"
      >
        <template #icon><X class="w-3 h-3" /></template>
        {{ t('files.clear') }}
      </BaseButton>
    </div>

  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useFileStore } from '@/features/files/model/file.store'
import { useSettingsStore } from '@/stores/settings.store'
import { BaseButton, BaseBadge } from '@/components/ui'
import { Check, ChevronDown, Zap, X } from 'lucide-vue-next'
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
  height: calc(56px * var(--ui-scale));
  background: #1c1f2e;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-xl);
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
  padding: 0 var(--space-3);
}

.limit-trigger {
  width: 100%;
  gap: var(--space-2);
}

.limit-value {
  font-size: var(--font-size-lg);
  font-family: ui-monospace, monospace;
  font-weight: 700;
  color: #e5e7eb;
  line-height: 1;
}

.limit-chevron {
  width: calc(14px * var(--ui-scale));
  height: calc(14px * var(--ui-scale));
  color: #6b7280;
  transition: all 0.2s;
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
  border-radius: var(--radius-md);
  border: none;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.15s ease-out;
  box-shadow: 0 4px 20px rgba(147, 51, 234, 0.35);
}

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
  border-radius: var(--radius-md);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

/* Bottom Glow */
.build-glow {
  position: absolute;
  bottom: calc(-10px * var(--ui-scale));
  left: 50%;
  transform: translateX(-50%);
  width: 70%;
  height: calc(24px * var(--ui-scale));
  background: linear-gradient(90deg, #db2777, #9333ea);
  filter: blur(calc(20px * var(--ui-scale)));
  opacity: 0.6;
  transition: opacity 0.15s ease-out;
}

.build-section:hover:not(:disabled) .build-glow {
  opacity: 0.9;
}

/* Content */
.build-content {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  height: 100%;
  color: white;
  z-index: 1;
}

.build-icon {
  width: calc(20px * var(--ui-scale));
  height: calc(20px * var(--ui-scale));
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
}

.build-spinner {
  width: calc(20px * var(--ui-scale));
  height: calc(20px * var(--ui-scale));
  animation: spin 1s linear infinite;
}

.build-text {
  font-size: var(--font-size-md);
  font-weight: 700;
  letter-spacing: 0.05em;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

/* FILE COUNTER */
.file-counter {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-xs);
  color: #6b7280;
}

.token-estimate {
  color: #6b7280;
  margin-left: var(--space-2);
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
