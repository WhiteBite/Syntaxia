<template>
  <Transition name="hud">
    <div v-if="visible" class="unified-hud">
      <!-- Tools Section -->
      <div class="hud-section hud-tools">
        <Tooltip :text="t('context.clear')" position="top">
          <button @click="$emit('clear')" class="hud-tool-btn hud-tool-btn--danger">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </Tooltip>
        <Tooltip :text="t('context.export')" position="top">
          <button @click="$emit('export')" class="hud-tool-btn">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
          </button>
        </Tooltip>
      </div>

      <!-- Divider -->
      <div class="hud-divider" />

      <!-- Navigation Section (only if chunking) -->
      <template v-if="showChunkNav">
        <div class="hud-section hud-nav">
          <button 
            class="hud-nav-btn" 
            :disabled="currentChunk <= 1"
            @click="$emit('prev-chunk')"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <div class="hud-counter-wrap">
            <div class="hud-counter">
              <span class="hud-counter-current">{{ currentChunk }}</span>
              <span class="hud-counter-sep">/</span>
              <span class="hud-counter-total">{{ totalChunks }}</span>
            </div>
            <!-- Progress dots -->
            <div class="hud-progress-dots">
              <span 
                v-for="i in totalChunks" 
                :key="i"
                class="hud-dot"
                :class="{ 
                  active: i === currentChunk,
                  copied: isChunkCopied(i - 1)
                }"
              />
            </div>
          </div>
          <button 
            class="hud-nav-btn" 
            :disabled="currentChunk >= totalChunks"
            @click="$emit('next-chunk')"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
        <div class="hud-divider" />
      </template>

      <!-- Main Action Button -->
      <button 
        class="hud-action-btn"
        :class="{ 'hud-action-btn--success': copySuccess }"
        @click="$emit('copy')"
      >
        <svg v-if="!copySuccess" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
        <span>{{ buttonText }}</span>
      </button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import Tooltip from '@/components/ui/Tooltip.vue'
import { useI18n } from '@/composables/useI18n'
import { computed } from 'vue'

const props = defineProps<{
  visible: boolean
  showChunkNav: boolean
  currentChunk: number
  totalChunks: number
  copySuccess: boolean
  isChunkCopied: (index: number) => boolean
}>()

defineEmits<{
  'clear': []
  'export': []
  'copy': []
  'prev-chunk': []
  'next-chunk': []
}>()

const { t } = useI18n()

const buttonText = computed(() => {
  if (props.copySuccess) return '✓'
  if (props.showChunkNav) return `#${props.currentChunk}`
  return t('context.copy')
})
</script>

<style scoped>
/* UNIFIED HUD BAR - Glassmorphism Island */
.unified-hud {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: rgba(15, 17, 26, 0.85);
  backdrop-filter: blur(24px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 24px;
  box-shadow: 
    0 16px 48px rgba(0, 0, 0, 0.6),
    0 4px 16px rgba(0, 0, 0, 0.4);
  z-index: 20;
}

.hud-section {
  display: flex;
  align-items: center;
  gap: 4px;
}

.hud-divider {
  width: 1px;
  height: 20px;
  background: rgba(255, 255, 255, 0.1);
  margin: 0 4px;
}

/* Tool buttons */
.hud-tool-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: none;
  border-radius: 8px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.hud-tool-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.hud-tool-btn--danger:hover {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
}

/* Navigation */
.hud-nav-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: rgba(255, 255, 255, 0.05);
  border: none;
  border-radius: 6px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.hud-nav-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.hud-nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.hud-counter-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 0 8px;
}

.hud-counter {
  display: flex;
  align-items: baseline;
  gap: 2px;
  font-family: ui-monospace, monospace;
}

.hud-progress-dots {
  display: flex;
  gap: 4px;
}

.hud-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.15);
  transition: all 0.15s ease-out;
}

.hud-dot.active {
  background: #8b5cf6;
  box-shadow: 0 0 8px rgba(139, 92, 246, 0.6);
}

.hud-dot.copied {
  background: #10b981;
}

.hud-counter-current {
  font-size: 14px;
  font-weight: 700;
  color: #e5e7eb;
}

.hud-counter-sep {
  font-size: 11px;
  color: #4b5563;
}

.hud-counter-total {
  font-size: 11px;
  color: #6b7280;
}

/* Main action button - Most prominent */
.hud-action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
  border: none;
  border-radius: 12px;
  color: white;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease-out;
  box-shadow: 
    0 4px 16px rgba(139, 92, 246, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
}

.hud-action-btn:hover {
  transform: translateY(-1px);
  box-shadow: 
    0 6px 24px rgba(139, 92, 246, 0.5),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.hud-action-btn--success {
  background: linear-gradient(135deg, #059669 0%, #10b981 100%);
  box-shadow: 
    0 4px 16px rgba(16, 185, 129, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
}

/* HUD Transition */
.hud-enter-active,
.hud-leave-active {
  transition: all 0.25s ease-out;
}

.hud-enter-from,
.hud-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(16px);
}
</style>
