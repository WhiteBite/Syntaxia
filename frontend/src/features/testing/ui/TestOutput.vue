<template>
  <div class="test-output">
    <!-- Header -->
    <div class="output-header">
      <h3 class="output-title">{{ t('testing.output') }}</h3>
      <BaseButton
        v-if="hasContent"
        variant="ghost"
        size="sm"
        icon-only
        :title="t('testing.copy')"
        @click="copyOutput"
      >
        <template #icon>
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none">
            <rect x="9" y="9" width="13" height="13" rx="2" stroke="currentColor" stroke-width="2"/>
            <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" stroke="currentColor" stroke-width="2"/>
          </svg>
        </template>
      </BaseButton>
    </div>

    <!-- Empty State -->
    <div v-if="!hasContent" class="output-empty">
      <svg class="empty-icon" viewBox="0 0 24 24" fill="none">
        <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <p class="empty-text">{{ t('testing.selectTestToSeeOutput') }}</p>
    </div>

    <!-- Output Content -->
    <div v-else class="output-content">
      <!-- Error Section -->
      <div v-if="error" class="output-section output-section--error">
        <div class="section-header">
          <svg class="section-icon" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
            <path d="M15 9l-6 6M9 9l6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <span class="section-title">{{ t('testing.error') }}</span>
        </div>
        <pre class="section-content">{{ error }}</pre>
      </div>

      <!-- Stack Trace Section -->
      <div v-if="stackTrace && stackTrace !== error" class="output-section output-section--stack">
        <div class="section-header">
          <svg class="section-icon" viewBox="0 0 24 24" fill="none">
            <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" stroke="currentColor" stroke-width="2"/>
            <path d="M8 12h8M8 16h5" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <span class="section-title">{{ t('testing.stackTrace') }}</span>
        </div>
        <pre class="section-content section-content--stack">{{ stackTrace }}</pre>
      </div>

      <!-- Stdout Section -->
      <div v-if="output" class="output-section">
        <div class="section-header">
          <svg class="section-icon" viewBox="0 0 24 24" fill="none">
            <path d="M4 17l6-6-6-6M12 19h8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <span class="section-title">{{ t('testing.stdout') }}</span>
        </div>
        <pre class="section-content">{{ output }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '@/components/ui/BaseButton.vue'
import { useI18n } from '@/composables/useI18n'
import { useUIStore } from '@/stores/ui.store'
import { computed } from 'vue'

const { t } = useI18n()
const uiStore = useUIStore()

const props = defineProps<{
  output?: string
  error?: string
  stackTrace?: string
}>()

const hasContent = computed(() => 
  Boolean(props.output || props.error || props.stackTrace)
)

function copyOutput(): void {
  const content = [
    props.error && `Error:\n${props.error}`,
    props.stackTrace && `Stack Trace:\n${props.stackTrace}`,
    props.output && `Output:\n${props.output}`
  ].filter(Boolean).join('\n\n')

  navigator.clipboard.writeText(content)
  uiStore.addToast(t('testing.copied'), 'success')
}
</script>

<style scoped>
.test-output {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-secondary, #161922);
  border-radius: 8px;
  overflow: hidden;
}

.output-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.output-title {
  font-size: 13px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.output-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  padding: 32px;
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

.output-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.output-section {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 8px;
  overflow: hidden;
}

.output-section--error {
  border-color: rgba(239, 68, 68, 0.3);
}

.output-section--stack {
  border-color: rgba(245, 158, 11, 0.3);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.output-section--error .section-header {
  background: rgba(239, 68, 68, 0.1);
}

.output-section--stack .section-header {
  background: rgba(245, 158, 11, 0.1);
}

.section-icon {
  width: 16px;
  height: 16px;
  color: #9ca3af;
}

.output-section--error .section-icon {
  color: #ef4444;
}

.output-section--stack .section-icon {
  color: #f59e0b;
}

.section-title {
  font-size: 12px;
  font-weight: 600;
  color: #d1d5db;
}

.section-content {
  padding: 12px;
  margin: 0;
  font-family: ui-monospace, monospace;
  font-size: 12px;
  line-height: 1.6;
  color: #d1d5db;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 30vh;
  overflow-y: auto;
}

.section-content--stack {
  color: #fbbf24;
}
</style>
