<template>
  <div class="sparklines flex items-center gap-2 h-4 px-2 select-none" :title="tooltip">
    <!-- Memory Sparkline -->
    <div class="flex items-center gap-1.5 h-full">
      <span class="text-[8px] font-bold text-gray-500 uppercase tracking-tighter">Mem</span>
      <svg class="w-12 h-3" preserveAspectRatio="none">
        <polyline
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          :class="memoryColorClass"
          :points="memoryPoints"
        />
      </svg>
    </div>

    <div class="w-px h-2 bg-white/10"></div>

    <!-- Token Sparkline -->
    <div class="flex items-center gap-1.5 h-full">
      <span class="text-[8px] font-bold text-gray-500 uppercase tracking-tighter">Tok</span>
      <svg class="w-12 h-3" preserveAspectRatio="none">
        <polyline
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="text-indigo-500"
          :points="tokenPoints"
        />
      </svg>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useMemoryMonitor } from '@/utils/memory-monitor'
import { useFileStore } from '@/features/files/model/file.store'

const memoryMonitor = useMemoryMonitor()
const fileStore = useFileStore()

const MAX_POINTS = 20
const memHistory = ref<number[]>([])
const tokenHistory = ref<number[]>([])
const currentMem = ref(0)

const memoryPoints = computed(() => generatePoints(memHistory.value, 48, 12))
const tokenPoints = computed(() => generatePoints(tokenHistory.value, 48, 12))

const memoryColorClass = computed(() => {
  if (currentMem.value > 80) return 'text-red-500'
  if (currentMem.value > 60) return 'text-amber-500'
  return 'text-emerald-500'
})

const tooltip = computed(() => {
  const mem = memHistory.value[memHistory.value.length - 1] || 0
  const tok = tokenHistory.value[tokenHistory.value.length - 1] || 0
  return `Memory: ${mem}% | Tokens: ${tok}k`
})

function generatePoints(data: number[], width: number, height: number): string {
  if (data.length < 2) return ''
  const max = Math.max(...data, 1)
  const step = width / (MAX_POINTS - 1)
  
  return data
    .map((val, i) => {
      const x = i * step
      const y = height - (val / max) * height
      return `${x},${y}`
    })
    .join(' ')
}

let interval: number | null = null

onMounted(() => {
  interval = window.setInterval(async () => {
    const stats = await memoryMonitor.getMemoryStats()
    if (stats) {
      currentMem.value = stats.percentage
      memHistory.value.push(stats.percentage)
      if (memHistory.value.length > MAX_POINTS) memHistory.value.shift()
    }
    
    const tokens = Math.round(fileStore.estimatedTokenCount / 1000)
    tokenHistory.value.push(tokens)
    if (tokenHistory.value.length > MAX_POINTS) tokenHistory.value.shift()
  }, 5000)
})

onUnmounted(() => {
  if (interval) clearInterval(interval)
})
</script>

<style scoped>
.sparklines svg {
  overflow: visible;
}
</style>
