<template>
  <Teleport to="body">
    <Transition name="explain-btn">
      <button
        v-if="visible && selectedCode"
        ref="buttonRef"
        class="explain-code-btn"
        :style="buttonStyle"
        :title="t('chat.explain.tooltip')"
        @click="handleExplain"
      >
        <LightBulbIcon class="w-4 h-4" />
        <span class="explain-code-btn__text">{{ t('chat.explain.button') }}</span>
      </button>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { LightBulbIcon } from '@heroicons/vue/24/outline'
import { computed } from 'vue'
import { useCodeExplanation } from '../composables/useCodeExplanation'

const props = defineProps<{
  visible: boolean
  selectedCode: string
  language?: string
  position: { x: number; y: number }
}>()

const emit = defineEmits<{
  (e: 'explained'): void
}>()

const { t } = useI18n()
const { explainCode } = useCodeExplanation()

const buttonStyle = computed(() => ({
  left: `${props.position.x}px`,
  top: `${props.position.y}px`
}))

async function handleExplain(): Promise<void> {
  if (!props.selectedCode.trim()) return
  
  await explainCode(props.selectedCode, props.language || 'text')
  emit('explained')
}
</script>

<style scoped>
.explain-code-btn {
  position: fixed;
  z-index: 1100;
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.625rem;
  background: var(--bg-2, #1f2937);
  border: 1px solid var(--border-1, #374151);
  border-radius: 0.375rem;
  color: var(--text-1, #f3f4f6);
  font-size: 0.75rem;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  transition: all 0.15s ease;
}

.explain-code-btn:hover {
  background: var(--bg-3, #374151);
  border-color: var(--accent, #8b5cf6);
}

.explain-code-btn__text {
  white-space: nowrap;
}

.explain-btn-enter-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.explain-btn-enter-from {
  opacity: 0;
  transform: translateY(-4px) scale(0.95);
}

.explain-btn-leave-active {
  transition: all 0.15s ease-out;
}

.explain-btn-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
