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
      <div class="relative w-full max-w-4xl max-h-[85vh] bg-gray-900 rounded-xl border border-gray-700 shadow-2xl flex flex-col overflow-hidden">
        <!-- Header -->
        <div class="flex items-center justify-between p-4 border-b border-gray-700">
          <div class="flex items-center gap-3 min-w-0">
            <span class="text-lg">{{ getFileIcon(fileName) }}</span>
            <div class="min-w-0">
              <h3 class="text-sm font-medium text-white truncate">{{ fileName }}</h3>
              <p class="text-xs text-gray-400 truncate">{{ filePath }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <span v-if="fileSize" class="text-xs text-gray-400">{{ formatSize(fileSize) }}</span>
            <button
              @click="copyContent"
              class="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors"
              :title="t('action.copy')"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
            <button
              @click="close"
              class="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
        
        <!-- Content -->
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
          <pre v-else class="p-4 text-sm text-gray-300 font-mono whitespace-pre-wrap break-words"><code>{{ content }}</code></pre>
        </div>
        
        <!-- Footer with line count -->
        <div class="px-4 py-2 border-t border-gray-700 flex items-center justify-between text-xs text-gray-400">
          <span>{{ lineCount }} {{ t('context.lines') }}</span>
          <span v-if="language">{{ language }}</span>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';
import { useUIStore } from '@/stores/ui.store';
import { getFileIcon } from '@/utils/fileIcons';
import { computed, watch } from 'vue';

const props = defineProps<{
  isOpen: boolean
  filePath: string
  content?: string
  isLoading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const { t } = useI18n()
const uiStore = useUIStore()

const fileName = computed(() => props.filePath.split('/').pop() || '')
const fileSize = computed(() => props.content?.length || 0)
const lineCount = computed(() => props.content?.split('\n').length || 0)

const language = computed(() => {
  const ext = fileName.value.split('.').pop()?.toLowerCase()
  const langMap: Record<string, string> = {
    ts: 'TypeScript', tsx: 'TypeScript', js: 'JavaScript', jsx: 'JavaScript',
    vue: 'Vue', go: 'Go', py: 'Python', java: 'Java', rs: 'Rust',
    cpp: 'C++', c: 'C', h: 'C/C++', cs: 'C#', rb: 'Ruby', php: 'PHP',
    swift: 'Swift', kt: 'Kotlin', css: 'CSS', scss: 'SCSS', html: 'HTML',
    json: 'JSON', yaml: 'YAML', yml: 'YAML', md: 'Markdown', xml: 'XML'
  }
  return ext ? langMap[ext] : undefined
})

function close() {
  emit('close')
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

async function copyContent() {
  if (!props.content) return
  try {
    await navigator.clipboard.writeText(props.content)
    uiStore.addToast(t('toast.contextCopied'), 'success')
  } catch {
    uiStore.addToast(t('toast.copyError'), 'error')
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
