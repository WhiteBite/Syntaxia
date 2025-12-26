<template>
  <div class="cut-line" :id="`chunk-${chunkEnd}`">
    <div class="cut-line-dashes"></div>
    <div class="cut-line-badge">
      <span class="cut-icon">✂️</span>
      <span class="cut-text">{{ t('chunks.endPart', { n: chunkEnd }) }} — {{ t('chunks.startPart', { n: chunkEnd + 1 }) }}</span>
    </div>
    <div class="cut-line-dashes"></div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'

defineProps<{
  chunkEnd: number
}>()

const { t } = useI18n()
</script>

<style scoped>
.cut-line {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 0;
  margin: 8px 0;
}

.cut-line-dashes {
  flex: 1;
  height: 1px;
  background: repeating-linear-gradient(
    90deg,
    rgba(168, 85, 247, 0.4) 0px,
    rgba(168, 85, 247, 0.4) 6px,
    transparent 6px,
    transparent 12px
  );
}

.cut-line-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  background: rgba(168, 85, 247, 0.1);
  border: 1px solid rgba(168, 85, 247, 0.25);
  border-radius: 20px;
  white-space: nowrap;
}

.cut-icon {
  font-size: 12px;
}

.cut-text {
  font-size: 10px;
  font-weight: 500;
  color: #c084fc;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
</style>
