<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-75"
        @click.self="$emit('close')" @keydown.esc="$emit('close')">
        <div
          class="bg-gray-900 rounded-lg shadow-2xl w-full max-w-5xl max-h-[90vh] flex flex-col border border-gray-700"
          @click.stop>
          <!-- Header -->
          <div class="flex items-center justify-between px-4 py-3 border-b border-gray-700">
            <div class="flex items-center gap-3 flex-1 min-w-0">
              <!-- File Icon -->
              <svg class="w-5 h-5 text-blue-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>

              <!-- File Info -->
              <div class="flex-1 min-w-0">
                <h3 class="text-sm font-semibold text-white truncate" :title="filePath">
                  {{ fileName }}
                </h3>
                <p class="text-xs text-gray-400 truncate" :title="filePath">
                  {{ filePath }}
                </p>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex items-center gap-2">
              <button v-if="!isLoading && !error" @click="copyToClipboard"
                class="p-2 hover:bg-gray-800 rounded transition-colors" title="Копировать">
                <svg class="w-4 h-4 text-gray-400 hover:text-white" fill="none" stroke="currentColor"
                  viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
              </button>

              <button @click="$emit('close')" class="p-2 hover:bg-gray-800 rounded transition-colors" title="Закрыть">
                <svg class="w-4 h-4 text-gray-400 hover:text-white" fill="none" stroke="currentColor"
                  viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Stats Bar -->
          <div v-if="stats && !isLoading" class="px-4 py-2 bg-gray-850 border-b border-gray-700">
            <div class="flex items-center gap-6 text-xs text-gray-400">
              <div class="flex items-center gap-2">
                <span class="text-gray-400">Размер:</span>
                <span class="text-white">{{ formatBytes(stats.size) }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="text-gray-400">Строк:</span>
                <span class="text-white">{{ formatNumber(stats.lines) }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="text-gray-400">Токены:</span>
                <span class="text-white">{{ formatNumber(stats.tokens) }}</span>
              </div>
              <div v-if="stats.language" class="flex items-center gap-2">
                <span class="text-gray-400">Язык:</span>
                <span class="text-white">{{ stats.language }}</span>
              </div>
            </div>
          </div>

          <!-- Content -->
          <div class="flex-1 overflow-hidden">
            <!-- Loading State -->
            <div v-if="isLoading" class="flex items-center justify-center h-full">
              <div class="flex flex-col items-center gap-3">
                <div class="animate-spin h-10 w-10 border-4 border-blue-500 border-t-transparent rounded-full"></div>
                <div class="text-gray-400 text-sm">Загрузка файла...</div>
              </div>
            </div>

            <!-- Error State -->
            <div v-else-if="error" class="flex items-center justify-center h-full">
              <div class="text-center text-red-400">
                <svg class="w-12 h-12 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <p class="text-sm font-semibold mb-1">Ошибка загрузки файла</p>
                <p class="text-xs text-gray-400">{{ error }}</p>
              </div>
            </div>

            <!-- File Content -->
            <div v-else class="h-full overflow-auto bg-gray-900">
              <pre
                class="p-4 text-sm font-mono text-gray-300 whitespace-pre-wrap break-words"><code>{{ content }}</code></pre>
            </div>
          </div>

          <!-- Footer -->
          <div class="flex items-center justify-between px-4 py-3 border-t border-gray-700">
            <div class="text-xs text-gray-400">
              ESC для закрытия
            </div>
            <div class="flex items-center gap-2">
              <button @click="close"
                class="px-4 py-2 text-sm bg-gray-800 hover:bg-gray-700 text-white rounded transition-colors">
                Закрыть
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useLogger } from '@/composables/useLogger'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { ReadFileContent } from '../../wailsjs/go/main/App'

const logger = useLogger('QuickLook')

interface Props {
  filePath?: string
}

interface FileStats {
  size: number
  lines: number
  tokens: number
  language?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const projectStore = useProjectStore()
const uiStore = useUIStore()

const isOpen = ref(false)
const isLoading = ref(false)
const error = ref<string | null>(null)
const content = ref<string>('')
const stats = ref<FileStats | null>(null)

const fileName = computed(() => {
  if (!props.filePath) return ''
  const parts = props.filePath.split(/[\\/]/)
  return parts[parts.length - 1]
})

// Watch for file path changes
watch(() => props.filePath, async (newPath) => {
  if (newPath && isOpen.value) {
    await loadFile()
  }
})

// Keyboard shortcuts
function handleKeyDown(event: KeyboardEvent) {
  if (event.key === 'Escape' && isOpen.value) {
    close()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
})

/**
 * Open QuickLook and load file
 */
async function open(filePath?: string) {
  isOpen.value = true

  // Use provided paths or props
  const targetPath = filePath || props.filePath
  const targetProject = projectStore.projectPath

  if (targetPath && targetProject) {
    await loadFile(targetPath)
  }
}

/**
 * Close QuickLook
 */
function close() {
  isOpen.value = false
  emit('close')
}

/**
 * Load file content from backend
 */
async function loadFile(filePath?: string) {
  const targetPath = filePath || props.filePath
  const targetProject = projectStore.projectPath

  if (!targetPath || !targetProject) {
    error.value = 'Путь к файлу не указан'
    return
  }

  isLoading.value = true
  error.value = null
  content.value = ''
  stats.value = null

  try {
    // Get relative path
    const relPath = targetPath.replace(targetProject, '').replace(/^[\\/]+/, '')

    // Load file content
    const fileContent = await ReadFileContent(targetProject, relPath)
    content.value = fileContent

    // Calculate stats
    const lines = fileContent.split('\n').length
    const size = new Blob([fileContent]).size
    const tokens = Math.floor(fileContent.length / 4) // Rough estimate

    // Detect language from extension
    const ext = targetPath.split('.').pop()?.toLowerCase()
    const language = detectLanguage(ext || '')

    stats.value = {
      size,
      lines,
      tokens,
      language
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить файл'
    logger.error('Failed to load file:', err)
  } finally {
    isLoading.value = false
  }
}

/**
 * Copy content to clipboard
 */
async function copyToClipboard() {
  if (!content.value) return

  try {
    await navigator.clipboard.writeText(content.value)
    uiStore.addToast('Content copied to clipboard', 'success')
  } catch (err) {
    logger.error('Failed to copy to clipboard:', err)
    uiStore.addToast('Failed to copy to clipboard', 'error')
  }
}

/**
 * Detect language from file extension
 */
function detectLanguage(ext: string): string {
  const langMap: Record<string, string> = {
    'ts': 'TypeScript',
    'tsx': 'TypeScript React',
    'js': 'JavaScript',
    'jsx': 'JavaScript React',
    'vue': 'Vue',
    'go': 'Go',
    'py': 'Python',
    'java': 'Java',
    'cpp': 'C++',
    'c': 'C',
    'cs': 'C#',
    'rs': 'Rust',
    'rb': 'Ruby',
    'php': 'PHP',
    'html': 'HTML',
    'css': 'CSS',
    'scss': 'SCSS',
    'json': 'JSON',
    'yaml': 'YAML',
    'yml': 'YAML',
    'xml': 'XML',
    'md': 'Markdown',
    'sql': 'SQL',
    'sh': 'Shell',
    'bash': 'Bash'
  }

  return langMap[ext] || ext.toUpperCase()
}

// Formatting utilities
function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

function formatNumber(num: number): string {
  if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`
  if (num >= 1000) return `${(num / 1000).toFixed(1)}K`
  return num.toString()
}

// Expose methods to parent
defineExpose({
  open,
  close
})
</script>

<style scoped>
.modal-enter-active>div,
.modal-leave-active>div {
  transition: transform 0.2s ease;
}

.modal-enter-from>div,
.modal-leave-to>div {
  transform: scale(0.95);
}

/* Scrollbar styling */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #1f2937;
}

::-webkit-scrollbar-thumb {
  background: #4b5563;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #6b7280;
}
</style>
/style>
