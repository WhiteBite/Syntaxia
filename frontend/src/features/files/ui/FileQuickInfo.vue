<template>
  <div class="file-quick-info" v-if="info || isLoading">
    <div v-if="isLoading" class="flex items-center justify-center py-2">
      <svg class="animate-spin h-4 w-4 text-indigo-400" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
    </div>
    
    <template v-else-if="info">
      <div class="info-header">
        <svg class="w-3.5 h-3.5 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <span class="text-xs font-medium text-gray-300">{{ t('files.quickInfo') }}</span>
      </div>
      
      <div class="info-grid">
        <div class="info-item">
          <span class="info-label">{{ t('files.symbols') }}</span>
          <span class="info-value">{{ info.symbolCount }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">{{ t('files.imports') }}</span>
          <span class="info-value">{{ info.importCount }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">{{ t('files.dependents') }}</span>
          <span class="info-value">{{ info.dependentCount }}</span>
        </div>
      </div>
      
      <div class="risk-section">
        <div class="flex items-center justify-between mb-1">
          <span class="info-label">{{ t('files.changeRisk') }}</span>
          <span class="text-xs" :class="getRiskColor(info.riskLevel)">
            {{ t(`files.risk.${info.riskLevel}`) }}
          </span>
        </div>
        <div class="risk-bar-bg">
          <div class="risk-bar" :class="getRiskBarClass(info.riskLevel)" 
            :style="{ width: getRiskWidth(info.changeRisk) }" />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useFileQuickInfo, type FileQuickInfo } from '@/composables/useFileQuickInfo';
import { useI18n } from '@/composables/useI18n';
import { computed, onMounted, ref, watch } from 'vue';

const props = defineProps<{
  projectPath: string
  filePath: string
}>()

const { t } = useI18n()
const { getInfo, isLoading: checkLoading, getRiskColor, getRiskWidth } = useFileQuickInfo()

const info = ref<FileQuickInfo | null>(null)
const isLoading = computed(() => checkLoading(props.projectPath, props.filePath))

async function loadInfo() {
  if (!props.projectPath || !props.filePath) return
  info.value = await getInfo(props.projectPath, props.filePath)
}

function getRiskBarClass(level: string): string {
  switch (level) {
    case 'low': return 'bg-green-500'
    case 'medium': return 'bg-amber-500'
    case 'high': return 'bg-red-500'
    default: return 'bg-gray-500'
  }
}

watch(() => [props.projectPath, props.filePath], loadInfo, { immediate: true })
onMounted(loadInfo)
</script>

<style scoped>
.file-quick-info {
  @apply p-2 rounded-lg;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  min-width: 160px;
}

.info-header {
  @apply flex items-center gap-1.5 mb-2 pb-1.5;
  border-bottom: 1px solid var(--border-default);
}

.info-grid {
  @apply grid grid-cols-3 gap-2 mb-2;
}

.info-item {
  @apply flex flex-col items-center;
}

.info-label {
  @apply text-xs;
  color: var(--text-secondary);
}

.info-value {
  @apply text-sm font-medium;
  color: var(--text-primary);
}

.risk-section {
  @apply pt-1.5;
  border-top: 1px solid var(--border-default);
}

.risk-bar-bg {
  @apply h-1.5 rounded-full overflow-hidden;
  background: var(--bg-3);
}

.risk-bar {
  @apply h-full rounded-full transition-all duration-300;
}
</style>
