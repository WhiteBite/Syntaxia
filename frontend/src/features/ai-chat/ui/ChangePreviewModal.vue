<template>
  <BaseModal
    :model-value="true"
    size="lg"
    :close-on-backdrop="true"
    @close="$emit('close')"
  >
    <template #header>
      <div class="flex items-center gap-3">
        <div class="section-icon section-icon-yellow">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <h2 class="text-lg font-semibold text-white">
          {{ t('sandbox.reviewChanges') }}
        </h2>
        <span class="badge badge-warning">
          {{ sandboxStore.changeCount }} {{ t('sandbox.files') }}
        </span>
      </div>
    </template>

    <!-- Content -->
    <div class="flex min-h-0" style="height: 60vh;">
      <!-- File list -->
      <div class="border-r border-gray-700 overflow-y-auto flex-shrink-0" style="width: min(16rem, 30%);">
        <div
          v-for="change in sandboxStore.changes"
          :key="change.path"
          @click="selectedFile = change.path"
          :class="[
            'p-3 cursor-pointer border-b border-gray-800 hover:bg-gray-800',
            selectedFile === change.path ? 'bg-gray-800' : ''
          ]"
        >
          <div class="flex items-center gap-2">
            <span :class="getOperationClass(change.operation)">
              {{ getOperationIcon(change.operation) }}
            </span>
            <span class="text-sm text-gray-300 truncate">{{ getFileName(change.path) }}</span>
          </div>
          <div class="text-xs text-gray-500 truncate mt-1">{{ change.path }}</div>
        </div>
      </div>

      <!-- Diff view -->
      <div class="flex-1 min-w-0 overflow-auto p-4">
        <div v-if="selectedFile && currentDiff" class="font-mono text-sm">
          <pre class="whitespace-pre-wrap"><code v-html="highlightDiff(currentDiff)"></code></pre>
        </div>
        <div v-else-if="isLoadingDiff" class="flex items-center justify-center h-full">
          <BaseSpinner size="lg" class="text-purple-500" />
        </div>
        <div v-else class="empty-state h-full">
          <p class="text-gray-400">{{ t('sandbox.selectFile') }}</p>
        </div>
      </div>
    </div>

    <template #footer>
      <button
        @click="handleDiscardAll"
        :disabled="sandboxStore.isLoading"
        class="btn btn-danger"
      >
        {{ t('sandbox.discardAll') }}
      </button>
      <div class="flex gap-2 ml-auto">
        <button @click="$emit('close')" class="btn btn-secondary">
          {{ t('common.cancel') }}
        </button>
        <button
          @click="handleApplyAll"
          :disabled="sandboxStore.isLoading"
          class="btn btn-primary"
        >
          {{ t('sandbox.applyAll') }}
        </button>
      </div>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { BaseModal, BaseSpinner } from '@/components/ui'
import { useI18n } from '@/composables/useI18n'
import { useSandboxStore } from '@/stores/sandbox.store'
import { ref, watch } from 'vue'

const emit = defineEmits<{
  close: []
  applied: []
}>()

const { t } = useI18n()
const sandboxStore = useSandboxStore()

const selectedFile = ref<string | null>(null)
const currentDiff = ref<string>('')
const isLoadingDiff = ref(false)

// Select first file by default
if (sandboxStore.changes.length > 0) {
  selectedFile.value = sandboxStore.changes[0].path
}

watch(selectedFile, async (path) => {
  if (!path) {
    currentDiff.value = ''
    return
  }
  
  isLoadingDiff.value = true
  try {
    currentDiff.value = await sandboxStore.getDiff(path)
  } finally {
    isLoadingDiff.value = false
  }
})

function getFileName(path: string): string {
  return path.split('/').pop() || path
}

function getOperationClass(op: string): string {
  switch (op) {
    case 'create': return 'text-green-400'
    case 'modify': return 'text-yellow-400'
    case 'delete': return 'text-red-400'
    default: return 'text-gray-400'
  }
}

function getOperationIcon(op: string): string {
  switch (op) {
    case 'create': return '+'
    case 'modify': return '~'
    case 'delete': return '-'
    default: return '?'
  }
}

function highlightDiff(diff: string): string {
  return diff
    .split('\n')
    .map(line => {
      if (line.startsWith('+') && !line.startsWith('+++')) {
        return `<span class="text-green-400">${escapeHtml(line)}</span>`
      }
      if (line.startsWith('-') && !line.startsWith('---')) {
        return `<span class="text-red-400">${escapeHtml(line)}</span>`
      }
      if (line.startsWith('@@')) {
        return `<span class="text-cyan-400">${escapeHtml(line)}</span>`
      }
      if (line.startsWith('diff') || line.startsWith('---') || line.startsWith('+++')) {
        return `<span class="text-gray-500">${escapeHtml(line)}</span>`
      }
      return escapeHtml(line)
    })
    .join('\n')
}

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
}

async function handleApplyAll() {
  const success = await sandboxStore.applyAll()
  if (success) {
    emit('applied')
  }
}

async function handleDiscardAll() {
  await sandboxStore.discardAll()
  emit('close')
}
</script>
