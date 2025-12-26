<template>
  <div class="model-limits-indicator">
    <!-- Compact Mode (default) -->
    <div 
      v-if="!expanded" 
      class="limits-compact"
      :class="[pulseClass, glowClass]"
      @click="expanded = true"
    >
      <div class="limits-bar-container">
        <div 
          class="limits-bar limits-bar-animate" 
          :style="{ width: usagePercent + '%' }"
          :class="usageClass"
        ></div>
      </div>
      <div class="limits-info">
        <span class="limits-used">{{ formatTokens(usedTokens) }}</span>
        <span class="limits-separator">/</span>
        <span class="limits-total">{{ formatTokens(contextLimit) }}</span>
        <span class="limits-label">{{ t('limits.tokens') }}</span>
        <span v-if="usagePercent >= 80" class="limits-percent" :class="usageClass">
          {{ usagePercent }}%
        </span>
      </div>
    </div>

    <!-- Expanded Mode -->
    <Transition name="limits-expand">
      <div v-if="expanded" class="limits-expanded" :class="glowClass">
        <div class="limits-header">
          <div class="limits-model-info">
            <span class="limits-model-name">{{ modelName }}</span>
            <span class="limits-provider-badge">{{ providerName }}</span>
          </div>
          <button class="icon-btn-sm" @click.stop="expanded = false">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
            </svg>
          </button>
        </div>

        <!-- Progress Ring -->
        <div class="limits-ring-container">
          <svg class="limits-ring" viewBox="0 0 100 100">
            <circle 
              class="limits-ring-bg" 
              cx="50" cy="50" r="42" 
              fill="none" 
              stroke-width="8"
            />
            <circle 
              class="limits-ring-progress limits-ring-animate" 
              :class="usageClass"
              cx="50" cy="50" r="42" 
              fill="none" 
              stroke-width="8"
              :stroke-dasharray="circumference"
              :stroke-dashoffset="dashOffset"
              stroke-linecap="round"
            />
          </svg>
          <div class="limits-ring-text">
            <span class="limits-ring-percent stats-counter">{{ usagePercent }}%</span>
            <span class="limits-ring-label">{{ t('limits.used') }}</span>
          </div>
        </div>

        <!-- Stats Grid -->
        <div class="limits-stats">
          <div class="limits-stat stats-counter">
            <span class="limits-stat-value">{{ formatTokens(usedTokens) }}</span>
            <span class="limits-stat-label">{{ t('limits.contextUsed') }}</span>
          </div>
          <div class="limits-stat stats-counter" style="animation-delay: 0.1s">
            <span class="limits-stat-value">{{ formatTokens(availableTokens) }}</span>
            <span class="limits-stat-label">{{ t('limits.available') }}</span>
          </div>
          <div class="limits-stat stats-counter" style="animation-delay: 0.2s">
            <span class="limits-stat-value">{{ formatTokens(contextLimit) }}</span>
            <span class="limits-stat-label">{{ t('limits.maxContext') }}</span>
          </div>
          <div class="limits-stat stats-counter" style="animation-delay: 0.3s">
            <span class="limits-stat-value limits-stat-cost">${{ estimatedCost }}</span>
            <span class="limits-stat-label">{{ t('limits.estCost') }}</span>
          </div>
        </div>

        <!-- Warning if near limit -->
        <Transition name="slide-fade">
          <div v-if="usagePercent >= 80" class="limits-warning" :class="[usageClass, pulseClass]">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <span>{{ usagePercent >= 95 ? t('limits.criticalWarning') : t('limits.nearLimit') }}</span>
          </div>
        </Transition>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';
import { MODEL_LIMITS, MODEL_PRICING } from '@/stores/ai.store';
import { computed, ref } from 'vue';

const { t } = useI18n()

const props = defineProps<{
  usedTokens: number
  modelName: string
  providerName: string
}>()

const expanded = ref(false)

const contextLimit = computed(() => {
  return MODEL_LIMITS[props.modelName] || 32000
})

const usagePercent = computed(() => {
  if (!props.usedTokens || !contextLimit.value || contextLimit.value === 0) return 0
  const percent = Math.round((props.usedTokens / contextLimit.value) * 100)
  if (!isFinite(percent) || isNaN(percent)) return 0
  return Math.min(Math.max(percent, 0), 100)
})

const availableTokens = computed(() => {
  return Math.max(0, contextLimit.value - props.usedTokens)
})

