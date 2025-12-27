<template>
  <div class="build-output">
    <!-- Header -->
    <div class="output-header">
      <h3 class="output-title">{{ t('build.output') }}</h3>
      <div class="output-actions">
        <button
          class="action-btn"
          :disabled="!output"
          @click="$emit('copy')"
          :title="t('build.copyOutput')"
        >
          <svg class="action-icon" viewBox="0 0 24 24" fill="none">
            <rect x="9" y="9" width="13" height="13" rx="2" stroke="currentColor" stroke-width="2"/>
            <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" stroke="currentColor" stroke-width="2"/>
          </svg>
        </button>
        <button
          class="action-btn"
          :disabled="!output"
          @click="$emit('clear')"
          :title="t('build.clearOutput')"
        >
          <svg class="action-icon" viewBox="0 0 24 24" fill="none">
            <path d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Output Content -->
    <div 
      ref="outputContainer"
      class="output-content"
      :class="{ 'output-content--empty': !output }"
    >
      <div v-if="!output" class="output-empty">
        <svg class="empty-icon" viewBox="0 0 24 24" fill="none">
          <path d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <p class="empty-text">{{ t('build.noOutput') }}</p>
      </div>
      <pre v-else class="output-text"><code>{{ output }}</code></pre>
    </div>

    <!-- Auto-scroll toggle -->
    <div v-if="output" class="output-footer">
      <label class="auto-scroll-toggle">
        <input
          type="checkbox"
          v-model="autoScrollEnabled"
          class="toggle-input"
        />
        <span class="toggle-label">{{ t('build.autoScroll') }}</span>
      </label>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { nextTick, ref, watch } from 'vue'

const { t } = useI18n()

const props = defineProps<{
  output: string
}>()

defineEmits<{
  (e: 'copy'): void
  (e: 'clear'): void
}>()

const outputContainer = ref<HTMLElement | null>(null)
const autoScrollEnabled = ref(true)

// Auto-scroll when output changes
watch(() => props.output, async () => {
  if (autoScrollEnabled.value && outputContainer.value) {
    await nextTick()
    outputContainer.value.scrollTop = outputContainer.value.scrollHeight
  }
})
</script>

<style scoped>
.build-output {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-primary, #0f1117);
  border-radius: 8px;
  overflow: hidden;
}

.output-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.output-title {
  font-size: 13px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.output-actions {
  display: flex;
  gap: 4px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: none;
  border: 1px solid transparent;
  border-radius: 6px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.action-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.06);
  color: #e5e7eb;
}

.action-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.action-icon {
  width: 16px;
  height: 16px;
}

.output-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 12px 14px;
  background: rgba(0, 0, 0, 0.2);
}

.output-content--empty {
  display: flex;
  align-items: center;
  justify-content: center;
}

.output-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.empty-icon {
  width: 40px;
  height: 40px;
  color: #4b5563;
  margin-bottom: 12px;
}

.empty-text {
  font-size: 13px;
  color: #6b7280;
  margin: 0;
}

.output-text {
  font-family: ui-monospace, 'SF Mono', 'Cascadia Code', 'Source Code Pro', Menlo, Consolas, monospace;
  font-size: 12px;
  line-height: 1.6;
  color: #d1d5db;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.output-footer {
  padding: 8px 14px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.auto-scroll-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.toggle-input {
  width: 14px;
  height: 14px;
  accent-color: #8b5cf6;
}

.toggle-label {
  font-size: 12px;
  color: #9ca3af;
}
</style>
