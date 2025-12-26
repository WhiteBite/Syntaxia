<template>
  <Teleport to="body">
    <Transition name="popover">
      <div v-if="visible" class="stats-popover-overlay" @click.self="$emit('close')">
        <div class="stats-popover" :style="popoverStyle">
          <!-- Header -->
          <div class="stats-popover-header">
            <div class="stats-popover-title">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
              <span>{{ t('stats.title') }}</span>
            </div>
            <button @click="$emit('close')" class="stats-popover-close">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Main Stats Grid -->
          <div class="stats-popover-grid">
            <div class="stats-popover-card">
              <div class="stats-popover-value">{{ fileCount }}</div>
              <div class="stats-popover-label">{{ t('context.files') }}</div>
            </div>
            <div class="stats-popover-card">
              <div class="stats-popover-value">{{ formatNumber(lineCount) }}</div>
              <div class="stats-popover-label">{{ t('context.lines') }}</div>
            </div>
            <div class="stats-popover-card">
              <div class="stats-popover-value stats-popover-value--accent">{{ formatNumber(tokenCount) }}</div>
              <div class="stats-popover-label">{{ t('action.tokens') }}</div>
            </div>
            <div class="stats-popover-card">
              <div class="stats-popover-value stats-popover-value--success">${{ estimatedCost.toFixed(4) }}</div>
              <div class="stats-popover-label">{{ t('action.cost') }}</div>
            </div>
          </div>

          <!-- File Types Section -->
          <div v-if="fileTypeStats.length > 0" class="stats-popover-section">
            <div class="stats-popover-section-title">{{ t('stats.byType') }}</div>
            <div class="stats-popover-list">
              <div v-for="stat in fileTypeStats" :key="stat.extension" class="stats-popover-item">
                <span class="stats-popover-item-icon">{{ stat.icon }}</span>
                <span class="stats-popover-item-name">.{{ stat.extension }}</span>
                <span class="stats-popover-item-count">{{ stat.count }}</span>
                <div class="stats-popover-item-bar">
                  <div class="stats-popover-item-fill" :class="stat.colorClass" :style="{ width: `${stat.percentage}%` }"></div>
                </div>
              </div>
            </div>
          </div>

          <!-- Top Folders Section -->
          <div v-if="folderStats.length > 0" class="stats-popover-section">
            <div class="stats-popover-section-title">{{ t('stats.byFolder') }}</div>
            <div class="stats-popover-list">
              <div v-for="stat in folderStats.slice(0, 5)" :key="stat.folder" class="stats-popover-item">
                <svg class="w-4 h-4 text-blue-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                </svg>
                <span class="stats-popover-item-name">{{ stat.folder || '/' }}</span>
                <span class="stats-popover-item-count">{{ stat.count }}</span>
                <div class="stats-popover-item-bar">
                  <div class="stats-popover-item-fill stats-popover-item-fill--folder" :style="{ width: `${stat.percentage}%` }"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useContextStore } from '@/features/context/model/context.store'
import { useFileStore } from '@/features/files/model/file.store'
import { FILE_TYPE_CONFIG } from '@/utils/fileIcons'
import { computed } from 'vue'

defineProps<{
  visible: boolean
  anchorRect?: DOMRect | null
}>()

defineEmits<{
  close: []
}>()

const { t } = useI18n()
const contextStore = useContextStore()
const fileStore = useFileStore()

const fileCount = computed(() => contextStore.fileCount)
const lineCount = computed(() => contextStore.lineCount)
const tokenCount = computed(() => contextStore.tokenCount)
const estimatedCost = computed(() => contextStore.estimatedCost)

const popoverStyle = computed(() => ({
  // Center in viewport
}))

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
    .slice(0, 6)
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

<style scoped>
.stats-popover-overlay {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 80px;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(2px);
}

.stats-popover {
  width: 100%;
  max-width: 380px;
  background: linear-gradient(180deg, #1c1f2e 0%, #161922 100%);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  box-shadow: 
    0 24px 48px rgba(0, 0, 0, 0.5),
    0 0 0 1px rgba(255, 255, 255, 0.05) inset;
  overflow: hidden;
}

.stats-popover-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.stats-popover-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  font-weight: 600;
  color: white;
}

.stats-popover-title svg {
  color: #8b5cf6;
}

.stats-popover-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: rgba(255, 255, 255, 0.05);
  border: none;
  border-radius: 8px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.stats-popover-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

/* Main Stats Grid */
.stats-popover-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 16px 20px;
}

.stats-popover-card {
  padding: 14px;
  background: rgba(0, 0, 0, 0.25);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  text-align: center;
}

.stats-popover-value {
  font-size: 20px;
  font-weight: 700;
  color: white;
  font-variant-numeric: tabular-nums;
}

.stats-popover-value--accent {
  color: #a78bfa;
}

.stats-popover-value--success {
  color: #34d399;
}

.stats-popover-label {
  font-size: 11px;
  font-weight: 500;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-top: 4px;
}

/* Sections */
.stats-popover-section {
  padding: 16px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.stats-popover-section-title {
  font-size: 10px;
  font-weight: 700;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  margin-bottom: 12px;
}

.stats-popover-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.stats-popover-item {
  display: grid;
  grid-template-columns: auto 1fr auto 60px;
  align-items: center;
  gap: 10px;
}

.stats-popover-item-icon {
  font-size: 14px;
}

.stats-popover-item-name {
  font-size: 12px;
  color: #e5e7eb;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.stats-popover-item-count {
  font-size: 11px;
  font-weight: 600;
  color: #9ca3af;
  font-variant-numeric: tabular-nums;
}

.stats-popover-item-bar {
  height: 4px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 2px;
  overflow: hidden;
}

.stats-popover-item-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.3s ease-out;
}

.stats-popover-item-fill--folder {
  background: linear-gradient(90deg, #3b82f6 0%, #60a5fa 100%);
}

/* Transitions */
.popover-enter-active,
.popover-leave-active {
  transition: opacity 0.15s ease-out;
}

.popover-enter-active .stats-popover,
.popover-leave-active .stats-popover {
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.15s ease-out;
}

.popover-enter-from,
.popover-leave-to {
  opacity: 0;
}

.popover-enter-from .stats-popover,
.popover-leave-to .stats-popover {
  transform: translateY(-10px) scale(0.98);
  opacity: 0;
}
</style>
