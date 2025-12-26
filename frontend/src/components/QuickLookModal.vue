<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="quicklook-overlay" @click.self="close" @keydown.escape="close">
        <div class="quicklook-modal">
          <!-- Header -->
          <div class="quicklook-header">
            <div class="flex items-center gap-2 min-w-0">
              <span class="text-lg">{{ getFileIcon(fileName) }}</span>
              <span class="quicklook-filename">{{ fileName }}</span>
              <span class="quicklook-path">{{ filePath }}</span>
            </div>
            <div class="flex items-center gap-2">
              <span v-if="fileSize" class="chip-unified chip-unified-accent">{{ formatSize(fileSize) }}</span>
              <button @click="close" class="icon-btn icon-btn-danger" :aria-label="t('common.close')">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Content -->
          <div class="quicklook-content scrollable-y">
            <div v-if="isLoading" class="quicklook-loading">
              <svg class="loading-spinner" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              <p>{{ t('files.loading') }}</p>
            </div>

            <div v-else-if="error" class="quicklook-error">
              <svg class="w-12 h-12 text-red-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <p>{{ error }}</p>
            </div>

            <div v-else-if="isBinary" class="quicklook-binary">
              <svg class="w-16 h-16 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <p class="text-lg font-medium text-gray-300">{{ t('files.binaryFile') }}</p>
              <p class="text-sm text-gray-400">{{ t('files.cannotPreview') }}</p>
            </div>

            <pre v-else class="quicklook-code"><code v-html="highlightedContent"></code></pre>
          </div>

          <!-- Footer -->
          <div class="quicklook-footer">
            <button @click="copyContent" class="btn-unified btn-unified-secondary" :disabled="!content">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              {{ t('context.copy') }}
            </button>
            <button @click="addToContext" class="btn-unified btn-unified-primary" :disabled="!content">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              {{ t('files.addToContext') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { getFileIcon } from '@/utils/fileIcons'
import { highlight } from '@/utils/highlighter'
import { computed, ref, watch } from 'vue'

const props = defineProps<{
  modelValue: boolean
  filePath: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'add-to-context': [path: string]
}>()

const { t } = useI18n()
const projectStore = useProjectStore()

const isOpen = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const content = ref('')
const isLoading = ref(false)
const error = ref<string | null>(null)
const fileSize = ref(0)

const fileName = computed(() => {
  const parts = props.filePath.split('/')
  return parts[parts.length - 1]
})

const fileExtension = computed(() => {
  const parts = fileName.value.split('.')
  return parts.length > 1 ? parts[parts.length - 1].toLowerCase() : ''
})

const isBinary = computed(() => {
  const binaryExtensions = ['png', 'jpg', 'jpeg', 'gif', 'webp', 'ico', 'svg', 'pdf', 'zip', 'tar', 'gz', 'exe', 'dll', 'so', 'dylib', 'woff', 'woff2', 'ttf', 'eot', 'mp3', 'mp4', 'wav', 'avi', 'mov']
  return binaryExtensions.includes(fileExtension.value)
})

const highlightedContent = computed(() => {
  if (!content.value || isBinary.value) return ''
  
  try {
    return highlight(content.value, fileExtension.value)
  } catch {
    return escapeHtml(content.value)
  }
})

function escapeHtml(text: string): string {
  return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

watch(() => [props.modelValue, props.filePath], async ([open, path]) => {
  if (open && path) {
    await loadContent()
  }
}, { immediate: true })

async function loadContent() {
  if (!props.filePath || !projectStore.currentPath) return
  
  isLoading.value = true
  error.value = null
  content.value = ''
  
  try {
    if (isBinary.value) {
      isLoading.value = false
      return
    }
    
    const result = await apiService.readFileContent(projectStore.currentPath, props.filePath)
    content.value = result
    fileSize.value = result.length
  } catch (err) {
    error.value = err instanceof Error ? err.message : t('error.loadFailed')
  } finally {
    isLoading.value = false
  }
}

function close() {
  isOpen.value = false
}

async function copyContent() {
  if (content.value) {
    await navigator.clipboard.writeText(content.value)
  }
}

function addToContext() {
  emit('add-to-context', props.filePath)
  close()
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}
</script>

<style scoped>
/* Modal transition - styles are in design-tokens.css */
.modal-enter-active,
.modal-leave-active {
  transition: all 200ms ease-out;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .quicklook-modal,
.modal-leave-to .quicklook-modal {
  transform: scale(0.95) translateY(-10px);
}
</style>
