<template>
  <div class="bg-gray-800/90 backdrop-blur p-4 rounded-lg border border-gray-600 max-w-sm shadow-xl">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-lg font-semibold text-white flex items-center gap-2">
        🧠 Memory Monitor
      </h3>
      <button @click="$emit('close')" class="text-gray-400 hover:text-white transition-colors" title="Close">
        ✕
      </button>
    </div>

    <div v-if="stats" class="space-y-3 text-sm">
      <!-- Heap Usage -->
      <div>
        <div class="flex justify-between text-gray-300 mb-1">
          <span>Heap Usage</span>
          <span class="font-mono">{{ stats.used }}MB / {{ stats.total }}MB</span>
        </div>
        <div class="w-full bg-gray-700 rounded-full h-2">
          <div :class="[
            'h-2 rounded-full transition-all duration-300',
            stats.percentage >= 90 ? 'bg-red-500' :
              stats.percentage >= 75 ? 'bg-yellow-500' :
                'bg-green-500'
          ]" :style="{ width: `${stats.percentage}%` }"></div>
        </div>
        <div class="text-xs text-gray-400 mt-1">{{ stats.percentage }}%</div>
      </div>

      <!-- Store Metrics -->
      <div v-if="storeMetrics" class="border-t border-gray-700 pt-3 space-y-2">
        <div class="text-gray-300">
          <div class="font-semibold mb-1">Store Metrics</div>

          <div class="flex justify-between text-xs">
            <span>File Store:</span>
            <span class="font-mono">{{ storeMetrics.fileStore.nodesCount }} nodes</span>
          </div>

          <div class="flex justify-between text-xs">
            <span>Context Cache:</span>
            <span class="font-mono">{{ formatBytes(storeMetrics.contextStore.cacheSize) }}</span>
          </div>

          <div class="flex justify-between text-xs">
            <span>API Cache:</span>
            <span class="font-mono">{{ storeMetrics.apiCache.entries }} entries</span>
          </div>
        </div>
      </div>

      <!-- Diagnostics Info -->
      <div v-if="diagnosticsInfo" class="border-t border-gray-700 pt-3 space-y-1 text-xs">
        <div class="text-gray-400">Diagnostics:</div>
        <div class="flex justify-between">
          <span>Snapshots:</span>
          <span class="font-mono">{{ diagnosticsInfo.snapshotCount }}</span>
        </div>
        <div class="flex justify-between">
          <span>Issues:</span>
          <span class="font-mono" :class="diagnosticsInfo.criticalIssues > 0 ? 'text-red-400' : 'text-gray-300'">
            {{ diagnosticsInfo.criticalIssues }} critical, {{ diagnosticsInfo.warnings }} warnings
          </span>
        </div>
      </div>

      <!-- Actions -->
      <div class="border-t border-gray-700 pt-3 space-y-2">
        <div class="flex gap-2">
          <button @click="dumpSnapshot"
            class="flex-1 px-3 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded text-xs font-medium transition-colors"
            title="Create heap snapshot for Chrome DevTools">
            📸 Dump Heap
          </button>
          <button @click="forceCleanup"
            class="flex-1 px-3 py-2 bg-orange-600 hover:bg-orange-500 text-white rounded text-xs font-medium transition-colors"
            title="Force garbage collection and cleanup">
            🧹 Cleanup
          </button>
        </div>
        <button @click="saveReport"
          class="w-full px-3 py-2 bg-purple-600 hover:bg-purple-500 text-white rounded text-xs font-medium transition-colors"
          title="Save diagnostic report for AI analysis">
          💾 Save Report (for AI)
        </button>
      </div>

      <!-- Last Update -->
      <div class="text-xs text-gray-400 text-center">
        Updated {{ timeSinceUpdate }}s ago
      </div>
    </div>

    <div v-else class="text-gray-400 text-sm text-center py-4">
      Memory stats unavailable
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLogger } from '@/composables/useLogger'
import { useUIStore } from '@/stores/ui.store'
import { useMemoryDiagnostics } from '@/utils/memory-diagnostics'
import type { MemoryStats } from '@/utils/memory-monitor'
import { useMemoryMonitor } from '@/utils/memory-monitor'
import { computed, onMounted, onUnmounted, ref } from 'vue'

const logger = useLogger('MemoryDashboard')

defineEmits<{
  close: []
}>()

const memoryMonitor = useMemoryMonitor()
const diagnostics = useMemoryDiagnostics()
const uiStore = useUIStore()

const stats = ref<MemoryStats | null>(null)
const lastUpdate = ref(Date.now())
const timeSinceUpdate = ref(0)

const storeMetrics = computed(() => stats.value?.storeMetrics)

const diagnosticsInfo = computed(() => {
  const issues = diagnostics.getIssues()
  return {
    snapshotCount: diagnostics.getLatestSnapshot() ? diagnostics['snapshots']?.length || 0 : 0,
    criticalIssues: issues.filter(i => i.severity === 'critical').length,
    warnings: issues.filter(i => i.severity === 'warning').length
  }
})

let updateInterval: number | null = null
let timeInterval: number | null = null

async function updateStats() {
  stats.value = await memoryMonitor.getMemoryStats()
  lastUpdate.value = Date.now()
  timeSinceUpdate.value = 0
}

function dumpSnapshot() {
  memoryMonitor.dumpHeapSnapshot('manual-dashboard')
}

function forceCleanup() {
  memoryMonitor.forceCleanup()
  // Update stats after cleanup
  setTimeout(() => void updateStats(), 1000)
}

async function saveReport() {
  try {
    const filename = await diagnostics.saveReportToFile()
    uiStore.addToast(`Diagnostic report saved: ${filename}`, 'success')
    logger.debug('Report saved for AI analysis')
    logger.debug(diagnostics.getStatusSummary())
  } catch (e) {
    uiStore.addToast('Failed to save diagnostic report', 'error')
    logger.error('Failed to save report:', e)
  }
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 10) / 10 + ' ' + sizes[i]
}

onMounted(() => {
  void updateStats()

  // Start diagnostics collection (protected from double start)
  // In dev mode, diagnostics will automatically use longer intervals
  diagnostics.startCollection(30000)

  // In dev mode, use longer intervals to reduce memory overhead
  const isDevMode = import.meta.env.DEV
  const statsInterval = isDevMode ? 15000 : 5000 // 15s in dev, 5s in prod
  const timeCounterInterval = isDevMode ? 5000 : 2000 // 5s in dev, 2s in prod

  // Update stats periodically
  updateInterval = window.setInterval(() => void updateStats(), statsInterval)

  // Update time counter
  timeInterval = window.setInterval(() => {
    timeSinceUpdate.value = Math.floor((Date.now() - lastUpdate.value) / 1000)
  }, timeCounterInterval)
})

onUnmounted(() => {
  if (updateInterval) clearInterval(updateInterval)
  if (timeInterval) clearInterval(timeInterval)
  // Stop diagnostics when dashboard is closed to save resources
  diagnostics.stopCollection()
})
</script>
