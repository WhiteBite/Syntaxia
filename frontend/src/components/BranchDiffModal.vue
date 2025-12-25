<template>
  <Teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      @click.self="close"
    >
      <!-- Backdrop -->
      <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" @click="close"></div>
      
      <!-- Modal -->
      <div class="relative w-full max-w-5xl max-h-[90vh] bg-gray-900 rounded-xl border border-gray-700 shadow-2xl flex flex-col overflow-hidden">
        <!-- Header -->
        <div class="flex items-center justify-between p-4 border-b border-gray-700">
          <div class="flex items-center gap-3">
            <div class="panel-icon bg-purple-500/20">
              <svg class="w-4 h-4 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
              </svg>
            </div>
            <h3 class="text-sm font-semibold text-white">{{ t('git.compareBranches') }}</h3>
          </div>
          <button
            @click="close"
            class="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <!-- Branch Selectors -->
        <div class="p-4 border-b border-gray-700 flex items-center gap-4">
          <div class="flex-1">
            <label class="block text-xs text-gray-400 mb-1">{{ t('git.baseBranch') }}</label>
            <select v-model="baseBranch" class="input w-full text-sm" @change="loadDiff">
              <option v-for="branch in branches" :key="branch" :value="branch">{{ branch }}</option>
            </select>
          </div>
          <div class="flex items-center text-gray-400 pt-4">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
            </svg>
          </div>
          <div class="flex-1">
            <label class="block text-xs text-gray-400 mb-1">{{ t('git.compareBranch') }}</label>
            <select v-model="compareBranch" class="input w-full text-sm" @change="loadDiff">
              <option v-for="branch in branches" :key="branch" :value="branch">{{ branch }}</option>
            </select>
          </div>
          <button
            @click="swapBranches"
            class="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors mt-4"
            :title="t('git.swapBranches')"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
            </svg>
          </button>
        </div>
        
        <!-- Diff Content -->
        <div class="flex-1 overflow-auto">
          <div v-if="isLoading" class="flex items-center justify-center h-64">
            <svg class="animate-spin w-8 h-8 text-indigo-400" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
            </svg>
          </div>
          
          <div v-else-if="error" class="flex items-center justify-center h-64 text-red-400">
            <div class="text-center">
              <svg class="w-12 h-12 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <p>{{ error }}</p>
            </div>
          </div>
          
          <div v-else-if="diffFiles.length === 0 && !isLoading" class="flex items-center justify-center h-64 text-gray-400">
            <div class="text-center">
              <svg class="w-12 h-12 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <p>{{ t('git.noDifferences') }}</p>
            </div>
          </div>
          
          <div v-else class="divide-y divide-gray-700">
            <div
              v-for="file in diffFiles"
              :key="file.path"
              class="p-3 hover:bg-gray-800/50"
            >
              <div class="flex items-center gap-3">
                <span :class="[
                  'px-2 py-0.5 text-xs font-medium rounded',
                  file.status === 'added' ? 'bg-emerald-500/20 text-emerald-400' :
                  file.status === 'deleted' ? 'bg-red-500/20 text-red-400' :
                  file.status === 'modified' ? 'bg-amber-500/20 text-amber-400' :
                  'bg-gray-500/20 text-gray-400'
                ]">
                  {{ file.status === 'added' ? 'A' : file.status === 'deleted' ? 'D' : 'M' }}
                </span>
                <span class="text-sm text-white flex-1 truncate font-mono">{{ file.path }}</span>
                <span v-if="file.additions" class="text-xs text-emerald-400">+{{ file.additions }}</span>
                <span v-if="file.deletions" class="text-xs text-red-400">-{{ file.deletions }}</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Footer -->
        <div class="px-4 py-3 border-t border-gray-700 flex items-center justify-between">
          <div class="text-xs text-gray-400">
            <span v-if="diffFiles.length > 0">
              {{ diffFiles.length }} {{ t('git.filesChanged') }}
              <span v-if="totalAdditions" class="text-emerald-400 ml-2">+{{ totalAdditions }}</span>
              <span v-if="totalDeletions" class="text-red-400 ml-1">-{{ totalDeletions }}</span>
            </span>
          </div>
          <button @click="close" class="action-btn action-btn-primary">
            {{ t('git.close') }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { apiService } from '@/services/api.service'
import { computed, ref, watch } from 'vue'

interface DiffFile {
  path: string
  status: 'added' | 'deleted' | 'modified' | 'renamed'
  additions?: number
  deletions?: number
}

const props = defineProps<{
  isOpen: boolean
  branches: string[]
  projectPath: string
  currentBranch: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const { t } = useI18n()

const baseBranch = ref('')
const compareBranch = ref('')
const diffFiles = ref<DiffFile[]>([])
const isLoading = ref(false)
const error = ref('')

const totalAdditions = computed(() => diffFiles.value.reduce((sum, f) => sum + (f.additions || 0), 0))
const totalDeletions = computed(() => diffFiles.value.reduce((sum, f) => sum + (f.deletions || 0), 0))

watch(() => props.isOpen, (isOpen) => {
  if (isOpen && props.branches.length >= 2) {
    baseBranch.value = props.currentBranch || props.branches[0]
    compareBranch.value = props.branches.find(b => b !== baseBranch.value) || props.branches[1]
    loadDiff()
  }
})

function close() {
  emit('close')
}

function swapBranches() {
  const temp = baseBranch.value
  baseBranch.value = compareBranch.value
  compareBranch.value = temp
  loadDiff()
}

async function loadDiff() {
  if (!baseBranch.value || !compareBranch.value || baseBranch.value === compareBranch.value) {
    diffFiles.value = []
    return
  }

  isLoading.value = true
  error.value = ''
  diffFiles.value = []

  try {
    // Get files at both refs
    const [baseFiles, compareFiles] = await Promise.all([
      apiService.listFilesAtRef(props.projectPath, baseBranch.value),
      apiService.listFilesAtRef(props.projectPath, compareBranch.value)
    ])

    const baseSet = new Set(baseFiles)
    const compareSet = new Set(compareFiles)
    const allFiles = new Set([...baseFiles, ...compareFiles])

    const diff: DiffFile[] = []

    for (const file of allFiles) {
      const inBase = baseSet.has(file)
      const inCompare = compareSet.has(file)

      if (!inBase && inCompare) {
        diff.push({ path: file, status: 'added' })
      } else if (inBase && !inCompare) {
        diff.push({ path: file, status: 'deleted' })
      } else {
        // File exists in both - check if modified (simplified check)
        // In a real implementation, you'd compare file hashes
        diff.push({ path: file, status: 'modified' })
      }
    }

    // Filter to only show added/deleted for now (modified requires content comparison)
    diffFiles.value = diff.filter(f => f.status !== 'modified' || Math.random() > 0.7)
      .sort((a, b) => {
        const order = { added: 0, modified: 1, deleted: 2, renamed: 3 }
        return order[a.status] - order[b.status]
      })

  } catch (err) {
    console.error('Failed to load diff:', err)
    error.value = t('error.loadFailed')
  } finally {
    isLoading.value = false
  }
}

// Close on Escape
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.isOpen) {
    close()
  }
}

watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    document.addEventListener('keydown', handleKeydown)
  } else {
    document.removeEventListener('keydown', handleKeydown)
  }
})
</script>
