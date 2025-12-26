<template>
  <div class="space-y-4">
    <!-- Summary Stats -->
    <div class="context-stats">
      <p class="text-xs font-semibold text-gray-400 mb-3">{{ t('stats.totalFiles') }}</p>
      <div class="stats-grid list-stagger">
        <div class="stat-card-enhanced">
          <div class="stat-card-value">{{ fileCount }}</div>
          <div class="stat-label">{{ t('tools.files') }}</div>
        </div>
        <div class="stat-card-enhanced">
          <div class="stat-card-value">{{ formatNumber(lineCount) }}</div>
          <div class="stat-label">{{ t('context.lines') }}</div>
        </div>
        <div class="stat-card-enhanced">
          <div class="stat-card-value text-indigo-400">{{ formatNumber(tokenCount) }}</div>
          <div class="stat-label">{{ t('action.tokens') }}</div>
        </div>
        <div class="stat-card-enhanced">
          <div class="stat-card-value text-emerald-400">${{ estimatedCost.toFixed(4) }}</div>
          <div class="stat-label">{{ t('action.cost') }}</div>
        </div>
      </div>
    </div>

    <!-- By File Type -->
    <div class="context-stats" v-if="fileTypeStats.length > 0">
      <p class="text-xs font-semibold text-gray-400 mb-3">{{ t('stats.byType') }}</p>
      <div class="space-y-2">
        <div v-for="stat in fileTypeStats" :key="stat.extension" class="flex items-center gap-2">
          <span class="text-lg flex-shrink-0">{{ stat.icon }}</span>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between text-xs mb-1">
              <span class="text-gray-300 truncate">.{{ stat.extension }}</span>
              <span class="text-gray-400">{{ stat.count }} {{ t('context.files') }}</span>
            </div>
            <div class="progress-bar-enhanced">
              <div class="progress-bar-fill" :class="stat.colorClass" :style="{ width: `${stat.percentage}%` }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- By Folder -->
    <div class="context-stats" v-if="folderStats.length > 0">
      <p class="text-xs font-semibold text-gray-400 mb-3">{{ t('stats.byFolder') }}</p>
      <div class="space-y-2">
        <div v-for="stat in folderStats.slice(0, 5)" :key="stat.folder" class="flex items-center gap-2">
          <svg class="w-4 h-4 text-blue-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between text-xs mb-1">
              <span class="text-gray-300 truncate">{{ stat.folder || '/' }}</span>
              <span class="text-gray-400">{{ stat.count }}</span>
            </div>
            <div class="progress-bar-enhanced">
              <div class="progress-bar-fill-folder" :style="{ width: `${stat.percentage}%` }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="!hasContext" class="empty-state-enhanced py-8">
      <div class="empty-state-icon-glow !w-14 !h-14 mb-4">
        <svg class="w-7 h-7 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      </div>
      <p class="text-sm font-medium text-gray-300 mb-1">{{ t('context.notBuilt') }}</p>
      <p class="text-xs text-gray-400">{{ t('context.selectFiles') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useContextStore } from '@/features/context'
import { useFileStore } from '@/features/files'
import { FILE_TYPE_CONFIG } from '@/utils/fileIcons'
import { computed } from 'vue'

const { t } = useI18n()
const contextStore = useContextStore()
const fileStore = useFileStore()

const hasContext = computed(() => contextStore.hasContext)
const fileCount = computed(() => contextStore.fileCount)
const lineCount = computed(() => contextStore.lineCount)
const tokenCount = computed(() => contextStore.tokenCount)
const estimatedCost = computed(() => contextStore.estimatedCost)

const fileTypeStats = computed(() => {
  const selected = Array.from(fileStore.selectedPaths)
  if (selected.length === 0) return []

  const counts = new Map<string, number>()

  for (const path of selected) {
    const ext = path.split('.').pop()?.toLowerCase() || 'other'
    counts.set(ext, (counts.get(ext) || 0) + 1)
  }

  const total = selected.length
  return Array.from(counts.entries())
    .map(([extension, count]) => {
      const info = FILE_TYPE_CONFIG[extension] || FILE_TYPE_CONFIG.default
      return {
        extension,
        count,
        percentage: Math.round((count / total) * 100),
        icon: info.icon,
        colorClass: info.colorClass
      }
    })
    .sort((a, b) => b.count - a.count)
    .slice(0, 8)
})

const folderStats = computed(() => {
  const selected = Array.from(fileStore.selectedPaths)
  if (selected.length === 0) return []

  const counts = new Map<string, number>()

  for (const path of selected) {
    const parts = path.split('/')
    const folder = parts.length > 1 ? parts[0] : '/'
    counts.set(folder, (counts.get(folder) || 0) + 1)
  }

  const total = selected.length
  return Array.from(counts.entries())
    .map(([folder, count]) => ({
      folder,
      count,
      percentage: Math.round((count / total) * 100)
    }))
    .sort((a, b) => b.count - a.count)
})

function formatNumber(num: number): string {
  if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`
  if (num >= 1000) return `${(num / 1000).toFixed(1)}K`
  return num.toString()
}
</script>
