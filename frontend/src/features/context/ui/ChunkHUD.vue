<template>
  <Transition name="hud">
    <div v-if="visible && totalChunks > 1" class="chunk-hud">
      <!-- Navigation -->
      <button 
        class="hud-nav-btn" 
        :disabled="currentChunk <= 1"
        @click="$emit('prev')"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>

      <!-- Progress Indicator -->
      <div class="hud-progress">
        <span class="hud-current">{{ currentChunk }}</span>
        <span class="hud-separator">/</span>
        <span class="hud-total">{{ totalChunks }}</span>
      </div>

      <!-- Copy Button -->
      <button 
        class="hud-copy-btn"
        :class="{ success: copySuccess }"
        @click="handleCopy"
      >
        <svg v-if="!copySuccess" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
        <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <span>{{ copySuccess ? t('chunks.copied') : t('chunks.copyPart', { n: currentChunk }) }}</span>
      </button>

      <!-- Navigation -->
      <button 
        class="hud-nav-btn" 
        :disabled="currentChunk >= totalChunks"
        @click="$emit('next')"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { ref } from 'vue'

defineProps<{
  visible: boolean
  currentChunk: number
  totalChunks: number
}>()

const emit = defineEmits<{
  (e: 'copy'): void
  (e: 'prev'): void
  (e: 'next'): void
}>()

const { t } = useI18n()
const copySuccess = ref(false)

function handleCopy() {
  emit('copy')
  copySuccess.value = true
  setTimeout(() => {
    copySuccess.value = false
  }, 1500)
}
</script>

<style scoped>
.chunk-hud {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(20, 24, 36, 0.92);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(0, 0, 0, 0.2);
  z-index: 20;
}

.hud-nav-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
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

.hud-progress {
  display: flex;
  align-items: baseline;
  gap: 2px;
  padding: 0 8px;
  font-family: ui-monospace, monospace;
}

.hud-current {
  font-size: 16px;
  font-weight: 700;
  color: #e5e7eb;
}

.hud-separator {
  font-size: 12px;
  color: #6b7280;
}

.hud-total {
  font-size: 12px;
  color: #6b7280;
}

.hud-copy-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: linear-gradient(135deg, #9333ea 0%, #a855f7 100%);
  border: none;
  border-radius: 10px;
  color: white;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease-out;
  box-shadow: 0 2px 8px rgba(147, 51, 234, 0.4);
}

.hud-copy-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(147, 51, 234, 0.5);
}

.hud-copy-btn.success {
  background: linear-gradient(135deg, #059669 0%, #10b981 100%);
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.4);
}

/* Transition */
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