const usageClass = computed(() => {
  if (usagePercent.value >= 95) return 'usage-critical'
  if (usagePercent.value >= 80) return 'usage-warning'
  return 'usage-safe'
})

const pulseClass = computed(() => {
  if (usagePercent.value >= 95) return 'limits-critical-pulse'
  if (usagePercent.value >= 80) return 'limits-warning-pulse'
  return ''
})

const glowClass = computed(() => {
  if (usagePercent.value >= 95) return 'limits-glow-critical'
  if (usagePercent.value >= 80) return 'limits-glow-warning'
  return ''
})

const estimatedCost = computed(() => {
  const pricePerK = MODEL_PRICING[props.modelName] || 0.001
  return ((props.usedTokens / 1000) * pricePerK).toFixed(4)
})

// SVG circle calculations
const circumference = 2 * Math.PI * 42
const dashOffset = computed(() => {
  return circumference - (usagePercent.value / 100) * circumference
})

function formatTokens(tokens: number): string {
  if (tokens >= 1000000) return (tokens / 1000000).toFixed(1) + 'M'
  if (tokens >= 1000) return (tokens / 1000).toFixed(1) + 'K'
  return tokens.toString()
}
</script>

<style scoped>
.model-limits-indicator {
  @apply mb-3;
}

/* Compact Mode */
.limits-compact {
  @apply flex items-center gap-3 p-2 bg-gray-800/40 rounded-lg border border-gray-700/30 cursor-pointer;
  @apply hover:bg-gray-800/60 hover:border-gray-600/50 transition-all;
}

.limits-bar-container {
  @apply flex-1 h-1.5 bg-gray-700/50 rounded-full overflow-hidden;
}

.limits-bar {
  @apply h-full rounded-full transition-all duration-500;
}

.limits-bar.usage-safe { @apply bg-emerald-500; }
.limits-bar.usage-warning { @apply bg-amber-500; }
.limits-bar.usage-critical { @apply bg-red-500; }

.limits-info {
  @apply flex items-center gap-1 text-[10px] text-gray-400 flex-shrink-0;
}

.limits-used { @apply text-white font-medium; }
.limits-separator { @apply text-gray-600; }
.limits-total { @apply text-gray-400; }
.limits-label { @apply text-gray-600 ml-0.5; }

.limits-percent {
  @apply ml-1 px-1 py-0.5 rounded text-[9px] font-bold;
}
.limits-percent.usage-warning { @apply bg-amber-500/20 text-amber-400; }
.limits-percent.usage-critical { @apply bg-red-500/20 text-red-400; }

/* Expanded Mode */
.limits-expanded {
  @apply p-3 bg-gray-800/50 rounded-xl border border-gray-700/30;
  @apply transition-all;
}

.limits-header {
  @apply flex items-center justify-between mb-3;
}

.limits-model-info { @apply flex items-center gap-2; }
.limits-model-name { @apply text-xs font-medium text-white; }
.limits-provider-badge { @apply px-1.5 py-0.5 text-[10px] bg-indigo-500/20 text-indigo-300 rounded; }

/* Progress Ring */
.limits-ring-container { @apply relative w-24 h-24 mx-auto mb-3; }
.limits-ring { @apply w-full h-full transform -rotate-90; }
.limits-ring-bg { @apply stroke-gray-700/50; }

.limits-ring-progress { @apply transition-all duration-700 ease-out; }
.limits-ring-progress.usage-safe { @apply stroke-emerald-500; }
.limits-ring-progress.usage-warning { @apply stroke-amber-500; }
.limits-ring-progress.usage-critical { @apply stroke-red-500; }

.limits-ring-text { @apply absolute inset-0 flex flex-col items-center justify-center; }
.limits-ring-percent { @apply text-lg font-bold text-white; }
.limits-ring-label { @apply text-[10px] text-gray-400; }

/* Stats Grid */
.limits-stats { @apply grid grid-cols-2 gap-2; }
.limits-stat { @apply text-center p-2 bg-gray-900/50 rounded-lg; }
.limits-stat-value { @apply block text-sm font-semibold text-white; }
.limits-stat-cost { @apply text-emerald-400; }
.limits-stat-label { @apply block text-[10px] text-gray-400 mt-0.5; }

/* Warning */
.limits-warning {
  @apply flex items-center gap-2 mt-3 p-2 rounded-lg text-xs;
}
.limits-warning.usage-warning { @apply bg-amber-500/10 text-amber-300 border border-amber-500/30; }
.limits-warning.usage-critical { @apply bg-red-500/10 text-red-300 border border-red-500/30; }
</style>
